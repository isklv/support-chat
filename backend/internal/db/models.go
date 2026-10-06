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
