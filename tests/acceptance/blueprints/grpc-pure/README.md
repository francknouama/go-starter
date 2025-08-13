# gRPC-Pure Blueprint ATDD Tests

This directory contains comprehensive Acceptance Test-Driven Development (ATDD) tests for the gRPC-Pure blueprint. These tests ensure that the generated gRPC services meet all requirements and function correctly across different configurations.

## Test Structure

### Core Test Files

- **`grpc_pure_atdd_test.go`** - Main ATDD test suite with basic generation and compilation testing
- **`grpc_pure_steps_test.go`** - BDD-style step definitions implementing Given-When-Then scenarios
- **`integration_test.go`** - End-to-end integration tests for complex scenarios
- **`features/grpc-pure.feature`** - Gherkin feature file defining test scenarios in natural language

### Test Categories

#### 1. Basic Generation Tests (`grpc_pure_atdd_test.go`)
- Basic gRPC Pure service generation
- Authentication variants (JWT, mTLS, OAuth2)
- Observability features (tracing, metrics)
- Database integration (PostgreSQL, MySQL, SQLite with GORM/SQLx)
- Service discovery (Consul, etcd, Kubernetes)
- Different logger implementations (slog, zap, logrus, zerolog)

#### 2. BDD Step Definitions (`grpc_pure_steps_test.go`)
Implements step definitions for the feature file scenarios:
- **Given** steps: Set up test preconditions
- **When** steps: Execute actions (generation commands, examinations)
- **Then** steps: Validate outcomes and assertions

#### 3. Integration Tests (`integration_test.go`)
- Production-ready configurations with all features
- Minimal setup validation
- High-security configurations with mTLS
- Multi-database setups
- Performance and load testing infrastructure
- Cross-platform compatibility
- Docker and Kubernetes deployment readiness

## Test Coverage

### Blueprint Variations Tested

| Feature Combination | Authentication | Observability | Service Discovery | Database | Logger |
|-------------------|---------------|---------------|------------------|----------|---------|
| **Basic** | None | None | None | None | slog |
| **JWT Auth** | JWT | None | None | None | slog |
| **Full Observability** | None | Tracing + Metrics | None | None | slog |
| **Service Discovery** | None | None | Consul | None | slog |
| **Database** | None | None | None | PostgreSQL + GORM | slog |
| **Alternative Logger** | None | None | None | None | zap |
| **mTLS** | mTLS | None | None | None | slog |
| **Production** | JWT | Tracing + Metrics | Consul | PostgreSQL + GORM | zap |

### Validation Areas

#### Protocol Buffer Validation
- ✅ Syntax correctness using `buf lint`
- ✅ Import validity
- ✅ Service definition annotations
- ✅ Message validation rules
- ✅ Code generation with `buf generate`

#### gRPC Server Validation
- ✅ Server configuration correctness
- ✅ Interceptor chain order and configuration
- ✅ Health check implementation
- ✅ TLS/mTLS configuration
- ✅ Service registration

#### Interceptor Chain Testing
- ✅ Authentication interceptor (JWT token validation)
- ✅ Logging interceptor (request/response logging)
- ✅ Metrics interceptor (Prometheus metrics collection)
- ✅ Recovery interceptor (panic handling)
- ✅ Rate limiting interceptor
- ✅ Tracing interceptor (OpenTelemetry integration)

#### Observability Validation
- ✅ OpenTelemetry tracing configuration
- ✅ Prometheus metrics setup
- ✅ Metrics collection accuracy
- ✅ Trace propagation
- ✅ Performance monitoring

#### Database Integration Testing
- ✅ Connection configuration
- ✅ Migration system
- ✅ Repository pattern implementation
- ✅ ORM integration (GORM/SQLx)
- ✅ Multiple database driver support

#### Service Discovery Testing
- ✅ Consul integration
- ✅ etcd integration
- ✅ Kubernetes service discovery
- ✅ Service registration and deregistration
- ✅ Health check integration

#### Compilation and Build Testing
- ✅ Go module resolution
- ✅ Dependency management
- ✅ Cross-platform compilation
- ✅ Docker containerization
- ✅ Kubernetes deployment readiness

## Test Execution

### Prerequisites

Install required testing tools:

```bash
# Protocol buffer tools
go install github.com/bufbuild/buf/cmd/buf@latest

# gRPC testing tools  
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Load testing for gRPC
go install github.com/bojand/ghz/cmd/ghz@latest
```

### Running Tests

```bash
# Run all gRPC-Pure ATDD tests
go test ./tests/acceptance/blueprints/grpc-pure -v

# Run only basic generation tests
go test ./tests/acceptance/blueprints/grpc-pure -v -run TestGRPCPureBlueprintBasicGeneration

# Run integration tests (longer running)
go test ./tests/acceptance/blueprints/grpc-pure -v -run TestGRPCPureIntegrationScenarios

# Run with short mode (skips integration tests)
go test ./tests/acceptance/blueprints/grpc-pure -v -short

# Run specific test case
go test ./tests/acceptance/blueprints/grpc-pure -v -run "TestGRPCPureBlueprintBasicGeneration/Basic_gRPC_Pure_Service"
```

