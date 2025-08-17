Feature: AI Advisor Edge Cases and Error Handling
  As a developer
  I want the AI advisor to handle edge cases gracefully
  So that I get helpful guidance even with unusual or incomplete requirements

  Background:
    Given the AI advisor is available
    And the error handling system is active

  Scenario: Empty or minimal requirements
    Given I provide minimal requirements
    And my project domain is "other"
    And my team experience is "mixed"
    When I request recommendations
    Then I should still get a default recommendation
    And the confidence should be low but positive
    And the reasoning should explain the default choice
    And the recommendation should suggest gathering more requirements

  Scenario: Conflicting requirements - junior team with complex needs
    Given I need a "api" project
    And my project domain is "fintech"
    And my team has "junior" experience level
    But I need "enterprise" compliance
    And I expect "massive" load
    And I have "distributed" database requirements
    When I request recommendations
    Then I should get a compromise recommendation
    And the reasoning should acknowledge the conflicts
    And the recommendation should suggest team augmentation
    And the alternatives should include both simple and complex options
    And the confidence should reflect the uncertainty

  Scenario: Invalid project type
    Given I specify an invalid project type "invalid-type"
    When I request recommendations through CLI
    Then I should get a clear error message
    And the error should list valid project types
    And the exit code should indicate failure
    And no partial recommendations should be generated

  Scenario: Invalid team experience level
    Given I need a "api" project
    And I specify invalid team experience "invalid-experience"
    When I request recommendations through CLI
    Then I should get a clear error message
    And the error should list valid experience levels
    And the exit code should indicate failure

  Scenario: Invalid domain
    Given I need a "api" project
    And I specify a domain with special characters "e-commerce@#$"
    When I request recommendations
    Then I should get sanitized domain handling
    And the recommendation should still be generated
    And the reasoning should mention domain normalization

  Scenario: Missing required CLI flags
    Given I run the advisor command
    But I don't provide required flags like "--type"
    When I execute the CLI command
    Then I should get helpful usage information
    And the error should specify which flags are missing
    And example commands should be provided
    And the exit code should indicate user error

  Scenario: Network or service unavailable
    Given the AI advisor service is temporarily unavailable
    When I request recommendations
    Then I should get a fallback recommendation system
    And the error message should be user-friendly
    And offline recommendations should be provided if possible
    And retry suggestions should be given

  Scenario: Malformed JSON output request
    Given I request recommendations in JSON format
    But the internal data is corrupted
    When I parse the output
    Then I should get a graceful JSON error response
    And the error field should describe the issue
    And the response should still be valid JSON
    And debugging information should be available

  Scenario: Extremely large file count estimates
    Given I need a complex "microservice" project
    And my requirements lead to unrealistic file counts
    When I request recommendations
    Then the file count estimate should be capped at reasonable limits
    And the reasoning should mention the complexity cap
    And alternatives with lower complexity should be suggested
    And a warning about project scope should be included

  Scenario: Zero confidence recommendations
    Given I provide completely contradictory requirements
    And no blueprint can reasonably satisfy them
    When I request recommendations
    Then I should still get a recommendation with zero confidence
    And the reasoning should explain why confidence is low
    And multiple alternatives should be provided
    And guidance on clarifying requirements should be given

  Scenario: Timeout handling for complex analysis
    Given I have extremely complex requirements
    And the analysis takes longer than expected
    When I request recommendations
    Then the system should timeout gracefully
    And I should get a partial recommendation if possible
    And the error should suggest simplifying requirements
    And retry options should be provided

  Scenario: Concurrent request handling
    Given multiple users request recommendations simultaneously
    When the system is under load
    Then each request should be handled independently
    And recommendations should remain consistent
    And no request should fail due to resource contention
    And response times should remain reasonable

  Scenario: Memory constraints with large requirements
    Given I provide extremely detailed requirements
    And the requirement object is very large
    When I request recommendations
    Then the system should handle memory efficiently
    And recommendations should still be generated
    And no memory leaks should occur
    And performance should remain acceptable

  Scenario: Invalid CLI output format
    Given I request recommendations
    And I specify an invalid output format "xml"
    When I execute the CLI command
    Then I should get an error about invalid format
    And valid formats should be listed
    And the default format should be suggested
    And the exit code should indicate user error

  Scenario: Blueprint registry corruption
    Given the blueprint registry is corrupted or missing
    When I request recommendations
    Then I should get an error about unavailable blueprints
    And recovery suggestions should be provided
    And basic fallback recommendations should be offered if possible
    And system health checks should be suggested

  Scenario: Circular dependency in recommendations
    Given the recommendation logic encounters circular dependencies
    When processing complex requirements
    Then the system should detect and break cycles
    And a clear recommendation should still be provided
    And the reasoning should explain the resolution
    And no infinite loops should occur

  Scenario: Unicode and special characters in project names
    Given I specify a project name with unicode characters "项目名称"
    And special characters in description "Project with émojis 🚀"
    When I request recommendations
    Then the system should handle unicode gracefully
    And generated project names should be safe for file systems
    And recommendations should not be affected by character encoding
    And sanitization should be transparent to the user

  Scenario: Version compatibility issues
    Given I specify an unsupported Go version "1.15"
    When I request recommendations
    Then I should get a warning about version support
    And alternative supported versions should be suggested
    And recommendations should adapt to version constraints
    And migration guidance should be provided

  Scenario: Resource exhaustion recovery
    Given the system is running low on resources
    When I request recommendations
    Then the system should prioritize essential functionality
    And reduced-complexity recommendations should be provided
    And resource usage should be optimized
    And graceful degradation should occur

  Scenario: Data validation edge cases
    Given I provide extremely long strings in requirements
    And numeric values outside expected ranges
    And deeply nested requirement structures
    When I request recommendations
    Then input validation should handle all edge cases
    And appropriate error messages should be provided
    And the system should remain stable
    And no security vulnerabilities should be exploited