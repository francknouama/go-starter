import { Switch } from '@headlessui/react'
import { 
  CodeBracketIcon, 
  QuestionMarkCircleIcon, 
  Cog6ToothIcon,
  AdjustmentsHorizontalIcon,
  SparklesIcon
} from '@heroicons/react/24/outline'
import type { DisclosureMode } from '../../types'

interface HeaderProps {
  disclosureMode: DisclosureMode
  onDisclosureModeChange: (mode: DisclosureMode) => void
}

export default function Header({ disclosureMode, onDisclosureModeChange }: HeaderProps) {
  return (
    <header className="bg-white/80 backdrop-blur-xl border-b border-white/20 shadow-lg relative z-20">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo and Title */}
          <div className="flex items-center space-x-2 md:space-x-3">
            <div className="flex items-center justify-center w-8 h-8 md:w-10 md:h-10 bg-gradient-to-r from-blue-600 to-purple-600 rounded-xl shadow-lg">
              <CodeBracketIcon className="w-5 h-5 md:w-6 md:h-6 text-white" />
            </div>
            <div>
              <h1 className="text-lg md:text-xl font-bold text-gray-900">Go Starter</h1>
              <p className="hidden md:block text-sm text-gray-500">Web Project Generator</p>
            </div>
          </div>

          {/* Navigation and Controls */}
          <div className="flex items-center space-x-2 md:space-x-4">
            {/* Status indicator */}
            <div className="hidden sm:flex items-center space-x-2 text-sm">
              <div className="flex items-center space-x-1">
                <div className={`w-2 h-2 rounded-full ${disclosureMode === 'basic' ? 'bg-blue-500' : 'bg-purple-500'}`} />
                <span className="text-gray-600 font-medium">
                  {disclosureMode === 'basic' ? 'Quick Start' : 'Full Control'} Mode
                </span>
              </div>
            </div>

            {/* Advanced Mode Toggle */}
            <button
              onClick={() => onDisclosureModeChange(disclosureMode === 'basic' ? 'advanced' : 'basic')}
              className={`flex items-center space-x-2 px-4 py-2 rounded-xl border transition-all duration-300 transform hover:scale-105 ${
                disclosureMode === 'advanced' 
                  ? 'bg-gradient-to-r from-purple-100/80 to-indigo-100/80 backdrop-blur-md border-purple-300/50 text-purple-700 shadow-lg' 
                  : 'bg-white/60 backdrop-blur-md border-gray-300/50 text-gray-700 hover:bg-white/80 shadow-md'
              }`}
              title={`Switch to ${disclosureMode === 'basic' ? 'Advanced' : 'Basic'} Mode`}
            >
              {disclosureMode === 'advanced' ? (
                <SparklesIcon className="w-4 h-4" />
              ) : (
                <AdjustmentsHorizontalIcon className="w-4 h-4" />
              )}
              <span className="text-sm font-medium hidden md:block">
                {disclosureMode === 'basic' ? 'Advanced' : 'Basic'}
              </span>
            </button>

            {/* Action Buttons - Hidden on mobile */}
            <div className="hidden md:flex items-center space-x-2">
              <button
                className="p-2 text-gray-500 hover:text-gray-600 rounded-xl hover:bg-white/60 backdrop-blur-md transition-all duration-200 hover:shadow-md"
                title="Help & Documentation"
              >
                <QuestionMarkCircleIcon className="w-5 h-5" />
              </button>
              
              <button
                className="p-2 text-gray-500 hover:text-gray-600 rounded-xl hover:bg-white/60 backdrop-blur-md transition-all duration-200 hover:shadow-md"
                title="Settings"
              >
                <Cog6ToothIcon className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>
  )
}