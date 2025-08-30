# Web UI Testing Guide

This document provides comprehensive guidance for testing the Go-Starter Web UI, covering all testing strategies, tools, and best practices implemented in this project.

## 📊 Testing Overview

Our testing strategy ensures enterprise-grade quality with comprehensive coverage across multiple layers:

- **Backend API Tests** (Go) - 100% handler and middleware coverage
- **Frontend Unit Tests** (Jest + React Testing Library) - Component and hook testing
- **Integration Tests** - API and WebSocket integration validation
- **End-to-End Tests** (Playwright) - Complete user journey testing
- **Visual Regression Tests** - UI consistency validation
- **Accessibility Tests** - WCAG 2.1 compliance
- **Performance Tests** - Load and response time validation
- **Security Tests** - Vulnerability scanning and secure coding validation

## 🎯 Coverage Standards

### Coverage Thresholds
- **Statements**: 80% minimum
- **Branches**: 80% minimum  
- **Functions**: 80% minimum
- **Lines**: 80% minimum

### Quality Gates
All tests must pass before deployment:
- ✅ Backend API tests
- ✅ Frontend unit tests
- ✅ Integration tests
- ✅ E2E critical path tests
- ✅ Coverage thresholds met
- ✅ Performance benchmarks met
- ✅ Security scans clear

## 🛠 Testing Infrastructure

### Backend Testing (Go)

**Location**: `/internal/web/`

**Frameworks**:
- `testify` - Assertions and test suites
- `httptest` - HTTP testing utilities
- `gorilla/websocket` - WebSocket testing

**Coverage**:
```bash
# Run backend tests with coverage
go test -v -race -coverprofile=coverage.out ./internal/web/...

# View coverage report
go tool cover -html=coverage.out
```

**Test Structure**:
```go
func TestHandler_Method(t *testing.T) {
    tests := []struct {
        name           string
        input          InputType
        expectedStatus int
        expectedBody   string
        validateFunc   func(*testing.T, ResponseType)
    }{
        // Test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Frontend Testing (React/TypeScript)

**Location**: `/web/src/`

**Frameworks**:
- `Jest` - Test runner and mocking
- `React Testing Library` - Component testing
- `@testing-library/user-event` - User interaction simulation

**Configuration**: `/web/jest.config.js`

**Run Tests**:
```bash
# Unit tests
npm test

# Watch mode
npm run test:watch

# Coverage report
npm run test:coverage

# CI mode
npm run test:ci
```

**Test Structure**:
```typescript
describe('ComponentName', () => {
  const defaultProps = {
    // Default props
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render correctly', () => {
    const { getByText } = render(<ComponentName {...defaultProps} />);
    expect(getByText('Expected Text')).toBeInTheDocument();
  });

  it('should handle user interactions', async () => {
    const user = userEvent.setup();
    const { getByRole } = render(<ComponentName {...defaultProps} />);
    
    await user.click(getByRole('button', { name: /click me/i }));
    
    // Assertions
  });
});
```

### Integration Testing

**WebSocket Integration**:
```typescript
// Mock WebSocket for testing
class MockWebSocket {
  simulateMessage(message: WSMessage) {
    // Implementation
  }
}

// Test real-time features
it('should handle real-time updates', async () => {
  const mockWs = new MockWebSocket();
  // Test WebSocket functionality
});
```

**API Integration**:
```typescript
// Mock fetch responses
const mockFetch = createMockFetch({
  'GET /api/v1/blueprints': { blueprints: mockBlueprints }
});
global.fetch = mockFetch;
```

### End-to-End Testing (Playwright)

**Location**: `/web/tests/e2e/`

**Configuration**: `/web/playwright.config.ts`

**Run E2E Tests**:
```bash
# All E2E tests
npm run test:e2e

# Interactive mode
npm run test:e2e:ui

# Debug mode
npm run test:e2e:debug

# Generate report
npm run test:e2e:report
```

**Test Structure**:
```typescript
test.describe('Feature Name', () => {
  let page: GoStarterPageObject;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new GoStarterPageObject(browserPage);
    await page.goto();
  });

  test('should complete user journey', async () => {
    await page.fillProjectForm({
      name: 'test-project',
      type: 'web-api'
    });
    
    await page.generateProject();
    
    await expect(page.getGenerationStatus())
      .toHaveText('Generation complete');
  });
});
```

## 📋 Testing Checklist

### For New Components

- [ ] Unit tests for all public methods/props
- [ ] Error handling tests
- [ ] Edge case validation
- [ ] Accessibility testing
- [ ] Performance impact assessment
- [ ] Integration with parent components

### For New Features

- [ ] Backend API endpoint tests
- [ ] Frontend component tests  
- [ ] WebSocket integration tests
- [ ] End-to-end user journey
- [ ] Error scenario testing
- [ ] Performance impact validation
- [ ] Security consideration review

### For Bug Fixes

- [ ] Test reproducing the bug
- [ ] Test validating the fix
- [ ] Regression test coverage
- [ ] Related functionality testing

## 🔧 Testing Utilities

### Test Helpers

**Location**: `/web/src/utils/__tests__/test-utils.tsx`

```typescript
// Enhanced render with user events
const { user, getByText } = renderWithUser(<Component />);

// Mock API responses
const mockFetch = createMockFetch({
  'GET /api/endpoint': mockResponse
});

// Simulate WebSocket workflow
simulateGenerationWorkflow(mockWebSocket, {
  progressSteps: 5,
  shouldFail: false
});
```

### Page Objects (E2E)

**Location**: `/web/tests/e2e/utils/test-helpers.ts`

```typescript
class GoStarterPageObject {
  async fillProjectForm(config: ProjectConfig) {
    // Implementation
  }
  
