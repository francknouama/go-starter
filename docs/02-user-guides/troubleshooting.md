# Troubleshooting Guide

Comprehensive solutions for common issues with go-starter. Get your projects working quickly with these proven fixes.

## Table of Contents

- [Quick Diagnostics](#quick-diagnostics)
- [Installation Issues](#installation-issues)
- [Project Generation Problems](#project-generation-problems)
- [Build and Compilation Errors](#build-and-compilation-errors)
- [Runtime Issues](#runtime-issues)
- [Logger-Specific Problems](#logger-specific-problems)
- [Database Connection Issues](#database-connection-issues)
- [Docker and Deployment Issues](#docker-and-deployment-issues)
- [Platform-Specific Issues](#platform-specific-issues)
- [Getting Help](#getting-help)

## Quick Diagnostics

### Run System Check

Save this as `diagnose.sh` and run to check your environment:

```bash
#!/bin/bash
echo "=== go-starter Diagnostic Report ==="
echo "Date: $(date)"
echo ""

echo "1. go-starter Installation:"
if command -v go-starter >/dev/null 2>&1; then
    go-starter version
    echo "✅ go-starter installed"
else
    echo "❌ go-starter NOT installed"
fi
echo ""

echo "2. Go Environment:"
if command -v go >/dev/null 2>&1; then
    go version
    echo "GOPATH: $(go env GOPATH)"
    echo "GOBIN: $(go env GOBIN)"
    echo "✅ Go installed"
else
    echo "❌ Go NOT installed"
fi
echo ""

echo "3. PATH Check:"
echo "PATH: $PATH"
if echo $PATH | grep -q "$(go env GOPATH)/bin"; then
    echo "✅ GOPATH/bin in PATH"
else
    echo "⚠️ GOPATH/bin NOT in PATH"
fi
echo ""

echo "4. Network Connectivity:"
if curl -s -o /dev/null -w "%{http_code}" https://proxy.golang.org 2>/dev/null | grep -q "200"; then
    echo "✅ Go proxy accessible"
else
    echo "❌ Go proxy NOT accessible"
fi
echo ""

echo "5. System Resources:"
echo "Disk space: $(df -h . | tail -1 | awk '{print $4}') available"
echo "Memory: $(free -h 2>/dev/null | grep Mem | awk '{print $7}' || echo 'N/A') available"
echo ""

echo "=== End Diagnostic Report ==="
```

Run with: `bash diagnose.sh`

## Installation Issues

### Command Not Found

**Problem**: `go-starter: command not found`

**Solution**:

```bash
# Check if Go is installed
go version

# If Go is missing, install it first
# Then install go-starter
go install github.com/francknouama/go-starter@latest

# Check GOPATH/bin is in PATH
echo $PATH | grep "$(go env GOPATH)/bin"

# If not in PATH, add it
export PATH="$PATH:$(go env GOPATH)/bin"

# Make permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
```

### Installation Succeeds but Binary Not Found

**Problem**: `go install` completes but `go-starter` not found

**Diagnostic**:
```bash
# Find where Go installs binaries
go env GOBIN
go env GOPATH

# Look for the binary
find $(go env GOPATH) -name "go-starter" -type f 2>/dev/null
which go-starter
```

**Solutions**:

1. **Set GOBIN explicitly**:
```bash
export GOBIN=$HOME/bin
mkdir -p $GOBIN
go install github.com/francknouama/go-starter@latest
export PATH=$GOBIN:$PATH
```

2. **Install to system path**:
```bash
sudo GOBIN=/usr/local/bin go install github.com/francknouama/go-starter@latest
```

3. **Manual installation**:
```bash
# Download and install manually
curl -L -o go-starter https://github.com/francknouama/go-starter/releases/latest/download/go-starter-$(uname -s | tr '[:upper:]' '[:lower:]')-amd64
chmod +x go-starter
sudo mv go-starter /usr/local/bin/
```

### Version Conflicts

**Problem**: Old version persists after update

**Solution**:
```bash
# Find all installations
which -a go-starter
type -a go-starter

# Remove old versions
sudo rm $(which go-starter)

# Clean module cache
go clean -modcache

# Reinstall latest
go install github.com/francknouama/go-starter@latest

# Verify new version
go-starter version
```

### Network/Proxy Issues

**Problem**: Installation fails with network errors

**Solution**:
```bash
# Check Go proxy settings
go env GOPROXY

# Configure proxy if behind corporate firewall
go env -w GOPROXY=https://proxy.golang.org,direct
go env -w GOSUMDB=sum.golang.org

# Or use direct mode (bypass proxy)
go env -w GOPROXY=direct

# For corporate environments
export GOPROXY=https://your-proxy.company.com
export GONOPROXY=github.com/francknouama/*
```

## Project Generation Problems

### Template Variables Not Replaced

**Problem**: Generated files contain `{{.ProjectName}}` literals

**Diagnostic**:
```bash
# Check for template literals in generated files
grep -r "{{" my-project/

# Generate with verbose output (if available)
go-starter new test-project --type=web-api --dry-run
```

**Solutions**:

1. **Check module path format**:
```bash
# ✅ Correct
go-starter new myapp --module=github.com/user/myapp

# ❌ Incorrect - invalid characters
go-starter new myapp --module="My App"
go-starter new myapp --module="github.com/user/my app"
```

2. **Verify project name**:
```bash
# ✅ Correct
go-starter new my-awesome-project

# ❌ Incorrect - invalid characters
go-starter new "my project"
go-starter new my_project!
```

3. **Use quotes for complex names**:
```bash
# If you must use spaces or special characters
go-starter new "my-project" --module=github.com/user/my-project
```

### Partial Generation Failure

**Problem**: Some files created, others missing

**Diagnostic**:
```bash
# Check disk space
df -h .

# Check file permissions
touch test-file && rm test-file || echo "No write permission"

# Monitor generation progress
watch -n 0.5 'find my-project -type f 2>/dev/null | wc -l'
```

**Recovery Steps**:
```bash
# 1. Clean partial generation
rm -rf my-project

# 2. Ensure adequate resources
# Need at least 100MB free space

# 3. Try in a different directory
cd /tmp
go-starter new test-project --type=web-api

# 4. If successful, move to desired location
mv test-project ~/projects/
```

### Interactive Mode Stuck

**Problem**: Interactive prompts don't respond

**Solution**:
```bash
# Use direct mode instead
go-starter new my-project \
  --type=web-api \
  --framework=gin \
  --logger=slog

# Or force non-interactive
go-starter new my-project --type=web-api --no-interactive

# Check terminal compatibility
echo $TERM
```

## Build and Compilation Errors

### Import Cycle Detected

**Problem**: 
```
import cycle not allowed
package github.com/user/project/internal/config
	imports github.com/user/project/internal/logger
	imports github.com/user/project/internal/config
```

**Diagnostic**:
```bash
# Visualize the dependency graph
go mod graph | grep -E "(config|logger)"

# Find circular dependencies
go list -f '{{.ImportPath}} {{.Imports}}' ./...
```

**Solutions**:

1. **Interface segregation**:
```go
// Create logger/interface.go
package logger

type Config interface {
    GetLogLevel() string
    GetLogFormat() string
}

// Remove direct config import from logger
```

2. **Separate configuration package**:
```go
// Create pkg/logconfig/config.go
package logconfig

type Config struct {
    Level  string `yaml:"level"`
    Format string `yaml:"format"`
}
```

3. **Dependency injection**:
```go
// Pass config to logger instead of importing
func NewLogger(level, format string) Logger {
    // Implementation
}
```

### Missing Dependencies

**Problem**:
```
cannot find module providing package go.uber.org/zap
cannot find module providing package github.com/gin-gonic/gin
```

**Solution**:
```bash
# 1. Verify go.mod exists
ls -la go.mod

# 2. Clean and reinitialize
rm go.mod go.sum
go mod init github.com/user/project

# 3. Add dependencies
go mod tidy

# 4. Download all modules
go mod download

# 5. Verify integrity
go mod verify

# 6. Test build
go build -v ./...
```

### CGO Compilation Issues

**Problem**: C compiler errors with certain dependencies

**Solution**:
```bash
# Disable CGO if not needed
CGO_ENABLED=0 go build ./...

# Or install build tools
# Ubuntu/Debian
sudo apt-get install build-essential

# CentOS/RHEL
sudo yum groupinstall "Development Tools"

# macOS
xcode-select --install

# Alpine Linux (Docker)
apk add --no-cache gcc musl-dev
```

## Runtime Issues

### Application Panics on Startup

**Problem**: Application crashes immediately with panic

**Diagnostic**:
```bash
# Run with race detector
go run -race cmd/server/main.go

# Enable debug logging
LOG_LEVEL=debug go run cmd/server/main.go

# Use Go debugger
dlv debug cmd/server/main.go
```

**Common Panic Causes**:

1. **Nil pointer dereference**:
```go
// Add nil checks
if config == nil {
    log.Fatal("Config is nil")
}
if logger == nil {
    logger = slog.Default()
}
```

2. **Missing configuration**:
```go
// Load config with error handling
config, err := LoadConfig()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

3. **Database connection failure**:
```go
// Add connection validation
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatalf("Database connection failed: %v", err)
}

// Test connection
sqlDB, err := db.DB()
if err != nil {
    log.Fatalf("Failed to get DB instance: %v", err)
}

if err := sqlDB.Ping(); err != nil {
    log.Fatalf("Database ping failed: %v", err)
}
```

### Port Already in Use

**Problem**: 
```
listen tcp :8080: bind: address already in use
```

**Solution**:
```bash
# Find process using port
lsof -ti:8080
netstat -tulpn | grep :8080
ss -tulpn | grep :8080

# Kill the process
kill -9 $(lsof -ti:8080)

# Or use different port
PORT=8081 go run cmd/server/main.go

# Configure in environment
echo "PORT=8081" > .env
```

### Configuration Not Loading

**Problem**: Application starts but uses default values

**Diagnostic**:
```bash
# Check config file location
ls -la configs/
ls -la config.yaml .env

# Verify file permissions
ls -la configs/config.yaml

# Test config loading
go run cmd/server/main.go --config configs/config.yaml
```

**Solution**:
```go
// Add config loading debug
func LoadConfig() (*Config, error) {
    configPaths := []string{
        "configs/config.yaml",
        "config.yaml",
        "./config.yaml",
    }
    
    for _, path := range configPaths {
        if _, err := os.Stat(path); err == nil {
            log.Printf("Loading config from: %s", path)
            return loadConfigFromFile(path)
        }
    }
    
    log.Println("No config file found, using defaults")
    return DefaultConfig(), nil
}
```

## Logger-Specific Problems

### Zap Logger Performance Issues

**Problem**: High memory usage or slow performance with zap

**Solution**:
```go
// Use production config
func NewZapLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    
    // Disable stack traces for performance
    config.DisableStacktrace = true
    
    // Enable sampling to reduce log volume
    config.Sampling = &zap.SamplingConfig{
        Initial:    100,
        Thereafter: 100,
    }
    
    logger, _ := config.Build()
    return logger
}

// Pre-allocate common fields
logger = logger.With(
    zap.String("service", "api"),
    zap.String("version", "1.0.0"),
)

// Avoid string formatting in logs
// ❌ Bad
logger.Info(fmt.Sprintf("User %s logged in", username))

// ✅ Good
logger.Info("User logged in", zap.String("user", username))
```

### Logrus Colors Not Working

**Problem**: Terminal output is monochrome

**Solution**:
```go
import "github.com/mattn/go-isatty"

// Force colors or detect terminal
if isatty.IsTerminal(os.Stdout.Fd()) || os.Getenv("FORCE_COLOR") == "true" {
    logrus.SetFormatter(&logrus.TextFormatter{
        ForceColors:   true,
        FullTimestamp: true,
    })
} else {
    logrus.SetFormatter(&logrus.JSONFormatter{})
}
```

### Zerolog Not Producing JSON

**Problem**: Console output instead of JSON in production

**Solution**:
```go
// Conditional output format
var logger zerolog.Logger

if os.Getenv("ENV") == "development" {
    // Human-readable in development
    logger = zerolog.New(zerolog.ConsoleWriter{
        Out:        os.Stderr,
        TimeFormat: time.RFC3339,
    }).With().Timestamp().Logger()
} else {
    // JSON in production
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
```

### slog Configuration Issues

**Problem**: slog not showing expected log levels

**Solution**:
```go
// Configure slog properly
var logger *slog.Logger

level := os.Getenv("LOG_LEVEL")
switch strings.ToLower(level) {
case "debug":
    logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
case "error":
    logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelError,
    }))
default:
    logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
}
```

## Database Connection Issues

### PostgreSQL Connection Refused

**Problem**: 
```
dial tcp [::1]:5432: connect: connection refused
```

**Diagnostic**:
```bash
# Check if PostgreSQL is running
ps aux | grep postgres
systemctl status postgresql

# Test connectivity
psql -h localhost -U postgres -c "SELECT 1"
telnet localhost 5432
nc -zv localhost 5432

# Check connection string
echo $DATABASE_URL
```

**Solutions**:

1. **Start PostgreSQL**:
```bash
# Ubuntu/Debian
sudo systemctl start postgresql

# macOS with Homebrew
brew services start postgresql

# Docker
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  postgres:15
```

2. **Fix connection string**:
```bash
# Correct format
DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"

# With SSL
DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=require"
```

3. **Add connection retry logic**:
```go
func ConnectWithRetry(dsn string, maxRetries int) (*gorm.DB, error) {
    var db *gorm.DB
    var err error
    
    for i := 0; i < maxRetries; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            sqlDB, _ := db.DB()
            if err := sqlDB.Ping(); err == nil {
                log.Printf("Database connected on attempt %d", i+1)
                return db, nil
            }
        }
        
        log.Printf("Database connection attempt %d failed: %v", i+1, err)
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
}
```

### Migration Failures

**Problem**: Database migrations fail or are partially applied

**Diagnostic**:
```bash
# Check migration status (if available)
go run main.go migrate status

# Check database schema
psql -d mydb -c "\dt"  # List tables
psql -d mydb -c "\d schema_migrations"  # Check migration table
```

**Solutions**:

1. **Manual migration check**:
```bash
# Check what migrations exist
ls -la migrations/

# Apply specific migration
psql -d mydb -f migrations/001_initial.up.sql
```

2. **Reset migrations (development only)**:
```bash
# Drop all tables
psql -d mydb -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

# Rerun migrations
go run main.go migrate up
```

3. **Fix migration code**:
```go
// Add transaction support to migrations
func runMigrations(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if err := tx.AutoMigrate(&User{}, &Product{}); err != nil {
            return err
        }
        
        // Add any additional migration logic
        return nil
    })
}
```

## Docker and Deployment Issues

### Docker Build Failures

**Problem**: Multi-stage builds fail or produce large images

**Diagnostic Dockerfile**:
```dockerfile
FROM golang:1.22-alpine AS builder

# Add debugging
RUN echo "Starting build..." && \
    apk add --no-cache git ca-certificates tzdata && \
    echo "Dependencies installed"

WORKDIR /app

# Debug: Check files
COPY go.mod go.sum ./
RUN echo "Copied go.mod and go.sum" && \
    cat go.mod && \
    echo "Running go mod download..." && \
    go mod download && \
    echo "Download complete"

COPY . .
RUN echo "Copied source files" && \
    find . -name "*.go" | head -10 && \
    echo "Building binary..." && \
    CGO_ENABLED=0 GOOS=linux go build -v -o server cmd/server/main.go && \
    echo "Build complete" && \
    ls -la server

# Final stage
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/server /server
COPY --from=builder /app/configs /configs

EXPOSE 8080
ENTRYPOINT ["/server"]
```

**Common Solutions**:

1. **Fix .dockerignore**:
```
# .dockerignore
.git
.gitignore
README.md
Dockerfile*
.dockerignore
node_modules
vendor
```

2. **Optimize build**:
```dockerfile
# Use specific Go version
FROM golang:1.22-alpine AS builder

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Then copy source
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o server cmd/server/main.go
```

### Container Crashes on Start

**Problem**: Container exits immediately with code 1

**Diagnostic**:
```bash
# Check container logs
docker logs <container-id>

# Run interactively for debugging
docker run -it --entrypoint sh myapp:latest

# Check binary permissions
docker run --rm myapp:latest ls -la /

# Test with minimal command
docker run myapp:latest /server --help
```

**Solutions**:

1. **Fix entrypoint**:
```dockerfile
# Ensure binary is executable
COPY --from=builder /app/server /server
RUN chmod +x /server

# Use exec form
ENTRYPOINT ["/server"]
```

2. **Add health check**:
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/server", "health"] || exit 1
```

3. **Debug startup**:
```bash
# Run with environment variables
docker run -e LOG_LEVEL=debug -e PORT=8080 myapp:latest

# Check port binding
docker run -p 8080:8080 myapp:latest
```

## Platform-Specific Issues

### macOS Issues

**Security warnings on first run**:
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/go-starter

# Or allow in System Preferences → Security & Privacy
```

**Homebrew conflicts**:
```bash
# Fix installation conflicts
brew unlink go-starter
brew link go-starter --force

# Or completely reinstall
brew uninstall go-starter
brew install go-starter
```

### Windows Issues

**Path separator problems**:
```go
// Use filepath for cross-platform paths
import "path/filepath"

configPath := filepath.Join("configs", "config.yaml")
templatePath := filepath.Join("templates", "api.go.tmpl")
```

**Line ending issues**:
```bash
# Configure Git
git config --global core.autocrlf true

# Or use .gitattributes
echo "* text=auto" > .gitattributes
echo "*.go text eol=lf" >> .gitattributes
```

**PowerShell execution policy**:
```powershell
# Allow script execution
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Linux Issues

**Permission denied**:
```bash
# Fix binary permissions
chmod +x go-starter

# Or install to system location
sudo cp go-starter /usr/local/bin/
sudo chmod +x /usr/local/bin/go-starter
```

**Missing dependencies**:
```bash
# Check binary dependencies
ldd go-starter

# Install missing libraries (Ubuntu/Debian)
sudo apt-get update
sudo apt-get install libc6 ca-certificates

# For Alpine Linux
apk add --no-cache ca-certificates tzdata
```

## Getting Help

### Before Reporting Issues

1. **Run diagnostics**:
```bash
# Generate diagnostic report
go-starter version > diagnostic.txt
go version >> diagnostic.txt
go env >> diagnostic.txt
uname -a >> diagnostic.txt
echo "PATH: $PATH" >> diagnostic.txt
```

2. **Create minimal reproduction**:
```bash
# Test with simple project
go-starter new test-minimal --type=cli --complexity=simple --dry-run
```

3. **Check existing issues**:
- Search [GitHub Issues](https://github.com/francknouama/go-starter/issues)
- Check [closed issues](https://github.com/francknouama/go-starter/issues?q=is%3Aissue+is%3Aclosed)

### Where to Get Help

1. **GitHub Issues**: [Report bugs](https://github.com/francknouama/go-starter/issues/new)
2. **GitHub Discussions**: [Community help](https://github.com/francknouama/go-starter/discussions)
3. **Documentation**: [User guides](../02-user-guides/)

### What to Include in Bug Reports

- **System Information**: Output of diagnostic script
- **Exact Command**: The command that failed
- **Expected Behavior**: What should happen
- **Actual Behavior**: What actually happened
- **Reproduction Steps**: Minimal steps to reproduce
- **Generated Files**: If applicable, contents of generated files

### Emergency Workarounds

If go-starter is completely broken:

1. **Use manual installation**:
```bash
# Download release binary directly
curl -L -o go-starter https://github.com/francknouama/go-starter/releases/latest/download/go-starter-$(uname -s | tr '[:upper:]' '[:lower:]')-amd64
chmod +x go-starter
./go-starter version
```

2. **Build from source**:
```bash
git clone https://github.com/francknouama/go-starter.git
cd go-starter
go build -o go-starter main.go
./go-starter version
```

3. **Use Docker**:
```bash
docker run --rm -v $(pwd):/workspace francknouama/go-starter:latest new my-project --type=web-api
```

---

**Remember**: Most issues have been encountered before. Search the documentation and existing issues first, and don't hesitate to ask for help in the community!

**Next**: Check the [FAQ](faq.md) for answers to common questions, or return to [user guides](../02-user-guides/) for more advanced usage patterns.