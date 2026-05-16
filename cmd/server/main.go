package main

import (
	"log"
	"net/http"

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

	mux := http.NewServeMux()
	h := handlers.New(store, cfg)
	router.Register(mux, h)

	loggedMux := middleware.Logging(middleware.CORS(mux))

	log.Println("Server running on :", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, loggedMux))
}
