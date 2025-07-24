Feature: Configuration Consistency Validation
  As a developer using go-starter
  I want all configuration files to be consistent with my selections
  So that my project has coherent configuration across all components

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized

  Scenario Outline: go.mod dependencies match selected features
    Given I generate a project with configuration:
      | framework       | <framework>       |
      | database_driver | <database_driver> |
      | orm             | <orm>             |
      | logger          | <logger>          |
      | auth_type       | <auth_type>       |
    When I examine the go.mod file
    Then it should contain the framework dependency "<expected_framework_dep>"
    And it should contain the database driver dependency "<expected_db_dep>"
    And it should contain the ORM dependency "<expected_orm_dep>"
    And it should contain the logger dependency "<expected_logger_dep>"
    And it should contain the auth dependency "<expected_auth_dep>"
    And it should not contain conflicting dependencies

    Examples:
      | framework | database_driver | orm  | logger | auth_type | expected_framework_dep   | expected_db_dep                | expected_orm_dep | expected_logger_dep        | expected_auth_dep    |
      | gin       | postgres        | gorm | slog   | jwt       | github.com/gin-gonic/gin | github.com/lib/pq              | gorm.io/gorm     |                            | golang.org/x/crypto  |
      | fiber     | mysql           |      | zap    | none      | github.com/gofiber/fiber | github.com/go-sql-driver/mysql|                  | go.uber.org/zap            |                      |
      | echo      | sqlite3         | gorm | logrus | api-key   | github.com/labstack/echo | github.com/mattn/go-sqlite3    | gorm.io/gorm     | github.com/sirupsen/logrus |                      |

  Scenario Outline: docker-compose.yml matches database selection
    Given I generate a project with database "<database>" and driver "<database_driver>"
    When I examine the docker-compose.yml file
    Then it should use the database service "<expected_db_service>"
    And it should not contain services for other databases
    And the database configuration should match the driver selection
    And environment variables should be consistent with the database type

    Examples:
      | database   | database_driver | expected_db_service |
      | postgresql | postgres        | postgres:           |
      | mysql      | mysql           | mysql:              |

  Scenario Outline: Configuration files have consistent database settings
    Given I generate a project with:
      | database_driver | <database_driver> |
      | orm             | <orm>             |
    When I examine all configuration files
    Then database connection strings should use the correct format
    And configuration files should not contain mixed database references
    And ORM-specific configurations should only exist when ORM is selected
    And raw SQL configurations should only exist when no ORM is selected

    Examples:
      | database_driver | orm  |
      | postgres        | gorm |
      | postgres        |      |
      | mysql           | gorm |
      | mysql           |      |
      | sqlite3         | gorm |
      | sqlite3         |      |

  Scenario: Cross-file configuration consistency
    Given I generate a project with complex configuration:
      | framework       | gin        |
      | database_driver | postgres   |
      | orm             | gorm       |
      | logger          | zap        |
      | auth_type       | jwt        |
      | architecture    | clean      |
    When I validate configuration consistency across all files
    Then go.mod should contain all required dependencies
    And docker-compose.yml should match the database selection
    And config files should use consistent naming and values
    And environment variables should be properly aligned
    And no configuration conflicts should exist between files

  Scenario Outline: Invalid configurations are rejected
    Given I attempt to generate a project with invalid configuration:
      | framework | database | orm    | issue                |
      | gin       | mongodb  | gorm   | unsupported database |
      | vue       | postgres | gorm   | unsupported framework|
      | gin       | postgres | django | unsupported orm      |
    When I attempt project generation
    Then the generation should fail with a clear error message
    And no partial configuration files should be created
    And the error should explain the configuration conflict

    Examples:
      | framework | database | orm    | issue                |
      | gin       | mongodb  | gorm   | unsupported database |
      | vue       | postgres | gorm   | unsupported framework|

  Scenario: Configuration consistency under different architectures
    Given I test configuration consistency with various architectures:
      | architecture | framework | database | orm  |
      | standard     | gin       | postgres | gorm |
      | clean        | fiber     | mysql    |      |
      | ddd          | echo      | sqlite3  | gorm |
      | hexagonal    | gin       | postgres |      |
    When I generate projects for each architecture
    Then all projects should have consistent configuration files
    And architecture-specific configurations should be properly set
    And no architecture should affect basic configuration consistency
    And dependency management should remain consistent across architectures