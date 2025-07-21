# Microservice Blueprint ATDD Test Suite

## Overview

This directory contains comprehensive Acceptance Test-Driven Development (ATDD) tests for microservice and serverless blueprints, specifically `microservice-standard` and `lambda-standard`. The test suite uses Behavior-Driven Development (BDD) with Gherkin feature files and testcontainers for realistic testing scenarios.

## Architecture

### Blueprint Types Covered

| Blueprint | Type | Key Features | Test Focus |
|-----------|------|-------------|------------|
| **microservice-standard** | Distributed service | gRPC, REST, GraphQL protocols | Protocol support, containerization, service mesh |
| **lambda-standard** | Serverless function | AWS Lambda integration | Event processing, cold start optimization |

### Test Structure

```
microservice/
├── features/
│   └── microservice.feature     # Gherkin scenarios (30+ test cases)
├── microservice_steps_test.go   # BDD step definitions with testcontainers
├── microservice_acceptance_test.go # High-level acceptance tests
└── README.md                    # This documentation
```

## Feature Coverage

### 🚀 Core Service Generation
- **Microservice Generation**: gRPC-based services with health checks
- **Lambda Generation**: AWS Lambda functions with SDK integration
- **Code Compilation**: All generated projects compile successfully
- **Container Readiness**: Docker support with multi-stage builds

### 🔗 Communication Protocols
- **gRPC**: Protocol Buffers, streaming, client/server generation
- **REST**: HTTP APIs with OpenAPI documentation
- **GraphQL**: Schema-first development with resolvers

### 🗄️ Database Integration
- **PostgreSQL**: GORM and SQLx ORM support
- **MySQL**: Connection pooling and migrations
- **MongoDB**: Document storage with MongoDB driver
- **Redis**: Caching and session storage

### 🔍 Service Discovery & Mesh
- **Service Registry**: Consul integration for service discovery
- **Health Checks**: Liveness and readiness probes
- **Load Balancing**: Client-side load balancing
- **Service Mesh**: Istio configuration and mTLS support

### 🛡️ Security & Authentication
- **JWT Authentication**: Token-based authentication
- **OAuth2**: Third-party authentication providers
- **mTLS**: Mutual TLS for service communication
- **API Key**: Simple API key authentication

### 📊 Observability
- **OpenTelemetry**: Distributed tracing integration
- **Prometheus**: Metrics collection and monitoring
- **Structured Logging**: Contextual logging with correlation IDs
- **Health Monitoring**: Application and infrastructure health

### ☁️ Cloud & Deployment
- **Kubernetes**: Production-ready manifests
- **Docker**: Multi-stage builds with security best practices
- **AWS Lambda**: Multiple trigger types (API Gateway, S3, DynamoDB)
- **Multi-region**: Geographic distribution strategies

### 🔄 Advanced Patterns
- **Circuit Breaker**: Resilience patterns for fault tolerance
- **Event Sourcing**: CQRS and event-driven architecture
- **Message Queues**: Asynchronous communication
- **Streaming**: Real-time data processing

## Test Scenarios

### Basic Generation Tests
```gherkin
Scenario: Generate microservice-standard application
  Given I want to create a gRPC-based microservice
  When I run the command "go-starter new my-service --type=microservice-standard --framework=grpc --database-driver=postgres --no-git"
  Then the generation should succeed
  And the project should contain all essential microservice components
  And the generated code should compile successfully
```

### Protocol Support Tests
```gherkin
Scenario: Microservice with different communication protocols
  Given I want to create microservices with various protocols
  When I generate a microservice with protocol "grpc"
  Then the project should support the "grpc" communication pattern
  And the service should include appropriate client libraries
  And the protocol-specific middleware should be configured
```

### Advanced Integration Tests
```gherkin
Scenario: Service with database integration
  Given I want to create services with database persistence
  When I generate a service with database "postgres" and ORM "gorm"
  Then the project should include database configuration
  And the service should support database migrations
  And the database connection should be testable with containers
```

## Running Tests

### Prerequisites
```bash
# Install required dependencies
go install github.com/cucumber/godog/cmd/godog@latest

# Ensure Docker is running for container tests
docker --version

# Ensure go-starter CLI is available
go-starter --help
```

