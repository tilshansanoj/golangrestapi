package routes

import (
	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"net/http"
)

func SetupBlogRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("POST /blogs/register", handler.CreateBlogHandler())
	mux.HandleFunc("GET /blogs", handler.ListBlogsHandler())
	mux.HandleFunc("GET /blogs/{id}", handler.GetBlogHandler())
}
