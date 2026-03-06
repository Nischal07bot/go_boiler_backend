# Go Production Boilerplate

A production-ready Go backend boilerplate built with clean architecture principles. Features observability, authentication, background jobs, and more.

```
+=========================================================================+
|                     GO PRODUCTION BOILERPLATE                           |
+=========================================================================+
|                                                                         |
|  +-------------------------------------------------------------------+  |
|  |                        PRESENTATION LAYER                         |  |
|  |  +-------------------------------------------------------------+  |  |
|  |  |  Echo HTTP Router  |  Middlewares  |  Typed Handlers        |  |  |
|  |  +-------------------------------------------------------------+  |  |
|  +-------------------------------------------------------------------+  |
|                                   |                                     |
|                                   v                                     |
|  +-------------------------------------------------------------------+  |
|  |                        BUSINESS LAYER                             |  |
|  |  +-------------------------------------------------------------+  |  |
|  |  |  Services  |  Validation  |  Error Handling                 |  |  |
|  |  +-------------------------------------------------------------+  |  |
|  +-------------------------------------------------------------------+  |
|                                   |                                     |
|                                   v                                     |
|  +-------------------------------------------------------------------+  |
|  |                         DATA LAYER                                |  |
|  |  +-------------------------------------------------------------+  |  |
|  |  |  Repositories  |  PostgreSQL  |  Redis  |  Migrations       |  |  |
|  |  +-------------------------------------------------------------+  |  |
|  +-------------------------------------------------------------------+  |
|                                   |                                     |
|                                   v                                     |
|  +-------------------------------------------------------------------+  |
|  |                      INFRASTRUCTURE                               |  |
|  |  +------------------+  +---------------+  +--------------------+  |  |
|  |  |   Background     |  |   Email       |  |   Observability    |  |  |
|  |  |   Jobs (Asynq)   |  |   (Resend)    |  |   (New Relic)      |  |  |
|  |  +------------------+  +---------------+  +--------------------+  |  |
|  +-------------------------------------------------------------------+  |
|                                                                         |
+=========================================================================+
```

---

## Features

- **Clean Architecture** - Separation of concerns with handler, service, repository layers
- **Dependency Injection** - Central server struct for easy testing and modularity
- **Echo Framework** - High-performance HTTP routing with typed handlers
- **PostgreSQL** - pgx/v5 driver with connection pooling
- **Redis** - Caching and job queue backend
- **Background Jobs** - Asynq-powered async task processing
- **Authentication** - Clerk SDK integration with JWT validation
- **Email Service** - Resend API with HTML templates
- **Observability** - New Relic APM, distributed tracing, log forwarding
- **Structured Logging** - Zerolog with context enrichment
- **Migrations** - Tern-based SQL migrations
- **Validation** - go-playground/validator with custom error messages
- **Error Handling** - Typed HTTP errors with SQL error mapping
- **Rate Limiting** - Configurable with custom metrics
- **Security** - CORS, secure headers, request ID tracking

---

## Tech Stack

| Category | Technology |
|----------|------------|
| Language | Go 1.21+ |
| Framework | Echo v4 |
| Database | PostgreSQL (pgx/v5) |
| Cache/Queue | Redis (go-redis/v9) |
| Jobs | Asynq |
| Auth | Clerk SDK |
| Email | Resend |
| Observability | New Relic |
| Logging | Zerolog |
| Validation | go-playground/validator |
| Config | Koanf |
| Migrations | Tern |
| Build | Taskfile |

---

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Task (taskfile.dev)

### Installation

```bash
# Clone the repository
git clone https://github.com/Nischal07bot/go_boiler_backend.git
cd go_boiler_backend/backend

# Install dependencies
go mod tidy

# Copy environment file
cp .env.example .env
# Edit .env with your configuration

# Run migrations
task migrations:up

# Start the server
task run
```

---

## Project Structure

