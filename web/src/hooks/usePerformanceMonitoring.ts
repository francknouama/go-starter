import { useEffect, useRef, useState } from 'react'

interface PerformanceMetrics {
  renderTime: number
  interactionTime: number
  searchResponseTime: number
  templatesLoaded: number
  totalTemplates: number
  memoryUsage: number | null
  fps: number | null
}

export function usePerformanceMonitoring() {
  const [metrics, setMetrics] = useState<PerformanceMetrics>({
    renderTime: 0,
    interactionTime: 0,
    searchResponseTime: 0,
    templatesLoaded: 0,
    totalTemplates: 0,
    memoryUsage: null,
    fps: null
  })

  const renderStartTime = useRef<number>(Date.now())
  const lastInteractionTime = useRef<number>(0)
  const searchStartTime = useRef<number>(0)
  const fpsCounter = useRef<number>(0)
  const fpsStartTime = useRef<number>(Date.now())

  // Track render performance
  useEffect(() => {
    const renderTime = Date.now() - renderStartTime.current
    setMetrics(prev => ({ ...prev, renderTime }))
  }, [])

  // Track FPS
  useEffect(() => {
    let animationId: number

    const trackFPS = () => {
      fpsCounter.current++
      const now = Date.now()
      const elapsed = now - fpsStartTime.current

      if (elapsed >= 1000) {
        const fps = Math.round((fpsCounter.current * 1000) / elapsed)
        setMetrics(prev => ({ ...prev, fps }))
        fpsCounter.current = 0
        fpsStartTime.current = now
      }

      animationId = requestAnimationFrame(trackFPS)
    }

    trackFPS()

    return () => {
      if (animationId) {
        cancelAnimationFrame(animationId)
      }
    }
  }, [])

  // Track memory usage (if available)
  useEffect(() => {
    const updateMemoryUsage = () => {
      if ('memory' in performance) {
        const memory = (performance as any).memory
        const memoryUsage = Math.round(memory.usedJSHeapSize / 1024 / 1024) // MB
        setMetrics(prev => ({ ...prev, memoryUsage }))
      }
    }

    updateMemoryUsage()
    const interval = setInterval(updateMemoryUsage, 5000) // Update every 5 seconds

    return () => clearInterval(interval)
  }, [])

  const trackInteraction = (interactionType: string) => {
    const now = Date.now()
    if (lastInteractionTime.current > 0) {
      const interactionTime = now - lastInteractionTime.current
      setMetrics(prev => ({ ...prev, interactionTime }))
    }
    lastInteractionTime.current = now

    // Log interaction for debugging in development
    if (process.env.NODE_ENV === 'development') {
      console.log(`🔍 Interaction: ${interactionType} - Response time: ${metrics.interactionTime}ms`)
    }
  }

  const trackSearchStart = () => {
    searchStartTime.current = Date.now()
  }

  const trackSearchComplete = () => {
    if (searchStartTime.current > 0) {
      const searchResponseTime = Date.now() - searchStartTime.current
      setMetrics(prev => ({ ...prev, searchResponseTime }))
      searchStartTime.current = 0
    }
  }

  const trackTemplatesLoaded = (loaded: number, total: number) => {
    setMetrics(prev => ({ 
      ...prev, 
      templatesLoaded: loaded, 
      totalTemplates: total 
    }))
  }

  const getPerformanceReport = () => {
    const report = {
      ...metrics,
      loadingProgress: metrics.totalTemplates > 0 ? 
        Math.round((metrics.templatesLoaded / metrics.totalTemplates) * 100) : 0,
      performanceGrade: getPerformanceGrade(),
      recommendations: getRecommendations()
    }

    if (process.env.NODE_ENV === 'development') {
      console.table(report)
    }

    return report
  }

  const getPerformanceGrade = (): 'A' | 'B' | 'C' | 'D' | 'F' => {
    let score = 100

    // Penalize slow render times
    if (metrics.renderTime > 100) score -= 10
    if (metrics.renderTime > 300) score -= 20
    if (metrics.renderTime > 500) score -= 30

    // Penalize slow interactions
    if (metrics.interactionTime > 100) score -= 10
    if (metrics.interactionTime > 200) score -= 15

    // Penalize slow search
    if (metrics.searchResponseTime > 50) score -= 5
    if (metrics.searchResponseTime > 150) score -= 15

    // Penalize low FPS
    if (metrics.fps !== null) {
      if (metrics.fps < 30) score -= 20
      if (metrics.fps < 20) score -= 40
    }

    if (score >= 90) return 'A'
    if (score >= 80) return 'B'
    if (score >= 70) return 'C'
    if (score >= 60) return 'D'
    return 'F'
  }

  const getRecommendations = (): string[] => {
    const recommendations: string[] = []

    if (metrics.renderTime > 300) {
      recommendations.push('Consider lazy loading components to improve initial render time')
    }

    if (metrics.interactionTime > 200) {
      recommendations.push('Optimize interaction handlers and consider debouncing user inputs')
    }

    if (metrics.searchResponseTime > 100) {
      recommendations.push('Implement search result caching or optimize search algorithm')
    }

    if (metrics.fps !== null && metrics.fps < 30) {
      recommendations.push('Reduce animations or complex CSS transitions to improve frame rate')
    }

    if (metrics.memoryUsage !== null && metrics.memoryUsage > 100) {
      recommendations.push('Consider implementing virtual scrolling for large lists')
    }

    if (recommendations.length === 0) {
      recommendations.push('Performance looks good! 🚀')
    }

    return recommendations
  }

  // Auto-report performance issues in development
  useEffect(() => {
    if (process.env.NODE_ENV === 'development') {
      const grade = getPerformanceGrade()
      if (grade === 'D' || grade === 'F') {
        console.warn('⚠️ Performance issues detected:', getRecommendations())
      }
    }
  }, [metrics])

  return {
    metrics,
    trackInteraction,
    trackSearchStart,
    trackSearchComplete,
    trackTemplatesLoaded,
    getPerformanceReport,
    performanceGrade: getPerformanceGrade()
  }
}