package utils

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

// Success mengirim response sukses dengan format konsisten:
// { "success": true, "data": ... }
func Success(w http.ResponseWriter, statusCode int, data interface{}) {
	writeJSON(w, statusCode, envelope{
		"success": true,
		"data":    data,
	})
}

// Error mengirim response error dengan format konsisten:
// { "success": false, "error": "..." }
func Error(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, envelope{
		"success": false,
		"error":   message,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

// DecodeJSON membaca body request JSON ke struct `dst`.
func DecodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}
