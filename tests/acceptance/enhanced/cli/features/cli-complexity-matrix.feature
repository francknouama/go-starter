Feature: CLI Complexity Matrix Testing
  As a go-starter user building command-line applications
  I want to ensure that different CLI complexity levels work correctly
  So that I can choose the appropriate complexity for my use case

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And I am testing CLI complexity combinations

  @critical @cli @simple @p1
  Scenario Outline: Simple CLI complexity validation
    When I generate a simple CLI tool with minimal configuration:
      | type           | cli        |
      | complexity     | simple     |
      | framework      | <framework>|
      | logger         | <logger>   |
      | go_version     | 1.23       |
      | expected_files | 12         |
    Then the project should compile successfully
    And the CLI tool should be executable
    And the project should have approximately 12 files
    And the project structure should be simple
    And simple CLI structure should be minimal and focused
    And Cobra framework should be properly integrated
    And <logger> logger should be properly integrated
    And progressive disclosure should work correctly
    And simple CLI should have minimal learning curve

    Examples:
      | framework | logger  |
      | cobra     | slog    |
      | cobra     | zap     |
      | cobra     | logrus  |
      | cobra     | zerolog |

  @critical @cli @standard @p1
  Scenario Outline: Standard CLI complexity validation
    When I generate a production-ready CLI application:
      | type           | cli        |
      | complexity     | standard   |
      | framework      | <framework>|
      | logger         | <logger>   |
      | go_version     | 1.23       |
      | expected_files | 27         |
    Then the project should compile successfully
    And the CLI tool should be executable
    And the project should have approximately 27 files
    And the project structure should be standard
    And standard CLI structure should be production-ready
    And Cobra framework should be properly integrated
    And <logger> logger should be properly integrated
    And command structure should follow Cobra conventions
    And subcommands should be properly organized
    And help system should be well-implemented
    And configuration management should work correctly
    And config files should be supported
    And environment variables should be handled
    And command-line flags should override config
    And testing framework should be properly set up
    And command tests should be available
    And documentation should be appropriate for complexity level
    And README should explain CLI usage
    And help text should be comprehensive
    And standard CLI should provide full functionality

    Examples:
      | framework | logger  |
      | cobra     | slog    |
      | cobra     | zap     |
      | cobra     | logrus  |
      | cobra     | zerolog |

  @integration @cli @comparison @p1
  Scenario Outline: CLI complexity comparison testing
    When I generate a CLI project with complexity configuration:
      | type           | cli           |
      | complexity     | <complexity>  |
      | framework      | cobra         |
      | logger         | slog          |
      | go_version     | 1.23         |
      | expected_files | <file_count>  |
    Then the project should compile successfully
    And the CLI tool should be executable
    And the project should have approximately <file_count> files
    And the project structure should be <complexity>
    And Cobra framework should be properly integrated
    And slog logger should be properly integrated
    And logging should be consistent across commands
    And log levels should be configurable
    And progressive disclosure should work correctly
    And test coverage should be appropriate for complexity
    And migration path from simple to standard should be clear

    Examples:
      | complexity | file_count |
      | simple     | 12         |
      | standard   | 27         |

  @performance @cli @optimization @p1
  Scenario Outline: CLI performance across complexity levels
    When I generate a CLI project with complexity configuration:
      | type           | cli           |
      | complexity     | <complexity>  |
      | framework      | cobra         |
      | logger         | <logger>      |
      | go_version     | 1.23         |
      | expected_files | <file_count>  |
    Then the project should compile successfully
    And the CLI tool should be executable
    And CLI startup time should be fast
    And memory usage should be minimal
    And binary size should be reasonable for complexity
    And build process should work correctly
    And Makefile should provide useful targets

    Examples:
      | complexity | file_count | logger  |
      | simple     | 12         | slog    |
      | simple     | 12         | zap     |
      | standard   | 27         | slog    |
      | standard   | 27         | zap     |
      | standard   | 27         | zerolog |

  @integration @cli @packaging @p1
  Scenario Outline: CLI packaging and distribution
    When I generate a production-ready CLI application:
      | type           | cli           |
      | complexity     | <complexity>  |
      | framework      | cobra         |
      | logger         | zap           |
      | go_version     | 1.23         |
      | expected_files | <file_count>  |
    Then the project should compile successfully
    And the CLI tool should be executable
    And build process should work correctly
    And Makefile should provide useful targets
    And Docker support should be included if appropriate
    And release process should be documented
    And documentation should be appropriate for complexity level

    Examples:
      | complexity | file_count |
      | simple     | 12         |
      | standard   | 27         |

  @cross-platform @cli @compatibility @p1
  Scenario Outline: CLI cross-platform compatibility
    When I generate a CLI project with complexity configuration:
      | type           | cli           |
      | complexity     | <complexity>  |
      | framework      | cobra         |
      | logger         | slog          |
      | go_version     | 1.23         |
      | expected_files | <file_count>  |
    Then the project should compile successfully
    And the CLI tool should be executable
    And CLI should work on multiple platforms
    And build targets should include common platforms
    And configuration paths should be OS-appropriate
    And build process should work correctly

    Examples:
      | complexity | file_count |
      | simple     | 12         |
      | standard   | 27         |

  @integration @cli @progressive-disclosure @p1
  Scenario: Progressive disclosure system validation
    When I generate a simple CLI tool with minimal configuration:
      | type           | cli     |
      | complexity     | simple  |
      | framework      | cobra   |
      | logger         | slog    |
      | go_version     | 1.23    |
      | expected_files | 8       |
    Then the project should compile successfully
    And the CLI tool should be executable
    And simple CLI structure should be minimal and focused
    And simple CLI should have minimal learning curve
    And progressive disclosure should work correctly
    
    When I generate a production-ready CLI application:
      | type           | cli      |
      | complexity     | standard |
      | framework      | cobra    |
      | logger         | slog     |
      | go_version     | 1.23     |
      | expected_files | 29       |
    Then the project should compile successfully
    And the CLI tool should be executable
    And standard CLI structure should be production-ready
    And standard CLI should provide full functionality
    And migration path from simple to standard should be clear

  @security @cli @validation @p1
  Scenario Outline: CLI security and validation
    When I generate a CLI project with complexity configuration:
      | type           | cli           |
      | complexity     | <complexity>  |
      | framework      | cobra         |
      | logger         | slog          |
      | go_version     | 1.23         |
      | expected_files | <file_count>  |
    Then the project should compile successfully
    And the CLI tool should be executable
    And configuration management should work correctly
    And environment variables should be handled
    And command-line flags should override config
    And help system should be well-implemented

    Examples:
      | complexity | file_count |
      | simple     | 12         |
      | standard   | 27         |

  @integration @cli @logging @consistency @p1
  Scenario Outline: Logging consistency across CLI complexity levels
    When I generate a CLI project with complexity configuration:
      | type           | cli           |
      | complexity     | <complexity>  |
      | framework      | cobra         |
      | logger         | <logger>      |
      | go_version     | 1.23         |
      | expected_files | <file_count>  |
    Then the project should compile successfully
    And the CLI tool should be executable
    And <logger> logger should be properly integrated
    And logging should be consistent across commands
    And log levels should be configurable
    And progressive disclosure should work correctly

    Examples:
      | complexity | file_count | logger  |
      | simple     | 12         | slog    |
      | simple     | 12         | zap     |
      | simple     | 12         | logrus  |
      | simple     | 12         | zerolog |
      | standard   | 27         | slog    |
      | standard   | 27         | zap     |
      | standard   | 27         | logrus  |
      | standard   | 27         | zerolog |

  @integration @cli @testing @coverage @p1
  Scenario Outline: Testing framework integration across complexity levels
    When I generate a CLI project with complexity configuration:
      | type           | cli           |
      | complexity     | <complexity>  |
      | framework      | cobra         |
      | logger         | slog          |
      | go_version     | 1.23         |
      | expected_files | <file_count>  |
    Then the project should compile successfully
    And the CLI tool should be executable
    And testing framework should be properly set up
    And command tests should be available
    And test coverage should be appropriate for complexity
    And build process should work correctly

    Examples:
      | complexity | file_count |
      | simple     | 12         |
      | standard   | 27         |