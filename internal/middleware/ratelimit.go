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
