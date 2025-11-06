package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/tilshansanoj/golangrestapi/dbconfig"
	"github.com/tilshansanoj/golangrestapi/internal/handlers"
	"github.com/tilshansanoj/golangrestapi/internal/routes"
	"github.com/tilshansanoj/golangrestapi/internal/store"
	"github.com/tilshansanoj/golangrestapi/serverconfig"
)

func main()  {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	// slog.Info("App is running...")

	//load server config
	config, err := serverconfig.GetConfig()
	if  err != nil {
		log.Fatalf("Failed to load server config: %v\n", err)
	}

	// Connect to the database
	db, err := dbconfig.ConnectDB(config.DatabaseURL)
	defer db.Close()

	// Connect to Redis
	rdb := dbconfig.ConnectRedis()
	defer func (rdb *redis.Client) {
		_ = rdb.Close()
	}(rdb)

	queries := store.New(db)
	handler := handlers.NewHandler(db, queries, rdb)
	//set up http server
	mux := http.NewServeMux()

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