Feature: Cross-Platform Compatibility Testing
  As a go-starter maintainer
  I want to ensure go-starter works consistently across all supported platforms
  So that users have a reliable experience regardless of their operating system

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And cross-platform testing is enabled

  Scenario: Platform detection and adaptation
    Given I am running on platform "<platform>"
    When I check platform-specific configurations
    Then the system should detect the platform correctly
    And platform-specific paths should be normalized
    And file permissions should be set appropriately for the platform
    And line ending handling should match platform conventions

    Examples:
      | platform |
      | windows  |
      | darwin   |
      | linux    |

  Scenario: File system compatibility across platforms
    Given I generate projects on different platforms
    When I test file system operations
    Then file paths should use correct separators for each platform
    And directory creation should work consistently
    And file permissions should be preserved appropriately
    And case sensitivity should be handled correctly
    And special characters in paths should be supported

  Scenario: Project generation consistency across platforms
    Given I use the same configuration on all platforms
    When I generate identical projects on:
      | Platform | Expected Files | Binary Extension |
      | windows  | 45            | .exe             |
      | darwin   | 45            | (none)           |
      | linux    | 45            | (none)           |
    Then all platforms should generate the same number of files
    And file contents should be identical across platforms
    And only platform-specific files should differ
    And binary outputs should have correct extensions

  Scenario: Compilation and execution across platforms
    Given I have generated projects on multiple platforms
    When I compile and run the projects
    Then compilation should succeed on all platforms
    And execution should produce consistent results
    And dependencies should resolve correctly
    And Go module handling should be consistent

  Scenario: Performance characteristics across platforms
    Given I benchmark operations on different platforms
    When I measure performance metrics
    Then performance variance should be within acceptable limits:
      | Metric            | Windows Variance | macOS Variance | Linux Baseline |
      | Generation Time   | < 20%           | < 10%          | 1.0x           |
      | Compilation Time  | < 30%           | < 15%          | 1.0x           |
      | Binary Size       | < 5%            | < 5%           | 1.0x           |
      | Memory Usage      | < 25%           | < 15%          | 1.0x           |
      | Startup Time      | < 20%           | < 10%          | 1.0x           |
    And performance degradation should be documented
    And platform-specific optimizations should be noted

  Scenario: Shell integration and scripting
    Given I test shell integration on different platforms
    When I run shell scripts and commands
    Then bash scripts should work on Unix-like systems
    And PowerShell scripts should work on Windows
    And batch files should execute correctly on Windows
    And shell completion should work appropriately
    And environment variable handling should be consistent

  Scenario: Path handling and normalization
    Given I work with various path formats
    When I process paths on different platforms:
      | Platform | Path Format                    | Expected Result                |
      | windows  | C:\Users\test\project         | C:\Users\test\project         |
      | windows  | /c/Users/test/project         | C:\Users\test\project         |
      | darwin   | /Users/test/project           | /Users/test/project           |
      | linux    | /home/test/project            | /home/test/project            |
      | all      | ./relative/path               | (platform-appropriate)        |
    Then paths should be normalized correctly for each platform
    And relative paths should resolve appropriately
    And path separators should be converted correctly
    And UNC paths should be handled on Windows

  Scenario: Unicode and character encoding support
    Given I use projects with international characters
    When I test character encoding across platforms
    Then Unicode filenames should be supported consistently
    And file contents with UTF-8 should be preserved
    And console output should display correctly
    And configuration files should maintain encoding

  Scenario: Docker and containerization compatibility
    Given I test containerized deployments
    When I build Docker images on different platforms
    Then multi-platform Docker builds should succeed
    And ARM64 and AMD64 architectures should be supported
    And container images should be consistent across platforms
    And Docker Compose should work on all platforms

  Scenario: CI/CD platform integration
    Given I run tests in different CI environments
    When I execute the test suite on:
      | CI Platform      | OS           | Expected Result |
      | GitHub Actions   | ubuntu-latest| Success         |
      | GitHub Actions   | windows-latest| Success        |
      | GitHub Actions   | macos-latest | Success         |
      | GitLab CI       | ubuntu       | Success         |
      | Azure DevOps    | windows      | Success         |
      | CircleCI        | linux        | Success         |
    Then all CI platforms should execute tests successfully
    And test results should be consistent
    And artifacts should be generated correctly

  Scenario: Network and proxy configuration
    Given I test network operations across platforms
    When I configure proxy settings and network access
    Then HTTP/HTTPS proxy settings should be respected
    And corporate firewall restrictions should be handled
    And DNS resolution should work consistently
    And certificate validation should be platform-appropriate

  Scenario: Archive and compression handling
    Given I work with compressed files and archives
    When I handle different archive formats
    Then ZIP files should be supported on all platforms
    And TAR.GZ files should work on Unix-like systems
    And extraction should preserve permissions and timestamps
    And compression ratios should be consistent

  Scenario: Temporary directory and cleanup
    Given I use temporary directories for testing
    When I create and clean up temporary resources
    Then temp directories should use platform-appropriate locations
    And cleanup should remove all temporary files
    And permissions should allow proper cleanup
    And concurrent access should be handled safely

  Scenario: Environment variable and configuration
    Given I work with environment variables and config
    When I set and read configuration values
    Then environment variables should be accessible consistently
    And configuration file locations should follow platform conventions
    And user-specific settings should be stored appropriately
    And system-wide settings should be supported

  Scenario: Binary distribution and packaging
    Given I package binaries for distribution
    When I create platform-specific packages
    Then Windows packages should include .exe files
    And Unix packages should have correct permissions
    And Package managers should be supported appropriately:
      | Platform | Package Manager | Package Format |
      | windows  | Chocolatey     | .nupkg         |
      | darwin   | Homebrew       | Formula        |
      | linux    | apt/yum        | .deb/.rpm      |
    And installation should work from package managers

  Scenario: Integration with platform-specific tools
    Given I integrate with platform-specific development tools
    When I test tool integration
    Then Visual Studio Code should work on all platforms
    And GoLand/IntelliJ should have consistent behavior
    And Platform-specific IDEs should be supported
    And Debugging tools should function correctly

  Scenario: Security and permissions model
    Given I work with different permission models
    When I handle security and permissions
    Then Unix permissions should be preserved and respected
    And Windows ACLs should be handled appropriately
    And Executable permissions should be set correctly
    And Security contexts should be maintained

  Scenario: Locale and internationalization
    Given I test with different system locales
    When I run go-starter with various locale settings
    Then Date and time formatting should be appropriate
    And Number formatting should follow locale conventions
    And Error messages should be displayed correctly
    And Sorting should respect locale collation rules

  Scenario: Platform-specific regression testing
    Given I maintain a regression test suite
    When I run platform-specific regression tests
    Then Previously fixed issues should not reoccur
    And Platform-specific bug fixes should be validated
    And Performance regressions should be detected
    And New platform versions should be tested automatically

  Scenario: Resource utilization and limits
    Given I monitor resource usage across platforms
    When I measure system resource consumption
    Then Memory usage should be within acceptable limits
    And CPU usage should be efficient on all platforms
    And Disk I/O should be optimized per platform
    And Network usage should be minimal and consistent

  Scenario: Error handling and diagnostics
    Given I test error conditions across platforms
    When I encounter errors and exceptions
    Then Error messages should be clear and platform-appropriate
    And Stack traces should be useful for debugging
    And Log files should be created in appropriate locations
    And Diagnostic information should be platform-specific when needed