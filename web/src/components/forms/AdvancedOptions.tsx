/**
 * Advanced Configuration Options Panel
 * Provides advanced configuration options for experienced users
 */

import React, { useState } from 'react'
import { Disclosure, Switch, Tab } from '@headlessui/react'
import { 
  ChevronDownIcon, 
  CogIcon,
  CloudIcon,
  ShieldCheckIcon,
  RocketLaunchIcon,
  WrenchScrewdriverIcon
} from '@heroicons/react/20/solid'
import type { ProjectConfig, Blueprint } from '../../types'
import HelpTooltip from '../common/HelpTooltip'
import CompactSelectionCard from './CompactSelectionCard'
import ValidatedInput from './ValidatedInput'
import { ValidationPatterns } from '../../utils/validation'

interface AdvancedOptionsProps {
  blueprint: Blueprint | null
  config: ProjectConfig
  onConfigChange: (config: ProjectConfig) => void
}

interface AdvancedFeature {
  id: string
  label: string
  description: string
  category: 'performance' | 'deployment' | 'security' | 'development' | 'testing'
  enabled: boolean
  options?: { [key: string]: any }
}

export default function AdvancedOptions({
  blueprint,
  config,
  onConfigChange
}: AdvancedOptionsProps) {
  // Advanced features state
  const [advancedFeatures, setAdvancedFeatures] = useState<AdvancedFeature[]>([
    {
      id: 'hot-reload',
      label: 'Hot Reload',
      description: 'Enable hot reloading for development',
      category: 'development',
      enabled: false,
      options: { watchFiles: ['*.go', '*.yaml'], excludeVendor: true }
    },
    {
      id: 'profiling',
      label: 'Runtime Profiling',
      description: 'Include pprof endpoints for performance analysis',
      category: 'performance',
      enabled: false,
      options: { endpoint: '/debug/pprof', enableCPU: true, enableMemory: true }
    },
    {
      id: 'graceful-shutdown',
      label: 'Graceful Shutdown',
      description: 'Handle shutdown signals gracefully',
      category: 'performance',
      enabled: true,
      options: { timeout: 30, signals: ['SIGTERM', 'SIGINT'] }
    },
    {
      id: 'health-checks',
      label: 'Health Check Endpoints',
      description: 'Kubernetes-compatible health check endpoints',
      category: 'deployment',
      enabled: true,
      options: { liveness: '/health/live', readiness: '/health/ready' }
    },
    {
      id: 'metrics-export',
      label: 'Metrics Export',
      description: 'Prometheus-compatible metrics endpoint',
      category: 'performance',
      enabled: false,
      options: { endpoint: '/metrics', namespace: '', customMetrics: true }
    },
    {
      id: 'distributed-tracing',
      label: 'Distributed Tracing',
      description: 'OpenTelemetry distributed tracing',
      category: 'performance',
      enabled: false,
      options: { provider: 'jaeger', endpoint: 'localhost:14268', samplingRatio: 0.1 }
    },
    {
      id: 'rate-limiting',
      label: 'Advanced Rate Limiting',
      description: 'Sophisticated rate limiting with different strategies',
      category: 'security',
      enabled: false,
      options: { strategy: 'sliding-window', rps: 100, burst: 200, storageType: 'memory' }
    },
    {
      id: 'request-validation',
      label: 'Request Validation',
      description: 'Comprehensive input validation and sanitization',
      category: 'security',
      enabled: true,
      options: { sanitizeHTML: true, validateJSON: true, maxBodySize: '10MB' }
    },
    {
      id: 'audit-logging',
      label: 'Audit Logging',
      description: 'Security audit trail logging',
      category: 'security',
      enabled: false,
      options: { logLevel: 'info', includeHeaders: false, includeBodies: false }
    },
    {
      id: 'container-optimized',
      label: 'Container Optimization',
      description: 'Optimize for containerized deployment',
      category: 'deployment',
      enabled: false,
      options: { multiStage: true, nonRoot: true, distroless: true, staticBinary: true }
    },
    {
      id: 'testing-framework',
      label: 'Enhanced Testing',
      description: 'Comprehensive testing setup with coverage reporting',
      category: 'testing',
      enabled: true,
      options: { framework: 'testify', coverage: true, benchmarks: true, fuzz: false }
    }
  ])

  const categories = [
    { id: 'performance', label: 'Performance', icon: RocketLaunchIcon, color: 'text-green-600' },
    { id: 'deployment', label: 'Deployment', icon: CloudIcon, color: 'text-blue-600' },
    { id: 'security', label: 'Security', icon: ShieldCheckIcon, color: 'text-red-600' },
    { id: 'development', label: 'Development', icon: WrenchScrewdriverIcon, color: 'text-purple-600' },
    { id: 'testing', label: 'Testing', icon: CogIcon, color: 'text-orange-600' }
  ]

  // Update feature state
  const toggleFeature = (featureId: string) => {
    setAdvancedFeatures(prev => 
      prev.map(feature => 
        feature.id === featureId 
          ? { ...feature, enabled: !feature.enabled }
          : feature
      )
    )
    
    // Update config
    const newConfig = { ...config }
    if (!newConfig.features) newConfig.features = {}
    if (!newConfig.features.advanced) newConfig.features.advanced = {}
    
    const feature = advancedFeatures.find(f => f.id === featureId)
    if (feature) {
      newConfig.features.advanced[featureId] = {
        enabled: !feature.enabled,
        options: feature.options
      }
      onConfigChange(newConfig)
    }
  }

  // Update feature options
  const updateFeatureOption = (featureId: string, optionKey: string, value: any) => {
    setAdvancedFeatures(prev => 
      prev.map(feature => 
        feature.id === featureId
          ? { 
              ...feature, 
              options: { ...feature.options, [optionKey]: value }
            }
          : feature
      )
    )

    // Update config
    const newConfig = { ...config }
    if (!newConfig.features) newConfig.features = {}
    if (!newConfig.features.advanced) newConfig.features.advanced = {}
    if (!newConfig.features.advanced[featureId]) {
      newConfig.features.advanced[featureId] = { enabled: false, options: {} }
    }
    
    newConfig.features.advanced[featureId].options = {
      ...newConfig.features.advanced[featureId].options,
      [optionKey]: value
    }
    
    onConfigChange(newConfig)
  }

  // Group features by category
  const featuresByCategory = categories.map(category => ({
    ...category,
    features: advancedFeatures.filter(feature => feature.category === category.id)
  }))

  // Render feature-specific options
  const renderFeatureOptions = (feature: AdvancedFeature) => {
    if (!feature.enabled || !feature.options) return null

    switch (feature.id) {
      case 'profiling':
        return (
          <div className="space-y-3 pl-6 border-l-2 border-gray-200">
            <ValidatedInput
              label="Profiling Endpoint"
              value={feature.options.endpoint || ''}
              onChange={(value) => updateFeatureOption(feature.id, 'endpoint', value)}
              placeholder="/debug/pprof"
              validation={{ pattern: ValidationPatterns.path }}
            />
            <div className="flex items-center space-x-4">
              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  checked={feature.options.enableCPU || false}
                  onChange={(e) => updateFeatureOption(feature.id, 'enableCPU', e.target.checked)}
                  className="rounded border-gray-300 text-blue-600"
                />
                <span className="text-sm text-gray-700">CPU Profiling</span>
              </label>
              <label className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  checked={feature.options.enableMemory || false}
                  onChange={(e) => updateFeatureOption(feature.id, 'enableMemory', e.target.checked)}
                  className="rounded border-gray-300 text-blue-600"
                />
                <span className="text-sm text-gray-700">Memory Profiling</span>
              </label>
            </div>
          </div>
        )

      case 'rate-limiting':
        return (
          <div className="space-y-3 pl-6 border-l-2 border-gray-200">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Strategy</label>
              <div className="grid grid-cols-2 gap-2">
                {['fixed-window', 'sliding-window', 'token-bucket'].map((strategy) => (
                  <CompactSelectionCard
                    key={strategy}
                    value={strategy}
                    label={strategy.split('-').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ')}
                    description=""
                    selected={feature.options.strategy === strategy}
                    onClick={() => updateFeatureOption(feature.id, 'strategy', strategy)}
                  />
                ))}
              </div>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <ValidatedInput
                label="Requests per Second"
                type="text"
                value={feature.options.rps || ''}
                onChange={(value) => updateFeatureOption(feature.id, 'rps', parseInt(value))}
                placeholder="100"
                />
              <ValidatedInput
                label="Burst Capacity"
                type="text"
                value={feature.options.burst || ''}
                onChange={(value) => updateFeatureOption(feature.id, 'burst', parseInt(value))}
                placeholder="200"
                />
            </div>
          </div>
        )

      case 'distributed-tracing':
        return (
          <div className="space-y-3 pl-6 border-l-2 border-gray-200">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Tracing Provider</label>
              <div className="grid grid-cols-2 gap-2">
                {['jaeger', 'zipkin', 'datadog', 'newrelic'].map((provider) => (
                  <CompactSelectionCard
                    key={provider}
                    value={provider}
                    label={provider.charAt(0).toUpperCase() + provider.slice(1)}
                    description=""
                    selected={feature.options.provider === provider}
                    onClick={() => updateFeatureOption(feature.id, 'provider', provider)}
                  />
                ))}
              </div>
            </div>
            <ValidatedInput
              label="Endpoint URL"
              value={feature.options.endpoint || ''}
              onChange={(value) => updateFeatureOption(feature.id, 'endpoint', value)}
              placeholder="localhost:14268"
            />
            <ValidatedInput
              label="Sampling Ratio (0.0-1.0)"
              type="text"
              value={feature.options.samplingRatio || ''}
              onChange={(value) => updateFeatureOption(feature.id, 'samplingRatio', parseFloat(value))}
              placeholder="0.1"
            />
          </div>
        )

      case 'container-optimized':
        return (
          <div className="space-y-3 pl-6 border-l-2 border-gray-200">
            <div className="grid grid-cols-2 gap-4">
              {[
                { key: 'multiStage', label: 'Multi-stage Build' },
                { key: 'nonRoot', label: 'Non-root User' },
                { key: 'distroless', label: 'Distroless Image' },
                { key: 'staticBinary', label: 'Static Binary' }
              ].map(({ key, label }) => (
                <label key={key} className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={feature.options[key] || false}
                    onChange={(e) => updateFeatureOption(feature.id, key, e.target.checked)}
                    className="rounded border-gray-300 text-blue-600"
                  />
                  <span className="text-sm text-gray-700">{label}</span>
                </label>
              ))}
            </div>
          </div>
        )

      case 'testing-framework':
        return (
          <div className="space-y-3 pl-6 border-l-2 border-gray-200">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">Testing Framework</label>
              <div className="grid grid-cols-2 gap-2">
                {['testify', 'ginkgo', 'goconvey', 'standard'].map((framework) => (
                  <CompactSelectionCard
                    key={framework}
                    value={framework}
                    label={framework === 'standard' ? 'Standard' : framework.charAt(0).toUpperCase() + framework.slice(1)}
                    description=""
                    selected={feature.options.framework === framework}
                    onClick={() => updateFeatureOption(feature.id, 'framework', framework)}
                  />
                ))}
              </div>
            </div>
            <div className="flex flex-wrap gap-4">
              {[
                { key: 'coverage', label: 'Coverage Reporting' },
                { key: 'benchmarks', label: 'Benchmark Tests' },
                { key: 'fuzz', label: 'Fuzz Testing' }
              ].map(({ key, label }) => (
                <label key={key} className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={feature.options[key] || false}
                    onChange={(e) => updateFeatureOption(feature.id, key, e.target.checked)}
                    className="rounded border-gray-300 text-blue-600"
                  />
                  <span className="text-sm text-gray-700">{label}</span>
                </label>
              ))}
            </div>
          </div>
        )

      default:
        return null
    }
  }

  return (
    <div className="space-y-6">
      <div className="bg-gradient-to-r from-purple-50 to-indigo-50 rounded-lg p-4 border border-purple-100">
        <div className="flex items-start gap-3">
          <CogIcon className="w-6 h-6 text-purple-600 flex-shrink-0 mt-1" />
          <div>
            <h3 className="font-semibold text-purple-900 mb-1">Advanced Configuration</h3>
            <p className="text-sm text-purple-700">
              Configure advanced features and optimizations for production environments.
              These options provide fine-grained control over your application's behavior.
            </p>
          </div>
        </div>
      </div>

      <Tab.Group>
        <Tab.List className="flex space-x-1 rounded-xl bg-gray-100 p-1">
          {categories.map((category) => (
            <Tab
              key={category.id}
              className={({ selected }) =>
                `w-full rounded-lg py-2.5 px-3 text-sm font-medium leading-5 transition-all
                 ${selected
                   ? 'bg-white shadow text-gray-900 ring-2 ring-blue-500 ring-opacity-60'
                   : 'text-gray-600 hover:bg-white/70 hover:text-gray-900'
                 }`
              }
            >
              <div className="flex items-center justify-center gap-2">
                <category.icon className={`w-4 h-4 ${category.color}`} />
                <span className="hidden sm:inline">{category.label}</span>
              </div>
            </Tab>
          ))}
        </Tab.List>

        <Tab.Panels className="mt-6">
          {featuresByCategory.map((category) => (
            <Tab.Panel key={category.id} className="space-y-4">
              <div className="mb-4">
                <h4 className="text-lg font-medium text-gray-900 flex items-center gap-2">
                  <category.icon className={`w-5 h-5 ${category.color}`} />
                  {category.label} Features
                </h4>
                <p className="text-sm text-gray-600 mt-1">
                  Configure {category.label.toLowerCase()} related features for your application.
                </p>
              </div>

              {category.features.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  <category.icon className={`w-12 h-12 mx-auto ${category.color} opacity-50 mb-3`} />
                  <p>No {category.label.toLowerCase()} features available yet.</p>
                  <p className="text-sm mt-1">More features coming soon!</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {category.features.map((feature) => (
                    <div key={feature.id} className="border border-gray-200 rounded-lg p-4 hover:border-gray-300 transition-colors">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-3 mb-2">
                            <h5 className="font-medium text-gray-900">{feature.label}</h5>
                            <Switch
                              checked={feature.enabled}
                              onChange={() => toggleFeature(feature.id)}
                              className={`${
                                feature.enabled ? 'bg-blue-600' : 'bg-gray-200'
                              } relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2`}
                            >
                              <span
                                className={`${
                                  feature.enabled ? 'translate-x-6' : 'translate-x-1'
                                } inline-block h-4 w-4 transform rounded-full bg-white transition-transform`}
                              />
                            </Switch>
                          </div>
                          <p className="text-sm text-gray-600 mb-3">{feature.description}</p>
                          {renderFeatureOptions(feature)}
                        </div>
                        <HelpTooltip 
                          content={`${feature.description}. This feature affects the ${feature.category} aspects of your application.`}
                          variant="info"
                          position="left"
                        />
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </Tab.Panel>
          ))}
        </Tab.Panels>
      </Tab.Group>

      <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <div className="flex items-start gap-3">
          <svg className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
          </svg>
          <div>
            <h4 className="font-medium text-yellow-800">Advanced Features Notice</h4>
            <p className="text-sm text-yellow-700 mt-1">
              These advanced features add complexity to your project. Enable only features you need 
              and understand. Some features may require additional infrastructure or configuration.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}