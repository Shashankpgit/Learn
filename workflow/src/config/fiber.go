package config

import (
	"app/src/constants"
	"app/src/utils"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

// FiberConfig returns the Fiber app configuration
// Takes Config as parameter for dependency injection
// Note: Using standard library JSON encoder/decoder for Go 1.24 compatibility
func FiberConfig(cfg *Config) fiber.Config {
	return fiber.Config{
		Prefork:       cfg.IsProd,
		CaseSensitive: true,
		ServerHeader:  constants.ServerHeaderName,
		AppName:       constants.AppName,
		ErrorHandler:  utils.ErrorHandler,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
	}
}
