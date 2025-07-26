# Test Utilities API Reference

## Overview

This document provides comprehensive API documentation for all test utilities, helper functions, and assertion methods available in the go-starter enhanced ATDD testing infrastructure.

## 📦 Package Organization

```
tests/helpers/
├── assertions.go          # Core assertion functions
├── test_utils.go         # Project generation and utilities
├── runtime.go            # Runtime testing helpers
├── security.go           # Security validation helpers
└── mocks/               # Mock implementations
    ├── filesystem_mock.go
    ├── prompter_mock.go
    └── template_registry_mock.go
```

## 🔧 Core Test Utilities

### Package: `tests/helpers`

#### Project Generation

##### `GenerateProject(t *testing.T, config TestConfig) string`

Generates a complete project based on the provided configuration.

**Parameters:**
- `t *testing.T`: Test instance for logging and cleanup
- `config TestConfig`: Project configuration specifying type, architecture, etc.

**Returns:**
- `string`: Absolute path to the generated project directory

**Example:**
```go
config := TestConfig{
    Type:         "web-api",
    Architecture: "hexagonal",
    Framework:    "gin", 
    Database:     "postgres",
    Logger:       "zap",
}
projectPath := helpers.GenerateProject(t, config)
```

**Cleanup:**
- Automatically registers cleanup function with `t.Cleanup()`
- Removes temporary directory after test completion

---

##### `GenerateProjectWithOptions(t *testing.T, config TestConfig, options GenerationOptions) string`

Extended project generation with additional options.

**Parameters:**
- `t *testing.T`: Test instance
- `config TestConfig`: Project configuration
- `options GenerationOptions`: Additional generation options

**GenerationOptions Fields:**
```go
type GenerationOptions struct {
    KeepTemporaryFiles bool          // Keep files after test for debugging
    DisableCleanup    bool           // Skip automatic cleanup
    CustomOutputDir   string         // Use specific output directory
    EnableDebugMode   bool           // Enable verbose debug logging
    PreGenerationHook func() error   // Called before generation
    PostGenerationHook func(string) error // Called after generation
}
```

**Example:**
```go
options := GenerationOptions{
    KeepTemporaryFiles: true,
    EnableDebugMode:    true,
    PostGenerationHook: func(projectPath string) error {
        return validateCustomRequirements(projectPath)
    },
}
projectPath := helpers.GenerateProjectWithOptions(t, config, options)
```

#### Test Configuration

##### `TestConfig` Structure

Complete configuration for project generation in tests.

```go
type TestConfig struct {
    // Required fields
    Type         string // "web-api", "cli", "library", "lambda", etc.
    Architecture string // "standard", "clean", "ddd", "hexagonal"
    Framework    string // "gin", "echo", "fiber", "chi", "cobra"
    
    // Optional fields
    Database     string // "postgres", "mysql", "sqlite", "mongodb"
    ORM          string // "gorm", "sqlx", "sqlc", "ent" 
    Auth         string // "jwt", "oauth2", "session", "api-key"
    Logger       string // "slog", "zap", "logrus", "zerolog"
    Complexity   string // "simple", "standard", "advanced"
    
    // Metadata
    Description  string // Human-readable test description
    Tags        []string // Test categorization tags
}
```

**Methods:**
```go
// Generate unique test name
func (tc TestConfig) Name() string

// Validate configuration consistency
func (tc TestConfig) Validate() error

// Apply default values for missing fields
func (tc TestConfig) WithDefaults() TestConfig

// Create configuration for specific use case
func (tc TestConfig) ForBlueprint(blueprint string) TestConfig
```

**Example:**
```go
config := TestConfig{
    Type:         "web-api",
    Architecture: "hexagonal",
    Framework:    "gin",
    Database:     "postgres",
    Description:  "Hexagonal architecture with Gin and PostgreSQL",
    Tags:        []string{"architecture", "database", "critical"},
}.WithDefaults()

if err := config.Validate(); err != nil {
    t.Fatalf("Invalid config: %v", err)
}
```