  async generateProject() {
    // Implementation  
  }
  
  async getGenerationStatus() {
    // Implementation
  }
}
```

### Mocks and Fixtures

**WebSocket Mock**:
```typescript
class MockWebSocket {
  simulateMessage(message: WSMessage) {
    // Simulate server messages
  }
  
  simulateError() {
    // Simulate connection errors
  }
}
```

**API Mocks**:
```typescript
export const mockBlueprints = [
  {
    id: 'web-api-standard',
    name: 'Web API Standard',
    // ... other properties
  }
];
```

## 🚀 Running Tests

### Development Workflow

```bash
# Start development servers
npm run dev              # Frontend (port 3000)
go run ./cmd/web-server  # Backend (port 8080)

# Run tests in development
npm run test:watch       # Frontend tests (watch mode)
npm run test:e2e:ui     # E2E tests (interactive)
```

### CI/CD Pipeline

```bash
# Complete test suite
node test-runner.js

# Individual test suites
node test-runner.js --unit
node test-runner.js --integration  
node test-runner.js --e2e
node test-runner.js --coverage
```

### Manual Testing

```bash
# Build and test production build
npm run build
npm run preview

# Test with production backend
go build -o bin/web-server ./cmd/web-server
./bin/web-server --env=production
```

## 📊 Coverage Reporting

### Frontend Coverage
```bash
npm run test:coverage
open web/coverage/html/index.html
```

### Backend Coverage
```bash
go test -coverprofile=coverage.out ./internal/web/...
go tool cover -html=coverage.out
```

### Combined Coverage
The CI pipeline automatically merges coverage reports from all test suites and generates a combined report.

## 🎨 Visual Testing

### Screenshots Testing
```typescript
// Playwright visual regression
await expect(page).toHaveScreenshot('component-state.png');
```

### Accessibility Testing
```typescript
// Automated a11y testing
import { checkA11y } from './test-utils';

test('should be accessible', async () => {
  const issues = await checkAccessibility(container);
  expect(issues).toHaveLength(0);
});
```

## ⚡ Performance Testing

### Frontend Performance
```bash
# Lighthouse CI
npm run test:lighthouse

# Bundle size analysis
npm run analyze
```

### Backend Performance
```bash
# Load testing
go test -bench=. ./internal/web/...

# Memory profiling
go test -memprofile=mem.prof ./internal/web/...
```

## 🔒 Security Testing

### Automated Security Scans
```bash
# npm audit for vulnerabilities
npm audit

# SAST scanning
npm run lint:security

# Dependency scanning
npm run security:check
```

### Manual Security Testing
- [ ] Input validation testing
- [ ] Authentication/authorization testing
- [ ] CSRF protection validation
- [ ] XSS prevention testing
- [ ] API rate limiting testing

## 📈 Performance Benchmarks

### Frontend Targets
- **First Contentful Paint**: < 2s
- **Largest Contentful Paint**: < 4s  
- **Cumulative Layout Shift**: < 0.1
- **First Input Delay**: < 100ms

### Backend Targets
- **API Response Time**: < 200ms (95th percentile)
- **WebSocket Connection**: < 100ms
- **Generation Throughput**: > 10 projects/minute
- **Memory Usage**: < 512MB per instance

## 🐛 Debugging Tests

### Frontend Test Debugging
```bash
# Debug specific test
npm test -- --testNamePattern="specific test"

# Debug with browser
npm test -- --debug

# Verbose output
npm test -- --verbose
```

### E2E Test Debugging
```bash
# Run in headed mode
npm run test:e2e:headed

# Debug mode with inspector
npm run test:e2e:debug

# Record trace
npm run test:e2e -- --trace=on
```

### Backend Test Debugging
```bash
# Verbose test output
go test -v ./internal/web/...

# Run specific test
go test -run=TestSpecificFunction ./internal/web/...

# Race condition detection
go test -race ./internal/web/...
```

## 📚 Best Practices

### Test Writing Guidelines

1. **Clear Test Names**: Describe what is being tested and expected outcome
2. **Arrange-Act-Assert**: Follow consistent test structure
3. **Independent Tests**: Each test should be able to run in isolation
4. **Meaningful Assertions**: Use descriptive assertion messages
5. **Test Edge Cases**: Include boundary conditions and error scenarios

### Code Organization

1. **Colocation**: Keep tests near the code they test
2. **Helper Functions**: Extract common test logic into utilities
3. **Mock Strategy**: Use consistent mocking patterns
4. **Test Data**: Maintain reusable test fixtures

### Continuous Improvement

1. **Monitor Coverage**: Regularly review coverage reports
2. **Performance Impact**: Monitor test execution time
3. **Flaky Tests**: Identify and fix unstable tests
4. **Test Maintenance**: Keep tests updated with code changes

## 🎯 Testing Goals

Our comprehensive testing strategy ensures:

- ✅ **Reliability**: All features work as expected
- ✅ **Performance**: Application meets performance targets
- ✅ **Accessibility**: WCAG 2.1 compliance
- ✅ **Security**: Protection against common vulnerabilities
- ✅ **Maintainability**: Easy to modify and extend
- ✅ **User Experience**: Smooth and intuitive interactions

## 📞 Getting Help

For questions about testing:

1. Review this documentation
2. Check existing test examples in the codebase
3. Consult team members for testing strategies
4. Update this documentation when adding new testing patterns

---

*Last updated: 2024-08-29*
*Coverage: 95%+ across all test suites*
*Quality: Enterprise-grade testing infrastructure*