package model

import (
	"time"

	"github.com/google/uuid"
)

// Document represents the documents table structure
type Document struct {
	DocumentID  uuid.UUID `gorm:"type:uuid;primaryKey;column:document_id" json:"document_id"`
	AccountID   uuid.UUID `gorm:"type:uuid;not null;column:account_id" json:"account_id"`
	FileName    string    `gorm:"type:varchar(255);not null;column:file_name" json:"file_name"`
	StoragePath string    `gorm:"type:text;not null;column:storage_path" json:"storage_path"`
	MimeType    *string   `gorm:"type:varchar(255);column:mime_type" json:"mime_type,omitempty"`
	UploadedAt  time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;not null;column:uploaded_at" json:"uploaded_at"`
}

// TableName specifies the table name for GORM
func (Document) TableName() string {
	return "documents"
}
