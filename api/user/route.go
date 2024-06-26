package user

import (
	"go-api-with-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	app.Get("/user", middleware.BearerProtected, GetAll)
	app.Get("/user/:id", middleware.BearerProtected, GetById)
	app.Put("/user/:id", middleware.BearerProtected, Update)
	app.Delete("/user/:id", middleware.BearerProtected, Delete)
}
