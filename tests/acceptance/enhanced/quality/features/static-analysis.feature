Feature: Static Analysis Validation
  As a developer using go-starter
  I want generated code to pass static analysis checks
  So that my project starts with clean, high-quality code

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized

  Scenario Outline: Generated code has no unused imports
    Given I generate a project with configuration:
      | framework       | <framework>       |
      | database_driver | <database_driver> |
      | orm             | <orm>             |
      | logger          | <logger>          |
      | auth_type       | <auth_type>       |
    When I run goimports analysis on all Go files
    Then there should be no unused import statements
    And goimports should report no formatting differences
    And all import statements should be properly organized
    And no import cycles should exist

    Examples:
      | framework | database_driver | orm  | logger | auth_type |
      | gin       | postgres        | gorm | slog   | jwt       |
      | fiber     | mysql           |      | zap    | none      |
      | echo      | sqlite3         | gorm | logrus | api-key   |

  Scenario Outline: Generated code has no unused variables
    Given I generate a project with configuration:
      | framework       | <framework>       |
      | database_driver | <database_driver> |
      | orm             | <orm>             |
      | logger          | <logger>          |
    When I run go vet analysis on the project
    Then there should be no "declared but not used" errors
    And there should be no "assigned but not used" errors
    And all variables should have meaningful usage
    And no dead code should be present

    Examples:
      | framework | database_driver | orm  | logger |
      | gin       | postgres        |      | slog   |
      | fiber     | mysql           | gorm | zap    |
      | echo      | sqlite3         |      | logrus |

  Scenario Outline: Problematic imports are only present when actually used
    Given I generate a project with configuration:
      | framework | <framework> |
      | orm       | <orm>       |
      | database  | <database>  |
    When I scan for problematic import patterns
    Then "fmt" package should only be imported when format functions are used
    And "os" package should only be imported when OS functions are used
    And "models" package should only be imported when ORM is "<orm>"
    And database-specific packages should only be imported when database is "<database>"
    And no phantom imports should exist

    Examples:
      | framework | orm  | database   |
      | gin       | gorm | postgresql |
      | gin       |      | postgresql |
      | fiber     | gorm | mysql      |
      | echo      |      | sqlite     |

  Scenario: Comprehensive static analysis passes
    Given I generate projects with various configurations:
      | name    | framework | database | orm  | logger | auth_type |
      | full    | gin       | postgres | gorm | zap    | jwt       |
      | minimal | gin       |          |      | slog   | none      |
      | mixed   | fiber     | mysql    |      | logrus | api-key   |
    When I run comprehensive static analysis on all projects
    Then all projects should pass go vet without errors
    And all projects should pass goimports without changes
    And all projects should pass golangci-lint basic checks
    And no code quality issues should be detected

  Scenario Outline: Code complexity remains manageable
    Given I generate a project with architecture "<architecture>"
    When I analyze code complexity metrics
    Then cyclomatic complexity should be within acceptable limits
    And function length should be reasonable
    And package coupling should be appropriate for the architecture
    And code duplication should be minimal

    Examples:
      | architecture |
      | standard     |
      | clean        |
      | ddd          |
      | hexagonal    |

  Scenario: Generated code follows Go best practices
    Given I generate a project with standard configuration
    When I analyze the code for Go best practices
    Then all exported functions should have documentation comments
    And error handling should follow Go conventions
    And naming conventions should be idiomatic
    And package organization should be logical
    And interface usage should be appropriate

  Scenario: Performance-related static analysis
    Given I generate projects with different logger configurations:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |
    When I analyze for performance anti-patterns
    Then there should be no obvious performance bottlenecks
    And logging should use appropriate levels
    And resource management should be proper
    And no memory leaks should be detectable through static analysis