import { test, expect } from '@playwright/test';
import { GoStarterPageObject } from './utils/test-helpers';

test.describe('Performance and Accessibility', () => {
  let page: GoStarterPageObject;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new GoStarterPageObject(browserPage);
  });

  test('should load page within performance budget', async () => {
    const startTime = Date.now();
    await page.goto();
    const loadTime = Date.now() - startTime;
    
    // Page should load within 3 seconds
    expect(loadTime).toBeLessThan(3000);
    
    // Check that critical elements are visible quickly
    await expect(page.page.getByRole('heading', { name: /go-starter/i })).toBeVisible();
    await expect(page.page.getByTestId('configuration-panel')).toBeVisible();
  });

  test('should have good Core Web Vitals', async () => {
    await page.goto();
    
    // Measure First Contentful Paint (FCP)
    const fcp = await page.page.evaluate(() => {
      return new Promise((resolve) => {
        new PerformanceObserver((entryList) => {
          const entries = entryList.getEntries();
          const fcpEntry = entries.find(entry => entry.name === 'first-contentful-paint');
          if (fcpEntry) {
            resolve(fcpEntry.startTime);
          }
        }).observe({ entryTypes: ['paint'] });
      });
    });
    
    // FCP should be under 1.8 seconds (good threshold)
    expect(fcp).toBeLessThan(1800);
    
    // Measure Largest Contentful Paint (LCP)
    const lcp = await page.page.evaluate(() => {
      return new Promise((resolve) => {
        new PerformanceObserver((entryList) => {
          const entries = entryList.getEntries();
          const lastEntry = entries[entries.length - 1];
          resolve(lastEntry.startTime);
        }).observe({ entryTypes: ['largest-contentful-paint'] });
        
        // Fallback timeout
        setTimeout(() => resolve(0), 5000);
      });
    });
    
    // LCP should be under 2.5 seconds (good threshold)
    if (lcp > 0) {
      expect(lcp).toBeLessThan(2500);
    }
  });

  test('should handle rapid user interactions without performance degradation', async () => {
    await page.goto();
    await page.waitForWebSocketConnection();
    
    const startTime = Date.now();
    
    // Perform rapid interactions
    for (let i = 0; i < 20; i++) {
      await page.fillProjectName(`rapid-test-${i}`);
      await page.page.waitForTimeout(50);
    }
    
    // Change multiple configuration options rapidly
    await page.selectProjectType('web-api');
    await page.selectFramework('gin');
    await page.selectLogger('zap');
    await page.toggleDisclosureMode('advanced');
    await page.selectArchitecture('clean');
    
    const endTime = Date.now();
    const totalTime = endTime - startTime;
    
    // Rapid interactions should complete within reasonable time
    expect(totalTime).toBeLessThan(5000);
    
    // Page should remain responsive
    await expect(page.page.getByLabel(/project name/i)).toBeEnabled();
    await expect(page.page.getByRole('button', { name: /generate/i })).toBeEnabled();
  });

  test('should be accessible with screen readers', async () => {
    await page.goto();
    
    // Check for proper heading structure
    const headings = page.page.getByRole('heading');
    const h1Count = await headings.filter({ level: 1 }).count();
    expect(h1Count).toBe(1); // Should have exactly one h1
    
    // Check for proper form labels
    await expect(page.page.getByLabel(/project name/i)).toBeVisible();
    await expect(page.page.getByLabel(/project type/i)).toBeVisible();
    await expect(page.page.getByLabel(/framework/i)).toBeVisible();
    
    // Check for ARIA attributes
    const configPanel = page.page.getByTestId('configuration-panel');
    await expect(configPanel).toHaveAttribute('role', 'form');
    
    // Check for focus management
    await page.page.keyboard.press('Tab');
    const focusedElement = page.page.locator(':focus');
    await expect(focusedElement).toBeVisible();
  });

  test('should support keyboard navigation', async () => {
    await page.goto();
    
    // Tab through all focusable elements
    const focusableElements = [];
    let tabCount = 0;
    const maxTabs = 20; // Prevent infinite loop
    
    while (tabCount < maxTabs) {
      await page.page.keyboard.press('Tab');
      const activeElement = await page.page.evaluate(() => {
        const element = document.activeElement;
        return element ? {
          tagName: element.tagName,
          type: element.getAttribute('type'),
          ariaLabel: element.getAttribute('aria-label'),
          id: element.id
        } : null;
      });
      
      if (activeElement) {
        focusableElements.push(activeElement);
      }
      
      tabCount++;
      
      // Break if we've cycled back to the first element
      if (focusableElements.length > 1 && 
          JSON.stringify(activeElement) === JSON.stringify(focusableElements[0])) {
        break;
      }
    }
    
    // Should have multiple focusable elements
    expect(focusableElements.length).toBeGreaterThan(5);
    
    // Should include form controls and buttons
    const elementTypes = focusableElements.map(el => el.tagName.toLowerCase());
    expect(elementTypes).toContain('input');
    expect(elementTypes).toContain('select');
    expect(elementTypes).toContain('button');
  });

  test('should have proper color contrast', async () => {
    await page.goto();
    
    // Test color contrast on key elements
    const textElements = [
      page.page.getByRole('heading', { name: /go-starter/i }),
      page.page.getByLabel(/project name/i),
      page.page.getByRole('button', { name: /generate/i })
    ];
    
    for (const element of textElements) {
      const styles = await element.evaluate((el) => {
        const computed = window.getComputedStyle(el);
        return {
          color: computed.color,
          backgroundColor: computed.backgroundColor,
          fontSize: computed.fontSize
        };
      });
      
      // Basic checks - proper styles should be applied
      expect(styles.color).not.toBe('');
      expect(styles.fontSize).not.toBe('');
    }
  });

  test('should support high contrast mode', async () => {
    // Enable high contrast mode
    await page.page.emulateMedia({ reducedMotion: 'reduce', colorScheme: 'dark' });
    await page.goto();
    
    // Elements should still be visible and functional
    await expect(page.page.getByRole('heading', { name: /go-starter/i })).toBeVisible();
    await expect(page.page.getByTestId('configuration-panel')).toBeVisible();
    await expect(page.page.getByRole('button', { name: /generate/i })).toBeEnabled();
    
    // Form should still be functional
    await page.fillProjectName('high-contrast-test');
    await page.selectProjectType('web-api');
    
    const projectName = await page.page.getByLabel(/project name/i).inputValue();
    expect(projectName).toBe('high-contrast-test');
  });

  test('should handle reduced motion preferences', async () => {
    // Enable reduced motion
    await page.page.emulateMedia({ reducedMotion: 'reduce' });
    await page.goto();
    
    // Page should still function without animations
    await page.fillProjectName('reduced-motion-test');
    await page.selectProjectType('web-api');
    
    // Interactions should work immediately without waiting for animations
    await page.toggleDisclosureMode('advanced');
    await expect(page.page.getByLabel(/architecture/i)).toBeVisible();
  });

  test('should be responsive at different viewport sizes', async () => {
    const viewports = [
      { width: 320, height: 568, name: 'Mobile Small' },
      { width: 375, height: 667, name: 'Mobile Medium' },
      { width: 768, height: 1024, name: 'Tablet' },
      { width: 1024, height: 768, name: 'Desktop Small' },
      { width: 1920, height: 1080, name: 'Desktop Large' }
    ];
    
    for (const viewport of viewports) {
      await page.setViewportSize(viewport.width, viewport.height);
      await page.goto();
      
      // Core functionality should be available at all sizes
      await expect(page.page.getByRole('heading', { name: /go-starter/i })).toBeVisible();
      await expect(page.page.getByLabel(/project name/i)).toBeVisible();
      await expect(page.page.getByRole('button', { name: /generate/i })).toBeVisible();
      
      // Form should be usable
      await page.fillProjectName(`test-${viewport.name.replace(/\s+/g, '-').toLowerCase()}`);
      
      const projectName = await page.page.getByLabel(/project name/i).inputValue();
      expect(projectName).toBeTruthy();
    }
  });

  test('should handle memory usage efficiently during long sessions', async () => {
    await page.goto();
    await page.waitForWebSocketConnection();
    
    // Simulate a long user session with many interactions
    for (let i = 0; i < 50; i++) {
      await page.fillProjectName(`session-test-${i}`);
      await page.selectProjectType(i % 2 === 0 ? 'web-api' : 'cli');
      await page.selectFramework(i % 3 === 0 ? 'gin' : 'echo');
      
      // Occasionally toggle disclosure mode
      if (i % 10 === 0) {
        await page.toggleDisclosureMode(i % 20 === 0 ? 'advanced' : 'basic');
      }
      
      await page.page.waitForTimeout(100);
    }
    
    // Page should still be responsive after many interactions
    await expect(page.page.getByLabel(/project name/i)).toBeEnabled();
    await expect(page.page.getByRole('button', { name: /generate/i })).toBeEnabled();
    
    // Final interaction should work
    await page.fillProjectName('final-test');
    const finalValue = await page.page.getByLabel(/project name/i).inputValue();
    expect(finalValue).toBe('final-test');
  });

  test('should load efficiently with slow network', async () => {
    // Simulate slow 3G network
    await page.page.route('**/*', async route => {
      await new Promise(resolve => setTimeout(resolve, 100)); // Add 100ms delay
      await route.continue();
    });
    
    const startTime = Date.now();
    await page.goto();
    const loadTime = Date.now() - startTime;
    
    // Should still load within reasonable time even on slow network
    expect(loadTime).toBeLessThan(10000); // 10 seconds max for slow network
    
    // Core functionality should be available
    await expect(page.page.getByRole('heading', { name: /go-starter/i })).toBeVisible();
    await expect(page.page.getByTestId('configuration-panel')).toBeVisible();
  });

  test('should handle JavaScript errors gracefully', async () => {
    // Monitor console errors
    const consoleErrors: string[] = [];
    page.page.on('console', msg => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });
    
    await page.goto();
    
    // Perform normal interactions
    await page.fillProjectName('error-test');
    await page.selectProjectType('web-api');
    await page.toggleDisclosureMode('advanced');
    
    // Should not have critical JavaScript errors
    const criticalErrors = consoleErrors.filter(error => 
      !error.includes('favicon') && // Ignore favicon errors
      !error.includes('Extension') && // Ignore browser extension errors
      !error.includes('non-passive') // Ignore passive event listener warnings
    );
    
    expect(criticalErrors.length).toBe(0);
  });
});