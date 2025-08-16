import type { ReactNode } from 'react'

interface CompactSelectionCardProps {
  value?: string
  icon?: ReactNode
  label: string
  description?: string
  selected: boolean
  onClick: () => void
}

export default function CompactSelectionCard({
  // value,
  icon,
  label,
  description,
  selected,
  onClick
}: CompactSelectionCardProps) {
  return (
    <button
      onClick={onClick}
      className={`
        w-full p-3 rounded-lg border transition-all duration-200 text-left
        ${selected 
          ? 'border-primary-500 bg-primary-50 shadow-sm' 
          : 'border-gray-200 hover:border-gray-300 hover:shadow-sm bg-white'
        }
      `}
      aria-pressed={selected}
    >
      <div className="flex items-center space-x-3">
        {icon && (
          <div className={`
            flex-shrink-0 w-8 h-8 rounded flex items-center justify-center
            ${selected ? 'bg-primary-100 text-primary-600' : 'bg-gray-100 text-gray-600'}
          `}>
            {icon}
          </div>
        )}
        <div className="flex-1">
          <p className={`text-sm font-medium ${selected ? 'text-primary-900' : 'text-gray-900'}`}>
            {label}
          </p>
          {description && (
            <p className={`text-xs ${selected ? 'text-primary-700' : 'text-gray-500'}`}>
              {description}
            </p>
          )}
        </div>
        {selected && (
          <div className="w-2 h-2 rounded-full bg-primary-600" />
        )}
      </div>
    </button>
  )
}