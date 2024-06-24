package user

import (
	"go-api-with-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	app.Get("/user", middleware.Protected, GetAll)
	app.Get("/user/:id", middleware.Protected, GetById)
	app.Put("/user/:id", middleware.Protected, Update)
	app.Delete("/user/:id", middleware.Protected, Delete)
}
