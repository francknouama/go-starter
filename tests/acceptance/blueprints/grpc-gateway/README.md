# gRPC Gateway Blueprint ATDD Tests

This directory contains comprehensive Acceptance Test-Driven Development (ATDD) tests for the gRPC Gateway blueprint, validating that generated projects meet user expectations and business requirements.

## Overview

The gRPC Gateway blueprint generates projects that expose gRPC services through REST APIs, providing dual protocol support for maximum client compatibility. These tests ensure generated projects compile, function correctly, and follow best practices.

## Test Structure

### Gherkin Feature Files (`features/`)
- **grpc-gateway.feature**: 22 comprehensive scenarios covering all gRPC Gateway capabilities
  - Basic gRPC Gateway generation with REST endpoints
  - Service type variations (user, order, notification, auth services)
  - Database integration with multiple ORMs (GORM, SQLx, MongoDB)
  - Authentication (JWT, OAuth2, mTLS, API key)
  - Observability (tracing, metrics, structured logging)
  - Rate limiting across protocols
  - Input validation and error handling
  - Multiple service support and routing
  - Streaming capabilities (server, client, bidirectional)
  - Custom HTTP mapping with protobuf annotations
  - Middleware chains and interceptors
  - Health checks and load balancing
  - Caching and performance optimization
  - API versioning and backward compatibility
  - Containerization and Kubernetes support
  - Comprehensive documentation generation

### Test Implementation Files
- **grpc_gateway_steps_test.go**: BDD step definitions implementing all Gherkin scenarios
  - Complete step implementation for all 22 scenarios
  - Testcontainers integration for realistic database testing
  - gRPC and REST protocol validation
  - Performance and security testing
  - Cleanup and resource management

- **grpc_gateway_acceptance_test.go**: High-level acceptance tests
  - End-to-end project generation validation
  - Feature variation testing
  - Integration scenario testing
  - Compilation and build verification

## Test Categories

### 1. Basic Generation Tests
Validates fundamental gRPC Gateway project structure:
- Go module initialization
- Protobuf definitions and code generation
- gRPC server implementation
- REST gateway configuration
- Docker and Kubernetes manifests
- Documentation and examples

### 2. Feature Integration Tests
Tests specific feature combinations:
- **Database Integration**: PostgreSQL, MySQL, SQLite with GORM/SQLx
- **Authentication**: JWT tokens, OAuth2, mTLS certificates, API keys
- **Observability**: OpenTelemetry tracing, Prometheus metrics, structured logging
- **Rate Limiting**: Per-method limits across gRPC and REST
- **Validation**: Protobuf validation rules and error formatting
- **Streaming**: Server/client/bidirectional streaming support
- **Load Balancing**: Client-side balancing and service discovery
- **Caching**: Response caching with TTL and invalidation

### 3. Advanced Pattern Tests
Validates complex architectural patterns:
- **Multi-Service Architecture**: Multiple protobuf services in single gateway
- **Custom HTTP Mapping**: Protobuf annotations for REST API design
- **API Versioning**: Backward compatibility and version routing
- **Middleware Chains**: gRPC interceptors and HTTP middleware
- **Error Handling**: Consistent error responses across protocols
- **Configuration Management**: Environment-based config and secrets

### 4. Quality Assurance Tests
Ensures production readiness:
- **Compilation**: All generated code compiles successfully
- **Testing Infrastructure**: Unit, integration, and e2e test scaffolding
- **Containerization**: Optimized Dockerfiles and security practices
- **Documentation**: Comprehensive README, API docs, and examples
- **Performance**: Benchmarks and optimization features

## Running the Tests

### Prerequisites
- Go 1.22+ installed
- Docker available (for testcontainers)
- `go-starter` CLI tool built and in PATH
- Sufficient disk space for container images

### Individual Test Execution
```bash
# Run all gRPC Gateway ATDD tests
go test -v ./tests/acceptance/blueprints/grpc-gateway/... -timeout 30m

# Run specific feature scenarios
go test -v ./tests/acceptance/blueprints/grpc-gateway/ -run "TestFeatures" -timeout 20m

# Run acceptance tests only
go test -v ./tests/acceptance/blueprints/grpc-gateway/ -run "TestGRPCGatewayBlueprintGeneration" -timeout 15m

# Run integration scenarios (requires more time and resources)
go test -v ./tests/acceptance/blueprints/grpc-gateway/ -run "TestGRPCGatewayIntegration" -timeout 25m
```

