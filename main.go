package main

import (
	"go-api-with-fiber/api"
	"go-api-with-fiber/config"
	"go-api-with-fiber/migration"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.DBConnect()

	migration.Migrate()

	api.Routes(app)

	app.Listen("127.0.0.1:1804")
}
