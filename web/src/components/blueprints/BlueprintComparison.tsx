import { useState } from 'react'
import { 
  XMarkIcon, 
  CheckIcon, 
  ClockIcon, 
  TagIcon,
  SparklesIcon,
  ChevronRightIcon
} from '@heroicons/react/20/solid'
import { ProjectTemplate } from '../../data/projectTemplates'
import Button from '../common/Button'

interface BlueprintComparisonProps {
  templates: ProjectTemplate[]
  onClose: () => void
  onSelect: (template: ProjectTemplate) => void
  className?: string
}

export default function BlueprintComparison({ 
  templates, 
  onClose, 
  onSelect, 
  className = '' 
}: BlueprintComparisonProps) {
  const [selectedTemplate, setSelectedTemplate] = useState<ProjectTemplate | null>(null)

  if (templates.length === 0) {
    return null
  }

  const comparisonRows = [
    {
      label: 'Type',
      key: 'category',
      getValue: (template: ProjectTemplate) => template.category.replace('-', ' ')
    },
    {
      label: 'Complexity',
      key: 'complexity',
      getValue: (template: ProjectTemplate) => template.complexity
    },
    {
      label: 'Setup Time',
      key: 'setupTime',
      getValue: (template: ProjectTemplate) => template.estimatedSetupTime
    },
    {
      label: 'File Count',
      key: 'fileCount',
      getValue: (template: ProjectTemplate) => {
        const match = template.quickStart.learnMore.match(/(\d+)\s+files?/)
        return match ? `~${match[1]} files` : 'Variable'
      }
    },
    {
      label: 'Main Framework',
      key: 'framework',
      getValue: (template: ProjectTemplate) => template.techStack[0] || 'Standard'
    },
    {
      label: 'Architecture',
      key: 'architecture',
      getValue: (template: ProjectTemplate) => template.config.architecture || 'Standard'
    },
    {
      label: 'Logger',
      key: 'logger',
      getValue: (template: ProjectTemplate) => template.config.logger || 'slog'
    }
  ]

  const getComplexityColor = (complexity: string) => {
    switch (complexity) {
      case 'beginner': return '#10B981'
      case 'intermediate': return '#3B82F6'
      case 'advanced': return '#F59E0B'
      case 'expert': return '#EF4444'
      default: return '#6B7280'
    }
  }

  return (
    <div className={`fixed inset-0 z-50 overflow-y-auto ${className}`}>
      <div className="min-h-screen px-4 text-center">
        {/* Background overlay */}
        <div className="fixed inset-0 bg-black bg-opacity-60 transition-opacity" onClick={onClose} />
        
        {/* Modal panel */}
        <div className="inline-block w-full max-w-6xl my-8 text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl relative">
          
          {/* Header */}
          <div className="px-6 py-4 border-b border-gray-200 bg-gradient-to-r from-blue-50 to-purple-50 rounded-t-2xl">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
                  <SparklesIcon className="w-6 h-6 text-purple-500" />
                  Blueprint Comparison
                </h2>
                <p className="text-sm text-gray-600 mt-1">
                  Compare {templates.length} blueprints side by side to find the perfect fit
                </p>
              </div>
              <Button variant="ghost" onClick={onClose} className="p-2">
                <XMarkIcon className="w-5 h-5" />
              </Button>
            </div>
          </div>

          {/* Comparison Table */}
          <div className="p-6 overflow-x-auto">
            <div className="min-w-full">
              {/* Template Headers */}
              <div className="grid gap-4 mb-6" style={{ gridTemplateColumns: '150px ' + 'repeat(' + templates.length + ', minmax(250px, 1fr))' }}>
                <div></div>
                {templates.map((template, index) => (
                  <div key={template.id} className="space-y-3">
                    {/* Template Card Header */}
                    <div 
                      className="rounded-xl p-4 border-2 transition-all duration-200 cursor-pointer hover:shadow-lg"
                      style={{
                        backgroundColor: `${template.color}08`,
                        borderColor: selectedTemplate?.id === template.id ? template.color : `${template.color}30`
                      }}
                      onClick={() => setSelectedTemplate(template)}
                    >
                      <div className="flex items-center gap-3 mb-2">
                        <div 
                          className="text-2xl flex items-center justify-center w-10 h-10 rounded-lg bg-white shadow-sm border"
                          style={{ color: template.color }}
                        >
                          {template.icon}
                        </div>
                        <div className="flex-1">
                          <h3 className="font-semibold text-gray-900 text-sm">
                            {template.name}
                          </h3>
                          <div className="flex items-center gap-2 mt-1">
                            <div 
                              className="w-2 h-2 rounded-full"
                              style={{ backgroundColor: getComplexityColor(template.complexity) }}
                            />
                            <span className="text-xs text-gray-500 capitalize">
                              {template.complexity}
                            </span>
                            <span className="text-xs text-gray-400">•</span>
                            <span className="text-xs text-gray-500 flex items-center gap-1">
                              <ClockIcon className="w-3 h-3" />
                              {template.estimatedSetupTime}
                            </span>
                          </div>
                        </div>
                      </div>
                      
                      <p className="text-xs text-gray-600 line-clamp-2 mb-3">
                        {template.description}
                      </p>

                      <Button
                        variant="primary"
                        size="sm"
                        fullWidth
                        onClick={(e) => {
                          e.stopPropagation()
                          onSelect(template)
                        }}
                        className="group"
                        style={{
                          backgroundColor: template.color,
                          borderColor: template.color
                        }}
                      >
                        <span className="flex items-center justify-center gap-1 text-xs font-medium">
                          Use Blueprint
                          <ChevronRightIcon className="w-3 h-3 transition-transform group-hover:translate-x-0.5" />
                        </span>
                      </Button>
                    </div>

                    {/* Production Ready Badge */}
                    <div className="text-center">
                      <div className="inline-flex items-center gap-1 bg-emerald-100 text-emerald-700 px-2 py-1 rounded-full text-xs font-medium border border-emerald-200">
                        <CheckIcon className="w-3 h-3" />
                        Production Ready
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* Comparison Rows */}
              <div className="space-y-2">
                {comparisonRows.map((row) => (
                  <div 
                    key={row.key}
                    className="grid gap-4 py-3 border-b border-gray-100 hover:bg-gray-50/50 transition-colors duration-200"
                    style={{ gridTemplateColumns: '150px ' + 'repeat(' + templates.length + ', minmax(250px, 1fr))' }}
                  >
                    <div className="font-medium text-gray-700 text-sm py-2">
                      {row.label}
                    </div>
                    {templates.map((template) => (
                      <div key={template.id} className="text-sm text-gray-600 py-2">
                        <span className="capitalize">
                          {row.getValue(template)}
                        </span>
                      </div>
                    ))}
                  </div>
                ))}

                {/* Use Cases Row */}
                <div 
                  className="grid gap-4 py-3 border-b border-gray-100"
                  style={{ gridTemplateColumns: '150px ' + 'repeat(' + templates.length + ', minmax(250px, 1fr))' }}
                >
                  <div className="font-medium text-gray-700 text-sm py-2">
                    Best For
                  </div>
                  {templates.map((template) => (
                    <div key={template.id} className="text-sm text-gray-600 py-2">
                      <p className="line-clamp-3">
                        {template.useCase}
                      </p>
                    </div>
                  ))}
                </div>

                {/* Tech Stack Row */}
                <div 
                  className="grid gap-4 py-3 border-b border-gray-100"
                  style={{ gridTemplateColumns: '150px ' + 'repeat(' + templates.length + ', minmax(250px, 1fr))' }}
                >
                  <div className="font-medium text-gray-700 text-sm py-2 flex items-center gap-1">
                    <TagIcon className="w-4 h-4" />
                    Tech Stack
                  </div>
                  {templates.map((template) => (
                    <div key={template.id} className="text-sm py-2">
                      <div className="flex flex-wrap gap-1">
                        {template.techStack.slice(0, 3).map((tech, idx) => (
                          <span 
                            key={idx}
                            className="inline-block px-2 py-1 rounded-md text-xs font-medium border"
                            style={{
                              backgroundColor: `${template.color}08`,
                              color: template.color,
                              borderColor: `${template.color}20`
                            }}
                          >
                            {tech}
                          </span>
                        ))}
                        {template.techStack.length > 3 && (
                          <span className="inline-block px-2 py-1 rounded-md bg-gray-100 text-gray-500 text-xs border border-gray-200">
                            +{template.techStack.length - 3}
                          </span>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>

          {/* Footer */}
          <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 rounded-b-2xl">
            <div className="flex items-center justify-between">
              <p className="text-sm text-gray-600">
                All blueprints are production-ready with comprehensive testing and documentation
              </p>
              <div className="flex gap-2">
                <Button variant="outline" onClick={onClose}>
                  Close Comparison
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}