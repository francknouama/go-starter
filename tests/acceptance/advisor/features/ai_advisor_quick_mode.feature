Feature: AI Advisor Quick Mode Recommendations
  As a developer
  I want to get quick project setup recommendations
  So that I can start development efficiently with the right architecture

  Background:
    Given the AI advisor is available
    And the knowledge base is loaded

  Scenario: E-commerce API recommendation for senior team
    Given I have a project requirement for "api"
    And my project domain is "e-commerce"
    And my team experience level is "senior"
    When I request quick mode recommendations
    Then I should get blueprint recommendations including "web-api-ddd", "web-api-clean", "web-api"
    And I should get framework recommendations including "gin", "echo", "fiber"
    And the confidence level should be at least 0.6
    And I should receive reasoning for the recommendations
    And I should receive alternatives
    And the estimated file count should be between 20 and 100

  Scenario: CLI tool recommendation for junior team
    Given I have a project requirement for "cli"
    And my project domain is "devtools"
    And my team experience level is "junior"
    When I request quick mode recommendations
    Then I should get blueprint recommendations including "cli-simple", "cli"
    And I should get framework recommendations including "cobra"
    And the confidence level should be at least 0.5
    And I should receive reasoning for the recommendations
    And I should receive alternatives
    And the estimated file count should be between 5 and 35

  Scenario: Fintech API recommendation for expert team
    Given I have a project requirement for "api"
    And my project domain is "fintech"
    And my team experience level is "expert"
    When I request quick mode recommendations
    Then I should get blueprint recommendations including "web-api-ddd", "web-api-hexagonal"
    And I should get framework recommendations including "gin", "echo"
    And I should get logger recommendations including "logrus", "zap"
    And the confidence level should be at least 0.7
    And I should receive reasoning for the recommendations
    And I should receive alternatives
    And the estimated file count should be between 50 and 150

  Scenario: IoT Lambda recommendation for mixed team
    Given I have a project requirement for "lambda"
    And my project domain is "iot"
    And my team experience level is "mixed"
    When I request quick mode recommendations
    Then I should get blueprint recommendations including "lambda-event-processing", "lambda"
    And I should get logger recommendations including "slog", "zap"
    And the confidence level should be at least 0.5
    And I should receive reasoning for the recommendations
    And I should receive alternatives
    And the estimated file count should be between 10 and 50

  Scenario: Microservice recommendation for expert team
    Given I have a project requirement for "microservice"
    And my project domain is "e-commerce"
    And my team experience level is "expert"
    When I request quick mode recommendations
    Then I should get blueprint recommendations including "microservice"
    And I should get framework recommendations including "gin", "fiber"
    And the confidence level should be at least 0.6
    And I should receive reasoning for the recommendations
    And I should receive alternatives
    And the estimated file count should be between 40 and 80

  Scenario Outline: Quick recommendations for various project types
    Given I have a project requirement for "<project_type>"
    And my project domain is "<domain>"
    And my team experience level is "<team_experience>"
    When I request quick mode recommendations
    Then I should get a valid blueprint recommendation
    And the confidence level should be at least <min_confidence>
    And I should receive reasoning for the recommendations
    And the estimated file count should be greater than 0

    Examples:
      | project_type | domain       | team_experience | min_confidence |
      | api          | healthcare   | senior          | 0.6            |
      | cli          | automation   | junior          | 0.5            |
      | library      | utilities    | mixed           | 0.5            |
      | lambda       | data-processing | expert        | 0.6            |
      | microservice | logistics    | senior          | 0.6            |

  Scenario: Quick recommendations with performance requirements
    Given I have a project requirement for "api"
    And my project domain is "fintech"
    And my team experience level is "senior"
    And I expect "high" load
    And I need "fast" response times
    When I request quick mode recommendations
    Then I should get performance-optimized blueprint recommendations
    And I should get high-performance framework recommendations
    And I should get performance-focused logger recommendations
    And the confidence level should be at least 0.7
    And the reasoning should mention performance considerations

  Scenario: Quick recommendations should be fast
    Given I have a project requirement for "api"
    And my project domain is "e-commerce"
    And my team experience level is "mixed"
    When I request 100 quick mode recommendations
    Then all recommendations should be generated within 1 second
    And each recommendation should be valid and complete