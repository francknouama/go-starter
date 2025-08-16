import React, { useEffect } from 'react'

function TestApp() {
  useEffect(() => {
    console.log('🚀 TestApp mounted successfully!')
    console.log('Root element:', document.getElementById('root'))
    document.title = 'Go Starter - React Working!'
  }, [])

  return (
    <div style={{ padding: '20px', backgroundColor: '#ff0000', minHeight: '100vh', color: 'white' }}>
      <h1 style={{ color: 'white', fontSize: '48px', fontWeight: 'bold' }}>🚀 REACT IS WORKING!</h1>
      <p style={{ color: 'white', fontSize: '24px', marginTop: '20px' }}>
        If you can see this RED background, React is rendering correctly!
      </p>
      <div style={{ 
        marginTop: '20px', 
        padding: '20px', 
        backgroundColor: '#00ff00', 
        color: 'black', 
        borderRadius: '8px',
        fontSize: '20px',
        fontWeight: 'bold'
      }}>
        ✅ SUCCESS: React rendering is functional
      </div>
      <div style={{ marginTop: '20px', fontSize: '18px', color: 'yellow' }}>
        Current Time: {new Date().toLocaleTimeString()}
      </div>
      <div style={{ marginTop: '20px', fontSize: '16px', color: '#ccc' }}>
        If you don't see this, check the browser console for errors.
      </div>
    </div>
  )
}

export default TestApp