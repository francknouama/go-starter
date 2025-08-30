import { useState, useEffect, useRef } from 'react'
import { 
  MagnifyingGlassIcon, 
  FunnelIcon, 
  StarIcon, 
  HeartIcon,
  ClockIcon,
  Squares2X2Icon,
  ListBulletIcon,
  ChevronLeftIcon,
  ChevronRightIcon
} from '@heroicons/react/20/solid'
import { HeartIcon as HeartOutlineIcon } from '@heroicons/react/24/outline'
import { PROJECT_TEMPLATES, TEMPLATE_CATEGORIES, type ProjectTemplate } from '../../data/projectTemplates'
import { useFavorites } from '../../hooks/useFavorites'
import TemplateCard from '../templates/TemplateCard'
import Button from '../common/Button'

type ViewMode = 'popular' | 'gallery' | 'favorites' | 'recent'
type LayoutMode = 'grid' | 'compact' | 'swipe'

interface MobileOptimizedGalleryProps {
  onSelectTemplate: (template: ProjectTemplate) => void
  onClose: () => void
  initialView?: ViewMode
}

export default function MobileOptimizedGallery({ 
  onSelectTemplate, 
  onClose, 
  initialView = 'popular' 
}: MobileOptimizedGalleryProps) {
  const [viewMode, setViewMode] = useState<ViewMode>(initialView)
  const [layoutMode, setLayoutMode] = useState<LayoutMode>('grid')
  const [selectedCategory, setSelectedCategory] = useState<string>('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [showFilters, setShowFilters] = useState(false)
  const [currentSwipeIndex, setCurrentSwipeIndex] = useState(0)
  
  const scrollRef = useRef<HTMLDivElement>(null)
  const favorites = useFavorites()

  // Detect mobile device
  const [isMobile, setIsMobile] = useState(false)
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768)
    }
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  // Filter templates based on view mode and filters
  const filteredTemplates = (() => {
    let templates = PROJECT_TEMPLATES

    switch (viewMode) {
      case 'popular':
        templates = PROJECT_TEMPLATES.filter(t => t.popularity >= 8).sort((a, b) => b.popularity - a.popularity)
        break
      case 'favorites':
        templates = favorites.getFavoriteTemplates(PROJECT_TEMPLATES)
        break
      case 'recent':
        templates = favorites.getRecentTemplates(PROJECT_TEMPLATES)
        break
      default:
        templates = PROJECT_TEMPLATES
    }

    // Apply category filter
    if (selectedCategory !== 'all') {
      templates = templates.filter(template => template.category === selectedCategory)
    }

    // Apply search filter
    if (searchQuery) {
      const searchTerm = searchQuery.toLowerCase()
      templates = templates.filter(template => 
        template.name.toLowerCase().includes(searchTerm) ||
        template.description.toLowerCase().includes(searchTerm) ||
        template.tags.some(tag => tag.toLowerCase().includes(searchTerm))
      )
    }

    return templates
  })()

  const handleTemplateSelect = (template: ProjectTemplate) => {
    favorites.addToRecent(template.id)
    onSelectTemplate(template)
    onClose()
  }

  const clearFilters = () => {
    setSelectedCategory('all')
    setSearchQuery('')
    setShowFilters(false)
  }

  // Swipe navigation for mobile
  const nextSlide = () => {
    if (currentSwipeIndex < Math.ceil(filteredTemplates.length / 2) - 1) {
      setCurrentSwipeIndex(prev => prev + 1)
    }
  }

  const prevSlide = () => {
    if (currentSwipeIndex > 0) {
      setCurrentSwipeIndex(prev => prev - 1)
    }
  }

  const viewModeOptions = [
    { 
      id: 'popular' as ViewMode, 
      name: 'Popular', 
      icon: StarIcon, 
      count: PROJECT_TEMPLATES.filter(t => t.popularity >= 8).length 
    },
    { 
      id: 'gallery' as ViewMode, 
      name: 'All', 
      icon: Squares2X2Icon, 
      count: PROJECT_TEMPLATES.length 
    },
    { 
      id: 'favorites' as ViewMode, 
      name: 'Saved', 
      icon: HeartIcon, 
      count: favorites.stats.favoritesCount 
    },
    { 
      id: 'recent' as ViewMode, 
      name: 'Recent', 
      icon: ClockIcon, 
      count: favorites.stats.recentCount 
    }
  ]

  return (
    <div className="fixed inset-0 z-50 overflow-hidden bg-white">
      {/* Mobile Header */}
      <div className="sticky top-0 z-40 bg-white border-b border-gray-200 shadow-sm">
        {/* Top Bar */}
        <div className="flex items-center justify-between px-4 py-3">
          <button
            onClick={onClose}
            className="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
          >
            <ChevronLeftIcon className="w-5 h-5" />
            <span className="font-medium">Back</span>
          </button>
          
          <h1 className="text-lg font-bold text-gray-900 text-center flex-1">
            Production Blueprints
          </h1>
          
          <div className="w-16"> {/* Spacer for centering */}</div>
        </div>

        {/* View Mode Tabs */}
        <div className="flex border-b border-gray-100 bg-gray-50">
          {viewModeOptions.map((option) => {
            const Icon = option.icon
            const isActive = viewMode === option.id
            
            return (
              <button
                key={option.id}
                onClick={() => {
                  setViewMode(option.id)
                  setCurrentSwipeIndex(0)
                }}
                className={`flex-1 flex items-center justify-center gap-1 py-3 text-sm font-medium transition-colors ${
                  isActive
                    ? 'text-emerald-600 bg-white border-b-2 border-emerald-500'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                <Icon className="w-4 h-4" />
                <span>{option.name}</span>
                {option.count > 0 && (
                  <span className={`ml-1 px-1.5 py-0.5 rounded-full text-xs ${
                    isActive
                      ? 'bg-emerald-100 text-emerald-600'
                      : 'bg-gray-200 text-gray-500'
                  }`}>
                    {option.count}
                  </span>
                )}
              </button>
            )
          })}
        </div>

        {/* Search and Filter */}
        <div className="p-4 space-y-3">
          {/* Search Bar */}
          <div className="relative">
            <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              type="text"
              placeholder="Search blueprints..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-emerald-500 focus:border-transparent bg-white text-sm"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
              >
                ✕
              </button>
            )}
          </div>

          {/* Quick Category Filter */}
          <div className="flex gap-2 overflow-x-auto pb-2">
            <button
              onClick={() => setSelectedCategory('all')}
              className={`flex-shrink-0 px-4 py-2 rounded-full text-sm font-medium transition-colors ${
                selectedCategory === 'all'
                  ? 'bg-emerald-100 text-emerald-700 border border-emerald-200'
                  : 'bg-white text-gray-600 border border-gray-200'
              }`}
            >
              All ({PROJECT_TEMPLATES.length})
            </button>
            
            {TEMPLATE_CATEGORIES.map((category) => (
              <button
                key={category.id}
                onClick={() => setSelectedCategory(category.id)}
                className={`flex-shrink-0 px-4 py-2 rounded-full text-sm font-medium transition-colors flex items-center gap-2 ${
                  selectedCategory === category.id
                    ? 'bg-emerald-100 text-emerald-700 border border-emerald-200'
                    : 'bg-white text-gray-600 border border-gray-200'
                }`}
              >
                {category.icon} {category.name} ({category.count})
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-y-auto bg-gray-50">
        {filteredTemplates.length === 0 ? (
          <div className="text-center py-16 px-4">
            <div className="text-gray-400 mb-6">
              <MagnifyingGlassIcon className="w-16 h-16 mx-auto" />
            </div>
            <h3 className="text-xl font-semibold text-gray-900 mb-3">
              {viewMode === 'favorites' ? 'No saved blueprints' : 
               viewMode === 'recent' ? 'No recent blueprints' : 'No blueprints found'}
            </h3>
            <p className="text-gray-600 mb-6 max-w-sm mx-auto">
              {viewMode === 'favorites' ? 'Save blueprints by tapping the heart icon to see them here.' :
               viewMode === 'recent' ? 'Blueprints you use will appear here for quick access.' :
               'Try adjusting your search terms or category selection.'}
            </p>
            {(searchQuery || selectedCategory !== 'all') && (
              <Button variant="primary" onClick={clearFilters}>
                Clear filters
              </Button>
            )}
          </div>
        ) : (
          <>
            {/* Swipe Mode for Mobile */}
            {isMobile && layoutMode === 'swipe' ? (
              <div className="relative">
                <div 
                  className="flex transition-transform duration-300 ease-out"
                  style={{ transform: `translateX(-${currentSwipeIndex * 100}%)` }}
                >
                  {Array.from({ length: Math.ceil(filteredTemplates.length / 2) }).map((_, slideIndex) => (
                    <div key={slideIndex} className="w-full flex-shrink-0 p-4 space-y-4">
                      {filteredTemplates.slice(slideIndex * 2, slideIndex * 2 + 2).map((template) => (
                        <div key={template.id} className="relative">
                          <TemplateCard
                            template={template}
                            onSelect={() => handleTemplateSelect(template)}
                            onPreview={() => {/* Handle preview */}}
                            displayMode="standard"
                            layoutMode="grid"
                          />
                          
                          {/* Floating Favorite Button */}
                          <button
                            onClick={() => favorites.toggleFavorite(template.id)}
                            className="absolute top-3 right-3 p-2 rounded-full bg-white shadow-lg border border-gray-200 hover:bg-gray-50 transition-colors z-10"
                          >
                            {favorites.isFavorite(template.id) ? (
                              <HeartIcon className="w-5 h-5 text-red-500" />
                            ) : (
                              <HeartOutlineIcon className="w-5 h-5 text-gray-400" />
                            )}
                          </button>
                        </div>
                      ))}
                    </div>
                  ))}
                </div>
                
                {/* Swipe Navigation */}
                {Math.ceil(filteredTemplates.length / 2) > 1 && (
                  <div className="flex items-center justify-between p-4">
                    <button
                      onClick={prevSlide}
                      disabled={currentSwipeIndex === 0}
                      className="p-3 rounded-full bg-white shadow-lg border border-gray-200 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <ChevronLeftIcon className="w-5 h-5 text-gray-600" />
                    </button>
                    
                    <div className="flex gap-2">
                      {Array.from({ length: Math.ceil(filteredTemplates.length / 2) }).map((_, index) => (
                        <button
                          key={index}
                          onClick={() => setCurrentSwipeIndex(index)}
                          className={`w-2 h-2 rounded-full transition-colors ${
                            index === currentSwipeIndex ? 'bg-emerald-500' : 'bg-gray-300'
                          }`}
                        />
                      ))}
                    </div>
                    
                    <button
                      onClick={nextSlide}
                      disabled={currentSwipeIndex === Math.ceil(filteredTemplates.length / 2) - 1}
                      className="p-3 rounded-full bg-white shadow-lg border border-gray-200 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <ChevronRightIcon className="w-5 h-5 text-gray-600" />
                    </button>
                  </div>
                )}
              </div>
            ) : (
              /* Grid Mode */
              <div className="p-4">
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                  {filteredTemplates.map((template) => (
                    <div key={template.id} className="relative">
                      <TemplateCard
                        template={template}
                        onSelect={() => handleTemplateSelect(template)}
                        onPreview={() => {/* Handle preview */}}
                        displayMode="standard"
                        layoutMode="grid"
                      />
                      
                      {/* Floating Favorite Button */}
                      <button
                        onClick={() => favorites.toggleFavorite(template.id)}
                        className="absolute top-3 right-3 p-2 rounded-full bg-white shadow-lg border border-gray-200 hover:bg-gray-50 transition-colors z-10"
                      >
                        {favorites.isFavorite(template.id) ? (
                          <HeartIcon className="w-4 h-4 text-red-500" />
                        ) : (
                          <HeartOutlineIcon className="w-4 h-4 text-gray-400" />
                        )}
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {/* Layout Mode Toggle (Mobile) */}
      {isMobile && filteredTemplates.length > 0 && (
        <div className="sticky bottom-4 left-1/2 transform -translate-x-1/2 z-40">
          <div className="bg-white rounded-full shadow-lg border border-gray-200 p-1 flex">
            <button
              onClick={() => setLayoutMode('grid')}
              className={`p-3 rounded-full transition-colors ${
                layoutMode === 'grid'
                  ? 'bg-emerald-100 text-emerald-600'
                  : 'text-gray-600 hover:bg-gray-100'
              }`}
            >
              <Squares2X2Icon className="w-5 h-5" />
            </button>
            <button
              onClick={() => setLayoutMode('swipe')}
              className={`p-3 rounded-full transition-colors ${
                layoutMode === 'swipe'
                  ? 'bg-emerald-100 text-emerald-600'
                  : 'text-gray-600 hover:bg-gray-100'
              }`}
            >
              <ChevronRightIcon className="w-5 h-5" />
            </button>
          </div>
        </div>
      )}
    </div>
  )
}