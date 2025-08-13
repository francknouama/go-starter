Feature: Fiber Web API Blueprint Generation
  As a Go developer
  I want to generate Fiber-based Web API applications
  So that I can build extremely fast RESTful services with Express.js-like syntax

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generate Fiber web API application
    Given I want to create a Fiber-based web API application
    When I run the command "go-starter new my-fiber-api --type=web-api-fiber --framework=fiber --database-driver=postgres --no-git"
    Then the generation should succeed
    And the project should contain Fiber-specific components
    And the generated code should compile successfully
    And the project should use Fiber framework patterns

  Scenario: Fiber app and middleware integration
    Given I have generated a Fiber web API application
    When I examine the server configuration
    Then the app should use fiber.New()
    And the middleware should be Fiber-compatible
    And the routing should use Fiber's handler patterns
    And the handlers should accept *fiber.Ctx

  Scenario: Fiber routing and groups
    Given I have generated a Fiber web API application
    When I examine the route definitions
    Then the routes should use Fiber's app.Get/Post/Put/Delete patterns
    And the routes should support route groups
    And the routes should include parameter binding
    And the route prefixes should be configurable

  Scenario: Fiber context and fast HTTP
    Given I want high-performance request handling
    When I generate a Fiber web API application
    Then the handlers should use *fiber.Ctx
    And request parsing should be optimized for speed
    And response generation should be efficient
    And memory allocation should be minimized

  Scenario: Fiber middleware stack
    Given I want comprehensive middleware with Fiber
    When I generate a Fiber web API application
    Then the middleware should include Fiber's built-in middleware
    And custom middleware should follow Fiber patterns
    And middleware should have access to fiber.Ctx
    And middleware should support next() pattern

  Scenario: Fiber request parsing and validation
    Given I want fast request processing
    When I generate a Fiber web API application
    Then JSON parsing should be extremely fast
    And form data should be efficiently processed
    And file uploads should be properly handled
    And request validation should be integrated

  Scenario: Fiber WebSocket integration
    Given I want real-time capabilities
    When I generate a Fiber web API with WebSocket support
    Then WebSocket endpoints should be properly configured
    And WebSocket handlers should use Fiber patterns
    And connection upgrading should be seamless
    And WebSocket middleware should be supported

  Scenario: Fiber compression and caching
    Given I want optimized responses
    When I generate a Fiber web API application
    Then response compression should be available
    And caching middleware should be configured
    And static file serving should be optimized
    And ETag support should be implemented

  Scenario: Fiber rate limiting and security
    Given I want secure and controlled access
    When I generate a Fiber web API with security features
    Then rate limiting should be efficiently implemented
    And CORS should be configured for Fiber
    And security headers should be set via middleware
    And DDoS protection should be available

  Scenario: Fiber authentication and authorization
    Given I want to secure my Fiber web API
    When I generate a Fiber web API with JWT authentication
    Then JWT middleware should integrate with Fiber
    And protected routes should use Fiber middleware
    And token validation should work with fiber.Ctx
    And authentication should be performant

  Scenario: Fiber error handling and recovery
    Given I want robust error handling
    When I generate a Fiber web API application
    Then error handling should use Fiber's error handling
    And panic recovery should be implemented
    And custom error pages should be supported
    And error responses should be properly formatted

  Scenario: Fiber templating engine
    Given I want templating support
    When I generate a Fiber web API with templating
    Then template engines should be configurable
    And template rendering should be fast
    And template caching should be enabled
    And multiple template formats should be supported

  Scenario: Fiber with database integration
    Given I want database connectivity
    When I generate a Fiber web API with database support
    Then database connections should be efficiently managed
    And repository patterns should work with Fiber context
    And database operations should be optimized
    And connection pooling should be properly configured

  Scenario: Fiber testing infrastructure
    Given I want to test my Fiber web API
    When I generate a Fiber web API application
    Then the test suite should include Fiber-specific tests
    And HTTP testing should use Fiber test patterns
    And performance testing should be included
    And integration tests should start a Fiber server

  Scenario: Fiber performance optimization
    Given I want maximum performance
    When I generate a Fiber web API application
    Then the server should be configured for speed
    And memory usage should be optimized
    And request processing should be minimal overhead
    And response times should be extremely fast

  Scenario: Fiber monitoring and metrics
    Given I want performance monitoring
    When I generate a Fiber web API application
    Then performance metrics should be collected
    And request timing should be measured
    And throughput monitoring should be available
    And resource usage should be tracked

  Scenario: Fiber graceful shutdown
    Given I want reliable deployments
    When I generate a Fiber web API application
    Then the Fiber app should support graceful shutdown
    And active connections should be properly closed
    And shutdown timeouts should be configurable
    And cleanup should be performed efficiently

  Scenario: Fiber production configuration
    Given I want production-ready Fiber web API
    When I generate a Fiber web API application
    Then production settings should be optimized
    And prefork mode should be configurable
    And process management should be included
    And monitoring endpoints should be available

  Scenario: Fiber API documentation
    Given I want documented APIs
    When I generate a Fiber web API application
    Then OpenAPI documentation should be generated
    And Fiber routes should be documented
    And Swagger integration should be available
    And API examples should be functional

  Scenario: Fiber logging and observability
    Given I want observable Fiber web API
    When I generate a Fiber web API with structured logging
    Then request logging should capture Fiber-specific data
    And log levels should be configurable
    And distributed tracing should be supported
    And performance logging should be included

  Scenario: Fiber custom middleware
    Given I want extensible middleware
    When I generate a Fiber web API application
    Then custom middleware should follow Fiber patterns
    And middleware should be easily composable
    And middleware should support configuration
    And middleware should be testable

  Scenario: Fiber static file optimization
    Given I want efficient static file serving
    When I generate a Fiber web API with static files
    Then static file serving should be highly optimized
    And file compression should be automatic
    And cache headers should be properly configured
    And directory browsing should be controllable

  Scenario: Fiber proxy and load balancing
    Given I want proxy capabilities
    When I generate a Fiber web API with proxy support
    Then proxy middleware should be available
    And load balancing should be configurable
    And reverse proxy should be optimized
    And health checks should be integrated