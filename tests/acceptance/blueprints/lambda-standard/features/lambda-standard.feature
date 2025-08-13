Feature: Lambda Standard Blueprint Generation
  As a Go developer
  I want to generate standard AWS Lambda functions
  So that I can build serverless applications with proper monitoring and logging

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate standard Lambda function
    Given I want to create a standard AWS Lambda function
    When I run the command "go-starter new my-lambda --type=lambda-standard --logger=slog --no-git"
    Then the generation should succeed
    And the project should contain Lambda-specific components
    And the generated code should compile successfully
    And the project should include AWS Lambda runtime

  Scenario: Lambda CloudWatch integration
    Given I have generated a standard Lambda function
    When I examine the observability configuration
    Then CloudWatch logging should be properly configured
    And metrics collection should be implemented
    And distributed tracing should be available
    And error tracking should be integrated

  Scenario: Lambda deployment configuration
    Given I have generated a standard Lambda function
    When I examine the deployment setup
    Then the deployment script should be available
    And SAM template should be properly configured
    And environment variables should be managed
    And IAM permissions should be defined

  Scenario: Lambda handler implementation
    Given I have generated a standard Lambda function
    When I examine the handler code
    Then the main handler should be properly structured
    And context handling should be implemented
    And error handling should be robust
    And response formatting should be correct

  Scenario: Lambda testing infrastructure
    Given I have generated a standard Lambda function
    When I examine the test setup
    Then unit tests should be included
    And integration tests should be available
    And local testing should be supported
    And test coverage should be measurable

  Scenario: Lambda performance optimization
    Given I have generated a standard Lambda function
    When I examine the performance configuration
    Then cold start optimization should be implemented
    And memory usage should be optimized
    And initialization should be efficient
    And runtime performance should be monitored