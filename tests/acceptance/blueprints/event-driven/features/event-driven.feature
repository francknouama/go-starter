Feature: Event-Driven Blueprint Generation
  As a developer building CQRS/Event Sourcing applications
  I want to generate event-driven Go projects with comprehensive CQRS support
  So that I can build scalable, maintainable event-driven systems

  Background:
    Given I have the go-starter CLI tool available
    And I am in a temporary working directory

  Scenario: Generate basic event-driven project
    Given I want to generate an event-driven blueprint
    When I run the generator with:
      | Parameter | Value |
      | type | event-driven |
      | name | basic-event-app |
      | module | github.com/test/basic-event-app |
      | framework | gin |
      | logger | slog |
      | database_type | postgres |
      | message_queue | redis |
    Then the generation should succeed
    And the generated project should have the event-driven structure:
      | Path | Type |
      | go.mod | file |
      | main.go | file |
      | internal/domain | directory |
      | internal/cqrs | directory |
      | internal/eventstore | directory |
      | internal/projections | directory |
      | internal/handlers | directory |
      | internal/config | directory |
      | cmd/api/main.go | file |
      | cmd/projector/main.go | file |
      | cmd/worker/main.go | file |
    And all modules should compile successfully

  Scenario: Generate event-driven project with different databases
    Given I want to generate an event-driven blueprint
    When I run the generator for each database:
      | Database |
      | postgres |
      | mysql |
      | mongodb |
    Then each generated project should:
      | Validation |
      | compile successfully |
      | include database-specific event store |
      | contain appropriate database drivers |
      | have correct connection configuration |

  Scenario: Generate event-driven project with different message queues
    Given I want to generate an event-driven blueprint
    When I run the generator for each message queue:
      | MessageQueue |
      | redis |
      | nats |
      | kafka |
      | rabbitmq |
    Then each generated project should:
      | Validation |
      | compile successfully |
      | include message queue event bus |
      | contain appropriate MQ client libraries |
      | have correct event publishing configuration |

  Scenario: Verify CQRS pattern implementation
    Given I have generated an event-driven project
    When I examine the CQRS implementation
    Then the project should include:
      | Component | Path |
      | Command interface | internal/cqrs/command.go |
      | Command bus | internal/cqrs/command_bus.go |
      | Query interface | internal/cqrs/query.go |
      | Query bus | internal/cqrs/query_bus.go |
      | Command handlers | internal/handlers/commands |
      | Query handlers | internal/handlers/queries |
    And the command bus should handle command dispatching
    And the query bus should handle query execution
    And commands and queries should be properly separated

  Scenario: Verify Event Sourcing implementation
    Given I have generated an event-driven project
    When I examine the Event Sourcing implementation
    Then the project should include:
      | Component | Path |
      | Aggregate root | internal/domain/aggregate.go |
      | Event interface | internal/domain/event.go |
      | Event store | internal/eventstore/store.go |
      | Event repository | internal/domain/repository.go |
      | Snapshot support | internal/eventstore/snapshots.go |
    And aggregates should apply events correctly
    And events should be persisted in the event store
    And snapshots should be supported for performance

  Scenario: Verify projection system
    Given I have generated an event-driven project
    When I examine the projection system
    Then the project should include:
      | Component | Path |
      | Projection interface | internal/projections/projection.go |
      | Projection manager | internal/projections/manager.go |
      | Event handlers | internal/handlers/events |
      | Read models | internal/models/read |
    And projections should update read models
    And event handlers should process domain events
    And the projector service should run continuously

  Scenario: Verify observability and monitoring
    Given I have generated an event-driven project
    When I examine the observability features
    Then the project should include:
      | Feature | Path |
      | Distributed tracing | internal/observability/tracing.go |
      | Metrics collection | internal/observability/metrics.go |
      | Structured logging | internal/logger |
      | Health checks | internal/health |
    And tracing should cover command/query flows
    And metrics should track performance
    And logging should be structured and contextual

  Scenario: Verify testing infrastructure
    Given I have generated an event-driven project
    When I examine the testing setup
    Then the project should include:
      | Test Type | Path |
      | Unit tests | *_test.go |
      | Integration tests | tests/integration |
      | CQRS tests | tests/cqrs |
      | Event sourcing tests | tests/eventsourcing |
      | End-to-end tests | tests/e2e |
    And unit tests should cover all components
    And integration tests should test event flows
    And CQRS tests should validate command/query separation
    And end-to-end tests should test complete scenarios

  Scenario: Verify API endpoints
    Given I have generated an event-driven project
    When I examine the API implementation
    Then the project should include:
      | Endpoint Type | Purpose |
      | Command endpoints | Execute commands |
      | Query endpoints | Execute queries |
      | Event webhooks | Receive external events |
      | Health endpoints | System health checks |
    And command endpoints should return 202 Accepted
    And query endpoints should return data immediately
    And webhooks should handle event ingestion
    And health checks should validate system status

  Scenario: Verify configuration management
    Given I have generated an event-driven project
    When I examine the configuration
    Then the project should support:
      | Configuration Type | Format |
      | Environment variables | .env |
      | Configuration files | YAML/JSON |
      | Command line flags | CLI args |
      | Kubernetes config | ConfigMaps |
    And configuration should include database settings
    And configuration should include message queue settings
    And configuration should include observability settings
    And configuration should be validated on startup

  Scenario: Verify deployment readiness
    Given I have generated an event-driven project
    When I examine the deployment configuration
    Then the project should include:
      | Deployment Type | Files |
      | Docker | Dockerfile, docker-compose.yml |
      | Kubernetes | k8s/*.yaml |
      | Helm | helm/chart.yaml |
      | CI/CD | .github/workflows/*.yml |
    And Docker containers should be optimized
    And Kubernetes manifests should be production-ready
    And Helm charts should support multiple environments
    And CI/CD should include testing and deployment

  Scenario: Performance and scalability validation
    Given I have generated an event-driven project
    When I test the performance characteristics
    Then the system should:
      | Metric | Requirement |
      | Command latency | < 50ms |
      | Query latency | < 10ms |
      | Event processing | > 1000 events/sec |
      | Memory usage | < 512MB baseline |
    And the system should handle concurrent commands
    And the system should scale horizontally
    And projections should handle event replay
    And snapshots should improve aggregate loading

  Scenario: Event versioning and migration
    Given I have generated an event-driven project
    When I examine event versioning support
    Then the project should include:
      | Feature | Implementation |
      | Event versioning | Version field in events |
      | Event migration | Migration handlers |
      | Schema evolution | Backward compatibility |
      | Upcasting | Event transformation |
    And old events should be readable
    And new event versions should be supported
    And migrations should be reversible
    And schema changes should not break existing data

  Scenario: Security and authorization
    Given I have generated an event-driven project
    When I examine the security implementation
    Then the project should include:
      | Security Feature | Implementation |
      | Authentication | JWT/OAuth2 |
      | Authorization | RBAC/ABAC |
      | Input validation | Command/query validation |
      | Audit logging | Security event tracking |
    And commands should be authorized
    And queries should respect permissions
    And sensitive data should be encrypted
    And audit trails should be maintained

  Scenario: Error handling and resilience
    Given I have generated an event-driven project
    When I examine error handling
    Then the project should include:
      | Resilience Pattern | Implementation |
      | Circuit breaker | Command/query protection |
      | Retry logic | Transient failure handling |
      | Bulkhead isolation | Service separation |
      | Graceful degradation | Fallback mechanisms |
    And errors should be properly logged
    And failures should not cascade
    And the system should recover automatically
    And error responses should be informative

  Scenario: Multi-service architecture support
    Given I want to generate a distributed event-driven system
    When I run the generator with microservices enabled:
      | Parameter | Value |
      | type | event-driven |
      | architecture | microservices |
      | services | user,order,payment,notification |
      | message_queue | kafka |
      | database_type | postgres |
    Then the generation should create:
      | Service | Purpose |
      | user-service | User management |
      | order-service | Order processing |
      | payment-service | Payment handling |
      | notification-service | Notifications |
    And each service should have its own event store
    And services should communicate via events
    And saga orchestration should be supported
    And service discovery should be configured