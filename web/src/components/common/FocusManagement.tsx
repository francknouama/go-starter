// WCAG 2.1 AA Focus Management System
// Comprehensive focus trapping, restoration, and navigation utilities

import { useEffect, useRef, useCallback, useState } from 'react'
import type { ReactNode } from 'react'

// Focusable element selector
const FOCUSABLE_SELECTOR = [
  'button',
  'input',
  'select', 
  'textarea',
  'a[href]',
  '[tabindex]:not([tabindex="-1"])',
  '[contenteditable="true"]',
  'audio[controls]',
  'video[controls]',
  'iframe',
  'embed',
  'object',
  'summary',
  '[role="button"]',
  '[role="link"]',
  '[role="menuitem"]',
  '[role="tab"]',
  '[role="checkbox"]',
  '[role="radio"]',
  '[role="slider"]',
  '[role="spinbutton"]',
  '[role="switch"]',
  '[role="textbox"]'
].join(', ')

/**
 * Get all focusable elements within a container
 */
export function getFocusableElements(container: HTMLElement): HTMLElement[] {
  const elements = Array.from(container.querySelectorAll(FOCUSABLE_SELECTOR)) as HTMLElement[]
  
  return elements.filter(element => {
    // Filter out disabled and hidden elements
    if (element.hasAttribute('disabled') || element.getAttribute('aria-disabled') === 'true') {
      return false
    }
    
    // Check if element is visible
    const style = window.getComputedStyle(element)
    if (style.display === 'none' || style.visibility === 'hidden') {
      return false
    }
    
    // Check if element has valid tabindex
    const tabIndex = element.getAttribute('tabindex')
    if (tabIndex === '-1') {
      return false
    }
    
    return true
  })
}

/**
 * Focus Management Hook for components that need focus control
 */
export function useFocusManagement(isActive: boolean = true) {
  const previousFocusRef = useRef<HTMLElement | null>(null)
  
  const saveFocus = useCallback(() => {
    previousFocusRef.current = document.activeElement as HTMLElement
  }, [])
  
  const restoreFocus = useCallback(() => {
    if (previousFocusRef.current && isActive) {
      try {
        previousFocusRef.current.focus()
      } catch (error) {
        console.warn('Failed to restore focus:', error)
      }
    }
  }, [isActive])
  
  const moveFocusToFirst = useCallback((container: HTMLElement) => {
    const focusableElements = getFocusableElements(container)
    if (focusableElements.length > 0) {
      focusableElements[0].focus()
    }
  }, [])
  
  const moveFocusToLast = useCallback((container: HTMLElement) => {
    const focusableElements = getFocusableElements(container)
    if (focusableElements.length > 0) {
      focusableElements[focusableElements.length - 1].focus()
    }
  }, [])
  
  return {
    saveFocus,
    restoreFocus,
    moveFocusToFirst,
    moveFocusToLast
  }
}

/**
 * Focus Trap Component
 * Traps focus within children and provides escape mechanisms
 */
interface FocusTrapProps {
  children: ReactNode
  isActive?: boolean
  autoFocus?: boolean
  restoreFocus?: boolean
  onEscape?: () => void
  className?: string
}

export function FocusTrap({
  children,
  isActive = true,
  autoFocus = true,
  restoreFocus = true,
  onEscape,
  className = ''
}: FocusTrapProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const { saveFocus, restoreFocus: restore, moveFocusToFirst } = useFocusManagement(isActive)
  
  // Setup focus trap
  useEffect(() => {
    if (!isActive || !containerRef.current) return
    
    // Save current focus
    if (restoreFocus) {
      saveFocus()
    }
    
    // Move focus to first focusable element
    if (autoFocus) {
      moveFocusToFirst(containerRef.current)
    }
    
    const handleKeyDown = (event: KeyboardEvent) => {
      if (!containerRef.current) return
      
      // Handle Escape key
      if (event.key === 'Escape' && onEscape) {
        event.preventDefault()
        onEscape()
        return
      }
      
      // Handle Tab key for focus trapping
      if (event.key === 'Tab') {
        const focusableElements = getFocusableElements(containerRef.current)
        
        if (focusableElements.length === 0) {
          event.preventDefault()
          return
        }
        
        const firstElement = focusableElements[0]
        const lastElement = focusableElements[focusableElements.length - 1]
        
        if (event.shiftKey) {
          // Shift + Tab (moving backward)
          if (document.activeElement === firstElement) {
            event.preventDefault()
            lastElement.focus()
          }
        } else {
          // Tab (moving forward)
          if (document.activeElement === lastElement) {
            event.preventDefault()
            firstElement.focus()
          }
        }
      }
    }
    
    document.addEventListener('keydown', handleKeyDown)
    
    return () => {
      document.removeEventListener('keydown', handleKeyDown)
      
      // Restore focus when unmounting
      if (restoreFocus) {
        restore()
      }
    }
  }, [isActive, autoFocus, restoreFocus, onEscape, saveFocus, restore, moveFocusToFirst])
  
  if (!isActive) {
    return <div className={className}>{children}</div>
  }
  
  return (
    <div ref={containerRef} className={className}>
      {children}
    </div>
  )
}

