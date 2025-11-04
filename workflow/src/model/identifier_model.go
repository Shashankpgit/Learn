package model

import (
	"app/src/constants"

	"github.com/google/uuid"
)

type Identifier struct {
	Identifier string    `gorm:"primaryKey;not null" json:"identifier"`
	EntityType string    `gorm:"not null" json:"entityType"`
	EntityID   uuid.UUID `gorm:"not null" json:"entityId"`

	// Relationships
	Actor *Actor `gorm:"foreignKey:EntityID" json:"actor,omitempty"`
}

// TableName overrides the table name used by Identifier to `identifiers`
func (Identifier) TableName() string {
	return constants.TableNameIdentifiers
}
