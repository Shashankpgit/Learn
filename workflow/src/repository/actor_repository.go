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

// ActorRepository defines the interface for actor data access operations
type ActorRepository interface {
	// Create creates a new actor in the database
	Create(ctx context.Context, tx *gorm.DB, actor *model.Actor) error

	// FindByID finds an actor by ID
	FindByID(ctx context.Context, tx *gorm.DB, actorID uuid.UUID) (*model.Actor, error)

	// FindByIDForUpdate finds an actor by ID for update operations (returns 401 instead of 404)
	FindByIDForUpdate(ctx context.Context, tx *gorm.DB, actorID uuid.UUID) (*model.Actor, error)

	// FindByEmail finds an actor by email
	FindByEmail(ctx context.Context, tx *gorm.DB, email string) (*model.Actor, error)

	// Update updates an actor in the database
	Update(ctx context.Context, tx *gorm.DB, actor *model.Actor) error

	// ExistsWithEmail checks if an actor with the given email exists
	ExistsWithEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error)

	// ExistsWithMasterPublicKey checks if an actor with the given master public key exists
	ExistsWithMasterPublicKey(ctx context.Context, tx *gorm.DB, masterPublicKey string) (bool, error)

	// ExistsWithPhoneNumber checks if an actor (excluding given actorID) with the phone number exists
	ExistsWithPhoneNumber(ctx context.Context, tx *gorm.DB, phoneNumber string, excludeActorID uuid.UUID) (bool, error)
}

type actorRepository struct {
	db *gorm.DB
}

// NewActorRepository creates a new instance of ActorRepository
func NewActorRepository(db *gorm.DB) ActorRepository {
	return &actorRepository{db: db}
}

func (r *actorRepository) Create(ctx context.Context, tx *gorm.DB, actor *model.Actor) error {
	if err := tx.WithContext(ctx).Create(actor).Error; err != nil {
		return fmt.Errorf("failed to create actor: %w", err)
	}
	return nil
}

func (r *actorRepository) FindByID(ctx context.Context, tx *gorm.DB, actorID uuid.UUID) (*model.Actor, error) {
	var actor model.Actor
	err := tx.WithContext(ctx).Where("actor_id = ?", actorID).First(&actor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, constants.ErrActorNotFound)
		}
		return nil, fmt.Errorf("failed to find actor: %w", err)
	}
	return &actor, nil
}

func (r *actorRepository) FindByIDForUpdate(ctx context.Context, tx *gorm.DB, actorID uuid.UUID) (*model.Actor, error) {
	var actor model.Actor
	err := tx.WithContext(ctx).Where("actor_id = ?", actorID).First(&actor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return 401 instead of 404 for authenticated operations
			// If JWT is valid but actor doesn't exist, token is stale/invalid
			return nil, fiber.NewError(fiber.StatusUnauthorized, constants.ErrUnauthorized)
		}
		return nil, fmt.Errorf("failed to find actor: %w", err)
	}
	return &actor, nil
}

func (r *actorRepository) FindByEmail(ctx context.Context, tx *gorm.DB, email string) (*model.Actor, error) {
	var actor model.Actor
	err := tx.WithContext(ctx).Where("email = ?", email).First(&actor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusNotFound, constants.ErrActorNotFound)
		}
		return nil, fmt.Errorf("failed to find actor by email: %w", err)
	}
	return &actor, nil
}

func (r *actorRepository) Update(ctx context.Context, tx *gorm.DB, actor *model.Actor) error {
	if err := tx.WithContext(ctx).Save(actor).Error; err != nil {
		return fmt.Errorf("failed to update actor: %w", err)
	}
	return nil
}

func (r *actorRepository) ExistsWithEmail(ctx context.Context, tx *gorm.DB, email string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&model.Actor{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}

func (r *actorRepository) ExistsWithMasterPublicKey(ctx context.Context, tx *gorm.DB, masterPublicKey string) (bool, error) {
	if masterPublicKey == "" {
		return false, nil
	}

	var count int64
	err := tx.WithContext(ctx).Model(&model.Actor{}).Where("master_public_key = ?", masterPublicKey).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check master public key existence: %w", err)
	}
	return count > 0, nil
}

func (r *actorRepository) ExistsWithPhoneNumber(ctx context.Context, tx *gorm.DB, phoneNumber string, excludeActorID uuid.UUID) (bool, error) {
	if phoneNumber == "" {
		return false, nil
	}

	var count int64
	err := tx.WithContext(ctx).Model(&model.Actor{}).
		Where("phone_number = ? AND actor_id != ?", phoneNumber, excludeActorID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check phone number existence: %w", err)
	}
	return count > 0, nil
}
