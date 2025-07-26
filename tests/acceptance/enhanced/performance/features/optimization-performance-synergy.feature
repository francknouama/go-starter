Feature: Optimization-Performance Synergy
  As a Go developer using go-starter's optimization system
  I want performance improvements to be measurable and predictable across different blueprint combinations
  So that I can make informed decisions about optimization strategies

  Background:
    Given I am using go-starter CLI
    And the optimization system is available
    And the performance monitoring system is available
    And all templates are properly initialized

  @performance @optimization @synergy
  Scenario: Measure optimization impact across blueprint architectures
    Given I have baseline projects for performance comparison:
      | blueprint | architecture | framework | baseline_compile_time | baseline_size_mb | baseline_complexity |
      | web-api   | standard     | gin       | 3.2s                 | 2.1             | 45                  |
      | web-api   | clean        | echo      | 4.1s                 | 2.8             | 52                  |
      | web-api   | hexagonal    | fiber     | 4.8s                 | 3.2             | 61                  |
      | web-api   | ddd          | chi       | 5.2s                 | 3.6             | 68                  |
    When I apply "standard" optimization to each architecture
    Then performance improvements should be measurable:
      | metric              | standard_improvement | clean_improvement | hexagonal_improvement | ddd_improvement |
      | Compile time        | 15-25%              | 20-30%           | 25-35%               | 30-40%         |
      | Binary size         | 5-10%               | 8-15%            | 12-20%               | 15-25%         |
      | Cyclomatic complexity| 10-20%              | 15-25%           | 20-30%               | 25-35%         |
    And architectural integrity should be maintained for all blueprints

  @performance @optimization @levels
  Scenario Outline: Performance scaling across optimization levels
    Given I have a "<architecture>" project with "<complexity_class>" code patterns
    When I apply optimization at different levels:
      | level      | expected_time_impact | expected_quality_gain | risk_level |
      | safe       | 0-5% improvement     | 5-15% gain           | minimal    |
      | standard   | 10-25% improvement   | 15-35% gain          | low        |
      | aggressive | 25-50% improvement   | 35-60% gain          | medium     |
      | expert     | 50-80% improvement   | 60-90% gain          | high       |
    Then performance should scale predictably with optimization intensity
    And quality gains should justify processing overhead
    And risk warnings should be provided for aggressive levels

    Examples:
      | architecture | complexity_class |
      | standard     | low              |
      | clean        | medium           |
      | hexagonal    | high             |
      | ddd          | very_high        |

  @performance @optimization @framework-specific
  Scenario: Framework-specific optimization performance characteristics
    Given I have projects using different frameworks with optimization applied:
      | framework | typical_performance | optimization_potential | framework_overhead |
      | gin       | high               | medium                | low               |
      | echo      | very_high          | low                   | minimal           |
      | fiber     | highest            | very_low              | minimal           |
      | chi       | medium             | high                  | medium            |
    When I measure framework-specific optimization benefits
    Then each framework should show characteristic performance patterns:
      | framework | compile_improvement | runtime_improvement | binary_size_impact |
      | gin       | 15-25%             | 5-10%              | 8-12%             |
      | echo      | 8-15%              | 3-7%               | 5-8%              |
      | fiber     | 5-10%              | 2-5%               | 3-5%              |
      | chi       | 20-30%             | 8-15%              | 10-18%            |
    And framework-specific optimizations should preserve middleware functionality

  @performance @optimization @database-impact
  Scenario: Database integration optimization performance analysis
    Given I have projects with different database configurations:
      | database  | orm   | connection_complexity | optimization_opportunities |
      | postgres  | gorm  | medium               | high                      |
      | postgres  | sqlx  | low                  | medium                    |
      | mysql     | gorm  | medium               | high                      |
      | mysql     | sqlx  | low                  | medium                    |
      | sqlite    | gorm  | low                  | medium                    |
      | mongodb   | -     | high                 | very_high                 |
    When I apply optimization focusing on database-related code
    Then database performance should improve measurably:
      | optimization_area        | expected_improvement |
      | Connection pool setup    | 20-40%              |
      | Query preparation time   | 15-30%              |
      | ORM initialization       | 25-50%              |
      | Migration script loading | 30-60%              |
    And database functionality should remain intact
    And connection stability should be maintained

  @performance @optimization @concurrent-processing
  Scenario: Concurrent optimization performance validation
    Given I have multiple projects for concurrent optimization:
      | project_size | file_count | complexity | expected_processing_time |
      | small        | 10-20      | low        | 2-5 seconds             |
      | medium       | 30-60      | medium     | 8-15 seconds            |
      | large        | 80-150     | high       | 20-45 seconds           |
      | very_large   | 200+       | very_high  | 45-90 seconds           |
    When I optimize projects concurrently with different worker counts:
      | workers | expected_speedup | memory_overhead | cpu_utilization |
      | 1       | 1x (baseline)    | minimal        | 25-40%         |
      | 2       | 1.6-1.8x        | low            | 50-70%         |
      | 4       | 2.8-3.2x        | medium         | 80-95%         |
      | 8       | 3.5-4.2x        | high           | 95-100%        |
    Then concurrent processing should scale efficiently
    And memory usage should remain bounded
    And no race conditions should occur between optimizations

  @performance @optimization @caching-effectiveness
  Scenario: AST parsing and analysis caching performance gains
    Given I have projects with similar code patterns for caching analysis
    When I optimize similar projects in sequence
    Then caching should provide measurable performance benefits:
      | cache_type          | first_run | subsequent_runs | cache_hit_ratio |
      | AST parsing cache   | baseline  | 40-60% faster  | 70-85%         |
      | Pattern analysis    | baseline  | 30-50% faster  | 60-80%         |
      | Import resolution   | baseline  | 50-70% faster  | 80-90%         |
      | Template analysis   | baseline  | 35-55% faster  | 65-75%         |
    And cache effectiveness should improve over time
    And memory usage for caching should be optimized
    And cache invalidation should work correctly when code changes

  @performance @optimization @memory-efficiency
  Scenario: Memory usage optimization during performance-intensive operations
    Given I enable memory profiling for optimization operations
    When I optimize projects of varying sizes:
      | project_size | files | lines_of_code | max_memory_mb | gc_frequency |
      | tiny         | 5     | 500          | 50           | minimal     |
      | small        | 15    | 2000         | 120          | low         |
      | medium       | 50    | 8000         | 300          | medium      |
      | large        | 150   | 25000        | 800          | high        |
      | huge         | 500   | 100000       | 2000         | frequent    |
    Then memory usage should scale predictably with project size
    And garbage collection should be effective
    And memory leaks should not occur during long-running optimizations
    And peak memory usage should not exceed reasonable thresholds

  @performance @optimization @regression-detection
  Scenario: Performance regression detection and alerting
    Given I have historical performance baselines for optimization operations
    When I run optimization with the current implementation
    Then performance should not regress beyond acceptable thresholds:
      | performance_metric    | acceptable_regression | alert_threshold |
      | Processing speed      | <15% slower          | 10% slower     |
      | Memory consumption    | <25% increase        | 20% increase   |
      | Cache hit ratio       | <10% decrease        | 5% decrease    |
      | Compilation time      | <20% slower          | 15% slower     |
    And any regressions should trigger automated alerts
    And regression analysis should provide actionable insights
    And performance trends should be tracked over time

  @performance @optimization @quality-vs-speed
  Scenario: Optimize quality vs processing speed trade-offs
    Given I have projects suitable for quality-speed trade-off analysis
    When I configure different optimization priorities:
      | priority_mode | quality_weight | speed_weight | expected_behavior           |
      | quality_first | 80%           | 20%         | thorough, slower processing |
      | balanced      | 50%           | 50%         | good balance                |
      | speed_first   | 20%           | 80%         | fast, basic optimization    |
    Then each mode should deliver predictable trade-offs:
      | mode          | processing_time | quality_improvement | user_satisfaction |
      | quality_first | 150-200%       | 80-95%             | high              |
      | balanced      | 100-130%       | 60-80%             | very_high         |
      | speed_first   | 60-80%         | 30-50%             | medium            |
    And users should be able to customize trade-off preferences
    And recommendations should adapt to project characteristics

  @performance @optimization @real-time-monitoring
  Scenario: Real-time performance monitoring during optimization
    Given I enable comprehensive real-time monitoring
    When I optimize a large project with detailed progress tracking
    Then real-time metrics should be available:
      | metric_category        | update_frequency | data_points                    |
      | Processing progress    | every 10 files   | files completed, remaining     |
      | Performance metrics    | every 30 seconds | speed, memory, CPU usage       |
      | Quality improvements   | per phase        | imports, variables, functions  |
      | Resource utilization   | continuous       | memory, CPU, disk I/O         |
      | Bottleneck detection   | real-time        | slow operations, blocking calls|
    And performance data should be exportable for analysis
    And alerts should trigger for performance anomalies
    And users should be able to monitor optimization progress effectively

  @performance @optimization @benchmarking-suite
  Scenario: Comprehensive optimization benchmarking suite
    Given I have a standardized benchmarking suite for optimization performance
    When I run the complete benchmark across all blueprint combinations
    Then I should get comprehensive performance profiles:
      | benchmark_category     | measured_aspects                           |
      | Blueprint Generation   | time, memory, complexity                   |
      | Optimization Pipeline  | processing speed, quality impact           |
      | Architecture Patterns  | optimization effectiveness per pattern     |
      | Framework Integration  | framework-specific performance             |
      | Database Optimization  | database-related improvement metrics       |
      | Concurrent Processing  | scalability, resource utilization         |
      | Caching Systems       | cache effectiveness, hit ratios            |
    And benchmark results should be reproducible
    And performance comparisons should be statistically significant
    And benchmarking data should support optimization strategy decisions

  @performance @optimization @adaptive-optimization
  Scenario: Adaptive optimization based on performance feedback
    Given I have performance feedback from previous optimization runs
    When the system encounters similar project patterns
    Then optimization strategies should adapt automatically:
      | adaptation_trigger     | adaptive_response                          |
      | Slow processing        | reduce optimization depth, increase speed  |
      | High memory usage      | enable incremental processing              |
      | Low quality gains      | increase optimization aggressiveness       |
      | Architecture complexity| use architecture-specific optimizations   |
    And adaptive strategies should improve over time
    And manual override should always be available
    And adaptation decisions should be transparent to users

  @performance @optimization @cross-platform-performance
  Scenario: Cross-platform optimization performance validation
    Given I have optimization workloads running on different platforms:
      | platform | cpu_architecture | memory_model | expected_characteristics    |
      | Linux    | x86_64          | unified     | high performance baseline   |
      | macOS    | arm64           | unified     | excellent memory efficiency |
      | Windows  | x86_64          | segmented   | good performance            |
    When I measure optimization performance across platforms
    Then performance should be consistent across platforms:
      | performance_aspect     | cross_platform_variance |
      | Processing speed       | <20% difference         |
      | Memory efficiency      | <15% difference         |
      | Quality improvements   | <5% difference          |
      | Resource utilization   | <25% difference         |
    And platform-specific optimizations should be applied where beneficial
    And performance characteristics should be documented per platform