package repository

import (
	"app/src/constants"
	"app/src/model"
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CredentialsRepository defines the interface for credentials data access operations
type CredentialsRepository interface {
	// Create creates a new credential token in the database
	Create(ctx context.Context, tx *gorm.DB, token *model.Token) error

	// FindByID finds a credential by token ID
	FindByID(ctx context.Context, tx *gorm.DB, tokenID uuid.UUID) (*model.Token, error)

	// FindAll retrieves all credentials
	FindAll(ctx context.Context, tx *gorm.DB) ([]model.Token, error)

	// Delete deletes a credential by token ID
	Delete(ctx context.Context, tx *gorm.DB, tokenID uuid.UUID) error
}

type credentialsRepository struct {
	db *gorm.DB
}

// NewCredentialsRepository creates a new instance of CredentialsRepository
func NewCredentialsRepository(db *gorm.DB) CredentialsRepository {
	return &credentialsRepository{db: db}
}

func (r *credentialsRepository) Create(ctx context.Context, tx *gorm.DB, token *model.Token) error {
	if err := tx.WithContext(ctx).Create(token).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fiber.NewError(fiber.StatusConflict, constants.ErrTokenAlreadyExists)
		}
		return fmt.Errorf("failed to create credential: %w", err)
	}
	return nil
}

func (r *credentialsRepository) FindByID(ctx context.Context, tx *gorm.DB, tokenID uuid.UUID) (*model.Token, error) {
	var token model.Token
	err := tx.WithContext(ctx).Where("token_id = ?", tokenID).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, constants.ErrCredentialNotFound)
		}
		return nil, fmt.Errorf("failed to find credential: %w", err)
	}
	return &token, nil
}

func (r *credentialsRepository) FindAll(ctx context.Context, tx *gorm.DB) ([]model.Token, error) {
	var tokens []model.Token
	err := tx.WithContext(ctx).Find(&tokens).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve credentials: %w", err)
	}
	return tokens, nil
}

func (r *credentialsRepository) Delete(ctx context.Context, tx *gorm.DB, tokenID uuid.UUID) error {
	result := tx.WithContext(ctx).Delete(&model.Token{}, "token_id = ?", tokenID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete credential: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, constants.ErrCredentialNotFound)
	}
	return nil
}
