Feature: gRPC Gateway Blueprint Generation
  As a developer building gRPC services with REST API compatibility
  I want to generate gRPC Gateway projects that expose gRPC services via REST
  So that I can serve both gRPC and HTTP clients from a single service

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate basic grpc-gateway project
    Given I want to create a gRPC service with REST gateway
    When I run the command "go-starter new my-grpc-gateway --type=grpc-gateway --module=github.com/example/grpc-gateway --no-git"
    Then the generation should succeed
    And the project should contain all essential grpc-gateway components
    And the generated code should compile successfully
    And the service should support both gRPC and REST interfaces
    And the gateway should include Protocol Buffer definitions

  Scenario: Gateway with different service definitions
    Given I want to create gRPC gateways with various service types
    When I generate a grpc-gateway with service type "<service_type>"
    Then the project should include appropriate protobuf definitions
    And the service should support the "<service_type>" pattern
    And the gateway should map gRPC methods to REST endpoints
    And the OpenAPI documentation should be generated
    And the client code should be generated for both protocols

    Examples:
      | service_type |
      | user-service |
      | order-service |
      | notification-service |
      | auth-service |

  Scenario: Gateway with different databases
    Given I want gRPC gateways with data persistence
    When I generate a grpc-gateway with database "<database>" and ORM "<orm>"
    Then the project should include database configuration
    And the service should support database migrations
    And the gRPC service should include CRUD operations
    And the REST endpoints should map to database operations
    And the repository layer should use the specified ORM

    Examples:
      | database  | orm   |
      | postgres  | gorm  |
      | postgres  | sqlx  |
      | mysql     | gorm  |
      | mongodb   | mongo |

  Scenario: Gateway with authentication
    Given I want secure gRPC gateways
    When I generate a grpc-gateway with authentication type "<auth_type>"
    Then the service should include gRPC authentication interceptors
    And the gateway should include HTTP authentication middleware
    And the service should support JWT token validation
    And the authentication should work for both gRPC and REST
    And the security configuration should be production-ready

    Examples:
      | auth_type |
      | jwt       |
      | oauth2    |
      | api-key   |
      | mtls      |

  Scenario: Gateway with observability features
    Given I want observable gRPC gateways
    When I generate a grpc-gateway with observability enabled
    Then the service should include gRPC tracing interceptors
    And the gateway should include HTTP tracing middleware
    And the service should support Prometheus metrics
    And the service should include structured logging
    And the observability should cover both protocols
    And the metrics should distinguish between gRPC and REST requests

  Scenario: Gateway with rate limiting
    Given I want rate-limited gRPC gateways
    When I generate a grpc-gateway with rate limiting enabled
    Then the service should include gRPC rate limiting interceptors
    And the gateway should include HTTP rate limiting middleware
    And the rate limits should be configurable per method
    And the rate limiting should work for both protocols
    And the service should return appropriate error responses

  Scenario: Gateway with validation
    Given I want input validation for gRPC gateways
    When I generate a grpc-gateway with validation enabled
    Then the protobuf definitions should include validation rules
    And the gRPC service should validate input messages
    And the REST gateway should validate HTTP requests
    And the validation errors should be properly formatted
    And the validation should be consistent across protocols

  Scenario: Gateway with multiple services
    Given I want gRPC gateways handling multiple services
    When I generate a grpc-gateway with multiple services
    Then the project should include multiple protobuf service definitions
    And the gateway should route to appropriate gRPC services
    And the REST API should include all service endpoints
    And the OpenAPI documentation should cover all services
    And the service discovery should support multiple backends

  Scenario: Gateway with streaming support
    Given I want gRPC gateways with streaming capabilities
    When I generate a grpc-gateway with streaming enabled
    Then the service should support server-side streaming
    And the service should support client-side streaming
    And the service should support bidirectional streaming
    And the gateway should handle streaming over HTTP
    And the streaming should include proper error handling

  Scenario: Gateway with custom HTTP mapping
    Given I want customized REST API design
    When I generate a grpc-gateway with custom HTTP mapping
    Then the protobuf definitions should include HTTP annotations
    And the REST endpoints should follow custom URL patterns
    And the HTTP methods should map correctly to gRPC methods
    And the request/response transformation should be configured
    And the OpenAPI specification should reflect custom mapping

  Scenario: Gateway with middleware chain
    Given I want extensible gRPC gateways
    When I generate a grpc-gateway with middleware support
    Then the service should include gRPC interceptor chain
    And the gateway should include HTTP middleware chain
    And the middleware should be configurable
    And the service should support custom interceptors
    And the middleware should be properly ordered

  Scenario: Gateway with health checks
    Given I want production-ready gRPC gateways
    When I generate a grpc-gateway with health checks
    Then the service should include gRPC health service
    And the gateway should include HTTP health endpoints
    And the health checks should verify database connectivity
    And the health checks should verify external dependencies
    And the service should support Kubernetes health probes

  Scenario: Gateway with load balancing
    Given I want scalable gRPC gateways
    When I generate a grpc-gateway with load balancing
    Then the service should support client-side load balancing
    And the gateway should support upstream service discovery
    And the service should include connection pooling
    And the load balancing should be configurable
    And the service should handle upstream failures gracefully

  Scenario: Gateway with caching
    Given I want performant gRPC gateways
    When I generate a grpc-gateway with caching enabled
    Then the service should include response caching
    And the cache should work for both gRPC and REST
    And the caching should be configurable per method
    And the cache should support TTL and invalidation
    And the caching should improve response times

  Scenario: Gateway with API versioning
    Given I want evolvable gRPC gateways
    When I generate a grpc-gateway with API versioning
    Then the protobuf definitions should support versioning
    And the REST API should include version prefixes
    And the gateway should route based on API version
    And the service should maintain backward compatibility
    And the documentation should cover version differences

  Scenario: Gateway with error handling
    Given I want robust gRPC gateways
    When I generate a grpc-gateway with error handling
    Then the service should include structured error responses
    And the errors should be consistent across protocols
    And the gateway should map gRPC errors to HTTP status codes
    And the service should include error recovery mechanisms
    And the error handling should be customizable

  Scenario: Gateway with configuration management
    Given I want configurable gRPC gateways
    When I generate a grpc-gateway with advanced configuration
    Then the service should support environment-based configuration
    And the configuration should be validated at startup
    And the service should support configuration hot-reloading
    And the configuration should include secure secret management
    And the configuration should be documented

  Scenario: Gateway with testing infrastructure
    Given I want well-tested gRPC gateways
    When I generate a grpc-gateway with test infrastructure
    Then the project should include unit tests for gRPC services
    And the project should include integration tests for REST endpoints
    And the project should include end-to-end tests
    And the tests should cover both protocols
    And the testing should include mock generation

  Scenario: Gateway with containerization
    Given I want containerized gRPC gateways
    When I generate a grpc-gateway with container support
    Then the project should include optimized Dockerfile
    And the container should support multi-stage builds
    And the service should include container health checks
    And the container should follow security best practices
    And the deployment should support Kubernetes

  Scenario: Gateway with documentation
    Given I want well-documented gRPC gateways
    When I generate a grpc-gateway with documentation
    Then the project should include comprehensive README
    And the project should include API documentation
    And the protobuf definitions should be well-documented
    And the service should include usage examples
    And the documentation should be up-to-date

  Scenario: Gateway with performance optimization
    Given I want high-performance gRPC gateways
    When I generate a grpc-gateway with performance optimizations
    Then the service should include connection pooling
    And the gateway should include response compression
    And the service should support HTTP/2
    And the service should include request buffering
    And the performance should be measurable with benchmarks