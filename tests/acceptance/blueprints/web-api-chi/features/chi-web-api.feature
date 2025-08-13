Feature: Chi Web API Blueprint Generation
  As a Go developer
  I want to generate Chi-based Web API applications
  So that I can build fast, lightweight RESTful services using the Chi router

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate Chi web API application
    Given I want to create a Chi-based web API application
    When I run the command "go-starter new my-chi-api --type=web-api-chi --framework=chi --database-driver=postgres --no-git"
    Then the generation should succeed
    And the project should contain Chi-specific components
    And the generated code should compile successfully
    And the project should use Chi router patterns

  Scenario: Chi router and middleware integration
    Given I have generated a Chi web API application
    When I examine the server configuration
    Then the router should use chi.NewRouter()
    And the middleware should be Chi-compatible
    And the routing should use Chi's sub-router patterns
    And the handlers should accept http.ResponseWriter and *http.Request

  Scenario: Chi-specific routing patterns
    Given I have generated a Chi web API application
    When I examine the route definitions
    Then the routes should use Chi's r.Route() patterns
    And the routes should support nested sub-routers
    And the routes should include proper HTTP method routing
    And the route groups should be logically organized

  Scenario: Chi middleware stack
    Given I want comprehensive middleware with Chi
    When I generate a Chi web API application
    Then the middleware should include Chi's built-in middleware
    And custom middleware should be Chi-compatible
    And the middleware order should be optimized
    And the middleware should support request context

  Scenario: Chi URL parameters and routing
    Given I want flexible URL routing
    When I generate a Chi web API application
    Then the routes should support URL parameters
    And the parameter extraction should use Chi patterns
    And the route matching should be performant
    And wildcard routes should be supported

  Scenario: Chi request context and values
    Given I want to pass data through request pipeline
    When I generate a Chi web API application
    Then the handlers should use request context
    And middleware should inject values into context
    And context values should be type-safe
    And request scoped data should be accessible

  Scenario: Chi authentication middleware
    Given I want to secure my Chi web API
    When I generate a Chi web API with JWT authentication
    Then the auth middleware should integrate with Chi
    And protected routes should use auth middleware
    And the JWT validation should work with Chi context
    And unauthorized requests should be properly handled

  Scenario: Chi error handling
    Given I want robust error handling
    When I generate a Chi web API application
    Then error handling should be Chi-compatible
    And error responses should use proper HTTP status codes
    And panic recovery should be implemented
    And error logging should include request context

  Scenario: Chi CORS configuration
    Given I want cross-origin support
    When I generate a Chi web API with CORS
    Then CORS should be properly configured for Chi
    And preflight requests should be handled
    And CORS headers should be set correctly
    And origin validation should be implemented

  Scenario: Chi performance optimizations
    Given I want optimal performance
    When I generate a Chi web API application
    Then the router should be configured for performance
    And middleware ordering should be optimized
    And request processing should be efficient
    And memory allocations should be minimized

  Scenario: Chi with database integration
    Given I want database connectivity
    When I generate a Chi web API with database support
    Then database connections should be managed properly
    And repository patterns should work with Chi handlers
    And database transactions should be request-scoped
    And connection pooling should be configured

  Scenario: Chi testing infrastructure
    Given I want to test my Chi web API
    When I generate a Chi web API application
    Then the test suite should include Chi-specific tests
    And HTTP testing should use Chi test patterns
    And middleware testing should be isolated
    And integration tests should start a Chi server

  Scenario: Chi graceful shutdown
    Given I want reliable deployments
    When I generate a Chi web API application
    Then the server should support graceful shutdown
    And in-flight requests should complete during shutdown
    And the shutdown should be properly logged
    And the Chi server should close cleanly

  Scenario: Chi production readiness
    Given I want production-ready Chi web API
    When I generate a Chi web API application
    Then the server should include production middleware
    And timeouts should be properly configured
    And the server should be container-ready
    And monitoring endpoints should be available

  Scenario: Chi API documentation
    Given I want documented APIs
    When I generate a Chi web API application
    Then OpenAPI documentation should be generated
    And Chi routes should be documented
    And API examples should work with Chi
    And documentation should be accessible via HTTP

  Scenario: Chi security headers
    Given I want secure Chi web API
    When I generate a Chi web API with security features
    Then security headers should be properly set
    And CSRF protection should be available
    And input validation should be implemented
    And security scanning should pass

  Scenario: Chi logging and observability
    Given I want observable Chi web API
    When I generate a Chi web API with structured logging
    Then request logging should capture Chi-specific data
    And log correlation should work across middleware
    And performance metrics should be collected
    And distributed tracing should be supported