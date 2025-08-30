import { SparklesIcon, TrophyIcon, CheckBadgeIcon } from '@heroicons/react/20/solid'

interface AchievementBannerProps {
  className?: string
}

export default function AchievementBanner({ className = '' }: AchievementBannerProps) {
  return (
    <div className={`relative overflow-hidden bg-gradient-to-r from-emerald-500 via-blue-500 to-purple-500 rounded-2xl p-6 ${className}`}>
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute inset-0" style={{
          backgroundImage: `url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.3'%3E%3Ccircle cx='30' cy='30' r='2'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")`,
          backgroundSize: '30px 30px'
        }} />
      </div>

      {/* Floating Elements */}
      <div className="absolute top-4 right-6">
        <div className="animate-bounce">
          <TrophyIcon className="w-8 h-8 text-yellow-300" />
        </div>
      </div>
      
      <div className="absolute bottom-4 right-12">
        <div className="animate-pulse">
          <SparklesIcon className="w-6 h-6 text-white opacity-60" />
        </div>
      </div>

      <div className="absolute top-6 left-1/3">
        <div className="animate-ping">
          <CheckBadgeIcon className="w-5 h-5 text-green-300 opacity-75" />
        </div>
      </div>

      {/* Main Content */}
      <div className="relative z-10">
        <div className="flex items-center gap-3 mb-3">
          <div className="bg-white/20 backdrop-blur-sm rounded-full p-2">
            <TrophyIcon className="w-6 h-6 text-yellow-300" />
          </div>
          <div className="bg-white/20 backdrop-blur-sm px-4 py-1.5 rounded-full">
            <span className="text-white text-sm font-bold tracking-wider">HISTORIC ACHIEVEMENT</span>
          </div>
        </div>
        
        <h2 className="text-3xl font-bold text-white mb-2">
          100% Production Coverage Achieved! 🎉
        </h2>
        
        <p className="text-white/90 text-lg mb-4 max-w-2xl leading-relaxed">
          All <strong>12 blueprints</strong> are now production-ready with comprehensive testing, 
          real-world patterns, and enterprise-grade features. This represents months of development 
          and the completion of our Extended Blueprint System initiative.
        </p>

        {/* Stats */}
        <div className="flex items-center gap-6 text-white/90">
          <div className="flex items-center gap-2">
            <CheckBadgeIcon className="w-5 h-5 text-green-300" />
            <span className="text-sm font-semibold">12 Production Blueprints</span>
          </div>
          <div className="flex items-center gap-2">
            <CheckBadgeIcon className="w-5 h-5 text-blue-300" />
            <span className="text-sm font-semibold">6 Categories Covered</span>
          </div>
          <div className="flex items-center gap-2">
            <CheckBadgeIcon className="w-5 h-5 text-purple-300" />
            <span className="text-sm font-semibold">Real-world Patterns</span>
          </div>
          <div className="flex items-center gap-2">
            <CheckBadgeIcon className="w-5 h-5 text-yellow-300" />
            <span className="text-sm font-semibold">Enterprise Grade</span>
          </div>
        </div>
      </div>

      {/* Shine Effect */}
      <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/10 to-transparent transform translate-x-[-100%] animate-[shimmer_3s_ease-in-out_infinite] pointer-events-none" />

      <style>{`
        @keyframes shimmer {
          0% { transform: translateX(-100%); }
          100% { transform: translateX(100%); }
        }
      `}</style>
    </div>
  )
}