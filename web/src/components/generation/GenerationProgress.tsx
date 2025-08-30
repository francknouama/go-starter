/**
 * Generation Progress Component
 * Real-time progress display with WebSocket integration
 */

import React, { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  SignalIcon,
  SignalSlashIcon,
  ComputerDesktopIcon,
  CloudIcon,
  CpuChipIcon,
  DocumentIcon,
  CheckIcon
} from '@heroicons/react/24/outline'

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

interface GenerationProgressProps {
  steps: WorkflowStep[]
  currentStep: number
  progress: number
  estimatedDuration: number
  startTime?: Date | null
  wsConnected?: boolean
  className?: string
}

interface ProgressMetric {
  label: string
  value: string | number
  unit?: string
  icon: React.ComponentType<any>
  color: string
}

export default function GenerationProgress({
  steps,
  currentStep,
  progress,
  estimatedDuration,
  startTime,
  wsConnected = false,
  className = ''
}: GenerationProgressProps) {
  const [animatedProgress, setAnimatedProgress] = useState(0)
  const [filesGenerated, setFilesGenerated] = useState(0)
  const [processingRate, setProcessingRate] = useState(0)
  const [memoryUsage, setMemoryUsage] = useState(0)

  // Animate progress changes
  useEffect(() => {
    const target = Math.max(progress, (currentStep / steps.length) * 100)
    const duration = 500
    const steps = 60
    const increment = (target - animatedProgress) / steps
    
    if (Math.abs(target - animatedProgress) > 0.1) {
      const timer = setInterval(() => {
        setAnimatedProgress(prev => {
          const next = prev + increment
          if ((increment > 0 && next >= target) || (increment < 0 && next <= target)) {
            clearInterval(timer)
            return target
          }
          return next
        })
      }, duration / steps)

      return () => clearInterval(timer)
    }
  }, [progress, currentStep, steps.length, animatedProgress])

  // Simulate metrics (in real app, these would come from WebSocket)
  useEffect(() => {
    if (wsConnected && currentStep < steps.length && steps[currentStep]?.status === 'active') {
      const interval = setInterval(() => {
        setFilesGenerated(prev => prev + Math.floor(Math.random() * 3))
        setProcessingRate(prev => {
          const variation = (Math.random() - 0.5) * 20
          return Math.max(0, Math.min(100, prev + variation))
        })
        setMemoryUsage(prev => {
          const variation = (Math.random() - 0.5) * 10
          return Math.max(20, Math.min(80, prev + variation))
        })
      }, 1000)

      return () => clearInterval(interval)
    }
  }, [wsConnected, currentStep, steps])

  // Calculate elapsed time
  const elapsedTime = startTime ? Math.floor((Date.now() - startTime.getTime()) / 1000) : 0
  const remainingTime = Math.max(0, estimatedDuration - elapsedTime)

  // Format time
  const formatTime = (seconds: number): string => {
    if (seconds < 60) return `${seconds}s`
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = seconds % 60
    return `${minutes}m ${remainingSeconds}s`
  }

  // Get current step info
  const activeStep = currentStep < steps.length ? steps[currentStep] : null

  // Metrics data
  const metrics: ProgressMetric[] = [
    {
      label: 'Files Generated',
      value: filesGenerated,
      icon: DocumentIcon,
      color: 'text-blue-600'
    },
    {
      label: 'Processing Rate',
      value: processingRate.toFixed(0),
      unit: 'files/min',
      icon: CpuChipIcon,
      color: 'text-green-600'
    },
    {
      label: 'Memory Usage',
      value: memoryUsage.toFixed(0),
      unit: 'MB',
      icon: ComputerDesktopIcon,
      color: 'text-orange-600'
    },
    {
      label: 'Elapsed Time',
      value: formatTime(elapsedTime),
      icon: CloudIcon,
      color: 'text-purple-600'
    }
  ]

  return (
    <div className={`bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden ${className}`}>
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-50 to-purple-50 px-6 py-4 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="relative">
              <div className="w-3 h-3 rounded-full bg-green-500 animate-pulse"></div>
              <div className="absolute inset-0 w-3 h-3 rounded-full bg-green-400 animate-ping opacity-25"></div>
            </div>
            <h3 className="text-lg font-semibold text-gray-900">Live Generation Progress</h3>
          </div>
          
          <div className="flex items-center gap-2 text-sm">
            {wsConnected ? (
              <>
                <SignalIcon className="w-4 h-4 text-green-600" />
                <span className="text-green-700">Connected</span>
              </>
            ) : (
              <>
                <SignalSlashIcon className="w-4 h-4 text-gray-500" />
                <span className="text-gray-600">Simulated</span>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Progress Bar */}
      <div className="px-6 py-4">
        <div className="mb-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-700">Overall Progress</span>
            <span className="text-sm font-bold text-gray-900">{animatedProgress.toFixed(1)}%</span>
          </div>
          
          <div className="relative w-full bg-gray-200 rounded-full h-3 overflow-hidden">
            <motion.div
              className="absolute inset-y-0 left-0 bg-gradient-to-r from-blue-500 via-purple-500 to-green-500 rounded-full"
              initial={{ width: 0 }}
              animate={{ width: `${animatedProgress}%` }}
              transition={{ 
                type: "spring", 
                stiffness: 100, 
                damping: 15,
                mass: 1
              }}
            />
            
            {/* Shimmer effect */}
            <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white to-transparent opacity-25 animate-pulse rounded-full"></div>
          </div>
        </div>

        {/* Current Step Status */}
        {activeStep && (
          <AnimatePresence mode="wait">
            <motion.div
              key={`step-${currentStep}`}
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              className="bg-blue-50 rounded-lg p-4 mb-4"
            >
              <div className="flex items-start gap-3">
                <div className="flex-shrink-0">
                  <div className="w-2 h-2 rounded-full bg-blue-500 animate-pulse mt-2"></div>
                </div>
                <div className="flex-1">
                  <h4 className="font-medium text-blue-900 mb-1">{activeStep.title}</h4>
                  <p className="text-sm text-blue-700 mb-2">{activeStep.description}</p>
                  
                  {/* Substep Progress */}
                  {activeStep.substeps && typeof activeStep.currentSubstep === 'number' && (
                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm">
                        <span className="text-blue-800 font-medium">
                          Step {activeStep.currentSubstep + 1} of {activeStep.substeps.length}:
                        </span>
                        <span className="text-blue-700">
                          {activeStep.substeps[activeStep.currentSubstep]}
                        </span>
                      </div>
                      
                      <div className="flex items-center gap-1">
                        {activeStep.substeps.map((_, index) => (
                          <div
                            key={index}
                            className={`flex-1 h-1.5 rounded-full transition-colors duration-300 ${
                              index <= activeStep.currentSubstep!
                                ? 'bg-blue-500'
                                : 'bg-blue-200'
                            }`}
                          >
                            {index === activeStep.currentSubstep && (
                              <div className="w-full h-full bg-blue-600 rounded-full animate-pulse"></div>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </motion.div>
          </AnimatePresence>
        )}

        {/* Metrics Grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
          {metrics.map((metric, index) => (
            <motion.div
              key={metric.label}
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ delay: index * 0.1 }}
              className="bg-gray-50 rounded-lg p-3 text-center"
            >
              <div className="flex items-center justify-center mb-2">
                <metric.icon className={`w-5 h-5 ${metric.color}`} />
              </div>
              <div className="text-lg font-bold text-gray-900">
                {metric.value}
                {metric.unit && <span className="text-sm font-normal text-gray-600 ml-1">{metric.unit}</span>}
              </div>
              <div className="text-xs text-gray-600">{metric.label}</div>
            </motion.div>
          ))}
        </div>

        {/* Time Estimation */}
        <div className="flex items-center justify-between text-sm text-gray-600">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-gray-400"></div>
            <span>Time elapsed: {formatTime(elapsedTime)}</span>
          </div>
          {remainingTime > 0 && (
            <div className="flex items-center gap-2">
              <span>Est. remaining: {formatTime(remainingTime)}</span>
              <div className="w-2 h-2 rounded-full bg-blue-400"></div>
            </div>
          )}
        </div>
      </div>

      {/* Real-time Activity Feed (if WebSocket connected) */}
      {wsConnected && (
        <div className="border-t border-gray-200 bg-gray-50 px-6 py-3">
          <div className="flex items-center gap-2 text-sm">
            <div className="w-2 h-2 rounded-full bg-green-400 animate-pulse"></div>
            <span className="text-gray-700 font-medium">Live Updates:</span>
            <span className="text-gray-600">Real-time progress from server</span>
          </div>
        </div>
      )}

      {/* Completed Steps Summary */}
      {steps.filter(s => s.status === 'completed').length > 0 && (
        <div className="border-t border-gray-200 px-6 py-3 bg-green-50">
          <div className="flex items-center gap-2 text-sm text-green-700">
            <CheckIcon className="w-4 h-4" />
            <span className="font-medium">
              Completed: {steps.filter(s => s.status === 'completed').map(s => s.title).join(', ')}
            </span>
          </div>
        </div>
      )}
    </div>
  )
}