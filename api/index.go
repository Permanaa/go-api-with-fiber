package api

import (
	"go-api-with-fiber/api/user"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	user.Route(app)
}
