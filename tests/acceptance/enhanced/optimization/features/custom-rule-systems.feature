Feature: Custom Rule Systems for Optimization
  As a Go developer with specific optimization requirements
  I want to define custom optimization rules and patterns for my projects
  So that I can tailor the optimization system to my unique needs and standards

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And the optimization system is available
    And custom rule system infrastructure is enabled

  @enhanced @custom-rules @optimization @critical
  Scenario: Define and apply custom optimization rules
    Given I want to create custom optimization rules for my project
    When I define custom rules using the rule definition DSL:
      | rule_name               | rule_type    | target_pattern           | action                    | priority |
      | RemoveDebugImports      | import       | ".*debug.*"             | remove_if_unused          | high     |
      | ConsolidateErrors       | variable     | "err[0-9]+"             | consolidate_to_single     | medium   |
      | SimplifyNestedIfs       | structure    | nested_if_depth > 3     | extract_to_function       | high     |
      | OptimizeStringConcat    | expression   | strings_concat_loop     | use_strings_builder       | medium   |
      | RemoveEmptyFunctions    | function     | body_lines == 0         | remove_if_private         | low      |
      | GroupRelatedImports     | import_group | similar_packages        | group_by_domain           | medium   |
      | ExtractMagicNumbers     | literal      | numeric_literal > 1     | extract_to_constant       | high     |
      | SimplifyBooleanExpr     | expression   | redundant_boolean       | simplify_expression       | medium   |
    Then the custom rules should be validated for correctness
    And the rules should be applicable to relevant code patterns
    And rule conflicts should be detected and reported
    And custom rules should integrate seamlessly with built-in rules

  @enhanced @custom-rules @rule-sets @organization
  Scenario: Create and manage custom rule sets for different scenarios
    Given I have different optimization scenarios requiring specific rules
    When I create custom rule sets for each scenario:
      | rule_set_name      | description                          | included_rules                                    | excluded_rules        |
      | PerformanceCritical| High-performance system optimization | AggressiveInlining, LoopOptimization, CacheFriendly| SafetyChecks         |
      | SecurityFocused    | Security-first optimization          | SecurePatterns, NoHardcodedSecrets, InputValidation| PerformanceInlining  |
      | Readability        | Code clarity and maintenance         | ExtractComplexLogic, AddDocComments, SimplifyNaming| AggressiveOptimization|
      | StartupOptimized   | Fast application startup             | LazyInit, MinimalImports, DeferExpensiveOps       | EagerLoading         |
      | MemoryConstrained  | Low memory footprint                 | PoolReuse, MinimalAllocations, CompactStructs     | Caching              |
      | TestFriendly       | Testability improvements             | DependencyInjection, InterfaceExtraction, MockableAPIs| GlobalState        |
    Then rule sets should be easily selectable during optimization
    And rule sets should be composable for complex scenarios
    And conflicting rules between sets should be handled gracefully
    And custom rule sets should be shareable across projects

  @enhanced @custom-rules @pattern-matching @advanced
  Scenario: Advanced pattern matching for custom optimization rules
    Given I need to match complex code patterns for optimization
    When I define rules with advanced pattern matching:
      | pattern_name        | pattern_type | matching_criteria                                | optimization_action    |
      | ErrorChainPattern   | ast         | "if err != nil { return err }" repeated 3+ times| consolidate_error_handling |
      | InterfaceOveruse    | type        | interface{} used in >50% of function params      | suggest_concrete_types |
      | DeepNesting         | structure   | nesting_level > 4 && cyclomatic_complexity > 10  | extract_nested_logic   |
      | StringBuilderCandidate| loop      | string concatenation in loop with 10+ iterations| convert_to_builder     |
      | UnusedStructFields  | struct      | field_usage_count == 0 across all references    | remove_unused_field    |
      | DuplicateLogic      | function    | similarity_score > 0.85 with another function    | extract_common_function|
      | SlicePreallocation  | allocation  | append in loop without pre-allocation            | preallocate_slice      |
      | MapKeyOptimization  | map_usage   | string keys that could be constants              | extract_key_constants  |
    Then pattern matching should accurately identify target code
    And false positives should be minimized through context analysis
    And pattern matches should provide detailed match information
    And optimization actions should be safely applicable

  @enhanced @custom-rules @rule-priority @conflict-resolution
  Scenario: Rule priority and conflict resolution system
    Given I have multiple custom rules that may conflict
    When I define rules with priority and conflict resolution:
      | rule_name          | priority | conflicts_with      | resolution_strategy    | conditions            |
      | InlineSmallFuncs   | 100      | ExtractFunctions   | priority_based        | func_lines < 5        |
      | ExtractFunctions   | 80       | InlineSmallFuncs   | context_aware         | complexity > 10       |
      | RemoveComments     | 20       | AddDocumentation   | exclude_both          | production_build      |
      | AddDocumentation   | 90       | RemoveComments     | priority_based        | public_api            |
      | OptimizeLoops      | 85       | MaintainReadability| balanced_approach     | performance_critical  |
      | SimplifyCode       | 70       | PerformanceOpt     | user_preference       | readability_focus     |
    Then conflict detection should identify all potential conflicts
    And resolution strategies should be applied consistently
    And users should be notified of conflict resolutions
    And manual override options should be available

  @enhanced @custom-rules @context-aware @intelligence
  Scenario: Context-aware custom rule application
    Given I have custom rules that should apply based on context
    When I define context-aware rules:
      | rule_name           | context_conditions                              | optimization_action        | skip_conditions           |
      | OptimizeHTTPHandler | in_http_handler && request_processing          | streamline_response_path   | has_middleware            |
      | DatabaseQueryOpt    | in_database_layer && query_complexity > medium | add_query_optimization     | has_cache_layer           |
      | ConcurrencyPattern  | uses_goroutines && shared_state_detected      | add_synchronization        | already_synchronized      |
      | ErrorPropagation    | in_service_layer && returns_error             | standardize_error_handling | custom_error_handler      |
      | LoggingOptimization | has_logging && log_calls > 10_per_function    | consolidate_logging        | debug_mode                |
      | CacheableResults    | pure_function && expensive_computation         | add_memoization            | dynamic_inputs            |
    Then context analysis should accurately determine rule applicability
    And context-aware rules should respect architectural boundaries
    And skip conditions should prevent inappropriate optimizations
    And context information should be available in rule execution

  @enhanced @custom-rules @templates @reusability
  Scenario: Custom rule templates and reusable patterns
    Given I want to create reusable optimization rule templates
    When I define rule templates with parameters:
      | template_name       | parameters                    | rule_pattern                           | description                |
      | LimitChecker        | limit_value, comparison_op    | value {{.comparison_op}} {{.limit_value}} | Generic limit validation  |
      | PatternReplacer     | search_pattern, replacement   | replace({{.search_pattern}}, {{.replacement}}) | Pattern substitution |
      | ConditionalOptimize | condition, true_action, false_action | if {{.condition}} then {{.true_action}} else {{.false_action}} | Conditional optimization |
      | MetricThreshold     | metric_name, threshold, action| when {{.metric_name}} > {{.threshold}} do {{.action}} | Metric-based rules |
      | ArchitectureRule    | layer, allowed_imports        | in_layer({{.layer}}) check imports({{.allowed_imports}}) | Architecture enforcement |
    Then templates should be instantiable with specific parameters
    And template validation should ensure parameter completeness
    And templates should be shareable and versionable
    And template instances should behave like regular rules

  @enhanced @custom-rules @testing @validation
  Scenario: Test and validate custom optimization rules
    Given I have defined custom optimization rules
    When I test the rules against sample code:
      | test_case          | input_code                          | rule_applied         | expected_output                    | validation_result |
      | DebugImportRemoval | import "debug/pprof" (unused)      | RemoveDebugImports   | // import removed                  | pass             |
      | ErrorConsolidation | multiple err1, err2, err3 checks   | ConsolidateErrors    | single err variable                | pass             |
      | NestedIfSimplify   | 5-level nested if statements       | SimplifyNestedIfs    | extracted helper functions         | pass             |
      | StringBuilderOpt   | for loop with string concatenation | OptimizeStringConcat | strings.Builder implementation     | pass             |
      | EdgeCase           | valid but unusual pattern          | CustomRule           | no change (correctly skipped)      | pass             |
    Then rule testing should provide clear pass/fail results
    And test coverage should include positive and negative cases
    And performance impact of rules should be measured
    And rule safety should be validated through testing

  @enhanced @custom-rules @integration @ecosystem
  Scenario: Integration with go-starter optimization ecosystem
    Given I have custom rules defined for my organization
    When I integrate custom rules with the optimization system:
      | integration_point    | custom_rule_behavior              | interaction_with_builtin     | expected_outcome           |
      | Rule Loading        | Load from .gostarter-rules.yaml   | Merge with built-in rules   | Combined rule set active   |
      | Priority Resolution | Custom rules priority 1-100       | Built-in rules 101-200      | Custom rules take precedence|
      | Profile Integration | Custom rules in optimization profiles | Extend existing profiles | Enhanced profiles available |
      | CLI Integration     | --custom-rules flag               | Works with all commands     | Seamless CLI experience    |
      | Reporting           | Custom rule metrics in reports    | Unified reporting format    | Comprehensive insights     |
      | Configuration       | Rule configuration in config file | Inherits global settings    | Flexible configuration     |
    Then custom rules should integrate without disrupting existing functionality
    And performance should not degrade with custom rules
    And custom rule metrics should be included in optimization reports
    And rollback should be possible if custom rules cause issues

  @enhanced @custom-rules @sharing @community
  Scenario: Share and import community optimization rules
    Given I want to share my custom rules with the community
    When I package and share my optimization rules:
      | sharing_method     | rule_package_format          | metadata_included           | distribution_channel    |
      | Export Package     | .gostarter-rules package     | author, version, license    | GitHub releases        |
      | Rule Registry      | YAML with validation schema  | description, examples       | go-starter registry    |
      | Git Repository     | Structured rule repository   | README, tests, examples     | Public Git repo        |
      | Organization Share | Internal rule library        | team, use cases, metrics    | Private registry       |
    Then rule packages should be easily importable by others
    And imported rules should be validated before use
    And rule versioning should prevent compatibility issues
    And community feedback should be incorporable

  @enhanced @custom-rules @performance @efficiency
  Scenario: Performance optimization of custom rule execution
    Given I have complex custom rules that need efficient execution
    When I optimize custom rule performance:
      | optimization_technique   | applied_to                    | performance_improvement    | trade_offs               |
      | Rule Caching            | AST pattern matching          | 50-70% faster matching    | Memory usage increase    |
      | Parallel Execution      | Independent rules             | 3-4x throughput           | Complexity increase      |
      | Early Termination       | Exclusive rule conditions     | 20-30% faster overall     | Rule order dependency    |
      | Pattern Precompilation  | Regex and AST patterns       | 40-60% faster execution   | Startup time increase    |
      | Incremental Processing  | Large codebases              | 80% faster on changes     | State management needed  |
      | Smart Scheduling        | Rule dependency ordering      | 15-25% efficiency gain    | Analysis overhead        |
    Then custom rule execution should be performant
    And performance should scale with codebase size
    And Resource usage should be predictable
    And Performance metrics should be available

  @enhanced @custom-rules @debugging @troubleshooting
  Scenario: Debug and troubleshoot custom optimization rules
    Given I need to debug why my custom rules aren't working as expected
    When I use debugging features for custom rules:
      | debug_feature        | capability                          | output_information              | use_case                   |
      | Rule Trace          | Step-by-step rule execution        | Match attempts, decisions       | Understanding rule behavior |
      | Pattern Debugger    | Interactive pattern matching       | AST visualization, matches      | Refining patterns          |
      | Dry Run Mode        | Preview changes without applying   | Planned modifications           | Safety verification        |
      | Conflict Analyzer   | Detailed conflict information      | Conflicting rules, reasons      | Resolving rule conflicts   |
      | Performance Profiler| Rule execution timing              | Time per rule, bottlenecks      | Optimization opportunities |
      | Match Explorer      | Browse all pattern matches         | Match locations, context        | Pattern effectiveness      |
    Then debugging should provide clear insights into rule behavior
    And troubleshooting should be efficient and informative
    And debug output should be actionable
    And issues should be quickly identifiable and fixable

  @enhanced @custom-rules @evolution @maintenance
  Scenario: Evolve and maintain custom rules over time
    Given I have custom rules that need updates as my codebase evolves
    When I maintain and evolve my custom rules:
      | maintenance_task     | trigger_condition              | update_action                  | validation_required      |
      | Rule Effectiveness  | Low match rate < 5%            | Review and update pattern      | Re-test on sample code   |
      | Pattern Drift       | Codebase patterns changed      | Adapt rules to new patterns    | Regression testing       |
      | Performance Tuning  | Rule execution > 100ms         | Optimize rule implementation   | Performance benchmarks   |
      | Deprecation         | Rule no longer needed          | Mark deprecated, plan removal  | Impact analysis          |
      | Enhancement         | New optimization opportunity   | Extend rule capabilities       | Compatibility check      |
      | Migration           | Major go-starter update        | Update rule syntax/API         | Full test suite run      |
    Then rule maintenance should be systematic and tracked
    And rule evolution should maintain backward compatibility
    And deprecated rules should have migration paths
    And rule lifecycle should be well-managed

  @enhanced @custom-rules @security @safety
  Scenario: Ensure security and safety of custom optimization rules
    Given I need to ensure custom rules don't introduce security issues
    When I implement security measures for custom rules:
      | security_measure     | protection_against             | implementation_method          | validation_approach      |
      | Input Validation    | Malicious rule definitions     | Schema validation, sanitization| Automated security scan  |
      | Sandbox Execution   | Harmful rule side effects      | Isolated rule execution        | Container/VM testing     |
      | Code Injection Prevention| Rule-based code injection  | AST-only modifications         | Static analysis          |
      | Access Control      | Unauthorized rule changes      | Role-based permissions         | Authentication required  |
      | Audit Trail         | Tracking rule modifications    | Change log, version history    | Compliance reporting     |
      | Safe Defaults       | Dangerous optimization         | Conservative default settings  | Manual opt-in for risky  |
    Then custom rules should not introduce security vulnerabilities
    And rule execution should be sandboxed and safe
    And malicious rules should be detected and rejected
    And security audit trails should be maintained