## ✅ Assertion Functions

### Compilation and Build Assertions

##### `AssertCompilationSuccess(t *testing.T, projectPath string)`

Verifies that the generated project compiles successfully.

**Implementation:**
- Runs `go build ./...` in the project directory
- Captures and reports build errors
- Validates all packages compile without warnings
- Checks for race conditions with `-race` flag

**Example:**
```go
projectPath := helpers.GenerateProject(t, config)
helpers.AssertCompilationSuccess(t, projectPath)
```

**Error Output:**
```
Compilation failed for project at /tmp/test-project-123:
go build output:
./internal/handlers/user.go:25:2: undefined: UnknownFunction
./internal/repository/user.go:15:1: imported and not used: "database/sql"
```

---

##### `AssertBuildTimeUnder(t *testing.T, projectPath string, maxDuration time.Duration)`

Ensures project builds within specified time limit.

**Parameters:**
- `projectPath string`: Path to project directory
- `maxDuration time.Duration`: Maximum allowed build time

**Example:**
```go
helpers.AssertBuildTimeUnder(t, projectPath, 30*time.Second)
```

---

##### `AssertBinaryExecutable(t *testing.T, binaryPath string)`

Validates that compiled binary is executable and functional.

**Checks:**
- File exists and has execute permissions
- Binary can be executed without errors
- Returns appropriate exit codes
- Responds to `--help` and `--version` flags (for CLI tools)

### Import and Dependency Assertions

##### `AssertNoUnusedImports(t *testing.T, projectPath string)`

Validates that all imports are used and no unused imports exist.

**Implementation:**
- Uses `go list -f '{{.GoFiles}}' ./...` to find all Go files
- Runs `goimports -d` to detect unused imports
- Validates import organization follows Go conventions

**Example:**
```go
helpers.AssertNoUnusedImports(t, projectPath) 
```

---

##### `AssertImportExists(t *testing.T, projectPath, packagePath, importPath string)`

Verifies specific import exists in a package.

**Parameters:**
- `packagePath string`: Relative path to package (e.g., "internal/handlers")
- `importPath string`: Expected import (e.g., "github.com/gin-gonic/gin")

**Example:**
```go
helpers.AssertImportExists(t, projectPath, "internal/handlers", "github.com/gin-gonic/gin")
```

---

##### `AssertDependencyVersion(t *testing.T, projectPath, module, version string)`

Validates specific dependency version in go.mod.

**Example:**
```go
helpers.AssertDependencyVersion(t, projectPath, "go.uber.org/zap", "v1.24.0")
```

### File System Assertions

##### `AssertFileExists(t *testing.T, projectPath, relativePath string)`

Verifies that a specific file exists in the project.

**Example:**
```go
helpers.AssertFileExists(t, projectPath, "internal/domain/user.go")
helpers.AssertFileExists(t, projectPath, "cmd/server/main.go")
```

---

##### `AssertFileCount(t *testing.T, projectPath string, expectedCount int)`

Validates the total number of generated files matches expectation.

**Implementation:**
- Recursively counts all files in project directory
- Excludes hidden files and directories (`.git`, `.DS_Store`, etc.)
- Provides detailed breakdown of file types in failure message

**Example:**
```go
// For simple CLI projects
helpers.AssertFileCount(t, projectPath, 8)

// For standard CLI projects  
helpers.AssertFileCount(t, projectPath, 29)
```

---

##### `AssertFileContains(t *testing.T, projectPath, relativePath, expectedContent string)`

Verifies file contains specific content.

**Example:**
```go
helpers.AssertFileContains(t, projectPath, "go.mod", "github.com/gin-gonic/gin")
helpers.AssertFileContains(t, projectPath, "README.md", "## Installation")
```

---

##### `AssertFilePermissions(t *testing.T, filePath string, expectedMode os.FileMode)`

