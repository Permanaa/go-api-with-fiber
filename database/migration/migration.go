package migration

import (
	"fmt"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"
)

func Migrate() {
	err := database.DB.AutoMigrate(
		&model.User{},
		&model.Product{},
	)

	if err != nil {
		fmt.Println("failed to run migration:", err)
		return
	}

	fmt.Println("database migrated")
}
