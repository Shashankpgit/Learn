package utils

import (
	"app/src/constants"
	"app/src/validation"
	"errors"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is a centralized error handler for the application
func ErrorHandler(c *fiber.Ctx, err error) error {
	responseBuilder := NewResponseBuilder()

	// Handle validation errors
	if errorsMap := validation.CustomErrorMessages(err); len(errorsMap) > 0 {
		return responseBuilder.ValidationError(c, errorsMap)
	}

	// Handle fiber errors
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return responseBuilder.Error(c, fiberErr.Code, fiberErr.Message)
	}

	// Handle unknown errors
	return responseBuilder.InternalServerError(c, constants.MsgInternalServerError)
}

// NotFoundHandler handles 404 errors
func NotFoundHandler(c *fiber.Ctx) error {
	return Response.NotFound(c, constants.MsgEndpointNotFound)
}
