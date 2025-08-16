import { Fragment, useState, useEffect, useRef } from 'react'
import type { ReactNode } from 'react'
import { Popover, Transition } from '@headlessui/react'
import { InformationCircleIcon, QuestionMarkCircleIcon } from '@heroicons/react/20/solid'

interface HelpTooltipProps {
  content: string | ReactNode
  trigger?: ReactNode
  position?: 'top' | 'bottom' | 'left' | 'right' | 'auto'
  maxWidth?: string
  showOnHover?: boolean
  variant?: 'info' | 'help'
  delay?: number
  className?: string
}

export default function HelpTooltip({ 
  content, 
  trigger,
  position = 'auto',
  maxWidth = 'max-w-xs',
  showOnHover = true,
  variant = 'info',
  delay = 200,
  className = ''
}: HelpTooltipProps) {
  const [isVisible, setIsVisible] = useState(false)
  const [actualPosition, setActualPosition] = useState(position)
  const triggerRef = useRef<HTMLDivElement>(null)
  const tooltipRef = useRef<HTMLDivElement>(null)
  const timeoutRef = useRef<NodeJS.Timeout | null>(null)

  // Auto-position tooltip based on viewport
  useEffect(() => {
    if (position === 'auto' && triggerRef.current && isVisible) {
      const triggerRect = triggerRef.current.getBoundingClientRect()
      const viewportWidth = window.innerWidth
      const viewportHeight = window.innerHeight
      
      // Determine best position
      let bestPosition: 'top' | 'bottom' | 'left' | 'right' = 'top'
      
      // Check space above and below
      const spaceAbove = triggerRect.top
      const spaceBelow = viewportHeight - triggerRect.bottom
      const spaceLeft = triggerRect.left
      const spaceRight = viewportWidth - triggerRect.right
      
      // Prefer top/bottom first, then left/right
      if (spaceBelow > 150) {
        bestPosition = 'bottom'
      } else if (spaceAbove > 150) {
        bestPosition = 'top'
      } else if (spaceRight > 200) {
        bestPosition = 'right'
      } else if (spaceLeft > 200) {
        bestPosition = 'left'
      }
      
      setActualPosition(bestPosition)
    } else if (position !== 'auto') {
      setActualPosition(position)
    }
  }, [position, isVisible])

  const handleMouseEnter = () => {
    if (showOnHover) {
      clearTimeout(timeoutRef.current)
      timeoutRef.current = setTimeout(() => setIsVisible(true), delay)
    }
  }

  const handleMouseLeave = () => {
    if (showOnHover) {
      clearTimeout(timeoutRef.current)
      timeoutRef.current = setTimeout(() => setIsVisible(false), 100)
    }
  }

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Escape' && isVisible) {
      setIsVisible(false)
    }
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault()
      setIsVisible(!isVisible)
    }
  }

  const handleFocus = () => {
    if (showOnHover) {
      clearTimeout(timeoutRef.current)
      timeoutRef.current = setTimeout(() => setIsVisible(true), delay)
    }
  }

  const handleBlur = () => {
    if (showOnHover) {
      clearTimeout(timeoutRef.current)
      timeoutRef.current = setTimeout(() => setIsVisible(false), 100)
    }
  }

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [])
  const Icon = variant === 'help' ? QuestionMarkCircleIcon : InformationCircleIcon
  
  const defaultTrigger = (
    <Icon 
      className="w-4 h-4 text-gray-500 hover:text-gray-600 cursor-help transition-colors" 
      aria-hidden="true"
    />
  )

  const getPositionClasses = () => {
    switch (actualPosition) {
      case 'top':
        return 'bottom-full left-1/2 transform -translate-x-1/2 mb-2'
      case 'bottom':
        return 'top-full left-1/2 transform -translate-x-1/2 mt-2'
      case 'left':
        return 'right-full top-1/2 transform -translate-y-1/2 mr-2'
      case 'right':
        return 'left-full top-1/2 transform -translate-y-1/2 ml-2'
      default:
        return 'bottom-full left-1/2 transform -translate-x-1/2 mb-2'
    }
  }

  const getArrowClasses = () => {
    switch (actualPosition) {
      case 'top':
        return 'top-full left-1/2 transform -translate-x-1/2 border-t-gray-800'
      case 'bottom':
        return 'bottom-full left-1/2 transform -translate-x-1/2 border-b-gray-800'
      case 'left':
        return 'left-full top-1/2 transform -translate-y-1/2 border-l-gray-800'
      case 'right':
        return 'right-full top-1/2 transform -translate-y-1/2 border-r-gray-800'
      default:
        return 'top-full left-1/2 transform -translate-x-1/2 border-t-gray-800'
    }
  }

  // Generate unique IDs for accessibility
  const tooltipId = `tooltip-${useRef(Math.random().toString(36).substr(2, 9)).current}`
  const triggerId = `trigger-${useRef(Math.random().toString(36).substr(2, 9)).current}`

  if (showOnHover) {
    return (
      <div className={`relative inline-flex ${className}`} ref={triggerRef}>
        <button
          id={triggerId}
          type="button"
          className="cursor-help inline-flex items-center justify-center p-0.5 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-1"
          onMouseEnter={handleMouseEnter}
          onMouseLeave={handleMouseLeave}
          onFocus={handleFocus}
          onBlur={handleBlur}
          onKeyDown={handleKeyDown}
          aria-describedby={isVisible ? tooltipId : undefined}
          aria-expanded={isVisible}
          aria-label={typeof content === 'string' ? content : 'Help information'}
        >
          {trigger || defaultTrigger}
        </button>
        
        {/* Accessible tooltip */}
        <Transition
          show={isVisible}
          enter="transition-all duration-200"
          enterFrom="opacity-0 scale-95"
          enterTo="opacity-100 scale-100"
          leave="transition-all duration-150"
          leaveFrom="opacity-100 scale-100"
          leaveTo="opacity-0 scale-95"
        >
          <div 
            id={tooltipId}
            ref={tooltipRef}
            role="tooltip"
            className={`
              absolute ${getPositionClasses()} z-50
            `}
          >
            <div className={`
              bg-gray-800 text-white text-sm rounded-lg p-3 shadow-lg border border-gray-700
              ${maxWidth} backdrop-blur-sm
            `}>
              {content}
              {/* Arrow */}
              <div className={`
                absolute w-0 h-0 
                border-l-4 border-l-transparent
                border-r-4 border-r-transparent
                border-t-4 ${getArrowClasses()}
              `} />
            </div>
          </div>
        </Transition>
      </div>
    )
  }

  // Click-based tooltip using Headless UI Popover
  return (
    <Popover className="relative inline-flex">
      {({ open }) => (
        <>
          <Popover.Button 
            className="cursor-help inline-flex items-center justify-center p-0.5 rounded focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-1"
            aria-expanded={open}
            aria-describedby={open ? tooltipId : undefined}
            aria-label={typeof content === 'string' ? content : 'Help information'}
          >
            {trigger || defaultTrigger}
          </Popover.Button>

          <Transition
            as={Fragment}
            enter="transition ease-out duration-200"
            enterFrom="opacity-0 translate-y-1"
            enterTo="opacity-100 translate-y-0"
            leave="transition ease-in duration-150"
            leaveFrom="opacity-100 translate-y-0"
            leaveTo="opacity-0 translate-y-1"
          >
            <Popover.Panel 
              id={tooltipId}
              role="tooltip"
              className={`absolute ${getPositionClasses()} z-50`}
            >
              <div className={`
                bg-gray-800 text-white text-sm rounded-lg p-3 shadow-lg
                ${maxWidth}
              `}>
                {content}
                {/* Arrow */}
                <div className={`
                  absolute w-0 h-0 
                  border-l-4 border-l-transparent
                  border-r-4 border-r-transparent
                  border-t-4 ${getArrowClasses()}
                `} />
              </div>
            </Popover.Panel>
          </Transition>
        </>
      )}
    </Popover>
  )
}