Feature: Critical Enterprise Combination Testing
  As a go-starter user building enterprise applications
  I want to ensure that complex architecture and feature combinations work correctly
  So that production deployments are reliable and maintain architectural integrity

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And I am testing enterprise-grade combinations

  @critical @enterprise
  Scenario Outline: Hexagonal architecture with enterprise features
    When I generate a web API project with configuration:
      | type         | web-api      |
      | architecture | hexagonal    |
      | framework    | <framework>  |
      | database     | <database>   |
      | orm          | <orm>        |
      | auth         | <auth>       |
      | logger       | <logger>     |
      | monitoring   | enabled      |
    Then the project should compile successfully
    And hexagonal architecture boundaries should be enforced
    And dependency injection should work correctly
    And <auth> authentication should be properly configured
    And <database> database integration should work
    And <logger> logging should be production-ready

    Examples:
      | framework | database  | orm   | auth    | logger  |
      | gin       | postgres  | gorm  | jwt     | zap     |
      | echo      | mysql     | sqlx  | oauth2  | zerolog |
      | fiber     | mongodb   | none  | api-key | slog    |
      | chi       | sqlite    | gorm  | session | logrus  |

  @critical @clean-architecture
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
    Then the project should compile successfully
    And clean architecture dependency rules should be enforced
    And inner layers should not depend on outer layers
    And use cases should be framework-independent
    And entities should be pure business logic
    And <asset_pipeline> assets should be properly integrated

    Examples:
      | framework | database | orm  | auth   | logger  | asset_pipeline |
      | gin       | postgres | gorm | jwt    | zap     | webpack        |
      | echo      | mysql    | sqlc | oauth2 | slog    | vite           |
      | fiber     | mongodb  | none | api-key| zerolog | esbuild        |

  @critical @ddd
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
    Then the project should compile successfully
    And domain model should be rich and expressive
    And bounded contexts should be properly defined
    And aggregates should enforce business invariants
    And domain events should be properly implemented
    And repository patterns should abstract persistence

    Examples:
      | framework | database | orm  | auth    | logger |
      | gin       | postgres | ent  | jwt     | zap    |
      | echo      | mysql    | gorm | oauth2  | slog   |

  @critical @workspace
  Scenario: Complex workspace with multiple blueprint types
    When I generate a workspace project with configuration:
      | type         | workspace |
      | name         | enterprise-system |
      | components   | api,worker,cli,functions |
    And I add component "api" with configuration:
      | type         | web-api   |
      | architecture | hexagonal |
      | framework    | gin       |
      | database     | postgres  |
      | auth         | jwt       |
    And I add component "worker" with configuration:
      | type         | microservice |
      | architecture | clean        |
      | framework    | grpc         |
      | database     | postgres     |
    And I add component "cli" with configuration:
      | type         | cli      |
      | complexity   | standard |
      | framework    | cobra    |
    And I add component "functions" with configuration:
      | type         | lambda   |
      | auth         | api-key  |
    Then the workspace should compile successfully
    And all components should share consistent go.mod dependencies
    And inter-component communication should be properly configured
    And shared database connections should be managed correctly
    And workspace-level configuration should be consistent

  @critical @multi-database
  Scenario Outline: Multi-database microservice patterns
    When I generate a microservice project with configuration:
      | type           | microservice  |
      | framework      | <framework>   |
      | primary_db     | <primary_db>  |
      | cache_db       | <cache_db>    |
      | auth           | <auth>        |
      | logger         | <logger>      |
      | monitoring     | enabled       |
    Then the project should compile successfully
    And primary database connections should be configured for <primary_db>
    And cache layer should be configured for <cache_db>
    And transaction management should work across databases
    And connection pooling should be optimized
    And health checks should monitor all database connections

    Examples:
      | framework | primary_db | cache_db | auth | logger  |
      | gin       | postgres   | redis    | jwt  | zap     |
      | echo      | mysql      | redis    | oauth2| slog   |
      | fiber     | mongodb    | redis    | api-key| zerolog|

  @critical @event-driven
  Scenario Outline: Event-driven architecture with CQRS
    When I generate an event-driven project with configuration:
      | type         | event-driven |
      | framework    | <framework>  |
      | database     | <database>   |
      | event_store  | <event_store>|
      | auth         | <auth>       |
      | logger       | <logger>     |
    Then the project should compile successfully
    And CQRS pattern should be properly implemented
    And event sourcing should work correctly
    And command handlers should be separated from query handlers
    And event store should be configured for <event_store>
    And projection rebuilding should be supported

    Examples:
      | framework | database | event_store | auth | logger |
      | gin       | postgres | postgres    | jwt  | zap    |
      | echo      | mongodb  | mongodb     | oauth2| slog  |

  @critical @performance
  Scenario Outline: High-performance configuration combinations
    When I generate a web API project optimized for performance:
      | type         | web-api     |
      | architecture | <architecture> |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | logger       | <logger>    |
      | cache        | enabled     |
      | monitoring   | enabled     |
    Then the project should compile successfully
    And database connections should be pooled and optimized
    And caching layer should be properly configured
    And logging should use structured, high-performance loggers
    And monitoring should have minimal performance overhead
    And memory allocations should be optimized

    Examples:
      | architecture | framework | database | orm  | logger  |
      | standard     | fiber     | postgres | sqlx | zerolog |
      | clean        | gin       | mysql    | gorm | zap     |
      | hexagonal    | echo      | postgres | ent  | slog    |

  @critical @security
  Scenario Outline: Security-focused combinations
    When I generate a web API project with security focus:
      | type         | web-api     |
      | architecture | <architecture> |
      | framework    | <framework> |
      | database     | <database>  |
      | auth         | <auth>      |
      | security     | hardened    |
      | tls          | enabled     |
      | audit        | enabled     |
    Then the project should compile successfully
    And authentication should be properly secured
    And database connections should use encrypted connections
    And all HTTP endpoints should enforce HTTPS
    And audit logging should capture security events
    And input validation should prevent injection attacks
    And CORS should be properly configured

    Examples:
      | architecture | framework | database | auth    |
      | clean        | gin       | postgres | jwt     |
      | hexagonal    | echo      | mysql    | oauth2  |
      | ddd          | fiber     | postgres | api-key |

  @critical @deployment
  Scenario Outline: Deployment-ready combinations
    When I generate a project ready for production deployment:
      | type           | <project_type> |
      | architecture   | <architecture> |
      | framework      | <framework>    |
      | database       | <database>     |
      | deployment     | <deployment>   |
      | monitoring     | enabled        |
      | health_checks  | enabled        |
      | metrics        | prometheus     |
    Then the project should compile successfully
    And <deployment> deployment configuration should be generated
    And health check endpoints should be available
    And metrics should be exposed for Prometheus
    And logging should be container-friendly
    And graceful shutdown should be implemented

    Examples:
      | project_type | architecture | framework | database | deployment |
      | web-api      | clean        | gin       | postgres | docker     |
      | microservice | hexagonal    | echo      | mysql    | kubernetes |
      | web-api      | ddd          | fiber     | mongodb  | lambda     |