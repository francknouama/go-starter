# AI Advisor BDD Test Suite

This directory contains comprehensive Behavior-Driven Development (BDD) tests for the AI-powered architecture advisor. The tests validate the advisor from a user acceptance perspective using Gherkin scenarios written in natural language.

## Overview

The BDD test suite provides:

- **Human-readable scenarios** written in Gherkin format
- **Comprehensive coverage** of advisor functionality 
- **Step definitions** that map Gherkin steps to Go code
- **Test helpers** for common validation patterns
- **Integration testing** with CLI commands

## File Structure

```
tests/acceptance/advisor/
├── README.md                           # This documentation
├── features/                           # Gherkin feature files
│   ├── ai_advisor_quick_mode.feature           # Quick recommendation scenarios
│   ├── ai_advisor_blueprint_selection.feature # Blueprint selection logic
│   ├── ai_advisor_framework_recommendations.feature # Framework suggestions
│   └── ai_advisor_edge_cases.feature          # Error handling scenarios
├── ai_advisor_bdd_steps.go            # Step definitions (Given/When/Then)
├── ai_advisor_bdd_test.go             # BDD test implementations
├── bdd_test_helpers.go               # Test utilities and fixtures
└── ai_advisor_atdd_test.go           # Original ATDD tests (maintained)
```

## Feature Coverage

### 1. Quick Mode Recommendations (`ai_advisor_quick_mode.feature`)

Tests the advisor's ability to provide fast, intelligent recommendations based on minimal input:

- **E-commerce API for senior teams** - Complex web applications
- **CLI tools for junior teams** - Simple command-line utilities  
- **Fintech APIs for expert teams** - High-security financial applications
- **IoT Lambda for mixed teams** - Event-driven serverless functions
- **Performance requirements** - Response time validation (< 1 second for 100 recommendations)

### 2. Blueprint Selection (`ai_advisor_blueprint_selection.feature`)

Validates intelligent blueprint selection based on project requirements:

- **Simple vs complex projects** - Appropriate architecture patterns
- **Domain-specific recommendations** - E-commerce, fintech, healthcare, etc.
- **Team experience matching** - Junior → simple, Expert → advanced
- **Conflicting requirements** - Graceful handling of contradictions
- **Migration scenarios** - Monolith-to-microservice patterns

### 3. Framework Recommendations (`ai_advisor_framework_recommendations.feature`)

Tests context-aware framework selection:

- **Performance optimization** - High-load scenarios → Gin, Fiber
- **Beginner-friendly choices** - Junior teams → Cobra for CLI
- **Enterprise features** - Security, middleware, routing capabilities
- **Ecosystem considerations** - Community support and integrations
- **Compatibility validation** - Framework-blueprint compatibility

### 4. Edge Cases (`ai_advisor_edge_cases.feature`)

Ensures robust error handling and edge case management:

- **Invalid inputs** - Malformed project types, team experiences
- **Missing requirements** - Graceful degradation with defaults
- **Conflicting needs** - Junior team + enterprise requirements
- **Performance constraints** - Timeout handling, memory limits
- **Unicode handling** - International project names and characters

## Running the Tests

### Prerequisites

1. **Build the CLI binary** (for integration tests):
   ```bash
   make build
   # or
   go build -o bin/go-starter main.go
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

### Run All BDD Tests

```bash
# Run all BDD tests
go test ./tests/acceptance/advisor/ -v

# Run specific BDD test categories
go test ./tests/acceptance/advisor/ -v -run "TestAIAdvisor_BDD_QuickMode"
go test ./tests/acceptance/advisor/ -v -run "TestAIAdvisor_BDD_Blueprint"
go test ./tests/acceptance/advisor/ -v -run "TestAIAdvisor_BDD_Framework"
go test ./tests/acceptance/advisor/ -v -run "TestAIAdvisor_BDD_EdgeCases"
```

### Run Original ATDD Tests

```bash
# Run original ATDD tests (still maintained)
go test ./tests/acceptance/advisor/ -v -run "TestAIAdvisor_ATDD"
```

### Run CLI Integration Tests

```bash
# Requires built binary
go test ./tests/acceptance/advisor/ -v -run "CLI"
```

## Test Architecture

### BDD Test Context

The `BDDTestContext` maintains state during scenario execution:

```go
type BDDTestContext struct {
    advisor            *advisor.ArchitectureAdvisor
    interactiveAdvisor *advisor.InteractiveAdvisor
    requirements       advisor.ProjectRequirements
    projectType        string
    recommendation     *advisor.ArchitectureRecommendation
    // ... additional state
}
```

### Step Definitions

Step definitions map Gherkin steps to Go functions:

```go
// Given steps - Set up test conditions
func (ctx *BDDTestContext) IHaveAProjectRequirementFor(projectType string) error
func (ctx *BDDTestContext) MyTeamExperienceLevelIs(experience string) error
func (ctx *BDDTestContext) MyProjectDomainIs(domain string) error

// When steps - Execute actions
func (ctx *BDDTestContext) IRequestQuickModeRecommendations() error
func (ctx *BDDTestContext) IAskForBlueprintRecommendations() error

// Then steps - Validate outcomes
func (ctx *BDDTestContext) IShouldGetBlueprintRecommendationsIncluding(blueprints string) error
func (ctx *BDDTestContext) TheConfidenceLevelShouldBeAtLeast(minConfidence string) error
```

### Test Helpers and Fixtures

The `BDDTestHelpers` provides:

- **Standardized test data** - Common project requirements
- **Validation functions** - Recommendation quality checks
- **Fluent assertions** - Chainable validation interface
- **Common scenarios** - Pre-built test cases

Example fluent assertions:
```go
err := NewBDDAssertions(t, recommendation).
    ShouldHaveBlueprintIn([]string{"web-api", "web-api-clean"}).
    ShouldHaveConfidenceAtLeast(0.6).
    ShouldHaveReasoning().
    ShouldHaveAlternatives().
    ShouldHaveFileCountBetween(20, 100).
    Assert()
