package adapter_test

import (
	"bytes"
	"context"
	"testing"

	"app/src/adapter"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStorageProviderInterface verifies all adapters implement the interface correctly
func TestStorageProviderInterface(t *testing.T) {
	var _ adapter.StorageProvider = (*adapter.MinIOAdapter)(nil)
	var _ adapter.StorageProvider = (*adapter.S3Adapter)(nil)
	var _ adapter.StorageProvider = (*adapter.GCSAdapter)(nil)
}

// testStorageOperations is a common test suite for all storage providers
func testStorageOperations(t *testing.T, storage adapter.StorageProvider) {
	ctx := context.Background()
	testKey := "test/example.txt"
	testContent := []byte("Hello, Storage!")

	t.Run("Upload file", func(t *testing.T) {
		reader := bytes.NewReader(testContent)
		opts := &adapter.UploadOptions{
			ContentType: "text/plain",
			Metadata: map[string]string{
				"test": "true",
			},
		}

		url, err := storage.Upload(ctx, testKey, reader, int64(len(testContent)), opts)
		require.NoError(t, err)
		assert.NotEmpty(t, url)
		t.Logf("File uploaded with URL: %s", url)
	})

	t.Run("Check file exists", func(t *testing.T) {
		exists, err := storage.Exists(ctx, testKey)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Delete file", func(t *testing.T) {
		err := storage.Delete(ctx, testKey)
		require.NoError(t, err)
	})

	t.Run("Verify file deleted", func(t *testing.T) {
		exists, err := storage.Exists(ctx, testKey)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}
