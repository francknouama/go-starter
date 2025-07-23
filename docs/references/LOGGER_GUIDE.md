# Logger Guide

Simple, unified logging for go-starter generated projects.

## Table of Contents

- [Overview](#overview)
- [Logger Comparison](#logger-comparison)
- [Selection Guide](#selection-guide)
- [Implementation](#implementation)
- [Configuration](#configuration)
- [Best Practices](#best-practices)
- [Migration from Complex Logger](#migration-from-complex-logger)

---

## Overview

Go-starter provides a **Simplified Logger System** that generates clean, maintainable logging code with minimal complexity while supporting multiple popular Go logging libraries.

### Key Benefits

- ✅ **Minimal Interface**: Simple, focused logging API (5 methods)
- ✅ **Single Implementation**: One file with conditional compilation
- ✅ **Reduced Complexity**: 60-90% fewer lines of logging code
- ✅ **Conditional Dependencies**: Only selected logger dependencies included
- ✅ **Production Ready**: All logger types validated and tested

---

## Logger Comparison

### Detailed Feature Matrix

| Feature | slog | zap | logrus | zerolog |
|---------|------|-----|--------|---------|
| **Performance** | Good | Excellent | Good | Excellent |
| **Memory Allocation** | Low | Zero* | Medium | Zero* |
| **Structured Logging** | ✅ | ✅ | ✅ | ✅ |
| **JSON Output** | ✅ | ✅ | ✅ | ✅ |
| **Custom Formatters** | Limited | ✅ | ✅ | ✅ |
| **Hooks/Extensions** | Limited | ✅ | ✅ | Limited |
| **Context Integration** | ✅ | ✅ | ✅ | ✅ |
| **Standard Library** | ✅ | ❌ | ❌ | ❌ |
| **Go Version** | 1.21+ | 1.17+ | 1.13+ | 1.15+ |
| **Community Size** | Growing | Large | Largest | Large |
| **Learning Curve** | Easy | Medium | Easy | Medium |

*Zero allocation in production mode

### Performance Benchmarks

Based on independent benchmarks (operations per second, higher is better):

```
BenchmarkInfo
slog:     1,000,000 ops/sec    (100 ns/op)
zap:      3,000,000 ops/sec    (33 ns/op)  ⭐ Fastest
logrus:     500,000 ops/sec    (200 ns/op)
zerolog:  2,500,000 ops/sec    (40 ns/op)

BenchmarkInfoWithFields  
slog:       800,000 ops/sec    (125 ns/op)
zap:      2,500,000 ops/sec    (40 ns/op)   ⭐ Fastest
logrus:     300,000 ops/sec    (333 ns/op)
zerolog:  2,200,000 ops/sec    (45 ns/op)
```

---

## Selection Guide

### Choose **slog** when:
- ✅ **Getting started** with Go or structured logging
- ✅ **Minimizing dependencies** (standard library only)
- ✅ **Simple applications** with moderate logging needs
- ✅ **Long-term stability** is critical (maintained by Go team)
- ✅ **Go 1.21+** compatibility is acceptable

**Example use cases:**
- CLI tools and utilities
- Simple web APIs
- Learning projects
- Corporate environments with strict dependency policies

### Choose **zap** when:
- ✅ **High-performance applications** with heavy logging
- ✅ **Zero allocation** is critical for performance
- ✅ **Complex logging requirements** with custom fields
- ✅ **Production applications** requiring optimal performance
- ✅ **Uber's proven track record** gives confidence

**Example use cases:**
- High-throughput web services
- Real-time data processing
- Microservices with heavy logging
- Performance-critical applications

### Choose **logrus** when:
- ✅ **Feature-rich logging** with hooks and formatters
- ✅ **Large ecosystem** of extensions and integrations
- ✅ **Team familiarity** with logrus patterns
- ✅ **Gradual migration** from older logging libraries
- ✅ **Custom formatters** and output destinations needed

**Example use cases:**
- Enterprise applications
- Applications with complex logging requirements
- Legacy system migrations
- Applications requiring custom log formatting

### Choose **zerolog** when:
- ✅ **Cloud-native applications** requiring JSON logging
- ✅ **High performance** with clean API
- ✅ **Zero allocation** logging in hot paths
- ✅ **Chainable API** appeals to your team
- ✅ **Minimal memory footprint** is important

**Example use cases:**
- Kubernetes applications
- Serverless functions (AWS Lambda)
- Container-based microservices
- Cloud-native applications

---

## Implementation

### Simplified Logger Interface

All generated projects include a minimal, focused logger interface:

```go
// internal/logger/interface.go
package logger

type Logger interface {
    Debug(msg string, fields ...Fields)
    Info(msg string, fields ...Fields)
    Warn(msg string, fields ...Fields)
    Error(msg string, fields ...Fields)
    Fatal(msg string, fields ...Fields)
    WithFields(fields Fields) func(string, ...Fields)
}

type Fields map[string]interface{}
```

### Single-File Implementation

Each project includes one `logger.go` file with conditional compilation based on the selected logger:

```go
// internal/logger/logger.go
package logger

import (
    // Imports only the selected logger
    {{- if eq .LoggerType "slog" }}
    "log/slog"
    {{- else if eq .LoggerType "zap" }}
    "go.uber.org/zap"
    {{- end}}
)

var logger *selectedLogger

func Initialize(level string) error {
    // Initialize only the selected logger implementation
}

func Debug(msg string, fields ...Fields) {
    // Implementation specific to selected logger
}

// ... other methods
```

### Usage Examples

#### Basic Logging
```go
// Initialize once at startup
logger.Initialize("info")

// Use throughout application
logger.Info("Server starting")
logger.Error("Database connection failed")

// With structured fields
logger.Info("Request processed", logger.Fields{
    "method":   "POST",
    "path":     "/api/users", 
    "duration": 45,
    "status":   201,
})
```

#### Structured Logging
```go
// Create scoped logger with persistent fields
userLogger := logger.WithFields(logger.Fields{
    "userID":    user.ID,
    "operation": "user.create",
})

userLogger("Starting user creation")
userLogger("Validation complete") 
userLogger("User saved to database")
```

### Logger-Specific Implementations

#### slog Implementation
```go
// internal/logger/slog.go (generated when slog selected)
package logger

import (
    "context"
    "log/slog"
    "os"
)

type slogLogger struct {
    logger *slog.Logger
}

func newSlogLogger() Logger {
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
        AddSource: true,
    })
    
    return &slogLogger{
        logger: slog.New(handler),
    }
}

func (l *slogLogger) Info(message string, fields ...interface{}) {
    l.logger.Info(message, fields...)
}
```

#### Zap Implementation
```go
// internal/logger/zap.go (generated when zap selected)
package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type zapLogger struct {
    logger *zap.Logger
}

func newZapLogger() Logger {
    config := zap.NewProductionConfig()
    config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
    config.OutputPaths = []string{"stdout"}
    config.ErrorOutputPaths = []string{"stderr"}
    
    logger, _ := config.Build(
        zap.AddCaller(),
        zap.AddStacktrace(zapcore.ErrorLevel),
    )
    
    return &zapLogger{logger: logger}
}
```

#### Logrus Implementation
```go
// internal/logger/logrus.go (generated when logrus selected)
package logger

import (
    "github.com/sirupsen/logrus"
)

type logrusLogger struct {
    logger *logrus.Logger
}

func newLogrusLogger() Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: "2006-01-02T15:04:05.000Z",
    })
    logger.SetLevel(logrus.InfoLevel)
    
    return &logrusLogger{logger: logger}
}
```

#### Zerolog Implementation
```go
// internal/logger/zerolog.go (generated when zerolog selected)
package logger

import (
    "os"
    "github.com/rs/zerolog"
)

type zerologLogger struct {
    logger zerolog.Logger
}

func newZerologLogger() Logger {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    logger := zerolog.New(os.Stdout).
        Level(zerolog.InfoLevel).
        With().
        Timestamp().
        Caller().
        Logger()
    
    return &zerologLogger{logger: logger}
}
```

---

## Configuration

### Environment Variables

All loggers support common configuration through environment variables:

```bash
# Log level (debug, info, warn, error)
LOG_LEVEL=info

# Output format (json, text, console)
LOG_FORMAT=json

# Output destination (stdout, stderr, file path)
LOG_OUTPUT=stdout

# Include caller information (true, false)
LOG_CALLER=true

# Include timestamps (true, false)  
LOG_TIMESTAMP=true
```

### YAML Configuration

```yaml
# configs/config.yaml
logger:
  level: "info"              # debug, info, warn, error
  format: "json"             # json, text, console
  output: "stdout"           # stdout, stderr, file path
  caller: true               # Include caller information
  timestamp: true            # Include timestamps
  
  # Logger-specific settings
  slog:
    add_source: true         # Include source code location
    
  zap:
    development: false       # Enable development mode
    disable_caller: false    # Disable caller information
    disable_stacktrace: false # Disable stack traces
    sampling:
      initial: 100           # Initial sampling rate
      thereafter: 100        # Subsequent sampling rate
      
  logrus:
    disable_colors: false    # Disable colored output
    full_timestamp: true     # Use full timestamp format
    
  zerolog:
    sampling:
      basic: 1               # Basic sampling rate
      burst: 5               # Burst sampling rate
    caller_skip_frame_count: 0 # Skip stack frames
```

### Programmatic Configuration

```go
// Configure logger based on environment
func configureLogger() logger.Logger {
    level := os.Getenv("LOG_LEVEL")
    if level == "" {
        level = "info"
    }
    
    format := os.Getenv("LOG_FORMAT")
    if format == "" {
        format = "json"
    }
    
    // Logger factory handles configuration internally
    return logger.New()
}
```

---

## Best Practices

### 1. Structured Logging

Always use structured logging with key-value pairs:

```go
// ✅ Good - Structured fields
logger.Info("User login successful", 
    "userId", user.ID,
    "email", user.Email,
    "ip", request.RemoteAddr,
    "userAgent", request.UserAgent(),
    "duration", time.Since(start).Milliseconds(),
)

// ❌ Avoid - String interpolation
logger.Info(fmt.Sprintf("User %s (%s) logged in from %s", 
    user.Email, user.ID, request.RemoteAddr))
```

### 2. Use Appropriate Log Levels

```go
// DEBUG: Detailed diagnostic information
logger.Debug("Database query executed", 
    "query", sql,
    "params", params,
    "duration", duration,
)

// INFO: General operational information
logger.Info("Server started", 
    "port", port,
    "env", environment,
    "version", version,
)

// WARN: Warning conditions that don't stop operation
logger.Warn("Deprecated API endpoint used",
    "endpoint", "/api/v1/users",
    "client", clientID,
    "migrate-to", "/api/v2/users",
)

// ERROR: Error conditions that affect operation
logger.Error("Database connection failed",
    "error", err,
    "host", dbHost,
    "port", dbPort,
    "retry-attempt", retryCount,
)
```

### 3. Context Propagation

```go
func processOrder(ctx context.Context, order *Order) error {
    // Create logger with context
    logger := logger.New().WithContext(ctx)
    
    // Add order-specific fields
    orderLogger := logger.With(
        "orderId", order.ID,
        "customerId", order.CustomerID,
        "amount", order.Total,
    )
    
    orderLogger.Info("Processing order")
    
    if err := validateOrder(ctx, order); err != nil {
        orderLogger.Error("Order validation failed", "error", err)
        return err
    }
    
    orderLogger.Info("Order processed successfully")
    return nil
}
```

### 4. Error Logging

```go
func processRequest(req *Request) error {
    logger := logger.New()
    
    if err := validateRequest(req); err != nil {
        // Log validation errors with context
        logger.Error("Request validation failed",
            "error", err,
            "requestId", req.ID,
            "endpoint", req.Endpoint,
            "validation-rule", "required-fields",
        )
        return fmt.Errorf("validation failed: %w", err)
    }
    
    if err := processData(req.Data); err != nil {
        // Log processing errors with stack trace context
        logger.Error("Data processing failed",
            "error", err,
            "requestId", req.ID,
            "dataSize", len(req.Data),
            "processor", "main",
        )
        return fmt.Errorf("processing failed: %w", err)
    }
    
    return nil
}
```

### 5. Performance Considerations

```go
// ✅ Use log levels to avoid expensive operations
if logger.IsDebugEnabled() {
    logger.Debug("Request details",
        "headers", formatHeaders(req.Headers), // Expensive formatting
        "body", string(req.Body),              // Memory allocation
    )
}

// ✅ Use defer for timing measurements
func processRequest(req *Request) error {
    start := time.Now()
    defer func() {
        logger.Info("Request completed",
            "duration", time.Since(start).Milliseconds(),
            "requestId", req.ID,
        )
    }()
    
    // Process request...
    return nil
}

// ✅ Reuse loggers with common fields
type UserService struct {
    logger logger.Logger
}

func NewUserService() *UserService {
    return &UserService{
        logger: logger.New().With("service", "user"),
    }
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) error {
    requestLogger := s.logger.WithContext(ctx).With(
        "operation", "create-user",
        "email", req.Email,
    )
    
    requestLogger.Info("Creating user")
    // ... implementation
}
```

---

## Migration from Complex Logger

### Migrating from Legacy go-starter Projects

If you have an existing go-starter project with the old complex logger system (factory pattern, multiple implementation files), here's how to migrate:

#### 1. Update Logger Files

**Remove old complex files:**
```bash
rm internal/logger/factory.go
rm internal/logger/slog.go internal/logger/zap.go internal/logger/logrus.go internal/logger/zerolog.go
```

**Generate new simplified files:**
```bash
# Regenerate your project with latest go-starter
go-starter new my-project-updated --type=your-type --logger=your-logger

# Copy the new logger files
cp my-project-updated/internal/logger/logger.go internal/logger/
cp my-project-updated/internal/logger/interface.go internal/logger/
```

#### 2. Update Code Usage

**Old complex interface:**
```go
// Before
factory := logger.NewFactory()
log, err := factory.Create(config)
log.InfoWith("Message", logger.Fields{"key": "value"})
```

**New simplified interface:**
```go  
// After
logger.Initialize("info")
logger.Info("Message", logger.Fields{"key": "value"})
```

#### 3. Update Initialization

**Old main.go:**
```go
factory := logger.NewFactory()
appLogger, err := factory.CreateFromProjectConfig("slog", "info", "json", true)
cmd.Execute(appLogger)
```

**New main.go:**
```go
logger.Initialize("info")
cmd.Execute()
```

### Complexity Reduction Results

- **CLI-Standard**: 1,051 → 98 lines (91% reduction)
- **Web-API-Standard**: 398 → 110 lines (72% reduction)  
- **Workspace**: 487 → 297 lines (39% reduction)
- **Lambda-Standard**: 316 → 282 lines (11% reduction)

---

## Troubleshooting

### Common Issues

#### 1. Logger Not Outputting Anything

**Problem**: Logger appears to work but no output is visible.

**Solutions**:
```go
// Check log level configuration
logger := logger.New()
logger.Debug("This won't show if level is INFO or higher")

// Use appropriate log level
logger.Info("This will show with INFO level")

// Check environment variables
LOG_LEVEL=debug go run main.go
```

#### 2. Performance Issues

**Problem**: Logging is slowing down the application.

**Solutions**:
```go
// Use async logging for high-throughput applications
// (Implementation varies by logger)

// Avoid expensive operations in log statements
// ❌ Bad
logger.Info("User data", "user", formatUserData(user))

// ✅ Good  
if logger.IsDebugEnabled() {
    logger.Debug("User data", "user", formatUserData(user))
}

// Use sampling for high-frequency logs
logger.Info("Request processed", 
    "requestId", req.ID,
    "sample", rand.Intn(100) < 10, // Log 10% of requests
)
```

#### 3. JSON Parsing Errors

**Problem**: Log output contains invalid JSON.

**Solutions**:
```go
// Ensure all logged values are JSON-serializable
// ❌ Bad - functions, channels, etc. cause issues
logger.Info("Handler", "func", someFunction)

// ✅ Good - use string representation
logger.Info("Handler", "funcName", "handleUser")

// Sanitize user input
logger.Info("User input", "data", sanitizeForLogging(userInput))
```

#### 4. Missing Context Information

**Problem**: Logs don't include request ID, trace ID, etc.

**Solutions**:
```go
// Ensure context is properly propagated
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := context.WithValue(r.Context(), "requestId", generateRequestID())
    
    // Pass context to all functions
    if err := processRequest(ctx, req); err != nil {
        logger.WithContext(ctx).Error("Request failed", "error", err)
        return
    }
}

// Use middleware to add context automatically
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := generateRequestID()
        ctx := context.WithValue(r.Context(), "requestId", requestID)
        
        logger := logger.New().WithContext(ctx)
        ctx = context.WithValue(ctx, "logger", logger)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Performance Debugging

#### Enable Profiling
```go
import _ "net/http/pprof"

// Add to main function
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

#### Benchmark Different Loggers
```bash
# Test current setup
go test -bench=BenchmarkLogging -benchmem

# Compare with different logger
go-starter new test-zap --type=web-api --logger=zap
cd test-zap
go test -bench=BenchmarkLogging -benchmem
```

### Configuration Debugging

#### Verify Configuration Loading
```go
func debugConfig() {
    config := loadConfig()
    logger := logger.New()
    
    logger.Info("Logger configuration",
        "level", config.Logger.Level,
        "format", config.Logger.Format,
        "output", config.Logger.Output,
    )
}
```

#### Test Different Environments
```bash
# Test with different log levels
LOG_LEVEL=debug go run main.go
LOG_LEVEL=error go run main.go

# Test with different formats
LOG_FORMAT=text go run main.go
LOG_FORMAT=json go run main.go
```

This completes the comprehensive logger selector guide. The logger system in go-starter provides maximum flexibility while maintaining consistency and performance across all supported logging libraries.