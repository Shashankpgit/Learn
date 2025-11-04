package model

import (
	"app/src/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActorIntegration struct {
	ActorID        uuid.UUID `gorm:"primaryKey;not null" json:"actorId"`
	Provider       string    `gorm:"column:provider;not null" json:"provider"`
	ExternalUserID string    `gorm:"column:external_user_id;not null" json:"externalUserId"`
	LinkedAt       time.Time `gorm:"column:linked_at;autoCreateTime:milli" json:"linkedAt"`

	// Relationships
	Actor *Actor `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

func (integration *ActorIntegration) BeforeCreate(_ *gorm.DB) error {
	integration.LinkedAt = time.Now()
	return nil
}

// TableName overrides the table name used by ActorIntegration to `actor_integrations`
func (ActorIntegration) TableName() string {
	return constants.TableNameActorIntegrations
}
