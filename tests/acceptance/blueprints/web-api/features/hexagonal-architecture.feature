Feature: Hexagonal Architecture Web API
  As a software architect
  I want to generate Hexagonal Architecture web API
  So that I can implement ports and adapters pattern

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Validate Hexagonal Architecture structure
    Given I want to create a Hexagonal Architecture web API
    When I run the command "go-starter new my-hex-api --type=web-api-hexagonal --framework=gin --database-driver=postgres --no-git"
    Then the generation should succeed
    And the project should have these layers:
      | Layer                | Directory                          | Purpose                        |
      | domain               | internal/domain/                   | Business entities & logic      |
      | application          | internal/application/              | Application services           |
      | primary_adapters     | internal/adapters/primary/         | Driving adapters (controllers) |
      | secondary_adapters   | internal/adapters/secondary/       | Driven adapters (persistence)  |
      | input_ports          | internal/application/ports/input/  | Primary ports                  |
      | output_ports         | internal/application/ports/output/ | Secondary ports                |
      | entities             | internal/domain/entities/          | Domain entities                |
      | services             | internal/application/services/     | Application services           |
      | persistence          | internal/adapters/secondary/persistence/ | Database adapters       |
      | http                 | internal/adapters/primary/http/    | HTTP adapters                  |
      | config               | internal/infrastructure/config/    | Configuration                  |
    And the generated code should compile successfully

  Scenario: Validate ports and adapters implementation
    Given I want to validate ports and adapters pattern
    When I generate a Hexagonal Architecture web API
    Then the project should have these ports:
      | Port Type    | Purpose                | Location                              |
      | input        | Primary ports          | internal/application/ports/input/     |
      | output       | Secondary ports        | internal/application/ports/output/    |
    And primary adapters should implement input ports
    And secondary adapters should implement output ports
    And application should depend only on ports
    And adapters should depend on application ports

  Scenario: Primary adapter validation (HTTP)
    Given I want to validate primary adapters
    When I generate a Hexagonal Architecture web API
    Then HTTP adapters should be in primary adapters
    And HTTP adapters should implement input ports
    And HTTP adapters should handle HTTP concerns only
    And HTTP adapters should not contain business logic
    And HTTP adapters should delegate to application services

  Scenario: Secondary adapter validation (Persistence)
    Given I want to validate secondary adapters
    When I generate a Hexagonal Architecture web API
    Then persistence adapters should be in secondary adapters
    And persistence adapters should implement output ports
    And persistence adapters should handle database concerns only
    And persistence adapters should not contain business logic
    And persistence adapters should be swappable

  Scenario: Application core isolation
    Given I want to validate application core isolation
    When I generate a Hexagonal Architecture web API
    Then the application core should not import adapters
    And the application core should define ports
    And the application core should contain business logic
    And the application core should be framework-independent
    And the application core should be database-independent

  Scenario: Port interface validation
    Given I want to validate port interfaces
    When I generate a Hexagonal Architecture web API
    Then input ports should define use cases
    And output ports should define external dependencies
    And ports should be technology-agnostic
    And ports should have clear contracts
    And ports should enable testability

  Scenario: Dependency direction validation
    Given I want to validate dependency directions
    When I generate a Hexagonal Architecture web API
    Then adapters should depend on application
    And application should depend on domain
    And domain should have no outward dependencies
    And ports should be defined by application
    And external concerns should be adapter responsibilities

  Scenario: Framework independence validation
    Given I want to validate framework independence
    When I generate a Hexagonal Architecture web API with framework "<framework>"
    Then the application core should not import framework
    And primary adapters should handle framework specifics
    And framework concerns should be isolated to adapters
    And business logic should be framework-agnostic
    And the generated code should compile successfully

    Examples:
      | framework |
      | gin       |
      | echo      |
      | fiber     |
      | chi       |

  Scenario: Database independence validation
    Given I want to validate database independence
    When I generate a Hexagonal Architecture web API with database "<database>" and ORM "<orm>"
    Then the application core should not import database libraries
    And secondary adapters should handle database concerns
    And output ports should abstract database operations
    And domain should be persistence-ignorant
    And the generated code should compile successfully

    Examples:
      | database | orm  |
      | postgres | gorm |
      | postgres | sqlx |
      | mysql    | gorm |
      | sqlite   | gorm |

  Scenario: Multiple adapter support
    Given I want to support multiple adapters for the same port
    When I generate a Hexagonal Architecture web API
    Then output ports should support multiple implementations
    And adapters should be easily swappable
    And configuration should determine active adapter
    And adapters should implement same interface
    And testing should support mock adapters

  Scenario: Testing strategy validation
    Given I want to validate testing approaches
    When I generate a Hexagonal Architecture web API
    Then unit tests should test application core in isolation
    And integration tests should test adapter implementations
    And acceptance tests should test through primary adapters
    And ports should enable easy mocking
    And business logic should be testable without adapters

  Scenario: Configuration and dependency injection
    Given I want to validate dependency injection
    When I generate a Hexagonal Architecture web API
    Then dependency injection should wire ports to adapters
    And application should receive injected dependencies
    And adapters should be configured externally
    And port implementations should be swappable
    And configuration should be centralized

  Scenario: Error handling across layers
    Given I want to validate error handling
    When I generate a Hexagonal Architecture web API
    Then domain errors should be expressed in business terms
    And application should handle use case errors
    And adapters should translate technical errors
    And ports should define error contracts
    And errors should flow correctly between layers

  Scenario: Cross-cutting concerns handling
    Given I want to validate cross-cutting concerns
    When I generate a Hexagonal Architecture web API
    Then logging should be handled by adapters
    And metrics should be adapter responsibility
    And security should be handled by primary adapters
    And transaction management should be secondary adapter concern
    And application core should remain pure

  Scenario: Domain-driven design integration
    Given I want to combine hexagonal architecture with DDD
    When I generate a Hexagonal Architecture web API with DDD patterns
    Then domain should contain entities and value objects
    And domain services should handle complex business logic
    And aggregates should maintain consistency boundaries
    And domain events should be supported
    And repositories should be output ports