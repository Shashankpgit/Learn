package test

import (
	"app/src/database"
	"app/src/middleware"
	"app/src/router"
	"app/src/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var App = fiber.New(fiber.Config{
	CaseSensitive: true,
	ErrorHandler:  utils.ErrorHandler,
})
var DB *gorm.DB
var Log = utils.Log

func init() {
	// TODO: You can modify host and database configuration for tests
	var err error
	DB, err = database.Connect("localhost", "testdb")
	if err != nil {
		Log.Fatalf("Failed to connect to test database: %v", err)
	}
	
	// Initialize Keycloak JWT validator
	validator, err := middleware.NewKeycloakJWTValidator(DB)
	if err != nil {
		Log.Fatalf("Failed to initialize Keycloak validator: %v", err)
	}
	
	// Create Keycloak auth middleware
	keycloakAuth := middleware.NewKeycloakAuthMiddleware(validator)
	
	if err := router.Routes(App, DB, keycloakAuth); err != nil {
		Log.Fatalf("Failed to setup routes: %v", err)
	}
	
	App.Use(utils.NotFoundHandler)
}
