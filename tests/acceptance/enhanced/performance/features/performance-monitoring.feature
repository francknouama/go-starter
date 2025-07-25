Feature: Performance Monitoring and Resource Usage
  As a go-starter maintainer
  I want to monitor performance and resource usage of generated projects
  So that we can ensure optimal performance across all blueprints

  Background:
    Given I have the go-starter CLI available
    And performance monitoring is enabled
    And resource tracking is configured

  Scenario: Blueprint generation performance benchmarking
    Given I measure generation time for each blueprint
    When I generate projects with various configurations
    Then generation times should meet performance targets:
      | Blueprint           | Simple Config | Complex Config | Max Time |
      | cli-simple          | < 1s         | < 2s          | 3s       |
      | cli-standard        | < 2s         | < 3s          | 5s       |
      | web-api-standard    | < 3s         | < 5s          | 8s       |
      | web-api-clean       | < 4s         | < 6s          | 10s      |
      | web-api-ddd         | < 4s         | < 7s          | 10s      |
      | web-api-hexagonal   | < 4s         | < 7s          | 10s      |
      | microservice        | < 5s         | < 8s          | 12s      |
      | workspace           | < 6s         | < 10s         | 15s      |
    And memory usage should not exceed 200MB during generation
    And CPU usage should not spike above 80%
    And disk I/O should be optimized with batching

  Scenario: Generated project compilation performance
    Given I generate projects with different configurations
    When I measure compilation performance
    Then compilation metrics should be within acceptable ranges:
      | Project Type    | First Build | Incremental Build | Binary Size | Memory Usage |
      | cli-simple      | < 5s        | < 2s              | < 10MB      | < 100MB      |
      | cli-standard    | < 10s       | < 3s              | < 20MB      | < 200MB      |
      | web-api         | < 20s       | < 5s              | < 30MB      | < 300MB      |
      | microservice    | < 30s       | < 8s              | < 40MB      | < 400MB      |
    And build cache should be effectively utilized
    And parallel compilation should be enabled
    And module download time should be tracked

  Scenario: Runtime performance characteristics
    Given I have generated projects with various architectures
    When I benchmark runtime performance
    Then runtime metrics should meet standards:
      | Architecture | Startup Time | Request Latency (p99) | Memory Footprint | Goroutines |
      | standard     | < 100ms      | < 10ms               | < 50MB           | < 100      |
      | clean        | < 150ms      | < 15ms               | < 60MB           | < 120      |
      | ddd          | < 200ms      | < 20ms               | < 70MB           | < 150      |
      | hexagonal    | < 200ms      | < 20ms               | < 70MB           | < 150      |
    And memory leaks should not occur under load
    And garbage collection pauses should be minimal
    And CPU usage should scale linearly with load

  Scenario: Database operation performance
    Given I have projects with different database configurations
    When I benchmark database operations
    Then database performance should be optimal:
      | Operation        | PostgreSQL | MySQL   | SQLite  | Target SLA |
      | Single Insert    | < 5ms      | < 5ms   | < 2ms   | 10ms       |
      | Bulk Insert (1k) | < 50ms     | < 60ms  | < 30ms  | 100ms      |
      | Simple Query     | < 2ms      | < 2ms   | < 1ms   | 5ms        |
      | Complex Join     | < 20ms     | < 25ms  | < 15ms  | 50ms       |
      | Transaction      | < 30ms     | < 35ms  | < 20ms  | 60ms       |
    And connection pooling should be properly configured
    And prepared statements should be utilized
    And query optimization should be validated

  Scenario: Load testing generated APIs
    Given I have generated web APIs with different frameworks
    When I perform load testing with standard scenarios
    Then APIs should handle load gracefully:
      | Framework | RPS (1 instance) | Latency p50 | Latency p99 | Error Rate | CPU Usage |
      | gin       | > 10000         | < 5ms       | < 50ms      | < 0.1%     | < 70%     |
      | fiber     | > 12000         | < 4ms       | < 40ms      | < 0.1%     | < 65%     |
      | echo      | > 9000          | < 6ms       | < 60ms      | < 0.1%     | < 75%     |
      | chi       | > 8000          | < 7ms       | < 70ms      | < 0.1%     | < 80%     |
    And throughput should degrade gracefully under overload
    And circuit breakers should activate appropriately
    And rate limiting should function correctly

  Scenario: Memory profiling and leak detection
    Given I run generated projects under sustained load
    When I profile memory usage over time
    Then memory characteristics should be stable:
      | Metric                    | Acceptable Range        |
      | Heap growth rate          | < 1MB/minute           |
      | Goroutine leak detection  | No persistent growth   |
      | Memory allocation rate    | < 100MB/second         |
      | GC pause time (p99)       | < 10ms                 |
      | Live object count         | Stable after warmup    |
    And memory profiling data should be collected
    And potential leak sources should be identified
    And optimization recommendations should be generated

  Scenario: Cross-platform performance validation
    Given I test on different operating systems
    When I measure performance metrics
    Then performance should be consistent across platforms:
      | Platform | Generation Time Variance | Runtime Performance Variance | Build Time Variance |
      | Linux    | Baseline                | Baseline                     | Baseline            |
      | macOS    | < 10%                   | < 15%                        | < 20%               |
      | Windows  | < 20%                   | < 25%                        | < 30%               |
    And platform-specific optimizations should be documented
    And file system operations should be optimized per platform
    And system resource usage should be appropriate

  Scenario: Performance regression detection
    Given I have historical performance baselines
    When I run performance tests on new changes
    Then regressions should be automatically detected:
      | Metric Type          | Regression Threshold | Action Required    |
      | Generation time      | > 10% increase      | Warning            |
      | Compilation time     | > 15% increase      | Investigation      |
      | Runtime performance  | > 20% degradation   | Block release      |
      | Memory usage         | > 25% increase      | Critical review    |
    And performance trends should be tracked
    And regression reports should be generated
    And root cause analysis should be facilitated

  Scenario: Resource usage optimization validation
    Given I analyze resource consumption patterns
    When I identify optimization opportunities
    Then optimizations should be validated:
      | Resource Type | Optimization Strategy         | Expected Improvement |
      | CPU          | Parallel processing           | 30-50% reduction     |
      | Memory       | Object pooling                | 20-30% reduction     |
      | Disk I/O     | Batch file operations         | 40-60% reduction     |
      | Network      | Connection pooling            | 50-70% reduction     |
    And optimization impact should be measured
    And trade-offs should be documented
    And best practices should be updated

  Scenario: Performance monitoring dashboard
    Given I collect performance metrics from all tests
    When I generate performance reports
    Then dashboards should provide insights:
      | Dashboard Section    | Key Metrics                          | Update Frequency |
      | Generation Speed     | Time per blueprint, success rate     | Per test run     |
      | Runtime Performance  | Latency, throughput, error rates     | Continuous       |
      | Resource Usage       | CPU, memory, disk, network           | Real-time        |
      | Trend Analysis       | Week-over-week changes               | Daily            |
      | Anomaly Detection    | Outliers and unexpected patterns     | Immediate        |
    And alerts should be configured for degradation
    And historical data should be retained
    And comparative analysis should be available