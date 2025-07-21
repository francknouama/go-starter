Feature: Standard Web API Blueprint
  As a developer
  I want to generate a standard web API project
  So that I can quickly build REST APIs with conventional patterns

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate standard web API with Gin
    Given I want a standard web API
    When I run the command "go-starter new my-gin-api --type=web-api-standard --framework=gin --no-git"
    Then the generation should succeed
    And the project should include Gin router setup
    And the project should have basic CRUD endpoints
    And the project should compile and run successfully

  Scenario: Generate standard web API with different frameworks
    Given I want a standard web API with different frameworks
    When I generate a standard web API with framework "<framework>"
    Then the project should use the "<framework>" web framework
    And the handlers should follow standard patterns
    And the middleware should be properly configured
    And the routing should be framework-appropriate
    And the generated code should compile successfully

    Examples:
      | framework |
      | gin       |
      | echo      |
      | fiber     |
      | chi       |
      | stdlib    |

  Scenario: Generate with database integration
    Given I want a standard web API with database
    When I generate a standard web API with database "<database>" and ORM "<orm>"
    Then the generation should succeed
    And database configuration should be included
    And ORM models should be provided
    And database connection should be testable
    And migration system should be configured
    And the generated code should compile successfully

    Examples:
      | database | orm  |
      | postgres | gorm |
      | postgres | sqlx |
      | mysql    | gorm |
      | mysql    | sqlx |
      | sqlite   | gorm |

  Scenario: Standard project structure validation
    Given I want to validate standard project structure
    When I generate a standard web API
    Then the project should have these directories:
      | Directory             | Purpose                  |
      | internal/handlers/    | HTTP request handlers    |
      | internal/services/    | Business logic           |
      | internal/repository/  | Data access layer        |
      | internal/models/      | Data models              |
      | internal/middleware/  | HTTP middleware          |
      | internal/config/      | Configuration            |
      | internal/database/    | Database setup           |
      | api/                  | API documentation        |
      | cmd/server/           | Application entry point  |
      | configs/              | Configuration files      |
      | migrations/           | Database migrations      |
      | scripts/              | Build and setup scripts  |

  Scenario: RESTful endpoint generation
    Given I want standard RESTful endpoints
    When I generate a standard web API
    Then the API should include these endpoints:
      | Method | Path           | Purpose           |
      | GET    | /api/health    | Health check      |
      | GET    | /api/users     | List users        |
      | POST   | /api/users     | Create user       |
      | GET    | /api/users/:id | Get user by ID    |
      | PUT    | /api/users/:id | Update user       |
      | DELETE | /api/users/:id | Delete user       |
    And endpoints should follow REST conventions
    And responses should use appropriate HTTP status codes

  Scenario: Middleware configuration
    Given I want proper middleware configuration
    When I generate a standard web API
    Then the API should include these middleware:
      | Middleware    | Purpose                  |
      | CORS          | Cross-origin requests    |
      | Logger        | Request logging          |
      | Recovery      | Panic recovery           |
      | RateLimit     | Rate limiting            |
      | Authentication| Auth validation          |
      | Validation    | Request validation       |
    And middleware should be properly ordered
    And middleware should be configurable

  Scenario: Authentication integration
    Given I want authentication in my standard API
    When I generate a standard web API with authentication type "<auth_type>"
    Then the API should include authentication middleware
    And protected routes should require authentication
    And authentication endpoints should be available
    And JWT/session handling should be implemented
    And the generated code should compile successfully

    Examples:
      | auth_type |
      | jwt       |
      | session   |
      | api-key   |

  Scenario: Request validation
    Given I want request validation
    When I generate a standard web API
    Then input validation should be implemented
    And validation errors should return appropriate responses
    And validation rules should be configurable
    And custom validators should be supportable
    And validation should prevent injection attacks

  Scenario: Error handling
    Given I want robust error handling
    When I generate a standard web API
    Then structured error responses should be implemented
    And different error types should have appropriate HTTP codes
    And error messages should be informative but secure
    And error logging should be implemented
    And panic recovery should be configured

  Scenario: Database connection and models
    Given I want database functionality
    When I generate a standard web API with database support
    Then database connection should be properly configured
    And connection pooling should be implemented
    And models should be defined with appropriate fields
    And database migrations should be available
    And health checks should include database status

  Scenario: Configuration management
    Given I want configurable applications
    When I generate a standard web API
    Then environment-based configuration should be supported
    And configuration should include these settings:
      | Setting       | Purpose                |
      | server.port   | Server listening port  |
      | server.host   | Server host address    |
      | database.url  | Database connection    |
      | auth.secret   | Authentication secret  |
      | cors.origins  | Allowed CORS origins   |
      | log.level     | Logging level          |
    And sensitive data should be externalized
    And configuration validation should be implemented

  Scenario: Logging integration
    Given I want comprehensive logging
    When I generate a standard web API with logger "<logger>"
    Then structured logging should be implemented
    And request/response logging should be available
    And log levels should be configurable
    And log correlation should be supported
    And the generated code should compile successfully

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  Scenario: Testing infrastructure
    Given I want comprehensive testing
    When I generate a standard web API
    Then unit tests should be included for:
      | Component    | Test Focus              |
      | handlers     | HTTP request handling   |
      | services     | Business logic          |
      | repository   | Data access             |
      | middleware   | HTTP middleware         |
    And integration tests should use testcontainers
    And test coverage should be measurable
    And test utilities should be provided

  Scenario: API documentation
    Given I want documented APIs
    When I generate a standard web API
    Then OpenAPI 3.0 specification should be generated
    And API endpoints should be documented
    And request/response schemas should be defined
    And documentation should be accessible via web interface
    And examples should be provided

  Scenario: Health check implementation
    Given I want monitoring capabilities
    When I generate a standard web API
    Then health check endpoint should be available
    And health checks should verify:
      | Component | Check                    |
      | database  | Connection status        |
      | cache     | Cache connectivity       |
      | external  | External service status  |
      | memory    | Memory usage             |
      | disk      | Disk space               |
    And health status should be appropriately formatted
    And health checks should support different formats

  Scenario: Performance and monitoring
    Given I want performance monitoring
    When I generate a standard web API
    Then metrics collection should be implemented
    And request timing should be measured
    And performance monitoring should include:
      | Metric           | Purpose                  |
      | request_duration | Response time tracking   |
      | request_count    | Request volume           |
      | error_rate       | Error percentage         |
      | concurrent_users | Active connections       |
    And metrics should be exportable

  Scenario: Security best practices
    Given I want secure APIs
    When I generate a standard web API
    Then security headers should be configured
    And CORS should be properly implemented
    And input sanitization should prevent injection
    And rate limiting should be available
    And authentication should be secure
    And sensitive data should be protected

  Scenario: Container deployment
    Given I want containerized deployment
    When I generate a standard web API
    Then Dockerfile should be optimized for production
    And container should support environment variables
    And health checks should work in containers
    And container size should be minimized
    And security scanning should be supported