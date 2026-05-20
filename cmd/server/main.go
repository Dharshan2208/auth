package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dharshan2208/auth/internal/config"
	"github.com/Dharshan2208/auth/internal/handlers"
	"github.com/Dharshan2208/auth/internal/middleware"
	"github.com/Dharshan2208/auth/internal/router"
	"github.com/Dharshan2208/auth/internal/storage"
)

// @title Auth Service API
// @version 1.0.0
// @description A secure authentication service with JWT-based access and refresh tokens. Provides user registration, login, token refresh, logout, password change, and role-based access control.
// @contact.name Dharshan
// @contact.url https://github.com/Dharshan2208/auth
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer:` prefix, e.g. "Bearer abcde12345".
func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))

	cfg := config.Load()

	store, err := storage.New(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to create store", "error", err)
		os.Exit(1)
	}
	defer store.DB.Close()

	mux := http.NewServeMux()
	h := handlers.New(store, cfg)
	rateLimiters := router.Register(mux, h)

	recoveryMux := middleware.Recovery(mux)
	loggedMux := middleware.Logging(middleware.CORS(recoveryMux))
	securedMux := middleware.SecureHeaders(loggedMux)

	// to prevent slowloris attacks
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           securedMux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")

	// rate limiter go routines need to be stopped
	for _, rl := range rateLimiters {
		rl.Stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited gracefully")
}
