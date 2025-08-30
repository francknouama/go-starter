import { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import type { UserEvent } from '@testing-library/user-event';

// Mock API responses for testing
export const mockBlueprints = [
  {
    id: 'cli-simple',
    name: 'CLI Simple',
    description: 'Simple command-line application',
    type: 'cli',
    complexity: 'simple',
    fileCount: 8,
    features: ['logging', 'testing'],
    dependencies: ['cobra']
  },
  {
    id: 'web-api-standard',
    name: 'Web API Standard',
    description: 'Standard REST API with middleware',
    type: 'web-api',
    complexity: 'standard',
    fileCount: 25,
    features: ['rest-api', 'middleware', 'database'],
    dependencies: ['gin', 'gorm']
  },
  {
    id: 'grpc-gateway',
    name: 'gRPC Gateway',
    description: 'Dual HTTP/gRPC API service',
    type: 'grpc',
    complexity: 'expert',
    fileCount: 45,
    features: ['grpc', 'gateway', 'protobuf'],
    dependencies: ['grpc-go', 'grpc-gateway']
  }
];

export const mockGenerationResponse = {
  id: 'gen-123',
  status: 'completed',
  fileCount: 25,
  downloadUrl: '/api/v1/download/gen-123',
  files: [
    { path: 'main.go', size: 1024 },
    { path: 'go.mod', size: 256 },
    { path: 'internal/handlers/health.go', size: 512 }
  ]
};

export const mockWebSocketMessages = {
  connected: {
    type: 'connected',
    data: { clientId: 'client-123', message: 'Connected to WebSocket' },
    timestamp: new Date().toISOString()
  },
  generationProgress: {
    type: 'generation_progress',
    data: { progress: 50, currentFile: 'internal/handlers/user.go', phase: 'Generating handlers' },
    timestamp: new Date().toISOString()
  },
  generationComplete: {
    type: 'generation_complete',
    data: { fileCount: 25, downloadUrl: '/api/v1/download/gen-123' },
    timestamp: new Date().toISOString()
  },
  generationError: {
    type: 'generation_error',
    data: { error: 'Template compilation failed', details: 'Missing variable: ProjectName' },
    timestamp: new Date().toISOString()
  }
};

// Mock fetch for API calls
export const createMockFetch = (responses: Record<string, any> = {}) => {
  return jest.fn((url: string, options?: RequestInit) => {
    const method = options?.method || 'GET';
    const key = `${method} ${url}`;
    
    if (responses[key]) {
      return Promise.resolve({
        ok: true,
        status: 200,
        json: () => Promise.resolve(responses[key]),
        text: () => Promise.resolve(JSON.stringify(responses[key]))
      });
    }
    
    // Default responses for common endpoints
    if (url.includes('/api/v1/blueprints') && method === 'GET') {
      return Promise.resolve({
        ok: true,
        status: 200,
        json: () => Promise.resolve({ blueprints: mockBlueprints })
      });
    }
    
    if (url.includes('/api/v1/generate') && method === 'POST') {
      return Promise.resolve({
        ok: true,
        status: 200,
        json: () => Promise.resolve(mockGenerationResponse)
      });
    }
    
    // Default 404 for unmatched endpoints
    return Promise.resolve({
      ok: false,
      status: 404,
      json: () => Promise.resolve({ error: 'Not found' })
    });
  });
};

// Mock WebSocket class for testing
export class MockWebSocket {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSING = 2;
  static CLOSED = 3;

  readyState = MockWebSocket.CONNECTING;
  url: string;
  onopen?: (event: Event) => void;
  onclose?: (event: CloseEvent) => void;
  onmessage?: (event: MessageEvent) => void;
  onerror?: (event: Event) => void;

  private eventListeners = new Map<string, Set<EventListener>>();
  private messageQueue: any[] = [];

  constructor(url: string) {
    this.url = url;
    
    // Simulate connection opening
    setTimeout(() => {
      this.readyState = MockWebSocket.OPEN;
      this._dispatchEvent(new Event('open'));
      
      // Send any queued messages
      this.messageQueue.forEach(message => {
        this._dispatchEvent(new MessageEvent('message', { data: JSON.stringify(message) }));
      });
      this.messageQueue = [];
    }, 0);
  }

  addEventListener(type: string, listener: EventListener) {
    if (!this.eventListeners.has(type)) {
      this.eventListeners.set(type, new Set());
    }
    this.eventListeners.get(type)!.add(listener);
  }

  removeEventListener(type: string, listener: EventListener) {
    const listeners = this.eventListeners.get(type);
    if (listeners) {
      listeners.delete(listener);
    }
  }

  send(data: string) {
    if (this.readyState !== MockWebSocket.OPEN) {
      throw new Error('WebSocket is not open');
    }
    // Echo back for testing
    const message = JSON.parse(data);
    setTimeout(() => {
      this.simulateMessage({ ...message, echo: true });
    }, 0);
  }

  close(code?: number, reason?: string) {
    this.readyState = MockWebSocket.CLOSED;
    setTimeout(() => {
      this._dispatchEvent(new CloseEvent('close', { code, reason }));
    }, 0);
  }

  // Test utilities
  simulateMessage(message: any) {
    if (this.readyState === MockWebSocket.OPEN) {
      this._dispatchEvent(new MessageEvent('message', { data: JSON.stringify(message) }));
    } else {
      this.messageQueue.push(message);
    }
  }

  simulateError(error?: string) {
    this.readyState = MockWebSocket.CLOSED;
    const errorEvent = new Event('error');
    (errorEvent as any).message = error || 'Connection error';
    this._dispatchEvent(errorEvent);
  }

  simulateClose(code = 1000, reason = 'Normal closure') {
    this.readyState = MockWebSocket.CLOSED;
    this._dispatchEvent(new CloseEvent('close', { code, reason }));
  }

  private _dispatchEvent(event: Event) {
    // Call direct event handlers
    if (event.type === 'open' && this.onopen) {
      this.onopen(event);
    } else if (event.type === 'close' && this.onclose) {
      this.onclose(event as CloseEvent);
    } else if (event.type === 'message' && this.onmessage) {
      this.onmessage(event as MessageEvent);
    } else if (event.type === 'error' && this.onerror) {
      this.onerror(event);
    }
    
    // Call addEventListener handlers
    const listeners = this.eventListeners.get(event.type);
    if (listeners) {
      listeners.forEach(listener => {
        try {
          listener(event);
        } catch (error) {
          console.error('Error in event listener:', error);
        }
      });
    }
  }
}

// Enhanced render function with common providers
interface CustomRenderOptions extends Omit<RenderOptions, 'queries'> {
  withUser?: boolean;
  mockWebSocket?: boolean;
  mockFetch?: Record<string, any>;
}

export function renderWithUser(
  ui: ReactElement,
  options: CustomRenderOptions = {}
): { user: UserEvent } & ReturnType<typeof render> {
  const { withUser = true, mockWebSocket = true, mockFetch, ...renderOptions } = options;
  
  // Setup mocks
  if (mockWebSocket) {
    global.WebSocket = MockWebSocket as any;
  }
  
  if (mockFetch) {
    global.fetch = createMockFetch(mockFetch) as any;
  }
  
  const user = userEvent.setup();
  const result = render(ui, renderOptions);
  
  return {
    user,
    ...result
  };
}

// Utility functions for common test scenarios
export const fillProjectForm = async (
  user: UserEvent,
  container: HTMLElement,
  config: {
    name?: string;
    type?: string;
    framework?: string;
    logger?: string;
    modulePath?: string;
  } = {}
) => {
  const {
    name = 'test-project',
    type = 'web-api',
    framework = 'gin',
    logger = 'slog',
    modulePath = 'github.com/user/test-project'
  } = config;

  // Fill project name
  const nameInput = container.querySelector('input[name="projectName"]') as HTMLInputElement;
  if (nameInput) {
    await user.clear(nameInput);
    await user.type(nameInput, name);
  }

  // Select project type
  const typeSelect = container.querySelector('select[name="projectType"]') as HTMLSelectElement;
  if (typeSelect) {
    await user.selectOptions(typeSelect, type);
  }

  // Select framework
  const frameworkSelect = container.querySelector('select[name="framework"]') as HTMLSelectElement;
  if (frameworkSelect) {
    await user.selectOptions(frameworkSelect, framework);
  }

  // Select logger
  const loggerSelect = container.querySelector('select[name="logger"]') as HTMLSelectElement;
  if (loggerSelect) {
    await user.selectOptions(loggerSelect, logger);
  }

  // Fill module path
  const modulePathInput = container.querySelector('input[name="modulePath"]') as HTMLInputElement;
  if (modulePathInput) {
    await user.clear(modulePathInput);
    await user.type(modulePathInput, modulePath);
  }
};

export const waitForWebSocketConnection = async (
  container: HTMLElement,
  timeout = 5000
): Promise<void> => {
  return new Promise((resolve, reject) => {
    const startTime = Date.now();
    const checkConnection = () => {
      const statusElement = container.querySelector('[data-testid="websocket-status"]');
      if (statusElement && statusElement.textContent === 'Connected') {
        resolve();
        return;
      }
      
      if (Date.now() - startTime > timeout) {
        reject(new Error('WebSocket connection timeout'));
        return;
      }
      
      setTimeout(checkConnection, 100);
    };
    
    checkConnection();
  });
};

export const simulateGenerationWorkflow = (
  mockWs: MockWebSocket,
  config: {
    progressSteps?: number;
    duration?: number;
    shouldFail?: boolean;
  } = {}
) => {
  const { progressSteps = 5, duration = 2000, shouldFail = false } = config;
  const stepDuration = duration / progressSteps;
  
  let currentStep = 0;
  
  const sendProgressUpdate = () => {
    if (shouldFail && currentStep === Math.floor(progressSteps / 2)) {
      mockWs.simulateMessage(mockWebSocketMessages.generationError);
      return;
    }
    
    if (currentStep < progressSteps) {
      const progress = (currentStep / progressSteps) * 100;
      mockWs.simulateMessage({
        ...mockWebSocketMessages.generationProgress,
        data: {
          progress,
          currentFile: `file-${currentStep}.go`,
          phase: `Step ${currentStep + 1}/${progressSteps}`
        }
      });
      
      currentStep++;
      setTimeout(sendProgressUpdate, stepDuration);
    } else {
      mockWs.simulateMessage(mockWebSocketMessages.generationComplete);
    }
  };
  
  // Start after a short delay
  setTimeout(sendProgressUpdate, 100);
};

// Test data generators
export const generateProjectConfig = (overrides: any = {}) => ({
  projectName: 'test-project',
  projectType: 'web-api',
  framework: 'gin',
  logger: 'slog',
  modulePath: 'github.com/user/test-project',
  goVersion: '1.21',
  architecture: 'standard',
  ...overrides
});

export const generateMockBlueprint = (overrides: any = {}) => ({
  id: 'test-blueprint',
  name: 'Test Blueprint',
  description: 'A test blueprint for testing',
  type: 'web-api',
  complexity: 'standard',
  fileCount: 20,
  features: ['rest-api', 'testing'],
  dependencies: ['gin', 'testify'],
  ...overrides
});

// Accessibility testing utilities
export const checkAccessibility = async (container: HTMLElement) => {
  // Check for common accessibility issues
  const issues: string[] = [];
  
  // Check for images without alt text
  const images = container.querySelectorAll('img');
  images.forEach((img, index) => {
    if (!img.getAttribute('alt')) {
      issues.push(`Image at index ${index} missing alt attribute`);
    }
  });
  
  // Check for buttons without accessible names
  const buttons = container.querySelectorAll('button');
  buttons.forEach((button, index) => {
    const hasAccessibleName = button.textContent?.trim() ||
                             button.getAttribute('aria-label') ||
                             button.getAttribute('aria-labelledby');
    if (!hasAccessibleName) {
      issues.push(`Button at index ${index} missing accessible name`);
    }
  });
  
  // Check for form inputs without labels
  const inputs = container.querySelectorAll('input, select, textarea');
  inputs.forEach((input, index) => {
    const id = input.getAttribute('id');
    const hasLabel = id && container.querySelector(`label[for="${id}"]`) ||
                    input.getAttribute('aria-label') ||
                    input.getAttribute('aria-labelledby');
    if (!hasLabel) {
      issues.push(`Form input at index ${index} missing label`);
    }
  });
  
  return issues;
};

// Performance testing utilities
export const measureRenderTime = async (renderFn: () => void): Promise<number> => {
  const start = performance.now();
  renderFn();
  const end = performance.now();
  return end - start;
};

export const createPerformanceObserver = (entryTypes: string[] = ['measure', 'navigation']) => {
  const entries: PerformanceEntry[] = [];
  
  if ('PerformanceObserver' in window) {
    const observer = new PerformanceObserver((list) => {
      entries.push(...list.getEntries());
    });
    
    observer.observe({ entryTypes });
    
    return {
      entries,
      disconnect: () => observer.disconnect()
    };
  }
  
  return {
    entries,
    disconnect: () => {}
  };
};

// Error boundary for testing error scenarios
export class TestErrorBoundary extends React.Component<
  { children: React.ReactNode; onError?: (error: Error) => void },
  { hasError: boolean; error?: Error }
> {
  constructor(props: any) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    this.props.onError?.(error);
  }

  render() {
    if (this.state.hasError) {
      return <div data-testid="error-boundary">Something went wrong: {this.state.error?.message}</div>;
    }

    return this.props.children;
  }
}

export { render, screen, waitFor } from '@testing-library/react';