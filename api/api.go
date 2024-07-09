package api

import (
	"go-api-with-fiber/api/auth"
	"go-api-with-fiber/api/product"
	"go-api-with-fiber/api/redis"
	"go-api-with-fiber/api/user"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	user.Route(app)
	auth.Route(app)
	redis.Route(app)
	product.Route(app)
}
