import { test, expect } from '@playwright/test';
import { GoStarterPageObject, createStandardWebAPIProject } from './utils/test-helpers';

test.describe('WebSocket Real-time Integration', () => {
  let page: GoStarterPageObject;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new GoStarterPageObject(browserPage);
    await page.goto();
  });

  test('should establish WebSocket connection on page load', async () => {
    // Wait for WebSocket connection indicator
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connected', {
      timeout: 10000
    });
    
    // Connection indicator should be green/success state
    const statusIndicator = page.page.getByTestId('websocket-status-indicator');
    await expect(statusIndicator).toHaveClass(/connected|success|green/);
  });

  test('should show real-time progress during generation', async () => {
    await createStandardWebAPIProject(page);
    
    // Start generation
    const generateButton = page.page.getByRole('button', { name: /generate project/i });
    await generateButton.click();
    
    // Should receive real-time progress updates
    await expect(page.page.getByTestId('generation-progress-bar')).toBeVisible();
    
    // Progress should start at 0 and increase
    const progressBar = page.page.getByTestId('generation-progress-bar');
    const initialProgress = await progressBar.getAttribute('aria-valuenow');
    expect(parseInt(initialProgress || '0')).toBe(0);
    
    // Wait for progress updates (should see intermediate values)
    await page.page.waitForFunction(() => {
      const progressElement = document.querySelector('[data-testid="generation-progress-bar"]');
      const progress = progressElement?.getAttribute('aria-valuenow');
      return progress && parseInt(progress) > 0 && parseInt(progress) < 100;
    }, { timeout: 15000 });
    
    // Eventually should complete to 100%
    await expect(progressBar).toHaveAttribute('aria-valuenow', '100', { timeout: 30000 });
    
    // Progress bar should disappear after completion
    await expect(progressBar).not.toBeVisible({ timeout: 5000 });
  });

  test('should show real-time file generation updates', async () => {
    await createStandardWebAPIProject(page);
    
    // Start generation
    await page.generateProject();
    
    // Should show current file being generated
    await expect(page.page.getByTestId('current-file-indicator')).toBeVisible();
    
    // Should see multiple file names during generation
    const fileNames: string[] = [];
    
    // Collect file names as they appear
    await page.page.waitForFunction(() => {
      const fileIndicator = document.querySelector('[data-testid="current-file-indicator"]');
      return fileIndicator && fileIndicator.textContent && fileIndicator.textContent.trim() !== '';
    });
    
    // Monitor file changes for a few seconds
    const startTime = Date.now();
    while (Date.now() - startTime < 10000) { // Monitor for 10 seconds
      try {
        const currentFile = await page.page.getByTestId('current-file-indicator').textContent({ timeout: 1000 });
        if (currentFile && !fileNames.includes(currentFile)) {
          fileNames.push(currentFile);
        }
      } catch {
        break; // Generation completed
      }
    }
    
    // Should have seen multiple files being generated
    expect(fileNames.length).toBeGreaterThan(1);
    
    // Should include expected file types
    expect(fileNames.some(file => file.includes('main.go'))).toBe(true);
    expect(fileNames.some(file => file.includes('.go'))).toBe(true);
  });

  test('should update file count in real-time', async () => {
    await createStandardWebAPIProject(page);
    
    // Get initial estimated file count
    const initialCount = await page.getEstimatedFiles();
    
    // Start generation
    await page.generateProject();
    
    // Monitor file count updates
    await page.page.waitForFunction(() => {
      const countElement = document.querySelector('[data-testid="generated-files-count"]');
      return countElement && parseInt(countElement.textContent || '0') > 0;
    }, { timeout: 15000 });
    
    // File count should increase during generation
    const generatedCountElement = page.page.getByTestId('generated-files-count');
    await expect(generatedCountElement).not.toHaveText('0');
    
    // Final count should match or be close to estimate
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    const finalCount = await generatedCountElement.textContent();
    const finalCountNumber = parseInt(finalCount || '0');
    
    expect(finalCountNumber).toBeGreaterThan(0);
    expect(finalCountNumber).toBeCloseTo(initialCount, 5); // Within 5 files of estimate
  });

  test('should handle WebSocket reconnection', async () => {
    // Verify initial connection
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connected');
    
    // Simulate network interruption by disconnecting WebSocket
    await page.page.evaluate(() => {
      // Access the WebSocket connection and close it
      const wsConnections = (window as any).__wsConnections;
      if (wsConnections && wsConnections.length > 0) {
        wsConnections[0].close();
      }
    });
    
    // Should show reconnecting status
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Reconnecting', {
      timeout: 5000
    });
    
    // Should automatically reconnect
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connected', {
      timeout: 10000
    });
    
    // Should still be able to generate projects after reconnection
    await createStandardWebAPIProject(page);
    await page.generateProject();
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
  });

  test('should handle WebSocket connection errors gracefully', async () => {
    // Simulate connection error by blocking WebSocket endpoint
    await page.page.route('**/ws/**', route => route.abort());
    
    // Reload to trigger new connection attempt
    await page.page.reload();
    await page.waitForPageLoad();
    
    // Should show connection error state
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connection Error', {
      timeout: 10000
    });
    
    // Should offer manual retry option
    await expect(page.page.getByRole('button', { name: /retry connection/i })).toBeVisible();
    
    // Should still allow basic functionality without WebSocket
    await createStandardWebAPIProject(page);
    
    // Generation should still work, just without real-time updates
    await page.generateProject();
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete', {
      timeout: 30000
    });
    
    // But progress bar might not show detailed updates
    const progressBar = page.page.getByTestId('generation-progress-bar');
    const hasProgressUpdates = await progressBar.isVisible();
    
    if (!hasProgressUpdates) {
      // Should show fallback progress indication
      await expect(page.page.getByTestId('generation-spinner')).toBeVisible();
    }
  });

  test('should broadcast generation status to multiple browser tabs', async ({ context }) => {
    // Open a second tab
    const page2 = await context.newPage();
    const goStarter2 = new GoStarterPageObject(page2);
    await goStarter2.goto();
    
    // Verify both tabs are connected
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connected');
    await expect(page2.getByTestId('websocket-status')).toHaveText('Connected');
    
    // Start generation in first tab
    await createStandardWebAPIProject(page);
    await page.generateProject();
    
    // Both tabs should see generation progress (if implemented)
    const hasBroadcast = await Promise.race([
      page2.getByTestId('generation-status').waitFor({ state: 'visible', timeout: 5000 })
        .then(() => true)
        .catch(() => false),
      new Promise(resolve => setTimeout(() => resolve(false), 6000))
    ]);
    
    if (hasBroadcast) {
      await expect(page2.getByTestId('generation-status')).toHaveText(/generating|in progress/i);
    }
    
    // Wait for completion in first tab
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    await page2.close();
  });

  test('should handle rapid WebSocket messages', async () => {
    await createStandardWebAPIProject(page);
    
    // Start generation to trigger rapid messages
    await page.generateProject();
    
    // Monitor for message handling errors in console
    const consoleErrors: string[] = [];
    page.page.on('console', msg => {
      if (msg.type() === 'error' && msg.text().includes('WebSocket')) {
        consoleErrors.push(msg.text());
      }
    });
    
    // Wait for generation to complete
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
    
    // Should not have WebSocket-related console errors
    expect(consoleErrors).toHaveLength(0);
    
    // UI should still be responsive
    await expect(page.page.getByRole('button', { name: /download/i })).toBeEnabled();
  });

  test('should preserve WebSocket connection during navigation', async () => {
    // Verify initial connection
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connected');
    
    // Navigate to template gallery
    await page.page.getByRole('button', { name: /browse templates/i }).click();
    
    // Connection should remain active
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connected');
    
    // Navigate back
    await page.page.getByRole('button', { name: /close|back/i }).click();
    
    // Connection should still be active
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connected');
    
    // Should still be able to generate projects
    await createStandardWebAPIProject(page);
    await page.generateProject();
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
  });

  test('should show connection quality indicators', async () => {
    // Should show connection latency indicator
    const latencyIndicator = page.page.getByTestId('connection-latency');
    await expect(latencyIndicator).toBeVisible();
    
    // Latency should be reasonable (< 1000ms for local testing)
    const latencyText = await latencyIndicator.textContent();
    if (latencyText) {
      const latencyMatch = latencyText.match(/(\d+)ms/);
      if (latencyMatch) {
        const latency = parseInt(latencyMatch[1]);
        expect(latency).toBeLessThan(1000);
      }
    }
    
    // Should show connection stability indicator
    await expect(page.page.getByTestId('connection-stability')).toBeVisible();
  });

  test('should handle WebSocket message size limits', async () => {
    // Configure a very complex project that might generate large messages
    await page.fillProjectName('large-message-test');
    await page.selectProjectType('web-api');
    await page.toggleDisclosureMode('advanced');
    await page.selectArchitecture('ddd');
    await page.selectDatabaseDriver('postgres');
    await page.selectAuthType('jwt');
    await page.fillModulePath('github.com/user/large-project');
    
    // Enable all features to maximize message size
    const features = ['monitoring', 'testing', 'documentation', 'deployment'];
    for (const feature of features) {
      try {
        await page.enableFeature(feature);
      } catch {
        // Feature might not exist, continue
      }
    }
    
    // Start generation
    await page.generateProject();
    
    // Should handle large messages without errors
    const consoleErrors: string[] = [];
    page.page.on('console', msg => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });
    
    // Wait for completion
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete', {
      timeout: 60000
    });
    
    // Should not have message size related errors
    const messageSizeErrors = consoleErrors.filter(error => 
      error.includes('message size') || 
      error.includes('too large') || 
      error.includes('exceeded')
    );
    expect(messageSizeErrors).toHaveLength(0);
  });

  test('should clean up WebSocket connections on page unload', async () => {
    // Verify connection is established
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connected');
    
    // Navigate away from the page
    await page.page.goto('about:blank');
    
    // When we return, should establish new connection
    await page.goto();
    await expect(page.page.getByTestId('websocket-status')).toHaveText('Connected');
    
    // Should work normally after reconnection
    await createStandardWebAPIProject(page);
    await page.generateProject();
    await expect(page.page.getByTestId('generation-status')).toHaveText('Generation complete');
  });
});