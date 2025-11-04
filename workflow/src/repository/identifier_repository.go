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

// IdentifierRepository defines the interface for identifier data access operations
type IdentifierRepository interface {
	// Create creates a new identifier in the database
	Create(ctx context.Context, tx *gorm.DB, identifier *model.Identifier) error

	// FindByValue finds an identifier by its value
	FindByValue(ctx context.Context, tx *gorm.DB, identifierValue string) (*model.Identifier, error)

	// FindByActorID finds an identifier by actor ID
	FindByActorID(ctx context.Context, tx *gorm.DB, actorID uuid.UUID) (*model.Identifier, error)

	// ExistsWithIdentifier checks if an identifier already exists
	ExistsWithIdentifier(ctx context.Context, tx *gorm.DB, identifier string) (bool, error)
}

type identifierRepository struct {
	db *gorm.DB
}

// NewIdentifierRepository creates a new instance of IdentifierRepository
func NewIdentifierRepository(db *gorm.DB) IdentifierRepository {
	return &identifierRepository{db: db}
}

func (r *identifierRepository) Create(ctx context.Context, tx *gorm.DB, identifier *model.Identifier) error {
	if err := tx.WithContext(ctx).Create(identifier).Error; err != nil {
		return fmt.Errorf("failed to create identifier: %w", err)
	}
	return nil
}

func (r *identifierRepository) FindByValue(ctx context.Context, tx *gorm.DB, identifierValue string) (*model.Identifier, error) {
	var identifier model.Identifier
	err := tx.WithContext(ctx).Where("identifier = ?", identifierValue).First(&identifier).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, constants.ErrResourceNotFound)
		}
		return nil, fmt.Errorf("failed to find identifier: %w", err)
	}
	return &identifier, nil
}

func (r *identifierRepository) FindByActorID(ctx context.Context, tx *gorm.DB, actorID uuid.UUID) (*model.Identifier, error) {
	var identifier model.Identifier
	err := tx.WithContext(ctx).
		Where("entity_id = ? AND entity_type = ?", actorID, constants.EntityTypeActor).
		First(&identifier).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil without error if not found (optional identifier)
		}
		return nil, fmt.Errorf("failed to find identifier for actor: %w", err)
	}
	return &identifier, nil
}

func (r *identifierRepository) ExistsWithIdentifier(ctx context.Context, tx *gorm.DB, identifier string) (bool, error) {
	if identifier == "" {
		return false, nil
	}

	var count int64
	err := tx.WithContext(ctx).Model(&model.Identifier{}).Where("identifier = ?", identifier).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check identifier existence: %w", err)
	}
	return count > 0, nil
}
