package main

import (
	"context"
	"log"
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

func main() {
	cfg := config.Load()

	store, err := storage.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer store.DB.Close()

	mux := http.NewServeMux()
	h := handlers.New(store, cfg)
	router.Register(mux, h)

	loggedMux := middleware.Logging(middleware.CORS(mux))

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: loggedMux,
	}

	go func() {
		log.Println("Server running on :", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
