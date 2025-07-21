# Template Integration Tests

This directory contains integration tests for the template engine, validating template compilation, parsing, generation, and processing functionality.

## Test Files

### 📄 **compilation_test.go** - Template Compilation Tests
Tests template compilation from source files to executable templates:
- Template syntax validation
- Variable binding verification
- Template inheritance and includes
- Cross-platform template compilation
- Error handling for malformed templates

### ⚙️ **generation_test.go** - Template Generation Tests  
Tests template generation with various configurations:
- Variable substitution accuracy
- Conditional generation logic
- File generation from templates
- Output formatting and structure
- Multi-template generation workflows

## Test Methodology

### Template Compilation Testing
1. **Syntax Validation**: Verify template syntax is correct
2. **Variable Binding**: Test variable definitions and usage
3. **Inheritance**: Test template extends and includes
4. **Error Handling**: Test malformed template handling
5. **Performance**: Test compilation speed and memory usage

### Template Generation Testing
1. **Variable Substitution**: Test all variable types and formats
2. **Conditional Logic**: Test if/else and loop constructs
3. **File Output**: Verify generated files are correct
4. **Multi-Configuration**: Test with different variable sets
5. **Edge Cases**: Test boundary conditions and error scenarios

## Running Template Tests

```bash
# Run all template tests
go test ./tests/integration/templates/... -v

# Run specific test files
go test ./tests/integration/templates/compilation_test.go -v
go test ./tests/integration/templates/generation_test.go -v

# Run with template debugging
go test ./tests/integration/templates/... -v -args -debug-templates

# Run with performance profiling
go test ./tests/integration/templates/... -v -cpuprofile=template.prof
```

## Template Test Configuration

### Environment Variables
- `TEMPLATE_DEBUG=1` - Enable template debugging output
- `TEMPLATE_CACHE_DISABLED=1` - Disable template caching for testing
- `TEMPLATE_STRICT_MODE=1` - Enable strict template validation

### Test Data
Template tests use fixtures located in:
- `testdata/templates/` - Template source files
- `testdata/configs/` - Test configuration files
- `testdata/expected/` - Expected output files

## Adding New Template Tests

When adding new template tests:

1. **Follow Naming Convention**: `Test[Component]_[Feature]_[Scenario]`
2. **Use Test Fixtures**: Create reusable template test data
3. **Include Error Cases**: Test both success and failure scenarios
4. **Document Complex Logic**: Add comments for complex template logic
5. **Test Multiple Formats**: Test different output formats and structures