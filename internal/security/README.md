# security Package

This package provides security validation and scanning functionality for go-starter, ensuring generated projects and blueprints are secure.

## Overview

The security package implements comprehensive security checks to prevent malicious code injection, path traversal attacks, and other security vulnerabilities in blueprints and generated projects.

## Key Components

### Types

- **`SecurityScanner`** - Main security scanning interface
- **`ScanResult`** - Results of security scan
- **`Vulnerability`** - Detected security issue
- **`SeverityLevel`** - Critical, High, Medium, Low

### Core Functions

- **`NewScanner(config ScannerConfig) *SecurityScanner`** - Create scanner
- **`ScanBlueprint(blueprint *Blueprint) (*ScanResult, error)`** - Scan blueprint
- **`ScanTemplate(content string) []Vulnerability`** - Scan template content
- **`ValidatePath(path string) error`** - Validate file paths
- **`SanitizeInput(input string) string`** - Sanitize user input

## Security Checks

### Path Traversal Protection
```go
// Prevents attempts to access files outside project directory
- Blocks ".." in paths
- Validates absolute vs relative paths
- Checks symlink targets
- Restricts file operations to project boundary
```

### Template Injection Prevention
```go
// Scans templates for malicious code
- Detects shell command execution
- Blocks system calls
- Validates template functions
- Restricts file I/O operations
```

### Input Validation
```go
// Validates all user inputs
- Project names: alphanumeric + dash/underscore
- Module paths: valid Go module syntax
- URLs: proper format, no javascript:
- File names: no special characters
```

### Blueprint Security
```go
// Ensures blueprint safety
- Scans for malicious scripts
- Validates dependencies
- Checks post-generation hooks
- Verifies file permissions
```

## Vulnerability Detection

### Categories

1. **Code Injection**
   - Shell command execution
   - Template function abuse
   - Script injection

2. **Path Traversal**
   - Directory traversal attempts
   - Symlink attacks
   - Absolute path usage

3. **Information Disclosure**
   - Hardcoded secrets
   - API keys in templates
   - Sensitive file exposure

4. **Dependency Risks**
   - Known vulnerable packages
   - Unverified dependencies
   - Outdated versions

## Usage Example

```go
import "github.com/yourusername/go-starter/internal/security"

// Create scanner
scanner := security.NewScanner(security.ScannerConfig{
    StrictMode: true,
    MaxDepth: 10,
})

// Scan blueprint
result, err := scanner.ScanBlueprint(blueprint)
if err != nil {
    log.Fatal(err)
}

// Check results
if result.HasCritical() {
    log.Fatal("Critical vulnerabilities found")
}

for _, vuln := range result.Vulnerabilities {
    fmt.Printf("[%s] %s: %s\n", vuln.Severity, vuln.Type, vuln.Description)
}
```

## Security Policies

### Default Policies
- No shell execution in templates
- Restricted file system access
- Validated external URLs
- Sanitized user inputs

### Configurable Rules
```yaml
security:
  policies:
    allow_shell_hooks: false
    max_template_depth: 5
    allowed_domains:
      - github.com
      - golang.org
    blocked_functions:
      - exec
      - system
      - eval
```

## Best Practices

1. **Regular Updates**
   - Update vulnerability database
   - Scan existing blueprints
   - Review security policies

2. **Blueprint Review**
   - Manual review for complex blueprints
   - Community blueprint vetting
   - Security scoring system

3. **User Education**
   - Security warnings in CLI
   - Documentation of risks
   - Safe blueprint guidelines

## Integration

### CLI Integration
- Automatic scanning during generation
- Security warnings and confirmations
- Bypass flags for trusted sources

### Web UI Integration
- Visual security indicators
- Detailed vulnerability reports
- Security score display

## Dependencies

- **github.com/google/go-safeweb/safehtml** - HTML sanitization
- **golang.org/x/tools/go/analysis** - Code analysis