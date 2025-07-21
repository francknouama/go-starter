Feature: Service Blueprint Generation
  As a distributed systems developer
  I want to generate modern, production-ready service applications
  So that I can quickly build microservices, serverless functions, and distributed systems with industry best practices

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory
    And I have Docker available for containerization testing

  Scenario: Generate microservice-standard application
    Given I want to create a gRPC-based microservice
    When I run the command "go-starter new my-service --type=microservice-standard --framework=grpc --database-driver=postgres --no-git"
    Then the generation should succeed
    And the project should contain all essential microservice components
    And the generated code should compile successfully
    And the service should support gRPC and HTTP interfaces
    And the microservice should be container-ready with health checks

  Scenario: Generate lambda-standard function
    Given I want to create an AWS Lambda serverless function
    When I run the command "go-starter new my-lambda --type=lambda-standard --runtime=aws --logger=cloudwatch --no-git"
    Then the generation should succeed
    And the project should contain all essential serverless components
    And the generated code should compile successfully
    And the lambda should support AWS SDK integration
    And the function should include observability features

  Scenario: Microservice with different communication protocols
    Given I want to create microservices with various protocols
    When I generate a microservice with protocol "<protocol>"
    Then the project should support the "<protocol>" communication pattern
    And the service should include appropriate client libraries
    And the protocol-specific middleware should be configured
    And the service should compile and run correctly

    Examples:
      | protocol |
      | grpc     |
      | rest     |
      | graphql  |

  Scenario: Service with different database integrations
    Given I want to create services with database persistence
    When I generate a service with database "<database>" and ORM "<orm>"
    Then the project should include database configuration
    And the service should support database migrations
    And the repository layer should use the specified ORM
    And the database connection should be testable with containers
    And the service should handle connection pooling

    Examples:
      | database  | orm   |
      | postgres  | gorm  |
      | postgres  | sqlx  |
      | mysql     | gorm  |
      | mongodb   | mongo |
      | redis     | redis |

  Scenario: Microservice with service discovery
    Given I want microservices that can discover each other
    When I generate a microservice with service discovery
    Then the service should include service registry integration
    And the service should support health check endpoints
    And the service should handle service registration
    And the service should support load balancing
    And the discovery configuration should be environment-aware

  Scenario: Service mesh integration
    Given I want microservices ready for service mesh deployment
    When I generate a microservice with service mesh support
    Then the project should include Istio configuration
    And the service should support mTLS communication
    And the service should include traffic management policies
    And the service should support distributed tracing
    And the mesh configuration should be production-ready

  Scenario: Circuit breaker and resilience patterns
    Given I want resilient microservices
    When I generate a microservice with resilience patterns
    Then the service should include circuit breaker middleware
    And the service should support retry mechanisms
    And the service should handle timeout configurations
    And the service should include rate limiting
    And the resilience patterns should be configurable

  Scenario: Observability and monitoring integration
    Given I want observable microservices
    When I generate a microservice with observability features
    Then the service should include OpenTelemetry integration
    And the service should support Prometheus metrics
    And the service should include distributed tracing
    And the service should support structured logging
    And the observability should cover all service layers

  Scenario: Authentication and authorization
    Given I want secure microservices
    When I generate a microservice with authentication type "<auth_type>"
    Then the service should include authentication middleware
    And the service should support JWT token validation
    And the service should include authorization policies
    And the service should handle authentication errors gracefully
    And the security configuration should be production-ready

    Examples:
      | auth_type |
      | jwt       |
      | oauth2    |
      | mtls      |
      | api-key   |

  Scenario: Container and Kubernetes deployment
    Given I want to deploy microservices in Kubernetes
    When I generate a microservice with Kubernetes support
    Then the project should include Dockerfile with multi-stage builds
    And the project should include Kubernetes manifests
    And the service should support horizontal pod autoscaling
    And the service should include resource limits and requests
    And the deployment should support rolling updates

  Scenario: AWS Lambda with different triggers
    Given I want Lambda functions with various triggers
    When I generate a Lambda with trigger type "<trigger>"
    Then the function should support the "<trigger>" event source
    And the function should include appropriate IAM permissions
    And the function should handle the event format correctly
    And the function should include error handling
    And the deployment should support CloudFormation/SAM

    Examples:
      | trigger    |
      | api-gateway|
      | s3         |
      | dynamodb   |
      | sqs        |
      | sns        |
      | eventbridge|

  Scenario: Serverless cold start optimization
    Given I want optimized Lambda functions
    When I generate a Lambda with performance optimizations
    Then the function should minimize cold start time
    And the function should include connection pooling
    And the function should support provisioned concurrency
    And the function should include memory optimization
    And the function should support custom runtime configurations

  Scenario: Multi-environment deployment strategies
    Given I want services deployable across environments
    When I generate a service with multi-environment support
    Then the service should support environment-specific configurations
    And the service should include staging and production configs
    And the service should support feature flags
    And the service should include blue-green deployment strategies
    And the configuration should support secrets management

  Scenario: API Gateway and Lambda integration
    Given I want serverless APIs
    When I generate a Lambda with API Gateway integration
    Then the function should support REST API patterns
    And the function should include request/response transformation
    And the function should support API versioning
    And the function should include CORS configuration
    And the API should support OpenAPI documentation

  Scenario: Event-driven architecture patterns
    Given I want event-driven microservices
    When I generate services with event sourcing patterns
    Then the services should support event streaming
    And the services should include event store integration
    And the services should support CQRS patterns
    And the services should handle eventual consistency
    And the services should include event replay capabilities

  Scenario: Performance and load testing
    Given I want performant microservices
    When I generate a microservice with performance testing
    Then the service should include load testing scripts
    And the service should support performance benchmarking
    And the service should include resource monitoring
    And the service should handle high-throughput scenarios
    And the performance tests should validate SLA requirements

  Scenario: Security scanning and compliance
    Given I want secure service deployments
    When I generate a service with security scanning
    Then the project should include vulnerability scanning
    And the service should support security compliance checks
    And the service should include dependency scanning
    And the service should support container image scanning
    And the security pipeline should prevent vulnerable deployments

  Scenario: Distributed tracing across services
    Given I want to trace requests across service boundaries
    When I generate microservices with distributed tracing
    Then the services should include trace correlation
    And the services should support trace propagation
    And the services should include span annotations
    And the services should integrate with tracing backends
    And the tracing should support sampling strategies

  Scenario: gRPC Gateway REST bridge
    Given I want to expose gRPC services via REST
    When I generate a service with gRPC Gateway
    Then the service should support dual gRPC/REST interfaces
    And the service should include Protocol Buffer definitions
    And the service should support automatic REST mapping
    And the service should include OpenAPI documentation
    And the gateway should handle protocol translation

  Scenario: Message queues and async communication
    Given I want asynchronous microservice communication
    When I generate services with message queue integration
    Then the services should support message publishing
    And the services should include message consumption
    And the services should handle message retry logic
    And the services should support dead letter queues
    And the messaging should be fault-tolerant

  Scenario: Database sharding and partitioning
    Given I want scalable data persistence
    When I generate services with database sharding
    Then the services should support horizontal partitioning
    And the services should include shard key management
    And the services should handle cross-shard transactions
    And the services should support read replicas
    And the data distribution should be configurable

  Scenario: Microservice testing strategies
    Given I want comprehensive service testing
    When I generate a microservice with testing infrastructure
    Then the project should include unit tests
    And the project should include integration tests
    And the project should include contract tests
    And the project should include end-to-end tests
    And the tests should use testcontainers for dependencies

  Scenario: Chaos engineering and resilience testing
    Given I want resilient distributed systems
    When I generate services with chaos engineering
    Then the services should include fault injection
    And the services should support failure scenarios
    And the services should include resilience testing
    And the services should handle network partitions
    And the chaos tests should validate system recovery

  Scenario: Configuration management and secrets
    Given I want secure configuration management
    When I generate services with advanced configuration
    Then the services should support external config sources
    And the services should include secrets management
    And the services should support configuration hot-reloading
    And the services should validate configuration schemas
    And the configuration should support environment overrides

  Scenario: Multi-region deployment patterns
    Given I want globally distributed services
    When I generate services for multi-region deployment
    Then the services should support region-aware routing
    And the services should include data replication strategies
    And the services should handle cross-region latency
    And the services should support disaster recovery
    And the deployment should minimize regional dependencies

  Scenario: Streaming data processing
    Given I want real-time data processing services
    When I generate services with streaming capabilities
    Then the services should support stream processing
    And the services should include windowing operations
    And the services should handle backpressure
    And the services should support exactly-once processing
    And the streaming should be horizontally scalable

  Scenario: Service versioning and backward compatibility
    Given I want evolvable service APIs
    When I generate services with versioning strategies
    Then the services should support API versioning
    And the services should maintain backward compatibility
    And the services should include deprecation handling
    And the services should support gradual rollouts
    And the versioning should include migration paths

  Scenario: Cost optimization and resource management
    Given I want cost-effective service deployments
    When I generate services with cost optimization
    Then the services should include resource monitoring
    And the services should support auto-scaling policies
    And the services should include cost tracking
    And the services should optimize resource utilization
    And the deployment should support spot instances

  Scenario: Compliance and audit logging
    Given I want compliant service architectures
    When I generate services with compliance features
    Then the services should include audit logging
    And the services should support compliance frameworks
    And the services should include data privacy controls
    And the services should support regulatory reporting
    And the compliance should be automated and verifiable