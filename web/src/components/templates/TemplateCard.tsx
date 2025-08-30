import { useState } from 'react'
import { 
  EyeIcon, 
  StarIcon, 
  ClockIcon, 
  TagIcon,
  ChevronRightIcon,
  SparklesIcon
} from '@heroicons/react/20/solid'
import { COMPLEXITY_LEVELS, type ProjectTemplate } from '../../data/projectTemplates'
import Button from '../common/Button'

interface TemplateCardProps {
  template: ProjectTemplate
  onSelect: () => void
  onPreview: () => void
  displayMode?: 'compact' | 'standard' | 'detailed'
  layoutMode?: 'grid' | 'compact'
}

export default function TemplateCard({ 
  template, 
  onSelect, 
  onPreview, 
  displayMode = 'standard',
  layoutMode = 'grid'
}: TemplateCardProps) {
  const [isHovered, setIsHovered] = useState(false)
  
  const complexityConfig = COMPLEXITY_LEVELS.find(c => c.id === template.complexity)
  const complexityColor = complexityConfig?.color || '#6B7280'

  const getTechStackDisplay = (techStack: string[], maxVisible = 3) => {
    const visible = techStack.slice(0, maxVisible)
    const remaining = techStack.length - maxVisible
    return { visible, remaining }
  }

  const { visible: visibleTech, remaining: remainingTech } = getTechStackDisplay(template.techStack)

  // Determine card styling based on layout and display mode
  const getCardClasses = () => {
    const baseClasses = "relative group cursor-pointer transition-all duration-300"
    
    if (layoutMode === 'compact') {
      return `${baseClasses} bg-white rounded-lg border border-gray-200 hover:border-purple-300 hover:shadow-md overflow-hidden`
    }
    
    // Enhanced grid layout with better desktop spacing and production-ready styling
    return `${baseClasses} transform ${
      isHovered ? 'scale-[1.02] shadow-2xl shadow-emerald-100/20' : 'hover:scale-[1.01] hover:shadow-xl hover:shadow-blue-100/20'
    } bg-white/98 backdrop-blur-xl rounded-2xl border border-gray-200/60 overflow-hidden shadow-lg hover:border-emerald-300`
  }

  // Conditional rendering based on display mode
  const shouldShowSection = (section: string) => {
    if (displayMode === 'compact') {
      return ['icon', 'title', 'category', 'primaryAction'].includes(section)
    }
    if (displayMode === 'standard') {
      return !['setupTime', 'useCase'].includes(section) || layoutMode === 'grid'
    }
    return true // detailed mode shows everything
  }

  // Compact layout for dense viewing
  if (layoutMode === 'compact') {
    return (
      <div className={getCardClasses()}>
        <div className="p-4 flex items-center gap-4">
          {/* Icon and basic info */}
          <div 
            className="text-2xl flex items-center justify-center w-10 h-10 rounded-lg border border-gray-200 bg-gray-50"
            style={{ color: template.color }}
          >
            {template.icon}
          </div>
          
          {/* Content */}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <h3 className="font-semibold text-gray-900 text-base truncate">
                {template.name}
              </h3>
              {template.popularity >= 8 && (
                <StarIcon className="w-4 h-4 text-amber-500 flex-shrink-0" />
              )}
            </div>
            <p className="text-sm text-gray-600 truncate">
              {template.description}
            </p>
            <div className="flex items-center gap-3 mt-2">
              <span className="text-xs text-gray-500 capitalize">
                {template.category.replace('-', ' ')}
              </span>
              <span className="text-xs text-gray-500">•</span>
              <span className="text-xs text-gray-500 capitalize">
                {template.complexity}
              </span>
            </div>
          </div>
          
          {/* Actions */}
          <div className="flex gap-2 flex-shrink-0">
            <Button variant="ghost" size="sm" onClick={onPreview}>
              <EyeIcon className="w-4 h-4" />
            </Button>
            <Button variant="primary" size="sm" onClick={onSelect}>
              Use
            </Button>
          </div>
        </div>
      </div>
    )
  }

  // Standard grid card layout (enhanced for desktop)
  return (
    <div 
      className={getCardClasses()}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Enhanced Card Header with production-ready badge */}
      <div 
        className="h-24 relative overflow-hidden border-b border-gray-100"
        style={{ 
          background: `linear-gradient(135deg, ${template.color}12, ${template.color}06)` 
        }}
      >
        {/* Production Ready Badge */}
        <div className="absolute top-3 right-3 z-10">
          <div className="bg-emerald-100 text-emerald-700 px-2 py-1 rounded-md text-xs font-semibold border border-emerald-200 flex items-center gap-1">
            <div className="w-1.5 h-1.5 bg-emerald-500 rounded-full animate-pulse"></div>
            Production Ready
          </div>
        </div>
        <div className="absolute inset-0 flex items-center justify-between p-4">
          <div className="flex items-center gap-3">
            <div 
              className="text-3xl flex items-center justify-center w-12 h-12 rounded-xl bg-white shadow-md border border-gray-100 group-hover:shadow-lg transition-shadow duration-200"
              style={{ color: template.color }}
            >
              {template.icon}
            </div>
            <div>
              <div className="text-xs font-semibold text-gray-600 uppercase tracking-wider mb-1">
                {template.category.replace('-', ' ')}
              </div>
              <div className="flex items-center gap-2">
                <div 
                  className="w-2 h-2 rounded-full shadow-sm"
                  style={{ backgroundColor: complexityColor }}
                />
                <span className="text-xs font-medium text-gray-600 capitalize">
                  {template.complexity}
                </span>
                <span className="text-xs text-gray-400">•</span>
                <span className="text-xs text-gray-500">
                  {template.estimatedSetupTime}
                </span>
              </div>
            </div>
          </div>
          
          {/* Popularity and file count indicator */}
          <div className="flex flex-col items-end gap-1">
            {template.popularity >= 8 && (
              <div className="flex items-center gap-1 text-amber-600 bg-amber-50 px-2 py-1 rounded-lg border border-amber-200 shadow-sm">
                <StarIcon className="w-3 h-3" />
                <span className="text-xs font-semibold">Popular</span>
              </div>
            )}
            <div className="text-xs text-gray-500 font-medium">
              {template.quickStart.learnMore.match(/\d+/)?.[0] || '~25'} files
            </div>
          </div>
        </div>
      </div>

      {/* Enhanced Card Content with better hierarchy */}
      <div className="p-5 space-y-4 flex-1 flex flex-col">
        {/* Title and Description */}
        <div>
          <h3 className="font-bold text-gray-900 text-lg mb-2 leading-tight group-hover:text-emerald-700 transition-colors duration-200">
            {template.name}
          </h3>
          <p className="text-sm text-gray-600 line-clamp-3 leading-relaxed">
            {template.description}
          </p>
        </div>

        {/* Use Case with enhanced styling */}
        {shouldShowSection('useCase') && (
          <div 
            className="rounded-lg p-3 border border-opacity-20 relative overflow-hidden"
            style={{
              backgroundColor: `${template.color}08`,
              borderColor: template.color
            }}
          >
            <div className="flex items-start gap-2">
              <SparklesIcon 
                className="w-4 h-4 mt-0.5 flex-shrink-0" 
                style={{ color: template.color }}
              />
              <div>
                <p className="text-xs font-semibold mb-1" style={{ color: template.color }}>
                  Perfect for:
                </p>
                <p className="text-xs text-gray-700 leading-relaxed">
                  {template.useCase}
                </p>
              </div>
            </div>
          </div>
        )}

        {/* Tech Stack with enhanced styling */}
        {shouldShowSection('techStack') && (
          <div className="space-y-2">
            <div className="flex items-center gap-1">
              <TagIcon className="w-3 h-3 text-gray-500" />
              <span className="text-xs font-semibold text-gray-700">Tech Stack</span>
            </div>
            <div className="flex flex-wrap gap-1.5">
              {visibleTech.slice(0, 3).map((tech, index) => (
                <span 
                  key={index}
                  className="inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-medium border transition-colors duration-200"
                  style={{
                    backgroundColor: `${template.color}10`,
                    color: template.color,
                    borderColor: `${template.color}30`
                  }}
                >
                  {tech}
                </span>
              ))}
              {template.techStack.length > 3 && (
                <span className="inline-flex items-center px-2.5 py-1 rounded-lg bg-gray-100 text-xs text-gray-500 border border-gray-200">
                  +{template.techStack.length - 3} more
                </span>
              )}
            </div>
          </div>
        )}

        {/* Setup Time - Only show in detailed mode */}
        {shouldShowSection('setupTime') && (
          <div className="flex items-center gap-2 text-xs text-gray-500 pt-1 border-t border-gray-100">
            <ClockIcon className="w-3 h-3" />
            <span>Setup: {template.estimatedSetupTime}</span>
          </div>
        )}
      </div>

      {/* Enhanced Card Actions with improved styling */}
      <div className="px-5 pb-5 space-y-2 mt-auto">
        <Button
          variant="primary"
          size="sm"
          fullWidth
          onClick={onSelect}
          className="group relative overflow-hidden font-semibold bg-gradient-to-r from-emerald-500 to-blue-500 hover:from-emerald-600 hover:to-blue-600 text-white border-0 shadow-lg hover:shadow-xl transition-all duration-200"
        >
          <span className="relative z-10 flex items-center justify-center gap-2">
            Use Blueprint
            <ChevronRightIcon className="w-4 h-4 transition-transform group-hover:translate-x-1" />
          </span>
          <div className="absolute inset-0 bg-gradient-to-r from-white/0 via-white/20 to-white/0 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-500"></div>
        </Button>
        
        <Button
          variant="ghost"
          size="sm"
          fullWidth
          onClick={onPreview}
          className="flex items-center justify-center gap-2 text-gray-600 hover:text-emerald-600 hover:bg-emerald-50 border border-gray-200 hover:border-emerald-200 transition-all duration-200"
        >
          <EyeIcon className="w-4 h-4" />
          Preview Details
        </Button>
      </div>

      {/* Enhanced Hover Effects */}
      {isHovered && (
        <>
          <div 
            className="absolute inset-0 rounded-2xl opacity-10 transition-opacity duration-300 pointer-events-none"
            style={{
              background: `radial-gradient(circle at center, ${template.color}20 0%, transparent 70%)`
            }}
          />
          <div 
            className="absolute inset-0 rounded-2xl opacity-5 transition-opacity duration-300 pointer-events-none"
            style={{
              boxShadow: `0 0 30px ${template.color}40`
            }}
          />
        </>
      )}
    </div>
  )
}