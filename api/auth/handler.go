package auth

import (
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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
			"message": err.Error(),
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
			"message": errCreate.Error(),
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

func LogIn(c *fiber.Ctx) error {
	logInRequest := new(LogInRequest)

	if err := c.BodyParser(logInRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	if err := validate.Struct(logInRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	var user model.User

	errFindByEmail := database.DB.Where("email = ?", logInRequest.Email).Find(&user).Error

	if errFindByEmail != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "wrong user email or password",
			"data":    nil,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(logInRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "wrong user email or password",
			"data":    nil,
		})
	}

	generateAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "scratching",
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	accessToken, errGenerateAccessToken := generateAccessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	if errGenerateAccessToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to log in",
			"data":    nil,
		})
	}

	logInResponse := LogInResponse{
		AccessToken: accessToken,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    logInResponse,
	})
}
