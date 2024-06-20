package migration

import (
	"go-api-with-fiber/database"
	"go-api-with-fiber/model"
	"log"
)

func Migrate() {
	err := database.DB.AutoMigrate(&model.User{})

	if err != nil {
		panic("failed to run migration")
	}

	log.Println("database migrated")
}
