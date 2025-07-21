# CLI Blueprint ATDD Test Suite

This directory contains comprehensive Acceptance Test-Driven Development (ATDD) tests for the CLI blueprints, ensuring they generate production-ready Go CLI applications that meet all user acceptance criteria across the two-tier approach (cli-simple and cli-standard).

## Overview

The ATDD test suite validates that the CLI blueprints:
- Generate complete, working CLI applications for both simple and standard tiers
- Support the progressive disclosure two-tier approach (8 vs 29 files)
- Support multiple logging libraries (slog, zap, logrus, zerolog)
- Include production-ready features (Docker, CI/CD, testing infrastructure)
- Follow CLI development best practices and patterns
- Provide comprehensive command-line interface functionality
- Implement proper error handling and validation
- Include performance optimizations and cross-platform support

## CLI Two-Tier Architecture

### Simple Tier (cli-simple)
- **Use Case**: Quick utilities, scripts, and learning projects
- **File Count**: 8 files (73% reduction from standard)
- **Structure**: Minimal, procedural approach
- **Logger**: slog only (standard library)
- **Features**: Essential CLI functionality (help, version, basic flags)
- **Target Users**: Beginners, prototyping, simple tools

### Standard Tier (cli-standard) 
- **Use Case**: Production CLI tools and enterprise applications
- **File Count**: ~29 files with layered architecture
- **Structure**: Internal packages with separation of concerns
- **Logger**: Multi-logger support (slog, zap, logrus, zerolog)
- **Features**: Advanced features (config management, interactive mode, CI/CD)
- **Target Users**: Production applications, team collaboration

## Test Structure

### BDD Test Files (NEW)

- **`cli_steps_test.go`** - Comprehensive BDD step definitions using godog with command execution testing
- **`cli_acceptance_test.go`** - High-level acceptance tests with tier-specific validation  
- **`features/cli.feature`** - Gherkin feature files defining 26+ user scenarios
- **`README.md`** - This comprehensive test documentation

### Legacy Test Files (EXISTING)

- **`standard_test.go`** - Standard CLI generation tests
- **`simple_test.go`** - Simple CLI tier specific tests
- **`cli_simple_atdd_test.go`** - Simple CLI ATDD validation
- **`enterprise_test.go`** - Enterprise feature testing
- **`runtime_integration_test.go`** - Runtime execution testing

### Test Categories

#### 1. Progressive Complexity Tests
- **Simple vs Standard Comparison**: Validates 8 vs 29 file difference
- **Feature Set Alignment**: Tests complexity-appropriate features
- **Migration Path Validation**: Ensures clear upgrade path from simple to standard
- **Performance Comparison**: Simple CLI should be faster and lighter

#### 2. Two-Tier Architecture Tests
- **Simple Tier Validation**: Minimal structure, procedural approach, slog only
- **Standard Tier Validation**: Layered architecture, internal packages, multi-logger
- **File Count Verification**: Exact file count matching (8 vs ~29)
- **Feature Distribution**: Essential vs advanced feature segregation

#### 3. Multi-Logger Integration Tests
- **slog**: Standard library structured logging (both tiers)
- **zap**: High-performance logging (standard tier only)
- **logrus**: Feature-rich logging (standard tier only)
- **zerolog**: Zero-allocation logging (standard tier only)

#### 4. Command Execution Tests
- **Help System**: Built-in help command validation
- **Version Command**: Version information display
- **Invalid Commands**: Error handling and helpful messages
- **Flag Processing**: Command-line flag parsing and validation
- **Exit Codes**: Appropriate status code handling

#### 5. Development Workflow Tests
- **Build Process**: Compilation and binary generation
- **Makefile Targets**: Essential development commands
- **CI/CD Integration**: GitHub Actions workflows
- **Testing Infrastructure**: Unit test generation and execution
- **Development Tools**: Linting, formatting, and quality tools

## BDD Test Scenarios (26 Comprehensive Scenarios)

### Feature: CLI Blueprint Generation
**As a Go developer**
**I want to generate modern, production-ready CLI applications** 
**So that I can quickly build command-line tools with industry best practices**

