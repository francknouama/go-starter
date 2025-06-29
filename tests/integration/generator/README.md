# Generator Test Suite

This directory contains comprehensive integration tests for the go-starter generator package. The tests are organized by functionality and follow consistent naming conventions to ensure maintainability and clarity.

## Test Organization

### Test Files Structure

```
tests/integration/generator/
├── README.md              # This file - comprehensive documentation
├── core_test.go          # Core generator functionality tests
├── validation_test.go    # Configuration validation tests
├── context_test.go       # Template context creation tests
├── go_version_test.go    # Go version handling tests
├── config_test.go        # Configuration processing tests
├── error_test.go         # Error handling and edge case tests
├── rollback_test.go      # Resource management and rollback tests
└── helpers_test.go       # Shared test utilities and helpers
```

### Test Naming Conventions

All tests follow the pattern: `TestGenerator_[Component]_[Scenario]`

- **Component**: The specific functionality being tested (e.g., `Generate`, `Preview`, `Validation`)
- **Scenario**: The specific test case or condition (e.g., `BasicFunctionality`, `InvalidConfig`, `DatabaseFeatures`)

Examples:
- `TestGenerator_Generate_BasicFunctionality`
- `TestGenerator_Validation_ProjectName`
- `TestGenerator_Context_LoggerConfiguration`
- `TestGenerator_GoVersion_CLI`

## Test Categories

### 1. Core Functionality Tests (`core_test.go`)

Tests the fundamental generator operations:

- **Generator Creation**: Tests `generator.New()` functionality
- **Basic Generation**: Tests project generation with various configurations
- **Preview Mode**: Tests preview functionality without file creation
- **Generation Results**: Tests the structure and content of generation results
- **Dry Run Mode**: Tests dry run functionality
- **Git Initialization**: Tests git repository setup

**Key Test Functions:**
- `TestGenerator_New`
- `TestGenerator_Generate_BasicFunctionality`
- `TestGenerator_Preview`
- `TestGenerator_GenerationResult`
- `TestGenerator_DryRun`
- `TestGenerator_GitInitialization`

### 2. Configuration Validation Tests (`validation_test.go`)

Tests all aspects of project configuration validation:

- **Project Name Validation**: Various project name formats and edge cases
- **Module Path Validation**: Go module path validation and formats
- **Project Type Validation**: Supported and unsupported project types
- **Go Version Validation**: Version format and compatibility checks
- **Framework Validation**: Framework compatibility with project types
- **Logger Validation**: Logger type validation and options
- **Architecture Validation**: Architecture pattern validation

**Key Test Functions:**
- `TestGenerator_Validation_ProjectName`
- `TestGenerator_Validation_ModulePath`
- `TestGenerator_Validation_ProjectType`
- `TestGenerator_Validation_GoVersion`
- `TestGenerator_Validation_Framework`
- `TestGenerator_Validation_Logger`
- `TestGenerator_Validation_Architecture`

### 3. Template Context Tests (`context_test.go`)

Tests template context creation and variable handling:

- **Basic Variables**: Standard template variables (ProjectName, ModulePath, etc.)
- **Logger Configuration**: Logger-specific context variables and flags
- **Database Configuration**: Database-related context variables and multi-database support
- **Authentication Configuration**: Authentication-related context variables
- **Feature Variables**: Feature-specific context creation
- **Template Variables**: Template-specific variable handling with defaults

**Key Test Functions:**
- `TestGenerator_Context_BasicVariables`
- `TestGenerator_Context_LoggerConfiguration`
- `TestGenerator_Context_DatabaseConfiguration`
- `TestGenerator_Context_AuthenticationConfiguration`
- `TestGenerator_Context_FeatureVariables`
- `TestGenerator_Context_TemplateVariables`

### 4. Go Version Handling Tests (`go_version_test.go`)

Tests Go version processing and CLI integration:

- **CLI Integration**: Tests `--go-version` flag functionality
- **Configuration Handling**: Tests Go version in generator configuration
- **Validation Logic**: Tests Go version validation and format checking
- **Context Variables**: Tests Go version in template context
- **Default Behavior**: Tests default Go version handling

**Key Test Functions:**
- `TestGenerator_GoVersion_CLI`
- `TestGenerator_GoVersion_Config`
- `TestGenerator_GoVersion_Validation`
- `TestGenerator_GoVersion_ContextVariables`
- `TestGenerator_GoVersion_DefaultBehavior`

### 5. Configuration Processing Tests (`config_test.go`)

Tests integration with the configuration system:

- **Profile Integration**: Tests generator integration with user profiles
- **Variable Handling**: Tests custom variable processing
- **Feature Configuration**: Tests feature configuration handling
- **Profile Defaults**: Tests profile default application and overrides

**Key Test Functions:**
- `TestGenerator_Config_Integration`
- `TestGenerator_Config_Variables`
- `TestGenerator_Config_Features`
- `TestGenerator_Config_ProfileDefaults`

### 6. Error Handling Tests (`error_test.go`)

Tests error conditions and edge cases:

- **Template Errors**: Missing templates, invalid template syntax
- **File System Errors**: Permission issues, disk space, path problems
- **Validation Errors**: Invalid configurations and error messages
- **Generation Failures**: Failed generation scenarios and recovery
- **Resource Cleanup**: Proper cleanup on failures

