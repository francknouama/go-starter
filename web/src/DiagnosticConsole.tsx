import React, { useEffect, useState } from 'react';

interface ConsoleError {
  message: string;
  stack?: string;
  timestamp: number;
  type: 'error' | 'warning' | 'log';
}

export function DiagnosticConsole() {
  const [errors, setErrors] = useState<ConsoleError[]>([]);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    // Capture console errors
    const originalError = console.error;
    const originalWarn = console.warn;
    const originalLog = console.log;

    console.error = (...args) => {
      originalError(...args);
      setErrors(prev => [...prev, {
        message: args.join(' '),
        timestamp: Date.now(),
        type: 'error'
      }]);
    };

    console.warn = (...args) => {
      originalWarn(...args);
      setErrors(prev => [...prev, {
        message: args.join(' '),
        timestamp: Date.now(),
        type: 'warning'
      }]);
    };

    // Capture uncaught errors
    const handleError = (event: ErrorEvent) => {
      setErrors(prev => [...prev, {
        message: event.message,
        stack: event.error?.stack,
        timestamp: Date.now(),
        type: 'error'
      }]);
    };

    const handleRejection = (event: PromiseRejectionEvent) => {
      setErrors(prev => [...prev, {
        message: `Unhandled Promise Rejection: ${event.reason}`,
        timestamp: Date.now(),
        type: 'error'
      }]);
    };

    window.addEventListener('error', handleError);
    window.addEventListener('unhandledrejection', handleRejection);

    // Log initial diagnostic info
    console.log('🔍 Diagnostic Console initialized');
    console.log('📍 Current URL:', window.location.href);
    console.log('🌐 User Agent:', navigator.userAgent);
    console.log('⚛️ React Version:', React.version);

    return () => {
      console.error = originalError;
      console.warn = originalWarn;
      console.log = originalLog;
      window.removeEventListener('error', handleError);
      window.removeEventListener('unhandledrejection', handleRejection);
    };
  }, []);

  if (!isVisible && errors.length === 0) {
    return (
      <button
        onClick={() => setIsVisible(true)}
        style={{
          position: 'fixed',
          bottom: '20px',
          right: '20px',
          zIndex: 9999,
          background: '#3b82f6',
          color: 'white',
          padding: '8px 16px',
          borderRadius: '6px',
          border: 'none',
          cursor: 'pointer',
          fontSize: '12px'
        }}
      >
        🔍 Show Console
      </button>
    );
  }

  return (
    <div style={{
      position: 'fixed',
      bottom: isVisible ? '20px' : '-400px',
      right: '20px',
      width: '400px',
      maxHeight: '300px',
      background: 'rgba(0, 0, 0, 0.9)',
      color: 'white',
      padding: '16px',
      borderRadius: '8px',
      fontSize: '12px',
      fontFamily: 'monospace',
      zIndex: 9999,
      transition: 'bottom 0.3s ease',
      overflow: 'auto',
      border: errors.some(e => e.type === 'error') ? '2px solid #ef4444' : '1px solid #6b7280'
    }}>
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: '12px'
      }}>
        <h3 style={{ margin: 0, fontSize: '14px' }}>
          🔍 Console Diagnostic ({errors.length} {errors.length === 1 ? 'entry' : 'entries'})
        </h3>
        <div>
          <button
            onClick={() => setErrors([])}
            style={{
              background: '#6b7280',
              color: 'white',
              border: 'none',
              padding: '4px 8px',
              borderRadius: '4px',
              cursor: 'pointer',
              marginRight: '8px',
              fontSize: '11px'
            }}
          >
            Clear
          </button>
          <button
            onClick={() => setIsVisible(false)}
            style={{
              background: '#ef4444',
              color: 'white',
              border: 'none',
              padding: '4px 8px',
              borderRadius: '4px',
              cursor: 'pointer',
              fontSize: '11px'
            }}
          >
            ✕
          </button>
        </div>
      </div>

      {errors.length === 0 ? (
        <div style={{ color: '#22c55e', textAlign: 'center', padding: '20px' }}>
          ✅ No console errors detected!
          <br />
          <small>Application appears to be running cleanly.</small>
        </div>
      ) : (
        <div style={{ maxHeight: '200px', overflow: 'auto' }}>
          {errors.slice(-10).map((error, index) => (
            <div
              key={index}
              style={{
                marginBottom: '8px',
                padding: '8px',
                background: error.type === 'error' ? 'rgba(239, 68, 68, 0.2)' : 
                          error.type === 'warning' ? 'rgba(245, 158, 11, 0.2)' : 
                          'rgba(107, 114, 128, 0.2)',
                borderRadius: '4px',
                borderLeft: `3px solid ${error.type === 'error' ? '#ef4444' : 
                                       error.type === 'warning' ? '#f59e0b' : '#6b7280'}`
              }}
            >
              <div style={{ 
                display: 'flex', 
                justifyContent: 'space-between',
                alignItems: 'center',
                marginBottom: '4px'
              }}>
                <span style={{ 
                  fontSize: '10px', 
                  color: error.type === 'error' ? '#ef4444' : 
                        error.type === 'warning' ? '#f59e0b' : '#6b7280',
                  fontWeight: 'bold'
                }}>
                  {error.type.toUpperCase()}
                </span>
                <span style={{ fontSize: '10px', color: '#9ca3af' }}>
                  {new Date(error.timestamp).toLocaleTimeString()}
                </span>
              </div>
              <div style={{ color: 'white', lineHeight: '1.4' }}>
                {error.message}
              </div>
              {error.stack && (
                <details style={{ marginTop: '4px' }}>
                  <summary style={{ 
                    cursor: 'pointer', 
                    fontSize: '10px', 
                    color: '#9ca3af' 
                  }}>
                    Stack trace
                  </summary>
                  <pre style={{ 
                    fontSize: '10px', 
                    margin: '4px 0 0 0', 
                    color: '#d1d5db',
                    whiteSpace: 'pre-wrap'
                  }}>
                    {error.stack}
                  </pre>
                </details>
              )}
            </div>
          ))}
        </div>
      )}

      <div style={{ 
        marginTop: '12px', 
        padding: '8px', 
        background: 'rgba(59, 130, 246, 0.2)', 
        borderRadius: '4px',
        fontSize: '10px',
        color: '#93c5fd'
      }}>
        💡 This diagnostic tool captures JavaScript console output in real-time.
        {errors.some(e => e.type === 'error') && (
          <><br />🚨 <strong>Errors detected!</strong> Check the entries above for details.</>
        )}
      </div>
    </div>
  );
}

