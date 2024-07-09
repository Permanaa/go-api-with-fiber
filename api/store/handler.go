package store

import (
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"
	"go-api-with-fiber/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func CreateStore(c *fiber.Ctx) error {
	createStoreRequest := new(CreateStoreRequest)

	if err := c.BodyParser(createStoreRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
		})
	}

	if err := validate.Struct(createStoreRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	bearerToken := c.Get("Authorization")

	tokenClaims, errParseToken := utils.ParseToken(bearerToken)

	if errParseToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errParseToken.Error(),
			"data":    nil,
		})
	}

	userIdString, errGetSubject := tokenClaims.GetSubject()

	if errGetSubject != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGetSubject.Error(),
			"data":    nil,
		})
	}

	userId, errConvertUserId := strconv.Atoi(userIdString)

	if errConvertUserId != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errConvertUserId.Error(),
			"data":    nil,
		})
	}

	var stores []model.Store

	errFindStoreByUserId := database.DB.Where("user_id = ?", userId).Find(&stores).Error

	if errFindStoreByUserId != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errFindStoreByUserId.Error(),
			"data":    nil,
		})
	}

	if len(stores) >= 1 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "one user one store",
			"data":    nil,
		})
	}

	newStore := model.Store{
		Name:   createStoreRequest.Name,
		UserID: uint(userId),
	}

	errCreateStore := database.DB.Create(&newStore).Error

	if errCreateStore != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errCreateStore.Error(),
			"data":    nil,
		})
	}

	createStoreResponse := StoreResponse{
		ID:        newStore.ID,
		Name:      newStore.Name,
		UserID:    newStore.UserID,
		CreatedAt: newStore.CreatedAt,
		UpdatedAt: newStore.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    createStoreResponse,
	})
}
