import type { ReactNode } from 'react'
import { CheckCircleIcon } from '@heroicons/react/20/solid'

interface ProjectTypeCardProps {
  value?: string
  icon: ReactNode
  title: string
  description: string
  features?: string[]
  selected: boolean
  onClick: () => void
  'aria-describedby'?: string
}

export default function ProjectTypeCard({
  value,
  icon,
  title,
  description,
  features,
  selected,
  onClick,
  'aria-describedby': ariaDescribedBy
}: ProjectTypeCardProps) {
  const cardId = value ? `project-type-${value}` : undefined
  const descriptionId = cardId ? `${cardId}-description` : undefined

  return (
    <div
      onClick={onClick}
      className={`
        relative rounded-lg border-2 cursor-pointer transition-all duration-200
        focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2
        ${selected 
          ? 'border-primary-500 bg-primary-50 shadow-sm' 
          : 'border-gray-200 hover:border-gray-300 hover:shadow-sm bg-white'
        }
      `}
      role="radio"
      aria-checked={selected}
      aria-labelledby={cardId ? `${cardId}-title` : undefined}
      aria-describedby={ariaDescribedBy || descriptionId}
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault()
          onClick()
        }
      }}
    >
      {/* Selection indicator */}
      {selected && (
        <div className="absolute top-3 right-3" aria-hidden="true">
          <CheckCircleIcon className="w-5 h-5 text-primary-600" />
        </div>
      )}
      
      <div className="p-4">
        <div className="flex items-start space-x-3">
          {/* Icon */}
          <div className={`
            flex-shrink-0 w-10 h-10 rounded-lg flex items-center justify-center
            ${selected ? 'bg-primary-100' : 'bg-gray-100'}
          `} aria-hidden="true">
            <div className={selected ? 'text-primary-600' : 'text-gray-600'}>
              {icon}
            </div>
          </div>
          
          {/* Content */}
          <div className="flex-1">
            <h3 
              id={cardId ? `${cardId}-title` : undefined}
              className={`
                text-sm font-medium mb-1
                ${selected ? 'text-primary-900' : 'text-gray-900'}
              `}
            >
              {title}
            </h3>
            <p 
              id={descriptionId}
              className={`
                text-xs leading-relaxed
                ${selected ? 'text-primary-700' : 'text-gray-600'}
              `}
            >
              {description}
              {selected && (
                <span className="sr-only"> (currently selected)</span>
              )}
            </p>
            
            {/* Features list */}
            {features && features.length > 0 && (
              <ul className="mt-2 space-y-1" aria-label={`${title} features`}>
                {features.map((feature, index) => (
                  <li key={index} className="flex items-center text-xs text-gray-500">
                    <span className="w-1 h-1 bg-gray-500 rounded-full mr-2" aria-hidden="true" />
                    {feature}
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}