```
golang_boilerplate/
+-- backend/
|   +-- cmd/
|   |   +-- golang-boilerplate/
|   |       +-- main.go                 # Application entry point
|   +-- internal/
|   |   +-- config/
|   |   |   +-- config.go               # Configuration loading (Koanf)
|   |   |   +-- observablity.go         # New Relic setup
|   |   +-- database/
|   |   |   +-- database.go             # PostgreSQL connection pool
|   |   |   +-- migrator.go             # Tern migrations runner
|   |   |   +-- migrations/             # SQL migration files
|   |   +-- errs/
|   |   |   +-- http.go                 # HTTP error types
|   |   |   +-- types.go                # Error constants
|   |   +-- handler/
|   |   |   +-- base.go                 # Generic typed handler pattern
|   |   |   +-- handlers.go             # Handler registry
|   |   |   +-- hello.go                # Example handler
|   |   +-- lib/
|   |   |   +-- email/                  # Resend email service
|   |   |   +-- job/                    # Asynq background jobs
|   |   |   +-- utils/                  # Shared utilities
|   |   +-- logger/
|   |   |   +-- logger.go               # Zerolog + New Relic integration
|   |   +-- middlewares/
|   |   |   +-- auth.go                 # Clerk authentication
|   |   |   +-- context.go              # Context enrichment
|   |   |   +-- global.go               # CORS, logging, recovery
|   |   |   +-- middleware.go           # Middleware registry
|   |   |   +-- rate_limit.go           # Rate limiting metrics
|   |   |   +-- request_id.go           # Request ID generation
|   |   |   +-- tracing.go              # New Relic tracing
|   |   +-- models/                     # Database models
|   |   +-- repository/                 # Data access layer
|   |   +-- router/
|   |   |   +-- router.go               # Route registration
|   |   |   +-- system.go               # System routes
|   |   +-- server/
|   |   |   +-- server.go               # Central DI container
|   |   +-- service/
|   |   |   +-- auth.go                 # Auth service
|   |   |   +-- services.go             # Service registry
|   |   +-- sqlerr/                     # SQL error handling
|   |   +-- templates/                  # HTML templates
|   |   +-- validation/                 # Request validation
|   +-- go.mod
|   +-- Taskfile.yml                    # Task runner commands
+-- packages/                           # Shared packages
+-- turbo.json                          # Turborepo config
+-- package.json                        # Monorepo root
```

---

## Architecture

### Request Flow

```
                           +-------------+
                           |   Client    |
                           +------+------+
                                  |
                                  | HTTP Request
                                  v
+-------------------------------------------------------------------------------------+
|                              ECHO ROUTER                                            |
|  +-------------------------------------------------------------------------------+  |
|  |                         MIDDLEWARE CHAIN                                      |  |
|  |  +--------+  +------+  +--------+  +-----------+  +--------+  +----------+   |  |
|  |  | Rate   |->| CORS |->| Secure |->| RequestID |->| NewRelic|->| Logger  |   |  |
|  |  | Limit  |  |      |  |        |  |           |  | Tracing |  |         |   |  |
|  |  +--------+  +------+  +--------+  +-----------+  +--------+  +----------+   |  |
|  +-------------------------------------------------------------------------------+  |
+-------------------------------------------------------------------------------------+
                                  |
                                  v
                    +-------------+-------------+
                    |        Handler            |
                    |  (Bind + Validate + Log)  |
                    +-------------+-------------+
                                  |
                                  v
                    +-------------+-------------+
                    |        Service            |
                    |    (Business Logic)       |
                    +-------------+-------------+
                                  |
                                  v
                    +-------------+-------------+
                    |       Repository          |
                    |     (Data Access)         |
                    +-------------+-------------+
                                  |
                  +---------------+---------------+
                  |                               |
                  v                               v
          +-------+-------+               +-------+-------+
          |  PostgreSQL   |               |     Redis     |
          |   (Primary)   |               |   (Cache)     |
          +---------------+               +---------------+
```

### Layer Details

```
+------------------------------------------------------------------+
|                          HANDLERS                                 |
|  +------------------------------------------------------------+  |
|  |   HelloHandler   |   AuthHandler   |   UserHandler   | ... |  |
|  +------------------------------------------------------------+  |
|  |              Base Handler (Server Access)                  |  |
|  +------------------------------------------------------------+  |
+------------------------------------------------------------------+
                                |
                                v
+------------------------------------------------------------------+
|                          SERVICES                                 |
|  +------------------------------------------------------------+  |
|  |              Business Logic Layer                          |  |
|  |          AuthService | UserService | ...                   |  |
|  +------------------------------------------------------------+  |
+------------------------------------------------------------------+
                                |
                                v
+------------------------------------------------------------------+
|                        REPOSITORIES                               |
|  +------------------------------------------------------------+  |
|  |              Data Access Layer                             |  |
|  |          UserRepo | ProductRepo | ...                      |  |
|  +------------------------------------------------------------+  |
+------------------------------------------------------------------+
                                |
                                v
+------------------------------------------------------------------+
|                        DATA STORES                                |
|  +------------------------+  +-----------------------------+     |
|  |      PostgreSQL        |  |           Redis             |     |
|  |   (Primary Database)   |  |   (Cache + Job Queue)       |     |
|  +------------------------+  +-----------------------------+     |
+------------------------------------------------------------------+
```

