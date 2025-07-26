Feature: Code Optimization Pipeline
  As a Go developer using go-starter
  I want to optimize my generated code automatically
  So that I can improve code quality and reduce maintenance overhead

  Background:
    Given I am using go-starter CLI
    And the optimization system is available

  @optimization @safe
  Scenario: Generate project with safe optimization level
    Given I want to create a new "web-api" project
    And I set the optimization level to "safe"
    When I generate the project "safe-optimized-api"
    Then the project should be created successfully
    And unused imports should be removed
    And imports should be organized alphabetically
    And no variables or functions should be removed
    And the project should compile without errors

  @optimization @standard
  Scenario: Generate project with standard optimization level
    Given I want to create a new "cli" project
    And I set the optimization level to "standard"
    When I generate the project "standard-optimized-cli"
    Then the project should be created successfully
    And unused imports should be removed
    And imports should be organized alphabetically
    And missing imports should be added carefully
    And no variables or functions should be removed
    And the project should compile without errors

  @optimization @aggressive
  Scenario: Generate project with aggressive optimization level
    Given I want to create a new "web-api" project
    And I set the optimization level to "aggressive"
    When I generate the project "aggressive-optimized-api"
    Then the project should be created successfully
    And unused imports should be removed
    And imports should be organized alphabetically
    And missing imports should be added
    And unused local variables should be removed
    And unused private functions should be removed
    And the project should compile without errors

  @optimization @expert
  Scenario: Generate project with expert optimization level
    Given I want to create a new "cli" project
    And I set the optimization level to "expert"
    When I generate the project "expert-optimized-cli"
    Then the project should be created successfully
    And all optimizations should be applied
    And high concurrency settings should be used
    And the project should compile without errors

  @optimization @dry-run
  Scenario: Preview optimizations with dry-run mode
    Given I want to create a new "web-api" project
    And I set the optimization level to "aggressive"
    And I enable dry-run mode
    When I generate the project "dry-run-test"
    Then optimization changes should be previewed
    And no files should be modified
    And warnings should be displayed for risky optimizations
    And the preview should show potential improvements

  @optimization @backup
  Scenario: Generate project with backup creation
    Given I want to create a new "cli" project
    And I set the optimization level to "aggressive"
    And I enable backup creation
    When I generate the project "backup-test-cli"
    Then the project should be created successfully
    And backup files should be created for modified files
    And optimization should be applied safely
    And the project should compile without errors

  @optimization @profile
  Scenario: Generate project using optimization profile
    Given I want to create a new "web-api" project
    And I set the optimization profile to "performance"
    When I generate the project "profile-optimized-api"
    Then the project should be created successfully
    And performance-focused optimizations should be applied
    And the project should compile without errors

  @optimization @integration
  Scenario: Optimize existing generated project
    Given I have an existing "web-api" project "existing-api"
    When I apply "standard" optimization to the existing project
    Then the project should be optimized in place
    And backup files should be created
    And the project should still compile without errors
    And optimization metrics should be reported

  @optimization @configuration
  Scenario: Validate optimization configuration
    Given I want to create a new "cli" project
    And I set invalid optimization settings
    When I attempt to generate the project "invalid-config"
    Then configuration validation should fail
    And helpful error messages should be displayed
    And suggestions for valid configurations should be provided

  @optimization @performance
  Scenario: Monitor optimization performance
    Given I want to create a new "web-api" project with many files
    And I set the optimization level to "expert"
    When I generate the project "performance-test-api"
    Then the optimization should complete within reasonable time
    And performance metrics should be reported
    And resource usage should be monitored
    And the project should compile without errors