Feature: Web API Blueprint Generation
  As a Go developer
  I want to generate modern, production-ready Web API applications
  So that I can quickly build scalable RESTful services with industry best practices

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate standard web API application
    Given I want to create a standard web API application
    When I run the command "go-starter new my-api --type=web-api --framework=gin --database-driver=postgres --no-git"
    Then the generation should succeed
    And the project should contain all essential web API components
    And the generated code should compile successfully
    And the API should expose OpenAPI documentation

  Scenario: Web API with different frameworks
    Given I want to create a web API application
    When I generate a web API with framework "<framework>"
    Then the project should use the "<framework>" web framework
    And the handlers should use "<framework>"-specific patterns
    And the middleware should be framework-compatible
    And the application should compile and serve requests

    Examples:
      | framework |
      | gin       |
      | echo      |
      | fiber     |
      | chi       |
      | stdlib    |

  Scenario: Web API with different architectures
    Given I want to create a web API with specific architecture
    When I generate a web API with architecture "<architecture>"
    Then the project should follow "<architecture>" patterns
    And the directory structure should reflect the architecture
    And the dependencies should flow in the correct direction
    And the code should be properly layered

    Examples:
      | architecture |
      | standard     |
      | clean        |
      | ddd          |
      | hexagonal    |

  Scenario: Database integration testing
    Given I want to create a web API with database support
    When I generate a web API with database "<database>" and ORM "<orm>"
    Then the project should include database configuration
    And the migration system should be properly configured
    And the repository layer should use the specified ORM
    And the database connection should be testable with containers

    Examples:
      | database  | orm  |
      | postgres  | gorm |
      | postgres  | sqlx |
      | mysql     | gorm |
      | mysql     | sqlx |
      | sqlite    | gorm |

  Scenario: Authentication and authorization
    Given I want to secure my web API
    When I generate a web API with authentication type "<auth_type>"
    Then the API should include authentication endpoints
    And the middleware should enforce authentication
    And the JWT/Session management should be secure
    And the protected routes should require valid credentials

    Examples:
      | auth_type |
      | jwt       |
      | session   |
      | oauth2    |
      | api-key   |

  Scenario: RESTful API endpoints
    Given I have generated a web API application
    When I examine the API endpoints
    Then the API should follow RESTful conventions
    And the endpoints should include proper HTTP verbs
    And the responses should use appropriate status codes
    And the API should handle CRUD operations correctly

  Scenario: API documentation and OpenAPI specification
    Given I have generated a web API application
    When I examine the API documentation
    Then the API should include OpenAPI 3.0 specification
    And the endpoints should be properly documented
    And the request/response schemas should be defined
    And the API should be testable via documentation

  Scenario: Error handling and validation
    Given I want robust error handling
    When I generate a web API application
    Then the API should include structured error responses
    And input validation should be implemented
    And error messages should be informative but secure
    And different error types should have appropriate HTTP codes

  Scenario: Testing infrastructure
    Given I want comprehensive testing
    When I generate a web API application
    Then the project should include unit tests
    And the project should include integration tests
    And the tests should use testcontainers for database testing
    And the test coverage should be measurable

  Scenario: Security best practices
    Given I want a secure web API
    When I generate a web API with security features
    Then the API should implement CORS properly
    And security headers should be configured
    And input sanitization should prevent injection attacks
    And rate limiting should be configurable

  Scenario: Performance and monitoring
    Given I want a production-ready web API
    When I generate a web API with monitoring
    Then the API should include health check endpoints
    And metrics collection should be implemented
    And request tracing should be available
    And performance monitoring should be configured

  Scenario: Container deployment
    Given I want to deploy my web API
    When I examine the deployment configuration
    Then the project should include Dockerfile
    And the container should be optimized for production
    And the deployment should support environment variables
    And the health checks should work in containers

  Scenario: CI/CD integration
    Given I want automated deployments
    When I generate a web API with CI/CD
    Then the project should include GitHub Actions workflows
    And the CI should run tests and security scans
    And the deployment should support multiple environments
    And the pipeline should include quality gates

  Scenario: Logging and observability
    Given I want observable web APIs
    When I generate a web API with logger "<logger>"
    Then the application should use structured logging
    And log levels should be configurable
    And request/response logging should be available
    And log correlation should be implemented

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  Scenario: Architecture pattern compliance
    Given I want to implement specific architecture patterns
    When I generate a web API with architecture "<architecture>"
    Then the project should follow the architectural principles
    And the code should be properly organized according to the pattern
    And the dependencies should flow in the correct direction
    And the architecture should support testability and maintainability

    Examples:
      | architecture |
      | standard     |
      | clean        |
      | ddd          |
      | hexagonal    |

  Scenario: Microservice readiness
    Given I want to build microservices
    When I generate a web API for microservice architecture
    Then the service should be independently deployable
    And the configuration should support service discovery
    And the API should include circuit breaker patterns
    And the service should support distributed tracing

  Scenario: Multi-environment configuration
    Given I need different deployment environments
    When I generate a web API with environment support
    Then the configuration should support dev/test/prod
    And environment variables should override defaults
    And sensitive data should be externalized
    And feature flags should be configurable

  Scenario: API versioning
    Given I need to maintain API compatibility
    When I generate a web API with versioning
    Then the API should support version prefixes
    And backward compatibility should be maintained
    And version deprecation should be handled gracefully
    And clients should receive version information

  Scenario: Content negotiation
    Given I want flexible API responses
    When I generate a web API with content negotiation
    Then the API should support JSON by default
    And alternative formats should be configurable
    And Accept headers should be respected
    And Content-Type should be properly set

  Scenario: Graceful shutdown
    Given I want reliable deployments
    When I generate a web API with graceful shutdown
    Then the server should handle shutdown signals
    And in-flight requests should complete
    And database connections should close cleanly
    And the shutdown should be logged appropriately