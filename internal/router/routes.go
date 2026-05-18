package router

import (
	"net/http"
	"time"

	"github.com/Dharshan2208/auth/internal/handlers"
	"github.com/Dharshan2208/auth/internal/middleware"
)

func Register(mux *http.ServeMux, h *handlers.Handler) []*middleware.RateLimiter {
	loginLimiter := middleware.NewRateLimiter(5, time.Minute)
	signupLimiter := middleware.NewRateLimiter(3, time.Minute)
	refreshLimiter := middleware.NewRateLimiter(3, time.Minute)

	mux.HandleFunc("/health", h.Health)

	mux.HandleFunc("/signup", signupLimiter.Limit(h.Signup))
	mux.HandleFunc("/login", loginLimiter.Limit(h.Login))
	mux.HandleFunc("/refresh", refreshLimiter.Limit(h.Refresh))
	mux.HandleFunc("/logout", h.Logout)

	mux.HandleFunc("/profile", middleware.Auth(h.Secret, h.Profile))
	mux.HandleFunc("/admin", middleware.Auth(h.Secret, h.Admin))

	return []*middleware.RateLimiter{loginLimiter, signupLimiter, refreshLimiter}
}
