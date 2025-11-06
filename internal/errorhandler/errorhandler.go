package errorhandler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func ResponseWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

func ResponseWithNotFound(w http.ResponseWriter) {
	ResponseWithError(w, http.StatusNotFound, "Resource not found")
}