import { test, expect } from '@playwright/test';
import { GoStarterPageObject, generateRandomProjectName } from './utils/test-helpers';

test.describe('Real-time Preview and WebSocket', () => {
  let page: GoStarterPageObject;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new GoStarterPageObject(browserPage);
    await page.goto();
  });

  test('should establish WebSocket connection', async () => {
    // Wait for WebSocket connection to be established
    await page.waitForWebSocketConnection();
    
    // Check WebSocket status indicator
    await expect(page.page.getByTestId('websocket-status')).toHaveText('connected');
  });

  test('should update preview in real-time when configuration changes', async () => {
    await page.waitForWebSocketConnection();
    
    // Initial state
    let previewContent = await page.getPreviewContent();
    const initialContent = previewContent;
    
    // Change project name
    const projectName = generateRandomProjectName();
    await page.fillProjectName(projectName);
    
    // Wait for real-time update
    await page.waitForRealTimeUpdate();
    previewContent = await page.getPreviewContent();
    
    // Preview should be updated and different from initial
    expect(previewContent).not.toBe(initialContent);
    expect(previewContent).toContain(projectName);
  });

  test('should update file explorer in real-time', async () => {
    await page.waitForWebSocketConnection();
    
    // Get initial file list
    const initialFiles = await page.getFileList();
    
    // Change project type which should affect file structure
    await page.selectProjectType('web-api');
    await page.fillProjectName('api-test');
    await page.fillModulePath('github.com/user/api-test');
    
    // Wait for update
    await page.waitForRealTimeUpdate();
    
    // File list should be updated
    const updatedFiles = await page.getFileList();
    expect(updatedFiles).not.toEqual(initialFiles);
    expect(updatedFiles.length).toBeGreaterThan(0);
    
    // Should contain typical web API files
    expect(updatedFiles.some(file => file.includes('main.go'))).toBe(true);
    expect(updatedFiles.some(file => file.includes('internal/'))).toBe(true);
  });

  test('should handle WebSocket reconnection gracefully', async () => {
    await page.waitForWebSocketConnection();
    
    // Simulate network interruption by blocking WebSocket
    await page.page.route('ws://localhost:8080/ws', route => route.abort());
    
    // WebSocket status should indicate disconnection
    await expect(page.page.getByTestId('websocket-status')).toHaveText('disconnected', { timeout: 5000 });
    
    // Re-enable WebSocket connection
    await page.page.unroute('ws://localhost:8080/ws');
    
    // Should reconnect automatically
    await page.waitForWebSocketConnection();
    await expect(page.page.getByTestId('websocket-status')).toHaveText('connected');
  });

  test('should update estimated file count in real-time', async () => {
    await page.waitForWebSocketConnection();
    
    // Start with minimal configuration
    await page.fillProjectName('count-test');
    await page.selectProjectType('cli');
    await page.fillModulePath('github.com/user/count-test');
    
    await page.waitForRealTimeUpdate();
    const initialCount = await page.getEstimatedFiles();
    
    // Switch to more complex configuration
    await page.selectProjectType('web-api');
    await page.toggleDisclosureMode('advanced');
    await page.selectArchitecture('clean');
    await page.selectDatabaseDriver('postgres');
    
    await page.waitForRealTimeUpdate();
    const updatedCount = await page.getEstimatedFiles();
    
    // File count should increase
    expect(updatedCount).toBeGreaterThan(initialCount);
  });

  test('should show loading states during preview updates', async () => {
    await page.waitForWebSocketConnection();
    
    // Make a change that triggers preview update
    await page.fillProjectName('loading-test');
    
    // Should show loading indicator briefly
    // Note: This might be fast, so we check it existed or preview is updated
    const loadingOrUpdated = await Promise.race([
      page.page.getByTestId('preview-loading').isVisible(),
      page.waitForRealTimeUpdate().then(() => true)
    ]);
    
    expect(loadingOrUpdated).toBe(true);
    
    // Loading should be gone after update
    await page.waitForRealTimeUpdate();
    await expect(page.page.getByTestId('preview-loading')).not.toBeVisible();
  });

  test('should handle multiple rapid configuration changes', async () => {
    await page.waitForWebSocketConnection();
    
    // Make multiple rapid changes
    const changes = [
      () => page.fillProjectName('rapid-test-1'),
      () => page.selectProjectType('web-api'),
      () => page.selectFramework('gin'),
      () => page.fillProjectName('rapid-test-2'),
      () => page.selectFramework('echo'),
      () => page.fillModulePath('github.com/user/rapid-test')
    ];
    
    // Execute changes rapidly
    for (const change of changes) {
      await change();
      await page.page.waitForTimeout(100); // Small delay between changes
    }
    
    // Wait for final update
    await page.waitForRealTimeUpdate();
    
    // Final state should reflect the last changes
    const previewContent = await page.getPreviewContent();
    expect(previewContent).toContain('rapid-test-2');
    expect(previewContent).toContain('echo');
  });

  test('should update preview when switching between architectures', async () => {
    await page.waitForWebSocketConnection();
    await page.toggleDisclosureMode('advanced');
    
    await page.fillProjectName('arch-preview-test');
    await page.selectProjectType('web-api');
    await page.fillModulePath('github.com/user/arch-preview-test');
    
    const architectures = ['standard', 'clean', 'ddd'];
    const previews: string[] = [];
    
    for (const architecture of architectures) {
      await page.selectArchitecture(architecture);
      await page.waitForRealTimeUpdate();
      
      const previewContent = await page.getPreviewContent();
      previews.push(previewContent);
      
      // Each architecture should produce different preview content
      expect(previewContent).toContain(architecture);
    }
    
    // All previews should be different
    expect(new Set(previews).size).toBe(architectures.length);
  });

  test('should handle file selection in explorer', async () => {
    await page.waitForWebSocketConnection();
    
    // Configure a project
    await page.fillProjectName('file-selection-test');
    await page.selectProjectType('web-api');
    await page.fillModulePath('github.com/user/file-selection-test');
    
    await page.waitForRealTimeUpdate();
    
    // Expand folders in file explorer
    await page.expandFolder('internal');
    await page.expandFolder('handlers');
    
    // Select a specific file
    await page.selectFile('main.go');
    
    // Preview should show the selected file's content
    const previewContent = await page.getPreviewContent();
    expect(previewContent).toContain('package main');
    expect(previewContent).toContain('func main()');
  });

  test('should debounce rapid updates to prevent excessive requests', async () => {
    await page.waitForWebSocketConnection();
    
    // Track network requests
    const requests: string[] = [];
    page.page.on('request', request => {
      if (request.url().includes('/api/generate') || request.url().includes('/api/preview')) {
        requests.push(request.url());
      }
    });
    
    // Make rapid changes to project name
    for (let i = 0; i < 10; i++) {
      await page.fillProjectName(`rapid-${i}`);
      await page.page.waitForTimeout(50); // Very rapid changes
    }
    
    // Wait for debouncing to settle
    await page.page.waitForTimeout(2000);
    
    // Should have made fewer requests than changes due to debouncing
    expect(requests.length).toBeLessThan(10);
    expect(requests.length).toBeGreaterThan(0);
  });

  test('should maintain preview state across WebSocket reconnections', async () => {
    await page.waitForWebSocketConnection();
    
    // Configure project
    await page.fillProjectName('reconnection-test');
    await page.selectProjectType('web-api');
    await page.selectFramework('gin');
    await page.fillModulePath('github.com/user/reconnection-test');
    
    await page.waitForRealTimeUpdate();
    const beforeReconnection = await page.getPreviewContent();
    
    // Simulate WebSocket disconnection and reconnection
    await page.page.route('ws://localhost:8080/ws', route => route.abort());
    await page.page.waitForTimeout(1000);
    await page.page.unroute('ws://localhost:8080/ws');
    
    // Wait for reconnection
    await page.waitForWebSocketConnection();
    await page.waitForRealTimeUpdate();
    
    const afterReconnection = await page.getPreviewContent();
    
    // Preview content should be maintained or restored
    expect(afterReconnection).toContain('reconnection-test');
    expect(afterReconnection).toContain('gin');
  });
});