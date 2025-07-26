# Contributing to Tests

## Welcome Contributors! 🎉

Thank you for your interest in contributing to go-starter's testing infrastructure! This guide will help you add new test scenarios, improve existing tests, and maintain our 100% blueprint validation coverage.

## 🎯 Before You Start

### Prerequisites
- Go 1.21 or later installed
- Git configured with your GitHub account
- Familiarity with Go testing (`go test`)
- Basic understanding of the go-starter blueprint system

### Understanding Our Testing Philosophy
- **Quality First**: We maintain 100% compilation success for all generated projects
- **Comprehensive Coverage**: Every blueprint combination must be tested
- **Performance Aware**: Tests should be efficient and fast
- **Cross-Platform**: All tests must work on Windows, macOS, and Linux

## 🧪 Types of Test Contributions

### 1. Adding New Blueprint Test Scenarios
**When**: New blueprint types or configurations are added
**Where**: `tests/acceptance/enhanced/`
**Difficulty**: Beginner to Intermediate

### 2. Improving Architecture Validation
**When**: New architecture patterns are introduced
**Where**: `tests/acceptance/enhanced/architecture/`
**Difficulty**: Intermediate to Advanced

### 3. Enhancing Performance Tests
**When**: Performance optimizations need validation
**Where**: `tests/acceptance/enhanced/performance/`
**Difficulty**: Intermediate

### 4. Cross-Platform Compatibility
**When**: Platform-specific issues are discovered
**Where**: `tests/acceptance/enhanced/platform/`
**Difficulty**: Intermediate

## 📝 Step-by-Step Contribution Process

### Step 1: Set Up Your Environment

```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/go-starter.git
cd go-starter

# Create a new branch for your contribution
git checkout -b test/add-new-scenario-name

# Install dependencies
go mod download

# Run existing tests to ensure everything works
go test -v ./tests/acceptance/enhanced/... -timeout=30m
```

### Step 2: Identify What to Test

#### Adding a New Blueprint Combination
```bash
# Check if your combination is already tested
grep -r "web-api.*hexagonal.*postgres.*jwt" tests/acceptance/enhanced/

# If not found, you can add it!
```

#### Finding Gaps in Coverage
```bash
# Generate coverage report
go test -v -coverprofile=coverage.out ./tests/acceptance/enhanced/...
go tool cover -html=coverage.out

# Look for untested combinations in the code
```

### Step 3: Write Your Test

#### Choose the Right Location

```
tests/acceptance/enhanced/
├── architecture/      # For architecture pattern tests
├── database/          # For database integration tests  
├── framework/         # For web framework tests
├── auth/              # For authentication tests
├── lambda/            # For serverless function tests
├── cli/               # For CLI application tests
├── enterprise/        # For enterprise pattern tests
└── matrix/            # For complex combination tests
```

#### Template for New Test

```go
// tests/acceptance/enhanced/YOUR_CATEGORY/your_test.go

package YOUR_CATEGORY

import (
    "testing"
    "github.com/francknouama/go-starter/tests/helpers"
)

func TestYourNewScenario(t *testing.T) {
    testCases := []helpers.TestConfig{
        {
            Type:         "web-api",
            Architecture: "hexagonal", 
            Framework:    "gin",
            Database:     "postgres",
            ORM:          "gorm",
            Auth:         "jwt",
            Logger:       "zap",
            Description:  "Hexagonal architecture with Gin, PostgreSQL, and JWT",
        },
        // Add more test cases as needed
    }
    
    for _, tc := range testCases {
        t.Run(tc.Description, func(t *testing.T) {
            // Generate the project
            projectPath := helpers.GenerateProject(t, tc)
            
            // Validate compilation
            helpers.AssertCompilationSuccess(t, projectPath)
            
            // Add your specific validations
            validateYourSpecificRequirement(t, projectPath, tc)
            
            // Clean up (handled automatically by helpers)
        })
    }
}

func validateYourSpecificRequirement(t *testing.T, projectPath string, config helpers.TestConfig) {
    // Your specific validation logic here
    // Examples:
    // - Check for specific files
    // - Validate configuration consistency
    // - Test runtime behavior
    // - Verify architecture patterns
}
```

