package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(id string) (string, error) {
	generate := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "scratching",
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 20)),
	})

	return generate.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET_KEY")))
}

func GenerateRefreshToken(id string) (string, error) {
	generate := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "scratching",
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
	})

	return generate.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET_KEY")))
}
