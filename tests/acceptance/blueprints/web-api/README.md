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
├── features/                           # Gherkin feature files (BDD scenarios)
│   ├── web-api.feature                 # General web API scenarios
│   ├── clean-architecture.feature     # Clean Architecture specific scenarios
│   ├── domain-driven-design.feature   # DDD specific scenarios  
│   ├── hexagonal-architecture.feature # Hexagonal Architecture scenarios
│   ├── standard-architecture.feature  # Standard Architecture scenarios
│   └── integration-testing.feature    # Cross-architecture integration scenarios
├── web_api_steps_test.go              # BDD step definitions for all scenarios
├── web_api_acceptance_test.go         # High-level acceptance tests
├── standard_test.go                   # Legacy: Standard architecture ATDD tests
├── clean_test.go                      # Legacy: Clean architecture ATDD tests  
├── ddd_test.go                        # Legacy: DDD architecture ATDD tests
├── hexagonal_test.go                  # Legacy: Hexagonal architecture ATDD tests
├── integration_test.go                # Legacy: Cross-architecture integration tests
└── README.md                          # This documentation
```

## Normalized BDD Structure ✨

**NEW**: The web-api ATDD tests have been normalized with a comprehensive BDD structure using Gherkin feature files and unified step definitions.

### Feature File Organization

The tests are now organized into dedicated feature files that follow BDD best practices:

#### 🏗️ **Architecture-Specific Features**
- **`clean-architecture.feature`** (9 scenarios): Clean Architecture patterns, dependency inversion, business logic isolation
- **`domain-driven-design.feature`** (15 scenarios): DDD patterns, entities, value objects, aggregates, domain events
- **`hexagonal-architecture.feature`** (15 scenarios): Ports & adapters, dependency direction, framework independence
- **`standard-architecture.feature`** (20 scenarios): Traditional layered architecture, RESTful endpoints, standard patterns

#### 🔗 **Cross-Cutting Features**  
- **`integration-testing.feature`** (15 scenarios): Database integration, logger integration, framework integration across all architectures
- **`web-api.feature`** (20+ scenarios): General web API scenarios, authentication, security, deployment, monitoring

### Unified Step Definitions

All feature files share a common set of step definitions in `web_api_steps_test.go`:

```go
// Architecture-specific Given steps
ctx.Given(`^I want to create a Clean Architecture web API$`, webApiCtx.iWantToCreateACleanArchitectureWebAPI)
ctx.Given(`^I want to create a DDD web API$`, webApiCtx.iWantToCreateADDDWebAPI)
ctx.Given(`^I want to create a Hexagonal Architecture web API$`, webApiCtx.iWantToCreateAHexagonalArchitectureWebAPI)

// Integration testing steps
ctx.When(`^I generate a web API with architecture "([^"]*)", database "([^"]*)", and ORM "([^"]*)"$`, 
         webApiCtx.iGenerateAWebAPIWithArchitectureDatabaseAndORM)

// Validation steps
ctx.Then(`^database configuration should follow architecture patterns$`, 
         webApiCtx.databaseConfigurationShouldFollowArchitecturePatterns)
```

### Benefits of Normalization

#### 🎯 **Improved Maintainability**
- Single source of step definitions across all architectures
- Consistent scenario structure and validation logic
- Reduced code duplication between architecture tests

#### 📋 **Enhanced Readability** 
- Clear separation of concerns between architectures
- Business-readable Gherkin scenarios
- Focused feature files for specific architectural concerns

#### 🔄 **Better Test Organization**
- Architecture-specific scenarios grouped logically
- Cross-cutting concerns in dedicated integration feature
- Easy to add new scenarios for specific architectures

#### 🚀 **Scalability**
- Easy to add new architectures by creating new feature files
- Shared step definitions reduce implementation effort
- Consistent patterns across all architecture tests

## Test Categories

### 1. BDD Feature Files (Primary)

Each feature file contains comprehensive Gherkin scenarios:

#### Clean Architecture (`clean-architecture.feature`)
- ✅ Layer structure validation (entities, use cases, adapters, frameworks)
- ✅ Dependency inversion principle enforcement  
- ✅ Business logic isolation from external concerns
- ✅ Framework abstraction validation
- ✅ Database integration with repository pattern
- ✅ Logger integration following Clean patterns
- ✅ Dependency injection configuration
- ✅ Interface-based design validation
- ✅ Architecture compliance validation

#### Domain-Driven Design (`domain-driven-design.feature`)
- ✅ Domain-centric structure validation
- ✅ Entities and value objects implementation
- ✅ Aggregate design and consistency boundaries
- ✅ Domain events and event handlers
- ✅ Domain services vs application services
- ✅ Repository pattern in domain context
- ✅ CQRS pattern implementation
- ✅ Ubiquitous language in code
- ✅ Bounded context separation
- ✅ Anti-corruption layers
- ✅ Business rule enforcement
- ✅ Domain model purity
- ✅ Infrastructure separation
- ✅ Strategic design patterns
- ✅ Tactical design patterns

#### Hexagonal Architecture (`hexagonal-architecture.feature`)
- ✅ Ports and adapters structure validation
- ✅ Primary adapters (driving) implementation
- ✅ Secondary adapters (driven) implementation  
- ✅ Application core isolation
- ✅ Port interface validation
- ✅ Dependency direction enforcement
- ✅ Framework independence validation
- ✅ Database independence validation
- ✅ Multiple adapter support
- ✅ Testing strategy validation
- ✅ Dependency injection configuration
- ✅ Error handling across layers
- ✅ Cross-cutting concerns handling
- ✅ DDD integration capabilities
- ✅ Adapter swappability

#### Standard Architecture (`standard-architecture.feature`)
- ✅ Traditional layered structure
- ✅ RESTful endpoint generation
- ✅ Middleware configuration
- ✅ Request validation
- ✅ Error handling
- ✅ Database connection and models
- ✅ Configuration management
- ✅ Authentication integration
- ✅ Testing infrastructure
- ✅ API documentation
- ✅ Health check implementation
- ✅ Performance monitoring
- ✅ Security best practices
- ✅ Container deployment
- ✅ Framework variations (gin, echo, fiber, chi, stdlib)
- ✅ Database variations (postgres, mysql, sqlite)
- ✅ ORM variations (gorm, sqlx)
- ✅ Logger variations (slog, zap, logrus, zerolog)
- ✅ Authentication types (jwt, session, api-key)

#### Integration Testing (`integration-testing.feature`)
- ✅ Database integration across architectures
- ✅ Logger integration across architectures
- ✅ Framework integration across architectures
- ✅ Authentication integration across architectures
- ✅ Multi-feature integration testing
- ✅ Database migration integration
- ✅ Container integration testing
- ✅ CI/CD integration validation
- ✅ Performance integration testing
- ✅ Security integration testing
- ✅ Error handling integration
- ✅ Testing infrastructure integration
- ✅ Configuration integration testing
- ✅ API documentation integration
- ✅ Monitoring and observability integration

### 2. Legacy Architecture-Specific Tests (Preserved)

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
# Run all web-api ATDD tests (BDD + Legacy)
go test ./tests/acceptance/blueprints/web-api/...
```

