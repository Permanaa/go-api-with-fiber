package migration

import (
	"fmt"
	"go-api-with-fiber/database"
	"go-api-with-fiber/database/model"
)

func Migrate() {
	err := database.DB.AutoMigrate(&model.User{})

	if err != nil {
		panic("failed to run migration")
	}

	fmt.Println("database migrated")
}
