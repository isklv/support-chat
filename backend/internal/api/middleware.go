package api

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"support-chat/internal/config"
)

func corsMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	allowAll := len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "*"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := false

		if allowAll {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			allowed = true
		} else if origin != "" {
			for _, o := range cfg.AllowedOrigins {
				if strings.EqualFold(o, origin) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
					allowed = true
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Support-Token, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			if allowed || origin == "" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Simple IP rate limiter
type ipRateLimiter struct {
	mu      sync.Mutex
	limits  map[string][]time.Time
	max     int
	window  time.Duration
	cleanup time.Time
}

func newRateLimiter(max int, window time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		limits: make(map[string][]time.Time),
		max:    max,
		window: window,
	}
}

func (rl *ipRateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// Periodic cleanup every minute
	if now.Sub(rl.cleanup) > time.Minute {
		cutoff := now.Add(-rl.window)
		for k, times := range rl.limits {
			valid := times[:0]
			for _, t := range times {
				if t.After(cutoff) {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.limits, k)
			} else {
				rl.limits[k] = valid
			}
		}
		rl.cleanup = now
	}

	cutoff := now.Add(-rl.window)
	times := rl.limits[ip]
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.max {
		rl.limits[ip] = valid
		return false
	}

	rl.limits[ip] = append(valid, now)
	return true
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}
	parts := strings.Split(r.RemoteAddr, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return r.RemoteAddr
}
