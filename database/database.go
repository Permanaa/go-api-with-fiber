package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client

func DBConnect() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}

	DB = db

	fmt.Println("connected to database")
}

func RedisConnect() {
	redisDBNameNumber, errConvertRedisDBName := strconv.Atoi(os.Getenv("REDIS_DB"))

	if errConvertRedisDBName != nil {
		panic(errConvertRedisDBName)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDBNameNumber,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("failed to connect redis:", err)
		return
	}

	RedisClient = client

	fmt.Println("connected to redis")
}
