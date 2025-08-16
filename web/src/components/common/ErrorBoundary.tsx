/**
 * React Error Boundary Components with Enhanced Error Handling
 * Provides comprehensive error catching, reporting, and recovery mechanisms
 */

import React, { Component } from 'react';
import type { ErrorInfo, ReactNode } from 'react';
import type { AppError } from '../../utils/error-handling';
import { ErrorFactory, errorStore, RetryManager } from '../../utils/error-handling';
import { PerformanceMonitor } from '../../utils/performance-optimization';
import { accessibility } from '../../styles/design-tokens';

// Error boundary props and state
interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: React.ComponentType<ErrorFallbackProps>;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
  enableRetry?: boolean;
  maxRetries?: number;
  isolate?: boolean; // Whether to isolate this boundary from parent boundaries
  level?: 'page' | 'section' | 'component';
  name?: string;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: AppError | null;
  errorInfo: ErrorInfo | null;
  retryCount: number;
  isRetrying: boolean;
  errorId: string | null;
}

interface ErrorFallbackProps {
  error: AppError;
  errorInfo: ErrorInfo | null;
  retry: () => void;
  canRetry: boolean;
  retryCount: number;
  level: 'page' | 'section' | 'component';
}

// Main Error Boundary Component
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  private retryTimeoutId: NodeJS.Timeout | null = null;
  private performanceMonitor = PerformanceMonitor.getInstance();

  constructor(props: ErrorBoundaryProps) {
    super(props);
    
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      retryCount: 0,
      isRetrying: false,
      errorId: null
    };
  }

  static getDerivedStateFromError(error: Error): Partial<ErrorBoundaryState> {
    const appError = ErrorFactory.createSystemError(
      `React Error Boundary: ${error.message}`,
      error
    );
    
    // Add to global error store
    errorStore.addError(appError);

    return {
      hasError: true,
      error: appError,
      errorId: appError.id
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    const { onError, name, level = 'component' } = this.props;
    
    // Enhanced error logging
    console.group(`🚨 Error Boundary (${name || level})`);
    console.error('Error:', error);
    console.error('Error Info:', errorInfo);
    console.error('Component Stack:', errorInfo.componentStack);
    console.groupEnd();

    // Update state with error info
    this.setState({ errorInfo });

    // Measure error impact on performance
    this.performanceMonitor.measureComponentRender(
      `error-boundary-${name || level}`,
      () => {
        // Custom error handling
        onError?.(error, errorInfo);
        
        // Report to external services (in production)
        this.reportError(error, errorInfo);
      }
    );
  }

  private reportError(error: Error, errorInfo: ErrorInfo): void {
    // In a real app, you would send this to your error reporting service
    if (process.env.NODE_ENV === 'production') {
      // Example: Sentry, LogRocket, Bugsnag, etc.
      console.log('Would report to error service:', {
        error: error.message,
        stack: error.stack,
        componentStack: errorInfo.componentStack,
        timestamp: new Date().toISOString(),
        userAgent: navigator.userAgent,
        url: window.location.href
      });
    }
  }

  private handleRetry = async (): Promise<void> => {
    const { maxRetries = 3, name = 'unnamed' } = this.props;
    const { retryCount } = this.state;

    if (retryCount >= maxRetries) {
      console.warn(`Max retries (${maxRetries}) reached for ${name}`);
      return;
    }

    this.setState({ isRetrying: true });

    try {
      // Use retry manager for intelligent retry logic
      await RetryManager.withRetry(
        async () => {
          // Simulate component recovery
          await new Promise(resolve => setTimeout(resolve, 100));
          return true;
        },
        {
          maxRetries: 1,
          onRetry: (attempt) => {
            console.log(`Retrying error boundary ${name}, attempt ${attempt}`);
          }
        }
      );

      // Clear error state after successful retry
      this.retryTimeoutId = setTimeout(() => {
        this.setState({
          hasError: false,
          error: null,
          errorInfo: null,
          isRetrying: false,
          retryCount: retryCount + 1
        });
      }, 250);

    } catch (retryError) {
      console.error(`Retry failed for ${name}:`, retryError);
      this.setState({ 
        isRetrying: false,
        retryCount: retryCount + 1 
      });
    }
  };

  componentWillUnmount() {
    if (this.retryTimeoutId) {
      clearTimeout(this.retryTimeoutId);
    }
  }

  render() {
    const { 
      children, 
      fallback: FallbackComponent = DefaultErrorFallback,
      enableRetry = true,
      maxRetries = 3,
      level = 'component'
    } = this.props;
    
    const { hasError, error, errorInfo, retryCount, isRetrying } = this.state;

    if (hasError && error) {
      const canRetry = enableRetry && retryCount < maxRetries && !isRetrying;

      return (
        <FallbackComponent
          error={error}
          errorInfo={errorInfo}
          retry={this.handleRetry}
          canRetry={canRetry}
          retryCount={retryCount}
          level={level}
        />
      );
    }

    return children;
  }
}

