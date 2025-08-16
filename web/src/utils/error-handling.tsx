/**
 * Comprehensive Error Handling System for go-starter Web UI
 * Provides centralized error management, user-friendly messages, and recovery mechanisms
 */

import React from 'react';
import { accessibility } from '../styles/design-tokens';

// Error types and severity levels
export type ErrorSeverity = 'low' | 'medium' | 'high' | 'critical';
export type ErrorCategory = 'validation' | 'network' | 'permission' | 'system' | 'user';

export interface AppError {
  id: string;
  code: string;
  message: string;
  userMessage: string;
  severity: ErrorSeverity;
  category: ErrorCategory;
  timestamp: Date;
  context?: Record<string, unknown>;
  stackTrace?: string;
  recoverable: boolean;
  retryable: boolean;
  actions?: ErrorAction[];
}

export interface ErrorAction {
  label: string;
  action: () => void | Promise<void>;
  primary?: boolean;
}

// Error boundary state
export interface ErrorBoundaryState {
  hasError: boolean;
  error: AppError | null;
  errorId: string | null;
  retryCount: number;
  isRecovering: boolean;
}

// Centralized error store
class ErrorStore {
  private errors: Map<string, AppError> = new Map();
  private listeners: Set<(errors: AppError[]) => void> = new Set();
  private maxErrors = 50; // Prevent memory leaks

  addError(error: AppError): void {
    this.errors.set(error.id, error);
    
    // Limit stored errors
    if (this.errors.size > this.maxErrors) {
      const oldestError = Array.from(this.errors.keys())[0];
      this.errors.delete(oldestError);
    }

    this.notifyListeners();
    this.logError(error);
  }

  removeError(id: string): void {
    this.errors.delete(id);
    this.notifyListeners();
  }

  getErrors(): AppError[] {
    return Array.from(this.errors.values()).sort(
      (a, b) => b.timestamp.getTime() - a.timestamp.getTime()
    );
  }

  getError(id: string): AppError | undefined {
    return this.errors.get(id);
  }

  clearErrors(): void {
    this.errors.clear();
    this.notifyListeners();
  }

  subscribe(listener: (errors: AppError[]) => void): () => void {
    this.listeners.add(listener);
    return () => this.listeners.delete(listener);
  }

  private notifyListeners(): void {
    const errors = this.getErrors();
    this.listeners.forEach(listener => listener(errors));
  }

  private logError(error: AppError): void {
    const logLevel = {
      low: 'info',
      medium: 'warn', 
      high: 'error',
      critical: 'error'
    }[error.severity] as 'info' | 'warn' | 'error';

    console[logLevel]('[ErrorStore]', {
      id: error.id,
      code: error.code,
      message: error.message,
      severity: error.severity,
      category: error.category,
      context: error.context,
      stackTrace: error.stackTrace
    });
  }
}

export const errorStore = new ErrorStore();

