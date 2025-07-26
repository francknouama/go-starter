Feature: Optimization Configuration Management
  As a Go developer using go-starter
  I want to manage optimization configurations effectively
  So that I can reuse and customize optimization settings for different projects

  Background:
    Given I am using go-starter CLI
    And the configuration system is available

  @configuration @profiles
  Scenario: Use predefined optimization profiles
    Given I want to create a new "web-api" project
    When I list available optimization profiles
    Then I should see profiles: "conservative", "balanced", "performance", "maximum"
    And each profile should have clear descriptions
    And profile settings should be documented

  @configuration @profiles @conservative
  Scenario: Generate project with conservative profile
    Given I want to create a new "cli" project
    And I set the optimization profile to "conservative"
    When I generate the project "conservative-cli"
    Then the project should be created successfully
    And only safe optimizations should be applied
    And dry-run mode should be enabled by default
    And the project should compile without errors

  @configuration @profiles @balanced
  Scenario: Generate project with balanced profile
    Given I want to create a new "web-api" project
    And I set the optimization profile to "balanced"
    When I generate the project "balanced-api"
    Then the project should be created successfully
    And standard level optimizations should be applied
    And backups should be created
    And the project should compile without errors

  @configuration @profiles @performance
  Scenario: Generate project with performance profile
    Given I want to create a new "cli" project
    And I set the optimization profile to "performance"
    When I generate the project "performance-cli"
    Then the project should be created successfully
    And aggressive optimizations should be applied
    And high concurrency settings should be used
    And the project should compile without errors

  @configuration @profiles @maximum
  Scenario: Generate project with maximum profile
    Given I want to create a new "web-api" project
    And I set the optimization profile to "maximum"
    When I generate the project "maximum-api"
    Then the project should be created successfully
    And expert level optimizations should be applied
    And all performance settings should be maximized
    And warnings should be displayed for risky settings
    And the project should compile without errors

  @configuration @custom
  Scenario: Create and use custom optimization profile
    Given I want to create a custom optimization profile "my-custom"
    And I set custom profile optimization level to "standard"
    And I enable unused variable removal
    And I disable missing import addition
    When I save the custom profile
    And I generate a project "custom-profile-test" using profile "my-custom"
    Then the project should be created successfully
    And custom optimization settings should be applied
    And the project should compile without errors

  @configuration @persistence
  Scenario: Save and load optimization configuration
    Given I have created a custom optimization profile "persistent-profile"
    When I save the configuration to a file
    And I load the configuration from the file
    Then the configuration should be restored correctly
    And all custom settings should be preserved
    And the profile should be available for use

  @configuration @validation
  Scenario: Validate optimization configuration
    Given I create an optimization configuration
    When I set invalid optimization level "invalid-level"
    Then configuration validation should fail
    And a helpful error message should be displayed
    And valid optimization levels should be suggested

  @configuration @validation @conflicting
  Scenario: Detect conflicting configuration options
    Given I create an optimization configuration
    And I enable write-optimized-files
    And I enable dry-run mode
    When I validate the configuration
    Then validation should fail with conflict error
    And the conflicting options should be identified
    And resolution suggestions should be provided

  @configuration @override
  Scenario: Override profile settings with explicit options
    Given I want to create a new "cli" project
    And I set the optimization profile to "conservative"
    And I override the dry-run setting to "false"
    When I generate the project "override-test-cli"
    Then the project should be created successfully
    And the conservative profile should be used as base
    And the dry-run override should be applied
    And the project should compile without errors

  @configuration @migration
  Scenario: Migrate configuration between versions
    Given I have an old optimization configuration file
    When I load the configuration with go-starter
    Then the configuration should be migrated automatically
    And all settings should be preserved where possible
    And migration warnings should be displayed if needed
    And the updated configuration should be saved

  @configuration @context-aware
  Scenario: Get context-aware optimization recommendations
    Given I want to create a new "web-api" project
    When I request optimization recommendations for "development" context
    Then the system should recommend "safe" optimization level
    And appropriate warning messages should be shown
    And alternative recommendations should be provided

  @configuration @multiple-projects
  Scenario: Apply same configuration to multiple projects
    Given I have saved an optimization configuration "multi-project-config"
    When I generate projects "api-1", "api-2", "cli-1" using the same configuration
    Then all projects should be created successfully
    And consistent optimization should be applied across projects
    And all projects should compile without errors

  @configuration @export-import
  Scenario: Export and import optimization profiles
    Given I have multiple custom optimization profiles
    When I export the profiles to a JSON file
    And I import the profiles in a different environment
    Then all profiles should be available
    And all custom settings should be preserved
    And profiles should work correctly in the new environment