### 7. Rollback and Resource Management Tests (`rollback_test.go`)

Tests transaction support and rollback functionality:

- **Transaction Creation**: Tests generation transaction creation
- **File Tracking**: Tests file and directory tracking for rollback
- **Rollback Scenarios**: Tests rollback on various failure conditions
- **Partial Failures**: Tests rollback when generation partially succeeds
- **Resource Cleanup**: Tests proper resource cleanup and memory management

## Test Utilities and Helpers

### Setup Functions

- `setupTestTemplates(t *testing.T)`: Sets up the template filesystem for testing
- `buildTestBinary(t *testing.T) string`: Builds the CLI binary for integration tests

### Helper Functions

- `isEmptyDir(t *testing.T, dir string) bool`: Checks if a directory is empty
- Standard assertion helpers for common validation patterns

## Running Tests

### Run All Generator Tests
```bash
go test -v ./tests/integration/generator/
```

### Run Specific Test Categories
```bash
# Core functionality tests
go test -v ./tests/integration/generator/ -run TestGenerator_Generate

# Validation tests
go test -v ./tests/integration/generator/ -run TestGenerator_Validation

# Context tests
go test -v ./tests/integration/generator/ -run TestGenerator_Context

# Go version tests
go test -v ./tests/integration/generator/ -run TestGenerator_GoVersion

# Configuration tests
go test -v ./tests/integration/generator/ -run TestGenerator_Config
```

### Run Individual Tests
```bash
# Run a specific test
go test -v ./tests/integration/generator/ -run TestGenerator_Generate_BasicFunctionality

# Run tests with coverage
go test -v -cover ./tests/integration/generator/
```

## Test Data and Fixtures

### Mock Configurations

Tests use various mock configurations to simulate different scenarios:

- **Minimal Config**: Basic required fields only
- **Complete Config**: All optional fields populated
- **Profile Config**: Configuration with profile defaults applied
- **Feature Config**: Configuration with various features enabled

### Template Handling

Tests are designed to handle the current template development state:

- **Template Not Found**: Tests gracefully skip when templates are not yet implemented
- **Template Validation**: Tests verify template structure without requiring actual generation
- **Mock Templates**: Some tests use mock template structures for validation

## Error Handling Strategy

### Graceful Degradation

Tests are designed to handle the ongoing template development:

```go
// Accept template not found errors
if err != nil {
    if _, ok := err.(*types.TemplateNotFoundError); ok {
        t.Skip("Skipping test as template is not yet implemented")
        return
    }
}
```

### Error Verification

Tests verify specific error conditions and messages:

```go
if tt.shouldFail {
    assert.Error(t, err, "Expected validation to fail")
    if tt.errorContains != "" {
        assert.Contains(t, err.Error(), tt.errorContains, 
            "Error should contain expected message")
    }
}
```

## Best Practices

### Test Structure

1. **Setup**: Use `setupTestTemplates(t)` for consistent template configuration
2. **Test Data**: Use table-driven tests for multiple scenarios
3. **Cleanup**: Use `t.TempDir()` for automatic cleanup
4. **Error Handling**: Gracefully handle template-not-found conditions
5. **Assertions**: Use specific assertions with descriptive messages

### Test Isolation

- Each test uses isolated temporary directories
- Tests don't depend on external state or previous test results
- Mock configurations are self-contained within each test

### Performance Considerations

- Use dry run mode when testing configuration validation
- Skip file system operations when testing logic only
- Use parallel execution where appropriate (`t.Parallel()`)

## Integration with CI/CD

### Required Environment

Tests require:
- Go 1.19+ (for testing latest Go version features)
- Git (for git initialization tests)
- Sufficient disk space for temporary directories

### Test Coverage

Generator tests aim for comprehensive coverage of:
- All public generator methods
- All validation logic
- All error conditions
- All configuration combinations
- All template context scenarios

### Future Enhancements

Planned test improvements:
- **Benchmark Tests**: Performance testing for large projects
- **Parallel Generation**: Tests for concurrent generation scenarios
- **Memory Usage**: Tests for memory efficiency and cleanup
- **Cross-Platform**: Enhanced testing for Windows/macOS/Linux compatibility

## Troubleshooting

### Common Issues

1. **Template Not Found**: Normal during development - tests skip gracefully
2. **Binary Build Failures**: Check Go installation and project dependencies
3. **Permission Errors**: Ensure test has write access to temporary directories
4. **Git Errors**: Git initialization tests require git to be installed

### Debugging Tests

```bash
# Run with verbose output
go test -v ./tests/integration/generator/

# Run specific test with detailed output
go test -v ./tests/integration/generator/ -run TestGenerator_Generate_BasicFunctionality

# Enable race detection
go test -race ./tests/integration/generator/
```

## Contributing

When adding new generator tests:

1. **Follow Naming Conventions**: Use the established pattern
2. **Add Documentation**: Update this README with new test categories
3. **Handle Template State**: Use graceful error handling for missing templates
4. **Write Table Tests**: Use table-driven tests for multiple scenarios
5. **Add Coverage**: Ensure new functionality is comprehensively tested

### Test Review Checklist

- [ ] Tests follow naming conventions
- [ ] Tests handle template-not-found gracefully
- [ ] Tests use isolated temporary directories
- [ ] Tests have descriptive assertions
- [ ] Tests cover both success and failure scenarios
- [ ] Tests are documented in this README