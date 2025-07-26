Feature: Optimization Matrix Integration
  As a Go developer using go-starter's enhanced system
  I want optimization to work seamlessly with matrix-generated projects
  So that I can optimize any blueprint combination with confidence

  Background:
    Given I am using go-starter CLI
    And the optimization system is available
    And the matrix testing system is available
    And all templates are properly initialized

  @integration @optimization @matrix
  Scenario: Apply optimization to matrix-generated web-api projects
    Given I generate projects using matrix configurations:
      | framework | database | architecture | auth    | logger |
      | gin       | postgres | hexagonal    | jwt     | slog   |
      | echo      | mysql    | clean        | oauth2  | zap    |
      | fiber     | sqlite   | ddd          | session | logrus |
      | chi       | postgres | standard     | jwt     | zerolog|
    When I apply "standard" optimization to each project
    Then all projects should maintain functionality
    And optimization should improve code quality metrics
    And compilation should succeed for all combinations

  @integration @optimization @matrix @levels
  Scenario Outline: Matrix projects with different optimization levels
    Given I generate a "<framework>" project with "<database>" and "<architecture>"
    When I apply "<optimization_level>" optimization
    Then the project should compile successfully
    And optimization metrics should show appropriate improvements:
      | level      | expected_imports_removed | expected_vars_removed | quality_improvement |
      | safe       | >= 1                    | 0                     | 5-10%              |
      | standard   | >= 2                    | >= 1                  | 10-25%             |
      | aggressive | >= 3                    | >= 2                  | 25-50%             |
    And architectural integrity should be maintained
    And framework-specific patterns should be preserved

    Examples:
      | framework | database | architecture | optimization_level |
      | gin       | postgres | hexagonal    | safe              |
      | gin       | postgres | hexagonal    | standard          |
      | gin       | postgres | hexagonal    | aggressive        |
      | echo      | mysql    | clean        | safe              |
      | echo      | mysql    | clean        | standard          |
      | fiber     | sqlite   | ddd          | safe              |
      | fiber     | sqlite   | ddd          | standard          |
      | chi       | postgres | standard     | safe              |

  @integration @optimization @matrix @performance
  Scenario: Performance impact measurement across matrix combinations
    Given I have baseline performance metrics for matrix combinations:
      | framework | database | architecture | baseline_compile_time | baseline_file_count |
      | gin       | postgres | hexagonal    | 5.2s                 | 45                  |
      | echo      | mysql    | clean        | 4.8s                 | 38                  |
      | fiber     | sqlite   | ddd          | 5.5s                 | 52                  |
    When I apply optimization and measure performance impact
    Then compilation time should improve by 10-30%
    And file count should decrease through code consolidation
    And quality metrics should show measurable improvement:
      | metric                 | improvement_range |
      | cyclomatic_complexity  | 5-15%            |
      | maintainability_index  | 10-25%           |
      | code_duplication      | 15-40%           |
    And no functional regressions should occur

  @integration @optimization @matrix @quality
  Scenario: Quality-optimization feedback loop
    Given I have quality analysis results for matrix-generated projects
    When optimization is applied based on quality findings
    Then specific quality issues should be resolved:
      | issue_type           | optimization_solution    | success_criteria        |
      | unused_imports       | safe level cleanup       | 100% import utilization |
      | unused_variables     | standard level cleanup   | 95% variable utilization|
      | complex_functions    | aggressive refactoring   | Reduced cyclomatic complexity |
      | duplicate_code       | expert level consolidation| <5% code duplication   |
    And quality scores should improve for all frameworks
    And architectural patterns should remain intact

  @integration @optimization @matrix @cross-validation
  Scenario: Cross-system validation and regression testing
    Given I have working matrix combinations with optimization applied
    When I run comprehensive validation across all systems
    Then no system should break another system's functionality
    And optimization should work consistently across:
      | system_aspect        | validation_criteria                    |
      | Framework patterns   | Framework-specific code preserved      |
      | Database connections | All database drivers functional        |
      | Architecture layers  | Architectural boundaries maintained    |
      | Authentication      | Auth systems work after optimization   |
      | Logger integration  | Logging works with optimized code      |
    And regression tests should pass for all combinations

  @integration @optimization @matrix @advanced
  Scenario: Advanced optimization-matrix synergies
    Given I have complex matrix configurations with multiple features
    When I apply context-aware optimization strategies:
      | context        | strategy              | focus_areas                |
      | Development    | conservative         | Import cleanup, safe vars  |
      | Testing        | balanced            | Testing patterns preserved |
      | Production     | performance         | Aggressive optimization    |
      | Maintenance    | maximum             | Code consolidation         |
    Then optimization should adapt to project context
    And matrix combinations should benefit from specialized optimization
    And performance gains should compound across optimizations

  @integration @optimization @matrix @metrics
  Scenario: Comprehensive metrics collection for optimization-matrix integration
    Given I apply optimization to a full matrix of project combinations
    When I collect detailed metrics across all systems
    Then I should have comprehensive data on:
      | metric_category      | specific_metrics                           |
      | Performance Impact   | Compile time, runtime performance, memory  |
      | Code Quality         | Complexity, maintainability, duplication  |
      | Optimization Results | Files modified, imports removed, vars cleaned |
      | System Integration   | Cross-system compatibility, regression tests |
      | User Experience     | Generation time, error rates, warnings    |
    And metrics should enable data-driven optimization decisions
    And trends should be trackable over time

  @integration @optimization @matrix @automation
  Scenario: Automated optimization recommendations for matrix projects
    Given I have historical data on optimization effectiveness by matrix combination
    When a new project is generated with specific matrix settings
    Then the system should automatically recommend:
      | recommendation_type  | criteria                                    |
      | Optimization Level   | Based on project complexity and context    |
      | Custom Profile       | Tailored to framework/architecture combo   |
      | Risk Assessment      | Warning for potentially risky combinations |
      | Performance Estimate| Expected improvement based on similar projects |
    And recommendations should improve over time with more data
    And users should be able to accept/reject recommendations

  @integration @optimization @matrix @error-handling
  Scenario: Robust error handling in optimization-matrix integration
    Given I have matrix projects with various edge cases and potential issues
    When optimization encounters problems during matrix integration
    Then the system should handle errors gracefully:
      | error_type              | expected_behavior                    |
      | Compilation failure     | Rollback optimization, preserve original |
      | Optimization conflicts  | Warning + safe fallback strategy     |
      | Resource exhaustion     | Graceful degradation with status    |
      | Template inconsistency  | Clear error + suggested resolution  |
    And error recovery should maintain system integrity
    And detailed error reports should assist troubleshooting

  @integration @optimization @matrix @scalability
  Scenario: Scalability validation for large matrix sets
    Given I have a large matrix of project combinations (50+ combinations)
    When I apply batch optimization across all combinations
    Then the system should scale efficiently:
      | scalability_aspect   | requirement                          |
      | Processing time      | Linear scaling with project count   |
      | Memory usage         | Bounded memory growth                |
      | Parallel processing  | Effective utilization of CPU cores  |
      | Error isolation      | Failures don't cascade to other projects |
      | Progress reporting   | Real-time status across batch       |
    And batch processing should be more efficient than individual runs
    And resource usage should remain within acceptable bounds