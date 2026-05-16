package router

import (
	"net/http"
	"time"

	"github.com/Dharshan2208/auth/internal/handlers"
	"github.com/Dharshan2208/auth/internal/middleware"
)

func Register(mux *http.ServeMux, h *handlers.Handler) {
	loginLimiter := middleware.NewRateLimiter(5, time.Minute)

	mux.HandleFunc("/health", h.Health)

	mux.HandleFunc("/signup", h.Signup)
	mux.HandleFunc("/login", loginLimiter.Limit(h.Login))
	mux.HandleFunc("/logout", h.Logout)
	mux.HandleFunc("/refresh", h.Refresh)

	mux.HandleFunc("/profile", middleware.Auth(h.Secret, h.Profile))
	mux.HandleFunc("/admin", middleware.Auth(h.Secret, h.Admin))
}
