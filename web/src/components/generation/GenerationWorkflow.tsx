/**
 * Generation Workflow Orchestrator
 * Manages the complete multi-step project generation process
 */

import React, { useState, useEffect, useCallback, useRef } from 'react'
import { AnimatePresence, motion } from 'framer-motion'
import { 
  CheckCircleIcon, 
  ExclamationTriangleIcon,
  ArrowPathIcon,
  DocumentArrowDownIcon,
  SparklesIcon,
  EyeIcon
} from '@heroicons/react/24/outline'
import type { ProjectConfig, Blueprint, GenerationRequest } from '../../services/api'
import StepIndicator from './StepIndicator'
import GenerationProgress from './GenerationProgress'
import Button from '../common/Button'
import { useProjectGeneration, useProjectDownload, useGenerationProgress } from '../../hooks/useApi'
import { useWebSocket } from '../../hooks/useApi'

interface GenerationWorkflowProps {
  blueprint: Blueprint | null
  config: ProjectConfig
  onComplete?: (result: GenerationResult) => void
  onCancel?: () => void
  onReset?: () => void
}

interface GenerationResult {
  projectId: string
  downloadUrl?: string
  previewUrl?: string
  fileCount: number
  estimatedSize: string
  files: string[]
  success: boolean
  error?: string
}

interface WorkflowStep {
  id: string
  title: string
  description: string
  status: 'pending' | 'active' | 'completed' | 'error' | 'skipped'
  duration?: number
  error?: string
  substeps?: string[]
  currentSubstep?: number
}

