Feature: Optimization Level Progression
  As a Go developer using go-starter's optimization system
  I want to ensure optimization levels progressively apply more optimizations
  So that I can choose the appropriate level of code transformation

  Background:
    Given I am using go-starter CLI
    And the optimization system is available
    And the multi-level optimization system is available

  @optimization @progression
  Scenario: Verify optimization level ordering
    Given I have a test project with various optimization opportunities
    When I apply each optimization level in sequence
    Then each level should apply progressively more optimizations:
      | Level      | Imports | Variables | Functions | Safety |
      | none       | 0       | 0         | 0         | Max    |
      | safe       | Yes     | 0         | 0         | High   |
      | standard   | Yes     | Local     | 0         | Medium |
      | aggressive | Yes     | All       | Private   | Low    |
      | expert     | Yes     | All       | All       | Min    |

  @optimization @incremental
  Scenario Outline: Progressive optimization effects
    Given I have a "web-api" project with optimization opportunities
    When I apply "<level>" optimization
    Then the optimization results should match "<level>" expectations
    And import optimizations should be "<imports>"
    And variable removal should be "<variables>"
    And function removal should be "<functions>"
    And the project should still compile successfully

    Examples:
      | level      | imports  | variables   | functions     |
      | none       | disabled | disabled    | disabled      |
      | safe       | enabled  | disabled    | disabled      |
      | standard   | enabled  | local-only  | disabled      |
      | aggressive | enabled  | all-unused  | private-only  |
      | expert     | enabled  | all-unused  | all-unused    |

  @optimization @safety-progression
  Scenario: Safety decreases with higher optimization levels
    Given I have projects with risky optimization opportunities
    When I apply different optimization levels
    Then safety warnings should follow this pattern:
      | Level      | Warnings           | Risk Level |
      | safe       | None              | Minimal    |
      | standard   | Informational     | Low        |
      | aggressive | Multiple warnings | Medium     |
      | expert     | Critical warnings | High       |

  @optimization @metrics
  Scenario: Optimization metrics increase with levels
    Given I have a project with measurable optimization opportunities:
      | Type              | Count |
      | Unused imports    | 10    |
      | Unused variables  | 15    |
      | Unused functions  | 8     |
      | Dead code blocks  | 5     |
    When I apply each optimization level
    Then the metrics should show progressive improvement:
      | Level      | Imports Removed | Variables Removed | Functions Removed | Total Impact |
      | none       | 0              | 0                 | 0                 | 0%           |
      | safe       | 10             | 0                 | 0                 | 26%          |
      | standard   | 10             | 8                 | 0                 | 47%          |
      | aggressive | 10             | 15                | 5                 | 79%          |
      | expert     | 10             | 15                | 8                 | 87%          |

  @optimization @validation
  Scenario: Validate level transition correctness
    Given I have a project optimized at "safe" level
    When I optimize it again at "standard" level
    Then standard optimizations should build upon safe optimizations
    And no safe optimizations should be reverted
    And additional standard-level optimizations should be applied

  @optimization @file-safety
  Scenario Outline: File modification safety by level
    Given I have a project with critical files
    When I apply "<level>" optimization
    Then file modification should respect "<level>" safety:
      | File Type        | none | safe | standard | aggressive | expert |
      | main.go          | No   | No   | Yes      | Yes        | Yes    |
      | interfaces       | No   | No   | No       | No         | Yes    |
      | public functions | No   | No   | No       | No         | Yes    |
      | private functions| No   | No   | No       | Yes        | Yes    |
      | test files       | No   | No   | No       | No         | Yes    |

    Examples:
      | level      |
      | none       |
      | safe       |
      | standard   |
      | aggressive |
      | expert     |

  @optimization @complexity
  Scenario: Complex code optimization progression
    Given I have a project with complex code patterns:
      """go
      // Complex nested conditions
      if condition1 {
          if condition2 {
              if condition3 {
                  // Deep nesting
              }
          }
      }
      
      // Unused complex functions
      func complexUnusedPrivate() {
          // Complex logic
      }
      
      func ComplexUnusedPublic() {
          // Complex public logic
      }
      """
    When I apply each optimization level
    Then optimization should handle complexity appropriately:
      | Level      | Nested Conditions | Private Functions | Public Functions |
      | safe       | Unchanged        | Unchanged         | Unchanged        |
      | standard   | Simplified       | Unchanged         | Unchanged        |
      | aggressive | Simplified       | Removed           | Unchanged        |
      | expert     | Simplified       | Removed           | Analyzed         |

  @optimization @rollback
  Scenario: Level downgrade handling
    Given I have a project optimized at "expert" level
    When I apply "safe" level optimization
    Then the system should warn about potential regression
    And suggest using the original unoptimized source
    And provide clear downgrade instructions

  @optimization @profile-interaction
  Scenario: Level and profile interaction
    Given I have optimization profiles configured
    When I set both level and profile
    Then profile settings should override level defaults
    And the effective configuration should be clearly reported
    And optimization should follow the merged settings

  @optimization @benchmark
  Scenario: Performance benchmark across levels
    Given I have a large project for benchmarking
    When I measure optimization performance at each level
    Then processing time should scale with optimization complexity:
      | Level      | Relative Time | Files/Second |
      | none       | 1.0x         | N/A          |
      | safe       | 1.2x         | 100+         |
      | standard   | 1.5x         | 80+          |
      | aggressive | 2.0x         | 50+          |
      | expert     | 3.0x         | 30+          |