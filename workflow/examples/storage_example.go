package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"app/src/adapter"
)

// Example demonstrating direct usage of the storage provider
func main() {
	// Example 1: Using MinIO
	fmt.Println("=== MinIO Example ===")
	minioExample()

	// Example 2: Using S3
	fmt.Println("\n=== AWS S3 Example ===")
	s3Example()

	// Example 3: Using GCS
	fmt.Println("\n=== Google Cloud Storage Example ===")
	gcsExample()

	// Example 4: Using Factory Pattern
	fmt.Println("\n=== Factory Pattern Example ===")
	factoryExample()
}

func minioExample() {
	// Configure MinIO
	config := adapter.MinIOConfig{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "test-bucket",
		UseSSL:    false,
	}

	// Create adapter
	storage, err := adapter.NewMinIOAdapter(config)
	if err != nil {
		log.Printf("Failed to create MinIO adapter: %v", err)
		return
	}

	// Upload a file
	ctx := context.Background()
	content := []byte("Hello, MinIO!")
	reader := bytes.NewReader(content)

	opts := &adapter.UploadOptions{
		ContentType: "text/plain",
		Metadata: map[string]string{
			"example": "minio",
		},
	}

	url, err := storage.Upload(ctx, "test/hello.txt", reader, int64(len(content)), opts)
	if err != nil {
		log.Printf("Failed to upload: %v", err)
		return
	}
	fmt.Printf("✓ Uploaded file, URL: %s\n", url)

	// Check existence
	exists, err := storage.Exists(ctx, "test/hello.txt")
	if err != nil {
		log.Printf("Failed to check existence: %v", err)
		return
	}
	fmt.Printf("✓ File exists: %v\n", exists)

	// Delete file
	err = storage.Delete(ctx, "test/hello.txt")
	if err != nil {
		log.Printf("Failed to delete: %v", err)
		return
	}
	fmt.Printf("✓ File deleted successfully\n")
}

func s3Example() {
	// Configure S3
	config := adapter.S3Config{
		Region:          "ap-south-1",
		Bucket:          "my-app-bucket",
		AccessKeyID:     "YOUR_ACCESS_KEY",
		SecretAccessKey: "YOUR_SECRET_KEY",
	}

	// Create adapter
	storage, err := adapter.NewS3Adapter(config)
	if err != nil {
		log.Printf("Failed to create S3 adapter: %v", err)
		return
	}

	// Upload a file
	ctx := context.Background()
	content := []byte("Hello, S3!")
	reader := bytes.NewReader(content)

	opts := &adapter.UploadOptions{
		ContentType: "text/plain",
		Metadata: map[string]string{
			"example": "s3",
		},
	}

	url, err := storage.Upload(ctx, "test/hello.txt", reader, int64(len(content)), opts)
	if err != nil {
		log.Printf("Failed to upload: %v", err)
		return
	}
	fmt.Printf("✓ Uploaded file, pre-signed URL: %s\n", url)
}

func gcsExample() {
	// Configure GCS
	config := adapter.GCSConfig{
		Bucket:          "my-gcs-bucket",
		CredentialsFile: "/path/to/service-account.json",
	}

	// Create adapter
	storage, err := adapter.NewGCSAdapter(config)
	if err != nil {
		log.Printf("Failed to create GCS adapter: %v", err)
		return
	}

	// Upload a file
	ctx := context.Background()
	content := []byte("Hello, GCS!")
	reader := bytes.NewReader(content)

	opts := &adapter.UploadOptions{
		ContentType: "text/plain",
		Metadata: map[string]string{
			"example": "gcs",
		},
	}

	url, err := storage.Upload(ctx, "test/hello.txt", reader, int64(len(content)), opts)
	if err != nil {
		log.Printf("Failed to upload: %v", err)
		return
	}
	fmt.Printf("✓ Uploaded file, signed URL: %s\n", url)
}

func factoryExample() {
	// Using factory pattern - switches provider based on config
	// Each config type knows how to create its own provider
	storageConfig := adapter.MinIOConfig{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Bucket:    "test-bucket",
		UseSSL:    false,
	}

	// Create storage provider (config creates the appropriate adapter)
	storage, err := adapter.NewStorageProviderFromConfig(storageConfig)
	if err != nil {
		log.Printf("Failed to create storage provider: %v", err)
		return
	}

	// Use the provider (same interface regardless of actual provider)
	ctx := context.Background()
	content := []byte("Hello from factory!")
	reader := bytes.NewReader(content)

	url, err := storage.Upload(ctx, "factory/test.txt", reader, int64(len(content)), nil)
	if err != nil {
		log.Printf("Failed to upload: %v", err)
		return
	}
	fmt.Printf("✓ Uploaded via factory pattern, URL: %s\n", url)

	// The beauty is: to switch providers, just change the config type
	// Use: adapter.MinIOConfig, adapter.S3Config, or adapter.GCSConfig
	// No other code changes needed!
	fmt.Println("✓ Same code works with any provider!")
}
