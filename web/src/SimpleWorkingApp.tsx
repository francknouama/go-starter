import { useState } from 'react'
import './index.css'

interface ProjectConfig {
  projectName: string
  moduleUrl: string
  projectType: string
  framework: string
  architecture: string
  logger: string
  goVersion: string
}

export default function SimpleWorkingApp() {
  const [config, setConfig] = useState<ProjectConfig>({
    projectName: '',
    moduleUrl: '',
    projectType: 'web-api',
    framework: 'gin',
    architecture: 'standard',
    logger: 'slog',
    goVersion: '1.21'
  })

  const [isGenerating, setIsGenerating] = useState(false)
  const [result, setResult] = useState<string | null>(null)

  const updateConfig = (field: keyof ProjectConfig, value: string) => {
    setConfig(prev => ({ ...prev, [field]: value }))
  }

  const generateProject = async () => {
    if (!config.projectName || !config.moduleUrl) {
      alert('Please fill in all required fields')
      return
    }

    setIsGenerating(true)
    setResult(null)

    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000))
      setResult(`✅ Successfully generated project "${config.projectName}"!`)
    } catch (error) {
      setResult(`❌ Error: ${error}`)
    } finally {
      setIsGenerating(false)
    }
  }

  const isFormValid = config.projectName.length > 0 && config.moduleUrl.length > 0

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 shadow-sm">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center space-x-3">
              <div className="flex items-center justify-center w-10 h-10 bg-blue-600 rounded-lg">
                <span className="text-white font-bold text-lg">🚀</span>
              </div>
              <div>
                <h1 className="text-xl font-bold text-gray-900">Go Starter</h1>
                <p className="text-sm text-gray-500">Web Project Generator</p>
              </div>
            </div>
            <div className="text-sm text-gray-600">
              <span className="inline-flex items-center px-2 py-1 rounded-full bg-green-100 text-green-800">
                ● Online
              </span>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Configuration Panel */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-6">Project Configuration</h2>
              
              <div className="space-y-6">
                {/* Project Name */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Project Name *
                  </label>
                  <input
                    type="text"
                    value={config.projectName}
                    onChange={(e) => updateConfig('projectName', e.target.value)}
                    placeholder="my-awesome-project"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                {/* Module URL */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Module Path *
                  </label>
                  <input
                    type="text"
                    value={config.moduleUrl}
                    onChange={(e) => updateConfig('moduleUrl', e.target.value)}
                    placeholder="github.com/user/my-awesome-project"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                </div>

                {/* Project Type */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Project Type
                  </label>
                  <select
                    value={config.projectType}
                    onChange={(e) => updateConfig('projectType', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="web-api">Web API</option>
                    <option value="cli">CLI Application</option>
                    <option value="library">Library</option>
                    <option value="lambda">AWS Lambda</option>
                  </select>
                </div>

                {/* Framework */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Framework
                  </label>
                  <select
                    value={config.framework}
                    onChange={(e) => updateConfig('framework', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="gin">Gin</option>
                    <option value="echo">Echo</option>
                    <option value="fiber">Fiber</option>
                    <option value="chi">Chi</option>
                  </select>
                </div>

                {/* Architecture */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Architecture
                  </label>
                  <select
                    value={config.architecture}
                    onChange={(e) => updateConfig('architecture', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="standard">Standard</option>
                    <option value="clean">Clean Architecture</option>
                    <option value="ddd">Domain-Driven Design</option>
                    <option value="hexagonal">Hexagonal</option>
                  </select>
                </div>

                {/* Logger */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Logger
                  </label>
                  <select
                    value={config.logger}
                    onChange={(e) => updateConfig('logger', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="slog">slog (Standard)</option>
                    <option value="zap">Zap</option>
                    <option value="logrus">Logrus</option>
                    <option value="zerolog">Zerolog</option>
                  </select>
                </div>

                {/* Generate Button */}
                <button
                  onClick={generateProject}
                  disabled={!isFormValid || isGenerating}
                  className={`w-full py-3 px-4 rounded-lg font-medium transition-all duration-200 ${
                    isFormValid && !isGenerating
                      ? 'bg-blue-600 hover:bg-blue-700 text-white shadow-lg hover:shadow-xl transform hover:scale-105'
                      : 'bg-gray-300 text-gray-500 cursor-not-allowed'
                  }`}
                >
                  {isGenerating ? '🔄 Generating...' : '🚀 Generate Project'}
                </button>

                {/* Result */}
                {result && (
                  <div className={`p-4 rounded-lg ${
                    result.includes('✅') ? 'bg-green-50 text-green-800 border border-green-200' : 'bg-red-50 text-red-800 border border-red-200'
                  }`}>
                    {result}
                  </div>
                )}
              </div>
            </div>

            {/* Preview Panel */}
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-6">Project Preview</h2>
              
              <div className="space-y-4">
                {/* Project Info */}
                <div className="bg-gray-50 rounded-lg p-4">
                  <h3 className="font-medium text-gray-900 mb-2">Project Information</h3>
                  <div className="space-y-1 text-sm text-gray-600">
                    <div><strong>Name:</strong> {config.projectName || 'Not set'}</div>
                    <div><strong>Module:</strong> {config.moduleUrl || 'Not set'}</div>
                    <div><strong>Type:</strong> {config.projectType}</div>
                    <div><strong>Framework:</strong> {config.framework}</div>
                    <div><strong>Architecture:</strong> {config.architecture}</div>
                    <div><strong>Logger:</strong> {config.logger}</div>
                  </div>
                </div>

                {/* File Structure Preview */}
                <div className="bg-gray-50 rounded-lg p-4">
                  <h3 className="font-medium text-gray-900 mb-2">Expected File Structure</h3>
                  <div className="text-sm text-gray-600 font-mono">
                    <div>📁 {config.projectName || 'project-name'}/</div>
                    <div className="ml-4">📄 main.go</div>
                    <div className="ml-4">📄 go.mod</div>
                    <div className="ml-4">📄 go.sum</div>
                    <div className="ml-4">📄 README.md</div>
                    <div className="ml-4">📁 cmd/</div>
                    <div className="ml-8">📄 server.go</div>
                    <div className="ml-4">📁 internal/</div>
                    <div className="ml-8">📄 config.go</div>
                    <div className="ml-8">📄 handler.go</div>
                    <div className="ml-8">📄 router.go</div>
                    <div className="ml-4">📁 tests/</div>
                    <div className="ml-8">📄 main_test.go</div>
                  </div>
                </div>

                {/* Quick Start */}
                <div className="bg-blue-50 rounded-lg p-4">
                  <h3 className="font-medium text-blue-900 mb-2">Quick Start Commands</h3>
                  <div className="text-sm text-blue-800 font-mono space-y-1">
                    <div>$ go mod tidy</div>
                    <div>$ go run main.go</div>
                    <div>$ go test ./...</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}