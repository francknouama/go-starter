import { XMarkIcon } from '@heroicons/react/24/outline'

interface KeyboardShortcutsOverlayProps {
  isOpen: boolean
  onClose: () => void
}

const shortcutCategories = [
  {
    category: 'Navigation',
    shortcuts: [
      { key: '?', description: 'Show this help dialog' },
      { key: 'Tab', description: 'Move to next field' },
      { key: 'Shift + Tab', description: 'Move to previous field' },
      { key: 'Esc', description: 'Close dialogs and modals' },
    ]
  },
  {
    category: 'Actions', 
    shortcuts: [
      { key: 'Ctrl/Cmd + Enter', description: 'Generate project' },
      { key: 'Ctrl/Cmd + R', description: 'Reset form' },
      { key: 'Ctrl/Cmd + D', description: 'Toggle dry run mode' },
      { key: 'Ctrl/Cmd + S', description: 'Save configuration' },
    ]
  },
  {
    category: 'Quick Selection',
    shortcuts: [
      { key: '1', description: 'Select CLI project type' },
      { key: '2', description: 'Select Web API project type' },
      { key: '3', description: 'Select Library project type' },
      { key: '4', description: 'Select Lambda project type' },
      { key: 'B', description: 'Switch to Basic mode' },
      { key: 'A', description: 'Switch to Advanced mode' },
    ]
  },
  {
    category: 'Framework Selection',
    shortcuts: [
      { key: 'G', description: 'Select Gin framework' },
      { key: 'E', description: 'Select Echo framework' },
      { key: 'F', description: 'Select Fiber framework' },
      { key: 'C', description: 'Select Chi framework' },
    ]
  },
  {
    category: 'Architecture Patterns',
    shortcuts: [
      { key: 'Shift + 1', description: 'Select Standard architecture' },
      { key: 'Shift + 2', description: 'Select Clean architecture' },
      { key: 'Shift + 3', description: 'Select DDD architecture' },
      { key: 'Shift + 4', description: 'Select Hexagonal architecture' },
    ]
  },
  {
    category: 'Help & Accessibility',
    shortcuts: [
      { key: 'H', description: 'Show Quick Start guide' },
      { key: 'I', description: 'Toggle help tooltips' },
      { key: 'Ctrl/Cmd + K', description: 'Open command palette' },
      { key: 'Alt + ?', description: 'Accessibility help' },
    ]
  }
]

const proTips = [
  'Hold Shift while pressing number keys to select architecture patterns',
  'Use Tab to navigate through form fields efficiently',
  'Press Esc to quickly close any open modal or dialog',
  'Ctrl/Cmd + Enter works from any field to generate your project',
  'Use dry run mode (Ctrl/Cmd + D) to preview without creating files',
]

export default function KeyboardShortcutsOverlay({ isOpen, onClose }: KeyboardShortcutsOverlayProps) {
  if (!isOpen) return null

  const isMac = typeof navigator !== 'undefined' && navigator.platform.toUpperCase().indexOf('MAC') >= 0

  const formatKey = (key: string) => {
    return key.replace('Ctrl/Cmd', isMac ? 'Cmd' : 'Ctrl')
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <div>
            <h2 className="text-xl font-semibold text-gray-900">Keyboard Shortcuts</h2>
            <p className="text-sm text-gray-600 mt-1">
              Speed up your workflow with these keyboard shortcuts
            </p>
          </div>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-600 transition-colors"
          >
            <XMarkIcon className="w-6 h-6" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6">
          {/* Shortcuts Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
            {shortcutCategories.map((category, index) => (
              <div key={index} className="space-y-3">
                <h3 className="font-semibold text-gray-900 text-sm uppercase tracking-wide">
                  {category.category}
                </h3>
                <div className="space-y-2">
                  {category.shortcuts.map((shortcut, shortcutIndex) => (
                    <div key={shortcutIndex} className="flex items-center justify-between">
                      <span className="text-sm text-gray-600">{shortcut.description}</span>
                      <kbd className="px-2 py-1 text-xs font-mono bg-gray-100 border border-gray-200 rounded">
                        {formatKey(shortcut.key)}
                      </kbd>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>

          {/* Pro Tips */}
          <div className="border-t border-gray-200 pt-6">
            <h3 className="font-semibold text-gray-900 mb-4 flex items-center gap-2">
              <span className="bg-yellow-100 text-yellow-800 text-xs font-medium px-2 py-1 rounded">
                PRO TIPS
              </span>
              Power User Shortcuts
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {proTips.map((tip, index) => (
                <div key={index} className="flex items-start gap-3 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                  <span className="flex-shrink-0 w-5 h-5 bg-yellow-500 text-white rounded-full flex items-center justify-center text-xs font-bold">
                    {index + 1}
                  </span>
                  <p className="text-sm text-yellow-800">{tip}</p>
                </div>
              ))}
            </div>
          </div>

          {/* Accessibility Note */}
          <div className="border-t border-gray-200 pt-6 mt-6">
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <h4 className="font-medium text-blue-900 mb-2">Accessibility Features</h4>
              <p className="text-sm text-blue-800 mb-3">
                go-starter is designed to be fully accessible. All functionality is available via keyboard navigation.
              </p>
              <ul className="text-sm text-blue-800 space-y-1">
                <li>• Screen reader compatible with ARIA labels</li>
                <li>• High contrast mode support</li>
                <li>• Focus indicators for all interactive elements</li>
                <li>• Alternative text for all images and icons</li>
              </ul>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between p-6 border-t border-gray-200 bg-gray-50">
          <div className="text-sm text-gray-600">
            Press <kbd className="px-2 py-1 text-xs font-mono bg-gray-200 border border-gray-300 rounded">?</kbd> anytime to reopen this help
          </div>
          <button
            onClick={onClose}
            className="bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700 transition-colors text-sm"
          >
            Got it!
          </button>
        </div>
      </div>
    </div>
  )
}