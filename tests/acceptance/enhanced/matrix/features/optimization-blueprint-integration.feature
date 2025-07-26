Feature: Optimization-Blueprint Integration Matrix
  As a Go developer using go-starter's comprehensive enhanced system
  I want optimization to work seamlessly across all blueprint types with measurable improvements
  So that I can generate high-quality, optimized projects regardless of complexity or architecture

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And the optimization system is available
    And the matrix testing system is available
    And all blueprint types are operational

  @enhanced @optimization @blueprint-matrix @critical
  Scenario: Comprehensive optimization across all blueprint types
    Given I have the complete blueprint matrix with optimization levels:
      | blueprint_type     | architecture | framework | database | auth_type | optimization_level |
      | web-api           | standard     | gin       | postgres | jwt       | safe              |
      | web-api           | clean        | echo      | mysql    | oauth2    | standard          |
      | web-api           | ddd          | fiber     | sqlite   | session   | aggressive        |
      | web-api           | hexagonal    | chi       | mongodb  | api-key   | expert            |
      | cli               | simple       | cobra     | -        | -         | safe              |
      | cli               | standard     | cobra     | -        | -         | standard          |
      | library           | standard     | -         | -        | -         | aggressive        |
      | lambda            | standard     | gin       | -        | jwt       | safe              |
      | lambda-proxy      | standard     | gin       | postgres | jwt       | standard          |
      | microservice      | clean        | echo      | postgres | oauth2    | aggressive        |
      | monolith          | ddd          | gin       | mysql    | session   | expert            |
      | workspace         | multi        | mixed     | postgres | mixed     | standard          |
    When I generate and optimize each blueprint configuration
    Then all projects should compile successfully after optimization
    And optimization impact should be measurable for each blueprint type
    And quality metrics should show improvement across all configurations
    And no functional regressions should occur

  @enhanced @optimization @blueprint-specific @quality
  Scenario: Blueprint-specific optimization effectiveness validation
    Given I want to measure optimization effectiveness per blueprint type
    When I apply different optimization levels to blueprint categories:
      | blueprint_category | complexity_level | baseline_metrics           | optimization_targets        |
      | Simple APIs        | low             | 15 files, 800 LOC         | 10% import reduction       |
      | Enterprise APIs    | high            | 45 files, 2500 LOC        | 25% code quality improvement|
      | CLI Tools          | medium          | 12 files, 600 LOC         | 15% unused code removal    |
      | Lambda Functions   | low             | 8 files, 400 LOC          | 20% performance improvement |
      | Microservices      | high            | 35 files, 2000 LOC        | 30% architectural cleanup  |
      | Monoliths          | very_high       | 60 files, 4000 LOC        | 35% complexity reduction   |
      | Libraries          | medium          | 10 files, 500 LOC         | 40% API clarity improvement|
      | Workspaces         | very_high       | 100+ files, 8000+ LOC     | 25% cross-module optimization|
    Then optimization effectiveness should meet category-specific targets
    And improvements should be sustainable across project lifecycle
    And optimization time should scale appropriately with project complexity

  @enhanced @optimization @architecture-aware @patterns
  Scenario: Architecture-aware optimization with pattern preservation
    Given I have projects with different architectural patterns requiring specific optimizations
    When I apply architecture-aware optimization to projects:
      | architecture | optimization_focus              | pattern_preservation        | expected_improvements       |
      | standard     | general cleanup                | loose coupling             | 15-25% code quality        |
      | clean        | layer boundary enforcement     | dependency direction       | 20-30% architectural clarity|
      | ddd          | domain model optimization      | bounded context integrity  | 25-35% domain clarity      |
      | hexagonal    | port-adapter optimization      | interface segregation      | 30-40% testability        |
      | event-driven | event flow optimization        | event sourcing patterns    | 35-45% event consistency   |
    Then architectural patterns should be preserved and enhanced
    And layer boundaries should be respected during optimization
    And domain-specific optimizations should improve architectural adherence
    And optimization should not violate architectural principles

  @enhanced @optimization @framework-integration @compatibility
  Scenario: Framework-specific optimization with deep integration
    Given I have projects using different frameworks with optimization applied
    When I optimize framework-specific patterns and integrations:
      | framework | optimization_areas                    | framework_patterns           | integration_depth          |
      | gin       | middleware chains, route organization | gin.Context usage           | deep handler optimization  |
      | echo      | context handling, error patterns     | echo.Context patterns       | middleware consolidation   |
      | fiber     | performance paths, memory allocation  | fiber.Ctx efficiency       | resource optimization      |
      | chi       | routing complexity, middleware stack  | chi.Router patterns         | route organization         |
      | cobra     | command structure, flag management    | cobra.Command hierarchy     | CLI pattern optimization   |
    Then framework-specific optimizations should preserve framework idioms
    And optimization should leverage framework strengths
    And framework integration should be enhanced, not compromised
    And performance improvements should be framework-appropriate

  @enhanced @optimization @database-integration @data-patterns
  Scenario: Database integration optimization with ORM awareness
    Given I have projects with various database and ORM combinations
    When I apply database-aware optimization to projects:
      | database | orm  | optimization_focus           | data_patterns            | expected_improvements      |
      | postgres | gorm | model relationships, queries | association patterns     | 20-30% query efficiency   |
      | postgres | sqlx | raw query optimization       | transaction patterns     | 25-35% performance gain   |
      | mysql    | gorm | connection management        | connection pooling       | 15-25% resource efficiency|
      | sqlite   | sqlx | embedded optimization        | file operation patterns  | 30-40% embedded performance|
      | mongodb  | -    | document structure, indexing | NoSQL patterns          | 25-35% document efficiency |
      | redis    | -    | caching patterns, key design | cache optimization       | 40-50% cache effectiveness |
    Then database integration patterns should be optimized appropriately
    And ORM usage should be enhanced with best practices
    And database-specific optimizations should improve data access patterns
    And connection management should be optimized per database type

  @enhanced @optimization @auth-integration @security-patterns
  Scenario: Authentication system optimization with security preservation
    Given I have projects with different authentication patterns
    When I optimize authentication and security implementations:
      | auth_type | security_patterns              | optimization_areas         | security_preservation      |
      | jwt       | token validation, claims       | token handling efficiency  | security boundary integrity|
      | oauth2    | provider integration, flows    | OAuth flow optimization    | security protocol adherence|
      | session   | session management, storage    | session lifecycle          | secure session handling    |
      | api-key   | key validation, rate limiting  | key management efficiency  | API security patterns      |
      | basic     | credential handling            | basic auth optimization    | credential security        |
    Then authentication patterns should be optimized without security compromise
    And security boundaries should be maintained and enhanced
    And optimization should improve authentication performance
    And security best practices should be enforced during optimization

  @enhanced @optimization @complexity-scaling @performance
  Scenario: Optimization performance scaling across project complexity levels
    Given I want to validate optimization performance across complexity scales
    When I apply optimization to projects of varying complexity:
      | complexity_level | project_characteristics        | optimization_time_target | memory_usage_target | success_criteria           |
      | simple          | 5-15 files, <1000 LOC         | <5 seconds              | <50MB              | 100% success, minimal impact|
      | moderate        | 15-30 files, 1000-3000 LOC    | 5-15 seconds            | 50-100MB           | 100% success, clear benefit |
      | complex         | 30-60 files, 3000-6000 LOC    | 15-45 seconds           | 100-200MB          | 100% success, significant benefit|
      | enterprise      | 60+ files, 6000+ LOC          | 45-120 seconds          | 200-400MB          | 100% success, major benefit |
    Then optimization should scale efficiently with project complexity
    And processing time should grow sub-linearly with project size
    And memory usage should be bounded and predictable
    And optimization quality should improve with larger codebases

  @enhanced @optimization @regression-prevention @stability
  Scenario: Comprehensive regression prevention across optimization levels
    Given I have established baselines for all blueprint-optimization combinations
    When I apply optimization and validate against regression benchmarks:
      | regression_category    | monitoring_areas              | detection_thresholds      | prevention_measures        |
      | functional_regression | compilation, test execution   | 0% tolerance              | automated rollback        |
      | performance_regression| build time, runtime performance| >10% degradation         | performance gates         |
      | quality_regression    | code metrics, complexity      | metric-specific thresholds| quality gates             |
      | security_regression   | security patterns, validation | 0% tolerance              | security validation       |
      | architectural_regression| pattern compliance, boundaries| pattern-specific rules   | architectural validation   |
    Then no regressions should be detected across any optimization level
    And regression detection should be automated and reliable
    And prevention measures should automatically activate when thresholds are exceeded
    And rollback mechanisms should restore project state when regressions occur

  @enhanced @optimization @metrics-collection @analytics
  Scenario: Comprehensive optimization metrics collection and analysis
    Given I want to collect detailed metrics on optimization effectiveness
    When I optimize projects and collect comprehensive analytics:
      | metric_category      | specific_metrics                    | collection_granularity | analysis_dimensions        |
      | Code Quality        | cyclomatic complexity, maintainability| per-file, per-function  | before/after, trend analysis|
      | Performance Impact  | processing time, memory usage       | per-operation, overall  | scaling analysis, efficiency|
      | Optimization Coverage| files processed, changes made        | per-blueprint, per-level| coverage completeness      |
      | Blueprint Effectiveness| success rate, improvement percentage| per-type, per-config   | comparative effectiveness  |
      | User Experience    | ease of use, configuration clarity  | per-interaction        | usability improvement      |
      | System Integration | cross-system compatibility          | per-integration-point  | integration health         |
    Then metrics should provide actionable insights for optimization improvement
    And analytics should identify optimization patterns and effectiveness trends
    And data should support evidence-based optimization strategy decisions
    And metrics collection should not impact optimization performance

  @enhanced @optimization @custom-profiles @extensibility
  Scenario: Custom optimization profiles for specialized use cases
    Given I want to create and validate custom optimization profiles
    When I define specialized optimization profiles for specific scenarios:
      | profile_name     | target_use_case          | optimization_rules            | custom_parameters           |
      | startup_optimized| fast startup applications| aggressive import cleanup     | startup_time_priority=high  |
      | memory_constrained| resource-limited environments| memory-focused optimization | memory_limit=128MB         |
      | ci_cd_optimized  | continuous integration   | build-time optimization       | build_speed_priority=high   |
      | security_hardened| security-critical apps   | security-first optimization   | security_level=paranoid     |
      | performance_first| high-performance systems | aggressive performance opts   | performance_target=p99<1ms  |
      | maintainability  | long-term maintenance    | readability-focused cleanup   | maintainability_score>8.5  |
    Then custom profiles should be validated across relevant blueprint types
    And profile-specific optimizations should achieve targeted improvements
    And custom parameters should influence optimization behavior appropriately
    And profiles should be composable and extensible for complex scenarios

  @enhanced @optimization @blueprint-evolution @future-proofing
  Scenario: Optimization compatibility with blueprint evolution and versioning
    Given I have blueprints that evolve over time with optimization applied
    When I validate optimization compatibility across blueprint versions:
      | evolution_scenario    | version_changes                | optimization_adaptation       | compatibility_requirements |
      | blueprint_updates    | new features, improved patterns| optimization rule updates     | backward compatibility     |
      | framework_upgrades   | framework version changes      | framework-specific adaptations| version-aware optimization |
      | architecture_migration| pattern transitions           | migration-aware optimization  | smooth transition support  |
      | dependency_updates   | library and tool updates       | dependency-aware optimization | dependency compatibility   |
      | platform_evolution   | new platform support          | platform-specific optimization| cross-platform consistency|
    Then optimization should adapt gracefully to blueprint evolution
    And optimization rules should be version-aware and adaptive
    And compatibility should be maintained across blueprint generations
    And optimization should support blueprint migration scenarios

  @enhanced @optimization @integration-testing @ecosystem
  Scenario: Deep integration testing across the complete enhanced ecosystem
    Given I have the complete enhanced ecosystem with optimization integration
    When I validate deep integration across all enhanced components:
      | integration_area     | component_interactions           | validation_requirements      | success_criteria           |
      | Matrix-Optimization | matrix generation + optimization | seamless workflow           | 100% integration success   |
      | Quality-Optimization| quality analysis + optimization  | feedback-driven improvement | measurable quality gains   |
      | Performance-Optimization| performance monitoring + opts  | performance-guided optimization| performance improvement    |
      | Validation-Optimization| cross-system validation + opts | comprehensive validation    | no integration conflicts   |
      | CLI-Optimization    | CLI workflows + optimization     | user experience continuity  | intuitive optimization UX  |
      | Template-Optimization| template engine + optimization  | template-aware optimization | template integrity         |
    Then all enhanced components should work together seamlessly
    And optimization should enhance rather than conflict with other systems
    And integration workflows should be intuitive and reliable
    And the complete ecosystem should provide synergistic benefits

  @enhanced @optimization @user-experience @workflow
  Scenario: Optimization user experience across different developer workflows
    Given I have developers with different workflows and optimization needs
    When I validate optimization user experience across workflow scenarios:
      | workflow_type        | developer_profile          | optimization_needs           | experience_requirements    |
      | rapid_prototyping   | startup developer          | fast, safe optimization      | minimal configuration      |
      | enterprise_development| enterprise team lead      | comprehensive, auditable     | detailed reporting         |
      | open_source_maintenance| OSS maintainer           | community-friendly, transparent| clear impact communication |
      | performance_engineering| performance engineer      | metrics-driven, detailed     | comprehensive analytics    |
      | security_focused    | security engineer          | security-aware, cautious     | security validation        |
      | educational_use     | student/learner            | educational, explanatory     | learning-oriented output   |
    Then optimization should adapt to different developer needs and contexts
    And user experience should be appropriate for each workflow type
    And optimization should provide value without overwhelming complexity
    And different profiles should have tailored optimization experiences