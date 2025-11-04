package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Token represents the tokens table structure
type Token struct {
	TokenID       uuid.UUID         `gorm:"column:token_id;type:uuid;primaryKey" json:"token_id"`
	AccountID     uuid.UUID         `gorm:"column:account_id;type:uuid" json:"account_id"`
	TokenType     string            `gorm:"column:token_type;type:varchar" json:"token_type"`
	IssuerDID     string            `gorm:"column:issuer_did;type:varchar" json:"issuer_did"`
	TokenStandard string            `gorm:"column:token_standard;type:varchar" json:"token_standard"`
	Status        string            `gorm:"column:status;type:varchar" json:"status"`
	Metadata      datatypes.JSONMap `gorm:"column:metadata;type:jsonb" json:"metadata"`
	CreatedAt     time.Time         `gorm:"column:created_at;type:timestamptz;default:now()" json:"created_at"`
}

// TableName specifies the table name for Token model
func (Token) TableName() string {
	return "tokens"
}

// VerifiableCredential represents the verifiable credential structure
type VerifiableCredential struct {
	Context           []string               `json:"@context" validate:"required"`
	ID                string                 `json:"id" validate:"required"`
	Type              []string               `json:"type" validate:"required"`
	Issuer            string                 `json:"issuer" validate:"required"`
	IssuanceDate      string                 `json:"issuanceDate" validate:"required"`
	CredentialSubject map[string]interface{} `json:"credentialSubject" validate:"required"`
	Proof             map[string]interface{} `json:"proof" validate:"required"`
}

// CredentialPayload represents the payload structure
type CredentialPayload struct {
	VerifiableCredential VerifiableCredential `json:"verifiableCredential" validate:"required"`
	DocumentID           string               `json:"documentId" validate:"required,uuid"`
}

// AddCredentialRequest represents the request structure for adding credentials
type AddCredentialRequest struct {
	VerificationType string            `json:"verificationType" validate:"required"`
	Payload          CredentialPayload `json:"payload" validate:"required"`
}
