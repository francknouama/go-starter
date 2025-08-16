import { useState } from 'react'
import Header from './components/layout/Header'
import ConfigurationPanel from './components/forms/ConfigurationPanel'
import FileExplorerPanel from './components/preview/FileExplorerPanel'
import ResponsiveLayout from './components/layout/ResponsiveLayout'
// import PreviewPanel from './components/preview/PreviewPanel'
// import HelpSystem from './components/help/HelpSystem'
// Error boundary and loading states temporarily disabled
// import { PageErrorBoundary } from './components/common/ErrorBoundary'
// import { LoadingOverlay } from './components/common/LoadingStates'
// import { useProjectWorkflow, connectWebSocket } from './hooks/useApi'
import type { ProjectConfig } from './services/api'
import type { DisclosureMode } from './types'

function App() {
  const [disclosureMode, setDisclosureMode] = useState<DisclosureMode>('basic')
  const [config, setConfig] = useState<ProjectConfig>({
    projectName: '',
    moduleUrl: '',
    goVersion: '1.21',
    projectType: 'web-api',
    framework: 'gin',
    architecture: 'standard',
    logger: 'slog',
  })
  // const [hasUsedTemplate, setHasUsedTemplate] = useState(false)

  // Temporarily disabled for debugging
  // const {
  //   defaultConfig,
  //   blueprints,
  //   preview,
  //   generation,
  //   generationProgress,
  //   generatePreview,
  //   generateProject,
  //   downloadProject,
  //   isLoading,
  // } = useProjectWorkflow()

  const isLoading = false
  const blueprints: any[] = [] // Placeholder for blueprints

  // Handlers for ConfigurationPanel
  const handleGenerate = (request: any) => {
    console.log('🚀 Generate project request:', request)
    // TODO: Implement project generation
  }

  const handleDownload = (request: any) => {
    console.log('📥 Download project request:', request)
    // TODO: Implement project download
  }

  // Initialize WebSocket connection and default config
  // useEffect(() => {
  //   // Connect to WebSocket for real-time updates
  //   connectWebSocket().catch(error => {
  //     console.warn('WebSocket connection failed:', error)
  //   })

  //   // Load default configuration
  //   if (defaultConfig) {
  //     setConfig(current => ({
  //       ...current, // Preserve any user changes
  //       ...defaultConfig, // Apply defaults for missing fields
  //     } as ProjectConfig))
  //   }
  // }, [defaultConfig])

  // Auto-generate preview when config changes
  // useEffect(() => {
  //   if (config.projectName && config.moduleUrl) {
  //     const timeoutId = setTimeout(() => {
  //       generatePreview(config).catch(error => {
  //         console.warn('Preview generation failed:', error)
  //       })
  //     }, 500) // Debounce preview generation

  //     return () => clearTimeout(timeoutId)
  //   }
  // }, [config, generatePreview])

  // Enhanced config change handler to track template usage
  // const handleConfigChange = (newConfig: ProjectConfig) => {
  //   setConfig(newConfig)
  //   // Simple heuristic to detect if a template was applied
  //   // (in a real app, you'd track this more explicitly)
  //   if (!hasUsedTemplate && 
  //       (newConfig.architecture !== 'standard' || 
  //        newConfig.framework !== 'gin' || 
  //        newConfig.logger !== 'slog')) {
  //     setHasUsedTemplate(true)
  //   }
  // }

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
      
      {/* Loading overlay with glass-morphism */}
      {isLoading && (
        <div className="fixed inset-0 bg-black bg-opacity-20 backdrop-blur-md z-50 flex items-center justify-center">
          <div className="bg-white/80 backdrop-blur-xl border border-white/20 rounded-2xl p-8 shadow-2xl flex items-center space-x-4">
            <div className="relative">
              <div className="animate-spin rounded-full h-8 w-8 border-3 border-transparent border-t-blue-600 border-r-purple-600"></div>
              <div className="absolute inset-0 animate-ping rounded-full h-8 w-8 border border-blue-400 opacity-25"></div>
            </div>
            <div>
              <span className="text-gray-800 font-semibold text-lg">Generating Project...</span>
              <p className="text-gray-600 text-sm mt-1">Building your Go application</p>
            </div>
          </div>
        </div>
      )}
      
      {/* Header Component */}
      <div id="header-navigation">
        <Header 
          disclosureMode={disclosureMode}
          onDisclosureModeChange={setDisclosureMode}
        />
      </div>
      
      <div id="main-content" className="container mx-auto py-8 relative z-10">
        <ResponsiveLayout
          configurationPanel={
            <ConfigurationPanel
              disclosureMode={disclosureMode}
              onModeChange={setDisclosureMode}
              config={config}
              onConfigChange={setConfig}
              blueprints={blueprints}
              onGenerate={handleGenerate}
              onDownload={handleDownload}
            />
          }
          fileExplorerPanel={
            <FileExplorerPanel
              config={config}
              preview={null}
            />
          }
        />
      </div>
      
      {/* Glass-morphism decorative elements */}
      <div className="fixed top-20 left-6 w-32 h-32 bg-gradient-to-r from-blue-200/30 to-purple-200/30 rounded-full backdrop-blur-sm border border-white/20 opacity-60 pointer-events-none hidden lg:block"></div>
      <div className="fixed bottom-40 right-12 w-24 h-24 bg-gradient-to-r from-pink-200/30 to-orange-200/30 rounded-full backdrop-blur-sm border border-white/20 opacity-60 pointer-events-none hidden lg:block"></div>
    </div>
  )
}

export default App
