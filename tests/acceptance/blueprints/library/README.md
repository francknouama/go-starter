# Library Blueprint ATDD Tests

This directory contains Acceptance Test Driven Development (ATDD) tests for Go library blueprints in the go-starter project.

## Overview

These tests validate that Library blueprints generate well-structured, reusable Go libraries that follow Go library development best practices and provide excellent developer experience.

## Test Structure

### Files

- **`standard_test.go`** - Core Library blueprint functionality tests
- **`integration_test.go`** - Cross-cutting concerns and integration tests
- **`README.md`** - This documentation file

### Test Categories

#### 1. Basic Generation Tests (`standard_test.go`)

**Public API Design**
- ✅ Clean and intuitive public API structure
- ✅ Descriptive function names and types
- ✅ Appropriate parameter and return types
- ✅ Consistent error handling patterns
- ✅ Go naming conventions compliance

**Package Structure**
- ✅ Proper package organization
- ✅ Public API in root package
- ✅ Internal implementation encapsulation
- ✅ Example usage in dedicated directory
- ✅ Comprehensive test coverage

**Documentation**
- ✅ Package-level documentation (doc.go)
- ✅ Godoc comments for all public functions
- ✅ Comprehensive README with examples
- ✅ API reference documentation
- ✅ Usage examples and tutorials

**Examples and Testing**
- ✅ Basic and advanced usage examples
- ✅ Runnable example tests
- ✅ Unit tests for core functionality
- ✅ Benchmark tests for performance
- ✅ Example test validation

**Logger Integration**
- ✅ Minimal internal logging for debugging
- ✅ Support for multiple logger types (slog, zap, logrus, zerolog)
- ✅ Non-interfering logging design
- ✅ Configurable logging levels
- ✅ Structured logging capabilities

**Client Initialization**
- ✅ Clear initialization patterns
- ✅ Straightforward configuration
- ✅ Sensible default values
- ✅ Documented initialization process
- ✅ Flexible client options

#### 2. Integration Tests (`integration_test.go`)

**Cross-Logger Integration**
- ✅ All logger types work with library components
- ✅ Minimal logging footprint
- ✅ No interference with user logging
- ✅ Performance impact assessment

**Compilation Validation**
- ✅ All library configurations compile successfully
- ✅ Example compilation verification
- ✅ Test suite execution
- ✅ Library importability validation

**Architecture Compliance**
- ✅ Library best practices adherence
- ✅ Package organization validation
- ✅ Dependency direction enforcement
- ✅ Clean architecture principles

**Security Validation**
- ✅ No hardcoded secrets or sensitive data
- ✅ Secure logging practices
- ✅ Minimal dependency footprint
- ✅ No security risks for library users

**Usability Validation**
- ✅ Intuitive and well-documented API
- ✅ Clear and comprehensive examples
- ✅ Helpful error messages
- ✅ Good default configurations

**Maintenance Validation**
- ✅ Well-organized and readable code
- ✅ Comprehensive test coverage
- ✅ Complete build system
- ✅ Maintainable documentation

**API Compatibility**
- ✅ Clearly defined public API
- ✅ Properly encapsulated internal APIs
- ✅ Go conventions compliance
- ✅ Versioning support

## Library Blueprint Features Tested

### Core Library Components

| Component | Feature | Test Coverage |
|-----------|---------|---------------|
| **Public API** | Client initialization | ✅ |
| **Public API** | Core functionality | ✅ |
| **Public API** | Error handling | ✅ |
| **Public API** | Documentation | ✅ |
| **Examples** | Basic usage | ✅ |
| **Examples** | Advanced usage | ✅ |
| **Examples** | Runnable tests | ✅ |
| **Testing** | Unit tests | ✅ |
| **Testing** | Benchmark tests | ✅ |
| **Testing** | Example tests | ✅ |

### Logger Integration

| Logger Type | Dependency | Internal Use | Non-Interfering |
|-------------|------------|--------------|-----------------|
| **slog** | ✅ Standard library | ✅ | ✅ |
| **zap** | ✅ go.uber.org/zap | ✅ | ✅ |
| **logrus** | ✅ github.com/sirupsen/logrus | ✅ | ✅ |
| **zerolog** | ✅ github.com/rs/zerolog | ✅ | ✅ |

### Project Structure Validation

```
Library Project Structure:
├── <library-name>.go           # Main library interface
├── <library-name>_test.go      # Unit tests
├── examples_test.go            # Example tests
├── doc.go                      # Package documentation
├── examples/
│   ├── basic/
│   │   └── main.go             # Basic usage example
│   └── advanced/
│       └── main.go             # Advanced usage example
├── internal/
│   └── logger/
│       └── logger.go           # Internal logging
├── README.md                   # Comprehensive documentation
├── Makefile                    # Build automation
└── go.mod                      # Module definition
```

