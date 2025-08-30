import React, { useState } from 'react'
import { ArrowPathIcon, EyeIcon, CodeBracketIcon } from '@heroicons/react/24/outline'
import type { ProjectConfig, WSFileContent } from '../../types'
import FileExplorerPanel from './FileExplorerPanel'
import CodePreview from './CodePreview'
import ProgressTracker from './ProgressTracker'
import { useRealtimePreview } from '../../hooks/useRealtimePreview'

interface PreviewPanelProps {
  config: ProjectConfig
  className?: string
  enableRealtimePreview?: boolean
}

export default function PreviewPanel({ 
  config, 
  className = '', 
  enableRealtimePreview = true 
}: PreviewPanelProps) {
  const [selectedFile, setSelectedFile] = useState<WSFileContent | null>(null)
  const [viewMode, setViewMode] = useState<'split' | 'tree' | 'code'>('split')
  
  const { 
    previewState, 
    connectionState,
    startPreview,
    clearPreview 
  } = useRealtimePreview({
    onPreviewComplete: (data) => {
      console.log('Preview complete:', data)
    },
    onError: (error) => {
      console.error('Preview error:', error)
    }
  })
  
  const handleFileSelect = (file: WSFileContent | null) => {
    setSelectedFile(file)
  }
  
  const handleStartPreview = () => {
    startPreview(config)
  }
  
  const handleRefreshPreview = () => {
    clearPreview()
    setTimeout(() => startPreview(config), 100)
  }
  
  return (
    <div className={`h-full bg-white border rounded-lg shadow-sm overflow-hidden flex flex-col ${className}`}>
      {/* Header */}
      <div className="border-b bg-gray-50 p-4 flex-shrink-0">
        <div className="flex items-center justify-between mb-3">
          <div>
            <h2 className="text-lg font-semibold text-gray-900">Project Preview</h2>
            <p className="text-sm text-gray-600 mt-1">
              {enableRealtimePreview 
                ? `Real-time preview for ${config.projectName}` 
                : `Static preview for ${config.projectName}`
              }
            </p>
          </div>
          
          <div className="flex items-center gap-2">
            {/* View Mode Toggle */}
            <div className="flex border rounded-lg">
              <button
                onClick={() => setViewMode('split')}
                className={`px-3 py-1 text-xs rounded-l-lg transition-colors ${
                  viewMode === 'split' 
                    ? 'bg-blue-600 text-white' 
                    : 'bg-white hover:bg-gray-50'
                }`}
              >
                Split
              </button>
              <button
                onClick={() => setViewMode('tree')}
                className={`px-3 py-1 text-xs border-l transition-colors ${
                  viewMode === 'tree' 
                    ? 'bg-blue-600 text-white' 
                    : 'bg-white hover:bg-gray-50'
                }`}
              >
                <EyeIcon className="h-3 w-3" />
              </button>
              <button
                onClick={() => setViewMode('code')}
                className={`px-3 py-1 text-xs rounded-r-lg border-l transition-colors ${
                  viewMode === 'code' 
                    ? 'bg-blue-600 text-white' 
                    : 'bg-white hover:bg-gray-50'
                }`}
              >
                <CodeBracketIcon className="h-3 w-3" />
              </button>
            </div>
            
            {enableRealtimePreview && (
              <>
                {!previewState.isGenerating ? (
                  <button 
                    onClick={handleStartPreview}
                    disabled={!connectionState.connected}
                    className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <EyeIcon className="h-4 w-4" />
                    Start Preview
                  </button>
                ) : (
                  <button 
                    onClick={handleRefreshPreview}
                    className="flex items-center gap-2 px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
                  >
                    <ArrowPathIcon className="h-4 w-4" />
                    Refresh
                  </button>
                )}
              </>
            )}
          </div>
        </div>
        
        {/* Progress Tracker - Compact Mode */}
        {enableRealtimePreview && (
          <ProgressTracker 
            connectionState={connectionState}
            previewState={previewState}
            variant="compact"
          />
        )}
      </div>
      
      {/* Content Area */}
      <div className="flex-1 flex overflow-hidden">
        {viewMode === 'split' && (
          <>
            {/* File Explorer */}
            <div className="w-1/2 border-r">
              <FileExplorerPanel 
                config={config}
                enableRealtimePreview={enableRealtimePreview}
                onFileSelect={handleFileSelect}
                onPreviewStart={() => {}}
                onPreviewComplete={() => {}}
              />
            </div>
            
            {/* Code Preview */}
            <div className="w-1/2">
              <CodePreview 
                file={selectedFile}
                isLoading={previewState.isGenerating && !selectedFile}
              />
            </div>
          </>
        )}
        
        {viewMode === 'tree' && (
          <div className="flex-1">
            <FileExplorerPanel 
              config={config}
              enableRealtimePreview={enableRealtimePreview}
              onFileSelect={handleFileSelect}
              onPreviewStart={() => {}}
              onPreviewComplete={() => {}}
            />
          </div>
        )}
        
        {viewMode === 'code' && (
          <div className="flex-1">
            <CodePreview 
              file={selectedFile}
              isLoading={previewState.isGenerating && !selectedFile}
            />
          </div>
        )}
      </div>
      
      {/* Footer - Detailed Progress when generating */}
      {enableRealtimePreview && previewState.isGenerating && (
        <div className="border-t bg-gray-50 p-4 flex-shrink-0">
          <ProgressTracker 
            connectionState={connectionState}
            previewState={previewState}
            variant="detailed"
            className="max-h-32 overflow-auto"
          />
        </div>
      )}
    </div>
  )
}