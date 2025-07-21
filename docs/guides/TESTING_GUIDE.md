# Testing Guide

This guide outlines the testing strategy and best practices for the `go-starter` project. We are committed to a Test-Driven Development (TDD) approach to ensure high code quality, reliability, and maintainability.

## 🧪 Test-Driven Development (TDD)

All new features and bug fixes in `go-starter` are developed following the Red-Green-Refactor TDD cycle:

1.  **Red:** Write a failing test that defines a new piece of functionality or a bug fix.
2.  **Green:** Write the minimum amount of code necessary to make the test pass.
3.  **Refactor:** Improve the code's design, readability, and performance while ensuring all tests still pass.

## 📊 Test Types

We categorize our tests into the following types:

### 1. Unit Tests

Unit tests focus on individual components or functions in isolation. They are fast, reliable, and provide immediate feedback on code changes.

*   **Location:** Typically reside in the same package as the code they test, in files ending with `_test.go` (e.g., `my_package/my_file_test.go`).
*   **Naming:** Test functions are prefixed with `Test` (e.g., `TestMyFunction`).
*   **Best Practices:**
    *   Test one thing at a time.
    *   Avoid external dependencies (use mocks or fakes).
    *   Ensure high code coverage for critical logic.

### 2. Integration Tests

Integration tests verify the interactions between different components or modules. They ensure that various parts of the system work correctly together.

*   **Location:** Reside in the `tests/integration/` directory.
*   **Naming:** Test functions are prefixed with `TestIntegration_` (e.g., `TestIntegration_ProjectGeneration`).
*   **Best Practices:**
    *   Test real interactions between components.
    *   May involve external dependencies (e.g., temporary files, mocked network calls).
    *   Use `tests/helpers/test_utils.go` for common setup/teardown tasks.

### 3. End-to-End (E2E) Tests

E2E tests simulate real user scenarios, testing the entire application flow from start to finish. These are typically run in CI/CD pipelines.

*   **Location:** Reside in the `tests/e2e/` directory (to be implemented).
*   **Naming:** Test functions are prefixed with `TestE2E_`.
*   **Best Practices:**
    *   Cover critical user journeys.
    *   Focus on high-level functionality.
    *   Should be stable and reliable.

### 4. Performance Benchmarks

Benchmarks measure the performance characteristics of specific code paths, helping to identify bottlenecks and track performance regressions.

*   **Location:** Reside in the `tests/benchmarks/` directory.
*   **Naming:** Benchmark functions are prefixed with `Benchmark` (e.g., `BenchmarkProjectGeneration`).
*   **Running Benchmarks:** Use `go test -bench=. -benchmem ./...`.

## 🚀 Running Tests

### Running All Tests

To run all unit and integration tests:

```bash
go test -v ./...
```

### Running Specific Tests

To run tests for a specific package:

```bash
go test -v ./internal/generator
```

To run a specific test function:

```bash
go test -v ./tests/integration/... -run TestIntegration_BasicProjectGeneration
```

### Running Benchmarks

To run all benchmarks:

```bash
go test -bench=. -benchmem ./...
```

To run a specific benchmark:

```bash
go test -bench=BenchmarkGenerator_GenerateWebAPI ./tests/benchmarks
```

## 📈 Code Coverage

We enforce a minimum code coverage threshold to ensure adequate testing of our codebase. The `scripts/check_coverage.sh` script is used to verify this.

### Checking Coverage Locally

To run tests and check coverage against the configured threshold:

```bash
./scripts/check_coverage.sh
```

### Viewing Detailed Coverage Report

After running tests with coverage, you can generate an HTML report to visualize covered and uncovered lines:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🛠️ Test Helpers and Mocks

### Test Helpers (`tests/helpers/`)

The `tests/helpers/` directory contains utility functions to simplify writing tests, such as:

*   `CreateTempDir(t *testing.T)`: Creates a temporary directory for tests and ensures cleanup.
*   `GenerateProject(t *testing.T, config types.ProjectConfig)`: Generates a project into a temporary directory for integration tests.
*   `AssertProjectGenerated(...)`: Asserts that expected files exist in a generated project.
*   `AssertCompilable(...)`: Asserts that a generated Go project compiles successfully.

### Mocks (`tests/helpers/mocks/`)

We use `github.com/stretchr/testify/mock` for creating mock implementations of interfaces, allowing us to isolate units under test from their dependencies.

*   **`MockPrompter`**: Mocks the interactive CLI prompts.
*   **`MockBlueprintRegistry`**: Mocks the blueprint registry for testing blueprint loading and selection logic.
*   **`MockFileSystem`**: Mocks file system operations for testing code that interacts with the file system without actual disk I/O.

## 🚀 CI/CD Integration

Our `.github/workflows/ci.yml` workflow automatically runs tests and checks coverage on every push and pull request. This ensures that all changes adhere to our quality standards before being merged.

*   The `test` job runs `scripts/check_coverage.sh` to enforce code coverage.
*   The `benchmark` job runs performance benchmarks to track performance metrics.

## ✅ TDD Development Commitment

As developers contributing to `go-starter`, we commit to:

*   Writing comprehensive tests before implementing each feature or bug fix.
*   Following the Red-Green-Refactor cycle diligently.
*   Ensuring high test coverage for all new and modified code.
*   Maintaining and improving the testing infrastructure.
*   Writing clear, concise, and maintainable tests.
