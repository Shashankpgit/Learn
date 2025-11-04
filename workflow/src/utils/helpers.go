package utils

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// IsNotFoundError checks if an error is a "not found" error
// Handles both gorm.ErrRecordNotFound and fiber 404 errors
func IsNotFoundError(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusNotFound {
		return true
	}

	return false
}

// StringPtr returns a pointer to a string if it's non-empty, otherwise nil
func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// ParseUUID parses a UUID string and returns a user-friendly error if invalid
func ParseUUID(id string, entityName string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid %s ID format", entityName))
	}
	return parsed, nil
}
