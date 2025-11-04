package middleware

import "github.com/gofiber/fiber/v2"

// MiddlewareProviders holds middleware handlers
// This struct is used to avoid dig's "same type" conflict when injecting multiple fiber.Handler instances
type MiddlewareProviders struct {
	Logger  fiber.Handler
	Recover fiber.Handler
}

// NewMiddlewareProviders creates middleware providers
// This is a constructor function for dependency injection
func NewMiddlewareProviders() *MiddlewareProviders {
	return &MiddlewareProviders{
		Logger:  LoggerConfig(),
		Recover: RecoverConfig(),
	}
}

