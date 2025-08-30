import { useState, useEffect } from 'react'
import { ProjectTemplate } from '../data/projectTemplates'

const FAVORITES_STORAGE_KEY = 'go-starter-favorites'
const RECENT_STORAGE_KEY = 'go-starter-recent'

interface FavoritesData {
  favorites: string[]
  recent: Array<{
    id: string
    timestamp: number
  }>
}

export function useFavorites() {
  const [favorites, setFavorites] = useState<string[]>([])
  const [recentlyUsed, setRecentlyUsed] = useState<Array<{ id: string; timestamp: number }>>([])

  // Load from localStorage on mount
  useEffect(() => {
    try {
      const storedFavorites = localStorage.getItem(FAVORITES_STORAGE_KEY)
      const storedRecent = localStorage.getItem(RECENT_STORAGE_KEY)
      
      if (storedFavorites) {
        setFavorites(JSON.parse(storedFavorites))
      }
      
      if (storedRecent) {
        const recent = JSON.parse(storedRecent)
        // Clean up old entries (older than 30 days)
        const thirtyDaysAgo = Date.now() - (30 * 24 * 60 * 60 * 1000)
        const cleanedRecent = recent.filter((item: any) => item.timestamp > thirtyDaysAgo)
        setRecentlyUsed(cleanedRecent)
      }
    } catch (error) {
      console.error('Failed to load favorites/recent from localStorage:', error)
    }
  }, [])

  // Save to localStorage whenever favorites change
  useEffect(() => {
    try {
      localStorage.setItem(FAVORITES_STORAGE_KEY, JSON.stringify(favorites))
    } catch (error) {
      console.error('Failed to save favorites to localStorage:', error)
    }
  }, [favorites])

  // Save to localStorage whenever recent changes
  useEffect(() => {
    try {
      localStorage.setItem(RECENT_STORAGE_KEY, JSON.stringify(recentlyUsed))
    } catch (error) {
      console.error('Failed to save recent to localStorage:', error)
    }
  }, [recentlyUsed])

  const addToFavorites = (templateId: string) => {
    setFavorites(prev => {
      if (prev.includes(templateId)) return prev
      return [...prev, templateId]
    })
  }

  const removeFromFavorites = (templateId: string) => {
    setFavorites(prev => prev.filter(id => id !== templateId))
  }

  const toggleFavorite = (templateId: string) => {
    if (favorites.includes(templateId)) {
      removeFromFavorites(templateId)
    } else {
      addToFavorites(templateId)
    }
  }

  const isFavorite = (templateId: string) => {
    return favorites.includes(templateId)
  }

  const addToRecent = (templateId: string) => {
    setRecentlyUsed(prev => {
      // Remove existing entry if present
      const filtered = prev.filter(item => item.id !== templateId)
      // Add to beginning with current timestamp
      const updated = [{ id: templateId, timestamp: Date.now() }, ...filtered]
      // Keep only last 10 items
      return updated.slice(0, 10)
    })
  }

  const getRecentlyUsedIds = () => {
    return recentlyUsed
      .sort((a, b) => b.timestamp - a.timestamp)
      .map(item => item.id)
  }

  const clearFavorites = () => {
    setFavorites([])
  }

  const clearRecent = () => {
    setRecentlyUsed([])
  }

  const getFavoriteTemplates = (allTemplates: ProjectTemplate[]) => {
    return allTemplates.filter(template => favorites.includes(template.id))
  }

  const getRecentTemplates = (allTemplates: ProjectTemplate[]) => {
    const recentIds = getRecentlyUsedIds()
    return recentIds
      .map(id => allTemplates.find(template => template.id === id))
      .filter(Boolean) as ProjectTemplate[]
  }

  const exportFavorites = () => {
    const data = {
      favorites,
      recentlyUsed,
      exportDate: new Date().toISOString(),
      version: '1.0'
    }
    
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'go-starter-preferences.json'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  const importFavorites = (fileContent: string) => {
    try {
      const data = JSON.parse(fileContent)
      
      if (data.favorites && Array.isArray(data.favorites)) {
        setFavorites(data.favorites)
      }
      
      if (data.recentlyUsed && Array.isArray(data.recentlyUsed)) {
        setRecentlyUsed(data.recentlyUsed)
      }
      
      return true
    } catch (error) {
      console.error('Failed to import favorites:', error)
      return false
    }
  }

  return {
    favorites,
    recentlyUsed: getRecentlyUsedIds(),
    addToFavorites,
    removeFromFavorites,
    toggleFavorite,
    isFavorite,
    addToRecent,
    clearFavorites,
    clearRecent,
    getFavoriteTemplates,
    getRecentTemplates,
    exportFavorites,
    importFavorites,
    stats: {
      favoritesCount: favorites.length,
      recentCount: recentlyUsed.length
    }
  }
}