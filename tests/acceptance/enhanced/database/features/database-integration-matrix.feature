Feature: Database Integration Matrix Testing
  As a go-starter user building data-driven applications
  I want to ensure that all database and ORM combinations work correctly
  So that data persistence layers are reliable and optimized

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And I am testing database integration combinations

  @critical @database @postgres @p0
  Scenario Outline: PostgreSQL database integration combinations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | postgres    |
      | orm          | <orm>       |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And PostgreSQL driver should be properly configured
    And <orm> ORM integration should work correctly
    And database connections should be properly managed
    And migration support should be available for <orm>
    And connection pooling should be optimized
    And database health checks should be implemented
    And transaction management should work correctly
    And <auth> authentication should integrate with database
    And database queries should be secure against injection

    Examples:
      | framework | orm  | auth    | logger  |
      | gin       | gorm | jwt     | zap     |
      | gin       | sqlx | jwt     | slog    |
      | echo      | gorm | oauth2  | zerolog |
      | echo      | sqlx | oauth2  | logrus  |
      | fiber     | gorm | session | zap     |
      | fiber     | sqlx | session | slog    |

  @critical @database @mysql @p0
  Scenario Outline: MySQL database integration combinations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | mysql       |
      | orm          | <orm>       |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And MySQL driver should be properly configured
    And <orm> ORM integration should work correctly
    And database connections should be properly managed
    And migration support should be available for <orm>
    And connection pooling should be optimized
    And database health checks should be implemented
    And transaction management should work correctly
    And charset should be properly configured for MySQL
    And timezone handling should work correctly

    Examples:
      | framework | orm  | auth    | logger  |
      | gin       | gorm | jwt     | zap     |
      | gin       | sqlx | jwt     | slog    |
      | echo      | gorm | oauth2  | zerolog |
      | echo      | sqlx | oauth2  | logrus  |
      | chi       | gorm | session | zap     |

  @critical @database @sqlite @p0
  Scenario Outline: SQLite database integration combinations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | sqlite      |
      | orm          | <orm>       |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And SQLite driver should be properly configured
    And <orm> ORM integration should work correctly
    And file-based database should be properly managed
    And migration support should be available for <orm>
    And WAL mode should be configured for performance
    And foreign key constraints should be enabled
    And transaction management should work correctly
    And database file permissions should be secure
    And backup strategies should be documented

    Examples:
      | framework | orm  | auth    | logger  |
      | gin       | gorm | jwt     | zap     |
      | gin       | sqlx | jwt     | slog    |
      | echo      | gorm | oauth2  | zerolog |
      | fiber     | sqlx | session | logrus  |

  @integration @database @gorm
  Scenario Outline: GORM ORM comprehensive testing
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | gorm        |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And GORM should be properly configured
    And model definitions should follow GORM conventions
    And migrations should be automatically generated
    And associations should be properly defined
    And query optimization should be enabled
    And soft deletes should be supported
    And hooks and callbacks should be available
    And connection pooling should be configured
    And performance monitoring should be enabled

    Examples:
      | framework | database | auth    | logger  |
      | gin       | postgres | jwt     | zap     |
      | echo      | mysql    | oauth2  | slog    |
      | fiber     | sqlite   | session | zerolog |

  @integration @database @sqlx
  Scenario Outline: SQLX query builder comprehensive testing
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | sqlx        |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And SQLX should be properly configured
    And raw SQL queries should be safely handled
    And prepared statements should be used
    And result scanning should work correctly
    And transaction management should be explicit
    And query logging should be available
    And performance should be optimized
    And SQL injection prevention should be implemented
    And custom types should be supported

    Examples:
      | framework | database | auth    | logger |
      | gin       | postgres | jwt     | zap    |
      | echo      | mysql    | oauth2  | slog   |
      | chi       | sqlite   | session | logrus |

  @integration @database @migration
  Scenario Outline: Database migration testing
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | gin         |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | jwt         |
      | logger       | slog        |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And migration system should be properly configured
    And initial migration should be created
    And migration versioning should be implemented
    And rollback capabilities should be available
    And migration status tracking should work
    And database schema should be version controlled
    And migration execution should be idempotent
    And data migrations should be supported

    Examples:
      | database | orm  |
      | postgres | gorm |
      | postgres | sqlx |
      | mysql    | gorm |
      | mysql    | sqlx |
      | sqlite   | gorm |

  @performance @database @connection-pooling
  Scenario Outline: Database connection pooling optimization
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | logger       | zap         |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And connection pool should be properly configured
    And max open connections should be optimized
    And max idle connections should be set
    And connection lifetime should be managed
    And connection retry logic should be implemented
    And pool metrics should be available
    And connection leaks should be prevented
    And performance should be monitored

    Examples:
      | framework | database | orm  |
      | gin       | postgres | gorm |
      | echo      | mysql    | sqlx |
      | fiber     | postgres | gorm |

  @security @database @injection-prevention
  Scenario Outline: SQL injection prevention testing
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | jwt         |
      | logger       | slog        |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And SQL injection prevention should be implemented
    And parameterized queries should be used
    And input validation should be strict
    And query sanitization should be applied
    And ORM should prevent SQL injection
    And error messages should not leak schema information
    And query logging should be secure
    And database permissions should be minimal

    Examples:
      | framework | database | orm  |
      | gin       | postgres | gorm |
      | echo      | mysql    | sqlx |
      | fiber     | sqlite   | gorm |

  @integration @database @transaction
  Scenario Outline: Transaction management testing
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | jwt         |
      | logger       | slog        |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And transaction management should be properly implemented
    And ACID properties should be maintained
    And nested transactions should be handled correctly
    And transaction timeouts should be configured
    And deadlock detection should be available
    And rollback on error should work correctly
    And transaction isolation levels should be configurable
    And distributed transactions should be supported if applicable

    Examples:
      | framework | database | orm  |
      | gin       | postgres | gorm |
      | echo      | mysql    | sqlx |
      | fiber     | postgres | gorm |

  @monitoring @database @health-checks
  Scenario Outline: Database health monitoring
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | logger       | zap         |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And database health checks should be implemented
    And connection status should be monitored
    And query performance should be tracked
    And slow query detection should be available
    And database metrics should be exposed
    And alerting should be configured for issues
    And health endpoint should include database status
    And recovery mechanisms should be implemented

    Examples:
      | framework | database | orm  |
      | gin       | postgres | gorm |
      | echo      | mysql    | sqlx |

  @integration @database @backup
  Scenario Outline: Database backup and recovery
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | gin         |
      | database     | <database>  |
      | orm          | gorm        |
      | logger       | slog        |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And backup strategies should be documented
    And point-in-time recovery should be supported
    And backup automation should be available
    And backup verification should be implemented
    And restore procedures should be documented
    And backup retention policies should be configured
    And disaster recovery plans should be included

    Examples:
      | database |
      | postgres |
      | mysql    |

  @performance @database @optimization
  Scenario Outline: Database performance optimization
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | logger       | zap         |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And query optimization should be implemented
    And indexing strategies should be documented
    And query caching should be available
    And performance profiling should be enabled
    And slow query logging should be configured
    And query plan analysis should be available
    And database statistics should be maintained
    And performance tuning guides should be included

    Examples:
      | framework | database | orm  |
      | gin       | postgres | gorm |
      | echo      | mysql    | sqlx |

  # P1: Enhanced Database Integration Matrix - Additional Coverage Scenarios
  
  @critical @database @framework-consistency @p1
  Scenario Outline: Framework consistency across database combinations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | jwt         |
      | logger       | slog        |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And <framework> framework should be properly integrated
    And <database> driver should be properly configured
    And <orm> ORM integration should work correctly
    And middleware integration should work consistently
    And routing patterns should follow <framework> conventions
    And error handling should be framework-specific
    And configuration loading should work with <framework>
    And health check endpoints should be properly configured

    Examples:
      | framework | database | orm  |
      | gin       | postgres | gorm |
      | gin       | mysql    | gorm |
      | gin       | sqlite   | gorm |
      | gin       | postgres | sqlx |
      | gin       | mysql    | sqlx |
      | gin       | sqlite   | sqlx |
      | echo      | postgres | gorm |
      | echo      | mysql    | gorm |
      | echo      | sqlite   | gorm |
      | echo      | postgres | sqlx |
      | echo      | mysql    | sqlx |
      | echo      | sqlite   | sqlx |
      | fiber     | postgres | gorm |
      | fiber     | mysql    | gorm |
      | fiber     | sqlite   | gorm |
      | fiber     | postgres | sqlx |
      | fiber     | mysql    | sqlx |
      | fiber     | sqlite   | sqlx |
      | chi       | postgres | gorm |
      | chi       | mysql    | gorm |
      | chi       | sqlite   | gorm |
      | chi       | postgres | sqlx |
      | chi       | mysql    | sqlx |
      | chi       | sqlite   | sqlx |

  @integration @database @architecture-matrix @p1 
  Scenario Outline: Database integration across different architectures
    When I generate a web API project with configuration:
      | type         | web-api        |
      | framework    | gin            |
      | architecture | <architecture> |
      | database     | <database>     |
      | orm          | <orm>          |
      | auth         | jwt            |
      | logger       | zap            |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And <architecture> architecture should be properly implemented
    And database layer should be correctly placed in architecture
    And dependency injection should work correctly
    And repository patterns should follow <architecture> principles
    And service layer should integrate properly with database
    And domain models should be architecture-appropriate
    And data access should respect architectural boundaries

    Examples:
      | architecture | database | orm  |
      | standard     | postgres | gorm |
      | standard     | mysql    | sqlx |
      | standard     | sqlite   | gorm |
      | clean        | postgres | gorm |
      | clean        | mysql    | sqlx |
      | clean        | sqlite   | gorm |
      | ddd          | postgres | gorm |
      | ddd          | mysql    | sqlx |
      | ddd          | sqlite   | gorm |
      | hexagonal    | postgres | gorm |
      | hexagonal    | mysql    | sqlx |
      | hexagonal    | sqlite   | gorm |

  @integration @database @comprehensive-auth @p1
  Scenario Outline: Database integration with comprehensive authentication
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And <auth> authentication should integrate with database
    And user session management should work with <database>
    And authentication middleware should be properly configured
    And password hashing should be secure and database-compatible
    And role-based access control should work with ORM
    And authentication tokens should be database-backed if applicable
    And user account management should use proper database patterns
    And security audit trails should be implemented

    Examples:
      | framework | database | orm  | auth    | logger  |
      | gin       | postgres | gorm | jwt     | zap     |
      | gin       | postgres | sqlx | jwt     | slog    |
      | gin       | mysql    | gorm | jwt     | zerolog |
      | gin       | mysql    | sqlx | jwt     | logrus  |
      | gin       | sqlite   | gorm | jwt     | zap     |
      | gin       | sqlite   | sqlx | jwt     | slog    |
      | echo      | postgres | gorm | oauth2  | zap     |
      | echo      | postgres | sqlx | oauth2  | slog    |
      | echo      | mysql    | gorm | oauth2  | zerolog |
      | echo      | mysql    | sqlx | oauth2  | logrus  |
      | echo      | sqlite   | gorm | oauth2  | zap     |
      | echo      | sqlite   | sqlx | oauth2  | slog    |
      | fiber     | postgres | gorm | session | zap     |
      | fiber     | postgres | sqlx | session | slog    |
      | fiber     | mysql    | gorm | session | zerolog |
      | fiber     | mysql    | sqlx | session | logrus  |
      | fiber     | sqlite   | gorm | session | zap     |
      | fiber     | sqlite   | sqlx | session | slog    |

  @performance @database @stress-testing @p1
  Scenario Outline: Database performance under stress conditions
    When I generate a web API project optimized for performance:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | logger       | zap         |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And connection pooling should handle high concurrency
    And query performance should be optimized for large datasets
    And memory usage should be efficiently managed
    And connection timeouts should be properly configured
    And resource cleanup should prevent memory leaks
    And performance metrics should be exposed
    And database bottlenecks should be identifiable
    And scaling strategies should be documented

    Examples:
      | framework | database | orm  |
      | gin       | postgres | gorm |
      | gin       | postgres | sqlx |
      | gin       | mysql    | gorm |
      | gin       | mysql    | sqlx |
      | gin       | sqlite   | gorm |
      | gin       | sqlite   | sqlx |
      | echo      | postgres | gorm |
      | echo      | mysql    | sqlx |
      | fiber     | postgres | gorm |
      | fiber     | mysql    | sqlx |

  @security @database @advanced-security @p1
  Scenario Outline: Advanced database security testing
    When I generate a web API project with security focus:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | jwt         |
      | logger       | slog        |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And database connections should use TLS encryption
    And SQL injection attacks should be prevented at all layers
    And database credentials should be securely managed
    And query logging should not expose sensitive data
    And database access should be role-based and minimal
    And data encryption at rest should be supported
    And audit logging should track all database operations
    And security headers should be properly configured
    And data validation should prevent malicious input

    Examples:
      | framework | database | orm  |
      | gin       | postgres | gorm |
      | gin       | postgres | sqlx |
      | gin       | mysql    | gorm |
      | gin       | mysql    | sqlx |
      | echo      | postgres | gorm |
      | echo      | mysql    | sqlx |
      | fiber     | postgres | gorm |
      | fiber     | mysql    | sqlx |

  @integration @database @production-readiness @p1
  Scenario Outline: Production deployment database readiness
    When I generate a project ready for production deployment:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | jwt         |
      | logger       | zap         |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And database migrations should be production-safe
    And connection pooling should be optimized for production load
    And monitoring and alerting should be comprehensive
    And backup and recovery procedures should be automated
    And database configuration should be environment-specific
    And health checks should detect database issues
    And performance monitoring should track key metrics
    And disaster recovery plans should be documented
    And scaling strategies should be implementation-ready

    Examples:
      | framework | database | orm  |
      | gin       | postgres | gorm |
      | gin       | postgres | sqlx |
      | gin       | mysql    | gorm |
      | gin       | mysql    | sqlx |
      | echo      | postgres | gorm |
      | echo      | mysql    | sqlx |
      | fiber     | postgres | gorm |
      | fiber     | sqlite   | gorm |

  @integration @database @cross-platform @p1
  Scenario Outline: Cross-platform database compatibility
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | gin         |
      | database     | <database>  |
      | orm          | <orm>       |
      | logger       | slog        |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And database drivers should be cross-platform compatible
    And file paths should work across operating systems
    And database configuration should be platform-agnostic
    And connection strings should handle platform differences
    And migration scripts should work on all platforms
    And performance characteristics should be documented per platform
    And installation instructions should cover all platforms

    Examples:
      | database | orm  |
      | postgres | gorm |
      | postgres | sqlx |
      | mysql    | gorm |
      | mysql    | sqlx |
      | sqlite   | gorm |
      | sqlite   | sqlx |