Feature: Template Logic Validation
  As a developer using go-starter
  I want template conditional generation to work correctly
  So that only the files I need are generated based on my configuration choices

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized

  Scenario Outline: ORM conditional file generation
    Given I specify ORM as "<orm>"
    And I specify database driver as "<database_driver>"
    When I generate a web-api project
    Then the ORM-specific files should be generated correctly:
      | orm_type | should_exist                    | should_not_exist               |
      | gorm     | internal/database/gorm.go       | internal/database/raw.go       |
      | gorm     | internal/models/user.go         | internal/database/sqlx.go      |
      |          | internal/database/raw.go        | internal/database/gorm.go      |
      |          | internal/database/migrations/   | internal/models/user.go        |
    And the go.mod dependencies should match the ORM selection

    Examples:
      | orm  | database_driver |
      | gorm | postgres        |
      |      | postgres        |
      | gorm | mysql           |
      |      | mysql           |

  Scenario Outline: Authentication conditional file generation
    Given I specify authentication type as "<auth_type>"
    When I generate a web-api project
    Then the authentication files should be generated correctly:
      | auth_type | should_exist                   | should_not_exist              |
      | jwt       | internal/middleware/auth.go    | internal/auth/oauth.go        |
      | jwt       | internal/auth/jwt.go           | internal/auth/session.go      |
      | api-key   | internal/middleware/auth.go    | internal/auth/jwt.go          |
      | api-key   | internal/auth/apikey.go        | internal/auth/oauth.go        |
      | none      |                                | internal/middleware/auth.go   |
      | none      |                                | internal/auth/jwt.go          |
    And the go.mod dependencies should match the authentication selection

    Examples:
      | auth_type |
      | jwt       |
      | api-key   |
      | none      |

  Scenario Outline: Logger conditional file generation
    Given I specify logger as "<logger>"
    When I generate a web-api project
    Then the logger-specific file "internal/logger/<logger>.go" should exist
    And other logger files should not exist:
      | logger  | should_not_exist               |
      | slog    | internal/logger/zap.go         |
      | slog    | internal/logger/logrus.go      |
      | slog    | internal/logger/zerolog.go     |
      | zap     | internal/logger/slog.go        |
      | zap     | internal/logger/logrus.go      |
      | zap     | internal/logger/zerolog.go     |
      | logrus  | internal/logger/slog.go        |
      | logrus  | internal/logger/zap.go         |
      | logrus  | internal/logger/zerolog.go     |
      | zerolog | internal/logger/slog.go        |
      | zerolog | internal/logger/zap.go         |
      | zerolog | internal/logger/logrus.go      |
    And the go.mod should contain the correct logger dependency
    And the logger implementation should use the correct import statements

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  Scenario Outline: Import consistency with conditional generation
    Given I generate a project with:
      | framework      | <framework>      |
      | database_driver| <database_driver>|  
      | orm            | <orm>            |
      | logger         | <logger>         |
    When I analyze the import statements in generated files
    Then imports should be consistent with the configuration:
      | condition                    | file_pattern           | should_contain_import        | should_not_contain_import   |
      | orm is gorm                  | internal/handlers/     | models package               | database/sql                |
      | orm is empty                 | internal/handlers/     | database/sql                 | models package              |
      | framework is gin             | cmd/server/main.go     | github.com/gin-gonic/gin     | github.com/gofiber/fiber    |
      | framework is fiber           | cmd/server/main.go     | github.com/gofiber/fiber     | github.com/gin-gonic/gin    |
      | logger is zap                | internal/logger/zap.go | go.uber.org/zap              | log/slog                    |
      | logger is slog               | internal/logger/slog.go| log/slog                     | go.uber.org/zap             |

    Examples:
      | framework | database_driver | orm  | logger |
      | gin       | postgres        | gorm | slog   |
      | gin       | postgres        |      | slog   |
      | fiber     | mysql           | gorm | zap    |
      | echo      | sqlite3         |      | logrus |

  Scenario: Dependency consistency across all features
    Given I generate projects with various feature combinations:
      | name      | framework | database_driver | orm  | logger  | auth_type |
      | full      | gin       | postgres        | gorm | zap     | jwt       |
      | minimal   | echo      |                 |      | slog    | none      |
      | mixed     | fiber     | mysql           |      | logrus  | api-key   |
    When I examine the go.mod file for each project
    Then each project should have only the required dependencies:
      | project | should_have                  | should_not_have              |
      | full    | github.com/gin-gonic/gin     | github.com/gofiber/fiber     |
      | full    | github.com/lib/pq            | github.com/go-sql-driver/mysql |
      | full    | gorm.io/gorm                 | github.com/jmoiron/sqlx      |
      | full    | go.uber.org/zap              | github.com/sirupsen/logrus   |
      | full    | golang.org/x/crypto          | golang.org/x/oauth2          |
      | minimal | github.com/labstack/echo     | github.com/gin-gonic/gin     |
      | minimal |                              | github.com/lib/pq            |
      | minimal |                              | gorm.io/gorm                 |
      | minimal |                              | go.uber.org/zap              |
      | mixed   | github.com/gofiber/fiber     | github.com/gin-gonic/gin     |
      | mixed   | github.com/go-sql-driver/mysql | github.com/lib/pq         |
      | mixed   | github.com/sirupsen/logrus   | go.uber.org/zap              |
    And no project should have conflicting or unnecessary dependencies

  Scenario: Template variable consistency
    Given I generate multiple projects with the same configuration
    When I compare the generated template variables across projects
    Then all projects should use consistent variable naming
    And all projects should have properly resolved template placeholders
    And no projects should contain unresolved template variables like "{{.Variable}}"