/**
 * Error View
 * Displays generation errors with recovery options
 */

import React from 'react'
import { motion } from 'framer-motion'
import { 
  ExclamationTriangleIcon,
  ArrowPathIcon,
  ArrowLeftIcon,
  HomeIcon,
  DocumentTextIcon
} from '@heroicons/react/24/outline'
import type { WorkflowSession } from '../../hooks/useGenerationWorkflow'
import Button from '../common/Button'

interface ErrorViewProps {
  error: string
  session: WorkflowSession | null
  onRetry: () => void
  onBack: () => void
  onNewProject: () => void
  className?: string
}

export default function ErrorView({
  error,
  session,
  onRetry,
  onBack,
  onNewProject,
  className = ''
}: ErrorViewProps) {
  // Determine error type for better messaging
  const getErrorType = (errorMessage: string) => {
    const message = errorMessage.toLowerCase()
    if (message.includes('network') || message.includes('connection')) {
      return 'network'
    } else if (message.includes('validation') || message.includes('invalid')) {
      return 'validation'
    } else if (message.includes('timeout')) {
      return 'timeout'
    } else if (message.includes('cancelled')) {
      return 'cancelled'
    }
    return 'unknown'
  }

  const errorType = getErrorType(error)

  const getErrorTitle = (type: string) => {
    switch (type) {
      case 'network':
        return 'Connection Error'
      case 'validation':
        return 'Configuration Error'
      case 'timeout':
        return 'Generation Timeout'
      case 'cancelled':
        return 'Generation Cancelled'
      default:
        return 'Generation Failed'
    }
  }

  const getErrorDescription = (type: string) => {
    switch (type) {
      case 'network':
        return 'Unable to connect to the server. Please check your internet connection and try again.'
      case 'validation':
        return 'There was an issue with your project configuration. Please review your settings and try again.'
      case 'timeout':
        return 'Project generation took too long to complete. This might be due to high server load.'
      case 'cancelled':
        return 'Project generation was cancelled. You can start over or modify your configuration.'
      default:
        return 'An unexpected error occurred during project generation. Please try again.'
    }
  }

  const getSuggestions = (type: string) => {
    switch (type) {
      case 'network':
        return [
          'Check your internet connection',
          'Try refreshing the page',
          'Wait a moment and retry',
          'Contact support if the problem persists'
        ]
      case 'validation':
        return [
          'Review your project name and module path',
          'Check for invalid characters or formatting',
          'Ensure all required fields are filled',
          'Try a simpler configuration first'
        ]
      case 'timeout':
        return [
          'Try generating a simpler project first',
          'Reduce advanced features temporarily',
          'Wait a few minutes and retry',
          'Check server status'
        ]
      case 'cancelled':
        return [
          'Review your project configuration',
          'Make any necessary adjustments',
          'Try generating again when ready'
        ]
      default:
        return [
          'Try again with the same configuration',
          'Simplify your project setup temporarily',
          'Check the error details below',
          'Contact support if the issue persists'
        ]
    }
  }

  return (
    <div className={`h-full flex items-center justify-center p-6 ${className}`}>
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="max-w-2xl w-full"
      >
        <div className="bg-white rounded-2xl shadow-xl border border-gray-200 overflow-hidden">
          {/* Header */}
          <div className="bg-red-50 px-8 py-6 border-b border-red-100">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
                <ExclamationTriangleIcon className="w-6 h-6 text-red-600" />
              </div>
              <div>
                <h1 className="text-xl font-semibold text-red-900">
                  {getErrorTitle(errorType)}
                </h1>
                <p className="text-red-700 mt-1">
                  {getErrorDescription(errorType)}
                </p>
              </div>
            </div>
          </div>

          {/* Content */}
          <div className="p-8">
            {/* Error Details */}
            <div className="mb-6">
              <h3 className="text-sm font-medium text-gray-900 mb-3">Error Details</h3>
              <div className="bg-gray-50 rounded-lg p-4 border">
                <code className="text-sm text-gray-800 break-words">
                  {error}
                </code>
              </div>
            </div>

            {/* Session Information */}
            {session && (
              <div className="mb-6">
                <h3 className="text-sm font-medium text-gray-900 mb-3">Session Information</h3>
                <div className="bg-gray-50 rounded-lg p-4 space-y-2 text-sm">
                  {session.blueprint && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">Blueprint:</span>
                      <span className="font-medium">{session.blueprint.name}</span>
                    </div>
                  )}
                  <div className="flex justify-between">
                    <span className="text-gray-600">Project Name:</span>
                    <span className="font-medium">{session.config.projectName || 'Not set'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Progress:</span>
                    <span className="font-medium">{session.progress}%</span>
                  </div>
                  {session.startTime && (
                    <div className="flex justify-between">
                      <span className="text-gray-600">Started:</span>
                      <span className="font-medium">
                        {new Date(session.startTime).toLocaleTimeString()}
                      </span>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Suggestions */}
            <div className="mb-8">
              <h3 className="text-sm font-medium text-gray-900 mb-3">What you can try:</h3>
              <ul className="space-y-2">
                {getSuggestions(errorType).map((suggestion, index) => (
                  <li key={index} className="flex items-start gap-3">
                    <div className="w-5 h-5 rounded-full bg-blue-100 flex items-center justify-center mt-0.5">
                      <span className="text-xs font-medium text-blue-600">{index + 1}</span>
                    </div>
                    <span className="text-gray-700">{suggestion}</span>
                  </li>
                ))}
              </ul>
            </div>

            {/* Actions */}
            <div className="flex flex-wrap gap-3 justify-center">
              <Button
                variant="primary"
                onClick={onRetry}
                className="bg-red-600 hover:bg-red-700"
              >
                <ArrowPathIcon className="w-4 h-4 mr-2" />
                Try Again
              </Button>

              <Button
                variant="outline"
                onClick={onBack}
                className="border-gray-300 text-gray-700 hover:bg-gray-50"
              >
                <ArrowLeftIcon className="w-4 h-4 mr-2" />
                Back to Configuration
              </Button>

              <Button
                variant="ghost"
                onClick={onNewProject}
                className="text-gray-600 hover:text-gray-900"
              >
                <HomeIcon className="w-4 h-4 mr-2" />
                Start Over
              </Button>
            </div>

            {/* Help Link */}
            <div className="mt-6 pt-6 border-t border-gray-200 text-center">
              <p className="text-sm text-gray-600 mb-3">
                Still having trouble? Check our troubleshooting guide or contact support.
              </p>
              <Button
                variant="ghost"
                size="sm"
                className="text-blue-600 hover:text-blue-700"
              >
                <DocumentTextIcon className="w-4 h-4 mr-2" />
                View Troubleshooting Guide
              </Button>
            </div>
          </div>
        </div>
      </motion.div>
    </div>
  )
}