/**
 * Focus Guard Component
 * Prevents focus from escaping a container
 */
interface FocusGuardProps {
  children: ReactNode
  onFocusOut?: () => void
}

export function FocusGuard({ children, onFocusOut }: FocusGuardProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  
  const handleFocus = useCallback((event: FocusEvent) => {
    if (!containerRef.current) return
    
    const isChildOfContainer = containerRef.current.contains(event.target as Node)
    
    if (!isChildOfContainer && onFocusOut) {
      onFocusOut()
    }
  }, [onFocusOut])
  
  useEffect(() => {
    document.addEventListener('focusin', handleFocus)
    
    return () => {
      document.removeEventListener('focusin', handleFocus)
    }
  }, [handleFocus])
  
  return <div ref={containerRef}>{children}</div>
}

/**
 * Skip Link Component
 * Provides quick navigation for keyboard users
 */
interface SkipLinkProps {
  href: string
  children: ReactNode
  className?: string
}

export function SkipLink({ href, children, className = '' }: SkipLinkProps) {
  const handleClick = useCallback((event: React.MouseEvent) => {
    const target = document.querySelector(href)
    if (target) {
      event.preventDefault()
      
      // Make target focusable if it isn't already
      const htmlTarget = target as HTMLElement
      const originalTabIndex = htmlTarget.getAttribute('tabindex')
      
      if (originalTabIndex === null) {
        htmlTarget.setAttribute('tabindex', '-1')
      }
      
      htmlTarget.focus()
      
      // Restore original tabindex after a short delay
      if (originalTabIndex === null) {
        setTimeout(() => {
          htmlTarget.removeAttribute('tabindex')
        }, 100)
      }
    }
  }, [href])
  
  return (
    <a
      href={href}
      onClick={handleClick}
      className={`
        skip-link sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 
        bg-primary-600 text-white px-4 py-2 rounded-md font-medium text-sm
        focus:z-50 transition-all duration-200
        ${className}
      `}
    >
      {children}
    </a>
  )
}

/**
 * Roving Tab Index Hook
 * For managing focus in lists and grids
 */
export function useRovingTabIndex<T extends HTMLElement>() {
  const elementsRef = useRef<T[]>([])
  const currentIndexRef = useRef(0)
  
  const setElements = useCallback((elements: T[]) => {
    elementsRef.current = elements
    
    // Set initial tabindex values
    elements.forEach((element, index) => {
      element.setAttribute('tabindex', index === 0 ? '0' : '-1')
    })
  }, [])
  
  const focusElement = useCallback((index: number) => {
    const elements = elementsRef.current
    if (index < 0 || index >= elements.length) return
    
    // Update tabindex values
    elements.forEach((element, i) => {
      element.setAttribute('tabindex', i === index ? '0' : '-1')
    })
    
    // Focus the element
    elements[index].focus()
    currentIndexRef.current = index
  }, [])
  
  const handleKeyDown = useCallback((event: KeyboardEvent) => {
    const elements = elementsRef.current
    if (elements.length === 0) return
    
    const currentIndex = currentIndexRef.current
    let newIndex = currentIndex
    
    switch (event.key) {
      case 'ArrowRight':
      case 'ArrowDown':
        event.preventDefault()
        newIndex = (currentIndex + 1) % elements.length
        break
        
      case 'ArrowLeft':
      case 'ArrowUp':
        event.preventDefault()
        newIndex = currentIndex === 0 ? elements.length - 1 : currentIndex - 1
        break
        
      case 'Home':
        event.preventDefault()
        newIndex = 0
        break
        
      case 'End':
        event.preventDefault()
        newIndex = elements.length - 1
        break
        
      default:
        return
    }
    
    focusElement(newIndex)
  }, [focusElement])
  
  return {
    setElements,
    focusElement,
    handleKeyDown,
    currentIndex: currentIndexRef.current
  }
}

/**
 * Focus Visible Polyfill Hook
 * Adds focus-visible support for older browsers
 */
