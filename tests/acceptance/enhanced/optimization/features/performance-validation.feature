Feature: Optimization Pipeline Performance Validation
  As a Go developer using go-starter's optimization system
  I want the optimization pipeline to meet performance benchmarks
  So that I can optimize projects efficiently without impacting my development workflow

  Background:
    Given I am using go-starter CLI
    And the optimization system is available
    And performance monitoring is enabled

  @performance @benchmarks
  Scenario: Meet optimization speed benchmarks
    Given I have projects of different sizes:
      | Project Size | Files | Lines of Code | Expected Processing |
      | Small        | 5-10  | < 1,000      | < 2 seconds        |
      | Medium       | 20-50 | 1,000-10,000 | < 10 seconds       |
      | Large        | 100+  | 10,000+      | < 60 seconds       |
    When I optimize each project type
    Then processing time should meet benchmark requirements
    And throughput should be at least 100 files/minute for small files
    And memory usage should remain under 500MB for large projects

  @performance @memory-efficiency
  Scenario: Memory usage efficiency validation
    Given I have a project that will test memory efficiency
    When I run optimization with memory profiling enabled
    Then peak memory usage should not exceed baseline + 200MB
    And memory should be released progressively during optimization
    And garbage collection should occur regularly
    And no memory leaks should be detected after completion

  @performance @concurrent-processing
  Scenario: Concurrent processing performance
    Given I configure concurrent file processing:
      | Setting              | Value |
      | MaxConcurrentFiles   | 8     |
      | EnableParallelism    | true  |
      | MemoryLimitMB        | 512   |
    When I optimize a large project with concurrent processing
    Then processing should utilize multiple CPU cores effectively
    And total processing time should be reduced compared to sequential
    And memory usage per thread should be controlled
    And no race conditions should occur

  @performance @scalability
  Scenario Outline: Performance scalability validation
    Given I have a project with "<file_count>" Go files
    And each file has approximately "<lines_per_file>" lines
    When I optimize the project with "<optimization_level>" level
    Then processing time should scale linearly with file count
    And memory usage should remain proportional to project size
    And optimization quality should not degrade with scale

    Examples:
      | file_count | lines_per_file | optimization_level |
      | 10         | 100           | standard          |
      | 50         | 200           | standard          |
      | 100        | 300           | standard          |
      | 200        | 400           | aggressive        |
      | 500        | 500           | expert            |

  @performance @resource-utilization
  Scenario: System resource utilization optimization
    Given I monitor system resources during optimization:
      | Resource    | Monitoring Metric        |
      | CPU         | Utilization percentage   |
      | Memory      | Peak and average usage   |
      | Disk I/O    | Read/write operations    |
      | File Handle | Open file descriptors    |
    When I run optimization on multiple project types
    Then CPU utilization should be efficient but not excessive
    And disk I/O should be minimized through smart caching
    And file handles should be properly managed and released
    And system responsiveness should remain good

  @performance @optimization-levels
  Scenario: Performance characteristics by optimization level
    Given I have the same project to optimize at different levels
    When I measure performance for each optimization level:
      | Level      | Expected Relative Time | Quality Impact |
      | safe       | 1.0x (baseline)       | Low           |
      | standard   | 1.5x                  | Medium        |
      | aggressive | 2.5x                  | High          |
      | expert     | 4.0x                  | Maximum       |
    Then performance should scale appropriately with optimization complexity
    And quality improvements should justify performance costs
    And users should be warned about performance implications

  @performance @caching-effectiveness
  Scenario: Optimization caching and reuse effectiveness
    Given I have projects with overlapping code patterns
    When I optimize similar projects sequentially
    Then AST parsing results should be cached where possible
    And repeated pattern analysis should be optimized
    And overall processing time should decrease for similar patterns
    And cache hit rates should be reported and monitored

  @performance @profile-specific-performance
  Scenario Outline: Profile-specific performance characteristics
    Given I have a "<project_type>" project suitable for "<profile>" profile
    When I optimize using the "<profile>" profile
    Then performance should match profile-specific expectations:
      | Profile      | Expected Speed | Memory Usage | Optimization Depth |
      | conservative | Fast          | Low          | Minimal           |
      | balanced     | Medium        | Medium       | Moderate          |
      | performance  | Slower        | Higher       | Comprehensive     |
      | maximum      | Slowest       | Highest      | Complete          |
    And profile recommendations should consider performance trade-offs

    Examples:
      | project_type | profile      |
      | web-api      | conservative |
      | web-api      | balanced     |
      | web-api      | performance  |
      | cli          | balanced     |
      | library      | conservative |
      | microservice | performance  |

  @performance @regression-testing
  Scenario: Performance regression detection
    Given I have baseline performance metrics from previous optimization runs
    When I run optimization with the current implementation
    Then performance should not regress beyond acceptable thresholds:
      | Metric                    | Acceptable Regression |
      | Processing time          | < 20% slower          |
      | Memory usage            | < 30% increase        |
      | Files processed/second   | < 15% decrease        |
      | Optimization quality     | No degradation        |
    And any performance regressions should be clearly reported
    And suggestions for performance improvement should be provided

  @performance @real-time-monitoring
  Scenario: Real-time performance monitoring during optimization
    Given I enable real-time performance monitoring
    When I run optimization on a large project
    Then progress should be reported with performance metrics:
      | Metric                  | Update Frequency |
      | Files processed         | Every 10 files   |
      | Current memory usage    | Every 30 seconds |
      | Estimated time remaining| Every minute     |
      | Processing rate         | Continuous       |
    And users should be able to monitor optimization progress
    And performance bottlenecks should be identified in real-time

  @performance @stress-testing
  Scenario: Performance under stress conditions
    Given I create stress test conditions:
      | Stress Factor           | Value              |
      | Concurrent optimizations| 5 projects         |
      | Large file sizes        | 5MB+ per file      |
      | Complex code patterns   | Deep nesting       |
      | Limited memory          | 1GB system limit   |
      | High CPU load          | 80% utilization    |
    When I run optimization under stress
    Then the system should handle stress gracefully
    And performance should degrade predictably
    And error handling should remain responsive
    And no system crashes should occur

  @performance @optimization-quality-vs-speed
  Scenario: Balance optimization quality vs processing speed
    Given I have projects with varying optimization opportunities
    When I optimize with different quality vs speed trade-offs:
      | Mode              | Quality Priority | Speed Priority | Expected Behavior        |
      | quality-focused   | High            | Low            | Thorough, slower         |
      | balanced          | Medium          | Medium         | Good balance             |
      | speed-focused     | Low             | High           | Fast, basic optimization |
    Then each mode should deliver expected quality-speed trade-offs
    And users should be able to choose appropriate modes
    And trade-offs should be clearly documented

  @performance @batch-processing
  Scenario: Batch processing multiple projects efficiently
    Given I have multiple projects to optimize:
      | Project    | Type        | Size   |
      | project-a  | web-api     | Medium |
      | project-b  | cli         | Small  |
      | project-c  | library     | Large  |
      | project-d  | microservice| Medium |
    When I optimize all projects in batch mode
    Then batch processing should be more efficient than individual runs
    And shared resources should be reused across projects
    And overall processing time should be optimized
    And progress should be reported for the entire batch

  @performance @performance-profiling
  Scenario: Detailed performance profiling and analysis
    Given I enable detailed performance profiling
    When I optimize a representative project
    Then detailed profiling data should be collected:
      | Profiling Data          | Details                    |
      | Function call times     | Time spent in each function |
      | Memory allocation       | Allocation patterns        |
      | Garbage collection      | GC frequency and duration  |
      | I/O operations          | File read/write patterns   |
      | CPU usage patterns      | Hot spots and bottlenecks  |
    And profiling results should identify optimization opportunities
    And performance recommendations should be generated

  @performance @network-efficiency
  Scenario: Network and I/O efficiency for remote projects
    Given I have projects stored in different locations:
      | Location        | Type              | Network Impact |
      | Local disk      | Standard access   | None          |
      | Network drive   | Remote filesystem | Medium        |
      | Cloud storage   | Remote API access | High          |
    When I optimize projects from different locations
    Then I/O operations should be minimized and optimized
    And network latency should be handled gracefully
    And caching should reduce repeated remote access
    And performance should remain acceptable for all locations

  @performance @performance-reporting
  Scenario: Comprehensive performance reporting
    Given I complete optimization of various projects
    When I generate performance reports
    Then reports should include comprehensive metrics:
      | Report Section          | Content                        |
      | Executive Summary       | Key performance indicators     |
      | Processing Times        | Breakdown by phase and file    |
      | Resource Usage          | Memory, CPU, I/O statistics    |
      | Optimization Impact     | Quality improvements achieved  |
      | Comparative Analysis    | Performance vs previous runs   |
      | Recommendations         | Performance improvement tips   |
    And reports should be exportable in multiple formats
    And historical performance trends should be tracked