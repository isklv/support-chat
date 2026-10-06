package support

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"support-chat/internal/config"
	"support-chat/internal/db"
	"support-chat/internal/mailer"
)

const (
	MaxTextLen    = 4000
	MaxImageBytes = 4 << 20 // 4MB
	maxImageSide  = 8000
	emailEvery    = 15 * time.Minute
)

var (
	ErrEmpty        = errors.New("support: empty message")
	ErrTooLong      = errors.New("support: message too long")
	ErrNeedEmail    = errors.New("support: an email is required for the first message")
	ErrNotFound     = errors.New("support: conversation not found")
	ErrInvalidImage = errors.New("support: invalid image format or size")
)

type Service struct {
	db   *gorm.DB
	cfg  *config.Config
	mail mailer.Mailer
	tg   *Telegram
	kick chan struct{}
}

func NewService(database *gorm.DB, cfg *config.Config, m mailer.Mailer, tg *Telegram) *Service {
	return &Service{
		db:   database,
		cfg:  cfg,
		mail: m,
		tg:   tg,
		kick: make(chan struct{}, 1),
	}
}

func (s *Service) Telegram() *Telegram {
	return s.tg
}

func (s *Service) DB() *gorm.DB {
	return s.db
}

type Visitor struct {
	UserID string
	Token  string
}

type View struct {
	Thread   *db.SupportThread   `json:"thread"`
	Messages []db.SupportMessage `json:"messages"`
	Token    string              `json:"token,omitempty"`
}

type PostInput struct {
	Text      string
	Email     string
	Name      string
	Page      string
	Lang      string
	UserAgent string
	Image     []byte
	Metadata  map[string]any
}

