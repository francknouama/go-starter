Feature: Optimization Configuration Profiles
  As a Go developer using go-starter's optimization system
  I want to manage and use optimization profiles effectively
  So that I can apply consistent optimization settings across projects

  Background:
    Given I am using go-starter CLI
    And the optimization system is available
    And the configuration system is available

  @configuration @profiles
  Scenario: Use predefined optimization profiles
    Given I have access to predefined optimization profiles
    When I list available profiles
    Then I should see the following profiles:
      | Profile      | Description                                          | Risk Level |
      | conservative | Minimal changes, maximum safety                      | Low        |
      | balanced     | Balanced optimization and safety                     | Medium     |
      | performance  | Performance-focused optimizations                    | Medium     |
      | aggressive   | Aggressive optimizations with some risk              | High       |
      | maximum      | Maximum optimization, minimal safety checks          | Very High  |
    And each profile should have appropriate default settings

  @configuration @profile-application
  Scenario Outline: Apply optimization profiles to projects
    Given I want to optimize a "<project_type>" project
    When I apply the "<profile>" optimization profile
    Then the project should be optimized according to "<profile>" settings
    And the optimization level should match the profile's configuration
    And profile-specific options should be enabled
    And the project should compile successfully

    Examples:
      | project_type | profile      |
      | web-api      | conservative |
      | web-api      | balanced     |
      | web-api      | performance  |
      | cli          | conservative |
      | cli          | performance  |
      | library      | balanced     |
      | microservice | aggressive   |

  @configuration @custom-profiles
  Scenario: Create custom optimization profile
    Given I want to create a custom optimization profile
    When I define a profile named "team-standard" with:
      | Setting                | Value      |
      | Level                  | standard   |
      | RemoveUnusedImports    | true       |
      | OrganizeImports        | true       |
      | RemoveUnusedVars       | false      |
      | RemoveUnusedFuncs      | false      |
      | CreateBackups          | true       |
      | MaxConcurrentFiles     | 8          |
    Then the profile should be saved successfully
    And I should be able to use "team-standard" profile
    And the profile should persist between sessions

  @configuration @profile-override
  Scenario: Override profile settings
    Given I am using the "balanced" profile
    When I override specific settings:
      | Setting             | Original | Override |
      | RemoveUnusedVars    | false    | true     |
      | MaxConcurrentFiles  | 4        | 16       |
    Then the overridden settings should take precedence
    And non-overridden settings should use profile defaults
    And the effective configuration should be clearly reported

  @configuration @profile-validation
  Scenario: Validate profile configurations
    Given I create a profile with invalid settings:
      | Setting            | Value    | Issue              |
      | Level              | "ultra"  | Invalid level      |
      | MaxConcurrentFiles | -1       | Negative value     |
      | MaxFileSize        | 0        | Zero not allowed   |
    When I attempt to use this profile
    Then validation should fail with clear error messages
    And suggestions for valid values should be provided
    And the profile should not be applied

  @configuration @profile-export-import
  Scenario: Export and import optimization profiles
    Given I have configured custom profiles:
      | Profile         | Level      | Key Settings                    |
      | project-alpha   | standard   | Conservative import handling    |
      | project-beta    | aggressive | Full optimization enabled       |
      | shared-library  | safe       | Minimal changes only           |
    When I export profiles to "optimization-profiles.json"
    Then the export should contain all custom profiles
    And profile settings should be preserved
    When I import profiles on another system
    Then all profiles should be available
    And settings should match the original

  @configuration @profile-precedence
  Scenario: Profile and flag precedence
    Given I have a "performance" profile configured
    When I run optimization with both profile and explicit flags:
      """
      go-starter optimize --profile=performance --level=safe --remove-unused-vars=false
      """
    Then explicit flags should override profile settings
    And the precedence order should be:
      | Priority | Source          | Example                    |
      | 1        | Command flags   | --level=safe              |
      | 2        | Profile         | performance profile        |
      | 3        | Defaults        | System defaults           |

  @configuration @profile-composition
  Scenario: Compose profiles from base configurations
    Given I have base profiles:
      | Profile    | Focus                |
      | imports    | Import optimization  |
      | variables  | Variable cleanup     |
      | functions  | Function removal     |
    When I create a composite profile "full-cleanup" that includes:
      | Base Profile | Settings Used        |
      | imports      | All import settings  |
      | variables    | Variable removal     |
      | functions    | Function analysis    |
    Then the composite profile should combine all settings
    And conflicts should be resolved by last-wins rule
    And the composition should be documented

  @configuration @profile-versioning
  Scenario: Handle profile version compatibility
    Given I have profiles from different go-starter versions:
      | Profile      | Version | Status      |
      | legacy-v1    | 1.0.0   | Deprecated  |
      | standard-v2  | 2.0.0   | Compatible  |
      | future-v3    | 3.0.0   | Unknown     |
    When I load these profiles
    Then compatible profiles should work normally
    And deprecated profiles should show migration warnings
    And incompatible profiles should be rejected with guidance

  @configuration @profile-conditions
  Scenario: Conditional profile application
    Given I have profiles with conditions:
      | Profile         | Condition                          |
      | ci-optimize     | Environment = CI                   |
      | dev-friendly    | Environment = Development          |
      | prod-aggressive | Environment = Production           |
      | go-1-21-plus    | Go Version >= 1.21                |
    When I run optimization in different environments
    Then the appropriate profile should be selected automatically
    And conditions should be evaluated correctly
    And manual override should still be possible

  @configuration @profile-recommendations
  Scenario: Get profile recommendations
    Given I have a "<project_type>" project with "<characteristics>"
    When I request profile recommendations
    Then the system should analyze the project
    And recommend appropriate profiles:
      | Project Type | Characteristics        | Recommended Profile |
      | web-api      | High traffic, critical | conservative        |
      | cli          | Internal tool          | balanced            |
      | library      | Public API             | conservative        |
      | microservice | Performance critical   | performance         |
    And provide reasoning for each recommendation

    Examples:
      | project_type | characteristics         |
      | web-api      | High traffic, critical  |
      | cli          | Internal tool           |
      | library      | Public API              |
      | microservice | Performance critical    |