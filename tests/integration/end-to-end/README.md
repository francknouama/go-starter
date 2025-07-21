# End-to-End Integration Tests

This directory contains comprehensive end-to-end integration tests that validate complete user workflows spanning multiple components of the go-starter system.

## Test Files

### 🎯 **basic_generation_test.go** - Basic Project Generation Workflow
Tests the fundamental project generation workflow:
- Blueprint selection and validation
- Configuration processing and validation
- Complete project generation process
- Generated project structure verification
- Compilation verification of generated projects

## Test Methodology

### End-to-End Workflow Testing
1. **User Input Simulation**: Simulate real user interactions and configurations
2. **Cross-Component Validation**: Test interactions between CLI, generator, templates, and blueprints
3. **Output Verification**: Validate complete generated project structure and functionality  
4. **Error Scenario Testing**: Test error handling across the entire workflow
5. **Performance Validation**: Test end-to-end performance characteristics

### Complete User Journey Testing
- **Blueprint Discovery**: Test blueprint listing and selection
- **Configuration Input**: Test configuration validation and processing
- **Generation Process**: Test complete generation workflow
- **Project Validation**: Test generated project compilation and basic functionality
- **Cleanup**: Test proper cleanup of temporary resources

## Test Scenarios

### Primary User Workflows
- **Simple CLI Generation**: End-to-end simple CLI project creation
- **Web API Generation**: Complete web API project with database integration
- **Multi-Blueprint Workflow**: Testing multiple blueprint types in sequence
- **Error Recovery**: Testing error handling and recovery scenarios

### Cross-Platform Testing
- **File System Compatibility**: Test on different file systems and path formats
- **Go Version Compatibility**: Test with different Go versions
- **Platform-Specific Generation**: Test platform-specific code generation

## Running End-to-End Tests

```bash
# Run all end-to-end tests
go test ./tests/integration/end-to-end/... -v

# Run with extended timeout for complex workflows
go test ./tests/integration/end-to-end/... -v -timeout=10m

# Run with detailed workflow logging
go test ./tests/integration/end-to-end/... -v -args -verbose-workflow

# Run with project preservation for debugging
go test ./tests/integration/end-to-end/... -v -args -keep-projects
```

## Performance Testing

End-to-end tests include performance validation:
- **Generation Time**: Measure complete generation workflow time
- **Memory Usage**: Monitor memory consumption during generation
- **File I/O Performance**: Test file creation and writing performance
- **Concurrent Generation**: Test multiple simultaneous generations

## Test Configuration

### Environment Variables
- `E2E_KEEP_PROJECTS=1` - Preserve generated projects for manual inspection
- `E2E_TIMEOUT=600s` - Set custom timeout for long-running workflows
- `E2E_VERBOSE=1` - Enable detailed workflow logging
- `E2E_PARALLEL_LIMIT=2` - Limit concurrent test execution

### Test Isolation
- Each test uses isolated temporary directories
- Tests clean up generated projects automatically (unless configured otherwise)
- No shared state between test executions
- Independent blueprint and configuration validation

## Adding New End-to-End Tests

When adding new end-to-end tests:

1. **Focus on Complete Workflows**: Test entire user journeys, not individual components
2. **Include Error Scenarios**: Test error handling and recovery workflows
3. **Use Realistic Configurations**: Use configurations similar to real-world usage
4. **Validate Complete Output**: Test generated project structure and functionality
5. **Consider Performance Impact**: Include performance assertions for critical workflows
6. **Document Test Scenarios**: Clearly document what user workflow is being tested