Validates file permissions are set correctly.

**Example:**
```go
helpers.AssertFilePermissions(t, projectPath+"/scripts/deploy.sh", 0755)
```

### Architecture Pattern Assertions

##### `AssertCleanArchitectureLayers(t *testing.T, projectPath string)`

Validates Clean Architecture pattern implementation.

**Validations:**
- Entity layer exists with no external dependencies
- Use case layer depends only on entities and interfaces
- Interface adapters implement defined interfaces
- Frameworks and drivers are in outermost layer
- Dependency inversion principle is followed

**Directory Structure Checked:**
```
internal/
├── domain/           # Entities (innermost layer)
├── usecases/        # Use cases (business logic)
├── interfaces/      # Interface adapters
└── infrastructure/  # Frameworks & drivers (outermost)
```

**Example:**
```go
if config.Architecture == "clean" {
    helpers.AssertCleanArchitectureLayers(t, projectPath)
}
```

---

##### `AssertHexagonalBoundaries(t *testing.T, projectPath string)`

Validates Hexagonal Architecture (Ports & Adapters) pattern.

**Validations:**
- Ports are defined as interfaces
- Primary adapters (controllers) use ports
- Secondary adapters (repositories) implement ports
- Core business logic is isolated from external concerns
- Dependencies point inward toward the core

**Directory Structure Checked:**
```
internal/
├── domain/          # Core business logic
├── ports/           # Interfaces (ports)
└── adapters/        # Implementations (adapters)
    ├── primary/     # Driving adapters (controllers)
    └── secondary/   # Driven adapters (repositories)
```

**Example:**
```go
if config.Architecture == "hexagonal" {
    helpers.AssertHexagonalBoundaries(t, projectPath)
}
```

---

##### `AssertDDDDomainModel(t *testing.T, projectPath string)`

Validates Domain-Driven Design pattern implementation.

**Validations:**
- Aggregates are properly defined with root entities
- Domain services exist for complex business logic
- Value objects are immutable
- Domain events are properly implemented
- Bounded contexts are clearly separated
- Repository interfaces are in domain layer

**Directory Structure Checked:**
```
internal/domain/
├── aggregates/      # Aggregate roots and entities
├── services/        # Domain services
├── valueobjects/    # Value objects
├── events/          # Domain events
└── repositories/    # Repository interfaces
```

**Example:**
```go
if config.Architecture == "ddd" {
    helpers.AssertDDDDomainModel(t, projectPath)
}
```

### Database Integration Assertions

##### `AssertDatabaseConfiguration(t *testing.T, projectPath string, dbType, orm string)`

Validates database configuration and integration.

**Parameters:**
- `dbType string`: Database type ("postgres", "mysql", "sqlite", etc.)
- `orm string`: ORM framework ("gorm", "sqlx", "sqlc", etc.)

**Validations:**
- Correct database driver dependencies in go.mod
- Database configuration files exist
- Connection logic is properly implemented
- Migration files are present (if applicable)
- Database initialization code exists

**Example:**
```go
helpers.AssertDatabaseConfiguration(t, projectPath, "postgres", "gorm")
```

---

##### `AssertMigrationSupport(t *testing.T, projectPath, dbType string)`

Validates database migration infrastructure.

**Validations:**
- Migration directory exists
- Initial migration files are present
- Migration runner code exists
- Up/down migration support
- Migration versioning system

**Example:**
```go
helpers.AssertMigrationSupport(t, projectPath, "postgres")
```

### Authentication Assertions

##### `AssertAuthenticationSystem(t *testing.T, projectPath, authType string)`

Validates authentication system implementation.

**Parameters:**
- `authType string`: Authentication type ("jwt", "oauth2", "session", "api-key")

**Validations per Auth Type:**

**JWT:**
- JWT token generation and validation
- Middleware for token verification
- Token refresh mechanisms
- Proper secret management

