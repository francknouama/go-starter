import type { ReactNode } from 'react'
// import { designTokens } from '../../styles/design-tokens'

// Loading Spinner Component
interface SpinnerProps {
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  color?: 'primary' | 'secondary' | 'success' | 'warning' | 'error'
  className?: string
}

export function Spinner({ 
  size = 'md', 
  color = 'primary',
  className = '' 
}: SpinnerProps) {
  const sizeStyles = {
    xs: 'w-3 h-3',
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8',
    xl: 'w-12 h-12',
  }

  const colorStyles = {
    primary: 'text-primary-600',
    secondary: 'text-gray-600',
    success: 'text-success-600',
    warning: 'text-warning-600',
    error: 'text-error-600',
  }

  return (
    <svg
      className={`animate-spin ${sizeStyles[size]} ${colorStyles[color]} ${className}`}
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        className="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        strokeWidth="4"
      />
      <path
        className="opacity-75"
        fill="currentColor"
        d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
  )
}

// Pulsing Dots Loader
interface DotsLoaderProps {
  size?: 'sm' | 'md' | 'lg'
  color?: 'primary' | 'secondary' | 'success' | 'warning' | 'error'
  className?: string
}

export function DotsLoader({ 
  size = 'md', 
  color = 'primary',
  className = '' 
}: DotsLoaderProps) {
  const sizeStyles = {
    sm: 'w-1 h-1',
    md: 'w-2 h-2',
    lg: 'w-3 h-3',
  }

  const colorStyles = {
    primary: 'bg-primary-600',
    secondary: 'bg-gray-600',
    success: 'bg-success-600',
    warning: 'bg-warning-600',
    error: 'bg-error-600',
  }

  const dotClass = `${sizeStyles[size]} ${colorStyles[color]} rounded-full animate-pulse`

  return (
    <div className={`flex space-x-1 ${className}`}>
      <div className={`${dotClass} animation-delay-0`} />
      <div className={`${dotClass} animation-delay-150`} />
      <div className={`${dotClass} animation-delay-300`} />
    </div>
  )
}

// Progress Bar Component
interface ProgressBarProps {
  value: number // 0-100
  size?: 'sm' | 'md' | 'lg'
  color?: 'primary' | 'success' | 'warning' | 'error'
  showLabel?: boolean
  label?: string
  className?: string
}

export function ProgressBar({
  value,
  size = 'md',
  color = 'primary',
  showLabel = false,
  label,
  className = ''
}: ProgressBarProps) {
  const sizeStyles = {
    sm: 'h-1',
    md: 'h-2',
    lg: 'h-3',
  }

  const colorStyles = {
    primary: 'bg-primary-600',
    success: 'bg-success-600',
    warning: 'bg-warning-600',
    error: 'bg-error-600',
  }

  const clampedValue = Math.max(0, Math.min(100, value))

  return (
    <div className={className}>
      {showLabel && (
        <div className="flex justify-between items-center mb-1">
          <span className="text-sm font-medium text-gray-700">
            {label}
          </span>
          <span className="text-sm text-gray-500">
            {Math.round(clampedValue)}%
          </span>
        </div>
      )}
      <div className={`w-full bg-gray-200 rounded-full ${sizeStyles[size]}`}>
        <div
          className={`${colorStyles[color]} ${sizeStyles[size]} rounded-full transition-all duration-300 ease-out`}
          style={{ width: `${clampedValue}%` }}
        />
      </div>
    </div>
  )
}

// Circular Progress Component
interface CircularProgressProps {
  value: number // 0-100
  size?: 'sm' | 'md' | 'lg' | 'xl'
  color?: 'primary' | 'success' | 'warning' | 'error'
  showLabel?: boolean
  strokeWidth?: number
  className?: string
}

export function CircularProgress({
  value,
  size = 'md',
  color = 'primary',
  showLabel = false,
  strokeWidth = 4,
  className = ''
}: CircularProgressProps) {
  const sizeMap = {
    sm: 32,
    md: 48,
    lg: 64,
    xl: 96,
  }

  const colorStyles = {
    primary: 'text-primary-600',
    success: 'text-success-600',
    warning: 'text-warning-600',
    error: 'text-error-600',
  }

  const dimension = sizeMap[size]
  const radius = (dimension - strokeWidth) / 2
  const circumference = radius * 2 * Math.PI
  const offset = circumference - (value / 100) * circumference

  return (
    <div className={`relative inline-flex items-center justify-center ${className}`}>
      <svg
        width={dimension}
        height={dimension}
        className="transform -rotate-90"
      >
        {/* Background circle */}
        <circle
          cx={dimension / 2}
          cy={dimension / 2}
          r={radius}
          stroke="currentColor"
          strokeWidth={strokeWidth}
          fill="none"
          className="text-gray-200"
        />
        {/* Progress circle */}
        <circle
          cx={dimension / 2}
          cy={dimension / 2}
          r={radius}
          stroke="currentColor"
          strokeWidth={strokeWidth}
          fill="none"
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          strokeLinecap="round"
          className={`${colorStyles[color]} transition-all duration-300 ease-out`}
        />
      </svg>
      {showLabel && (
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="text-sm font-medium text-gray-700">
            {Math.round(value)}%
          </span>
        </div>
      )}
    </div>
  )
}

