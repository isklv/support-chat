package support

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"support-chat/internal/db"
)

// RunPoller runs the Telegram long-polling loop.
func (s *Service) RunPoller(ctx context.Context) {
	if s.tg == nil {
		return
	}

	log.Println("support: deleting webhook before starting long polling...")
	if err := s.tg.DeleteWebhook(ctx); err != nil {
		log.Printf("support: failed to delete webhook: %v", err)
	}

	log.Println("support: long-polling worker started")
	var offset int64

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		pollCtx, cancel := context.WithTimeout(ctx, 35*time.Second)
		updates, err := s.tg.GetUpdates(pollCtx, offset, 25)
		cancel()

		if err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("support: telegram getUpdates error: %v", err)
				time.Sleep(3 * time.Second)
			}
			continue
		}

		for _, u := range updates {
			if u.UpdateID >= offset {
				offset = u.UpdateID + 1
			}
			if err := s.HandleUpdate(ctx, u); err != nil {
				log.Printf("support: handle telegram update %d failed: %v", u.UpdateID, err)
			}
		}
	}
}

// HandleUpdate processes an update received either via Webhook or Polling.
func (s *Service) HandleUpdate(ctx context.Context, u Update) error {
	m := u.Message
	if s.tg == nil || m == nil || m.Chat.ID != s.tg.ChatID || m.From == nil || m.From.IsBot || m.ForumTopicCreated != nil {
		return nil
	}

	t, err := s.threadForMessage(ctx, m)
	if err != nil || t == nil {
		return err
	}

	operator := displayName(m.From)
	text := strings.TrimSpace(m.Text)
	if text == "" {
		text = strings.TrimSpace(m.Caption)
	}

	switch command(text) {
	case "/close":
		_, err := s.SetStatus(ctx, t.ID, false, operator)
		return err
	case "/reopen":
		_, err := s.SetStatus(ctx, t.ID, true, operator)
		return err
	}

	var img []byte
	contentType := ""
	fileID := ""
	if len(m.Photo) > 0 {
		fileID = m.Photo[len(m.Photo)-1].FileID // largest photo
	} else if m.Document != nil && strings.HasPrefix(m.Document.MimeType, "image/") {
		fileID = m.Document.FileID
	}

	if fileID != "" {
		data, err := s.tg.Download(ctx, fileID, MaxImageBytes)
		if err == nil {
			if ct, err := detectImage(data); err == nil {
				img, contentType = data, ct
			}
		}
		if contentType == "" {
			log.Printf("support: operator image in thread #%d was rejected or unreadable", t.Number)
		}
	}

	if text == "" && contentType == "" {
		return nil
	}

	mid := m.MessageID
	_, err = s.AddOperatorMessage(ctx, t, operator, clip(text, MaxTextLen), img, contentType, &mid)
	return err
}

func command(text string) string {
	if !strings.HasPrefix(text, "/") {
		return ""
	}
	word := strings.Fields(text)[0]
	return strings.ToLower(strings.SplitN(word, "@", 2)[0])
}

func (s *Service) threadForMessage(ctx context.Context, m *tgMessage) (*db.SupportThread, error) {
	var t db.SupportThread
	if m.IsTopicMessage && m.MessageThreadID != 0 {
		err := s.db.WithContext(ctx).
			First(&t, "telegram_chat_id = ? AND telegram_topic_id = ?", m.Chat.ID, m.MessageThreadID).Error
		if err == nil {
			return &t, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	if m.ReplyTo != nil {
		var msg db.SupportMessage
		err := s.db.WithContext(ctx).Where("telegram_message_id = ?", m.ReplyTo.MessageID).First(&msg).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		if err := s.db.WithContext(ctx).First(&t, "id = ?", msg.ThreadID).Error; err != nil {
			return nil, err
		}
		return &t, nil
	}

	return nil, nil
}
