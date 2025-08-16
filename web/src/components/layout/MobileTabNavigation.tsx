import { Bars3Icon, FolderOpenIcon } from '@heroicons/react/24/outline'
import type { TabId } from './ResponsiveLayout'

interface MobileTabNavigationProps {
  activeTab: TabId
  onTabChange: (tab: TabId) => void
}

export interface Tab {
  id: TabId
  label: string
  icon: typeof Bars3Icon
}

const tabs: Tab[] = [
  { id: 'configuration', label: 'Configure', icon: Bars3Icon },
  { id: 'files', label: 'Files', icon: FolderOpenIcon }
]

export default function MobileTabNavigation({ 
  activeTab, 
  onTabChange 
}: MobileTabNavigationProps) {
  return (
    <nav 
      className="fixed bottom-0 left-0 right-0 bg-white/95 backdrop-blur-xl border-t border-gray-200 shadow-2xl z-50 safe-area-bottom"
      role="tablist"
      aria-label="Main navigation tabs"
    >
      <div className="grid grid-cols-2 h-16">
        {tabs.map((tab) => {
          const Icon = tab.icon
          const isActive = activeTab === tab.id
          
          return (
            <button
              key={tab.id}
              onClick={() => onTabChange(tab.id)}
              className={`
                flex flex-col items-center justify-center space-y-1 
                transition-all duration-200 relative
                ${isActive 
                  ? 'text-blue-600' 
                  : 'text-gray-500 hover:text-gray-700'
                }
              `}
              role="tab"
              aria-selected={isActive}
              aria-controls={`${tab.id}-panel`}
              id={`${tab.id}-tab`}
              tabIndex={isActive ? 0 : -1}
              onKeyDown={(e) => {
                if (e.key === 'ArrowLeft') {
                  e.preventDefault()
                  const currentIndex = tabs.findIndex(t => t.id === activeTab)
                  const prevIndex = currentIndex > 0 ? currentIndex - 1 : tabs.length - 1
                  onTabChange(tabs[prevIndex].id)
                } else if (e.key === 'ArrowRight') {
                  e.preventDefault()
                  const currentIndex = tabs.findIndex(t => t.id === activeTab)
                  const nextIndex = currentIndex < tabs.length - 1 ? currentIndex + 1 : 0
                  onTabChange(tabs[nextIndex].id)
                }
              }}
            >
              {/* Active indicator */}
              {isActive && (
                <div className="absolute inset-x-0 top-0 h-0.5 bg-blue-600" />
              )}
              
              <Icon className={`w-5 h-5 ${isActive ? 'stroke-2' : ''}`} aria-hidden="true" />
              <span className="text-xs font-medium">
                {tab.label}
                <span className="sr-only">
                  {isActive ? ' (current tab)' : ' tab'}
                </span>
              </span>
            </button>
          )
        })}
      </div>
    </nav>
  )
}