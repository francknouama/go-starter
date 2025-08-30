/**
 * Fixed Workflow Manager Component
 * Addresses infinite render loop issues with proper state management
 */

import React, { useCallback, useMemo, useState, useRef } from 'react'
import { AnimatePresence, motion } from 'framer-motion'
import { useBlueprints } from '../../hooks/useApi'

// Import workflow step components
import BlueprintSelectionView from './BlueprintSelectionView'
import ConfigurationView from './ConfigurationView'
import GenerationView from './GenerationView'
import SuccessView from './SuccessView'
import ErrorView from './ErrorView'

// Import common components
import { LoadingOverlay } from '../common/LoadingStates'
import { ErrorBoundary } from '../common/ErrorBoundary'

export type AppView = 
  | 'blueprint-selection'
  | 'configuration'
  | 'generation'
  | 'success'
  | 'error'

interface WorkflowManagerProps {
  className?: string
}

export default function WorkflowManagerFixed({ className = '' }: WorkflowManagerProps) {
  // Local state to avoid Zustand store infinite loops
  const [currentView, setCurrentView] = useState<AppView>('blueprint-selection')
  const [selectedBlueprint, setSelectedBlueprint] = useState<any>(null)
  const [projectConfig, setProjectConfig] = useState({
    projectName: '',
    moduleUrl: '',
    goVersion: '1.21',
    projectType: 'web-api',
    framework: 'gin',
    architecture: 'standard',
    logger: 'slog'
  })
  const [isLoading, setIsLoading] = useState(false)
  const [generationResult, setGenerationResult] = useState<any>(null)

  // Load available blueprints (this is safe as it's read-only)
  const { blueprints, loading: blueprintsLoading, error: blueprintsError } = useBlueprints()

  // Navigation without store dependencies
  const navigateTo = useCallback((view: AppView) => {
    setCurrentView(view)
  }, [])

  // Blueprint selection handler
  const handleBlueprintSelect = useCallback(async (blueprint: any) => {
    setSelectedBlueprint(blueprint)
    navigateTo('configuration')
  }, [navigateTo])

  // Configuration completion handler
  const handleConfigurationComplete = useCallback(async () => {
    if (!selectedBlueprint) {
      console.error('No blueprint selected')
      return
    }

    setIsLoading(true)
    navigateTo('generation')

    try {
      // Simulate generation process
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      const mockResult = {
        projectId: `project-${Date.now()}`,
        fileCount: 25,
        estimatedSize: '2.5MB',
        files: [],
        success: true,
        generationTime: 2
      }
      
      setGenerationResult(mockResult)
      setIsLoading(false)
      navigateTo('success')
    } catch (error) {
      const errorResult = {
        projectId: '',
        fileCount: 0,
        estimatedSize: '0MB',
        files: [],
        success: false,
        error: error instanceof Error ? error.message : 'Generation failed',
        generationTime: 0
      }
      setGenerationResult(errorResult)
      setIsLoading(false)
      navigateTo('error')
    }
  }, [selectedBlueprint, navigateTo])

  // Other handlers
  const handleGenerationCancel = useCallback(() => {
    setIsLoading(false)
    navigateTo('configuration')
  }, [navigateTo])

  const handleGenerationRetry = useCallback(() => {
    handleConfigurationComplete()
  }, [handleConfigurationComplete])

  const handleDownload = useCallback(async () => {
    console.log('Download requested')
  }, [])

  const handleNewProject = useCallback(() => {
    setCurrentView('blueprint-selection')
    setSelectedBlueprint(null)
    setProjectConfig({
      projectName: '',
      moduleUrl: '',
      goVersion: '1.21',
      projectType: 'web-api',
      framework: 'gin',
      architecture: 'standard',
      logger: 'slog'
    })
    setGenerationResult(null)
    setIsLoading(false)
  }, [])

  const handleGoBack = useCallback(() => {
    switch (currentView) {
      case 'configuration':
        navigateTo('blueprint-selection')
        break
      case 'generation':
        setIsLoading(false)
        navigateTo('configuration')
        break
      case 'success':
      case 'error':
        navigateTo('configuration')
        break
      default:
        break
    }
  }, [currentView, navigateTo])

  // Handle blueprint loading errors
  if (blueprintsError) {
    return (
      <ErrorBoundary
        level="page"
        enableRetry={true}
        onError={(error) => {
          console.error('Blueprint loading error:', error)
        }}
      >
        <div className="flex items-center justify-center h-full">
          <div className="text-center p-6">
            <h2 className="text-xl font-semibold mb-2">Blueprint Loading Error</h2>
            <p className="text-gray-600 mb-4">{blueprintsError}</p>
            <button 
              onClick={() => window.location.reload()}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Refresh Page
            </button>
          </div>
        </div>
      </ErrorBoundary>
    )
  }

  // View transition animations
  const viewTransition = {
    initial: { opacity: 0, x: 20 },
    animate: { opacity: 1, x: 0 },
    exit: { opacity: 0, x: -20 },
    transition: { duration: 0.3, ease: 'easeInOut' }
  }

  return (
    <div className={`h-full flex flex-col ${className}`}>
      <ErrorBoundary 
        level="section"
        name={`${currentView}-view`}
        enableRetry={true}
        onError={(error) => {
          console.error(`Error in ${currentView} view:`, error)
        }}
      >
        <div className="flex-1 relative overflow-hidden">
          <AnimatePresence mode="wait">
            {/* Blueprint Selection View */}
            {currentView === 'blueprint-selection' && (
              <motion.div
                key="blueprint-selection"
                {...viewTransition}
                className="absolute inset-0"
              >
                <BlueprintSelectionView
                  blueprints={blueprints}
                  loading={blueprintsLoading}
                  onSelect={handleBlueprintSelect}
                  selectedBlueprint={selectedBlueprint}
                />
              </motion.div>
            )}

            {/* Configuration View */}
            {currentView === 'configuration' && (
              <motion.div
                key="configuration"
                {...viewTransition}
                className="absolute inset-0"
              >
                <ConfigurationView
                  blueprint={selectedBlueprint}
                  config={projectConfig}
                  isValid={projectConfig.projectName.length > 0}
                  onConfigChange={(updates) => setProjectConfig(prev => ({ ...prev, ...updates }))}
                  onGenerate={handleConfigurationComplete}
                  onBack={handleGoBack}
                  preferences={{
                    theme: 'system',
                    disclosureMode: 'basic',
                    autoSave: true,
                    showHelpTooltips: true,
                    defaultGoVersion: '1.21'
                  }}
                />
              </motion.div>
            )}

            {/* Generation View */}
            {currentView === 'generation' && (
              <motion.div
                key="generation"
                {...viewTransition}
                className="absolute inset-0"
              >
                <GenerationView
                  session={{
                    id: 'test-session',
                    config: projectConfig,
                    state: 'generating',
                    progress: 0,
                    currentStep: 0,
                    estimatedDuration: 15
                  }}
                  blueprint={selectedBlueprint}
                  config={projectConfig}
                  onCancel={handleGenerationCancel}
                  onComplete={(result) => {
                    setGenerationResult(result)
                    setIsLoading(false)
                    navigateTo('success')
                  }}
                  onError={(error) => {
                    setGenerationResult({
                      projectId: '',
                      fileCount: 0,
                      estimatedSize: '0MB',
                      files: [],
                      success: false,
                      error: error.message || 'Generation failed',
                      generationTime: 0
                    })
                    setIsLoading(false)
                    navigateTo('error')
                  }}
                  wsConnected={false}
                  progressData={null}
                />
              </motion.div>
            )}

            {/* Success View */}
            {currentView === 'success' && generationResult?.success && (
              <motion.div
                key="success"
                {...viewTransition}
                className="absolute inset-0"
              >
                <SuccessView
                  result={generationResult}
                  blueprint={selectedBlueprint}
                  config={projectConfig}
                  onDownload={handleDownload}
                  onNewProject={handleNewProject}
                  onShowPreview={() => {
                    console.log('Show preview')
                  }}
                  downloading={false}
                  downloadError={null}
                />
              </motion.div>
            )}

            {/* Error View */}
            {currentView === 'error' && (
              <motion.div
                key="error"
                {...viewTransition}
                className="absolute inset-0"
              >
                <ErrorView
                  error={generationResult?.error || 'An unknown error occurred'}
                  session={{
                    id: 'test-session',
                    config: projectConfig,
                    state: 'error',
                    progress: 0,
                    currentStep: 0,
                    estimatedDuration: 15
                  }}
                  onRetry={handleGenerationRetry}
                  onBack={handleGoBack}
                  onNewProject={handleNewProject}
                />
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        {/* Global Loading Overlay */}
        {isLoading && (
          <LoadingOverlay
            message="Generating Project..."
            description="Building your Go application"
          />
        )}
      </ErrorBoundary>
    </div>
  )
}

// Export WorkflowManager as the fixed component
export { WorkflowManagerFixed as WorkflowManager }