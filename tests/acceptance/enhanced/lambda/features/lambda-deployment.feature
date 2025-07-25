Feature: Lambda Deployment Scenario Testing
  As a go-starter user building serverless applications
  I want to ensure that Lambda functions are properly configured for deployment
  So that they can be successfully deployed and run on AWS Lambda

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And I am testing Lambda deployment scenarios

  @critical @lambda @basic @p1
  Scenario Outline: Basic Lambda function deployment readiness
    When I generate a Lambda function with configuration:
      | type       | lambda      |
      | framework  | <framework> |
      | logger     | <logger>    |
      | go_version | 1.23        |
    Then the project should compile successfully
    And the Lambda binary should be created
    And the binary should be optimized for Lambda runtime
    And cross-compilation should work correctly
    And AWS Lambda handler should be properly implemented
    And Lambda context handling should work correctly
    And CloudWatch logging should be integrated
    And <logger> logger should work with CloudWatch
    And SAM template should be properly configured
    And Lambda function resources should be defined
    And environment variables should be configurable
    And IAM roles should be properly configured
    And deployment configuration should be complete

    Examples:
      | framework | logger  |
      | none      | slog    |
      | none      | zap     |
      | none      | logrus  |
      | none      | zerolog |

  @critical @lambda @api-gateway @p1
  Scenario Outline: Lambda API Gateway proxy deployment
    When I generate a Lambda API Gateway proxy with configuration:
      | type       | lambda-proxy |
      | framework  | <framework>  |
      | logger     | <logger>     |
      | go_version | 1.23         |
    Then the project should compile successfully
    And the Lambda binary should be created
    And AWS Lambda handler should be properly implemented
    And API Gateway integration should work correctly
    And request/response mapping should be implemented
    And CORS headers should be properly configured
    And error responses should follow API Gateway format
    And <logger> logger should work with CloudWatch
    And SAM template should be properly configured
    And API Gateway resources should be defined in SAM
    And deployment configuration should be complete

    Examples:
      | framework | logger |
      | none      | slog   |
      | none      | zap    |

  @performance @lambda @cold-start @p1
  Scenario Outline: Lambda cold start optimization
    When I generate a Lambda function optimized for cold starts:
      | type       | lambda   |
      | framework  | none     |
      | logger     | <logger> |
      | go_version | 1.23     |
    Then the project should compile successfully
    And the Lambda binary should be created
    And cold start optimization should be implemented
    And minimal dependencies should be included
    And binary size should be optimized
    And memory usage should be efficient
    And the binary should be optimized for Lambda runtime
    And <logger> logger should work with CloudWatch

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | zerolog |

  @integration @lambda @testing @p1
  Scenario Outline: Lambda local testing and development
    When I generate a Lambda function with configuration:
      | type       | lambda   |
      | framework  | none     |
      | logger     | <logger> |
      | go_version | 1.23     |
    Then the project should compile successfully
    And local testing should be supported
    And unit tests should cover Lambda handlers
    And SAM local testing should work
    And event samples should be provided
    And Makefile should have deployment targets
    And deployment documentation should be comprehensive

    Examples:
      | logger |
      | slog   |
      | zap    |

  @integration @lambda @aws-sdk @p1
  Scenario Outline: Lambda with AWS SDK integration
    When I generate a Lambda function with AWS SDK integration:
      | type       | lambda   |
      | framework  | none     |
      | logger     | <logger> |
      | go_version | 1.23     |
    Then the project should compile successfully
    And the Lambda binary should be created
    And AWS SDK v2 should be properly integrated
    And service clients should be properly initialized
    And AWS credentials handling should be secure
    And X-Ray tracing should be available
    And <logger> logger should work with CloudWatch
    And error handling should be comprehensive

    Examples:
      | logger |
      | slog   |
      | zap    |

  @deployment @lambda @cicd @p1
  Scenario Outline: Lambda deployment automation
    When I generate a Lambda function with configuration:
      | type       | lambda   |
      | framework  | none     |
      | logger     | slog     |
      | go_version | 1.23     |
    Then the project should compile successfully
    And deployment scripts should be included
    And Makefile should have deployment targets
    And CI/CD templates should be available
    And SAM template should be properly configured
    And deployment documentation should be comprehensive
    And multiple environments should be supported

  @monitoring @lambda @observability @p1
  Scenario Outline: Lambda monitoring and observability
    When I generate a Lambda function with configuration:
      | type       | lambda   |
      | framework  | none     |
      | logger     | <logger> |
      | go_version | 1.23     |
    Then the project should compile successfully
    And CloudWatch logging should be integrated
    And <logger> logger should work with CloudWatch
    And structured logging should be used
    And error handling should be comprehensive
    And panic recovery should be implemented
    And metrics should be exported to CloudWatch
    And X-Ray tracing should be available

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  @security @lambda @configuration @p1
  Scenario: Lambda security and configuration management
    When I generate a Lambda function with configuration:
      | type       | lambda |
      | framework  | none   |
      | logger     | slog   |
      | go_version | 1.23   |
    Then the project should compile successfully
    And environment variables should be configurable
    And IAM roles should be properly configured
    And AWS credentials handling should be secure
    And secrets management should be implemented
    And parameter store integration should work
    And environment-specific configs should exist

  @integration @lambda @complete @p1
  Scenario Outline: Complete Lambda deployment validation
    When I generate a Lambda function with configuration:
      | type       | <type>   |
      | framework  | none     |
      | logger     | <logger> |
      | go_version | 1.23     |
    Then the project should compile successfully
    And the Lambda binary should be created
    And the binary should be optimized for Lambda runtime
    And AWS Lambda handler should be properly implemented
    And Lambda context handling should work correctly
    And CloudWatch logging should be integrated
    And <logger> logger should work with CloudWatch
    And SAM template should be properly configured
    And deployment configuration should be complete
    And Makefile should have deployment targets
    And error handling should be comprehensive
    And deployment documentation should be comprehensive

    Examples:
      | type         | logger  |
      | lambda       | slog    |
      | lambda       | zap     |
      | lambda       | logrus  |
      | lambda       | zerolog |
      | lambda-proxy | slog    |
      | lambda-proxy | zap     |

  @cross-platform @lambda @build @p1
  Scenario: Lambda cross-platform build validation
    When I generate a Lambda function with configuration:
      | type       | lambda |
      | framework  | none   |
      | logger     | slog   |
      | go_version | 1.23   |
    Then the project should compile successfully
    And cross-compilation should work correctly
    And the Lambda binary should be created
    And the binary should be optimized for Lambda runtime
    And Makefile should have deployment targets
    And build process should support multiple architectures