### Performance Testing

```bash
# Run performance validation tests
go test ./tests/acceptance/blueprints/grpc-pure -v -run TestGRPCPurePerformanceScenarios

# Run cross-environment compatibility tests
go test ./tests/acceptance/blueprints/grpc-pure -v -run TestGRPCPureCrossEnvironmentCompatibility
```

## Helper Functions

### Available Utilities (`tests/helpers/grpc_utils.go`)

- **`ValidateProtocolBuffers(projectPath)`** - Validates proto files using buf
- **`GenerateProtocolBuffers(projectPath)`** - Generates Go code from proto files
- **`TestGRPCServiceHealth(projectPath, port)`** - Tests gRPC health check endpoint
- **`StartGRPCService(projectPath, port)`** - Starts service for testing
- **`LoadTestGRPCService(projectPath, port, requests, concurrency)`** - Performs load testing
- **`ValidateGRPCInterceptors(projectPath, expectedInterceptors)`** - Validates interceptor configuration
- **`CheckGRPCDependencies(projectPath, expectedDeps)`** - Verifies go.mod dependencies
- **`CompileAndTestGRPCProject(projectPath)`** - Full compilation and test cycle

## Test Configuration

### Environment Variables

```bash
# Skip integration tests
export SKIP_INTEGRATION_TESTS=true

# Set custom gRPC port for testing
export TEST_GRPC_PORT=50051

# Enable verbose output
export VERBOSE_TESTS=true
```

### Test Timeouts

- **Basic Tests**: 2 minutes per test case
- **Compilation Tests**: 5 minutes for complex configurations
- **Integration Tests**: 10 minutes for full scenarios
- **Performance Tests**: 15 minutes including load testing

## Troubleshooting

### Common Issues

1. **buf tool not found**
   ```bash
   go install github.com/bufbuild/buf/cmd/buf@latest
   ```

2. **Docker build fails in CI**
   - This is expected in some CI environments
   - Tests log warnings but continue

3. **Cross-platform build failures**
   - Check Go version compatibility
   - Ensure CGO requirements are met for database drivers

4. **Load test failures**
   - Verify ghz tool is installed
   - Check firewall settings for test ports

### Test Debugging

Enable verbose output:
```bash
go test ./tests/acceptance/blueprints/grpc-pure -v -args -verbose
```

Check test logs:
```bash
go test ./tests/acceptance/blueprints/grpc-pure -v 2>&1 | tee test.log
```

## Success Criteria

### Test Passing Requirements

- ✅ All blueprint variants generate successfully
- ✅ Generated projects compile without errors
- ✅ Protocol buffers validate with buf lint
- ✅ All conditional dependencies are correctly included
- ✅ Interceptor chains are properly configured
- ✅ Authentication mechanisms work as expected
- ✅ Observability features are functional
- ✅ Database integrations are properly set up
- ✅ Service discovery mechanisms are configured
- ✅ Docker containers build successfully
- ✅ Cross-platform builds work correctly

### Performance Benchmarks

- **Generation Time**: < 30 seconds for complex configurations
- **Compilation Time**: < 2 minutes for full project
- **Docker Build**: < 5 minutes with all dependencies
- **Test Execution**: < 10 minutes for full test suite

## Contributing

### Adding New Test Cases

1. **Feature File**: Add scenarios to `features/grpc-pure.feature`
2. **Step Definitions**: Implement steps in `grpc_pure_steps_test.go`
3. **Integration Tests**: Add complex scenarios to `integration_test.go`
4. **Validation Functions**: Add helpers to validate new features

### Test Naming Conventions

- Test functions: `TestGRPCPure[Feature][Scenario]`
- Validation functions: `validate[Component][Aspect]`
- Helper functions: `[action][Component]`

### Code Quality

- All test functions must include `t.Helper()` for helper functions
- Use `require` for setup failures, `assert` for validations  
- Include clear error messages with context
- Clean up resources in defer statements
- Document complex test logic with comments

## Integration with CI/CD

These tests are designed to run in continuous integration pipelines:

- **GitHub Actions**: Integrated with existing workflow
- **Parallel Execution**: Tests can run concurrently
- **Artifact Collection**: Generated projects saved for debugging
- **Failure Reporting**: Detailed logs and error messages
- **Cross-Platform**: Tests run on Linux, macOS, and Windows

## Future Enhancements

### Planned Test Additions

- [ ] gRPC streaming tests (bidirectional, client/server streaming)
- [ ] Load balancing configuration validation
- [ ] Circuit breaker integration testing
- [ ] Multi-region deployment scenarios
- [ ] Performance regression testing
- [ ] Security vulnerability scanning
- [ ] API compatibility testing across versions