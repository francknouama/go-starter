import { test, expect } from '@playwright/test';
import { GoStarterPageObject, createStandardWebAPIProject, createAdvancedProject, generateRandomProjectName, generateRandomModulePath } from './utils/test-helpers';

test.describe('Project Configuration', () => {
  let page: GoStarterPageObject;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new GoStarterPageObject(browserPage);
    await page.goto();
  });

  test('should configure a standard web API project', async () => {
    await createStandardWebAPIProject(page);
    
    // Verify configuration is applied
    expect(await page.page.getByLabel(/project name/i).inputValue()).toBe('my-web-api');
    expect(await page.page.getByLabel(/project type/i).inputValue()).toBe('web-api');
    expect(await page.page.getByLabel(/framework/i).inputValue()).toBe('gin');
    expect(await page.page.getByLabel(/logger/i).inputValue()).toBe('slog');
    expect(await page.page.getByLabel(/module path/i).inputValue()).toBe('github.com/user/my-web-api');
    
    // Preview should update
    await page.waitForPreviewUpdate();
    const previewContent = await page.getPreviewContent();
    expect(previewContent).toContain('my-web-api');
    expect(previewContent).toContain('gin');
  });

  test('should configure advanced project with all options', async () => {
    await createAdvancedProject(page);
    
    // Verify all advanced configuration is applied
    expect(await page.page.getByLabel(/project name/i).inputValue()).toBe('enterprise-api');
    expect(await page.page.getByLabel(/architecture/i).inputValue()).toBe('clean');
    expect(await page.page.getByLabel(/database driver/i).inputValue()).toBe('postgres');
    expect(await page.page.getByLabel(/authentication/i).inputValue()).toBe('jwt');
    
    // Should show more estimated files for advanced configuration
    const estimatedFiles = await page.getEstimatedFiles();
    expect(estimatedFiles).toBeGreaterThan(30); // Advanced projects have more files
  });

  test('should validate required fields', async () => {
    // Try to generate without required fields
    await page.generateProject();
    
    // Should show validation errors
    const errors = await page.getValidationErrors();
    expect(errors.length).toBeGreaterThan(0);
    expect(errors.some(error => error.includes('Project name'))).toBe(true);
    expect(errors.some(error => error.includes('Module path'))).toBe(true);
  });

  test('should validate module path format', async () => {
    await page.fillProjectName('test-project');
    
    // Test invalid module paths
    const invalidPaths = [
      'invalid-path',
      'http://github.com/user/project',
      'github.com/',
      'github.com/user/',
      'GITHUB.COM/USER/PROJECT' // uppercase
    ];
    
    for (const invalidPath of invalidPaths) {
      await page.fillModulePath(invalidPath);
      await page.generateProject();
      
      const hasError = await page.hasValidationError('module-path');
      expect(hasError).toBe(true);
    }
    
    // Test valid module path
    await page.fillModulePath('github.com/user/valid-project');
    await page.generateProject();
    const hasError = await page.hasValidationError('module-path');
    expect(hasError).toBe(false);
  });

  test('should handle different project types correctly', async () => {
    const projectTypes = [
      { type: 'web-api', expectedFrameworks: ['gin', 'echo', 'fiber'] },
      { type: 'cli', expectedFrameworks: ['cobra'] },
      { type: 'library', expectedFrameworks: [] }, // Libraries might not have frameworks
      { type: 'lambda', expectedFrameworks: [] } // Lambdas might not have traditional frameworks
    ];

    for (const { type, expectedFrameworks } of projectTypes) {
      await page.selectProjectType(type);
      await page.fillProjectName(`test-${type}`);
      await page.fillModulePath(`github.com/user/test-${type}`);
      
      // Check that appropriate frameworks are available
      if (expectedFrameworks.length > 0) {
        for (const framework of expectedFrameworks) {
          await expect(page.page.getByRole('option', { name: framework })).toBeVisible();
        }
      }
      
      // Preview should update based on project type
      await page.waitForPreviewUpdate();
      const previewContent = await page.getPreviewContent();
      expect(previewContent).toContain(type);
    }
  });

  test('should update file count based on configuration', async () => {
    // Start with minimal configuration
    await page.fillProjectName('minimal-project');
    await page.selectProjectType('cli');
    await page.fillModulePath('github.com/user/minimal');
    
    const minimalFiles = await page.getEstimatedFiles();
    
    // Switch to more complex configuration
    await page.selectProjectType('web-api');
    await page.toggleDisclosureMode('advanced');
    await page.selectArchitecture('clean');
    await page.selectDatabaseDriver('postgres');
    await page.selectAuthType('jwt');
    
    await page.waitForPreviewUpdate();
    const complexFiles = await page.getEstimatedFiles();
    
    // Complex configuration should have more files
    expect(complexFiles).toBeGreaterThan(minimalFiles);
  });

  test('should preserve configuration during disclosure mode changes', async () => {
    // Configure in basic mode
    await page.fillProjectName('persistent-project');
    await page.selectProjectType('web-api');
    await page.selectFramework('gin');
    await page.selectLogger('zap');
    await page.fillModulePath('github.com/user/persistent');
    
    // Switch to advanced mode and add more configuration
    await page.toggleDisclosureMode('advanced');
    await page.selectArchitecture('ddd');
    await page.selectDatabaseDriver('postgres');
    
    // Switch back to basic mode
    await page.toggleDisclosureMode('basic');
    
    // Verify basic configuration is preserved
    expect(await page.page.getByLabel(/project name/i).inputValue()).toBe('persistent-project');
    expect(await page.page.getByLabel(/project type/i).inputValue()).toBe('web-api');
    expect(await page.page.getByLabel(/framework/i).inputValue()).toBe('gin');
    expect(await page.page.getByLabel(/logger/i).inputValue()).toBe('zap');
    
    // Switch back to advanced and verify advanced configuration is preserved
    await page.toggleDisclosureMode('advanced');
    expect(await page.page.getByLabel(/architecture/i).inputValue()).toBe('ddd');
    expect(await page.page.getByLabel(/database driver/i).inputValue()).toBe('postgres');
  });

  test('should handle random project generation', async () => {
    const projectName = generateRandomProjectName();
    const modulePath = generateRandomModulePath();
    
    await page.fillProjectName(projectName);
    await page.selectProjectType('web-api');
    await page.selectFramework('gin');
    await page.fillModulePath(modulePath);
    
    await page.waitForPreviewUpdate();
    const previewContent = await page.getPreviewContent();
    
    // Preview should contain the random values
    expect(previewContent).toContain(projectName);
    expect(previewContent).toContain(modulePath.split('/').pop()); // project name from module path
  });

  test('should support complex architecture configurations', async () => {
    const architectures = ['standard', 'clean', 'ddd', 'hexagonal'];
    
    await page.toggleDisclosureMode('advanced');
    await page.fillProjectName('arch-test');
    await page.selectProjectType('web-api');
    await page.fillModulePath('github.com/user/arch-test');
    
    for (const architecture of architectures) {
      await page.selectArchitecture(architecture);
      await page.waitForPreviewUpdate();
      
      const estimatedFiles = await page.getEstimatedFiles();
      const previewContent = await page.getPreviewContent();
      
      // Each architecture should have different file counts and structure
      expect(estimatedFiles).toBeGreaterThan(10);
      expect(previewContent).toContain(architecture);
      
      // More complex architectures should have more files
      if (architecture === 'ddd' || architecture === 'hexagonal') {
        expect(estimatedFiles).toBeGreaterThan(40);
      }
    }
  });

  test('should handle logger selection correctly', async () => {
    const loggers = ['slog', 'zap', 'logrus', 'zerolog'];
    
    await page.fillProjectName('logger-test');
    await page.selectProjectType('web-api');
    await page.fillModulePath('github.com/user/logger-test');
    
    for (const logger of loggers) {
      await page.selectLogger(logger);
      await page.waitForPreviewUpdate();
      
      const previewContent = await page.getPreviewContent();
      expect(previewContent).toContain(logger);
    }
  });

  test('should show appropriate defaults for different project types', async () => {
    // CLI projects should default to cobra
    await page.selectProjectType('cli');
    expect(await page.page.getByLabel(/framework/i).inputValue()).toBe('cobra');
    
    // Web API projects should default to a web framework
    await page.selectProjectType('web-api');
    const webFramework = await page.page.getByLabel(/framework/i).inputValue();
    expect(['gin', 'echo', 'fiber']).toContain(webFramework);
    
    // All projects should default to slog logger
    expect(await page.page.getByLabel(/logger/i).inputValue()).toBe('slog');
  });
});