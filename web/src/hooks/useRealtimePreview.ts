import { useCallback, useEffect, useRef, useState } from 'react'
import { useWebSocket } from './useWebSocket'
import type { 
  ProjectConfig, 
  RealtimePreviewState,
  WSProgressData,
  WSFileTreeNode,
  WSFileContent,
  WSPreviewCompleteData,
  WSStatusData
} from '../types/index'

interface UseRealtimePreviewOptions {
  onPreviewComplete?: (data: WSPreviewCompleteData) => void
  onError?: (error: string) => void
}

interface UseRealtimePreviewReturn {
  previewState: RealtimePreviewState
  startPreview: (config: ProjectConfig) => void
  selectFile: (path: string) => WSFileContent | null
  getFileContent: (path: string) => string | null
  connectionState: any
  clearPreview: () => void
}

export function useRealtimePreview(options: UseRealtimePreviewOptions = {}): UseRealtimePreviewReturn {
  const { onPreviewComplete, onError } = options
  
  const [previewState, setPreviewState] = useState<RealtimePreviewState>({
    isGenerating: false,
    progress: null,
    fileTree: null,
    selectedFile: null,
    files: new Map(),
    status: 'Ready'
  })
  
  const currentRequestIdRef = useRef<string | null>(null)
  
  const updatePreviewState = useCallback((updates: Partial<RealtimePreviewState>) => {
    setPreviewState(prev => ({ ...prev, ...updates }))
  }, [])
  
  const { connectionState, sendMessage, subscribe } = useWebSocket({
    onConnect: () => {
      console.log('Real-time preview WebSocket connected')
    },
    onDisconnect: () => {
      console.log('Real-time preview WebSocket disconnected')
      updatePreviewState({
        isGenerating: false,
        status: 'Disconnected'
      })
    },
    onError: (error) => {
      console.error('Real-time preview WebSocket error:', error)
      updatePreviewState({
        isGenerating: false,
        status: 'Connection Error',
        error: 'WebSocket connection failed'
      })
      onError?.('WebSocket connection failed')
    }
  })
  
  // Subscribe to progress updates
  useEffect(() => {
    const unsubscribe = subscribe('progress', (data: WSProgressData) => {
      updatePreviewState({
        progress: data,
        status: data.message
      })
    })
    
    return unsubscribe
  }, [subscribe, updatePreviewState])
  
  // Subscribe to status updates
  useEffect(() => {
    const unsubscribe = subscribe('status', (data: WSStatusData) => {
      updatePreviewState({
        status: data.message
      })
    })
    
    return unsubscribe
  }, [subscribe, updatePreviewState])
  
  // Subscribe to file tree updates
  useEffect(() => {
    const unsubscribe = subscribe('file_tree', (data: WSFileTreeNode) => {
      updatePreviewState({
        fileTree: data
      })
    })
    
    return unsubscribe
  }, [subscribe, updatePreviewState])
  
  // Subscribe to file content updates
  useEffect(() => {
    const unsubscribe = subscribe('file_content', (data: WSFileContent) => {
      setPreviewState(prev => {
        const newFiles = new Map(prev.files)
        newFiles.set(data.path, data)
        
        return {
          ...prev,
          files: newFiles
        }
      })
    })
    
    return unsubscribe
  }, [subscribe])
  
  // Subscribe to completion events
  useEffect(() => {
    const unsubscribe = subscribe('complete', (data: WSPreviewCompleteData) => {
      updatePreviewState({
        isGenerating: false,
        status: 'Preview Complete',
        progress: {
          stage: 'complete',
          progress: 1.0,
          message: 'Preview generation complete',
          totalFiles: data.totalFiles,
          processedFiles: data.totalFiles
        }
      })
      
      onPreviewComplete?.(data)
    })
    
    return unsubscribe
  }, [subscribe, updatePreviewState, onPreviewComplete])
  
  // Subscribe to error events
  useEffect(() => {
    const unsubscribe = subscribe('error', (errorMessage: string) => {
      updatePreviewState({
        isGenerating: false,
        status: 'Error',
        error: errorMessage,
        progress: null
      })
      
      onError?.(errorMessage)
    })
    
    return unsubscribe
  }, [subscribe, updatePreviewState, onError])
  
  const startPreview = useCallback((config: ProjectConfig) => {
    // Generate unique request ID
    const requestId = `preview_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
    currentRequestIdRef.current = requestId
    
    // Clear previous state
    setPreviewState({
      isGenerating: true,
      progress: null,
      fileTree: null,
      selectedFile: null,
      files: new Map(),
      status: 'Initializing preview...'
    })
    
    // Convert config to the format expected by the backend
    const request = {
      name: config.projectName,
      modulePath: config.moduleUrl,
      type: config.projectType,
      architecture: config.architecture,
      framework: config.framework,
      logger: config.logger,
      goVersion: config.goVersion,
      complexity: 'standard', // Default complexity
      advanced: false,
      memoryMode: true, // Always use memory mode for preview
      includeExamples: false
    }
    
    // Send preview request via WebSocket
    sendMessage('preview_request', request, requestId)
  }, [sendMessage])
  
  const selectFile = useCallback((path: string): WSFileContent | null => {
    const file = previewState.files.get(path)
    if (file) {
      updatePreviewState({ selectedFile: file })
      return file
    }
    return null
  }, [previewState.files, updatePreviewState])
  
  const getFileContent = useCallback((path: string): string | null => {
    const file = previewState.files.get(path)
    return file?.content || null
  }, [previewState.files])
  
  const clearPreview = useCallback(() => {
    currentRequestIdRef.current = null
    setPreviewState({
      isGenerating: false,
      progress: null,
      fileTree: null,
      selectedFile: null,
      files: new Map(),
      status: 'Ready'
    })
  }, [])
  
  return {
    previewState,
    startPreview,
    selectFile,
    getFileContent,
    connectionState,
    clearPreview
  }
}