The Gherkin feature file (`features/cli.feature`) includes 26+ comprehensive scenarios covering:

#### Core Generation Scenarios
1. **Simple CLI Generation**: 8-file minimal structure for quick utilities
2. **Standard CLI Generation**: 29-file comprehensive structure for production
3. **Progressive Complexity**: Comparison and migration between tiers
4. **Multi-Logger Support**: slog, zap, logrus, zerolog integration

#### Framework and Integration Scenarios  
5. **Cobra Framework**: Command structure and subcommand support
6. **Configuration Management**: Viper integration and YAML configs
7. **Interactive Features**: Survey-based prompts and user input
8. **Output Formatting**: JSON, text, and quiet/verbose modes

#### Production Readiness Scenarios
9. **Error Handling**: Graceful error messages and exit codes
10. **Testing Infrastructure**: Unit tests and testify integration
11. **CI/CD Integration**: GitHub Actions workflows and automation
12. **Container Support**: Docker containerization and deployment

#### Advanced Feature Scenarios
13. **Cross-Platform**: Windows, macOS, Linux compatibility
14. **Performance**: Fast startup times and efficient execution
15. **Security**: Input validation and secure defaults
16. **Documentation**: Built-in help and comprehensive README

#### Example BDD Scenario (Simple CLI):
```gherkin
Scenario: Generate simple CLI application
  Given I want to create a simple CLI application for quick utilities
  When I run the command "go-starter new my-tool --type=cli --complexity=simple --no-git"
  Then the generation should succeed
  And the project should contain essential CLI components for simple tier
  And the generated code should compile successfully  
  And the CLI should execute basic commands
  And the project structure should be minimal with 8 files
```

#### Example BDD Scenario (Progressive Complexity):
```gherkin
Scenario: Progressive complexity demonstration
  Given I want to showcase progressive complexity
  When I compare simple and standard CLI blueprints
  Then simple CLI should have 8 files vs standard 29 files
  And simple CLI should focus on essential features only
  And standard CLI should include advanced capabilities
  And both should compile and execute successfully
  And migration path should be clear from simple to standard
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

### Prerequisites

- Go 1.21+
- Docker (for CLI container testing)
- Make (optional, for convenience commands)

### Command Examples

```bash
# Run all CLI ATDD tests (both BDD and legacy)
go test ./tests/acceptance/blueprints/cli/... -v

# Run new BDD tests with Gherkin scenarios
go test ./tests/acceptance/blueprints/cli/ -run TestCLIBDD -v

# Run comprehensive acceptance tests
go test ./tests/acceptance/blueprints/cli/ -run TestCLIAcceptance -v

# Run specific tier tests
go test ./tests/acceptance/blueprints/cli/ -run TestCLIAcceptance_SimpleTierGeneration -v
go test ./tests/acceptance/blueprints/cli/ -run TestCLIAcceptance_StandardTierGeneration -v

# Run progressive complexity comparison
go test ./tests/acceptance/blueprints/cli/ -run TestCLIAcceptance_ProgressiveComplexityComparison -v

# Run multi-logger tests
go test ./tests/acceptance/blueprints/cli/ -run TestCLIAcceptance_MultiLoggerSupport -v

# Run performance tests
go test ./tests/acceptance/blueprints/cli/ -run TestCLIAcceptance_PerformanceAndEfficiency -v

# Run legacy CLI tests
go test ./tests/acceptance/blueprints/cli/ -run TestStandard_CLI -v
go test ./tests/acceptance/blueprints/cli/ -run TestSimple_CLI -v

# Skip slow tests during development
go test ./tests/acceptance/blueprints/cli/ -short

# Run tests in parallel (faster execution)
go test ./tests/acceptance/blueprints/cli/ -v -parallel 4

