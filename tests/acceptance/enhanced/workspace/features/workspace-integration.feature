Feature: Workspace Integration Testing
  As a go-starter user building complex multi-component systems
  I want to ensure that workspace projects with multiple blueprints integrate correctly
  So that enterprise monorepo and microservice architectures work reliably

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And I am testing workspace integration scenarios

  @critical @workspace @p0
  Scenario: Enterprise workspace with web API, CLI, and Lambda components
    When I generate a workspace project with configuration:
      | type         | workspace           |
      | name         | enterprise-system   |
      | go_version   | 1.23               |
      | output_dir   | ./test-workspace    |
    And I add component "api-service" with configuration:
      | type         | web-api    |
      | architecture | hexagonal  |
      | framework    | gin        |
      | database     | postgres   |
      | orm          | gorm       |
      | auth         | jwt        |
      | logger       | zap        |
    And I add component "admin-cli" with configuration:
      | type         | cli        |
      | complexity   | standard   |
      | framework    | cobra      |
      | logger       | zap        |
    And I add component "worker-functions" with configuration:
      | type         | lambda     |
      | auth         | jwt        |
      | logger       | zap        |
    Then the workspace should be generated successfully
    And all components should share a consistent root go.mod
    And shared dependencies should be properly managed
    And component go.mod files should be properly configured
    And workspace compilation should succeed
    And inter-component imports should work correctly
    And shared configuration should be available to all components

  @critical @workspace @p0
  Scenario: Microservice workspace with multiple databases
    When I generate a workspace project with configuration:
      | type         | workspace              |
      | name         | microservice-platform  |
      | go_version   | 1.23                  |
    And I add component "user-service" with configuration:
      | type         | microservice |
      | architecture | clean        |
      | framework    | echo         |
      | database     | postgres     |
      | orm          | gorm         |
      | auth         | jwt          |
      | logger       | slog         |
    And I add component "payment-service" with configuration:
      | type         | microservice |
      | architecture | hexagonal    |
      | framework    | gin          |
      | database     | mysql        |
      | orm          | sqlx         |
      | auth         | oauth2       |
      | logger       | slog         |
    And I add component "notification-service" with configuration:
      | type         | microservice |
      | architecture | standard     |
      | framework    | fiber        |
      | database     | postgres     |
      | orm          | gorm         |
      | auth         | jwt          |
      | logger       | slog         |
    Then the workspace should be generated successfully
    And each service should have independent database configurations
    And shared authentication libraries should be consistent
    And logger implementations should be consistent across services
    And workspace compilation should succeed
    And service-to-service communication patterns should be available

  @critical @workspace @p0
  Scenario: Monolith workspace with modular components
    When I generate a workspace project with configuration:
      | type         | workspace    |
      | name         | modular-app  |
      | go_version   | 1.23        |
    And I add component "core-api" with configuration:
      | type         | web-api     |
      | architecture | ddd         |
      | framework    | gin         |
      | database     | postgres    |
      | orm          | gorm        |
      | auth         | session     |
      | logger       | zerolog     |
    And I add component "admin-panel" with configuration:
      | type         | web-api     |
      | architecture | standard    |
      | framework    | gin         |
      | database     | postgres    |
      | orm          | gorm        |
      | auth         | session     |
      | logger       | zerolog     |
    And I add component "utility-cli" with configuration:
      | type         | cli         |
      | complexity   | simple      |
      | framework    | cobra       |
      | logger       | zerolog     |
    Then the workspace should be generated successfully
    And shared database connections should be properly managed
    And session management should be consistent across web components
    And logger configuration should be shared
    And workspace compilation should succeed
    And modular structure should be maintained

  @critical @workspace @p0
  Scenario: Event-driven workspace with CQRS components
    When I generate a workspace project with configuration:
      | type         | workspace       |
      | name         | event-platform  |
      | go_version   | 1.23           |
    And I add component "command-service" with configuration:
      | type         | event-driven |
      | framework    | gin          |
      | database     | postgres     |
      | orm          | gorm         |
      | auth         | jwt          |
      | logger       | zap          |
    And I add component "query-service" with configuration:
      | type         | web-api      |
      | architecture | clean        |
      | framework    | echo         |
      | database     | postgres     |
      | orm          | sqlx         |
      | auth         | jwt          |
      | logger       | zap          |
    And I add component "event-processor" with configuration:
      | type         | lambda       |
      | auth         | jwt          |
      | logger       | zap          |
    Then the workspace should be generated successfully
    And event sourcing patterns should be properly implemented
    And CQRS separation should be maintained
    And event store configuration should be shared
    And workspace compilation should succeed
    And event flow between components should be configured

  @critical @workspace @p0
  Scenario: Complex workspace with all blueprint types
    When I generate a workspace project with configuration:
      | type         | workspace              |
      | name         | comprehensive-platform |
      | go_version   | 1.23                  |
    And I add component "web-api" with configuration:
      | type         | web-api    |
      | architecture | hexagonal  |
      | framework    | gin        |
      | database     | postgres   |
      | orm          | gorm       |
      | auth         | oauth2     |
      | logger       | logrus     |
    And I add component "grpc-service" with configuration:
      | type         | grpc-gateway |
      | database     | mysql        |
      | orm          | sqlx         |
      | auth         | jwt          |
      | logger       | logrus       |
    And I add component "management-cli" with configuration:
      | type         | cli        |
      | complexity   | standard   |
      | framework    | cobra      |
      | logger       | logrus     |
    And I add component "shared-library" with configuration:
      | type         | library    |
      | logger       | logrus     |
    And I add component "lambda-functions" with configuration:
      | type         | lambda-proxy |
      | auth         | oauth2       |
      | logger       | logrus       |
    And I add component "event-system" with configuration:
      | type         | event-driven |
      | framework    | echo         |
      | database     | postgres     |
      | orm          | gorm         |
      | auth         | oauth2       |
      | logger       | logrus       |
    Then the workspace should be generated successfully
    And all blueprint types should be properly integrated
    And shared dependencies should be managed efficiently
    And logger configuration should be consistent across all components
    And authentication should be compatible across components
    And workspace compilation should succeed
    And inter-component communication should be properly configured
    And shared library should be usable by all components

  @integration @workspace
  Scenario: Workspace dependency management validation
    When I generate a workspace with multiple components having different dependency requirements
    Then go.mod dependency resolution should work correctly
    And no version conflicts should exist
    And shared dependencies should use consistent versions
    And workspace should build without dependency errors

  @integration @workspace
  Scenario: Workspace configuration consistency validation
    When I generate a workspace with multiple components
    Then configuration patterns should be consistent across components
    And environment variable handling should be standardized
    And logging configuration should be centralized
    And database connection management should be optimized

  @performance @workspace
  Scenario: Large workspace performance validation
    When I generate a workspace with 8+ components
    Then generation time should be reasonable
    And memory usage should remain within acceptable limits
    And compilation time should be optimized
    And workspace structure should be maintainable