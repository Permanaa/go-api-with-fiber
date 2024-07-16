package main

import (
	"fmt"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	errLoadEnv := godotenv.Load()
	if errLoadEnv != nil {
		log.Fatalf("Error loading .env file: %s", errLoadEnv.Error())
	}

	database.DBConnect()

	var product model.Product

	product.Name = "Product Seeding Testing"
	product.Price = 10000

	errSave := database.DB.Save(&product).Error

	if errSave != nil {
		log.Fatalf("failed to seed: %s", errSave.Error())
	}

	fmt.Println("seeding product testing was successful")
}