### CI/CD Integration
Tests are automatically executed in GitHub Actions workflows:
- **ci.yml**: Includes dedicated `acceptance-tests` job (25-minute timeout)
- **tdd-enforcement.yml**: Validates test coverage and quality
- Both workflows support parallel execution across Go versions

## Test Scenarios Detailed

### Core Scenarios

#### Basic gRPC Gateway (Scenario 1)
- **Given**: Developer wants gRPC service with REST gateway
- **When**: `go-starter new my-grpc-gateway --type=grpc-gateway`
- **Then**: Complete dual-protocol service is generated
- **Validates**: Project structure, compilation, protocol support, protobuf definitions

#### Service Types (Scenarios 2-3)
Tests service-specific patterns:
- **User Service**: User management, authentication integration
- **Order Service**: Transaction processing, database operations  
- **Notification Service**: Event handling, streaming capabilities
- **Auth Service**: Security patterns, token management

#### Database Integration (Scenarios 4-6)
- **PostgreSQL + GORM**: Full ORM with migrations and transactions
- **PostgreSQL + SQLx**: Raw SQL with type safety
- **MySQL + GORM**: Cross-database compatibility
- **MongoDB**: NoSQL document operations

### Advanced Scenarios

#### Authentication (Scenarios 7-10)
- **JWT Authentication**: Token validation across protocols
- **OAuth2**: Third-party authentication integration
- **mTLS**: Certificate-based mutual authentication
- **API Key**: Simple key-based authentication

#### Observability (Scenarios 11-12)
- **Distributed Tracing**: OpenTelemetry integration
- **Metrics Collection**: Prometheus metrics for both protocols
- **Structured Logging**: JSON logging with correlation IDs
- **Protocol Distinction**: Separate metrics for gRPC vs REST

#### Advanced Features (Scenarios 13-22)
- **Rate Limiting**: Configurable limits per method/protocol
- **Input Validation**: Protobuf validation rules
- **Multi-Service**: Multiple services in single gateway
- **Streaming**: All streaming types with error handling
- **Custom Mapping**: Protobuf HTTP annotations
- **Middleware Chains**: Extensible interceptor/middleware system
- **Health Checks**: gRPC health service + HTTP endpoints
- **Load Balancing**: Client-side balancing and service discovery
- **Caching**: Response caching with protocol awareness
- **API Versioning**: Backward-compatible version management
- **Error Handling**: Consistent error mapping
- **Configuration**: Environment-based config with validation
- **Testing**: Comprehensive test infrastructure
- **Containerization**: Production-ready Docker/K8s setup
- **Documentation**: Auto-generated API docs and examples
- **Performance**: HTTP/2, compression, benchmarking

## Expected Project Structure

Generated gRPC Gateway projects follow this structure:
```
my-grpc-gateway/
├── go.mod                          # Go module with gRPC dependencies
├── main.go                         # Application entry point
├── cmd/
│   ├── server/main.go             # Server command
│   └── client/main.go             # Client examples
├── proto/
│   ├── service.proto              # Service definitions
│   ├── gateway.yaml               # HTTP mapping rules
│   └── generated/                 # Generated Go code
├── internal/
│   ├── server/grpc.go            # gRPC server implementation
│   ├── gateway/rest.go           # REST gateway implementation
│   ├── service/                  # Business logic
│   ├── middleware/               # gRPC interceptors
│   ├── auth/                     # Authentication
│   ├── metrics/                  # Observability
│   └── config/                   # Configuration
├── pkg/
│   ├── api/v1/                   # Generated API
│   └── client/                   # Client libraries
├── config/
│   ├── config.yaml               # Application configuration
│   └── database.go               # Database setup
├── migrations/                    # Database migrations
├── tests/
│   ├── unit/                     # Unit tests
│   ├── integration/              # Integration tests
│   └── e2e/                      # End-to-end tests
├── docs/
│   ├── api/                      # Generated API documentation
│   └── examples/                 # Usage examples
├── scripts/
│   ├── proto-gen.sh              # Protobuf generation
│   └── build.sh                  # Build automation
├── k8s/                          # Kubernetes manifests
├── docker-compose.yml            # Development setup
├── Dockerfile                    # Production container
├── Makefile                      # Build automation
└── README.md                     # Project documentation
```

