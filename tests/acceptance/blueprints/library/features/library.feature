Feature: Library Blueprint Generation
  As a Go developer
  I want to generate professional, reusable Go libraries
  So that I can create well-structured packages with proper documentation, testing, and distribution

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate basic library-standard project
    Given I want to create a reusable Go library
    When I run the command "go-starter new my-library --type=library-standard --module=github.com/example/my-library --no-git"
    Then the generation should succeed
    And the project should contain all essential library components
    And the generated code should compile successfully
    And the library should follow Go package best practices
    And the library should include comprehensive documentation

  Scenario: Library with different logging implementations
    Given I want to create libraries with various logging options
    When I generate a library with logger "<logger>"
    Then the project should support the "<logger>" logging interface
    And the library should use dependency injection for logging
    And the logger should be optional and not forced on consumers
    And the library should compile without the logger dependency

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  Scenario: Library documentation and examples
    Given I want a well-documented library
    When I generate a library with comprehensive documentation
    Then the project should include a detailed README
    And the project should include a package documentation file
    And the project should include usage examples
    And the examples should be executable and testable
    And the documentation should follow Go documentation standards

  Scenario: Library testing structure
    Given I want a thoroughly tested library
    When I generate a library with test infrastructure
    Then the project should include unit tests
    And the project should include example tests
    And the project should include benchmark tests
    And the project should include test coverage configuration
    And the tests should follow Go testing conventions

  Scenario: Library API design patterns
    Given I want a library with clean API design
    When I generate a library following best practices
    Then the library should use functional options pattern
    And the library should have minimal public API surface
    And the library should provide clear error types
    And the library should support context for cancellation
    And the library should be thread-safe

  Scenario: Library versioning and releases
    Given I want a library with proper versioning
    When I generate a library with release management
    Then the project should include semantic versioning
    And the project should include a CHANGELOG file
    And the project should include GitHub release workflow
    And the project should support Go module versioning
    And the releases should be automated via CI/CD

  Scenario: Library configuration patterns
    Given I want a configurable library
    When I generate a library with configuration options
    Then the library should support functional options
    And the library should provide sensible defaults
    And the library should validate configuration
    And the library should support environment variables
    And the configuration should be immutable after initialization

  Scenario: Library error handling
    Given I want robust error handling
    When I generate a library with error management
    Then the library should define custom error types
    And the errors should support wrapping and unwrapping
    And the errors should include context information
    And the errors should be comparable
    And the error messages should be clear and actionable

  Scenario: Library performance optimization
    Given I want a performant library
    When I generate a library with performance features
    Then the library should include benchmark tests
    And the library should minimize allocations
    And the library should support connection pooling
    And the library should implement caching where appropriate
    And the performance should be measurable

  Scenario: Library dependency management
    Given I want minimal dependencies
    When I generate a library with dependency constraints
    Then the library should have minimal external dependencies
    And the library should use Go standard library where possible
    And the dependencies should be clearly documented
    And the library should support dependency injection
    And the go.mod should be kept clean and minimal

  Scenario: Library examples directory
    Given I want comprehensive usage examples
    When I generate a library with example programs
    Then the project should include a separate examples directory
    And the examples should demonstrate basic usage
    And the examples should demonstrate advanced usage
    And the examples should be self-contained
    And the examples should include their own go.mod

  Scenario: Library continuous integration
    Given I want automated quality checks
    When I generate a library with CI/CD pipelines
    Then the project should include GitHub Actions workflows
    And the CI should run tests on multiple Go versions
    And the CI should check code formatting
    And the CI should run linting with golangci-lint
    And the CI should report test coverage

  Scenario: Library code quality
    Given I want high code quality standards
    When I generate a library with quality configurations
    Then the project should include golangci-lint configuration
    And the project should follow Go code conventions
    And the project should include pre-commit hooks
    And the code should pass go vet checks
    And the code should be properly formatted

  Scenario: Library licensing
    Given I want proper licensing
    When I generate a library with license "<license>"
    Then the project should include the "<license>" license file
    And the license should be referenced in README
    And the license should be included in file headers
    And the go.mod should reflect the license
    And the license should be compatible with dependencies

    Examples:
      | license      |
      | MIT          |
      | Apache-2.0   |
      | GPL-3.0      |
      | BSD-3-Clause |

  Scenario: Library metrics and observability
    Given I want observable library behavior
    When I generate a library with metrics support
    Then the library should support metrics collection
    And the metrics should be optional
    And the library should expose key performance indicators
    And the metrics should follow Prometheus conventions
    And the observability should have minimal overhead

  Scenario: Library retry and resilience
    Given I want resilient library operations
    When I generate a library with resilience patterns
    Then the library should support configurable retries
    And the library should implement exponential backoff
    And the library should handle transient failures
    And the library should support circuit breaker pattern
    And the resilience should be configurable

  Scenario: Library compatibility
    Given I want cross-platform compatibility
    When I generate a library for multiple platforms
    Then the library should work on Linux
    And the library should work on macOS
    And the library should work on Windows
    And the library should support different architectures
    And the compatibility should be tested in CI

  Scenario: Library security considerations
    Given I want a secure library
    When I generate a library with security features
    Then the library should validate all inputs
    And the library should handle sensitive data properly
    And the library should follow security best practices
    And the library should include security documentation
    And the dependencies should be scanned for vulnerabilities

  Scenario: Library resource management
    Given I want proper resource management
    When I generate a library with resources
    Then the library should implement Close() methods
    And the library should handle cleanup properly
    And the library should prevent resource leaks
    And the library should support graceful shutdown
    And the resources should be documented

  Scenario: Library interface design
    Given I want extensible interfaces
    When I generate a library with interfaces
    Then the library should define small interfaces
    And the interfaces should be consumer-focused
    And the library should accept interfaces
    And the library should return concrete types
    And the interfaces should be well-documented

  Scenario: Library concurrency support
    Given I want concurrent operations
    When I generate a library with concurrency
    Then the library should be thread-safe
    And the library should support concurrent operations
    And the library should handle synchronization properly
    And the library should avoid race conditions
    And the concurrency should be documented

  Scenario: Library migration support
    Given I want to support library evolution
    When I generate a library with migration support
    Then the library should support backward compatibility
    And the library should provide migration guides
    And the library should use deprecation notices
    And the library should follow semantic versioning
    And the breaking changes should be documented

  Scenario: Library distribution
    Given I want easy library distribution
    When I generate a library for distribution
    Then the library should be go-gettable
    And the library should work with Go modules
    And the library should include installation instructions
    And the library should support vendoring
    And the distribution should be automated

  Scenario: Library documentation generation
    Given I want generated documentation
    When I generate a library with doc generation
    Then the library should support godoc
    And the library should include doc.go
    And the library should have example functions
    And the documentation should be searchable
    And the docs should be published automatically