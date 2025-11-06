package utils

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	
)

type Claims struct {
	UserID   int64 `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(userID int64, username string, secretKey []byte) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "your_app_name",
		},
	}	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}