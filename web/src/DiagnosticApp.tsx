import React from 'react'

function DiagnosticApp() {
  console.log('DiagnosticApp is rendering')
  
  return (
    <div className="p-8 bg-gray-100 min-h-screen">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold mb-4">Diagnostic Test App</h1>
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-2">Status Check</h2>
          <p className="text-green-600">✅ React is working correctly</p>
          <p className="text-green-600">✅ TypeScript is compiling</p>
          <p className="text-green-600">✅ Tailwind CSS is loaded</p>
          <p className="text-green-600">✅ No obvious runtime errors</p>
          
          <div className="mt-4">
            <button
              onClick={() => alert('Button click works!')}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Test Button
            </button>
          </div>
          
          <div className="mt-4 text-sm text-gray-600">
            <p>If you can see this and the button works, the basic setup is fine.</p>
            <p>The issue is likely in the complex App component structure.</p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default DiagnosticApp