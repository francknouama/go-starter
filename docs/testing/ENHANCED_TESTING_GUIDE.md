# Enhanced Testing Guide

## Overview

This guide provides comprehensive information about the enhanced testing infrastructure in go-starter, including the ATDD system, test execution, debugging, and best practices for maintaining our 100% blueprint validation coverage.

## 🧪 Enhanced ATDD System

### What is ATDD?

Acceptance Test-Driven Development (ATDD) is a collaborative approach where acceptance criteria are defined before implementation. Our enhanced ATDD system provides:

- **114+ comprehensive scenarios** covering all blueprint combinations
- **100% compilation guarantee** for generated projects
- **Architecture validation** using AST parsing
- **Cross-platform compatibility** testing

### Test Categories

#### 1. P0 Critical Tests (Production Blockers)
- **Cross-Blueprint Integration**: 5 scenarios
- **Enterprise Architecture Matrix**: 15 scenarios  
- **Database Integration Matrix**: 15 scenarios
- **Authentication System Matrix**: 15 scenarios

#### 2. P1 High-Priority Tests (Quality Assurance)
- **Lambda Deployment Scenarios**: 26 scenarios
- **Framework Consistency Validation**: 12 scenarios
- **CLI Complexity Testing**: 8 scenarios
- **Database Matrix Expansion**: 18 scenarios

## 🚀 Running Tests

### Quick Commands

```bash
# Run all tests (unit + integration + ATDD)
go test -v ./...

# Run only ATDD tests
go test -v ./tests/acceptance/enhanced/...

# Run specific test category
go test -v ./tests/acceptance/enhanced/architecture/...
go test -v ./tests/acceptance/enhanced/database/...
go test -v ./tests/acceptance/enhanced/lambda/...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Advanced Test Execution

#### Parallel Execution
```bash
# Run tests in parallel (faster execution)
go test -v -parallel=4 ./tests/acceptance/enhanced/...

# Limit parallelism for resource-intensive tests
go test -v -parallel=2 ./tests/acceptance/enhanced/compilation/...
```

#### Filtering Tests
```bash
# Run tests matching a pattern
go test -v -run="TestEnterpriseArchitecture" ./tests/acceptance/enhanced/...

# Run specific blueprint combination
go test -v -run="TestWebAPI.*Hexagonal.*Postgres" ./tests/acceptance/enhanced/...

# Run only compilation tests
go test -v -run=".*Compilation.*" ./tests/acceptance/enhanced/...
```

#### Test Timeouts
```bash
# Set custom timeout for long-running tests
go test -v -timeout=30m ./tests/acceptance/enhanced/...

# Quick smoke test (shorter timeout)
go test -v -timeout=5m -run="TestBasic.*" ./tests/acceptance/enhanced/...
```

## 🔍 Test Structure and Organization

### Feature File Format

Our tests use Gherkin-style feature files for readability:

```gherkin
# tests/acceptance/enhanced/architecture/features/architecture-validation.feature

Feature: Architecture Pattern Validation
  As a developer using go-starter
  I want generated projects to follow proper architectural patterns
  So that my code is maintainable and follows best practices

  @critical @architecture
  Scenario: Clean Architecture boundaries are enforced
    Given I generate a web-api project with clean architecture
    When I analyze the generated code structure
    Then the dependency rules should be properly enforced
    And no layer should depend on outer layers
    And interfaces should be properly defined

  @critical @architecture
  Scenario: Hexagonal Architecture ports and adapters
    Given I generate a web-api project with hexagonal architecture
    When I analyze the generated code structure  
    Then ports should be defined as interfaces
    And adapters should implement the ports
    And the core domain should be isolated
```

### Test Implementation Structure

```go
// tests/acceptance/enhanced/architecture/architecture_test.go

