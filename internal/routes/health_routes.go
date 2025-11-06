package routes

import (
	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"net/http"
)

func SetupHealthRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("/health", handler.HealthCheckHandler())
}