#### Feature File (Optional but Recommended)

```gherkin
# tests/acceptance/enhanced/YOUR_CATEGORY/features/your-scenario.feature

Feature: Your New Scenario Description
  As a developer using go-starter
  I want to generate projects with specific configurations
  So that my projects follow the intended patterns

  @critical @your-category
  Scenario: Specific configuration works correctly
    Given I generate a project with the configuration:
      | type         | web-api    |
      | architecture | hexagonal  |
      | framework    | gin        |
      | database     | postgres   |
    When I compile the generated project
    Then it should compile successfully
    And it should follow hexagonal architecture patterns
    And it should include proper database integration
```

### Step 4: Add Validation Logic

#### Common Validation Patterns

```go
// File existence validation
func validateRequiredFiles(t *testing.T, projectPath string, expectedFiles []string) {
    for _, file := range expectedFiles {
        helpers.AssertFileExists(t, projectPath, file)
    }
}

// Architecture pattern validation
func validateHexagonalPattern(t *testing.T, projectPath string) {
    // Check for ports (interfaces)
    helpers.AssertFileExists(t, projectPath, "internal/domain/ports/user_repository.go")
    
    // Check for adapters
    helpers.AssertFileExists(t, projectPath, "internal/adapters/repository/user_postgres.go")
    
    // Validate dependency directions using AST
    helpers.AssertNoDependencyViolations(t, projectPath, "hexagonal")
}

// Configuration consistency validation
func validateConfigConsistency(t *testing.T, projectPath string, config helpers.TestConfig) {
    // Check that selected database driver is used
    if config.Database == "postgres" {
        helpers.AssertFileContains(t, projectPath, "go.mod", "github.com/lib/pq")
    }
    
    // Check that selected framework is used
    if config.Framework == "gin" {
        helpers.AssertFileContains(t, projectPath, "go.mod", "github.com/gin-gonic/gin")
    }
}
```

#### Using Test Helpers

```go
// Use existing helpers whenever possible
helpers.AssertCompilationSuccess(t, projectPath)
helpers.AssertNoUnusedImports(t, projectPath)
helpers.AssertFileCount(t, projectPath, 25) // Expected file count
helpers.AssertArchitecturePattern(t, projectPath, "clean")
helpers.AssertDatabaseIntegration(t, projectPath, "postgres", "gorm")
```

### Step 5: Test Your Contribution

```bash
# Run your specific test
go test -v -run="TestYourNewScenario" ./tests/acceptance/enhanced/YOUR_CATEGORY/...

# Run related tests to ensure no regressions
go test -v ./tests/acceptance/enhanced/YOUR_CATEGORY/...

# Run full suite to ensure compatibility (optional but recommended)
go test -v ./tests/acceptance/enhanced/... -timeout=30m
```

### Step 6: Add Documentation

#### Update Test Documentation
```markdown
# Add to appropriate section in docs/testing/ENHANCED_TESTING_GUIDE.md

#### Your New Test Category
- **Description**: What your test validates
- **Coverage**: Which blueprint combinations are tested
- **Validation**: What specific checks are performed
- **Usage**: How to run your tests
```

#### Add Comments to Your Code
```go
// TestYourNewScenario validates that projects generated with specific
// configurations compile successfully and follow the intended patterns.
// This test covers the gap identified in issue #XXX where certain
// combinations were not being validated.
func TestYourNewScenario(t *testing.T) {
    // Implementation with clear comments
}
```

## 🔍 Code Review Checklist

Before submitting your PR, ensure:

### ✅ Test Quality
- [ ] Tests follow existing naming conventions
- [ ] Test cases cover edge cases and common scenarios
- [ ] All assertions use helper functions when available
- [ ] Error messages are clear and actionable
- [ ] Tests clean up properly (no resource leaks)

### ✅ Performance
- [ ] Tests complete in reasonable time (<5 minutes per scenario)
- [ ] No unnecessary file I/O or compilation
- [ ] Proper use of parallel execution where applicable
- [ ] Memory usage is reasonable

