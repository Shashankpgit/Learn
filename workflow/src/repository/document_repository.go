package repository

import (
	"app/src/model"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentRepository defines the interface for document data access operations
type DocumentRepository interface {
	// Create creates a new document record in the database
	Create(ctx context.Context, tx *gorm.DB, document *model.Document) error

	// FindByID finds a document by ID
	FindByID(ctx context.Context, tx *gorm.DB, documentID uuid.UUID) (*model.Document, error)

	// FindByPath finds a document by its path
	FindByPath(ctx context.Context, tx *gorm.DB, path string) (*model.Document, error)

	// Delete deletes a document by ID
	Delete(ctx context.Context, tx *gorm.DB, documentID uuid.UUID) error
}

type documentRepository struct {
	db *gorm.DB
}

// NewDocumentRepository creates a new instance of DocumentRepository
func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(ctx context.Context, tx *gorm.DB, document *model.Document) error {
	if err := tx.WithContext(ctx).Create(document).Error; err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}
	return nil
}

func (r *documentRepository) FindByID(ctx context.Context, tx *gorm.DB, documentID uuid.UUID) (*model.Document, error) {
	var document model.Document
	if err := tx.WithContext(ctx).Where("document_id = ?", documentID).First(&document).Error; err != nil {
		return nil, fmt.Errorf("failed to find document by ID: %w", err)
	}
	return &document, nil
}

func (r *documentRepository) FindByPath(ctx context.Context, tx *gorm.DB, path string) (*model.Document, error) {
	var document model.Document
	if err := tx.WithContext(ctx).Where("storage_path = ?", path).First(&document).Error; err != nil {
		return nil, fmt.Errorf("failed to find document by path: %w", err)
	}
	return &document, nil
}

func (r *documentRepository) Delete(ctx context.Context, tx *gorm.DB, documentID uuid.UUID) error {
	if err := tx.WithContext(ctx).Where("document_id = ?", documentID).Delete(&model.Document{}).Error; err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}