### Run BDD Feature Tests (Primary)
```bash
# Run comprehensive BDD tests with all feature files
go test ./tests/acceptance/blueprints/web-api/ -run TestWebAPIBDD

# Run with verbose output to see all scenarios
go test -v ./tests/acceptance/blueprints/web-api/ -run TestWebAPIBDD
```

### Run Specific Feature Files
```bash
# Test specific architecture using tags (if implemented)
go test ./tests/acceptance/blueprints/web-api/ -tags=clean-architecture
go test ./tests/acceptance/blueprints/web-api/ -tags=hexagonal-architecture

# Or run with godog directly for specific features
cd tests/acceptance/blueprints/web-api/
godog features/clean-architecture.feature
godog features/hexagonal-architecture.feature
godog features/integration-testing.feature
```

### Run Legacy Architecture Tests
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
# All tests with verbose output
go test -v ./tests/acceptance/blueprints/web-api/...

# BDD tests only with scenario details
go test -v ./tests/acceptance/blueprints/web-api/ -run TestWebAPIBDD
```

### Run Specific Scenarios
```bash
# Legacy specific scenario
go test -v ./tests/acceptance/blueprints/web-api/ -run TestStandard_WebAPI_BasicGeneration_WithGin

# BDD specific scenario (depends on godog filtering)
go test -v ./tests/acceptance/blueprints/web-api/ -run TestWebAPIBDD | grep "Clean Architecture"
```

### Debug Feature Files
```bash
# Validate feature file syntax
cd tests/acceptance/blueprints/web-api/
godog --no-colors --format=progress features/

# Run single feature file
godog features/clean-architecture.feature

# Run with specific tags
godog --tags="@database" features/
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

## Normalization Summary

The web-api ATDD tests have been successfully normalized from embedded scenarios to a comprehensive BDD structure:

### What Was Normalized

**Before**: Embedded Gherkin scenarios in Go test comments
```go
// Scenario: Clean Architecture web API generation
//   Given I want to create a Clean Architecture web API
//   When I generate a web API with Clean Architecture
//   Then the project should follow Clean Architecture principles
```

**After**: Dedicated Gherkin feature files with unified step definitions
```gherkin
# features/clean-architecture.feature
Feature: Clean Architecture Web API
  As a software architect
  I want to generate Clean Architecture web API
  So that I can maintain separation of concerns

Scenario: Clean Architecture web API generation
  Given I want to create a Clean Architecture web API
  When I generate a web API with Clean Architecture
  Then the project should follow Clean Architecture principles
```

### Files Created During Normalization

1. **`features/clean-architecture.feature`** - 9 Clean Architecture scenarios
2. **`features/domain-driven-design.feature`** - 15 DDD scenarios  
3. **`features/hexagonal-architecture.feature`** - 15 Hexagonal Architecture scenarios
4. **`features/standard-architecture.feature`** - 20 Standard Architecture scenarios
5. **`features/integration-testing.feature`** - 15 Integration testing scenarios
6. **Updated `web_api_steps_test.go`** - Unified step definitions for all feature files

### Extraction Sources

Scenarios were extracted from embedded comments in:
- `clean_test.go` → `clean-architecture.feature`
- `ddd_test.go` → `domain-driven-design.feature`  
- `hexagonal_test.go` → `hexagonal-architecture.feature`
- `standard_test.go` → `standard-architecture.feature`
- `integration_test.go` → `integration-testing.feature`

### Benefits Achieved

✅ **Improved Organization**: Architecture-specific concerns properly separated
✅ **Enhanced Readability**: Business stakeholders can read feature files directly  
✅ **Reduced Duplication**: Single source of step definitions across all architectures
✅ **Better Maintainability**: Changes to step logic only need to be made in one place
✅ **Scalable Structure**: Easy to add new architectures or extend existing ones
✅ **Standards Compliance**: Follows BDD/Gherkin best practices
✅ **Tool Integration**: Compatible with standard Gherkin tooling and IDEs

This normalization work represents a significant improvement in the web-api ATDD test structure, moving from embedded comments to a professional BDD implementation that serves as both comprehensive testing and living documentation of the system's capabilities.