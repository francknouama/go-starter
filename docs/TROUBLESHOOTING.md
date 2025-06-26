# Comprehensive Troubleshooting Guide

This guide provides detailed solutions for common issues and advanced troubleshooting techniques for go-starter.

## Table of Contents

- [Quick Diagnostics](#quick-diagnostics)
- [Installation Problems](#installation-problems)
- [Generation Issues](#generation-issues)
- [Build and Compilation Errors](#build-and-compilation-errors)
- [Runtime Issues](#runtime-issues)
- [Logger-Specific Problems](#logger-specific-problems)
- [Database Connection Issues](#database-connection-issues)
- [Docker and Deployment Issues](#docker-and-deployment-issues)
- [Performance Troubleshooting](#performance-troubleshooting)
- [Platform-Specific Issues](#platform-specific-issues)
- [Advanced Debugging](#advanced-debugging)

## Quick Diagnostics

Run this diagnostic script to check your environment:

```bash
#!/bin/bash
echo "=== go-starter Diagnostic Report ==="
echo "Date: $(date)"
echo ""
echo "1. go-starter Version:"
go-starter version 2>/dev/null || echo "NOT INSTALLED"
echo ""
echo "2. Go Version:"
go version || echo "Go NOT INSTALLED"
echo ""
echo "3. Go Environment:"
go env GOPATH
go env GOBIN
go env GOPROXY
echo ""
echo "4. PATH Check:"
echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "GOPATH/bin in PATH ✓" || echo "GOPATH/bin NOT in PATH ✗"
echo ""
echo "5. Network Connectivity:"
curl -s -o /dev/null -w "%{http_code}" https://proxy.golang.org && echo "Go proxy accessible ✓" || echo "Go proxy NOT accessible ✗"
echo ""
echo "6. Disk Space:"
df -h . | tail -1
echo "=== End Diagnostic Report ==="
```

## Installation Problems

### Issue: Installation Fails Silently

**Symptoms:**
- `go install` completes but binary not found
- No error messages displayed

**Diagnostic Steps:**
```bash
# Check where Go installs binaries
go env GOBIN
go env GOPATH

# Look for the binary
find $(go env GOPATH) -name "go-starter" -type f 2>/dev/null

# Check installation logs
go install -v github.com/francknouama/go-starter@latest
```

**Solutions:**

1. **Set GOBIN explicitly:**
```bash
export GOBIN=$HOME/bin
mkdir -p $GOBIN
go install github.com/francknouama/go-starter@latest
export PATH=$GOBIN:$PATH
```

2. **Direct installation to system path:**
```bash
sudo GOBIN=/usr/local/bin go install github.com/francknouama/go-starter@latest
```

3. **Build from source:**
```bash
git clone https://github.com/francknouama/go-starter.git
cd go-starter
make build
sudo cp bin/go-starter /usr/local/bin/
```

### Issue: Version Conflicts

**Symptoms:**
- Old version persists after update
- Multiple versions installed

**Solutions:**

```bash
# Find all installations
which -a go-starter
type -a go-starter

# Remove old versions
rm $(which go-starter)

# Clean module cache
go clean -modcache

# Reinstall latest
go install github.com/francknouama/go-starter@latest
```

## Generation Issues

### Issue: Template Variables Not Replaced

**Symptoms:**
- Generated files contain `{{.ProjectName}}` literals
- Template syntax visible in output

**Diagnostic Steps:**
```bash
# Check template processing
go-starter new test-project --type=web-api --verbose

# Examine generated files
grep -r "{{" test-project/
```

**Solutions:**

1. **Verify template syntax:**
```go
// Correct template usage
package {{.PackageName}}

// Incorrect (won't be processed)
package {{ .PackageName }}  // Extra spaces
```

2. **Check module path format:**
```bash
# Correct
go-starter new myapp --module=github.com/user/myapp

# Incorrect - will cause template issues
go-starter new myapp --module="My App"
```

### Issue: Partial Generation Failure

**Symptoms:**
- Some files created, others missing
- Generation stops mid-process

**Advanced Diagnostics:**
```bash
# Enable debug mode (if available)
export GO_STARTER_DEBUG=true
go-starter new test-project --type=web-api

# Check system resources during generation
watch -n 1 'df -h; echo "---"; free -h'

# Monitor file creation
watch -n 0.5 'find test-project -type f | wc -l'
```

**Recovery Steps:**
```bash
# 1. Clean partial generation
rm -rf test-project

# 2. Check disk space
df -h .
# Need at least 100MB free

# 3. Check file permissions
touch test-file && rm test-file || echo "No write permission"

# 4. Try in temp directory
cd /tmp
go-starter new test-project --type=web-api

# 5. If successful, move to desired location
mv test-project ~/projects/
```

## Build and Compilation Errors

### Issue: Import Cycle Detected

**Symptoms:**
```
import cycle not allowed
package github.com/user/project/internal/config
	imports github.com/user/project/internal/logger
	imports github.com/user/project/internal/config
```

**Solutions:**

1. **Identify the cycle:**
```bash
go mod graph | grep -E "(config|logger)"
```

2. **Break the cycle:**
```go
// Option 1: Interface segregation
// logger/interface.go
package logger

type Config interface {
    GetLogLevel() string
    GetLogFormat() string
}

// Option 2: Separate package
// create pkg/logconfig/config.go
package logconfig

type LogConfig struct {
    Level  string
    Format string
}
```

### Issue: Missing Dependencies

**Symptoms:**
```
cannot find module providing package go.uber.org/zap
```

**Comprehensive Fix:**
```bash
# 1. Clean module cache
go clean -modcache

# 2. Reset go.mod
rm go.mod go.sum
go mod init github.com/user/project

# 3. Re-download dependencies
go mod tidy
go mod download

# 4. Verify all dependencies
go mod verify

# 5. Build with verbose output
go build -v ./...
```

## Runtime Issues

### Issue: Panic on Startup

**Symptoms:**
- Application crashes immediately
- Panic stack trace displayed

**Debugging Steps:**
```bash
# 1. Run with race detector
go run -race cmd/server/main.go

# 2. Enable verbose logging
LOG_LEVEL=debug go run cmd/server/main.go

# 3. Use debugger
dlv debug cmd/server/main.go

# 4. Add recovery middleware
```

**Common Panic Causes:**

1. **Nil pointer dereference:**
```go
// Add nil checks
if logger == nil {
    logger = slog.Default()
}
```

2. **Configuration not loaded:**
```go
// Ensure config is loaded before use
config, err := LoadConfig()
if err != nil {
    log.Fatal("Failed to load config:", err)
}
```

### Issue: Port Already in Use

**Symptoms:**
```
listen tcp :8080: bind: address already in use
```

**Solutions:**

```bash
# 1. Find process using port
lsof -ti:8080
netstat -tulpn | grep :8080
ss -tulpn | grep :8080

# 2. Kill the process
kill -9 $(lsof -ti:8080)

# 3. Use different port
PORT=8081 go run cmd/server/main.go

# 4. Configure port in config
echo "server:\n  port: 8081" > configs/config.yaml
```

## Logger-Specific Problems

### Issue: Zap Logger Performance Degradation

**Symptoms:**
- Slow application performance
- High memory usage
- Logger allocations in profiles

**Solutions:**

```go
// 1. Use production config
func NewZapLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    config.DisableStacktrace = true
    config.Sampling = &zap.SamplingConfig{
        Initial:    100,
        Thereafter: 100,
    }
    
    logger, _ := config.Build()
    return logger
}

// 2. Pre-allocate fields
logger := logger.With(
    zap.String("service", "api"),
    zap.String("version", "1.0.0"),
)

// 3. Avoid string formatting
// Bad
logger.Info(fmt.Sprintf("User %s logged in", username))

// Good
logger.Info("User logged in", zap.String("user", username))
```

### Issue: Logrus Not Showing Colors

**Symptoms:**
- Terminal output is monochrome
- Colors work elsewhere but not in app

**Solutions:**

```go
// Force color output
logrus.SetFormatter(&logrus.TextFormatter{
    ForceColors:   true,
    FullTimestamp: true,
})

// Or detect terminal
if isatty.IsTerminal(os.Stdout.Fd()) {
    logrus.SetFormatter(&logrus.TextFormatter{
        ForceColors: true,
    })
}
```

### Issue: Zerolog Output Not JSON

**Symptoms:**
- Console output instead of JSON
- Human-readable format in production

**Solutions:**

```go
// Ensure JSON output
zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

// Conditional formatting
if os.Getenv("ENV") == "development" {
    logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
```

## Database Connection Issues

### Issue: GORM Connection Timeout

**Symptoms:**
- Application hangs on startup
- "connection refused" errors
- Timeout after 30 seconds

**Debugging Steps:**

```bash
# 1. Test database connectivity
psql -h localhost -U postgres -d testdb -c "SELECT 1"
mysql -h localhost -u root -p -e "SELECT 1"

# 2. Check connection string
echo $DATABASE_URL

# 3. Verify network
telnet localhost 5432
nc -zv localhost 5432
```

**Solutions:**

```go
// 1. Add connection retry logic
func ConnectWithRetry(dsn string, maxRetries int) (*gorm.DB, error) {
    var db *gorm.DB
    var err error
    
    for i := 0; i < maxRetries; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
            Logger: logger.Default.LogMode(logger.Info),
        })
        
        if err == nil {
            sqlDB, _ := db.DB()
            if err := sqlDB.Ping(); err == nil {
                return db, nil
            }
        }
        
        log.Printf("Database connection attempt %d failed: %v", i+1, err)
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
}

// 2. Configure connection pool
sqlDB, err := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### Issue: Migration Failures

**Symptoms:**
- "relation does not exist" errors
- Migrations partially applied
- Schema out of sync

**Solutions:**

```bash
# 1. Check migration status
go run main.go migrate status

# 2. Force migration
go run main.go migrate up --force

# 3. Reset database (development only)
go run main.go migrate down --all
go run main.go migrate up

# 4. Manual intervention
psql -d mydb -f migrations/001_initial.up.sql
```

## Docker and Deployment Issues

### Issue: Docker Build Fails

**Symptoms:**
- Build errors in multi-stage Dockerfile
- Missing dependencies in final image
- Large image sizes

**Advanced Dockerfile Debugging:**

```dockerfile
# Add debugging stages
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata

# Debug: List files
RUN ls -la

WORKDIR /app
COPY go.mod go.sum ./

# Debug: Check module files
RUN cat go.mod

RUN go mod download

# Debug: Verify download
RUN go mod verify

COPY . .

# Debug: List all files
RUN find . -type f -name "*.go" | head -20

RUN CGO_ENABLED=0 GOOS=linux go build -v -o server cmd/server/main.go

# Debug: Check binary
RUN ls -la server
RUN file server

# Final stage
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/server /server

ENTRYPOINT ["/server"]
```

### Issue: Container Crashes on Start

**Symptoms:**
- Container exits immediately
- No logs produced
- Status shows "Exited (1)"

**Debugging Steps:**

```bash
# 1. Check container logs
docker logs <container-id>

# 2. Run with shell for debugging
docker run -it --entrypoint sh myapp:latest

# 3. Check file permissions
docker run --rm myapp:latest ls -la /

# 4. Run with environment variables
docker run -e LOG_LEVEL=debug -e PORT=8080 myapp:latest

# 5. Use docker-compose for easier debugging
docker-compose up --build
docker-compose logs -f app
```

## Performance Troubleshooting

### Issue: High Memory Usage

**Symptoms:**
- Memory continuously growing
- OOM kills
- Slow garbage collection

**Profiling Steps:**

```go
// Add profiling endpoints
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

```bash
# Capture heap profile
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Monitor in real-time
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/heap
```

**Common Memory Leaks:**

1. **Goroutine leaks:**
```go
// Add goroutine monitoring
func monitorGoroutines() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        log.Printf("Active goroutines: %d", runtime.NumGoroutine())
    }
}
```

2. **Connection leaks:**
```go
// Always close connections
defer db.Close()
defer resp.Body.Close()
defer file.Close()
```

### Issue: Slow Request Performance

**Symptoms:**
- High latency
- Timeouts under load
- CPU spikes

**Performance Analysis:**

```bash
# 1. Run benchmarks
go test -bench=. -benchmem ./...

# 2. CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# 3. Trace analysis
go test -trace=trace.out -bench=.
go tool trace trace.out

# 4. Load testing
hey -n 10000 -c 100 http://localhost:8080/api/users
ab -n 10000 -c 100 http://localhost:8080/api/users
```

## Platform-Specific Issues

### macOS Issues

**Issue: Security warnings on first run**

```bash
# Solution: Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/go-starter
```

**Issue: Homebrew installation conflicts**

```bash
# Fix: Unlink and relink
brew unlink go-starter
brew link go-starter --force
```

### Windows Issues

**Issue: Path separators in generated code**

```go
// Use filepath package for cross-platform paths
import "path/filepath"

configPath := filepath.Join("configs", "config.yaml")
```

**Issue: Line ending problems**

```bash
# Configure Git
git config --global core.autocrlf true

# Or use .gitattributes
echo "* text=auto" > .gitattributes
```

### Linux Issues

**Issue: Permission denied on binary**

```bash
# Fix: Add execute permission
chmod +x go-starter
```

**Issue: Missing shared libraries**

```bash
# Check dependencies
ldd go-starter

# Install missing libraries (Ubuntu/Debian)
sudo apt-get install libc6
```

## Advanced Debugging

### Enable Trace Logging

```go
// Add trace logging to your application
func init() {
    if os.Getenv("TRACE") == "true" {
        log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
        log.SetPrefix("[TRACE] ")
    }
}

// Use throughout code
if os.Getenv("TRACE") == "true" {
    log.Printf("Function called with args: %+v", args)
}
```

### Remote Debugging

```bash
# 1. Build with debugging symbols
go build -gcflags="all=-N -l" -o app cmd/server/main.go

# 2. Run with Delve
dlv exec ./app

# 3. Or attach to running process
dlv attach $(pgrep app)
```

### Systematic Debugging Approach

1. **Isolate the problem:**
   - Minimal reproduction case
   - Remove dependencies one by one
   - Test in clean environment

2. **Gather information:**
   - Enable all logging
   - Capture system state
   - Record exact steps

3. **Form hypothesis:**
   - What changed recently?
   - What's different from working state?
   - What assumptions are made?

4. **Test systematically:**
   - Change one variable at a time
   - Document what works/fails
   - Build up from working state

5. **Document solution:**
   - What was the root cause?
   - How to prevent recurrence?
   - Update documentation

## Getting Additional Help

If these solutions don't resolve your issue:

1. **Collect diagnostic information:**
   ```bash
   go-starter version > diagnostic.txt
   go version >> diagnostic.txt
   go env >> diagnostic.txt
   uname -a >> diagnostic.txt
   ```

2. **Create minimal reproduction:**
   - Smallest possible example
   - Clear steps to reproduce
   - Expected vs actual behavior

3. **Report issue:**
   - [GitHub Issues](https://github.com/francknouama/go-starter/issues)
   - Include diagnostic.txt
   - Add reproduction steps
   - Specify template and logger used

Remember: Most issues have been encountered before. Check closed issues and discussions for additional solutions.