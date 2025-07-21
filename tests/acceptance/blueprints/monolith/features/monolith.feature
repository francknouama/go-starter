Feature: Monolith Blueprint Generation
  As a Go developer
  I want to generate a complete monolith web application
  So that I can quickly start building production-ready web services

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate basic monolith application
    Given I want to create a monolith application
    When I run the command "go-starter new my-monolith --type=monolith --framework=gin --database-driver=postgres --no-git"
    Then the generation should succeed
    And the project should contain all essential monolith components
    And the generated code should compile successfully

  Scenario: Monolith with different web frameworks
    Given I want to create a monolith application
    When I generate a monolith with framework "<framework>"
    Then the project should use the "<framework>" web framework
    And the code should include "<framework>"-specific imports
    And the application should compile and run

    Examples:
      | framework |
      | gin       |
      | echo      |
      | fiber     |
      | chi       |

  Scenario: Monolith with different database drivers
    Given I want to create a monolith application
    When I generate a monolith with database "<database>" and ORM "<orm>"
    Then the project should include "<database>" database configuration
    And the migration files should support "<database>"
    And the ORM integration should use "<orm>"

    Examples:
      | database  | orm  |
      | postgres  | gorm |
      | mysql     | gorm |
      | sqlite    | gorm |
      | postgres  | sqlx |

  Scenario: Monolith with authentication
    Given I want to create a secure monolith application
    When I generate a monolith with authentication type "<auth_type>"
    Then the project should include "<auth_type>" authentication setup
    And the session management should be properly configured
    And the security headers should be implemented

    Examples:
      | auth_type |
      | session   |
      | jwt       |
      | oauth2    |

  Scenario: Production-ready monolith
    Given I want to create a production-ready monolith
    When I generate a monolith with all production features
    Then the project should include Docker configuration
    And the project should include CI/CD pipelines
    And the project should include Kubernetes deployment
    And the project should include monitoring and health checks
    And the project should include comprehensive testing

  Scenario: Monolith asset pipeline integration
    Given I want to create a monolith with frontend assets
    When I generate a monolith with asset pipeline enabled
    Then the project should include Vite configuration
    And the project should include Tailwind CSS setup
    And the asset build process should be integrated with the backend

  Scenario: Modular monolith architecture
    Given I want to create a well-structured monolith
    When I generate a monolith application
    Then the code should follow modular monolith patterns
    And the layers should be properly separated
    And the dependencies should flow in the correct direction
    And the code should be easily testable

  Scenario: Database migration system
    Given I have generated a monolith with database support
    When I examine the migration system
    Then the project should include migration scripts
    And the migration commands should work correctly
    And the database schema should be properly versioned

  Scenario: Comprehensive testing setup
    Given I want to ensure code quality
    When I generate a monolith application
    Then the project should include unit tests
    And the project should include integration tests
    And the project should include benchmark tests
    And the tests should use proper mocking
    And the test coverage should be measurable

  Scenario: Development workflow
    Given I have generated a monolith application
    When I use the development tools
    Then the Makefile should provide essential commands
    And the development server should support hot reload
    And the linting should pass without errors
    And the security scanning should be configured

  Scenario: Multi-logger support
    Given I want flexible logging options
    When I generate a monolith with logger "<logger>"
    Then the application should use the "<logger>" logging library
    And the logging should follow consistent interface patterns
    And the log levels should be properly configured

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  Scenario: Performance optimization
    Given I want a high-performance monolith
    When I generate a monolith with performance optimizations
    Then the database should use connection pooling
    And the application should include caching strategies
    And the asset pipeline should optimize for production
    And benchmark tests should validate performance requirements

  Scenario: Security hardening
    Given I want a secure monolith application
    When I generate a monolith with security features
    Then the application should implement OWASP security headers
    And the sessions should use secure configurations
    And the input validation should prevent common attacks
    And the security scanning should be automated

  Scenario: Deployment readiness
    Given I want to deploy my monolith to production
    When I examine the deployment configuration
    Then the project should include multiple deployment targets
    And the CI/CD should support staging and production environments
    And the rollback procedures should be documented
    And the health checks should validate system status