package dbconfig

import (
	"context"
	"os"
	"strconv"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RedisClient *redis.Client


func ConnectRedis() *redis.Client {

	host := getEnv("REDIS_HOST", "localhost:6379")
	password := getEnv("REDIS_PASSWORD", "")
	db := getEnvAsInt("REDIS_DB", 0)

	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	RedisClient = rdb

	fmt.Println("Connected to Redis!")

	return rdb
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue	
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if valueStr == "" {
        return defaultValue
    }
    
    value, err := strconv.Atoi(valueStr)
    if err != nil {
        return defaultValue
    }
    return value
}	