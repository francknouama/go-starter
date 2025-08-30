/**
 * Download and Success Flow Component
 * Handles project download and displays success information
 */

import React, { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  CheckCircleIcon,
  DocumentArrowDownIcon,
  FolderIcon,
  DocumentTextIcon,
  ClipboardDocumentIcon,
  ShareIcon,
  ArrowPathIcon,
  EyeIcon,
  InformationCircleIcon,
  ExclamationTriangleIcon
} from '@heroicons/react/24/outline'
import { 
  CheckCircleIcon as CheckCircleIconSolid 
} from '@heroicons/react/24/solid'
import type { ProjectConfig, Blueprint } from '../../services/api'
import Button from '../common/Button'
import { useProjectDownload } from '../../hooks/useApi'

interface DownloadSuccessFlowProps {
  projectId: string
  projectName: string
  blueprint: Blueprint | null
  config: ProjectConfig
  fileCount: number
  estimatedSize: string
  generationTime: number
  onNewProject?: () => void
  onShareConfig?: () => void
  onShowPreview?: () => void
  className?: string
}

interface ProjectSummary {
  structure: FileStructure[]
  features: string[]
  technologies: string[]
  nextSteps: string[]
  quickCommands: QuickCommand[]
}

interface FileStructure {
  name: string
  type: 'folder' | 'file'
  description: string
  children?: FileStructure[]
}

interface QuickCommand {
  title: string
  command: string
  description: string
  icon: React.ComponentType<any>
}

