package main

import (
	"go-api-with-fiber/api"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/migration"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	database.DBConnect()
	database.RedisConnect()

	migration.Migrate()

	app.Use(limiter.New())

	api.Routes(app)

	app.Listen("127.0.0.1:1804")
}