# Run with extended timeout for comprehensive tests
go test ./tests/acceptance/blueprints/cli/ -v -timeout 15m
```

### Environment Variables

- `TEST_TIMEOUT` - Override default test timeout (default: 30s)
- `VERBOSE_OUTPUT` - Enable verbose test output
- `KEEP_TEST_ARTIFACTS` - Don't cleanup generated CLI projects
- `SKIP_DOCKER_TESTS` - Skip container-based tests

### CI/CD Integration

The CLI ATDD tests are integrated into the CI pipeline:

```yaml
- name: Run CLI ATDD Tests
  run: |
    go test ./tests/acceptance/blueprints/cli/... -v -timeout 15m
    go test ./tests/acceptance/blueprints/cli/ -run TestCLIBDD -v
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

## Validation Methods

### CLI Test Suite Methods

#### CLITestContext (BDD)
- **`generateCLIProject()`**: Generates CLI projects with specified configuration
- **`buildAndTestCLI()`**: Compiles and executes CLI for functional testing
- **`compileCLIProject()`**: Validates project compilation
- **`countFiles()`**: Verifies file count matches tier expectations
- **`checkFileExists()`** & **`checkFileContains()`**: File structure validation

#### CLIAcceptanceTestSuite
- **`generateCLIProject()`**: High-level project generation with validation
- **`compileCLIProject()`**: Compilation verification with timing
- **`testCLIExecution()`**: Command execution testing with output capture
- **Performance Tracking**: Build time, execution time, file count metrics

### Legacy Validator Methods
- **`ValidateCobraSetup()`**: Cobra framework integration validation
- **`ValidateLoggerFunctionality()`**: Multi-logger testing and configuration
- **`ValidateViperConfiguration()`**: Configuration management validation
- **`ValidateDockerSupport()`**: Container deployment validation

## Test Scenarios Coverage

### Two-Tier Architecture Validation

#### Simple Tier User Stories
1. **As a beginner developer, I want minimal CLI structure**
   - ✅ 8 files total (main.go, go.mod, README.md, Makefile, config.go, cmd/, .gitignore)
   - ✅ slog logging only (standard library)
   - ✅ Basic Cobra framework integration
   - ✅ Essential commands (help, version)
   - ✅ Fast compilation and execution

#### Standard Tier User Stories  
2. **As a professional developer, I want production-ready CLI tools**
   - ✅ ~29 files with layered architecture
   - ✅ Multi-logger support (slog, zap, logrus, zerolog)
   - ✅ Advanced configuration management (Viper, YAML)
   - ✅ CI/CD integration (GitHub Actions)
   - ✅ Container support (Docker, Dockerfile)
   - ✅ Interactive features (survey prompts)
   - ✅ Testing infrastructure (testify, unit tests)

#### Progressive Complexity User Stories
3. **As a developer, I want clear progression from simple to advanced**
   - ✅ 73% file reduction in simple tier (8 vs 29 files)
   - ✅ Feature alignment with complexity level
   - ✅ Migration documentation and upgrade paths
   - ✅ Performance comparison (simple is faster)

### Framework and Integration User Stories

4. **As a CLI developer, I want robust command-line interfaces**
   - ✅ Cobra framework with subcommands
   - ✅ Built-in help system and usage information
   - ✅ Version command implementation
   - ✅ Flag processing and argument validation
   - ✅ Error handling with appropriate exit codes

5. **As a developer, I want flexible logging options**
   - ✅ slog (standard library, both tiers)
   - ✅ zap (high performance, standard only)
   - ✅ logrus (feature-rich, standard only) 
   - ✅ zerolog (zero allocation, standard only)

### Production Readiness User Stories

6. **As a DevOps engineer, I want production-ready CLI tools**
   - ✅ Docker containerization with multi-stage builds
   - ✅ CI/CD pipelines with automated testing
   - ✅ Cross-platform compatibility
   - ✅ Performance optimization and fast startup
   - ✅ Comprehensive documentation and help

7. **As a developer, I want comprehensive testing**
   - ✅ Unit tests with testify framework
   - ✅ Command execution validation
   - ✅ Error scenario testing
   - ✅ Performance benchmarking

## Quality Gates

Each test validates specific quality gates:

