package repository

import (
	"app/src/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// ActorIntegrationRepository defines the interface for actor integration data access operations
type ActorIntegrationRepository interface {
	// Create creates a new actor integration in the database
	Create(ctx context.Context, tx *gorm.DB, integration *model.ActorIntegration) error
	// FindByActorIDAndProvider finds an actor integration by actor ID and provider
	FindByActorIDAndProvider(ctx context.Context, db *gorm.DB, actorID interface{}, provider string) (*model.ActorIntegration, error)
}

type actorIntegrationRepository struct {
	db *gorm.DB
}

// NewActorIntegrationRepository creates a new instance of ActorIntegrationRepository
func NewActorIntegrationRepository(db *gorm.DB) ActorIntegrationRepository {
	return &actorIntegrationRepository{db: db}
}

func (r *actorIntegrationRepository) Create(ctx context.Context, tx *gorm.DB, integration *model.ActorIntegration) error {
	if err := tx.WithContext(ctx).Create(integration).Error; err != nil {
		return fmt.Errorf("failed to create actor integration: %w", err)
	}
	return nil
}

func (r *actorIntegrationRepository) FindByActorIDAndProvider(ctx context.Context, db *gorm.DB, actorID interface{}, provider string) (*model.ActorIntegration, error) {
	var integration model.ActorIntegration
	if err := db.WithContext(ctx).Where("actor_id = ? AND provider = ?", actorID, provider).First(&integration).Error; err != nil {
		return nil, fmt.Errorf("failed to find actor integration: %w", err)
	}
	return &integration, nil
}
