package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Dharshan2208/auth/internal/httpx"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.Ping(r.Context()); err != nil {
		slog.Warn("health check failed", "error", err)
		httpx.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "unavailable",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "running",
	})
}
