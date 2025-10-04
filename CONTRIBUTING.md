# Contributing to go-starter

Thank you for your interest in contributing to go-starter! This project follows **strict Test-Driven Development (TDD) principles** to ensure high code quality and maintainability.

## 🧪 TDD is Mandatory

**All code contributions must follow Test-Driven Development (TDD) practices.** This is not optional.

### What is TDD?

Test-Driven Development is a development process that follows the **Red-Green-Refactor** cycle:

1. **Red**: Write a failing test first
2. **Green**: Write the minimal code to make the test pass
3. **Refactor**: Improve the code while keeping tests green

### Why TDD?

- **Quality Assurance**: Ensures every feature works as intended
- **Design Improvement**: Writing tests first leads to better API design
- **Regression Prevention**: Comprehensive tests prevent future bugs
- **Documentation**: Tests serve as living documentation
- **Confidence**: Enables safe refactoring and feature additions

## 📋 Before You Start

### 1. Read the Development Requirements

- **Minimum Test Coverage**: >70% for all new code
- **Project Coverage**: Must maintain >30% overall coverage
- **Test Types**: Unit tests, integration tests, and edge case testing required
- **TDD Process**: Must follow Red-Green-Refactor cycle with commit evidence

### 2. Set Up Your Development Environment

```bash
# Clone the repository
git clone https://github.com/francknouama/go-starter.git
cd go-starter

# Install dependencies
go mod download

# Verify your setup
make test

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 3. Progressive Disclosure Development Requirements ✨

**NEW**: Working with progressive disclosure features requires additional setup and testing:

#### Progressive Disclosure Test Requirements
- **Complexity Level Testing**: All complexity levels (Simple, Standard, Advanced, Expert) must be tested
- **Help System Testing**: Both basic and advanced help modes must be validated
- **Blueprint Selection**: Complexity-aware blueprint mapping must be tested
- **Interactive Prevention**: Sufficient flags must prevent prompting
- **Default Application**: Smart defaults must be tested for CLI blueprints

#### Testing Progressive Disclosure Features
```bash
# Test progressive disclosure unit tests
go test ./internal/prompts/progressive_test.go -v

# Test CLI complexity levels
go test ./tests/acceptance/cli/progressive_disclosure_test.go -v

# Test blueprint selection with complexity
go test ./tests/acceptance/blueprints/cli/cli_simple_atdd_test.go -v

# Validate all complexity levels generate working code
for complexity in simple standard advanced expert; do
  go-starter new test-$complexity --type=cli --complexity=$complexity --dry-run
done
```

### 4. Understanding the Codebase

- Read the [README.md](README.md) for project overview
- Review existing tests to understand testing patterns
- Check [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for technical details
- Examine the current test coverage: `make test`

## 🚀 Development Workflow

### Step 1: Create an Issue

Before starting development, create an issue using our TDD-enforced templates:

- **Feature Request**: Use `.github/ISSUE_TEMPLATE/feature_request.yml`
- **Development Task**: Use `.github/ISSUE_TEMPLATE/development_task.yml`

Both templates require:
- Detailed test plan
- TDD commitment checkboxes
- Acceptance criteria
- Definition of Done

### Step 2: Follow the TDD Process

#### 🔴 Red Phase: Write Failing Tests

1. **Create a new branch**: `git checkout -b feature/your-feature-name`
2. **Write failing tests first**: Create test file(s) before implementation
3. **Run tests to confirm they fail**: `go test -v ./...`
4. **Commit the failing tests**: `git commit -m "Add failing tests for [feature]"`

```go
// Example: writing a failing test first
func TestNewFeature_ShouldReturnExpectedValue(t *testing.T) {
    // Arrange
    input := "test-input"
    expected := "expected-output"
    
    // Act
    result := NewFeature(input) // This function doesn't exist yet!
    
    // Assert
    if result != expected {
        t.Errorf("NewFeature() = %v, want %v", result, expected)
    }
}
```

#### 🟢 Green Phase: Make Tests Pass

1. **Implement minimal code**: Write only enough code to make tests pass
2. **Run tests**: `go test -v ./...`
3. **Ensure tests pass**: All new tests should now be green
4. **Commit implementation**: `git commit -m "Implement [feature] to pass tests"`

```go
// Example: minimal implementation
func NewFeature(input string) string {
    return "expected-output" // Minimal implementation
}
```

#### 🔄 Refactor Phase: Improve Code Quality

1. **Improve implementation**: Enhance code quality, performance, error handling
2. **Keep tests green**: Ensure tests continue to pass during refactoring
3. **Add more tests**: Add edge cases, error scenarios, integration tests
4. **Commit improvements**: `git commit -m "Refactor [feature] for better quality"`

```go
// Example: improved implementation
func NewFeature(input string) string {
    if input == "" {
        return "" // Handle edge case
    }
    // More sophisticated logic
    return processInput(input)
}

