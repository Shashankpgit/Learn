package adapter

import (
	"app/src/constants"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOConfig holds MinIO-specific configuration
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// CreateProvider implements StorageConfig interface
func (c MinIOConfig) CreateProvider() (StorageProvider, error) {
	return NewMinIOAdapter(c)
}

// MinIOAdapter implements StorageProvider for MinIO
type MinIOAdapter struct {
	client *minio.Client
	bucket string
}

// NewMinIOAdapter creates a new MinIO storage adapter
func NewMinIOAdapter(config MinIOConfig) (*MinIOAdapter, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constants.ErrFailedToCreateMinIOClient, err)
	}

	// Check if bucket exists and create if it doesn't
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constants.ErrFailedToCheckBucketExistence, err)
	}
	if !exists {
		err = client.MakeBucket(ctx, config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", constants.ErrFailedToCreateBucket, err)
		}
	}

	return &MinIOAdapter{
		client: client,
		bucket: config.Bucket,
	}, nil
}

// Upload uploads a file to MinIO and returns a pre-signed URL
func (m *MinIOAdapter) Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *UploadOptions) (string, error) {
	uploadOpts := minio.PutObjectOptions{}

	if opts != nil {
		if opts.ContentType != "" {
			uploadOpts.ContentType = opts.ContentType
		}
		if opts.Metadata != nil {
			uploadOpts.UserMetadata = opts.Metadata
		}
	}

	_, err := m.client.PutObject(ctx, m.bucket, key, reader, size, uploadOpts)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrFailedToUploadFileToMinIO, err)
	}

	// Generate a pre-signed URL valid for 24 hours
	url, err := m.client.PresignedGetObject(ctx, m.bucket, key, time.Duration(constants.StorageURLExpiration)*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrFailedToGeneratePreSignedURL, err)
	}

	return url.String(), nil
}

// Delete removes a file from MinIO
func (m *MinIOAdapter) Delete(ctx context.Context, key string) error {
	err := m.client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", constants.ErrFailedToDeleteFileFromMinIO, err)
	}

	return nil
}

// Exists checks if a file exists in MinIO
func (m *MinIOAdapter) Exists(ctx context.Context, key string) (bool, error) {
	_, err := m.client.StatObject(ctx, m.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == constants.ErrNoSuchKey {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", constants.ErrFailedToCheckFileExistence, err)
	}

	return true, nil
}
