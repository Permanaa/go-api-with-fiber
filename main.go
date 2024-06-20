package main

import (
	"go-api-with-fiber/api"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/migration"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	database.DBConnect()

	migration.Migrate()

	api.Routes(app)

	app.Listen("127.0.0.1:1804")
}
