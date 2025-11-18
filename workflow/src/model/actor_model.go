package model

import (
	"app/src/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Actor struct {
	ActorID                uuid.UUID `gorm:"primaryKey;not null" json:"id"`
	DID                    uuid.UUID `gorm:"column:did;type:uuid;uniqueIndex;not null" json:"did"`
	Email                  string    `gorm:"uniqueIndex;not null" json:"email"`
	FirstName              string    `gorm:"column:first_name;not null" json:"firstName"`
	LastName               string    `gorm:"column:last_name;not null" json:"lastName"`
	PhoneNumber            *string   `gorm:"column:phone_number;uniqueIndex" json:"phoneNumber,omitempty"`
	MasterPublicKey        string    `gorm:"column:master_public_key;uniqueIndex;not null" json:"masterPublicKey"`
	EntityType             string    `gorm:"column:entity_type;type:varchar(50);not null" json:"entityType"`
	Nationality            *string   `gorm:"column:nationality;type:varchar(10)" json:"nationality,omitempty"`
	CountryOfResidence     *string   `gorm:"column:country_of_residence;type:varchar(10)" json:"countryOfResidence,omitempty"`
	CountryOfIncorporation *string   `gorm:"column:country_of_incorporation;type:varchar(10)" json:"countryOfIncorporation,omitempty"`
	VerificationLevel      string    `gorm:"column:verification_level;type:varchar(50);not null;default:'Tier0_Unverified'" json:"verificationLevel"`
	CreatedAt              time.Time `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`

	// Relationships
	Identifiers      []Identifier      `gorm:"foreignKey:EntityID" json:"identifiers,omitempty"`
	ActorIntegration *ActorIntegration `gorm:"foreignKey:ActorID" json:"actorIntegration,omitempty"`
}

func (actor *Actor) BeforeCreate(_ *gorm.DB) error {
	actorID, err := uuid.NewV7()

	if err != nil {
		return err
	}

	actor.ActorID = actorID
	actor.DID = uuid.New()
	return nil
}

// TableName overrides the table name used by Actor to `actors`
func (Actor) TableName() string {
	return constants.TableNameActors
}
