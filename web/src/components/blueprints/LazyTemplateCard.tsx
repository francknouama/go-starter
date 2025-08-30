import { useState, useRef, useEffect } from 'react'
import { ProjectTemplate } from '../../data/projectTemplates'
import TemplateCard from '../templates/TemplateCard'

interface LazyTemplateCardProps {
  template: ProjectTemplate
  onSelect: () => void
  onPreview: () => void
  displayMode?: 'compact' | 'standard' | 'detailed'
  layoutMode?: 'grid' | 'compact'
  className?: string
}

export default function LazyTemplateCard({ 
  template,
  onSelect,
  onPreview,
  displayMode = 'standard',
  layoutMode = 'grid',
  className = ''
}: LazyTemplateCardProps) {
  const [isVisible, setIsVisible] = useState(false)
  const [isLoaded, setIsLoaded] = useState(false)
  const elementRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        const entry = entries[0]
        if (entry.isIntersecting) {
          setIsVisible(true)
          // Load the actual component after a small delay to prevent layout shift
          setTimeout(() => setIsLoaded(true), 50)
        }
      },
      {
        rootMargin: '50px 0px', // Load when element is 50px away from viewport
        threshold: 0.1
      }
    )

    if (elementRef.current) {
      observer.observe(elementRef.current)
    }

    return () => {
      if (elementRef.current) {
        observer.unobserve(elementRef.current)
      }
    }
  }, [])

  // Skeleton loader
  const SkeletonCard = () => (
    <div className={`bg-white rounded-2xl border border-gray-200 overflow-hidden animate-pulse ${className}`}>
      {/* Header skeleton */}
      <div className="h-24 bg-gradient-to-r from-gray-100 to-gray-200">
        <div className="p-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-gray-300 rounded-xl" />
            <div className="space-y-2">
              <div className="w-16 h-3 bg-gray-300 rounded" />
              <div className="w-20 h-3 bg-gray-300 rounded" />
            </div>
          </div>
          <div className="w-16 h-6 bg-gray-300 rounded-lg" />
        </div>
      </div>

      {/* Content skeleton */}
      <div className="p-5 space-y-4">
        <div className="space-y-2">
          <div className="w-3/4 h-5 bg-gray-300 rounded" />
          <div className="w-full h-4 bg-gray-300 rounded" />
          <div className="w-2/3 h-4 bg-gray-300 rounded" />
        </div>

        <div className="space-y-2">
          <div className="w-20 h-3 bg-gray-300 rounded" />
          <div className="flex gap-2">
            <div className="w-16 h-6 bg-gray-300 rounded-lg" />
            <div className="w-20 h-6 bg-gray-300 rounded-lg" />
            <div className="w-12 h-6 bg-gray-300 rounded-lg" />
          </div>
        </div>
      </div>

      {/* Actions skeleton */}
      <div className="px-5 pb-5 space-y-2">
        <div className="w-full h-10 bg-gray-300 rounded-lg" />
        <div className="w-full h-8 bg-gray-300 rounded-lg" />
      </div>
    </div>
  )

  return (
    <div ref={elementRef} className={className}>
      {isVisible && isLoaded ? (
        <TemplateCard
          template={template}
          onSelect={onSelect}
          onPreview={onPreview}
          displayMode={displayMode}
          layoutMode={layoutMode}
        />
      ) : (
        <SkeletonCard />
      )}
    </div>
  )
}