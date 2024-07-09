package product

import (
	"errors"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

func CreateProduct(c *fiber.Ctx) error {
	createProductRequest := new(ProductRequest)

	if err := c.BodyParser(createProductRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	if err := validate.Struct(createProductRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	newProduct := model.Product{
		Name:  createProductRequest.Name,
		Price: createProductRequest.Price,
	}

	errCreateProduct := database.DB.Create(&newProduct).Error

	if errCreateProduct != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": errCreateProduct.Error(),
			"data":    nil,
		})
	}

	createProductResponse := ProductResponse{
		ID:        newProduct.ID,
		Name:      newProduct.Name,
		Price:     newProduct.Price,
		CreatedAt: newProduct.CreatedAt,
		UpdatedAt: newProduct.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    createProductResponse,
	})
}

func GetAllProduct(c *fiber.Ctx) error {
	var products []model.Product

	errGetAll := database.DB.Find(&products).Error

	if errGetAll != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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
		"message": "success",
		"data":    productsResponse,
	})
}

func GetByIdProduct(c *fiber.Ctx) error {
	var product model.Product

	err := database.DB.First(&product, c.Params("id")).Error

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

	productResponse := ProductResponse{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    productResponse,
	})
}
