package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

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
		Password string `json:"password"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if req.Username == "" || req.Password == "" {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "username and password are required"})
		return
	}

	_, err := h.Store.GetUserByUsername(req.Username)
	if err == nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "user already exists"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not hash password"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := h.Store.CreateUser(user); err != nil {
		log.Printf("[SIGNUP] FAILED | Username: %s | Error: %+v", req.Username, err)
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create user"})
		return
	}

	log.Printf("[SIGNUP] SUCCESS | Username: %s | IP: %s | Time: %v", req.Username, r.RemoteAddr, time.Since(start))
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
		Password string `json:"password"`
	}

	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	user, err := h.Store.GetUserByUsername(req.Username)
	if err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

	if err := auth.CheckPassword(user.Password, req.Password); err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

	log.Printf("[LOGIN]  SUCCESS | Username: %s | Role: %s | IP: %s | Time: %v", req.Username, user.Role, r.RemoteAddr, time.Since(start))

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

	err = h.Store.SaveRefreshToken(user.ID, hashed, time.Now().Add(h.Cfg.RefreshTokenTTL))
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

	err := h.Store.DeleteRefreshToken(hashed)
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

	userID, err := h.Store.GetUserIDByRefreshToken(hashed)
	if err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "refresh token revoked"})
		return
	}

	user, err := h.Store.GetUserByID(userID)
	if err != nil {
		httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "user not found"})
		return
	}

	_ = h.Store.DeleteRefreshToken(hashed)

	newRefreshToken, err := auth.GenerateRefreshToken(user, h.Secret, h.Cfg.RefreshTokenTTL)
	if err != nil {
		httpx.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
		return
	}

	newHash := auth.HashToken(newRefreshToken)

	err = h.Store.SaveRefreshToken(user.ID, newHash, time.Now().Add(h.Cfg.RefreshTokenTTL))
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