### Test Execution

#### Run All Microservice Tests
```bash
# Full test suite with BDD scenarios
go test -v ./tests/acceptance/blueprints/microservice/

# Run with testcontainers (requires Docker)
go test -v ./tests/acceptance/blueprints/microservice/ -tags integration
```

#### Run Specific Test Categories
```bash
# Basic generation tests only
go test -v -run TestMicroserviceStandardGeneration ./tests/acceptance/blueprints/microservice/

# Protocol support tests
go test -v -run TestProtocolSupport ./tests/acceptance/blueprints/microservice/

# Database integration tests
go test -v -run TestDatabaseIntegration ./tests/acceptance/blueprints/microservice/

# Service discovery tests
go test -v -run TestServiceDiscovery ./tests/acceptance/blueprints/microservice/
```

#### Run BDD Scenarios with Godog
```bash
# Run specific feature scenarios
godog run features/microservice.feature

# Run with specific tags
godog run features/microservice.feature -t @microservice

# Generate test reports
godog run features/microservice.feature -f junit:reports/microservice-test-results.xml
```

#### Performance Testing
```bash
# Benchmark microservice generation
go test -bench=BenchmarkServiceGeneration ./tests/acceptance/blueprints/microservice/

# Benchmark compilation speed
go test -bench=BenchmarkServiceCompilation ./tests/acceptance/blueprints/microservice/
```

### Test Configuration

#### Environment Variables
```bash
# Test timeouts and retries
export TEST_TIMEOUT=30s
export TEST_RETRY_ATTEMPTS=3

# Database testing
export TEST_DATABASE_URL="postgres://test:test@localhost:5432/testdb"

# Container testing
export TESTCONTAINER_RYUK_DISABLED=true  # For CI environments

# Service endpoints
export SERVICE_BASE_URL="http://localhost:8080"
```

#### CI/CD Integration
```bash
# Short test mode (no containers)
go test -short ./tests/acceptance/blueprints/microservice/

# Parallel execution
go test -parallel 4 ./tests/acceptance/blueprints/microservice/

# Coverage reporting
go test -coverprofile=coverage.out ./tests/acceptance/blueprints/microservice/
```

## Test Implementation Details

### Testcontainer Integration

The test suite uses testcontainers for realistic testing:

```go
// Database containers for integration testing
func (ctx *ServiceTestContext) startDatabaseContainer(dbType string) {
    switch dbType {
    case "postgres":
        req = testcontainers.ContainerRequest{
            Image: "postgres:15-alpine",
            ExposedPorts: []string{"5432/tcp"},
            Env: map[string]string{
                "POSTGRES_DB": "testdb",
                "POSTGRES_USER": "testuser",
                "POSTGRES_PASSWORD": "testpass",
            },
            WaitingFor: wait.ForListeningPort("5432/tcp"),
        }
    }
}
```

### BDD Step Definitions

Key step definition patterns:

```go
// Generation steps
sc.When(`^I run the command "([^"]*)"$`, ctx.iRunTheCommand)
sc.Then(`^the generation should succeed$`, ctx.theGenerationShouldSucceed)

// Validation steps
sc.Then(`^the project should contain all essential microservice components$`, 
    ctx.theProjectShouldContainAllEssentialMicroserviceComponents)

// Integration steps
sc.Then(`^the database connection should be testable with containers$`, 
    ctx.theDatabaseConnectionShouldBeTestableWithContainers)
