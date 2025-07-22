Feature: Lambda Proxy Blueprint Generation
  As a developer
  I want to generate serverless Lambda proxy applications
  So that I can build AWS API Gateway integrated REST APIs

  Background:
    Given I have the go-starter CLI tool available
    And I am in a temporary working directory

  Scenario: Generate basic lambda-proxy with Gin framework
    Given I want to generate a lambda-proxy blueprint
    When I run the generator with:
      | parameter    | value                           |
      | type         | lambda-proxy                    |
      | name         | test-lambda-api                 |
      | module       | github.com/test/lambda-api      |
      | framework    | gin                             |
      | auth_type    | none                            |
      | logger_type  | slog                            |
    Then the generation should succeed
    And the generated project should have the following structure:
      | file                                      | type |
      | main.go                                   | file |
      | handler.go                                | file |
      | go.mod                                    | file |
      | template.yaml                             | file |
      | internal/config/config.go                 | file |
      | internal/handlers/health.go               | file |
      | internal/handlers/api.go                  | file |
      | internal/middleware/cors.go               | file |
      | internal/middleware/logging.go            | file |
      | internal/middleware/recovery.go           | file |
      | internal/observability/logger.go          | file |
      | internal/observability/metrics.go         | file |
      | internal/observability/tracing.go         | file |
      | scripts/deploy.sh                         | file |
      | scripts/local-dev.sh                      | file |
    And the project should compile successfully
    And the main.go should contain Gin framework imports
    And the SAM template should be valid

  Scenario: Generate lambda-proxy with JWT authentication
    Given I want to generate a lambda-proxy blueprint
    When I run the generator with:
      | parameter    | value                           |
      | type         | lambda-proxy                    |
      | name         | secure-lambda-api               |
      | module       | github.com/test/secure-api      |
      | framework    | echo                            |
      | auth_type    | jwt                             |
      | logger_type  | zap                             |
    Then the generation should succeed
    And the generated project should include authentication files:
      | file                                      |
      | internal/handlers/auth.go                 |
      | internal/middleware/auth.go               |
      | internal/services/auth.go                 |
    And the project should compile successfully
    And the auth middleware should contain JWT validation logic
    And the SAM template should include custom authorizers

  Scenario: Generate lambda-proxy with Cognito authentication
    Given I want to generate a lambda-proxy blueprint
    When I run the generator with:
      | parameter    | value                           |
      | type         | lambda-proxy                    |
      | name         | cognito-lambda-api              |
      | module       | github.com/test/cognito-api     |
      | framework    | fiber                           |
      | auth_type    | cognito                         |
      | logger_type  | logrus                          |
    Then the generation should succeed
    And the generated project should include Cognito integration:
      | file                                      |
      | internal/services/auth.go                 |
      | internal/middleware/auth.go               |
    And the project should compile successfully
    And the auth service should contain Cognito user pool validation
    And the SAM template should reference Cognito user pools

  Scenario: Generate lambda-proxy with all frameworks
    Given I want to generate a lambda-proxy blueprint
    When I run the generator for each framework:
      | framework |
      | gin       |
      | echo      |
      | fiber     |
      | chi       |
      | stdlib    |
    Then each generated project should:
      | validation                                |
      | compile successfully                      |
      | contain framework-specific imports        |
      | have consistent API structure             |
      | include proper middleware integration     |
      | generate valid SAM templates             |

  Scenario: Generate lambda-proxy with Terraform infrastructure
    Given I want to generate a lambda-proxy blueprint
    When I run the generator with:
      | parameter    | value                           |
      | type         | lambda-proxy                    |
      | name         | terraform-lambda-api            |
      | module       | github.com/test/terraform-api   |
      | framework    | gin                             |
    Then the generation should succeed
    And the generated project should include Terraform files:
      | file                                      |
      | terraform/main.tf                         |
      | terraform/variables.tf                    |
      | terraform/outputs.tf                      |
    And the Terraform configuration should be valid
    And the project should compile successfully

  Scenario: Generate lambda-proxy with observability stack
    Given I want to generate a lambda-proxy blueprint
    When I run the generator with:
      | parameter    | value                           |
      | type         | lambda-proxy                    |
      | name         | observable-lambda-api           |
      | module       | github.com/test/observable-api  |
      | framework    | echo                            |
      | logger_type  | zerolog                         |
    Then the generation should succeed
    And the generated project should include observability components:
      | file                                      |
      | internal/observability/logger.go          |
      | internal/observability/metrics.go         |
      | internal/observability/tracing.go         |
    And the project should compile successfully
    And the observability components should integrate with AWS services
    And the logger should use zerolog implementation

  Scenario: Generate lambda-proxy with CI/CD workflows
    Given I want to generate a lambda-proxy blueprint
    When I run the generator with:
      | parameter    | value                           |
      | type         | lambda-proxy                    |
      | name         | cicd-lambda-api                 |
      | module       | github.com/test/cicd-api        |
      | framework    | chi                             |
    Then the generation should succeed
    And the generated project should include CI/CD workflows:
      | file                                      |
      | .github/workflows/ci.yml                  |
      | .github/workflows/deploy.yml              |
      | .github/workflows/security.yml            |
      | .github/workflows/release.yml             |
    And the CI/CD workflows should be valid YAML
    And the workflows should include multi-platform builds
    And the security workflow should include comprehensive scanning

  Scenario: Validate lambda-proxy runtime architecture
    Given I have generated a lambda-proxy project
    When I analyze the generated architecture
    Then the project should demonstrate proper lambda proxy patterns:
      | pattern                                   |
      | API Gateway event processing              |
      | Framework adapter integration             |
      | Cold start optimization                   |
      | Error handling and recovery               |
      | Structured logging with correlation IDs   |
      | Health check endpoints                    |
      | CORS handling                             |
      | Request/response transformation           |

  Scenario: Validate lambda-proxy deployment readiness
    Given I have generated a lambda-proxy project with all components
    When I verify deployment readiness
    Then the project should be deployment-ready:
      | aspect                                    |
      | SAM template validates successfully       |
      | Terraform configuration is syntactically correct |
      | All Go code compiles without errors      |
      | Dependencies are properly managed         |
      | Deployment scripts are executable        |
      | Environment variables are documented     |
      | IAM permissions are least-privilege      |

  Scenario: Validate lambda-proxy performance characteristics
    Given I have generated a lambda-proxy project
    When I analyze performance characteristics
    Then the project should meet performance requirements:
      | requirement                               | target   |
      | Cold start time                          | < 500ms  |
      | Memory usage optimization                | enabled  |
      | Binary size optimization                 | enabled  |
      | Framework adapter efficiency            | optimal  |
      | Logging performance impact              | minimal  |

  Scenario Outline: Generate lambda-proxy with different logger types
    Given I want to generate a lambda-proxy blueprint
    When I run the generator with logger type "<logger>"
    Then the generated project should compile successfully
    And the logger implementation should be "<logger>"
    And the logging should be consistent across all components

    Examples:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |

  Scenario: Validate lambda-proxy security implementation
    Given I have generated a lambda-proxy project with JWT auth
    When I analyze the security implementation
    Then the project should demonstrate security best practices:
      | practice                                  |
      | Input validation on all endpoints        |
      | Proper error handling without leaks      |
      | Secure JWT token validation              |
      | API Gateway custom authorizer setup      |
      | CORS policy enforcement                   |
      | Request size limits                       |
      | Rate limiting consideration               |

  Scenario: Validate lambda-proxy multi-architecture build
    Given I have generated a lambda-proxy project
    When I build the project for multiple architectures
    Then the build should succeed for:
      | architecture |
      | linux/amd64  |
      | linux/arm64  |
    And the generated binaries should be optimized for Lambda runtime
    And the deployment packages should be correctly formatted