- **Functionality**: Generated CLIs compile and execute correctly
- **Architecture**: Follows tier-appropriate architectural patterns  
- **Performance**: Simple CLI < 100ms startup, reasonable build times
- **Usability**: Comprehensive help system and error messages
- **Maintainability**: Clean code structure and comprehensive testing
- **Portability**: Cross-platform compatibility and containerization
- **Documentation**: Built-in help and comprehensive README files

## Expected Outcomes

All tests should:
1. **Generate valid CLI projects** that compile without errors for both tiers
2. **Produce working command-line applications** with proper Cobra structure
3. **Include appropriate files** matching tier complexity (8 vs 29 files)
4. **Execute successfully** with help, version, and basic command functionality
5. **Support multi-logger configurations** with conditional generation
6. **Provide production features** (Docker, CI/CD, testing) for standard tier
7. **Demonstrate clear progression** from simple to standard complexity
8. **Pass all internal tests** in generated projects
9. **Achieve performance targets** (< 100ms startup for simple tier)
10. **Include comprehensive documentation** and development workflow

## Troubleshooting

### Common Issues and Solutions

1. **CLI generation fails**
   ```
   Error: Project generation failed
   Solution: Check blueprint templates exist and are valid
   Debug: Use KEEP_TEST_ARTIFACTS=1 to inspect generated code
   ```

2. **Compilation errors in generated CLI**
   ```
   Error: Generated project doesn't compile
   Solution: Verify go.mod dependencies and template syntax
   Check: Logger imports and framework dependencies
   ```

3. **File count doesn't match tier expectations**
   ```
   Error: Expected 8 files, got 10 files
   Solution: Check for hidden files or blueprint modifications
   Verify: Template conditional logic and file generation rules
   ```

4. **CLI execution fails**
   ```
   Error: CLI binary doesn't execute properly
   Solution: Check Cobra framework setup and command registration
   Debug: Run CLI manually with --help and --version
   ```

5. **Performance tests fail**
   ```
   Error: CLI startup time exceeds 100ms
   Solution: Profile the CLI binary and check for slow imports
   Optimize: Remove unnecessary dependencies from simple tier
   ```

### Debug Mode

Enable comprehensive debugging:

```bash
# Verbose test execution with artifact preservation
VERBOSE_OUTPUT=1 KEEP_TEST_ARTIFACTS=1 go test ./tests/acceptance/blueprints/cli/ -v -run TestCLIAcceptance_SimpleTierGeneration

# Analyze generated CLI projects
KEEP_TEST_ARTIFACTS=1 go test ./tests/acceptance/blueprints/cli/ -v
ls -la /tmp/cli-bdd-*/      # Inspect generated CLIs
./tmp/cli-bdd-*/test-cli-simple/test-cli-simple --help  # Test CLI manually
```

### Performance Analysis

```bash
# Measure CLI performance
time ./generated-cli --version
time ./generated-cli --help

# Compare tiers
time ./simple-cli --version    # Should be < 50ms
time ./standard-cli --version  # Should be < 100ms
```

## Contributing

When contributing to the CLI ATDD test suite:

### Guidelines
1. **Follow BDD patterns**: Use Given-When-Then structure in tests
2. **Test both tiers**: Ensure simple and standard tier coverage
3. **Validate tier characteristics**: File counts, features, performance
4. **Include command testing**: Execute generated CLIs, don't just compile
5. **Test multi-logger support**: Verify conditional logger generation
6. **Update documentation**: Add new scenarios to feature files
7. **Consider migration paths**: Test upgrade from simple to standard

### Code Quality
- Use descriptive test names indicating tier and feature tested
- Include performance assertions for CLI execution speed
- Validate both positive and negative command scenarios
- Test cross-platform considerations where applicable
- Ensure thread-safety for parallel test execution

The CLI ATDD test suite ensures that both simple and standard CLI blueprints generate high-quality, functional command-line applications that meet user expectations across the progressive disclosure spectrum. The comprehensive BDD approach provides business-readable documentation while maintaining technical rigor.