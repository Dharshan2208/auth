package middleware

import (
	"net"
	"net/http"
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
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*clientInfo),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

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
			http.Error(w, "too many attempts", http.StatusTooManyRequests)
			return
		}

		client.Attempts++

		next(w, r)
	}
}
