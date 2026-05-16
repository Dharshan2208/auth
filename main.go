package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type Server struct {
	users         map[string]User
	secret        []byte
	refreshTokens map[string]string
}

type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
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

	if _, exists := s.users[req.Username]; exists {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not hash password"})
		return
	}

	role := "user"
	if req.Username == "admin" {
		role = "admin"
	}

	user := User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     role,
	}

	s.users[req.Username] = user

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

	user, exists := s.users[req.Username]
	if !exists {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid creds"})
		return
	}

	log.Printf("[LOGIN]  SUCCESS | Username: %s | Role: %s | IP: %s | Time: %v", req.Username, user.Role, r.RemoteAddr, time.Since(start))
	// writeJSON(w, http.StatusOK, map[string]string{"message": "login successful"})

	accessToken, err := s.generateJWT(user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	s.refreshTokens[refreshToken] = user.Username

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

	if _, exists := s.refreshTokens[req.RefreshToken]; !exists {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}

	delete(s.refreshTokens, req.RefreshToken)
	writeJSON(w, http.StatusOK, map[string]string{"message": "looged out successfully.."})
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

	username, exists := s.refreshTokens[req.RefreshToken]
	if !exists {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "refresh token revoked"})
		return
	}

	user := s.users[username]

	accessToken, err := s.generateJWT(user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"access_token": accessToken})
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return s.secret, nil
		})

		if err != nil || !token.Valid {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		r.Header.Set("username", claims["username"].(string))
		r.Header.Set("role", claims["role"].(string))

		next(w, r)
	}
}

// Logging middleware
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := uuid.New().String()[:8]

		// Attaching requestID to context and headers
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		r = r.WithContext(ctx)
		w.Header().Set("X-Request-ID", requestID)

		rw := &responseWriter{ResponseWriter: w, statusCode: 0, wroteHeader: false}

		log.Printf(`[REQ] %s | %s %s | IP: %s | RequestID: %s`,
			r.Proto, r.Method, r.URL.Path, r.RemoteAddr, requestID)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Printf(`[RES] %s | %s %s | Status: %d | Duration: %v | RequestID: %s`,
			r.Proto, r.Method, r.URL.Path, rw.statusCode, duration, requestID)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
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
	role := r.Header.Get("role")

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

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return // prevent multiple WriteHeader calls
	}
	rw.statusCode = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK) // default only if nothing was set
	}
	return rw.ResponseWriter.Write(b)
}

func (s *Server) generateJWT(user User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// this creates cryptographic proof
	return token.SignedString(s.secret)
}

func (s *Server) generateRefreshToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type":     "refresh",
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(s.secret)
}

func main() {
	server := &Server{
		users:         make(map[string]User),
		secret:        []byte("for-now-random-key-hardcoded"),
		refreshTokens: make(map[string]string),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", server.health)
	mux.HandleFunc("/signup", server.signup)
	mux.HandleFunc("/login", server.login)
	mux.HandleFunc("/logout", server.logout)
	mux.HandleFunc("/refresh", server.refresh)
	mux.HandleFunc("/profile", server.authMiddleware(server.profile))
	mux.HandleFunc("/admin", server.authMiddleware(server.admin))

	loggedMux := server.loggingMiddleware(server.corsMiddleware(mux))

	log.Println("Server running on : 8080")
	log.Fatal(http.ListenAndServe(":8080", loggedMux))
}
