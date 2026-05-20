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

// Signup godoc
// @Summary Register a new user
// @Description Create a new user account with username, email, and password. The username and email are normalized to lowercase. Password must be 12-72 characters with uppercase, lowercase, digit, and symbol.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body SignupRequest true "User registration details"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Router /signup [post]
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Username = strings.ToLower(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	if req.Username == "" || req.Email == "" || req.Password == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "missing required fields"})
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
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "username already taken"})
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	if _, err := h.Store.GetUserByUsernameOrEmail(r.Context(), req.Email); err == nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email already registered"})
		return
	} else if !errors.Is(err, pgx.ErrNoRows) {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
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
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
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

// Login godoc
// @Summary Authenticate a user
// @Description Authenticate with username/email and password. Returns a pair of access and refresh tokens on success.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} TokenPairResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Router /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	req.Identifier = strings.TrimSpace(req.Identifier)
	req.Identifier = strings.ToLower(req.Identifier)

	if req.Identifier == "" || req.Password == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "missing required fields"})
		return
	}

	user, err := h.Store.GetUserByUsernameOrEmail(r.Context(), req.Identifier)
	if err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

	if err := auth.CheckPassword(user.PasswordHash, req.Password); err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

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

// Logout godoc
// @Summary Logout and revoke a refresh token
// @Description Revoke the provided refresh token, effectively logging out the session.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LogoutRequest true "Refresh token to revoke"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
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

// Refresh godoc
// @Summary Refresh an access token
// @Description Exchange a valid refresh token for a new access token and a rotated refresh token. If a refresh token is reused (i.e., the old one was already rotated), all sessions for that device are revoked (token theft protection).
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Current refresh token"
// @Success 200 {object} TokenPairResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Router /refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
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

// ChangePassword godoc
// @Summary Change the authenticated user's password
// @Description Change the password for the currently authenticated user. The old password must be confirmed twice. On success, all other sessions are revoked and the user must re-authenticate on other devices.
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Router /password/change [post]
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok || userID <= 0 {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req struct {
		OldPassword        string `json:"old_password"`
		ConfirmOldPassword string `json:"confirm_old_password"`
		NewPassword        string `json:"new_password"`
	}
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.OldPassword == "" || req.ConfirmOldPassword == "" || req.NewPassword == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "missing required fields"})
		return
	}
	if req.OldPassword != req.ConfirmOldPassword {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "old password confirmation does not match"})
		return
	}
	if req.NewPassword == req.OldPassword {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "new password must be different"})
		return
	}
	if err := auth.ValidatePassword(req.NewPassword); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	user, err := h.Store.GetUserByID(r.Context(), userID)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	if err := auth.CheckPassword(user.PasswordHash, req.OldPassword); err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

	newHash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	if err := h.Store.UpdateUserPasswordHash(r.Context(), userID, newHash); err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}

	// Revoke refresh tokens so any other logged-in devices must re-authenticate.
	_ = h.Store.RevokeAllSessionsForUser(r.Context(), userID)

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "password updated"})
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
