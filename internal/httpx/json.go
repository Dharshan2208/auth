package httpx

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func DecodeJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
