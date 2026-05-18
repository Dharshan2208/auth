package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"

	"github.com/Dharshan2208/auth/internal/auth"
	"github.com/Dharshan2208/auth/internal/httpx"
	"github.com/Dharshan2208/auth/internal/models"
)

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		httpx.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Username == "" || req.Email == "" || req.Password == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email and password are required"})
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid email"})
		return
	}
	if len(req.Password) < 8 {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "password must be at least 8 characters"})
		return
	}

	if _, err := h.Store.GetUserByUsernameOrEmail(r.Context(), req.Username); err == nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "username already exists"})
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not validate username"})
		return
	}

	if _, err := h.Store.GetUserByUsernameOrEmail(r.Context(), req.Email); err == nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email already exists"})
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not validate email"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not hash password"})
		return
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	if err := h.Store.CreateUser(r.Context(), user); err != nil {
		slog.Error("signup failed",
			"username", req.Username,
			"email", req.Email,
			"error", err,
		)
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create user"})
		return
	}

	slog.Info("signup success",
		"username", req.Username,
		"email", req.Email,
		"ip", r.RemoteAddr,
		"duration", time.Since(start),
	)
	httpx.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		httpx.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	identifier := req.Username
	if req.Email != "" {
		identifier = req.Email
	}
	if identifier == "" || req.Password == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "username/email and password are required"})
		return
	}

	user, err := h.Store.GetUserByUsernameOrEmail(r.Context(), identifier)
	if err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

	if err := auth.CheckPassword(user.PasswordHash, req.Password); err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

	slog.Info("login success",
		"username", user.Username,
		"email", user.Email,
		"role", user.Role,
		"ip", r.RemoteAddr,
		"duration", time.Since(start),
	)

	accessToken, err := auth.GenerateAccessToken(user, h.Secret, h.Cfg.AccessTokenTTL)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user, h.Secret, h.Cfg.RefreshTokenTTL)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	hashed := auth.HashToken(refreshToken)

	err = h.Store.SaveRefreshToken(r.Context(), user.ID, hashed, time.Now().Add(h.Cfg.RefreshTokenTTL))
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not save refresh token"})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	hashed := auth.HashToken(req.RefreshToken)

	err := h.Store.DeleteRefreshToken(r.Context(), hashed)
	if err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	token, err := jwt.Parse(req.RefreshToken,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return h.Secret, nil
		},
	)

	if err != nil || !token.Valid {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "wrong token type"})
		return
	}

	hashed := auth.HashToken(req.RefreshToken)

	// Atomically claim the token ...DELETE ... RETURNING is a single statement,
	// so only one concurrent request can win the race.
	userID, err := h.Store.ConsumeRefreshToken(r.Context(), hashed)
	if err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "refresh token revoked"})
		return
	}

	user, err := h.Store.GetUserByID(r.Context(), userID)
	if err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "user not found"})
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken(user, h.Secret, h.Cfg.RefreshTokenTTL)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
		return
	}

	newHash := auth.HashToken(newRefreshToken)

	err = h.Store.SaveRefreshToken(r.Context(), user.ID, newHash, time.Now().Add(h.Cfg.RefreshTokenTTL))
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not save refresh token"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(user, h.Secret, h.Cfg.AccessTokenTTL)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}
