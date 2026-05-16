package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Dharshan2208/auth/internal/auth"
	"github.com/Dharshan2208/auth/internal/config"
	"github.com/Dharshan2208/auth/internal/middleware"
	"github.com/Dharshan2208/auth/internal/models"
	"github.com/Dharshan2208/auth/internal/storage"
)

type Server struct {
	secret []byte
	store  *storage.Store
	cfg    *config.Config
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if req.Username == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username and password are required"})
		return
	}

	_, err := s.store.GetUserByUsername(req.Username)

	if err == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user already exists"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not hash password"})
		return
	}

	role := "user"
	if req.Username == "admin" {
		role = "admin"
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.store.CreateUser(user); err != nil {
		log.Printf("[SIGNUP] FAILED | Username: %s | Error: %+v", req.Username, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create user"})
		return
	}

	log.Printf("[SIGNUP] SUCCESS | Username: %s | IP: %s | Time: %v", req.Username, r.RemoteAddr, time.Since(start))
	writeJSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	user, err := s.store.GetUserByUsername(
		req.Username,
	)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

	log.Printf("[LOGIN]  SUCCESS | Username: %s | Role: %s | IP: %s | Time: %v", req.Username, user.Role, r.RemoteAddr, time.Since(start))

	accessToken, err := auth.GenerateAccessToken(user, s.secret, s.cfg.AccessTokenTTL)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user, s.secret, s.cfg.RefreshTokenTTL)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	hashed := auth.HashToken(refreshToken)

	err = s.store.SaveRefreshToken(user.ID, hashed, time.Now().Add(s.cfg.RefreshTokenTTL))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not save refresh token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	hashed := auth.HashToken(req.RefreshToken)

	err := s.store.DeleteRefreshToken(hashed)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

func (s *Server) refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	token, err := jwt.Parse(req.RefreshToken,
		func(token *jwt.Token) (any, error) {
			return s.secret, nil
		},
	)

	if err != nil || !token.Valid {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "wrong token type"})
		return
	}

	hashed := auth.HashToken(req.RefreshToken)

	userID, err := s.store.GetUserIDByRefreshToken(hashed)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "refresh token revoked"})
		return
	}

	user, err := s.store.GetUserByID(userID)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "user not found"})
		return
	}

	_ = s.store.DeleteRefreshToken(hashed)

	newRefreshToken, err := auth.GenerateRefreshToken(user, s.secret, s.cfg.AccessTokenTTL)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
		return
	}

	newHash := auth.HashToken(newRefreshToken)

	err = s.store.SaveRefreshToken(user.ID, newHash, time.Now().Add(s.cfg.RefreshTokenTTL))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not save refresh token"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(user, s.secret, s.cfg.AccessTokenTTL)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "running",
	})
}

func (s *Server) profile(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "welcome to profile",
	})
}

func (s *Server) admin(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role").(string)

	if role != "admin" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "welcome admin",
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	cfg := config.Load()

	store, err := storage.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	server := &Server{
		store:  store,
		cfg:    cfg,
		secret: []byte(cfg.JWTSecret),
	}

	mux := http.NewServeMux()
	loginLimiter := middleware.NewRateLimiter(5, time.Minute)

	mux.HandleFunc("/health", server.health)
	mux.HandleFunc("/signup", server.signup)
	mux.HandleFunc("/login", loginLimiter.Limit(server.login))
	mux.HandleFunc("/logout", server.logout)
	mux.HandleFunc("/refresh", server.refresh)

	mux.HandleFunc("/profile", middleware.Auth(server.secret, server.profile))
	mux.HandleFunc("/admin", middleware.Auth(server.secret, server.admin))

	loggedMux := middleware.Logging(middleware.CORS(mux))

	log.Println("Server running on :", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, loggedMux))
}
