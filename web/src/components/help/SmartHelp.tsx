import { useState } from 'react'
import { XMarkIcon, LightBulbIcon, ExclamationTriangleIcon } from '@heroicons/react/24/outline'
import type { ProjectType, Architecture, Framework } from '../../types'

interface SmartHelpProps {
  projectType: ProjectType
  architecture: Architecture
  framework: Framework
  disclosureMode: 'basic' | 'advanced'
}

interface Suggestion {
  id: string
  type: 'tip' | 'warning' | 'recommendation'
  title: string
  content: string
  action?: {
    label: string
    onClick: () => void
  }
}

export default function SmartHelp({ 
  projectType, 
  architecture, 
  framework, 
  disclosureMode 
}: SmartHelpProps) {
  const [dismissedSuggestions, setDismissedSuggestions] = useState<string[]>([])

  const generateSuggestions = (): Suggestion[] => {
    const suggestions: Suggestion[] = []

    // Architecture complexity warnings
    if (disclosureMode === 'basic' && (architecture === 'hexagonal' || architecture === 'ddd')) {
      suggestions.push({
        id: 'complex-architecture-basic-mode',
        type: 'warning',
        title: 'Complex Architecture in Basic Mode',
        content: `${architecture === 'hexagonal' ? 'Hexagonal' : 'DDD'} architecture is advanced. Consider switching to Advanced mode for full configuration options.`,
        action: {
          label: 'Switch to Advanced',
          onClick: () => console.log('Switch to advanced mode')
        }
      })
    }

    // Framework + Project Type compatibility
    if (projectType === 'cli' && (framework === 'gin' || framework === 'echo')) {
      suggestions.push({
        id: 'web-framework-cli-project',
        type: 'warning',
        title: 'Web Framework for CLI Project',
        content: `${framework} is a web framework, but you're building a CLI application. Consider using the Cobra framework instead.`,
        action: {
          label: 'Use Cobra',
          onClick: () => console.log('Switch to Cobra')
        }
      })
    }

    // Performance recommendations
    if (framework === 'fiber' && architecture === 'standard') {
      suggestions.push({
        id: 'fiber-performance-tip',
        type: 'tip',
        title: 'Fiber Performance Tip',
        content: 'Fiber excels with high-concurrency workloads. Consider using Clean Architecture to better organize your high-performance application.',
        action: {
          label: 'Use Clean Architecture',
          onClick: () => console.log('Switch to Clean')
        }
      })
    }

    // Progressive complexity suggestions
    if (projectType === 'microservice' && disclosureMode === 'basic') {
      suggestions.push({
        id: 'microservice-needs-advanced',
        type: 'recommendation',
        title: 'Microservice Configuration',
        content: 'Microservices typically need database, authentication, and deployment configuration. Advanced mode provides these options.',
        action: {
          label: 'Enable Advanced Mode',
          onClick: () => console.log('Enable advanced mode')
        }
      })
    }

    // Best practice recommendations
    if (projectType === 'web-api' && architecture === 'standard' && framework === 'gin') {
      suggestions.push({
        id: 'gin-clean-architecture',
        type: 'tip',
        title: 'API Architecture Suggestion',
        content: 'For maintainable web APIs, consider Clean Architecture. It provides better separation of concerns and testability.',
        action: {
          label: 'Try Clean Architecture',
          onClick: () => console.log('Switch to Clean')
        }
      })
    }

    // Learning path suggestions
    if (disclosureMode === 'basic' && architecture === 'standard') {
      suggestions.push({
        id: 'learning-path',
        type: 'tip',
        title: 'Ready for More?',
        content: 'Once you\'re comfortable with the standard structure, explore Clean Architecture for better code organization and testing.',
      })
    }

    // Filter out dismissed suggestions
    return suggestions.filter(suggestion => !dismissedSuggestions.includes(suggestion.id))
  }

  const suggestions = generateSuggestions()

  const dismissSuggestion = (id: string) => {
    setDismissedSuggestions(prev => [...prev, id])
  }

  const getIconForType = (type: Suggestion['type']) => {
    switch (type) {
      case 'warning':
        return <ExclamationTriangleIcon className="w-5 h-5 text-amber-600" />
      case 'tip':
      case 'recommendation':
        return <LightBulbIcon className="w-5 h-5 text-blue-600" />
      default:
        return <LightBulbIcon className="w-5 h-5 text-blue-600" />
    }
  }

  const getStylesForType = (type: Suggestion['type']) => {
    switch (type) {
      case 'warning':
        return 'bg-amber-50 border-amber-200'
      case 'tip':
        return 'bg-blue-50 border-blue-200'
      case 'recommendation':
        return 'bg-green-50 border-green-200'
      default:
        return 'bg-blue-50 border-blue-200'
    }
  }

  if (suggestions.length === 0) {
    return null
  }

  return (
    <div className="space-y-3">
      {suggestions.map((suggestion) => (
        <div
          key={suggestion.id}
          className={`p-4 border rounded-lg ${getStylesForType(suggestion.type)}`}
        >
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 mt-0.5">
              {getIconForType(suggestion.type)}
            </div>
            <div className="flex-1 min-w-0">
              <h4 className="text-sm font-medium text-gray-900 mb-1">
                {suggestion.title}
              </h4>
              <p className="text-sm text-gray-700 mb-3">
                {suggestion.content}
              </p>
              <div className="flex items-center gap-2">
                {suggestion.action && (
                  <button
                    onClick={suggestion.action.onClick}
                    className="text-xs font-medium text-primary-700 hover:text-primary-800 underline"
                  >
                    {suggestion.action.label} →
                  </button>
                )}
                <button
                  onClick={() => dismissSuggestion(suggestion.id)}
                  className="text-xs text-gray-500 hover:text-gray-700"
                >
                  Dismiss
                </button>
              </div>
            </div>
            <button
              onClick={() => dismissSuggestion(suggestion.id)}
              className="flex-shrink-0 text-gray-500 hover:text-gray-600"
            >
              <XMarkIcon className="w-4 h-4" />
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}