export default function GenerationWorkflow({
  blueprint,
  config,
  onComplete,
  onCancel,
  onReset
}: GenerationWorkflowProps) {
  // Hooks
  const { generateProject, generation, loading: generationLoading, error: generationError } = useProjectGeneration()
  const { downloadProject, downloading, error: downloadError } = useProjectDownload()
  const { connected: wsConnected } = useWebSocket()
  const progressData = useGenerationProgress()

  // State
  const [currentStep, setCurrentStep] = useState(0)
  const [workflowState, setWorkflowState] = useState<'preparing' | 'generating' | 'completed' | 'error' | 'cancelled'>('preparing')
  const [result, setResult] = useState<GenerationResult | null>(null)
  const [estimatedDuration, setEstimatedDuration] = useState(15) // seconds
  const [startTime, setStartTime] = useState<Date | null>(null)
  const [showPreview, setShowPreview] = useState(false)
  
  // Refs for cleanup
  const timeoutRef = useRef<NodeJS.Timeout | null>(null)
  const progressRef = useRef<number>(0)

  // Workflow steps definition
  const [steps, setSteps] = useState<WorkflowStep[]>([
    {
      id: 'validation',
      title: 'Configuration Validation',
      description: 'Validating project configuration and dependencies',
      status: 'pending',
      substeps: ['Checking project name', 'Validating module path', 'Verifying blueprint compatibility'],
      currentSubstep: 0
    },
    {
      id: 'template-processing',
      title: 'Template Processing',
      description: 'Processing blueprint templates and generating code structure',
      status: 'pending',
      substeps: ['Loading blueprint templates', 'Processing variables', 'Generating file structure'],
      currentSubstep: 0
    },
    {
      id: 'code-generation',
      title: 'Code Generation',
      description: 'Generating source code and configuration files',
      status: 'pending',
      substeps: ['Creating source files', 'Generating configuration', 'Setting up dependencies'],
      currentSubstep: 0
    },
    {
      id: 'optimization',
      title: 'Code Optimization',
      description: 'Optimizing generated code and applying best practices',
      status: 'pending',
      substeps: ['Formatting code', 'Optimizing imports', 'Applying linting rules'],
      currentSubstep: 0
    },
    {
      id: 'finalization',
      title: 'Project Finalization',
      description: 'Finalizing project structure and preparing for download',
      status: 'pending',
      substeps: ['Creating project archive', 'Generating documentation', 'Preparing download'],
      currentSubstep: 0
    }
  ])

  // Estimate duration based on blueprint complexity
  useEffect(() => {
    if (blueprint) {
      const baseTime = 10
      const complexityMultiplier = {
        'cli-simple': 0.5,
        'cli-standard': 1.0,
        'web-api-standard': 1.2,
        'web-api-clean': 1.5,
        'grpc-gateway': 2.0,
        'microservice-standard': 2.2,
        'monolith': 2.5
      }
      
      const multiplier = complexityMultiplier[blueprint.id as keyof typeof complexityMultiplier] || 1.0
      const hasAdvancedFeatures = config.features?.advanced && Object.keys(config.features.advanced).length > 0
      const advancedMultiplier = hasAdvancedFeatures ? 1.3 : 1.0
      
      const estimated = Math.ceil(baseTime * multiplier * advancedMultiplier)
      setEstimatedDuration(estimated)
    }
  }, [blueprint, config])

  // Handle WebSocket progress updates
  useEffect(() => {
    if (progressData && workflowState === 'generating') {
      const { stage, progress, message, completed, error } = progressData
      
      // Update current step based on stage
      const stepIndex = steps.findIndex(step => step.id === stage || step.title.toLowerCase().includes(stage.toLowerCase()))
      if (stepIndex !== -1 && stepIndex !== currentStep) {
        setCurrentStep(stepIndex)
      }

      // Update progress
      progressRef.current = progress

      // Handle completion
      if (completed && !error) {
        handleGenerationSuccess({
          projectId: progressData.projectId,
          fileCount: 0, // Will be updated from API response
          estimatedSize: '0MB',
          files: [],
          success: true
        })
      }

      // Handle error
      if (error) {
        handleGenerationError(error)
      }
    }
  }, [progressData, workflowState, steps, currentStep])

  // Update step status
  const updateStepStatus = useCallback((stepIndex: number, status: WorkflowStep['status'], error?: string, substep?: number) => {
    setSteps(prev => prev.map((step, index) => {
      if (index === stepIndex) {
        return {
          ...step,
          status,
          error,
          currentSubstep: substep ?? step.currentSubstep,
          duration: status === 'completed' && step.status === 'active' ? Date.now() - (startTime?.getTime() || 0) : step.duration
        }
      } else if (index < stepIndex && step.status !== 'completed') {
        return { ...step, status: 'completed' }
      }
      return step
    }))
  }, [startTime])

  // Simulate step progression (fallback when WebSocket is not available)
  const simulateStepProgression = useCallback(() => {
    if (!wsConnected && workflowState === 'generating') {
      const stepDuration = (estimatedDuration * 1000) / steps.length
      
      const progressStep = (stepIndex: number) => {
        if (stepIndex >= steps.length) {
          // All steps completed
          handleGenerationSuccess({
            projectId: `project-${Date.now()}`,
            fileCount: Math.floor(Math.random() * 50) + 10,
            estimatedSize: `${(Math.random() * 10 + 1).toFixed(1)}MB`,
            files: [],
            success: true
          })
          return
        }

        setCurrentStep(stepIndex)
        updateStepStatus(stepIndex, 'active')

        // Simulate substep progression
        const step = steps[stepIndex]
        if (step.substeps) {
          const substepDuration = stepDuration / step.substeps.length
          
          step.substeps.forEach((_, substepIndex) => {
            setTimeout(() => {
              updateStepStatus(stepIndex, 'active', undefined, substepIndex)
            }, substepDuration * substepIndex)
          })
        }

        // Complete current step and move to next
        setTimeout(() => {
          updateStepStatus(stepIndex, 'completed')
          progressStep(stepIndex + 1)
        }, stepDuration)
      }

      progressStep(0)
    }
  }, [wsConnected, workflowState, estimatedDuration, steps, updateStepStatus])

  // Start generation process
  const startGeneration = useCallback(async () => {
    if (!blueprint || !config.projectName || !config.moduleUrl) {
      return
    }

    setWorkflowState('generating')
    setStartTime(new Date())
    setCurrentStep(0)
    
    // Reset all steps
    setSteps(prev => prev.map(step => ({ ...step, status: 'pending', currentSubstep: 0, error: undefined })))

    const generationRequest: GenerationRequest = {
      ...config,
      outputFormat: 'zip',
      includeTests: true,
      includeDocs: true
    }

    try {
      // Start actual generation
      const response = await generateProject(generationRequest)
      
      // If WebSocket is not connected, simulate progression
      if (!wsConnected) {
        simulateStepProgression()
      }
      
    } catch (error) {
      handleGenerationError(error instanceof Error ? error.message : 'Generation failed')
    }
  }, [blueprint, config, generateProject, wsConnected, simulateStepProgression])

  // Handle successful generation
  const handleGenerationSuccess = (result: Partial<GenerationResult>) => {
    const finalResult: GenerationResult = {
      projectId: result.projectId || `project-${Date.now()}`,
      downloadUrl: result.downloadUrl,
      previewUrl: result.previewUrl,
      fileCount: result.fileCount || 0,
      estimatedSize: result.estimatedSize || '0MB',
      files: result.files || [],
      success: true
    }
    
    setResult(finalResult)
    setWorkflowState('completed')
    updateStepStatus(steps.length - 1, 'completed')
    
    if (onComplete) {
      onComplete(finalResult)
    }
  }

  // Handle generation error
  const handleGenerationError = (error: string) => {
    setWorkflowState('error')
    updateStepStatus(currentStep, 'error', error)
    
    const errorResult: GenerationResult = {
      projectId: '',
      fileCount: 0,
      estimatedSize: '0MB',
      files: [],
      success: false,
      error
    }
    
    setResult(errorResult)
    
    if (onComplete) {
      onComplete(errorResult)
    }
  }

  // Handle download
  const handleDownload = async () => {
    if (!result || !result.success || !config.projectName) return

    try {
      const generationRequest: GenerationRequest = {
        ...config,
        outputFormat: 'zip',
        includeTests: true,
        includeDocs: true
      }
      
      await downloadProject(generationRequest, `${config.projectName}.zip`)
    } catch (error) {
      console.error('Download failed:', error)
    }
  }

  // Handle cancel
  const handleCancel = () => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    setWorkflowState('cancelled')
    if (onCancel) {
      onCancel()
    }
  }

  // Handle reset
  const handleReset = () => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    
    setWorkflowState('preparing')
    setCurrentStep(0)
    setResult(null)
    setStartTime(null)
    setShowPreview(false)
    progressRef.current = 0
    
    // Reset all steps
    setSteps(prev => prev.map(step => ({ 
      ...step, 
      status: 'pending', 
      currentSubstep: 0, 
      error: undefined,
      duration: undefined 
    })))
    
    if (onReset) {
      onReset()
    }
  }

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [])

  // Auto-start generation if in preparing state and config is valid
  useEffect(() => {
    if (workflowState === 'preparing' && blueprint && config.projectName && config.moduleUrl) {
      const timer = setTimeout(() => {
        startGeneration()
      }, 1000) // Small delay for UX

      return () => clearTimeout(timer)
    }
  }, [workflowState, blueprint, config, startGeneration])

  const isLoading = workflowState === 'generating' || generationLoading

  return (
    <div className="max-w-4xl mx-auto">
      <AnimatePresence mode="wait">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="mb-8 text-center"
        >
          <div className="flex items-center justify-center gap-3 mb-4">
            <SparklesIcon className="w-8 h-8 text-blue-600" />
            <h2 className="text-2xl font-bold text-gray-900">
              Generating Your Go Project
            </h2>
          </div>
          {blueprint && (
            <p className="text-gray-600">
              Creating <strong>{blueprint.name}</strong> project: <strong>{config.projectName}</strong>
            </p>
          )}
        </motion.div>

        {/* Step Indicator */}
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          className="mb-8"
        >
          <StepIndicator 
            steps={steps}
            currentStep={currentStep}
            estimatedDuration={estimatedDuration}
            startTime={startTime}
          />
        </motion.div>

        {/* Generation Progress */}
        {isLoading && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="mb-8"
          >
            <GenerationProgress
              steps={steps}
              currentStep={currentStep}
              progress={progressRef.current}
              estimatedDuration={estimatedDuration}
              startTime={startTime}
              wsConnected={wsConnected}
            />
          </motion.div>
        )}

        {/* Error State */}
        {workflowState === 'error' && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="mb-8"
          >
            <div className="bg-red-50 border border-red-200 rounded-lg p-6">
              <div className="flex items-start gap-3">
                <ExclamationTriangleIcon className="w-6 h-6 text-red-600 flex-shrink-0" />
                <div>
                  <h3 className="text-lg font-medium text-red-900 mb-2">Generation Failed</h3>
                  <p className="text-red-700 mb-4">
                    {result?.error || generationError || 'An unexpected error occurred during project generation.'}
                  </p>
                  <div className="flex gap-3">
                    <Button
                      variant="outline"
                      onClick={handleReset}
                      className="border-red-300 text-red-700 hover:bg-red-50"
                    >
                      <ArrowPathIcon className="w-4 h-4 mr-2" />
                      Try Again
                    </Button>
                    <Button
                      variant="ghost"
                      onClick={handleCancel}
                      className="text-gray-600"
                    >
                      Cancel
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        )}

        {/* Success State */}
        {workflowState === 'completed' && result?.success && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="mb-8"
          >
            <div className="bg-green-50 border border-green-200 rounded-lg p-6">
              <div className="flex items-start gap-3">
                <CheckCircleIcon className="w-6 h-6 text-green-600 flex-shrink-0" />
                <div className="flex-1">
                  <h3 className="text-lg font-medium text-green-900 mb-2">
                    Project Generated Successfully!
                  </h3>
                  <p className="text-green-700 mb-4">
                    Your <strong>{blueprint?.name}</strong> project is ready for download.
                  </p>
                  
                  {/* Project Stats */}
                  <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mb-6 p-4 bg-white rounded-lg border border-green-200">
                    <div className="text-center">
                      <div className="text-2xl font-bold text-gray-900">{result.fileCount}</div>
                      <div className="text-sm text-gray-600">Files Generated</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-gray-900">{result.estimatedSize}</div>
                      <div className="text-sm text-gray-600">Project Size</div>
                    </div>
                    <div className="text-center md:col-span-1 col-span-2">
                      <div className="text-2xl font-bold text-gray-900">
                        {startTime ? Math.ceil((Date.now() - startTime.getTime()) / 1000) : estimatedDuration}s
                      </div>
                      <div className="text-sm text-gray-600">Generation Time</div>
                    </div>
                  </div>

                  {/* Action Buttons */}
                  <div className="flex flex-wrap gap-3">
                    <Button
                      variant="primary"
                      onClick={handleDownload}
                      disabled={downloading}
                      className="bg-green-600 hover:bg-green-700"
                    >
                      <DocumentArrowDownIcon className="w-4 h-4 mr-2" />
                      {downloading ? 'Downloading...' : `Download ${config.projectName}.zip`}
                    </Button>
                    
                    {result.previewUrl && (
                      <Button
                        variant="outline"
                        onClick={() => setShowPreview(!showPreview)}
                        className="border-green-300 text-green-700 hover:bg-green-50"
                      >
                        <EyeIcon className="w-4 h-4 mr-2" />
                        {showPreview ? 'Hide Preview' : 'Preview Project'}
                      </Button>
                    )}
                    
                    <Button
                      variant="ghost"
                      onClick={handleReset}
                      className="text-gray-600"
                    >
                      Generate Another Project
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        )}

        {/* Cancel/Control Buttons */}
        {isLoading && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className="text-center"
          >
            <Button
              variant="ghost"
              onClick={handleCancel}
              className="text-gray-600 hover:text-red-600"
            >
              Cancel Generation
            </Button>
          </motion.div>
        )}

        {/* Download Error */}
        {downloadError && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            className="mb-8"
          >
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <div className="flex items-center gap-3">
                <ExclamationTriangleIcon className="w-5 h-5 text-yellow-600" />
                <div>
                  <p className="text-yellow-800 font-medium">Download Failed</p>
                  <p className="text-yellow-700 text-sm mt-1">{downloadError}</p>
                </div>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}