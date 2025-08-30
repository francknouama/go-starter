import { useMemo } from 'react'
import { 
  CheckCircleIcon, 
  ExclamationTriangleIcon, 
  InformationCircleIcon,
  WifiIcon,
  ClockIcon,
  DocumentIcon
} from '@heroicons/react/24/outline'
import type { WSProgressData, WSConnectionState, RealtimePreviewState } from '../../types'

interface ProgressTrackerProps {
  connectionState: WSConnectionState
  previewState: RealtimePreviewState
  className?: string
  variant?: 'compact' | 'detailed'
}

const ProgressStage = {
  INITIALIZING: 'initializing',
  LOADING_BLUEPRINT: 'loading',
  GENERATING: 'generating',
  BUILDING_TREE: 'building',
  COMPLETE: 'complete',
  ERROR: 'error'
}

const StageMessages = {
  [ProgressStage.INITIALIZING]: 'Initializing preview generation...',
  [ProgressStage.LOADING_BLUEPRINT]: 'Loading blueprint...',
  [ProgressStage.GENERATING]: 'Generating files...',
  [ProgressStage.BUILDING_TREE]: 'Building file tree...',
  [ProgressStage.COMPLETE]: 'Preview generation complete',
  [ProgressStage.ERROR]: 'Preview generation failed'
}

function getProgressColor(stage: string): string {
  switch (stage) {
    case ProgressStage.COMPLETE:
      return 'bg-green-600'
    case ProgressStage.ERROR:
      return 'bg-red-600'
    case ProgressStage.GENERATING:
      return 'bg-blue-600'
    default:
      return 'bg-gray-400'
  }
}

function getConnectionStatusIcon(connected: boolean, error?: string) {
  if (error) {
    return <ExclamationTriangleIcon className="h-4 w-4 text-red-500" />
  }
  if (connected) {
    return <WifiIcon className="h-4 w-4 text-green-500" />
  }
  return <WifiIcon className="h-4 w-4 text-gray-400" />
}

