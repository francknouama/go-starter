import { useState } from 'react'
import { PlayIcon, CheckCircleIcon, ExclamationTriangleIcon, ArrowPathIcon } from '@heroicons/react/20/solid'

interface PreviewPanelProps {
  preview?: any
  generation?: any
  generationProgress?: any
  onStartGeneration?: () => void
}

export default function PreviewPanel({ preview, generation, generationProgress, onStartGeneration }: PreviewPanelProps) {
  const getGenerationStatus = () => {
    if (generationProgress) {
      return {
        status: generationProgress.completed ? 'completed' : 'generating',
        progress: generationProgress.progress || 0,
        filesGenerated: Math.floor((generationProgress.progress || 0) / 100 * (preview?.fileCount || 0)),
        totalFiles: preview?.fileCount || 0,
        currentFile: generationProgress.message,
        error: generationProgress.error,
      }
    }
    
    if (generation) {
      return {
        status: 'completed',
        progress: 100,
        filesGenerated: generation.fileCount || 0,
        totalFiles: generation.fileCount || 0,
        currentFile: undefined,
        error: undefined,
      }
    }
    
    if (preview) {
      return {
        status: 'idle',
        progress: 0,
        filesGenerated: 0,
        totalFiles: preview.fileCount || 0,
        currentFile: undefined,
        error: undefined,
      }
    }
    
    return {
      status: 'idle',
      progress: 0,
      filesGenerated: 0,
      totalFiles: 0,
      currentFile: undefined,
      error: undefined,
    }
  }
  
  const status = getGenerationStatus()

  const renderFileStructure = () => {
    if (!preview?.fileStructure) return null
    
    const renderNode = (node: any, depth = 0) => {
      const indent = depth * 16
      const icon = node.type === 'directory' ? '📁' : '📄'
      
      return (
        <div key={node.path} style={{ marginLeft: indent }}>
          <div className="flex items-center space-x-2 py-1">
            <span className="text-gray-500">{icon}</span>
            <span className="text-xs font-mono text-gray-700">{node.name}</span>
          </div>
          {node.children && node.children.map((child: any) => renderNode(child, depth + 1))}
        </div>
      )
    }
    
    return (
      <div className="space-y-1">
        {preview.fileStructure.map((node: any) => renderNode(node))}
      </div>
    )
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
        return <PlayIcon className="w-5 h-5 text-gray-500" />
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
    <div className="bg-white/70 backdrop-blur-lg rounded-2xl shadow-xl border border-white/30 p-4 md:p-6 h-full">
      <div className="mb-4 md:mb-6">
        <h2 className="text-base md:text-lg font-semibold text-gray-900 mb-1 md:mb-2">Live Preview</h2>
        <p className="text-xs md:text-sm text-gray-600">Real-time project generation preview</p>
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
          <div className="bg-gradient-to-br from-blue-50/80 to-indigo-50/80 backdrop-blur-md rounded-xl p-3 border border-white/20">
            <p className="text-xs font-medium text-gray-500 uppercase tracking-wide">Files Generated</p>
            <p className="text-2xl font-bold text-gray-900">{status.filesGenerated}</p>
          </div>
          <div className="bg-gradient-to-br from-purple-50/80 to-pink-50/80 backdrop-blur-md rounded-xl p-3 border border-white/20">
            <p className="text-xs font-medium text-gray-500 uppercase tracking-wide">Progress</p>
            <p className="text-2xl font-bold text-gray-900">{Math.round(status.progress)}%</p>
          </div>
        </div>
      </div>

      {/* Preview Content */}
      <div className="flex-1 bg-gradient-to-br from-white/40 to-gray-50/40 backdrop-blur-md rounded-xl p-4 min-h-[300px] border border-white/20">
        {status.status === 'idle' && (
          <div className="flex flex-col items-center justify-center h-full text-center">
            <div className="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center mb-4">
              <PlayIcon className="w-8 h-8 text-gray-500" />
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">Ready to Generate</h3>
            <p className="text-sm text-gray-500 mb-4">
              Configure your project settings and click generate to see a live preview.
            </p>
            <button 
              onClick={onStartGeneration}
              className="px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-xl hover:from-blue-700 hover:to-purple-700 transition-all duration-200 font-medium shadow-lg transform hover:scale-105"
              disabled={!onStartGeneration}
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
              {generation?.downloadUrl && (
                <a 
                  href={generation.downloadUrl}
                  download
                  className="inline-block mt-2 text-xs text-green-700 underline hover:text-green-800"
                >
                  Download Project ZIP
                </a>
              )}
            </div>
            
            <div className="bg-white rounded-lg border border-gray-200 p-4 max-h-64 overflow-y-auto">
              <h4 className="text-sm font-medium text-gray-900 mb-2">Generated Files Preview:</h4>
              {preview ? renderFileStructure() : (
                <p className="text-xs text-gray-500">File structure not available</p>
              )}
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