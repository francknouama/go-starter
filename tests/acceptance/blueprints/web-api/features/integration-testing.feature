Feature: Web API Integration Testing
  As a developer
  I want comprehensive integration testing for web API blueprints
  So that I can ensure database and logger integration works across architectures

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Database integration across architectures
    Given I want to validate database integration
    When I generate a web API with architecture "<architecture>", database "<database>", and ORM "<orm>"
    Then the generation should succeed
    And database configuration should follow architecture patterns
    And data access should be properly abstracted
    And transactions should be handled correctly
    And migrations should be architecture-appropriate
    And the generated code should compile successfully
    And database tests should use testcontainers

    Examples:
      | architecture | database | orm  |
      | standard     | postgres | gorm |
      | standard     | postgres | sqlx |
      | standard     | mysql    | gorm |
      | standard     | sqlite   | gorm |
      | clean        | postgres | gorm |
      | clean        | mysql    | gorm |
      | clean        | sqlite   | gorm |
      | ddd          | postgres | gorm |
      | ddd          | mysql    | gorm |
      | hexagonal    | postgres | gorm |
      | hexagonal    | postgres | sqlx |
      | hexagonal    | mysql    | gorm |

  Scenario: Logger integration across architectures
    Given I want to validate logger integration
    When I generate a web API with architecture "<architecture>" and logger "<logger>"
    Then the generation should succeed
    And logging should follow architecture patterns
    And log configuration should be architecture-appropriate
    And structured logging should be properly implemented
    And log levels should be configurable
    And the generated code should compile successfully

    Examples:
      | architecture | logger  |
      | standard     | slog    |
      | standard     | zap     |
      | standard     | logrus  |
      | standard     | zerolog |
      | clean        | slog    |
      | clean        | zap     |
      | clean        | logrus  |
      | clean        | zerolog |
      | ddd          | slog    |
      | ddd          | zap     |
      | ddd          | logrus  |
      | hexagonal    | slog    |
      | hexagonal    | zap     |
      | hexagonal    | logrus  |
      | hexagonal    | zerolog |

  Scenario: Framework integration across architectures
    Given I want to validate framework integration
    When I generate a web API with architecture "<architecture>" and framework "<framework>"
    Then the generation should succeed
    And framework integration should follow architecture patterns
    And HTTP handling should be properly layered
    And routing should be architecture-appropriate
    And middleware should be correctly organized
    And the generated code should compile successfully

    Examples:
      | architecture | framework |
      | standard     | gin       |
      | standard     | echo      |
      | standard     | fiber     |
      | standard     | chi       |
      | clean        | gin       |
      | clean        | echo      |
      | clean        | fiber     |
      | ddd          | gin       |
      | ddd          | echo      |
      | hexagonal    | gin       |
      | hexagonal    | echo      |
      | hexagonal    | fiber     |

  Scenario: Authentication integration across architectures
    Given I want to validate authentication integration
    When I generate a web API with architecture "<architecture>" and authentication "<auth_type>"
    Then the generation should succeed
    And authentication should follow architecture patterns
    And security concerns should be properly layered
    And authentication logic should be appropriately placed
    And authorization should be architecture-compliant
    And the generated code should compile successfully

    Examples:
      | architecture | auth_type |
      | standard     | jwt       |
      | standard     | session   |
      | standard     | api-key   |
      | clean        | jwt       |
      | clean        | session   |
      | ddd          | jwt       |
      | ddd          | session   |
      | hexagonal    | jwt       |
      | hexagonal    | session   |
      | hexagonal    | api-key   |

  Scenario: Multi-feature integration testing
    Given I want to validate complete feature integration
    When I generate a web API with:
      | Feature       | Value    |
      | architecture  | clean    |
      | framework     | gin      |
      | database      | postgres |
      | orm           | gorm     |
      | logger        | zap      |
      | auth          | jwt      |
    Then the generation should succeed
    And all features should work together harmoniously
    And architecture patterns should be maintained
    And no feature should violate architecture boundaries
    And the generated code should compile successfully
    And integration tests should pass

  Scenario: Database migration integration
    Given I want to validate database migrations
    When I generate a web API with database support
    Then migration files should be generated
    And migration system should be properly configured
    And migrations should be executable
    And rollback functionality should be available
    And migration status should be trackable

  Scenario: Container integration testing
    Given I want to validate container integration
    When I generate a web API with containerization
    Then Docker configuration should be optimized
    And container should include all dependencies
    And container should support environment configuration
    And health checks should work in containers
    And container should be production-ready

  Scenario: CI/CD integration validation
    Given I want to validate CI/CD integration
    When I generate a web API with CI/CD configuration
    Then GitHub Actions workflows should be included
    And workflows should run tests and security scans
    And deployment should support multiple environments
    And quality gates should be properly configured
    And artifacts should be properly managed

  Scenario: Performance integration testing
    Given I want to validate performance features
    When I generate a web API with performance monitoring
    Then metrics collection should be properly integrated
    And performance monitoring should not violate architecture
    And metrics should be exportable
    And performance overhead should be minimal
    And monitoring should be configurable

  Scenario: Security integration testing
    Given I want to validate security integration
    When I generate a web API with comprehensive security
    Then security headers should be properly configured
    And CORS should be appropriately implemented
    And input validation should prevent common attacks
    And rate limiting should be properly integrated
    And security should not violate architecture patterns

  Scenario: Error handling integration
    Given I want to validate error handling integration
    When I generate a web API with comprehensive error handling
    Then error responses should be consistent across architectures
    And error logging should be properly integrated
    And error recovery should be architecture-appropriate
    And custom errors should be supportable
    And error handling should not leak implementation details

  Scenario: Testing infrastructure integration
    Given I want to validate testing infrastructure
    When I generate a web API with comprehensive testing
    Then unit tests should be architecture-specific
    And integration tests should use testcontainers
    And acceptance tests should cover user scenarios
    And test utilities should be provided
    And test coverage should be measurable

  Scenario: Configuration integration testing
    Given I want to validate configuration integration
    When I generate a web API with environment configuration
    Then configuration should support multiple environments
    And configuration should be architecture-appropriate
    And sensitive data should be properly externalized
    And configuration validation should be implemented
    And configuration should be easily changeable

  Scenario: API documentation integration
    Given I want to validate API documentation integration
    When I generate a web API with documentation
    Then OpenAPI specification should be generated
    And documentation should be architecture-aware
    And API examples should be provided
    And documentation should be automatically updated
    And documentation should be easily accessible

  Scenario: Monitoring and observability integration
    Given I want to validate monitoring integration
    When I generate a web API with observability features
    Then monitoring should be architecture-appropriate
    And observability should not violate layer boundaries
    And metrics should provide valuable insights
    And monitoring should be easily configurable
    And observability should support debugging