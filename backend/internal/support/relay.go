package support

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"support-chat/internal/db"
)

func (s *Service) Kick() {
	select {
	case s.kick <- struct{}{}:
	default:
	}
}

func (s *Service) RunRelay(ctx context.Context) {
	if s.tg == nil {
		log.Println("support: telegram bot not configured; relay idle")
		return
	}
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	log.Println("support: relay worker started")
	for {
		if n, err := s.RelayPending(ctx); err != nil {
			log.Printf("support: relay pending failed: %v", err)
		} else if n > 0 {
			log.Printf("support: relayed %d messages to Telegram", n)
		}

		select {
		case <-ctx.Done():
			return
		case <-s.kick:
		case <-ticker.C:
		}
	}
}

func (s *Service) RelayPending(ctx context.Context) (int, error) {
	if s.tg == nil {
		return 0, nil
	}
	var pending []db.SupportMessage
	if err := s.db.WithContext(ctx).
		Where("author = ? AND relayed_at IS NULL", db.AuthorUser).
		Order("created_at, id").Limit(50).Find(&pending).Error; err != nil {
		return 0, err
	}

	blocked := map[string]bool{}
	sent := 0
	for _, m := range pending {
		if blocked[m.ThreadID] {
			continue
		}
		if err := s.relay(ctx, m); err != nil {
			log.Printf("support: relay message %s failed: %v", m.ID, err)
			blocked[m.ThreadID] = true
			if strings.Contains(strings.ToLower(err.Error()), "thread not found") || strings.Contains(err.Error(), "TOPIC_DELETED") {
				s.db.WithContext(ctx).Model(&db.SupportThread{}).Where("id = ?", m.ThreadID).
					Updates(map[string]any{"telegram_chat_id": nil, "telegram_topic_id": nil})
			}
			continue
		}
		sent++
	}
	return sent, nil
}

func (s *Service) relay(ctx context.Context, m db.SupportMessage) error {
	var t db.SupportThread
	if err := s.db.WithContext(ctx).First(&t, "id = ?", m.ThreadID).Error; err != nil {
		return err
	}
	topic, err := s.topicFor(ctx, &t)
	if err != nil {
		return err
	}

	who := t.Name
	if who == "" {
		who = t.Email
	}
	text := fmt.Sprintf("👤 %s: %s", who, m.Text)
	if topic == 0 {
		text = fmt.Sprintf("#%d %s", t.Number, text)
	}

	var mid int64
	if m.HasAttachment {
		var a db.SupportAttachment
		if err := s.db.WithContext(ctx).First(&a, "message_id = ?", m.ID).Error; err != nil {
			return err
		}
		mid, err = s.tg.SendPhoto(ctx, topic, text, a.Data, a.ContentType)
	} else {
		mid, err = s.tg.Send(ctx, topic, text)
	}
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	return s.db.WithContext(ctx).Model(&db.SupportMessage{}).Where("id = ?", m.ID).
		Updates(map[string]any{"telegram_message_id": mid, "relayed_at": &now}).Error
}

func (s *Service) topicFor(ctx context.Context, t *db.SupportThread) (int64, error) {
	if t.TelegramTopicID != nil && t.TelegramChatID != nil && *t.TelegramChatID == s.tg.ChatID {
		return *t.TelegramTopicID, nil
	}

	who := t.Name
	if who == "" {
		who = t.Email
	}
	topicName := fmt.Sprintf("#%d %s", t.Number, who)
	topic, err := s.tg.CreateTopic(ctx, topicName)
	if err != nil {
		log.Printf("support: could not create forum topic for %s (%v); sending to chat", topicName, err)
		topic = 0
	}

	cardID, err := s.tg.Send(ctx, topic, s.card(t))
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	sysMsg := db.SupportMessage{
		ID:                fmt.Sprintf("card-%s", t.ID),
		ThreadID:          t.ID,
		Author:            db.AuthorSystem,
		Text:              "card",
		TelegramMessageID: &cardID,
		RelayedAt:         &now,
		CreatedAt:         now,
	}
	_ = s.db.WithContext(ctx).Create(&sysMsg).Error

	if topic != 0 {
		chatID := s.tg.ChatID
		if err := s.db.WithContext(ctx).Model(&db.SupportThread{}).Where("id = ?", t.ID).
			Updates(map[string]any{"telegram_chat_id": chatID, "telegram_topic_id": topic}).Error; err != nil {
			return 0, err
		}
		t.TelegramChatID = &chatID
		t.TelegramTopicID = &topic
	}
	return topic, nil
}

func (s *Service) card(t *db.SupportThread) string {
	var c struct {
		Page      string `json:"page"`
		Lang      string `json:"lang"`
		UserAgent string `json:"user_agent"`
	}
	_ = json.Unmarshal([]byte(t.Context), &c)

	var b strings.Builder
	fmt.Fprintf(&b, "🆕 Обращение #%d\n", t.Number)
	if t.Name != "" {
		fmt.Fprintf(&b, "👤 Имя: %s\n", t.Name)
	}
	fmt.Fprintf(&b, "✉️ Email: %s\n", t.Email)
	if c.Page != "" {
		fmt.Fprintf(&b, "🌐 Страница: %s\n", c.Page)
	}
	if c.Lang != "" {
		fmt.Fprintf(&b, "🌍 Язык: %s\n", c.Lang)
	}
	if c.UserAgent != "" {
		fmt.Fprintf(&b, "📱 Устройство: %s\n", c.UserAgent)
	}
	b.WriteString("\n💬 Отвечайте в этом топике прямо через Telegram.\n/close — закрыть тикет\n/reopen — открыть заново")
	return b.String()
}
