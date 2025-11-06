package serverconfig

import (
	"github.com/joho/godotenv"
	"os"
	"fmt"
	"strconv"
)

type serverConfig struct {
	ServerPort string
	DatabaseURL string
	Environment string
	LogLevel string
	RedisHost string
	RedisPassword string
	RedisDatabase int
}


func GetConfig() (*serverConfig, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return &serverConfig{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		RedisHost:       getEnv("REDIS_HOST", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDatabase:   getEnvAsInt("REDIS_DB", 0),
	}, nil
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