```

## Example Scenarios

### Gherkin Scenario Example

```gherkin
Scenario: E-commerce API recommendation for senior team
  Given I have a project requirement for "api"
  And my project domain is "e-commerce"
  And my team experience level is "senior"
  When I request quick mode recommendations
  Then I should get blueprint recommendations including "web-api-ddd", "web-api-clean", "web-api"
  And I should get framework recommendations including "gin", "echo", "fiber"
  And the confidence level should be at least 0.6
  And I should receive reasoning for the recommendations
  And I should receive alternatives
  And the estimated file count should be between 20 and 100
```

### Go Test Implementation

```go
RunBDDScenario(t, "E-commerce API recommendation for senior team", func(ctx *BDDTestContext) error {
    // Given
    if err := ctx.TheAIAdvisorIsAvailable(); err != nil {
        return err
    }
    if err := ctx.IHaveAProjectRequirementFor("api"); err != nil {
        return err
    }
    if err := ctx.MyProjectDomainIs("e-commerce"); err != nil {
        return err
    }
    if err := ctx.MyTeamExperienceLevelIs("senior"); err != nil {
        return err
    }
    
    // When
    if err := ctx.IRequestQuickModeRecommendations(); err != nil {
        return err
    }
    
    // Then
    if err := ctx.IShouldGetBlueprintRecommendationsIncluding("web-api-ddd, web-api-clean, web-api"); err != nil {
        return err
    }
    // ... additional validations
    
    return ctx.AssertRecommendationQuality()
})
```

## Writing New BDD Tests

### 1. Add Gherkin Scenarios

Create new scenarios in the appropriate `.feature` file:

```gherkin
Scenario: Your new scenario name
  Given some precondition
  When some action occurs  
  Then some outcome should happen
```

### 2. Implement Step Definitions

Add step definition functions in `ai_advisor_bdd_steps.go`:

```go
func (ctx *BDDTestContext) YourNewGivenStep(parameter string) error {
    // Set up test state
    return nil
}

func (ctx *BDDTestContext) YourNewWhenStep() error {
    // Execute the action
    return nil
}

func (ctx *BDDTestContext) YourNewThenStep(expected string) error {
    // Validate the outcome
    return nil
}
```

### 3. Create Test Implementation

Add the test in `ai_advisor_bdd_test.go`:

```go
RunBDDScenario(t, "Your scenario name", func(ctx *BDDTestContext) error {
    // Map to step definitions
    if err := ctx.YourNewGivenStep("parameter"); err != nil {
        return err
    }
    if err := ctx.YourNewWhenStep(); err != nil {
        return err
    }
    return ctx.YourNewThenStep("expected")
})
```

## Best Practices

### Scenario Writing

- **Use descriptive names** that explain business value
- **Follow Given/When/Then structure** consistently
- **Keep scenarios focused** on single behaviors
- **Use scenario outlines** for data-driven tests
- **Write from user perspective** not implementation details

### Step Definitions

- **Keep steps reusable** across scenarios
- **Handle errors gracefully** with clear messages
- **Validate inputs** before processing
- **Maintain test state** in BDDTestContext
- **Use helper functions** for complex validations

### Test Data

- **Use fixtures** for common test data
- **Isolate tests** - no shared state between scenarios
- **Validate edge cases** - empty, null, invalid inputs
- **Test realistic scenarios** based on actual use cases
- **Include performance tests** for critical paths

## Integration with CI/CD

The BDD tests integrate with the existing CI/CD pipeline:

```bash
# In GitHub Actions or similar
- name: Run BDD Tests
  run: |
    make build
    go test ./tests/acceptance/advisor/ -v -timeout=5m
```

## Troubleshooting

### Common Issues

1. **CLI binary not found**
   ```bash
   make build  # Build the binary first
   ```

2. **Test timeouts**
   ```bash
   go test ./tests/acceptance/advisor/ -v -timeout=10m
   ```

3. **Import errors**
   ```bash
   go mod tidy  # Update dependencies
   ```

### Debugging Tests

1. **Enable verbose output**:
   ```bash
   go test ./tests/acceptance/advisor/ -v -run "YourSpecificTest"
   ```

2. **Add debug logging** in step definitions:
   ```go
   fmt.Printf("DEBUG: recommendation = %+v\n", ctx.recommendation)
   ```

3. **Use table-driven tests** for systematic testing:
   ```go
   tests := []struct {
       name     string
       input    string
       expected string
   }{
       // test cases
   }
   ```

## Benefits of BDD Approach

### For Stakeholders

- **Human-readable specifications** that document system behavior
- **Living documentation** that stays up-to-date with implementation
- **Clear acceptance criteria** for feature development
- **Shared understanding** between developers, QA, and product teams

### For Developers

- **Test-driven development** with clear requirements
- **Regression prevention** through comprehensive scenarios
- **Refactoring confidence** with behavior preservation
- **Documentation** that explains why code exists

### For QA Engineers

- **Acceptance criteria validation** in automated tests
- **Edge case coverage** through systematic scenario design
- **Performance validation** with built-in timing tests
- **Integration testing** across CLI and library interfaces

The BDD test suite provides a robust foundation for validating the AI advisor's behavior while maintaining readability and maintainability for all team members.