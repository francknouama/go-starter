# CLI Blueprint ATDD Tests

This directory contains Acceptance Test Driven Development (ATDD) tests for CLI application blueprints in the go-starter project.

## Overview

These tests validate that CLI blueprints generate functional, well-structured command-line applications that follow Go and CLI development best practices.

## Test Structure

### Files

- **`standard_test.go`** - Core CLI blueprint functionality tests
- **`integration_test.go`** - Cross-cutting concerns and integration tests
- **`README.md`** - This documentation file

### Test Categories

#### 1. Basic Generation Tests (`standard_test.go`)

**CLI Framework Integration**
- ✅ Cobra framework setup and configuration
- ✅ Root command definition and structure
- ✅ Subcommand generation (version, config, serve)
- ✅ Command-line flag handling
- ✅ Help text generation and formatting

**Configuration Management**
- ✅ Viper integration for configuration
- ✅ YAML configuration file support
- ✅ Environment variable binding
- ✅ Configuration precedence (flags > env > config file)
- ✅ Configuration validation and error handling

**Logger Integration**
- ✅ Support for multiple logger types (slog, zap, logrus, zerolog)
- ✅ Logger factory pattern implementation
- ✅ Structured logging capabilities
- ✅ Log level configuration and flags
- ✅ Logger-specific dependency management

**Binary Execution**
- ✅ Successful compilation to executable binary
- ✅ Command execution and help system
- ✅ Error handling and graceful shutdown
- ✅ Version command functionality

**Docker Support**
- ✅ Dockerfile generation
- ✅ Container build capability
- ✅ Multi-stage build optimization

#### 2. Integration Tests (`integration_test.go`)

**Cross-Logger Integration**
- ✅ All logger types work across CLI components
- ✅ Logger configuration consistency
- ✅ Binary execution with different loggers
- ✅ Performance validation across loggers

**Compilation Validation**
- ✅ All CLI configurations compile successfully
- ✅ Dependency resolution and management
- ✅ Binary executable generation
- ✅ Cross-platform compatibility

**Architecture Compliance**
- ✅ CLI best practices adherence
- ✅ Package organization and structure
- ✅ Dependency direction validation
- ✅ Architectural linting compliance

**Security Validation**
- ✅ No hardcoded secrets or sensitive data
- ✅ Secure configuration handling
- ✅ Safe logging practices
- ✅ Input validation and sanitization

**Runtime Validation**
- ✅ CLI execution without errors
- ✅ Help system functionality
- ✅ Configuration loading and validation
- ✅ Logging system operation

## CLI Blueprint Features Tested

### Core CLI Components

| Component | Feature | Test Coverage |
|-----------|---------|---------------|
| **Cobra Framework** | Root command setup | ✅ |
| **Cobra Framework** | Subcommand generation | ✅ |
| **Cobra Framework** | Flag handling | ✅ |
| **Cobra Framework** | Help text generation | ✅ |
| **Viper Configuration** | YAML file support | ✅ |
| **Viper Configuration** | Environment variables | ✅ |
| **Viper Configuration** | Flag binding | ✅ |
| **Viper Configuration** | Configuration precedence | ✅ |

### Logger Integration

| Logger Type | Dependency | Configuration | Binary Execution |
|-------------|------------|---------------|------------------|
| **slog** | ✅ Standard library | ✅ | ✅ |
| **zap** | ✅ go.uber.org/zap | ✅ | ✅ |
| **logrus** | ✅ github.com/sirupsen/logrus | ✅ | ✅ |
| **zerolog** | ✅ github.com/rs/zerolog | ✅ | ✅ |

### Project Structure Validation

```
CLI Project Structure:
├── main.go                     # Entry point
├── cmd/
│   ├── root.go                 # Root command
│   ├── version.go              # Version subcommand
│   ├── config.go               # Config subcommand
│   └── serve.go                # Serve subcommand
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   └── logger/
│       ├── interface.go        # Logger interface
│       ├── factory.go          # Logger factory
│       └── <logger>.go         # Logger implementations
├── configs/
│   └── config.yaml             # Configuration file
├── Dockerfile                  # Container support
├── Makefile                    # Build automation
└── go.mod                      # Module definition
```

## Test Execution

### Running CLI Tests

```bash
# Run all CLI tests
go test -v ./tests/acceptance/blueprints/cli/...

# Run specific test
go test -v ./tests/acceptance/blueprints/cli/ -run TestCLI_BasicGeneration

# Run with race detection
go test -v -race ./tests/acceptance/blueprints/cli/...

# Run with coverage
go test -v -coverprofile=coverage.out ./tests/acceptance/blueprints/cli/...
```

### Test Validation Process

1. **Project Generation**: Generate CLI project with specified configuration
2. **Structure Validation**: Verify correct file and directory structure
3. **Content Validation**: Check file contents for expected patterns
4. **Compilation Test**: Ensure project compiles without errors
5. **Binary Execution**: Test that generated binary runs correctly
6. **Integration Test**: Validate cross-cutting concerns

