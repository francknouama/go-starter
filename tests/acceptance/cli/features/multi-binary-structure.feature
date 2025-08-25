Feature: Multi-Binary Structure Support
  As a developer using go-starter
  I want the new multi-binary structure to work correctly
  So that I can use the appropriate binary for my needs

  Background:
    Given I have the go-starter project with multi-binary structure
    And all binaries can be built successfully

  Scenario: All binaries compile independently
    When I build the CLI binary from "cmd/go-starter"
    Then the build should succeed
    And the binary should be executable
    And the binary size should be reasonable

    When I build the dev server from "cmd/go-starter-dev" 
    Then the build should succeed
    And the binary should be executable
    And the binary size should be reasonable

    When I build the web server from "web/cmd/web-server"
    Then the build should succeed  
    And the binary should be executable
    And the binary size should be reasonable

    When I build the legacy binary from "."
    Then the build should succeed
    And the binary should be executable
    And the binary size should be reasonable

  Scenario: Installation paths work correctly
    When I install the CLI tool using "go install ./cmd/go-starter"
    Then the installation should succeed
    And "go-starter" binary should be available in GOPATH/bin
    And the binary should be functional

    When I install the dev server using "go install ./cmd/go-starter-dev"
    Then the installation should succeed
    And "go-starter-dev" binary should be available in GOPATH/bin

    When I install using legacy method "go install ."
    Then the installation should succeed
    And "go-starter" binary should be available in GOPATH/bin
    And the binary should work with possible deprecation warning

  Scenario: Legacy binary shows deprecation warning
    Given I have built the legacy binary from root directory
    When I run the binary with "--help" flag
    Then I should see a deprecation warning
    And the warning should mention new binary locations
    And it should show "go build -o go-starter ./cmd/go-starter"
    And it should show "go build -o go-starter-web ./web/cmd/web-server"  
    And it should show "go build -o go-starter-dev ./cmd/go-starter-dev"
    And the CLI functionality should still work

  Scenario: CLI tool functionality works correctly
    Given I have built the CLI binary
    When I run "version" command
    Then I should see version information
    And the command should succeed

    When I run "list" command
    Then I should see available blueprints
    And the list should include "web-api", "cli", "library"
    And the command should succeed

    When I run "new test-project --type=cli --dry-run"
    Then I should see files to be generated
    And the list should include "main.go", "go.mod"
    And the command should succeed

  Scenario: Development server starts correctly
    Given I have built the dev server binary
    When I start the development server
    Then the server should start without immediate errors
    And it should be able to handle termination gracefully

  Scenario: Production web server starts correctly  
    Given I have built the web server binary
    When I start the web server
    Then the server should start without immediate errors
    And it should be able to handle termination gracefully

  Scenario: Embedded assets work without filesystem access
    Given I have built the CLI binary with embedded blueprints
    When I run the CLI from an isolated directory without blueprint files
    And I execute "list" command
    Then the command should succeed
    And I should see embedded blueprints listed
    And the list should include "web-api", "cli", "library"

  Scenario: Generated projects compile successfully
    Given I have the CLI binary with embedded blueprints
    When I generate a project using "new test-project --type=cli --complexity=simple"
    Then the project generation should succeed
    And the generated project should compile with "go build ./..."
    And all dependencies should be resolved correctly

  Scenario: Cross-platform compatibility
    Given I am running on the current platform
    When I build and run any binary
    Then path handling should work correctly for the platform
    And binary formats should be appropriate for the platform
    And file operations should use platform-specific conventions

  Scenario: Migration experience is smooth
    Given I am upgrading from the legacy single-binary structure
    When I see the deprecation warning
    Then the migration instructions should be clear and actionable
    And following the instructions should result in working binaries
    And no existing functionality should be broken
    And all existing commands should continue to work

  Scenario: Performance remains acceptable
    Given the new multi-binary structure
    When I build each binary
    Then build times should be under 30 seconds each
    And binary sizes should be reasonable (5MB - 50MB)
    
    When I run the CLI for quick operations
    Then startup time should be under 2 seconds
    And response time should be acceptable for user interaction

  Scenario: Backward compatibility is maintained
    Given existing users with established workflows
    When they use the new CLI binary
    Then all existing commands should work identically
    And command syntax should remain unchanged
    And output format should be consistent
    And no breaking changes should be introduced

  Scenario: Documentation examples remain valid
    Given the documentation contains usage examples
    When those examples are executed with new binaries
    Then all examples should work correctly
    And help text should be consistent
    And error messages should be clear and helpful