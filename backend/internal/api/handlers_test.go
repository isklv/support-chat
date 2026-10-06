package api_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"support-chat/internal/api"
	"support-chat/internal/config"
	"support-chat/internal/db"
	"support-chat/internal/mailer"
	"support-chat/internal/support"
)

func setupTestServer(t *testing.T) http.Handler {
	t.Helper()
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	database, err := db.Init("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	cfg := &config.Config{
		AllowedOrigins: []string{"*"},
	}
	mailClient := mailer.New(cfg)
	service := support.NewService(database, cfg, mailClient, nil)
	srv := api.NewServer(cfg, service)
	return srv.Routes()
}

func TestAPI_PostAndFetchMessages(t *testing.T) {
	handler := setupTestServer(t)

	// 1. Post message with multipart form
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	_ = w.WriteField("text", "Can you help me with billing?")
	_ = w.WriteField("email", "john@example.com")
	_ = w.WriteField("name", "John Doe")
	_ = w.WriteField("page", "/dashboard")
	_ = w.WriteField("lang", "en")
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/support/messages", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var res struct {
		Thread *struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"thread"`
		Messages []struct {
			Text string `json:"text"`
		} `json:"messages"`
		Token string `json:"token"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Token == "" {
		t.Fatal("expected token in response")
	}
	if len(res.Messages) != 1 || res.Messages[0].Text != "Can you help me with billing?" {
		t.Errorf("unexpected message response")
	}

	// 2. Fetch current conversation using X-Support-Token
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/support", nil)
	getReq.Header.Set("X-Support-Token", res.Token)
	getRec := httptest.NewRecorder()

	handler.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", getRec.Code, getRec.Body.String())
	}

	// 3. Check honeypot field drops bot messages quietly
	var botBody bytes.Buffer
	bw := multipart.NewWriter(&botBody)
	_ = bw.WriteField("text", "Buy cheap pills!")
	_ = bw.WriteField("email", "spammer@bot.net")
	_ = bw.WriteField("website", "http://spam.com") // Honeypot!
	_ = bw.Close()

	botReq := httptest.NewRequest(http.MethodPost, "/api/v1/support/messages", &botBody)
	botReq.Header.Set("Content-Type", bw.FormDataContentType())
	botRec := httptest.NewRecorder()

	handler.ServeHTTP(botRec, botReq)
	if botRec.Code != http.StatusOK {
		t.Errorf("expected 200 fake ok for honeypot, got %d", botRec.Code)
	}
}
