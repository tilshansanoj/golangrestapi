package routes

import (
	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"net/http"
)

func RegisterBlogRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("POST /blogs/register", handler.CreateBlogHandler())
}

func GetBlogRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("GET /blogs", handler.GetBlogsHandler())
}
