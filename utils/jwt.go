package utils

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(id string) (string, error) {
	generate := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "scratching",
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
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

func ParseToken(bearerToken string) (jwt.MapClaims, error) {
	splitBearer := strings.Split(bearerToken, " ")
	tokenString := splitBearer[1]

	parseToken, errParseToken := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET_KEY")), nil
	})

	if errParseToken != nil {
		return nil, errParseToken
	}

	tokenClaims := parseToken.Claims.(jwt.MapClaims)

	return tokenClaims, nil
}
