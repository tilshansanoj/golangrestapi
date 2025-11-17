package middlewares

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"github.com/tilshansanoj/golangrestapi/dbconfig"
	"github.com/tilshansanoj/golangrestapi/internal/auth"
	"github.com/tilshansanoj/golangrestapi/internal/errorhandler"
)

type contextKey string

const UserClaimKey contextKey = "user_claims"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Retrieves the authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorhandler.ResponseWithError(w, http.StatusUnauthorized, "No token provided")
			return
		}

		// Bearer token parsing
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &auth.Claims{}

		// Check redis for blacklisted token
		isBlacklisted, err := dbconfig.RedisClient.Get(context.Background(), tokenString).Result()
		if err == nil && isBlacklisted == "blacklisted" {
			errorhandler.ResponseWithError(w, http.StatusUnauthorized, "Token has been revoked")
			slog.Warn("Blacklisted token access attempt", "token", tokenString)
			return
		} else if err != nil && err != redis.Nil {
			errorhandler.ResponseWithError(w, http.StatusInternalServerError, "Internal server error")
			slog.Error("Redis error while checking blacklisted token", "error", err)
			return
		}

		// parse the token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		// handle validation errors
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				errorhandler.ResponseWithError(w, http.StatusUnauthorized, "Invalid token signature")
				return
			}
			errorhandler.ResponseWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		if token.Valid {
			ctx := context.WithValue(r.Context(), UserClaimKey, claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}else {
			errorhandler.ResponseWithError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
	}
}