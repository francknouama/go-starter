# Multi-Binary Structure ATDD Test Suite

This directory contains comprehensive Acceptance Test-Driven Development (ATDD) tests for the new multi-binary structure of go-starter. These tests ensure that the restructuring maintains all functionality while properly supporting the new binary architecture.

## Overview

The go-starter project has been restructured to support multiple binaries:

- **`cmd/go-starter/`** - Main CLI tool
- **`cmd/go-starter-dev/`** - Development web server  
- **`web/cmd/web-server/`** - Production web server
- **`main.go`** - Legacy wrapper (deprecated)

## Test Files

### Core Test Suites

#### `multi_binary_test.go`
Comprehensive test suite covering all aspects of the multi-binary structure:

- **Multi-Binary Compilation Tests**: Verifies all binaries build independently
- **Installation Path Tests**: Tests various installation methods (`go install ./cmd/...`)
- **Backward Compatibility Tests**: Ensures existing workflows continue to work
- **Binary Functionality Tests**: Validates core functionality of each binary
- **Embedded Assets Tests**: Verifies embedded blueprints work correctly
- **Cross-Platform Tests**: Tests behavior across Windows/macOS/Linux
- **Migration Tests**: Validates smooth user migration experience

#### `multi_binary_steps_test.go`
BDD step definitions for Gherkin-style testing:

- Implements step definitions for the feature files
- Provides reusable test context and utilities
- Supports cucumber/godog integration
- Enables readable, business-focused test scenarios

#### `embedded_assets_test.go`
Focused tests for embedded blueprint functionality:

- **Blueprint Embedding**: Verifies all blueprints are embedded correctly
- **Project Generation**: Tests that embedded blueprints generate valid projects
- **Template Variable Processing**: Ensures template variables are processed correctly
- **Conditional Generation**: Tests conditional file generation logic
- **Validation**: Verifies blueprint validation works with embedded assets
- **Isolation Testing**: Confirms CLI works without filesystem access to blueprints

#### `cross_platform_test.go`
Platform-specific compatibility tests:

- **Binary Extensions**: Tests Windows `.exe` vs Unix binary naming
- **Path Separators**: Verifies correct path handling per platform
- **File Operations**: Tests file creation/permissions across platforms
- **Environment Variables**: Tests platform-specific environment handling
- **Cross-Compilation**: Verifies binaries can be built for different platforms
- **Unicode Support**: Tests character encoding handling

#### `performance_test.go`
Performance and scalability tests:

- **Build Performance**: Measures and validates build times
- **Startup Performance**: Tests CLI startup and response times
- **Memory Usage**: Validates reasonable memory consumption
- **Generation Performance**: Tests project generation speed
- **Concurrent Operations**: Validates concurrent usage scenarios
- **Resource Cleanup**: Ensures no resource leaks

### Feature Files

#### `features/multi-binary-structure.feature`
Gherkin-style feature definitions using Given-When-Then format:

```gherkin
Feature: Multi-Binary Structure Support
  As a developer using go-starter
  I want the new multi-binary structure to work correctly
  So that I can use the appropriate binary for my needs

Scenario: All binaries compile independently
  When I build the CLI binary from "cmd/go-starter"
  Then the build should succeed
  And the binary should be executable
  And the binary size should be reasonable
```

## Test Categories

### 1. Compilation Tests
- **Independent Building**: Each binary builds without dependencies on others
- **Cross-Platform**: Binaries build correctly on Windows, macOS, Linux
- **Binary Size**: Generated binaries are within reasonable size limits (5MB-60MB)
- **Executable Format**: Correct extensions and permissions per platform

### 2. Installation Tests
- **CLI Tool**: `go install ./cmd/go-starter` works correctly
- **Dev Server**: `go install ./cmd/go-starter-dev` installs properly
- **Legacy Support**: `go install .` still works with deprecation warning
- **PATH Integration**: Installed binaries are accessible from PATH

### 3. Functionality Tests
- **CLI Commands**: All existing commands work identically
- **Server Startup**: Web servers start and shutdown gracefully
- **API Endpoints**: Web server endpoints function correctly
- **Blueprint Access**: All blueprints remain accessible

### 4. Backward Compatibility Tests
- **Deprecation Warnings**: Legacy usage shows clear migration guidance
- **Command Compatibility**: All existing command syntax works
- **Output Format**: Consistent output formatting maintained
- **Documentation Examples**: Existing examples continue to work

### 5. Embedded Assets Tests
- **Blueprint Availability**: All 20+ blueprints accessible without filesystem
- **Template Processing**: Variables and conditionals work correctly
- **Project Generation**: Generated projects compile successfully
- **Isolation**: CLI works in environments without source code access

### 6. Performance Tests
- **Build Times**: Each binary builds under 60 seconds
- **Startup Speed**: CLI responds under 5 seconds
- **Memory Usage**: Reasonable memory consumption during operations
- **Concurrent Usage**: Multiple simultaneous operations work correctly

