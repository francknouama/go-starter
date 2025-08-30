/**
 * Simple Error Boundary Component
 * Provides basic error catching without complex dependencies
 */

import React, { Component, ErrorInfo, ReactNode } from 'react';

interface Props {
  children: ReactNode;
  level?: 'page' | 'section' | 'component';
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
  enableRetry?: boolean;
}

interface State {
  hasError: boolean;
  error: Error | null;
  retryCount: number;
}

export class SimpleErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null, retryCount: 0 };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error, retryCount: 0 };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Error Boundary caught an error:', error, errorInfo);
    this.props.onError?.(error, errorInfo);
  }

  handleRetry = () => {
    this.setState(prevState => ({
      hasError: false,
      error: null,
      retryCount: prevState.retryCount + 1
    }));
  };

  render() {
    if (this.state.hasError) {
      const { level = 'component', enableRetry = true } = this.props;
      
      return (
        <div className="flex items-center justify-center min-h-[200px] w-full">
          <div className="max-w-md w-full bg-white rounded-lg shadow-lg border p-6">
            <div className="flex items-start mb-4">
              <div className="flex-shrink-0 text-red-500">
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.268 18.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
              </div>
              <div className="ml-3">
                <h3 className="text-lg font-semibold text-gray-900">
                  {level === 'page' ? 'Page Error' : 
                   level === 'section' ? 'Section Error' : 
                   'Component Error'}
                </h3>
                <p className="text-sm text-gray-600 mt-1">
                  Something went wrong while rendering this {level}.
                </p>
              </div>
            </div>

            <div className="mb-6">
              <p className="text-gray-800 mb-2">
                An unexpected error occurred. Please try refreshing or retry the operation.
              </p>
              {process.env.NODE_ENV === 'development' && (
                <details className="mt-3">
                  <summary className="text-sm text-gray-500 cursor-pointer">Error Details</summary>
                  <pre className="text-xs text-red-600 mt-2 p-2 bg-red-50 rounded">
                    {this.state.error?.message}
                  </pre>
                </details>
              )}
            </div>

            <div className="flex space-x-3">
              {enableRetry && (
                <button
                  onClick={this.handleRetry}
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
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

// Export as ErrorBoundary for compatibility
export { SimpleErrorBoundary as ErrorBoundary };