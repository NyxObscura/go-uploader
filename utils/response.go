package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse adalah struktur standar untuk semua respons JSON.
type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON mengirimkan respons dalam format JSON.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

