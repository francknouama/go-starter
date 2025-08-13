Feature: gRPC Pure Service Blueprint Generation
  As a Go developer
  I want to generate pure gRPC services with advanced features
  So that I can build high-performance microservices with comprehensive observability

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate basic gRPC service
    Given I want to create a gRPC service
    When I run "go-starter new my-grpc --type=grpc-pure --module=github.com/example/my-grpc --no-git"
    Then the generation should succeed
    And the project should contain gRPC-specific components
    And the generated code should compile successfully
    And protocol buffers should be syntactically valid

  Scenario: Generate gRPC service with JWT authentication
    Given I want to create a gRPC service with authentication
    When I run "go-starter new secure-grpc --type=grpc-pure --auth-type=jwt --module=github.com/example/secure-grpc --no-git"
    Then the generation should succeed
    And JWT authentication interceptors should be included
    And auth-related dependencies should be present in go.mod
    And the gRPC server should include auth middleware
    And the project should compile successfully

  Scenario: Generate gRPC service with observability
    Given I want a gRPC service with full observability
    When I run "go-starter new obs-grpc --type=grpc-pure --tracing=true --metrics=true --module=github.com/example/obs-grpc --no-git"
    Then the generation should succeed
    And OpenTelemetry tracing should be configured
    And Prometheus metrics should be included
    And observability middleware should be present
    And the project should compile successfully

  Scenario: Generate gRPC service with service discovery
    Given I want a gRPC service with service discovery
    When I run "go-starter new discovery-grpc --type=grpc-pure --service-discovery=consul --module=github.com/example/discovery-grpc --no-git"
    Then the generation should succeed
    And Consul service discovery should be configured
    And service registration should be implemented
    And health checking should be integrated
    And the project should compile successfully

  Scenario: Generate gRPC service with database integration
    Given I want a gRPC service with database support
    When I run "go-starter new db-grpc --type=grpc-pure --database-driver=postgres --database-orm=gorm --module=github.com/example/db-grpc --no-git"
    Then the generation should succeed
    And PostgreSQL database configuration should be included
    And GORM integration should be present
    And database migrations should be included
    And the project should compile successfully

  Scenario: Generate gRPC service with multiple loggers
    Given I want to test different logging implementations
    When I run "go-starter new zap-grpc --type=grpc-pure --logger=zap --module=github.com/example/zap-grpc --no-git"
    Then the generation should succeed
    And Zap logger implementation should be included
    And logger dependencies should be present in go.mod
    And the project should compile successfully

  Scenario: Generate gRPC service with mTLS authentication
    Given I want a gRPC service with mutual TLS
    When I run "go-starter new mtls-grpc --type=grpc-pure --auth-type=mtls --module=github.com/example/mtls-grpc --no-git"
    Then the generation should succeed
    And mTLS authentication should be configured
    And TLS certificate handling should be included
    And the project should compile successfully

  Scenario: Protocol buffer validation
    Given I have generated a gRPC service
    When I examine the protocol buffer files
    Then all proto files should be syntactically correct
    And proto imports should be valid
    And service definitions should include proper annotations
    And message validation rules should be present

  Scenario: gRPC interceptor chain validation
    Given I have generated a gRPC service with interceptors
    When I examine the server configuration
    Then the interceptor chain should be properly ordered
    And auth interceptors should validate JWT tokens
    And logging interceptors should capture request/response data
    And metrics interceptors should track performance
    And recovery interceptors should handle panics gracefully

  Scenario: Health check integration
    Given I have generated a gRPC service
    When I examine the health check configuration
    Then gRPC health probe should be implemented
    And health status should be configurable
    And dependency health checks should be included
    And the health service should be registered

  Scenario: Load testing framework
    Given I have generated a gRPC service
    When I examine the testing infrastructure
    Then load testing scripts should be included
    And performance benchmarks should be available
    And concurrent client testing should be supported
    And metrics collection should be automated

  Scenario: Complete integration with all features
    Given I want a production-ready gRPC service
    When I run "go-starter new full-grpc --type=grpc-pure --auth-type=jwt --tracing=true --metrics=true --service-discovery=consul --database-driver=postgres --database-orm=gorm --logger=zap --module=github.com/example/full-grpc --no-git"
    Then the generation should succeed
    And all advanced features should be properly integrated
    And no feature conflicts should exist
    And the project should compile successfully
    And all tests should pass
    And the service should be production-ready

  Scenario: gRPC reflection support for development
    Given I want a gRPC service with development tools
    When I run "go-starter new dev-grpc --type=grpc-pure --reflection=true --module=github.com/example/dev-grpc --no-git"
    Then the generation should succeed
    And gRPC reflection should be enabled
    And development tools should be included
    And the project should compile successfully

  Scenario: Cross-platform compatibility testing
    Given I have generated a gRPC service
    When I test on different platforms
    Then the service should compile on Linux
    And the service should compile on macOS
    And the service should compile on Windows
    And all dependencies should be cross-platform compatible

  Scenario: Docker and Kubernetes deployment readiness
    Given I have generated a gRPC service
    When I examine the deployment configuration
    Then Dockerfile should be optimized for gRPC services
    And Kubernetes manifests should be included
    And Health check endpoints should be properly configured
    And Resource limits should be sensible defaults
    And the container should build successfully

  Scenario: Performance and scalability validation
    Given I have generated a gRPC service with metrics
    When I run performance tests
    Then the service should handle concurrent connections
    And metrics should be collected accurately
    And memory usage should be within acceptable limits
    And response times should meet performance criteria