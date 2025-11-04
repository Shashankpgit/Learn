# Refactored Finternet UNITS API - Source

This directory contains the refactored version of the Finternet UNITS API, implementing **Uber's dig** dependency injection framework along with Clean Code and SOLID principles.

## 🎯 Overview

The refactored codebase maintains **100% API compatibility** with the original implementation while significantly improving:

- **Maintainability**: Clear dependency management through constructor injection
- **Testability**: Easy to mock dependencies for unit testing
- **Scalability**: Modular architecture that's easy to extend
- **Code Quality**: Adherence to SOLID principles and Clean Code practices

## 🏗️ Architecture

### Dependency Injection with Uber dig

This project uses [Uber's dig](https://github.com/uber-go/dig) for dependency injection, providing:

- **Constructor-based injection**: All dependencies are injected through constructors
- **Type-safe wiring**: Compile-time safety for dependency resolution
- **Automatic lifecycle management**: dig handles object creation and lifecycle
- **Clear dependency graph**: Easy to understand component relationships

### Project Structure

```
source/
├── adapter/              # Storage adapters (MinIO, S3, GCS)
├── config/               # Configuration management
├── constants/            # Application constants
├── container/            # Dependency injection container
├── controller/           # HTTP request handlers
├── database/             # Database connection and management
├── middleware/           # HTTP middleware (auth, logging, recovery)
├── model/                # Data models
├── repository/           # Data access layer
├── response/             # Response structures
├── router/               # HTTP routing
├── service/              # Business logic layer
├── utils/                # Utility functions
├── validation/           # Request validation
└── main.go               # Application entry point
```

## 🔑 Key Improvements

### 1. **Technology-Agnostic Naming**

Keycloak-specific names have been replaced with generic authentication terminology:

| Original              | Refactored           |
|-----------------------|----------------------|
| `KeycloakService`     | `AuthService`        |
| `KeycloakAuthMiddleware` | `AuthMiddleware` |
| `KeycloakJWTValidator` | `AuthJWTValidator`  |
| `KeycloakClaims`      | `AuthClaims`         |
| `KeycloakTokenResponse` | `AuthTokenResponse` |

This makes the codebase more maintainable and easier to swap authentication providers if needed.

### 2. **Dependency Injection Container**

All dependencies are registered and resolved through the DI container (`container/container.go`):

```go
// Example: Services are automatically wired with their dependencies
service.NewActorService(
    log,                    // injected
    db,                     // injected
    validate,               // injected
    authService,            // injected
    actorRepo,              // injected
    identifierRepo,         // injected
    actorIntegrationRepo,   // injected
)
```

### 3. **Constructor-Based Injection**

Every component follows the `NewXxx` constructor pattern:

```go
// Config
func NewConfig() (*Config, error)

// Database
func NewDatabase(cfg *Config, log *logrus.Logger) (*gorm.DB, error)

// Services
func NewAuthService(cfg *Config, log *logrus.Logger) AuthService
func NewActorService(log, db, validate, authService, repos...) ActorService

// Repositories
func NewActorRepository(db *gorm.DB) ActorRepository

// Controllers
func NewActorController(actorService service.ActorService) *ActorController
```

### 4. **SOLID Principles**

#### Single Responsibility Principle (SRP)
- Each component has a focused, single purpose
- Services handle business logic
- Repositories handle data access
- Controllers handle HTTP concerns
- Middleware handles cross-cutting concerns

#### Open/Closed Principle (OCP)
- Code is open for extension through interfaces
- Closed for modification through dependency injection
- New features can be added without changing existing code

#### Liskov Substitution Principle (LSP)
- All implementations follow their interface contracts
- Interfaces can be swapped without breaking functionality

#### Interface Segregation Principle (ISP)
- Focused interfaces for each component
- No client depends on methods it doesn't use

#### Dependency Inversion Principle (DIP)
- High-level modules depend on abstractions (interfaces)
- Not on concrete implementations
- Exemplified throughout the codebase with interface-based design

### 5. **Clean Code Practices**

- **Meaningful names**: Self-documenting function and variable names
- **Small functions**: Each function does one thing well
- **No side effects**: Functions are predictable and testable
- **Error handling**: Proper error propagation and handling
- **Comments**: Only where necessary to explain "why", not "what"

## 🚀 Running the Application

### Prerequisites

- Go 1.24.0 or higher
- PostgreSQL database
- Authentication provider (Keycloak or compatible OIDC provider)
- Storage provider (MinIO, S3, or GCS)

### Configuration

Copy `.env.example` to `.env` and configure:

```bash
# Application
APP_HOST=localhost
APP_PORT=3000
APP_ENV=dev

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=finternet
DB_USER=postgres
DB_PASSWORD=password

# Authentication (Keycloak or any OIDC provider)
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=finternet
KEYCLOAK_CLIENT_ID=finternet-api
KEYCLOAK_CLIENT_SECRET=your-secret
KEYCLOAK_ADMIN_USER=admin
KEYCLOAK_ADMIN_PASSWORD=admin

# Storage Provider
STORAGE_PROVIDER=minio
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=finternet
MINIO_USE_SSL=false
```

### Build and Run

```bash
# Navigate to source directory
cd source/

# Build the application
go build -o bin/server

# Run the application
./bin/server
```

Or run directly:

```bash
go run main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./service/...
```

## 📝 How Dependency Injection Works

### 1. Container Registration

In `container/container.go`, all dependencies are registered in order:

```go
providers := []interface{}{
    // Core infrastructure
    config.NewConfig,
    utils.NewLogger,
    database.NewDatabase,
    validation.NewValidator,
    
    // Repositories
    repository.NewActorRepository,
    // ... more repositories
    
    // Services
    service.NewAuthService,
    service.NewActorService,
    // ... more services
    
    // Controllers
    controller.NewActorController,
    // ... more controllers
    
    // Router
    router.NewRouter,
}
```

### 2. Automatic Resolution

dig automatically resolves dependencies based on parameter types:

```go
// dig sees that NewActorService needs:
// - *logrus.Logger
// - *gorm.DB
// - *validator.Validate
// - AuthService
// - ActorRepository
// - etc.

// dig automatically finds and injects all these dependencies
func NewActorService(
    log *logrus.Logger,
    db *gorm.DB,
    validate *validator.Validate,
    authService AuthService,
    actorRepo repository.ActorRepository,
    // ... more dependencies
) ActorService {
    return &actorService{
        log:          log,
        db:           db,
        validate:     validate,
        authService:  authService,
        actorRepo:    actorRepo,
    }
}
```

### 3. Application Startup

In `main.go`, the container invokes the run function with all dependencies:

```go
func main() {
    c, err := container.NewContainer()
    if err != nil {
        // handle error
    }
    
    // dig injects all dependencies into run()
    if err := c.Invoke(run); err != nil {
        // handle error
    }
}

func run(
    cfg *config.Config,      // injected by dig
    log *logrus.Logger,      // injected by dig
    db *gorm.DB,             // injected by dig
    app *fiber.App,          // injected by dig
    r *router.Router,        // injected by dig
) error {
    // Application logic
}
```

## 🧪 Testing Benefits

The refactored architecture makes testing much easier:

### Example: Testing ActorService

```go
func TestActorService_RegisterActor(t *testing.T) {
    // Create mocks for dependencies
    mockDB := setupMockDB(t)
    mockAuthService := &MockAuthService{}
    mockActorRepo := &MockActorRepository{}
    // ... more mocks
    
    // Inject mocks through constructor
    service := NewActorService(
        logrus.New(),
        mockDB,
        validator.New(),
        mockAuthService,
        mockActorRepo,
        // ... more mocks
    )
    
    // Test the service
    actor, err := service.RegisterActor(ctx, request)
    assert.NoError(t, err)
    assert.NotNil(t, actor)
}
```

## 🔄 Migration from Original Code

All APIs remain unchanged:

- ✅ Same endpoints
- ✅ Same request/response formats
- ✅ Same authentication flow
- ✅ Same business logic

Only the internal architecture has been improved.

## 📚 Additional Resources

- [Uber dig Documentation](https://pkg.go.dev/go.uber.org/dig)
- [Clean Code Principles](https://github.com/ryanmcdermott/clean-code-javascript)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Go Best Practices](https://golang.org/doc/effective_go)

## 🤝 Contributing

When adding new features:

1. Create a constructor function (`NewXxx`)
2. Accept dependencies as parameters
3. Return an interface, not a concrete type
4. Register the constructor in `container/container.go`
5. Let dig handle the wiring

## 📄 License

MIT License - See LICENSE file for details

---

**Note**: This refactored codebase is functionally identical to the original but provides a much more maintainable and testable foundation for future development.

