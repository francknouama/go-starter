import { test as base, expect } from '@playwright/test';
import { GoStarterPageObject } from './utils/test-helpers';

// Extend base test to include our page object
export const test = base.extend<{ goStarterPage: GoStarterPageObject }>({
  goStarterPage: async ({ page }, use) => {
    const goStarterPage = new GoStarterPageObject(page);
    await use(goStarterPage);
  },
});

export { expect };

// Global test configuration
test.beforeEach(async ({ page }) => {
  // Set up common test environment
  await page.addInitScript(() => {
    // Disable animations for consistent testing
    const style = document.createElement('style');
    style.textContent = `
      *, *::before, *::after {
        animation-duration: 0.01ms !important;
        animation-delay: 0.01ms !important;
        transition-duration: 0.01ms !important;
        transition-delay: 0.01ms !important;
      }
    `;
    document.head.appendChild(style);
    
    // Mock console for error tracking
    window.__testErrors = [];
    const originalError = console.error;
    console.error = (...args) => {
      window.__testErrors.push(args.join(' '));
      originalError(...args);
    };
  });

  // Set up test data attributes for better element selection
  await page.addInitScript(() => {
    // Add test IDs to elements that might not have them
    window.addEventListener('DOMContentLoaded', () => {
      // Configuration panel
      const configPanel = document.querySelector('[data-testid="configuration-panel"]') ||
                         document.querySelector('form') ||
                         document.querySelector('.configuration');
      if (configPanel && !configPanel.getAttribute('data-testid')) {
        configPanel.setAttribute('data-testid', 'configuration-panel');
      }

      // Preview panel
      const previewPanel = document.querySelector('[data-testid="preview-panel"]') ||
                          document.querySelector('.preview') ||
                          document.querySelector('.code-preview');
      if (previewPanel && !previewPanel.getAttribute('data-testid')) {
        previewPanel.setAttribute('data-testid', 'preview-panel');
      }

      // File explorer panel
      const filePanel = document.querySelector('[data-testid="file-explorer-panel"]') ||
                       document.querySelector('.file-explorer') ||
                       document.querySelector('.files');
      if (filePanel && !filePanel.getAttribute('data-testid')) {
        filePanel.setAttribute('data-testid', 'file-explorer-panel');
      }

      // Generate button
      const generateBtn = document.querySelector('[data-testid="generate-button"]') ||
                         document.querySelector('button[type="submit"]') ||
                         document.querySelector('button:contains("Generate")');
      if (generateBtn && !generateBtn.getAttribute('data-testid')) {
        generateBtn.setAttribute('data-testid', 'generate-button');
      }
    });
  });
});

// Global test cleanup
test.afterEach(async ({ page }) => {
  // Check for console errors during the test
  const errors = await page.evaluate(() => window.__testErrors || []);
  if (errors.length > 0) {
    console.warn('Console errors during test:', errors);
  }

  // Clean up any downloads or temporary files
  // (Playwright handles this automatically, but we can add custom cleanup here)
});

// Custom matchers for go-starter specific assertions
expect.extend({
  async toBeValidGoModule(received: string) {
    const goModuleRegex = /^[a-zA-Z0-9.-]+\/[a-zA-Z0-9._/-]+$/;
    const pass = goModuleRegex.test(received);
    
    return {
      message: () => pass 
        ? `Expected "${received}" not to be a valid Go module path`
        : `Expected "${received}" to be a valid Go module path`,
      pass,
    };
  },

  async toContainValidFileStructure(received: string[]) {
    const requiredFiles = ['main.go', 'go.mod'];
    const hasRequired = requiredFiles.every(file => 
      received.some(f => f.includes(file))
    );
    
    return {
      message: () => hasRequired
        ? `Expected file list not to contain required Go files`
        : `Expected file list to contain required Go files: ${requiredFiles.join(', ')}`,
      pass: hasRequired,
    };
  },

  async toHaveValidProjectStructure(received: string[]) {
    const structureChecks = {
      hasMainFile: received.some(f => f.includes('main.go')),
      hasGoMod: received.some(f => f.includes('go.mod')),
      hasInternalDir: received.some(f => f.includes('internal/')),
      hasReadme: received.some(f => f.includes('README.md'))
    };

    const pass = Object.values(structureChecks).every(check => check);
    const failedChecks = Object.entries(structureChecks)
      .filter(([_, passed]) => !passed)
      .map(([check, _]) => check);

    return {
      message: () => pass
        ? `Expected project structure to be invalid`
        : `Expected valid project structure. Missing: ${failedChecks.join(', ')}`,
      pass,
    };
  }
});

// Test utilities for common scenarios
export const TestScenarios = {
  BASIC_WEB_API: {
    name: 'basic-web-api',
    type: 'web-api',
    framework: 'gin',
    logger: 'slog',
    module: 'github.com/user/basic-web-api'
  },

  ADVANCED_ENTERPRISE: {
    name: 'enterprise-api',
    type: 'web-api',
    architecture: 'clean',
    framework: 'gin',
    logger: 'zap',
    database: 'postgres',
    auth: 'jwt',
    module: 'github.com/enterprise/api'
  },

  SIMPLE_CLI: {
    name: 'simple-cli',
    type: 'cli',
    framework: 'cobra',
    logger: 'slog',
    module: 'github.com/user/simple-cli'
  },

  MICROSERVICE: {
    name: 'user-service',
    type: 'microservice',
    framework: 'gin',
    logger: 'zap',
    database: 'postgres',
    module: 'github.com/company/user-service'
  }
};

// Performance thresholds
export const PerformanceThresholds = {
  PAGE_LOAD_TIME: 3000, // 3 seconds
  PREVIEW_UPDATE_TIME: 2000, // 2 seconds
  GENERATION_TIME: 30000, // 30 seconds
  WEBSOCKET_CONNECT_TIME: 5000, // 5 seconds
  FCP_TIME: 1800, // 1.8 seconds
  LCP_TIME: 2500 // 2.5 seconds
};

// Accessibility requirements
export const AccessibilityRequirements = {
  MIN_COLOR_CONTRAST: 4.5, // WCAG AA standard
  MAX_TAB_STOPS: 50, // Reasonable number of focusable elements
  REQUIRED_ARIA_LABELS: ['project-name', 'project-type', 'framework'],
  REQUIRED_HEADINGS: { h1: 1, h2: 3 } // At least 1 h1, 3 h2s
};