## Test Execution

### Running Library Tests

```bash
# Run all Library tests
go test -v ./tests/acceptance/blueprints/library/...

# Run specific test
go test -v ./tests/acceptance/blueprints/library/ -run TestLibrary_BasicGeneration

# Run with race detection
go test -v -race ./tests/acceptance/blueprints/library/...

# Run with coverage
go test -v -coverprofile=coverage.out ./tests/acceptance/blueprints/library/...
```

### Test Validation Process

1. **Project Generation**: Generate library project with specified configuration
2. **Structure Validation**: Verify correct file and directory structure
3. **API Validation**: Check public API design and documentation
4. **Example Validation**: Ensure examples compile and run correctly
5. **Compilation Test**: Verify library compiles without errors
6. **Test Execution**: Run unit tests and example tests
7. **Integration Test**: Validate cross-cutting concerns

## Common Test Patterns

### Gherkin-Style BDD Testing

```go
func TestLibrary_BasicGeneration(t *testing.T) {
    // Scenario: Generate basic library
    // Given I want a Go library
    // When I generate a library project
    // Then the project should have a clear public API
    // And the project should include example usage
    // And the project should have comprehensive tests
    // And the project should include proper documentation
}
```

### Validator Pattern

```go
validator := NewLibraryValidator(projectPath)
validator.ValidatePublicAPI(t)
validator.ValidateExampleUsage(t)
validator.ValidateCompilation(t)
validator.ValidateDocumentation(t)
```

### Logger Integration Testing

```go
loggers := []string{"slog", "zap", "logrus", "zerolog"}
for _, logger := range loggers {
    t.Run("Logger_"+logger, func(t *testing.T) {
        // Test minimal internal logging
    })
}
```

## LibraryValidator Methods

### Core Validation Methods

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidatePublicAPI(t)` | Public API structure | ✅ |
| `ValidateExampleUsage(t)` | Example generation | ✅ |
| `ValidateComprehensiveTests(t)` | Test coverage | ✅ |
| `ValidateDocumentation(t)` | Documentation quality | ✅ |
| `ValidateCompilation(t)` | Project compilation | ✅ |

### API Design Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateIntuitiveAPI(t)` | API intuitiveness | ✅ |
| `ValidateDescriptiveFunctionNames(t)` | Function naming | ✅ |
| `ValidateParameterTypes(t)` | Parameter design | ✅ |
| `ValidateReturnValueConventions(t)` | Return values | ✅ |
| `ValidateConsistentErrorHandling(t)` | Error handling | ✅ |

### Example and Testing Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateBasicAndAdvancedExamples(t)` | Example completeness | ✅ |
| `ValidateTestedExamples(t)` | Example testing | ✅ |
| `ValidateUnitTestCoverage(t)` | Unit test coverage | ✅ |
| `ValidateBenchmarkTests(t)` | Performance testing | ✅ |
| `ValidateRunnableExamples(t)` | Example execution | ✅ |

### Logger Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateMinimalLogging(t, logger)` | Minimal logging | ✅ |
| `ValidateConfigurableLogging(t)` | Logger configuration | ✅ |
| `ValidateNonInterfering(t)` | Non-interference | ✅ |
| `ValidateStructuredLogging(t, logger)` | Structured logging | ✅ |

### Documentation Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidatePackageDocumentation(t)` | Package docs | ✅ |
| `ValidateGodocComments(t)` | Godoc comments | ✅ |
| `ValidateDocumentationExamples(t)` | Doc examples | ✅ |
| `ValidateComprehensiveREADME(t)` | README quality | ✅ |

### Architecture Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateLibraryArchitecture(t)` | Library architecture | ✅ |
| `ValidateLibraryDependencyDirections(t)` | Dependency flow | ✅ |
| `ValidateLibraryPackageOrganization(t)` | Package structure | ✅ |
| `ValidateLibraryArchitecturalCompliance(t)` | Architecture compliance | ✅ |

### Security Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateLibrarySecurityPractices(t)` | Security practices | ✅ |
| `ValidateSecureLibraryLogging(t)` | Secure logging | ✅ |
| `ValidateSecureDependencies(t)` | Dependency security | ✅ |
| `ValidateNoUserSecurityRisks(t)` | User security | ✅ |

