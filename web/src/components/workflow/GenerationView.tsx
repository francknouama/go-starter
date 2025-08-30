/**
 * Generation View  
 * Displays the project generation workflow with real-time progress
 */

import React from 'react'
import type { Blueprint, ProjectConfig } from '../../services/api'
import type { WorkflowSession, GenerationResult } from '../../hooks/useGenerationWorkflow'
import GenerationWorkflow from '../generation/GenerationWorkflow'

interface GenerationViewProps {
  session: WorkflowSession | null
  blueprint: Blueprint | null
  config: ProjectConfig
  onCancel: () => void
  onComplete: (result: GenerationResult) => void
  onError: (error: Error) => void
  wsConnected: boolean
  progressData: any
  className?: string
}

export default function GenerationView({
  session,
  blueprint,
  config,
  onCancel,
  onComplete,
  onError,
  wsConnected,
  progressData,
  className = ''
}: GenerationViewProps) {
  return (
    <div className={`h-full p-6 ${className}`}>
      <GenerationWorkflow
        blueprint={blueprint}
        config={config}
        onComplete={onComplete}
        onCancel={onCancel}
        onReset={onCancel}
      />
    </div>
  )
}