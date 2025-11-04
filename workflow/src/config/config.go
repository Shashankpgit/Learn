package config

import (
	"app/src/adapter"
	"app/src/constants"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	IsProd            bool
	AppHost           string
	AppPort           int
	DBHost            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBPort            int
	AuthURL           string
	AuthRealm         string
	AuthClientID      string
	AuthSecret        string
	AuthAdminUser     string
	AuthAdminPassword string
	StorageConfig     adapter.StorageConfig
}

// NewConfig creates and initializes a new Config instance
func NewConfig() (*Config, error) {
	if err := loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	cfg := &Config{
		IsProd:            viper.GetString(constants.EnvAppEnv) == constants.EnvProduction,
		AppHost:           viper.GetString(constants.EnvAppHost),
		AppPort:           viper.GetInt(constants.EnvAppPort),
		DBHost:            viper.GetString(constants.EnvDBHost),
		DBUser:            viper.GetString(constants.EnvDBUser),
		DBPassword:        viper.GetString(constants.EnvDBPassword),
		DBName:            viper.GetString(constants.EnvDBName),
		DBPort:            viper.GetInt(constants.EnvDBPort),
		AuthURL:           viper.GetString(constants.EnvKeycloakURL),
		AuthRealm:         viper.GetString(constants.EnvKeycloakRealm),
		AuthClientID:      viper.GetString(constants.EnvKeycloakClientID),
		AuthSecret:        viper.GetString(constants.EnvKeycloakClientSecret),
		AuthAdminUser:     viper.GetString(constants.EnvKeycloakAdminUser),
		AuthAdminPassword: viper.GetString(constants.EnvKeycloakAdminPassword),
		StorageConfig:     loadStorageConfig(),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadConfig() error {
	// Enable automatic environment variable reading
	// This reads all environment variables automatically, so no config file is needed
	viper.AutomaticEnv()
	return nil
}

func (c *Config) validate() error {
	// Validate required string fields
	requiredFields := map[string]string{
		constants.EnvAppHost:               c.AppHost,
		constants.EnvDBHost:                c.DBHost,
		constants.EnvDBUser:                c.DBUser,
		constants.EnvDBName:                c.DBName,
		constants.EnvKeycloakURL:           c.AuthURL,
		constants.EnvKeycloakRealm:         c.AuthRealm,
		constants.EnvKeycloakClientID:      c.AuthClientID,
		constants.EnvKeycloakClientSecret:  c.AuthSecret,
		constants.EnvKeycloakAdminUser:     c.AuthAdminUser,
		constants.EnvKeycloakAdminPassword: c.AuthAdminPassword,
	}

	var missing []string
	for key, value := range requiredFields {
		if value == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required configuration: %s", strings.Join(missing, ", "))
	}

	return nil
}

func loadStorageConfig() adapter.StorageConfig {
	providerStr := viper.GetString("STORAGE_PROVIDER")

	switch providerStr {
	case "minio":
		return adapter.MinIOConfig{
			Endpoint:  viper.GetString("MINIO_ENDPOINT"),
			AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
			SecretKey: viper.GetString("MINIO_SECRET_KEY"),
			Bucket:    viper.GetString("MINIO_BUCKET"),
			UseSSL:    viper.GetBool("MINIO_USE_SSL"),
		}
	case "s3", "aws":
		return adapter.S3Config{
			Region:          viper.GetString("AWS_REGION"),
			Bucket:          viper.GetString("AWS_BUCKET"),
			AccessKeyID:     viper.GetString("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: viper.GetString("AWS_SECRET_ACCESS_KEY"),
		}
	case "gcs", "gcp":
		return adapter.GCSConfig{
			Bucket:          viper.GetString("GCP_BUCKET"),
			CredentialsFile: viper.GetString("GCP_CREDENTIALS_FILE"),
		}
	default:
		// Default to MinIO
		return adapter.MinIOConfig{
			Endpoint:  viper.GetString("MINIO_ENDPOINT"),
			AccessKey: viper.GetString("MINIO_ACCESS_KEY"),
			SecretKey: viper.GetString("MINIO_SECRET_KEY"),
			Bucket:    viper.GetString("MINIO_BUCKET"),
			UseSSL:    viper.GetBool("MINIO_USE_SSL"),
		}
	}
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(constants.DBDSNFormat,
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
		constants.DBSSLMode, constants.DBTimeZone)
}

// GetServerAddress returns the server address in host:port format
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.AppHost, c.AppPort)
}