// Additional test for edge case
func TestNewFeature_WithEmptyInput_ShouldReturnEmpty(t *testing.T) {
    result := NewFeature("")
    if result != "" {
        t.Errorf("NewFeature(\"\") = %v, want empty string", result)
    }
}
```

### Step 3: Ensure Comprehensive Testing

#### Test Coverage Requirements

- **Line Coverage**: >70% for all new code
- **Branch Coverage**: >80% for critical paths
- **Error Coverage**: 100% for error handling paths

```bash
# Check coverage locally
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out

# Focus on package coverage
go test -coverprofile=coverage.out ./internal/your-package/
go tool cover -func=coverage.out
```

#### Test Types Required

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **Edge Case Tests**: Test boundary conditions
4. **Error Tests**: Test error handling and failure modes
5. **Table-Driven Tests**: Test multiple scenarios efficiently

#### Testing Best Practices

```go
// ✅ Good: Table-driven test
func TestValidateProjectName(t *testing.T) {
    tests := []struct {
        name        string
        projectName string
        wantErr     bool
        errorMsg    string
    }{
        {
            name:        "valid simple name",
            projectName: "my-project",
            wantErr:     false,
        },
        {
            name:        "empty name should error",
            projectName: "",
            wantErr:     true,
            errorMsg:    "project name cannot be empty",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateProjectName(tt.projectName)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateProjectName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if tt.wantErr && !strings.Contains(err.Error(), tt.errorMsg) {
                t.Errorf("ValidateProjectName() error = %v, want error containing %v", err, tt.errorMsg)
            }
        })
    }
}

// ✅ Good: Testing error scenarios
func TestGenerator_Generate_InvalidConfig(t *testing.T) {
    generator := New()
    invalidConfig := types.ProjectConfig{} // Missing required fields
    
    result, err := generator.Generate(invalidConfig, types.GenerationOptions{})
    
    // Test that error is returned
    if err == nil {
        t.Error("Expected error for invalid config, got nil")
    }
    
    // Test that result indicates failure
    if result.Success {
        t.Error("Expected result.Success to be false for invalid config")
    }
}
```

### Step 4: Code Quality Standards

Before submitting your PR, ensure:

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run vet
go vet ./...

# Run tests with race detection
go test -race ./...

# Check test coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### Step 5: Submit Pull Request

1. **Push your branch**: `git push origin feature/your-feature-name`
2. **Create PR**: Use the provided PR template
3. **Fill TDD sections**: Provide evidence of TDD compliance
4. **Include coverage report**: Show test coverage statistics
5. **Link to issue**: Reference the original issue

#### Required PR Information

Your PR must include:

- **TDD Evidence**: Commit hashes showing Red-Green-Refactor progression
- **Coverage Report**: Current coverage statistics
- **Test Files**: List of test files added/modified
- **TDD Compliance Declaration**: Confirmation you followed TDD

## 🔍 Code Review Process

### What Reviewers Look For

1. **TDD Compliance**: Evidence that TDD was followed
2. **Test Quality**: Comprehensive, well-structured tests
3. **Coverage**: Adequate test coverage for new code
4. **Code Quality**: Clean, readable, maintainable code
5. **Documentation**: Clear comments and updated docs

### Common Review Feedback

- **Insufficient Tests**: "Please add tests for error case X"
- **Low Coverage**: "Coverage is X%, need >70% for new code"
- **Missing Edge Cases**: "Please test boundary condition Y"
- **No TDD Evidence**: "Please provide commit showing tests were written first"

### Addressing Review Feedback

When addressing feedback:

1. **Continue TDD**: Write tests for missing scenarios first
2. **Update coverage**: Ensure new tests improve coverage
3. **Commit incrementally**: Show your TDD process in commits
4. **Re-run checks**: Verify all automated checks pass

## 🚨 Common TDD Violations

### ❌ Writing Implementation First

```go
// DON'T DO THIS: Implementation without tests
func NewFeature(input string) string {
    return "result"
}
// Then writing tests after...
```

### ❌ Insufficient Test Coverage

```go
// DON'T DO THIS: Only testing happy path
func TestNewFeature_OnlyHappyPath(t *testing.T) {
    result := NewFeature("valid-input")
    if result != "expected" {
        t.Error("Failed")
    }
    // Missing: error cases, edge cases, invalid inputs
}
```

### ❌ Non-Descriptive Tests

```go
// DON'T DO THIS: Unclear test names
func TestFeature(t *testing.T) {
    // What does this test actually verify?
}

