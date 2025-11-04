package router

import (
	"app/src/config"
	"app/src/constants"
	"app/src/controller"
	"app/src/middleware"
	"app/src/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

// Router manages all application routes
type Router struct {
	app                   *fiber.App
	cfg                   *config.Config
	actorController       *controller.ActorController
	credentialsController *controller.CredentialController
	healthCheckController *controller.HealthCheckController
	authMiddleware        *middleware.AuthMiddleware
}

// NewRouter creates a new router instance with all dependencies injected
func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	actorController *controller.ActorController,
	credentialsController *controller.CredentialController,
	healthCheckController *controller.HealthCheckController,
	authMiddleware *middleware.AuthMiddleware,
	middlewareProviders *middleware.MiddlewareProviders,
) *Router {
	r := &Router{
		app:                   app,
		cfg:                   cfg,
		actorController:       actorController,
		credentialsController: credentialsController,
		healthCheckController: healthCheckController,
		authMiddleware:        authMiddleware,
	}

	r.setupMiddleware(middlewareProviders)
	r.setupRoutes()

	return r
}

// setupMiddleware configures global middleware
func (r *Router) setupMiddleware(providers *middleware.MiddlewareProviders) {
	r.app.Use(providers.Logger)
	r.app.Use(helmet.New())
	r.app.Use(compress.New())
	r.app.Use(cors.New())
	r.app.Use(providers.Recover)
}

// setupRoutes configures all application routes
func (r *Router) setupRoutes() {
	v1 := r.app.Group(constants.RouteGroupV1)

	r.setupHealthCheckRoutes(v1)
	r.setupActorRoutes(v1)
	r.setupCredentialsRoutes(v1)

	if !r.cfg.IsProd {
		r.setupDocsRoutes(v1)
	}

	r.app.Use(utils.NotFoundHandler)
}

// setupHealthCheckRoutes sets up health check routes
func (r *Router) setupHealthCheckRoutes(v1 fiber.Router) {
	v1.Group(constants.RouteHealthCheck).Get("/", r.healthCheckController.Check)
}

// setupActorRoutes sets up actor routes
func (r *Router) setupActorRoutes(v1 fiber.Router) {
	actor := v1.Group("/actor")
	auth := r.authMiddleware.Authenticate()

	// Public routes
	actor.Post("/create", r.actorController.RegisterActor)
	actor.Post("/login", r.actorController.Login)
	actor.Post("/forgotPassword", r.actorController.ForgotPassword)
	actor.Post("/resolve", r.actorController.ResolveUniversalIdentifier)

	// Protected routes
	actor.Post("/update", auth, r.actorController.UpdateActor)
	actor.Post("/getProfile", auth, r.actorController.GetProfile)
	actor.Post("/signout", auth, r.actorController.Signout)
}

// setupCredentialsRoutes sets up credentials routes (all protected)
func (r *Router) setupCredentialsRoutes(v1 fiber.Router) {
	credentials := v1.Group("/credentials", r.authMiddleware.Authenticate())

	credentials.Post("/add", r.credentialsController.AddCredential)
	credentials.Post("/list", r.credentialsController.ListCredentials)
	credentials.Post("/get", r.credentialsController.GetCredential)
	credentials.Post("/delete", r.credentialsController.DeleteCredential)
	credentials.Post("/upload", r.credentialsController.UploadFile)
}

// setupDocsRoutes sets up API documentation routes
func (r *Router) setupDocsRoutes(v1 fiber.Router) {
	v1.Group(constants.RouteDocs).Get(constants.RouteDocsWildcard, func(c *fiber.Ctx) error {
		return c.SendString("API Documentation")
	})
}
