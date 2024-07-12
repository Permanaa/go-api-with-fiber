package store

import (
	"errors"
	"fmt"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"
	"go-api-with-fiber/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
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

	tokenClaims, errParseToken := utils.ParseAccessToken(bearerToken)

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
			"message": "user can only create one store",
			"data":    nil,
		})
	}

	newStore := model.Store{
		Name:   createStoreRequest.Name,
		Slug:   slug.Make(createStoreRequest.Name),
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
		Slug:      newStore.Slug,
		UserID:    newStore.UserID,
		CreatedAt: newStore.CreatedAt,
		UpdatedAt: newStore.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    createStoreResponse,
	})
}

func DeleteStoreBySlug(c *fiber.Ctx) error {
	var store model.Store

	errGetStore := database.DB.Where("slug = ?", c.Params("slug")).First(&store).Error

	if errors.Is(errGetStore, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": errGetStore.Error(),
			"data":    nil,
		})
	}

	if errGetStore != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGetStore.Error(),
			"data":    nil,
		})
	}

	errDeleteStore := database.DB.Delete(&store).Error

	if errDeleteStore != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errDeleteStore,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    store,
	})
}

func GetStoreBySlug(c *fiber.Ctx) error {
	var store model.Store

	errGetStore := database.DB.Where("slug = ?", c.Params("slug")).First(&store).Error

	if errors.Is(errGetStore, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": errGetStore.Error(),
			"data":    nil,
		})
	}

	if errGetStore != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGetStore.Error(),
			"data":    nil,
		})
	}

	storeResponse := StoreResponse{
		ID:        store.ID,
		Name:      store.Name,
		Slug:      store.Slug,
		UserID:    store.UserID,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    storeResponse,
	})
}

func GetAllStore(c *fiber.Ctx) error {
	pageQuery := c.Query("page", "1")
	limitQuery := c.Query("limit", "10")
	orderByQuery := c.Query("orderBy", "created_at")
	sortQuery := c.Query("sort", "desc")

	page, errConvertPage := strconv.Atoi(pageQuery)

	if errConvertPage != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errConvertPage.Error(),
			"data":    nil,
		})
	}

	limit, errConvertLimit := strconv.Atoi(limitQuery)

	if errConvertLimit != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errConvertLimit.Error(),
			"data":    nil,
		})
	}

	errValidateQuery := validate.Struct(GetAllStoreQuery{
		Page:    page,
		Limit:   limit,
		OrderBy: orderByQuery,
		Sort:    sortQuery,
	})

	if errValidateQuery != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": errValidateQuery.Error(),
			"data":    nil,
		})
	}

	order := fmt.Sprintf("%s %s", orderByQuery, sortQuery)

	var stores []model.Store

	errGetAllStore := database.DB.Limit(limit).Offset((page - 1) * limit).Order(order).Find(&stores).Error

	if errGetAllStore != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errGetAllStore.Error(),
			"data":    nil,
		})
	}

	storesResponse := []StoreResponse{}

	for _, store := range stores {
		storesResponse = append(storesResponse, StoreResponse{
			ID:        store.ID,
			Name:      store.Name,
			Slug:      store.Slug,
			UserID:    store.UserID,
			CreatedAt: store.CreatedAt,
			UpdatedAt: store.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    storesResponse,
	})
}
