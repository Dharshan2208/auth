package handlers

import (
	"net/http"

	"github.com/Dharshan2208/auth/internal/httpx"
)

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "welcome to profile",
	})
}

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role").(string)

	if role != "admin" {
		httpx.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "welcome admin",
	})
}
