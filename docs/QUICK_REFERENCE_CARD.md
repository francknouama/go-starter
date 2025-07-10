# go-starter Quick Reference

## Installation

```bash
# Homebrew (macOS/Linux)
brew tap francknouama/tap
brew install go-starter

# Go install
go install github.com/francknouama/go-starter@latest

# Verify installation
go-starter version
```

## Project Generation

### Interactive Mode (Recommended for beginners)
```bash
go-starter new my-project
```

### Direct Mode Commands

#### Web API
```bash
# Basic web API with default logger (slog)
go-starter new my-api --type=web-api

# High-performance API with zap logger
go-starter new fast-api --type=web-api --logger=zap

# API with database support
go-starter new shop-api --type=web-api --database=postgres,redis

# API with specific module path
go-starter new my-api --type=web-api --module=github.com/company/my-api
```

#### CLI Application
```bash
# Basic CLI with default logger
go-starter new my-cli --type=cli

# CLI with rich output using logrus
go-starter new devtool --type=cli --logger=logrus

# CLI with custom module path
go-starter new tool --type=cli --module=github.com/myorg/tool
```

#### Go Library
```bash
# Basic library (minimal logger interface)
go-starter new my-lib --type=library

# Library with specific logger
go-starter new sdk --type=library --logger=slog
```

#### AWS Lambda
```bash
# Basic Lambda function
go-starter new my-lambda --type=lambda

# Lambda with CloudWatch-optimized logging
go-starter new processor --type=lambda --logger=zerolog
```

### Common Flags
```bash
--type        # Project type: web-api, cli, library, lambda
--logger      # Logger type: slog, zap, logrus, zerolog
--database    # Database(s): postgres, mysql, mongodb, redis, sqlite
--module      # Go module path
--force       # Overwrite existing directory
```

## Blueprint & Logger Quick Decision

### By Use Case

| Use Case | Blueprint | Logger | Command |
|----------|----------|--------|---------|
| REST API Service | web-api | zap | `go-starter new api --type=web-api --logger=zap` |
| Developer Tool | cli | logrus | `go-starter new tool --type=cli --logger=logrus` |
| Shared Package | library | slog | `go-starter new lib --type=library --logger=slog` |
| Serverless Function | lambda | zerolog | `go-starter new func --type=lambda --logger=zerolog` |
| Standard Web App | web-api | slog | `go-starter new app --type=web-api` |

### Logger Comparison

| Logger | Best For | Performance | Key Feature |
|--------|----------|-------------|-------------|
| **slog** | Default choice | Good | Built-in to Go |
| **zap** | High-traffic APIs | Excellent | Zero allocation |
| **logrus** | CLI tools | Moderate | Rich formatting |
| **zerolog** | Cloud-native | Excellent | Clean JSON |

## Generated Project Commands

### All Blueprints Include Makefile

```bash
make help         # Show all available commands
make run          # Run the application
make test         # Run tests with coverage
make lint         # Run golangci-lint
make build        # Build production binary
make clean        # Clean build artifacts
```

### Web API Specific
```bash
make docker       # Build Docker image
make migrate-up   # Run database migrations
make migrate-down # Rollback migrations
make swagger      # Generate API documentation
```

### CLI Specific
```bash
make install      # Install CLI globally
make release      # Build for all platforms
```

### Lambda Specific
```bash
make build-lambda # Build for AWS Lambda
make deploy-dev   # Deploy to development
make deploy-prod  # Deploy to production
make logs         # View CloudWatch logs
```

## Configuration Examples

### Web API Configuration
```yaml
# configs/config.yaml
server:
  port: 8080
  timeout: 30s

database:
  driver: postgres
  host: localhost
  port: 5432
  name: myapp

logger:
  level: info
  format: json
```

### Environment Variables
```bash
# Server
PORT=8080
HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=myapp

# Logger
LOG_LEVEL=debug
LOG_FORMAT=json
```

## Logger Usage Examples

### slog (Default)
```go
logger.Info("Server started", 
    slog.String("port", "8080"),
    slog.String("env", "production"))

logger.Error("Operation failed",
    slog.String("error", err.Error()),
    slog.Int("retries", 3))
```

### zap (High Performance)
```go
logger.Info("Request processed",
    zap.String("method", "GET"),
    zap.String("path", "/api/users"),
    zap.Duration("latency", time.Since(start)))

logger.Error("Database error",
    zap.Error(err),
    zap.String("query", query))
```

