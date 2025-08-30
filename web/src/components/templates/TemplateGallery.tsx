import { useState, useMemo, useCallback } from 'react'
import { MagnifyingGlassIcon, FunnelIcon, StarIcon, ClockIcon, Squares2X2Icon, ListBulletIcon } from '@heroicons/react/20/solid'
import { 
  PROJECT_TEMPLATES, 
  TEMPLATE_CATEGORIES, 
  COMPLEXITY_LEVELS,
  getTemplatesByCategory, 
  getTemplatesByComplexity,
  searchTemplates,
  getPopularTemplates,
  type ProjectTemplate 
} from '../../data/projectTemplates'
import TemplateCard from './TemplateCard'
import TemplatePreviewModal from './TemplatePreviewModal'
import Button from '../common/Button'

type ViewMode = 'popular' | 'gallery'
type LayoutMode = 'grid' | 'compact'
type CardDisplayMode = 'compact' | 'standard' | 'detailed'

interface TemplateGalleryProps {
  onSelectTemplate: (template: ProjectTemplate) => void
  onClose: () => void
}

export default function TemplateGallery({ onSelectTemplate, onClose }: TemplateGalleryProps) {
  const [selectedCategory, setSelectedCategory] = useState<string>('all')
  const [selectedComplexity, setSelectedComplexity] = useState<string>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [showFilters, setShowFilters] = useState(false)
  const [previewTemplate, setPreviewTemplate] = useState<ProjectTemplate | null>(null)
  const [viewMode, setViewMode] = useState<ViewMode>('popular')
  const [layoutMode, setLayoutMode] = useState<LayoutMode>('grid')
  const [cardDisplayMode, setCardDisplayMode] = useState<CardDisplayMode>('standard')

  // Filter templates based on selected criteria
  const filteredTemplates = useMemo(() => {
    let templates = viewMode === 'popular' ? getPopularTemplates(12) : PROJECT_TEMPLATES

    // Apply category filter
    if (selectedCategory !== 'all') {
      templates = templates.filter(template => template.category === selectedCategory)
    }

    // Apply complexity filter
    if (selectedComplexity !== 'all') {
      templates = templates.filter(template => template.complexity === selectedComplexity)
    }

    // Apply search filter
    if (searchQuery) {
      const searchResults = searchTemplates(searchQuery)
      templates = templates.filter(template => 
        searchResults.some(result => result.id === template.id)
      )
    }

    return templates
  }, [selectedCategory, selectedComplexity, searchQuery, viewMode])

  const handleTemplateSelect = (template: ProjectTemplate) => {
    onSelectTemplate(template)
    onClose()
  }

  const clearFilters = useCallback(() => {
    setSelectedCategory('all')
    setSelectedComplexity('all')
    setSearchQuery('')
  }, [])

  // Responsive grid classes based on layout mode and screen size
  const getGridClasses = () => {
    if (layoutMode === 'compact') {
      return 'grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-4'
    }
    // Enhanced desktop-focused grid (3 columns max for better readability)
    return 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-3 2xl:grid-cols-4 gap-6'
  }

  const hasActiveFilters = selectedCategory !== 'all' || selectedComplexity !== 'all' || searchQuery

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="min-h-screen px-4 text-center">
        {/* Background overlay */}
        <div className="fixed inset-0 bg-black bg-opacity-50 transition-opacity" onClick={onClose} />
        
        {/* Enhanced Modal panel with better desktop sizing */}
        <div className="inline-block w-full max-w-7xl my-4 md:my-8 text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl relative"
             style={{ minHeight: 'min(90vh, 900px)' }}>
          
          {/* Enhanced Header with Achievement Banner */}
          <div className="px-6 py-6 border-b border-gray-200 bg-gradient-to-r from-emerald-50 via-blue-50 to-purple-50 rounded-t-2xl">
            <div className="flex items-center justify-between">
              <div>
                <div className="flex items-center gap-3 mb-2">
                  <h2 className="text-3xl font-bold text-gray-900">Production-Ready Blueprints</h2>
                  <div className="bg-emerald-100 text-emerald-800 px-3 py-1 rounded-full text-sm font-semibold border border-emerald-200">
                    🎉 100% Coverage
                  </div>
                </div>
                <p className="text-lg text-gray-700 mb-2">Historic achievement: <strong>12 production-ready blueprints</strong> - choose the perfect foundation for your Go project</p>
                <div className="flex items-center gap-4 text-sm text-gray-600">
                  <span className="flex items-center gap-1">
                    <span className="w-2 h-2 bg-emerald-500 rounded-full"></span>
                    All blueprints production-ready
                  </span>
                  <span className="flex items-center gap-1">
                    <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                    Comprehensive testing included
                  </span>
                  <span className="flex items-center gap-1">
                    <span className="w-2 h-2 bg-purple-500 rounded-full"></span>
                    Real-world patterns
                  </span>
                </div>
              </div>
              <Button variant="ghost" onClick={onClose} className="p-2">
                <span className="sr-only">Close</span>
                ✕
              </Button>
            </div>
            
            {/* Enhanced Stats and Controls */}
            <div className="flex flex-col lg:flex-row gap-6 mt-6">
              {/* Achievement Stats */}
              <div className="flex items-center gap-6">
                <div className="text-center">
                  <div className="text-2xl font-bold text-emerald-600">{PROJECT_TEMPLATES.length}</div>
                  <div className="text-xs text-gray-600">Blueprints</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-blue-600">{TEMPLATE_CATEGORIES.length}</div>
                  <div className="text-xs text-gray-600">Categories</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-purple-600">100%</div>
                  <div className="text-xs text-gray-600">Production</div>
                </div>
              </div>
              
              {/* View mode and layout controls */}
              <div className="flex flex-col sm:flex-row gap-4 flex-1 justify-end">
              {/* View mode toggle */}
              <div className="flex gap-2">
                <button
                  onClick={() => setViewMode('popular')}
                  className={`px-5 py-2.5 rounded-lg text-sm font-medium transition-all duration-200 ${
                    viewMode === 'popular'
                      ? 'bg-gradient-to-r from-emerald-100 to-blue-100 text-emerald-700 border border-emerald-200 shadow-sm'
                      : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50 hover:border-gray-300'
                  }`}
                >
                  <StarIcon className="w-4 h-4 inline-block mr-2" />
                  Most Popular
                </button>
                <button
                  onClick={() => setViewMode('gallery')}
                  className={`px-5 py-2.5 rounded-lg text-sm font-medium transition-all duration-200 ${
                    viewMode === 'gallery'
                      ? 'bg-gradient-to-r from-emerald-100 to-blue-100 text-emerald-700 border border-emerald-200 shadow-sm'
                      : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50 hover:border-gray-300'
                  }`}
                >
                  All {PROJECT_TEMPLATES.length} Blueprints
                </button>
              </div>

                {/* Layout mode toggle */}
                <div className="flex gap-1 bg-gray-100 p-1 rounded-lg">
                  <button
                    onClick={() => setLayoutMode('grid')}
                    className={`p-2.5 rounded-md transition-all duration-200 ${
                      layoutMode === 'grid'
                        ? 'bg-white text-emerald-600 shadow-sm'
                        : 'text-gray-600 hover:text-gray-900'
                    }`}
                    title="Grid view"
                  >
                    <Squares2X2Icon className="w-4 h-4" />
                  </button>
                  <button
                    onClick={() => setLayoutMode('compact')}
                    className={`p-2.5 rounded-md transition-all duration-200 ${
                      layoutMode === 'compact'
                        ? 'bg-white text-emerald-600 shadow-sm'
                        : 'text-gray-600 hover:text-gray-900'
                    }`}
                    title="Compact view"
                  >
                    <ListBulletIcon className="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          {/* Enhanced Search and Category Navigation */}
          <div className="px-6 py-5 space-y-5 border-b border-gray-100 bg-gradient-to-r from-gray-50 to-white">
            {/* Search and main controls */}
            <div className="flex gap-4 items-center">
              {/* Enhanced search with instant feedback */}
              <div className="flex-1 relative">
                <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-500" />
                <input
                  type="text"
                  placeholder="Search templates by name, category, or technology..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white"
                />
                {searchQuery && (
                  <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                    <button
                      onClick={() => setSearchQuery('')}
                      className="text-gray-500 hover:text-gray-600 text-lg"
                    >
                      ✕
                    </button>
                  </div>
                )}
              </div>

              {/* Advanced filter toggle */}
              <Button
                variant="outline"
                onClick={() => setShowFilters(!showFilters)}
                className="flex items-center gap-2 whitespace-nowrap"
              >
                <FunnelIcon className="w-4 h-4" />
                {showFilters ? 'Hide Filters' : 'More Filters'}
                {hasActiveFilters && (
                  <span className="w-2 h-2 bg-purple-500 rounded-full" />
                )}
              </Button>

              {/* Clear all filters */}
              {hasActiveFilters && (
                <Button variant="ghost" onClick={clearFilters} className="text-sm whitespace-nowrap">
                  Clear All
                </Button>
              )}
            </div>

            {/* Enhanced Category Navigation */}
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <span className="text-sm font-semibold text-gray-800">Browse by Category:</span>
                <span className="text-xs text-gray-500">All categories are production-ready ✅</span>
              </div>
              
              <div className="flex flex-wrap gap-3 items-center">
              
              {/* Popular categories as always-visible pills */}
                <button
                  onClick={() => setSelectedCategory('all')}
                  className={`px-4 py-2 rounded-xl text-sm font-medium transition-all duration-200 flex items-center gap-2 ${
                    selectedCategory === 'all'
                      ? 'bg-gradient-to-r from-emerald-100 to-blue-100 text-emerald-700 border border-emerald-200 shadow-sm'
                      : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50 hover:border-gray-300'
                  }`}
                >
                  🎯 All ({PROJECT_TEMPLATES.length})
                </button>
              
                {TEMPLATE_CATEGORIES.map((category) => (
                  <button
                    key={category.id}
                    onClick={() => setSelectedCategory(category.id)}
                    className={`px-4 py-2 rounded-xl text-sm font-medium transition-all duration-200 flex items-center gap-2 ${
                      selectedCategory === category.id
                        ? 'bg-gradient-to-r from-emerald-100 to-blue-100 text-emerald-700 border border-emerald-200 shadow-sm'
                        : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50 hover:border-gray-300'
                    }`}
                  >
                    {category.icon} {category.name} ({category.count})
                  </button>
                ))}
              </div>
              
              
              {/* Complexity Level Filter */}
              <div className="space-y-2 pt-2 border-t border-gray-200">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-semibold text-gray-800">Filter by Complexity:</span>
                  <span className="text-xs text-gray-500">Choose your experience level</span>
                </div>
                
                <div className="flex flex-wrap gap-2">
                  <button
                    onClick={() => setSelectedComplexity('all')}
                    className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-200 ${
                      selectedComplexity === 'all'
                        ? 'bg-gradient-to-r from-gray-100 to-gray-200 text-gray-700 border border-gray-300 shadow-sm'
                        : 'bg-white text-gray-500 border border-gray-200 hover:bg-gray-50'
                    }`}
                  >
                    All Levels
                  </button>
                  {COMPLEXITY_LEVELS.filter(c => c.count > 0).map((complexity) => (
                    <button
                      key={complexity.id}
                      onClick={() => setSelectedComplexity(complexity.id)}
                      className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-200 flex items-center gap-1 ${
                        selectedComplexity === complexity.id
                          ? 'text-white border border-opacity-20 shadow-sm'
                          : 'bg-white border border-gray-200 hover:bg-gray-50'
                      }`}
                      style={{
                        backgroundColor: selectedComplexity === complexity.id ? complexity.color : undefined,
                        color: selectedComplexity === complexity.id ? 'white' : complexity.color
                      }}
                    >
                      <div 
                        className="w-2 h-2 rounded-full"
                        style={{ backgroundColor: complexity.color }}
                      />
                      {complexity.name} ({complexity.count})
                    </button>
                  ))
                }
                </div>
              </div>
            </div>
          </div>

          {/* Advanced Filter panels */}
          {showFilters && (
            <div className="px-6 pb-4">
              <div className="p-4 bg-white rounded-lg border space-y-4">
                {/* Category filter */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Category</label>
                  <div className="flex flex-wrap gap-2">
                    <button
                      onClick={() => setSelectedCategory('all')}
                      className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                        selectedCategory === 'all'
                          ? 'bg-purple-100 text-purple-700 border border-purple-200'
                          : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                      }`}
                    >
                      All
                    </button>
                    {TEMPLATE_CATEGORIES.map((category) => (
                      <button
                        key={category.id}
                        onClick={() => setSelectedCategory(category.id)}
                        className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                          selectedCategory === category.id
                            ? 'bg-purple-100 text-purple-700 border border-purple-200'
                            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                        }`}
                      >
                        {category.icon} {category.name}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Complexity filter */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">Complexity</label>
                  <div className="flex flex-wrap gap-2">
                    <button
                      onClick={() => setSelectedComplexity('all')}
                      className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                        selectedComplexity === 'all'
                          ? 'bg-purple-100 text-purple-700 border border-purple-200'
                          : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                      }`}
                    >
                      All Levels
                    </button>
                    {COMPLEXITY_LEVELS.map((complexity) => (
                      <button
                        key={complexity.id}
                        onClick={() => setSelectedComplexity(complexity.id)}
                        className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                          selectedComplexity === complexity.id
                            ? 'bg-purple-100 text-purple-700 border border-purple-200'
                            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                        }`}
                      >
                        {complexity.name}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          )}

          {/* Enhanced Templates Grid with better desktop experience */}
          <div className="px-6 py-6 flex-1 overflow-y-auto" style={{ maxHeight: 'calc(90vh - 320px)' }}>
            {filteredTemplates.length === 0 ? (
              <div className="text-center py-16">
                <div className="text-gray-400 mb-6">
                  <MagnifyingGlassIcon className="w-16 h-16 mx-auto" />
                </div>
                <h3 className="text-xl font-semibold text-gray-900 mb-3">No blueprints found</h3>
                <p className="text-gray-600 mb-6 max-w-md mx-auto leading-relaxed">
                  We couldn't find any blueprints matching your current filters. Try adjusting your search terms or category selection.
                </p>
                <div className="flex gap-3 justify-center">
                  <Button variant="outline" onClick={clearFilters}>
                    Clear all filters
                  </Button>
                  <Button variant="primary" onClick={() => setViewMode('gallery')}>
                    Browse all blueprints
                  </Button>
                </div>
              </div>
            ) : (
              <div className={getGridClasses()}>
                {filteredTemplates.map((template) => (
                  <TemplateCard
                    key={template.id}
                    template={template}
                    onSelect={() => handleTemplateSelect(template)}
                    onPreview={() => setPreviewTemplate(template)}
                    displayMode={cardDisplayMode}
                    layoutMode={layoutMode}
                  />
                ))}
              </div>
            )}

            {/* Enhanced Results Summary */}
            {filteredTemplates.length > 0 && (
              <div className="mt-8 pt-6 border-t border-gray-100">
                <div className="flex items-center justify-between text-sm">
                  <div className="text-gray-600">
                    Showing <strong className="text-gray-900">{filteredTemplates.length}</strong> of <strong className="text-gray-900">{PROJECT_TEMPLATES.length}</strong> production-ready blueprints
                  </div>
                  <div className="flex items-center gap-4 text-xs text-gray-500">
                    <span className="flex items-center gap-1">
                      <div className="w-2 h-2 bg-emerald-500 rounded-full"></div>
                      Production Ready
                    </span>
                    <span className="flex items-center gap-1">
                      <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                      Comprehensive Tests
                    </span>
                    <span className="flex items-center gap-1">
                      <div className="w-2 h-2 bg-purple-500 rounded-full"></div>
                      Real-world Patterns
                    </span>
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Template Preview Modal */}
          {previewTemplate && (
            <TemplatePreviewModal
              template={previewTemplate}
              onClose={() => setPreviewTemplate(null)}
              onSelect={() => handleTemplateSelect(previewTemplate)}
            />
          )}
        </div>
      </div>
    </div>
  )
}