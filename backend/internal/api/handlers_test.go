package api_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
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
		StaticDir:      "../../static",
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

func TestAPI_AdminBotsCRUD(t *testing.T) {
	handler := setupTestServer(t)

	// 1. Initial list should be empty
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/bots", nil)
	listRec := httptest.NewRecorder()
	handler.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", listRec.Code)
	}
	var listResp struct {
		Bots []db.Bot `json:"bots"`
	}
	if err := json.NewDecoder(listRec.Body).Decode(&listResp); err != nil {
		t.Fatalf("failed to decode list: %v", err)
	}
	if len(listResp.Bots) != 0 {
		t.Fatalf("expected 0 bots, got %d", len(listResp.Bots))
	}

	// 2. Create Bot
	createPayload := []byte(`{
		"name": "Тестовый Бот",
		"token": "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		"chat_id": -100123456789,
		"mode": "polling",
		"is_active": true,
		"is_default": true,
		"description": "Основной бот для тестов"
	}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/bots", bytes.NewReader(createPayload))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	handler.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d: %s", createRec.Code, createRec.Body.String())
	}

	var createResp struct {
		Bot db.Bot `json:"bot"`
	}
	if err := json.NewDecoder(createRec.Body).Decode(&createResp); err != nil {
		t.Fatalf("failed to decode created bot: %v", err)
	}
	botID := createResp.Bot.ID
	if botID == "" || createResp.Bot.Name != "Тестовый Бот" {
		t.Fatalf("invalid created bot: %+v", createResp.Bot)
	}
	if !createResp.Bot.IsDefault {
		t.Errorf("expected first created bot to be default")
	}

	// 3. Get Bot
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/bots/"+botID, nil)
	getRec := httptest.NewRecorder()
	handler.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	// 4. Update Bot
	updatePayload := []byte(`{
		"name": "Обновленный Бот",
		"description": "Новое описание"
	}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/admin/bots/"+botID, bytes.NewReader(updatePayload))
	updateReq.Header.Set("Content-Type", "application/json")
	updateRec := httptest.NewRecorder()
	handler.ServeHTTP(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", updateRec.Code, updateRec.Body.String())
	}
	var updateResp struct {
		Bot db.Bot `json:"bot"`
	}
	_ = json.NewDecoder(updateRec.Body).Decode(&updateResp)
	if updateResp.Bot.Name != "Обновленный Бот" {
		t.Errorf("expected name 'Обновленный Бот', got %s", updateResp.Bot.Name)
	}

	// 5. Toggle active
	toggleReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/bots/"+botID+"/toggle", nil)
	toggleRec := httptest.NewRecorder()
	handler.ServeHTTP(toggleRec, toggleReq)
	if toggleRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", toggleRec.Code)
	}
	var toggleResp struct {
		Bot db.Bot `json:"bot"`
	}
	_ = json.NewDecoder(toggleRec.Body).Decode(&toggleResp)
	if toggleResp.Bot.IsActive {
		t.Errorf("expected bot to be inactive after toggle")
	}

	// 6. Delete Bot
	delReq := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/bots/"+botID, nil)
	delRec := httptest.NewRecorder()
	handler.ServeHTTP(delRec, delReq)
	if delRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", delRec.Code, delRec.Body.String())
	}

	// Verify it's gone
	getDeletedReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/bots/"+botID, nil)
	getDeletedRec := httptest.NewRecorder()
	handler.ServeHTTP(getDeletedRec, getDeletedReq)
	if getDeletedRec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for deleted bot, got %d", getDeletedRec.Code)
	}
}

func TestAPI_AdminWebRoute(t *testing.T) {
	handler := setupTestServer(t)

	// 1. GET /admin should redirect to /admin/
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Errorf("expected 302 Found, got %d", rec.Code)
	}

	// 2. GET /admin/ should return 200 OK and HTML
	reqSlash := httptest.NewRequest(http.MethodGet, "/admin/", nil)
	recSlash := httptest.NewRecorder()
	handler.ServeHTTP(recSlash, reqSlash)
	if recSlash.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recSlash.Code)
	}
	body := recSlash.Body.String()
	if !strings.Contains(body, "Управление ботами") {
		t.Errorf("expected admin page title in response, got %s", body)
	}
}
