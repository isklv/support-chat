package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"support-chat/internal/config"
	"support-chat/internal/support"
)

type Server struct {
	cfg         *config.Config
	service     *support.Service
	rateLimiter *ipRateLimiter
}

func NewServer(cfg *config.Config, service *support.Service) *Server {
	return &Server{
		cfg:         cfg,
		service:     service,
		rateLimiter: newRateLimiter(60, 10*time.Minute),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/healthz", s.handleHealthz)

	// Visitor endpoints
	mux.HandleFunc("/api/v1/support", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.handleSupport(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	})

	mux.HandleFunc("/api/v1/support/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.handleSupportPost(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	})

	mux.HandleFunc("/api/v1/support/seen", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.handleSupportSeen(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	})

	mux.HandleFunc("/api/v1/support/unread", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.handleSupportUnread(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	})

	mux.HandleFunc("/api/v1/support/attachments/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.handleSupportAttachment(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	})

	// Telegram Webhook
	mux.HandleFunc("/api/v1/support/telegram/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.handleTelegramWebhook(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	})

	// Operator / Admin endpoints
	mux.HandleFunc("/api/v1/admin/support", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.handleAdminSupportList(w, r)
		} else {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	})

	mux.HandleFunc("/api/v1/admin/support/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/messages") && r.Method == http.MethodPost:
			s.handleAdminSupportReply(w, r)
		case strings.HasSuffix(path, "/status") && r.Method == http.MethodPost:
			s.handleAdminSupportStatus(w, r)
		case strings.Contains(path, "/attachments/") && r.Method == http.MethodGet:
			s.handleAdminSupportAttachment(w, r)
		case r.Method == http.MethodGet:
			s.handleAdminSupportThread(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		}
	})

	// Static files: widget.js, styles, demo page
	staticDir := s.cfg.StaticDir
	if info, err := os.Stat(staticDir); err == nil && info.IsDir() {
		fileServer := http.FileServer(http.Dir(staticDir))

		// Specific helper routes for widget.js and widget.css
		mux.HandleFunc("/widget.js", func(w http.ResponseWriter, r *http.Request) {
			target := filepath.Join(staticDir, "widget.js")
			if _, err := os.Stat(target); err == nil {
				w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
				w.Header().Set("Cache-Control", "no-cache")
				http.ServeFile(w, r, target)
				return
			}
			http.NotFound(w, r)
		})

		mux.HandleFunc("/widget.css", func(w http.ResponseWriter, r *http.Request) {
			target := filepath.Join(staticDir, "widget.css")
			if _, err := os.Stat(target); err == nil {
				w.Header().Set("Content-Type", "text/css; charset=utf-8")
				w.Header().Set("Cache-Control", "no-cache")
				http.ServeFile(w, r, target)
				return
			}
			http.NotFound(w, r)
		})

		mux.Handle("/", fileServer)
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				fmt.Fprint(w, `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Support Chat API</title></head>
<body style="font-family: sans-serif; padding: 40px; line-height: 1.6;">
  <h1>Support Chat Backend</h1>
  <p>Status: Running</p>
  <p>API Endpoint: <code>/api/v1/support</code></p>
</body>
</html>`)
				return
			}
			http.NotFound(w, r)
		})
	}

	return corsMiddleware(s.cfg, mux)
}
