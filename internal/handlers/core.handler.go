package handlers

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
	"github.com/tilshansanoj/golangrestapi/internal/store"
)

type Handler struct{
	// DB instance
	DB *sql.DB
	// Query stores
	Queries *store.Queries
	// Redis instance
	Redis *redis.Client
}

func NewHandler(db *sql.DB, queries *store.Queries, redis *redis.Client) *Handler {
	return &Handler{
		DB:      db,
		Queries: queries,
		Redis:   redis,
	}
}