package auth

import (
	"encoding/base64"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"
	"go-api-with-fiber/utils"
	"os"
	"strconv"
	"strings"
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

	accessToken, errGenerateAccessToken := utils.GenerateAccessToken(strconv.Itoa(int(user.ID)))

	if errGenerateAccessToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGenerateAccessToken.Error(),
			"data":    nil,
		})
	}

	refreshToken, errGenerateRefreshtoken := utils.GenerateRefreshToken(strconv.Itoa(int(user.ID)))

	if errGenerateRefreshtoken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGenerateRefreshtoken.Error(),
			"data":    nil,
		})
	}

	logInResponse := TokenResponse{
		AccessToken:  accessToken,
		Refreshtoken: refreshToken,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    logInResponse,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	refreshTokenRequest := new(RefreshTokenRequest)

	if err := c.BodyParser(refreshTokenRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	if err := validate.Struct(refreshTokenRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	parseToken, errParseToken := jwt.Parse(refreshTokenRequest.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET_KEY")), nil
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

	accessToken, errGenerateAccessToken := utils.GenerateAccessToken(strconv.Itoa(int(user.ID)))

	if errGenerateAccessToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGenerateAccessToken.Error(),
			"data":    nil,
		})
	}

	refreshToken, errGenerateRefreshtoken := utils.GenerateRefreshToken(strconv.Itoa(int(user.ID)))

	if errGenerateRefreshtoken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGenerateRefreshtoken.Error(),
			"data":    nil,
		})
	}

	refreshTokenResponse := TokenResponse{
		AccessToken:  accessToken,
		Refreshtoken: refreshToken,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    refreshTokenResponse,
	})
}

func AnonymousToken(c *fiber.Ctx) error {
	basicToken := c.Get("Authorization")

	if basicToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"data":    nil,
		})
	}

	splitBearer := strings.Split(basicToken, " ")
	tokenString := splitBearer[1]

	basicByte, errDecode := base64.StdEncoding.DecodeString(tokenString)

	if errDecode != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errDecode.Error(),
			"data":    nil,
		})
	}

	splitUserPass := strings.Split(string(basicByte), ":")

	username := splitUserPass[0]
	password := splitUserPass[1]

	if username != os.Getenv("BASIC_AUTH_USERNAME") || password != os.Getenv("BASIC_AUTH_PASSWORD") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "unauthorized",
			"data":    nil,
		})
	}

	generateAnonToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "scratching",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
	})

	anonToken, errGenerateAnonToken := generateAnonToken.SignedString([]byte(os.Getenv("JWT_ANONYMOUS_TOKEN_SECRET_KEY")))

	if errGenerateAnonToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGenerateAnonToken.Error(),
			"data":    nil,
		})
	}

	AnonymousResponse := AnonymousResponse{
		AnonymousToken: anonToken,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    AnonymousResponse,
	})
}
