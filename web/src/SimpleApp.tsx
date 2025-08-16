import { useState, useEffect } from 'react'
import { 
  SparklesIcon, 
  CogIcon, 
  EyeIcon, 
  DocumentIcon,
  PlayIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
  ArrowPathIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  CodeBracketIcon,
  ClipboardDocumentIcon,
  ArrowDownTrayIcon,
  AdjustmentsHorizontalIcon,
  LightBulbIcon,
  CommandLineIcon,
  GlobeAltIcon,
  BookOpenIcon,
  CloudIcon
} from '@heroicons/react/24/outline'
import { 
  CheckCircleIcon as CheckCircleIconSolid,
  ExclamationTriangleIcon as ExclamationTriangleIconSolid 
} from '@heroicons/react/24/solid'

interface ProjectConfig {
  projectName: string
  moduleUrl: string
  projectType: string
  architecture: string
  framework: string
  logger: string
  goVersion: string
  database: {
    enabled: boolean
    drivers: string[]
  }
  authentication: {
    enabled: boolean
    methods: string[]
  }
}

interface FileNode {
  name: string
  type: 'file' | 'directory'
  path: string
  size?: number
  language?: string
  children?: FileNode[]
}

function SimpleApp() {
  const [disclosureMode, setDisclosureMode] = useState<'basic' | 'advanced'>('basic')
  const [config, setConfig] = useState<ProjectConfig>({
    projectName: '',
    moduleUrl: '',
    projectType: 'web-api',
    architecture: 'standard',
    framework: 'gin',
    logger: 'slog',
    goVersion: '1.21',
    database: {
      enabled: false,
      drivers: []
    },
    authentication: {
      enabled: false,
      methods: []
    }
  })
  const [generationStatus, setGenerationStatus] = useState<'idle' | 'generating' | 'completed' | 'error'>('idle')
  const [generationProgress, setGenerationProgress] = useState(0)
  const [expandedFiles, setExpandedFiles] = useState<Set<string>>(new Set())
  const [selectedFile, setSelectedFile] = useState<string | null>(null)
  const [validationErrors, setValidationErrors] = useState<string[]>([])

  // Sample file structure for preview - dynamically generated based on configuration
  const generateFileStructure = (): FileNode[] => {
    const baseFiles: FileNode[] = [
      { name: 'main.go', type: 'file', path: '/main.go', size: 1024, language: 'go' },
      { name: 'go.mod', type: 'file', path: '/go.mod', size: 256, language: 'mod' },
      { name: 'go.sum', type: 'file', path: '/go.sum', size: 512, language: 'sum' },
      { name: 'README.md', type: 'file', path: '/README.md', size: 2048, language: 'markdown' },
      { name: 'Dockerfile', type: 'file', path: '/Dockerfile', size: 512, language: 'dockerfile' }
    ]

    const cmdChildren: FileNode[] = [
      { name: 'server.go', type: 'file', path: '/cmd/server.go', size: 2048, language: 'go' }
    ]

    const internalChildren: FileNode[] = [
      { name: 'config.go', type: 'file', path: '/internal/config.go', size: 1024, language: 'go' },
      { name: 'handler.go', type: 'file', path: '/internal/handler.go', size: 3072, language: 'go' },
      { name: 'middleware.go', type: 'file', path: '/internal/middleware.go', size: 1536, language: 'go' },
      { name: 'logger.go', type: 'file', path: '/internal/logger.go', size: 1024, language: 'go' }
    ]

    // Add database files if database is enabled
    if (config.database.enabled && config.database.drivers.length > 0) {
      internalChildren.push({
        name: 'database',
        type: 'directory',
        path: '/internal/database',
        children: [
          { name: 'connection.go', type: 'file', path: '/internal/database/connection.go', size: 1536, language: 'go' },
          { name: 'migrations.go', type: 'file', path: '/internal/database/migrations.go', size: 2048, language: 'go' },
          { name: 'models.go', type: 'file', path: '/internal/database/models.go', size: 2560, language: 'go' }
        ]
      })
      
      // Add driver-specific files
      config.database.drivers.forEach(driver => {
        const dbDir = internalChildren.find(child => child.name === 'database')
        if (dbDir && dbDir.children) {
          dbDir.children.push({
            name: `${driver}.go`,
            type: 'file',
            path: `/internal/database/${driver}.go`,
            size: 1024,
            language: 'go'
          })
        }
      })
    }

    // Add authentication files if authentication is enabled
    if (config.authentication.enabled && config.authentication.methods.length > 0) {
      internalChildren.push({
        name: 'auth',
        type: 'directory',
        path: '/internal/auth',
        children: [
          { name: 'middleware.go', type: 'file', path: '/internal/auth/middleware.go', size: 2048, language: 'go' },
          { name: 'service.go', type: 'file', path: '/internal/auth/service.go', size: 3072, language: 'go' }
        ]
      })

      // Add method-specific files
      config.authentication.methods.forEach(method => {
        const authDir = internalChildren.find(child => child.name === 'auth')
        if (authDir && authDir.children) {
          authDir.children.push({
            name: `${method}.go`,
            type: 'file',
            path: `/internal/auth/${method}.go`,
            size: 1536,
            language: 'go'
          })
        }
      })
    }

    const testsChildren: FileNode[] = [
      { name: 'main_test.go', type: 'file', path: '/tests/main_test.go', size: 1024, language: 'go' },
      { name: 'integration_test.go', type: 'file', path: '/tests/integration_test.go', size: 2048, language: 'go' }
    ]

    // Add database tests if database is enabled
    if (config.database.enabled) {
      testsChildren.push({
        name: 'database_test.go',
        type: 'file',
        path: '/tests/database_test.go',
        size: 1536,
        language: 'go'
      })
    }

    // Add auth tests if authentication is enabled
    if (config.authentication.enabled) {
      testsChildren.push({
        name: 'auth_test.go',
        type: 'file',
        path: '/tests/auth_test.go',
        size: 2048,
        language: 'go'
      })
    }

    return [
      {
        name: config.projectName || 'my-go-project',
        type: 'directory',
        path: '/',
        children: [
          ...baseFiles,
          {
            name: 'cmd',
            type: 'directory',
            path: '/cmd',
            children: cmdChildren
          },
          {
            name: 'internal',
            type: 'directory',
            path: '/internal',
            children: internalChildren
          },
          {
            name: 'tests',
            type: 'directory',
            path: '/tests',
            children: testsChildren
          }
        ]
      }
    ]
  }

  const fileStructure: FileNode[] = generateFileStructure()

  // Validation
  useEffect(() => {
    const errors: string[] = []
    if (!config.projectName) errors.push('Project name is required')
    if (config.projectName && !/^[a-zA-Z][\w-]*$/.test(config.projectName)) {
      errors.push('Project name must start with a letter and contain only letters, numbers, hyphens, and underscores')
    }
    if (!config.moduleUrl) errors.push('Module URL is required')
    if (config.moduleUrl && !/^[a-zA-Z0-9][a-zA-Z0-9.-]*[a-zA-Z0-9]\/[a-zA-Z0-9][a-zA-Z0-9.-]*\/[a-zA-Z][\w-]*$/.test(config.moduleUrl)) {
      errors.push('Module URL must be a valid Go module path (e.g., github.com/user/project)')
    }
    setValidationErrors(errors)
  }, [config])

  const isFormValid = validationErrors.length === 0 && config.projectName && config.moduleUrl

  const handleGenerate = async () => {
    if (!isFormValid) return
    
    setGenerationStatus('generating')
    setGenerationProgress(0)
    
    // Simulate generation progress
    const steps = [
      { message: 'Creating project structure...', progress: 20 },
      { message: 'Generating Go modules...', progress: 40 },
      { message: 'Building application code...', progress: 60 },
      { message: 'Adding configuration files...', progress: 80 },
      { message: 'Finalizing project...', progress: 100 }
    ]
    
    for (const step of steps) {
      await new Promise(resolve => setTimeout(resolve, 800))
      setGenerationProgress(step.progress)
    }
    
    setGenerationStatus('completed')
  }

  const toggleFileExpansion = (path: string) => {
    const newExpanded = new Set(expandedFiles)
    if (newExpanded.has(path)) {
      newExpanded.delete(path)
    } else {
      newExpanded.add(path)
    }
    setExpandedFiles(newExpanded)
  }

  const renderFileTree = (nodes: FileNode[], depth = 0) => {
    return nodes.map((node) => {
      const isExpanded = expandedFiles.has(node.path)
      const isSelected = selectedFile === node.path
      
      return (
        <div key={node.path}>
          <div 
            className={`flex items-center py-1 px-2 rounded cursor-pointer hover:bg-gray-100 transition-colors ${
              isSelected ? 'bg-blue-50 border-l-2 border-blue-500' : ''
            }`}
            style={{ paddingLeft: `${depth * 16 + 8}px` }}
            onClick={() => {
              if (node.type === 'directory') {
                toggleFileExpansion(node.path)
              } else {
                setSelectedFile(node.path)
              }
            }}
          >
            {node.type === 'directory' && (
              <div className="mr-1">
                {isExpanded ? (
                  <ChevronDownIcon className="w-4 h-4 text-gray-500" />
                ) : (
                  <ChevronRightIcon className="w-4 h-4 text-gray-500" />
                )}
              </div>
            )}
            <div className="mr-2">
              {node.type === 'directory' ? (
                <span className="text-blue-500">📁</span>
              ) : (
                <span className="text-gray-500">
                  {node.language === 'go' ? '🔧' : 
                   node.language === 'markdown' ? '📝' : 
                   node.language === 'dockerfile' ? '🐳' : '📄'}
                </span>
              )}
            </div>
            <span className="text-sm font-mono text-gray-700 flex-1">{node.name}</span>
            {node.size && (
              <span className="text-xs text-gray-500 ml-2">
                {(node.size / 1024).toFixed(1)}KB
              </span>
            )}
          </div>
          {node.type === 'directory' && isExpanded && node.children && (
            <div>
              {renderFileTree(node.children, depth + 1)}
            </div>
          )}
        </div>
      )
    })
  }

  const projectTypes = [
    { value: 'web-api', label: 'Web API', icon: GlobeAltIcon, description: 'REST API server with various architectures' },
    { value: 'cli', label: 'CLI Tool', icon: CommandLineIcon, description: 'Command-line tools with Cobra framework' },
    { value: 'library', label: 'Library', icon: BookOpenIcon, description: 'Reusable Go packages' },
    { value: 'lambda', label: 'Lambda', icon: CloudIcon, description: 'AWS Lambda serverless functions' }
  ]

  const architectures = [
    { value: 'standard', label: 'Standard', description: 'Simple layered architecture' },
    { value: 'clean', label: 'Clean Architecture', description: 'Layered with dependency inversion' },
    { value: 'ddd', label: 'Domain-Driven Design', description: 'Domain-focused approach' },
    { value: 'hexagonal', label: 'Hexagonal', description: 'Ports and adapters pattern' }
  ]

  const frameworks = [
    { value: 'gin', label: 'Gin', description: 'Fast HTTP framework' },
    { value: 'echo', label: 'Echo', description: 'High performance framework' },
    { value: 'fiber', label: 'Fiber', description: 'Express-inspired framework' },
    { value: 'chi', label: 'Chi', description: 'Lightweight router' }
  ]

  const loggers = [
    { value: 'slog', label: 'slog', description: 'Standard library structured logging' },
    { value: 'zap', label: 'Zap', description: 'High-performance logger' },
    { value: 'logrus', label: 'Logrus', description: 'Popular structured logger' },
    { value: 'zerolog', label: 'Zerolog', description: 'Zero allocation logger' }
  ]

  const databaseOptions = [
    { value: 'postgresql', label: 'PostgreSQL', description: 'Powerful relational database', icon: '🐘' },
    { value: 'mysql', label: 'MySQL', description: 'Popular relational database', icon: '🐬' },
    { value: 'mongodb', label: 'MongoDB', description: 'Document-based NoSQL database', icon: '🍃' },
    { value: 'sqlite', label: 'SQLite', description: 'Lightweight embedded database', icon: '🪶' }
  ]

  const authenticationOptions = [
    { value: 'jwt', label: 'JWT', description: 'JSON Web Tokens', icon: '🔑' },
    { value: 'oauth2', label: 'OAuth2', description: 'Third-party authentication', icon: '🔐' },
    { value: 'session', label: 'Session', description: 'Server-side sessions', icon: '🎫' },
    { value: 'apikey', label: 'API Key', description: 'Simple API key authentication', icon: '🗝️' }
  ]

  // Helper functions for advanced options
  const toggleDatabaseOption = (option: string) => {
    setConfig(prev => {
      const drivers = [...prev.database.drivers]
      const index = drivers.indexOf(option)
      
      if (index > -1) {
        drivers.splice(index, 1)
      } else {
        drivers.push(option)
      }
      
      return {
        ...prev,
        database: {
          enabled: drivers.length > 0,
          drivers
        }
      }
    })
  }

  const toggleAuthenticationOption = (option: string) => {
    setConfig(prev => {
      const methods = [...prev.authentication.methods]
      const index = methods.indexOf(option)
      
      if (index > -1) {
        methods.splice(index, 1)
      } else {
        methods.push(option)
      }
      
      return {
        ...prev,
        authentication: {
          enabled: methods.length > 0,
          methods
        }
      }
    })
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 via-blue-50 to-indigo-50">
      {/* Modern Header */}
      <div className="bg-white/80 backdrop-blur-sm border-b border-gray-200/50 sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center">
                <SparklesIcon className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="text-2xl font-bold bg-gradient-to-r from-gray-900 to-gray-600 bg-clip-text text-transparent">
                  Go Starter Studio
                </h1>
                <p className="text-sm text-gray-600">Professional Go project generator</p>
              </div>
            </div>
            
            {/* Mode Toggle */}
            <div className="flex items-center space-x-2">
              <button
                onClick={() => setDisclosureMode(disclosureMode === 'basic' ? 'advanced' : 'basic')}
                className={`flex items-center space-x-2 px-4 py-2 rounded-lg border transition-all ${
                  disclosureMode === 'advanced' 
                    ? 'bg-purple-100 border-purple-300 text-purple-700' 
                    : 'bg-gray-100 border-gray-300 text-gray-700 hover:bg-gray-200'
                }`}
              >
                <AdjustmentsHorizontalIcon className="w-4 h-4" />
                <span className="text-sm font-medium">
                  {disclosureMode === 'basic' ? 'Basic Mode' : 'Advanced Mode'}
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          
          {/* Configuration Panel */}
          <div className="lg:col-span-4">
            <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-6 sticky top-24">
              <div className="flex items-center space-x-3 mb-6">
                <div className="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center">
                  <CogIcon className="w-5 h-5 text-blue-600" />
                </div>
                <div>
                  <h2 className="text-lg font-semibold text-gray-900">Configuration</h2>
                  <p className="text-sm text-gray-600">
                    {disclosureMode === 'basic' ? 'Essential settings' : 'All available options'}
                  </p>
                </div>
              </div>

              {/* Validation Summary */}
              {validationErrors.length > 0 && (
                <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
                  <div className="flex items-center space-x-2 mb-2">
                    <ExclamationTriangleIconSolid className="w-5 h-5 text-red-500" />
                    <span className="text-sm font-medium text-red-800">Please fix the following errors:</span>
                  </div>
                  <ul className="space-y-1">
                    {validationErrors.map((error, index) => (
                      <li key={index} className="text-sm text-red-700">• {error}</li>
                    ))}
                  </ul>
                </div>
              )}

              {isFormValid && (
                <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
                  <div className="flex items-center space-x-2">
                    <CheckCircleIconSolid className="w-5 h-5 text-green-500" />
                    <span className="text-sm font-medium text-green-800">Configuration is valid</span>
                  </div>
                </div>
              )}

              <div className="space-y-6">
                {/* Basic Settings */}
                <div>
                  <h3 className="text-sm font-semibold text-gray-900 mb-4 flex items-center">
                    <span>Basic Settings</span>
                    <span className="ml-2 w-2 h-2 bg-blue-500 rounded-full"></span>
                  </h3>
                  
                  <div className="space-y-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Project Name</label>
                      <input
                        type="text"
                        value={config.projectName}
                        onChange={(e) => setConfig({ ...config, projectName: e.target.value })}
                        placeholder="my-awesome-project"
                        className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Module URL</label>
                      <input
                        type="text"
                        value={config.moduleUrl}
                        onChange={(e) => setConfig({ ...config, moduleUrl: e.target.value })}
                        placeholder="github.com/user/my-awesome-project"
                        className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-3">Project Type</label>
                      <div className="grid grid-cols-2 gap-2">
                        {projectTypes.map((type) => {
                          const IconComponent = type.icon
                          return (
                            <button
                              key={type.value}
                              onClick={() => setConfig({ ...config, projectType: type.value })}
                              className={`p-3 border rounded-lg text-left transition-all hover:shadow-md ${
                                config.projectType === type.value
                                  ? 'border-blue-500 bg-blue-50 shadow-md'
                                  : 'border-gray-200 hover:border-gray-300'
                              }`}
                            >
                              <div className="flex items-center space-x-2 mb-1">
                                <IconComponent className="w-4 h-4 text-gray-600" />
                                <span className="font-medium text-sm">{type.label}</span>
                              </div>
                              <p className="text-xs text-gray-500">{type.description}</p>
                            </button>
                          )
                        })}
                      </div>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">Go Version</label>
                      <div className="flex space-x-2">
                        {['1.21', '1.20', '1.19'].map((version) => (
                          <button
                            key={version}
                            onClick={() => setConfig({ ...config, goVersion: version })}
                            className={`px-4 py-2 rounded-lg border font-medium text-sm transition-all ${
                              config.goVersion === version
                                ? 'border-blue-500 bg-blue-50 text-blue-700'
                                : 'border-gray-300 bg-white text-gray-700 hover:border-gray-400'
                            }`}
                          >
                            Go {version}
                          </button>
                        ))}
                      </div>
                    </div>
                  </div>
                </div>

                {/* Framework Settings */}
                {(config.projectType === 'web-api' || disclosureMode === 'advanced') && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-900 mb-4 flex items-center">
                      <span>Framework & Architecture</span>
                      <span className="ml-2 w-2 h-2 bg-purple-500 rounded-full"></span>
                    </h3>
                    
                    <div className="space-y-4">
                      {config.projectType === 'web-api' && (
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-2">Framework</label>
                          <select
                            value={config.framework}
                            onChange={(e) => setConfig({ ...config, framework: e.target.value })}
                            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          >
                            {frameworks.map((framework) => (
                              <option key={framework.value} value={framework.value}>
                                {framework.label} - {framework.description}
                              </option>
                            ))}
                          </select>
                        </div>
                      )}

                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Architecture</label>
                        <select
                          value={config.architecture}
                          onChange={(e) => setConfig({ ...config, architecture: e.target.value })}
                          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                          {architectures.map((arch) => (
                            <option key={arch.value} value={arch.value}>
                              {arch.label} - {arch.description}
                            </option>
                          ))}
                        </select>
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">Logger</label>
                        <select
                          value={config.logger}
                          onChange={(e) => setConfig({ ...config, logger: e.target.value })}
                          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        >
                          {loggers.map((logger) => (
                            <option key={logger.value} value={logger.value}>
                              {logger.label} - {logger.description}
                            </option>
                          ))}
                        </select>
                      </div>
                    </div>
                  </div>
                )}

                {/* Advanced Settings */}
                {disclosureMode === 'advanced' && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-900 mb-4 flex items-center">
                      <span>Advanced Options</span>
                      <span className="ml-2 w-2 h-2 bg-purple-500 rounded-full"></span>
                    </h3>
                    
                    <div className="space-y-4">
                      {/* Database Options */}
                      <div className="p-4 bg-purple-50 border border-purple-200 rounded-lg">
                        <div className="flex items-center justify-between mb-3">
                          <h4 className="text-sm font-medium text-purple-900 flex items-center">
                            🗄️ Database
                            {config.database.enabled && (
                              <span className="ml-2 px-2 py-1 text-xs bg-purple-200 text-purple-800 rounded-full">
                                {config.database.drivers.length} selected
                              </span>
                            )}
                          </h4>
                          {config.database.enabled && (
                            <button
                              onClick={() => setConfig(prev => ({
                                ...prev,
                                database: { enabled: false, drivers: [] }
                              }))}
                              className="text-xs text-purple-600 hover:text-purple-800 underline"
                            >
                              Clear all
                            </button>
                          )}
                        </div>
                        <div className="grid grid-cols-2 gap-2">
                          {databaseOptions.map((db) => {
                            const isSelected = config.database.drivers.includes(db.value)
                            return (
                              <button
                                key={db.value}
                                onClick={() => toggleDatabaseOption(db.value)}
                                className={`px-3 py-3 text-sm border rounded-lg transition-all duration-200 text-left ${
                                  isSelected
                                    ? 'border-purple-500 bg-purple-100 text-purple-900 shadow-md transform scale-105'
                                    : 'border-purple-300 bg-white text-purple-700 hover:bg-purple-50 hover:border-purple-400'
                                }`}
                              >
                                <div className="flex items-center space-x-2">
                                  <span className="text-lg">{db.icon}</span>
                                  <div className="flex-1">
                                    <div className="font-medium flex items-center">
                                      {db.label}
                                      {isSelected && (
                                        <CheckCircleIconSolid className="w-4 h-4 text-purple-600 ml-1" />
                                      )}
                                    </div>
                                    <div className="text-xs text-purple-600 mt-1">{db.description}</div>
                                  </div>
                                </div>
                              </button>
                            )
                          })}
                        </div>
                        {config.database.enabled && (
                          <div className="mt-3 p-2 bg-purple-100 rounded-md">
                            <p className="text-xs text-purple-800">
                              💡 This will add database connection, migrations, and model files to your project
                            </p>
                          </div>
                        )}
                      </div>

                      {/* Authentication Options */}
                      <div className="p-4 bg-indigo-50 border border-indigo-200 rounded-lg">
                        <div className="flex items-center justify-between mb-3">
                          <h4 className="text-sm font-medium text-indigo-900 flex items-center">
                            🔐 Authentication
                            {config.authentication.enabled && (
                              <span className="ml-2 px-2 py-1 text-xs bg-indigo-200 text-indigo-800 rounded-full">
                                {config.authentication.methods.length} selected
                              </span>
                            )}
                          </h4>
                          {config.authentication.enabled && (
                            <button
                              onClick={() => setConfig(prev => ({
                                ...prev,
                                authentication: { enabled: false, methods: [] }
                              }))}
                              className="text-xs text-indigo-600 hover:text-indigo-800 underline"
                            >
                              Clear all
                            </button>
                          )}
                        </div>
                        <div className="grid grid-cols-2 gap-2">
                          {authenticationOptions.map((auth) => {
                            const isSelected = config.authentication.methods.includes(auth.value)
                            return (
                              <button
                                key={auth.value}
                                onClick={() => toggleAuthenticationOption(auth.value)}
                                className={`px-3 py-3 text-sm border rounded-lg transition-all duration-200 text-left ${
                                  isSelected
                                    ? 'border-indigo-500 bg-indigo-100 text-indigo-900 shadow-md transform scale-105'
                                    : 'border-indigo-300 bg-white text-indigo-700 hover:bg-indigo-50 hover:border-indigo-400'
                                }`}
                              >
                                <div className="flex items-center space-x-2">
                                  <span className="text-lg">{auth.icon}</span>
                                  <div className="flex-1">
                                    <div className="font-medium flex items-center">
                                      {auth.label}
                                      {isSelected && (
                                        <CheckCircleIconSolid className="w-4 h-4 text-indigo-600 ml-1" />
                                      )}
                                    </div>
                                    <div className="text-xs text-indigo-600 mt-1">{auth.description}</div>
                                  </div>
                                </div>
                              </button>
                            )
                          })}
                        </div>
                        {config.authentication.enabled && (
                          <div className="mt-3 p-2 bg-indigo-100 rounded-md">
                            <p className="text-xs text-indigo-800">
                              🛡️ This will add authentication middleware, services, and security features
                            </p>
                          </div>
                        )}
                      </div>

                      {/* Configuration Summary */}
                      {(config.database.enabled || config.authentication.enabled) && (
                        <div className="p-4 bg-gradient-to-r from-green-50 to-blue-50 border border-green-200 rounded-lg">
                          <h4 className="text-sm font-medium text-green-900 mb-2 flex items-center">
                            📋 Advanced Configuration Summary
                          </h4>
                          <div className="space-y-2 text-sm">
                            {config.database.enabled && (
                              <div className="flex items-center space-x-2">
                                <span className="text-green-600">•</span>
                                <span className="text-green-800">
                                  Database: {config.database.drivers.join(', ')}
                                </span>
                              </div>
                            )}
                            {config.authentication.enabled && (
                              <div className="flex items-center space-x-2">
                                <span className="text-blue-600">•</span>
                                <span className="text-blue-800">
                                  Auth: {config.authentication.methods.join(', ')}
                                </span>
                              </div>
                            )}
                            <div className="mt-2 pt-2 border-t border-green-200">
                              <p className="text-xs text-gray-600">
                                ⚡ Your project structure will be automatically updated to include the selected features
                              </p>
                            </div>
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                )}

                {/* Generate Button */}
                <div className="pt-6 border-t border-gray-200">
                  <button
                    onClick={handleGenerate}
                    disabled={!isFormValid || generationStatus === 'generating'}
                    className={`w-full py-4 px-6 rounded-lg font-semibold text-white transition-all ${
                      isFormValid && generationStatus !== 'generating'
                        ? 'bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 shadow-lg hover:shadow-xl'
                        : 'bg-gray-500 cursor-not-allowed'
                    }`}
                  >
                    {generationStatus === 'generating' ? (
                      <div className="flex items-center justify-center space-x-2">
                        <ArrowPathIcon className="w-5 h-5 animate-spin" />
                        <span>Generating...</span>
                      </div>
                    ) : (
                      <div className="flex items-center justify-center space-x-2">
                        <SparklesIcon className="w-5 h-5" />
                        <span>Generate Project</span>
                      </div>
                    )}
                  </button>
                </div>
              </div>
            </div>
          </div>

          {/* Preview & Files Panel */}
          <div className="lg:col-span-8">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              
              {/* Preview Panel */}
              <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-6">
                <div className="flex items-center space-x-3 mb-6">
                  <div className="w-8 h-8 bg-green-100 rounded-lg flex items-center justify-center">
                    <EyeIcon className="w-5 h-5 text-green-600" />
                  </div>
                  <div>
                    <h2 className="text-lg font-semibold text-gray-900">Live Preview</h2>
                    <p className="text-sm text-gray-600">Real-time generation status</p>
                  </div>
                </div>

                {/* Generation Status */}
                <div className="mb-6">
                  <div className="flex items-center space-x-3 mb-4">
                    {generationStatus === 'generating' && <ArrowPathIcon className="w-5 h-5 text-blue-500 animate-spin" />}
                    {generationStatus === 'completed' && <CheckCircleIcon className="w-5 h-5 text-green-500" />}
                    {generationStatus === 'error' && <ExclamationTriangleIcon className="w-5 h-5 text-red-500" />}
                    {generationStatus === 'idle' && <PlayIcon className="w-5 h-5 text-gray-500" />}
                    
                    <div className="flex-1">
                      <p className="text-sm font-medium text-gray-900">
                        {generationStatus === 'generating' && 'Generating project...'}
                        {generationStatus === 'completed' && 'Generation completed!'}
                        {generationStatus === 'error' && 'Generation failed'}
                        {generationStatus === 'idle' && 'Ready to generate'}
                      </p>
                    </div>
                  </div>

                  {/* Progress Bar */}
                  {generationStatus === 'generating' && (
                    <div className="mb-4">
                      <div className="flex items-center justify-between text-xs text-gray-600 mb-1">
                        <span>Progress</span>
                        <span>{generationProgress}%</span>
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-gradient-to-r from-blue-500 to-indigo-500 h-2 rounded-full transition-all duration-300"
                          style={{ width: `${generationProgress}%` }}
                        />
                      </div>
                    </div>
                  )}

                  {/* Statistics */}
                  <div className="grid grid-cols-2 gap-4">
                    <div className="bg-gray-50 rounded-lg p-3">
                      <p className="text-xs font-medium text-gray-500 uppercase tracking-wide">Files</p>
                      <p className="text-xl font-bold text-gray-900">
                        {(() => {
                          const countFiles = (nodes: FileNode[]): number => {
                            return nodes.reduce((count, node) => {
                              if (node.type === 'file') {
                                return count + 1
                              } else if (node.children) {
                                return count + countFiles(node.children)
                              }
                              return count
                            }, 0)
                          }
                          const totalFiles = countFiles(fileStructure)
                          return generationStatus === 'completed' ? totalFiles : Math.floor(generationProgress / 100 * totalFiles)
                        })()}
                      </p>
                    </div>
                    <div className="bg-gray-50 rounded-lg p-3">
                      <p className="text-xs font-medium text-gray-500 uppercase tracking-wide">Size</p>
                      <p className="text-xl font-bold text-gray-900">
                        {(() => {
                          const calculateSize = (nodes: FileNode[]): number => {
                            return nodes.reduce((size, node) => {
                              if (node.type === 'file' && node.size) {
                                return size + node.size
                              } else if (node.children) {
                                return size + calculateSize(node.children)
                              }
                              return size
                            }, 0)
                          }
                          const totalSize = calculateSize(fileStructure)
                          const sizeInKB = Math.round(totalSize / 1024)
                          return generationStatus === 'completed' ? `${sizeInKB}KB` : '0KB'
                        })()}
                      </p>
                    </div>
                  </div>
                </div>

                {/* Action Buttons */}
                {generationStatus === 'completed' && (
                  <div className="space-y-3">
                    <button className="w-full bg-green-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-green-700 transition-colors flex items-center justify-center space-x-2">
                      <ArrowDownTrayIcon className="w-5 h-5" />
                      <span>Download Project</span>
                    </button>
                    <button className="w-full bg-gray-100 text-gray-700 py-3 px-4 rounded-lg font-medium hover:bg-gray-200 transition-colors flex items-center justify-center space-x-2">
                      <ClipboardDocumentIcon className="w-5 h-5" />
                      <span>Copy Project URL</span>
                    </button>
                  </div>
                )}
              </div>

              {/* File Explorer Panel */}
              <div className="bg-white/90 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-6">
                <div className="flex items-center space-x-3 mb-6">
                  <div className="w-8 h-8 bg-amber-100 rounded-lg flex items-center justify-center">
                    <DocumentIcon className="w-5 h-5 text-amber-600" />
                  </div>
                  <div>
                    <h2 className="text-lg font-semibold text-gray-900">File Explorer</h2>
                    <p className="text-sm text-gray-600">Project structure preview</p>
                  </div>
                </div>

                <div className="bg-gray-50 rounded-lg p-4 max-h-96 overflow-y-auto">
                  {renderFileTree(fileStructure)}
                </div>

                {/* File Preview */}
                {selectedFile && (
                  <div className="mt-4 p-4 bg-gray-900 rounded-lg">
                    <div className="flex items-center space-x-2 mb-2">
                      <CodeBracketIcon className="w-4 h-4 text-gray-500" />
                      <span className="text-sm font-mono text-gray-300">{selectedFile}</span>
                    </div>
                    <div className="text-xs font-mono text-green-400">
                      {selectedFile.endsWith('.go') && (
                        <div>
                          <span className="text-purple-400">package</span> <span className="text-yellow-300">main</span><br/>
                          <br/>
                          <span className="text-purple-400">import</span> (<br/>
                          &nbsp;&nbsp;<span className="text-green-300">"fmt"</span><br/>
                          )<br/>
                          <br/>
                          <span className="text-purple-400">func</span> <span className="text-blue-300">main</span>() {'{'}<br/>
                          &nbsp;&nbsp;fmt.Println(<span className="text-green-300">"Hello, World!"</span>)<br/>
                          {'}'}
                        </div>
                      )}
                      {selectedFile.endsWith('.md') && (
                        <div>
                          <span className="text-blue-300"># {config.projectName || 'My Go Project'}</span><br/>
                          <br/>
                          <span className="text-gray-300">A modern Go application generated with Go Starter.</span><br/>
                          <br/>
                          <span className="text-blue-300">## Features</span><br/>
                          <span className="text-gray-300">- Clean architecture</span><br/>
                          <span className="text-gray-300">- Modern logging</span><br/>
                          <span className="text-gray-300">- Comprehensive testing</span>
                        </div>
                      )}
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Quick Tips */}
            <div className="mt-6 bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-2xl p-6">
              <div className="flex items-start space-x-3">
                <div className="w-8 h-8 bg-blue-100 rounded-lg flex items-center justify-center flex-shrink-0">
                  <LightBulbIcon className="w-5 h-5 text-blue-600" />
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-blue-900 mb-2">Pro Tips</h3>
                  <ul className="space-y-2 text-sm text-blue-800">
                    <li>• Use descriptive project names that reflect your application's purpose</li>
                    <li>• Module URLs should match your repository structure for easy importing</li>
                    <li>• Switch to Advanced Mode to enable database and authentication features</li>
                    <li>• Multiple database drivers can be selected for maximum flexibility</li>
                    <li>• The file explorer updates in real-time as you add advanced features</li>
                    <li>• Watch the file count change as you select different options</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default SimpleApp