func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func GenerateToken() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Service) current(ctx context.Context, v Visitor) (*db.SupportThread, error) {
	var t db.SupportThread
	q := s.db.WithContext(ctx)
	var err error

	switch {
	case v.UserID != "" && v.UserID != "guest":
		err = q.Where("email = ? OR id IN (SELECT thread_id FROM support_tokens WHERE token_hash = ?)", v.UserID, HashToken(v.Token)).
			Order("created_at DESC").First(&t).Error
	case v.Token != "":
		err = q.Joins("JOIN support_tokens st ON st.thread_id = support_threads.id").
			Where("st.token_hash = ?", HashToken(v.Token)).Order("support_threads.created_at DESC").First(&t).Error
	default:
		return nil, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func (s *Service) messages(ctx context.Context, threadID string) ([]db.SupportMessage, error) {
	out := []db.SupportMessage{}
	err := s.db.WithContext(ctx).
		Where("thread_id = ? AND NOT (author = ? AND text = ?)", threadID, db.AuthorSystem, "card").
		Order("created_at, id").Find(&out).Error
	return out, err
}

func (s *Service) Current(ctx context.Context, v Visitor) (*View, error) {
	t, err := s.current(ctx, v)
	if err != nil || t == nil {
		return &View{Messages: []db.SupportMessage{}}, err
	}
	msgs, err := s.messages(ctx, t.ID)
	return &View{Thread: t, Messages: msgs}, err
}

func (s *Service) Unread(ctx context.Context, v Visitor) (int, error) {
	t, err := s.current(ctx, v)
	if err != nil || t == nil {
		return 0, err
	}
	var n int64
	q := s.db.WithContext(ctx).Model(&db.SupportMessage{}).
		Where("thread_id = ? AND author = ?", t.ID, db.AuthorOperator)
	if t.UserSeenAt != nil {
		q = q.Where("created_at > ?", *t.UserSeenAt)
	}
	err = q.Count(&n).Error
	return int(n), err
}

func (s *Service) MarkSeen(ctx context.Context, v Visitor) error {
	t, err := s.current(ctx, v)
	if err != nil || t == nil {
		return err
	}
	now := time.Now().UTC()
	return s.db.WithContext(ctx).Model(&db.SupportThread{}).Where("id = ?", t.ID).Update("user_seen_at", &now).Error
}

func detectImage(data []byte) (string, error) {
	if len(data) == 0 || len(data) > MaxImageBytes {
		return "", ErrInvalidImage
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width == 0 || cfg.Height == 0 || cfg.Width > maxImageSide || cfg.Height > maxImageSide {
		return "", ErrInvalidImage
	}
	switch format {
	case "jpeg":
		return "image/jpeg", nil
	case "png":
		return "image/png", nil
	case "webp":
		return "image/webp", nil
	}
	return "", ErrInvalidImage
}

func clip(s string, n int) string {
	s = strings.TrimSpace(s)
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

func normalizeLocale(l string) string {
	parts := strings.SplitN(strings.ToLower(strings.TrimSpace(l)), "-", 2)
	switch parts[0] {
	case "ru":
		return "ru"
	case "es":
		return "es"
	}
	return "en"
}

func normalizeEmail(s string) (string, error) {
	addr, err := mail.ParseAddress(strings.TrimSpace(s))
	if err != nil {
		return "", err
	}
	return strings.ToLower(addr.Address), nil
}

func (s *Service) Post(ctx context.Context, v Visitor, in PostInput) (*View, error) {
	text := strings.TrimSpace(in.Text)
	if utf8.RuneCountInString(text) > MaxTextLen {
		return nil, ErrTooLong
	}
	contentType := ""
	if len(in.Image) > 0 {
		ct, err := detectImage(in.Image)
		if err != nil {
			return nil, err
		}
		contentType = ct
	}
	if text == "" && contentType == "" {
		return nil, ErrEmpty
	}

	thread, err := s.current(ctx, v)
	if err != nil {
		return nil, err
	}
	newToken := ""
	now := time.Now().UTC()

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if thread == nil || thread.Status == db.SupportClosed {
			if thread != nil {
				if strings.TrimSpace(in.Email) == "" {
					in.Email = thread.Email
				}
				if strings.TrimSpace(in.Name) == "" {
					in.Name = thread.Name
				}
			}
			t, token, err := s.newThread(ctx, tx, v, in)
			if err != nil {
				return err
			}
			thread, newToken = t, token
		}

		msgID := uuid.New().String()
		msg := db.SupportMessage{
			ID:            msgID,
			ThreadID:      thread.ID,
			Author:        db.AuthorUser,
			Text:          text,
			HasAttachment: contentType != "",
			CreatedAt:     now,
		}
		if err := tx.Create(&msg).Error; err != nil {
			return fmt.Errorf("store message: %w", err)
		}

		if contentType != "" {
			att := db.SupportAttachment{
				MessageID:   msgID,
				ContentType: contentType,
				Data:        in.Image,
			}
			if err := tx.Create(&att).Error; err != nil {
				return fmt.Errorf("store attachment: %w", err)
			}
		}

		return tx.Model(&db.SupportThread{}).Where("id = ?", thread.ID).
			Updates(map[string]any{"last_message_at": now, "user_seen_at": &now}).Error
	})
	if err != nil {
		return nil, err
	}

	s.Kick()
	msgs, err := s.messages(ctx, thread.ID)
	if err != nil {
		return nil, err
	}
	thread.LastMessageAt = now
	return &View{Thread: thread, Messages: msgs, Token: newToken}, nil
}

func (s *Service) newThread(ctx context.Context, tx *gorm.DB, v Visitor, in PostInput) (*db.SupportThread, string, error) {
	ctxInfo := map[string]any{
		"page":       clip(in.Page, 250),
		"lang":       clip(in.Lang, 10),
		"user_agent": clip(in.UserAgent, 300),
		"created_at": time.Now().UTC().Format(time.RFC3339),
	}
	if in.Metadata != nil {
		for k, val := range in.Metadata {
			ctxInfo[k] = val
		}
	}

	email, err := normalizeEmail(in.Email)
	if err != nil {
		return nil, "", ErrNeedEmail
	}

	rawJSON, _ := json.Marshal(ctxInfo)
	threadID := uuid.New().String()
	now := time.Now().UTC()

	t := db.SupportThread{
		ID:            threadID,
		Email:         email,
		Name:          clip(in.Name, 100),
		Locale:        normalizeLocale(in.Lang),
		Status:        db.SupportOpen,
		Context:       string(rawJSON),
		LastMessageAt: now,
		CreatedAt:     now,
	}

	if err := tx.Create(&t).Error; err != nil {
		return nil, "", fmt.Errorf("create thread: %w", err)
	}

	// Fetch back auto-incremented number if available
	_ = tx.Select("number").First(&t, "id = ?", t.ID).Error

	token := ""
	if v.Token != "" {
		res := tx.Model(&db.SupportToken{}).Where("token_hash = ?", HashToken(v.Token)).Update("thread_id", t.ID)
		if res.Error != nil {
			return nil, "", res.Error
		}
		if res.RowsAffected == 1 {
			return &t, "", nil
		}
	}

	rawToken, err := GenerateToken()
	if err != nil {
		return nil, "", err
	}
	if err := tx.Create(&db.SupportToken{TokenHash: HashToken(rawToken), ThreadID: t.ID, CreatedAt: now}).Error; err != nil {
		return nil, "", fmt.Errorf("store token: %w", err)
	}
	token = rawToken

	return &t, token, nil
}

func (s *Service) Attachment(ctx context.Context, v Visitor, operator bool, messageID string) (*db.SupportAttachment, error) {
	var msg db.SupportMessage
	if err := s.db.WithContext(ctx).First(&msg, "id = ?", messageID).Error; err != nil {
		return nil, ErrNotFound
	}
	if !operator {
		t, err := s.current(ctx, v)
		if err != nil || t == nil || t.ID != msg.ThreadID {
			return nil, ErrNotFound
		}
	}
	var a db.SupportAttachment
	if err := s.db.WithContext(ctx).First(&a, "message_id = ?", messageID).Error; err != nil {
		return nil, ErrNotFound
	}
	return &a, nil
}

type Summary struct {
	db.SupportThread
	LastText   string `json:"last_text"`
	LastAuthor string `json:"last_author"`
	Messages   int    `json:"messages"`
}

func (s *Service) List(ctx context.Context, status string, limit int) ([]Summary, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	var threads []db.SupportThread
	q := s.db.WithContext(ctx).Order("last_message_at DESC").Limit(limit)
	if status == db.SupportOpen || status == db.SupportClosed {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&threads).Error; err != nil {
		return nil, err
	}
	out := make([]Summary, 0, len(threads))
	for _, t := range threads {
		var last db.SupportMessage
		s.db.WithContext(ctx).Where("thread_id = ? AND author <> ?", t.ID, db.AuthorSystem).
			Order("created_at DESC").First(&last)
		var n int64
		s.db.WithContext(ctx).Model(&db.SupportMessage{}).Where("thread_id = ?", t.ID).Count(&n)
		out = append(out, Summary{
			SupportThread: t,
			LastText:      clip(last.Text, 140),
			LastAuthor:    last.Author,
			Messages:      int(n),
		})
	}
	return out, nil
}

func (s *Service) Get(ctx context.Context, id string) (*View, error) {
	var t db.SupportThread
	if err := s.db.WithContext(ctx).First(&t, "id = ?", id).Error; err != nil {
		return nil, ErrNotFound
	}
	msgs, err := s.messages(ctx, t.ID)
	return &View{Thread: &t, Messages: msgs}, err
}

func (s *Service) Reply(ctx context.Context, id, operator, text string, img []byte) (*View, error) {
	text = strings.TrimSpace(text)
	if utf8.RuneCountInString(text) > MaxTextLen {
		return nil, ErrTooLong
	}
	contentType := ""
	if len(img) > 0 {
		ct, err := detectImage(img)
		if err != nil {
			return nil, err
		}
		contentType = ct
	}
	if text == "" && contentType == "" {
		return nil, ErrEmpty
	}
	var t db.SupportThread
	if err := s.db.WithContext(ctx).First(&t, "id = ?", id).Error; err != nil {
		return nil, ErrNotFound
	}
	msg, err := s.AddOperatorMessage(ctx, &t, operator, text, img, contentType, nil)
	if err != nil {
		return nil, err
	}
	if s.InThisChat(&t) {
		go func() {
			c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			caption := fmt.Sprintf("💬 %s: %s", operator, text)
			var mid int64
			var err error
			if contentType != "" {
				mid, err = s.tg.SendPhoto(c, *t.TelegramTopicID, caption, img, contentType)
			} else {
				mid, err = s.tg.Send(c, *t.TelegramTopicID, caption)
			}
			if err != nil {
				log.Printf("support: console reply to telegram failed: %v", err)
				return
			}
			s.db.Model(&db.SupportMessage{}).Where("id = ?", msg.ID).Update("telegram_message_id", mid)
		}()
	}
	return s.Get(ctx, id)
}

func (s *Service) AddOperatorMessage(ctx context.Context, t *db.SupportThread, operator, text string, img []byte, contentType string, tgID *int64) (*db.SupportMessage, error) {
	now := time.Now().UTC()
	msgID := uuid.New().String()
	msg := db.SupportMessage{
		ID:                msgID,
		ThreadID:          t.ID,
		Author:            db.AuthorOperator,
		OperatorName:      clip(operator, 80),
		Text:              text,
		HasAttachment:     contentType != "",
		TelegramMessageID: tgID,
		RelayedAt:         &now,
		CreatedAt:         now,
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&msg).Error; err != nil {
			return fmt.Errorf("store answer: %w", err)
		}
		if contentType != "" {
			if err := tx.Create(&db.SupportAttachment{MessageID: msgID, ContentType: contentType, Data: img}).Error; err != nil {
				return err
			}
		}
		updates := map[string]any{"last_message_at": now}
		if t.Status == db.SupportClosed {
			updates["status"] = db.SupportOpen
			updates["closed_at"] = nil
		}
		return tx.Model(&db.SupportThread{}).Where("id = ?", t.ID).Updates(updates).Error
	})
	if err != nil {
		return nil, err
	}
	t.Status = db.SupportOpen
	go s.notifyEmail(*t, msg)
	return &msg, nil
}