## Test Data and Scenarios

### Feature Flags Tested
Tests validate all gRPC Gateway blueprint flags:
- `--database-driver`: postgres, mysql, mongodb, sqlite
- `--database-orm`: gorm, sqlx, mongo
- `--auth-type`: jwt, oauth2, api-key, mtls  
- `--observability`: tracing, metrics, logging
- `--rate-limiting`: per-method configuration
- `--validation`: protobuf validation rules
- `--multi-service`: multiple service definitions
- `--streaming`: streaming capabilities
- `--custom-mapping`: HTTP annotations
- `--middleware`: interceptor chains
- `--health-checks`: health endpoints
- `--load-balancing`: service discovery
- `--caching`: response caching
- `--api-versioning`: version management
- `--containerization`: Docker/K8s support

### Performance Benchmarks
Tests include performance validation:
- **Startup Time**: < 2 seconds for basic service
- **Memory Usage**: < 50MB baseline memory
- **Request Latency**: < 10ms for simple operations
- **Throughput**: > 1000 requests/second basic load
- **Container Size**: < 20MB optimized Docker image

## Debugging and Troubleshooting

### Common Test Failures

1. **Container Startup Issues**:
   ```bash
   # Check Docker daemon
   docker version
   
   # Clean up containers
   docker system prune -f
   ```

2. **Compilation Failures**:
   ```bash
   # Verify protobuf tools
   protoc --version
   
   # Check Go module
   go mod verify
   go mod tidy
   ```

3. **Network Issues**:
   ```bash
   # Check port availability
   netstat -tulpn | grep :8080
   
   # Use different ports
   export GRPC_PORT=9090
   export HTTP_PORT=8080
   ```

4. **Resource Constraints**:
   ```bash
   # Check disk space
   df -h
   
   # Check memory
   free -h
   
   # Reduce parallel tests
   go test -p 1 -v ./...
   ```

### Test Environment Variables
- `GRPC_GATEWAY_TEST_TIMEOUT`: Test timeout (default: 30m)
- `GRPC_GATEWAY_SKIP_DOCKER`: Skip Docker-dependent tests
- `GRPC_GATEWAY_VERBOSE`: Enable verbose output
- `GRPC_GATEWAY_PARALLEL`: Number of parallel tests

### Logging and Debug Output
Tests produce detailed logs for debugging:
- **Command Output**: All CLI command outputs captured
- **Generated Files**: Lists of all generated files
- **Compilation Logs**: Build errors and warnings
- **Container Logs**: Docker container startup logs
- **Network Traces**: gRPC/HTTP request/response logs

## Contributing

When adding new gRPC Gateway features:

1. **Add Gherkin Scenarios**: Update `features/grpc-gateway.feature`
2. **Implement Steps**: Add step definitions in `grpc_gateway_steps_test.go`
3. **Add Acceptance Tests**: Update `grpc_gateway_acceptance_test.go`
4. **Validate CI**: Ensure tests pass in GitHub Actions
5. **Update Documentation**: Update this README and examples

### Scenario Writing Guidelines
- Use Given-When-Then format consistently
- Focus on user value and business requirements
- Include both positive and negative scenarios
- Validate end-to-end functionality
- Test error conditions and edge cases

### Step Implementation Guidelines
- Use testcontainers for realistic integration testing
- Implement proper cleanup and resource management
- Provide detailed error messages for failures
- Support both quick tests and thorough validation
- Follow existing patterns for consistency

## References

- [gRPC-Gateway Documentation](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
- [Testcontainers Go](https://testcontainers.com/modules/golang/)
- [Godog BDD Framework](https://github.com/cucumber/godog)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Prometheus Metrics](https://prometheus.io/docs/guides/go-application/)

## Test Metrics

Current test coverage for gRPC Gateway blueprint:
- **Scenarios**: 22 comprehensive scenarios
- **Step Definitions**: 100+ step implementations
- **Feature Flags**: 15+ blueprint options tested
- **Integration Points**: Database, Auth, Observability, Containerization
- **Platforms**: Linux, macOS, Windows (CI/CD)
- **Go Versions**: 1.22, 1.23, 1.24

The gRPC Gateway ATDD suite ensures generated projects are production-ready, well-documented, and follow industry best practices for dual-protocol service architecture.