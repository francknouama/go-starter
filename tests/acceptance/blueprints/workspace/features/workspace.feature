Feature: Workspace Blueprint Generation
  As a developer
  I want to generate Go multi-module workspace projects
  So that I can build monorepo applications with shared libraries and services

  Background:
    Given I have the go-starter CLI tool available
    And I am in a temporary working directory

  Scenario: Generate basic workspace with all modules enabled
    Given I want to generate a workspace blueprint
    When I run the generator with:
      | parameter               | value                           |
      | type                    | workspace                       |
      | name                    | test-workspace                  |
      | module                  | github.com/test/workspace       |
      | go_version              | 1.21                            |
      | framework               | gin                             |
      | database_type           | postgres                        |
      | message_queue           | redis                           |
      | logger_type             | slog                            |
      | enable_web_api          | true                            |
      | enable_cli              | true                            |
      | enable_worker           | true                            |
      | enable_microservices    | true                            |
      | enable_docker           | true                            |
      | enable_kubernetes       | true                            |
    Then the generation should succeed
    And the generated workspace should have the Go workspace structure:
      | file                                      | type      |
      | go.work                                   | file      |
      | workspace.yaml                            | file      |
      | Makefile                                  | file      |
      | pkg/shared/                               | directory |
      | pkg/models/                               | directory |
      | pkg/storage/                              | directory |
      | pkg/events/                               | directory |
      | cmd/api/                                  | directory |
      | cmd/cli/                                  | directory |
      | cmd/worker/                               | directory |
      | services/user-service/                    | directory |
      | services/notification-service/            | directory |
    And all modules should compile successfully
    And the workspace should sync without errors

  Scenario: Generate workspace with minimal configuration
    Given I want to generate a workspace blueprint
    When I run the generator with:
      | parameter               | value                           |
      | type                    | workspace                       |
      | name                    | minimal-workspace               |
      | module                  | github.com/test/minimal         |
      | enable_web_api          | false                           |
      | enable_cli              | true                            |
      | enable_worker           | false                           |
      | enable_microservices    | false                           |
      | database_type           | none                            |
      | message_queue           | none                            |
      | enable_docker           | false                           |
      | enable_kubernetes       | false                           |
    Then the generation should succeed
    And the generated workspace should have minimal structure:
      | file                                      | type      |
      | go.work                                   | file      |
      | pkg/shared/                               | directory |
      | pkg/models/                               | directory |
      | cmd/cli/                                  | directory |
    And the workspace should not include:
      | file                                      |
      | pkg/storage/                              |
      | pkg/events/                               |
      | cmd/api/                                  |
      | cmd/worker/                               |
      | services/                                 |
      | docker-compose.yml                        |
      | deployments/k8s/                          |
    And all modules should compile successfully

  Scenario: Generate workspace with different frameworks
    Given I want to generate a workspace blueprint
    When I run the generator for each framework:
      | framework |
      | gin       |
      | echo      |
      | fiber     |
      | chi       |
    Then each generated workspace should:
      | validation                                |
      | compile successfully                      |
      | contain framework-specific imports        |
      | have consistent module structure          |
      | include proper dependency management      |

  Scenario: Generate workspace with different database types
    Given I want to generate a workspace blueprint
    When I run the generator for each database:
      | database  |
      | postgres  |
      | mysql     |
      | mongodb   |
      | sqlite    |
    Then each generated workspace should:
      | validation                                |
      | include database-specific storage package |
      | contain appropriate database drivers      |
      | have correct connection configuration     |
      | compile successfully                      |

  Scenario: Generate workspace with different message queues
    Given I want to generate a workspace blueprint
    When I run the generator for each message queue:
      | message_queue |
      | redis         |
      | nats          |
      | kafka         |
      | rabbitmq      |
    Then each generated workspace should:
      | validation                                |
      | include message queue events package     |
      | contain appropriate MQ client libraries  |
      | have correct connection configuration     |
      | compile successfully                      |

  Scenario: Generate workspace with different logger types
    Given I want to generate a workspace blueprint
    When I run the generator for each logger:
      | logger  |
      | slog    |
      | zap     |
      | logrus  |
      | zerolog |
    Then each generated workspace should:
      | validation                                |
      | include logger-specific implementation    |
      | have consistent logging interface         |
      | compile successfully                      |

  Scenario: Generate workspace with Docker support
    Given I want to generate a workspace blueprint
    When I run the generator with:
      | parameter               | value                           |
      | type                    | workspace                       |
      | name                    | docker-workspace                |
      | module                  | github.com/test/docker          |
      | enable_docker           | true                            |
      | enable_web_api          | true                            |
      | enable_worker           | true                            |
      | enable_microservices    | true                            |
    Then the generation should succeed
    And the generated workspace should include Docker files:
      | file                                      |
      | docker-compose.yml                        |
      | docker-compose.dev.yml                    |
      | cmd/api/Dockerfile                        |
      | cmd/worker/Dockerfile                     |
      | services/user-service/Dockerfile          |
      | services/notification-service/Dockerfile  |
    And the Docker configurations should be valid

  Scenario: Generate workspace with Kubernetes support
    Given I want to generate a workspace blueprint
    When I run the generator with:
      | parameter               | value                           |
      | type                    | workspace                       |
      | name                    | k8s-workspace                   |
      | module                  | github.com/test/k8s             |
      | enable_kubernetes       | true                            |
      | enable_web_api          | true                            |
      | enable_worker           | true                            |
      | enable_microservices    | true                            |
    Then the generation should succeed
    And the generated workspace should include Kubernetes manifests:
      | file                                              |
      | deployments/k8s/namespace.yaml                    |
      | deployments/k8s/configmap.yaml                    |
      | deployments/k8s/secrets.yaml                      |
      | deployments/k8s/api-deployment.yaml               |
      | deployments/k8s/worker-deployment.yaml            |
      | deployments/k8s/user-service-deployment.yaml      |
      | deployments/k8s/notification-service-deployment.yaml |
    And the Kubernetes manifests should be valid YAML

  Scenario: Validate workspace build system
    Given I have generated a complete workspace
    When I test the build system
    Then the Makefile should provide all expected targets:
      | target          | description                    |
      | build-all       | Build all modules              |
      | test-all        | Run all tests                  |
      | lint-all        | Lint all modules               |
      | clean-all       | Clean build artifacts          |
      | deps-update     | Update dependencies            |
      | docker-build    | Build Docker images            |
      | k8s-deploy      | Deploy to Kubernetes           |
    And all build targets should execute successfully
    And the build should respect dependency order

  Scenario: Validate workspace module dependencies
    Given I have generated a workspace with all modules
    When I analyze module dependencies
    Then the dependency graph should be correct:
      | dependent_module                | dependencies                    |
      | pkg/models                      | pkg/shared                      |
      | pkg/storage                     | pkg/shared, pkg/models          |
      | pkg/events                      | pkg/shared, pkg/models          |
      | cmd/api                         | pkg/shared, pkg/models, pkg/storage |
      | cmd/cli                         | pkg/shared, pkg/models          |
      | cmd/worker                      | pkg/shared, pkg/models, pkg/events |
      | services/user-service           | pkg/shared, pkg/models, pkg/storage |
      | services/notification-service   | pkg/shared, pkg/models, pkg/events |
    And there should be no circular dependencies
    And all module replacements should be correct

  Scenario: Validate workspace testing infrastructure
    Given I have generated a workspace with testing enabled
    When I run the workspace tests
    Then the test infrastructure should be comprehensive:
      | component                               |
      | Unit tests for all shared packages      |
      | Integration tests for services          |
      | End-to-end tests for API endpoints      |
      | Test helpers and utilities              |
      | Coverage reporting                      |
    And all tests should pass
    And code coverage should meet minimum thresholds

  Scenario: Validate workspace CI/CD integration
    Given I have generated a workspace
    When I examine the CI/CD configuration
    Then the workspace should include CI/CD workflows:
      | file                                      |
      | .github/workflows/ci.yml                  |
      | .github/workflows/release.yml             |
      | .github/workflows/security.yml            |
    And the workflows should handle multi-module projects:
      | feature                                   |
      | Build all modules in dependency order     |
      | Run tests for all modules                 |
      | Generate combined coverage report         |
      | Security scanning for all modules         |
      | Multi-platform builds                     |
      | Release management                        |

  Scenario: Validate workspace observability integration
    Given I have generated a workspace with observability enabled
    When I analyze the observability stack
    Then the workspace should include observability components:
      | component                               |
      | Structured logging across all modules   |
      | Metrics collection with Prometheus      |
      | Distributed tracing with OpenTelemetry  |
      | Health checks for all services          |
      | Correlation IDs for request tracking    |
    And the observability should be consistent across modules

  Scenario: Validate workspace configuration management
    Given I have generated a workspace
    When I examine the configuration system
    Then the workspace should provide centralized configuration:
      | feature                                   |
      | Shared configuration utilities            |
      | Environment-specific configuration       |
      | Configuration validation                  |
      | Secret management                         |
      | Configuration reloading                   |
    And configuration should be consistent across modules

  Scenario: Validate workspace development tools
    Given I have generated a workspace
    When I examine the development tools
    Then the workspace should include development utilities:
      | tool                                      |
      | Build scripts for all modules             |
      | Test scripts with coverage                |
      | Linting and formatting tools              |
      | Dependency management scripts             |
      | Local development setup                   |
      | Debug helpers                             |
    And all tools should be properly configured

  Scenario: Validate workspace scalability patterns
    Given I have generated a workspace with microservices
    When I analyze the architecture
    Then the workspace should demonstrate scalability patterns:
      | pattern                                   |
      | Service independence                      |
      | Shared library reuse                      |
      | Event-driven communication               |
      | Database per service                      |
      | Configuration externalization            |
      | Health check endpoints                    |
      | Graceful shutdown handling               |

  Scenario: Validate workspace security best practices
    Given I have generated a workspace
    When I analyze the security implementation
    Then the workspace should demonstrate security best practices:
      | practice                                  |
      | Input validation across all services     |
      | Secure configuration management           |
      | Dependency vulnerability scanning         |
      | Container security practices              |
      | Network security policies                 |
      | Secrets management                        |
      | Audit logging                            |

  Scenario: Validate workspace performance optimization
    Given I have generated a workspace
    When I analyze the performance characteristics
    Then the workspace should include performance optimizations:
      | optimization                              |
      | Efficient build processes                 |
      | Optimized Docker images                   |
      | Connection pooling                        |
      | Caching strategies                        |
      | Resource management                       |
      | Monitoring and profiling tools            |