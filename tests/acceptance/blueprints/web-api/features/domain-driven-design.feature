Feature: Domain-Driven Design Web API
  As a domain expert
  I want to generate DDD-structured web API
  So that code reflects business concepts

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Domain structure validation
    Given I want to create a DDD web API
    When I run the command "go-starter new my-ddd-api --type=web-api-ddd --framework=gin --database-driver=postgres --no-git"
    Then the generation should succeed
    And the project should have domain-centric structure:
      | Component         | Directory                          | Purpose                        |
      | domain            | internal/domain/                   | Domain models and logic        |
      | application       | internal/application/              | Application services           |
      | presentation      | internal/presentation/             | Controllers and DTOs           |
      | infrastructure    | internal/infrastructure/           | External concerns              |
      | shared            | internal/shared/                   | Shared utilities               |
      | entities          | internal/domain/entities/          | Domain entities                |
      | value_objects     | internal/domain/valueobjects/      | Value objects                  |
      | domain_services   | internal/domain/services/          | Domain services                |
      | repositories      | internal/domain/repositories/      | Repository interfaces          |
      | events            | internal/domain/events/            | Domain events                  |
    And the generated code should compile successfully

  Scenario: Business rule enforcement
    Given I want to create a DDD web API with business rules
    When I generate a DDD web API
    Then domain entities should enforce business rules
    And value objects should be immutable
    And domain services should handle complex logic
    And repositories should abstract data access
    And domain events should be properly implemented

  Scenario: Domain entity validation
    Given I want to validate domain entities
    When I generate a DDD web API
    Then entities should have identity
    And entities should encapsulate business logic
    And entities should maintain data consistency
    And entities should not depend on infrastructure
    And entities should communicate through domain events

  Scenario: Value object validation
    Given I want to validate value objects
    When I generate a DDD web API
    Then value objects should be immutable
    And value objects should validate their state
    And value objects should be comparable by value
    And value objects should not have identity
    And value objects should express domain concepts

  Scenario: Aggregate design validation
    Given I want to validate aggregate design
    When I generate a DDD web API
    Then aggregates should have clear boundaries
    And aggregate roots should control access
    And aggregates should maintain consistency
    And aggregates should communicate via domain events
    And cross-aggregate references should use IDs only

  Scenario: Domain service validation
    Given I want to validate domain services
    When I generate a DDD web API
    Then domain services should contain business logic
    And domain services should be stateless
    And domain services should coordinate entities
    And domain services should not depend on infrastructure
    And domain services should express domain concepts

  Scenario: Repository pattern validation
    Given I want to validate repository patterns
    When I generate a DDD web API
    Then repository interfaces should be in domain
    And repository implementations should be in infrastructure
    And repositories should abstract data access
    And repositories should work with aggregates
    And repositories should support domain queries

  Scenario: Application service validation
    Given I want to validate application services
    When I generate a DDD web API
    Then application services should orchestrate use cases
    And application services should be transaction boundaries
    And application services should not contain business logic
    And application services should coordinate domain objects
    And application services should handle cross-cutting concerns

  Scenario: Domain event validation
    Given I want to validate domain events
    When I generate a DDD web API
    Then domain events should be immutable
    And domain events should express business events
    And domain events should contain relevant data
    And domain events should be publishable
    And domain events should support eventual consistency

  Scenario: Anti-corruption layer validation
    Given I want to protect my domain from external systems
    When I generate a DDD web API with external integrations
    Then anti-corruption layers should translate external models
    And domain models should remain pure
    And external dependencies should be isolated
    And translation should preserve domain invariants
    And integration should not leak into domain

  Scenario: Bounded context validation
    Given I want to validate bounded contexts
    When I generate a DDD web API
    Then bounded contexts should have clear boundaries
    And contexts should have their own models
    And cross-context communication should be explicit
    And shared concepts should be carefully managed
    And context mapping should be documented

  Scenario: Ubiquitous language validation
    Given I want to ensure ubiquitous language
    When I generate a DDD web API
    Then code should reflect business terminology
    And method names should use domain language
    And class names should express domain concepts
    And comments should align with business understanding
    And documentation should use consistent terminology

  Scenario: Domain model validation with different databases
    Given I want to validate domain persistence
    When I generate a DDD web API with database "<database>" and ORM "<orm>"
    Then domain models should be persistence-ignorant
    And repositories should abstract database details
    And domain logic should not depend on ORM
    And data mapping should preserve domain rules
    And the generated code should compile successfully

    Examples:
      | database | orm  |
      | postgres | gorm |
      | postgres | sqlx |
      | mysql    | gorm |
      | mongodb  | mongo |

  Scenario: DDD with different frameworks
    Given I want to validate framework independence
    When I generate a DDD web API with framework "<framework>"
    Then domain layer should not import framework
    And presentation layer should adapt to framework
    And application layer should be framework-agnostic
    And infrastructure should handle framework specifics
    And the generated code should compile successfully

    Examples:
      | framework |
      | gin       |
      | echo      |
      | fiber     |
      | chi       |