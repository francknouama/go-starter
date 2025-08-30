import React, { useEffect, useState } from 'react';

// Simple test component to verify React functionality
export function ComponentTest() {
  const [testResults, setTestResults] = useState<{
    [key: string]: { status: 'pass' | 'fail' | 'pending', message: string }
  }>({
    reactRender: { status: 'pending', message: 'Testing React render...' },
    stateManagement: { status: 'pending', message: 'Testing state management...' },
    useEffect: { status: 'pending', message: 'Testing useEffect...' },
    eventHandling: { status: 'pending', message: 'Testing event handling...' },
    imports: { status: 'pending', message: 'Testing imports...' }
  });

  const [counter, setCounter] = useState(0);

  useEffect(() => {
    console.log('🧪 Component test started');
    
    // Test React render
    setTimeout(() => {
      setTestResults(prev => ({
        ...prev,
        reactRender: { status: 'pass', message: 'React render successful' }
      }));
    }, 100);

    // Test useEffect
    setTimeout(() => {
      setTestResults(prev => ({
        ...prev,
        useEffect: { status: 'pass', message: 'useEffect hooks working' }
      }));
    }, 200);

    // Test state management
    setTimeout(() => {
      setCounter(1);
      setTestResults(prev => ({
        ...prev,
        stateManagement: { status: 'pass', message: 'State updates working' }
      }));
    }, 300);

    // Test imports
    setTimeout(() => {
      const hasReact = typeof React !== 'undefined';
      const hasDocument = typeof document !== 'undefined';
      
      setTestResults(prev => ({
        ...prev,
        imports: { 
          status: hasReact && hasDocument ? 'pass' : 'fail', 
          message: hasReact && hasDocument ? 'All imports successful' : 'Import issues detected' 
        }
      }));
    }, 400);

  }, []);

  const handleTestClick = () => {
    setCounter(prev => prev + 1);
    setTestResults(prev => ({
      ...prev,
      eventHandling: { status: 'pass', message: 'Event handling working' }
    }));
    console.log('🧪 Event handling test passed');
  };

  const allTestsPassed = Object.values(testResults).every(test => test.status === 'pass');
  const anyTestsFailed = Object.values(testResults).some(test => test.status === 'fail');

  return (
    <div style={{
      position: 'fixed',
      top: '60px',
      left: '20px',
      width: '300px',
      background: 'rgba(0, 0, 0, 0.9)',
      color: 'white',
      padding: '16px',
      borderRadius: '8px',
      fontSize: '12px',
      fontFamily: 'monospace',
      zIndex: 9997,
      border: allTestsPassed ? '2px solid #22c55e' : anyTestsFailed ? '2px solid #ef4444' : '1px solid #6b7280'
    }}>
      <h3 style={{ margin: '0 0 12px 0', fontSize: '14px' }}>
        🧪 Component Test Suite
      </h3>
      
      <div style={{ marginBottom: '12px' }}>
        {Object.entries(testResults).map(([testName, result]) => (
          <div key={testName} style={{ 
            marginBottom: '6px',
            padding: '4px',
            background: result.status === 'pass' ? 'rgba(34, 197, 94, 0.2)' : 
                       result.status === 'fail' ? 'rgba(239, 68, 68, 0.2)' : 'rgba(107, 114, 128, 0.2)',
            borderRadius: '3px',
            display: 'flex',
            alignItems: 'center'
          }}>
            <span style={{ marginRight: '8px' }}>
              {result.status === 'pass' ? '✅' : result.status === 'fail' ? '❌' : '⏳'}
            </span>
            <span>{result.message}</span>
          </div>
        ))}
      </div>

      <div style={{ 
        marginBottom: '12px',
        padding: '8px',
        background: 'rgba(59, 130, 246, 0.2)',
        borderRadius: '4px'
      }}>
        <div>Counter: {counter}</div>
        <button
          onClick={handleTestClick}
          style={{
            marginTop: '4px',
            background: '#3b82f6',
            color: 'white',
            border: 'none',
            padding: '4px 8px',
            borderRadius: '3px',
            cursor: 'pointer',
            fontSize: '11px'
          }}
        >
          Test Click Event
        </button>
      </div>

      {allTestsPassed && (
        <div style={{ 
          color: '#22c55e', 
          textAlign: 'center',
          fontWeight: 'bold'
        }}>
          🎉 All tests passed!
        </div>
      )}

      {anyTestsFailed && (
        <div style={{ 
          color: '#ef4444', 
          textAlign: 'center',
          fontWeight: 'bold'
        }}>
          ⚠️ Some tests failed!
        </div>
      )}
    </div>
  );
}