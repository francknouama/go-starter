/**
 * Configuration View
 * Main configuration interface with enhanced workflow integration
 */

import React from 'react'
import { motion } from 'framer-motion'
import { ArrowLeftIcon, CogIcon } from '@heroicons/react/24/outline'
import type { Blueprint, ProjectConfig } from '../../services/api'
import type { UserPreferences } from '../../stores/generationStore'
import Button from '../common/Button'
import ConfigurationPanel from '../forms/ConfigurationPanel'
import BlueprintSpecificOptions from '../forms/BlueprintSpecificOptions'
import AdvancedOptions from '../forms/AdvancedOptions'
import FileExplorerPanel from '../preview/FileExplorerPanel'

interface ConfigurationViewProps {
  blueprint: Blueprint | null
  config: ProjectConfig
  isValid: boolean
  onConfigChange: (config: Partial<ProjectConfig>) => void
  onGenerate: () => void
  onBack: () => void
  preferences: UserPreferences
  className?: string
}

export default function ConfigurationView({
  blueprint,
  config,
  isValid,
  onConfigChange,
  onGenerate,
  onBack,
  preferences,
  className = ''
}: ConfigurationViewProps) {
  const [activeTab, setActiveTab] = React.useState<'basic' | 'blueprint' | 'advanced'>('basic')

  return (
    <div className={`h-full flex flex-col ${className}`}>
      {/* Header */}
      <div className="px-6 py-4 bg-white border-b border-gray-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="sm" onClick={onBack}>
              <ArrowLeftIcon className="w-4 h-4 mr-2" />
              Back to Blueprints
            </Button>
            
            {blueprint && (
              <div className="flex items-center gap-3">
                <CogIcon className="w-6 h-6 text-blue-600" />
                <div>
                  <h1 className="text-xl font-semibold text-gray-900">
                    Configure {blueprint.name}
                  </h1>
                  <p className="text-sm text-gray-600">{blueprint.description}</p>
                </div>
              </div>
            )}
          </div>

          <Button
            variant="primary"
            disabled={!isValid}
            onClick={onGenerate}
            className="bg-gradient-to-r from-green-600 to-blue-600 hover:from-green-700 hover:to-blue-700"
          >
            Generate Project
          </Button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        {/* Configuration Panel */}
        <div className="flex-1 overflow-y-auto">
          <div className="p-6">
            {/* Tab Navigation */}
            <div className="flex space-x-1 mb-6 bg-gray-100 p-1 rounded-lg max-w-fit">
              {[
                { id: 'basic', label: 'Basic' },
                { id: 'blueprint', label: 'Blueprint Options' },
                { id: 'advanced', label: 'Advanced' }
              ].map((tab) => (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id as any)}
                  className={`px-4 py-2 text-sm font-medium rounded-md transition-all ${
                    activeTab === tab.id
                      ? 'bg-white text-gray-900 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </div>

            {/* Tab Content */}
            <motion.div
              key={activeTab}
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ duration: 0.2 }}
            >
              {activeTab === 'basic' && (
                <ConfigurationPanel
                  disclosureMode={preferences.disclosureMode}
                  config={config}
                  onConfigChange={onConfigChange}
                  blueprints={blueprint ? [blueprint] : []}
                  onGenerate={(request) => {
                    // Update config with any changes from the request
                    if (request.config) {
                      onConfigChange(request.config)
                    }
                    onGenerate()
                  }}
                  onDownload={() => {}} // Not used in this context
                />
              )}

              {activeTab === 'blueprint' && (
                <BlueprintSpecificOptions
                  blueprint={blueprint}
                  config={config}
                  onConfigChange={onConfigChange}
                  isAdvanced={preferences.disclosureMode === 'advanced'}
                />
              )}

              {activeTab === 'advanced' && (
                <AdvancedOptions
                  blueprint={blueprint}
                  config={config}
                  onConfigChange={onConfigChange}
                />
              )}
            </motion.div>
          </div>
        </div>

        {/* Preview Panel */}
        <div className="w-96 border-l border-gray-200 bg-gray-50">
          <FileExplorerPanel
            config={config}
            preview={null}
          />
        </div>
      </div>
    </div>
  )
}