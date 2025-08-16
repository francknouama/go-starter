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
        <div className="inline-block w-full max-w-6xl my-4 md:my-8 text-left align-middle transition-all transform bg-white shadow-xl rounded-2xl relative"
             style={{ minHeight: 'min(85vh, 800px)' }}>
          
          {/* Header */}
          <div className="px-6 py-4 border-b border-gray-200 bg-gradient-to-r from-purple-50 to-indigo-50 rounded-t-2xl">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-bold text-gray-900">Project Templates</h2>
                <p className="text-sm text-gray-600 mt-1">Choose a real-world template to jumpstart your project</p>
              </div>
              <Button variant="ghost" onClick={onClose} className="p-2">
                <span className="sr-only">Close</span>
                ✕
              </Button>
            </div>
            
            {/* View mode and layout controls */}
            <div className="flex flex-col sm:flex-row gap-4 mt-4">
              {/* View mode toggle */}
              <div className="flex gap-2">
                <button
                  onClick={() => setViewMode('popular')}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                    viewMode === 'popular'
                      ? 'bg-purple-100 text-purple-700 border border-purple-200'
                      : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
                  }`}
                >
                  <StarIcon className="w-4 h-4 inline-block mr-1" />
                  Popular Templates
                </button>
                <button
                  onClick={() => setViewMode('gallery')}
                  className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
                    viewMode === 'gallery'
                      ? 'bg-purple-100 text-purple-700 border border-purple-200'
                      : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-50'
                  }`}
                >
                  All Templates ({PROJECT_TEMPLATES.length})
                </button>
              </div>

              {/* Layout mode toggle */}
              <div className="flex gap-1 ml-auto">
                <button
                  onClick={() => setLayoutMode('grid')}
                  className={`p-2 rounded-lg transition-colors ${
                    layoutMode === 'grid'
                      ? 'bg-purple-100 text-purple-700'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`}
                  title="Grid view"
                >
                  <Squares2X2Icon className="w-4 h-4" />
                </button>
                <button
                  onClick={() => setLayoutMode('compact')}
                  className={`p-2 rounded-lg transition-colors ${
                    layoutMode === 'compact'
                      ? 'bg-purple-100 text-purple-700'
                      : 'text-gray-600 hover:bg-gray-100'
                  }`}
                  title="Compact view"
                >
                  <ListBulletIcon className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>

          {/* Enhanced Search and Always-Visible Filters */}
          <div className="px-6 py-4 space-y-4 border-b border-gray-100 bg-gray-50">
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

            {/* Always-visible quick filters */}
            <div className="flex flex-wrap gap-2 items-center">
              <span className="text-sm font-medium text-gray-700 mr-2">Quick filters:</span>
              
              {/* Popular categories as always-visible pills */}
              <button
                onClick={() => setSelectedCategory('all')}
                className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                  selectedCategory === 'all'
                    ? 'bg-purple-100 text-purple-700 border border-purple-200'
                    : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-100'
                }`}
              >
                All
              </button>
              
              {TEMPLATE_CATEGORIES.slice(0, 4).map((category) => (
                <button
                  key={category.id}
                  onClick={() => setSelectedCategory(category.id)}
                  className={`px-3 py-1 rounded-full text-sm font-medium transition-colors ${
                    selectedCategory === category.id
                      ? 'bg-purple-100 text-purple-700 border border-purple-200'
                      : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-100'
                  }`}
                >
                  {category.icon} {category.name}
                </button>
              ))}
              
              {/* Complexity quick filter */}
              <div className="border-l border-gray-300 pl-3 ml-2">
                {COMPLEXITY_LEVELS.map((complexity) => (
                  <button
                    key={complexity.id}
                    onClick={() => setSelectedComplexity(complexity.id)}
                    className={`px-3 py-1 ml-1 rounded-full text-sm font-medium transition-colors ${
                      selectedComplexity === complexity.id
                        ? 'bg-purple-100 text-purple-700 border border-purple-200'
                        : 'bg-white text-gray-600 border border-gray-200 hover:bg-gray-100'
                    }`}
                  >
                    {complexity.name}
                  </button>
                ))}
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
          <div className="px-6 py-6 flex-1 overflow-y-auto" style={{ maxHeight: 'calc(85vh - 280px)' }}>
            {filteredTemplates.length === 0 ? (
              <div className="text-center py-12">
                <div className="text-gray-500 mb-4">
                  <MagnifyingGlassIcon className="w-12 h-12 mx-auto" />
                </div>
                <h3 className="text-lg font-medium text-gray-900 mb-2">No templates found</h3>
                <p className="text-gray-600 mb-4">
                  Try adjusting your search or filter criteria
                </p>
                <Button variant="outline" onClick={clearFilters}>
                  Clear all filters
                </Button>
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

            {/* Results summary */}
            {filteredTemplates.length > 0 && (
              <div className="mt-6 text-center text-sm text-gray-600">
                Showing {filteredTemplates.length} of {PROJECT_TEMPLATES.length} templates
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