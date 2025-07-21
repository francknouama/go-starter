Feature: CLI Blueprint Generation
  As a Go developer
  I want to generate modern, production-ready CLI applications
  So that I can quickly build command-line tools with industry best practices

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate simple CLI application
    Given I want to create a simple CLI application for quick utilities
    When I run the command "go-starter new my-tool --type=cli --complexity=simple --no-git"
    Then the generation should succeed
    And the project should contain essential CLI components for simple tier
    And the generated code should compile successfully
    And the CLI should execute basic commands
    And the project structure should be minimal with 8 files

  Scenario: Generate standard CLI application
    Given I want to create a production-ready CLI application
    When I run the command "go-starter new my-cli --type=cli --complexity=standard --no-git"
    Then the generation should succeed
    And the project should contain all essential CLI components for standard tier
    And the generated code should compile successfully
    And the CLI should support advanced features
    And the project structure should be comprehensive with 29 files

  Scenario: CLI with different complexity levels
    Given I want to create a CLI with specific complexity
    When I generate a CLI with complexity "<complexity>"
    Then the project should follow "<tier>" tier patterns
    And the file count should match the complexity level
    And the feature set should align with the complexity
    And the CLI should compile and execute correctly

    Examples:
      | complexity | tier     |
      | simple     | simple   |
      | standard   | standard |

  Scenario: Simple CLI validation
    Given I want to build quick command-line utilities
    When I generate a simple CLI application
    Then the CLI should have minimal structure
    And the CLI should use slog for logging
    And the CLI should support basic flags (help, version, quiet)
    And the CLI should have Cobra framework integration
    And the configuration should be minimal
    And the commands should be straightforward

  Scenario: Standard CLI validation  
    Given I want to build production CLI tools
    When I generate a standard CLI application
    Then the CLI should have layered architecture
    And the CLI should support multiple loggers (slog, zap, logrus, zerolog)
    And the CLI should have comprehensive flag support
    And the CLI should include configuration management
    And the CLI should support interactive mode
    And the CLI should have testing infrastructure
    And the CLI should include CI/CD integration

  Scenario: Progressive complexity demonstration
    Given I want to showcase progressive complexity
    When I compare simple and standard CLI blueprints
    Then simple CLI should have 8 files vs standard 29 files
    And simple CLI should focus on essential features only
    And standard CLI should include advanced capabilities
    And both should compile and execute successfully
    And migration path should be clear from simple to standard

  Scenario: CLI with different logger types
    Given I want to use different logging libraries in CLI
    When I generate a CLI with logger "<logger>"
    Then the CLI should use the "<logger>" logging library
    And the logger should be properly configured
    And the CLI should support structured logging
    And the log levels should be configurable
    And the CLI should compile with the logger dependency

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  Scenario: CLI framework integration
    Given I want to ensure CLI framework works properly
    When I generate a CLI application
    Then the CLI should use Cobra framework
    And the CLI should support subcommands
    And the CLI should have built-in help system
    And the CLI should support command completion
    And the CLI should handle flags and arguments correctly

  Scenario: Command execution and behavior
    Given I have generated a CLI application
    When I execute CLI commands
    Then the help command should display usage information
    And the version command should show version details
    And invalid commands should show helpful error messages
    And the CLI should exit with appropriate status codes
    And command output should be properly formatted

  Scenario: Configuration management
    Given I want configurable CLI applications
    When I generate a CLI with configuration support
    Then the CLI should support configuration files
    And the CLI should support environment variables
    And the CLI should have configuration precedence
    And the CLI should validate configuration values
    And the CLI should provide configuration examples

  Scenario: Output formatting and verbosity
    Given I want flexible CLI output
    When I generate a CLI application
    Then the CLI should support multiple output formats
    And the CLI should support quiet mode
    And the CLI should support verbose mode
    And the CLI should handle JSON output
    And the CLI should format output appropriately

  Scenario: Error handling and validation
    Given I want robust CLI applications
    When I generate a CLI with error handling
    Then the CLI should handle invalid arguments gracefully
    And the CLI should provide clear error messages
    And the CLI should validate input parameters
    And the CLI should handle system errors properly
    And the CLI should use appropriate exit codes

  Scenario: Interactive mode capabilities
    Given I want interactive CLI applications
    When I generate a standard CLI with interactive features
    Then the CLI should support interactive prompts
    And the CLI should validate user input interactively
    And the CLI should provide selection menus
    And the CLI should handle user cancellation
    And the CLI should maintain session state

  Scenario: Testing infrastructure for CLI
    Given I want well-tested CLI applications
    When I generate a CLI with testing support
    Then the CLI should include unit tests
    And the CLI should have command testing examples
    And the CLI should support integration testing
    And the CLI should include test coverage reporting
    And the tests should validate CLI behavior

  Scenario: Development workflow and tooling
    Given I want efficient CLI development
    When I generate a CLI application
    Then the project should include Makefile with essential targets
    And the project should have development scripts
    And the project should support hot reloading for development
    And the project should include linting configuration
    And the project should have pre-commit hooks

  Scenario: Distribution and packaging
    Given I want to distribute CLI applications
    When I generate a CLI for distribution
    Then the project should include build scripts
    And the project should support cross-compilation
    And the project should have release automation
    And the project should include installation instructions
    And the project should support package managers

  Scenario: CI/CD integration for CLI
    Given I want automated CLI workflows
    When I generate a CLI with CI/CD support
    Then the project should include GitHub Actions workflows
    And the CI should run tests and linting
    And the CI should build for multiple platforms
    And the CI should create release artifacts
    And the deployment should support multiple environments

  Scenario: Container support for CLI
    Given I want containerized CLI applications
    When I generate a CLI with container support
    Then the project should include Dockerfile
    And the container should be optimized for CLI usage
    And the container should support entrypoint configuration
    And the container should handle signals properly
    And the container should support volume mounting

  Scenario: Security and hardening
    Given I want secure CLI applications
    When I generate a CLI with security features
    Then the CLI should validate all inputs
    And the CLI should handle sensitive data securely
    And the CLI should avoid exposing secrets in logs
    And the CLI should use secure defaults
    And the CLI should follow security best practices

  Scenario: Performance and efficiency
    Given I want high-performance CLI applications
    When I generate a CLI with performance optimizations
    Then the CLI should have fast startup time
    And the CLI should handle large datasets efficiently
    And the CLI should support concurrent operations
    And the CLI should minimize memory usage
    And the CLI should provide performance metrics

  Scenario: Documentation and help system
    Given I want well-documented CLI applications
    When I generate a CLI with comprehensive documentation
    Then the CLI should have built-in help for all commands
    And the CLI should include usage examples
    And the CLI should have man page generation
    And the CLI should support help formatting
    And the CLI should provide troubleshooting guidance

  Scenario: Completion and shell integration
    Given I want CLI applications with shell integration
    When I generate a CLI with completion support
    Then the CLI should support bash completion
    And the CLI should support zsh completion
    And the CLI should support fish completion
    And the CLI should support PowerShell completion
    And the completion should be context-aware

  Scenario: Migration from simple to standard
    Given I have a simple CLI that needs more features
    When I want to migrate to standard tier
    Then the migration path should be documented
    And the core functionality should be preserved
    And the configuration should be upgradeable
    And the commands should remain compatible
    And the migration should be reversible

  Scenario: Cross-platform compatibility
    Given I want CLI applications that work everywhere
    When I generate a CLI for cross-platform use
    Then the CLI should work on Windows, macOS, and Linux
    And the CLI should handle path separators correctly
    And the CLI should respect platform conventions
    And the CLI should handle different terminal capabilities
    And the CLI should support platform-specific features

  Scenario: Plugin and extension system
    Given I want extensible CLI applications
    When I generate a CLI with plugin support
    Then the CLI should support plugin architecture
    And the CLI should have plugin discovery
    And the CLI should handle plugin lifecycle
    And the CLI should validate plugin compatibility
    And the CLI should provide plugin APIs

  Scenario: Monitoring and observability
    Given I want observable CLI applications
    When I generate a CLI with monitoring features
    Then the CLI should support metrics collection
    And the CLI should provide execution tracing
    And the CLI should log performance data
    And the CLI should support health checks
    And the CLI should integrate with monitoring systems

  Scenario: Internationalization and localization
    Given I want CLI applications for global use
    When I generate a CLI with i18n support
    Then the CLI should support multiple languages
    And the CLI should handle locale-specific formatting
    And the CLI should support right-to-left text
    And the CLI should handle currency and dates
    And the CLI should provide translation infrastructure

  Scenario: Advanced command patterns
    Given I want sophisticated CLI interfaces
    When I generate a CLI with advanced patterns
    Then the CLI should support command chaining
    And the CLI should handle command pipelines
    And the CLI should support command composition
    And the CLI should validate command relationships
    And the CLI should provide command history

  # Simple CLI Specific Scenarios (Issue #149)
  # Based on CLI audit findings addressing complexity and learning curve

  Scenario: Simple CLI with minimal complexity (Issue #149)
    Given I want a simple CLI application with minimal complexity
    When I generate a CLI with architecture "simple"
    Then the project should have less than 10 files
    And the project should use only essential dependencies
    And the project should have complexity level 3/10
    And the project should compile and run successfully
    And the learning curve should be minimal for beginners

  Scenario: Simple CLI includes only essential features
    Given I want a CLI with essential features only
    When I generate a simple CLI project
    Then the project should have help command support
    And the project should have version command support
    And the project should have quiet flag support
    And the project should have output format support
    And the project should NOT have complex business logic
    And the project should NOT have factory patterns
    And the project should NOT have create/update/delete commands

  Scenario: Simple CLI uses slog (standard library) only
    Given I want a CLI with simple logging
    When I generate a simple CLI project with slog logger
    Then the project should use slog from standard library
    And the project should NOT have logger factory
    And the project should NOT have logger interfaces
    And the project should have direct slog usage
    And logging should work correctly
    And the project should NOT have external logger dependencies

  Scenario: Simple CLI supports shell completion
    Given I want a CLI with shell completion
    When I generate a simple CLI project
    Then the project should support bash completion
    And the project should support zsh completion
    And the project should support fish completion
    And the project should support powershell completion
    And completion should work without additional setup
    And completion configuration should be minimal

  Scenario: Simple CLI has minimal configuration
    Given I want a CLI with simple configuration
    When I generate a simple CLI project
    Then the project should have environment variable support
    And the project should have simple validation
    And the project should NOT use Viper
    And the project should NOT have complex config files
    And configuration should be straightforward
    And the project should use standard library for configuration

  Scenario: Simple CLI has beginner-friendly structure
    Given I am a beginner Go developer
    When I generate a simple CLI project
    Then the project structure should be intuitive
    And each file should have a clear purpose
    And the code should be easy to understand
    And the documentation should be comprehensive
    And the learning curve should be minimal
    And main.go should be simple and under 30 lines
    And each file should have clear comments explaining purpose

  Scenario: Simple CLI performs better than standard CLI
    Given I want to compare simple vs standard CLI
    When I generate both simple and standard CLI architectures
    Then simple CLI should have fewer files than standard CLI
    And simple CLI should have fewer dependencies than standard CLI
    And simple CLI should compile faster than standard CLI
    And simple CLI should be easier to maintain
    And simple CLI should have smaller binary size
    And simple CLI should start faster

  Scenario: Simple CLI complexity validation
    Given I want to ensure CLI simplicity
    When I generate a simple CLI project
    Then the project should have complexity level 3/10
    And the project should have learning curve 3/10
    And the project should NOT have factory patterns
    And the project should NOT have complex interfaces
    And the project should NOT have deep directory nesting
    And the maximum directory depth should be 3 levels
    And each file should have at most 1 interface definition

  Scenario: Simple CLI essential dependencies only
    Given I want minimal external dependencies
    When I generate a simple CLI project
    Then the project should only have cobra as direct dependency
    And the project should NOT have Viper dependency
    And the project should NOT have external logger dependencies
    And the project should use standard library where possible
    And the go.mod should have minimal dependencies
    And dependency count should be less than standard CLI

  Scenario: Simple CLI file structure validation
    Given I want an intuitive project structure
    When I generate a simple CLI project
    Then the project should have clear main.go
    And the project should have cmd directory for commands
    And the project should have simple config.go
    And the project should NOT have internal directory structure
    And the project should NOT have logger factory directory
    And total file count should be under 10 files