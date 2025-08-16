import { LightBulbIcon } from '@heroicons/react/24/outline'
import type { ProjectType } from '../../types'

interface ModeSuggestionProps {
  projectType: ProjectType
  currentMode: 'basic' | 'advanced'
  onSwitchToAdvanced: () => void
}

const ADVANCED_PROJECT_TYPES: ProjectType[] = ['microservice', 'event-driven', 'monolith']

export default function ModeSuggestion({ 
  projectType, 
  currentMode, 
  onSwitchToAdvanced 
}: ModeSuggestionProps) {
  // Only show suggestion if using basic mode with an advanced project type
  if (currentMode === 'advanced' || !ADVANCED_PROJECT_TYPES.includes(projectType)) {
    return null
  }

  return (
    <div className="mb-4 p-3 bg-amber-50 border border-amber-200 rounded-lg">
      <div className="flex items-start space-x-2">
        <LightBulbIcon className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
        <div className="flex-1">
          <p className="text-sm font-medium text-amber-800">
            Advanced mode recommended
          </p>
          <p className="text-xs text-amber-700 mt-1">
            {projectType === 'microservice' && 'Microservices often need database, authentication, and deployment configurations.'}
            {projectType === 'event-driven' && 'Event-driven architectures typically require message brokers and advanced patterns.'}
            {projectType === 'monolith' && 'Monolithic applications benefit from database and authentication setup.'}
          </p>
          <button
            onClick={onSwitchToAdvanced}
            className="mt-2 text-xs font-medium text-amber-700 hover:text-amber-800 underline"
          >
            Switch to Advanced Mode →
          </button>
        </div>
      </div>
    </div>
  )
}