package routes

import (
	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	// Setup health routes
	SetupHealthRoutes(mux, handler)
	SetupUserRoutes(mux, handler)
	RegisterBlogRoutes(mux, handler)
	GetBlogRoutes(mux, handler)
}