// Skeleton Loader for content loading
interface SkeletonProps {
  className?: string
  width?: string | number
  height?: string | number
  rounded?: boolean
  lines?: number
}

export function Skeleton({
  className = '',
  width = '100%',
  height = '1rem',
  rounded = false,
  lines = 1
}: SkeletonProps) {
  const skeletonClass = `
    bg-gray-200 animate-pulse
    ${rounded ? 'rounded-full' : 'rounded'}
    ${className}
  `.trim()

  if (lines === 1) {
    return (
      <div
        className={skeletonClass}
        style={{ width, height }}
      />
    )
  }

  return (
    <div className="space-y-2">
      {Array.from({ length: lines }).map((_, index) => (
        <div
          key={index}
          className={skeletonClass}
          style={{ 
            width: index === lines - 1 ? '75%' : width, 
            height 
          }}
        />
      ))}
    </div>
  )
}

// Loading Overlay for async operations
interface LoadingOverlayProps {
  isLoading: boolean
  children: ReactNode
  spinner?: boolean
  blur?: boolean
  message?: string
  className?: string
}

export function LoadingOverlay({
  isLoading,
  children,
  spinner = true,
  blur = true,
  message,
  className = ''
}: LoadingOverlayProps) {
  return (
    <div className={`relative ${className}`}>
      {children}
      {isLoading && (
        <div className={`
          absolute inset-0 flex items-center justify-center
          bg-white bg-opacity-75 z-10
          ${blur ? 'backdrop-blur-sm' : ''}
        `}>
          <div className="text-center">
            {spinner && <Spinner size="lg" className="mx-auto mb-2" />}
            {message && (
              <p className="text-sm text-gray-600 font-medium">
                {message}
              </p>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

// Skeleton Card for loading content blocks
interface SkeletonCardProps {
  showAvatar?: boolean
  lines?: number
  className?: string
}

export function SkeletonCard({
  showAvatar = false,
  lines = 3,
  className = ''
}: SkeletonCardProps) {
  return (
    <div className={`p-4 border border-gray-200 rounded-lg ${className}`}>
      <div className="animate-pulse">
        {showAvatar && (
          <div className="flex items-center space-x-3 mb-4">
            <Skeleton width={40} height={40} rounded />
            <div className="flex-1">
              <Skeleton width="60%" height="1rem" />
              <Skeleton width="40%" height="0.75rem" className="mt-1" />
            </div>
          </div>
        )}
        <div className="space-y-2">
          {Array.from({ length: lines }).map((_, index) => (
            <Skeleton
              key={index}
              width={index === lines - 1 ? '75%' : '100%'}
              height="0.875rem"
            />
          ))}
        </div>
        <div className="flex space-x-2 mt-4">
          <Skeleton width={80} height={32} />
          <Skeleton width={80} height={32} />
        </div>
      </div>
    </div>
  )
}

// Pulse Animation Component
interface PulseProps {
  children: ReactNode
  color?: 'primary' | 'success' | 'warning' | 'error'
  intensity?: 'subtle' | 'normal' | 'strong'
  className?: string
}

export function Pulse({
  children,
  color = 'primary',
  intensity = 'normal',
  className = ''
}: PulseProps) {
  const colorStyles = {
    primary: 'shadow-primary-500/25',
    success: 'shadow-success-500/25',
    warning: 'shadow-warning-500/25',
    error: 'shadow-error-500/25',
  }

  const intensityStyles = {
    subtle: 'animate-pulse-gentle',
    normal: 'animate-pulse',
    strong: 'animate-ping',
  }

  return (
    <div className={`
      ${intensityStyles[intensity]}
      ${colorStyles[color]}
      ${className}
    `}>
      {children}
    </div>
  )
}

// LoadingStates Component - Main loading state indicator
interface LoadingStatesProps {
  message?: string
  description?: string
  type?: 'spinner' | 'dots' | 'progress'
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

export function LoadingStates({
  message = 'Loading...',
  description,
  type = 'spinner',
  size = 'md',
  className = ''
}: LoadingStatesProps) {
  const renderLoader = () => {
    switch (type) {
      case 'dots':
        return <DotsLoader size={size} className="mb-4" />
      case 'progress':
        return <ProgressBar value={50} size={size} className="mb-4" />
      default:
        return <Spinner size={size === 'sm' ? 'md' : 'lg'} className="mb-4" />
    }
  }

  return (
    <div className={`flex flex-col items-center justify-center text-center ${className}`}>
      {renderLoader()}
      {message && (
        <h3 className="text-lg font-medium text-gray-900 mb-2">
          {message}
        </h3>
      )}
      {description && (
        <p className="text-sm text-gray-600 max-w-sm">
          {description}
        </p>
      )}
    </div>
  )
}