# Architecture Testing Suite

This directory contains comprehensive tests for validating the different architecture patterns supported by the go-starter project generator.

## Overview

The go-starter generator supports multiple architectural patterns, each with distinct organizational principles and dependency flows. These tests ensure that generated projects correctly implement their intended architecture pattern and follow established best practices.

## Supported Architectures

### 1. Standard Architecture (`standard_test.go`)
- **Pattern**: Traditional layered architecture with simple, straightforward organization
- **Structure**: Flat internal structure with handlers, services, repository, models
- **Dependencies**: Straightforward flow: handlers → services → repository → database
- **Use Case**: Rapid prototyping, simple applications, teams new to Go
- **Template**: `web-api-standard`

**Key Validations:**
- File structure follows standard Go project layout
- Dependencies flow in expected direction
- No unnecessary abstraction layers
- Framework-specific handlers are generated correctly
- Database integration is simple and direct

### 2. Clean Architecture (`clean_test.go`)
- **Pattern**: Clean Architecture with strict dependency inversion
- **Structure**: Layered with entities, use cases, interface adapters, frameworks & drivers
- **Dependencies**: Dependencies point inward toward business logic
- **Use Case**: Enterprise applications requiring high testability and maintainability
- **Template**: `web-api-clean`

**Key Validations:**
- Four distinct layers: entities, use cases, interface adapters, frameworks & drivers
- Dependency Inversion Principle (DIP) compliance
- Business logic isolation in entities and use cases
- Interfaces (ports) defined in inner layers
- Infrastructure implementations in outer layers
- Framework abstraction through adapters

### 3. Domain-Driven Design (`ddd_test.go`)
- **Pattern**: Domain-centric design with rich domain models
- **Structure**: Domain, application, presentation, infrastructure layers + shared kernel
- **Dependencies**: Focus on domain modeling and business logic
- **Use Case**: Complex business domains with rich domain logic
- **Template**: `web-api-ddd`

**Key Validations:**
- Domain-centric organization by business concepts
- Rich domain models with business logic
- CQRS pattern implementation (commands/queries)
- Domain events and event handling
- Value objects and specifications
- Shared kernel for common domain concepts
- Application services orchestrating domain operations

### 4. Hexagonal Architecture (`hexagonal_test.go`)
- **Pattern**: Ports and Adapters (Hexagonal) architecture
- **Structure**: Core (domain + ports + services), primary adapters, secondary adapters
- **Dependencies**: Core defines interfaces, adapters implement them
- **Use Case**: Highly testable applications with multiple interfaces
- **Template**: `web-api-hexagonal` (may not be fully implemented yet)

**Key Validations:**
- Clear separation between core and adapters
- Port interfaces defined in core
- Primary adapters (driving side) - HTTP, CLI, etc.
- Secondary adapters (driven side) - database, external APIs
- Framework and technology independence of core
- Dependency inversion through port interfaces

## Test Structure

Each architecture test file follows a consistent pattern:

### Test Categories

1. **Basic Generation Tests**
   - Validate that projects generate without errors
   - Ensure core files are created
   - Verify project compiles successfully

2. **Structure Validation Tests**
   - Verify correct directory organization
   - Ensure architecture-specific patterns are followed
   - Validate that inappropriate patterns are NOT present

3. **Dependency Direction Tests**
   - Verify dependencies flow in the correct direction
   - Ensure architecture principles are not violated
   - Validate abstraction boundaries

4. **Feature Integration Tests**
   - Test database integration within the architecture
   - Validate authentication implementation
   - Ensure logging integration follows architecture patterns

5. **Framework/Logger Variation Tests**
   - Test different web frameworks (gin, echo, fiber, chi)
   - Test different logging libraries (slog, zap, logrus, zerolog)
   - Ensure architecture abstraction works with all variants

6. **Business Logic Isolation Tests**
   - Verify business logic is properly separated
   - Ensure external concerns don't leak into core
   - Validate testing strategies specific to each architecture

## Test Naming Convention

Tests follow the pattern: `TestArchitecture_[ArchitectureName]_[Scenario]`

