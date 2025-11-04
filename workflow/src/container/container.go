package container

import (
	"app/src/adapter"
	"app/src/config"
	"app/src/controller"
	"app/src/database"
	"app/src/middleware"
	"app/src/repository"
	"app/src/router"
	"app/src/service"
	"app/src/utils"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

// Container holds the dependency injection container
type Container struct {
	*dig.Container
}

// NewContainer creates and configures a new dependency injection container
func NewContainer() (*Container, error) {
	c := dig.New()

	// Register all providers in dependency order
	providers := []interface{}{
		// Core infrastructure
		config.NewConfig,
		utils.NewLogger,
		utils.NewResponseBuilder,
		database.NewDatabase,
		validation.NewValidator,
		ProvideStorageFactory,

		// Repositories
		repository.NewActorRepository,
		repository.NewIdentifierRepository,
		repository.NewActorIntegrationRepository,
		repository.NewCredentialsRepository,
		repository.NewDocumentRepository,

		// Services
		service.NewAuthService,
		service.NewActorService,
		service.NewCredentialsService,
		service.NewHealthCheckService,

		// Middleware
		middleware.NewAuthJWTValidator,
		middleware.NewAuthMiddleware,
		middleware.NewMiddlewareProviders,

		// Controllers
		controller.NewActorController,
		controller.NewCredentialsController,
		controller.NewHealthCheckController,

		// Router
		router.NewRouter,

		// Fiber app
		NewFiberApp,
	}

	for _, provider := range providers {
		if err := c.Provide(provider); err != nil {
			return nil, err
		}
	}

	return &Container{Container: c}, nil
}

// ProvideStorageFactory creates a storage factory from configuration
func ProvideStorageFactory(cfg *config.Config) *adapter.StorageFactory {
	return adapter.NewStorageFactory(cfg.StorageConfig)
}

// NewFiberApp creates a new Fiber application
func NewFiberApp(cfg *config.Config) *fiber.App {
	return fiber.New(config.FiberConfig(cfg))
}
