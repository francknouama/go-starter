# Monolith Blueprint ATDD Test Suite

This directory contains comprehensive Acceptance Test-Driven Development (ATDD) tests for the monolith blueprint, ensuring it generates production-ready Go web applications that meet all user acceptance criteria.

## Overview

The ATDD test suite validates that the monolith blueprint:
- Generates complete, working web applications
- Supports multiple frameworks, databases, and authentication methods
- Includes production-ready features (Docker, CI/CD, monitoring)
- Follows modular monolith architectural patterns
- Provides comprehensive testing infrastructure
- Implements security best practices
- Includes performance optimizations

## Test Structure

### Test Files

- **`monolith_atdd_test.go`** - Core ATDD tests using Go's standard testing framework
- **`monolith_acceptance_test.go`** - High-level acceptance tests with comprehensive scenarios
- **`monolith_steps_test.go`** - BDD-style step definitions with Gherkin syntax
- **`features/monolith.feature`** - Gherkin feature files defining user stories

### Test Categories

#### 1. Generation and Compilation Tests
- Verifies blueprint is available and generates successfully
- Ensures generated code compiles without errors
- Validates essential file structure and components

#### 2. Multi-Framework Support Tests
- Tests generation with Gin, Echo, Fiber, and Chi frameworks
- Verifies framework-specific imports and configurations
- Ensures each framework variant compiles and runs

#### 3. Database Driver Tests
- Tests PostgreSQL, MySQL, and SQLite support
- Verifies GORM and sqlx ORM integrations
- Uses testcontainers for real database testing
- Validates migration files and database configurations

#### 4. Authentication Tests
- Tests session, JWT, and OAuth2 authentication types
- Verifies secure session configuration (HttpOnly, Secure, SameSite)
- Validates authentication routes and middleware

#### 5. Production Readiness Tests
- Docker configuration and multi-stage builds
- CI/CD pipelines with comprehensive testing
- Kubernetes deployment configurations
- Security scanning and vulnerability checks
- Health checks and monitoring endpoints

#### 6. Asset Pipeline Tests
- Vite configuration for modern asset building
- Tailwind CSS integration
- Production optimization settings
- Asset compilation and minification

#### 7. Architectural Pattern Tests
- Modular monolith structure validation
- Layer separation and dependency flow
- Interface-based design patterns
- Testability and mocking capabilities

#### 8. Database Migration Tests
- Migration script functionality
- Schema versioning and rollbacks
- Database connection testing
- Migration execution validation

#### 9. Testing Infrastructure Tests
- Unit test generation with testify
- Integration test setup
- Benchmark test implementation
- Mock generation and usage

#### 10. Security Tests
- OWASP security header implementation
- Input validation and sanitization
- Session security configuration
- Automated security scanning setup

#### 11. Performance Tests
- Database connection pooling
- Caching strategy implementation
- Asset optimization
- Benchmark test validation

#### 12. Development Workflow Tests
- Makefile target validation
- Development tool configuration
- Hot reload setup (Air)
- Linting configuration (golangci-lint)

### Using Testcontainers

The test suite leverages testcontainers for realistic database testing:

```go
// PostgreSQL container for integration tests
postgresContainer := testcontainers.ContainerRequest{
    Image:        "postgres:16-alpine",
    ExposedPorts: []string{"5432/tcp"},
    Env: map[string]string{
        "POSTGRES_DB":       "testdb",
        "POSTGRES_USER":     "testuser", 
        "POSTGRES_PASSWORD": "testpass",
    },
    WaitingFor: wait.ForLog("database system is ready to accept connections"),
}
```

Benefits of using testcontainers:
- Real database environment testing
- Isolation between test runs
- Automatic cleanup
- Support for multiple database types
- Consistent test environments

## Running the Tests

### Prerequisites

- Go 1.21+
- Docker (for testcontainers)
- Make (optional, for convenience commands)

### Command Examples

```bash
# Run all ATDD tests
go test ./tests/acceptance/blueprints/monolith/... -v

# Run specific test category
go test ./tests/acceptance/blueprints/monolith/ -run TestMonolithAcceptance_BasicGeneration -v

# Run BDD tests with Gherkin features
go test ./tests/acceptance/blueprints/monolith/ -run TestMonolithBDD -v

# Run tests with database integration (requires Docker)
go test ./tests/acceptance/blueprints/monolith/ -run TestMonolithAcceptance_DatabaseDriverSupport -v

# Run performance and load tests
go test ./tests/acceptance/blueprints/monolith/ -run Performance -v -timeout 10m

# Skip slow tests
go test ./tests/acceptance/blueprints/monolith/ -short

# Run tests in parallel
go test ./tests/acceptance/blueprints/monolith/ -v -parallel 4
```

### Environment Variables

- `TEST_TIMEOUT` - Override default test timeout (default: 5m)
- `SKIP_DOCKER_TESTS` - Skip testcontainer-based tests
- `VERBOSE_OUTPUT` - Enable verbose test output
- `KEEP_TEST_ARTIFACTS` - Don't cleanup generated test projects

### CI/CD Integration

The ATDD tests are integrated into the CI pipeline:

```yaml
- name: Run Monolith ATDD Tests
  run: |
    go test ./tests/acceptance/blueprints/monolith/... -v -timeout 10m
    go test ./tests/acceptance/blueprints/monolith/ -run TestMonolithBDD -v
```

