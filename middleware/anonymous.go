package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AnonymousProtected(c *fiber.Ctx) error {
	bearerToken := c.Get("Authorization")

	if bearerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"data":    nil,
		})
	}

	splitBearer := strings.Split(bearerToken, " ")
	tokenString := splitBearer[1]

	parseToken, errParseToken := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ANONYMOUS_TOKEN_SECRET_KEY")), nil
	})

	if errParseToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errParseToken.Error(),
			"data":    nil,
		})
	}

	tokenClaims := parseToken.Claims.(jwt.MapClaims)

	expirationTime, errGetExp := tokenClaims.GetExpirationTime()

	if errGetExp != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGetExp.Error(),
			"data":    nil,
		})
	}

	if expirationTime.Unix() < time.Now().Unix() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"data":    nil,
		})
	}

	return c.Next()
}
