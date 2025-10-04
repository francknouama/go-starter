# End-to-End (E2E) Testing with Playwright

This document provides comprehensive guidance for running and maintaining E2E tests for the go-starter web interface.

## Overview

Our E2E testing suite uses [Playwright](https://playwright.dev/) to test the complete user journey of the go-starter web application. The tests cover:

- **Application Navigation**: Page layout, responsive design, disclosure modes
- **Project Configuration**: Form validation, progressive disclosure, real-time updates
- **Real-time Preview**: WebSocket connectivity, live preview updates, file explorer
- **Project Generation**: End-to-end project creation and download
- **AI Advisor Integration**: Architecture recommendations and form integration
- **Performance & Accessibility**: Core Web Vitals, WCAG compliance, keyboard navigation

## Quick Start

### Prerequisites

- Node.js 18+ installed
- Go 1.21+ installed
- Web server running on port 8080
- Frontend dev server running on port 3000

### Installation

```bash
cd web
npm install
npx playwright install
```

### Running Tests

```bash
# Run all E2E tests
npm run test:e2e

# Run tests with UI (interactive mode)
npm run test:e2e:ui

# Run tests in headed mode (see browser)
npm run test:e2e:headed

# Run specific test file
npx playwright test app-navigation.spec.ts

# Run tests for specific browser
npx playwright test --project=chromium

# Debug tests
npm run test:e2e:debug
```

## Test Structure

### Test Files

```
tests/e2e/
├── utils/
│   └── test-helpers.ts          # Page objects and utilities
├── app-navigation.spec.ts       # Navigation and layout tests
├── project-configuration.spec.ts # Form and configuration tests
├── real-time-preview.spec.ts    # WebSocket and preview tests
├── project-generation.spec.ts   # Generation and download tests
├── ai-advisor.spec.ts          # AI advisor integration tests
├── performance-accessibility.spec.ts # Performance and a11y tests
└── test-setup.ts               # Global test configuration
```

### Page Object Model

We use the Page Object Model pattern with a comprehensive `GoStarterPageObject` class:

```typescript
import { GoStarterPageObject } from './utils/test-helpers';

test('should configure project', async ({ page: browserPage }) => {
  const page = new GoStarterPageObject(browserPage);
  await page.goto();
  
  await page.fillProjectName('my-project');
  await page.selectProjectType('web-api');
  await page.selectFramework('gin');
  
  const previewContent = await page.getPreviewContent();
  expect(previewContent).toContain('my-project');
});
```

## Test Categories

### 1. Navigation Tests (`app-navigation.spec.ts`)

Tests core application navigation and layout:
- Page loading and panel visibility
- Progressive disclosure mode switching
- Responsive design across viewports
- State preservation during mode changes
- Accessibility and keyboard navigation

```typescript
test('should toggle between basic and advanced modes', async () => {
  await page.toggleDisclosureMode('advanced');
  expect(await page.getCurrentDisclosureMode()).toBe('advanced');
  
  const visibleFields = await page.getVisibleFields();
  expect(visibleFields).toContain('Database Driver');
});
```

### 2. Configuration Tests (`project-configuration.spec.ts`)

Tests form functionality and validation:
- Project type configuration
- Field validation and error handling
- Complex architecture selections
- Configuration persistence

```typescript
test('should validate module path format', async () => {
  await page.fillModulePath('invalid-path');
  await page.generateProject();
  
  const hasError = await page.hasValidationError('module-path');
  expect(hasError).toBe(true);
});
```

### 3. Real-time Preview Tests (`real-time-preview.spec.ts`)

Tests WebSocket functionality and live updates:
- WebSocket connection establishment
- Real-time preview updates
- File explorer synchronization
- Connection recovery

```typescript
test('should update preview in real-time', async () => {
  await page.waitForWebSocketConnection();
  
  const projectName = generateRandomProjectName();
  await page.fillProjectName(projectName);
  
  await page.waitForRealTimeUpdate();
  const previewContent = await page.getPreviewContent();
  expect(previewContent).toContain(projectName);
});
```

### 4. Generation Tests (`project-generation.spec.ts`)

Tests end-to-end project generation:
- Project generation process
- Download functionality
- Progress indicators
- Error handling

```typescript
test('should generate and download project', async () => {
  await createStandardWebAPIProject(page);
  await page.generateProject();
  
  await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
  
  const download = await page.downloadProject();
  expect(download.suggestedFilename()).toMatch(/my-web-api.*\.zip/);
});
```

### 5. AI Advisor Tests (`ai-advisor.spec.ts`)

Tests AI advisor integration:
- Advisor modal functionality
- Requirement gathering
- Recommendation display
- Configuration application

```typescript
test('should provide e-commerce recommendation', async () => {
  await page.openAIAdvisor();
  await page.fillAdvisorRequirements({
    domain: 'e-commerce',
    teamSize: 5,
    experience: 'senior'
  });
  
  const recommendation = await page.getAIRecommendation();
  expect(['web-api-ddd', 'web-api-clean']).toContain(recommendation.blueprint);
});
```

### 6. Performance & Accessibility Tests (`performance-accessibility.spec.ts`)

Tests performance and accessibility requirements:
- Page load performance
- Core Web Vitals
- WCAG compliance
- Keyboard navigation

```typescript
test('should meet Core Web Vitals', async () => {
  const loadTime = await page.measurePageLoadTime();
  expect(loadTime).toBeLessThan(3000);
  
  const fcp = await page.page.evaluate(() => /* FCP measurement */);
  expect(fcp).toBeLessThan(1800);
});
```

## Utilities and Helpers

### Page Object Methods

The `GoStarterPageObject` provides comprehensive methods:

```typescript
// Navigation
await page.goto()
await page.waitForPageLoad()

// Configuration
await page.fillProjectName(name)
await page.selectProjectType(type)
await page.toggleDisclosureMode('advanced')

// Preview and Generation
await page.waitForPreviewUpdate()
await page.getPreviewContent()
await page.generateProject()
await page.downloadProject()

// AI Advisor
await page.openAIAdvisor()
await page.fillAdvisorRequirements(requirements)
await page.getAIRecommendation()

// Validation
await page.getValidationErrors()
await page.hasValidationError(field)
```

### Test Data Generators

```typescript
import { generateRandomProjectName, generateRandomModulePath } from './utils/test-helpers';

const projectName = generateRandomProjectName();
const modulePath = generateRandomModulePath();
```

### Common Scenarios

```typescript
import { createStandardWebAPIProject, createAdvancedProject } from './utils/test-helpers';

await createStandardWebAPIProject(page);
await createAdvancedProject(page);
```

## Performance Thresholds

Our tests enforce the following performance requirements:

```typescript
export const PerformanceThresholds = {
  PAGE_LOAD_TIME: 3000,        // 3 seconds
  PREVIEW_UPDATE_TIME: 2000,   // 2 seconds
  GENERATION_TIME: 30000,      // 30 seconds
  FCP_TIME: 1800,              // 1.8 seconds (First Contentful Paint)
  LCP_TIME: 2500               // 2.5 seconds (Largest Contentful Paint)
};
```

## Accessibility Requirements

Tests verify WCAG 2.1 AA compliance:

- Color contrast ratios
- Keyboard navigation
- Screen reader compatibility
- ARIA labels and roles
- Focus management

## CI/CD Integration

E2E tests run automatically on:
- Pull requests to main/develop
- Pushes to main/develop (web-related changes)

The CI pipeline includes:
- Cross-browser testing (Chrome, Firefox, Safari)
- Mobile device testing
- Visual regression testing
- Performance monitoring
- Accessibility validation

## Debugging Tests

### Interactive Debugging

```bash
# Run with UI for interactive debugging
npm run test:e2e:ui

# Run in debug mode
npm run test:e2e:debug

# Run specific test in headed mode
npx playwright test app-navigation.spec.ts --headed --debug
```

### Test Artifacts

Failed tests automatically generate:
- Screenshots
- Videos
- Trace files
- Console logs

Access artifacts in the `test-results/` directory.

### Common Issues

1. **WebSocket Connection Failures**
   ```bash
   # Ensure web server is running
   go run cmd/web-server/main.go
   ```

2. **Port Conflicts**
   ```bash
   # Check if ports 3000 and 8080 are available
   lsof -i :3000
   lsof -i :8080
   ```

3. **Timing Issues**
   ```typescript
   // Use proper waiting strategies
   await page.waitForSelector('[data-testid="element"]');
   await page.waitForLoadState('networkidle');
   ```

## Writing New Tests

### Test Structure

```typescript
import { test, expect } from '@playwright/test';
import { GoStarterPageObject } from './utils/test-helpers';

test.describe('Feature Name', () => {
  let page: GoStarterPageObject;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new GoStarterPageObject(browserPage);
    await page.goto();
  });

  test('should do something', async () => {
    // Arrange
    await page.fillProjectName('test-project');
    
    // Act
    await page.selectProjectType('web-api');
    
    // Assert
    const previewContent = await page.getPreviewContent();
    expect(previewContent).toContain('web-api');
  });
});
```

### Best Practices

1. **Use Data Test IDs**: Prefer `data-testid` attributes over CSS selectors
2. **Wait for Updates**: Always wait for async operations to complete
3. **Isolate Tests**: Each test should be independent
4. **Clear Assertions**: Use descriptive assertion messages
5. **Mock External Services**: Mock APIs and external dependencies

### Test Naming

Use descriptive test names that explain the expected behavior:

```typescript
test('should show validation error for invalid module path')
test('should update preview when project type changes')
test('should preserve configuration during disclosure mode switch')
```

## Continuous Improvement

### Metrics We Track

- Test execution time
- Flakiness rates
- Coverage of user journeys
- Performance regression detection

### Regular Maintenance

- Update Playwright and dependencies monthly
- Review and update performance thresholds quarterly
- Add tests for new features
- Refactor tests to reduce maintenance burden

## Getting Help

- **Playwright Documentation**: https://playwright.dev/docs/intro
- **Project Issues**: Submit GitHub issues for test-related problems
- **Team Support**: Reach out to the frontend team for test writing guidance

## Contributing

When adding new features to the web interface:

1. Write E2E tests before implementation
2. Ensure tests pass in all browsers
3. Update test documentation
4. Add performance and accessibility considerations
5. Review test coverage with the team