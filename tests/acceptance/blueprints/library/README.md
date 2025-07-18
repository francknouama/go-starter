# Library Blueprint Acceptance Tests

This directory contains comprehensive Acceptance Test-Driven Development (ATDD) tests for the Library blueprint. These tests validate that generated Go libraries work correctly across different configurations and provide proper public APIs.

## Test Overview

The ATDD tests are written in Gherkin-style BDD format with Given-When-Then scenarios for business readability and follow the established patterns from the Web API blueprint tests.

## Test Structure

### `standard_test.go`
- **TestStandard_Library_BasicGeneration**: Core library generation with public API
- **TestStandard_Library_WithDifferentLoggers**: Logger compatibility testing (slog, zap, logrus, zerolog)
- **TestStandard_Library_Documentation**: Documentation and Go doc support
- **TestStandard_Library_ExampleUsage**: Example usage and demonstration code
- **TestStandard_Library_TestSupport**: Test infrastructure validation
- **TestStandard_Library_PublicAPI**: Public API design and naming conventions

## Test Scenarios

### Feature: Standard Library Blueprint
**As a developer**  
**I want to generate a Go library project**  
**So that I can quickly build reusable packages**

#### Scenario 1: Basic Library Generation
```gherkin
Given I want a standard Go library
When I generate the project
Then the project should include library structure
And the project should have public API
And the project should compile and run successfully
And the project should have working examples
```

#### Scenario 2: Logger Configuration
```gherkin
Given I want a library with configurable logging
When I generate with different loggers (slog, zap, logrus, zerolog)
Then the project should include the selected logger
And the project should compile successfully
And logging should work as expected
```

#### Scenario 3: Documentation Support
```gherkin
Given I want a library with documentation
When I generate the project
Then the project should include Go documentation
And the project should have README with examples
And the project should have doc.go file
```

#### Scenario 4: Example Usage
```gherkin
Given I want a library with examples
When I generate the project
Then the project should include basic examples
And the project should include advanced examples
And the examples should compile and run successfully
```

#### Scenario 5: Test Support
```gherkin
Given I want a library with comprehensive tests
When I generate the project
Then the project should include test files
And the project should include example tests
And the project should use testify for assertions
And the tests should run successfully
```

#### Scenario 6: Public API Design
```gherkin
Given I want a library with clean public API
When I generate the project
Then the project should export main functionality
And the project should hide internal implementation
And the project should have proper Go naming conventions
```

## Validations Performed

### LibraryValidator Methods
- **ValidateLibraryStructure**: Verifies library file structure
- **ValidatePublicAPI**: Ensures proper public API design
- **ValidateCompilation**: Confirms project compiles successfully
- **ValidateExamples**: Tests example code compilation and execution
- **ValidateLogger**: Validates logger implementation and configuration
- **ValidateLoggerFunctionality**: Tests logger functionality
- **ValidateDocumentation**: Validates Go documentation setup
- **ValidateREADME**: Tests README content and examples
- **ValidateDocFile**: Validates doc.go file
- **ValidateBasicExample**: Tests basic usage example
- **ValidateAdvancedExample**: Tests advanced usage example
- **ValidateExampleExecution**: Tests example compilation
- **ValidateTestFiles**: Ensures test files exist and are properly structured
- **ValidateExampleTests**: Validates example tests (ExampleXxx functions)
- **ValidateTestifyUsage**: Validates testify framework integration
- **ValidateTestExecution**: Tests that all tests run successfully
- **ValidateInternalPackages**: Validates internal package structure
- **ValidateNamingConventions**: Tests Go naming conventions

## File Structure Validation

The tests validate the following generated file structure:
```
project/
├── project-name.go           # Main library file with public API
├── project-name_test.go      # Main library tests
├── go.mod                    # Go module definition
├── Makefile                  # Build targets
├── README.md                 # Project documentation with examples
├── doc.go                    # Package documentation
├── examples_test.go          # Example tests (ExampleXxx functions)
├── examples/
│   ├── basic/
│   │   └── main.go          # Basic usage example
│   └── advanced/
│       └── main.go          # Advanced usage example
├── internal/
│   └── logger/
│       └── logger.go        # Internal logger (minimal for libraries)
├── .gitignore               # Git ignore file
└── .github/
    └── workflows/
        └── ci.yml           # CI/CD pipeline
```

## Library Design Principles

The tests validate adherence to Go library best practices:

### Public API Design
- **Exported functions** start with uppercase letters
- **Unexported functions** start with lowercase letters
- **Clear and consistent naming** following Go conventions
- **Minimal surface area** exposing only necessary functionality
- **Proper error handling** with meaningful error messages

### Internal Implementation
- **Internal packages** for implementation details
- **Minimal dependencies** to reduce library footprint
- **Optional logging** for debugging without forcing dependencies on consumers
- **Clean separation** between public API and internal implementation

### Documentation
- **Package documentation** in doc.go
- **Function documentation** following Go doc conventions
- **Usage examples** in README.md
- **Working examples** in examples/ directory
- **Example tests** (ExampleXxx functions) for go doc

### Testing
- **Comprehensive test coverage** for public API
- **Example tests** for documentation
- **Testify framework** for assertions
- **Test data** organization and management

## Running the Tests

To run the Library blueprint acceptance tests:

```bash
# Run all Library acceptance tests
go test -v ./tests/acceptance/blueprints/library/...

# Run specific test scenarios
go test -v ./tests/acceptance/blueprints/library/ -run TestStandard_Library_BasicGeneration

# Run with timeout for longer tests
go test -v -timeout=10m ./tests/acceptance/blueprints/library/...
```

## Integration with CI/CD

These tests are integrated into the CI/CD pipeline and run automatically on:
- Pull requests to main/develop branches
- Pushes to main/develop branches
- Part of the dedicated ATDD test job with 25-minute timeout

## Test Dependencies

The tests use the following testing frameworks and utilities:
- **testify/assert**: Assertion library
- **helpers package**: Shared testing utilities
- **types package**: Project configuration types
- **generator package**: Project generation service

## Expected Outcomes

All tests should:
1. **Generate valid libraries** that compile without errors
2. **Produce clean public APIs** following Go conventions
3. **Include comprehensive documentation** with examples
4. **Provide working examples** for library usage
5. **Pass all internal tests** in the generated project
6. **Support all logger types** with proper conditional generation
7. **Hide internal implementation** details from consumers
8. **Follow Go best practices** for library design
9. **Include example tests** for go doc integration
10. **Maintain minimal dependencies** to reduce library footprint

## Troubleshooting

If tests fail:
1. Check that all blueprint template files exist and are valid
2. Verify that the generator service is working correctly
3. Ensure all dependencies are properly configured in template.yaml
4. Check that post-generation hooks (go mod tidy, go fmt) are working
5. Validate that public API functions are properly exported
6. Ensure example code compiles and runs correctly
7. Check that internal packages are properly structured
8. Verify that documentation follows Go doc conventions
9. Ensure example tests (ExampleXxx functions) work correctly
10. Validate that the CI/CD environment has all required tools (Go, etc.)

## Library Quality Checklist

Generated libraries should meet these quality standards:
- [ ] **Compiles without errors** with `go build`
- [ ] **All tests pass** with `go test ./...`
- [ ] **Examples compile** in examples/ directory
- [ ] **Documentation is complete** with go doc
- [ ] **Public API is clean** and follows Go conventions
- [ ] **Internal implementation is hidden** from consumers
- [ ] **Minimal dependencies** to reduce library footprint
- [ ] **Example tests work** for go doc integration
- [ ] **README includes usage examples**
- [ ] **License and author information** are properly configured