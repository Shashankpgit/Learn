package adapter

import (
	"app/src/constants"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// GCSConfig holds GCP Cloud Storage-specific configuration
type GCSConfig struct {
	Bucket          string
	CredentialsFile string // Path to service account JSON file
}

// CreateProvider implements StorageConfig interface
func (c GCSConfig) CreateProvider() (StorageProvider, error) {
	return NewGCSAdapter(c)
}

// GCSAdapter implements StorageProvider for Google Cloud Storage
type GCSAdapter struct {
	client *storage.Client
	bucket string
}

// NewGCSAdapter creates a new GCS storage adapter
func NewGCSAdapter(config GCSConfig) (*GCSAdapter, error) {
	ctx := context.Background()

	var client *storage.Client
	var err error

	if config.CredentialsFile != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(config.CredentialsFile))
	} else {
		// Use default credentials (from GOOGLE_APPLICATION_CREDENTIALS env var or ADC)
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", constants.ErrFailedToCreateGCSClient, err)
	}

	return &GCSAdapter{
		client: client,
		bucket: config.Bucket,
	}, nil
}

// Upload uploads a file to GCS and returns a signed URL
func (g *GCSAdapter) Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *UploadOptions) (string, error) {
	bucket := g.client.Bucket(g.bucket)
	obj := bucket.Object(key)
	writer := obj.NewWriter(ctx)

	if opts != nil {
		if opts.ContentType != "" {
			writer.ContentType = opts.ContentType
		}
		if opts.Metadata != nil {
			writer.Metadata = opts.Metadata
		}
	}

	// Copy the data
	written, err := io.Copy(writer, reader)
	if err != nil {
		writer.Close()
		return "", fmt.Errorf("%s: %w", constants.ErrFailedToWriteFileToGCS, err)
	}

	if written != size && size > 0 {
		writer.Close()
		return "", fmt.Errorf(constants.ErrSizeMismatch, written, size)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrFailedToCloseGCSWriter, err)
	}

	// Generate a signed URL valid for 24 hours
	signOpts := &storage.SignedURLOptions{
		Method:  constants.HTTPMethodGET,
		Expires: time.Now().Add(time.Duration(constants.StorageURLExpiration) * time.Hour),
		Scheme:  storage.SigningSchemeV4,
	}

	url, err := bucket.SignedURL(key, signOpts)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrFailedToGenerateSignedURL, err)
	}

	return url, nil
}

// Delete removes a file from GCS
func (g *GCSAdapter) Delete(ctx context.Context, key string) error {
	bucket := g.client.Bucket(g.bucket)
	obj := bucket.Object(key)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("%s: %w", constants.ErrFailedToDeleteFileFromGCS, err)
	}

	return nil
}

// Exists checks if a file exists in GCS
func (g *GCSAdapter) Exists(ctx context.Context, key string) (bool, error) {
	bucket := g.client.Bucket(g.bucket)
	obj := bucket.Object(key)

	_, err := obj.Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", constants.ErrFailedToCheckFileExistence, err)
	}

	return true, nil
}
