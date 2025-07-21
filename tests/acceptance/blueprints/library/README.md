# Library Blueprint ATDD Test Suite

## Overview

This directory contains comprehensive Acceptance Test-Driven Development (ATDD) tests for the `library-standard` blueprint. The test suite uses Behavior-Driven Development (BDD) with Gherkin feature files to ensure generated Go libraries follow best practices for reusable packages.

## Architecture

### Library Blueprint Philosophy

The library blueprint focuses on creating professional, reusable Go packages with:
- **Minimal API Surface**: Clean, well-designed public interfaces
- **Zero Forced Dependencies**: Optional logging through dependency injection
- **Production-Ready**: Comprehensive testing, documentation, and CI/CD
- **Best Practices**: Functional options, context support, thread safety

### Test Structure

```
library/
├── features/
│   └── library.feature          # Gherkin scenarios (25+ test cases)
├── library_steps_test.go        # BDD step definitions
├── library_acceptance_test.go   # High-level acceptance tests
└── README.md                    # This documentation
```

## Feature Coverage

### 📚 Core Library Generation
- **Essential Components**: Main library file, tests, documentation
- **Project Structure**: Examples directory, CI/CD workflows
- **Go Module Support**: Proper module initialization and versioning
- **Package Best Practices**: Proper naming, documentation, exports

### 🔧 API Design Patterns
- **Functional Options**: Configuration through option functions
- **Minimal Public API**: Limited exported types and functions
- **Error Handling**: Custom error types with wrapping support
- **Context Support**: Cancellation and timeout handling
- **Thread Safety**: Concurrent access protection

### 📝 Documentation
- **README**: Comprehensive with installation, usage, and examples
- **GoDoc**: Package documentation following Go standards
- **Examples**: Executable example programs in separate directory
- **CHANGELOG**: Version history with semantic versioning

### 🧪 Testing Infrastructure
- **Unit Tests**: Table-driven tests with good coverage
- **Example Tests**: Testable documentation examples
- **Benchmark Tests**: Performance measurement
- **Coverage Reports**: Integration with CI for coverage tracking

### 🔌 Logger Integration
- **Optional Logging**: Never forced on library consumers
- **Dependency Injection**: Logger provided through options
- **Multiple Loggers**: Support for slog, zap, logrus, zerolog
- **Zero Logger Coupling**: Library works without any logger

### 🚀 CI/CD and Distribution
- **GitHub Actions**: Automated testing and releases
- **Multi-Version Testing**: Test on multiple Go versions
- **Release Automation**: Semantic versioning with tags
- **Quality Checks**: Linting, formatting, security scanning

### 📄 Licensing Options
- **MIT License**: Permissive open source
- **Apache 2.0**: Patent protection included
- **GPL 3.0**: Copyleft licensing
- **BSD 3-Clause**: Simple permissive license

## Test Scenarios

### Basic Library Generation
```gherkin
Scenario: Generate basic library-standard project
  Given I want to create a reusable Go library
  When I run the command "go-starter new my-library --type=library-standard --module=github.com/example/my-library --no-git"
  Then the generation should succeed
  And the project should contain all essential library components
  And the library should follow Go package best practices
```

### Logger Integration Testing
```gherkin
Scenario: Library with different logging implementations
  Given I want to create libraries with various logging options
  When I generate a library with logger "zap"
  Then the library should use dependency injection for logging
  And the logger should be optional and not forced on consumers
```

### API Design Validation
```gherkin
Scenario: Library API design patterns
  Given I want a library with clean API design
  When I generate a library following best practices
  Then the library should use functional options pattern
  And the library should have minimal public API surface
  And the library should be thread-safe
```

## Running Tests

### Prerequisites
```bash
# Install required dependencies
go install github.com/cucumber/godog/cmd/godog@latest

# Ensure go-starter CLI is available
go-starter --help

# Verify Go version (1.21+ recommended)
go version
```

### Test Execution

#### Run All Library Tests
```bash
# Full test suite with BDD scenarios
go test -v ./tests/acceptance/blueprints/library/

# Run with race detection
go test -race -v ./tests/acceptance/blueprints/library/
```

#### Run Specific Test Categories
```bash
# Basic generation tests
go test -v -run TestBasicLibraryGeneration ./tests/acceptance/blueprints/library/

# Logger integration tests
go test -v -run TestLoggerIntegration ./tests/acceptance/blueprints/library/

# Documentation quality tests
go test -v -run TestDocumentationQuality ./tests/acceptance/blueprints/library/

# API design tests
go test -v -run TestAPIDesignPatterns ./tests/acceptance/blueprints/library/
```

#### Run BDD Scenarios
```bash
# Run all feature scenarios
godog run features/library.feature

# Run with specific tags
godog run features/library.feature -t @logging

# Generate test reports
godog run features/library.feature -f junit:reports/library-test-results.xml
```

#### Performance Testing
```bash
# Benchmark library generation
go test -bench=BenchmarkLibraryGeneration ./tests/acceptance/blueprints/library/

# Benchmark compilation speed
go test -bench=BenchmarkLibraryCompilation ./tests/acceptance/blueprints/library/
```

### Test Configuration

#### Environment Variables
```bash
# Test configuration
export TEST_TIMEOUT=30s

# Temporary directory for test projects
export TMPDIR=/tmp/library-tests

# Skip cleanup for debugging
export SKIP_CLEANUP=true
```

#### CI/CD Integration
```bash
# Short test mode (skip long-running tests)
go test -short ./tests/acceptance/blueprints/library/

# Parallel execution
go test -parallel 4 ./tests/acceptance/blueprints/library/

# Coverage reporting
go test -coverprofile=coverage.out ./tests/acceptance/blueprints/library/
go tool cover -html=coverage.out -o coverage.html
```

