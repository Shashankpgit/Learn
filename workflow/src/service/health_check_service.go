package service

import (
	"app/src/constants"
	"errors"
	"runtime"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// HealthCheckService defines the interface for health check operations
type HealthCheckService interface {
	GormCheck() error
	MemoryHeapCheck() error
}

// healthCheckService implements HealthCheckService with dependency injection
type healthCheckService struct {
	log *logrus.Logger
	db  *gorm.DB
}

// NewHealthCheckService creates a new health check service instance
// This is a constructor function for dependency injection
func NewHealthCheckService(log *logrus.Logger, db *gorm.DB) HealthCheckService {
	return &healthCheckService{
		log: log,
		db:  db,
	}
}

func (s *healthCheckService) GormCheck() error {
	sqlDB, errDB := s.db.DB()
	if errDB != nil {
		s.log.Errorf("failed to access the database connection pool: %v", errDB)
		return errDB
	}

	if err := sqlDB.Ping(); err != nil {
		s.log.Errorf("failed to ping the database: %v", err)
		return err
	}

	return nil
}

// MemoryHeapCheck checks if heap memory usage exceeds a threshold
func (s *healthCheckService) MemoryHeapCheck() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats) // Collect memory statistics

	heapAlloc := memStats.HeapAlloc // Heap memory currently allocated
	heapThreshold := uint64(constants.HealthHeapThreshold)

	s.log.Infof("Heap Memory Allocation: %v bytes", heapAlloc)

	// If the heap allocation exceeds the threshold, return an error
	if heapAlloc > heapThreshold {
		s.log.Errorf("Heap memory usage exceeds threshold: %v bytes", heapAlloc)
		return errors.New(constants.HealthHeapErrorMsg)
	}

	return nil
}