func TestArchitectureValidation(t *testing.T) {
    testCases := []TestConfig{
        {
            Type: "web-api",
            Architecture: "clean", 
            Framework: "gin",
            Database: "postgres",
            Description: "Clean Architecture with Gin and PostgreSQL",
        },
        {
            Type: "web-api",
            Architecture: "hexagonal",
            Framework: "echo", 
            Database: "mysql",
            Description: "Hexagonal Architecture with Echo and MySQL",
        },
        // ... more test cases
    }
    
    for _, tc := range testCases {
        t.Run(tc.Description, func(t *testing.T) {
            // Test implementation
            projectPath := generateProject(t, tc)
            validateArchitecture(t, projectPath, tc)
            assertCompilationSuccess(t, projectPath)
        })
    }
}
```

## 🛠️ Test Utilities and Helpers

### Core Test Helpers

#### Project Generation
```go
// Generate a project with specific configuration
projectPath := helpers.GenerateProject(t, helpers.TestConfig{
    Type:         "web-api",
    Architecture: "hexagonal", 
    Framework:    "gin",
    Database:     "postgres",
    Logger:       "zap",
})
```

#### Validation Helpers
```go
// Validate project compilation
helpers.AssertCompilationSuccess(t, projectPath)

// Validate imports
helpers.AssertNoUnusedImports(t, projectPath)

// Validate file structure
helpers.AssertFileExists(t, projectPath, "internal/domain/user.go")
helpers.AssertFileCount(t, projectPath, 25) // For standard CLI

// Validate architecture patterns
helpers.AssertCleanArchitectureLayers(t, projectPath)
helpers.AssertHexagonalBoundaries(t, projectPath)
helpers.AssertDDDDomainModel(t, projectPath)
```

#### Performance Helpers
```go
// Track build performance
buildTime := helpers.MeasureBuildTime(t, projectPath)
helpers.AssertBuildTimeUnder(t, buildTime, 30*time.Second)

// Monitor resource usage
resources := helpers.MonitorResourceUsage(t, func() {
    helpers.AssertCompilationSuccess(t, projectPath)
})
helpers.AssertMemoryUsageUnder(t, resources.Memory, 500*1024*1024) // 500MB
```

### Custom Assertions

#### Architecture-Specific Assertions
```go
// Clean Architecture validation
func AssertCleanArchitectureLayers(t *testing.T, projectPath string) {
    // Verify entity layer exists and has no external dependencies
    // Verify use case layer depends only on entities
    // Verify adapter layer implements interfaces correctly
    // Verify dependency injection setup
}

// Hexagonal Architecture validation  
func AssertHexagonalBoundaries(t *testing.T, projectPath string) {
    // Verify ports are defined as interfaces
    // Verify adapters implement ports
    // Verify core domain isolation
    // Verify dependency inversion
}

// DDD validation
func AssertDDDDomainModel(t *testing.T, projectPath string) {
    // Verify aggregates are properly defined
    // Verify domain services exist
    // Verify bounded contexts are isolated
    // Verify domain events are implemented
}
```

## 🐛 Debugging Test Failures

### Common Failure Scenarios

#### 1. Compilation Failures
```bash
# Debug compilation issues
go test -v -run="TestCompilation" ./tests/acceptance/enhanced/... -args -debug

# Check specific blueprint compilation
cd tests/tmp/generated-project-xyz
go build ./...
```

#### 2. Import Issues
```bash
# Analyze import problems
go test -v -run="TestImports" ./tests/acceptance/enhanced/... -args -verbose

# Manual import check
goimports -d .
go mod tidy
```

#### 3. Template Variable Issues
```bash
# Debug template rendering
go test -v -run="TestVariables" ./tests/acceptance/enhanced/... -args -dump-config

# Check template syntax
go run cmd/debug-templates/main.go --blueprint=web-api-hexagonal
```

### Debug Flags and Options

```bash
# Enable verbose logging
go test -v ./tests/acceptance/enhanced/... -args -debug -verbose

# Keep temporary files for inspection
go test -v ./tests/acceptance/enhanced/... -args -keep-temp

# Dump generated project structure
go test -v ./tests/acceptance/enhanced/... -args -dump-structure

