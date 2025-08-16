import { useState } from 'react'
import { XMarkIcon, ChevronRightIcon, CheckCircleIcon } from '@heroicons/react/24/outline'
import { ProjectTypeIcons } from '../icons/ProjectIcons'

interface QuickStartGuideProps {
  isOpen: boolean
  onClose: () => void
  onProjectTypeSelect: (type: string) => void
  onModeSelect: (mode: 'basic' | 'advanced') => void
}

const quickStartSteps = [
  {
    id: 1,
    title: 'Choose Your Project Type',
    description: 'Start with the type of Go application you want to build',
    recommendations: [
      { type: 'cli', label: 'CLI Application', good_for: 'Command-line tools, scripts, utilities' },
      { type: 'web-api', label: 'Web API', good_for: 'REST APIs, microservices, web backends' },
      { type: 'library', label: 'Library', good_for: 'Reusable packages, shared code' },
    ]
  },
  {
    id: 2,
    title: 'Select Complexity Level',
    description: 'Choose the right complexity for your needs',
    content: {
      basic: 'Quick Start - Essential options only, perfect for beginners',
      advanced: 'Full Control - All customization options available'
    }
  },
  {
    id: 3,
    title: 'Configure Essentials',
    description: 'Set up your project name and module path',
    tips: [
      'Project name: Use lowercase, numbers, and hyphens',
      'Module path: Follow Go conventions (github.com/user/project)',
      'Framework: Gin is a safe default for web APIs'
    ]
  },
  {
    id: 4,
    title: 'Review & Generate',
    description: 'Preview your project structure before generation',
    features: [
      'File structure preview',
      'Dependency overview',
      'Configuration summary'
    ]
  },
  {
    id: 5,
    title: 'You\'re Ready!',
    description: 'Your Go project is ready for development',
    next_steps: [
      'Run go mod download to install dependencies',
      'Start developing with your chosen architecture',
      'Use the generated README for project-specific instructions'
    ]
  }
]

