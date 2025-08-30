import React, { useState } from 'react'
import Header from './components/layout/Header'
import WorkflowManager from './components/workflow/WorkflowManagerSimplified'
import { ErrorBoundary } from './components/common/SimpleErrorBoundary'
import { DiagnosticConsole, DiagnosticStatus } from './DiagnosticConsole'
import { ComponentTest } from './ComponentTest'

function App() {
  // Local state to avoid Zustand store issues
  const [disclosureMode, setDisclosureMode] = useState<'basic' | 'advanced'>('basic')

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50 to-purple-100 relative overflow-hidden">
      {/* Skip Links for Keyboard Navigation */}
      <nav className="sr-only focus-within:not-sr-only" aria-label="Skip links">
        <div className="flex space-x-2 p-2 bg-blue-600">
          <a
            href="#main-content"
            className="inline-block px-3 py-1 bg-white text-blue-600 rounded text-sm font-medium focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-blue-600"
          >
            Skip to main content
          </a>
          <a
            href="#header-navigation"
            className="inline-block px-3 py-1 bg-white text-blue-600 rounded text-sm font-medium focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-blue-600"
          >
            Skip to navigation
          </a>
        </div>
      </nav>

      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-gradient-to-r from-purple-400 to-pink-400 rounded-full mix-blend-multiply filter blur-xl opacity-30 animate-pulse"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-gradient-to-r from-blue-400 to-cyan-400 rounded-full mix-blend-multiply filter blur-xl opacity-30 animate-pulse animation-delay-2000"></div>
        <div className="absolute top-40 left-1/2 transform -translate-x-1/2 w-80 h-80 bg-gradient-to-r from-indigo-400 to-purple-400 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse animation-delay-4000"></div>
      </div>
      
      {/* Header Component */}
      <div id="header-navigation">
        <Header 
          disclosureMode={disclosureMode}
          onDisclosureModeChange={setDisclosureMode}
        />
      </div>
      
      {/* Main Content - Fixed Workflow Manager */}
      <div id="main-content" className="h-[calc(100vh-4rem)] relative z-10">
        <ErrorBoundary 
          onError={(error) => console.error('Application error:', error)}
          enableRetry={true}
          level="page"
        >
          <WorkflowManager className="h-full" />
        </ErrorBoundary>
      </div>
      
      {/* Glass-morphism decorative elements */}
      <div className="fixed top-20 left-6 w-32 h-32 bg-gradient-to-r from-blue-200/30 to-purple-200/30 rounded-full backdrop-blur-sm border border-white/20 opacity-60 pointer-events-none hidden lg:block"></div>
      <div className="fixed bottom-40 right-12 w-24 h-24 bg-gradient-to-r from-pink-200/30 to-orange-200/30 rounded-full backdrop-blur-sm border border-white/20 opacity-60 pointer-events-none hidden lg:block"></div>
      
      {/* Diagnostic Tools */}
      <DiagnosticStatus />
      <DiagnosticConsole />
      <ComponentTest />
    </div>
  )
}

export default App