// Default Error Fallback Component
const DefaultErrorFallback: React.FC<ErrorFallbackProps> = ({
  error,
  errorInfo,
  retry,
  canRetry,
  retryCount,
  level
}) => {
  const isHighContrast = accessibility.prefersHighContrast();
  const shouldReduceMotion = accessibility.shouldReduceMotion();

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'text-red-600 bg-red-50 border-red-200';
      case 'high':
        return 'text-orange-600 bg-orange-50 border-orange-200';
      case 'medium':
        return 'text-yellow-600 bg-yellow-50 border-yellow-200';
      default:
        return 'text-blue-600 bg-blue-50 border-blue-200';
    }
  };

  const getLevelIcon = () => {
    switch (level) {
      case 'page':
        return (
          <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        );
      case 'section':
        return (
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        );
      default:
        return (
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.268 18.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
        );
    }
  };

  const containerPadding = {
    page: 'p-8',
    section: 'p-6',
    component: 'p-4'
  }[level];

  const titleSize = {
    page: 'text-2xl',
    section: 'text-xl',
    component: 'text-lg'
  }[level];

  return (
    <div 
      className={`
        flex items-center justify-center min-h-[200px] w-full
        ${isHighContrast ? 'high-contrast' : ''}
        ${shouldReduceMotion ? 'motion-reduced' : ''}
      `}
      role="alert"
      aria-live="assertive"
    >
      <div className={`max-w-md w-full bg-white rounded-lg shadow-lg border ${containerPadding}`}>
        {/* Error Header */}
        <div className="flex items-start mb-4">
          <div className={`flex-shrink-0 ${getSeverityColor(error.severity)}`}>
            {getLevelIcon()}
          </div>
          <div className="ml-3 flex-1">
            <h3 className={`${titleSize} font-semibold text-gray-900`}>
              {level === 'page' ? 'Page Error' : 
               level === 'section' ? 'Section Error' : 
               'Component Error'}
            </h3>
            <p className="text-sm text-gray-600 mt-1">
              Severity: {error.severity} • ID: {error.id}
            </p>
          </div>
        </div>

        {/* Error Message */}
        <div className="mb-6">
          <p className="text-gray-800 mb-3">{error.userMessage}</p>
          
          {error.severity === 'critical' && (
            <div className="p-3 bg-red-50 border border-red-200 rounded-md">
              <p className="text-sm text-red-800">
                <strong>Critical Error:</strong> This error may affect the stability of the application. 
                Please refresh the page or contact support.
              </p>
            </div>
          )}

          {retryCount > 0 && (
            <div className="mt-3 p-2 bg-gray-50 border border-gray-200 rounded-md">
              <p className="text-xs text-gray-600">
                Previous retry attempts: {retryCount}
              </p>
            </div>
          )}
        </div>

        {/* Actions */}
        <div className="space-y-3">
          {/* Primary Actions */}
          <div className="flex space-x-3">
            {canRetry && (
              <button
                onClick={retry}
                className="
                  flex-1 bg-blue-600 text-white px-4 py-2 rounded-md 
                  hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
                  disabled:opacity-50 disabled:cursor-not-allowed
                  transition-colors duration-200
                "
                aria-describedby="retry-description"
              >
                Try Again
              </button>
            )}
            
            <button
              onClick={() => window.location.reload()}
              className="
                flex-1 bg-gray-600 text-white px-4 py-2 rounded-md 
                hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2
                transition-colors duration-200
              "
            >
              Reload {level === 'page' ? 'Page' : 'App'}
            </button>
          </div>

          {/* Secondary Actions */}
          {error.actions && error.actions.length > 0 && (
            <div className="pt-3 border-t border-gray-200">
              <h4 className="text-sm font-medium text-gray-900 mb-2">Additional Actions:</h4>
              <div className="space-y-2">
                {error.actions.map((action, index) => (
                  <button
                    key={index}
                    onClick={action.action}
                    className={`
                      w-full text-left px-3 py-2 text-sm rounded-md transition-colors duration-200
                      ${action.primary 
                        ? 'bg-blue-50 text-blue-700 hover:bg-blue-100 focus:ring-blue-500' 
                        : 'bg-gray-50 text-gray-700 hover:bg-gray-100 focus:ring-gray-500'
                      }
                      focus:outline-none focus:ring-2 focus:ring-offset-1
                    `}
                  >
                    {action.label}
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Hidden descriptions for screen readers */}
        <div className="sr-only">
          <div id="retry-description">
            {canRetry 
              ? `Retry the failed operation. ${3 - retryCount} attempts remaining.`
              : 'Retry is not available. Maximum retry attempts reached.'
            }
          </div>
        </div>

        {/* Debug Info (Development Only) */}
        {process.env.NODE_ENV === 'development' && errorInfo && (
          <details className="mt-6 text-xs">
            <summary className="cursor-pointer text-gray-500 hover:text-gray-700">
              Debug Information
            </summary>
            <div className="mt-2 p-3 bg-gray-50 border rounded-md">
              <div className="mb-2">
                <strong>Component Stack:</strong>
                <pre className="whitespace-pre-wrap text-gray-600 mt-1">
                  {errorInfo.componentStack}
                </pre>
              </div>
              {error.context && (
                <div>
                  <strong>Context:</strong>
                  <pre className="whitespace-pre-wrap text-gray-600 mt-1">
                    {JSON.stringify(error.context, null, 2)}
                  </pre>
                </div>
              )}
            </div>
          </details>
        )}
      </div>
    </div>
  );
};

// Specialized Error Boundaries

// Page-level error boundary
export const PageErrorBoundary: React.FC<{
  children: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}> = ({ children, onError }) => (
  <ErrorBoundary
    level="page"
    name="page-boundary"
    enableRetry={true}
    maxRetries={2}
    onError={onError}
  >
    {children}
  </ErrorBoundary>
);

// Section-level error boundary  
export const SectionErrorBoundary: React.FC<{
  children: ReactNode;
  sectionName?: string;
}> = ({ children, sectionName }) => (
  <ErrorBoundary
    level="section"
    name={sectionName}
    enableRetry={true}
    maxRetries={3}
  >
    {children}
  </ErrorBoundary>
);

// Component-level error boundary
export const ComponentErrorBoundary: React.FC<{
  children: ReactNode;
  componentName?: string;
  fallback?: React.ComponentType<ErrorFallbackProps>;
}> = ({ children, componentName, fallback }) => (
  <ErrorBoundary
    level="component"
    name={componentName}
    enableRetry={true}
    maxRetries={5}
    fallback={fallback}
  >
    {children}
  </ErrorBoundary>
);

// HOC for adding error boundaries to components
export function withErrorBoundary<P extends object>(
  Component: React.ComponentType<P>,
  boundaryProps?: Partial<ErrorBoundaryProps>
) {
  const WrappedComponent = React.forwardRef<any, P>((props, ref) => (
    <ErrorBoundary {...boundaryProps}>
      {/* @ts-expect-error - Complex generic type issue with forwardRef */}
      <Component {...props} ref={ref} />
    </ErrorBoundary>
  ));
  
  WrappedComponent.displayName = `withErrorBoundary(${Component.displayName || Component.name})`;
  
  return WrappedComponent;
}

// Hook for error reporting within components
export function useErrorReporting() {
  const reportError = React.useCallback((error: Error, context?: Record<string, unknown>) => {
    const appError = ErrorFactory.createSystemError(
      `Component Error: ${error.message}`,
      error,
      context
    );
    errorStore.addError(appError);
  }, []);

  const reportValidationError = React.useCallback((field: string, message: string) => {
    const error = ErrorFactory.createValidationError(field, message);
    errorStore.addError(error);
    return error;
  }, []);

  const reportNetworkError = React.useCallback((endpoint: string, status?: number, message?: string) => {
    const error = ErrorFactory.createNetworkError(endpoint, status, message);
    errorStore.addError(error);
    return error;
  }, []);

  return {
    reportError,
    reportValidationError,
    reportNetworkError
  };
}