**OAuth2:**
- OAuth2 flow implementation
- Provider configuration
- Token exchange logic
- User profile retrieval

**Session:**
- Session store configuration
- Session middleware
- Session lifecycle management
- CSRF protection

**Example:**
```go
if config.Auth != "" {
    helpers.AssertAuthenticationSystem(t, projectPath, config.Auth)
}
```

### Framework-Specific Assertions

##### `AssertFrameworkIntegration(t *testing.T, projectPath, framework string)`

Validates web framework integration and usage.

**Framework-Specific Validations:**

**Gin:**
- Router configuration
- Middleware setup
- Handler function signatures
- JSON binding usage

**Echo:**
- Echo instance configuration
- Context usage patterns
- Middleware registration
- Error handling

**Fiber:**
- Fiber app configuration
- Route registration
- Context methods usage
- Performance optimizations

**Example:**
```go
helpers.AssertFrameworkIntegration(t, projectPath, config.Framework)
```

---

##### `AssertNoFrameworkCrossContamination(t *testing.T, projectPath, expectedFramework string)`

Ensures only the selected framework is used.

**Validations:**
- No imports from other web frameworks
- No conflicting middleware
- Consistent context usage
- No mixed router patterns

**Example:**
```go
helpers.AssertNoFrameworkCrossContamination(t, projectPath, "gin")
```

### Logger Integration Assertions

##### `AssertLoggerIntegration(t *testing.T, projectPath, loggerType string)`

Validates logging framework integration.

**Logger-Specific Validations:**

**slog (Standard Library):**
- Proper slog handler configuration
- Structured logging usage
- Log level management
- Context-aware logging

**zap:**
- Zap logger initialization
- Field-based logging
- Performance optimizations
- Error handling

**logrus:**
- Logrus configuration
- Hook integrations
- Field and format management
- Level-based logging

**zerolog:**
- Zero-allocation patterns
- Chained logging calls
- Context integration
- Performance optimization

**Example:**
```go
helpers.AssertLoggerIntegration(t, projectPath, config.Logger)
```

### Performance Assertions

##### `MeasureBuildTime(t *testing.T, projectPath string) time.Duration`

Measures and returns project build time.

**Example:**
```go
buildTime := helpers.MeasureBuildTime(t, projectPath)
if buildTime > 30*time.Second {
    t.Errorf("Build time %v exceeds 30 seconds", buildTime)
}
```

---

##### `AssertMemoryUsageUnder(t *testing.T, projectPath string, maxMemory int64)`

Validates build memory usage stays under limit.

**Example:**
```go
helpers.AssertMemoryUsageUnder(t, projectPath, 512*1024*1024) // 512MB
```

---

##### `BenchmarkProjectGeneration(b *testing.B, config TestConfig)`

Benchmarks project generation performance.

**Example:**
```go
func BenchmarkWebAPIGeneration(b *testing.B) {
    config := TestConfig{Type: "web-api", Architecture: "clean"}
    helpers.BenchmarkProjectGeneration(b, config)
}
```

## 🎭 Mock Implementations

### Mock Prompter

##### `MockPrompter` Structure

Mock implementation for interactive CLI prompts.

```go
type MockPrompter struct {
    Responses map[string]interface{}
    CallLog   []string
}

// Methods
func (m *MockPrompter) SelectBlueprint(blueprints []string) (string, error)
func (m *MockPrompter) GetProjectName() (string, error)
func (m *MockPrompter) GetFramework(frameworks []string) (string, error)
func (m *MockPrompter) ConfirmGeneration(summary string) (bool, error)
```

**Example:**
```go
mockPrompter := &MockPrompter{
    Responses: map[string]interface{}{
        "blueprint": "web-api",
        "framework": "gin",
        "database":  "postgres",
        "confirm":   true,
    },
}

// Use in tests
generator := NewGenerator(WithPrompter(mockPrompter))
```

### Mock File System

##### `MockFileSystem` Structure

Mock implementation for file system operations.

