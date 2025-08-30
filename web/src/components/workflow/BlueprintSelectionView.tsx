/**
 * Blueprint Selection View
 * Displays available blueprints and allows user selection
 */

import React, { useState, useMemo } from 'react'
import { motion } from 'framer-motion'
import { 
  MagnifyingGlassIcon,
  FunnelIcon,
  SparklesIcon,
  CheckCircleIcon
} from '@heroicons/react/24/outline'
import type { Blueprint } from '../../services/api'
import { useBlueprintSelection } from '../../stores/generationStore'
import Button from '../common/Button'
import { LoadingStates } from '../common/LoadingStates'
import TemplateGallery from '../templates/TemplateGallery'
import AchievementBanner from '../blueprints/AchievementBanner'

interface BlueprintSelectionViewProps {
  blueprints: Blueprint[]
  loading: boolean
  onSelect: (blueprint: Blueprint) => void
  selectedBlueprint?: Blueprint | null
  className?: string
}

export default function BlueprintSelectionView({
  blueprints,
  loading,
  onSelect,
  selectedBlueprint,
  className = ''
}: BlueprintSelectionViewProps) {
  const { blueprintFilters, setBlueprintFilters } = useBlueprintSelection()
  const [showFilters, setShowFilters] = useState(false)

  // Filter blueprints based on search and filters
  const filteredBlueprints = useMemo(() => {
    let filtered = [...blueprints]

    // Search filter
    if (blueprintFilters.search) {
      const searchTerm = blueprintFilters.search.toLowerCase()
      filtered = filtered.filter(bp =>
        bp.name.toLowerCase().includes(searchTerm) ||
        bp.description.toLowerCase().includes(searchTerm) ||
        bp.features.some(feature => feature.toLowerCase().includes(searchTerm)) ||
        bp.tags.some(tag => tag.toLowerCase().includes(searchTerm))
      )
    }

    // Category filter
    if (blueprintFilters.category) {
      filtered = filtered.filter(bp => bp.category === blueprintFilters.category)
    }

    // Difficulty filter
    if (blueprintFilters.difficulty) {
      filtered = filtered.filter(bp => bp.difficulty === blueprintFilters.difficulty)
    }

    return filtered
  }, [blueprints, blueprintFilters])

  // Get unique categories and difficulties for filters
  const availableCategories = useMemo(() => 
    [...new Set(blueprints.map(bp => bp.category))].sort()
  , [blueprints])

  const availableDifficulties = useMemo(() => 
    [...new Set(blueprints.map(bp => bp.difficulty))].sort()
  , [blueprints])

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setBlueprintFilters({ search: e.target.value })
  }

  const handleCategoryFilter = (category: string) => {
    setBlueprintFilters({ 
      category: blueprintFilters.category === category ? undefined : category 
    })
  }

  const handleDifficultyFilter = (difficulty: string) => {
    setBlueprintFilters({ 
      difficulty: blueprintFilters.difficulty === difficulty ? undefined : difficulty 
    })
  }

  const clearFilters = () => {
    setBlueprintFilters({ search: '', category: undefined, difficulty: undefined })
  }

  if (loading) {
    return (
      <div className={`flex items-center justify-center h-full ${className}`}>
        <LoadingStates 
          message="Loading Blueprints"
          description="Discovering available project templates..."
        />
      </div>
    )
  }

  return (
    <div className={`h-full flex flex-col ${className}`}>
      {/* Header */}
      <div className="px-6 py-8 bg-gradient-to-r from-blue-50 via-purple-50 to-pink-50">
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          className="max-w-6xl mx-auto"
        >
          <div className="flex items-center gap-3 mb-4">
            <SparklesIcon className="w-8 h-8 text-blue-600" />
            <h1 className="text-3xl font-bold text-gray-900">
              Choose Your Blueprint
            </h1>
          </div>
          <p className="text-lg text-gray-600 max-w-3xl">
            Select from our collection of production-ready Go project templates. 
            Each blueprint includes best practices, proper architecture, and everything you need to get started.
          </p>

          {/* Achievement Banner */}
          <div className="mt-6">
            <AchievementBanner 
              totalBlueprints={blueprints.length}
              productionReady={blueprints.filter(bp => bp.tags?.includes('production-ready')).length}
            />
          </div>
        </motion.div>
      </div>

      {/* Search and Filters */}
      <div className="px-6 py-4 bg-white border-b border-gray-200">
        <div className="max-w-6xl mx-auto">
          <div className="flex flex-col md:flex-row gap-4 items-start md:items-center justify-between">
            {/* Search */}
            <div className="relative flex-1 max-w-md">
              <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
              <input
                type="text"
                placeholder="Search blueprints..."
                value={blueprintFilters.search || ''}
                onChange={handleSearchChange}
                className="pl-10 pr-4 py-2 w-full border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>

            {/* Filter Toggle and Results Count */}
            <div className="flex items-center gap-4">
              <span className="text-sm text-gray-600">
                {filteredBlueprints.length} of {blueprints.length} blueprints
              </span>
              
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowFilters(!showFilters)}
                className="flex items-center gap-2"
              >
                <FunnelIcon className="w-4 h-4" />
                Filters
              </Button>
            </div>
          </div>

          {/* Filter Controls */}
          {showFilters && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              className="mt-4 p-4 bg-gray-50 rounded-lg"
            >
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Category Filter */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Category
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {availableCategories.map((category) => (
                      <button
                        key={category}
                        onClick={() => handleCategoryFilter(category)}
                        className={`px-3 py-1 text-sm rounded-full border transition-colors ${
                          blueprintFilters.category === category
                            ? 'bg-blue-100 border-blue-300 text-blue-800'
                            : 'bg-white border-gray-300 text-gray-700 hover:bg-gray-50'
                        }`}
                      >
                        {category}
                      </button>
                    ))}
                  </div>
                </div>

                {/* Difficulty Filter */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Difficulty
                  </label>
                  <div className="flex flex-wrap gap-2">
                    {availableDifficulties.map((difficulty) => (
                      <button
                        key={difficulty}
                        onClick={() => handleDifficultyFilter(difficulty)}
                        className={`px-3 py-1 text-sm rounded-full border transition-colors ${
                          blueprintFilters.difficulty === difficulty
                            ? 'bg-purple-100 border-purple-300 text-purple-800'
                            : 'bg-white border-gray-300 text-gray-700 hover:bg-gray-50'
                        }`}
                      >
                        {difficulty}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              {/* Clear Filters */}
              <div className="mt-4 pt-4 border-t border-gray-200">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={clearFilters}
                  className="text-gray-600"
                >
                  Clear all filters
                </Button>
              </div>
            </motion.div>
          )}
        </div>
      </div>

      {/* Blueprint Gallery */}
      <div className="flex-1 overflow-y-auto px-6 py-6">
        <div className="max-w-6xl mx-auto">
          {filteredBlueprints.length === 0 ? (
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              className="text-center py-12"
            >
              <div className="w-24 h-24 mx-auto mb-6 rounded-full bg-gray-100 flex items-center justify-center">
                <MagnifyingGlassIcon className="w-12 h-12 text-gray-400" />
              </div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">
                No blueprints found
              </h3>
              <p className="text-gray-600 mb-6">
                Try adjusting your search terms or filters to find the perfect blueprint for your project.
              </p>
              <Button variant="outline" onClick={clearFilters}>
                Clear filters
              </Button>
            </motion.div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredBlueprints.map((blueprint, index) => (
                <motion.div
                  key={blueprint.id}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.1 }}
                  className={`group relative bg-white rounded-xl border-2 transition-all duration-200 hover:shadow-lg cursor-pointer ${
                    selectedBlueprint?.id === blueprint.id
                      ? 'border-blue-500 ring-4 ring-blue-100'
                      : 'border-gray-200 hover:border-blue-300'
                  }`}
                  onClick={() => onSelect(blueprint)}
                >
                  {/* Selected Indicator */}
                  {selectedBlueprint?.id === blueprint.id && (
                    <div className="absolute -top-2 -right-2 w-6 h-6 bg-blue-500 rounded-full flex items-center justify-center">
                      <CheckCircleIcon className="w-4 h-4 text-white" />
                    </div>
                  )}

                  {/* Content */}
                  <div className="p-6">
                    {/* Header */}
                    <div className="flex items-start justify-between mb-3">
                      <div>
                        <h3 className="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">
                          {blueprint.name}
                        </h3>
                        <div className="flex items-center gap-2 mt-1">
                          <span className={`px-2 py-1 text-xs rounded-full ${
                            blueprint.difficulty === 'beginner' 
                              ? 'bg-green-100 text-green-800'
                              : blueprint.difficulty === 'intermediate'
                              ? 'bg-yellow-100 text-yellow-800'
                              : 'bg-red-100 text-red-800'
                          }`}>
                            {blueprint.difficulty}
                          </span>
                          <span className="px-2 py-1 text-xs bg-gray-100 text-gray-800 rounded-full">
                            {blueprint.category}
                          </span>
                        </div>
                      </div>
                      
                      <div className="text-right">
                        <div className="text-sm font-medium text-gray-900">
                          ~{blueprint.estimatedFiles} files
                        </div>
                      </div>
                    </div>

                    {/* Description */}
                    <p className="text-gray-600 text-sm mb-4 line-clamp-2">
                      {blueprint.description}
                    </p>

                    {/* Features */}
                    <div className="mb-4">
                      <div className="flex flex-wrap gap-1">
                        {blueprint.features.slice(0, 3).map((feature) => (
                          <span
                            key={feature}
                            className="px-2 py-1 text-xs bg-blue-50 text-blue-700 rounded"
                          >
                            {feature}
                          </span>
                        ))}
                        {blueprint.features.length > 3 && (
                          <span className="px-2 py-1 text-xs bg-gray-50 text-gray-600 rounded">
                            +{blueprint.features.length - 3} more
                          </span>
                        )}
                      </div>
                    </div>

                    {/* Technologies */}
                    <div className="flex items-center justify-between text-sm text-gray-500">
                      <div className="flex items-center gap-2">
                        {blueprint.frameworks.slice(0, 2).map((framework) => (
                          <span key={framework} className="font-mono text-xs bg-gray-100 px-2 py-1 rounded">
                            {framework}
                          </span>
                        ))}
                      </div>
                      {blueprint.tags?.includes('production-ready') && (
                        <span className="text-green-600 font-medium text-xs">
                          ✓ Production Ready
                        </span>
                      )}
                    </div>
                  </div>

                  {/* Hover Effect */}
                  <div className="absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-500 opacity-0 group-hover:opacity-5 transition-opacity rounded-xl" />
                </motion.div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Selection Actions */}
      {selectedBlueprint && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="px-6 py-4 bg-white border-t border-gray-200"
        >
          <div className="max-w-6xl mx-auto flex items-center justify-between">
            <div className="flex items-center gap-3">
              <CheckCircleIcon className="w-6 h-6 text-green-600" />
              <div>
                <div className="font-semibold text-gray-900">
                  Selected: {selectedBlueprint.name}
                </div>
                <div className="text-sm text-gray-600">
                  {selectedBlueprint.description}
                </div>
              </div>
            </div>
            
            <Button
              variant="primary"
              onClick={() => onSelect(selectedBlueprint)}
              className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
            >
              Configure Project
            </Button>
          </div>
        </motion.div>
      )}
    </div>
  )
}