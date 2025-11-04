package database

import (
	"app/src/config"
	"app/src/constants"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase establishes a connection to the database
func NewDatabase(cfg *config.Config, log *logrus.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		TranslateError:         true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pooling
	sqlDB.SetMaxIdleConns(constants.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(constants.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(constants.DBConnMaxLifetime * time.Minute)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connection established successfully")
	return db, nil
}

// Close closes the database connection gracefully
func Close(db *gorm.DB, log *logrus.Logger) error {
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("Error getting database instance: %v", err)
		return err
	}

	if err := sqlDB.Close(); err != nil {
		log.Errorf("Error closing database connection: %v", err)
		return err
	}

	log.Info("Database connection closed successfully")
	return nil
}
