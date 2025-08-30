import React from 'react';
import { renderWithUser, MockWebSocket, mockWebSocketMessages, simulateGenerationWorkflow, generateProjectConfig } from '../../../utils/__tests__/test-utils';
import GenerationWorkflow from '../GenerationWorkflow';
import { waitFor } from '@testing-library/react';

// Mock the API service
jest.mock('../../../services/api', () => ({
  generateProject: jest.fn(),
  downloadProject: jest.fn(),
  validateConfiguration: jest.fn(),
}));

// Mock the WebSocket hook
const mockWebSocketHook = {
  connectionState: { connected: true, connecting: false, reconnectAttempts: 0 },
  sendMessage: jest.fn(),
  connect: jest.fn(),
  disconnect: jest.fn(),
  subscribe: jest.fn(() => jest.fn()), // Return unsubscribe function
};

jest.mock('../../../hooks/useWebSocket', () => ({
  useWebSocket: () => mockWebSocketHook,
}));

describe('GenerationWorkflow Integration', () => {
  let mockWebSocket: MockWebSocket;
  
  beforeEach(() => {
    jest.clearAllMocks();
    mockWebSocket = new MockWebSocket('ws://localhost:8080/ws');
    global.WebSocket = MockWebSocket as any;
  });

  afterEach(() => {
    if (mockWebSocket) {
      mockWebSocket.close();
    }
  });

  describe('Complete Generation Workflow', () => {
    it('should complete full generation workflow successfully', async () => {
      const config = generateProjectConfig({
        projectName: 'integration-test',
        projectType: 'web-api',
        framework: 'gin'
      });

      const { user, getByTestId, getByRole, queryByTestId } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Should show initial configuration
      expect(getByTestId('project-config')).toBeInTheDocument();
      expect(getByTestId('project-name')).toHaveValue('integration-test');

      // Generate button should be enabled
      const generateButton = getByRole('button', { name: /generate project/i });
      expect(generateButton).toBeEnabled();

      // Start generation
      await user.click(generateButton);

      // Should show loading state
      expect(getByTestId('generation-status')).toHaveTextContent('Initializing generation...');
      expect(generateButton).toBeDisabled();

      // Should show progress bar
      expect(getByTestId('generation-progress-bar')).toBeInTheDocument();

      // Simulate WebSocket progress updates
      simulateGenerationWorkflow(mockWebSocket);

      // Wait for progress updates
      await waitFor(() => {
        const progressBar = getByTestId('generation-progress-bar');
        expect(progressBar).toHaveAttribute('aria-valuenow');
        const progress = parseInt(progressBar.getAttribute('aria-valuenow') || '0');
        expect(progress).toBeGreaterThan(0);
      });

      // Should show current file being generated
      await waitFor(() => {
        expect(getByTestId('current-file-indicator')).toBeInTheDocument();
      });

      // Wait for completion
      await waitFor(() => {
        expect(getByTestId('generation-status')).toHaveTextContent('Generation complete');
      }, { timeout: 10000 });

      // Progress bar should be hidden
      expect(queryByTestId('generation-progress-bar')).not.toBeInTheDocument();

      // Download button should be enabled
      const downloadButton = getByRole('button', { name: /download/i });
      expect(downloadButton).toBeEnabled();

      // Should show file count
      const fileCountElement = getByTestId('generated-files-count');
      expect(fileCountElement).toHaveTextContent(/\d+/);
      const fileCount = parseInt(fileCountElement.textContent || '0');
      expect(fileCount).toBeGreaterThan(0);
    });

    it('should handle generation errors gracefully', async () => {
      const config = generateProjectConfig();

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Start generation
      const generateButton = getByRole('button', { name: /generate project/i });
      await user.click(generateButton);

      // Simulate error during generation
      simulateGenerationWorkflow(mockWebSocket, { shouldFail: true });

      // Wait for error state
      await waitFor(() => {
        expect(getByTestId('generation-status')).toHaveTextContent(/error|failed/i);
      });

      // Should show error message
      expect(getByTestId('error-message')).toBeInTheDocument();
      expect(getByTestId('error-message')).toHaveTextContent(/template compilation failed/i);

      // Generate button should be enabled again
      expect(generateButton).toBeEnabled();

      // Download button should be disabled
      const downloadButton = getByRole('button', { name: /download/i });
      expect(downloadButton).toBeDisabled();

      // Should show retry option
      const retryButton = getByRole('button', { name: /try again|retry/i });
      expect(retryButton).toBeEnabled();
    });

    it('should support regeneration with different configurations', async () => {
      const initialConfig = generateProjectConfig({
        projectName: 'regen-test',
        projectType: 'cli'
      });

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={initialConfig} />
      );

      // Complete initial generation
      const generateButton = getByRole('button', { name: /generate project/i });
      await user.click(generateButton);

      simulateGenerationWorkflow(mockWebSocket);

      await waitFor(() => {
        expect(getByTestId('generation-status')).toHaveTextContent('Generation complete');
      });

      const initialFileCount = parseInt(getByTestId('generated-files-count').textContent || '0');

      // Change configuration
      const projectTypeSelect = getByTestId('project-type-select') as HTMLSelectElement;
      await user.selectOptions(projectTypeSelect, 'web-api');

      // Should show configuration changed indicator
      await waitFor(() => {
        expect(getByTestId('configuration-changed-indicator')).toBeInTheDocument();
      });

      // Regenerate button should be available
      const regenerateButton = getByRole('button', { name: /regenerate|generate new/i });
      expect(regenerateButton).toBeEnabled();

      // Start regeneration
      await user.click(regenerateButton);

      simulateGenerationWorkflow(mockWebSocket);

      await waitFor(() => {
        expect(getByTestId('generation-status')).toHaveTextContent('Generation complete');
      });

      // File count should be different
      const newFileCount = parseInt(getByTestId('generated-files-count').textContent || '0');
      expect(newFileCount).not.toBe(initialFileCount);
      expect(newFileCount).toBeGreaterThan(initialFileCount); // Web API has more files than CLI
    });

    it('should maintain WebSocket connection throughout workflow', async () => {
      const config = generateProjectConfig();

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Should show connected status
      expect(getByTestId('websocket-status')).toHaveTextContent('Connected');

      // Start generation
      await user.click(getByRole('button', { name: /generate project/i }));

      // Should maintain connection during generation
      expect(getByTestId('websocket-status')).toHaveTextContent('Connected');

      simulateGenerationWorkflow(mockWebSocket);

      // Wait for completion while checking connection
      await waitFor(() => {
        expect(getByTestId('generation-status')).toHaveTextContent('Generation complete');
        expect(getByTestId('websocket-status')).toHaveTextContent('Connected');
      });

      // Connection should still be active after completion
      expect(getByTestId('websocket-status')).toHaveTextContent('Connected');
    });

    it('should handle WebSocket disconnection during generation', async () => {
      const config = generateProjectConfig();

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Start generation
      await user.click(getByRole('button', { name: /generate project/i }));

      // Simulate connection loss
      mockWebSocket.simulateClose(1006, 'Connection lost');

      // Should show reconnecting status
      await waitFor(() => {
        expect(getByTestId('websocket-status')).toHaveTextContent(/reconnecting|disconnected/i);
      });

      // Should handle reconnection
      setTimeout(() => {
        mockWebSocket = new MockWebSocket('ws://localhost:8080/ws');
        mockWebSocket.simulateMessage(mockWebSocketMessages.connected);
      }, 1000);

      // Should reconnect and continue
      await waitFor(() => {
        expect(getByTestId('websocket-status')).toHaveTextContent('Connected');
      }, { timeout: 5000 });
    });
  });

  describe('Real-time Progress Updates', () => {
    it('should show detailed progress information', async () => {
      const config = generateProjectConfig();

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Start generation
      await user.click(getByRole('button', { name: /generate project/i }));

      // Should show progress bar
      const progressBar = getByTestId('generation-progress-bar');
      expect(progressBar).toBeInTheDocument();
      expect(progressBar).toHaveAttribute('aria-valuenow', '0');

      // Simulate detailed progress updates
      const progressUpdates = [
        { progress: 20, currentFile: 'main.go', phase: 'Generating main files' },
        { progress: 40, currentFile: 'handlers/user.go', phase: 'Generating handlers' },
        { progress: 60, currentFile: 'models/user.go', phase: 'Generating models' },
        { progress: 80, currentFile: 'tests/user_test.go', phase: 'Generating tests' }
      ];

      for (const update of progressUpdates) {
        mockWebSocket.simulateMessage({
          ...mockWebSocketMessages.generationProgress,
          data: update
        });

        await waitFor(() => {
          expect(progressBar).toHaveAttribute('aria-valuenow', update.progress.toString());
        });

        // Should show current file
        expect(getByTestId('current-file-indicator')).toHaveTextContent(update.currentFile);

        // Should show current phase
        expect(getByTestId('generation-phase')).toHaveTextContent(update.phase);
      }

      // Complete generation
      mockWebSocket.simulateMessage(mockWebSocketMessages.generationComplete);

      await waitFor(() => {
        expect(getByTestId('generation-status')).toHaveTextContent('Generation complete');
      });
    });

    it('should update file count in real-time', async () => {
      const config = generateProjectConfig();

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Start generation
      await user.click(getByRole('button', { name: /generate project/i }));

      // Should start with 0 files generated
      const fileCountElement = getByTestId('generated-files-count');
      expect(fileCountElement).toHaveTextContent('0');

      // Simulate file generation updates
      const fileUpdates = [5, 10, 15, 20, 25];

      for (const fileCount of fileUpdates) {
        mockWebSocket.simulateMessage({
          type: 'files_generated',
          data: { count: fileCount },
          timestamp: new Date().toISOString()
        });

        await waitFor(() => {
          expect(fileCountElement).toHaveTextContent(fileCount.toString());
        });
      }
    });

    it('should show estimated time remaining', async () => {
      const config = generateProjectConfig();

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Start generation
      await user.click(getByRole('button', { name: /generate project/i }));

      // Should show estimated time
      await waitFor(() => {
        expect(getByTestId('estimated-time-remaining')).toBeInTheDocument();
      });

      // Simulate progress with time estimates
      const progressWithTime = [
        { progress: 25, estimatedTimeRemaining: 120 }, // 2 minutes
        { progress: 50, estimatedTimeRemaining: 60 },  // 1 minute
        { progress: 75, estimatedTimeRemaining: 30 },  // 30 seconds
      ];

      for (const update of progressWithTime) {
        mockWebSocket.simulateMessage({
          ...mockWebSocketMessages.generationProgress,
          data: update
        });

        await waitFor(() => {
          const timeElement = getByTestId('estimated-time-remaining');
          expect(timeElement).toHaveTextContent(/\d+/);
        });
      }
    });
  });

  describe('User Interactions During Generation', () => {
    it('should prevent configuration changes during generation', async () => {
      const config = generateProjectConfig();

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Start generation
      await user.click(getByRole('button', { name: /generate project/i }));

      // Configuration inputs should be disabled
      const projectNameInput = getByTestId('project-name') as HTMLInputElement;
      const projectTypeSelect = getByTestId('project-type-select') as HTMLSelectElement;
      const frameworkSelect = getByTestId('framework-select') as HTMLSelectElement;

      expect(projectNameInput).toBeDisabled();
      expect(projectTypeSelect).toBeDisabled();
      expect(frameworkSelect).toBeDisabled();

      // Complete generation
      simulateGenerationWorkflow(mockWebSocket);

      await waitFor(() => {
        expect(getByTestId('generation-status')).toHaveTextContent('Generation complete');
      });

      // Inputs should be enabled again
      expect(projectNameInput).toBeEnabled();
      expect(projectTypeSelect).toBeEnabled();
      expect(frameworkSelect).toBeEnabled();
    });

    it('should allow canceling generation', async () => {
      const config = generateProjectConfig();

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Start generation
      await user.click(getByRole('button', { name: /generate project/i }));

      // Should show cancel button
      const cancelButton = getByRole('button', { name: /cancel|stop/i });
      expect(cancelButton).toBeEnabled();

      // Simulate partial progress
      mockWebSocket.simulateMessage({
        ...mockWebSocketMessages.generationProgress,
        data: { progress: 30, currentFile: 'partial.go', phase: 'Generating...' }
      });

      await waitFor(() => {
        const progressBar = getByTestId('generation-progress-bar');
        expect(progressBar).toHaveAttribute('aria-valuenow', '30');
      });

      // Cancel generation
      await user.click(cancelButton);

      // Should show canceled status
      await waitFor(() => {
        expect(getByTestId('generation-status')).toHaveTextContent(/canceled|stopped/i);
      });

      // Generate button should be enabled again
      const generateButton = getByRole('button', { name: /generate project/i });
      expect(generateButton).toBeEnabled();

      // Download button should be disabled
      const downloadButton = getByRole('button', { name: /download/i });
      expect(downloadButton).toBeDisabled();
    });

    it('should show helpful messages during long generation', async () => {
      const config = generateProjectConfig({
        projectType: 'web-api',
        architecture: 'ddd' // Complex architecture for longer generation
      });

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Start generation
      await user.click(getByRole('button', { name: /generate project/i }));

      // Should show helpful tips during generation
      await waitFor(() => {
        expect(getByTestId('generation-tips')).toBeInTheDocument();
      });

      const tipsElement = getByTestId('generation-tips');
      expect(tipsElement).toHaveTextContent(/tip|hint|while you wait/i);

      // Tips should change over time
      setTimeout(() => {
        mockWebSocket.simulateMessage({
          type: 'generation_tip',
          data: { tip: 'Complex architectures include dependency injection and clean separation of concerns.' },
          timestamp: new Date().toISOString()
        });
      }, 1000);

      await waitFor(() => {
        expect(tipsElement).toHaveTextContent(/dependency injection/i);
      });
    });
  });

  describe('Accessibility and UX', () => {
    it('should provide proper ARIA labels and roles', async () => {
      const config = generateProjectConfig();

      const { getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Progress bar should have proper ARIA attributes
      const progressBar = getByTestId('generation-progress-bar');
      expect(progressBar).toHaveAttribute('role', 'progressbar');
      expect(progressBar).toHaveAttribute('aria-valuenow');
      expect(progressBar).toHaveAttribute('aria-valuemin', '0');
      expect(progressBar).toHaveAttribute('aria-valuemax', '100');

      // Generation status should be announced to screen readers
      const statusElement = getByTestId('generation-status');
      expect(statusElement).toHaveAttribute('aria-live', 'polite');

      // Buttons should have descriptive labels
      const generateButton = getByRole('button', { name: /generate project/i });
      expect(generateButton).toHaveAccessibleName('Generate Project');
    });

    it('should announce progress updates to screen readers', async () => {
      const config = generateProjectConfig();

      const { user, getByTestId, getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // Start generation
      await user.click(getByRole('button', { name: /generate project/i }));

      // Progress announcements should have aria-live
      const progressAnnouncement = getByTestId('progress-announcement');
      expect(progressAnnouncement).toHaveAttribute('aria-live', 'polite');

      // Simulate progress update
      mockWebSocket.simulateMessage({
        ...mockWebSocketMessages.generationProgress,
        data: { progress: 50, currentFile: 'handlers/user.go', phase: 'Generating handlers' }
      });

      await waitFor(() => {
        expect(progressAnnouncement).toHaveTextContent('50% complete: Generating handlers');
      });
    });

    it('should provide keyboard navigation support', async () => {
      const config = generateProjectConfig();

      const { getByRole } = renderWithUser(
        <GenerationWorkflow initialConfig={config} />
      );

      // All interactive elements should be focusable
      const generateButton = getByRole('button', { name: /generate project/i });
      const downloadButton = getByRole('button', { name: /download/i });

      generateButton.focus();
      expect(generateButton).toHaveFocus();

      // Tab navigation should work
      const event = new KeyboardEvent('keydown', { key: 'Tab' });
      generateButton.dispatchEvent(event);
    });
  });
});