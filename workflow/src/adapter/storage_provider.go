package adapter

import (
	"context"
	"io"
)

// UploadOptions contains optional parameters for file upload
type UploadOptions struct {
	ContentType string
	Metadata    map[string]string
}

// StorageProvider defines the interface for file storage operations
type StorageProvider interface {
	// Upload uploads a file to the storage provider
	// Returns the accessible URL (pre-signed for private storage) and any error
	Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *UploadOptions) (string, error)

	// Delete removes a file from storage
	Delete(ctx context.Context, key string) error

	// Exists checks if a file exists in storage
	Exists(ctx context.Context, key string) (bool, error)
}
