import { useState } from 'react'
import { PlayIcon, CheckCircleIcon, ExclamationTriangleIcon, ArrowPathIcon } from '@heroicons/react/20/solid'

interface GenerationStatus {
  status: 'idle' | 'generating' | 'completed' | 'error'
  progress: number
  filesGenerated: number
  totalFiles: number
  currentFile?: string
  error?: string
}

export default function PreviewPanel() {
  const [status, setStatus] = useState<GenerationStatus>({
    status: 'idle',
    progress: 0,
    filesGenerated: 0,
    totalFiles: 0,
  })

  const startGeneration = async () => {
    setStatus({
      status: 'generating',
      progress: 0,
      filesGenerated: 0,
      totalFiles: 25,
    })

    // Simulate generation progress
    for (let i = 1; i <= 25; i++) {
      await new Promise(resolve => setTimeout(resolve, 100))
      setStatus(prev => ({
        ...prev,
        progress: (i / 25) * 100,
        filesGenerated: i,
        currentFile: `Generating ${getRandomFileName()}...`,
      }))
    }

    setStatus(prev => ({
      ...prev,
      status: 'completed',
      currentFile: undefined,
    }))
  }

  const getRandomFileName = () => {
    const files = [
      'main.go',
      'internal/config/config.go',
      'internal/handlers/user.go',
      'internal/middleware/auth.go',
      'cmd/server/main.go',
      'internal/database/connection.go',
      'go.mod',
      'Dockerfile',
      'README.md',
      'Makefile',
    ]
    return files[Math.floor(Math.random() * files.length)]
  }

  const getStatusIcon = () => {
    switch (status.status) {
      case 'generating':
        return <ArrowPathIcon className="w-5 h-5 text-blue-500 animate-spin" />
      case 'completed':
        return <CheckCircleIcon className="w-5 h-5 text-green-500" />
      case 'error':
        return <ExclamationTriangleIcon className="w-5 h-5 text-red-500" />
      default:
        return <PlayIcon className="w-5 h-5 text-gray-400" />
    }
  }

  const getStatusText = () => {
    switch (status.status) {
      case 'generating':
        return 'Generating project...'
      case 'completed':
        return 'Generation completed!'
      case 'error':
        return 'Generation failed'
      default:
        return 'Ready to generate'
    }
  }

  return (
    <div className="card h-full">
      <div className="mb-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-2">Live Preview</h2>
        <p className="text-sm text-gray-600">Real-time project generation preview</p>
      </div>

      {/* Generation Status */}
      <div className="mb-6">
        <div className="flex items-center space-x-3 mb-4">
          {getStatusIcon()}
          <div className="flex-1">
            <p className="text-sm font-medium text-gray-900">{getStatusText()}</p>
            {status.currentFile && (
              <p className="text-xs text-gray-500">{status.currentFile}</p>
            )}
          </div>
        </div>

        {/* Progress Bar */}
        {status.status === 'generating' && (
          <div className="mb-4">
            <div className="flex items-center justify-between text-xs text-gray-600 mb-1">
              <span>Progress</span>
              <span>{status.filesGenerated}/{status.totalFiles} files</span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div 
                className="bg-blue-600 h-2 rounded-full transition-all duration-200"
                style={{ width: `${status.progress}%` }}
              />
            </div>
          </div>
        )}

        {/* Statistics */}
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="bg-gray-50 rounded-lg p-3">
            <p className="text-xs font-medium text-gray-500 uppercase tracking-wide">Files Generated</p>
            <p className="text-2xl font-bold text-gray-900">{status.filesGenerated}</p>
          </div>
          <div className="bg-gray-50 rounded-lg p-3">
            <p className="text-xs font-medium text-gray-500 uppercase tracking-wide">Progress</p>
            <p className="text-2xl font-bold text-gray-900">{Math.round(status.progress)}%</p>
          </div>
        </div>
      </div>

      {/* Preview Content */}
      <div className="flex-1 bg-gray-50 rounded-lg p-4 min-h-[300px]">
        {status.status === 'idle' && (
          <div className="flex flex-col items-center justify-center h-full text-center">
            <div className="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center mb-4">
              <PlayIcon className="w-8 h-8 text-gray-400" />
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">Ready to Generate</h3>
            <p className="text-sm text-gray-500 mb-4">
              Configure your project settings and click generate to see a live preview.
            </p>
            <button 
              onClick={startGeneration}
              className="btn-primary"
            >
              Start Generation Preview
            </button>
          </div>
        )}

        {status.status === 'generating' && (
          <div className="space-y-3">
            <div className="bg-white rounded p-3 border-l-4 border-blue-500">
              <p className="text-sm font-medium text-gray-900">Creating project structure...</p>
              <p className="text-xs text-gray-500 mt-1">Setting up directories and base files</p>
            </div>
            <div className="bg-white rounded p-3 border-l-4 border-blue-500">
              <p className="text-sm font-medium text-gray-900">Generating Go modules...</p>
              <p className="text-xs text-gray-500 mt-1">Creating go.mod and dependency management</p>
            </div>
            <div className="bg-white rounded p-3 border-l-4 border-blue-500">
              <p className="text-sm font-medium text-gray-900">Building application code...</p>
              <p className="text-xs text-gray-500 mt-1">Generating handlers, middleware, and business logic</p>
            </div>
          </div>
        )}

        {status.status === 'completed' && (
          <div className="space-y-3">
            <div className="bg-green-50 border border-green-200 rounded-lg p-4">
              <div className="flex items-center">
                <CheckCircleIcon className="w-5 h-5 text-green-400 mr-2" />
                <h3 className="text-sm font-medium text-green-800">Generation Completed Successfully!</h3>
              </div>
              <p className="text-xs text-green-700 mt-1">
                Your Go project has been generated with {status.filesGenerated} files.
              </p>
            </div>
            
            <div className="bg-white rounded-lg border border-gray-200 p-4">
              <h4 className="text-sm font-medium text-gray-900 mb-2">Generated Files Preview:</h4>
              <div className="space-y-1 text-xs font-mono text-gray-600">
                <div className="flex items-center space-x-2">
                  <span className="text-gray-400">📁</span>
                  <span>cmd/server/main.go</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-gray-400">📁</span>
                  <span>internal/config/config.go</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-gray-400">📁</span>
                  <span>internal/handlers/</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-gray-400">📄</span>
                  <span>go.mod</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span className="text-gray-400">📄</span>
                  <span>README.md</span>
                </div>
                <div className="text-gray-400 text-center mt-2">
                  ... and {status.filesGenerated - 5} more files
                </div>
              </div>
            </div>
          </div>
        )}

        {status.status === 'error' && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4">
            <div className="flex items-center">
              <ExclamationTriangleIcon className="w-5 h-5 text-red-400 mr-2" />
              <h3 className="text-sm font-medium text-red-800">Generation Failed</h3>
            </div>
            <p className="text-xs text-red-700 mt-1">
              {status.error || 'An error occurred during project generation.'}
            </p>
          </div>
        )}
      </div>
    </div>
  )
}