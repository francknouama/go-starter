import { useState } from 'react'
import { Disclosure } from '@headlessui/react'
import { ChevronDownIcon, SparklesIcon } from '@heroicons/react/20/solid'
import type { DisclosureMode, ProjectConfig, ProjectType, Architecture, Framework, LoggerType } from '../../types'
import type { ProjectTemplate } from '../../data/projectTemplates'
import ProjectTypeCard from './ProjectTypeCard'
import SelectionGrid from './SelectionGrid'
import CompactSelectionCard from './CompactSelectionCard'
import ValidatedInput from './ValidatedInput'
import FormValidationSummary from './FormValidationSummary'
import ModeSelector from './ModeSelector'
import ModeSuggestion from './ModeSuggestion'
import SmartHelp from '../help/SmartHelp'
import TemplateGallery from '../templates/TemplateGallery'
import { ProjectTypeIcons, ArchitectureIcons, FrameworkIcons, LoggerIcons } from '../icons/ProjectIcons'
import { ValidationPatterns, CustomValidators } from '../../utils/validation'
import HelpTooltip from '../common/HelpTooltip'
import Button from '../common/Button'
import { getHelpContent } from '../../content/helpContent'

interface ConfigurationPanelProps {
  disclosureMode: DisclosureMode
  onModeChange?: (mode: DisclosureMode) => void
  config: ProjectConfig
  onConfigChange: (config: ProjectConfig) => void
  blueprints: any[]
  onGenerate: (request: any) => void
  onDownload: (request: any) => void
}