# Enable performance profiling
go test -v ./tests/acceptance/enhanced/... -args -profile -cpuprofile=cpu.prof
```

### Debugging Tips

1. **Check Temporary Directories**: Use `-keep-temp` to inspect generated projects
2. **Enable Debug Logging**: Use `-debug` flag for detailed execution logs
3. **Run Single Test**: Isolate failing tests with specific `-run` patterns
4. **Check Dependencies**: Ensure all required tools are installed
5. **Verify Environment**: Check Go version, OS compatibility, and PATH

## 📊 Coverage and Metrics

### Coverage Reports

```bash
# Generate detailed coverage report
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Coverage by package
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Show only uncovered lines
go tool cover -html=coverage.out -o coverage.html
# Open coverage.html and look for red lines
```

### Performance Metrics

```bash
# Run benchmarks with memory stats
go test -bench=. -benchmem ./tests/benchmarks/...

# Profile memory usage
go test -bench=. -memprofile=mem.prof ./tests/benchmarks/...
go tool pprof mem.prof

# Profile CPU usage  
go test -bench=. -cpuprofile=cpu.prof ./tests/benchmarks/...
go tool pprof cpu.prof
```

### Test Metrics Dashboard

The enhanced ATDD system provides metrics on:

- **Success Rates**: Per blueprint combination
- **Build Times**: Compilation performance tracking
- **Resource Usage**: Memory and CPU consumption
- **Coverage Metrics**: Line and branch coverage
- **Failure Analysis**: Common failure patterns

## 🚀 Performance Optimization

### Test Execution Speed

#### Current Performance
- **Full Suite**: ~15-20 minutes (before optimization)
- **Parallel Execution**: Up to 5 concurrent suites
- **Cache Hit Rate**: 70%+ for repeated configurations

#### Optimization Strategies
```bash
# Use test caching
go test -v -cache ./tests/acceptance/enhanced/...

# Enable parallel execution
go test -v -parallel=4 ./tests/acceptance/enhanced/...

# Run only modified tests
go test -v -short ./tests/acceptance/enhanced/...

# Use build cache
export GOCACHE=$(pwd)/.gocache
go test -v ./tests/acceptance/enhanced/...
```

### Memory Optimization
```bash
# Limit memory usage
export GOMEMLIMIT=2GiB
go test -v ./tests/acceptance/enhanced/...

# Use smaller test subsets for development
go test -v -run="TestBasic.*" ./tests/acceptance/enhanced/...
```

## 🔧 CI/CD Integration

### GitHub Actions Configuration

```yaml
# .github/workflows/enhanced-atdd.yml
name: Enhanced ATDD Tests

on: [push, pull_request]

jobs:
  atdd-tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: [1.21, 1.22]
        test-suite:
          - architecture
          - database  
          - lambda
          - enterprise
    
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
      
      - name: Run ATDD Tests
        run: |
          go test -v -timeout=30m \
            ./tests/acceptance/enhanced/${{ matrix.test-suite }}/...
```

### Quality Gates

- **All tests must pass**: Zero tolerance for failing tests
- **Coverage threshold**: Minimum 85% coverage maintained
- **Performance regression**: Build times must not exceed thresholds
- **Cross-platform validation**: Tests must pass on all supported platforms

## 📈 Best Practices

### Writing New Tests

1. **Follow the Pattern**: Use existing test structure as template
2. **Clear Descriptions**: Write descriptive test names and scenarios
3. **Proper Cleanup**: Always clean up temporary resources
4. **Error Handling**: Provide clear error messages for failures
5. **Documentation**: Document complex test logic

### Test Maintenance

1. **Regular Updates**: Keep tests in sync with blueprint changes
2. **Performance Monitoring**: Track and optimize slow tests
3. **Dependency Management**: Keep test dependencies up to date
4. **Refactoring**: Regularly refactor test code for maintainability

### Contributing Guidelines

1. **Test-First**: Write tests before implementing features
2. **Comprehensive Coverage**: Ensure new features are fully tested
3. **Review Process**: All test changes require peer review
4. **Documentation**: Update documentation when adding new test categories

This enhanced testing system ensures that go-starter maintains the highest quality standards while providing comprehensive validation of all blueprint combinations.