Examples:
- `TestArchitecture_Standard_BasicGeneration`
- `TestArchitecture_Clean_DependencyInversionPrinciple`
- `TestArchitecture_DDD_DomainEvents`
- `TestArchitecture_Hexagonal_PortsAndAdapters`

## Architecture-Specific Validations

### Standard Architecture
- ✅ Simple, flat internal structure
- ✅ Direct framework usage in handlers
- ✅ Straightforward repository pattern
- ❌ No complex abstraction layers
- ❌ No dependency inversion containers

### Clean Architecture
- ✅ Four distinct layers with proper separation
- ✅ Dependency Inversion Principle compliance
- ✅ Business logic in entities and use cases
- ✅ Framework abstraction through adapters
- ✅ Dependency injection container
- ❌ No direct framework imports in domain layer

### Domain-Driven Design
- ✅ Domain-centric organization
- ✅ Rich domain models with business methods
- ✅ CQRS pattern (commands/queries)
- ✅ Domain events and specifications
- ✅ Value objects for complex data
- ✅ Shared kernel for common concepts
- ❌ No generic CRUD operations without business meaning

### Hexagonal Architecture
- ✅ Core business logic completely isolated
- ✅ Port interfaces defined in core
- ✅ Adapters implement port interfaces
- ✅ Technology independence
- ✅ Symmetric design (primary/secondary adapters)
- ❌ No external framework dependencies in core

## Running the Tests

```bash
# Run all architecture tests
go test ./tests/integration/architectures/...

# Run specific architecture tests
go test ./tests/integration/architectures/standard_test.go
go test ./tests/integration/architectures/clean_test.go
go test ./tests/integration/architectures/ddd_test.go
go test ./tests/integration/architectures/hexagonal_test.go

# Run with verbose output
go test -v ./tests/integration/architectures/...

# Run specific test case
go test -v ./tests/integration/architectures/ -run TestArchitecture_Clean_DependencyInversionPrinciple
```

## Common Failure Scenarios

### Template Not Implemented
Some architecture templates may not be fully implemented yet. Tests handle this gracefully by:
- Checking for directory/file existence before validation
- Using `t.Skip()` when core components are missing
- Logging warnings for missing expected components

### Compilation Failures
All tests verify that generated projects compile successfully. Common issues:
- Missing imports in generated code
- Incorrect template variable substitution
- Circular dependencies
- Missing interface implementations

### Architecture Violations
Tests specifically check for architecture principle violations:
- Wrong import directions (e.g., domain importing infrastructure)
- Missing abstraction layers
- Incorrect dependency flows
- Framework leakage into business logic

## Development Guidelines

When adding new architecture tests:

1. **Follow the naming convention** for consistency
2. **Test the positive path** (what should be present)
3. **Test the negative path** (what should NOT be present)
4. **Validate compilation** of generated projects
5. **Handle missing templates gracefully** with appropriate skips/warnings
6. **Test multiple variations** (frameworks, loggers, features)
7. **Focus on architecture-specific concerns** rather than generic functionality

## Integration with CI/CD

These tests are designed to run in continuous integration environments:
- No external dependencies required
- Generate projects in temporary directories
- Clean up after themselves
- Provide clear failure messages
- Support parallel execution

## Future Enhancements

Planned improvements to the architecture testing suite:

1. **Performance Testing**: Validate that generated projects meet performance benchmarks
2. **Security Testing**: Ensure generated code follows security best practices
3. **Migration Testing**: Validate that projects can be upgraded between architecture patterns
4. **Template Validation**: Pre-validate templates before generation testing
5. **Cross-Platform Testing**: Ensure generated projects work on different operating systems
6. **Integration Testing**: Test generated projects with real databases and external services

## Contributing

When contributing to architecture tests:

1. Understand the architecture pattern being tested
2. Review existing tests for patterns and conventions
3. Ensure tests are deterministic and reliable
4. Add appropriate documentation for complex validations
5. Test edge cases and error conditions
6. Verify tests work with incomplete template implementations

This comprehensive test suite ensures that the go-starter generator produces high-quality, architecturally sound Go projects that follow established patterns and best practices.