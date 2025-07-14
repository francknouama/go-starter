# utils Package

This package provides common utility functions used throughout the go-starter project.

## Overview

The utils package contains helper functions, common operations, and shared utilities that don't belong to a specific domain but are used across multiple packages.

## Key Components

### File Operations

- **`FileExists(path string) bool`** - Check if file exists
- **`IsDirectory(path string) bool`** - Check if path is directory
- **`CreateDirectory(path string) error`** - Create directory with parents
- **`CopyFile(src, dst string) error`** - Copy file with permissions
- **`WriteFileAtomic(path string, data []byte) error`** - Atomic file write

### String Utilities

- **`ToKebabCase(s string) string`** - Convert to kebab-case
- **`ToCamelCase(s string) string`** - Convert to CamelCase
- **`ToSnakeCase(s string) string`** - Convert to snake_case
- **`Pluralize(s string) string`** - Simple pluralization
- **`Truncate(s string, length int) string`** - Truncate with ellipsis

### Path Utilities

- **`NormalizePath(path string) string`** - Cross-platform path normalization
- **`ExpandHome(path string) string`** - Expand ~ to home directory
- **`RelativePath(base, target string) string`** - Calculate relative path
- **`IsSubPath(parent, child string) bool`** - Check if child is under parent

### Validation Helpers

- **`IsValidGoIdentifier(s string) bool`** - Validate Go identifier
- **`IsValidURL(s string) bool`** - Validate URL format
- **`IsValidEmail(s string) bool`** - Validate email address
- **`IsValidSemver(s string) bool`** - Validate semantic version

### System Utilities

- **`GetHomeDir() string`** - Get user home directory
- **`GetWorkingDir() string`** - Get current working directory
- **`IsTerminal() bool`** - Check if running in terminal
- **`GetTerminalWidth() int`** - Get terminal width

### Error Handling

- **`Must(err error)`** - Panic on error (for init functions)
- **`Retry(fn func() error, attempts int) error`** - Retry with backoff
- **`WrapError(err error, msg string) error`** - Wrap error with context
- **`IsRetryableError(err error) bool`** - Check if error is retryable

## Usage Examples

### File Operations
```go
import "github.com/yourusername/go-starter/internal/utils"

// Check file existence
if utils.FileExists("config.yaml") {
    // Process config
}

// Atomic write
data := []byte("content")
err := utils.WriteFileAtomic("output.txt", data)
```

### String Manipulation
```go
// Convert naming conventions
tableName := utils.ToSnakeCase("UserAccount")  // "user_account"
structName := utils.ToCamelCase("user-account") // "UserAccount"
urlSlug := utils.ToKebabCase("UserAccount")     // "user-account"
```

### Path Operations
```go
// Expand home directory
configPath := utils.ExpandHome("~/.config/app.yaml")

// Check path relationship
if utils.IsSubPath("/app", "/app/data/file.txt") {
    // Path is safe
}
```

### Validation
```go
// Validate inputs
if !utils.IsValidGoIdentifier(varName) {
    return fmt.Errorf("invalid identifier: %s", varName)
}

if utils.IsValidURL(endpoint) {
    // Process URL
}
```

### Error Handling
```go
// Retry operation
err := utils.Retry(func() error {
    return connectToService()
}, 3)

// Wrap errors
if err != nil {
    return utils.WrapError(err, "failed to process request")
}
```

## Best Practices

1. **Pure Functions** - Utils should be stateless and side-effect free
2. **Error Handling** - Return errors instead of panicking
3. **Cross-Platform** - Ensure compatibility across OS
4. **Performance** - Optimize for common use cases
5. **Testing** - Comprehensive unit tests for all utilities

## Common Patterns

### Safe Operations
```go
// Safe file operations with cleanup
func SafeWriteFile(path string, data []byte) error {
    dir := filepath.Dir(path)
    if !utils.FileExists(dir) {
        if err := utils.CreateDirectory(dir); err != nil {
            return err
        }
    }
    return utils.WriteFileAtomic(path, data)
}
```

### Input Sanitization
```go
// Sanitize user input
projectName := utils.ToKebabCase(
    utils.Truncate(userInput, 50),
)
```

## Performance Considerations

- Cached regex compilation
- Efficient string builders
- Minimal allocations
- Lazy initialization

## Dependencies

- **golang.org/x/sys** - System calls
- **github.com/mitchellh/go-homedir** - Home directory
- **github.com/spf13/afero** - File system abstraction