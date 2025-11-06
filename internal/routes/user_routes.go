package routes

import (
	"net/http"

	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"github.com/tilshansanoj/golangrestapi/internal/middlewares"
)

func SetupUserRoutes(mux *http.ServeMux, handler *handlers.Handler) {

	mux.HandleFunc("POST /user/register", handler.CreateUserHandler())
	mux.HandleFunc("POST /user/login", handler.LoginUserHandler())
	// mux.HandleFunc("GET /user/profile", handler.GetUserProfileHandler())
	mux.HandleFunc("GET /user/get-all", handler.GetAllUsersHandler())
	mux.HandleFunc("GET /user/profile", middlewares.AuthMiddleware(http.HandlerFunc(handler.GetUserHandler())))
	mux.HandleFunc("PATCH /user/profile", middlewares.AuthMiddleware(http.HandlerFunc(handler.UpdateUser())))
	mux.Handle("/user/", http.StripPrefix("/user", mux))
}