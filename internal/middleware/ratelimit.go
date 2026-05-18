package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type clientInfo struct {
	Attempts int
	ResetAt  time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientInfo
	limit   int
	window  time.Duration
	done    chan struct{}
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientInfo),
		limit:   limit,
		window:  window,
		done:    make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Stop() {
	close(rl.done)
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for ip, client := range rl.clients {
				if time.Now().After(client.ResetAt) {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.done:
			return
		}
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)

		rl.mu.Lock()
		defer rl.mu.Unlock()

		client, exists := rl.clients[ip]

		if !exists || time.Now().After(client.ResetAt) {
			client = &clientInfo{
				Attempts: 0,
				ResetAt:  time.Now().Add(rl.window),
			}
			rl.clients[ip] = client
		}

		if client.Attempts >= rl.limit {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "too many requests"})
			return
		}

		client.Attempts++

		next(w, r)
	}
}