```go
type MockFileSystem struct {
    Files       map[string][]byte
    Directories map[string]bool
    Operations  []string
}

// Methods
func (m *MockFileSystem) WriteFile(path string, data []byte) error
func (m *MockFileSystem) ReadFile(path string) ([]byte, error)
func (m *MockFileSystem) MkdirAll(path string) error
func (m *MockFileSystem) Exists(path string) bool
```

**Example:**
```go
mockFS := &MockFileSystem{
    Files: map[string][]byte{
        "/tmp/project/main.go": []byte("package main\n\nfunc main() {}"),
    },
    Directories: map[string]bool{
        "/tmp/project": true,
    },
}
```

### Mock Template Registry

##### `MockTemplateRegistry` Structure

Mock implementation for blueprint template loading.

```go
type MockTemplateRegistry struct {
    Templates map[string]*Template
    LoadCalls []string
}

// Methods
func (m *MockTemplateRegistry) LoadTemplate(name string) (*Template, error)
func (m *MockTemplateRegistry) ListTemplates() ([]string, error)
func (m *MockTemplateRegistry) ValidateTemplate(name string) error
```

## 🔍 Debugging Utilities

### Debug Helpers

##### `EnableDebugMode(t *testing.T)`

Enables verbose debug logging for a test.

**Example:**
```go
func TestComplexScenario(t *testing.T) {
    helpers.EnableDebugMode(t)
    // Test implementation with detailed logging
}
```

---

##### `DumpProjectStructure(t *testing.T, projectPath string)`

Prints detailed project directory structure for debugging.

**Example:**
```go
if testing.Verbose() {
    helpers.DumpProjectStructure(t, projectPath)
}
```

---

##### `KeepTemporaryFiles(t *testing.T, projectPath string)`

Prevents cleanup of generated files for manual inspection.

**Example:**
```go
func TestDebugScenario(t *testing.T) {
    projectPath := helpers.GenerateProject(t, config)
    helpers.KeepTemporaryFiles(t, projectPath)
    // Files will remain after test completion
}
```

## 📊 Metrics and Reporting

### Test Metrics

##### `CollectTestMetrics(t *testing.T, projectPath string) TestMetrics`

Collects comprehensive metrics about generated project.

```go
type TestMetrics struct {
    FileCount       int
    LineCount       int
    BuildTime       time.Duration
    BinarySize      int64
    Dependencies    int
    TestCoverage    float64
    ComplexityScore int
}
```

**Example:**
```go
metrics := helpers.CollectTestMetrics(t, projectPath)
t.Logf("Generated project metrics: %+v", metrics)
```

### Performance Monitoring

##### `MonitorResourceUsage(t *testing.T, fn func()) ResourceUsage`

Monitors CPU and memory usage during function execution.

```go
type ResourceUsage struct {
    MaxMemory    int64
    CPUTime      time.Duration
    SystemCalls  int64
    FileHandles  int
}
```

**Example:**
```go
usage := helpers.MonitorResourceUsage(t, func() {
    helpers.AssertCompilationSuccess(t, projectPath) 
})
t.Logf("Resource usage: %+v", usage)
```

## 🔐 Security Validation

### Security Helpers

##### `AssertNoSecretsInCode(t *testing.T, projectPath string)`

Scans generated code for accidentally committed secrets.

**Detects:**
- API keys and tokens
- Database passwords
- Private keys
- Hardcoded credentials

**Example:**
```go
helpers.AssertNoSecretsInCode(t, projectPath)
```

---

##### `AssertSecureDefaults(t *testing.T, projectPath string)`

Validates security-related default configurations.

**Checks:**
- HTTPS by default
- Secure session configurations
- Input validation presence
- SQL injection prevention
- XSS protection headers

**Example:**
```go
helpers.AssertSecureDefaults(t, projectPath)
```

This comprehensive API reference provides all the tools needed to create robust, maintainable tests for the go-starter project generator.