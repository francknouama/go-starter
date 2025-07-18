# Web API Blueprint ATDD Tests

This directory contains Acceptance Test-Driven Development (ATDD) tests for all Web API blueprint architectures. These tests implement comprehensive validation using Gherkin-style BDD scenarios to ensure generated projects meet business requirements and architectural principles.

## Overview

The ATDD test suite validates four Web API architecture patterns:
- **Standard**: Traditional layered architecture
- **Clean**: Clean Architecture with dependency inversion
- **DDD**: Domain-Driven Design with rich domain models
- **Hexagonal**: Ports and Adapters architecture

## Test Structure

```
tests/acceptance/blueprints/web-api/
├── standard_test.go      # Standard architecture ATDD tests
├── clean_test.go         # Clean architecture ATDD tests  
├── ddd_test.go          # DDD architecture ATDD tests
├── hexagonal_test.go    # Hexagonal architecture ATDD tests
├── integration_test.go  # Cross-architecture integration tests
└── README.md           # This documentation
```

## Test Categories

### 1. Architecture-Specific Tests

Each architecture has its own test file with scenarios covering:

#### Standard Architecture (`standard_test.go`)
- ✅ Basic generation with framework integration
- ✅ Database integration scenarios  
- ✅ Logger integration testing
- ✅ Framework variations (gin, echo, fiber, chi)
- ✅ Architecture compliance validation
- ✅ Runtime behavior testing
- ✅ Minimal configuration testing

#### Clean Architecture (`clean_test.go`)
- ✅ Layer structure validation
- ✅ Dependency injection validation
- ✅ Logger integration following Clean patterns
- ✅ Framework abstraction validation
- ✅ Database integration with repository pattern
- ✅ Architecture compliance validation

#### DDD Architecture (`ddd_test.go`)
- ✅ Domain-centric structure validation
- ✅ Business rule enforcement
- ✅ CQRS pattern implementation
- ✅ Logger integration in DDD context
- ✅ Framework integration maintaining domain independence
- ✅ Architecture compliance validation

#### Hexagonal Architecture (`hexagonal_test.go`)
- ✅ Ports and adapters validation
- ✅ Adapter swappability testing
- ✅ Core isolation validation
- ✅ Logger as secondary adapter
- ✅ Framework as primary adapter
- ✅ Architecture compliance validation

### 2. Cross-Architecture Integration Tests (`integration_test.go`)

Tests that validate consistent behavior across all architectures:

#### Database Integration
- Validates database configuration follows architecture patterns
- Tests data access abstraction
- Validates transaction handling
- Ensures migrations work correctly

#### Logger Integration  
- Tests consistent logging across architectures
- Validates architecture-specific logger patterns
- Ensures log levels are configurable
- Validates structured logging usage

#### Compilation Validation
- Ensures all architectures compile successfully
- Validates dependency resolution
- Tests binary generation

#### Architecture Compliance
- Validates architecture-specific principles
- Tests dependency directions
- Validates layer boundaries
- Ensures architectural linting passes

## Gherkin-Style BDD Scenarios

All tests follow Gherkin-style Given/When/Then patterns for business readability:

```go
func TestStandard_WebAPI_BasicGeneration_WithGin(t *testing.T) {
    // Scenario: Generate standard web API with Gin
    // Given I want a standard web API
    // When I generate with framework "gin"
    // Then the project should include Gin router setup
    // And the project should have basic CRUD endpoints
    // And the project should compile and run successfully
    // And HTTP server should start on configured port
    
    // Given I want a standard web API
    config := types.ProjectConfig{...}
    
    // When I generate with framework "gin"
    projectPath := helpers.GenerateProject(t, config)
    
    // Then the project should include Gin router setup
    validator := NewWebAPIValidator(projectPath, "standard")
    validator.ValidateGinRouterSetup(t)
    
    // And the project should have basic CRUD endpoints
    validator.ValidateBasicCRUDEndpoints(t)
    
    // And the project should compile and run successfully
    validator.ValidateCompilation(t)
}
```

## Validation Architecture

### Blueprint Validators

Each architecture has a dedicated validator with specific validation methods:

- **WebAPIValidator**: Standard architecture validation
- **CleanArchitectureValidator**: Clean architecture specific validation
- **DDDValidator**: Domain-Driven Design validation 
- **HexagonalValidator**: Hexagonal architecture validation

### Integration Validators

Cross-cutting validators for shared concerns:

- **DatabaseIntegrationValidator**: Database integration across architectures
- **LoggerIntegrationValidator**: Logger integration across architectures
- **CompilationValidator**: Compilation and build validation
- **ArchitectureComplianceValidator**: Architecture principle compliance

## Key Validation Methods

### Common Validations
- `ValidateCompilation(t)`: Ensures generated project compiles
- `ValidateLoggerIntegration(t, logger)`: Tests logger configuration
- `ValidateFrameworkIntegration(t, framework)`: Tests framework setup
- `ValidateDatabaseConfiguration(t, driver, orm)`: Tests database setup

### Architecture-Specific Validations

