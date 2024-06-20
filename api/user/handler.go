package user

import (
	"errors"
	"go-api-with-fiber/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

func GetAll(c *fiber.Ctx) error {
	var users []User

	result := database.DB.Find(&users)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed get users",
			"data":    nil,
		})
	}

	var usersResponse []UserResponse

	for _, user := range users {
		usersResponse = append(usersResponse, UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    usersResponse,
	})
}

func GetById(c *fiber.Ctx) error {
	var user User

	err := database.DB.First(&user, c.Params("id")).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "record not found",
			"data":    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get user",
			"data":    nil,
		})
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    userResponse,
	})
}

func Create(c *fiber.Ctx) error {
	userBody := new(UserRequest)

	if err := c.BodyParser(userBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"data":    nil,
		})
	}

	if err := validate.Struct(userBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	newUser := User{
		Name: userBody.Name,
	}

	err := database.DB.Create(&newUser).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create user",
			"data":    nil,
		})
	}

	userResponse := UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    userResponse,
	})
}

func Update(c *fiber.Ctx) error {
	userRequest := new(UserRequest)

	if err := c.BodyParser(userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"data":    nil,
		})
	}

	if err := validate.Struct(userRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	var user User

	errGetById := database.DB.First(&user, c.Params("id")).Error

	if errors.Is(errGetById, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "record not found",
			"data":    nil,
		})
	}

	if errGetById != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get user",
			"data":    nil,
		})
	}

	errUpdate := database.DB.Model(&user).Updates(User{
		Name: userRequest.Name,
	}).Error

	if errUpdate != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to update user",
			"data":    nil,
		})
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    userResponse,
	})
}

func Delete(c *fiber.Ctx) error {
	var user User

	err := database.DB.First(&user, c.Params("id")).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "record not found",
			"data":    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get user",
			"data":    nil,
		})
	}

	errDelete := database.DB.Delete(&user).Error

	if errDelete != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete user",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    user,
	})
}
