package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Dharshan2208/auth/internal/httpx"
)

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	user, err := h.Store.GetUserByID(r.Context(), userID)
	if err != nil {
		slog.Error("profile: failed to fetch user", "user_id", userID, "error", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	resp := map[string]any{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
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