### Background Processing

```
+------------------------------------------------------------------+
|                    BACKGROUND PROCESSING                          |
|  +------------------------------------------------------------+  |
|  |                    Asynq Job Queue                         |  |
|  |  +----------+  +----------+  +----------+                  |  |
|  |  | Critical |  | Default  |  |   Low    |  <- Priority     |  |
|  |  | (Weight 6)| | (Weight 3)| | (Weight 1)|     Queues      |  |
|  |  +----------+  +----------+  +----------+                  |  |
|  +------------------------------------------------------------+  |
|  |                    Job Handlers                            |  |
|  |            EmailTask | NotificationTask | ...              |  |
|  +------------------------------------------------------------+  |
+------------------------------------------------------------------+
```

### Observability

```
+------------------------------------------------------------------+
|                      OBSERVABILITY                                |
|  +--------------------+  +-------------------+  +--------------+  |
|  |     New Relic      |  |     Zerolog       |  |   Metrics    |  |
|  |  - APM             |  |  - Structured     |  |  - Rate Limit|  |
|  |  - Distributed     |  |  - Context-aware  |  |  - Latency   |  |
|  |    Tracing         |  |  - Trace IDs      |  |  - Errors    |  |
|  |  - Log Forwarding  |  |  - Stack Traces   |  |              |  |
|  +--------------------+  +-------------------+  +--------------+  |
+------------------------------------------------------------------+
```

---

## Configuration

All configuration is done via environment variables with `BOILERPLATE_` prefix:

```env
# Server
BOILERPLATE_SERVER_PORT=8080
BOILERPLATE_SERVER_CORS_ALLOWED_ORIGINS=http://localhost:3000

# Database
BOILERPLATE_DATABASE_NAME=mydb
BOILERPLATE_DATABASE_HOST=localhost
BOILERPLATE_DATABASE_PORT=5432
BOILERPLATE_DATABASE_USER=postgres
BOILERPLATE_DATABASE_PASSWORD=secret
BOILERPLATE_DATABASE_SSL_MODE=disable

# Redis
BOILERPLATE_REDIS_ADDRESS=localhost:6379

# Auth (Clerk)
BOILERPLATE_AUTH_CLERK_SECRET_KEY=sk_test_xxx

# Observability (New Relic)
BOILERPLATE_OBSERVABILITY_NEWRELIC_LICENSE_KEY=xxx
BOILERPLATE_OBSERVABILITY_NEWRELIC_APP_NAME=my-app
BOILERPLATE_OBSERVABILITY_ENABLED=true

# Email (Resend)
BOILERPLATE_INTEGRATION_RESEND_API_KEY=re_xxx
```

---

## API Routes

| Method | Path | Description |
|--------|------|-------------|
| GET | `/hello` | Health check / greeting endpoint |

---

## Task Commands

```bash
task run              # Run the application
task tidy             # Clean up dependencies
task migrations:new   # Create new migration
task migrations:up    # Run migrations
```

---

## Handler Pattern

This boilerplate uses a typed handler pattern with generics for type-safe request/response handling:

```go
// Define request with validation
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
}

func (r CreateUserRequest) Validate() error {
    return validation.Validate(r)
}

// Handler with typed request/response
func (h *UserHandler) CreateUser(c echo.Context, req CreateUserRequest) (*User, error) {
    // req is already validated
    return h.service.CreateUser(req)
}
```

---

## Error Handling

Structured error responses with proper HTTP status codes:

```go
// Returns 400 Bad Request
errs.NewBadRequestError("INVALID_INPUT", "Invalid email format", "email")

// Returns 404 Not Found  
errs.NewNotFoundError("USER_NOT_FOUND", "User does not exist")

// Returns 401 Unauthorized
errs.NewUnauthorizedError("TOKEN_EXPIRED", "Authentication token has expired")
```

SQL errors are automatically mapped to appropriate HTTP errors:
- Unique violation → 409 Conflict
- Foreign key violation → 400 Bad Request
- Not null violation → 400 Bad Request

---

## License

MIT