export default function ConfigurationPanel({ 
  disclosureMode, 
  onModeChange, 
  config, 
  onConfigChange, 
  blueprints, 
  onGenerate, 
  onDownload 
}: ConfigurationPanelProps) {
  const [showTemplateGallery, setShowTemplateGallery] = useState(false)
  const [appliedTemplate, setAppliedTemplate] = useState<ProjectTemplate | null>(null)
  
  // Prevent unused variable warnings
  void blueprints
  void onDownload
  
  // Form validation state
  const isFormValid = () => {
    // Check required fields
    if (!config.projectName || !config.moduleUrl) return false
    
    // Check format validations
    if (!ValidationPatterns.projectName.test(config.projectName)) return false
    if (!ValidationPatterns.goModule.test(config.moduleUrl)) return false
    
    return true
  }

  const projectTypes: Array<{ value: ProjectType; label: string; description: string }> = [
    { value: 'cli', label: 'CLI Application', description: 'Command-line tools with Cobra framework' },
    { value: 'web-api', label: 'Web API', description: 'REST API server with various architectures' },
    { value: 'library', label: 'Library', description: 'Reusable Go packages' },
    { value: 'lambda', label: 'AWS Lambda', description: 'Serverless functions' },
    { value: 'microservice', label: 'Microservice', description: 'Distributed service with gRPC' },
  ]

  const architectures: Array<{ value: Architecture; label: string; description: string }> = [
    { value: 'standard', label: 'Standard', description: 'Simple layered architecture' },
    { value: 'clean', label: 'Clean Architecture', description: 'Layered with dependency inversion' },
    { value: 'ddd', label: 'Domain-Driven Design', description: 'Domain-focused approach' },
    { value: 'hexagonal', label: 'Hexagonal', description: 'Ports and adapters pattern' },
  ]

  const frameworks: Array<{ value: Framework; label: string; description: string }> = [
    { value: 'gin', label: 'Gin', description: 'Fast HTTP framework' },
    { value: 'echo', label: 'Echo', description: 'High performance framework' },
    { value: 'fiber', label: 'Fiber', description: 'Express-inspired framework' },
    { value: 'chi', label: 'Chi', description: 'Lightweight router' },
  ]

  const loggers: Array<{ value: LoggerType; label: string; description: string }> = [
    { value: 'slog', label: 'slog', description: 'Standard library structured logging' },
    { value: 'zap', label: 'Zap', description: 'High-performance logger' },
    { value: 'logrus', label: 'Logrus', description: 'Popular structured logger' },
    { value: 'zerolog', label: 'Zerolog', description: 'Zero allocation logger' },
  ]

  const updateConfig = (field: keyof ProjectConfig, value: any) => {
    onConfigChange({ ...config, [field]: value })
  }

  const applyTemplate = (template: ProjectTemplate) => {
    // Apply template configuration to current config
    const templateConfig = { ...config, ...template.config }
    
    // Show visual feedback
    console.log(`Applied template: ${template.name}`)
    
    // Ensure projectName is preserved if already set
    if (config.projectName && !template.config.projectName) {
      templateConfig.projectName = config.projectName
    }
    
    // Ensure moduleUrl is preserved if already set
    if (config.moduleUrl && !template.config.moduleUrl) {
      templateConfig.moduleUrl = config.moduleUrl
    }
    
    // Auto-generate a project name if none exists
    if (!templateConfig.projectName) {
      templateConfig.projectName = template.name.toLowerCase().replace(/[^a-z0-9]+/g, '-')
    }
    
    // Auto-generate module URL if none exists
    if (!templateConfig.moduleUrl) {
      templateConfig.moduleUrl = `github.com/user/${templateConfig.projectName}`
    }
    
    onConfigChange(templateConfig)
    setAppliedTemplate(template)
    setShowTemplateGallery(false)
  }

  return (
    <div className="bg-white/70 backdrop-blur-lg rounded-2xl shadow-xl border border-white/30 p-4 md:p-6 h-full overflow-y-auto">
      <div className="mb-4 md:mb-6">
        <h2 className="text-base md:text-lg font-semibold text-gray-900 mb-1 md:mb-2">Project Configuration</h2>
        <p className="text-xs md:text-sm text-gray-600">
          {disclosureMode === 'basic' 
            ? 'Configure essential project settings' 
            : 'Configure all available project options'
          }
        </p>
      </div>

      <div className="space-y-6">
        {/* Template Selection - New Feature */}
        <div className="bg-gradient-to-r from-purple-50 to-indigo-50 rounded-lg p-4 border border-purple-100">
          <div className="flex items-center justify-between mb-3">
            <div className="flex items-center gap-2">
              <SparklesIcon className="w-5 h-5 text-purple-600" />
              <div>
                <h3 className="font-medium text-purple-900">Start with a Template</h3>
                <p className="text-sm text-purple-700">Choose a real-world template to jumpstart your project</p>
              </div>
            </div>
            <Button 
              variant="primary"
              size="sm"
              onClick={() => setShowTemplateGallery(true)}
              className="bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-700 hover:to-indigo-700"
            >
              Browse Templates
            </Button>
          </div>
          
          {appliedTemplate && (
            <div className="mt-3 p-3 bg-white rounded-lg border border-purple-200">
              <div className="flex items-center gap-3">
                <div 
                  className="text-xl flex items-center justify-center w-8 h-8 rounded-lg"
                  style={{ color: appliedTemplate.color }}
                >
                  {appliedTemplate.icon}
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <span className="font-medium text-gray-900">{appliedTemplate.name}</span>
                    <span className="text-xs px-2 py-1 bg-green-100 text-green-800 rounded-full">Applied</span>
                  </div>
                  <p className="text-sm text-gray-600">{appliedTemplate.description}</p>
                </div>
                <Button 
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setAppliedTemplate(null)
                    // Reset to basic config
                    onConfigChange({
                      projectName: config.projectName,
                      moduleUrl: config.moduleUrl,
                      goVersion: '1.21',
                      projectType: 'web-api',
                      framework: 'gin',
                      architecture: 'standard',
                      logger: 'slog',
                    })
                  }}
                >
                  Remove
                </Button>
              </div>
            </div>
          )}
        </div>

        {/* Mode Selector - Always visible */}
        <ModeSelector 
          mode={disclosureMode} 
          onChange={(newMode) => {
            if (onModeChange) {
              onModeChange(newMode)
            }
          }}
        />
        
        {/* Mode Suggestion */}
        <ModeSuggestion
          projectType={config.projectType as ProjectType}
          currentMode={disclosureMode}
          onSwitchToAdvanced={() => {
            if (onModeChange) {
              onModeChange('advanced')
            }
          }}
        />
        
        {/* Smart Help - Context-aware suggestions */}
        <SmartHelp
          projectType={config.projectType as ProjectType}
          architecture={config.architecture as Architecture}
          framework={config.framework as Framework}
          disclosureMode={disclosureMode}
        />
        
        {/* Form Validation Summary */}
        <FormValidationSummary
          isValid={isFormValid()}
          projectName={config.projectName}
          moduleUrl={config.moduleUrl}
        />
        
        {/* Basic Configuration */}
        <Disclosure defaultOpen>
          {({ open }) => (
            <>
              <Disclosure.Button className="flex w-full justify-between rounded-lg bg-gray-50 px-4 py-2 text-left text-sm font-medium text-gray-900 hover:bg-gray-100 focus:outline-none focus-visible:ring focus-visible:ring-primary-500 focus-visible:ring-opacity-75">
                <span>Basic Settings</span>
                <ChevronDownIcon
                  className={`${open ? 'rotate-180 transform' : ''} h-5 w-5 text-gray-500`}
                />
              </Disclosure.Button>
              <Disclosure.Panel className="px-4 pt-4 pb-2 space-y-4">
                {/* Project Name - Validated */}
                <div>
                  <div className="flex items-center gap-2 mb-2">
                    <HelpTooltip 
                      content={getHelpContent('tips', 'project-name')}
                      variant="help"
                      position="right"
                    />
                  </div>
                  <ValidatedInput
                    label="Project Name"
                    value={config.projectName}
                    onChange={(value) => updateConfig('projectName', value)}
                    placeholder="my-awesome-project"
                    helpText="This will be used as your project directory name"
                    validation={{
                      required: true,
                      pattern: ValidationPatterns.projectName,
                      custom: CustomValidators.projectName
                    }}
                    showSuccessState={true}
                  />
                </div>

                {/* Module URL - Validated */}
                <div>
                  <div className="flex items-center gap-2 mb-2">
                    <HelpTooltip 
                      content={getHelpContent('concepts', 'module-path')}
                      variant="help"
                      position="right"
                    />
                  </div>
                  <ValidatedInput
                    label="Module Path"
                    value={config.moduleUrl}
                    onChange={(value) => updateConfig('moduleUrl', value)}
                    placeholder="github.com/user/my-awesome-project"
                    helpText="Go module path for your project"
                    validation={{
                      required: true,
                      pattern: ValidationPatterns.goModule,
                      custom: CustomValidators.goModule
                    }}
                    showSuccessState={true}
                  />
                </div>

                {/* Project Type - Card Selection */}
                <fieldset>
                  <legend className="flex items-center gap-2 mb-3">
                    <span className="label">Project Type</span>
                    <HelpTooltip 
                      content="Choose the type of Go application you want to create. Each type provides optimized templates and dependencies."
                      variant="help"
                      position="right"
                    />
                  </legend>
                  <div role="radiogroup" aria-labelledby="project-type-legend">
                    <span id="project-type-legend" className="sr-only">Project Type Selection</span>
                    <SelectionGrid columns={2}>
                      {projectTypes.map((type) => {
                        const Icon = ProjectTypeIcons[type.value as keyof typeof ProjectTypeIcons] || ProjectTypeIcons['web-api']
                        return (
                          <div key={type.value} className="relative">
                            <ProjectTypeCard
                              value={type.value}
                              icon={<Icon className="w-6 h-6" />}
                              title={type.label}
                              description={type.description}
                              selected={config.projectType === type.value}
                              onClick={() => updateConfig('projectType', type.value)}
                            />
                            <div className="absolute top-2 right-2">
                              <HelpTooltip 
                                content={getHelpContent('projectTypes', type.value)}
                                variant="info"
                                position="left"
                                maxWidth="max-w-sm"
                              />
                            </div>
                          </div>
                        )
                      })}
                    </SelectionGrid>
                  </div>
                </fieldset>

                {/* Go Version - Inline Selection */}
                <div>
                  <div className="flex items-center gap-2 mb-2">
                    <label className="label">Go Version</label>
                    <HelpTooltip 
                      content={getHelpContent('concepts', 'go-version')}
                      variant="help"
                      position="right"
                    />
                  </div>
                  <div className="flex space-x-2">
                    {['1.21', '1.20', '1.19'].map((version) => (
                      <button
                        key={version}
                        onClick={() => updateConfig('goVersion', version)}
                        className={`
                          px-4 py-2 rounded-lg border font-medium text-sm transition-all
                          ${config.goVersion === version
                            ? 'border-primary-500 bg-primary-50 text-primary-700'
                            : 'border-gray-300 bg-white text-gray-700 hover:border-gray-400'
                          }
                        `}
                      >
                        Go {version}
                      </button>
                    ))}
                  </div>
                </div>
              </Disclosure.Panel>
            </>
          )}
        </Disclosure>

        {/* Framework Configuration */}
        <Disclosure defaultOpen>
          {({ open }) => (
            <>
              <Disclosure.Button className="flex w-full justify-between rounded-lg bg-gray-50 px-4 py-2 text-left text-sm font-medium text-gray-900 hover:bg-gray-100 focus:outline-none focus-visible:ring focus-visible:ring-primary-500 focus-visible:ring-opacity-75">
                <span>Framework & Architecture</span>
                <ChevronDownIcon
                  className={`${open ? 'rotate-180 transform' : ''} h-5 w-5 text-gray-500`}
                />
              </Disclosure.Button>
              <Disclosure.Panel className="px-4 pt-4 pb-2 space-y-4">
                {/* Framework - Compact Cards */}
                {(config.projectType === 'web-api' || config.projectType === 'microservice') && (
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <label className="label">Framework</label>
                      <HelpTooltip 
                        content="HTTP framework for handling requests and routing. Each framework has different performance characteristics and API styles."
                        variant="help"
                        position="right"
                      />
                    </div>
                    <div className="grid grid-cols-2 gap-2">
                      {frameworks.map((framework) => {
                        const IconComponent = FrameworkIcons[framework.value as keyof typeof FrameworkIcons]
                        return (
                          <div key={framework.value} className="relative">
                            <CompactSelectionCard
                              value={framework.value}
                              icon={IconComponent ? <IconComponent /> : undefined}
                              label={framework.label}
                              description={framework.description}
                              selected={config.framework === framework.value}
                              onClick={() => updateConfig('framework', framework.value)}
                            />
                            <div className="absolute top-1 right-1">
                              <HelpTooltip 
                                content={getHelpContent('frameworks', framework.value)}
                                variant="info"
                                position="left"
                                maxWidth="max-w-xs"
                              />
                            </div>
                          </div>
                        )
                      })}
                    </div>
                  </div>
                )}

                {/* Architecture - Card Selection */}
                <div>
                  <div className="flex items-center gap-2 mb-3">
                    <label className="label">Architecture Pattern</label>
                    <HelpTooltip 
                      content={getHelpContent('tips', 'architecture-choice')}
                      variant="help"
                      position="right"
                    />
                  </div>
                  <SelectionGrid columns={2}>
                    {architectures.map((arch) => {
                      const Icon = ArchitectureIcons[arch.value as keyof typeof ArchitectureIcons] || ArchitectureIcons['standard']
                      return (
                        <div key={arch.value} className="relative">
                          <ProjectTypeCard
                            value={arch.value}
                            icon={<Icon className="w-6 h-6" />}
                            title={arch.label}
                            description={arch.description}
                            selected={config.architecture === arch.value}
                            onClick={() => updateConfig('architecture', arch.value)}
                          />
                          <div className="absolute top-2 right-2">
                            <HelpTooltip 
                              content={getHelpContent('architectures', arch.value)}
                              variant="info"
                              position="left"
                              maxWidth="max-w-sm"
                            />
                          </div>
                        </div>
                      )
                    })}
                  </SelectionGrid>
                </div>

                {/* Logger - Compact Cards */}
                <div>
                  <div className="flex items-center gap-2 mb-3">
                    <label className="label">Logger Type</label>
                    <HelpTooltip 
                      content="Logging library for your application. Consider performance requirements and feature needs when choosing."
                      variant="help"
                      position="right"
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-2">
                    {loggers.map((logger) => {
                      const IconComponent = LoggerIcons[logger.value as keyof typeof LoggerIcons]
                      return (
                        <div key={logger.value} className="relative">
                          <CompactSelectionCard
                            value={logger.value}
                            icon={IconComponent ? <IconComponent /> : undefined}
                            label={logger.label}
                            description={logger.description}
                            selected={config.logger === logger.value}
                            onClick={() => updateConfig('logger', logger.value)}
                          />
                          <div className="absolute top-1 right-1">
                            <HelpTooltip 
                              content={getHelpContent('loggers', logger.value)}
                              variant="info"
                              position="left"
                              maxWidth="max-w-xs"
                            />
                          </div>
                        </div>
                      )
                    })}
                  </div>
                </div>
              </Disclosure.Panel>
            </>
          )}
        </Disclosure>

        {/* Advanced Configuration (only in advanced mode) */}
        {disclosureMode === 'advanced' && (
          <>
            <Disclosure>
              {({ open }) => (
                <>
                  <Disclosure.Button className="flex w-full justify-between rounded-lg bg-purple-50 px-4 py-2 text-left text-sm font-medium text-purple-900 hover:bg-purple-100 focus:outline-none focus-visible:ring focus-visible:ring-primary-500 focus-visible:ring-opacity-75">
                    <div className="flex items-center gap-2">
                      <span>Database Configuration</span>
                      <HelpTooltip 
                        content="Configure database connectivity and ORM for data persistence. Skip if your project doesn't need a database."
                        variant="help"
                        position="right"
                        maxWidth="max-w-sm"
                      />
                    </div>
                    <ChevronDownIcon
                      className={`${open ? 'rotate-180 transform' : ''} h-5 w-5 text-purple-500`}
                    />
                  </Disclosure.Button>
                  <Disclosure.Panel className="px-4 pt-4 pb-2 space-y-4">
                    <div className="space-y-4">
                      <div>
                        <div className="flex items-center gap-2 mb-2">
                          <label className="label">Database Driver</label>
                          <HelpTooltip 
                            content="Database system for data persistence. Choose based on your data requirements and scalability needs."
                            variant="help"
                            position="right"
                          />
                        </div>
                        <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
                          {[
                            { value: 'postgres', label: 'PostgreSQL' },
                            { value: 'mysql', label: 'MySQL' },
                            { value: 'mongodb', label: 'MongoDB' },
                            { value: 'sqlite', label: 'SQLite' },
                            { value: 'redis', label: 'Redis' },
                            { value: '', label: 'None' }
                          ].map((db) => {
                            const isSelected = config.features?.database?.driver === db.value
                            return (
                              <button
                                key={db.value}
                                onClick={() => {
                                  const newConfig = { ...config }
                                  if (db.value === '') {
                                    // Remove database feature
                                    if (newConfig.features) {
                                      delete newConfig.features.database
                                      if (Object.keys(newConfig.features).length === 0) {
                                        delete newConfig.features
                                      }
                                    }
                                  } else {
                                    // Add/update database feature
                                    if (!newConfig.features) newConfig.features = {}
                                    newConfig.features.database = {
                                      driver: db.value as any,
                                      orm: 'gorm' // Default ORM
                                    }
                                  }
                                  onConfigChange(newConfig)
                                }}
                                className={`px-3 py-2 text-sm border rounded-lg transition-colors ${
                                  isSelected
                                    ? 'border-purple-500 bg-purple-50 text-purple-700'
                                    : 'border-gray-300 hover:border-purple-300'
                                }`}
                              >
                                {db.label}
                              </button>
                            )
                          })}
                        </div>
                      </div>
                      <div>
                        <div className="flex items-center gap-2 mb-2">
                          <label className="label">ORM/Database Library</label>
                          <HelpTooltip 
                            content={getHelpContent('concepts', 'orm')}
                            variant="help"
                            position="right"
                          />
                        </div>
                        <div className="grid grid-cols-2 gap-2">
                          {[
                            { value: 'gorm', label: 'GORM' },
                            { value: 'sqlx', label: 'sqlx' },
                            { value: 'sqlc', label: 'sqlc' },
                            { value: 'ent', label: 'ent' }
                          ].map((orm) => {
                            const isSelected = config.features?.database?.orm === orm.value
                            const isDisabled = !config.features?.database?.driver
                            return (
                              <button
                                key={orm.value}
                                onClick={() => {
                                  if (!isDisabled && config.features?.database) {
                                    const newConfig = { ...config }
                                    newConfig.features!.database!.orm = orm.value as any
                                    onConfigChange(newConfig)
                                  }
                                }}
                                disabled={isDisabled}
                                className={`px-3 py-2 text-sm border rounded-lg transition-colors ${
                                  isDisabled
                                    ? 'border-gray-200 bg-gray-50 text-gray-400 cursor-not-allowed'
                                    : isSelected
                                    ? 'border-purple-500 bg-purple-50 text-purple-700'
                                    : 'border-gray-300 hover:border-purple-300'
                                }`}
                              >
                                {orm.label}
                              </button>
                            )
                          })}
                        </div>
                      </div>
                    </div>
                  </Disclosure.Panel>
                </>
              )}
            </Disclosure>

            <Disclosure>
              {({ open }) => (
                <>
                  <Disclosure.Button className="flex w-full justify-between rounded-lg bg-purple-50 px-4 py-2 text-left text-sm font-medium text-purple-900 hover:bg-purple-100 focus:outline-none focus-visible:ring focus-visible:ring-primary-500 focus-visible:ring-opacity-75">
                    <div className="flex items-center gap-2">
                      <span>Authentication & Security</span>
                      <HelpTooltip 
                        content="Add user authentication and security features to your application. Choose based on your authentication requirements."
                        variant="help"
                        position="right"
                        maxWidth="max-w-sm"
                      />
                    </div>
                    <ChevronDownIcon
                      className={`${open ? 'rotate-180 transform' : ''} h-5 w-5 text-purple-500`}
                    />
                  </Disclosure.Button>
                  <Disclosure.Panel className="px-4 pt-4 pb-2 space-y-4">
                    <div className="space-y-4">
                      <div>
                        <div className="flex items-center gap-2 mb-2">
                          <label className="label">Authentication Type</label>
                          <HelpTooltip 
                            content="Authentication method for securing your application. Choose based on your use case and security requirements."
                            variant="help"
                            position="right"
                          />
                        </div>
                        <div className="grid grid-cols-2 gap-2">
                          {[
                            { value: 'jwt', label: 'JWT' },
                            { value: 'oauth2', label: 'OAuth2' },
                            { value: 'session', label: 'Session' },
                            { value: 'api-key', label: 'API Key' },
                            { value: '', label: 'None' }
                          ].map((auth) => {
                            const isSelected = config.features?.authentication?.type === auth.value
                            return (
                              <button
                                key={auth.value}
                                onClick={() => {
                                  const newConfig = { ...config }
                                  if (auth.value === '') {
                                    // Remove authentication feature
                                    if (newConfig.features) {
                                      delete newConfig.features.authentication
                                      if (Object.keys(newConfig.features).length === 0) {
                                        delete newConfig.features
                                      }
                                    }
                                  } else {
                                    // Add/update authentication feature
                                    if (!newConfig.features) newConfig.features = {}
                                    newConfig.features.authentication = {
                                      type: auth.value as any,
                                      providers: [] // Default providers
                                    }
                                  }
                                  onConfigChange(newConfig)
                                }}
                                className={`px-3 py-2 text-sm border rounded-lg transition-colors ${
                                  isSelected
                                    ? 'border-purple-500 bg-purple-50 text-purple-700'
                                    : 'border-gray-300 hover:border-purple-300'
                                }`}
                              >
                                {auth.label}
                              </button>
                            )
                          })}
                        </div>
                      </div>
                    </div>
                  </Disclosure.Panel>
                </>
              )}
            </Disclosure>

            <Disclosure>
              {({ open }) => (
                <>
                  <Disclosure.Button className="flex w-full justify-between rounded-lg bg-purple-50 px-4 py-2 text-left text-sm font-medium text-purple-900 hover:bg-purple-100 focus:outline-none focus-visible:ring focus-visible:ring-primary-500 focus-visible:ring-opacity-75">
                    <div className="flex items-center gap-2">
                      <span>Deployment Options</span>
                      <HelpTooltip 
                        content={getHelpContent('concepts', 'deployment')}
                        variant="help"
                        position="right"
                        maxWidth="max-w-sm"
                      />
                    </div>
                    <ChevronDownIcon
                      className={`${open ? 'rotate-180 transform' : ''} h-5 w-5 text-purple-500`}
                    />
                  </Disclosure.Button>
                  <Disclosure.Panel className="px-4 pt-4 pb-2 space-y-4">
                    <div className="space-y-4">
                      <div>
                        <div className="flex items-center gap-2 mb-2">
                          <label className="label">Deployment Targets</label>
                          <HelpTooltip 
                            content={getHelpContent('concepts', 'deployment')}
                            variant="help"
                            position="right"
                          />
                        </div>
                        <div className="space-y-2">
                          {[
                            { id: 'docker', label: 'Docker', description: 'Dockerfile and docker-compose' },
                            { id: 'kubernetes', label: 'Kubernetes', description: 'K8s manifests and Helm charts' },
                            { id: 'serverless', label: 'Serverless', description: 'AWS Lambda, Vercel, etc.' },
                            { id: 'traditional', label: 'Traditional', description: 'SystemD, binary deployment' }
                          ].map((target) => (
                            <label key={target.id} className="flex items-start space-x-3 cursor-pointer">
                              <input type="checkbox" className="mt-1 rounded border-gray-300 text-purple-600 focus:ring-purple-500" />
                              <div>
                                <p className="text-sm font-medium text-gray-900">{target.label}</p>
                                <p className="text-xs text-gray-600">{target.description}</p>
                              </div>
                            </label>
                          ))}
                        </div>
                      </div>
                    </div>
                  </Disclosure.Panel>
                </>
              )}
            </Disclosure>
          </>
        )}

        {/* Action Buttons */}
        <div className="pt-4 border-t border-gray-200">
          <div className="flex space-x-3">
            <Button
              variant={isFormValid() ? 'primary' : 'secondary'}
              size="md"
              fullWidth
              disabled={!isFormValid()}
              onClick={() => {
                if (isFormValid()) {
                  console.log('🚀 Generating project with config:', config)
                  onGenerate({
                    projectName: config.projectName,
                    config: config,
                    options: {
                      memoryMode: true,
                      includeExamples: true
                    }
                  })
                }
              }}
              title={!isFormValid() ? 'Please fill in all required fields correctly' : 'Generate your Go project'}
              className={`
                ${isFormValid() 
                  ? 'bg-gradient-to-r from-green-600 to-blue-600 hover:from-green-700 hover:to-blue-700 text-white shadow-lg transform transition-all duration-200 hover:scale-105 active:scale-95' 
                  : ''
                }
              `}
            >
              {isFormValid() ? '🚀 Generate Project' : '⚠️ Complete Required Fields'}
            </Button>
            <Button
              variant="outline"
              size="md"
              onClick={() => {
                onConfigChange({
                  projectName: '',
                  moduleUrl: '',
                  goVersion: '1.21',
                  projectType: 'web-api',
                  framework: 'gin',
                  architecture: 'standard',
                  logger: 'slog',
                })
              }}
            >
              Reset
            </Button>
          </div>
        </div>
      </div>

      {/* Template Gallery Modal */}
      {showTemplateGallery && (
        <TemplateGallery
          onSelectTemplate={applyTemplate}
          onClose={() => setShowTemplateGallery(false)}
        />
      )}
    </div>
  )
}