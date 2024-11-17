package utils

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string, key string, expiredAt time.Time) (string, error) {
	generate := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "scratching",
		Subject:   id,
		ExpiresAt: jwt.NewNumericDate(expiredAt),
	})

	return generate.SignedString([]byte(key))
}

func ParseAccessToken(bearerToken string) (jwt.MapClaims, error) {
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
