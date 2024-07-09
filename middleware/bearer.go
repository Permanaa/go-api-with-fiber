package middleware

import (
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"
	"go-api-with-fiber/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func BearerProtected(c *fiber.Ctx) error {
	bearerToken := c.Get("Authorization")

	if bearerToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"data":    nil,
		})
	}

	tokenClaims, errParseToken := utils.ParseToken(bearerToken)

	if errParseToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errParseToken.Error(),
			"data":    nil,
		})
	}

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

	userId, errGetSubject := tokenClaims.GetSubject()

	if errGetSubject != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGetSubject.Error(),
			"data":    nil,
		})
	}

	var user model.User

	errGetUser := database.DB.First(&user, userId).Error

	if errGetUser != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"data":    nil,
		})
	}

	return c.Next()
}
