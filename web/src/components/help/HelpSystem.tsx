import { useState, useEffect, useCallback } from 'react'
import { QuestionMarkCircleIcon } from '@heroicons/react/24/outline'
import QuickStartGuide from './QuickStartGuide'
import KeyboardShortcutsOverlay from './KeyboardShortcutsOverlay'
import SmartHelp from './SmartHelp'
import type { ProjectConfig, DisclosureMode } from '../../types'

interface HelpSystemProps {
  config: ProjectConfig
  disclosureMode: DisclosureMode
  onModeChange?: (mode: DisclosureMode) => void
  onConfigChange?: (config: ProjectConfig) => void
}

export default function HelpSystem({ 
  config, 
  disclosureMode, 
  onModeChange,
  onConfigChange 
}: HelpSystemProps) {
  const [showQuickStart, setShowQuickStart] = useState(false)
  const [showShortcuts, setShowShortcuts] = useState(false)
  const [helpTooltipsVisible, setHelpTooltipsVisible] = useState(false)
  const [isFirstVisit, setIsFirstVisit] = useState(false)

  // Check if this is user's first visit
  useEffect(() => {
    const hasVisited = localStorage.getItem('go-starter-visited')
    if (!hasVisited) {
      setIsFirstVisit(true)
      setShowQuickStart(true)
      localStorage.setItem('go-starter-visited', 'true')
    }
  }, [])

  // Keyboard shortcuts handler
  const handleKeyPress = useCallback((event: KeyboardEvent) => {
    // Don't trigger shortcuts when typing in inputs
    if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
      return
    }

    switch (event.key) {
      case '?':
        event.preventDefault()
        setShowShortcuts(true)
        break
      case 'Escape':
        event.preventDefault()
        setShowQuickStart(false)
        setShowShortcuts(false)
        break
    }

    // Shortcuts with modifiers
    if (event.altKey) {
      switch (event.key) {
        case 'h':
          event.preventDefault()
          setShowQuickStart(true)
          break
        case 'b':
          event.preventDefault()
          onModeChange?.('basic')
          break
        case 'a':
          event.preventDefault()
          onModeChange?.('advanced')
          break
        case 't':
          event.preventDefault()
          setHelpTooltipsVisible(!helpTooltipsVisible)
          break
      }
    }

    if (event.ctrlKey || event.metaKey) {
      switch (event.key) {
        case 'Enter':
          event.preventDefault()
          // Trigger form submission if valid
          // This would be handled by parent component
          break
        case 'r':
          event.preventDefault()
          // Reset form
          onConfigChange?.({
            projectName: '',
            moduleUrl: '',
            goVersion: '1.21',
            projectType: 'web-api',
            framework: 'gin',
            architecture: 'standard',
            logger: 'slog',
          })
          break
      }
    }

    // Quick project type selection (1-5)
    if (!event.ctrlKey && !event.altKey && !event.metaKey) {
      const projectTypes = ['cli', 'web-api', 'library', 'lambda', 'microservice'] as const
      const numKey = parseInt(event.key)
      if (numKey >= 1 && numKey <= 5) {
        event.preventDefault()
        onConfigChange?.({
          ...config,
          projectType: projectTypes[numKey - 1]
        })
      }
    }
  }, [config, disclosureMode, onModeChange, onConfigChange, helpTooltipsVisible])

  useEffect(() => {
    document.addEventListener('keydown', handleKeyPress)
    return () => document.removeEventListener('keydown', handleKeyPress)
  }, [handleKeyPress])

  // Add CSS class for tooltip visibility
  useEffect(() => {
    if (helpTooltipsVisible) {
      document.body.classList.add('show-all-tooltips')
    } else {
      document.body.classList.remove('show-all-tooltips')
    }
  }, [helpTooltipsVisible])

  return (
    <>
      {/* Floating help button */}
      <div className="fixed bottom-6 right-6 z-40 flex flex-col gap-3">
        {/* Smart help panel - contextual tips */}
        <div className="max-w-sm">
          <SmartHelp 
            projectType={config.projectType as any}
            architecture={config.architecture as any}
            framework={config.framework as any}
            disclosureMode={disclosureMode}
          />
        </div>

        {/* Help menu */}
        <div className="flex flex-col gap-2">
          <button
            onClick={() => setShowQuickStart(true)}
            className="bg-blue-600 hover:bg-blue-700 text-white p-3 rounded-full shadow-lg transition-all duration-200 hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            title="Quick Start Guide (Alt+H)"
            aria-label="Open quick start guide"
          >
            <QuestionMarkCircleIcon className="w-6 h-6" />
          </button>

          <button
            onClick={() => setShowShortcuts(true)}
            className="bg-gray-600 hover:bg-gray-700 text-white p-2 rounded-full shadow-lg transition-all duration-200 hover:scale-105 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2"
            title="Keyboard Shortcuts (?)"
            aria-label="Show keyboard shortcuts"
          >
            <kbd className="text-xs font-bold">?</kbd>
          </button>

          {/* Tooltips toggle */}
          <button
            onClick={() => setHelpTooltipsVisible(!helpTooltipsVisible)}
            className={`p-2 rounded-full shadow-lg transition-all duration-200 hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 ${
              helpTooltipsVisible
                ? 'bg-yellow-500 hover:bg-yellow-600 text-white focus:ring-yellow-500'
                : 'bg-gray-300 hover:bg-gray-400 text-gray-700 focus:ring-gray-500'
            }`}
            title="Toggle Help Tooltips (Alt+T)"
            aria-label="Toggle help tooltips visibility"
          >
            <span className="text-xs font-bold">💡</span>
          </button>
        </div>
      </div>

      {/* First visit welcome message */}
      {isFirstVisit && (
        <div className="fixed top-4 left-1/2 transform -translate-x-1/2 z-50 bg-blue-600 text-white px-6 py-3 rounded-lg shadow-lg">
          <p className="text-sm font-medium">
            👋 Welcome to go-starter! The quick start guide is opening to help you get started.
          </p>
        </div>
      )}

      {/* Keyboard shortcut hint */}
      <div className="fixed bottom-6 left-6 z-40">
        <div className="bg-gray-800 text-white text-xs px-3 py-2 rounded-lg shadow-lg opacity-75">
          Press <kbd className="bg-gray-700 px-1 py-0.5 rounded">?</kbd> for shortcuts
        </div>
      </div>

      {/* Help overlays */}
      <QuickStartGuide
        isOpen={showQuickStart}
        onClose={() => setShowQuickStart(false)}
        onProjectTypeSelect={(type) => {}}
        onModeSelect={onModeChange || (() => {})}
      />

      <KeyboardShortcutsOverlay
        isOpen={showShortcuts}
        onClose={() => setShowShortcuts(false)}
      />

      {/* Global CSS for tooltip visibility */}
      <style>{`
        .show-all-tooltips [data-tooltip]:hover::after {
          opacity: 1 !important;
          visibility: visible !important;
        }
        
        .show-all-tooltips .group:hover .group-hover\\:visible {
          visibility: visible !important;
        }
        
        .show-all-tooltips .group:hover .group-hover\\:opacity-100 {
          opacity: 1 !important;
        }
      `}</style>
    </>
  )
}