// Jest setup file for testing environment configuration
import '@testing-library/jest-dom';

// Mock ResizeObserver which is not available in jsdom
global.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

// Mock IntersectionObserver which is not available in jsdom
global.IntersectionObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}));

// Mock matchMedia which is not available in jsdom
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(), // deprecated
    removeListener: jest.fn(), // deprecated
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

// Mock WebSocket for testing WebSocket functionality
global.WebSocket = jest.fn().mockImplementation(() => ({
  close: jest.fn(),
  send: jest.fn(),
  readyState: WebSocket.CONNECTING,
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
})) as any;

// Mock localStorage
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
};
global.localStorage = localStorageMock as any;

// Mock sessionStorage
const sessionStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
};
global.sessionStorage = sessionStorageMock as any;

// Mock URL.createObjectURL for file handling tests
global.URL.createObjectURL = jest.fn(() => 'mock-object-url');
global.URL.revokeObjectURL = jest.fn();

// Mock fetch for API testing
global.fetch = jest.fn();

// Setup custom Jest matchers for better assertions
expect.extend({
  toHaveAccessibleName(received: HTMLElement, expectedName: string) {
    const accessibleName = received.getAttribute('aria-label') || 
                          received.getAttribute('aria-labelledby') ||
                          received.textContent;
    
    const pass = accessibleName === expectedName;
    
    return {
      pass,
      message: () => 
        pass 
          ? `Expected element not to have accessible name "${expectedName}"`
          : `Expected element to have accessible name "${expectedName}", but got "${accessibleName}"`,
    };
  },
  
  toBeVisuallyHidden(received: HTMLElement) {
    const style = window.getComputedStyle(received);
    const isHidden = style.display === 'none' || 
                    style.visibility === 'hidden' ||
                    style.opacity === '0' ||
                    received.hasAttribute('hidden');
    
    return {
      pass: isHidden,
      message: () => 
        isHidden 
          ? `Expected element not to be visually hidden`
          : `Expected element to be visually hidden`,
    };
  },
});

// Extend Jest matchers interface
declare global {
  namespace jest {
    interface Matchers<R> {
      toHaveAccessibleName(expectedName: string): R;
      toBeVisuallyHidden(): R;
    }
  }
}

// Silence console errors in tests unless explicitly testing them
const originalError = console.error;
beforeAll(() => {
  console.error = (...args: any[]) => {
    if (
      typeof args[0] === 'string' &&
      args[0].includes('Warning: ReactDOM.render is deprecated')
    ) {
      return;
    }
    originalError.call(console, ...args);
  };
});

afterAll(() => {
  console.error = originalError;
});

// Clean up after each test
afterEach(() => {
  jest.clearAllMocks();
  localStorage.clear();
  sessionStorage.clear();
});