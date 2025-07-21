# CLI Command Integration Tests

This directory contains comprehensive integration tests for the go-starter CLI tool command interface, organized into focused test files that cover all aspects of CLI functionality.

**Note**: For blueprint-specific CLI tests (testing generated CLI projects), see `../blueprints/cli/` directory.

## Test Organization

### Test Files Structure

- **`help_test.go`** - Tests help command functionality and documentation
- **`version_test.go`** - Tests version command and version display
- **`list_test.go`** - Tests template listing functionality
- **`new_test.go`** - Tests project creation command (interactive and non-interactive modes)
- **`validation_test.go`** - Tests input validation and error handling
- **`flags_test.go`** - Tests global flags and command-line options
- **`completion_test.go`** - Tests shell completion functionality
- **`timeout_test.go`** - Tests that commands don't hang and complete within reasonable time
- **`environment_test.go`** - Tests environment variable handling and system integration
- **`utils_test.go`** - Shared utilities and helper functions for all tests

## Test Naming Convention

All tests follow the pattern: `TestCLI_[Category]_[Scenario]`

Examples:
- `TestCLI_Help_RootCommand` - Tests root help command
- `TestCLI_Version_Format` - Tests version output format
- `TestCLI_Validation_InvalidCommand` - Tests invalid command handling

## Test Categories

### 1. Help System Tests (`help_test.go`)
- Root command help (`--help`, `-h`)
- Subcommand help (`new --help`, `list --help`, etc.)
- Help command structure and formatting
- Help exit codes

### 2. Version Command Tests (`version_test.go`)
- Version command functionality
- Version flag variants (`--version`, `-v`)
- Version output format validation
- Version command consistency across invocation methods

### 3. List Command Tests (`list_test.go`)
- Basic list functionality
- Output format validation
- Verbose mode with list command
- Template information display
- Handling of empty template scenarios

### 4. New Command Tests (`new_test.go`)
- Interactive mode testing
- Non-interactive mode with flags
- Parameter validation
- Error handling for invalid inputs
- Help for new command

### 5. Input Validation Tests (`validation_test.go`)
- Invalid command handling
- Invalid flag handling
- Required argument validation
- Argument value validation
- Exit code consistency
- Graceful error handling

### 6. Flags and Options Tests (`flags_test.go`)
- Global flags (`--verbose`, `--config`)
- Invalid flags handling
- Flag order independence
- Config flag functionality
- Verbose flag behavior
- Help flag variants

### 7. Shell Completion Tests (`completion_test.go`)
- Supported shells (bash, zsh, fish, powershell)
- Unsupported shell handling
- Completion without arguments
- Completion help functionality
- Output format validation
- Command completion integration

### 8. Timeout and Reliability Tests (`timeout_test.go`)
- Non-interactive command timeouts
- Interactive command timeout handling
- Stdin input handling
- Resource-intensive operation timeouts
- Graceful shutdown testing

### 9. Environment Integration Tests (`environment_test.go`)
- `GO_STARTER_CONFIG` environment variable
- `GO_STARTER_VERBOSE` environment variable
- Multiple environment variables
- Environment variable precedence
- System environment variable interaction
- Default behavior in clean environment

## Key Testing Principles

### 1. Comprehensive Coverage
- All CLI commands and subcommands
- All flags and options
- Error conditions and edge cases
- Different invocation patterns

### 2. Reliability Testing
- Commands should not hang indefinitely
- Graceful error handling without panics
- Consistent exit codes
- Proper cleanup of resources

### 3. User Experience Testing
- Help messages are informative and accurate
- Error messages are helpful and actionable
- Output formatting is consistent
- Interactive modes handle input appropriately

### 4. Cross-Platform Compatibility
- Tests work on different operating systems
- Environment variable handling is consistent
- File path handling is platform-aware

### 5. Performance Considerations
- Commands complete within reasonable time
- Memory usage is appropriate
- Resource cleanup is proper

## Test Utilities

The `utils_test.go` file provides shared utilities:

- **`buildTestBinary(t)`** - Builds CLI binary for testing
- **`cleanupBinary(t, binary)`** - Cleans up test binary
- **`runCommand(t, binary, args...)`** - Executes commands consistently
- **`runCommandWithInput(t, binary, input, args...)`** - Executes with stdin
- **`expectSuccess(t, err, output, command)`** - Validates successful execution
- **`expectFailure(t, err, output, command)`** - Validates expected failures
- **`checkNotPanic(t, output, context)`** - Ensures no panic conditions
- **`checkContains(t, output, expected, context)`** - Validates output content
- **`checkNotEmpty(t, output, context)`** - Ensures non-empty output

## Running Tests

### Run All CLI Tests
```bash
go test -v ./tests/integration/cli/
```

### Run Specific Test Categories
```bash
# Test help functionality
go test -v ./tests/integration/cli/ -run TestCLI_Help

# Test version functionality
go test -v ./tests/integration/cli/ -run TestCLI_Version

# Test validation
go test -v ./tests/integration/cli/ -run TestCLI_Validation
```

### Run Individual Test Files
```bash
# Test only help functionality
go test -v ./tests/integration/cli/help_test.go ./tests/integration/cli/utils_test.go

# Test only completion functionality
go test -v ./tests/integration/cli/completion_test.go ./tests/integration/cli/utils_test.go
```

## Test Configuration

### Environment Variables
Tests may use these environment variables:
- `GO_STARTER_CONFIG` - Custom config file path
- `GO_STARTER_VERBOSE` - Enable verbose output

### Temporary Files
Tests create temporary files in system temp directories and clean them up automatically.

### Binary Building
Each test suite builds its own test binary to ensure isolation and consistency.

## Expected Behaviors

### Successful Commands
- Exit with code 0
- Produce meaningful output
- Complete within reasonable time (typically < 30 seconds)
- Handle flags and arguments correctly

### Failed Commands
- Exit with non-zero code
- Provide helpful error messages
- Not produce panic or stack traces
- Handle errors gracefully

### Interactive Commands
- Respond to stdin input appropriately
- Timeout gracefully when no input provided
- Handle partial input scenarios
- Provide clear prompts and guidance

## Maintenance

### Adding New Tests
1. Choose appropriate test file based on functionality
2. Follow naming convention: `TestCLI_[Category]_[Scenario]`
3. Use helper functions from `utils_test.go`
4. Include comprehensive error checking
5. Test both success and failure scenarios

### Updating Tests
1. Ensure backward compatibility
2. Update related test cases when CLI behavior changes
3. Maintain comprehensive coverage
4. Update this documentation when adding new test categories

### Debugging Test Failures
1. Check test output for specific error messages
2. Verify binary building succeeds
3. Run individual tests to isolate issues
4. Use verbose mode for detailed output
5. Check environment variable settings

This test suite ensures the CLI tool provides a reliable, user-friendly interface that handles all scenarios gracefully and provides clear feedback to users.