export function useFocusVisible() {
  useEffect(() => {
    let hadKeyboardEvent = true
    
    const onPointerDown = () => {
      hadKeyboardEvent = false
    }
    
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.metaKey || event.altKey || event.ctrlKey) {
        return
      }
      hadKeyboardEvent = true
    }
    
    const onFocus = (event: FocusEvent) => {
      const target = event.target as HTMLElement
      
      if (hadKeyboardEvent || target.matches(':focus-visible')) {
        target.classList.add('focus-visible')
      }
    }
    
    const onBlur = (event: FocusEvent) => {
      const target = event.target as HTMLElement
      target.classList.remove('focus-visible')
    }
    
    document.addEventListener('keydown', onKeyDown, true)
    document.addEventListener('mousedown', onPointerDown, true)
    document.addEventListener('pointerdown', onPointerDown, true)
    document.addEventListener('touchstart', onPointerDown, true)
    document.addEventListener('focus', onFocus, true)
    document.addEventListener('blur', onBlur, true)
    
    return () => {
      document.removeEventListener('keydown', onKeyDown, true)
      document.removeEventListener('mousedown', onPointerDown, true)
      document.removeEventListener('pointerdown', onPointerDown, true)
      document.removeEventListener('touchstart', onPointerDown, true)
      document.removeEventListener('focus', onFocus, true)
      document.removeEventListener('blur', onBlur, true)
    }
  }, [])
}

/**
 * Focus Restoration Utility
 * Automatically saves and restores focus for dynamic content
 */
export class FocusRestorer {
  private previousFocus: HTMLElement | null = null
  
  save(): void {
    this.previousFocus = document.activeElement as HTMLElement
  }
  
  restore(): void {
    if (this.previousFocus) {
      try {
        this.previousFocus.focus()
      } catch (error) {
        console.warn('Failed to restore focus:', error)
      }
    }
  }
  
  clear(): void {
    this.previousFocus = null
  }
}

// Screen Reader Only Component
interface ScreenReaderOnlyProps {
  children: ReactNode
  as?: keyof React.JSX.IntrinsicElements
  id?: string
}

export function ScreenReaderOnly({ children, as: Component = 'div', id }: ScreenReaderOnlyProps) {
  const TagName = Component as React.ElementType
  return (
    <TagName id={id} className="sr-only">
      {children}
    </TagName>
  )
}

// Keyboard Navigation Utilities
// Live Region for announcements
export function LiveRegion({ children, politeness = 'polite' }: { children: ReactNode; politeness?: 'polite' | 'assertive' }) {
  return (
    <div 
      role="status" 
      aria-live={politeness}
      aria-atomic="true"
      className="sr-only"
    >
      {children}
    </div>
  )
}

// Announcement Region for dynamic messages
export function AnnouncementRegion({ message, visible = false }: { message: string; visible?: boolean }) {
  return (
    <div 
      role="alert" 
      aria-live="assertive"
      aria-atomic="true"
      className={visible ? '' : 'sr-only'}
    >
      {message}
    </div>
  )
}

// High contrast preference hook
export function useHighContrast() {
  const [isHighContrast, setIsHighContrast] = useState(false)
  
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-contrast: high)')
    setIsHighContrast(mediaQuery.matches)
    
    const handler = (e: MediaQueryListEvent) => setIsHighContrast(e.matches)
    mediaQuery.addEventListener('change', handler)
    
    return () => mediaQuery.removeEventListener('change', handler)
  }, [])
  
  return isHighContrast
}

// Reduced motion preference hook
export function useReducedMotion() {
  const [prefersReducedMotion, setPrefersReducedMotion] = useState(false)
  
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
    setPrefersReducedMotion(mediaQuery.matches)
    
    const handler = (e: MediaQueryListEvent) => setPrefersReducedMotion(e.matches)
    mediaQuery.addEventListener('change', handler)
    
    return () => mediaQuery.removeEventListener('change', handler)
  }, [])
  
  return prefersReducedMotion
}

export const keyboardNavigation = {
  // Add basic keyboard navigation functions
  handleArrowKeys: (event: KeyboardEvent, items: HTMLElement[], currentIndex: number) => {
    switch (event.key) {
      case 'ArrowDown':
      case 'ArrowRight':
        event.preventDefault()
        const nextIndex = Math.min(currentIndex + 1, items.length - 1)
        items[nextIndex]?.focus()
        return nextIndex
      case 'ArrowUp':
      case 'ArrowLeft':
        event.preventDefault()
        const prevIndex = Math.max(currentIndex - 1, 0)
        items[prevIndex]?.focus()
        return prevIndex
    }
    return currentIndex
  }
}