Feature: Configuration Matrix Testing
  As a developer using go-starter
  I want all valid configuration combinations to work correctly
  So that I can confidently choose any supported combination of features

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized

  Scenario Outline: Critical configuration combinations work flawlessly
    Given I use the critical configuration combination:
      | framework      | <framework>      |
      | database       | <database>       |
      | database_driver| <database_driver>|
      | orm            | <orm>            |
      | logger         | <logger>         |
      | auth_type      | <auth_type>      |
      | architecture   | <architecture>   |
    When I generate a web-api project with this configuration
    Then the project should generate successfully
    And the project should compile without errors
    And all framework-specific code should be consistent
    And all database configuration should be consistent  
    And all logger implementation should be consistent
    And all authentication setup should be consistent
    And the architecture structure should be correct

    Examples: Critical Combinations (Must Always Work)
      | framework | database   | database_driver | orm  | logger   | auth_type | architecture |
      | gin       | postgresql | postgres        | gorm | slog     | jwt       | standard     |
      | gin       | postgresql | postgres        |      | slog     | jwt       | standard     |
      | echo      | mysql      | mysql           | gorm | zap      | jwt       | standard     |
      | fiber     | postgresql | postgres        |      | zerolog  | none      | standard     |

  Scenario Outline: High priority combinations work correctly
    Given I use the high priority configuration combination:
      | framework      | <framework>      |
      | database       | <database>       |
      | database_driver| <database_driver>|
      | orm            | <orm>            |
      | logger         | <logger>         |
      | auth_type      | <auth_type>      |
      | architecture   | <architecture>   |
    When I generate a web-api project with this configuration
    Then the project should generate successfully
    And the project should compile without errors
    And the advanced architecture should be implemented correctly
    And all component integrations should work properly

    Examples: High Priority Combinations (Common Advanced Usage)
      | framework | database   | database_driver | orm  | logger | auth_type | architecture |
      | gin       | postgresql | postgres        | gorm | zap    | jwt       | clean        |
      | echo      | postgresql | postgres        | gorm | slog   | jwt       | clean        |
      | fiber     | mysql      | mysql           | gorm | logrus | api-key   | standard     |
      | gin       | sqlite     | sqlite3         | gorm | slog   | none      | standard     |

  Scenario Outline: Framework consistency across all database combinations  
    Given I use framework "<framework>"
    And I test it with all supported database combinations:
      | database   | database_driver | orm  |
      | postgresql | postgres        | gorm |
      | postgresql | postgres        |      |
      | mysql      | mysql           | gorm |
      | mysql      | mysql           |      |
      | sqlite     | sqlite3         | gorm |
      | sqlite     | sqlite3         |      |
    When I generate projects for each database combination
    Then all projects should use "<framework>" consistently
    And no project should contain references to other frameworks
    And the "<framework>" initialization should be correct in all projects
    And framework-specific middleware should be properly configured

    Examples:
      | framework |
      | gin       |
      | fiber     |
      | echo      |

  Scenario Outline: Database consistency across all framework combinations
    Given I use database "<database>" with driver "<database_driver>"
    And I test it with all supported framework combinations:
      | framework | logger | auth_type |
      | gin       | slog   | jwt       |
      | fiber     | zap    | none      |
      | echo      | logrus | api-key   |
    When I generate projects for each framework combination
    Then all projects should use database "<database>" consistently
    And the database driver "<database_driver>" should be configured correctly
    And database connection strings should match the database type
    And docker-compose.yml should use the correct database service

    Examples:
      | database   | database_driver |
      | postgresql | postgres        |
      | mysql      | mysql           |
      | sqlite     | sqlite3         |

  Scenario Outline: Logger consistency across all combinations
    Given I use logger "<logger>"
    And I test it with various configuration combinations:
      | framework | database | orm  | auth_type | architecture |
      | gin       | postgres | gorm | jwt       | standard     |
      | fiber     | mysql    |      | none      | clean        |
      | echo      | sqlite   | gorm | api-key   | ddd          |
    When I generate projects for each combination
    Then all projects should use logger "<logger>" consistently
    And the logger implementation file should exist
    And other logger files should not exist
    And go.mod should contain the correct logger dependency
    And logger initialization should be correct

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  Scenario Outline: Architecture consistency validation
    Given I use architecture "<architecture>"
    And I test it with various feature combinations:
      | framework | database | orm  | logger | auth_type |
      | gin       | postgres | gorm | slog   | jwt       |
      | fiber     | mysql    |      | zap    | none      |
      | echo      | sqlite   | gorm | logrus | api-key   |
    When I generate projects for each combination
    Then all projects should implement "<architecture>" architecture correctly
    And the directory structure should match the architecture pattern
    And component relationships should follow architecture principles
    And dependency flow should be correct for the architecture

    Examples:
      | architecture |
      | standard     |
      | clean        |
      | ddd          |
      | hexagonal    |

  Scenario: Full matrix validation (Performance Test)
    Given I have the complete configuration matrix:
      | dimension    | values                           |
      | frameworks   | gin, fiber, echo                 |
      | databases    | postgresql, mysql, sqlite        |
      | orms         | "", gorm, sqlx                   |
      | loggers      | slog, zap, logrus, zerolog       |
      | auth_types   | none, jwt, api-key               |
      | architectures| standard, clean, ddd, hexagonal  |
    When I run the matrix test in short mode
    Then only critical priority combinations should be tested
    And all critical combinations should pass
    And the test execution should complete within acceptable time limits
    And no memory leaks should occur during testing

  Scenario: Configuration validation prevents invalid combinations
    Given I attempt to use invalid configuration combinations:
      | framework | database | orm    | issue                    |
      | gin       | mongodb  | gorm   | unsupported database     |
      | vue       | postgres | gorm   | unsupported framework    |
      | gin       | postgres | django | unsupported orm          |
    When I attempt to generate projects with these configurations
    Then the generation should fail gracefully
    And clear error messages should explain the invalid combinations
    And no partial projects should be created

  Scenario: Matrix test performance characteristics
    Given I run the configuration matrix tests
    When I measure the execution time and resource usage
    Then critical tests should complete within 30 seconds
    And high priority tests should complete within 2 minutes
    And memory usage should remain under 500MB
    And all tests should clean up temporary resources
    And no zombie processes should remain after test completion