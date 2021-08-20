package helpers

import (
	"encoding/json"
	"net/http"
)

func RespondStatusCode(w http.ResponseWriter, status int, message string) {
	RespondJson(w, status, map[string]interface{}{"status_code": status, "message": message})
}

func RespondJson(w http.ResponseWriter, status_code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	json.NewEncoder(w).Encode(data)
}
