package handlers

import (
	"github.com/Dharshan2208/auth/internal/config"
	"github.com/Dharshan2208/auth/internal/storage"
)

type Handler struct {
	Secret []byte
	Store  *storage.Store
	Cfg    *config.Config
}

func New(store *storage.Store, cfg *config.Config) *Handler {
	return &Handler{
		Store:  store,
		Cfg:    cfg,
		Secret: []byte(cfg.JWTSecret),
	}
}
