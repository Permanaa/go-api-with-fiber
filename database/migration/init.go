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

	errMigrate := database.DB.AutoMigrate(
		&model.Product{},
		&model.Store{},
		&model.User{},
	)

	if errMigrate != nil {
		log.Fatalf("failed to run migration: %s", errMigrate.Error())
	}

	fmt.Println("database migrated")
}
