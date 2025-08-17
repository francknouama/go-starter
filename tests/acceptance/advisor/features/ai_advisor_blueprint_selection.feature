Feature: AI Advisor Blueprint Selection
  As a developer
  I want intelligent blueprint recommendations based on my project requirements
  So that I choose the right architecture pattern for my specific needs

  Background:
    Given the AI advisor is available
    And the knowledge base contains blueprint patterns
    And the blueprint registry is loaded

  Scenario: Simple CLI tool for junior developers
    Given I need a "cli" project
    And my team has "junior" experience level
    And my project domain is "devtools"
    And I prefer "simple" development style
    When I ask for blueprint recommendations
    Then I should get "cli-simple" as the primary recommendation
    And the recommended complexity should be "simple"
    And the reasoning should mention beginner-friendly approach
    And the alternative should include "cli" for future growth

  Scenario: Enterprise API with complex business logic
    Given I need a "api" project
    And my team has "expert" experience level
    And my project domain is "fintech"
    And I have "complex" database requirements
    And I need "enterprise" compliance
    When I ask for blueprint recommendations
    Then I should get "web-api-ddd" or "web-api-hexagonal" as primary recommendation
    And the recommended architecture should support domain modeling
    And the reasoning should mention separation of concerns
    And the estimated file count should be at least 50

  Scenario: MVP development with tight timeline
    Given I need a "api" project
    And my team has "mixed" experience level
    And my time to market is "mvp"
    And my budget is "startup"
    When I ask for blueprint recommendations
    Then I should get "web-api" as the primary recommendation
    And the recommended architecture should be "standard"
    And the reasoning should mention quick development
    And the estimated development time should be optimized

  Scenario: High-performance microservice
    Given I need a "microservice" project
    And I expect "massive" load
    And I need "realtime" response times
    And my team has "expert" experience level
    When I ask for blueprint recommendations
    Then I should get "microservice" as the primary recommendation
    And the recommended framework should support high performance
    And the recommended logger should be performance-optimized
    And the reasoning should mention performance characteristics

  Scenario: Serverless function for event processing
    Given I need a "lambda" project
    And my project domain is "iot"
    And I have "event-driven" data patterns
    And my deployment target is "cloud"
    When I ask for blueprint recommendations
    Then I should get "lambda-event-processing" as the primary recommendation
    And the recommended architecture should support event handling
    And the reasoning should mention serverless benefits
    And the estimated files should be minimal

  Scenario: Library for code reuse
    Given I need a "library" project
    And I want to publish to the community
    And my team has "senior" experience level
    And I need comprehensive documentation
    When I ask for blueprint recommendations
    Then I should get "library" as the primary recommendation
    And the recommended structure should include examples
    And the reasoning should mention public API design
    And the alternatives should suggest different library patterns

  Scenario Outline: Blueprint selection based on domain expertise
    Given I need a "<project_type>" project
    And my project domain is "<domain>"
    And my team has "<experience>" experience level
    When I ask for blueprint recommendations
    Then I should get domain-appropriate blueprint recommendations
    And the confidence should reflect domain match
    And the reasoning should mention domain-specific considerations

    Examples:
      | project_type | domain          | experience | 
      | api          | healthcare      | expert     |
      | api          | e-commerce      | senior     |
      | api          | content-mgmt    | mixed      |
      | cli          | automation      | junior     |
      | cli          | deployment      | senior     |
      | lambda       | data-processing | expert     |
      | lambda       | webhooks        | mixed      |
      | microservice | distributed     | expert     |

  Scenario: Blueprint selection with conflicting requirements
    Given I need a "api" project
    And my project domain is "fintech"
    And my team has "junior" experience level
    But I need "enterprise" compliance
    And I expect "high" load
    When I ask for blueprint recommendations
    Then I should get a compromise recommendation
    And the reasoning should acknowledge the complexity mismatch
    And the alternatives should include both simple and complex options
    And the recommendation should suggest team training or consultation

  Scenario: Blueprint selection for monolith migration
    Given I need a "microservice" project
    And I'm migrating from a monolith
    And my team has "mixed" experience level
    And I want to start with one service
    When I ask for blueprint recommendations
    Then I should get migration-friendly recommendations
    And the reasoning should mention incremental migration
    And the alternatives should include monolith-to-microservice patterns
    And the recommendation should include migration guidance

  Scenario: Blueprint validation and availability
    Given I request blueprint recommendations
    When the advisor suggests a blueprint
    Then the blueprint should exist in the registry
    And the blueprint should be valid and parseable
    And the blueprint should support the recommended features
    And the template files should be complete

  Scenario: Blueprint selection consistency
    Given I have consistent project requirements
    When I request blueprint recommendations multiple times
    Then I should get consistent primary recommendations
    And the confidence scores should be stable
    And the reasoning should be consistent
    And the alternatives should remain similar