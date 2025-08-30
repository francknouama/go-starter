/**
 * Generation Workflow State Management Hook
 * Orchestrates the complete project generation workflow
 */

import { useState, useCallback, useRef, useEffect } from 'react'
import type { ProjectConfig, Blueprint, GenerationRequest } from '../services/api'
import { useConfiguration } from './useConfiguration'
import { useProjectGeneration, useProjectDownload, useGenerationProgress, useWebSocket } from './useApi'

export type WorkflowState = 
  | 'idle'           // Initial state, ready to start
  | 'configuring'    // User is configuring project
  | 'validating'     // Validating configuration
  | 'generating'     // Project generation in progress
  | 'completed'      // Generation completed successfully
  | 'error'          // Generation failed
  | 'cancelled'      // User cancelled generation

export interface GenerationResult {
  projectId: string
  downloadUrl?: string
  previewUrl?: string
  fileCount: number
  estimatedSize: string
  files: string[]
  success: boolean
  error?: string
  generationTime: number
}

export interface WorkflowSession {
  id: string
  blueprint?: Blueprint | null
  config: ProjectConfig
  state: WorkflowState
  result?: GenerationResult
  startTime?: Date
  progress: number
  currentStep: number
  estimatedDuration: number
}

interface UseGenerationWorkflowOptions {
  autoSave?: boolean
  persistSessions?: boolean
  maxSessions?: number
}

