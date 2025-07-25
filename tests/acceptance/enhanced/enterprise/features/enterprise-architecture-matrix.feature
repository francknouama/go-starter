Feature: Enterprise Architecture Matrix Testing
  As a go-starter user building enterprise applications
  I want to ensure that complex architecture patterns work correctly with all supported configurations
  So that production enterprise deployments are reliable and maintainable

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And I am testing enterprise architecture combinations

  @critical @enterprise @hexagonal @p0
  Scenario Outline: Hexagonal architecture with all database combinations
    When I generate a web API project with configuration:
      | type         | web-api      |
      | architecture | hexagonal    |
      | framework    | <framework>  |
      | database     | <database>   |
      | orm          | <orm>        |
      | auth         | <auth>       |
      | logger       | <logger>     |
      | go_version   | 1.23         |
    Then the project should compile successfully
    And hexagonal architecture boundaries should be enforced
    And domain layer should be isolated from infrastructure
    And ports and adapters pattern should be implemented correctly
    And dependency injection should work with <framework>
    And <database> database should be properly abstracted
    And <auth> authentication should integrate with domain layer
    And <logger> logging should be available in all layers
    And repository interfaces should be defined in domain layer
    And infrastructure adapters should implement domain interfaces

    Examples:
      | framework | database | orm  | auth    | logger  |
      | gin       | postgres | gorm | jwt     | zap     |
      | echo      | mysql    | sqlx | oauth2  | zerolog |
      | fiber     | sqlite   | gorm | session | logrus  |
      | chi       | postgres | sqlx | jwt     | slog    |
      | gin       | mysql    | gorm | oauth2  | zap     |

  @critical @enterprise @clean @p0
  Scenario Outline: Clean architecture with complex features
    When I generate a web API project with configuration:
      | type         | web-api          |
      | architecture | clean            |
      | framework    | <framework>      |
      | database     | <database>       |
      | orm          | <orm>            |
      | auth         | <auth>           |
      | logger       | <logger>         |
      | asset_pipeline | <asset_pipeline> |
      | go_version   | 1.23            |
    Then the project should compile successfully
    And clean architecture dependency rules should be enforced
    And entities should not depend on external frameworks
    And use cases should be framework-independent
    And interface adapters should bridge frameworks and use cases
    And frameworks and drivers should be in outermost layer
    And dependency inversion principle should be applied
    And <asset_pipeline> assets should be properly integrated
    And <auth> authentication should follow clean architecture patterns
    And database layer should not leak into business logic

    Examples:
      | framework | database | orm  | auth    | logger  | asset_pipeline |
      | gin       | postgres | gorm | jwt     | zap     | webpack        |
      | echo      | mysql    | sqlx | oauth2  | slog    | vite           |
      | fiber     | sqlite   | gorm | session | zerolog | esbuild        |
      | chi       | postgres | sqlx | jwt     | logrus  | embedded       |

  @critical @enterprise @ddd @p0
  Scenario Outline: Domain-Driven Design with enterprise patterns
    When I generate a web API project with configuration:
      | type         | web-api     |
      | architecture | ddd         |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | complexity   | expert      |
      | go_version   | 1.23        |
    Then the project should compile successfully
    And domain model should be rich and expressive
    And bounded contexts should be properly defined
    And aggregates should enforce business invariants
    And domain events should be properly implemented
    And value objects should be immutable
    And domain services should encapsulate domain logic
    And repository patterns should abstract persistence
    And application services should orchestrate use cases
    And infrastructure should be separated from domain
    And <database> persistence should support domain model

    Examples:
      | framework | database | orm  | auth    | logger |
      | gin       | postgres | gorm | jwt     | zap    |
      | echo      | postgres | sqlx | oauth2  | slog   |
      | fiber     | mysql    | gorm | session | zerolog|

  @critical @enterprise @standard @p0
  Scenario Outline: Standard architecture with production optimizations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | architecture | standard    |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23        |
    Then the project should compile successfully
    And standard layered architecture should be implemented
    And handlers should delegate to services
    And services should coordinate business operations
    And repositories should handle data persistence
    And middleware should be properly configured
    And error handling should be consistent across layers
    And <framework> routing should be properly organized
    And <database> connections should be optimized
    And <auth> authentication should be production-ready

    Examples:
      | framework | database | orm  | auth    | logger |
      | gin       | postgres | gorm | jwt     | zap    |
      | echo      | mysql    | sqlx | oauth2  | slog   |
      | fiber     | sqlite   | gorm | session | zerolog|
      | chi       | postgres | sqlx | jwt     | logrus |
      | gin       | mysql    | gorm | oauth2  | zap    |

  @integration @enterprise @microservice
  Scenario Outline: Microservice architecture patterns
    When I generate a microservice project with configuration:
      | type         | microservice |
      | architecture | <architecture> |
      | framework    | <framework>    |
      | database     | <database>     |
      | orm          | <orm>          |
      | auth         | <auth>         |
      | logger       | <logger>       |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And microservice patterns should be implemented
    And service discovery should be configured
    And health checks should be available
    And metrics collection should be set up
    And distributed tracing should be configured
    And circuit breaker patterns should be available
    And <architecture> architecture should be properly applied
    And API documentation should be generated

    Examples:
      | architecture | framework | database | orm  | auth   | logger |
      | clean        | gin       | postgres | gorm | jwt    | zap    |
      | hexagonal    | echo      | mysql    | sqlx | oauth2 | slog   |
      | standard     | fiber     | postgres | gorm | session| zerolog|

  @integration @enterprise @event-driven
  Scenario Outline: Event-driven architecture with CQRS
    When I generate an event-driven project with configuration:
      | type         | event-driven |
      | framework    | <framework>  |
      | database     | <database>   |
      | orm          | <orm>        |
      | auth         | <auth>       |
      | logger       | <logger>     |
      | go_version   | 1.23        |
    Then the project should compile successfully
    And CQRS pattern should be properly implemented
    And command handlers should be separated from query handlers
    And event sourcing should work correctly
    And event store should be configured for <database>
    And projection rebuilding should be supported
    And saga patterns should be available
    And eventual consistency should be handled
    And event versioning should be supported

    Examples:
      | framework | database | orm  | auth   | logger |
      | gin       | postgres | gorm | jwt    | zap    |
      | echo      | postgres | sqlx | oauth2 | slog   |

  @integration @enterprise @grpc
  Scenario Outline: gRPC Gateway enterprise patterns
    When I generate a gRPC Gateway project with configuration:
      | type         | grpc-gateway |
      | database     | <database>   |
      | orm          | <orm>        |
      | auth         | <auth>       |
      | logger       | <logger>     |
      | go_version   | 1.23        |
    Then the project should compile successfully
    And gRPC services should be properly defined
    And protobuf definitions should be valid
    And Gateway should proxy HTTP to gRPC correctly
    And authentication should work for both HTTP and gRPC
    And both REST and gRPC endpoints should be available
    And service mesh compatibility should be maintained
    And load balancing should be configured
    And SSL/TLS should be properly configured

    Examples:
      | database | orm  | auth   | logger |
      | postgres | gorm | jwt    | zap    |
      | mysql    | sqlx | oauth2 | slog   |

  @performance @enterprise @architecture
  Scenario Outline: High-performance architecture combinations
    When I generate a web API project optimized for performance:
      | type           | web-api         |
      | architecture   | <architecture>  |
      | framework      | <framework>     |
      | database       | <database>      |
      | orm            | <orm>           |
      | logger         | <logger>        |
      | go_version     | 1.23           |
    Then the project should compile successfully
    And performance optimizations should be implemented
    And database connections should be pooled efficiently
    And logging should have minimal performance overhead
    And memory allocations should be optimized
    And response times should be within acceptable limits
    And concurrent request handling should be optimized
    And resource usage should be monitored

    Examples:
      | architecture | framework | database | orm  | logger  |
      | standard     | fiber     | postgres | sqlx | zerolog |
      | clean        | gin       | mysql    | gorm | zap     |
      | hexagonal    | echo      | postgres | sqlx | slog    |

  @security @enterprise @architecture
  Scenario Outline: Security-focused architecture combinations
    When I generate a web API project with security focus:
      | type         | web-api         |
      | architecture | <architecture>  |
      | framework    | <framework>     |
      | database     | <database>      |
      | auth         | <auth>          |
      | logger       | <logger>        |
      | go_version   | 1.23           |
    Then the project should compile successfully
    And security best practices should be implemented
    And authentication should be properly secured
    And authorization should be role-based
    And input validation should prevent injection attacks
    And HTTPS should be enforced
    And security headers should be configured
    And audit logging should capture security events
    And sensitive data should be properly handled
    And CORS should be properly configured

    Examples:
      | architecture | framework | database | auth    | logger |
      | clean        | gin       | postgres | jwt     | zap    |
      | hexagonal    | echo      | mysql    | oauth2  | slog   |
      | ddd          | fiber     | postgres | session | zerolog|

  @deployment @enterprise @architecture
  Scenario Outline: Deployment-ready architecture combinations
    When I generate a project ready for production deployment:
      | type           | web-api         |
      | architecture   | <architecture>  |
      | framework      | <framework>     |
      | database       | <database>      |
      | auth           | <auth>          |
      | logger         | <logger>        |
      | go_version     | 1.23           |
    Then the project should compile successfully
    And deployment configurations should be generated
    And Docker configuration should be optimized
    And health check endpoints should be available
    And metrics should be exposed for monitoring
    And logging should be container-friendly
    And graceful shutdown should be implemented
    And configuration management should be environment-aware
    And database migrations should be automated

    Examples:
      | architecture | framework | database | auth   | logger |
      | clean        | gin       | postgres | jwt    | zap    |
      | hexagonal    | echo      | mysql    | oauth2 | slog   |
      | ddd          | fiber     | postgres | session| zerolog|