# Blueprint Unit Tests

This directory contains comprehensive unit tests for individual blueprints in the go-starter project generator. These tests validate blueprint structure, template syntax, generated code quality, and compliance with architectural principles.

## Test Coverage

### CLI-Simple Blueprint (`cli_simple_test.go`)

Comprehensive unit tests for the CLI-Simple blueprint that validate:

#### 1. **Template Metadata Tests**
- ✅ `template.yaml` structure and validity
- ✅ Required variables definition
- ✅ Essential files mapping
- ✅ Minimal dependencies (Cobra only)
- ✅ Simplicity constraints (< 10 files)

#### 2. **Template Syntax Tests**
- ✅ All `.tmpl` files have valid Go template syntax
- ✅ Template execution with test data
- ✅ Sprig function support (e.g., `upper`)
- ✅ Error-free template processing

#### 3. **Individual Template Tests**
- ✅ **main.go**: Entry point, slog usage, proper imports
- ✅ **cmd/root.go**: Cobra command, CLI flags, JSON output support
- ✅ **cmd/version.go**: Version command implementation
- ✅ **config.go**: Environment variable configuration, validation

#### 4. **Generated Project Tests**
- ✅ **Structure validation**: File count, directory structure
- ✅ **Compilation**: Generated code compiles successfully
- ✅ **Dependencies**: Only essential dependencies included
- ✅ **Simplicity**: No complex patterns or deep nesting

#### 5. **Feature-Specific Tests**
- ✅ **Variable substitution**: ProjectName, ModulePath, etc.
- ✅ **Slog integration**: Standard library logging usage
- ✅ **Complexity validation**: File count, interfaces, depth metrics

#### 6. **Integration Tests**
- ✅ **Template loader**: Blueprint loading and parsing
- ✅ **Filesystem structure**: Blueprint file organization

## Test Philosophy

### Unit vs Integration vs Acceptance
- **Unit Tests** (this directory): Test blueprint templates and structure
- **Integration Tests**: Test generation process and tool integration
- **Acceptance Tests** (ATDD): Test business requirements and user scenarios

### Complexity Validation
The CLI-Simple blueprint tests enforce simplicity constraints:
- ≤ 10 total files
- ≤ 6 Go files
- ≤ 2 directory levels deep
- ≤ 1 direct dependency
- ≤ 2 interface definitions (allowing `map[string]interface{}`)

### Template Validation Approach
1. **Syntax validation**: Parse templates with Go's `text/template`
2. **Execution validation**: Execute templates with test data
3. **Content validation**: Check generated content for expected patterns
4. **Anti-pattern validation**: Ensure complex patterns are absent

## Running Tests

### Individual Blueprint Tests
```bash
# Run CLI-Simple blueprint tests
go test ./tests/unit/blueprints/cli_simple_test.go -v

# Run all blueprint unit tests
go test ./tests/unit/blueprints/ -v
```

### With Coverage
```bash
go test ./tests/unit/blueprints/ -v -cover
```

### Specific Test Functions
```bash
# Test template metadata only
go test ./tests/unit/blueprints/ -v -run TestCLISimpleBlueprint_TemplateMetadata

# Test compilation only
go test ./tests/unit/blueprints/ -v -run TestCLISimpleBlueprint_GeneratedCodeCompilation
```

## Test Data and Fixtures

### Standard Test Data
```go
testData := map[string]interface{}{
    "ProjectName": "test-cli",
    "ModulePath":  "github.com/test/test-cli",
    "Author":      "Test Author",
    "GoVersion":   "1.21",
}
```

### Template Functions
Tests use Sprig template functions via `sprig.FuncMap()`:
- `upper`: Converts to uppercase (preserves hyphens)
- `lower`: Converts to lowercase
- Other Sprig functions as needed

## Adding New Blueprint Tests

To add tests for a new blueprint:

1. **Create test file**: `{blueprint_name}_test.go`
2. **Add init function**: Set up templates filesystem
3. **Implement test categories**:
   - Template metadata validation
   - Template syntax validation
   - Individual template tests
   - Generated project validation
   - Feature-specific tests
   - Integration tests

### Template for New Blueprint Tests
```go
package blueprints_test

import (
    // Standard imports
)

func init() {
    // Setup templates filesystem
}

func TestNewBlueprint_TemplateMetadata(t *testing.T) {
    // Test template.yaml structure
}

func TestNewBlueprint_TemplateFileSyntax(t *testing.T) {
    // Test all .tmpl files
}

// Add more specific tests...
```

## Integration with CI/CD

These unit tests run as part of the standard test suite:
- **Local development**: `go test ./...`
- **CI pipeline**: Automated on pull requests
- **Release validation**: Required before releases

## Debugging Test Failures

### Common Issues and Solutions

1. **Template syntax errors**:
   - Check for malformed Go template syntax
   - Verify variable names match template.yaml
   - Ensure proper Sprig function usage

2. **File path issues**:
   - Blueprint path resolution depends on working directory
   - Tests use `findBlueprintPath()` helper for cross-platform compatibility

3. **Generated project compilation failures**:
   - Check dependencies in go.mod template
   - Verify import paths are correct
   - Ensure all referenced packages exist

4. **Complexity validation failures**:
   - Review simplicity constraints for the specific blueprint
   - Consider if the constraint is appropriate or needs adjustment

### Debug Helpers
- Set `GOTESTS_DEBUG=1` for verbose output
- Use `t.Logf()` for debugging information
- Check temporary directories for generated projects

This comprehensive test suite ensures blueprint quality, consistency, and adherence to architectural principles while supporting rapid development and refactoring.