export default function QuickStartGuide({ 
  isOpen, 
  onClose, 
  onProjectTypeSelect, 
  onModeSelect 
}: QuickStartGuideProps) {
  const [currentStep, setCurrentStep] = useState(1)
  const [completedSteps, setCompletedSteps] = useState<number[]>([])

  if (!isOpen) return null

  const handleStepComplete = (stepId: number) => {
    if (!completedSteps.includes(stepId)) {
      setCompletedSteps([...completedSteps, stepId])
    }
    if (stepId < quickStartSteps.length) {
      setCurrentStep(stepId + 1)
    }
  }

  const handleProjectTypeSelection = (type: string) => {
    onProjectTypeSelect(type)
    handleStepComplete(1)
  }

  const handleModeSelection = (mode: 'basic' | 'advanced') => {
    onModeSelect(mode)
    handleStepComplete(2)
  }

  const currentStepData = quickStartSteps.find(step => step.id === currentStep)

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <div>
            <h2 className="text-xl font-semibold text-gray-900">Quick Start Guide</h2>
            <p className="text-sm text-gray-600 mt-1">Get started with go-starter in 5 easy steps</p>
          </div>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-600 transition-colors"
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        {/* Progress Bar */}
        <div className="px-6 py-4 bg-gray-50">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-700">
              Step {currentStep} of {quickStartSteps.length}
            </span>
            <span className="text-sm text-gray-500">
              {Math.round((currentStep / quickStartSteps.length) * 100)}% complete
            </span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className="bg-primary-600 h-2 rounded-full transition-all duration-300"
              style={{ width: `${(currentStep / quickStartSteps.length) * 100}%` }}
            />
          </div>
        </div>

        {/* Content */}
        <div className="p-6">
          {currentStepData && (
            <div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                {currentStepData.title}
              </h3>
              <p className="text-gray-600 mb-4">
                {currentStepData.description}
              </p>

              {/* Step 1: Project Type Selection */}
              {currentStep === 1 && currentStepData.recommendations && (
                <div className="space-y-3">
                  {currentStepData.recommendations.map((rec) => {
                    const Icon = ProjectTypeIcons[rec.type as keyof typeof ProjectTypeIcons]
                    return (
                      <button
                        key={rec.type}
                        onClick={() => handleProjectTypeSelection(rec.type)}
                        className="w-full p-4 border border-gray-200 rounded-lg hover:border-primary-300 hover:bg-primary-50 transition-all text-left"
                      >
                        <div className="flex items-start gap-3">
                          {Icon && <Icon className="w-6 h-6 text-primary-600 mt-1" />}
                          <div className="flex-1">
                            <h4 className="font-medium text-gray-900">{rec.label}</h4>
                            <p className="text-sm text-gray-600 mt-1">{rec.good_for}</p>
                          </div>
                          <ChevronRightIcon className="w-5 h-5 text-gray-500" />
                        </div>
                      </button>
                    )
                  })}
                </div>
              )}

              {/* Step 2: Complexity Selection */}
              {currentStep === 2 && currentStepData.content && (
                <div className="space-y-3">
                  <button
                    onClick={() => handleModeSelection('basic')}
                    className="w-full p-4 border border-gray-200 rounded-lg hover:border-primary-300 hover:bg-primary-50 transition-all text-left"
                  >
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="font-medium text-gray-900">Quick Start Mode</h4>
                        <p className="text-sm text-gray-600 mt-1">{currentStepData.content.basic}</p>
                      </div>
                      <ChevronRightIcon className="w-5 h-5 text-gray-500" />
                    </div>
                  </button>
                  <button
                    onClick={() => handleModeSelection('advanced')}
                    className="w-full p-4 border border-gray-200 rounded-lg hover:border-primary-300 hover:bg-primary-50 transition-all text-left"
                  >
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="font-medium text-gray-900">Advanced Mode</h4>
                        <p className="text-sm text-gray-600 mt-1">{currentStepData.content.advanced}</p>
                      </div>
                      <ChevronRightIcon className="w-5 h-5 text-gray-500" />
                    </div>
                  </button>
                </div>
              )}

              {/* Step 3: Configuration Tips */}
              {currentStep === 3 && currentStepData.tips && (
                <div className="space-y-3">
                  {currentStepData.tips.map((tip, index) => (
                    <div key={index} className="flex items-start gap-3 p-3 bg-blue-50 rounded-lg">
                      <CheckCircleIcon className="w-5 h-5 text-blue-600 mt-0.5" />
                      <p className="text-sm text-blue-800">{tip}</p>
                    </div>
                  ))}
                  <button
                    onClick={() => handleStepComplete(3)}
                    className="w-full mt-4 bg-primary-600 text-white py-2 px-4 rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    Continue to Review
                  </button>
                </div>
              )}

              {/* Step 4: Review Features */}
              {currentStep === 4 && currentStepData.features && (
                <div className="space-y-3">
                  {currentStepData.features.map((feature, index) => (
                    <div key={index} className="flex items-center gap-3 p-3 bg-green-50 rounded-lg">
                      <CheckCircleIcon className="w-5 h-5 text-green-600" />
                      <p className="text-sm text-green-800">{feature}</p>
                    </div>
                  ))}
                  <button
                    onClick={() => handleStepComplete(4)}
                    className="w-full mt-4 bg-primary-600 text-white py-2 px-4 rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    Generate Project
                  </button>
                </div>
              )}

              {/* Step 5: Completion */}
              {currentStep === 5 && currentStepData.next_steps && (
                <div className="space-y-3">
                  <div className="text-center mb-6">
                    <CheckCircleIcon className="w-16 h-16 text-green-600 mx-auto mb-4" />
                    <h4 className="text-lg font-semibold text-gray-900">Congratulations!</h4>
                    <p className="text-gray-600">Your Go project has been generated successfully.</p>
                  </div>
                  <div className="space-y-2">
                    <h5 className="font-medium text-gray-900">Next Steps:</h5>
                    {currentStepData.next_steps.map((step, index) => (
                      <div key={index} className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg">
                        <span className="flex-shrink-0 w-6 h-6 bg-primary-600 text-white rounded-full flex items-center justify-center text-xs font-medium">
                          {index + 1}
                        </span>
                        <p className="text-sm text-gray-700">{step}</p>
                      </div>
                    ))}
                  </div>
                  <button
                    onClick={onClose}
                    className="w-full mt-6 bg-primary-600 text-white py-2 px-4 rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    Get Started Coding!
                  </button>
                </div>
              )}
            </div>
          )}
        </div>

        {/* Navigation */}
        <div className="flex items-center justify-between p-6 border-t border-gray-200 bg-gray-50">
          <button
            onClick={() => setCurrentStep(Math.max(1, currentStep - 1))}
            disabled={currentStep === 1}
            className="text-sm text-gray-600 hover:text-gray-800 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            ← Previous
          </button>
          <div className="flex gap-2">
            {quickStartSteps.map((step) => (
              <button
                key={step.id}
                onClick={() => setCurrentStep(step.id)}
                className={`w-2 h-2 rounded-full transition-colors ${
                  step.id === currentStep 
                    ? 'bg-primary-600' 
                    : completedSteps.includes(step.id)
                      ? 'bg-green-600'
                      : 'bg-gray-300'
                }`}
              />
            ))}
          </div>
          <button
            onClick={() => setCurrentStep(Math.min(quickStartSteps.length, currentStep + 1))}
            disabled={currentStep === quickStartSteps.length}
            className="text-sm text-gray-600 hover:text-gray-800 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next →
          </button>
        </div>
      </div>
    </div>
  )
}