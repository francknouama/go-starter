# Blueprint Usage Guide

Complete guide to using go-starter blueprints for different project types.

## 📊 Implementation Status (v2.0+)

### ✅ Phase 2 Complete - All 12 Blueprints Production Ready

#### Core Web API Blueprints (4/4) ✅
| Blueprint | Status | Loggers | Architecture | Release |
|----------|--------|---------|--------------|---------|
| **Web API Standard** | ✅ Production Ready | slog, zap, logrus, zerolog | Standard layered | v1.0.0+ |
| **Web API Clean Architecture** | ✅ Production Ready | slog, zap, logrus, zerolog | Clean Architecture | v2.0.0+ |
| **Web API DDD** | ✅ Production Ready | slog, zap, logrus, zerolog | Domain-Driven Design | v2.0.0+ |
| **Web API Hexagonal** | ✅ Production Ready | slog, zap, logrus, zerolog | Ports & Adapters | v2.0.0+ |

#### CLI Application Blueprints (2/2) ✅
| Blueprint | Status | Loggers | Files | Complexity |
|----------|--------|---------|--------|------------|
| **CLI Simple** | ✅ Production Ready | slog, zap, logrus, zerolog | 8 files | Beginner |
| **CLI Standard** | ✅ Production Ready | slog, zap, logrus, zerolog | 29 files | Professional |

#### Enterprise & Cloud-Native Blueprints (4/4) ✅
| Blueprint | Status | Loggers | Key Features | Release |
|----------|--------|---------|--------------|---------|
| **gRPC Gateway** | ✅ Production Ready | slog, zap, logrus, zerolog | Dual HTTP/gRPC, TLS | v2.0.0+ |
| **Event-Driven** | ✅ Production Ready | slog, zap, logrus, zerolog | CQRS, Event Sourcing | v2.0.0+ |
| **Microservice** | ✅ Production Ready | slog, zap, logrus, zerolog | Service mesh, K8s | v2.0.0+ |
| **Monolith** | ✅ Production Ready | slog, zap, logrus, zerolog | Full-stack web app | v2.0.0+ |

#### Serverless & Tools Blueprints (2/2) ✅
| Blueprint | Status | Loggers | Runtime | Release |
|----------|--------|---------|---------|---------|
| **AWS Lambda** | ✅ Production Ready | slog, zap, logrus, zerolog | AWS Lambda Go | v1.0.0+ |
| **Lambda Proxy** | ✅ Production Ready | slog, zap, logrus, zerolog | API Gateway proxy | v2.0.0+ |
| **Library** | ✅ Production Ready | slog, zap, loggers, zerolog | Clean API design | v1.0.0+ |
| **Go Workspace** | ✅ Production Ready | slog, zap, logrus, zerolog | Multi-module monorepo | v2.0.0+ |

**🎉 Phase 2 Achievement**: 12/12 blueprints (100% complete) - **All enterprise architecture patterns implemented!**  
**Total Combinations Available**: 48+ (12 blueprints × 4 loggers × architecture variants) - All tested and validated ✅

### 🚧 Phase 3 - Web Interface & Enhanced Features (In Progress)

## Table of Contents

