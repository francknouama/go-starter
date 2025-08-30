import { useState } from 'react'
import { 
  AcademicCapIcon, 
  BeakerIcon, 
  BuildingLibraryIcon, 
  RocketLaunchIcon,
  ChevronRightIcon,
  InformationCircleIcon
} from '@heroicons/react/20/solid'

type ExperienceLevel = 'beginner' | 'intermediate' | 'advanced' | 'expert'

interface ProgressiveDisclosureProps {
  onLevelChange: (level: ExperienceLevel) => void
  currentLevel: ExperienceLevel
  className?: string
}

export default function ProgressiveDisclosure({ 
  onLevelChange, 
  currentLevel, 
  className = '' 
}: ProgressiveDisclosureProps) {
  const [showDetails, setShowDetails] = useState(false)

  const experienceLevels = [
    {
      id: 'beginner' as ExperienceLevel,
      name: 'Beginner',
      description: 'New to Go development',
      icon: AcademicCapIcon,
      color: '#10B981',
      bgColor: 'bg-emerald-50',
      borderColor: 'border-emerald-200',
      textColor: 'text-emerald-700',
      iconColor: 'text-emerald-500',
      features: [
        'Simple project structures',
        'Basic patterns and concepts',
        'Guided documentation',
        'Learning-focused examples'
      ],
      blueprintsShown: 'Simple CLI, Library, Lambda Function',
      recommendedFor: 'Learning Go, quick prototypes, simple utilities'
    },
    {
      id: 'intermediate' as ExperienceLevel,
      name: 'Intermediate',
      description: 'Familiar with Go basics',
      icon: BeakerIcon,
      color: '#3B82F6',
      bgColor: 'bg-blue-50',
      borderColor: 'border-blue-200',
      textColor: 'text-blue-700',
      iconColor: 'text-blue-500',
      features: [
        'Production-ready patterns',
        'Database integration',
        'Middleware and routing',
        'Testing and CI/CD'
      ],
      blueprintsShown: 'Web APIs, Standard CLI, Serverless, Monolith',
      recommendedFor: 'Building APIs, production services, full applications'
    },
    {
      id: 'advanced' as ExperienceLevel,
      name: 'Advanced',
      description: 'Experienced Go developer',
      icon: BuildingLibraryIcon,
      color: '#F59E0B',
      bgColor: 'bg-amber-50',
      borderColor: 'border-amber-200',
      textColor: 'text-amber-700',
      iconColor: 'text-amber-500',
      features: [
        'Clean architecture patterns',
        'Microservices design',
        'Enterprise features',
        'Complex integrations'
      ],
      blueprintsShown: 'Clean Architecture, gRPC Gateway, Microservices',
      recommendedFor: 'Enterprise systems, distributed architectures, scalable services'
    },
    {
      id: 'expert' as ExperienceLevel,
      name: 'Expert',
      description: 'Go architecture specialist',
      icon: RocketLaunchIcon,
      color: '#EF4444',
      bgColor: 'bg-red-50',
      borderColor: 'border-red-200',
      textColor: 'text-red-700',
      iconColor: 'text-red-500',
      features: [
        'All blueprint types',
        'Advanced customization',
        'Performance optimization',
        'Custom patterns'
      ],
      blueprintsShown: 'All 12 production blueprints with advanced options',
      recommendedFor: 'Complex systems, performance-critical apps, custom solutions'
    }
  ]

  const selectedLevel = experienceLevels.find(level => level.id === currentLevel)

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold text-gray-900">Experience Level</h3>
          <p className="text-sm text-gray-600">
            Choose your Go experience to see relevant blueprints
          </p>
        </div>
        <button
          onClick={() => setShowDetails(!showDetails)}
          className="flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 transition-colors duration-200"
        >
          <InformationCircleIcon className="w-4 h-4" />
          {showDetails ? 'Hide details' : 'Show details'}
        </button>
      </div>

      {/* Level Selection Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
        {experienceLevels.map((level) => {
          const Icon = level.icon
          const isSelected = currentLevel === level.id
          
          return (
            <button
              key={level.id}
              onClick={() => onLevelChange(level.id)}
              className={`
                relative p-4 rounded-xl border-2 text-left transition-all duration-200 hover:scale-105
                ${isSelected 
                  ? `${level.bgColor} ${level.borderColor} shadow-lg` 
                  : 'bg-white border-gray-200 hover:border-gray-300 hover:shadow-md'
                }
              `}
            >
              {/* Selection indicator */}
              {isSelected && (
                <div className="absolute top-2 right-2">
                  <div className={`w-3 h-3 rounded-full ${level.iconColor.replace('text-', 'bg-')} animate-pulse`} />
                </div>
              )}

              <div className="flex items-center gap-3 mb-2">
                <div className={`p-2 rounded-lg ${isSelected ? 'bg-white shadow-sm' : level.bgColor}`}>
                  <Icon className={`w-5 h-5 ${isSelected ? level.iconColor : 'text-gray-500'}`} />
                </div>
                <div>
                  <div className={`font-semibold ${isSelected ? level.textColor : 'text-gray-700'}`}>
                    {level.name}
                  </div>
                  <div className="text-xs text-gray-500">
                    {level.description}
                  </div>
                </div>
              </div>

              {isSelected && (
                <div className="mt-2">
                  <ChevronRightIcon className={`w-4 h-4 ${level.iconColor}`} />
                </div>
              )}
            </button>
          )
        })}
      </div>

      {/* Detailed Information */}
      {showDetails && selectedLevel && (
        <div className={`${selectedLevel.bgColor} ${selectedLevel.borderColor} border rounded-xl p-6 space-y-4 transition-all duration-300`}>
          <div className="flex items-center gap-3">
            <selectedLevel.icon className={`w-6 h-6 ${selectedLevel.iconColor}`} />
            <h4 className={`text-lg font-semibold ${selectedLevel.textColor}`}>
              {selectedLevel.name} Level Details
            </h4>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            {/* Features */}
            <div>
              <h5 className={`font-medium ${selectedLevel.textColor} mb-2`}>
                What you'll get:
              </h5>
              <ul className="space-y-1">
                {selectedLevel.features.map((feature, index) => (
                  <li key={index} className="flex items-center gap-2 text-sm text-gray-700">
                    <div className={`w-1.5 h-1.5 rounded-full ${selectedLevel.iconColor.replace('text-', 'bg-')}`} />
                    {feature}
                  </li>
                ))}
              </ul>
            </div>

            {/* Blueprint Focus */}
            <div>
              <h5 className={`font-medium ${selectedLevel.textColor} mb-2`}>
                Recommended blueprints:
              </h5>
              <p className="text-sm text-gray-700 mb-3">
                {selectedLevel.blueprintsShown}
              </p>
              <h5 className={`font-medium ${selectedLevel.textColor} mb-2`}>
                Perfect for:
              </h5>
              <p className="text-sm text-gray-700">
                {selectedLevel.recommendedFor}
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}