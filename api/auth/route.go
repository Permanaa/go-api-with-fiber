package auth

import "github.com/gofiber/fiber/v2"

func Route(app *fiber.App) {
	app.Post("/auth/register", Register)
	app.Post("/auth/login", LogIn)
	app.Post("/auth/refresh-token", RefreshToken)
}
