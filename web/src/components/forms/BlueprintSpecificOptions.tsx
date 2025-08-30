/**
 * Blueprint-Specific Configuration Options
 * Dynamically renders configuration options based on selected blueprint
 */

import React, { useMemo } from 'react'
import { Disclosure } from '@headlessui/react'
import { ChevronDownIcon } from '@heroicons/react/20/solid'
import type { ProjectConfig, Blueprint } from '../../types'
import HelpTooltip from '../common/HelpTooltip'
import CompactSelectionCard from './CompactSelectionCard'

interface BlueprintSpecificOptionsProps {
  blueprint: Blueprint | null
  config: ProjectConfig
  onConfigChange: (config: ProjectConfig) => void
  isAdvanced?: boolean
}

interface OptionGroup {
  id: string
  title: string
  description: string
  options: ConfigOption[]
  condition?: (blueprint: Blueprint, config: ProjectConfig) => boolean
}

interface ConfigOption {
  id: string
  label: string
  description: string
  type: 'select' | 'multi-select' | 'boolean' | 'text' | 'number'
  options?: { value: string; label: string; description?: string }[]
  default?: any
  required?: boolean
  validation?: (value: any) => string | null
  helpContent?: string
}

export default function BlueprintSpecificOptions({
  blueprint,
  config,
  onConfigChange,
  isAdvanced = false
}: BlueprintSpecificOptionsProps) {
  // Define blueprint-specific option groups
  const optionGroups: OptionGroup[] = useMemo(() => [
    {
      id: 'web-api-features',
      title: 'Web API Features',
      description: 'Configure REST API specific features',
      condition: (bp) => bp.category === 'web-api' || bp.id.includes('web-api'),
      options: [
        {
          id: 'middleware',
          label: 'Middleware Stack',
          description: 'Select middleware components',
          type: 'multi-select',
          options: [
            { value: 'cors', label: 'CORS', description: 'Cross-Origin Resource Sharing' },
            { value: 'auth', label: 'Authentication', description: 'JWT/OAuth middleware' },
            { value: 'rate-limit', label: 'Rate Limiting', description: 'Request throttling' },
            { value: 'metrics', label: 'Metrics', description: 'Prometheus metrics' },
            { value: 'compression', label: 'Compression', description: 'Gzip compression' },
            { value: 'request-id', label: 'Request ID', description: 'Request tracing' }
          ]
        },
        {
          id: 'api-docs',
          label: 'API Documentation',
          description: 'Generate OpenAPI/Swagger documentation',
          type: 'select',
          options: [
            { value: 'swagger', label: 'Swagger/OpenAPI', description: 'Interactive API docs' },
            { value: 'redoc', label: 'ReDoc', description: 'Clean API documentation' },
            { value: 'none', label: 'None', description: 'No API documentation' }
          ],
          default: 'swagger'
        },
        {
          id: 'validation',
          label: 'Request Validation',
          description: 'Input validation strategy',
          type: 'select',
          options: [
            { value: 'validator', label: 'go-playground/validator', description: 'Tag-based validation' },
            { value: 'manual', label: 'Manual validation', description: 'Custom validation logic' },
            { value: 'schema', label: 'JSON Schema', description: 'Schema-based validation' }
          ],
          default: 'validator'
        }
      ]
    },
    {
      id: 'grpc-features',
      title: 'gRPC Configuration',
      description: 'Configure gRPC service features',
      condition: (bp) => bp.category === 'grpc' || bp.id.includes('grpc'),
      options: [
        {
          id: 'grpc-gateway',
          label: 'gRPC-Gateway',
          description: 'Enable HTTP/JSON API alongside gRPC',
          type: 'boolean',
          default: false,
          helpContent: 'Automatically generates REST API from gRPC definitions'
        },
        {
          id: 'reflection',
          label: 'gRPC Reflection',
          description: 'Enable server reflection for debugging',
          type: 'boolean',
          default: true,
          helpContent: 'Allows tools like grpcurl to introspect your service'
        },
        {
          id: 'interceptors',
          label: 'gRPC Interceptors',
          description: 'Select gRPC middleware interceptors',
          type: 'multi-select',
          options: [
            { value: 'logging', label: 'Logging', description: 'Request/response logging' },
            { value: 'metrics', label: 'Metrics', description: 'Prometheus metrics' },
            { value: 'auth', label: 'Authentication', description: 'JWT/OAuth validation' },
            { value: 'rate-limit', label: 'Rate Limiting', description: 'Request throttling' },
            { value: 'recovery', label: 'Panic Recovery', description: 'Error handling' },
            { value: 'tracing', label: 'Distributed Tracing', description: 'OpenTelemetry traces' }
          ]
        },
        {
          id: 'health-check',
          label: 'Health Check',
          description: 'gRPC health check service',
          type: 'boolean',
          default: true,
          helpContent: 'Implements standard gRPC health checking protocol'
        }
      ]
    },
    {
      id: 'cli-features',
      title: 'CLI Features',
      description: 'Configure command-line application features',
      condition: (bp) => bp.category === 'cli' || bp.id.includes('cli'),
      options: [
        {
          id: 'completion',
          label: 'Shell Completion',
          description: 'Generate shell completion scripts',
          type: 'multi-select',
          options: [
            { value: 'bash', label: 'Bash', description: 'Bash completion' },
            { value: 'zsh', label: 'Zsh', description: 'Zsh completion' },
            { value: 'fish', label: 'Fish', description: 'Fish completion' },
            { value: 'powershell', label: 'PowerShell', description: 'PowerShell completion' }
          ]
        },
        {
          id: 'config-format',
          label: 'Configuration Format',
          description: 'Default configuration file format',
          type: 'select',
          options: [
            { value: 'yaml', label: 'YAML', description: 'Human-readable format' },
            { value: 'json', label: 'JSON', description: 'Machine-readable format' },
            { value: 'toml', label: 'TOML', description: 'Configuration format' },
            { value: 'env', label: 'Environment', description: 'Environment variables only' }
          ],
          default: 'yaml'
        },
        {
          id: 'subcommands',
          label: 'Command Structure',
          description: 'CLI command organization',
          type: 'select',
          options: [
            { value: 'single', label: 'Single Command', description: 'Simple CLI with flags' },
            { value: 'multi', label: 'Multi-Command', description: 'Subcommand structure' },
            { value: 'plugin', label: 'Plugin System', description: 'Extensible plugin architecture' }
          ],
          default: 'multi'
        }
      ]
    },
    {
      id: 'lambda-features',
      title: 'Lambda Configuration',
      description: 'Configure AWS Lambda specific features',
      condition: (bp) => bp.category === 'lambda' || bp.id.includes('lambda'),
      options: [
        {
          id: 'runtime',
          label: 'Lambda Runtime',
          description: 'Go runtime version for Lambda',
          type: 'select',
          options: [
            { value: 'go1.x', label: 'Go 1.x (Amazon Linux)', description: 'Standard Go runtime' },
            { value: 'provided.al2', label: 'Custom Runtime (AL2)', description: 'Custom runtime on Amazon Linux 2' }
          ],
          default: 'provided.al2'
        },
        {
          id: 'triggers',
          label: 'Event Triggers',
          description: 'Lambda trigger configuration',
          type: 'multi-select',
          options: [
            { value: 'api-gateway', label: 'API Gateway', description: 'HTTP API triggers' },
            { value: 's3', label: 'S3 Events', description: 'S3 bucket events' },
            { value: 'dynamodb', label: 'DynamoDB', description: 'DynamoDB streams' },
            { value: 'sqs', label: 'SQS', description: 'SQS queue messages' },
            { value: 'sns', label: 'SNS', description: 'SNS topic messages' },
            { value: 'eventbridge', label: 'EventBridge', description: 'Custom events' }
          ]
        },
        {
          id: 'tracing',
          label: 'X-Ray Tracing',
          description: 'Enable AWS X-Ray distributed tracing',
          type: 'boolean',
          default: false,
          helpContent: 'Adds distributed tracing for Lambda functions'
        }
      ]
    },
    {
      id: 'microservice-features',
      title: 'Microservice Features',
      description: 'Configure microservice architecture features',
      condition: (bp) => bp.category === 'microservice' || bp.id.includes('microservice'),
      options: [
        {
          id: 'service-discovery',
          label: 'Service Discovery',
          description: 'Service registration and discovery',
          type: 'select',
          options: [
            { value: 'consul', label: 'HashiCorp Consul', description: 'Full-featured service mesh' },
            { value: 'etcd', label: 'etcd', description: 'Distributed key-value store' },
            { value: 'kubernetes', label: 'Kubernetes DNS', description: 'Native K8s discovery' },
            { value: 'none', label: 'None', description: 'No service discovery' }
          ],
          default: 'none'
        },
        {
          id: 'circuit-breaker',
          label: 'Circuit Breaker',
          description: 'Fault tolerance patterns',
          type: 'boolean',
          default: true,
          helpContent: 'Prevents cascade failures in distributed systems'
        },
        {
          id: 'monitoring',
          label: 'Observability Stack',
          description: 'Monitoring and observability tools',
          type: 'multi-select',
          options: [
            { value: 'prometheus', label: 'Prometheus', description: 'Metrics collection' },
            { value: 'jaeger', label: 'Jaeger', description: 'Distributed tracing' },
            { value: 'grafana', label: 'Grafana', description: 'Visualization dashboards' },
            { value: 'elk', label: 'ELK Stack', description: 'Log aggregation' }
          ]
        }
      ]
    },
    {
      id: 'database-advanced',
      title: 'Database Advanced Options',
      description: 'Advanced database configuration',
      condition: (bp, config) => !!config.features?.database?.driver && isAdvanced,
      options: [
        {
          id: 'migrations',
          label: 'Database Migrations',
          description: 'Database schema migration strategy',
          type: 'select',
          options: [
            { value: 'golang-migrate', label: 'golang-migrate', description: 'Popular migration tool' },
            { value: 'goose', label: 'Goose', description: 'Database migration tool' },
            { value: 'atlas', label: 'Atlas', description: 'Modern migration tool' },
            { value: 'manual', label: 'Manual', description: 'Custom migration logic' }
          ],
          default: 'golang-migrate'
        },
        {
          id: 'connection-pool',
          label: 'Connection Pooling',
          description: 'Database connection pool configuration',
          type: 'boolean',
          default: true,
          helpContent: 'Optimizes database connections for better performance'
        },
        {
          id: 'read-replicas',
          label: 'Read Replicas',
          description: 'Support for read replica databases',
          type: 'boolean',
          default: false,
          helpContent: 'Separates read and write database operations'
        },
        {
          id: 'caching',
          label: 'Database Caching',
          description: 'Query result caching strategy',
          type: 'select',
          options: [
            { value: 'redis', label: 'Redis', description: 'In-memory caching' },
            { value: 'memcached', label: 'Memcached', description: 'Distributed caching' },
            { value: 'in-memory', label: 'In-Memory', description: 'Local application cache' },
            { value: 'none', label: 'None', description: 'No caching' }
          ],
          default: 'none'
        }
      ]
    },
    {
      id: 'security-advanced',
      title: 'Security & Compliance',
      description: 'Advanced security and compliance features',
      condition: () => isAdvanced,
      options: [
        {
          id: 'security-headers',
          label: 'Security Headers',
          description: 'HTTP security headers middleware',
          type: 'multi-select',
          options: [
            { value: 'hsts', label: 'HSTS', description: 'HTTP Strict Transport Security' },
            { value: 'csp', label: 'CSP', description: 'Content Security Policy' },
            { value: 'nosniff', label: 'X-Content-Type-Options', description: 'MIME sniffing protection' },
            { value: 'frame-options', label: 'X-Frame-Options', description: 'Clickjacking protection' },
            { value: 'xss-protection', label: 'X-XSS-Protection', description: 'XSS filtering' }
          ]
        },
        {
          id: 'input-validation',
          label: 'Input Sanitization',
          description: 'Advanced input validation and sanitization',
          type: 'boolean',
          default: true,
          helpContent: 'Protects against injection attacks and malformed input'
        },
        {
          id: 'audit-logging',
          label: 'Audit Logging',
          description: 'Security event logging for compliance',
          type: 'boolean',
          default: false,
          helpContent: 'Logs security events for compliance and forensics'
        },
        {
          id: 'secrets-management',
          label: 'Secrets Management',
          description: 'Secure secrets handling',
          type: 'select',
          options: [
            { value: 'env', label: 'Environment Variables', description: 'Basic environment variables' },
            { value: 'vault', label: 'HashiCorp Vault', description: 'Enterprise secrets management' },
            { value: 'aws-secrets', label: 'AWS Secrets Manager', description: 'AWS-native secrets' },
            { value: 'k8s-secrets', label: 'Kubernetes Secrets', description: 'K8s native secrets' }
          ],
          default: 'env'
        }
      ]
    }
  ], [isAdvanced])

  // Filter option groups based on blueprint and conditions
  const relevantGroups = useMemo(() => {
    if (!blueprint) return []
    
    return optionGroups.filter(group => 
      !group.condition || group.condition(blueprint, config)
    )
  }, [blueprint, config, optionGroups])

  // Update configuration with new option value
  const updateOption = (groupId: string, optionId: string, value: any) => {
    const newConfig = { ...config }
    
    // Initialize features object if it doesn't exist
    if (!newConfig.features) {
      newConfig.features = {}
    }

    // Create nested structure for blueprint-specific options
    if (!newConfig.features.blueprintOptions) {
      newConfig.features.blueprintOptions = {}
    }

    if (!newConfig.features.blueprintOptions[groupId]) {
      newConfig.features.blueprintOptions[groupId] = {}
    }

    newConfig.features.blueprintOptions[groupId][optionId] = value
    onConfigChange(newConfig)
  }

  // Get current option value
  const getOptionValue = (groupId: string, optionId: string, defaultValue?: any) => {
    return config.features?.blueprintOptions?.[groupId]?.[optionId] ?? defaultValue
  }

  // Render different option types
  const renderOption = (group: OptionGroup, option: ConfigOption) => {
    const currentValue = getOptionValue(group.id, option.id, option.default)

    switch (option.type) {
      case 'select':
        return (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
            {option.options?.map((opt) => (
              <CompactSelectionCard
                key={opt.value}
                value={opt.value}
                label={opt.label}
                description={opt.description || ''}
                selected={currentValue === opt.value}
                onClick={() => updateOption(group.id, option.id, opt.value)}
              />
            ))}
          </div>
        )

      case 'multi-select':
        const selectedValues = Array.isArray(currentValue) ? currentValue : []
        return (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
            {option.options?.map((opt) => (
              <CompactSelectionCard
                key={opt.value}
                value={opt.value}
                label={opt.label}
                description={opt.description || ''}
                selected={selectedValues.includes(opt.value)}
                onClick={() => {
                  const newValues = selectedValues.includes(opt.value)
                    ? selectedValues.filter(v => v !== opt.value)
                    : [...selectedValues, opt.value]
                  updateOption(group.id, option.id, newValues)
                }}
              />
            ))}
          </div>
        )

      case 'boolean':
        return (
          <label className="flex items-center space-x-3 cursor-pointer">
            <input
              type="checkbox"
              checked={!!currentValue}
              onChange={(e) => updateOption(group.id, option.id, e.target.checked)}
              className="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
            />
            <div className="flex-1">
              <span className="text-sm font-medium text-gray-900">{option.label}</span>
              {option.helpContent && (
                <p className="text-xs text-gray-600 mt-1">{option.helpContent}</p>
              )}
            </div>
          </label>
        )

      case 'text':
        return (
          <input
            type="text"
            value={currentValue || ''}
            onChange={(e) => updateOption(group.id, option.id, e.target.value)}
            placeholder={option.label}
            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
          />
        )

      case 'number':
        return (
          <input
            type="number"
            value={currentValue || ''}
            onChange={(e) => updateOption(group.id, option.id, parseInt(e.target.value))}
            placeholder={option.label}
            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
          />
        )

      default:
        return null
    }
  }

  if (relevantGroups.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        <p>No blueprint-specific options available.</p>
        <p className="text-sm mt-1">Options will appear here based on your selected blueprint.</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div className="text-sm text-gray-600 mb-4">
        <p>Configure options specific to the <strong>{blueprint?.name}</strong> blueprint.</p>
      </div>

      {relevantGroups.map((group) => (
        <Disclosure key={group.id}>
          {({ open }) => (
            <>
              <Disclosure.Button 
                className="flex w-full justify-between rounded-lg bg-blue-50 px-4 py-3 text-left text-sm font-medium text-blue-900 hover:bg-blue-100 focus:outline-none focus-visible:ring focus-visible:ring-blue-500 focus-visible:ring-opacity-75"
                aria-expanded={open}
                aria-controls={`${group.id}-panel`}
              >
                <div>
                  <span>{group.title}</span>
                  <p className="text-xs text-blue-700 mt-1">{group.description}</p>
                </div>
                <ChevronDownIcon
                  className={`${open ? 'rotate-180 transform' : ''} h-5 w-5 text-blue-500 flex-shrink-0`}
                  aria-hidden="true"
                />
              </Disclosure.Button>
              <Disclosure.Panel 
                id={`${group.id}-panel`}
                className="px-4 pt-4 pb-2 space-y-6"
              >
                {group.options.map((option) => (
                  <div key={option.id}>
                    <div className="flex items-center gap-2 mb-3">
                      <label className="text-sm font-medium text-gray-900">
                        {option.label}
                        {option.required && <span className="text-red-500 ml-1">*</span>}
                      </label>
                      {option.helpContent && (
                        <HelpTooltip 
                          content={option.helpContent}
                          position="right"
                        />
                      )}
                    </div>
                    <p className="text-xs text-gray-600 mb-3">{option.description}</p>
                    {renderOption(group, option)}
                  </div>
                ))}
              </Disclosure.Panel>
            </>
          )}
        </Disclosure>
      ))}
    </div>
  )
}