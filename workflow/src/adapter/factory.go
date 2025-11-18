package adapter

import (
	"app/src/constants"
	"fmt"
)

// StorageConfig is an interface that each provider config must implement
type StorageConfig interface {
	// CreateProvider creates and returns the storage provider instance
	CreateProvider() (StorageProvider, error)
}

// StorageFactory creates storage providers based on configuration
// This is a factory pattern implementation for dependency injection
type StorageFactory struct {
	storageConfig StorageConfig
}

// NewStorageFactory creates a new storage factory instance
// This is a constructor function for dependency injection
func NewStorageFactory(cfg StorageConfig) *StorageFactory {
	return &StorageFactory{
		storageConfig: cfg,
	}
}

// NewProvider creates and returns a storage provider based on the configuration
func (f *StorageFactory) NewProvider() (StorageProvider, error) {
	if f.storageConfig == nil {
		return nil, fmt.Errorf(constants.ErrStorageConfigNil)
	}
	return f.storageConfig.CreateProvider()
}

// NewStorageProviderFromConfig creates the appropriate storage provider based on configuration
// This is kept for backward compatibility but StorageFactory should be used with DI
func NewStorageProviderFromConfig(cfg StorageConfig) (StorageProvider, error) {
	if cfg == nil {
		return nil, fmt.Errorf(constants.ErrStorageConfigNil)
	}
	return cfg.CreateProvider()
}
