package store

import (
	"go-api-with-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	app.Post("/store", middleware.BearerProtected, CreateStore)
	app.Delete("/store/:slug", middleware.BearerProtected, DeleteStoreBySlug)
}
