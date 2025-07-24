Feature: Framework Consistency Validation
  As a developer using go-starter
  I want to ensure no framework cross-contamination occurs
  So that my generated project uses only the selected framework without conflicts

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized

  Scenario Outline: Framework isolation prevents cross-contamination
    Given I generate a project with framework "<framework>"
    When I scan all generated Go files for framework references
    Then the project should only contain "<framework>" framework imports
    And the project should not contain "<forbidden_framework_1>" framework imports
    And the project should not contain "<forbidden_framework_2>" framework imports
    And main.go should use the correct framework initialization pattern
    And go.mod should contain only the "<framework>" dependency

    Examples:
      | framework | forbidden_framework_1 | forbidden_framework_2 |
      | gin       | fiber                 | echo                  |
      | fiber     | gin                   | echo                  |
      | echo      | gin                   | fiber                 |

  Scenario Outline: Framework-specific patterns are correctly implemented
    Given I generate a project with:
      | framework | <framework> |
      | logger    | slog        |
      | database  | postgresql  |
    When I examine the framework implementation
    Then the main.go should use "<expected_initialization>"
    And handlers should use "<expected_handler_pattern>"
    And middleware should use "<expected_middleware_pattern>"
    And no conflicting framework patterns should exist

    Examples:
      | framework | expected_initialization | expected_handler_pattern | expected_middleware_pattern |
      | gin       | gin.Default()          | gin.Context               | gin.HandlerFunc             |
      | fiber     | fiber.New()            | fiber.Ctx                 | fiber.Handler               |
      | echo      | echo.New()             | echo.Context              | echo.MiddlewareFunc         |

  Scenario: Multiple projects maintain framework isolation
    Given I generate multiple projects with different frameworks:
      | project_name | framework |
      | gin_project  | gin       |
      | fiber_project| fiber     |
      | echo_project | echo      |
    When I validate framework consistency across all projects
    Then each project should use only its designated framework
    And no project should contain references to other frameworks
    And all framework-specific patterns should be correctly implemented
    And go.mod files should contain only the correct framework dependencies