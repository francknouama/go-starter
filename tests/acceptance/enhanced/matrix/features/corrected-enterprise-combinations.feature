Feature: Corrected Enterprise Combination Testing
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
      | fiber     | sqlite    | gorm  | session | logrus  |
      | chi       | postgres  | sqlx  | jwt     | slog    |

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
      | echo      | mysql    | sqlx | oauth2 | slog    | vite           |
      | fiber     | sqlite   | gorm | session| zerolog | esbuild        |

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
      | gin       | postgres | gorm | jwt     | zap    |
      | echo      | mysql    | sqlx | oauth2  | slog   |

  @critical @microservice
  Scenario Outline: Microservice architecture combinations
    When I generate a microservice project with configuration:
      | type         | microservice  |
      | framework    | <framework>   |
      | database     | <database>    |
      | orm          | <orm>         |
      | auth         | <auth>        |
      | logger       | <logger>      |
    Then the project should compile successfully
    And microservice patterns should be implemented
    And service communication should be properly configured
    And database connections should be optimized
    And health checks should be available

    Examples:
      | framework | database | orm  | auth | logger  |
      | gin       | postgres | gorm | jwt  | zap     |
      | echo      | mysql    | sqlx | oauth2| slog   |

  @critical @event-driven
  Scenario Outline: Event-driven architecture with CQRS
    When I generate an event-driven project with configuration:
      | type         | event-driven |
      | framework    | <framework>  |
      | database     | <database>   |
      | orm          | <orm>        |
      | auth         | <auth>       |
      | logger       | <logger>     |
    Then the project should compile successfully
    And CQRS pattern should be properly implemented
    And event sourcing should work correctly
    And command handlers should be separated from query handlers
    And event store should be configured properly
    And projection rebuilding should be supported

    Examples:
      | framework | database | orm  | auth | logger |
      | gin       | postgres | gorm | jwt  | zap    |
      | echo      | postgres | sqlx | oauth2| slog  |

  @critical @lambda
  Scenario Outline: AWS Lambda combinations
    When I generate a lambda project with configuration:
      | type         | <lambda_type> |
      | auth         | <auth>        |
      | logger       | <logger>      |
      | go_version   | <go_version>  |
    Then the project should compile successfully
    And lambda handler should be properly configured
    And AWS SDK dependencies should be correct
    And cold start should be optimized
    And logging should be CloudWatch compatible

    Examples:
      | lambda_type  | auth    | logger  | go_version |
      | lambda       | jwt     | zap     | 1.21       |
      | lambda-proxy | oauth2  | slog    | 1.22       |
      | lambda       | session | zerolog | 1.23       |

  @critical @workspace
  Scenario: Complex workspace with multiple blueprint types
    When I generate a workspace project with configuration:
      | type         | workspace |
      | name         | enterprise-system |
    And I add component "api" with configuration:
      | type         | web-api   |
      | architecture | hexagonal |
      | framework    | gin       |
      | database     | postgres  |
      | auth         | jwt       |
    And I add component "cli" with configuration:
      | type         | cli      |
      | complexity   | standard |
      | framework    | cobra    |
    And I add component "functions" with configuration:
      | type         | lambda   |
      | auth         | jwt      |
    Then the workspace should compile successfully
    And all components should share consistent go.mod dependencies
    And workspace-level configuration should be consistent

  @critical @cli-complexity
  Scenario Outline: CLI complexity combinations
    When I generate a CLI project with configuration:
      | type         | cli           |
      | complexity   | <complexity>  |
      | framework    | cobra         |
      | logger       | <logger>      |
      | go_version   | <go_version>  |
    Then the project should compile successfully
    And CLI structure should match <complexity> level
    And file count should be appropriate for <complexity>
    And command structure should be properly organized

    Examples:
      | complexity | logger  | go_version | expected_files |
      | simple     | slog    | 1.21       | ~8             |
      | standard   | zap     | 1.22       | ~29            |
      | simple     | zerolog | 1.23       | ~8             |
      | standard   | logrus  | 1.21       | ~29            |

  @critical @library
  Scenario Outline: Library project combinations
    When I generate a library project with configuration:
      | type       | library     |
      | logger     | <logger>    |
      | go_version | <go_version>|
    Then the project should compile successfully
    And public API should be well-defined
    And examples should be provided
    And documentation should be comprehensive
    And testing structure should be in place

    Examples:
      | logger  | go_version |
      | slog    | 1.21       |
      | zap     | 1.22       |
      | zerolog | 1.23       |
      | logrus  | 1.21       |

  @critical @monolith
  Scenario Outline: Monolith application combinations
    When I generate a monolith project with configuration:
      | type         | monolith    |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | <auth>      |
      | logger       | <logger>    |
    Then the project should compile successfully
    And modular structure should be maintained
    And all components should be integrated
    And shared dependencies should be managed correctly

    Examples:
      | framework | database | orm  | auth   | logger |
      | gin       | postgres | gorm | jwt    | zap    |
      | echo      | mysql    | sqlx | oauth2 | slog   |

  @critical @grpc-gateway
  Scenario Outline: gRPC Gateway combinations
    When I generate a gRPC Gateway project with configuration:
      | type         | grpc-gateway |
      | database     | <database>   |
      | orm          | <orm>        |
      | auth         | <auth>       |
      | logger       | <logger>     |
    Then the project should compile successfully
    And gRPC services should be properly defined
    And Gateway should proxy HTTP to gRPC correctly
    And protobuf definitions should be valid
    And both HTTP and gRPC endpoints should work

    Examples:
      | database | orm  | auth | logger |
      | postgres | gorm | jwt  | zap    |
      | mysql    | sqlx | oauth2| slog  |

  @critical @database-combinations
  Scenario Outline: Database and ORM combinations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | gin         |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | jwt         |
      | logger       | slog        |
    Then the project should compile successfully
    And database connection should be properly configured
    And ORM integration should work correctly
    And migrations should be available if supported
    And connection pooling should be optimized

    Examples:
      | database | orm  | notes                    |
      | postgres | gorm | Full ORM with migrations |
      | postgres | sqlx | Query builder approach   |
      | mysql    | gorm | Full ORM with migrations |
      | mysql    | sqlx | Query builder approach   |
      | sqlite   | gorm | File-based database      |
      | sqlite   | sqlx | File-based database      |

  @critical @auth-combinations
  Scenario Outline: Authentication system combinations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | postgres    |
      | orm          | gorm        |
      | auth         | <auth>      |
      | logger       | zap         |
    Then the project should compile successfully
    And <auth> authentication should be implemented
    And middleware should be properly configured
    And security headers should be set
    And token validation should work correctly

    Examples:
      | framework | auth    | notes                        |
      | gin       | jwt     | JSON Web Token authentication|
      | echo      | oauth2  | OAuth2 flow implementation   |
      | fiber     | session | Session-based authentication |
      | chi       | jwt     | JSON Web Token authentication|

  @critical @asset-pipeline
  Scenario Outline: Asset pipeline combinations
    When I generate a web API project with configuration:
      | type           | web-api          |
      | framework      | gin              |
      | asset_pipeline | <asset_pipeline> |
      | logger         | slog             |
    Then the project should compile successfully
    And asset pipeline should be configured for <asset_pipeline>
    And build scripts should be available
    And static assets should be served correctly

    Examples:
      | asset_pipeline | notes                           |
      | embedded       | Assets embedded in binary       |
      | webpack        | Webpack build configuration     |
      | vite           | Vite build configuration        |
      | esbuild        | ESBuild configuration           |