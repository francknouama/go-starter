# logger Package

This package provides a unified logging interface for go-starter and the projects it generates, supporting multiple logging libraries.

## Overview

The logger package implements a flexible logging system that allows developers to choose from multiple high-performance logging libraries while maintaining a consistent interface across all generated code.

## Supported Loggers

| Logger | Package | Performance | Features |
|--------|---------|-------------|----------|
| **slog** | `log/slog` | Good | Standard library, structured logging |
| **zap** | `go.uber.org/zap` | Excellent | Zero allocation, high performance |
| **logrus** | `github.com/sirupsen/logrus` | Good | Feature-rich, hooks, formatters |
| **zerolog** | `github.com/rs/zerolog` | Excellent | Zero allocation, JSON output |

## Key Components

### Interface

```go
type Logger interface {
    Debug(msg string, fields ...any)
    Info(msg string, fields ...any)
    Warn(msg string, fields ...any)
    Error(msg string, fields ...any)
    Fatal(msg string, fields ...any)
    With(fields ...any) Logger
}
```

### Factory Functions

- **`NewLogger(loggerType string, config LogConfig) (Logger, error)`** - Create logger instance
- **`NewSlogLogger(config LogConfig) Logger`** - Create slog logger
- **`NewZapLogger(config LogConfig) Logger`** - Create zap logger
- **`NewLogrusLogger(config LogConfig) Logger`** - Create logrus logger
- **`NewZerologLogger(config LogConfig) Logger`** - Create zerolog logger

## Configuration

```go
type LogConfig struct {
    Level      string // debug, info, warn, error
    Format     string // json, text, console
    Output     string // stdout, stderr, file path
    TimeFormat string // RFC3339, Unix, etc.
}
```

## Usage Examples

### Creating a Logger

```go
import "github.com/yourusername/go-starter/internal/logger"

// Create slog logger (default)
log, err := logger.NewLogger("slog", logger.LogConfig{
    Level:  "info",
    Format: "json",
})

// Create zap logger for performance
log, err := logger.NewLogger("zap", logger.LogConfig{
    Level:  "debug",
    Format: "console",
})
```

### Using the Logger

```go
// Basic logging
log.Info("Server starting", "port", 8080, "env", "production")

// With context
requestLogger := log.With("request_id", "123-456")
requestLogger.Debug("Processing request", "method", "GET", "path", "/api/users")

// Error logging
if err != nil {
    log.Error("Failed to connect to database", "error", err, "retry", 3)
}
```

## Logger-Specific Features

### slog
- Builtin to Go 1.21+
- Handler customization
- Context support

### zap
- Sampling for high-volume logs
- Caller information
- Stack traces for errors

### logrus
- Hook system for integrations
- Multiple formatters
- Field-based logging

### zerolog
- Fastest JSON logger
- UNIX-friendly CLI output
- Context chaining

## Best Practices

1. Use structured logging with key-value pairs
2. Choose appropriate log levels
3. Include context (request ID, user ID, etc.)
4. Avoid logging sensitive information
5. Use sampling for high-frequency logs

## Performance Comparison

| Logger | Allocations | Time/Op | Use Case |
|--------|------------|---------|----------|
| slog | Low | Good | General purpose |
| zap | Zero | Fastest | High performance |
| logrus | Medium | Good | Feature-rich apps |
| zerolog | Zero | Fastest | JSON-heavy apps |