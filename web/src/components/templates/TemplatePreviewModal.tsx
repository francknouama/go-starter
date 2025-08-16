import { useState } from 'react'
import { 
  XMarkIcon, 
  ChevronRightIcon,
  ClockIcon,
  StarIcon,
  TagIcon,
  PlayIcon,
  DocumentTextIcon,
  CubeIcon,
  LightBulbIcon,
  CommandLineIcon
} from '@heroicons/react/20/solid'
import { COMPLEXITY_LEVELS, type ProjectTemplate } from '../../data/projectTemplates'
import Button from '../common/Button'

interface TemplatePreviewModalProps {
  template: ProjectTemplate
  onClose: () => void
  onSelect: () => void
}

export default function TemplatePreviewModal({ template, onClose, onSelect }: TemplatePreviewModalProps) {
  const [activeTab, setActiveTab] = useState<'overview' | 'architecture' | 'quickstart'>('overview')
  
  const complexityConfig = COMPLEXITY_LEVELS.find(c => c.id === template.complexity)

  const tabs = [
    { id: 'overview', name: 'Overview', icon: DocumentTextIcon },
    { id: 'architecture', name: 'Architecture', icon: CubeIcon },
    { id: 'quickstart', name: 'Quick Start', icon: PlayIcon }
  ]

  return (
    <div className="fixed inset-0 z-60 overflow-hidden">
      <div className="min-h-screen px-4 text-center">
        {/* Background overlay */}
        <div className="fixed inset-0 bg-black bg-opacity-50 transition-opacity" onClick={onClose} />
        
        {/* Modal panel */}
        <div className="inline-block w-full max-w-4xl my-8 text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl relative">
          
          {/* Modal Header */}
          <div 
            className="px-6 py-6 border-b border-gray-200 bg-gradient-to-r rounded-t-2xl relative overflow-hidden"
            style={{ 
              background: `linear-gradient(135deg, ${template.color}15, ${template.color}05)` 
            }}
          >
            {/* Background Pattern */}
            <div className="absolute inset-0 opacity-5">
              <div className="absolute inset-0" style={{
                backgroundImage: `radial-gradient(circle at 30% 70%, ${template.color} 15%, transparent 16%)`,
                backgroundSize: '60px 60px'
              }} />
            </div>
            
            <div className="relative flex items-start justify-between">
              <div className="flex items-start gap-4">
                {/* Template Icon */}
                <div 
                  className="text-4xl flex items-center justify-center w-16 h-16 rounded-xl border-2 border-white bg-white shadow-lg"
                  style={{ color: template.color }}
                >
                  {template.icon}
                </div>
                
                {/* Template Info */}
                <div className="space-y-2">
                  <div className="flex items-center gap-3">
                    <h2 className="text-2xl font-bold text-gray-900">{template.name}</h2>
                    {template.popularity >= 8 && (
                      <div className="flex items-center gap-1 text-amber-600 bg-amber-100 px-3 py-1 rounded-full">
                        <StarIcon className="w-4 h-4" />
                        <span className="text-sm font-medium">Popular</span>
                      </div>
                    )}
                  </div>
                  
                  <p className="text-gray-700 leading-relaxed max-w-lg">
                    {template.description}
                  </p>
                  
                  <div className="flex items-center gap-6 text-sm text-gray-600">
                    <div className="flex items-center gap-2">
                      <div 
                        className="w-3 h-3 rounded-full"
                        style={{ backgroundColor: complexityConfig?.color }}
                      />
                      <span className="capitalize font-medium">{template.complexity}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <ClockIcon className="w-4 h-4" />
                      <span>{template.estimatedSetupTime}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <TagIcon className="w-4 h-4" />
                      <span className="capitalize">{template.category}</span>
                    </div>
                  </div>
                </div>
              </div>
              
              {/* Close Button */}
              <Button variant="ghost" onClick={onClose} className="p-2">
                <XMarkIcon className="w-5 h-5" />
              </Button>
            </div>
          </div>

          {/* Tab Navigation */}
          <div className="border-b border-gray-200 bg-gray-50">
            <nav className="flex space-x-8 px-6" aria-label="Tabs">
              {tabs.map((tab) => {
                const Icon = tab.icon
                const isActive = activeTab === tab.id
                return (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id as typeof activeTab)}
                    className={`
                      group inline-flex items-center gap-2 py-4 px-1 border-b-2 font-medium text-sm
                      transition-colors duration-200
                      ${isActive
                        ? 'border-purple-500 text-purple-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                      }
                    `}
                  >
                    <Icon className="w-5 h-5" />
                    {tab.name}
                  </button>
                )
              })}
            </nav>
          </div>

          {/* Tab Content */}
          <div className="px-6 py-6 max-h-[60vh] overflow-y-auto">
            
            {/* Overview Tab */}
            {activeTab === 'overview' && (
              <div className="space-y-6">
                {/* Use Case */}
                <div>
                  <div className="flex items-center gap-2 mb-3">
                    <LightBulbIcon className="w-5 h-5 text-purple-500" />
                    <h3 className="font-semibold text-gray-900">Perfect For</h3>
                  </div>
                  <div className="bg-purple-50 rounded-lg p-4">
                    <p className="text-gray-700">{template.useCase}</p>
                  </div>
                </div>

                {/* Tech Stack */}
                <div>
                  <h3 className="font-semibold text-gray-900 mb-3">Technology Stack</h3>
                  <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                    {template.techStack.map((tech, index) => (
                      <div key={index} className="flex items-center gap-2 p-3 bg-gray-50 rounded-lg">
                        <div className="w-2 h-2 rounded-full bg-green-500" />
                        <span className="text-sm font-medium text-gray-700">{tech}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Tags */}
                <div>
                  <h3 className="font-semibold text-gray-900 mb-3">Features & Tags</h3>
                  <div className="flex flex-wrap gap-2">
                    {template.tags.map((tag, index) => (
                      <span 
                        key={index}
                        className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>

                {/* Recommended For */}
                <div>
                  <h3 className="font-semibold text-gray-900 mb-3">Recommended For</h3>
                  <div className="space-y-2">
                    {template.recommendedFor.map((item, index) => (
                      <div key={index} className="flex items-center gap-2">
                        <ChevronRightIcon className="w-4 h-4 text-green-500" />
                        <span className="text-gray-700">{item}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            )}

            {/* Architecture Tab */}
            {activeTab === 'architecture' && (
              <div className="space-y-6">
                {/* Architecture Diagram */}
                <div>
                  <h3 className="font-semibold text-gray-900 mb-3">System Architecture</h3>
                  <div className="bg-gray-900 rounded-lg p-6 overflow-x-auto">
                    <pre className="text-green-400 text-sm font-mono whitespace-pre leading-relaxed">
                      {template.architecture.diagram}
                    </pre>
                  </div>
                </div>

                {/* Key Components */}
                <div>
                  <h3 className="font-semibold text-gray-900 mb-3">Key Components</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    {template.architecture.components.map((component, index) => (
                      <div key={index} className="flex items-center gap-3 p-3 bg-blue-50 rounded-lg">
                        <div className="w-2 h-2 rounded-full bg-blue-500" />
                        <span className="text-sm font-medium text-gray-700">{component}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Design Patterns */}
                <div>
                  <h3 className="font-semibold text-gray-900 mb-3">Design Patterns</h3>
                  <div className="space-y-2">
                    {template.architecture.patterns.map((pattern, index) => (
                      <div key={index} className="flex items-center gap-2 p-2">
                        <div className="w-1.5 h-1.5 rounded-full bg-purple-500" />
                        <span className="text-gray-700 font-medium">{pattern}</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            )}

            {/* Quick Start Tab */}
            {activeTab === 'quickstart' && (
              <div className="space-y-6">
                {/* Setup Commands */}
                <div>
                  <div className="flex items-center gap-2 mb-3">
                    <CommandLineIcon className="w-5 h-5 text-green-500" />
                    <h3 className="font-semibold text-gray-900">Setup Commands</h3>
                  </div>
                  <div className="bg-gray-900 rounded-lg p-4 space-y-2">
                    {template.quickStart.commands.map((command, index) => (
                      <div key={index} className="flex items-center gap-3">
                        <span className="text-gray-500 text-sm font-mono">$</span>
                        <code className="text-green-400 text-sm font-mono">{command}</code>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Next Steps */}
                <div>
                  <h3 className="font-semibold text-gray-900 mb-3">Next Steps</h3>
                  <div className="space-y-3">
                    {template.quickStart.nextSteps.map((step, index) => (
                      <div key={index} className="flex items-start gap-3 p-3 bg-gray-50 rounded-lg">
                        <div className="flex items-center justify-center w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-sm font-bold mt-0.5">
                          {index + 1}
                        </div>
                        <span className="text-gray-700 text-sm leading-relaxed">{step}</span>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Learn More */}
                <div className="bg-purple-50 rounded-lg p-4">
                  <h4 className="font-medium text-purple-900 mb-2">Learn More</h4>
                  <p className="text-purple-700 text-sm">{template.quickStart.learnMore}</p>
                </div>
              </div>
            )}
          </div>

          {/* Modal Footer */}
          <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 rounded-b-2xl">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-4 text-sm text-gray-600">
                <span>Estimated setup time: <strong>{template.estimatedSetupTime}</strong></span>
                <span>•</span>
                <span>Complexity: <strong className="capitalize">{template.complexity}</strong></span>
              </div>
              <div className="flex items-center gap-3">
                <Button variant="outline" onClick={onClose}>
                  Cancel
                </Button>
                <Button 
                  variant="primary" 
                  onClick={onSelect}
                  className="group relative overflow-hidden"
                >
                  <span className="relative z-10 flex items-center gap-2">
                    Use This Template
                    <ChevronRightIcon className="w-4 h-4 transition-transform group-hover:translate-x-0.5" />
                  </span>
                  <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full group-hover:translate-x-full transition-transform duration-700" />
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}