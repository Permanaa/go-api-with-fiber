package product

import (
	"go-api-with-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	app.Post("/product", middleware.BearerProtected, CreateProduct)
	app.Get("/product", middleware.BearerProtected, GetAllProduct)
}