### logrus (Feature Rich)
```go
logger.WithFields(logrus.Fields{
    "user_id": userID,
    "action": "login",
    "ip": clientIP,
}).Info("User authenticated")

logger.WithError(err).Error("Payment processing failed")
```

### zerolog (Zero Allocation)
```go
logger.Info().
    Str("service", "api").
    Str("version", "1.0.0").
    Msg("Service initialized")

logger.Error().
    Err(err).
    Str("user_id", userID).
    Str("operation", "create_order").
    Msg("Operation failed")
```

## Common Development Workflows

### Starting a New API Project
```bash
# 1. Generate project
go-starter new awesome-api --type=web-api --logger=zap

# 2. Navigate to project
cd awesome-api

# 3. Install dependencies
go mod tidy

# 4. Set up database (if using Docker)
docker-compose up -d

# 5. Run migrations
make migrate-up

# 6. Start development server
make run
```

### Building for Production
```bash
# 1. Run tests
make test

# 2. Run linter
make lint

# 3. Build binary
make build

# 4. Build Docker image
make docker

# 5. Run production binary
./bin/server --config=configs/config.prod.yaml
```

### Working with Databases

#### PostgreSQL Connection
```go
dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)
```

#### MySQL Connection
```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
    config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
```

#### Multi-Database Setup
```bash
# Generate with multiple databases
go-starter new api --type=web-api --database=postgres,redis,mongodb

# Docker Compose will include all selected databases
docker-compose up -d
```

## Testing Patterns

### Unit Test Structure
```go
func TestUserService_Create(t *testing.T) {
    // Arrange
    mockRepo := mocks.NewMockUserRepository(t)
    service := NewUserService(mockRepo, logger.New())
    
    mockRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil)
    
    // Act
    user, err := service.Create(context.Background(), CreateUserRequest{
        Email: "test@example.com",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email)
}
```

### Integration Test
```go
func TestAPI_CreateUser_Integration(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer cleanupDB(db)
    
    router := setupRouter(db)
    
    // Test
    w := httptest.NewRecorder()
    body := `{"email":"test@example.com","username":"testuser"}`
    req := httptest.NewRequest("POST", "/api/users", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    
    // Verify
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response UserResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", response.Email)
}
```

## Deployment

### Docker Deployment
```bash
# Build image
docker build -t myapp:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PASSWORD=secret \
  --name myapp \
  myapp:latest

# View logs
docker logs -f myapp
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: postgres-service
        - name: LOG_LEVEL
          value: info
```

### AWS Lambda Deployment
```bash
# Build for Lambda
make build-lambda

# Deploy with SAM
sam deploy --guided

# Or using Makefile
make deploy-prod
```

## Troubleshooting Quick Fixes

### Installation Issues
```bash
# Can't find go-starter after install
export PATH=$PATH:$(go env GOPATH)/bin

# Permission denied
sudo go install github.com/francknouama/go-starter@latest
```

### Generation Issues
```bash
# Blueprint not found
go-starter list  # Check available blueprints

# Invalid module path
# Use proper format: github.com/username/project
```

### Build Issues
```bash
# Missing dependencies
go mod tidy
go mod download

# Clear module cache
go clean -modcache
```

### Runtime Issues
```bash
# Port in use
lsof -ti:8080 | xargs kill -9

# Database connection failed
docker-compose up -d  # Start database
```

## Useful Aliases

Add to your shell configuration:

```bash
# ~/.bashrc or ~/.zshrc

# go-starter aliases
alias gsn='go-starter new'
alias gsnapi='go-starter new --type=web-api'
alias gsncli='go-starter new --type=cli'
alias gsnlib='go-starter new --type=library'
alias gsnlambda='go-starter new --type=lambda'

# Quick project creation
function new-api() {
    go-starter new "$1" --type=web-api --logger=zap --module=github.com/$(git config user.name)/"$1"
}

function new-cli() {
    go-starter new "$1" --type=cli --logger=logrus --module=github.com/$(git config user.name)/"$1"
}
```

## Links & Resources

- **Repository**: [github.com/francknouama/go-starter](https://github.com/francknouama/go-starter)
- **Issues**: [Report bugs or request features](https://github.com/francknouama/go-starter/issues)
- **Documentation**: [Full documentation](https://github.com/francknouama/go-starter/tree/main/docs)

## Version Information

```bash
# Check version
go-starter version

# Update to latest
go install github.com/francknouama/go-starter@latest
```

---

**Pro Tip**: Start with interactive mode (`go-starter new project-name`) to learn available options, then use direct mode with flags for faster project creation once you know what you need.