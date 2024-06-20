package auth

import (
	"go-api-with-fiber/database"
	"go-api-with-fiber/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Register(c *fiber.Ctx) error {
	registerRequest := new(RegisterRequest)

	if err := c.BodyParser(registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"data":    nil,
		})
	}

	if err := validate.Struct(registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	if registerRequest.Password != registerRequest.PasswordConfirmation {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "password is not the same as password confirmation",
			"data":    nil,
		})
	}

	hashPassword, errHash := hashPassword(registerRequest.Password)

	if errHash != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "couldn't hash password",
			"data":    nil,
		})
	}

	newUser := model.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: hashPassword,
	}

	errCreate := database.DB.Create(&newUser).Error

	if errCreate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to register",
			"data":    nil,
		})
	}

	registerResponse := RegisterResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    registerResponse,
	})
}
