package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type Server struct {
	users map[string]User
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

	user := User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     "user",
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
	writeJSON(w, http.StatusOK, map[string]string{"message": "login successful"})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "running",
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	server := &Server{
		users: make(map[string]User),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", server.health)
	mux.HandleFunc("/signup", server.signup)
	mux.HandleFunc("/login", server.login)

	log.Println("Server running on : 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
