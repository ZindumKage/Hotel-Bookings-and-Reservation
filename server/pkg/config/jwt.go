package config

import (
	"os"
	"time"
	
	"github.com/golang-jwt/jwt/v5"
)



func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GenerateAccessToken(userID uint) (string, error) {

		claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), 
		"type":    "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{ 
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), 
		"type":    "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}	