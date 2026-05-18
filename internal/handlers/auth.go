package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
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
	req.Username = strings.ToLower(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Username == "" || req.Email == "" || req.Password == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email and password are required"})
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid email"})
		return
	}

	if err := auth.ValidateUsername(req.Username); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := auth.ValidatePassword(req.Password); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
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
	req.Username = strings.ToLower(req.Username)
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

	err = h.Store.CreateSession(r.Context(), user.ID, hashed, time.Now().Add(h.Cfg.RefreshTokenTTL), clientIP(r), r.UserAgent())
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

	userIDFromClaims, ok := claimUserID(claims)
	if !ok {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}

	user, err := h.Store.GetUserByID(r.Context(), userIDFromClaims)
	if err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "user not found"})
		return
	}

	// Rotate the refresh token in DB first; if the old token is missing but still JWT-valid,
	// treat it as reuse and revoke sessions for this device.
	now := time.Now()
	newRefreshToken, err := auth.GenerateRefreshToken(user, h.Secret, h.Cfg.RefreshTokenTTL)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
		return
	}

	newHash := auth.HashToken(newRefreshToken)

	_, err = h.Store.RotateSession(r.Context(), hashed, newHash, now.Add(h.Cfg.RefreshTokenTTL), now, clientIP(r), r.UserAgent())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_ = h.Store.RevokeAllSessionsForUserDevice(r.Context(), userIDFromClaims, r.UserAgent())
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
			return
		}
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not rotate session"})
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

func claimUserID(claims jwt.MapClaims) (int, bool) {
	raw, ok := claims["user_id"]
	if !ok {
		return 0, false
	}
	// jwt.MapClaims unmarshals numbers as float64.
	if f, ok := raw.(float64); ok {
		return int(f), true
	}
	if i, ok := raw.(int); ok {
		return i, true
	}
	return 0, false
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
