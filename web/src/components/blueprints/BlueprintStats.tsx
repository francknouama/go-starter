import { CheckBadgeIcon, ChartBarIcon, CubeIcon, SparklesIcon } from '@heroicons/react/20/solid'
import { TEMPLATE_CATEGORIES, COMPLEXITY_LEVELS, PROJECT_TEMPLATES } from '../../data/projectTemplates'

interface BlueprintStatsProps {
  className?: string
}

export default function BlueprintStats({ className = '' }: BlueprintStatsProps) {
  const totalBlueprints = PROJECT_TEMPLATES.length
  const productionReady = PROJECT_TEMPLATES.length // All are production ready!
  const categories = TEMPLATE_CATEGORIES.length
  const avgPopularity = Math.round(PROJECT_TEMPLATES.reduce((sum, t) => sum + t.popularity, 0) / PROJECT_TEMPLATES.length * 10) / 10

  const stats = [
    {
      label: 'Production Blueprints',
      value: totalBlueprints,
      subtext: 'All production-ready',
      icon: CheckBadgeIcon,
      color: 'emerald',
      percentage: 100
    },
    {
      label: 'Categories Covered',
      value: categories,
      subtext: '100% coverage',
      icon: CubeIcon,
      color: 'blue',
      percentage: 100
    },
    {
      label: 'Average Rating',
      value: `${avgPopularity}/10`,
      subtext: 'Community approved',
      icon: ChartBarIcon,
      color: 'purple',
      percentage: avgPopularity * 10
    },
    {
      label: 'Enterprise Features',
      value: '100%',
      subtext: 'Real-world patterns',
      icon: SparklesIcon,
      color: 'amber',
      percentage: 100
    }
  ]

  const getColorClasses = (color: string) => {
    switch (color) {
      case 'emerald':
        return {
          bg: 'bg-emerald-50',
          border: 'border-emerald-200',
          text: 'text-emerald-600',
          icon: 'text-emerald-500',
          bar: 'bg-emerald-500'
        }
      case 'blue':
        return {
          bg: 'bg-blue-50',
          border: 'border-blue-200',
          text: 'text-blue-600',
          icon: 'text-blue-500',
          bar: 'bg-blue-500'
        }
      case 'purple':
        return {
          bg: 'bg-purple-50',
          border: 'border-purple-200',
          text: 'text-purple-600',
          icon: 'text-purple-500',
          bar: 'bg-purple-500'
        }
      case 'amber':
        return {
          bg: 'bg-amber-50',
          border: 'border-amber-200',
          text: 'text-amber-600',
          icon: 'text-amber-500',
          bar: 'bg-amber-500'
        }
      default:
        return {
          bg: 'bg-gray-50',
          border: 'border-gray-200',
          text: 'text-gray-600',
          icon: 'text-gray-500',
          bar: 'bg-gray-500'
        }
    }
  }

  return (
    <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 ${className}`}>
      {stats.map((stat, index) => {
        const colors = getColorClasses(stat.color)
        const Icon = stat.icon
        
        return (
          <div 
            key={index}
            className={`${colors.bg} ${colors.border} border rounded-2xl p-6 transition-all duration-300 hover:shadow-lg hover:scale-105 relative overflow-hidden`}
          >
            {/* Background Pattern */}
            <div className="absolute top-0 right-0 w-20 h-20 opacity-5">
              <Icon className={`w-full h-full ${colors.icon}`} />
            </div>
            
            {/* Content */}
            <div className="relative z-10">
              <div className="flex items-center justify-between mb-3">
                <div className={`p-2 rounded-lg ${colors.bg} ${colors.border} border shadow-sm`}>
                  <Icon className={`w-5 h-5 ${colors.icon}`} />
                </div>
                {stat.percentage === 100 && (
                  <div className="bg-green-100 text-green-700 px-2 py-1 rounded-full text-xs font-bold flex items-center gap-1">
                    <div className="w-1.5 h-1.5 bg-green-500 rounded-full animate-pulse"></div>
                    Complete
                  </div>
                )}
              </div>
              
              <div className="space-y-1">
                <div className={`text-2xl font-bold ${colors.text}`}>
                  {stat.value}
                </div>
                <div className="text-sm font-semibold text-gray-700">
                  {stat.label}
                </div>
                <div className="text-xs text-gray-500">
                  {stat.subtext}
                </div>
              </div>
              
              {/* Progress bar */}
              <div className="mt-4 bg-white rounded-full h-1.5 overflow-hidden shadow-inner">
                <div 
                  className={`h-full ${colors.bar} transition-all duration-1000 ease-out rounded-full`}
                  style={{ width: `${stat.percentage}%` }}
                />
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}