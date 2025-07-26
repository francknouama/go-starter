Feature: Optimization Configuration Persistence
  As a Go developer using go-starter's optimization system
  I want my optimization configurations to persist across sessions
  So that I can maintain consistent settings and share them with my team

  Background:
    Given I am using go-starter CLI
    And the optimization system is available
    And the configuration system is available

  @configuration @persistence
  Scenario: Save optimization configuration to project
    Given I have a project that needs optimization
    When I configure optimization settings:
      | Setting                | Value       |
      | Level                  | standard    |
      | RemoveUnusedImports    | true        |
      | OrganizeImports        | true        |
      | RemoveUnusedVars       | true        |
      | CreateBackups          | true        |
      | SkipTestFiles          | true        |
    And I save the configuration to ".go-starter-optimize.json"
    Then the configuration file should be created
    And it should contain all specified settings
    And the file should be properly formatted

  @configuration @auto-load
  Scenario: Auto-load project configuration
    Given I have a project with ".go-starter-optimize.json" configuration
    When I run optimization without specifying settings
    Then the configuration should be loaded automatically
    And the loaded settings should be applied
    And a message should confirm configuration source

  @configuration @config-search
  Scenario: Configuration file search order
    Given I have optimization configs in multiple locations:
      | Location                          | Priority |
      | ./.go-starter-optimize.json       | 1        |
      | ./config/optimization.json        | 2        |
      | ~/.go-starter/optimization.json   | 3        |
      | /etc/go-starter/optimization.json | 4        |
    When I run optimization
    Then configs should be loaded in priority order
    And first found config should be used
    And search locations should be logged in verbose mode

  @configuration @config-merge
  Scenario: Merge multiple configuration sources
    Given I have a global config in "~/.go-starter/optimization.json":
      """json
      {
        "defaultLevel": "safe",
        "author": "TeamLead",
        "createBackups": true
      }
      """
    And I have a project config in ".go-starter-optimize.json":
      """json
      {
        "level": "standard",
        "removeUnusedVars": true
      }
      """
    When I run optimization with CLI flags:
      """
      --level=aggressive --dry-run
      """
    Then configurations should merge with precedence:
      | Source        | Level      | Backups | UnusedVars | DryRun |
      | Global        | safe       | true    | -          | -      |
      | Project       | standard   | true    | true       | -      |
      | CLI Flags     | aggressive | true    | true       | true   |
      | Final Result  | aggressive | true    | true       | true   |

  @configuration @config-schema
  Scenario: Validate configuration schema
    Given I have a configuration file with:
      """json
      {
        "level": "standard",
        "options": {
          "removeUnusedImports": true,
          "organizeImports": true,
          "removeUnusedVars": false
        },
        "performance": {
          "maxConcurrentFiles": 8,
          "maxFileSize": 2097152
        },
        "safety": {
          "createBackups": true,
          "dryRun": false
        }
      }
      """
    When I load this configuration
    Then schema validation should pass
    And all settings should be properly typed
    And nested options should be correctly parsed

  @configuration @config-migration
  Scenario: Migrate old configuration formats
    Given I have an old format configuration:
      """json
      {
        "optimizationLevel": 2,
        "removeImports": true,
        "backupFiles": true
      }
      """
    When I load this configuration
    Then it should be automatically migrated to new format:
      """json
      {
        "level": "standard",
        "options": {
          "removeUnusedImports": true
        },
        "safety": {
          "createBackups": true
        }
      }
      """
    And a migration report should be generated
    And the user should be prompted to save the new format

  @configuration @config-templates
  Scenario: Use configuration templates
    Given go-starter provides configuration templates
    When I run "go-starter optimize init"
    Then I should see available templates:
      | Template     | Description                           |
      | minimal      | Minimal safe optimizations            |
      | recommended  | Recommended balanced settings         |
      | ci-cd        | Settings optimized for CI/CD          |
      | development  | Developer-friendly settings           |
      | production   | Production-ready optimizations        |
    When I select "recommended" template
    Then a configuration file should be created with template settings

  @configuration @config-validation-errors
  Scenario: Handle invalid configuration gracefully
    Given I have a malformed configuration file:
      """json
      {
        "level": "ultra-extreme",
        "options": {
          "removeUnusedImports": "yes",
          "maxConcurrentFiles": "many"
        }
      }
      """
    When I attempt to load this configuration
    Then validation should fail with specific errors:
      | Field                      | Error                           | Suggestion                    |
      | level                      | Invalid optimization level      | Use: none,safe,standard,...   |
      | options.removeUnusedImports| Expected boolean, got string    | Use: true or false           |
      | options.maxConcurrentFiles | Expected number, got string     | Use: numeric value > 0       |
    And the system should fall back to defaults
    And offer to create a valid configuration

  @configuration @config-comments
  Scenario: Support configuration with comments
    Given I want documented configuration
    When I generate a configuration with comments
    Then the file should include helpful comments:
      ```jsonc
      {
        // Optimization level: none, safe, standard, aggressive, expert
        "level": "standard",
        
        "options": {
          // Import optimization settings
          "removeUnusedImports": true,  // Remove imports not used in code
          "organizeImports": true,       // Sort and group imports
          
          // Code optimization settings  
          "removeUnusedVars": false,     // Risky: might remove vars with side effects
          "removeUnusedFuncs": false     // Very risky: might break external APIs
        }
      }
      ```

  @configuration @config-env-vars
  Scenario: Override configuration with environment variables
    Given I have a configuration file with standard settings
    And I set environment variables:
      | Variable                          | Value      |
      | GO_STARTER_OPT_LEVEL             | aggressive |
      | GO_STARTER_OPT_DRY_RUN           | true       |
      | GO_STARTER_OPT_CREATE_BACKUPS    | false      |
    When I run optimization
    Then environment variables should override file settings
    And the override source should be logged
    And a warning should be shown for security-sensitive overrides

  @configuration @config-lock
  Scenario: Lock configuration for team consistency
    Given I have a team-agreed configuration
    When I add a "configLock" section:
      """json
      {
        "level": "standard",
        "configLock": {
          "enabled": true,
          "lockedBy": "tech-lead",
          "lockedAt": "2024-01-15T10:00:00Z",
          "reason": "Approved team configuration - do not modify",
          "allowedOverrides": ["dryRun", "verbose"]
        }
      }
      """
    Then most settings should be immutable
    And only allowed overrides should work
    And attempts to change locked settings should show lock info