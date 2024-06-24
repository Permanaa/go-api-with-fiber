package user

import (
	"errors"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

func GetAll(c *fiber.Ctx) error {
	var users []model.User

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
			Email:     user.Email,
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
	var user model.User

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
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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

	var user model.User

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

	errUpdate := database.DB.Model(&user).Updates(model.User{
		Name:  userRequest.Name,
		Email: userRequest.Email,
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
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    userResponse,
	})
}

func Delete(c *fiber.Ctx) error {
	var user model.User

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

	userResponse := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    userResponse,
	})
}