## Common Test Patterns

### Gherkin-Style BDD Testing

```go
func TestCLI_BasicGeneration_WithCobra(t *testing.T) {
    // Scenario: Generate basic CLI application
    // Given I want a CLI application
    // When I generate a CLI project
    // Then the project should use Cobra framework
    // And the project should have a root command
    // And the project should include example subcommands
    // And the project should compile to a working binary
}
```

### Validator Pattern

```go
validator := NewCLIValidator(projectPath)
validator.ValidateCobraSetup(t)
validator.ValidateRootCommand(t)
validator.ValidateCompilation(t)
validator.ValidateBinaryExecution(t)
```

### Logger Integration Testing

```go
loggers := []string{"slog", "zap", "logrus", "zerolog"}
for _, logger := range loggers {
    t.Run("Logger_"+logger, func(t *testing.T) {
        // Test logger-specific functionality
    })
}
```

## CLIValidator Methods

### Core Validation Methods

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateCobraSetup(t)` | Cobra framework integration | ✅ |
| `ValidateRootCommand(t)` | Root command structure | ✅ |
| `ValidateSubcommands(t, commands)` | Subcommand generation | ✅ |
| `ValidateCompilation(t)` | Project compilation | ✅ |
| `ValidateBinaryExecution(t)` | Binary execution | ✅ |

### Configuration Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateYAMLConfiguration(t)` | YAML config support | ✅ |
| `ValidateEnvironmentVariables(t)` | Environment variable binding | ✅ |
| `ValidateCommandLineFlags(t)` | Flag handling | ✅ |
| `ValidateConfigurationPrecedence(t)` | Config precedence | ✅ |

### Logger Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateLoggerIntegration(t, logger)` | Logger integration | ✅ |
| `ValidateLoggerFlags(t)` | Logger command flags | ✅ |
| `ValidateStructuredLogging(t, logger)` | Structured logging | ✅ |
| `ValidateLoggerConfiguration(t)` | Logger configuration | ✅ |

### Architecture Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateCLIArchitecture(t)` | CLI architecture principles | ✅ |
| `ValidateDependencyDirections(t)` | Dependency flow | ✅ |
| `ValidatePackageOrganization(t)` | Package structure | ✅ |
| `ValidateArchitecturalCompliance(t)` | Architecture compliance | ✅ |

### Security Validation

| Method | Purpose | Coverage |
|--------|---------|-----------|
| `ValidateSecurityPractices(t)` | Security best practices | ✅ |
| `ValidateSecureConfiguration(t)` | Secure config handling | ✅ |
| `ValidateSecureLogging(t)` | Safe logging practices | ✅ |

## Test Configuration Examples

### Basic CLI Configuration

```go
config := types.ProjectConfig{
    Name:      "test-cli-basic",
    Module:    "github.com/test/test-cli-basic",
    Type:      "cli-standard",
    GoVersion: "1.21",
    Framework: "cobra",
    Logger:    "slog",
}
```

### CLI with Different Loggers

```go
config := types.ProjectConfig{
    Name:      "test-cli-zap",
    Module:    "github.com/test/test-cli-zap",
    Type:      "cli-standard",
    GoVersion: "1.21",
    Framework: "cobra",
    Logger:    "zap",
}
```

## Known Limitations

1. **Runtime Testing**: Limited runtime testing due to test environment constraints
2. **Interactive Testing**: Cannot test interactive CLI features automatically
3. **Platform Testing**: Tests run on CI platform, may not cover all target platforms
4. **Network Testing**: Network-dependent CLI features are not tested

## Contributing

When adding new CLI blueprint tests:

1. **Follow BDD Pattern**: Use Gherkin-style comments for scenarios
2. **Test All Loggers**: Ensure new features work with all logger types
3. **Validate Compilation**: Always include compilation validation
4. **Document Changes**: Update this README with new test categories
5. **Use Validators**: Extend CLIValidator with new validation methods

## Future Enhancements

- [ ] **Interactive Testing**: Test interactive CLI features
- [ ] **Performance Testing**: Benchmark CLI startup and execution time
- [ ] **Multi-platform Testing**: Test on different operating systems
- [ ] **Plugin System**: Test CLI plugin architecture
- [ ] **Advanced Configuration**: Test complex configuration scenarios
- [ ] **Monitoring Integration**: Test observability features

## Related Documentation

- [Web API Blueprint Tests](../web-api/README.md)
- [Library Blueprint Tests](../library/README.md)
- [ATDD Testing Guide](../../README.md)
- [Blueprint Development Guide](../../../../docs/blueprints.md)

---

*This documentation is maintained alongside the CLI blueprint tests and should be updated when new test categories or validation methods are added.*