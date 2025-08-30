/**
 * Step Indicator Component
 * Visual progress indicator for the generation workflow
 */

import React from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  CheckCircleIcon,
  ExclamationCircleIcon,
  ClockIcon,
  CogIcon
} from '@heroicons/react/24/outline'
import { 
  CheckCircleIcon as CheckCircleIconSolid,
  ExclamationCircleIcon as ExclamationCircleIconSolid
} from '@heroicons/react/24/solid'

interface WorkflowStep {
  id: string
  title: string
  description: string
  status: 'pending' | 'active' | 'completed' | 'error' | 'skipped'
  duration?: number
  error?: string
  substeps?: string[]
  currentSubstep?: number
}

interface StepIndicatorProps {
  steps: WorkflowStep[]
  currentStep: number
  estimatedDuration?: number
  startTime?: Date | null
  className?: string
}

export default function StepIndicator({
  steps,
  currentStep,
  estimatedDuration = 15,
  startTime,
  className = ''
}: StepIndicatorProps) {
  // Calculate overall progress
  const completedSteps = steps.filter(step => step.status === 'completed').length
  const totalSteps = steps.length
  const overallProgress = totalSteps > 0 ? (completedSteps / totalSteps) * 100 : 0

  // Calculate elapsed and remaining time
  const elapsedTime = startTime ? Math.floor((Date.now() - startTime.getTime()) / 1000) : 0
  const remainingTime = Math.max(0, estimatedDuration - elapsedTime)

  // Format time display
  const formatTime = (seconds: number): string => {
    if (seconds < 60) return `${seconds}s`
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = seconds % 60
    return `${minutes}m ${remainingSeconds}s`
  }

  // Get step icon
  const getStepIcon = (step: WorkflowStep, index: number) => {
    switch (step.status) {
      case 'completed':
        return <CheckCircleIconSolid className="w-6 h-6 text-green-600" />
      case 'error':
        return <ExclamationCircleIconSolid className="w-6 h-6 text-red-600" />
      case 'active':
        return (
          <div className="relative">
            <CogIcon className="w-6 h-6 text-blue-600 animate-spin" />
            <div className="absolute inset-0 rounded-full border-2 border-blue-200 animate-pulse"></div>
          </div>
        )
      case 'pending':
        return (
          <div className="w-6 h-6 rounded-full border-2 border-gray-300 bg-white flex items-center justify-center">
            <span className="text-xs font-medium text-gray-500">{index + 1}</span>
          </div>
        )
      case 'skipped':
        return (
          <div className="w-6 h-6 rounded-full border-2 border-gray-200 bg-gray-100 flex items-center justify-center">
            <span className="text-xs text-gray-400">—</span>
          </div>
        )
      default:
        return (
          <div className="w-6 h-6 rounded-full border-2 border-gray-300 bg-white flex items-center justify-center">
            <span className="text-xs font-medium text-gray-500">{index + 1}</span>
          </div>
        )
    }
  }

  // Get step color classes
  const getStepClasses = (step: WorkflowStep) => {
    switch (step.status) {
      case 'completed':
        return 'text-green-900 border-green-200 bg-green-50'
      case 'error':
        return 'text-red-900 border-red-200 bg-red-50'
      case 'active':
        return 'text-blue-900 border-blue-200 bg-blue-50 ring-2 ring-blue-500 ring-opacity-20'
      case 'pending':
        return 'text-gray-700 border-gray-200 bg-gray-50'
      case 'skipped':
        return 'text-gray-500 border-gray-200 bg-gray-100'
      default:
        return 'text-gray-700 border-gray-200 bg-gray-50'
    }
  }

  return (
    <div className={`w-full ${className}`}>
      {/* Progress Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-3">
            <h3 className="text-lg font-semibold text-gray-900">Generation Progress</h3>
            <div className="flex items-center gap-2 text-sm text-gray-600">
              <ClockIcon className="w-4 h-4" />
              {startTime ? (
                <span>
                  {elapsedTime > 0 && `${formatTime(elapsedTime)} elapsed`}
                  {remainingTime > 0 && elapsedTime > 0 && ' • '}
                  {remainingTime > 0 && `~${formatTime(remainingTime)} remaining`}
                </span>
              ) : (
                <span>Estimated: {formatTime(estimatedDuration)}</span>
              )}
            </div>
          </div>
          <div className="text-right">
            <div className="text-sm font-medium text-gray-900">
              {completedSteps} of {totalSteps} complete
            </div>
            <div className="text-xs text-gray-600">
              {overallProgress.toFixed(0)}% done
            </div>
          </div>
        </div>

        {/* Overall Progress Bar */}
        <div className="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
          <motion.div
            className="h-full bg-gradient-to-r from-blue-500 to-green-500 rounded-full"
            initial={{ width: 0 }}
            animate={{ width: `${overallProgress}%` }}
            transition={{ duration: 0.5, ease: 'easeOut' }}
          />
        </div>
      </div>

      {/* Step List - Desktop View */}
      <div className="hidden md:block">
        <div className="space-y-4">
          {steps.map((step, index) => (
            <motion.div
              key={step.id}
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: index * 0.1 }}
              className={`relative flex items-start gap-4 p-4 rounded-lg border transition-all duration-300 ${getStepClasses(step)}`}
            >
              {/* Step Icon */}
              <div className="flex-shrink-0 mt-1">
                {getStepIcon(step, index)}
              </div>

              {/* Step Content */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center justify-between mb-1">
                  <h4 className="text-sm font-medium truncate">
                    {step.title}
                  </h4>
                  {step.duration && (
                    <span className="text-xs text-gray-500 ml-2">
                      {formatTime(Math.floor(step.duration / 1000))}
                    </span>
                  )}
                </div>
                
                <p className="text-sm text-gray-600 mb-2">
                  {step.description}
                </p>

                {/* Error Message */}
                {step.error && (
                  <div className="text-sm text-red-600 mb-2 p-2 bg-red-100 rounded">
                    {step.error}
                  </div>
                )}

                {/* Substeps Progress */}
                {step.substeps && step.status === 'active' && (
                  <div className="mt-3">
                    <div className="flex items-center gap-2 mb-2">
                      <div className="text-xs font-medium text-gray-700">
                        Step {(step.currentSubstep || 0) + 1} of {step.substeps.length}:
                      </div>
                      <div className="text-xs text-gray-600">
                        {step.substeps[step.currentSubstep || 0]}
                      </div>
                    </div>
                    
                    <div className="w-full bg-gray-200 rounded-full h-1.5">
                      <motion.div
                        className="h-full bg-blue-400 rounded-full"
                        initial={{ width: 0 }}
                        animate={{ 
                          width: `${((step.currentSubstep || 0) + 1) / step.substeps.length * 100}%` 
                        }}
                        transition={{ duration: 0.3 }}
                      />
                    </div>
                  </div>
                )}
              </div>

              {/* Connection Line */}
              {index < steps.length - 1 && (
                <div className="absolute left-7 top-12 w-0.5 h-6 bg-gray-300">
                  <AnimatePresence>
                    {step.status === 'completed' && (
                      <motion.div
                        initial={{ height: 0 }}
                        animate={{ height: '100%' }}
                        exit={{ height: 0 }}
                        className="w-full bg-green-400 rounded"
                        transition={{ duration: 0.5 }}
                      />
                    )}
                  </AnimatePresence>
                </div>
              )}
            </motion.div>
          ))}
        </div>
      </div>

      {/* Step List - Mobile View */}
      <div className="md:hidden">
        <div className="flex items-center justify-between overflow-x-auto pb-4 gap-4">
          {steps.map((step, index) => (
            <motion.div
              key={step.id}
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: index * 0.1 }}
              className="flex-shrink-0 text-center"
            >
              <div className="mb-2">
                {getStepIcon(step, index)}
              </div>
              
              <div className="text-xs font-medium text-gray-700 mb-1 max-w-20 truncate">
                {step.title.split(' ')[0]}
              </div>
              
              {step.status === 'active' && step.substeps && (
                <div className="text-xs text-gray-500">
                  {step.currentSubstep! + 1}/{step.substeps.length}
                </div>
              )}
              
              {step.error && (
                <div className="text-xs text-red-600 mt-1">
                  Error
                </div>
              )}

              {/* Connection line for mobile */}
              {index < steps.length - 1 && (
                <div className="absolute top-3 left-full w-8 h-0.5 bg-gray-300 transform translate-x-2">
                  <AnimatePresence>
                    {step.status === 'completed' && (
                      <motion.div
                        initial={{ width: 0 }}
                        animate={{ width: '100%' }}
                        exit={{ width: 0 }}
                        className="h-full bg-green-400 rounded"
                        transition={{ duration: 0.5 }}
                      />
                    )}
                  </AnimatePresence>
                </div>
              )}
            </motion.div>
          ))}
        </div>

        {/* Current Step Details - Mobile */}
        {currentStep < steps.length && (
          <motion.div
            key={`mobile-step-${currentStep}`}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            className={`p-4 rounded-lg border mt-4 ${getStepClasses(steps[currentStep])}`}
          >
            <h4 className="font-medium mb-2">{steps[currentStep].title}</h4>
            <p className="text-sm text-gray-600 mb-3">{steps[currentStep].description}</p>
            
            {steps[currentStep].error && (
              <div className="text-sm text-red-600 p-2 bg-red-100 rounded mb-3">
                {steps[currentStep].error}
              </div>
            )}
            
            {steps[currentStep].substeps && steps[currentStep].status === 'active' && (
              <div>
                <div className="text-sm text-gray-700 mb-2">
                  {steps[currentStep].substeps![steps[currentStep].currentSubstep || 0]}
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <motion.div
                    className="h-full bg-blue-400 rounded-full"
                    initial={{ width: 0 }}
                    animate={{ 
                      width: `${((steps[currentStep].currentSubstep || 0) + 1) / steps[currentStep].substeps!.length * 100}%` 
                    }}
                    transition={{ duration: 0.3 }}
                  />
                </div>
              </div>
            )}
          </motion.div>
        )}
      </div>
    </div>
  )
}