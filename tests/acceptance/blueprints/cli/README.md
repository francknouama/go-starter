# CLI Blueprint Acceptance Tests

This directory contains comprehensive Acceptance Test-Driven Development (ATDD) tests for the CLI blueprint. These tests validate that generated CLI applications work correctly across different configurations and scenarios.

## Test Overview

The ATDD tests are written in Gherkin-style BDD format with Given-When-Then scenarios for business readability and follow the established patterns from the Web API blueprint tests.

## Test Structure

### `standard_test.go`
- **TestStandard_CLI_BasicGeneration_WithCobra**: Core CLI generation with Cobra framework
- **TestStandard_CLI_WithDifferentLoggers**: Logger compatibility testing (slog, zap, logrus, zerolog)
- **TestStandard_CLI_ConfigurationSupport**: Viper configuration testing
- **TestStandard_CLI_DockerSupport**: Docker containerization support
- **TestStandard_CLI_TestSupport**: Test infrastructure validation

## Test Scenarios

### Feature: Standard CLI Blueprint
**As a developer**  
**I want to generate a CLI application project**  
**So that I can quickly build command-line tools**

#### Scenario 1: Basic CLI Generation
```gherkin
Given I want a standard CLI application
When I generate with framework "cobra"
Then the project should include Cobra command setup
And the project should have root and version commands
And the project should compile and run successfully
And CLI should show help and version output
```

#### Scenario 2: Logger Configuration
```gherkin
Given I want a CLI application with configurable logging
When I generate with different loggers (slog, zap, logrus, zerolog)
Then the project should include the selected logger
And the project should compile successfully
And logging should work as expected
```

#### Scenario 3: Configuration Support
```gherkin
Given I want a CLI application with configuration
When I generate the project
Then the project should include Viper configuration
And the project should have config file support
And the project should compile and run successfully
```

#### Scenario 4: Docker Support
```gherkin
Given I want a CLI application with Docker
When I generate the project
Then the project should include Dockerfile
And the project should have Makefile with Docker targets
And the project should compile and run successfully
```

#### Scenario 5: Test Support
```gherkin
Given I want a CLI application with tests
When I generate the project
Then the project should include test files
And the project should use testify for assertions
And the tests should run successfully
```

## Validations Performed

### CLIValidator Methods
- **ValidateCobraSetup**: Verifies Cobra framework integration
- **ValidateRootCommand**: Ensures root command implementation
- **ValidateVersionCommand**: Validates version command functionality
- **ValidateCompilation**: Confirms project compiles successfully
- **ValidateHelpOutput**: Tests CLI help command output
- **ValidateVersionOutput**: Tests CLI version command output
- **ValidateLogger**: Validates logger implementation and configuration
- **ValidateLoggerFunctionality**: Tests logger functionality
- **ValidateViperConfiguration**: Validates Viper configuration setup
- **ValidateConfigFileSupport**: Tests configuration file support
- **ValidateDockerSupport**: Validates Docker integration
- **ValidateMakefileTargets**: Tests Makefile build targets
- **ValidateTestFiles**: Ensures test files exist and are properly structured
- **ValidateTestifyUsage**: Validates testify framework integration
- **ValidateTestExecution**: Tests that all tests run successfully

## File Structure Validation

The tests validate the following generated file structure:
```
project/
├── main.go                    # CLI entry point
├── go.mod                     # Go module definition
├── Makefile                   # Build targets
├── README.md                  # Project documentation
├── Dockerfile                 # Docker configuration
├── cmd/
│   ├── root.go               # Root command implementation
│   ├── root_test.go          # Root command tests
│   └── version.go            # Version command
├── internal/
│   ├── config/
│   │   ├── config.go         # Configuration management
│   │   └── config_test.go    # Configuration tests
│   └── logger/
│       ├── interface.go      # Logger interface
│       ├── factory.go        # Logger factory
│       └── [logger].go       # Logger implementation
├── configs/
│   └── config.yaml           # Configuration file
└── .github/
    └── workflows/
        ├── ci.yml            # CI/CD pipeline
        └── release.yml       # Release automation
```

## Running the Tests

To run the CLI blueprint acceptance tests:

```bash
# Run all CLI acceptance tests
go test -v ./tests/acceptance/blueprints/cli/...

# Run specific test scenarios
go test -v ./tests/acceptance/blueprints/cli/ -run TestStandard_CLI_BasicGeneration

# Run with timeout for longer tests
go test -v -timeout=10m ./tests/acceptance/blueprints/cli/...
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
1. **Generate valid projects** that compile without errors
2. **Produce working CLI applications** with proper command structure
3. **Include all required files** for the selected configuration
4. **Pass all internal tests** in the generated project
5. **Support all logger types** with proper conditional generation
6. **Provide Docker support** with valid Dockerfile and Makefile targets
7. **Include comprehensive test coverage** with testify framework integration

## Troubleshooting

If tests fail:
1. Check that all blueprint template files exist and are valid
2. Verify that the generator service is working correctly
3. Ensure all dependencies are properly configured in template.yaml
4. Check that post-generation hooks (go mod tidy, go fmt) are working
5. Validate that the CI/CD environment has all required tools (Go, Docker, etc.)