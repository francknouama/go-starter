# Enhanced ATDD Tests

This directory contains enhanced ATDD tests that focus on template generation quality and configuration consistency, implementing Phase 1 of the enhanced testing strategy.

## Test Structure

### `/quality/`
- **`features/`**: Gherkin feature files defining quality validation scenarios
  - `code-quality.feature`: Compilation, unused imports/variables validation
  - `framework-consistency.feature`: Framework isolation and cross-contamination prevention
  - `configuration-consistency.feature`: Configuration file alignment validation
  - `static-analysis.feature`: Code quality and static analysis checks
  - `template-logic.feature`: Conditional generation validation
- **`quality_test.go`**: Godog test runner for quality validation features
- **`enhanced_steps_test.go`**: BDD step definitions implementing all Gherkin scenarios

### `/configuration/`
- **`features/`**: Gherkin feature files for configuration matrix testing
  - `matrix-testing.feature`: Priority-based configuration combination testing
- **`configuration_test.go`**: Godog test runner for configuration matrix features

## Test Categories

### Quality Validation Tests
- **Compilation Success**: Ensures all generated projects compile without errors
- **Unused Import Detection**: Identifies unused imports in generated code
- **Unused Variable Detection**: Catches unused variables via `go vet`
- **Configuration Consistency**: Validates go.mod, docker-compose.yml, and config files match selections
- **Framework Consistency**: Prevents cross-contamination between frameworks (gin/fiber/echo)

### Template Logic Tests
- **Conditional File Generation**: Validates files are generated/omitted based on configuration
- **Import Consistency**: Ensures imports match conditional generation logic
- **Dependency Consistency**: Validates go.mod dependencies match selected features

### Configuration Matrix Tests
- **Priority-Based Testing**: Critical/High/Medium/Low priority combinations
- **Comprehensive Coverage**: Tests all permutations of frameworks, databases, ORMs, loggers, auth types
- **Performance Optimization**: Uses priority levels to control test execution time

## Key Improvements Over Legacy Tests

1. **Quality Focus**: Catches template generation issues that functional tests miss
2. **Static Analysis**: Uses Go tooling (goimports, go vet) for validation
3. **Configuration Validation**: Ensures consistency across all generated files
4. **Matrix Testing**: Systematic coverage of configuration combinations
5. **Priority Management**: Focuses on critical combinations first

## Usage

### Run All Enhanced Tests
```bash
go test ./tests/acceptance/enhanced/...
```

### Run Specific Test Categories
```bash
# Quality validation tests (using godog)
go test ./tests/acceptance/enhanced/quality/

# Configuration matrix tests (using godog)
go test ./tests/acceptance/enhanced/configuration/
```

### Run with Specific Godog Options
```bash
# Run with verbose output
go test ./tests/acceptance/enhanced/quality/ -godog.format=pretty

# Run specific feature file
go test ./tests/acceptance/enhanced/quality/ -godog.paths=features/code-quality.feature

# Run with short mode (critical priority only)
go test -short ./tests/acceptance/enhanced/configuration/
```

## Integration with CI

These tests are designed to run alongside existing ATDD tests, providing additional quality validation without replacing functional coverage. They can be integrated into CI as additional jobs:

```yaml
jobs:
  enhanced-quality:
    runs-on: ubuntu-latest
    steps:
      - name: Run Enhanced Quality Tests
        run: go test ./tests/acceptance/enhanced/quality/
        
  enhanced-matrix:
    runs-on: ubuntu-latest
    steps:
      - name: Run Enhanced Matrix Tests (Critical)
        run: go test -short ./tests/acceptance/enhanced/configuration/
```

## Test Configuration Examples

### Critical Combinations (Always Tested)
- gin + postgresql + gorm + slog + jwt + standard
- echo + mysql + gorm + zap + jwt + standard  
- fiber + postgresql + raw + zerolog + none + standard

### Template Logic Validation
- ORM vs Raw SQL conditional generation
- JWT vs no-auth conditional generation
- Logger-specific file generation
- Framework-specific dependency inclusion

This enhanced testing approach ensures template quality while maintaining fast feedback cycles through priority-based execution.