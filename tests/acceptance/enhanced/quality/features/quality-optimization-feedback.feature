Feature: Quality-Optimization Feedback Loop
  As a Go developer using go-starter's enhanced system
  I want optimization to continuously improve based on quality analysis feedback
  So that code quality gets better with each optimization iteration

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And the optimization system is available
    And the quality analysis system is available

  @quality @optimization @feedback
  Scenario: Basic quality-optimization feedback loop
    Given I generate a project with quality issues:
      | issue_type          | severity | file_pattern    | description                    |
      | unused_imports      | high     | *.go           | Multiple unused import statements |
      | unused_variables    | medium   | handlers/*.go  | Unused variable declarations   |
      | complex_functions   | low      | services/*.go  | Functions with high complexity |
    When I analyze the project for quality issues
    Then I should identify all quality problems accurately
    When I apply optimization based on quality feedback
    Then optimization should target identified quality issues:
      | optimization_area | expected_improvements           | success_criteria          |
      | Import cleanup   | Remove 90-100% unused imports  | Zero unused imports       |
      | Variable cleanup | Remove 80-95% unused variables | <5% unused variables      |
      | Function refactor| Reduce complexity by 20-40%    | Functions under 15 CC     |
    And the project should compile successfully after optimization
    And quality metrics should show measurable improvement

  @quality @optimization @iterative
  Scenario Outline: Iterative quality improvement cycles
    Given I have a "<complexity_level>" project with "<architecture>" architecture
    When I run quality analysis to establish baseline metrics
    Then baseline quality should be recorded for:
      | metric_category      | measurement_criteria              |
      | Code Complexity      | Cyclomatic complexity per function|
      | Import Efficiency    | Unused vs total imports ratio     |
      | Variable Usage       | Unused vs declared variables      |
      | Function Cohesion    | Function length and responsibility|
      | Architecture Adherence| Layer separation violations      |
    When I apply optimization iteration "<iteration>" with focus "<focus_area>"
    Then quality should improve incrementally:
      | iteration | complexity_improvement | import_improvement | variable_improvement |
      | 1         | 10-25%                | 60-80%            | 40-60%              |
      | 2         | 25-45%                | 80-95%            | 60-85%              |
      | 3         | 45-70%                | 95-100%           | 85-100%             |
    And each iteration should maintain architectural integrity
    And no regressions should occur in previously optimized areas

    Examples:
      | complexity_level | architecture | iteration | focus_area        |
      | moderate         | standard     | 1         | imports          |
      | moderate         | standard     | 2         | variables        |
      | moderate         | standard     | 3         | functions        |
      | high             | clean        | 1         | architecture     |
      | high             | clean        | 2         | complexity       |
      | very_high        | hexagonal    | 1         | dependencies     |

  @quality @optimization @adaptive
  Scenario: Adaptive optimization based on quality patterns
    Given I have projects with different quality patterns:
      | project_type | dominant_issues              | optimization_priority |
      | web-api      | import bloat, unused vars   | cleanup-focused      |
      | cli-tool     | complex functions, duplication| refactoring-focused  |
      | microservice | coupling issues, dependencies | architecture-focused |
      | library      | API consistency, documentation| interface-focused    |
    When quality analysis identifies dominant quality patterns
    Then optimization strategy should adapt to quality findings:
      | issue_pattern        | adaptive_strategy              | expected_focus_areas     |
      | High import waste    | Aggressive import cleanup      | Unused imports, wildcards|
      | Function complexity  | Function decomposition priority| Cyclomatic complexity    |
      | Architecture violations| Structure-first optimization   | Layer boundaries        |
      | Code duplication     | Consolidation and extraction   | Common patterns         |
    And optimization parameters should adjust based on pattern severity
    And feedback should influence future optimization decisions

  @quality @optimization @framework-specific
  Scenario: Framework-specific quality-optimization synergy
    Given I generate projects with different framework quality characteristics:
      | framework | typical_quality_issues           | optimization_opportunities |
      | gin       | middleware complexity, routing   | Handler consolidation     |
      | echo      | context handling, error patterns | Error flow optimization   |
      | fiber     | performance patterns, memory     | Resource optimization     |
      | chi       | routing complexity, middleware   | Route organization        |
    When I apply framework-aware quality optimization
    Then optimization should address framework-specific quality issues:
      | framework | optimization_areas                    | quality_improvements      |
      | gin       | Middleware chain, handler organization| Reduced routing complexity|
      | echo      | Context usage, error handling patterns| Improved error flow      |
      | fiber     | Memory allocation, performance paths  | Resource efficiency      |
      | chi       | Route definitions, middleware stack   | Cleaner routing structure|
    And framework best practices should be automatically enforced
    And framework-specific quality patterns should improve over time

  @quality @optimization @database-integration
  Scenario: Database-related quality optimization feedback
    Given I have projects with database integration quality issues:
      | database  | orm   | quality_issues                        | optimization_targets     |
      | postgres  | gorm  | connection leaks, query complexity   | Connection management    |
      | postgres  | sqlx  | error handling, transaction patterns | Error flow optimization  |
      | mysql     | gorm  | model complexity, relationship issues| Model simplification     |
      | mongodb   | -     | document structure, query efficiency | Query optimization       |
    When quality analysis examines database integration patterns
    Then database-specific optimization should be applied:
      | optimization_area     | quality_improvements                 | success_metrics          |
      | Connection Management | Proper connection pooling patterns  | Zero connection leaks    |
      | Query Optimization    | Efficient query patterns            | <100ms query times       |
      | Model Organization    | Clean model definitions             | Reduced model complexity |
      | Transaction Handling  | Proper transaction scoping         | No hanging transactions  |
    And database integration quality should be continuously monitored
    And optimization should prevent common database antipatterns

  @quality @optimization @cross-architecture
  Scenario: Cross-architecture quality optimization consistency
    Given I have projects using different architectures with similar quality issues:
      | architecture | layer_violations | dependency_issues | complexity_hotspots |
      | standard     | 5               | 8                | 12                 |
      | clean        | 12              | 15               | 18                 |
      | hexagonal    | 18              | 22               | 25                 |
      | ddd          | 25              | 30               | 35                 |
    When I apply architecture-aware quality optimization
    Then optimization should respect architectural boundaries:
      | architecture | layer_integrity | dependency_direction | separation_clarity |
      | standard     | maintained      | correct             | improved          |
      | clean        | enforced        | unidirectional      | clear boundaries  |
      | hexagonal    | strict ports    | inward only         | explicit contracts|
      | ddd          | domain focus    | domain-driven       | bounded contexts  |
    And quality improvements should align with architectural principles
    And architectural violations should be detected and corrected

  @quality @optimization @performance-quality
  Scenario: Performance-quality optimization balance
    Given I have projects where performance and quality optimizations may conflict
    When quality optimization suggests changes that might impact performance
    Then the system should provide balanced optimization recommendations:
      | scenario                  | quality_benefit | performance_impact | recommendation        |
      | Extract complex function  | High           | Minimal           | Proceed with extraction|
      | Inline simple function    | Low            | Moderate gain     | Keep performance      |
      | Add error handling        | High           | Small overhead    | Proceed with safety   |
      | Remove debug logging      | Low            | Small gain        | Context-dependent     |
    And users should be informed of trade-offs
    And optimization should allow user preference configuration
    And performance regressions should be prevented

  @quality @optimization @continuous-improvement
  Scenario: Continuous quality improvement over multiple projects
    Given I have historical quality data from previous optimization runs
    When I generate and optimize new projects
    Then the system should learn from quality patterns:
      | learning_area          | improvement_mechanism               | expected_benefits        |
      | Common issue detection | Pattern recognition from history   | Faster issue identification|
      | Optimization effectiveness| Success rate tracking            | Better optimization choices|
      | Framework best practices| Usage pattern analysis            | Framework-specific advice |
      | Architecture guidelines | Violation pattern learning        | Preventive recommendations|
    And quality recommendations should improve over time
    And optimization strategies should become more targeted
    And false positive rates should decrease with more data

  @quality @optimization @validation
  Scenario: Quality optimization validation and rollback
    Given I have a project with mixed quality issues and working functionality
    When optimization is applied with quality feedback
    Then optimization validation should ensure:
      | validation_aspect      | validation_criteria                    | rollback_triggers        |
      | Functional integrity   | All tests pass after optimization     | Test failures           |
      | Compilation success    | Project compiles without errors       | Compilation errors      |
      | Performance baseline   | No significant performance degradation| >20% performance loss   |
      | Quality improvements   | Measurable quality gains achieved     | No quality improvement  |
      | Architecture compliance| Architectural patterns maintained     | Architecture violations |
    And rollback should be automatic when validation fails
    And rollback should preserve original project state
    And users should receive detailed validation reports

  @quality @optimization @metrics-collection
  Scenario: Comprehensive quality-optimization metrics collection
    Given I optimize multiple projects with quality feedback enabled
    When I collect comprehensive quality-optimization metrics
    Then I should have detailed data on:
      | metric_category        | specific_metrics                        | collection_frequency |
      | Quality Improvements   | Before/after quality scores            | Per optimization run |
      | Optimization Effectiveness| Success rates by issue type          | Per project type     |
      | Time Performance       | Optimization processing times          | Real-time           |
      | User Satisfaction      | Developer feedback on quality gains    | Post-optimization   |
      | False Positive Rates   | Incorrect quality issue identification | Continuous          |
      | Architectural Impact   | Architecture adherence changes         | Pre/post analysis   |
    And metrics should enable optimization strategy improvements
    And quality trends should be trackable over time
    And metrics should support evidence-based optimization decisions

  @quality @optimization @integration-testing
  Scenario: Quality-optimization integration with existing workflows
    Given I have existing development workflows with quality gates
    When quality-optimization feedback is integrated into the workflow
    Then integration should work seamlessly with:
      | integration_point      | expected_behavior                     | success_criteria         |
      | CI/CD pipelines       | Automatic quality optimization        | No pipeline failures     |
      | Code review process   | Quality improvement suggestions       | Enhanced review insights |
      | IDE integrations      | Real-time quality feedback           | In-editor recommendations|
      | Git hooks             | Pre-commit quality optimization      | Improved commit quality  |
      | Build systems         | Quality-aware build optimization     | Faster, cleaner builds   |
    And existing tools should not be disrupted
    And quality improvements should be visible in existing metrics
    And developer workflow efficiency should improve

  @quality @optimization @error-handling
  Scenario: Robust error handling in quality-optimization feedback
    Given I have projects with various edge cases and potential issues
    When quality-optimization encounters problems during analysis or optimization
    Then error handling should be robust:
      | error_type                | expected_behavior                  | recovery_action          |
      | Parse errors in Go files  | Skip problematic files, continue  | Report parsing issues    |
      | Optimization conflicts    | Provide safe fallback strategies  | Manual resolution prompt |
      | Quality analysis failures | Use degraded analysis mode        | Limited optimization     |
      | Resource exhaustion       | Graceful degradation              | Partial optimization     |
      | Corrupted project state   | Automatic rollback to safe state  | Full state restoration   |
    And error recovery should maintain project integrity
    And detailed error reports should assist troubleshooting
    And users should receive clear guidance on resolution steps