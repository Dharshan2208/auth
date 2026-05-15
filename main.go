package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct{}

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
	server := &Server{}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", server.health)

	log.Println("Server running on : 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
