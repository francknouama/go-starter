Feature: Echo Web API Blueprint Generation
  As a Go developer
  I want to generate Echo-based Web API applications
  So that I can build high-performance RESTful services with Echo's rich feature set

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate Echo web API application
    Given I want to create an Echo-based web API application
    When I run the command "go-starter new my-echo-api --type=web-api-echo --framework=echo --database-driver=postgres --no-git"
    Then the generation should succeed
    And the project should contain Echo-specific components
    And the generated code should compile successfully
    And the project should use Echo framework patterns

  Scenario: Echo instance and middleware integration
    Given I have generated an Echo web API application
    When I examine the server configuration
    Then the server should use echo.New()
    And the middleware should be Echo-compatible
    And the routing should use Echo's handler patterns
    And the handlers should accept echo.Context

  Scenario: Echo routing and groups
    Given I have generated an Echo web API application
    When I examine the route definitions
    Then the routes should use Echo's e.GET/POST/PUT/DELETE patterns
    And the routes should support route groups
    And the routes should include parameter binding
    And the route middleware should be configurable per route

  Scenario: Echo context and data binding
    Given I want flexible request handling
    When I generate an Echo web API application
    Then the handlers should use echo.Context
    And request data should bind automatically to structs
    And path parameters should be extractable via context
    And query parameters should be easily accessible

  Scenario: Echo middleware stack
    Given I want comprehensive middleware with Echo
    When I generate an Echo web API application
    Then the middleware should include Echo's built-in middleware
    And custom middleware should follow Echo patterns
    And middleware should have access to echo.Context
    And middleware chaining should be properly configured

  Scenario: Echo request validation
    Given I want validated inputs
    When I generate an Echo web API application
    Then request validation should be integrated
    And validation errors should return proper responses
    And custom validators should be supported
    And validation middleware should be configurable

  Scenario: Echo error handling
    Given I want robust error handling
    When I generate an Echo web API application
    Then error handling should use Echo's error handling
    And custom error handlers should be implemented
    And HTTP error responses should be properly formatted
    And error logging should include request context

  Scenario: Echo authentication and authorization
    Given I want to secure my Echo web API
    When I generate an Echo web API with JWT authentication
    Then JWT middleware should integrate with Echo
    And protected routes should use Echo middleware
    And token validation should work with echo.Context
    And unauthorized access should be properly handled

  Scenario: Echo CORS and security
    Given I want cross-origin and security support
    When I generate an Echo web API with security features
    Then CORS should be configured using Echo middleware
    And security headers should be set via Echo middleware
    And rate limiting should be implemented
    And CSRF protection should be available

  Scenario: Echo templating and responses
    Given I want flexible response handling
    When I generate an Echo web API application
    Then JSON responses should be easily generated
    And response headers should be configurable
    And status codes should be properly set
    And content negotiation should be supported

  Scenario: Echo with database integration
    Given I want database connectivity
    When I generate an Echo web API with database support
    Then database connections should be injected into handlers
    And repository patterns should work with Echo context
    And database errors should be properly handled
    And transaction management should be request-scoped

  Scenario: Echo WebSocket support
    Given I want real-time features
    When I generate an Echo web API with WebSocket support
    Then WebSocket endpoints should be properly configured
    And WebSocket handlers should use Echo patterns
    And WebSocket middleware should be supported
    And connection management should be implemented

  Scenario: Echo testing infrastructure
    Given I want to test my Echo web API
    When I generate an Echo web API application
    Then the test suite should include Echo-specific tests
    And HTTP testing should use Echo test patterns
    And context testing should be properly implemented
    And integration tests should start an Echo server

  Scenario: Echo performance features
    Given I want optimal performance
    When I generate an Echo web API application
    Then the server should be configured for performance
    And request processing should be efficient
    And middleware should be optimized
    And response compression should be available

  Scenario: Echo graceful shutdown
    Given I want reliable deployments
    When I generate an Echo web API application
    Then the Echo server should support graceful shutdown
    And active connections should be handled during shutdown
    And the shutdown process should be logged
    And cleanup should be performed properly

  Scenario: Echo production configuration
    Given I want production-ready Echo web API
    When I generate an Echo web API application
    Then timeouts should be properly configured
    And the server should include production middleware
    And logging should be structured and configurable
    And health checks should be implemented

  Scenario: Echo API documentation
    Given I want documented APIs
    When I generate an Echo web API application
    Then OpenAPI documentation should be generated
    And Echo routes should be automatically documented
    And Swagger UI should be accessible
    And API examples should be functional

  Scenario: Echo logging and observability
    Given I want observable Echo web API
    When I generate an Echo web API with structured logging
    Then request logging should capture Echo-specific data
    And middleware should support request tracing
    And metrics collection should be implemented
    And distributed tracing should be configured

  Scenario: Echo custom middleware
    Given I want extensible middleware
    When I generate an Echo web API application
    Then custom middleware should follow Echo patterns
    And middleware should have access to full context
    And middleware should be easily testable
    And middleware should support configuration

  Scenario: Echo static file serving
    Given I want to serve static content
    When I generate an Echo web API with static file support
    Then static file routes should be configured
    And file serving should be optimized
    And directory listing should be controllable
    And cache headers should be properly set