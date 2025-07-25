Feature: Edge Cases and Error Handling Combination Testing
  As a go-starter maintainer
  I want to ensure that edge cases and invalid combinations are handled gracefully
  So that users get helpful error messages and the CLI doesn't crash

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And I am testing edge cases and error scenarios

  @edge-case @error-handling
  Scenario Outline: Invalid combination error handling
    When I attempt to generate a project with invalid configuration:
      | type      | <type>      |
      | framework | <framework> |
      | database  | <database>  |
    Then generation should fail gracefully
    And error message should be informative and helpful
    And error message should suggest valid alternatives
    And CLI should exit with appropriate error code

    Examples:
      | type    | framework | database | expected_error                                    |
      | cli     | gin       | any      | CLI projects cannot use web frameworks            |
      | library | echo      | any      | Library projects cannot use web frameworks        |
      | lambda  | cobra     | any      | Lambda functions cannot use CLI frameworks        |
      | cli     | none      | postgres | CLI projects typically don't need databases       |

  @edge-case @minimal-config
  Scenario Outline: Minimal configuration combinations
    When I generate a project with minimal configuration:
      | type | <type> |
    Then the project should be generated with sensible defaults
    And the project should compile successfully
    And default configuration should be appropriate for <type>
    And generated project should follow best practices

    Examples:
      | type         | expected_framework | expected_logger |
      | web-api      | gin               | slog            |
      | cli          | cobra             | slog            |
      | library      | none              | slog            |
      | microservice | gin               | slog            |

  @edge-case @maximum-complexity
  Scenario Outline: Maximum complexity combinations
    When I generate a project with all advanced features enabled:
      | type           | <type>           |
      | architecture   | <architecture>   |
      | framework      | <framework>      |
      | database       | postgres         |
      | orm            | gorm             |
      | auth           | oauth2           |
      | logger         | zap              |
      | monitoring     | enabled          |
      | metrics        | prometheus       |
      | tracing        | jaeger           |
      | cache          | redis            |
      | queue          | rabbitmq         |
      | search         | elasticsearch    |
      | asset_pipeline | webpack          |
      | testing        | comprehensive    |
      | documentation  | enabled          |
      | ci_cd          | github-actions   |
      | deployment     | kubernetes       |
    Then the project should handle maximum complexity gracefully
    And compilation should succeed despite complexity
    And all features should integrate correctly
    And performance should remain acceptable

    Examples:
      | type    | architecture |
      | web-api | hexagonal    |
      | web-api | ddd          |

  @edge-case @progressive-disclosure
  Scenario Outline: Progressive disclosure edge cases
    When I use progressive disclosure with configuration:
      | mode       | <mode>       |
      | complexity | <complexity> |
      | type       | <type>       |
    And I transition from basic to advanced mode
    Then configuration should be preserved correctly
    And additional options should become available
    And no configuration conflicts should occur
    And user experience should be smooth

    Examples:
      | mode     | complexity | type    |
      | basic    | simple     | cli     |
      | advanced | expert     | web-api |

  @edge-case @concurrent-generation
  Scenario: Concurrent project generation
    When I generate multiple projects simultaneously:
      | name     | type    | framework |
      | project1 | web-api | gin       |
      | project2 | cli     | cobra     |
      | project3 | library | none      |
    Then all projects should be generated successfully
    And no file conflicts should occur
    And temporary directories should be properly isolated
    And resource usage should be reasonable

  @edge-case @special-characters
  Scenario Outline: Special characters in project names and paths
    When I generate a project with special characters:
      | name       | <project_name> |
      | output_dir | <output_dir>   |
      | type       | <type>         |
    Then the project should handle special characters gracefully
    And file paths should be properly escaped
    And generated code should compile correctly
    And module paths should be valid Go syntax

    Examples:
      | project_name    | output_dir        | type    |
      | my-awesome-api  | ./test-output     | web-api |
      | my_cli_tool     | ./test_output     | cli     |
      | project.v2      | ./test.output     | library |
      | café-api        | ./café-output     | web-api |

  @edge-case @large-workspace
  Scenario: Large workspace with many components
    When I generate a workspace with many components:
      | component_count | 10      |
      | types          | web-api,cli,library,microservice |
    Then workspace generation should complete successfully
    And memory usage should remain reasonable
    And all components should be properly configured
    And go.mod dependency resolution should work correctly
    And build times should be acceptable

  @edge-case @dependency-conflicts
  Scenario Outline: Dependency conflict resolution
    When I generate a project with potentially conflicting dependencies:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | <auth>      |
    Then dependency conflicts should be resolved automatically
    And go.mod should contain compatible versions
    And the project should compile without version conflicts
    And security vulnerabilities should be avoided

    Examples:
      | framework | database | orm  | auth   |
      | gin       | postgres | gorm | jwt    |
      | echo      | mysql    | sqlx | oauth2 |

  @edge-case @template-inheritance
  Scenario Outline: Complex template inheritance scenarios
    When I generate a project that uses complex template inheritance:
      | type         | <type>         |
      | architecture | <architecture> |
      | shared_templates | enabled    |
    Then template inheritance should work correctly
    And no circular dependencies should exist in templates
    And shared templates should be applied consistently
    And custom templates should override base templates properly

    Examples:
      | type         | architecture |
      | web-api      | clean        |
      | microservice | hexagonal    |

  @edge-case @cross-platform-paths
  Scenario Outline: Cross-platform path handling edge cases
    When I generate a project on platform "<platform>":
      | type       | web-api      |
      | output_dir | <output_dir> |
    Then paths should be handled correctly for <platform>
    And file separators should be platform-appropriate
    And generated scripts should work on <platform>
    And file permissions should be set correctly

    Examples:
      | platform | output_dir                    |
      | windows  | C:\temp\test-project         |
      | unix     | /tmp/test-project            |
      | windows  | \\server\share\test-project  |

  @edge-case @rollback-scenarios
  Scenario Outline: Generation failure and rollback
    When project generation fails during <failure_stage>:
      | type         | web-api         |
      | architecture | clean           |
      | failure_point| <failure_stage> |
    Then the CLI should handle the failure gracefully
    And partial files should be cleaned up
    And error message should indicate the failure point
    And user should be able to retry generation
    And no corrupted state should remain

    Examples:
      | failure_stage    |
      | template_parsing |
      | file_generation  |
      | dependency_resolution |
      | post_generation_hooks |

  @edge-case @memory-constraints
  Scenario: Generation under memory constraints
    When I generate projects under memory pressure:
      | available_memory | low       |
      | project_count   | multiple  |
      | project_size    | large     |
    Then generation should complete successfully
    And memory usage should be optimized
    And garbage collection should be effective
    And no memory leaks should occur

  @edge-case @concurrent-template-access
  Scenario: Concurrent template access patterns
    When multiple generation processes access templates simultaneously
    Then template loading should be thread-safe
    And no race conditions should occur
    And template caching should work correctly
    And performance should not degrade significantly

  @edge-case @invalid-go-version
  Scenario Outline: Invalid Go version handling
    When I specify an invalid Go version:
      | go_version | <go_version> |
      | type       | web-api      |
    Then the CLI should validate the Go version
    And provide helpful error message about supported versions
    And suggest the nearest valid version
    And not generate incompatible code

    Examples:
      | go_version |
      | 1.17       |
      | 1.25       |
      | 2.0        |
      | invalid    |

  @edge-case @network-dependency-failure
  Scenario: Network dependency resolution failure
    When network access is limited during generation:
      | network_access | restricted |
      | type          | web-api    |
      | dependencies  | external   |
    Then generation should handle network failures gracefully
    And provide informative error messages
    And suggest offline alternatives when possible
    And not leave projects in broken state