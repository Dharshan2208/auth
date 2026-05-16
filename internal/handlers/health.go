package handlers

import (
	"net/http"

	"github.com/Dharshan2208/auth/internal/httpx"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "running",
	})
}
