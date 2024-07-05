package redis

import (
	"go-api-with-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	app.Get("/redis/set", middleware.BearerProtected, SetRedis)
	app.Get("/redis/get/:key", middleware.BearerProtected, GetRedis)
}
