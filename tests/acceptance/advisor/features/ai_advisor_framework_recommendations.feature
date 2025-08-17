Feature: AI Advisor Framework Recommendations
  As a developer
  I want context-aware framework recommendations
  So that I choose the best tools for my specific project requirements

  Background:
    Given the AI advisor is available
    And the framework knowledge base is loaded
    And the compatibility matrix is available

  Scenario: High-performance API framework selection
    Given I need an "api" project
    And I expect "massive" load
    And I need "realtime" response times
    And my team has "expert" experience level
    When I ask for framework recommendations
    Then I should get high-performance frameworks like "gin", "fiber"
    And the reasoning should mention performance benchmarks
    And the recommendation should include optimization tips
    And the confidence should be high for performance scenarios

  Scenario: Beginner-friendly CLI framework selection
    Given I need a "cli" project
    And my team has "junior" experience level
    And I prefer "simple" development approach
    When I ask for framework recommendations
    Then I should get "cobra" as the primary recommendation
    And the reasoning should mention ease of use
    And the alternatives should include framework-less approach
    And the recommendation should include learning resources

  Scenario: Enterprise API with complex routing
    Given I need an "api" project
    And I have complex routing requirements
    And I need middleware support
    And my team has "senior" experience level
    When I ask for framework recommendations
    Then I should get frameworks with advanced routing like "gin", "echo"
    And the reasoning should mention middleware ecosystem
    And the recommendation should include routing patterns
    And the alternatives should cover different middleware approaches

  Scenario: Lightweight microservice framework
    Given I need a "microservice" project
    And I want minimal overhead
    And I need fast startup times
    And my deployment target is "kubernetes"
    When I ask for framework recommendations
    Then I should get lightweight frameworks like "fiber", "chi"
    And the reasoning should mention startup performance
    And the recommendation should include container optimization
    And the resource usage should be minimal

  Scenario: Framework selection for specific domains
    Given I need an "api" project
    And my project domain is "fintech"
    And I need strong security features
    And I require audit trails
    When I ask for framework recommendations
    Then I should get security-focused framework recommendations
    And the reasoning should mention security features
    And the recommendation should include security middleware
    And the compliance considerations should be addressed

  Scenario Outline: Framework compatibility with project types
    Given I need a "<project_type>" project
    And my team has "<experience>" experience level
    When I ask for framework recommendations
    Then I should get frameworks compatible with "<project_type>"
    And the framework should support the project architecture
    And the recommendation should be appropriate for "<experience>" level

    Examples:
      | project_type | experience |
      | api          | junior     |
      | api          | senior     |
      | api          | expert     |
      | cli          | junior     |
      | cli          | senior     |
      | microservice | expert     |
      | library      | mixed      |

  Scenario: Framework ecosystem considerations
    Given I need an "api" project
    And I want rich ecosystem support
    And I plan to use third-party integrations
    When I ask for framework recommendations
    Then I should get frameworks with strong ecosystems
    And the reasoning should mention community support
    And the recommendation should include popular middleware
    And the alternatives should compare ecosystem sizes

  Scenario: Framework migration considerations
    Given I'm migrating from "express.js" to Go
    And I need similar features and patterns
    And my team has "mixed" Go experience
    When I ask for framework recommendations
    Then I should get migration-friendly frameworks
    And the reasoning should mention familiar patterns
    And the recommendation should include migration guides
    And the learning curve should be addressed

  Scenario: Framework selection for testing requirements
    Given I need an "api" project
    And I require comprehensive testing
    And I want test-friendly architecture
    And my team follows TDD practices
    When I ask for framework recommendations
    Then I should get test-friendly frameworks
    And the reasoning should mention testing capabilities
    And the recommendation should include testing patterns
    And the mocking support should be addressed

  Scenario: Framework performance comparison
    Given I need performance-critical application
    And I want to compare framework options
    When I ask for framework recommendations with benchmarks
    Then I should get performance comparison data
    And the reasoning should include benchmark results
    And the recommendation should mention trade-offs
    And the use case suitability should be explained

  Scenario: Framework recommendation validation
    Given I receive a framework recommendation
    When I validate the recommendation
    Then the framework should be actively maintained
    And the framework should have stable API
    And the framework should support Go version compatibility
    And the framework should have good documentation

  Scenario: Custom framework needs
    Given I have very specific requirements
    And standard frameworks don't fit perfectly
    And my team has "expert" experience level
    When I ask for framework recommendations
    Then I should get custom solution suggestions
    And the reasoning should explain why custom is better
    And the recommendation should include implementation guidance
    And the trade-offs should be clearly explained

  Scenario: Framework selection for team productivity
    Given I prioritize team productivity over performance
    And my team has "mixed" experience level
    And I want quick development cycles
    When I ask for framework recommendations
    Then I should get productivity-focused frameworks
    And the reasoning should mention development speed
    And the recommendation should include productivity tools
    And the learning resources should be comprehensive

  Scenario: Logger framework integration
    Given I need an "api" project
    And I want structured logging
    And I need performance monitoring
    When I ask for framework and logger recommendations
    Then the framework and logger should be compatible
    And the reasoning should mention integration benefits
    And the recommendation should include configuration examples
    And the observability features should be highlighted