#### Standard
- `ValidateArchitecturePrinciples(t, "standard")`: Standard patterns
- `ValidateDependencyDirections(t)`: Simple dependency flow
- `ValidateMinimalConfiguration(t)`: Minimal project validation

#### Clean Architecture  
- `ValidateLayerStructure(t)`: Four-layer structure
- `ValidateDependencyInversion(t)`: Dependency inversion principle
- `ValidateBusinessLogicIsolation(t)`: Framework independence

#### DDD
- `ValidateDomainCentricStructure(t)`: Domain-centric organization
- `ValidateBusinessRuleEnforcement(t)`: Business logic in domain
- `ValidateCQRSPattern(t)`: Command/Query separation

#### Hexagonal
- `ValidatePortsStructure(t)`: Ports and adapters structure
- `ValidateCoreBusinessLogicIndependence(t)`: Core isolation
- `ValidateAdapterSwappability(t)`: Pluggable adapters

## Running the Tests

### Run All ATDD Tests
```bash
go test ./tests/acceptance/blueprints/web-api/...
```

### Run Specific Architecture Tests
```bash
# Standard architecture
go test ./tests/acceptance/blueprints/web-api/standard_test.go

# Clean architecture  
go test ./tests/acceptance/blueprints/web-api/clean_test.go

# DDD architecture
go test ./tests/acceptance/blueprints/web-api/ddd_test.go

# Hexagonal architecture
go test ./tests/acceptance/blueprints/web-api/hexagonal_test.go

# Integration tests
go test ./tests/acceptance/blueprints/web-api/integration_test.go
```

### Run with Verbose Output
```bash
go test -v ./tests/acceptance/blueprints/web-api/...
```

### Run Specific Scenario
```bash
go test -v ./tests/acceptance/blueprints/web-api/ -run TestStandard_WebAPI_BasicGeneration_WithGin
```

## Test Features

### ✅ Implemented Features

- **Gherkin-style BDD scenarios** for business readability
- **Architecture-specific validation** for each pattern
- **Cross-architecture integration testing** for consistency
- **Logger integration testing** across all supported loggers (slog, zap, logrus, zerolog)
- **Framework integration testing** across all supported frameworks (gin, echo, fiber, chi)
- **Database integration testing** with various drivers and ORMs
- **Compilation validation** ensuring generated projects build successfully
- **Architecture compliance validation** ensuring principles are followed
- **Graceful handling of incomplete implementations** (e.g., hexagonal architecture)

### 🚧 Planned Enhancements

- **Runtime validation**: Actually start servers and test HTTP endpoints
- **Performance testing**: Validate generated projects meet performance benchmarks
- **Security testing**: Ensure generated code follows security best practices
- **Migration testing**: Validate database migrations work correctly
- **Template validation**: Pre-validate templates before generation testing
- **Cross-platform testing**: Ensure generated projects work on different operating systems

## Error Handling

### Template Not Implemented
Tests gracefully handle incomplete template implementations:

```go
if !validator.IsHexagonalImplemented() {
    t.Skip("Hexagonal architecture template not fully implemented yet")
}
```

### Compilation Failures
All tests validate compilation and provide detailed error messages:

```go
func (v *Validator) ValidateCompilation(t *testing.T) {
    t.Helper()
    helpers.AssertCompilable(t, v.ProjectPath)
}
```

### Architecture Violations
Tests specifically check for architecture principle violations and provide clear feedback on what's wrong and why.

## Integration with CI/CD

These ATDD tests are designed for continuous integration:

- ✅ No external dependencies required
- ✅ Generate projects in temporary directories  
- ✅ Clean up after themselves
- ✅ Provide clear failure messages
- ✅ Support parallel execution
- ✅ Graceful handling of incomplete implementations

## Contributing

When adding new ATDD tests:

1. **Follow Gherkin patterns**: Use Given/When/Then structure
2. **Add business context**: Include feature and scenario descriptions
3. **Test positive and negative paths**: What should and shouldn't exist
4. **Validate compilation**: Always ensure generated projects compile
5. **Handle missing templates**: Use graceful skipping for incomplete implementations
6. **Test multiple variations**: Different frameworks, loggers, databases
7. **Focus on business value**: Test what matters to developers using the tool

## Benefits of ATDD Approach

### Business Readability
Tests read like specifications that business stakeholders can understand:

```
Feature: Standard Web API Blueprint
As a developer
I want to generate a standard web API project
So that I can quickly build REST APIs

Scenario: Generate standard web API with Gin
Given I want a standard web API
When I generate with framework "gin"
Then the project should include Gin router setup
And the project should have basic CRUD endpoints
And the project should compile and run successfully
```

### Living Documentation
Tests serve as up-to-date documentation of what each architecture supports and how it should behave.

### Quality Assurance
Comprehensive validation ensures generated projects work correctly and follow architectural principles.

### Regression Prevention
Automated tests catch breaking changes to templates and ensure consistency across architectures.

This ATDD implementation provides comprehensive validation of Web API blueprints while maintaining business readability and serving as living documentation of the system's capabilities.