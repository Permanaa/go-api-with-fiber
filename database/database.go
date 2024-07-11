package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client

func DBConnect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		panic(err.Error())
	}

	DB = db

	fmt.Println("connected to database")
}

func RedisConnect() {
	redisDBNameNumber, errConvertRedisDBName := strconv.Atoi(os.Getenv("REDIS_DB"))

	if errConvertRedisDBName != nil {
		fmt.Println("failed to connect redis:", errConvertRedisDBName.Error())
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDBNameNumber,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		fmt.Println("failed to connect redis:", err.Error())
		return
	}

	RedisClient = client

	fmt.Println("connected to redis")
}