## Test Implementation Details

### Code Analysis

The test suite includes Go AST parsing for deep code analysis:

```go
// Parse and analyze generated Go code
file, err := parser.ParseFile(fset, "library.go", content, parser.ParseComments)

// Check for exported vs unexported members
ast.Inspect(file, func(n ast.Node) bool {
    switch node := n.(type) {
    case *ast.FuncDecl:
        if ast.IsExported(node.Name.Name) {
            // Validate exported function
        }
    }
    return true
})
```

### Documentation Validation

Comprehensive documentation checks:

```go
// Verify README sections
requiredSections := []string{
    "## Installation",
    "## Usage",
    "## API",
    "## Examples",
    "## Contributing",
    "## License",
}
```

### Logger Independence Testing

Validates that loggers are truly optional:

```go
// Test compilation without logger dependency
// Remove logger from go.mod and verify library still compiles
// This ensures logger is not tightly coupled
```

## Test Categories

### 1. Generation Tests
- **Focus**: Project structure and file generation
- **Coverage**: All required files and directories
- **Validation**: go.mod correctness, file permissions

### 2. Code Quality Tests
- **Focus**: Go best practices and conventions
- **Coverage**: Package naming, exports, documentation
- **Validation**: AST parsing, go vet, golangci-lint

### 3. API Design Tests
- **Focus**: Library interface design
- **Coverage**: Functional options, error types, context usage
- **Validation**: Minimal exports, thread safety

### 4. Documentation Tests
- **Focus**: User-facing documentation
- **Coverage**: README, GoDoc, examples
- **Validation**: Completeness, accuracy, executability

### 5. Testing Infrastructure Tests
- **Focus**: Test setup and execution
- **Coverage**: Unit tests, benchmarks, examples
- **Validation**: Coverage, performance, conventions

### 6. Integration Tests
- **Focus**: Logger and dependency integration
- **Coverage**: All supported loggers
- **Validation**: Optional dependencies, clean interfaces

### 7. Distribution Tests
- **Focus**: Publishing and versioning
- **Coverage**: CI/CD, releases, module support
- **Validation**: go get compatibility, versioning

## Troubleshooting

### Common Issues

#### Generation Failures
```bash
# Check go-starter installation
which go-starter
go-starter version

# Verify command syntax
go-starter new --help
```

#### Compilation Errors
```bash
# Clean module cache
go clean -modcache

# Update dependencies
cd generated-library && go mod tidy

# Check Go version compatibility
go version
```

#### Test Failures
```bash
# Run with verbose output
go test -v -run TestName

# Check for race conditions
go test -race

# Examine generated files
ls -la generated-library/
```

### Debug Mode

Enable detailed debugging:

```bash
# Verbose test output
export DEBUG=true
go test -v

# Keep generated projects
export SKIP_CLEANUP=true

# Print generation output
go test -v -args -test.v
```

## Performance Expectations

### Generation Performance

| Metric | Expected | Maximum |
|--------|----------|---------|
| Generation Time | < 2s | 5s |
| Compilation Time | < 5s | 10s |
| Test Execution | < 10s | 30s |
| Files Generated | 18+ | - |

### Resource Usage

- **CPU**: Low during generation
- **Memory**: < 100MB typical usage
- **Disk**: < 1MB per generated library

## Best Practices for Libraries

### API Design Principles

1. **Accept Interfaces, Return Structs**: Maximum flexibility
2. **Minimal Public API**: Export only what's necessary
3. **Functional Options**: Clean configuration pattern
4. **Context First**: Support cancellation in long operations
5. **Error Wrapping**: Rich error context with unwrapping

### Documentation Standards

1. **Package Documentation**: Clear overview in doc.go
2. **Function Documentation**: Start with function name
3. **Examples**: Executable examples with output
4. **README Structure**: Installation, usage, API reference
5. **Version History**: Semantic versioning in CHANGELOG

### Testing Philosophy

1. **Table-Driven Tests**: Comprehensive test cases
2. **Example Tests**: Documentation that runs
3. **Benchmarks**: Performance baselines
4. **Coverage Goals**: 80%+ for libraries
5. **Race Detection**: Always test for concurrency issues

## Contributing

### Adding New Test Scenarios

1. **Update Feature File**: Add scenarios to `features/library.feature`
2. **Implement Steps**: Add step definitions to `library_steps_test.go`
3. **Create Tests**: Add acceptance tests to `library_acceptance_test.go`
4. **Update Documentation**: Document new test scenarios

### Test Development Guidelines

1. **User Perspective**: Write scenarios from library user viewpoint
2. **Comprehensive Validation**: Check generation, compilation, and usage
3. **Performance Aware**: Keep tests fast and focused
4. **Real-World Usage**: Test practical library patterns
5. **Documentation**: Keep README and tests in sync

## Related Documentation

- [Monolith Blueprint ATDD](../monolith/README.md)
- [Web API Blueprint ATDD](../web-api/README.md)
- [CLI Blueprint ATDD](../cli/README.md)
- [Service Blueprint ATDD](../service/README.md)
- [Blueprint Development Guide](../../../../docs/blueprint-development.md)

## Support

For issues with library blueprint testing:

1. Check the [troubleshooting section](#troubleshooting)
2. Review generated library structure
3. Verify Go module configuration
4. Open an issue with reproduction steps

---

**Note**: The library blueprint emphasizes creating reusable Go packages that are a joy to use. The comprehensive test suite ensures that generated libraries follow Go community best practices while maintaining flexibility for different use cases and logging preferences.