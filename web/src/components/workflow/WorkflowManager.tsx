/**
 * Workflow Manager Component
 * Orchestrates the complete end-to-end project generation workflow
 */

import React, { useEffect, useCallback } from 'react'
import { AnimatePresence, motion } from 'framer-motion'
import { useGenerationStore, useCurrentWorkflow } from '../../stores/generationStore'
import { useGenerationWorkflow } from '../../hooks/useGenerationWorkflow'
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

interface WorkflowManagerProps {
  className?: string
}

export default function WorkflowManager({ className = '' }: WorkflowManagerProps) {
  // Global state
  const {
    currentView,
    selectedBlueprint,
    projectConfig,
    isConfigValid,
    currentSession,
    generationResult,
    isLoading: globalLoading
  } = useCurrentWorkflow()

  const {
    navigateTo,
    setSelectedBlueprint,
    setAvailableBlueprints,
    updateProjectConfig,
    setCurrentSession,
    setGenerationResult,
    setLoading,
    resetApp,
    preferences
  } = useGenerationStore()

  // Workflow management
  const workflowManager = useGenerationWorkflow({
    autoSave: preferences.autoSave,
    persistSessions: true,
    maxSessions: 10
  })

  // Load available blueprints
  const { blueprints, loading: blueprintsLoading, error: blueprintsError } = useBlueprints()

  // Initialize blueprints
  useEffect(() => {
    if (blueprints && blueprints.length > 0) {
      setAvailableBlueprints(blueprints)
    }
  }, [blueprints, setAvailableBlueprints])

  // Sync workflow manager session with global state (prevent infinite loops)
  useEffect(() => {
    if (workflowManager.session && workflowManager.session.id !== currentSession?.id) {
      setCurrentSession(workflowManager.session)
    }
  }, [workflowManager.session?.id, currentSession?.id, setCurrentSession])

  // Sync generation result (prevent infinite loops)
  useEffect(() => {
    const workflowResult = workflowManager.session?.result
    if (workflowResult && workflowResult.projectId !== generationResult?.projectId) {
      setGenerationResult(workflowResult)
    }
  }, [workflowManager.session?.result?.projectId, generationResult?.projectId, setGenerationResult])

  // Auto-navigate based on workflow state
  useEffect(() => {
    if (workflowManager.session) {
      switch (workflowManager.session.state) {
        case 'generating':
        case 'validating':
          if (currentView !== 'generation') {
            navigateTo('generation')
          }
          break
        case 'completed':
          if (currentView !== 'success') {
            navigateTo('success')
          }
          break
        case 'error':
        case 'cancelled':
          if (currentView !== 'error') {
            navigateTo('error')
          }
          break
        case 'configuring':
          if (currentView === 'generation' || currentView === 'success' || currentView === 'error') {
            navigateTo('configuration')
          }
          break
        default:
          break
      }
    }
  }, [workflowManager.session?.state, currentView, navigateTo])

  // Set global loading state (debounced to prevent render loops)
  useEffect(() => {
    const loading = blueprintsLoading || workflowManager.isLoading
    if (loading !== globalLoading) {
      setLoading(loading)
    }
  }, [blueprintsLoading, workflowManager.isLoading, globalLoading, setLoading])

  // Handlers for workflow actions
  const handleBlueprintSelect = useCallback(async (blueprint: any) => {
    setSelectedBlueprint(blueprint)
    
    // Start a new workflow session
    workflowManager.startNewSession(blueprint)
    
    // Navigate to configuration
    navigateTo('configuration')
  }, [setSelectedBlueprint, workflowManager, navigateTo])

  const handleConfigurationComplete = useCallback(async () => {
    if (!selectedBlueprint || !isConfigValid) {
      console.error('Blueprint or configuration not ready')
      return
    }

    try {
      // Start the generation process
      await workflowManager.startGeneration(selectedBlueprint, projectConfig)
      // Navigation will be handled by the workflow state effect
    } catch (error) {
      console.error('Failed to start generation:', error)
      navigateTo('error')
    }
  }, [selectedBlueprint, isConfigValid, projectConfig, workflowManager, navigateTo])

  const handleGenerationCancel = useCallback(() => {
    workflowManager.cancelGeneration()
    navigateTo('configuration')
  }, [workflowManager, navigateTo])

  const handleGenerationRetry = useCallback(() => {
    if (selectedBlueprint && isConfigValid) {
      handleConfigurationComplete()
    } else {
      navigateTo('configuration')
    }
  }, [selectedBlueprint, isConfigValid, handleConfigurationComplete, navigateTo])

  const handleDownload = useCallback(async () => {
    try {
      await workflowManager.handleDownload()
    } catch (error) {
      console.error('Download failed:', error)
      // Show download error but don't navigate away from success page
    }
  }, [workflowManager])

  const handleNewProject = useCallback(() => {
    workflowManager.resetSession()
    resetApp()
  }, [workflowManager, resetApp])

  const handleGoBack = useCallback(() => {
    switch (currentView) {
      case 'configuration':
        navigateTo('blueprint-selection')
        break
      case 'generation':
        workflowManager.cancelGeneration()
        navigateTo('configuration')
        break
      case 'success':
      case 'error':
        navigateTo('configuration')
        break
      default:
        break
    }
  }, [currentView, navigateTo, workflowManager])

  // Handle errors from blueprints API
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
                  isValid={isConfigValid}
                  onConfigChange={updateProjectConfig}
                  onGenerate={handleConfigurationComplete}
                  onBack={handleGoBack}
                  preferences={preferences}
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
                  session={currentSession}
                  blueprint={selectedBlueprint}
                  config={projectConfig}
                  onCancel={handleGenerationCancel}
                  onComplete={(result) => {
                    setGenerationResult(result)
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
                    navigateTo('error')
                  }}
                  wsConnected={workflowManager.wsConnected}
                  progressData={workflowManager.progressData}
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
                    // TODO: Implement preview functionality
                    console.log('Show preview')
                  }}
                  downloading={workflowManager.isLoading}
                  downloadError={workflowManager.downloadError}
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
                  session={currentSession}
                  onRetry={handleGenerationRetry}
                  onBack={handleGoBack}
                  onNewProject={handleNewProject}
                />
              </motion.div>
            )}
          </AnimatePresence>
        </div>

        {/* Global Loading Overlay */}
        {globalLoading && (
          <LoadingOverlay
            message="Loading..."
            description="Please wait while we prepare your workspace"
          />
        )}
      </ErrorBoundary>
    </div>
  )
}

// Export WorkflowManager as the default component
export { WorkflowManager }

// Export view components for individual use if needed
export {
  BlueprintSelectionView,
  ConfigurationView,
  GenerationView,
  SuccessView,
  ErrorView
}