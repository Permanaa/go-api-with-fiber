package redis

import (
	"go-api-with-fiber/database"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func SetRedis(c *fiber.Ctx) error {
	setRedisRequest := new(SetRedisRequest)

	if err := c.BodyParser(setRedisRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"data":    nil,
		})
	}

	if err := validate.Struct(setRedisRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	setRedis := database.RedisClient.Set(c.Context(), setRedisRequest.Key, setRedisRequest.Value, time.Minute*5)

	if setRedis.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": setRedis.Err().Error(),
			"data":    nil,
		})
	}

	value, err := setRedis.Result()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    value,
	})
}

func GetRedis(c *fiber.Ctx) error {
	getRedis := database.RedisClient.Get(c.Context(), c.Params("key"))

	if getRedis.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": getRedis.Err(),
			"data":    nil,
		})
	}

	value, err := getRedis.Result()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    value,
	})
}
