package auth

import "github.com/gofiber/fiber/v2"

func Route(app *fiber.App) {
	app.Post("/auth/register", Register)
}
