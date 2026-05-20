package router

import (
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Dharshan2208/auth/internal/handlers"
	"github.com/Dharshan2208/auth/internal/middleware"

	_ "github.com/Dharshan2208/auth/docs" // swagger generated docs
)

func Register(mux *http.ServeMux, h *handlers.Handler) []*middleware.RateLimiter {
	loginLimiter := middleware.NewRateLimiter(5, time.Minute)
	signupLimiter := middleware.NewRateLimiter(3, time.Minute)
	refreshLimiter := middleware.NewRateLimiter(3, time.Minute)
	logoutLimiter := middleware.NewRateLimiter(3, time.Minute)
	changePasswordLimiter := middleware.NewRateLimiter(3, time.Minute)

	// API v1 routes
	mux.HandleFunc("GET /api/v1/health", h.Health)

	mux.HandleFunc("POST /api/v1/signup", signupLimiter.Limit(h.Signup))
	mux.HandleFunc("POST /api/v1/login", loginLimiter.Limit(h.Login))
	mux.HandleFunc("POST /api/v1/refresh", refreshLimiter.Limit(h.Refresh))
	mux.HandleFunc("POST /api/v1/logout", logoutLimiter.Limit(h.Logout))

	mux.HandleFunc("GET /api/v1/profile", middleware.Auth(h.Secret, h.Profile))
	mux.HandleFunc("GET /api/v1/admin", middleware.Auth(h.Secret, h.Admin))
	mux.HandleFunc("POST /api/v1/password/change", changePasswordLimiter.Limit(middleware.Auth(h.Secret, h.ChangePassword)))

	// Swagger documentation UI
	mux.Handle("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	return []*middleware.RateLimiter{loginLimiter, signupLimiter, refreshLimiter, changePasswordLimiter}
}
