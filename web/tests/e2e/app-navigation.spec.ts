import { test, expect } from '@playwright/test';
import { GoStarterPageObject } from './utils/test-helpers';

test.describe('App Navigation and Layout', () => {
  let page: GoStarterPageObject;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new GoStarterPageObject(browserPage);
    await page.goto();
  });

  test('should load the main page with all panels visible', async () => {
    // Check header is visible
    await expect(page.page.getByRole('heading', { name: /go-starter/i })).toBeVisible();
    
    // Check all three main panels are visible
    await expect(page.page.getByTestId('configuration-panel')).toBeVisible();
    await expect(page.page.getByTestId('preview-panel')).toBeVisible();
    await expect(page.page.getByTestId('file-explorer-panel')).toBeVisible();
  });

  test('should toggle between basic and advanced disclosure modes', async () => {
    // Start in basic mode
    expect(await page.getCurrentDisclosureMode()).toBe('basic');
    
    // Switch to advanced mode
    await page.toggleDisclosureMode('advanced');
    expect(await page.getCurrentDisclosureMode()).toBe('advanced');
    
    // Check that advanced fields are now visible
    const visibleFields = await page.getVisibleFields();
    expect(visibleFields).toContain('Database Driver');
    expect(visibleFields).toContain('Authentication Type');
    
    // Switch back to basic mode
    await page.toggleDisclosureMode('basic');
    expect(await page.getCurrentDisclosureMode()).toBe('basic');
    
    // Check that advanced fields are hidden
    const basicFields = await page.getVisibleFields();
    expect(basicFields).not.toContain('Database Driver');
    expect(basicFields).not.toContain('Authentication Type');
  });

  test('should be responsive on mobile devices', async () => {
    // Test mobile viewport
    await page.setViewportSize(375, 667); // iPhone SE size
    
    // Panels should stack vertically on mobile
    const configPanel = page.page.getByTestId('configuration-panel');
    const previewPanel = page.page.getByTestId('preview-panel');
    
    const configBox = await configPanel.boundingBox();
    const previewBox = await previewPanel.boundingBox();
    
    // On mobile, preview should be below config (higher y coordinate)
    expect(previewBox!.y).toBeGreaterThan(configBox!.y);
  });

  test('should maintain state when switching disclosure modes', async () => {
    // Fill some basic information
    await page.fillProjectName('test-project');
    await page.selectProjectType('web-api');
    await page.selectFramework('gin');
    
    // Switch to advanced mode
    await page.toggleDisclosureMode('advanced');
    
    // Check that basic information is preserved
    expect(await page.page.getByLabel(/project name/i).inputValue()).toBe('test-project');
    expect(await page.page.getByLabel(/project type/i).inputValue()).toBe('web-api');
    expect(await page.page.getByLabel(/framework/i).inputValue()).toBe('gin');
    
    // Fill advanced information
    await page.selectDatabaseDriver('postgres');
    await page.selectAuthType('jwt');
    
    // Switch back to basic mode
    await page.toggleDisclosureMode('basic');
    
    // Check that all information is preserved
    expect(await page.page.getByLabel(/project name/i).inputValue()).toBe('test-project');
    expect(await page.page.getByLabel(/project type/i).inputValue()).toBe('web-api');
    expect(await page.page.getByLabel(/framework/i).inputValue()).toBe('gin');
    
    // Switch back to advanced and check advanced fields are preserved
    await page.toggleDisclosureMode('advanced');
    expect(await page.page.getByLabel(/database driver/i).inputValue()).toBe('postgres');
    expect(await page.page.getByLabel(/authentication/i).inputValue()).toBe('jwt');
  });

  test('should show/hide appropriate fields based on project type', async () => {
    // Select CLI project type
    await page.selectProjectType('cli');
    
    // Framework should default to cobra and not show web frameworks
    await expect(page.page.getByRole('option', { name: 'cobra' })).toBeVisible();
    
    // Switch to web-api
    await page.selectProjectType('web-api');
    
    // Should show web frameworks
    await expect(page.page.getByRole('option', { name: 'gin' })).toBeVisible();
    await expect(page.page.getByRole('option', { name: 'echo' })).toBeVisible();
    await expect(page.page.getByRole('option', { name: 'fiber' })).toBeVisible();
  });

  test('should handle page refresh gracefully', async () => {
    // Fill form data
    await page.fillProjectName('persistent-project');
    await page.selectProjectType('web-api');
    await page.fillModulePath('github.com/user/persistent');
    
    // Refresh the page
    await page.page.reload();
    await page.waitForPageLoad();
    
    // Check if the app loads correctly after refresh
    await expect(page.page.getByRole('heading', { name: /go-starter/i })).toBeVisible();
    await expect(page.page.getByTestId('configuration-panel')).toBeVisible();
    
    // Note: Form data persistence depends on implementation
    // This test ensures the app doesn't crash on refresh
  });

  test('should have proper ARIA labels and accessibility features', async () => {
    // Check for proper ARIA labels on form controls
    await expect(page.page.getByLabel(/project name/i)).toBeVisible();
    await expect(page.page.getByLabel(/project type/i)).toBeVisible();
    await expect(page.page.getByLabel(/framework/i)).toBeVisible();
    
    // Check for keyboard navigation
    await page.page.keyboard.press('Tab');
    const focusedElement = await page.page.evaluate(() => document.activeElement?.tagName);
    expect(['INPUT', 'SELECT', 'BUTTON']).toContain(focusedElement);
    
    // Check for semantic structure
    await expect(page.page.getByRole('main')).toBeVisible();
    await expect(page.page.getByRole('heading', { level: 1 })).toBeVisible();
  });

  test('should load within acceptable time limits', async () => {
    const loadTime = await page.measurePageLoadTime();
    
    // Page should load within 5 seconds
    expect(loadTime).toBeLessThan(5000);
    
    // Preview updates should be responsive
    const updateTime = await page.measurePreviewUpdateTime();
    expect(updateTime).toBeLessThan(2000);
  });
});