Feature: Architecture Pattern Validation
  As a developer using go-starter
  I want to ensure that different architecture patterns are correctly implemented
  So that generated projects follow the intended architectural principles

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized

  Scenario: Standard Architecture Validation
    Given I generate a project with architecture "standard"
    When I analyze the project structure
    Then the project should follow standard layered architecture:
      | Layer       | Location                    | Purpose                           |
      | handlers    | internal/handlers/          | HTTP request handling             |
      | services    | internal/services/          | Business logic                    |
      | repository  | internal/repository/        | Data access layer                 |
      | models      | internal/models/            | Domain models                     |
      | middleware  | internal/middleware/        | Cross-cutting concerns            |
      | config      | internal/config/            | Configuration management          |
    And dependencies should flow downward only
    And handlers should depend on services
    And services should depend on repositories
    And repositories should not depend on handlers or services

  Scenario: Clean Architecture Validation
    Given I generate a project with architecture "clean"
    When I analyze the project structure
    Then the project should follow clean architecture principles:
      | Layer          | Location                      | Dependency Rule                   |
      | entities       | internal/domain/entities/     | No dependencies                   |
      | usecases       | internal/domain/usecases/     | Depends only on entities          |
      | controllers    | internal/adapters/controllers/| Depends on usecases               |
      | repositories   | internal/infrastructure/      | Implements domain interfaces      |
      | presenters     | internal/adapters/presenters/ | Formats output data               |
    And the dependency rule should be strictly enforced
    And business logic should be in usecases layer
    And framework dependencies should only exist in outer layers
    And domain entities should be pure Go structs

  Scenario: Domain-Driven Design (DDD) Architecture Validation
    Given I generate a project with architecture "ddd"
    When I analyze the project structure
    Then the project should follow DDD principles:
      | Component        | Location                        | Characteristics                   |
      | aggregates       | internal/domain/*/entity.go     | Root entities with invariants     |
      | value objects    | internal/domain/*/value_objects.go | Immutable domain concepts      |
      | domain services  | internal/domain/*/service.go    | Domain logic spanning entities    |
      | specifications   | internal/domain/*/specifications.go | Business rules               |
      | repositories     | internal/domain/*/repository.go | Domain interfaces                 |
      | application      | internal/application/           | Use cases and DTOs                |
      | events           | internal/domain/*/events.go     | Domain events                     |
    And aggregates should encapsulate business invariants
    And value objects should be immutable
    And domain services should contain cross-entity logic
    And application layer should orchestrate domain objects

  Scenario: Hexagonal Architecture Validation
    Given I generate a project with architecture "hexagonal"
    When I analyze the project structure
    Then the project should follow hexagonal architecture (ports and adapters):
      | Component      | Location                          | Role                              |
      | domain         | internal/domain/                  | Core business logic               |
      | ports          | internal/application/ports/       | Interfaces (primary & secondary)  |
      | primary adapters | internal/adapters/primary/      | HTTP, CLI, gRPC handlers          |
      | secondary adapters | internal/adapters/secondary/  | DB, external services             |
      | application    | internal/application/             | Use cases implementation          |
    And domain should have no external dependencies
    And all external interactions should go through ports
    And adapters should implement port interfaces
    And the core should be testable in isolation

  Scenario Outline: Architecture-specific import validation
    Given I generate a project with architecture "<architecture>"
    When I check import statements across the project
    Then imports should follow "<architecture>" dependency rules
    And there should be no circular dependencies
    And external framework dependencies should be isolated
    And domain/core layers should not import infrastructure

    Examples:
      | architecture |
      | standard     |
      | clean        |
      | ddd          |
      | hexagonal    |

  Scenario Outline: Architecture consistency with different frameworks
    Given I generate a project with architecture "<architecture>"
    And framework "<framework>"
    When I analyze the integration points
    Then the framework should be properly isolated in the appropriate layer
    And business logic should not depend on the framework
    And framework-specific code should be in adapters/handlers layer
    And switching frameworks should not affect core business logic

    Examples:
      | architecture | framework |
      | clean        | gin       |
      | clean        | echo      |
      | clean        | fiber     |
      | hexagonal    | gin       |
      | hexagonal    | echo      |
      | ddd          | gin       |
      | ddd          | fiber     |

  Scenario: Architecture metrics and quality checks
    Given I generate projects with all architecture patterns
    When I measure architecture quality metrics
    Then each architecture should meet its quality criteria:
      | Architecture | Max Package Coupling | Max Cyclomatic Complexity | Test Coverage Target |
      | standard     | 5                   | 10                        | 70%                  |
      | clean        | 3                   | 8                         | 80%                  |
      | ddd          | 3                   | 8                         | 85%                  |
      | hexagonal    | 2                   | 7                         | 85%                  |
    And dependency graphs should be acyclic
    And package cohesion should be high
    And architectural fitness functions should pass

  Scenario: Cross-architecture feature implementation comparison
    Given I implement the same feature across all architectures:
      | Feature         | Description                          |
      | User Creation   | Create user with validation          |
      | Authentication  | JWT-based authentication             |
      | Data Retrieval  | Get user by ID with caching          |
    When I compare the implementations
    Then each architecture should implement the feature correctly
    And the code organization should reflect architectural principles
    And test patterns should match the architecture style
    And performance characteristics should be documented

  Scenario: Architecture migration readiness
    Given I have a project with architecture "standard"
    When I assess migration readiness to other architectures
    Then I should identify refactoring requirements:
      | Target Architecture | Major Changes Required               | Risk Level |
      | clean              | Extract interfaces, separate layers   | Medium     |
      | ddd                | Identify aggregates, add value objects| High       |
      | hexagonal          | Define ports, create adapters        | Medium     |
    And migration guides should be available
    And breaking changes should be documented
    And incremental migration paths should exist