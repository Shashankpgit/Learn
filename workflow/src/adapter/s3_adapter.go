package adapter

import (
	"app/src/constants"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Config holds AWS S3-specific configuration
type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

// CreateProvider implements StorageConfig interface
func (c S3Config) CreateProvider() (StorageProvider, error) {
	return NewS3Adapter(c)
}

// S3Adapter implements StorageProvider for AWS S3
type S3Adapter struct {
	client *s3.Client
	bucket string
}

// NewS3Adapter creates a new AWS S3 storage adapter
func NewS3Adapter(cfg S3Config) (*S3Adapter, error) {
	ctx := context.Background()

	// Create AWS config with static credentials
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constants.ErrFailedToLoadAWSConfig, err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Adapter{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

// Upload uploads a file to S3 and returns a pre-signed URL
func (s *S3Adapter) Upload(ctx context.Context, key string, reader io.Reader, size int64, opts *UploadOptions) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	if opts != nil {
		if opts.ContentType != "" {
			input.ContentType = aws.String(opts.ContentType)
		}
		if opts.Metadata != nil {
			input.Metadata = opts.Metadata
		}
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrFailedToUploadFileToS3, err)
	}

	// Generate a pre-signed URL for accessing the file
	presignClient := s3.NewPresignClient(s.client)
	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrFailedToGeneratePreSignedURL, err)
	}

	return request.URL, nil
}

// Delete removes a file from S3
func (s *S3Adapter) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", constants.ErrFailedToDeleteFileFromS3, err)
	}

	return nil
}

// Exists checks if a file exists in S3
func (s *S3Adapter) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Check if it's a not found error
		if err.Error() == constants.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", constants.ErrFailedToCheckFileExistence, err)
	}

	return true, nil
}