// DO THIS: Clear, descriptive test names
func TestNewFeature_WithInvalidInput_ReturnsError(t *testing.T) {
    // Clear what's being tested
}
```

## 🎯 Project-Specific Guidelines

### Progressive Disclosure System Testing ✨

**NEW**: All progressive disclosure features require comprehensive TDD coverage.

When working on progressive disclosure features:

```go
func TestProgressiveDisclosure_ComplexitySelection(t *testing.T) {
    tests := []struct {
        name         string
        complexity   string
        blueprintType string
        expected     string
        wantErr      bool
    }{
        {
            name:         "CLI simple complexity",
            complexity:   "simple",
            blueprintType: "cli",
            expected:     "cli-simple",
            wantErr:      false,
        },
        {
            name:         "CLI standard complexity",
            complexity:   "standard", 
            blueprintType: "cli",
            expected:     "cli",
            wantErr:      false,
        },
        // More test cases for all blueprint types...
    }
    // Implementation tests...
}

func TestHelpSystem_FlagFiltering(t *testing.T) {
    // Test basic mode shows only essential flags
    // Test advanced mode shows all flags
    // Test flag deduplication
    // Test context-aware help hints
}

func TestInteractivePrompts_Prevention(t *testing.T) {
    // Test sufficient flags prevent prompts
    // Test smart defaults application
    // Test module path generation for testing
}
```

### Blueprint Testing

When working on blueprints:

```go
func TestBlueprint_Generation(t *testing.T) {
    // Test blueprint parsing
    // Test variable substitution
    // Test file generation
    // Test generated project compiles
    
    // NEW: Test complexity-aware generation
    // Test progressive disclosure integration
    // Test logger simplification
}
```

### Generator Testing

When working on the generator:

```go
func TestGenerator_ProcessTemplate(t *testing.T) {
    // Test template processing logic
    // Test error handling
    // Test file creation
    // Test directory structure
}
```

### CLI Testing

When working on CLI commands:

```go
func TestCommand_Execute(t *testing.T) {
    // Test command parsing
    // Test flag handling
    // Test interactive prompts (mock user input)
    // Test output formatting
}
```

## 📁 Project Organization Guidelines

### Documentation Structure

The project follows a structured documentation organization:

```
docs/
├── 01-getting-started/    # User onboarding and installation
├── 02-user-guides/        # User documentation and tutorials
├── 03-blueprints/         # Blueprint catalog and selection guides
├── 04-developers/         # Developer documentation and guides
├── 05-development/        # Development processes and workflows
├── 06-architecture/       # Architecture design and decisions
├── 07-reports/            # Project reports and analysis
│   └── archive/          # Historical reports (timestamped)
├── 08-agents/            # Agent documentation and specifications
├── 09-workspace-migration/ # Workspace migration planning
├── design/               # Design systems and brand identity
└── releases/             # Release notes and setup guides
```

### File Naming Conventions

1. **Directories**: Use `kebab-case` for directory names
   - ✅ `my-feature/`, `web-ui/`, `user-guides/`
   - ❌ `MyFeature/`, `Web_UI/`, `user_guides/`

2. **Major Documentation**: Use `SCREAMING_SNAKE_CASE` for important docs
   - ✅ `CONTRIBUTING.md`, `BLUEPRINT_STATUS_GUIDE.md`, `README.md`
   - ❌ `contributing.md`, `blueprint-status-guide.md`

3. **Regular Documentation**: Use `kebab-case` for standard docs
   - ✅ `web-ui-user-guide.md`, `quick-start.md`, `installation.md`
   - ❌ `web_ui_user_guide.md`, `QuickStart.md`

4. **Archive Files**: Include timestamp suffix for historical files
   - ✅ `phase4-validation-2025-08.md`, `web-ui-fix-summary-2025-08.md`
   - ❌ `phase4-validation.md`, `web-ui-fix-summary.md`

### Development Artifact Guidelines

**Build Artifacts** (excluded by .gitignore):
- `/bin/` - Compiled binaries
- `/dist/` - Distribution builds
- `/build/` - Build outputs
- `web/bin/` - Web server binaries

**Coverage Reports** (excluded by .gitignore):
- `**/coverage-reports/` - Coverage report directories
- `coverage-*.json` - JSON coverage files
- `coverage-*.xml` - XML coverage files

**Web UI Artifacts** (excluded by .gitignore):
- `web/bin/` - Web binaries
- `web/dist/` - Production builds
- `web/build/` - Development builds
- `web/node_modules/` - Node dependencies
- `web/.next/` - Next.js cache
- `web/.cache/` - Build cache

**Spec System Artifacts** (excluded by .gitignore):
- `specs/*/research.md` - Generated research
- `specs/*/plan.md` - Generated plans
- `specs/*/tasks.md` - Generated tasks
- `.specify/memory/` - Specification memory
- `.specify/templates/` - Specification templates

**Development Tools** (excluded by .gitignore):
- `scripts/*.log` - Script logs
- `scripts/output/` - Script outputs
- `validation-cli/bin/` - Validation binaries

### Where to Put New Files

**User-facing documentation** → `docs/02-user-guides/`
**Developer guides** → `docs/04-developers/`
**Architecture decisions** → `docs/06-architecture/`
**Project reports** → `docs/07-reports/` (or `docs/07-reports/archive/` for historical)
**Blueprint documentation** → `docs/03-blueprints/`
**Development processes** → `docs/05-development/`

### Git Best Practices

**File Moves**:
- Always use `git mv` to preserve history
- Verify history preservation with `git log --follow <file>`

**Commits**:
- Use conventional commit format: `type(scope): description`
- Common types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`
- Include Claude Code footer when AI-assisted

**Example**:
```bash
# Moving a file
git mv old/path/file.md new/path/file.md

# Committing the move
git commit -m "docs: reorganize documentation structure

Move user guides to docs/02-user-guides/ for better organization.

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

## 🆘 Getting Help

If you need help with TDD or testing:

1. **Review existing tests**: Look at current test files for patterns
2. **Check documentation**: Read Go testing best practices
3. **Ask questions**: Create a discussion or comment on issues
4. **Pair programming**: Request a review call for complex features

### Useful Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Table-Driven Tests in Go](https://github.com/golang/go/wiki/TableDrivenTests)
- [Go Test Coverage](https://golang.org/doc/tutorial/add-a-test)
- [TDD by Example (Kent Beck)](https://www.amazon.com/Test-Driven-Development-Kent-Beck/dp/0321146530)

## 📝 Issue and PR Templates

We provide TDD-enforced templates:

- **Feature Request** (`.github/ISSUE_TEMPLATE/feature_request.yml`): For user-facing features
- **Development Task** (`.github/ISSUE_TEMPLATE/development_task.yml`): For internal development
- **Pull Request** (`.github/PULL_REQUEST_TEMPLATE.md`): TDD compliance verification

All templates include mandatory TDD requirements and commitments.

## 🔧 Automated Enforcement

Our CI/CD pipeline automatically enforces:

- **Test Coverage**: Fails if coverage drops below thresholds
- **Test Quality**: Verifies all Go files have corresponding tests
- **Code Quality**: Runs linting and formatting checks
- **TDD Compliance**: Comments on PRs with coverage reports

## 📊 Current Project Status

- **Overall Coverage**: 31.6%
- **Target Coverage**: 85%
- **TDD Compliance**: Mandatory for all new code
- **Coverage Trend**: Improving with each release

## 🏆 Recognition

Contributors who consistently follow TDD practices will be:

- Recognized in release notes
- Given priority for code review
- Considered for maintainer roles
- Featured in project documentation

## ⚠️ Enforcement Policy

**Code that doesn't follow TDD practices will be rejected.** This includes:

- PRs without corresponding tests
- Tests written after implementation (without TDD evidence)
- Insufficient test coverage (<70% for new code)
- Poor test quality (only happy path testing)

We enforce TDD strictly because it's fundamental to the project's quality and maintainability.

---

**Thank you for contributing to go-starter with high-quality, test-driven code!** 🧪✨

Your commitment to TDD helps us build a reliable, maintainable project that serves the Go community well.