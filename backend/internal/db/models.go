package db

import (
	"time"
)

const (
	SupportOpen   = "open"
	SupportClosed = "closed"

	AuthorUser     = "user"
	AuthorOperator = "operator"
	AuthorSystem   = "system"
)

// SupportThread represents a conversation thread.
type SupportThread struct {
	ID              string     `gorm:"primaryKey;size:36" json:"id"`
	Number          int64      `gorm:"index;autoIncrement" json:"number"`
	BotID           *string    `gorm:"size:36;index" json:"bot_id,omitempty"`
	Email           string     `gorm:"index;size:255" json:"email"`
	Name            string     `gorm:"size:255" json:"name"`
	Locale          string     `gorm:"size:10;default:'en'" json:"locale"`
	Status          string     `gorm:"size:20;default:'open';index" json:"status"`
	Context         string     `gorm:"type:text" json:"context"` // JSON string
	TelegramChatID  *int64     `json:"-"`
	TelegramTopicID *int64     `gorm:"index" json:"-"`
	UserSeenAt      *time.Time `json:"user_seen_at,omitempty"`
	EmailedAt       *time.Time `json:"-"`
	LastMessageAt   time.Time  `gorm:"index" json:"last_message_at"`
	CreatedAt       time.Time  `gorm:"index" json:"created_at"`
	ClosedAt        *time.Time `json:"closed_at,omitempty"`
}

// Bot represents a Telegram bot configuration for support chat.
type Bot struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Token       string     `gorm:"size:255;not null" json:"token"`
	Username    string     `gorm:"size:255" json:"username"`
	FirstName   string     `gorm:"size:255" json:"first_name"`
	ChatID      int64      `json:"chat_id"`
	Mode        string     `gorm:"size:20;default:'polling'" json:"mode"` // polling or webhook
	WebhookURL  string     `gorm:"size:500" json:"webhook_url"`
	BaseURL     string     `gorm:"size:500" json:"base_url"`
	Proxy       string     `gorm:"size:500" json:"proxy"`
	IsDefault   bool       `gorm:"default:false;index" json:"is_default"`
	IsActive    bool       `gorm:"default:true;index" json:"is_active"`
	Description string     `gorm:"size:500" json:"description"`
	CreatedAt   time.Time  `gorm:"index" json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastSeenAt  *time.Time `json:"last_seen_at,omitempty"`
}

// SupportMessage represents a single message in a thread.
type SupportMessage struct {
	ID                string     `gorm:"primaryKey;size:36" json:"id"`
	ThreadID          string     `gorm:"index;size:36" json:"thread_id"`
	Author            string     `gorm:"size:20;index" json:"author"` // user, operator, system
	OperatorName      string     `gorm:"size:100" json:"operator_name,omitempty"`
	Text              string     `gorm:"type:text" json:"text"`
	HasAttachment     bool       `json:"has_attachment"`
	TelegramMessageID *int64     `gorm:"index" json:"-"`
	RelayedAt         *time.Time `gorm:"index" json:"-"`
	CreatedAt         time.Time  `gorm:"index" json:"created_at"`
}

// SupportAttachment holds binary screenshot / image data.
type SupportAttachment struct {
	MessageID   string `gorm:"primaryKey;size:36"`
	ContentType string `gorm:"size:100"`
	Data        []byte `gorm:"type:blob"`
}

// SupportToken links a visitor token to their thread.
type SupportToken struct {
	TokenHash string    `gorm:"primaryKey;size:64"`
	ThreadID  string    `gorm:"index;size:36"`
	CreatedAt time.Time `json:"created_at"`
}