```

### File Validation

Comprehensive file structure validation:

```go
func (ctx *ServiceTestContext) checkRequiredFiles(files []string) error {
    var missingFiles []string
    for _, file := range files {
        filePath := filepath.Join(ctx.projectPath, file)
        if _, err := os.Stat(filePath); os.IsNotExist(err) {
            missingFiles = append(missingFiles, file)
        }
    }
    if len(missingFiles) > 0 {
        return fmt.Errorf("missing required files: %v", missingFiles)
    }
    return nil
}
```

## Test Categories

### 1. Generation Tests
- **Focus**: Basic project generation and file creation
- **Coverage**: All service blueprint types
- **Validation**: File existence, go.mod structure, compilation

### 2. Protocol Tests
- **Focus**: Communication protocol support
- **Coverage**: gRPC, REST, GraphQL implementations
- **Validation**: Client/server code, middleware, dependencies

### 3. Database Tests
- **Focus**: Data persistence and ORM integration
- **Coverage**: Multiple databases and ORMs
- **Validation**: Connection pooling, migrations, repositories

### 4. Container Tests
- **Focus**: Containerization and deployment
- **Coverage**: Docker, Kubernetes, health checks
- **Validation**: Multi-stage builds, security practices

### 5. Integration Tests
- **Focus**: Service-to-service communication
- **Coverage**: Service discovery, load balancing, mesh
- **Validation**: Registry integration, health endpoints

### 6. Security Tests
- **Focus**: Authentication and authorization
- **Coverage**: JWT, OAuth2, mTLS, API keys
- **Validation**: Middleware, token handling, encryption

### 7. Observability Tests
- **Focus**: Monitoring and tracing
- **Coverage**: Metrics, logs, traces, alerts
- **Validation**: OpenTelemetry, Prometheus integration

## Troubleshooting

### Common Issues

#### Generation Failures
```bash
# Check go-starter availability
which go-starter
go-starter --version

# Verify command syntax
go-starter new test --type=microservice-standard --dry-run
```

#### Container Issues
```bash
# Check Docker daemon
docker info

# Clean up containers
docker container prune -f

# Check testcontainer logs
export TESTCONTAINERS_DEBUG=true
```

#### Compilation Errors
```bash
# Clean module cache
go clean -modcache

# Verify Go version
go version  # Should be 1.21+

# Check generated go.mod
cd generated-project && go mod tidy
```

### Debug Mode

Enable detailed logging:

```bash
# Enable debug output
export DEBUG=true

# Verbose test output
go test -v -args -test.v

# Container debugging
export TESTCONTAINERS_DEBUG=true
```

## Performance Benchmarks

### Expected Performance Metrics

| Test Category | Expected Duration | Resource Usage |
|---------------|------------------|----------------|
| Basic Generation | < 5s | Low CPU/Memory |
| Protocol Tests | < 10s per protocol | Medium CPU |
| Database Tests | < 30s per database | High Memory (containers) |
| Integration Tests | < 60s | High CPU/Memory |

### Optimization Tips

1. **Parallel Execution**: Use `-parallel` flag for independent tests
2. **Container Reuse**: Reuse containers across related test scenarios
3. **Short Mode**: Skip container tests in development with `-short`
4. **Selective Testing**: Run specific test categories during development

## Contributing

### Adding New Test Scenarios

1. **Add Gherkin Scenario**: Update `features/microservice.feature`
2. **Implement Step Definition**: Add to `microservice_steps_test.go`
3. **Create Integration Test**: Add to `microservice_acceptance_test.go`
4. **Update Documentation**: Update this README

### Test Development Guidelines

1. **Use Testcontainers**: For realistic database/service testing
2. **BDD Approach**: Write scenarios from user perspective
3. **Comprehensive Validation**: Check file structure, compilation, runtime
4. **Performance Aware**: Keep test execution time reasonable
5. **CI/CD Ready**: Ensure tests work in containerized environments

### Example New Scenario

```gherkin
Scenario: Service with custom middleware
  Given I want to create a service with custom middleware
  When I generate a service with middleware "rate-limiting,cors,auth"
  Then the service should include rate limiting middleware
  And the service should include CORS middleware  
  And the service should include auth middleware
  And the middleware should be properly configured
```

## Related Documentation

- [Monolith Blueprint ATDD](../monolith/README.md)
- [Web API Blueprint ATDD](../web-api/README.md)
- [CLI Blueprint ATDD](../cli/README.md)
- [Blueprint Development Guide](../../../../docs/blueprint-development.md)
- [Testing Strategy](../../../../docs/testing-strategy.md)

## Support

For issues with microservice blueprint testing:

1. Check the [troubleshooting section](#troubleshooting)
2. Review the [CI/CD logs](../../../../.github/workflows/)
3. Open an issue with test output and environment details
4. Consult the main project documentation

---

**Note**: This test suite provides comprehensive coverage for distributed systems and serverless patterns. The combination of BDD scenarios and testcontainer integration ensures reliable validation of complex service architectures while maintaining fast feedback cycles for development teams.