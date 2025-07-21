import { useState } from 'react'
import { Disclosure } from '@headlessui/react'
import { ChevronDownIcon, InformationCircleIcon } from '@heroicons/react/20/solid'
import type { DisclosureMode, ProjectConfig, ProjectType, Architecture, Framework, LoggerType } from '../../types'

interface ConfigurationPanelProps {
  disclosureMode: DisclosureMode
}

export default function ConfigurationPanel({ disclosureMode }: ConfigurationPanelProps) {
  const [config, setConfig] = useState<ProjectConfig>({
    projectName: '',
    moduleUrl: '',
    goVersion: '1.21',
    projectType: 'web-api',
    framework: 'gin',
    architecture: 'standard',
    logger: 'slog',
  })

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
    setConfig(prev => ({ ...prev, [field]: value }))
  }

  return (
    <div className="card h-full overflow-y-auto">
      <div className="mb-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-2">Project Configuration</h2>
        <p className="text-sm text-gray-600">
          {disclosureMode === 'basic' 
            ? 'Configure essential project settings' 
            : 'Configure all available project options'
          }
        </p>
      </div>

      <div className="space-y-6">
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
                {/* Project Name */}
                <div>
                  <label className="label">
                    Project Name
                    <InformationCircleIcon className="inline w-4 h-4 ml-1 text-gray-400" title="Name of your Go project" />
                  </label>
                  <input
                    type="text"
                    className="input"
                    placeholder="my-awesome-project"
                    value={config.projectName}
                    onChange={(e) => updateConfig('projectName', e.target.value)}
                  />
                </div>

                {/* Module URL */}
                <div>
                  <label className="label">
                    Module URL
                    <InformationCircleIcon className="inline w-4 h-4 ml-1 text-gray-400" title="Go module path (e.g., github.com/user/project)" />
                  </label>
                  <input
                    type="text"
                    className="input"
                    placeholder="github.com/user/my-awesome-project"
                    value={config.moduleUrl}
                    onChange={(e) => updateConfig('moduleUrl', e.target.value)}
                  />
                </div>

                {/* Project Type */}
                <div>
                  <label className="label">Project Type</label>
                  <select
                    className="input"
                    value={config.projectType}
                    onChange={(e) => updateConfig('projectType', e.target.value as ProjectType)}
                  >
                    {projectTypes.map((type) => (
                      <option key={type.value} value={type.value}>
                        {type.label} - {type.description}
                      </option>
                    ))}
                  </select>
                </div>

                {/* Go Version */}
                <div>
                  <label className="label">Go Version</label>
                  <select
                    className="input"
                    value={config.goVersion}
                    onChange={(e) => updateConfig('goVersion', e.target.value)}
                  >
                    <option value="1.21">Go 1.21</option>
                    <option value="1.20">Go 1.20</option>
                    <option value="1.19">Go 1.19</option>
                  </select>
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
                {/* Framework */}
                {(config.projectType === 'web-api' || config.projectType === 'microservice') && (
                  <div>
                    <label className="label">Framework</label>
                    <select
                      className="input"
                      value={config.framework}
                      onChange={(e) => updateConfig('framework', e.target.value as Framework)}
                    >
                      {frameworks.map((framework) => (
                        <option key={framework.value} value={framework.value}>
                          {framework.label} - {framework.description}
                        </option>
                      ))}
                    </select>
                  </div>
                )}

                {/* Architecture */}
                <div>
                  <label className="label">Architecture Pattern</label>
                  <select
                    className="input"
                    value={config.architecture}
                    onChange={(e) => updateConfig('architecture', e.target.value as Architecture)}
                  >
                    {architectures.map((arch) => (
                      <option key={arch.value} value={arch.value}>
                        {arch.label} - {arch.description}
                      </option>
                    ))}
                  </select>
                </div>

                {/* Logger */}
                <div>
                  <label className="label">Logger Type</label>
                  <select
                    className="input"
                    value={config.logger}
                    onChange={(e) => updateConfig('logger', e.target.value as LoggerType)}
                  >
                    {loggers.map((logger) => (
                      <option key={logger.value} value={logger.value}>
                        {logger.label} - {logger.description}
                      </option>
                    ))}
                  </select>
                </div>
              </Disclosure.Panel>
            </>
          )}
        </Disclosure>

        {/* Advanced Configuration (only in advanced mode) */}
        {disclosureMode === 'advanced' && (
          <Disclosure>
            {({ open }) => (
              <>
                <Disclosure.Button className="flex w-full justify-between rounded-lg bg-gray-50 px-4 py-2 text-left text-sm font-medium text-gray-900 hover:bg-gray-100 focus:outline-none focus-visible:ring focus-visible:ring-primary-500 focus-visible:ring-opacity-75">
                  <span>Advanced Settings</span>
                  <ChevronDownIcon
                    className={`${open ? 'rotate-180 transform' : ''} h-5 w-5 text-gray-500`}
                  />
                </Disclosure.Button>
                <Disclosure.Panel className="px-4 pt-4 pb-2 space-y-4">
                  <div className="text-sm text-gray-600 bg-blue-50 p-3 rounded-lg">
                    <p className="font-medium text-blue-800 mb-1">Advanced Features</p>
                    <p>Database configuration, authentication, and deployment options will be available here.</p>
                  </div>
                </Disclosure.Panel>
              </>
            )}
          </Disclosure>
        )}

        {/* Action Buttons */}
        <div className="pt-4 border-t border-gray-200">
          <div className="flex space-x-3">
            <button className="btn-primary flex-1">
              Generate Project
            </button>
            <button className="btn-secondary">
              Reset
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}