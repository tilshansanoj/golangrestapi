package routes

import (
	"net/http"

	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"github.com/tilshansanoj/golangrestapi/internal/middlewares"
)

func SetupUserRoutes(mux *http.ServeMux, handler *handlers.Handler) {
	mux.HandleFunc("POST /api/users/register", handler.CreateUserHandler())
	mux.HandleFunc("POST /api/users/login", handler.LoginUserHandler())
	mux.HandleFunc("GET /api/users", handler.GetAllUsersHandler())
	mux.HandleFunc("GET /api/users/me", middlewares.AuthMiddleware(http.HandlerFunc(handler.GetUserByIdHandler())))
	mux.HandleFunc("PATCH /api/users/{id}", handler.UpdateUserHandler())
	mux.HandleFunc("DELETE /api/users/{id}", handler.DeleteUserHandler())
	mux.Handle("/api/users/", http.StripPrefix("/api/users", mux))
}