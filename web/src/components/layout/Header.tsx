import { Switch } from '@headlessui/react'
import { CodeBracketIcon, QuestionMarkCircleIcon, Cog6ToothIcon } from '@heroicons/react/24/outline'
import type { DisclosureMode } from '../../types'

interface HeaderProps {
  disclosureMode: DisclosureMode
  onDisclosureModeChange: (mode: DisclosureMode) => void
}

export default function Header({ disclosureMode, onDisclosureModeChange }: HeaderProps) {
  return (
    <header className="bg-white border-b border-gray-200 shadow-sm">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo and Title */}
          <div className="flex items-center space-x-3">
            <div className="flex items-center justify-center w-10 h-10 bg-primary-600 rounded-lg">
              <CodeBracketIcon className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-900">Go Starter</h1>
              <p className="text-sm text-gray-500">Web Project Generator</p>
            </div>
          </div>

          {/* Navigation and Controls */}
          <div className="flex items-center space-x-6">
            {/* Mode Toggle */}
            <div className="flex items-center space-x-3">
              <span className="text-sm font-medium text-gray-700">
                {disclosureMode === 'basic' ? 'Basic Mode' : 'Advanced Mode'}
              </span>
              <Switch
                checked={disclosureMode === 'advanced'}
                onChange={(checked) => onDisclosureModeChange(checked ? 'advanced' : 'basic')}
                className={`${
                  disclosureMode === 'advanced' ? 'bg-primary-600' : 'bg-gray-200'
                } relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2`}
              >
                <span className="sr-only">Toggle advanced mode</span>
                <span
                  className={`${
                    disclosureMode === 'advanced' ? 'translate-x-6' : 'translate-x-1'
                  } inline-block h-4 w-4 transform rounded-full bg-white transition-transform`}
                />
              </Switch>
            </div>

            {/* Action Buttons */}
            <div className="flex items-center space-x-2">
              <button
                className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"
                title="Help & Documentation"
              >
                <QuestionMarkCircleIcon className="w-5 h-5" />
              </button>
              
              <button
                className="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-100 transition-colors"
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