### ✅ Compatibility
- [ ] Tests work on Windows, macOS, and Linux
- [ ] No hardcoded paths or OS-specific behavior
- [ ] Proper handling of different Go versions
- [ ] No external dependencies without justification

### ✅ Documentation
- [ ] Feature files are clear and readable
- [ ] Code comments explain complex logic
- [ ] Test purpose is documented
- [ ] Related issues are referenced

## 🚀 Advanced Contribution Patterns

### Adding New Test Categories

```bash
# Create new test category directory
mkdir -p tests/acceptance/enhanced/your-new-category/features

# Follow the established pattern
cp tests/acceptance/enhanced/architecture/architecture_test.go \
   tests/acceptance/enhanced/your-new-category/your_category_test.go

# Update package name and test logic
```

### Creating Custom Validation Helpers

```go
// tests/helpers/your_custom_helpers.go

package helpers

import (
    "testing"
    "path/filepath"
)

// AssertYourCustomPattern validates your specific requirements
func AssertYourCustomPattern(t *testing.T, projectPath string, expectedPattern string) {
    // Your custom validation logic
    // Remember to:
    // - Use t.Helper() to mark as helper function
    // - Provide clear error messages
    // - Handle edge cases
    
    t.Helper()
    
    // Implementation here
    if !validationPassed {
        t.Errorf("Custom pattern validation failed: expected %s, got %s", 
                 expectedPattern, actualPattern)
    }
}
```

### Performance Testing Contributions

```go
func BenchmarkYourScenario(b *testing.B) {
    config := helpers.TestConfig{
        Type: "web-api",
        Framework: "gin",
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        projectPath := helpers.GenerateProject(b, config)
        // Don't measure cleanup time
        b.StopTimer()
        os.RemoveAll(projectPath)
        b.StartTimer()
    }
}
```

## 🐛 Troubleshooting Common Issues

### Test Failures
```bash
# Run with verbose output
go test -v -run="YourTest" ./tests/acceptance/enhanced/... -args -debug

# Keep temporary files for inspection
go test -v -run="YourTest" ./tests/acceptance/enhanced/... -args -keep-temp

# Check logs
tail -f /tmp/go-starter-test-*.log
```

### Performance Issues
```bash
# Profile your test
go test -run="YourTest" -cpuprofile=cpu.prof ./tests/acceptance/enhanced/...
go tool pprof cpu.prof

# Check memory usage
go test -run="YourTest" -memprofile=mem.prof ./tests/acceptance/enhanced/...
go tool pprof mem.prof
```

### Platform Compatibility
```bash
# Test on different platforms using Docker
docker run --rm -v $(pwd):/workspace -w /workspace golang:1.21 \
    go test -v ./tests/acceptance/enhanced/YOUR_CATEGORY/...

# Use GitHub Actions for cross-platform testing
```

## 📞 Getting Help

### Community Resources
- **GitHub Discussions**: Ask questions and get help
- **Issue Tracker**: Report bugs or request features
- **Discord/Slack**: Real-time chat with maintainers (if available)

### Maintainer Contact
- Create an issue with the `question` label
- Reference this contribution guide in your issue
- Provide clear context about what you're trying to achieve

## 🎉 After Your Contribution

### What Happens Next
1. **Automated Testing**: Your PR will trigger CI/CD tests
2. **Code Review**: Maintainers will review your code
3. **Feedback**: You may receive suggestions for improvements
4. **Merge**: Once approved, your contribution will be merged
5. **Recognition**: You'll be added to the contributors list!

### Staying Involved
- **Watch for Issues**: Look for `good first issue` and `testing` labels
- **Improve Documentation**: Help other contributors
- **Share Knowledge**: Write blog posts or tutorials about testing go-starter
- **Mentor Others**: Help new contributors get started

## 🏆 Recognition

Contributors who help improve our testing infrastructure will be:
- Added to the CONTRIBUTORS.md file
- Mentioned in release notes
- Invited to join the testing team (for significant contributions)
- Given priority in feature request discussions

Thank you for helping make go-starter the most reliable Go project generator! 🚀