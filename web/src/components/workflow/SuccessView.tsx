/**
 * Success View
 * Displays successful generation results with download options
 */

import React from 'react'
import type { Blueprint, ProjectConfig } from '../../services/api'
import type { GenerationResult } from '../../hooks/useGenerationWorkflow'
import DownloadSuccessFlow from '../generation/DownloadSuccessFlow'

interface SuccessViewProps {
  result: GenerationResult
  blueprint: Blueprint | null
  config: ProjectConfig
  onDownload: () => void
  onNewProject: () => void
  onShowPreview?: () => void
  downloading: boolean
  downloadError?: string | null
  className?: string
}

export default function SuccessView({
  result,
  blueprint,
  config,
  onDownload,
  onNewProject,
  onShowPreview,
  downloading,
  downloadError,
  className = ''
}: SuccessViewProps) {
  return (
    <div className={`h-full overflow-y-auto p-6 ${className}`}>
      <DownloadSuccessFlow
        projectId={result.projectId}
        projectName={config.projectName}
        blueprint={blueprint}
        config={config}
        fileCount={result.fileCount}
        estimatedSize={result.estimatedSize}
        generationTime={result.generationTime}
        onNewProject={onNewProject}
        onShowPreview={onShowPreview}
      />
    </div>
  )
}