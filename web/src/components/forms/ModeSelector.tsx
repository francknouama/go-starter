import { RadioGroup } from '@headlessui/react'
import { SparklesIcon, RocketLaunchIcon, CheckIcon } from '@heroicons/react/24/outline'
import type { DisclosureMode } from '../../types'

interface ModeSelectorProps {
  mode: DisclosureMode
  onChange: (mode: DisclosureMode) => void
}

interface ModeOption {
  id: DisclosureMode
  title: string
  subtitle: string
  description: string
  icon: typeof SparklesIcon
  features: string[]
  recommended?: boolean
}

export default function ModeSelector({ mode, onChange }: ModeSelectorProps) {
  const modeOptions: ModeOption[] = [
    {
      id: 'basic',
      title: 'Quick Start',
      subtitle: 'Essential options only',
      description: 'Perfect for getting started quickly with sensible defaults. Ideal for beginners or when you need a standard project setup.',
      icon: RocketLaunchIcon,
      features: [
        'Essential configuration options',
        'Popular project templates',
        'Standard architectures',
        'Quick project generation',
        'Minimal decision fatigue'
      ],
      recommended: true
    },
    {
      id: 'advanced',
      title: 'Full Control',
      subtitle: 'All configuration options',
      description: 'Complete customization for experienced developers. Access all architectures, databases, authentication, and deployment options.',
      icon: SparklesIcon,
      features: [
        'All architecture patterns',
        'Database configuration',
        'Authentication systems',
        'Deployment options',
        'Enterprise features'
      ]
    }
  ]

  return (
    <div className="mb-6">
      <div className="mb-4">
        <h3 className="text-sm font-medium text-gray-900 mb-1">Choose Your Experience</h3>
        <p className="text-xs text-gray-600">
          Select how you want to configure your project
        </p>
      </div>

      <RadioGroup value={mode} onChange={onChange}>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
          {modeOptions.map((option) => (
            <RadioGroup.Option
              key={option.id}
              value={option.id}
              className={({ active, checked }) =>
                `${
                  active
                    ? 'ring-2 ring-primary-500 ring-offset-2'
                    : ''
                }
                ${
                  checked
                    ? 'bg-primary-50 border-primary-500'
                    : 'bg-white border-gray-200 hover:border-gray-300'
                }
                relative rounded-lg border-2 p-4 cursor-pointer focus:outline-none transition-all duration-200`
              }
            >
              {({ checked }) => (
                <>
                  {/* Recommended badge */}
                  {option.recommended && !checked && (
                    <div className="absolute -top-2 -right-2">
                      <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 border border-green-200">
                        Recommended
                      </span>
                    </div>
                  )}

                  {/* Selected indicator */}
                  {checked && (
                    <div className="absolute top-3 right-3">
                      <CheckIcon className="h-5 w-5 text-primary-600" />
                    </div>
                  )}

                  <div className="flex items-start space-x-3">
                    {/* Icon */}
                    <div className={`
                      flex-shrink-0 w-10 h-10 rounded-lg flex items-center justify-center
                      ${checked ? 'bg-primary-100' : 'bg-gray-100'}
                    `}>
                      <option.icon className={`
                        w-6 h-6
                        ${checked ? 'text-primary-600' : 'text-gray-600'}
                      `} />
                    </div>

                    {/* Content */}
                    <div className="flex-1">
                      <RadioGroup.Label
                        as="h4"
                        className={`text-sm font-medium ${
                          checked ? 'text-primary-900' : 'text-gray-900'
                        }`}
                      >
                        {option.title}
                      </RadioGroup.Label>
                      <RadioGroup.Description
                        as="p"
                        className={`text-xs mt-0.5 ${
                          checked ? 'text-primary-700' : 'text-gray-500'
                        }`}
                      >
                        {option.subtitle}
                      </RadioGroup.Description>

                      {/* Expanded description */}
                      <p className={`
                        text-xs mt-2 leading-relaxed
                        ${checked ? 'text-primary-700' : 'text-gray-600'}
                      `}>
                        {option.description}
                      </p>

                      {/* Features list */}
                      <ul className="mt-3 space-y-1">
                        {option.features.map((feature, index) => (
                          <li key={index} className="flex items-start text-xs">
                            <CheckIcon className={`
                              w-3 h-3 mr-2 mt-0.5 flex-shrink-0
                              ${checked ? 'text-primary-600' : 'text-gray-500'}
                            `} />
                            <span className={checked ? 'text-primary-700' : 'text-gray-600'}>
                              {feature}
                            </span>
                          </li>
                        ))}
                      </ul>
                    </div>
                  </div>
                </>
              )}
            </RadioGroup.Option>
          ))}
        </div>
      </RadioGroup>

      {/* Mode explanation */}
      <div className={`
        mt-4 p-3 rounded-lg border
        ${mode === 'basic' 
          ? 'bg-blue-50 border-blue-200' 
          : 'bg-purple-50 border-purple-200'
        }
      `}>
        <p className={`
          text-xs leading-relaxed
          ${mode === 'basic' ? 'text-blue-700' : 'text-purple-700'}
        `}>
          {mode === 'basic' 
            ? '💡 Perfect choice! Basic mode shows only essential options to get you started quickly. You can always switch to advanced mode later if needed.'
            : '⚡ Advanced mode unlocked! You now have access to all configuration options including databases, authentication, and deployment settings.'
          }
        </p>
      </div>
    </div>
  )
}