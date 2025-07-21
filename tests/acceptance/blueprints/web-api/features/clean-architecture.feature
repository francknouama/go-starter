Feature: Clean Architecture Web API
  As a software architect
  I want to generate Clean Architecture web API
  So that I can maintain separation of concerns

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Validate Clean Architecture layers
    Given I want to create a Clean Architecture web API
    When I run the command "go-starter new my-clean-api --type=web-api-clean --framework=gin --database-driver=postgres --no-git"
    Then the generation should succeed
    And the project should have these layers:
      | Layer        | Directory                          | Purpose                    |
      | entities     | internal/domain/entities/          | Business entities          |
      | usecases     | internal/domain/usecases/          | Business logic             |
      | ports        | internal/domain/ports/             | Contracts and adapters     |
      | controllers  | internal/adapters/controllers/     | Interface adapters         |
      | presenters   | internal/adapters/presenters/      | Output adapters            |
      | persistence  | internal/infrastructure/persistence/ | Data access layer          |
      | web          | internal/infrastructure/web/       | Web framework layer        |
      | services     | internal/infrastructure/services/  | External services          |
      | logger       | internal/infrastructure/logger/    | Logging infrastructure     |
      | config       | internal/infrastructure/config/    | Configuration              |
      | container    | internal/infrastructure/container/ | Dependency injection       |
    And dependencies should only point inward
    And business logic should be framework-independent
    And interfaces should define contracts clearly

  Scenario: Dependency injection validation
    Given I want to create a Clean Architecture web API with dependency injection
    When I run the command "go-starter new clean-di-api --type=web-api-clean --framework=gin --database-driver=postgres --no-git"
    Then the generation should succeed
    And dependency injection should be configured
    And repositories should be interfaces
    And use cases should depend on interfaces only
    And frameworks should implement interfaces
    And the generated code should compile successfully

  Scenario: Logger integration follows Clean Architecture patterns
    Given I want to create a Clean Architecture web API with logger integration
    When I generate a web-api-clean with logger "<logger>"
    Then the generation should succeed
    And the logger should be in infrastructure layer
    And the logger should be injected through interfaces
    And business logic should not depend on concrete logger
    And the generated code should compile successfully

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  Scenario: Framework abstraction validation
    Given I want to create a Clean Architecture web API with framework abstraction
    When I generate a web-api-clean with framework "<framework>"
    Then the generation should succeed
    And the framework should be abstracted in infrastructure
    And business logic should not import framework
    And the generated code should compile successfully

    Examples:
      | framework |
      | gin       |
      | echo      |
      | fiber     |
      | chi       |

  Scenario: Database integration follows repository pattern
    Given I want to create a Clean Architecture web API with database integration
    When I generate a web-api-clean with database "<database>" and ORM "<orm>"
    Then the generation should succeed
    And the database should be in infrastructure layer
    And repository interfaces should be in domain
    And business logic should not depend on database specifics
    And the generated code should compile successfully

    Examples:
      | database | orm  |
      | postgres | gorm |
      | postgres | sqlx |
      | mysql    | gorm |

  Scenario: Architecture compliance validation
    Given I want to create a Clean Architecture web API that follows all principles
    When I run the command "go-starter new clean-compliant-api --type=web-api-clean --framework=gin --database-driver=postgres --no-git"
    Then the generation should succeed
    And the code should follow Clean Architecture principles
    And dependency directions should be correct
    And layer boundaries should be enforced
    And the architecture should be independently testable
    And the generated code should compile successfully

  Scenario: Business logic isolation
    Given I want to ensure business logic is isolated from external concerns
    When I generate a Clean Architecture web API
    Then entities should not import from outer layers
    And use cases should not import infrastructure or adapters
    And business logic should be framework-independent
    And business logic should be database-independent
    And business logic should be UI-independent

  Scenario: Interface contracts validation
    Given I want to ensure proper interface design
    When I generate a Clean Architecture web API
    Then ports directory should contain interface definitions
    And interfaces should define clear contracts
    And use cases should depend on interfaces only
    And infrastructure should implement interfaces
    And interface segregation should be applied

  Scenario: Dependency inversion principle
    Given I want to validate dependency inversion
    When I generate a Clean Architecture web API
    Then high-level modules should not depend on low-level modules
    And both should depend on abstractions
    And abstractions should not depend on details
    And details should depend on abstractions
    And dependency injection should be configured properly