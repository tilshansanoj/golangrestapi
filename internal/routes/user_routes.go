package routes

import (
	"net/http"

	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"github.com/tilshansanoj/golangrestapi/internal/middlewares"
)

func SetupUserRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("POST /api/user/register", handler.CreateUserHandler())
	mux.HandleFunc("POST /api/user/login", handler.LoginUserHandler())
	mux.HandleFunc("GET /api/user", handler.GetAllUsersHandler())
	mux.HandleFunc("GET /api/user/me", middlewares.AuthMiddleware(http.HandlerFunc(handler.GetUserByIdHandler())))
	mux.HandleFunc("PATCH /api/user/{id}", handler.UpdateUserHandler())
	mux.HandleFunc("DELETE /api/user/{id}", handler.DeleteUserHandler())
	mux.Handle("/api/user/", http.StripPrefix("/api/user", mux))
}