// Error factory functions
export class ErrorFactory {
  static createValidationError(
    field: string,
    message: string,
    context?: Record<string, unknown>
  ): AppError {
    return {
      id: this.generateId(),
      code: `VALIDATION_${field.toUpperCase()}`,
      message: `Validation failed for ${field}: ${message}`,
      userMessage: message,
      severity: 'medium',
      category: 'validation',
      timestamp: new Date(),
      context: { field, ...context },
      recoverable: true,
      retryable: false,
      actions: [
        {
          label: 'Fix and retry',
          action: () => {
            // Focus the problematic field
            const element = document.getElementById(field);
            if (element) {
              element.focus();
              element.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
          },
          primary: true
        }
      ]
    };
  }

  static createNetworkError(
    endpoint: string,
    status?: number,
    message?: string,
    context?: Record<string, unknown>
  ): AppError {
    const isRetryable = !status || status >= 500 || status === 408 || status === 429;
    
    return {
      id: this.generateId(),
      code: `NETWORK_${status || 'UNKNOWN'}`,
      message: `Network request failed: ${endpoint} (${status || 'unknown'})`,
      userMessage: this.getNetworkErrorMessage(status, message),
      severity: status && status < 500 ? 'medium' : 'high',
      category: 'network',
      timestamp: new Date(),
      context: { endpoint, status, ...context },
      recoverable: true,
      retryable: isRetryable,
      actions: this.getNetworkErrorActions(endpoint, isRetryable)
    };
  }

  static createPermissionError(
    resource: string,
    action: string,
    context?: Record<string, unknown>
  ): AppError {
    return {
      id: this.generateId(),
      code: 'PERMISSION_DENIED',
      message: `Permission denied: ${action} on ${resource}`,
      userMessage: `You don't have permission to ${action} ${resource}. Please contact your administrator.`,
      severity: 'high',
      category: 'permission',
      timestamp: new Date(),
      context: { resource, action, ...context },
      recoverable: false,
      retryable: false,
      actions: [
        {
          label: 'Contact support',
          action: () => {
            // Could open support modal or redirect to help
            console.log('Contact support for permission issue');
          }
        }
      ]
    };
  }

  static createSystemError(
    message: string,
    error?: Error,
    context?: Record<string, unknown>
  ): AppError {
    return {
      id: this.generateId(),
      code: 'SYSTEM_ERROR',
      message: `System error: ${message}`,
      userMessage: 'An unexpected error occurred. Please try again or contact support if the problem persists.',
      severity: 'critical',
      category: 'system',
      timestamp: new Date(),
      context,
      stackTrace: error?.stack,
      recoverable: true,
      retryable: true,
      actions: [
        {
          label: 'Reload page',
          action: () => window.location.reload(),
          primary: true
        },
        {
          label: 'Report issue',
          action: () => {
            // Could open bug report modal
            console.error('User reported system error:', { message, error, context });
          }
        }
      ]
    };
  }

  static createUserError(
    message: string,
    suggestion: string,
    context?: Record<string, unknown>
  ): AppError {
    return {
      id: this.generateId(),
      code: 'USER_ERROR',
      message: `User error: ${message}`,
      userMessage: `${message} ${suggestion}`,
      severity: 'low',
      category: 'user',
      timestamp: new Date(),
      context,
      recoverable: true,
      retryable: false
    };
  }

  private static generateId(): string {
    return `error_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  private static getNetworkErrorMessage(status?: number, message?: string): string {
    if (message) return message;

    switch (status) {
      case 400:
        return 'Invalid request. Please check your input and try again.';
      case 401:
        return 'You need to sign in to continue.';
      case 403:
        return 'You don\'t have permission to perform this action.';
      case 404:
        return 'The requested resource was not found.';
      case 408:
        return 'Request timed out. Please try again.';
      case 429:
        return 'Too many requests. Please wait a moment and try again.';
      case 500:
        return 'Server error. Please try again later.';
      case 502:
      case 503:
        return 'Service temporarily unavailable. Please try again later.';
      default:
        return 'Network error. Please check your connection and try again.';
    }
  }

  private static getNetworkErrorActions(endpoint: string, isRetryable: boolean): ErrorAction[] {
    const actions: ErrorAction[] = [];

    if (isRetryable) {
      actions.push({
        label: 'Retry',
        action: async () => {
          // This would trigger a retry of the original request
          console.log('Retrying request to:', endpoint);
        },
        primary: true
      });
    }

    actions.push({
      label: 'Refresh page',
      action: () => window.location.reload()
    });

    return actions;
  }
}

// Error handler hook
export function useErrorHandler() {
  const [errors, setErrors] = React.useState<AppError[]>([]);

  React.useEffect(() => {
    const unsubscribe = errorStore.subscribe(setErrors);
    setErrors(errorStore.getErrors());
    return unsubscribe;
  }, []);

  const addError = React.useCallback((error: AppError) => {
    errorStore.addError(error);
  }, []);

  const removeError = React.useCallback((id: string) => {
    errorStore.removeError(id);
  }, []);

  const clearErrors = React.useCallback(() => {
    errorStore.clearErrors();
  }, []);

  const handleValidationError = React.useCallback((field: string, message: string) => {
    const error = ErrorFactory.createValidationError(field, message);
    addError(error);
    return error;
  }, [addError]);

  const handleNetworkError = React.useCallback((endpoint: string, status?: number, message?: string) => {
    const error = ErrorFactory.createNetworkError(endpoint, status, message);
    addError(error);
    return error;
  }, [addError]);

  const handleSystemError = React.useCallback((message: string, error?: Error) => {
    const appError = ErrorFactory.createSystemError(message, error);
    addError(appError);
    return appError;
  }, [addError]);

  return {
    errors,
    addError,
    removeError,
    clearErrors,
    handleValidationError,
    handleNetworkError,
    handleSystemError
  };
}

// Global error boundary
export class GlobalErrorBoundary extends React.Component<
  { children: React.ReactNode; fallback?: React.ComponentType<{ error: AppError; retry: () => void }> },
  ErrorBoundaryState
> {
  private retryTimeoutId: NodeJS.Timeout | null = null;

  constructor(props: any) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorId: null,
      retryCount: 0,
      isRecovering: false
    };
  }

  static getDerivedStateFromError(error: Error): Partial<ErrorBoundaryState> {
    const appError = ErrorFactory.createSystemError(
      'Component error boundary triggered',
      error
    );
    
    errorStore.addError(appError);

    return {
      hasError: true,
      error: appError,
      errorId: appError.id
    };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Error boundary caught error:', error, errorInfo);
  }

  handleRetry = () => {
    this.setState(prevState => ({
      isRecovering: true,
      retryCount: prevState.retryCount + 1
    }));

    // Clear the error after a brief delay to trigger re-render
    this.retryTimeoutId = setTimeout(() => {
      this.setState({
        hasError: false,
        error: null,
        errorId: null,
        isRecovering: false
      });
    }, 100);
  };

  componentWillUnmount() {
    if (this.retryTimeoutId) {
      clearTimeout(this.retryTimeoutId);
    }
  }

  render() {
    if (this.state.hasError && this.state.error) {
      const FallbackComponent = this.props.fallback || DefaultErrorFallback;
      return <FallbackComponent error={this.state.error} retry={this.handleRetry} />;
    }

    return this.props.children;
  }
}

// Default error fallback component
interface ErrorFallbackProps {
  error: AppError;
  retry: () => void;
}

const DefaultErrorFallback: React.FC<ErrorFallbackProps> = ({ error, retry }) => {
  const isHighContrast = accessibility.prefersHighContrast();

  return (
    <div className={`min-h-screen flex items-center justify-center bg-gray-50 ${isHighContrast ? 'high-contrast' : ''}`}>
      <div className="max-w-md w-full bg-white shadow-lg rounded-lg p-6">
        <div className="flex items-center mb-4">
          <div className="flex-shrink-0">
            <svg className="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.268 18.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <div className="ml-3">
            <h3 className="text-lg font-medium text-gray-900">
              Something went wrong
            </h3>
            <p className="text-sm text-gray-500">
              Error ID: {error.id}
            </p>
          </div>
        </div>

        <div className="mb-6">
          <p className="text-gray-700">{error.userMessage}</p>
          {error.severity === 'critical' && (
            <div className="mt-3 p-3 bg-red-50 border border-red-200 rounded-md">
              <p className="text-sm text-red-800">
                This is a critical error. Please contact support if the problem persists.
              </p>
            </div>
          )}
        </div>

        <div className="flex space-x-3">
          {error.recoverable && (
            <button
              onClick={retry}
              className="flex-1 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            >
              Try Again
            </button>
          )}
          
          <button
            onClick={() => window.location.reload()}
            className="flex-1 bg-gray-600 text-white px-4 py-2 rounded-md hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
          >
            Reload Page
          </button>
        </div>

        {error.actions && error.actions.length > 0 && (
          <div className="mt-4">
            <h4 className="text-sm font-medium text-gray-900 mb-2">Additional Actions:</h4>
            <div className="space-y-2">
              {error.actions.map((action, index) => (
                <button
                  key={index}
                  onClick={action.action}
                  className={`w-full text-left px-3 py-2 text-sm rounded-md ${
                    action.primary 
                      ? 'bg-blue-50 text-blue-700 hover:bg-blue-100' 
                      : 'bg-gray-50 text-gray-700 hover:bg-gray-100'
                  } focus:outline-none focus:ring-2 focus:ring-blue-500`}
                >
                  {action.label}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

// Retry mechanisms
export class RetryManager {
  private static retryAttempts = new Map<string, number>();
  private static readonly MAX_RETRIES = 3;
  private static readonly RETRY_DELAYS = [1000, 2000, 4000]; // Exponential backoff

  static async withRetry<T>(
    operation: () => Promise<T>,
    options: {
      maxRetries?: number;
      delays?: number[];
      shouldRetry?: (error: any) => boolean;
      onRetry?: (attempt: number, error: any) => void;
    } = {}
  ): Promise<T> {
    const {
      maxRetries = this.MAX_RETRIES,
      delays = this.RETRY_DELAYS,
      shouldRetry = (error) => !error.response || error.response.status >= 500,
      onRetry
    } = options;

    let lastError: any;

    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        const result = await operation();
        return result;
      } catch (error) {
        lastError = error;

        if (attempt === maxRetries || !shouldRetry(error)) {
          throw error;
        }

        const delay = delays[attempt] || delays[delays.length - 1];
        onRetry?.(attempt + 1, error);

        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }

    throw lastError;
  }

  static resetRetryCount(key: string): void {
    this.retryAttempts.delete(key);
  }

  static getRetryCount(key: string): number {
    return this.retryAttempts.get(key) || 0;
  }

  static incrementRetryCount(key: string): number {
    const current = this.getRetryCount(key);
    const newCount = current + 1;
    this.retryAttempts.set(key, newCount);
    return newCount;
  }
}

// Performance error detection
export class PerformanceErrorDetector {
  private static readonly SLOW_REQUEST_THRESHOLD = 5000; // 5 seconds
  private static readonly MEMORY_THRESHOLD = 100 * 1024 * 1024; // 100MB

  static monitorNetworkRequests(): void {
    // Monitor fetch requests
    const originalFetch = window.fetch;
    window.fetch = async function(...args) {
      const startTime = performance.now();
      
      try {
        const response = await originalFetch.apply(this, args);
        const duration = performance.now() - startTime;
        
        if (duration > PerformanceErrorDetector.SLOW_REQUEST_THRESHOLD) {
          const url = args[0] instanceof Request ? args[0].url : args[0];
          errorStore.addError(ErrorFactory.createSystemError(
            `Slow network request detected: ${url} took ${Math.round(duration)}ms`,
            undefined,
            { url, duration, threshold: PerformanceErrorDetector.SLOW_REQUEST_THRESHOLD }
          ));
        }
        
        return response;
      } catch (error) {
        const duration = performance.now() - startTime;
        const url = args[0] instanceof Request ? args[0].url : args[0];
        
        errorStore.addError(ErrorFactory.createNetworkError(
          String(url),
          undefined,
          error instanceof Error ? error.message : 'Network request failed',
          { duration }
        ));
        
        throw error;
      }
    };
  }

  static monitorMemoryUsage(): void {
    if ('memory' in performance) {
      setInterval(() => {
        const memory = (performance as any).memory;
        if (memory.usedJSHeapSize > this.MEMORY_THRESHOLD) {
          errorStore.addError(ErrorFactory.createSystemError(
            `High memory usage detected: ${Math.round(memory.usedJSHeapSize / 1024 / 1024)}MB`,
            undefined,
            { memoryUsage: memory }
          ));
        }
      }, 30000); // Check every 30 seconds
    }
  }

  static initialize(): void {
    this.monitorNetworkRequests();
    this.monitorMemoryUsage();
  }
}

// Initialize performance monitoring
if (typeof window !== 'undefined') {
  PerformanceErrorDetector.initialize();
}