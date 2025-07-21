# Integration Tests

This directory contains comprehensive integration tests for the go-starter project, organized by component and test type for optimal maintainability and clarity.

## Directory Structure

### 📦 **blueprints/** - Blueprint-Specific Integration Tests
Blueprint tests validate that generated code compiles and functions correctly for each blueprint type.

- `web-api/` - Web API blueprint tests (all architectures)
- `cli/` - CLI blueprint tests (simple and standard) 
- `shared/` - Shared utilities for blueprint testing

### 🖥️ **cli/** - CLI Command Interface Tests  
Tests for the command-line interface, validating flags, commands, and user interactions.

### ⚙️ **generator/** - Generator Engine Tests
Tests for the core project generation engine, validation, and error handling.

### 📄 **templates/** - Template Engine Integration Tests
Tests for template compilation, parsing, and generation functionality.

### 🔄 **end-to-end/** - Full Workflow Integration Tests
Complete user workflow tests that span multiple components and validate entire generation processes.

## Running Integration Tests

```bash
# Run all integration tests
go test ./tests/integration/... -v

# Run specific component tests
go test ./tests/integration/blueprints/... -v
go test ./tests/integration/cli/... -v
go test ./tests/integration/generator/... -v
go test ./tests/integration/templates/... -v
go test ./tests/integration/end-to-end/... -v

# Run with race detection
go test ./tests/integration/... -v -race

# Run with coverage
go test ./tests/integration/... -v -cover
```

## Test Categories

- **Blueprint Integration**: Validates generated code compiles and functions
- **CLI Integration**: Tests command-line interface behavior
- **Generator Integration**: Tests core generation engine functionality  
- **Template Integration**: Tests template compilation and processing
- **End-to-End Integration**: Tests complete user workflows

## Contributing

When adding new integration tests:

1. Place tests in the appropriate directory based on what they test
2. Use descriptive file names with `_test.go` suffix
3. Follow the naming pattern: `Test[Component]_[Feature]_[Scenario]`
4. Include documentation for complex test scenarios
5. Ensure tests are independent and can run in parallel