export function useGenerationWorkflow({
  autoSave = true,
  persistSessions = true,
  maxSessions = 5
}: UseGenerationWorkflowOptions = {}) {
  // Configuration management
  const configManager = useConfiguration({
    enablePersistence: autoSave,
    autoSave
  })

  // API hooks
  const { generateProject, generation, loading: generationLoading, error: generationError } = useProjectGeneration()
  const { downloadProject, downloading, error: downloadError } = useProjectDownload()
  const { connected: wsConnected } = useWebSocket()
  const progressData = useGenerationProgress()

  // Workflow state
  const [currentSession, setCurrentSession] = useState<WorkflowSession>({
    id: `session-${Date.now()}`,
    config: configManager.config,
    state: 'idle',
    progress: 0,
    currentStep: 0,
    estimatedDuration: 15
  })

  const [sessionHistory, setSessionHistory] = useState<WorkflowSession[]>([])
  const [isInitialized, setIsInitialized] = useState(false)

  // Refs for cleanup and state management
  const abortControllerRef = useRef<AbortController | null>(null)
  const sessionTimeoutRef = useRef<NodeJS.Timeout | null>(null)

  // Initialize workflow
  useEffect(() => {
    if (!isInitialized) {
      loadSessionHistory()
      setIsInitialized(true)
    }
  }, [isInitialized])

  // Sync configuration changes (prevent infinite loops)
  useEffect(() => {
    if ((currentSession.state === 'configuring' || currentSession.state === 'idle') && 
        JSON.stringify(currentSession.config) !== JSON.stringify(configManager.config)) {
      setCurrentSession(prev => ({
        ...prev,
        config: configManager.config
      }))
    }
  }, [configManager.config, currentSession.state]) // Remove currentSession.config to prevent loops

  // Handle WebSocket progress updates (prevent infinite loops)
  const currentProgressRef = useRef(currentSession.progress)
  const currentStepRef = useRef(currentSession.currentStep)
  
  useEffect(() => {
    currentProgressRef.current = currentSession.progress
    currentStepRef.current = currentSession.currentStep
  })
  
  useEffect(() => {
    if (progressData && currentSession.state === 'generating') {
      const { stage, progress, message, completed, error, projectId } = progressData
      
      // Only update if progress or step actually changed
      const newProgress = progress || currentProgressRef.current
      const newStep = Math.max(currentStepRef.current, getStepFromStage(stage))
      
      if (newProgress !== currentProgressRef.current || newStep !== currentStepRef.current) {
        setCurrentSession(prev => ({
          ...prev,
          progress: newProgress,
          currentStep: newStep
        }))
      }

      if (completed && !error && currentSession.state === 'generating') {
        handleGenerationSuccess({
          projectId: projectId || `project-${Date.now()}`,
          fileCount: 0, // Will be updated from API response
          estimatedSize: '0MB',
          files: [],
          success: true,
          generationTime: currentSession.startTime ? Math.ceil((Date.now() - currentSession.startTime.getTime()) / 1000) : 15
        })
      }

      if (error && currentSession.state === 'generating') {
        handleGenerationError(error)
      }
    }
  }, [progressData, currentSession.state, currentSession.startTime]) // Remove progress/step deps to prevent loops

  // Session persistence (debounced to prevent excessive saves)
  const saveSession = useCallback((session: WorkflowSession) => {
    if (!persistSessions) return

    try {
      const sessions = [...sessionHistory.slice(-(maxSessions - 1)), session]
      
      // Only update history if it actually changed
      if (JSON.stringify(sessions) !== JSON.stringify(sessionHistory)) {
        setSessionHistory(sessions)
      }
      
      // Debounce localStorage writes
      const timeoutId = setTimeout(() => {
        localStorage.setItem('go-starter-sessions', JSON.stringify({
          current: session,
          history: sessions,
          timestamp: Date.now()
        }))
      }, 100)
      
      return () => clearTimeout(timeoutId)
    } catch (error) {
      console.error('Failed to save session:', error)
    }
  }, [sessionHistory, persistSessions, maxSessions])

  const loadSessionHistory = useCallback(() => {
    if (!persistSessions) return

    try {
      const saved = localStorage.getItem('go-starter-sessions')
      if (saved) {
        const data = JSON.parse(saved)
        
        // Check if data is not too old (1 hour)
        if (Date.now() - data.timestamp < 60 * 60 * 1000) {
          setSessionHistory(data.history || [])
          
          // Only restore if session was in progress
          if (data.current && ['generating', 'validating'].includes(data.current.state)) {
            setCurrentSession(prev => ({
              ...prev,
              ...data.current,
              state: 'error' as WorkflowState, // Mark as error since we lost connection
              result: {
                ...data.current.result,
                success: false,
                error: 'Session interrupted - please try again'
              }
            }))
          }
        }
      }
    } catch (error) {
      console.error('Failed to load session history:', error)
    }
  }, [persistSessions])

  // Utility function to map stages to step numbers
  const getStepFromStage = (stage: string): number => {
    const stageMap: { [key: string]: number } = {
      validation: 0,
      'template-processing': 1,
      'code-generation': 2,
      optimization: 3,
      finalization: 4
    }
    return stageMap[stage] || 0
  }

  // Start new session
  const startNewSession = useCallback((blueprint?: Blueprint | null) => {
    const newSession: WorkflowSession = {
      id: `session-${Date.now()}`,
      blueprint,
      config: configManager.config,
      state: 'configuring',
      progress: 0,
      currentStep: 0,
      estimatedDuration: calculateEstimatedDuration(blueprint, configManager.config)
    }

    setCurrentSession(newSession)
    saveSession(newSession)

    return newSession.id
  }, [configManager.config, saveSession])

  // Calculate estimated duration based on blueprint and config
  const calculateEstimatedDuration = useCallback((blueprint?: Blueprint | null, config?: ProjectConfig): number => {
    let baseTime = 10 // seconds

    if (blueprint) {
      const complexityMultiplier: { [key: string]: number } = {
        'cli-simple': 0.5,
        'cli-standard': 1.0,
        'web-api-standard': 1.2,
        'web-api-clean': 1.5,
        'grpc-gateway': 2.0,
        'microservice-standard': 2.2,
        'monolith': 2.5
      }
      
      const multiplier = complexityMultiplier[blueprint.id] || 1.0
      baseTime *= multiplier
    }

    if (config?.features) {
      // Add time for additional features
      if (config.features.database) baseTime += 3
      if (config.features.authentication) baseTime += 2
      if (config.features.advanced) {
        const advancedCount = Object.values(config.features.advanced).filter(f => f.enabled).length
        baseTime += advancedCount * 1.5
      }
    }

    return Math.ceil(baseTime)
  }, [])

  // Update session state
  const updateSessionState = useCallback((updates: Partial<WorkflowSession>) => {
    setCurrentSession(prev => {
      const updated = { ...prev, ...updates }
      saveSession(updated)
      return updated
    })
  }, [saveSession])

  // Start generation process
  const startGeneration = useCallback(async (blueprint: Blueprint, config: ProjectConfig): Promise<string> => {
    // Validate configuration first
    const isValid = await configManager.validateConfig()
    if (!isValid) {
      throw new Error('Configuration validation failed')
    }

    // Create abort controller for cancellation
    abortControllerRef.current = new AbortController()

    const session: WorkflowSession = {
      ...currentSession,
      blueprint,
      config,
      state: 'validating',
      startTime: new Date(),
      progress: 0,
      currentStep: 0,
      estimatedDuration: calculateEstimatedDuration(blueprint, config)
    }

    setCurrentSession(session)
    saveSession(session)

    try {
      // Move to generating state
      updateSessionState({ state: 'generating' })

      const generationRequest: GenerationRequest = {
        ...config,
        outputFormat: 'zip',
        includeTests: true,
        includeDocs: true
      }

      const result = await generateProject(generationRequest)

      // If we get here, generation was successful
      await handleGenerationSuccess({
        projectId: result.projectId,
        downloadUrl: result.downloadUrl,
        previewUrl: result.previewUrl,
        fileCount: result.fileCount,
        estimatedSize: result.estimatedSize,
        files: [],
        success: true,
        generationTime: session.startTime ? Math.ceil((Date.now() - session.startTime.getTime()) / 1000) : session.estimatedDuration
      })

      return result.projectId
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Generation failed'
      await handleGenerationError(errorMessage)
      throw error
    }
  }, [currentSession, configManager, calculateEstimatedDuration, updateSessionState, saveSession, generateProject])

  // Handle successful generation
  const handleGenerationSuccess = useCallback(async (result: Omit<GenerationResult, 'generationTime'> & { generationTime?: number }) => {
    const finalResult: GenerationResult = {
      ...result,
      generationTime: result.generationTime || (currentSession.startTime ? Math.ceil((Date.now() - currentSession.startTime.getTime()) / 1000) : currentSession.estimatedDuration)
    }

    updateSessionState({
      state: 'completed',
      result: finalResult,
      progress: 100,
      currentStep: 4 // Final step
    })

    return finalResult
  }, [currentSession, updateSessionState])

  // Handle generation error
  const handleGenerationError = useCallback(async (error: string) => {
    const errorResult: GenerationResult = {
      projectId: '',
      fileCount: 0,
      estimatedSize: '0MB',
      files: [],
      success: false,
      error,
      generationTime: currentSession.startTime ? Math.ceil((Date.now() - currentSession.startTime.getTime()) / 1000) : 0
    }

    updateSessionState({
      state: 'error',
      result: errorResult
    })

    return errorResult
  }, [currentSession, updateSessionState])

  // Cancel generation
  const cancelGeneration = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    updateSessionState({
      state: 'cancelled',
      result: {
        projectId: '',
        fileCount: 0,
        estimatedSize: '0MB',
        files: [],
        success: false,
        error: 'Generation cancelled by user',
        generationTime: currentSession.startTime ? Math.ceil((Date.now() - currentSession.startTime.getTime()) / 1000) : 0
      }
    })
  }, [currentSession, updateSessionState])

  // Reset session
  const resetSession = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    const resetSession: WorkflowSession = {
      id: `session-${Date.now()}`,
      config: configManager.config,
      state: 'idle',
      progress: 0,
      currentStep: 0,
      estimatedDuration: 15
    }

    setCurrentSession(resetSession)
    saveSession(resetSession)
  }, [configManager.config, saveSession])

  // Download project
  const handleDownload = useCallback(async (): Promise<void> => {
    if (!currentSession.result?.success || !currentSession.config.projectName) {
      throw new Error('No successful generation to download')
    }

    const generationRequest: GenerationRequest = {
      ...currentSession.config,
      outputFormat: 'zip',
      includeTests: true,
      includeDocs: true
    }

    await downloadProject(generationRequest, `${currentSession.config.projectName}.zip`)
  }, [currentSession, downloadProject])

  // Load previous session
  const loadSession = useCallback((sessionId: string) => {
    const session = sessionHistory.find(s => s.id === sessionId)
    if (session) {
      setCurrentSession(session)
      configManager.loadConfig(session.config)
    }
  }, [sessionHistory, configManager])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (abortControllerRef.current) {
        abortControllerRef.current.abort()
      }
      if (sessionTimeoutRef.current) {
        clearTimeout(sessionTimeoutRef.current)
      }
    }
  }, [])

  return {
    // Current state
    session: currentSession,
    sessionHistory,
    configManager,
    
    // Status flags
    isGenerating: currentSession.state === 'generating',
    isCompleted: currentSession.state === 'completed',
    isError: currentSession.state === 'error',
    isCancelled: currentSession.state === 'cancelled',
    isLoading: generationLoading || downloading,
    
    // WebSocket status
    wsConnected,
    progressData,
    
    // Errors
    generationError,
    downloadError,
    
    // Actions
    startNewSession,
    startGeneration,
    cancelGeneration,
    resetSession,
    handleDownload,
    loadSession,
    updateSessionState,
    
    // Utilities
    calculateEstimatedDuration
  }
}