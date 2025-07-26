Feature: Multi-Level Optimization System
  As a Go developer using go-starter
  I want to choose appropriate optimization levels for different scenarios
  So that I can balance code quality improvements with safety considerations

  Background:
    Given I am using go-starter CLI
    And the multi-level optimization system is available

  @levels @none
  Scenario: Generate project with no optimizations
    Given I want to create a new "web-api" project
    And I set the optimization level to "none"
    When I generate the project "no-optimization-api"
    Then the project should be created successfully
    And no optimizations should be applied
    And original code structure should be preserved
    And the project should compile without errors

  @levels @safe
  Scenario: Generate project with safe optimizations only
    Given I want to create a new "cli" project
    And I set the optimization level to "safe"
    When I generate the project "safe-optimization-cli"
    Then the project should be created successfully
    And unused imports should be removed
    And imports should be organized
    And no code structure changes should be made
    And the project should compile without errors

  @levels @standard
  Scenario: Generate project with standard optimizations
    Given I want to create a new "web-api" project
    And I set the optimization level to "standard"
    When I generate the project "standard-optimization-api"
    Then the project should be created successfully
    And unused imports should be removed
    And imports should be organized
    And missing imports should be added carefully
    And variable and function structure should be preserved
    And the project should compile without errors

  @levels @aggressive
  Scenario: Generate project with aggressive optimizations
    Given I want to create a new "cli" project
    And I set the optimization level to "aggressive"
    When I generate the project "aggressive-optimization-cli"
    Then the project should be created successfully
    And all import optimizations should be applied
    And unused local variables should be removed
    And unused private functions should be removed
    And conditional optimizations should be applied
    And higher concurrency should be used
    And the project should compile without errors

  @levels @expert
  Scenario: Generate project with expert optimizations
    Given I want to create a new "web-api" project
    And I set the optimization level to "expert"
    When I generate the project "expert-optimization-api"
    Then the project should be created successfully
    And all available optimizations should be applied
    And maximum performance settings should be used
    And the project should compile without errors

  @levels @progression
  Scenario: Test optimization level progression
    Given I want to create projects with different optimization levels
    When I generate "test-none" with level "none"
    And I generate "test-safe" with level "safe"
    And I generate "test-standard" with level "standard"  
    And I generate "test-aggressive" with level "aggressive"
    And I generate "test-expert" with level "expert"
    Then each project should have progressively more optimizations
    And optimization coverage should increase with each level
    And all projects should compile without errors

  @levels @comparison
  Scenario: Compare optimization levels side by side
    Given I want to create the same project with different optimization levels
    When I generate "api-none" with level "none"
    And I generate "api-safe" with level "safe"
    And I generate "api-expert" with level "expert"
    Then I should be able to compare the differences
    And optimization metrics should show progression
    And file counts should reflect optimization impact
    And all projects should have equivalent functionality

  @levels @context-appropriate
  Scenario: Use context-appropriate optimization levels
    Given I want to optimize for different development contexts
    When I use "safe" level for "development" context
    And I use "standard" level for "testing" context
    And I use "aggressive" level for "production" context
    And I use "expert" level for "maintenance" context
    Then each context should get appropriate optimizations
    And warnings should be displayed for risky combinations
    And recommendations should be provided

  @levels @validation
  Scenario: Validate optimization level constraints
    Given I want to create a new "cli" project
    When I set an invalid optimization level "super-expert"
    Then the system should reject the invalid level
    And valid optimization levels should be listed
    And helpful guidance should be provided

  @levels @warnings
  Scenario: Display warnings for risky optimization levels
    Given I want to create a new "web-api" project
    And I set the optimization level to "aggressive"
    And I disable dry-run mode
    When I generate the project "risky-optimization-api"
    Then warnings should be displayed about risky optimizations
    And safety recommendations should be provided
    And the user should be able to proceed with confirmation
    And the project should compile without errors

  @levels @metrics
  Scenario: Track optimization metrics across levels
    Given I want to measure optimization effectiveness
    When I generate projects with different optimization levels
    Then metrics should be collected for each level
    And performance improvements should be measured
    And optimization impact should be quantified
    And metrics should be available for analysis

  @levels @upgrade
  Scenario: Upgrade optimization level for existing project
    Given I have an existing project with "safe" optimization level
    When I upgrade the optimization level to "standard"
    Then additional optimizations should be applied
    And existing optimizations should be preserved
    And backup files should be created
    And the project should still compile without errors

  @levels @downgrade
  Scenario: Downgrade optimization level with preservation
    Given I have an existing project with "aggressive" optimization level
    When I downgrade the optimization level to "safe"
    Then risky optimizations should be reverted safely
    And safe optimizations should be preserved
    And project functionality should remain intact
    And the project should compile without errors

  @levels @batch
  Scenario: Apply optimization levels to multiple projects
    Given I have multiple projects with different optimization needs
    When I apply "safe" level to development projects
    And I apply "standard" level to testing projects  
    And I apply "aggressive" level to production projects
    Then each project should receive appropriate optimizations
    And batch operation should complete successfully
    And all projects should compile without errors

  @levels @custom-rules
  Scenario: Define custom rules for optimization levels
    Given I want to customize optimization level behavior
    When I define custom rules for "aggressive" level
    And I specify which optimizations to include/exclude
    Then the custom rules should be applied correctly
    And the project should reflect the custom optimization logic
    And the project should compile without errors