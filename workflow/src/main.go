package main

import (
	"app/src/config"
	"app/src/container"
	"app/src/database"
	"app/src/router"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// @title Finternet UNITS APIs Documentation
// @version 1.0.0
// @license.name MIT
// @host localhost:3000
// @BasePath /v1
func main() {
	// Create the dependency injection container
	c, err := container.NewContainer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create container: %v\n", err)
		os.Exit(1)
	}

	// Start the application using the container
	if err := c.Invoke(run); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}
}

// run is the main application function with all dependencies injected by dig
func run(
	cfg *config.Config,
	log *logrus.Logger,
	db *gorm.DB,
	app *fiber.App,
	r *router.Router, // Router instance - ensures routes are initialized
) error {
	log.Info("Application starting...")

	// Get server address
	address := cfg.GetServerAddress()

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	
	go func() {
		log.Infof("Starting server on %s", address)
		if err := app.Listen(address); err != nil {
			serverErrors <- fmt.Errorf("error starting server: %w", err)
		}
	}()

	// Handle graceful shutdown
	return handleGracefulShutdown(log, db, app, serverErrors)
}

// handleGracefulShutdown handles graceful shutdown of the application
func handleGracefulShutdown(
	log *logrus.Logger,
	db *gorm.DB,
	app *fiber.App,
	serverErrors <-chan error,
) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Errorf("Server error: %v", err)
		return err
	case sig := <-quit:
		log.Infof("Received signal: %v. Shutting down server...", sig)

		// Close database connection
		if err := database.Close(db, log); err != nil {
			log.Errorf("Error closing database: %v", err)
		}

		// Shutdown fiber app
		if err := app.Shutdown(); err != nil {
			log.Errorf("Error during server shutdown: %v", err)
			return err
		}

		log.Info("Server exited gracefully")
		return nil
	}
}
