package routes

import (
	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"net/http"
)

func SetupBlogRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("POST /api/blogs/register", handler.CreateBlogHandler())
	mux.HandleFunc("GET /api/blogs", handler.ListBlogsHandler())
	mux.HandleFunc("GET /api/blogs/{id}", handler.GetBlogHandler())
	mux.HandleFunc("PATCH /api/blogs/{id}", handler.UpdateBlogHandler())
	mux.HandleFunc("DELETE /api/blogs/{id}", handler.DeleteBlogHandler())
	mux.Handle("/api/blogs/", http.StripPrefix("/api/blogs", mux))
}