### Usability Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateIntuitiveAndDocumentedAPI(t)` | API usability | ✅ |
| `ValidateClearAndComprehensiveExamples(t)` | Example clarity | ✅ |
| `ValidateHelpfulErrorMessages(t)` | Error messages | ✅ |
| `ValidateGoodDefaults(t)` | Default values | ✅ |

## Test Configuration Examples

### Basic Library Configuration

```go
config := types.ProjectConfig{
    Name:      "test-library-basic",
    Module:    "github.com/test/test-library-basic",
    Type:      "library-standard",
    GoVersion: "1.21",
    Logger:    "slog",
}
```

### Library with Different Loggers

```go
config := types.ProjectConfig{
    Name:      "test-library-zap",
    Module:    "github.com/test/test-library-zap",
    Type:      "library-standard",
    GoVersion: "1.21",
    Logger:    "zap",
}
```

## API Design Principles Tested

### 1. Intuitive Interface Design

```go
// Good: Clear, descriptive function names
func NewClient(config Config) (*Client, error)
func (c *Client) Process(ctx context.Context, data []byte) (*Result, error)

// Good: Consistent error handling
func (c *Client) Start() error
func (c *Client) Stop() error
```

### 2. Minimal Logging

```go
// Internal logging should be minimal and non-interfering
logger := internal.NewLogger(internal.LogLevelWarn)
logger.Debug("internal operation completed")
```

### 3. Configuration Patterns

```go
// Good: Clear configuration with sensible defaults
type Config struct {
    Timeout    time.Duration `yaml:"timeout" default:"30s"`
    MaxRetries int          `yaml:"max_retries" default:"3"`
    Debug      bool         `yaml:"debug" default:"false"`
}
```

## Library Quality Metrics

### Code Quality Indicators

| Metric | Expected | Validated |
|--------|----------|-----------|
| **Public API Coverage** | All functions documented | ✅ |
| **Example Coverage** | Basic + Advanced examples | ✅ |
| **Test Coverage** | Unit + Integration + Benchmark | ✅ |
| **Documentation** | Package + Function + README | ✅ |
| **Dependencies** | Minimal external dependencies | ✅ |

### Performance Considerations

| Aspect | Requirement | Validation |
|--------|-------------|------------|
| **Startup Time** | Minimal initialization overhead | ✅ |
| **Memory Usage** | Efficient memory allocation | ✅ |
| **CPU Usage** | Optimized processing | ✅ |
| **Benchmark Tests** | Performance regression detection | ✅ |

## Known Limitations

1. **Runtime Performance**: Limited performance testing in CI environment
2. **Memory Profiling**: Basic memory usage validation only
3. **Concurrent Usage**: Limited concurrency testing
4. **Integration Testing**: Cannot test with real external services

## Contributing

When adding new Library blueprint tests:

1. **Follow BDD Pattern**: Use Gherkin-style comments for scenarios
2. **Test All Loggers**: Ensure new features work with all logger types
3. **Validate Examples**: Always test that examples compile and run
4. **Document API**: Ensure all public functions have documentation
5. **Test Usability**: Consider end-user experience in tests
6. **Update README**: Document new test categories and methods

## Future Enhancements

- [ ] **Performance Testing**: Comprehensive benchmarking and profiling
- [ ] **Concurrency Testing**: Test thread-safety and concurrent usage
- [ ] **Integration Testing**: Test with real external dependencies
- [ ] **Compatibility Testing**: Test with different Go versions
- [ ] **Memory Testing**: Advanced memory leak detection
- [ ] **API Stability**: Test backward compatibility

## Best Practices Validated

### 1. Library Design Principles

- **Single Responsibility**: Each library has a clear, focused purpose
- **Minimal Dependencies**: Libraries avoid unnecessary external dependencies
- **Stable API**: Public API is well-designed and stable
- **Good Documentation**: Comprehensive documentation and examples

### 2. Go Library Conventions

- **Package Naming**: Clear, descriptive package names
- **Export Patterns**: Proper use of exported vs unexported identifiers
- **Error Handling**: Consistent error handling throughout
- **Testing**: Comprehensive test coverage including examples

### 3. Developer Experience

- **Easy Installation**: Simple `go get` installation
- **Clear Examples**: Working examples for common use cases
- **Good Defaults**: Sensible default configurations
- **Helpful Errors**: Descriptive error messages

## Related Documentation

- [CLI Blueprint Tests](../cli/README.md)
- [Web API Blueprint Tests](../web-api/README.md)
- [ATDD Testing Guide](../../README.md)
- [Blueprint Development Guide](../../../../docs/blueprints.md)

---

*This documentation is maintained alongside the Library blueprint tests and should be updated when new test categories or validation methods are added.*