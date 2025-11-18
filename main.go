package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/tilshansanoj/golangrestapi/dbconfig"
	_ "github.com/tilshansanoj/golangrestapi/docs"
	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"github.com/tilshansanoj/golangrestapi/internal/routes"
	"github.com/tilshansanoj/golangrestapi/internal/store"
	"github.com/tilshansanoj/golangrestapi/serverconfig"
)

// @title Golang REST API
// @version 1.0
// @description This is a sample API created by Golang.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	// slog.Info("App is running...")

	//load server config
	config, err := serverconfig.GetConfig()
	if err != nil {
		log.Fatalf("Failed to load server config: %v\n", err)
	}

	// Connect to the database
	db, err := dbconfig.ConnectDB(config.DatabaseURL)
	defer db.Close()

	// Connect to Redis
	rdb := dbconfig.ConnectRedis()
	defer func(rdb *redis.Client) {
		_ = rdb.Close()
	}(rdb)

	queries := store.New(db)
	handler := handlers.NewHandler(db, queries, rdb)
	//set up http server
	mux := http.NewServeMux()
	//setup swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	//setup routes
	routes.SetupRoutes(mux, handler)

	//server instance
	serverAddress := fmt.Sprintf(":%s", config.ServerPort)
	server := &http.Server{
		Addr:    serverAddress,
		Handler: mux,
	}

	fmt.Printf("Server is running on port %s\n", config.ServerPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}

}
