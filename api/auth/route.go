package auth

import (
	"go-api-with-fiber/middleware"

	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	app.Post("/auth/register", middleware.AnonymousProtected, Register)
	app.Post("/auth/login", middleware.AnonymousProtected, LogIn)
	app.Post("/auth/refresh-token", middleware.AnonymousProtected, RefreshToken)
	app.Post("/auth/anonymous-token", AnonymousToken)
}