## Test Scenarios Coverage

### User Stories Validated

1. **As a Go developer, I want to generate a monolith web application**
   - ✅ Basic project generation
   - ✅ Essential file structure
   - ✅ Compilation success

2. **As a developer, I want to choose my preferred web framework**
   - ✅ Gin framework support
   - ✅ Echo framework support  
   - ✅ Fiber framework support
   - ✅ Chi framework support

3. **As a developer, I want to use different databases**
   - ✅ PostgreSQL with GORM/sqlx
   - ✅ MySQL with GORM/sqlx
   - ✅ SQLite with GORM
   - ✅ Connection pooling
   - ✅ Migration system

4. **As a developer, I want secure authentication**
   - ✅ Session-based authentication
   - ✅ JWT authentication
   - ✅ OAuth2 integration
   - ✅ Secure session configuration

5. **As a developer, I want production-ready deployment**
   - ✅ Docker containerization
   - ✅ CI/CD pipelines
   - ✅ Kubernetes deployment
   - ✅ Health checks
   - ✅ Security scanning

6. **As a developer, I want modern frontend assets**
   - ✅ Vite integration
   - ✅ Tailwind CSS setup
   - ✅ Asset optimization
   - ✅ Legacy browser support

7. **As a developer, I want maintainable architecture**
   - ✅ Modular monolith patterns
   - ✅ Layer separation
   - ✅ Interface-based design
   - ✅ Comprehensive testing

8. **As a developer, I want database management**
   - ✅ Migration scripts
   - ✅ Schema versioning
   - ✅ Rollback capabilities
   - ✅ Seed data support

9. **As a developer, I want quality assurance**
   - ✅ Unit test generation
   - ✅ Integration tests
   - ✅ Benchmark tests
   - ✅ Code coverage

10. **As a developer, I want security best practices**
    - ✅ OWASP security headers
    - ✅ Input validation
    - ✅ Secure sessions
    - ✅ Automated scanning

### Quality Gates

Each test validates specific quality gates:

- **Functionality**: Generated code compiles and runs
- **Architecture**: Follows modular monolith patterns
- **Security**: Implements OWASP best practices
- **Performance**: Includes optimization strategies
- **Maintainability**: Comprehensive test coverage
- **Operations**: Production deployment ready

## Extending the Tests

### Adding New Scenarios

1. Create a new test function in `monolith_acceptance_test.go`
2. Follow the naming convention: `TestMonolithAcceptance_FeatureName`
3. Use the `MonolithAcceptanceTestSuite` for setup/teardown
4. Add corresponding Gherkin scenarios in `features/monolith.feature`

### Adding BDD Steps

1. Add new step definitions in `monolith_steps_test.go`
2. Register steps in `InitializeMonolithScenario`
3. Implement step functions following the pattern:
   - Given: Preconditions and setup
   - When: Actions and operations
   - Then: Assertions and validations

### Testing New Database Types

1. Add testcontainer support for the new database
2. Create setup/cleanup methods
3. Add database-specific test scenarios
4. Update the database driver test matrix

### Performance Testing

1. Add benchmark functions with `func BenchmarkMonolith...`
2. Use `b.ResetTimer()` and `b.StopTimer()` appropriately
3. Include performance assertions
4. Test under realistic load conditions

## Best Practices

### Test Organization
- Group related tests in the same file
- Use descriptive test names
- Include both positive and negative test cases
- Test edge cases and error conditions

### Test Data Management
- Use testcontainers for database state
- Clean up after each test scenario
- Use realistic test data
- Avoid hard-coded values when possible

### Performance Considerations
- Use `testing.Short()` to skip slow tests
- Run database tests in parallel when safe
- Clean up containers promptly
- Use appropriate timeouts

### Error Handling
- Provide detailed error messages
- Include context in assertions
- Log intermediate steps for debugging
- Use `require` for fatal errors, `assert` for validation

## Troubleshooting

### Common Issues

1. **Docker not available**
   ```
   Error: Cannot connect to Docker daemon
   Solution: Ensure Docker is running and accessible
   ```

2. **Port conflicts**
   ```
   Error: Port already in use
   Solution: Testcontainers automatically handles port allocation
   ```

3. **Timeout errors**
   ```
   Error: Container startup timeout
   Solution: Increase timeout or check Docker performance
   ```

4. **Build failures**
   ```
   Error: Generated project doesn't compile
   Solution: Check go.mod dependencies and template syntax
   ```

### Debug Mode

Enable verbose output for debugging:

```bash
VERBOSE_OUTPUT=1 go test ./tests/acceptance/blueprints/monolith/ -v -run TestMonolithAcceptance_BasicGeneration
```

Keep test artifacts for inspection:

```bash
KEEP_TEST_ARTIFACTS=1 go test ./tests/acceptance/blueprints/monolith/ -v
```

## Contributing

When contributing to the ATDD test suite:

1. Follow the existing test patterns and naming conventions
2. Add comprehensive test coverage for new features
3. Include both positive and negative test scenarios
4. Update documentation for new test categories
5. Ensure tests are deterministic and don't depend on external state
6. Use testcontainers for realistic integration testing
7. Add appropriate cleanup and error handling

The ATDD test suite is crucial for ensuring the monolith blueprint generates high-quality, production-ready Go applications. Comprehensive test coverage helps maintain quality and prevents regressions as the blueprint evolves.