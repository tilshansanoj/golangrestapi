package routes

import (
	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"net/http"
)

func SetupBlogRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("POST /api/blog/register", handler.CreateBlogHandler())
	mux.HandleFunc("GET /api/blog", handler.ListBlogsHandler())
	mux.HandleFunc("GET /api/blog/{id}", handler.GetBlogHandler())
	mux.Handle("/api/blog/", http.StripPrefix("/api/blog", mux))
}
