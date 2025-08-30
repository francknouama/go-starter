/**
 * Simplified Workflow Manager Component
 * For debugging and isolating runtime errors
 */

import React, { useState, useCallback, useEffect } from 'react'
import { AnimatePresence, motion } from 'framer-motion'
import { ErrorBoundary } from '../common/SimpleErrorBoundary'
import type { Blueprint } from '../../services/api'

export type AppView = 
  | 'blueprint-selection'
  | 'configuration'
  | 'generation'
  | 'success'
  | 'error'

interface WorkflowManagerSimplifiedProps {
  className?: string
}

export default function WorkflowManagerSimplified({ className = '' }: WorkflowManagerSimplifiedProps) {
  const [currentView, setCurrentView] = useState<AppView>('blueprint-selection')
  const [selectedBlueprint, setSelectedBlueprint] = useState<any>(null)
  const [blueprints, setBlueprints] = useState<Blueprint[]>([])
  const [blueprintsLoading, setBlueprintsLoading] = useState(true)
  const [blueprintsError, setBlueprintsError] = useState<string | null>(null)

  // Simple blueprint loading without hooks to avoid async listener issues
  useEffect(() => {
    const loadBlueprints = async () => {
      try {
        setBlueprintsLoading(true)
        setBlueprintsError(null)
        
        // Simple delay and mock data to avoid promise issues
        await new Promise(resolve => setTimeout(resolve, 500))
        
        const mockBlueprints: Blueprint[] = [
          {
            id: 'web-api-standard',
            name: 'Web API Standard',
            description: 'Standard REST API with Go and Gin framework',
            category: 'web-api',
            difficulty: 'beginner',
            frameworks: ['gin'],
            architectures: ['standard'],
            features: ['rest-api', 'middleware', 'cors'],
            tags: ['api', 'rest', 'web'],
            estimatedFiles: 25
          },
          {
            id: 'cli-simple',
            name: 'CLI Simple',
            description: 'Simple command-line application',
            category: 'cli',
            difficulty: 'beginner',
            frameworks: ['cobra'],
            architectures: ['standard'],
            features: ['commands', 'flags'],
            tags: ['cli', 'tool'],
            estimatedFiles: 8
          },
          {
            id: 'grpc-gateway',
            name: 'gRPC Gateway',
            description: 'Dual HTTP/gRPC API with gateway',
            category: 'grpc',
            difficulty: 'advanced',
            frameworks: ['grpc'],
            architectures: ['standard'],
            features: ['grpc', 'http', 'protobuf'],
            tags: ['grpc', 'api', 'gateway'],
            estimatedFiles: 45
          }
        ]
        
        setBlueprints(mockBlueprints)
      } catch (error) {
        setBlueprintsError(error instanceof Error ? error.message : 'Failed to load blueprints')
      } finally {
        setBlueprintsLoading(false)
      }
    }

    loadBlueprints()
  }, [])

  console.log('WorkflowManagerSimplified render:', { currentView, blueprintsLoading, blueprintsError, blueprintCount: blueprints.length })

  const navigateTo = useCallback((view: AppView) => {
    console.log('Navigating to:', view)
    setCurrentView(view)
  }, [])

  const handleBlueprintSelect = useCallback((blueprint: any) => {
    console.log('Blueprint selected:', blueprint)
    setSelectedBlueprint(blueprint)
    navigateTo('configuration')
  }, [navigateTo])

  // Handle blueprint loading errors
  if (blueprintsError) {
    return (
      <div className={`h-full flex flex-col ${className}`}>
        <div className="flex items-center justify-center h-full">
          <div className="text-center p-6">
            <h2 className="text-xl font-semibold mb-2 text-red-600">Loading Error</h2>
            <p className="text-gray-600 mb-4">Failed to load blueprints: {blueprintsError}</p>
            <button 
              onClick={() => window.location.reload()}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Retry
            </button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <ErrorBoundary 
      level="page"
      onError={(error) => console.error('WorkflowManagerSimplified error:', error)}
      enableRetry={true}
    >
      <div className={`h-full flex flex-col ${className}`}>
        <div className="flex-1 relative overflow-hidden">
          <AnimatePresence mode="wait">
            {/* Blueprint Selection View */}
            {currentView === 'blueprint-selection' && (
              <motion.div
                key="blueprint-selection"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3, ease: 'easeInOut' }}
                className="absolute inset-0 p-6"
              >
                <div className="max-w-4xl mx-auto">
                  <h1 className="text-3xl font-bold mb-6">Select Blueprint</h1>
                  
                  {blueprintsLoading ? (
                    <div className="flex items-center justify-center h-64">
                      <div className="text-center">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                        <p className="text-gray-600">Loading blueprints...</p>
                      </div>
                    </div>
                  ) : (
                    <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
                      {blueprints && blueprints.length > 0 ? (
                        blueprints.map((blueprint) => (
                          <div 
                            key={blueprint.id}
                            className="p-4 border rounded-lg hover:shadow-md transition-shadow cursor-pointer"
                            onClick={() => handleBlueprintSelect(blueprint)}
                          >
                            <h3 className="font-semibold mb-2">{blueprint.name}</h3>
                            <p className="text-sm text-gray-600 mb-2">{blueprint.description}</p>
                            <div className="text-xs text-gray-500">
                              {blueprint.estimatedFiles} files • {blueprint.difficulty}
                            </div>
                          </div>
                        ))
                      ) : (
                        <div className="col-span-full text-center py-8">
                          <p className="text-gray-600">No blueprints available</p>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              </motion.div>
            )}

            {/* Configuration View */}
            {currentView === 'configuration' && (
              <motion.div
                key="configuration"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3, ease: 'easeInOut' }}
                className="absolute inset-0 p-6"
              >
                <div className="max-w-4xl mx-auto">
                  <div className="flex items-center mb-6">
                    <button
                      onClick={() => navigateTo('blueprint-selection')}
                      className="mr-4 p-2 hover:bg-gray-100 rounded"
                    >
                      ← Back
                    </button>
                    <h1 className="text-3xl font-bold">Configure Project</h1>
                  </div>
                  
                  {selectedBlueprint ? (
                    <div className="bg-white p-6 rounded-lg shadow">
                      <h2 className="text-xl font-semibold mb-4">
                        Configuring: {selectedBlueprint.name}
                      </h2>
                      <p className="text-gray-600 mb-4">{selectedBlueprint.description}</p>
                      
                      <div className="space-y-4">
                        <div>
                          <label className="block text-sm font-medium mb-1">Project Name</label>
                          <input
                            type="text"
                            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="my-go-project"
                          />
                        </div>
                        
                        <div className="pt-4">
                          <button
                            onClick={() => navigateTo('generation')}
                            className="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                          >
                            Generate Project
                          </button>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <p className="text-gray-600">No blueprint selected</p>
                      <button
                        onClick={() => navigateTo('blueprint-selection')}
                        className="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                      >
                        Select Blueprint
                      </button>
                    </div>
                  )}
                </div>
              </motion.div>
            )}

            {/* Generation View */}
            {currentView === 'generation' && (
              <motion.div
                key="generation"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3, ease: 'easeInOut' }}
                className="absolute inset-0 p-6"
              >
                <div className="max-w-4xl mx-auto">
                  <h1 className="text-3xl font-bold mb-6">Generating Project</h1>
                  <div className="bg-white p-6 rounded-lg shadow">
                    <div className="flex items-center justify-center h-64">
                      <div className="text-center">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                        <p className="text-gray-600">Generating your project...</p>
                        <div className="mt-4">
                          <button
                            onClick={() => navigateTo('success')}
                            className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 mr-2"
                          >
                            Simulate Success
                          </button>
                          <button
                            onClick={() => navigateTo('error')}
                            className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
                          >
                            Simulate Error
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </motion.div>
            )}

            {/* Success View */}
            {currentView === 'success' && (
              <motion.div
                key="success"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3, ease: 'easeInOut' }}
                className="absolute inset-0 p-6"
              >
                <div className="max-w-4xl mx-auto">
                  <h1 className="text-3xl font-bold mb-6 text-green-600">Success!</h1>
                  <div className="bg-white p-6 rounded-lg shadow">
                    <p className="text-gray-600 mb-4">Your project has been generated successfully!</p>
                    <button
                      onClick={() => navigateTo('blueprint-selection')}
                      className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                    >
                      Start New Project
                    </button>
                  </div>
                </div>
              </motion.div>
            )}

            {/* Error View */}
            {currentView === 'error' && (
              <motion.div
                key="error"
                initial={{ opacity: 0, x: 20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: -20 }}
                transition={{ duration: 0.3, ease: 'easeInOut' }}
                className="absolute inset-0 p-6"
              >
                <div className="max-w-4xl mx-auto">
                  <h1 className="text-3xl font-bold mb-6 text-red-600">Error</h1>
                  <div className="bg-white p-6 rounded-lg shadow">
                    <p className="text-gray-600 mb-4">An error occurred during project generation.</p>
                    <button
                      onClick={() => navigateTo('configuration')}
                      className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                    >
                      Try Again
                    </button>
                  </div>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </ErrorBoundary>
  )
}

// Export both names
export { WorkflowManagerSimplified as WorkflowManager }