Feature: Code Quality Validation
  As a developer using go-starter
  I want generated projects to have high code quality
  So that I can build upon clean, maintainable code without quality issues

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized

  Scenario Outline: Generated projects compile successfully
    Given I generate a project with configuration:
      | framework      | <framework>      |
      | database       | <database>       |
      | database_driver| <database_driver>|
      | orm            | <orm>            |
      | logger         | <logger>         |
      | auth_type      | <auth_type>      |
      | architecture   | <architecture>   |
    When I attempt to compile the generated project
    Then the compilation should succeed without errors
    And the build output should not contain warnings

    Examples:
      | framework | database     | database_driver | orm   | logger | auth_type | architecture |
      | gin       | postgresql   | postgres        |       | slog   | jwt       | standard     |
      | fiber     | mysql        | mysql           | gorm  | zap    | none      | standard     |
      | echo      | sqlite       | sqlite3         |       | logrus | api-key   | clean        |

  Scenario Outline: Generated projects have no unused imports
    Given I generate a project with configuration:
      | framework      | <framework>      |
      | database       | <database>       |
      | database_driver| <database_driver>|
      | orm            | <orm>            |
      | logger         | <logger>         |
      | auth_type      | <auth_type>      |
    When I analyze the generated Go files for unused imports
    Then there should be no unused import statements
    And goimports should report no formatting issues
    And common problematic imports like "fmt", "os", "models" should only be present when used

    Examples:
      | framework | database     | database_driver | orm   | logger | auth_type |
      | gin       | postgresql   | postgres        | gorm  | slog   | jwt       |
      | fiber     | mysql        | mysql           |       | zap    | none      |
      | echo      | sqlite       | sqlite3         | gorm  | logrus | api-key   |

  Scenario Outline: Generated projects have no unused variables
    Given I generate a project with configuration:
      | framework      | <framework>      |
      | database       | <database>       |
      | database_driver| <database_driver>|
      | orm            | <orm>            |
      | logger         | <logger>         |
    When I run static analysis on the generated project
    Then go vet should report no unused variables
    And there should be no "declared but not used" errors
    And there should be no "assigned but not used" errors

    Examples:
      | framework | database     | database_driver | orm   | logger |
      | gin       | postgresql   | postgres        |       | slog   |
      | fiber     | mysql        | mysql           | gorm  | zap    |
      | echo      | sqlite       | sqlite3         |       | logrus |

  Scenario Outline: Configuration files are consistent with selections
    Given I generate a project with configuration:
      | framework      | <framework>      |
      | database       | <database>       |
      | database_driver| <database_driver>|
      | orm            | <orm>            |
      | logger         | <logger>         |
      | auth_type      | <auth_type>      |
    When I examine the generated configuration files
    Then go.mod should contain the correct framework dependency "<expected_framework_dep>"
    And go.mod should contain the correct database driver dependency "<expected_db_dep>"
    And go.mod should contain the correct logger dependency "<expected_logger_dep>"
    And docker-compose.yml should use the correct database service "<expected_db_service>"
    And configuration files should not contain contradictory settings

    Examples:
      | framework | database   | database_driver | orm  | logger | auth_type | expected_framework_dep        | expected_db_dep                      | expected_logger_dep           | expected_db_service |
      | gin       | postgresql | postgres        | gorm | slog   | jwt       | github.com/gin-gonic/gin      | github.com/lib/pq                    |                               | postgres:           |
      | fiber     | mysql      | mysql           |      | zap    | none      | github.com/gofiber/fiber      | github.com/go-sql-driver/mysql      | go.uber.org/zap               | mysql:              |
      | echo      | sqlite     | sqlite3         | gorm | logrus | api-key   | github.com/labstack/echo      | github.com/mattn/go-sqlite3         | github.com/sirupens/logrus    |                     |

  Scenario Outline: No framework cross-contamination occurs
    Given I generate a project with framework "<framework>"
    When I scan all generated Go files for framework references
    Then the project should only contain "<framework>" framework imports
    And the project should not contain "<forbidden_framework_1>" framework imports  
    And the project should not contain "<forbidden_framework_2>" framework imports
    And main.go should use the correct framework initialization

    Examples:
      | framework | forbidden_framework_1 | forbidden_framework_2 |
      | gin       | fiber                 | echo                  |
      | fiber     | gin                   | echo                  |
      | echo      | gin                   | fiber                 |

  Scenario: Multiple projects maintain quality consistency
    Given I generate multiple projects with different configurations:
      | name          | framework | database   | orm  | logger |
      | gin-project   | gin       | postgresql | gorm | slog   |
      | fiber-project | fiber     | mysql      |      | zap    |
      | echo-project  | echo      | sqlite     | gorm | logrus |
    When I validate the code quality of all projects
    Then each project should compile successfully
    And each project should have no unused imports
    And each project should have no unused variables
    And each project should have consistent configuration files
    And no project should have framework cross-contamination