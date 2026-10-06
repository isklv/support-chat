package support

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Telegram struct {
	Token  string
	ChatID int64
	Base   string
	Client *http.Client
}

func NewTelegram(token string, chatID int64, base, proxy string) *Telegram {
	if strings.TrimSpace(token) == "" || chatID == 0 {
		return nil
	}
	if base == "" {
		base = "https://api.telegram.org"
	}
	client := &http.Client{Timeout: 35 * time.Second}
	if proxy != "" {
		if u, err := url.Parse(proxy); err == nil {
			transport := http.DefaultTransport.(*http.Transport).Clone()
			transport.Proxy = http.ProxyURL(u)
			client.Transport = transport
		}
	}
	return &Telegram{
		Token:  token,
		ChatID: chatID,
		Base:   strings.TrimRight(base, "/"),
		Client: client,
	}
}

func (t *Telegram) WebhookSecret() string {
	sum := sha256.Sum256([]byte("support-chat:" + t.Token))
	return hex.EncodeToString(sum[:16])
}

// Telegram Bot API data models
type tgUser struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type tgChat struct {
	ID int64 `json:"id"`
}

type tgPhoto struct {
	FileID   string `json:"file_id"`
	FileSize int64  `json:"file_size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type tgDocument struct {
	FileID   string `json:"file_id"`
	MimeType string `json:"mime_type"`
	FileSize int64  `json:"file_size"`
}

type tgMessage struct {
	MessageID         int64            `json:"message_id"`
	MessageThreadID   int64            `json:"message_thread_id"`
	IsTopicMessage    bool             `json:"is_topic_message"`
	From              *tgUser          `json:"from"`
	Chat              tgChat           `json:"chat"`
	Text              string           `json:"text"`
	Caption           string           `json:"caption"`
	Photo             []tgPhoto        `json:"photo"`
	Document          *tgDocument      `json:"document"`
	ReplyTo           *tgMessage       `json:"reply_to_message"`
	ForumTopicCreated *json.RawMessage `json:"forum_topic_created"`
}

type Update struct {
	UpdateID int64      `json:"update_id"`
	Message  *tgMessage `json:"message"`
}

type apiResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
}

func (t *Telegram) call(ctx context.Context, method string, params map[string]any, out any) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/bot%s/%s", t.Base, t.Token, method), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return t.do(req, method, out)
}

func (t *Telegram) do(req *http.Request, method string, out any) error {
	resp, err := t.Client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram %s: %w", method, err)
	}
	defer resp.Body.Close()
	var r apiResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&r); err != nil {
		return fmt.Errorf("telegram %s: %s: %w", method, resp.Status, err)
	}
	if !r.OK {
		return fmt.Errorf("telegram %s: %s", method, r.Description)
	}
	if out != nil {
		return json.Unmarshal(r.Result, out)
	}
	return nil
}

// CreateTopic opens a forum topic in a supergroup and returns its id.
func (t *Telegram) CreateTopic(ctx context.Context, name string) (int64, error) {
	var out struct {
		MessageThreadID int64 `json:"message_thread_id"`
	}
	runes := []rune(name)
	if len(runes) > 128 {
		name = string(runes[:128])
	}
	err := t.call(ctx, "createForumTopic", map[string]any{"chat_id": t.ChatID, "name": name}, &out)
	return out.MessageThreadID, err
}

// SetTopicOpen closes or reopens a topic.
func (t *Telegram) SetTopicOpen(ctx context.Context, topic int64, open bool) error {
	method := "closeForumTopic"
	if open {
		method = "reopenForumTopic"
	}
	return t.call(ctx, method, map[string]any{"chat_id": t.ChatID, "message_thread_id": topic}, nil)
}

// Send posts text into a topic (0: main chat) and returns message id.
func (t *Telegram) Send(ctx context.Context, topic int64, text string) (int64, error) {
	params := map[string]any{
		"chat_id":                  t.ChatID,
		"text":                     text,
		"disable_web_page_preview": true,
	}
	if topic != 0 {
		params["message_thread_id"] = topic
	}
	var out struct {
		MessageID int64 `json:"message_id"`
	}
	err := t.call(ctx, "sendMessage", params, &out)
	return out.MessageID, err
}

// SendPhoto posts an image with a caption into a topic.
func (t *Telegram) SendPhoto(ctx context.Context, topic int64, caption string, data []byte, contentType string) (int64, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("chat_id", strconv.FormatInt(t.ChatID, 10))
	if topic != 0 {
		_ = w.WriteField("message_thread_id", strconv.FormatInt(topic, 10))
	}
	if caption != "" {
		runes := []rune(caption)
		if len(runes) > 1000 {
			caption = string(runes[:1000]) + "…"
		}
		_ = w.WriteField("caption", caption)
	}
	ext := "jpg"
	if strings.Contains(contentType, "png") {
		ext = "png"
	} else if strings.Contains(contentType, "webp") {
		ext = "webp"
	}
	part, err := w.CreateFormFile("photo", "image."+ext)
	if err != nil {
		return 0, err
	}
	if _, err := part.Write(data); err != nil {
		return 0, err
	}
	if err := w.Close(); err != nil {
		return 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/bot%s/sendPhoto", t.Base, t.Token), &buf)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	var out struct {
		MessageID int64 `json:"message_id"`
	}
	err = t.do(req, "sendPhoto", &out)
	return out.MessageID, err
}

// Download fetches a file an operator sent, up to max bytes.
func (t *Telegram) Download(ctx context.Context, fileID string, max int64) ([]byte, error) {
	var file struct {
		FilePath string `json:"file_path"`
		FileSize int64  `json:"file_size"`
	}
	if err := t.call(ctx, "getFile", map[string]any{"file_id": fileID}, &file); err != nil {
		return nil, err
	}
	if file.FileSize > max {
		return nil, fmt.Errorf("telegram file: %d bytes is too big", file.FileSize)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/file/bot%s/%s", t.Base, t.Token, file.FilePath), nil)
	if err != nil {
		return nil, err
	}
	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("telegram file: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram file: status %s", resp.Status)
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, max+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > max {
		return nil, fmt.Errorf("telegram file: exceeds max allowed size")
	}
	return data, nil
}

// SetWebhook points Telegram at our endpoint with secret token.
func (t *Telegram) SetWebhook(ctx context.Context, url string) error {
	return t.call(ctx, "setWebhook", map[string]any{
		"url":             url,
		"secret_token":    t.WebhookSecret(),
		"allowed_updates": []string{"message"},
	}, nil)
}

// DeleteWebhook removes active webhook (required before polling).
func (t *Telegram) DeleteWebhook(ctx context.Context) error {
	return t.call(ctx, "deleteWebhook", map[string]any{
		"drop_pending_updates": false,
	}, nil)
}

// GetUpdates polls Telegram for updates with long-polling timeout.
func (t *Telegram) GetUpdates(ctx context.Context, offset int64, timeoutSec int) ([]Update, error) {
	params := map[string]any{
		"timeout":         timeoutSec,
		"allowed_updates": []string{"message"},
	}
	if offset > 0 {
		params["offset"] = offset
	}
	var updates []Update
	err := t.call(ctx, "getUpdates", params, &updates)
	return updates, err
}

func displayName(u *tgUser) string {
	if u == nil {
		return ""
	}
	name := strings.TrimSpace(u.FirstName + " " + u.LastName)
	if name == "" {
		name = u.Username
	}
	return name
}
