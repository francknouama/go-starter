# Pull Request

## Description
<!-- Provide a brief description of the changes in this PR -->

## Type of Change
<!-- Mark the relevant option with an "x" -->
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Template addition/modification
- [ ] Refactoring (no functional changes)
- [ ] Performance improvement
- [ ] Other (please describe):

## Related Issues
<!-- Link to any related issues -->
Closes #(issue number)
Fixes #(issue number)
Related to #(issue number)

## Changes Made
<!-- List the specific changes made in this PR -->
- 
- 
- 

## Template Changes (if applicable)
<!-- If this PR modifies templates, please specify which ones and how -->
- Template(s) affected: 
- Changes made: 
- Backward compatibility: 

## 🧪 Test-Driven Development (TDD) Verification
<!-- This section is MANDATORY for all code changes -->

### TDD Implementation Evidence
<!-- Provide evidence that TDD was followed -->
- [ ] **Red Phase**: I wrote failing tests first (provide commit hash where tests were added before implementation)
- [ ] **Green Phase**: I implemented minimal code to make tests pass (provide commit hash)
- [ ] **Refactor Phase**: I improved code quality while keeping tests green (provide commit hash if applicable)

**Commit demonstrating TDD approach**: `[commit-hash]` <!-- Replace with actual commit hash -->

### Test Coverage Analysis
<!-- Provide detailed test coverage information -->
- [ ] **New Code Coverage**: >70% line coverage for all new code
- [ ] **Overall Project Coverage**: No reduction in overall coverage percentage
- [ ] **Branch Coverage**: >80% for critical code paths
- [ ] **Error Path Coverage**: 100% coverage for error handling

**Coverage Report**: 
```
Package                          Coverage
internal/[package]               XX.X%
internal/[package]_test          XX.X%
Overall Project                  XX.X%
```

### Test Quality Verification
<!-- Verify comprehensive testing approach -->
- [ ] **Unit Tests**: All new functions/methods have corresponding unit tests
- [ ] **Integration Tests**: Cross-component functionality tested
- [ ] **Edge Cases**: Boundary conditions and edge cases tested
- [ ] **Error Scenarios**: All error paths and failure modes tested
- [ ] **Table-Driven Tests**: Used for multiple input scenarios
- [ ] **Mock Dependencies**: External dependencies properly mocked

### Test Execution Results
<!-- Confirm all tests pass -->
- [ ] **Local Test Run**: All tests pass locally (`go test -v ./...`)
- [ ] **Race Detection**: Tests pass with race detection (`go test -race ./...`)
- [ ] **Existing Tests**: No existing tests broken by changes
- [ ] **CI Pipeline**: All CI checks pass (or link to passing build)

## Testing Details
<!-- Describe the specific testing performed -->
- [ ] All existing tests pass
- [ ] New tests added following TDD principles
- [ ] Manual testing completed for user-facing changes
- [ ] Template generation tested with various configurations (if applicable)
- [ ] Generated projects compile successfully (if applicable)

### Test Files Added/Modified
<!-- List specific test files -->
- `[file_test.go]` - Added tests for [functionality]
- `[file_test.go]` - Modified tests for [changes]

### Test Configuration (if applicable)
<!-- If testing templates, list the configurations tested -->
- [ ] web-api-standard
- [ ] cli-standard  
- [ ] library-standard
- [ ] lambda-standard
- [ ] Different logger types (slog, zap, logrus, zerolog)
- [ ] Different database options
- [ ] Different architecture patterns

## Quality Assurance Checklist
<!-- Mark completed items with an "x" -->
- [ ] **Code Style**: My code follows the project's style guidelines
- [ ] **Code Review**: I have performed a thorough self-review of my code
- [ ] **Documentation**: I have commented complex code and updated documentation
- [ ] **Linting**: `golangci-lint run` passes with no new issues
- [ ] **Formatting**: `go fmt` applied to all modified files
- [ ] **Vet**: `go vet` passes with no issues
- [ ] **Dependencies**: `go mod tidy` run and go.mod/go.sum updated if needed

## TDD Compliance Declaration
<!-- REQUIRED: Confirm TDD compliance -->
- [ ] **I followed Test-Driven Development principles** (Red-Green-Refactor)
- [ ] **I wrote tests before implementation code** (not after)
- [ ] **All new code has comprehensive test coverage** (>70%)
- [ ] **I tested both happy paths and error scenarios**
- [ ] **I understand this PR will be rejected if TDD compliance cannot be verified**

## Template Testing (if applicable)
<!-- If this PR affects templates, confirm template testing -->
- [ ] Generated projects compile without errors
- [ ] All logger implementations work correctly
- [ ] Database integrations function properly
- [ ] Generated tests pass
- [ ] Docker builds succeed (if applicable)
- [ ] Makefiles execute without errors

## Breaking Changes
<!-- If this introduces breaking changes, describe them here -->
None

## Additional Notes
<!-- Any additional information, context, or screenshots -->

## Screenshots (if applicable)
<!-- Add screenshots to help explain your changes -->

---

**For maintainers:**
- [ ] Label applied
- [ ] Milestone assigned
- [ ] Ready for review