export default function DownloadSuccessFlow({
  projectId,
  projectName,
  blueprint,
  config,
  fileCount,
  estimatedSize,
  generationTime,
  onNewProject,
  onShareConfig,
  onShowPreview,
  className = ''
}: DownloadSuccessFlowProps) {
  const { downloadProject, downloading, error: downloadError } = useProjectDownload()
  const [downloadComplete, setDownloadComplete] = useState(false)
  const [showProjectSummary, setShowProjectSummary] = useState(false)
  const [copiedCommand, setCopiedCommand] = useState<string | null>(null)

  // Auto-show project summary after a delay
  useEffect(() => {
    const timer = setTimeout(() => {
      setShowProjectSummary(true)
    }, 2000)

    return () => clearTimeout(timer)
  }, [])

  // Generate project summary based on configuration
  const projectSummary: ProjectSummary = React.useMemo(() => {
    const structure: FileStructure[] = [
      {
        name: 'cmd',
        type: 'folder',
        description: 'Application entry points',
        children: [
          { name: 'main.go', type: 'file', description: 'Main application entry point' }
        ]
      },
      {
        name: 'internal',
        type: 'folder',
        description: 'Private application code',
        children: [
          { name: 'handlers', type: 'folder', description: 'HTTP handlers' },
          { name: 'models', type: 'folder', description: 'Data models' },
          { name: 'services', type: 'folder', description: 'Business logic' }
        ]
      },
      {
        name: 'pkg',
        type: 'folder',
        description: 'Public libraries',
        children: []
      },
      { name: 'go.mod', type: 'file', description: 'Go module definition' },
      { name: 'go.sum', type: 'file', description: 'Go module checksums' },
      { name: 'README.md', type: 'file', description: 'Project documentation' },
      { name: 'Makefile', type: 'file', description: 'Build automation' },
      { name: '.gitignore', type: 'file', description: 'Git ignore rules' }
    ]

    // Add config-specific folders
    if (config.features?.database) {
      structure[1].children?.push(
        { name: 'database', type: 'folder', description: 'Database connections' },
        { name: 'migrations', type: 'folder', description: 'Database migrations' }
      )
      structure.push(
        { name: 'docker-compose.yml', type: 'file', description: 'Database services' }
      )
    }

    if (config.features?.authentication) {
      structure[1].children?.push(
        { name: 'auth', type: 'folder', description: 'Authentication logic' },
        { name: 'middleware', type: 'folder', description: 'HTTP middleware' }
      )
    }

    const features: string[] = []
    const technologies: string[] = [`Go ${config.goVersion}`, config.framework, config.logger]

    // Add features based on config
    if (config.features?.database) {
      features.push(`${config.features.database.driver} Database`)
      features.push(`${config.features.database.orm} ORM`)
      technologies.push(config.features.database.driver, config.features.database.orm)
    }

    if (config.features?.authentication) {
      features.push(`${config.features.authentication.type} Authentication`)
      technologies.push(config.features.authentication.type)
    }

    if (config.features?.advanced) {
      Object.entries(config.features.advanced).forEach(([key, feature]) => {
        if (feature.enabled) {
          features.push(key.split('-').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' '))
        }
      })
    }

    features.push(`${config.architecture} Architecture`)
    features.push(`${config.projectType} Template`)

    const nextSteps: string[] = [
      'Extract the downloaded ZIP file',
      'Navigate to the project directory',
      'Install dependencies with `go mod download`',
      'Review the README.md file for project-specific instructions',
      'Run `make help` to see available commands',
      'Start development with `make dev` or `go run main.go`'
    ]

    // Add config-specific next steps
    if (config.features?.database) {
      nextSteps.splice(3, 0, 'Set up your database connection in config.yaml')
      nextSteps.splice(4, 0, 'Run database migrations with `make migrate`')
    }

    const quickCommands: QuickCommand[] = [
      {
        title: 'Download Dependencies',
        command: 'go mod download',
        description: 'Install all Go module dependencies',
        icon: DocumentArrowDownIcon
      },
      {
        title: 'Build Project',
        command: 'make build',
        description: 'Compile the project binary',
        icon: FolderIcon
      },
      {
        title: 'Run Tests',
        command: 'make test',
        description: 'Execute all project tests',
        icon: CheckCircleIcon
      },
      {
        title: 'Start Development',
        command: config.projectType === 'cli' ? 'go run main.go --help' : 'make dev',
        description: config.projectType === 'cli' ? 'Show CLI help' : 'Start development server',
        icon: ArrowPathIcon
      }
    ]

    return {
      structure,
      features,
      technologies: [...new Set(technologies)],
      nextSteps,
      quickCommands
    }
  }, [config])

  // Handle download
  const handleDownload = async () => {
    try {
      // Use the download hook to handle the actual download
      const generationRequest = {
        ...config,
        outputFormat: 'zip' as const,
        includeTests: true,
        includeDocs: true
      }
      
      await downloadProject(generationRequest, `${projectName}.zip`)
      setDownloadComplete(true)
    } catch (error) {
      console.error('Download failed:', error)
    }
  }

  // Copy command to clipboard
  const copyCommand = async (command: string) => {
    try {
      await navigator.clipboard.writeText(command)
      setCopiedCommand(command)
      setTimeout(() => setCopiedCommand(null), 2000)
    } catch (error) {
      console.error('Failed to copy command:', error)
    }
  }

  return (
    <div className={`max-w-6xl mx-auto space-y-8 ${className}`}>
      {/* Success Header */}
      <motion.div
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        className="text-center"
      >
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ delay: 0.2, type: "spring", stiffness: 150 }}
          className="w-20 h-20 mx-auto mb-6 rounded-full bg-green-100 flex items-center justify-center"
        >
          <CheckCircleIconSolid className="w-12 h-12 text-green-600" />
        </motion.div>
        
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Project Generated Successfully!
        </h1>
        <p className="text-lg text-gray-600 mb-6">
          Your <strong>{blueprint?.name}</strong> project "<strong>{projectName}</strong>" is ready to download.
        </p>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-2xl mx-auto">
          <div className="text-center">
            <div className="text-3xl font-bold text-blue-600">{fileCount}</div>
            <div className="text-sm text-gray-600">Files Generated</div>
          </div>
          <div className="text-center">
            <div className="text-3xl font-bold text-purple-600">{estimatedSize}</div>
            <div className="text-sm text-gray-600">Project Size</div>
          </div>
          <div className="text-center">
            <div className="text-3xl font-bold text-green-600">{generationTime}s</div>
            <div className="text-sm text-gray-600">Generation Time</div>
          </div>
        </div>
      </motion.div>

      {/* Download Section */}
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ delay: 0.4 }}
        className="bg-gradient-to-r from-green-50 to-blue-50 rounded-2xl p-8 border border-green-200"
      >
        <div className="text-center">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Download Your Project
          </h2>
          
          {!downloadComplete ? (
            <div className="space-y-4">
              <p className="text-gray-600">
                Your project is packaged and ready for download as a ZIP file.
              </p>
              
              <Button
                variant="primary"
                size="lg"
                onClick={handleDownload}
                disabled={downloading}
                className="bg-gradient-to-r from-green-600 to-blue-600 hover:from-green-700 hover:to-blue-700 text-white shadow-lg"
              >
                <DocumentArrowDownIcon className="w-5 h-5 mr-2" />
                {downloading ? 'Preparing Download...' : `Download ${projectName}.zip`}
              </Button>

              {downloadError && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-4 mt-4">
                  <div className="flex items-center gap-3">
                    <ExclamationTriangleIcon className="w-5 h-5 text-red-600" />
                    <div className="text-left">
                      <p className="text-red-800 font-medium">Download Failed</p>
                      <p className="text-red-700 text-sm">{downloadError}</p>
                    </div>
                  </div>
                  <Button
                    variant="outline"
                    onClick={handleDownload}
                    className="mt-3 border-red-300 text-red-700 hover:bg-red-50"
                  >
                    Try Again
                  </Button>
                </div>
              )}
            </div>
          ) : (
            <motion.div
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              className="space-y-4"
            >
              <div className="flex items-center justify-center gap-3 text-green-700">
                <CheckCircleIconSolid className="w-6 h-6" />
                <span className="font-semibold">Download Complete!</span>
              </div>
              <p className="text-gray-600">
                Your project has been downloaded. Check your Downloads folder for the ZIP file.
              </p>
            </motion.div>
          )}

          {/* Action Buttons */}
          <div className="flex flex-wrap justify-center gap-3 mt-6">
            {onShowPreview && (
              <Button
                variant="outline"
                onClick={onShowPreview}
                className="border-blue-300 text-blue-700 hover:bg-blue-50"
              >
                <EyeIcon className="w-4 h-4 mr-2" />
                Preview Project
              </Button>
            )}
            
            {onShareConfig && (
              <Button
                variant="ghost"
                onClick={onShareConfig}
                className="text-gray-600"
              >
                <ShareIcon className="w-4 h-4 mr-2" />
                Share Configuration
              </Button>
            )}
            
            {onNewProject && (
              <Button
                variant="ghost"
                onClick={onNewProject}
                className="text-gray-600"
              >
                <ArrowPathIcon className="w-4 h-4 mr-2" />
                Generate Another Project
              </Button>
            )}
          </div>
        </div>
      </motion.div>

      {/* Project Summary */}
      <AnimatePresence>
        {showProjectSummary && (
          <motion.div
            initial={{ opacity: 0, y: 40 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -40 }}
            className="grid grid-cols-1 lg:grid-cols-2 gap-8"
          >
            {/* Project Features */}
            <div className="bg-white rounded-xl shadow-lg border border-gray-200 p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
                <InformationCircleIcon className="w-5 h-5 text-blue-600" />
                Project Features
              </h3>
              <div className="space-y-3">
                {projectSummary.features.map((feature, index) => (
                  <motion.div
                    key={feature}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="flex items-center gap-3"
                  >
                    <CheckCircleIcon className="w-4 h-4 text-green-600" />
                    <span className="text-gray-700">{feature}</span>
                  </motion.div>
                ))}
              </div>
              
              <div className="mt-6 pt-6 border-t border-gray-200">
                <h4 className="font-medium text-gray-900 mb-3">Technologies Used</h4>
                <div className="flex flex-wrap gap-2">
                  {projectSummary.technologies.map((tech) => (
                    <span
                      key={tech}
                      className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm"
                    >
                      {tech}
                    </span>
                  ))}
                </div>
              </div>
            </div>

            {/* Quick Commands */}
            <div className="bg-white rounded-xl shadow-lg border border-gray-200 p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
                <DocumentTextIcon className="w-5 h-5 text-purple-600" />
                Quick Commands
              </h3>
              <div className="space-y-3">
                {projectSummary.quickCommands.map((cmd, index) => (
                  <motion.div
                    key={cmd.command}
                    initial={{ opacity: 0, x: 20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ delay: index * 0.1 }}
                    className="group"
                  >
                    <div className="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
                      <div className="flex items-center gap-3">
                        <cmd.icon className="w-4 h-4 text-gray-600" />
                        <div>
                          <div className="font-medium text-gray-900">{cmd.title}</div>
                          <div className="text-sm text-gray-600">{cmd.description}</div>
                        </div>
                      </div>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => copyCommand(cmd.command)}
                        className="opacity-0 group-hover:opacity-100 transition-opacity"
                      >
                        {copiedCommand === cmd.command ? (
                          <CheckCircleIcon className="w-4 h-4 text-green-600" />
                        ) : (
                          <ClipboardDocumentIcon className="w-4 h-4 text-gray-500" />
                        )}
                      </Button>
                    </div>
                    <div className="mt-2 ml-7 font-mono text-sm text-gray-600 bg-gray-900 text-green-400 px-3 py-1 rounded">
                      $ {cmd.command}
                    </div>
                  </motion.div>
                ))}
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Next Steps */}
      <AnimatePresence>
        {showProjectSummary && (
          <motion.div
            initial={{ opacity: 0, y: 40 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -40 }}
            transition={{ delay: 0.2 }}
            className="bg-white rounded-xl shadow-lg border border-gray-200 p-6"
          >
            <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <ArrowPathIcon className="w-5 h-5 text-orange-600" />
              Next Steps
            </h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h4 className="font-medium text-gray-900 mb-3">Getting Started</h4>
                <ol className="space-y-3">
                  {projectSummary.nextSteps.slice(0, Math.ceil(projectSummary.nextSteps.length / 2)).map((step, index) => (
                    <motion.li
                      key={step}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: index * 0.1 }}
                      className="flex items-start gap-3"
                    >
                      <div className="flex-shrink-0 w-6 h-6 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center text-sm font-medium">
                        {index + 1}
                      </div>
                      <span className="text-gray-700">{step}</span>
                    </motion.li>
                  ))}
                </ol>
              </div>
              
              <div>
                <h4 className="font-medium text-gray-900 mb-3">Development</h4>
                <ol className="space-y-3" start={Math.ceil(projectSummary.nextSteps.length / 2) + 1}>
                  {projectSummary.nextSteps.slice(Math.ceil(projectSummary.nextSteps.length / 2)).map((step, index) => (
                    <motion.li
                      key={step}
                      initial={{ opacity: 0, x: 20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: (index + Math.ceil(projectSummary.nextSteps.length / 2)) * 0.1 }}
                      className="flex items-start gap-3"
                    >
                      <div className="flex-shrink-0 w-6 h-6 bg-green-100 text-green-600 rounded-full flex items-center justify-center text-sm font-medium">
                        {index + Math.ceil(projectSummary.nextSteps.length / 2) + 1}
                      </div>
                      <span className="text-gray-700">{step}</span>
                    </motion.li>
                  ))}
                </ol>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}