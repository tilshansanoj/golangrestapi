package utils

import (
	"log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:")
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(storedPassword, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	return err == nil
}