// Diagnostic status component
export function DiagnosticStatus() {
  const [status, setStatus] = useState<string>('Checking...');

  useEffect(() => {
    const checkStatus = () => {
      try {
        // Basic React functionality test
        const reactWorking = typeof React !== 'undefined';
        const domWorking = typeof document !== 'undefined';
        const windowWorking = typeof window !== 'undefined';
        
        console.log('📊 System Check:');
        console.log('⚛️ React:', reactWorking ? '✅' : '❌');
        console.log('🌐 DOM:', domWorking ? '✅' : '❌'); 
        console.log('🪟 Window:', windowWorking ? '✅' : '❌');
        
        if (reactWorking && domWorking && windowWorking) {
          setStatus('✅ System operational');
          console.log('✅ All systems operational');
        } else {
          setStatus('❌ System issues detected');
          console.error('❌ System issues detected');
        }
      } catch (error) {
        setStatus('❌ Critical error');
        console.error('❌ Diagnostic check failed:', error);
      }
    };

    checkStatus();
    
    // Log component mount
    console.log('🔍 DiagnosticStatus mounted');
    
    return () => {
      console.log('🔍 DiagnosticStatus unmounted');
    };
  }, []);

  return (
    <div style={{
      position: 'fixed',
      top: '20px',
      right: '20px',
      background: 'rgba(0, 0, 0, 0.8)',
      color: 'white',
      padding: '8px 12px',
      borderRadius: '6px',
      fontSize: '12px',
      fontFamily: 'monospace',
      zIndex: 9998
    }}>
      {status}
    </div>
  );
}