### 7. Migration Tests
- **Clear Guidance**: Deprecation messages provide actionable instructions
- **Working Examples**: Migration commands actually work
- **No Breaking Changes**: No functionality is lost in transition
- **Smooth Transition**: Easy upgrade path for existing users

## Running the Tests

### Full Test Suite
```bash
# Run all multi-binary ATDD tests
go test ./tests/acceptance/cli/... -v

# Run without performance tests (faster)
go test ./tests/acceptance/cli/... -v -short

# Run specific test suite
go test ./tests/acceptance/cli/multi_binary_test.go -v
```

### BDD Tests
```bash
# Run Gherkin/BDD tests
go test ./tests/acceptance/cli/multi_binary_steps_test.go -v
```

### Specific Categories
```bash
# Test embedded assets only
go test ./tests/acceptance/cli/embedded_assets_test.go -v

# Test cross-platform compatibility
go test ./tests/acceptance/cli/cross_platform_test.go -v

# Test performance characteristics
go test ./tests/acceptance/cli/performance_test.go -v
```

### CI/CD Integration
```bash
# Optimized for CI environments
go test ./tests/acceptance/cli/... -v -short -timeout=10m
```

## Test Requirements

### Prerequisites
- Go 1.21+
- Project source code available
- Write permissions for temporary directories
- Network access for Go module downloads

### Environment Setup
Tests automatically:
- Create temporary directories for isolation
- Build required binaries from source
- Clean up resources after completion
- Handle cross-platform differences

### Platform Support
Tests are designed to work on:
- **Linux** (amd64, arm64)
- **macOS** (amd64, arm64) 
- **Windows** (amd64)

## Acceptance Criteria

### ✅ Multi-Binary Compilation
- [ ] All 4 binaries build independently
- [ ] Build times under 60 seconds each
- [ ] Binary sizes 5MB-60MB range
- [ ] Correct file extensions per platform

### ✅ Installation Methods
- [ ] `go install ./cmd/go-starter` works
- [ ] `go install ./cmd/go-starter-dev` works
- [ ] `go install ./web/cmd/web-server` works
- [ ] Legacy `go install .` works with warning

### ✅ Functionality Preservation
- [ ] All CLI commands work identically
- [ ] All command-line flags function correctly
- [ ] All output formats remain consistent
- [ ] All existing workflows continue

### ✅ Embedded Assets
- [ ] All blueprints accessible without filesystem
- [ ] Generated projects compile successfully
- [ ] Template variables process correctly
- [ ] Conditional generation works

### ✅ Backward Compatibility
- [ ] No breaking changes introduced
- [ ] Clear deprecation warnings shown
- [ ] Migration instructions work
- [ ] Documentation examples valid

### ✅ Performance Standards
- [ ] CLI startup under 5 seconds
- [ ] Project generation under 30 seconds
- [ ] Memory usage reasonable
- [ ] Concurrent operations supported

### ✅ Cross-Platform Support
- [ ] Windows binary has .exe extension
- [ ] Unix binaries have correct permissions
- [ ] Path handling works per platform
- [ ] Cross-compilation succeeds

## Troubleshooting

### Common Issues

**Build Failures**
```bash
# Ensure you're in the project root
cd /path/to/go-starter
go mod tidy
go test ./tests/acceptance/cli/multi_binary_test.go -v
```

**Permission Errors**
```bash
# On Unix systems, ensure execute permissions
chmod +x /path/to/binary
```

**Path Issues**
```bash
# Verify project structure
ls cmd/go-starter/main.go
ls cmd/go-starter-dev/main.go
ls web/cmd/web-server/main.go
```

**Test Timeouts**
```bash
# Use shorter tests for CI
go test ./tests/acceptance/cli/... -short -timeout=5m
```

### Debug Mode
```bash
# Run with verbose output
go test ./tests/acceptance/cli/... -v -args -test.v

# Run single test for debugging
go test ./tests/acceptance/cli/multi_binary_test.go -run TestMultiBinaryCompilation -v
```

## Contributing

### Adding New Tests
1. Follow the existing test structure and naming conventions
2. Use ATDD methodology (Given-When-Then)
3. Include both positive and negative test cases
4. Test across all supported platforms
5. Add corresponding BDD scenarios if appropriate

### Test Guidelines
- **Isolation**: Each test should be independent
- **Cleanup**: Tests should clean up after themselves  
- **Timing**: Be generous with timeouts for CI environments
- **Logging**: Include helpful debug output for failures
- **Documentation**: Update this README for new test categories

### Quality Standards
- All tests must pass on supported platforms
- Performance tests should have reasonable benchmarks
- Error messages should be clear and actionable
- Tests should be maintainable and readable

## Integration with CI/CD

These tests are designed to integrate with GitHub Actions and other CI/CD systems:

```yaml
- name: Run Multi-Binary ATDD Tests
  run: |
    go test ./tests/acceptance/cli/... -v -short -timeout=10m
    
- name: Run Performance Tests
  run: |
    go test ./tests/acceptance/cli/performance_test.go -v -timeout=15m
```

The test suite provides comprehensive validation that the multi-binary restructure maintains all functionality while adding new capabilities.