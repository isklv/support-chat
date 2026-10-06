package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"support-chat/internal/support"
)

func writeJSON(w http.ResponseWriter, code int, val any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(val)
}

func writeError(w http.ResponseWriter, code int, errCode string) {
	writeJSON(w, code, map[string]string{"error": errCode})
}

func writeSupportFailure(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, support.ErrEmpty):
		writeError(w, http.StatusBadRequest, "empty_message")
	case errors.Is(err, support.ErrTooLong):
		writeError(w, http.StatusBadRequest, "message_too_long")
	case errors.Is(err, support.ErrNeedEmail):
		writeError(w, http.StatusBadRequest, "email_required")
	case errors.Is(err, support.ErrInvalidImage):
		writeError(w, http.StatusBadRequest, "invalid_image")
	case errors.Is(err, support.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found")
	default:
		log.Printf("api: internal error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal_error")
	}
}

func visitorFromRequest(r *http.Request) support.Visitor {
	token := strings.TrimSpace(r.Header.Get("X-Support-Token"))
	if token == "" {
		token = strings.TrimSpace(r.URL.Query().Get("t"))
	}
	userID := strings.TrimSpace(r.Header.Get("X-User-Id"))
	return support.Visitor{
		UserID: userID,
		Token:  token,
	}
}

func readSupportForm(w http.ResponseWriter, r *http.Request) (text string, img []byte, form func(string) string, ok bool) {
	r.Body = http.MaxBytesReader(w, r.Body, support.MaxImageBytes+64<<10)
	if err := r.ParseMultipartForm(support.MaxImageBytes + 64<<10); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_body")
		return "", nil, nil, false
	}
	if file, _, err := r.FormFile("image"); err == nil {
		defer file.Close()
		data, err := io.ReadAll(io.LimitReader(file, support.MaxImageBytes+1))
		if err != nil || len(data) > support.MaxImageBytes {
			writeError(w, http.StatusBadRequest, "invalid_image")
			return "", nil, nil, false
		}
		img = data
	}
	return r.FormValue("text"), img, r.FormValue, true
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "ok",
		"telegram": s.service.Telegram() != nil,
	})
}

func (s *Server) handleSupport(w http.ResponseWriter, r *http.Request) {
	view, err := s.service.Current(r.Context(), visitorFromRequest(r))
	if err != nil {
		writeSupportFailure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (s *Server) handleSupportPost(w http.ResponseWriter, r *http.Request) {
	if !s.rateLimiter.Allow(clientIP(r)) {
		w.Header().Set("Retry-After", "60")
		writeError(w, http.StatusTooManyRequests, "rate_limited")
		return
	}

	text, img, form, ok := readSupportForm(w, r)
	if !ok {
		return
	}

	// Honeypot field for anti-spam bots
	if form("website") != "" {
		writeJSON(w, http.StatusOK, map[string]any{"messages": []any{}})
		return
	}

	var metadata map[string]any
	if metaRaw := form("metadata"); metaRaw != "" {
		_ = json.Unmarshal([]byte(metaRaw), &metadata)
	}

	view, err := s.service.Post(r.Context(), visitorFromRequest(r), support.PostInput{
		Text:      text,
		Email:     form("email"),
		Name:      form("name"),
		Page:      form("page"),
		Lang:      form("lang"),
		UserAgent: r.UserAgent(),
		Image:     img,
		Metadata:  metadata,
	})
	if err != nil {
		writeSupportFailure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (s *Server) handleSupportSeen(w http.ResponseWriter, r *http.Request) {
	if err := s.service.MarkSeen(r.Context(), visitorFromRequest(r)); err != nil {
		writeSupportFailure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleSupportUnread(w http.ResponseWriter, r *http.Request) {
	n, err := s.service.Unread(r.Context(), visitorFromRequest(r))
	if err != nil {
		writeSupportFailure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"unread": n})
}

func (s *Server) handleSupportAttachment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/support/attachments/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid_id")
		return
	}
	a, err := s.service.Attachment(r.Context(), visitorFromRequest(r), false, id)
	if err != nil {
		writeSupportFailure(w, err)
		return
	}
	w.Header().Set("Content-Type", a.ContentType)
	w.Header().Set("Cache-Control", "private, max-age=86400")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(a.Data)
}

func (s *Server) handleTelegramWebhook(w http.ResponseWriter, r *http.Request) {
	tg := s.service.Telegram()
	if tg == nil || r.Header.Get("X-Telegram-Bot-Api-Secret-Token") != tg.WebhookSecret() {
		writeError(w, http.StatusForbidden, "forbidden")
		return
	}
	var u support.Update
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&u); err != nil {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
		return
	}
	if err := s.service.HandleUpdate(r.Context(), u); err != nil {
		log.Printf("api: telegram update %d failed: %v", u.UpdateID, err)
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// Operator console handlers
func (s *Server) handleAdminSupportList(w http.ResponseWriter, r *http.Request) {
	list, err := s.service.List(r.Context(), r.URL.Query().Get("status"), 100)
	if err != nil {
		writeSupportFailure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"threads":  list,
		"telegram": s.service.Telegram() != nil,
	})
}

func (s *Server) handleAdminSupportThread(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "invalid_path")
		return
	}
	id := parts[3]
	view, err := s.service.Get(r.Context(), id)
	if err != nil {
		writeSupportFailure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (s *Server) handleAdminSupportReply(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "invalid_path")
		return
	}
	id := parts[3]
	text, img, form, ok := readSupportForm(w, r)
	if !ok {
		return
	}
	opName := form("operator")
	if opName == "" {
		opName = "Support"
	}
	view, err := s.service.Reply(r.Context(), id, opName, text, img)
	if err != nil {
		writeSupportFailure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (s *Server) handleAdminSupportStatus(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "invalid_path")
		return
	}
	id := parts[3]
	var req struct {
		Open     bool   `json:"open"`
		Operator string `json:"operator"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_body")
		return
	}
	op := req.Operator
	if op == "" {
		op = "Support"
	}
	view, err := s.service.SetStatus(r.Context(), id, req.Open, op)
	if err != nil {
		writeSupportFailure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, view)
}

func (s *Server) handleAdminSupportAttachment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/support/attachments/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid_id")
		return
	}
	a, err := s.service.Attachment(r.Context(), support.Visitor{}, true, id)
	if err != nil {
		writeSupportFailure(w, err)
		return
	}
	w.Header().Set("Content-Type", a.ContentType)
	w.Header().Set("Cache-Control", "private, max-age=86400")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(a.Data)
}
