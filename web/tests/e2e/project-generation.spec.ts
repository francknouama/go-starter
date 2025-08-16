import { test, expect } from '@playwright/test';
import { GoStarterPageObject, createStandardWebAPIProject, createCLIProject, generateRandomProjectName } from './utils/test-helpers';

test.describe('Project Generation and Download', () => {
  let page: GoStarterPageObject;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new GoStarterPageObject(browserPage);
    await page.goto();
  });

  test('should generate and download a web API project', async () => {
    await createStandardWebAPIProject(page);
    
    // Generate the project
    await page.generateProject();
    
    // Wait for generation to complete
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    // Download should be available
    await expect(page.page.getByRole('button', { name: /download/i })).toBeEnabled();
    
    // Download the project
    const download = await page.downloadProject();
    
    // Verify download
    expect(download.suggestedFilename()).toMatch(/my-web-api.*\.zip/);
    expect(await download.path()).toBeTruthy();
  });

  test('should generate a CLI project with correct structure', async () => {
    await createCLIProject(page);
    
    // Generate the project
    await page.generateProject();
    
    // Wait for generation to complete
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    // Verify the generated file structure shows CLI-specific files
    const fileList = await page.getFileList();
    expect(fileList.some(file => file.includes('cmd/'))).toBe(true);
    expect(fileList.some(file => file.includes('main.go'))).toBe(true);
    expect(fileList.some(file => file.includes('cobra'))).toBe(true);
  });

  test('should handle generation errors gracefully', async () => {
    // Try to generate without required fields
    await page.generateProject();
    
    // Should show validation errors instead of proceeding
    const errors = await page.getValidationErrors();
    expect(errors.length).toBeGreaterThan(0);
    
    // Generation status should indicate error
    await expect(page.page.getByTestId('generation-status')).toHaveText(/error|failed/i);
    
    // Download button should be disabled
    await expect(page.page.getByRole('button', { name: /download/i })).toBeDisabled();
  });

  test('should show generation progress', async () => {
    await createStandardWebAPIProject(page);
    
    // Start generation
    const generateButton = page.page.getByRole('button', { name: /generate project/i });
    await generateButton.click();
    
    // Should show progress indicators
    await expect(page.page.getByTestId('generation-progress')).toBeVisible();
    await expect(page.page.getByTestId('generation-status')).toHaveText(/generating|in progress/i);
    
    // Generate button should be disabled during generation
    await expect(generateButton).toBeDisabled();
    
    // Wait for completion
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete', {
      timeout: 30000 // Allow time for generation
    });
    
    // Progress indicator should be hidden
    await expect(page.page.getByTestId('generation-progress')).not.toBeVisible();
    
    // Generate button should be enabled again
    await expect(generateButton).toBeEnabled();
  });

  test('should generate different file counts for different complexities', async () => {
    const projectName = generateRandomProjectName();
    
    // Test simple CLI
    await page.fillProjectName(`${projectName}-simple`);
    await page.selectProjectType('cli');
    await page.fillModulePath(`github.com/user/${projectName}-simple`);
    
    await page.generateProject();
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    const simpleFiles = await page.getEstimatedFiles();
    
    // Reset and test complex web API
    await page.page.reload();
    await page.waitForPageLoad();
    
    await page.fillProjectName(`${projectName}-complex`);
    await page.selectProjectType('web-api');
    await page.toggleDisclosureMode('advanced');
    await page.selectArchitecture('clean');
    await page.selectDatabaseDriver('postgres');
    await page.selectAuthType('jwt');
    await page.fillModulePath(`github.com/user/${projectName}-complex`);
    
    await page.generateProject();
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    const complexFiles = await page.getEstimatedFiles();
    
    // Complex project should have significantly more files
    expect(complexFiles).toBeGreaterThan(simpleFiles * 2);
  });

  test('should preserve generation results during page interactions', async () => {
    await createStandardWebAPIProject(page);
    await page.generateProject();
    
    // Wait for generation to complete
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    // Make configuration changes
    await page.fillProjectName('modified-name');
    await page.selectFramework('echo');
    
    // Generation status should indicate that changes were made
    await expect(page.page.getByTestId('generation-status')).toHaveText(/modified|outdated/i);
    
    // Download button should indicate it's for the previous generation
    await expect(page.page.getByRole('button', { name: /download previous/i })).toBeVisible();
    
    // New generation should be available
    await expect(page.page.getByRole('button', { name: /regenerate|generate new/i })).toBeEnabled();
  });

  test('should handle large project generation', async () => {
    // Configure a large, complex project
    await page.fillProjectName('enterprise-project');
    await page.selectProjectType('web-api');
    await page.toggleDisclosureMode('advanced');
    await page.selectArchitecture('ddd');
    await page.selectFramework('gin');
    await page.selectLogger('zap');
    await page.selectDatabaseDriver('postgres');
    await page.selectAuthType('jwt');
    await page.fillModulePath('github.com/enterprise/large-project');
    
    // Enable multiple features
    await page.enableFeature('monitoring');
    await page.enableFeature('testing');
    await page.enableFeature('documentation');
    
    // Generate the project
    await page.generateProject();
    
    // Should handle large generation (might take longer)
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete', {
      timeout: 60000 // Allow up to 1 minute for large projects
    });
    
    // Should have many files
    const fileCount = await page.getEstimatedFiles();
    expect(fileCount).toBeGreaterThan(50);
    
    // Download should still work
    const download = await page.downloadProject();
    expect(download.suggestedFilename()).toMatch(/enterprise-project.*\.zip/);
  });

  test('should support regeneration with different configurations', async () => {
    // Initial generation
    await page.fillProjectName('regen-test');
    await page.selectProjectType('cli');
    await page.fillModulePath('github.com/user/regen-test');
    
    await page.generateProject();
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    const initialFiles = await page.getEstimatedFiles();
    
    // Change configuration
    await page.selectProjectType('web-api');
    await page.selectFramework('gin');
    
    // Regenerate
    await page.generateProject();
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    const newFiles = await page.getEstimatedFiles();
    
    // File count should be different
    expect(newFiles).not.toBe(initialFiles);
    expect(newFiles).toBeGreaterThan(initialFiles); // Web API typically has more files than CLI
  });

  test('should validate module path before generation', async () => {
    await page.fillProjectName('validation-test');
    await page.selectProjectType('web-api');
    
    // Test various invalid module paths
    const invalidPaths = [
      'invalid',
      'github.com/',
      'http://github.com/user/project',
      'github.com/user/-invalid',
      'github.com/user/invalid-'
    ];
    
    for (const invalidPath of invalidPaths) {
      await page.fillModulePath(invalidPath);
      await page.generateProject();
      
      // Should show validation error
      const hasError = await page.hasValidationError('module-path');
      expect(hasError).toBe(true);
      
      // Should not proceed with generation
      await expect(page.page.getByTestId('generation-status')).not.toHaveText('Generation complete');
    }
    
    // Valid path should work
    await page.fillModulePath('github.com/user/valid-project');
    await page.generateProject();
    
    // Should not have validation error
    const hasError = await page.hasValidationError('module-path');
    expect(hasError).toBe(false);
    
    // Should proceed with generation
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
  });

  test('should handle concurrent generation requests', async () => {
    await createStandardWebAPIProject(page);
    
    // Start first generation
    const generateButton = page.page.getByRole('button', { name: /generate project/i });
    await generateButton.click();
    
    // Try to start another generation while first is in progress
    await expect(generateButton).toBeDisabled();
    
    // Should show that generation is in progress
    await expect(page.page.getByTestId('generation-status')).toHaveText(/generating|in progress/i);
    
    // Wait for completion
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    // Button should be enabled again
    await expect(generateButton).toBeEnabled();
  });

  test('should provide download with correct filename and content', async () => {
    const projectName = 'download-test-project';
    
    await page.fillProjectName(projectName);
    await page.selectProjectType('web-api');
    await page.selectFramework('gin');
    await page.fillModulePath(`github.com/user/${projectName}`);
    
    await page.generateProject();
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    // Download the project
    const download = await page.downloadProject();
    
    // Verify filename contains project name
    expect(download.suggestedFilename()).toContain(projectName);
    expect(download.suggestedFilename()).toMatch(/\.zip$/);
    
    // Verify download completed
    const downloadPath = await download.path();
    expect(downloadPath).toBeTruthy();
    
    // Verify file size is reasonable (not empty, not too small)
    const fs = require('fs');
    if (downloadPath) {
      const stats = fs.statSync(downloadPath);
      expect(stats.size).toBeGreaterThan(1000); // At least 1KB
      expect(stats.size).toBeLessThan(10 * 1024 * 1024); // Less than 10MB
    }
  });
});