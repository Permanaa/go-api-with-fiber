package product

import (
	"errors"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"
	"go-api-with-fiber/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

func CreateProduct(c *fiber.Ctx) error {
	createProductRequest := new(ProductRequest)

	if err := c.BodyParser(createProductRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
	}

	if err := validate.Struct(createProductRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
			"data":    nil,
		})
	}

	bearerToken := c.Get("Authorization")

	tokenClaims, errParseToken := utils.ParseAccessToken(bearerToken)

	if errParseToken != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errParseToken.Error(),
			"data":    nil,
		})
	}

	userIdString, errGetSubject := tokenClaims.GetSubject()

	if errGetSubject != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errGetSubject.Error(),
			"data":    nil,
		})
	}

	userId, errConvertUserId := strconv.Atoi(userIdString)

	if errConvertUserId != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errConvertUserId.Error(),
			"data":    nil,
		})
	}

	var store model.Store

	errFindStoreByUserId := database.DB.Where("user_id = ?", userId).Preload("Products").First(&store).Error

	if errFindStoreByUserId != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errFindStoreByUserId.Error(),
			"data":    nil,
		})
	}

	for _, product := range store.Products {
		if product.Name == createProductRequest.Name {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"code":    fiber.StatusConflict,
				"message": "duplicate product name",
				"data":    nil,
			})
		}
	}

	newProduct := model.Product{
		Name:    createProductRequest.Name,
		Price:   createProductRequest.Price,
		StoreID: store.ID,
	}

	errCreateProduct := database.DB.Create(&newProduct).Error

	if errCreateProduct != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": errCreateProduct.Error(),
			"data":    nil,
		})
	}

	createProductResponse := ProductResponse{
		ID:        newProduct.ID,
		Name:      newProduct.Name,
		Price:     newProduct.Price,
		StoreID:   newProduct.StoreID,
		CreatedAt: newProduct.CreatedAt,
		UpdatedAt: newProduct.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "success",
		"data":    createProductResponse,
	})
}

func GetAllProduct(c *fiber.Ctx) error {
	var products []model.Product

	errGetAll := database.DB.Find(&products).Error

	if errGetAll != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "failed get users",
			"data":    nil,
		})
	}

	var productsResponse []ProductResponse

	for _, product := range products {
		productsResponse = append(productsResponse, ProductResponse{
			ID:        product.ID,
			Name:      product.Name,
			Price:     product.Price,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "success",
		"data":    productsResponse,
	})
}

func GetByIdProduct(c *fiber.Ctx) error {
	var product model.Product

	err := database.DB.First(&product, c.Params("id")).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"message": "record not found",
			"data":    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "failed to get user",
			"data":    nil,
		})
	}

	productResponse := ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code":    fiber.StatusOK,
		"message": "success",
		"data":    productResponse,
	})
}