function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`
  return `${(ms / 60000).toFixed(1)}m`
}

export default function ProgressTracker({ 
  connectionState, 
  previewState, 
  className = '',
  variant = 'detailed'
}: ProgressTrackerProps) {
  const { progress, status, error, isGenerating, fileTree, files } = previewState
  
  const progressPercentage = useMemo(() => {
    if (!progress) return 0
    return Math.round(progress.progress * 100)
  }, [progress?.progress])
  
  const elapsedTime = useMemo(() => {
    if (!isGenerating || !progress) return null
    // This would ideally come from the backend, for now we simulate
    return Date.now() - (progress as any).startTime || Date.now()
  }, [isGenerating, progress])
  
  const fileStats = useMemo(() => {
    const fileCount = files.size
    const totalSize = Array.from(files.values())
      .reduce((sum, file) => sum + file.size, 0)
    
    return { fileCount, totalSize }
  }, [files])
  
  if (variant === 'compact') {
    return (
      <div className={`flex items-center gap-3 p-2 bg-gray-50 border rounded ${className}`}>
        {getConnectionStatusIcon(connectionState.connected, connectionState.error)}
        
        {isGenerating && progress && (
          <>
            <div className="flex-1 min-w-0">
              <div className="flex justify-between text-xs text-gray-600 mb-1">
                <span className="truncate">{progress.message}</span>
                <span>{progressPercentage}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-1">
                <div 
                  className={`h-1 rounded-full transition-all duration-300 ${getProgressColor(progress.stage)}`}
                  style={{ width: `${progressPercentage}%` }}
                />
              </div>
            </div>
            
            {progress.currentFile && (
              <div className="text-xs text-gray-500 truncate max-w-32">
                {progress.currentFile.split('/').pop()}
              </div>
            )}
          </>
        )}
        
        {!isGenerating && status && (
          <span className="text-sm text-gray-600 truncate">{status}</span>
        )}
        
        {error && (
          <ExclamationTriangleIcon className="h-4 w-4 text-red-500 flex-shrink-0" />
        )}
      </div>
    )
  }
  
  return (
    <div className={`bg-white border rounded-lg shadow-sm ${className}`}>
      {/* Header */}
      <div className="border-b bg-gray-50 px-4 py-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <h3 className="font-medium text-gray-900">Preview Generation</h3>
            {isGenerating && (
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600" />
            )}
          </div>
          
          <div className="flex items-center gap-3 text-sm">
            {/* Connection Status */}
            <div className="flex items-center gap-2">
              {getConnectionStatusIcon(connectionState.connected, connectionState.error)}
              <span className={connectionState.connected ? 'text-green-600' : 'text-red-600'}>
                {connectionState.connected ? 'Connected' : 'Disconnected'}
              </span>
            </div>
            
            {/* Elapsed Time */}
            {isGenerating && elapsedTime && (
              <div className="flex items-center gap-1 text-gray-500">
                <ClockIcon className="h-4 w-4" />
                <span>{formatDuration(elapsedTime)}</span>
              </div>
            )}
          </div>
        </div>
      </div>
      
      {/* Progress Content */}
      <div className="p-4">
        {/* Status Message */}
        <div className="flex items-start gap-3 mb-4">
          {error ? (
            <ExclamationTriangleIcon className="h-5 w-5 text-red-500 flex-shrink-0 mt-0.5" />
          ) : isGenerating ? (
            <InformationCircleIcon className="h-5 w-5 text-blue-500 flex-shrink-0 mt-0.5" />
          ) : (
            <CheckCircleIcon className="h-5 w-5 text-green-500 flex-shrink-0 mt-0.5" />
          )}
          
          <div className="flex-1 min-w-0">
            <div className="text-sm font-medium text-gray-900">
              {error ? 'Generation Failed' : status || 'Ready'}
            </div>
            {error && (
              <div className="text-sm text-red-600 mt-1">{error}</div>
            )}
          </div>
        </div>
        
        {/* Progress Bar */}
        {progress && isGenerating && (
          <div className="mb-4">
            <div className="flex justify-between text-sm text-gray-600 mb-2">
              <span>{progress.message}</span>
              <span>{progressPercentage}%</span>
            </div>
            
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div 
                className={`h-2 rounded-full transition-all duration-300 ${getProgressColor(progress.stage)}`}
                style={{ width: `${progressPercentage}%` }}
              />
            </div>
            
            {/* Progress Details */}
            <div className="flex justify-between text-xs text-gray-500 mt-2">
              <div>
                {progress.processedFiles !== undefined && progress.totalFiles !== undefined && (
                  <span>
                    {progress.processedFiles} of {progress.totalFiles} files
                  </span>
                )}
              </div>
              
              <div>
                Stage: {progress.stage}
              </div>
            </div>
            
            {/* Current File */}
            {progress.currentFile && (
              <div className="mt-2 p-2 bg-blue-50 rounded text-xs">
                <div className="flex items-center gap-2">
                  <DocumentIcon className="h-3 w-3 text-blue-500" />
                  <span className="text-blue-700 font-mono">{progress.currentFile}</span>
                </div>
              </div>
            )}
          </div>
        )}
        
        {/* File Statistics */}
        {fileStats.fileCount > 0 && (
          <div className="grid grid-cols-2 gap-4 pt-4 border-t">
            <div className="text-center">
              <div className="text-lg font-bold text-gray-900">{fileStats.fileCount}</div>
              <div className="text-xs text-gray-500">Files Generated</div>
            </div>
            
            <div className="text-center">
              <div className="text-lg font-bold text-gray-900">
                {fileStats.totalSize > 1024 
                  ? `${(fileStats.totalSize / 1024).toFixed(1)}KB`
                  : `${fileStats.totalSize}B`
                }
              </div>
              <div className="text-xs text-gray-500">Total Size</div>
            </div>
          </div>
        )}
        
        {/* Connection Issues */}
        {connectionState.reconnectAttempts > 0 && (
          <div className="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded">
            <div className="flex items-center gap-2">
              <ExclamationTriangleIcon className="h-4 w-4 text-yellow-600" />
              <div className="text-sm text-yellow-800">
                Reconnection attempt {connectionState.reconnectAttempts}
                {connectionState.lastReconnectAttempt && (
                  <span className="ml-2 text-xs">
                    ({formatDuration(Date.now() - connectionState.lastReconnectAttempt)} ago)
                  </span>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}