func (s *Service) SetStatus(ctx context.Context, id string, open bool, by string) (*View, error) {
	var t db.SupportThread
	if err := s.db.WithContext(ctx).First(&t, "id = ?", id).Error; err != nil {
		return nil, ErrNotFound
	}
	status, note := db.SupportClosed, "closed"
	updates := map[string]any{"status": db.SupportClosed, "closed_at": time.Now().UTC()}
	if open {
		status, note = db.SupportOpen, "reopened"
		updates = map[string]any{"status": db.SupportOpen, "closed_at": nil}
	}
	if t.Status != status {
		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&db.SupportThread{}).Where("id = ?", t.ID).Updates(updates).Error; err != nil {
				return err
			}
			now := time.Now().UTC()
			sysMsg := db.SupportMessage{
				ID:           uuid.New().String(),
				ThreadID:     t.ID,
				Author:       db.AuthorSystem,
				OperatorName: clip(by, 80),
				Text:         note,
				RelayedAt:    &now,
				CreatedAt:    now,
			}
			return tx.Create(&sysMsg).Error
		})
		if err != nil {
			return nil, err
		}
		if s.InThisChat(&t) {
			topic := *t.TelegramTopicID
			go func() {
				c, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
				label := "✅ Closed by "
				if open {
					label = "↩️ Reopened by "
				}
				_, _ = s.tg.Send(c, topic, label+by)
				_ = s.tg.SetTopicOpen(c, topic, open)
			}()
		}
	}
	return s.Get(ctx, id)
}

