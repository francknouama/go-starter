# Blueprint Integration Tests

This directory contains integration tests for all blueprint types, validating that generated code compiles correctly and produces functional applications.

## Directory Structure

### 🌐 **web-api/** - Web API Blueprint Tests
Tests for all web API architecture patterns:
- `standard_test.go` - Standard layered architecture
- `clean_test.go` - Clean Architecture pattern  
- `ddd_test.go` - Domain-Driven Design pattern
- `hexagonal_test.go` - Hexagonal Architecture pattern
- `database_integration_test.go` - Database integration tests
- `error_handling_test.go` - Error handling validation
- `runtime_integration_test.go` - Runtime compilation tests

### 🖥️ **cli/** - CLI Blueprint Tests
Tests for CLI application blueprints:
- `simple_test.go` - Simple CLI blueprint (8 files, minimal complexity)
- `standard_test.go` - Standard CLI blueprint (29 files, full features)

### 🔧 **shared/** - Shared Blueprint Utilities
Common utilities for blueprint testing:
- `blueprint_helpers.go` - Shared test helpers and validation functions

## Test Methodology

### Blueprint Validation Process
1. **Generation**: Generate project using specific blueprint configuration
2. **Compilation**: Verify generated code compiles without errors
3. **Structure**: Validate file structure matches blueprint expectations  
4. **Architecture**: Test architecture pattern compliance
5. **Functionality**: Basic smoke tests for key functionality
6. **Cleanup**: Remove generated test projects

### Architecture Pattern Testing
Each web API architecture is tested for:
- **Dependency Direction**: Proper dependency flow
- **Layer Separation**: Clear boundaries between layers
- **Interface Compliance**: Proper interface definitions
- **Framework Integration**: Multi-framework support (Gin, Echo, Fiber, Chi)
- **Logger Integration**: Multi-logger support (slog, zap, logrus, zerolog)

### CLI Blueprint Testing  
CLI blueprints are tested for:
- **File Count**: Correct number of generated files
- **Structure Validation**: Expected directory structure
- **Compilation**: Generated code compiles successfully
- **Command Functionality**: Basic command execution
- **Two-Tier Approach**: Simple vs Standard complexity validation

## Running Blueprint Tests

```bash
# Run all blueprint tests
go test ./tests/integration/blueprints/... -v

# Run specific blueprint type
go test ./tests/integration/blueprints/web-api/... -v
go test ./tests/integration/blueprints/cli/... -v

# Run with architecture validation
go test ./tests/integration/blueprints/web-api/ -v -args -validate-architecture

# Run with performance profiling
go test ./tests/integration/blueprints/... -v -cpuprofile=cpu.prof -memprofile=mem.prof
```

## Adding New Blueprint Tests

When adding tests for new blueprints:

1. **Create Blueprint Directory**: `mkdir blueprints/{new-blueprint}`
2. **Follow Naming Convention**: `{blueprint-name}_test.go`
3. **Use Standard Test Structure**:
   ```go
   func TestBlueprint_{Name}_{Feature}(t *testing.T) {
       // Test implementation
   }
   ```
4. **Include Architecture Validation**: Validate generated structure
5. **Test Multiple Configurations**: Test with different options
6. **Add Compilation Verification**: Ensure generated code compiles
7. **Document Complex Scenarios**: Add comments for complex tests

## Test Configuration

### Environment Variables
- `INTEGRATION_TESTS_KEEP_PROJECTS=1` - Keep generated projects for debugging
- `INTEGRATION_TESTS_TIMEOUT=300s` - Set custom timeout for tests
- `INTEGRATION_TESTS_PARALLEL=4` - Set parallel test execution limit

### Test Flags
- `-validate-architecture` - Enable deep architecture validation
- `-keep-projects` - Keep generated projects after tests
- `-verbose-generation` - Enable verbose generation logging