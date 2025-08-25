Feature: gRPC Gateway Production Readiness (Post-Audit)
  As a developer deploying gRPC Gateway services to production
  I want to ensure all audit findings are addressed
  So that my service meets 2024-2025 microservices standards

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory
    And I have the project audit requirements

  @critical @P0
  Scenario: Health Check Protocol Implementation
    Given I want production-ready gRPC Gateway services
    When I generate a grpc-gateway project with health checks enabled
    Then the service should implement gRPC health checking protocol
    And the health service should be registered with the gRPC server
    And the health endpoints should support both gRPC and REST protocols
    And the health checks should verify database connectivity
    And the health checks should support Kubernetes liveness probes
    And the health checks should support Kubernetes readiness probes
    
    Examples:
      | health_check_type  | grpc_endpoint                    | rest_endpoint     |
      | liveness           | grpc.health.v1.Health/Check      | /health/live     |
      | readiness          | grpc.health.v1.Health/Check      | /health/ready    |
      | database           | grpc.health.v1.Health/Check      | /health/db       |

  @critical @P0
  Scenario: Distributed Tracing Integration
    Given I want observable gRPC Gateway services
    When I generate a grpc-gateway with OpenTelemetry tracing
    Then the service should include gRPC tracing interceptors
    And the gateway should include HTTP tracing middleware
    And the service should generate trace spans for all requests
    And the tracing should include request correlation IDs
    And the trace context should propagate between gRPC and REST layers
    And the tracing should integrate with Jaeger or Zipkin
    
    Examples:
      | tracer_provider | span_attributes                    | propagation_format |
      | jaeger          | grpc.method,grpc.service,user.id   | b3                |
      | zipkin          | http.method,http.status_code       | jaeger            |
      | otel            | service.name,service.version       | tracecontext      |

  @critical @P0
  Scenario: Metrics Collection Implementation
    Given I want metrics-enabled gRPC Gateway services
    When I generate a grpc-gateway with Prometheus metrics
    Then the service should include gRPC metrics interceptors
    And the gateway should include HTTP metrics middleware
    And the metrics should distinguish between gRPC and REST requests
    And the service should export standard gRPC metrics
    And the service should export custom business metrics
    And the metrics endpoint should be secure but accessible
    
    Examples:
      | metric_type        | metric_name                    | labels                      |
      | grpc_requests      | grpc_server_requests_total     | method,status               |
      | http_requests      | http_requests_total            | method,status,endpoint      |
      | request_duration   | request_duration_seconds       | protocol,method             |
      | active_connections | active_connections             | protocol                    |

  @high @P1
  Scenario: Enhanced Error Mapping
    Given I want robust error handling in gRPC Gateway
    When I generate a grpc-gateway with enhanced error mapping
    Then the service should map domain errors to appropriate gRPC codes
    And the gateway should translate gRPC errors to HTTP status codes
    And the error responses should be consistent across protocols
    And the errors should include proper error context
    And the service should support structured error responses
    
    Examples:
      | domain_error      | grpc_code          | http_status | error_context          |
      | UserNotFound      | NOT_FOUND          | 404        | user_id,resource_type  |
      | ValidationError   | INVALID_ARGUMENT   | 400        | field_name,rule        |
      | DuplicateUser     | ALREADY_EXISTS     | 409        | conflicting_field      |
      | DatabaseError     | INTERNAL           | 500        | operation,table        |

  @high @P1
  Scenario: Advanced Authentication & Authorization
    Given I want secure gRPC Gateway services
    When I generate a grpc-gateway with RBAC authentication
    Then the service should include gRPC authentication interceptors
    And the gateway should include HTTP authentication middleware
    And the service should support role-based access control
    And the authentication should validate JWT tokens properly
    And the service should support granular authorization rules
    And the authentication should work consistently across protocols
    
    Examples:
      | auth_method | token_validation     | authorization_scope    | protocol_compatibility |
      | jwt         | HS256,RS256,ES256   | read:users,write:users | grpc+rest             |
      | oauth2      | bearer_token        | admin,user,guest       | grpc+rest             |
      | api_key     | header_validation   | service:internal       | grpc+rest             |

  @high @P1
  Scenario: Rate Limiting Implementation
    Given I want rate-limited gRPC Gateway services
    When I generate a grpc-gateway with intelligent rate limiting
    Then the service should include gRPC rate limiting interceptors
    And the gateway should include HTTP rate limiting middleware
    And the rate limits should be configurable per method
    And the rate limits should support different algorithms
    And the service should return appropriate rate limit headers
    And the rate limiting should be consistent across protocols
    
    Examples:
      | algorithm      | scope         | limit_per_minute | protocol | error_response    |
      | token_bucket   | per_user      | 100             | grpc     | RESOURCE_EXHAUSTED |
      | sliding_window | per_ip        | 1000            | http     | 429               |
      | fixed_window   | per_method    | 50              | both     | appropriate       |

  @medium @P2
  Scenario: Service Mesh Integration
    Given I want service-mesh-ready gRPC Gateway services
    When I generate a grpc-gateway with service mesh features
    Then the service should support client-side load balancing
    And the service should integrate with service discovery
    And the service should implement circuit breaker patterns
    And the service should support connection pooling
    And the service should handle upstream failures gracefully
    And the configuration should be externalized for mesh management
    
    Examples:
      | service_mesh | load_balancing | discovery_method | circuit_breaker | connection_pool |
      | istio        | round_robin    | k8s_dns         | enabled         | 10_connections  |
      | consul       | least_conn     | consul_catalog  | enabled         | 20_connections  |
      | linkerd      | random         | linkerd_proxy   | disabled        | 5_connections   |

  @medium @P2
  Scenario: Input Validation Enhancement
    Given I want validated gRPC Gateway services
    When I generate a grpc-gateway with comprehensive validation
    Then the protobuf definitions should include validation rules
    And the gRPC service should validate input messages
    And the REST gateway should validate HTTP requests
    And the validation errors should be properly formatted
    And the validation should be consistent across protocols
    And the service should support custom validation rules
    
    Examples:
      | field_type | validation_rule              | grpc_validation        | http_validation     |
      | email      | format=email                | protoc-gen-validate    | middleware          |
      | uuid       | format=uuid                 | custom_validator       | path_param          |
      | range      | min=1,max=100              | field_constraint       | query_param         |
      | required   | required=true              | field_presence         | request_body        |

  @medium @P2
  Scenario: Configuration Management Enhancement
    Given I want configurable gRPC Gateway services
    When I generate a grpc-gateway with advanced configuration
    Then the service should support environment-based configuration
    And the configuration should be validated at startup
    And the service should support configuration hot-reloading
    And the configuration should include secure secret management
    And the configuration should be well-documented
    And the configuration should support different deployment environments
    
    Examples:
      | config_source | validation_level | hot_reload | secret_management | environment |
      | env_vars      | strict          | enabled    | k8s_secrets      | production  |
      | config_files  | lenient         | disabled   | hashicorp_vault  | staging     |
      | consul_kv     | strict          | enabled    | aws_ssm          | development |

  @low @P3
  Scenario: Performance Optimization Features
    Given I want high-performance gRPC Gateway services
    When I generate a grpc-gateway with performance optimizations
    Then the service should include connection pooling
    And the gateway should support response compression
    And the service should optimize for HTTP/2 usage
    And the service should include request buffering
    And the performance should be measurable with benchmarks
    And the service should support streaming optimization
    
    Examples:
      | optimization_type | feature              | expected_improvement | measurement_method  |
      | connection_pool   | reuse_connections   | 30% latency         | benchmark_test     |
      | compression       | gzip_responses      | 60% bandwidth       | response_size      |
      | http2             | multiplexing        | 40% throughput      | concurrent_requests |
      | buffering         | batch_processing    | 25% cpu_usage       | profiling          |

  @integration @P1
  Scenario: End-to-End Protocol Compatibility
    Given I have a complete gRPC Gateway service
    When I test both gRPC and REST protocols simultaneously
    Then both protocols should provide identical functionality
    And the data consistency should be maintained across protocols
    And the authentication should work seamlessly for both
    And the error handling should be consistent
    And the performance characteristics should be acceptable
    And the observability should cover both protocols equally
    
    Examples:
      | test_scenario     | grpc_endpoint                    | rest_endpoint        | expected_behavior |
      | create_user       | UserService.CreateUser          | POST /api/v1/users   | identical_result  |
      | get_user          | UserService.GetUser             | GET /api/v1/users/id | identical_result  |
      | auth_failure      | any_authenticated_endpoint      | any_rest_endpoint    | same_error_code   |
      | validation_error  | invalid_request_data            | invalid_json         | consistent_error  |

  @deployment @P1
  Scenario: Container and Kubernetes Readiness
    Given I want to deploy gRPC Gateway to Kubernetes
    When I generate a grpc-gateway with container support
    Then the project should include optimized Dockerfile
    And the container should support multi-stage builds
    And the service should include proper health check endpoints
    And the container should follow security best practices
    And the Kubernetes manifests should be production-ready
    And the deployment should support horizontal scaling
    
    Examples:
      | deployment_target | container_features           | k8s_resources              | security_features    |
      | production        | multi_stage,health_checks   | deployment,service,ingress | non_root,read_only   |
      | staging           | single_stage,basic_health   | deployment,service         | basic_security       |
      | development       | debug_enabled,hot_reload    | deployment                 | minimal_security     |

  @security @P0
  Scenario: Security Hardening Compliance
    Given I need production-grade security
    When I generate a grpc-gateway with security hardening
    Then the service should implement defense-in-depth
    And all communications should be encrypted with TLS 1.3
    And the service should validate all input rigorously
    And the authentication should follow OWASP guidelines
    And the service should implement proper CORS policies
    And the security configuration should be audit-ready
    
    Examples:
      | security_layer | implementation               | compliance_standard | audit_requirement  |
      | transport      | tls_1.3,certificate_pinning | NIST                | certificate_mgmt   |
      | application    | input_validation,sanitization | OWASP               | vulnerability_scan |
      | authentication | jwt,oauth2,rate_limiting     | OAuth2.1            | token_audit        |
      | authorization  | rbac,resource_permissions    | custom              | access_logs        |