func (s *Service) notifyEmail(t db.SupportThread, msg db.SupportMessage) {
	if s.mail == nil || t.Email == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	res := s.db.WithContext(ctx).Model(&db.SupportThread{}).
		Where("id = ? AND (emailed_at IS NULL OR emailed_at < ?)", t.ID, time.Now().UTC().Add(-emailEvery)).
		Update("emailed_at", time.Now().UTC())
	if res.Error != nil || res.RowsAffected == 0 {
		return
	}

	rawToken, err := GenerateToken()
	if err != nil {
		return
	}
	if err := s.db.WithContext(ctx).Create(&db.SupportToken{TokenHash: HashToken(rawToken), ThreadID: t.ID, CreatedAt: time.Now().UTC()}).Error; err != nil {
		return
	}

	link := fmt.Sprintf("%s/?support=%s", s.cfg.AppURL, rawToken)
	subject := "Ответ от службы поддержки"
	if t.Locale == "en" {
		subject = "Reply from Support"
	} else if t.Locale == "es" {
		subject = "Respuesta de Soporte"
	}

	body := fmt.Sprintf(`
<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; max-width: 520px; margin: 0 auto; padding: 24px; color: #1e293b;">
  <h2 style="margin-top: 0; color: #0f172a;">%s</h2>
  <div style="background: #f1f5f9; padding: 16px; border-radius: 8px; margin: 20px 0; border-left: 4px solid #3b82f6;">
    <strong style="display: block; margin-bottom: 6px; color: #334155;">%s:</strong>
    <p style="margin: 0; white-space: pre-wrap; font-size: 15px;">%s</p>
  </div>
  <p style="margin-top: 24px;">
    <a href="%s" style="display: inline-block; background: #2563eb; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 6px; font-weight: 600;">
      Открыть диалог / Open chat
    </a>
  </p>
</div>`, subject, msg.OperatorName, msg.Text, link)

	_ = s.mail.Send(t.Email, subject, body)
}

func (s *Service) InThisChat(t *db.SupportThread) bool {
	return s.tg != nil && t.TelegramTopicID != nil && t.TelegramChatID != nil && *t.TelegramChatID == s.tg.ChatID
}
