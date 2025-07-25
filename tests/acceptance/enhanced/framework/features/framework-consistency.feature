Feature: Framework Consistency Across Architectures
  As a go-starter user building web applications with different architectural patterns
  I want to ensure that web frameworks work consistently across all supported architectures
  So that I can choose architecture patterns freely without framework limitations

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And I am testing framework consistency across architectures

  @critical @framework @gin @p1
  Scenario Outline: Gin framework consistency across architectures
    When I generate a web API project with framework consistency configuration:
      | type         | web-api        |
      | framework    | gin            |
      | architecture | <architecture> |
      | logger       | slog           |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And gin framework should be properly integrated
    And <architecture> architecture should be properly implemented
    And framework structure should be consistent with <architecture> patterns
    And routing should follow gin conventions in <architecture> architecture
    And middleware integration should be architecture-aware
    And dependency injection should work with gin and <architecture>
    And error handling should be consistent across layers
    And configuration management should integrate properly
    And health checks should be implemented consistently
    And logging integration should work across all layers
    And testing patterns should be architecture-appropriate
    And documentation should reflect architecture choices

    Examples:
      | architecture |
      | standard     |
      | clean        |
      | ddd          |
      | hexagonal    |

  @critical @framework @echo @p1
  Scenario Outline: Echo framework consistency across architectures
    When I generate a web API project with framework consistency configuration:
      | type         | web-api        |
      | framework    | echo           |
      | architecture | <architecture> |
      | logger       | zap            |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And echo framework should be properly integrated
    And <architecture> architecture should be properly implemented
    And framework structure should be consistent with <architecture> patterns
    And routing should follow echo conventions in <architecture> architecture
    And middleware integration should be architecture-aware
    And dependency injection should work with echo and <architecture>
    And error handling should be consistent across layers
    And configuration management should integrate properly
    And health checks should be implemented consistently
    And logging integration should work across all layers
    And testing patterns should be architecture-appropriate

    Examples:
      | architecture |
      | standard     |
      | clean        |
      | ddd          |
      | hexagonal    |

  @critical @framework @fiber @p1
  Scenario Outline: Fiber framework consistency across architectures
    When I generate a web API project with framework consistency configuration:
      | type         | web-api        |
      | framework    | fiber          |
      | architecture | <architecture> |
      | logger       | zerolog        |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And fiber framework should be properly integrated
    And <architecture> architecture should be properly implemented
    And framework structure should be consistent with <architecture> patterns
    And routing should follow fiber conventions in <architecture> architecture
    And middleware integration should be architecture-aware
    And dependency injection should work with fiber and <architecture>
    And error handling should be consistent across layers
    And configuration management should integrate properly
    And health checks should be implemented consistently
    And logging integration should work across all layers
    And testing patterns should be architecture-appropriate

    Examples:
      | architecture |
      | standard     |
      | clean        |
      | ddd          |
      | hexagonal    |

  @critical @framework @chi @p1
  Scenario Outline: Chi framework consistency across architectures
    When I generate a web API project with framework consistency configuration:
      | type         | web-api        |
      | framework    | chi            |
      | architecture | <architecture> |
      | logger       | logrus         |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And chi framework should be properly integrated
    And <architecture> architecture should be properly implemented
    And framework structure should be consistent with <architecture> patterns
    And routing should follow chi conventions in <architecture> architecture
    And middleware integration should be architecture-aware
    And dependency injection should work with chi and <architecture>
    And error handling should be consistent across layers
    And configuration management should integrate properly
    And health checks should be implemented consistently
    And logging integration should work across all layers
    And testing patterns should be architecture-appropriate

    Examples:
      | architecture |
      | standard     |
      | clean        |
      | ddd          |
      | hexagonal    |

  @performance @framework @optimization @p1
  Scenario Outline: Framework performance consistency across architectures
    When I generate a web API project with framework consistency configuration:
      | type         | web-api        |
      | framework    | <framework>    |
      | architecture | <architecture> |
      | logger       | zap            |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And framework performance should be optimized for <architecture>
    And memory usage should be efficient
    And startup time should be reasonable
    And logging integration should work across all layers

    Examples:
      | framework | architecture |
      | gin       | standard     |
      | gin       | clean        |
      | gin       | hexagonal    |
      | echo      | standard     |
      | echo      | clean        |
      | echo      | ddd          |
      | fiber     | standard     |
      | fiber     | hexagonal    |
      | chi       | standard     |
      | chi       | clean        |

  @security @framework @middleware @p1
  Scenario Outline: Security middleware consistency across frameworks and architectures
    When I generate a web API project with framework consistency configuration:
      | type         | web-api        |
      | framework    | <framework>    |
      | architecture | <architecture> |
      | logger       | slog           |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And security middleware should be properly integrated
    And CORS configuration should be framework-specific
    And input validation should be implemented consistently
    And error handling should be consistent across layers
    And logging integration should work across all layers

    Examples:
      | framework | architecture |
      | gin       | standard     |
      | gin       | clean        |
      | gin       | ddd          |
      | gin       | hexagonal    |
      | echo      | standard     |
      | echo      | clean        |
      | echo      | hexagonal    |
      | fiber     | standard     |
      | fiber     | ddd          |
      | chi       | standard     |
      | chi       | clean        |

  @integration @framework @detailed @p1
  Scenario Outline: Detailed framework-architecture integration testing
    When I generate a web API project with framework consistency configuration:
      | type         | web-api        |
      | framework    | <framework>    |
      | architecture | <architecture> |
      | logger       | <logger>       |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And <framework> framework should be properly integrated
    And <architecture> architecture should be properly implemented
    And framework structure should be consistent with <architecture> patterns
    And routing should follow <framework> conventions in <architecture> architecture
    And middleware integration should be architecture-aware
    And dependency injection should work with <framework> and <architecture>
    And error handling should be consistent across layers
    And configuration management should integrate properly
    And health checks should be implemented consistently
    And logging integration should work across all layers
    And testing patterns should be architecture-appropriate
    And documentation should reflect architecture choices

    Examples:
      | framework | architecture | logger  |
      | gin       | standard     | slog    |
      | gin       | clean        | zap     |
      | gin       | ddd          | zerolog |
      | gin       | hexagonal    | logrus  |
      | echo      | standard     | zap     |
      | echo      | clean        | slog    |
      | echo      | ddd          | zerolog |
      | echo      | hexagonal    | logrus  |
      | fiber     | standard     | zerolog |
      | fiber     | clean        | zap     |
      | fiber     | ddd          | slog    |
      | fiber     | hexagonal    | logrus  |
      | chi       | standard     | logrus  |
      | chi       | clean        | slog    |
      | chi       | ddd          | zap     |
      | chi       | hexagonal    | zerolog |

  @architecture @framework @validation @p1
  Scenario Outline: Architecture-specific framework integration validation
    When I generate a web API project with framework consistency configuration:
      | type         | web-api        |
      | framework    | <framework>    |
      | architecture | <architecture> |
      | logger       | slog           |
      | go_version   | 1.23          |
    Then the project should compile successfully
    And <architecture>-specific validation should pass for <framework>

    Examples:
      | framework | architecture |
      | gin       | standard     |
      | gin       | clean        |
      | gin       | ddd          |
      | gin       | hexagonal    |
      | echo      | standard     |
      | echo      | clean        |
      | echo      | ddd          |
      | echo      | hexagonal    |
      | fiber     | standard     |
      | fiber     | clean        |
      | fiber     | ddd          |
      | fiber     | hexagonal    |
      | chi       | standard     |
      | chi       | clean        |
      | chi       | ddd          |
      | chi       | hexagonal    |