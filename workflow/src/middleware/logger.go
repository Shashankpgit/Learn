package middleware

import (
	"app/src/constants"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func LoggerConfig() fiber.Handler {
	return logger.New(logger.Config{
		Format:     constants.LoggerFormat,
		TimeFormat: constants.LoggerTimeFormat,
	})
}