### Core Blueprints ✅
- [Web API Blueprint](#web-api-blueprint) ✅
- [CLI Application Blueprint](#cli-application-blueprint) ✅  
- [Go Library Blueprint](#go-library-blueprint) ✅
- [AWS Lambda Blueprint](#aws-lambda-blueprint) ✅

### Advanced Architecture Blueprints ✅
- [Clean Architecture Web API](#clean-architecture-web-api) ✅
- [DDD Web API](#ddd-web-api) ✅
- [Hexagonal Architecture Web API](#hexagonal-architecture-web-api) ✅

### Enterprise & Cloud-Native ✅
- [gRPC Gateway Blueprint](#grpc-gateway-blueprint) ✅
- [Event-Driven Architecture Blueprint](#event-driven-architecture-blueprint) ✅
- [Microservice Blueprint](#microservice-blueprint) ✅
- [Monolith Blueprint](#monolith-blueprint) ✅

### Serverless & Tools ✅
- [Lambda Proxy Blueprint](#lambda-proxy-blueprint) ✅
- [Go Workspace Blueprint](#go-workspace-blueprint) ✅

### System Features ✅
- [Progressive Disclosure System](#progressive-disclosure-system) ✅
- [Logger Integration](#logger-integration) ✅
- [Best Practices](#best-practices) ✅

---

## Web API Blueprint ✅

**Status**: ✅ Production Ready | **Framework**: Gin | **Architectures**: Standard

### Overview
The Web API blueprint creates a production-ready REST API using the Gin framework with best practices for structure, middleware, and deployment. Currently implements the standard architecture pattern with plans for Clean Architecture, DDD, and Hexagonal patterns in future releases.

### Quick Start
```bash
# Interactive mode
go-starter new my-api --type=web-api

# Direct mode with specific logger
go-starter new my-api --type=web-api --framework=gin --logger=zap
```

### Generated Structure
```
my-api/
├── go.mod                          # Module definition with logger dependencies
├── go.sum                          # Dependency checksums
├── Dockerfile                      # Multi-stage production build
├── Makefile                        # Development and deployment commands
├── README.md                       # Project documentation
├── .gitignore                      # Git ignore patterns
├── .golangci.yml                   # Linting configuration
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   ├── config.go               # Configuration loading
│   │   └── config_test.go          # Configuration tests
│   ├── handlers/
│   │   ├── health.go               # Health check endpoints
│   │   ├── health_test.go          # Handler tests
│   │   └── users.go                # User management endpoints
│   ├── middleware/
│   │   ├── cors.go                 # CORS middleware
│   │   ├── logger.go               # Request logging
│   │   ├── recovery.go             # Panic recovery
│   │   └── auth.go                 # Authentication (if enabled)
│   ├── models/
│   │   ├── user.go                 # Data models
│   │   └── response.go             # API response types
│   ├── services/
│   │   ├── user.go                 # Business logic
│   │   └── user_test.go            # Service tests
│   ├── repository/
│   │   ├── user.go                 # Data access layer
│   │   └── user_test.go            # Repository tests
│   ├── database/
│   │   ├── database.go             # Database connection
│   │   └── migrations/             # Database migrations
│   └── logger/
│       └── factory.go              # Logger factory (selected type)
├── configs/
│   ├── config.yaml                 # Default configuration
│   ├── config.dev.yaml             # Development config
│   └── config.prod.yaml            # Production config
├── docker/
│   └── docker-compose.yml          # Local development with database
├── .github/
│   └── workflows/
│       └── ci.yml                  # GitHub Actions CI/CD
└── tests/
    ├── integration/
    │   └── api_test.go              # Integration tests
    ├── unit/
    │   └── services_test.go         # Unit tests
    └── testdata/
        └── fixtures.json            # Test data
```

### Key Features

#### 🔧 Built-in Middleware
- **CORS**: Cross-origin resource sharing support
- **Logger**: Request/response logging with selected logger
- **Recovery**: Panic recovery with graceful error handling
- **Auth**: JWT authentication (configurable)

#### 📊 Health Checks
```go
GET /health      # Basic health status
GET /health/db   # Database connectivity check
GET /metrics     # Prometheus metrics (if enabled)
```

#### 🗃️ Database Integration
- **GORM ORM**: Object-relational mapping
- **Migrations**: Automatic database schema management
- **Connection pooling**: Production-ready database configuration
- **Multiple drivers**: PostgreSQL, MySQL, SQLite support

#### 🚀 Development Commands
```bash
make run          # Start development server
make test         # Run all tests
make test-unit    # Run unit tests only
make test-integration  # Run integration tests only
make lint         # Run golangci-lint
make build        # Build production binary
make docker       # Build Docker image
make clean        # Clean build artifacts
```

### Configuration

#### Environment Variables
```bash
# Server configuration
PORT=8080
HOST=localhost
ENV=development

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=myuser
DB_PASSWORD=mypassword
DB_NAME=myapi

# Logger configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

#### YAML Configuration
```yaml
# configs/config.yaml
server:
  port: 8080
  host: "localhost"
  env: "development"
  
database:
  host: "localhost"
  port: 5432
  user: "myuser"
  password: "mypassword"
  name: "myapi"
  ssl_mode: "disable"
  
logger:
  level: "info"
  format: "json"
  output: "stdout"
```

### Deployment

#### Docker Deployment
```bash
# Build and run locally
make docker
docker run -p 8080:8080 my-api:latest

# With environment variables
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=secret \
  my-api:latest
```

#### Production Deployment
```bash
# Build for production
make build

# Run binary
./bin/server --config=configs/config.prod.yaml
```

---

## CLI Application Blueprint ✅

**Status**: ✅ Production Ready | **Framework**: Cobra | **Architectures**: Standard

### Overview
Creates a powerful command-line application using the Cobra framework with subcommands, configuration management, and professional CLI patterns.

### Quick Start
```bash
# Interactive mode
go-starter new my-tool --type=cli

# Direct mode with specific logger
go-starter new my-tool --type=cli --framework=cobra --logger=logrus
```

### Generated Structure
```
my-tool/
├── go.mod                          # Module definition
├── main.go                         # Application entry point
├── Dockerfile                      # Container support
├── Makefile                        # Development commands
├── README.md                       # Documentation
├── .gitignore                      # Git ignore patterns
├── .golangci.yml                   # Linting configuration
├── cmd/
│   ├── root.go                     # Root command definition  
│   ├── root_test.go                # Command tests
│   ├── version.go                  # Version subcommand
│   └── completion.go               # Shell completion
├── internal/
│   ├── config/
│   │   ├── config.go               # Configuration management
│   │   └── config_test.go          # Config tests
│   └── logger/
│       └── factory.go              # Logger factory
├── configs/
│   └── config.yaml                 # Default configuration
└── .github/
    └── workflows/
        └── ci.yml                  # CI/CD pipeline
```

### Key Features

#### 🎛️ Command Structure
```bash
my-tool --help                      # Show help
my-tool version                     # Show version
my-tool completion bash             # Generate bash completion
my-tool config validate             # Validate configuration
```

#### ⚙️ Configuration Management
- **Viper integration**: Configuration from files, environment, flags
- **Multiple formats**: YAML, JSON, TOML support
- **Environment variable binding**: Automatic env var mapping
- **Flag precedence**: Command flags override config files

#### 🔧 Built-in Commands
```go
// Root command with global flags
var rootCmd = &cobra.Command{
    Use:   "my-tool",
    Short: "A powerful CLI tool",
    Long:  `A comprehensive CLI application with best practices.`,
}

// Version command
var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print version information",
    Run:   func(cmd *cobra.Command, args []string) { ... },
}
```

### Adding Custom Commands

#### Create New Command
```go
// cmd/deploy.go
package cmd

import (
    "github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Deploy application to target environment",
    RunE:  runDeploy,
}

func init() {
    rootCmd.AddCommand(deployCmd)
    deployCmd.Flags().StringP("env", "e", "dev", "Target environment")
    deployCmd.Flags().BoolP("dry-run", "d", false, "Perform dry run")
}

func runDeploy(cmd *cobra.Command, args []string) error {
    logger := logger.New()
    
    env, _ := cmd.Flags().GetString("env")
    dryRun, _ := cmd.Flags().GetBool("dry-run")
    
    logger.Info("Starting deployment", "env", env, "dry-run", dryRun)
    
    // Your deployment logic here
    
    return nil
}
```

### Development Commands
```bash
make build        # Build binary
make install      # Install globally
make test         # Run tests
make lint         # Run linting
make clean        # Clean build artifacts
make release      # Build release binaries
```

---

## Go Library Blueprint ✅

**Status**: ✅ Production Ready | **Type**: Library Package | **Architectures**: Standard

### Overview
Creates a well-structured Go library with comprehensive documentation, examples, and testing setup suitable for open source distribution.

### Quick Start
```bash
# Interactive mode
go-starter new awesome-lib --type=library

# Direct mode
go-starter new awesome-lib --type=library --logger=slog
```

### Generated Structure
```
awesome-lib/
├── go.mod                          # Module definition
├── awesome_lib.go                  # Main library interface
├── types.go                        # Public types and constants
├── errors.go                       # Error definitions
├── options.go                      # Configuration options
├── Makefile                        # Development commands
├── README.md                       # Documentation
├── LICENSE                         # License file
├── .gitignore                      # Git ignore patterns
├── .golangci.yml                   # Linting configuration
├── internal/
│   ├── client/
│   │   ├── client.go               # Internal implementation
│   │   └── client_test.go          # Internal tests
│   └── logger/
│       └── factory.go              # Minimal logger interface
├── examples/
│   ├── basic/
│   │   └── main.go                 # Basic usage example
│   └── advanced/
│       └── main.go                 # Advanced usage example
├── docs/
│   ├── API.md                      # API documentation
│   └── EXAMPLES.md                 # Usage examples
└── .github/
    └── workflows/
        └── ci.yml                  # CI/CD pipeline
```

### Key Features

#### 📦 Clean Public API
```go
// Public interface
type Client interface {
    Connect(ctx context.Context, opts ...Option) error
    Process(ctx context.Context, data []byte) (*Result, error)
    Close() error
}

// Constructor function
func New(opts ...Option) Client {
    return &client{
        logger: logger.New(),
        config: defaultConfig(),
    }
}

// Functional options pattern
func WithTimeout(timeout time.Duration) Option {
    return func(c *client) {
        c.timeout = timeout
    }
}
```

#### 📚 Documentation
- **GoDoc comments**: Comprehensive API documentation
- **Examples**: Runnable examples for all public functions
- **README**: Installation, usage, and contribution guidelines
- **API reference**: Detailed API documentation

#### 🧪 Testing Setup
```go
// Comprehensive test coverage
func TestClient_Connect(t *testing.T) {
    tests := []struct {
        name    string
        opts    []Option
        wantErr bool
    }{
        {"default options", nil, false},
        {"with timeout", []Option{WithTimeout(5 * time.Second)}, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := New(tt.opts...)
            err := client.Connect(context.Background())
            if (err != nil) != tt.wantErr {
                t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

// Benchmark tests
func BenchmarkClient_Process(b *testing.B) {
    client := New()
    data := make([]byte, 1024)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = client.Process(context.Background(), data)
    }
}
```

### Development Commands
```bash
make test         # Run tests with coverage
make benchmark    # Run benchmark tests
make lint         # Run linting
make docs         # Generate documentation
make examples     # Run all examples
make release      # Prepare for release
```

---

## AWS Lambda Blueprint ✅

**Status**: ✅ Production Ready | **Runtime**: AWS Lambda Go | **Architectures**: Standard

### Overview
Creates an AWS Lambda function optimized for serverless deployment with API Gateway integration and CloudWatch logging.

### Quick Start
```bash
# Interactive mode
go-starter new my-function --type=lambda

# Direct mode with CloudWatch-optimized logger
go-starter new my-function --type=lambda --logger=zerolog
```

### Generated Structure
```
my-function/
├── go.mod                          # Module definition
├── main.go                         # Lambda entry point
├── template.yaml                   # SAM deployment template
├── Makefile                        # Build and deployment commands
├── README.md                       # Documentation
├── .gitignore                      # Git ignore patterns
├── internal/
│   ├── handler/
│   │   ├── handler.go              # Lambda handler logic
│   │   └── handler_test.go         # Handler tests
│   ├── logger/
│   │   └── factory.go              # CloudWatch-optimized logger
│   └── response/
│       ├── response.go             # API Gateway response helpers
│       └── response_test.go        # Response tests
├── events/
│   └── api-gateway.json            # Test event data
└── .github/
    └── workflows/
        └── deploy.yml               # Deployment pipeline
```

### Key Features

#### ⚡ Lambda Handler
```go
func main() {
    logger := logger.New()
    handler := &Handler{logger: logger}
    
    lambda.Start(handler.HandleRequest)
}

func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    h.logger.Info("Processing request", 
        "method", request.HTTPMethod,
        "path", request.Path,
        "requestId", request.RequestContext.RequestID,
    )
    
    // Your business logic here
    
    return response.Success(map[string]interface{}{
        "message": "Hello from Lambda!",
        "timestamp": time.Now().Unix(),
    }), nil
}
```

#### 🌐 API Gateway Integration
- **Request handling**: HTTP request parsing and validation
- **Response formatting**: Proper API Gateway response structure
- **Error handling**: Structured error responses
- **CORS support**: Cross-origin resource sharing

#### 📊 CloudWatch Optimized Logging
```go
// Structured logging optimized for CloudWatch
logger.Info("Function started",
    "function", "my-function",
    "version", "$LATEST",
    "requestId", ctx.Value("requestId"),
)

logger.Error("Database error",
    "error", err.Error(),
    "operation", "user.create",
    "duration", time.Since(start).Milliseconds(),
)
```

### SAM Template
```yaml
# template.yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Parameters:
  Environment:
    Type: String
    Default: dev
    AllowedValues: [dev, staging, prod]

Globals:
  Function:
    Timeout: 30
    Runtime: go1.x
    Environment:
      Variables:
        LOG_LEVEL: info
        ENVIRONMENT: !Ref Environment

Resources:
  MyFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: .
      Handler: bootstrap
      Events:
        ApiEvent:
          Type: Api
          Properties:
            Path: /{proxy+}
            Method: any
```

### Development Commands
```bash
make build-lambda    # Cross-compile for Linux
make test-local      # Test locally with SAM CLI  
make deploy-dev      # Deploy to development
make deploy-prod     # Deploy to production
make logs            # View CloudWatch logs
make invoke          # Test function locally
make clean           # Clean build artifacts
```

### Deployment
```bash
# Development deployment
make deploy-dev

# Production deployment
make deploy-prod ENVIRONMENT=prod

# View logs
make logs ENVIRONMENT=prod

# Test function
make invoke EVENT=events/api-gateway.json
```

---

## Logger Integration

### Overview
All blueprints include sophisticated logger integration with consistent interfaces across four popular Go logging libraries.

### Supported Loggers

#### slog (Default)
```go
// Standard library structured logging
logger.Info("Server started", "port", 8080, "env", "production")
logger.Error("Database error", "error", err, "table", "users")
```

#### Zap (High Performance)
```go
// Ultra-fast logging with zero allocations
logger.Info("Server started", zap.Int("port", 8080), zap.String("env", "production"))
logger.Error("Database error", zap.Error(err), zap.String("table", "users"))
```

#### Logrus (Feature Rich)
```go
// Popular structured logger with rich features
logger.WithFields(logrus.Fields{
    "port": 8080,
    "env": "production",
}).Info("Server started")
```

#### Zerolog (Zero Allocation)
```go
// Zero allocation JSON logger
logger.Info().Int("port", 8080).Str("env", "production").Msg("Server started")
logger.Error().Err(err).Str("table", "users").Msg("Database error")
```

### Logger Factory Pattern
Each blueprint includes a logger factory that provides a consistent interface:

```go
// internal/logger/factory.go
package logger

import "context"

type Logger interface {
    Debug(message string, fields ...interface{})
    Info(message string, fields ...interface{})
    Warn(message string, fields ...interface{})
    Error(message string, fields ...interface{})
    With(fields ...interface{}) Logger
    WithContext(ctx context.Context) Logger
}

func New() Logger {
    // Returns the selected logger implementation
    // (slog, zap, logrus, or zerolog)
}
```

### Configuration
```yaml
# Logger configuration in config.yaml
logger:
  level: "info"          # debug, info, warn, error
  format: "json"         # json, text, console
  output: "stdout"       # stdout, stderr, file path
  caller: true           # Include caller information
  timestamp: true        # Include timestamps
```

---

## Best Practices

### Project Structure
- **cmd/**: Application entry points
- **internal/**: Private application code
- **pkg/**: Public library code (for libraries)
- **configs/**: Configuration files
- **docs/**: Documentation
- **tests/**: Test files and test data

### Error Handling
```go
// Use structured errors with context
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   interface{} `json:"value,omitempty"`
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}

// Log errors with context
logger.Error("Validation failed",
    "error", err,
    "field", "email",
    "value", userInput.Email,
    "requestId", ctx.Value("requestId"),
)
```

### Testing
```go
// Table-driven tests
func TestUserService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateUserRequest
        want    *User
        wantErr bool
    }{
        {
            name: "valid user",
            input: CreateUserRequest{
                Name:  "John Doe",
                Email: "john@example.com",
            },
            want: &User{
                Name:  "John Doe", 
                Email: "john@example.com",
            },
            wantErr: false,
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := NewUserService(mockRepo, logger.New())
            got, err := service.Create(context.Background(), tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("Create() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Configuration Management
```go
// Use environment-specific configs
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Logger   LoggerConfig   `yaml:"logger"`
}

// Load with environment override
func LoadConfig() (*Config, error) {
    config := &Config{}
    
    // Load base config
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    // Environment-specific overrides
    viper.SetEnvPrefix("APP")
    viper.AutomaticEnv()
    
    if err := viper.Unmarshal(config); err != nil {
        return nil, err
    }
    
    return config, nil
}
```

### Security
```go
// Input validation
func validateCreateUserRequest(req CreateUserRequest) error {
    var errors []ValidationError
    
    if req.Email == "" {
        errors = append(errors, ValidationError{
            Field:   "email",
            Message: "email is required",
        })
    } else if !isValidEmail(req.Email) {
        errors = append(errors, ValidationError{
            Field:   "email", 
            Message: "invalid email format",
            Value:   req.Email,
        })
    }
    
    if len(errors) > 0 {
        return MultiValidationError{Errors: errors}
    }
    
    return nil
}

// Sanitize inputs
func sanitizeUserInput(input string) string {
    // Remove dangerous characters
    cleaned := strings.TrimSpace(input)
    cleaned = html.EscapeString(cleaned)
    return cleaned
}
```

### Performance
```go
// Use context for timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := service.Process(ctx, data)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        logger.Warn("Operation timed out", "operation", "process", "timeout", "5s")
    }
    return err
}

// Connection pooling for databases
db, err := sql.Open("postgres", dsn)
if err != nil {
    return err
}

db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(time.Hour)
```

This completes the comprehensive blueprint usage guide. Each blueprint is designed with production-ready practices and can be customized for specific needs while maintaining consistency across all logger implementations.

## Architecture-Specific Limitations

### Clean Architecture Web API

**Important**: When using the Clean Architecture pattern, authentication features require a database to be configured. This is because:

- Authentication use cases depend on user entities and repositories
- User entities are only generated when a database driver is selected
- The authentication system needs persistent storage for users and sessions

**Valid Configurations**:
- ✅ Clean Architecture + Database + Authentication
- ✅ Clean Architecture + Database (no auth)
- ✅ Clean Architecture (no database, no auth)
- ❌ Clean Architecture + Authentication (no database) - Will not compile

This design ensures proper separation of concerns and follows Clean Architecture principles where business logic depends on data persistence abstractions.