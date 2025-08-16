import type { ReactNode } from 'react'

interface SelectionGridProps {
  columns?: 1 | 2 | 3
  gap?: 'sm' | 'md' | 'lg'
  children: ReactNode
}

export default function SelectionGrid({ 
  columns = 2, 
  gap = 'md',
  children 
}: SelectionGridProps) {
  const getGridClass = () => {
    switch (columns) {
      case 1:
        return 'grid-cols-1'
      case 2:
        return 'grid-cols-1 md:grid-cols-2'
      case 3:
        return 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3'
      default:
        return 'grid-cols-1 md:grid-cols-2'
    }
  }

  const getGapClass = () => {
    switch (gap) {
      case 'sm':
        return 'gap-2'
      case 'md':
        return 'gap-3 md:gap-4'
      case 'lg':
        return 'gap-4 md:gap-6'
      default:
        return 'gap-3 md:gap-4'
    }
  }

  return (
    <div className={`grid ${getGridClass()} ${getGapClass()}`}>
      {children}
    </div>
  )
}