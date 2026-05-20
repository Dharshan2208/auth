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
	logoutLimiter := middleware.NewRateLimiter(3, time.Minute)
	changePasswordLimiter := middleware.NewRateLimiter(3, time.Minute)

	mux.HandleFunc("GET /health", h.Health)

	mux.HandleFunc("POST /signup", signupLimiter.Limit(h.Signup))
	mux.HandleFunc("POST /login", loginLimiter.Limit(h.Login))
	mux.HandleFunc("POST /refresh", refreshLimiter.Limit(h.Refresh))
	mux.HandleFunc("POST /logout", logoutLimiter.Limit(h.Logout))

	mux.HandleFunc("GET /profile", middleware.Auth(h.Secret, h.Profile))
	mux.HandleFunc("GET /admin", middleware.Auth(h.Secret, h.Admin))
	mux.HandleFunc("POST /password/change", changePasswordLimiter.Limit(middleware.Auth(h.Secret, h.ChangePassword)))

	return []*middleware.RateLimiter{loginLimiter, signupLimiter, refreshLimiter, changePasswordLimiter}
}
