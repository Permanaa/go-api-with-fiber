package migration

import (
	"go-api-with-fiber/api/user"
	"go-api-with-fiber/config"
	"log"
)

func Migrate() {
	err := config.DB.AutoMigrate(&user.User{})

	if err != nil {
		panic("failed to run migration")
	}

	log.Println("database migrated")
}
