import { forwardRef } from 'react'
import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { accessibility } from '../../styles/design-tokens'

export type ButtonVariant = 
  | 'primary' 
  | 'secondary' 
  | 'outline' 
  | 'ghost' 
  | 'success' 
  | 'warning' 
  | 'error'

export type ButtonSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant
  size?: ButtonSize
  loading?: boolean
  leftIcon?: ReactNode
  rightIcon?: ReactNode
  fullWidth?: boolean
  children: ReactNode
  // Accessibility props
  'aria-label'?: string
  'aria-describedby'?: string
  'aria-pressed'?: boolean
  loadingText?: string
  // High contrast mode support
  highContrast?: boolean
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(({
  variant = 'primary',
  size = 'md',
  loading = false,
  leftIcon,
  rightIcon,
  fullWidth = false,
  disabled,
  children,
  className = '',
  loadingText = 'Loading...',
  highContrast = false,
  ...props
}, ref) => {
  // Check for user preferences
  const shouldReduceMotion = accessibility.shouldReduceMotion()
  const prefersHighContrast = accessibility.prefersHighContrast() || highContrast

  // Base button styles using design tokens with accessibility considerations
  const baseStyles = [
    'inline-flex items-center justify-center',
    'font-medium',
    // Conditional transitions based on motion preference
    shouldReduceMotion ? '' : 'transition-all duration-200',
    'border rounded-md',
    // Enhanced focus styles for accessibility
    'focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
    'disabled:opacity-50 disabled:cursor-not-allowed',
    // Conditional animations
    shouldReduceMotion ? '' : 'active:scale-[0.98]',
    fullWidth ? 'w-full' : '',
    // Ensure minimum touch target size (44px for WCAG 2.5.5)
    size === 'xs' ? 'min-h-[44px]' : '',
  ].filter(Boolean).join(' ')

  // Size styles with minimum touch target compliance
  const sizeStyles = {
    xs: 'h-11 px-2 text-xs', // Increased from h-6 to meet 44px minimum
    sm: 'h-11 px-3 text-sm', // Increased from h-8 to meet 44px minimum
    md: 'h-10 px-4 text-sm',
    lg: 'h-12 px-6 text-base',
    xl: 'h-14 px-8 text-lg',
  }

  // Variant styles with high contrast support
  const getVariantStyles = () => {
    if (prefersHighContrast) {
      // High contrast mode uses simpler Tailwind classes without template literals
      return {
        primary: [
          'bg-black border-black text-white',
          'hover:bg-gray-900 hover:border-gray-900',
          'focus-visible:ring-black',
        ].join(' '),
        secondary: [
          'bg-white border-black text-black',
          'hover:bg-gray-100 hover:border-black',
          'focus-visible:ring-black',
        ].join(' '),
        outline: [
          'bg-transparent border-black text-black',
          'hover:bg-gray-100 hover:border-black',
          'focus-visible:ring-black',
        ].join(' '),
        ghost: [
          'bg-transparent border-transparent text-black',
          'hover:bg-gray-100',
          'focus-visible:ring-black',
        ].join(' '),
        success: [
          'bg-black border-black text-white',
          'hover:bg-gray-900 hover:border-gray-900',
          'focus-visible:ring-black',
        ].join(' '),
        warning: [
          'bg-black border-black text-white',
          'hover:bg-gray-900 hover:border-gray-900',
          'focus-visible:ring-black',
        ].join(' '),
        error: [
          'bg-black border-black text-white',
          'hover:bg-gray-900 hover:border-gray-900',
          'focus-visible:ring-black',
        ].join(' '),
      }
    }

    return {
      primary: [
        'bg-primary-600 border-primary-600 text-white',
        'hover:bg-primary-700 hover:border-primary-700',
        'focus-visible:ring-primary-500',
        'active:bg-primary-800',
      ].join(' '),
    
      secondary: [
        'bg-gray-100 border-gray-300 text-gray-900',
        'hover:bg-gray-200 hover:border-gray-400',
        'focus-visible:ring-gray-500',
        'active:bg-gray-300',
      ].join(' '),
      
      outline: [
        'bg-transparent border-gray-300 text-gray-700',
        'hover:bg-gray-50 hover:border-gray-400',
        'focus-visible:ring-primary-500',
        'active:bg-gray-100',
      ].join(' '),
      
      ghost: [
        'bg-transparent border-transparent text-gray-700',
        'hover:bg-gray-100',
        'focus-visible:ring-primary-500',
        'active:bg-gray-200',
      ].join(' '),
      
      success: [
        'bg-success-600 border-success-600 text-white',
        'hover:bg-success-700 hover:border-success-700',
        'focus-visible:ring-success-500',
        'active:bg-success-800',
      ].join(' '),
      
      warning: [
        'bg-warning-600 border-warning-600 text-white',
        'hover:bg-warning-700 hover:border-warning-700',
        'focus-visible:ring-warning-500',
        'active:bg-warning-800',
      ].join(' '),
      
      error: [
        'bg-error-600 border-error-600 text-white',
        'hover:bg-error-700 hover:border-error-700',
        'focus-visible:ring-error-500',
        'active:bg-error-800',
      ].join(' '),
    }
  }

  const variantStyles = getVariantStyles()

  // Loading spinner component with accessibility
  const LoadingSpinner = ({ size }: { size: ButtonSize }) => {
    const spinnerSize = {
      xs: 'w-3 h-3',
      sm: 'w-4 h-4',
      md: 'w-4 h-4',
      lg: 'w-5 h-5',
      xl: 'w-6 h-6',
    }

    return (
      <>
        <svg
          className={`${shouldReduceMotion ? '' : 'animate-spin'} ${spinnerSize[size]}`}
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          role="img"
          aria-label="Loading"
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
        <span className="sr-only">{loadingText}</span>
      </>
    )
  }

  // Icon spacing
  const iconSpacing = {
    xs: 'gap-1',
    sm: 'gap-1.5',
    md: 'gap-2',
    lg: 'gap-2',
    xl: 'gap-3',
  }

  // Combined styles
  const combinedClassName = [
    baseStyles,
    sizeStyles[size],
    variantStyles[variant],
    iconSpacing[size],
    className,
  ].filter(Boolean).join(' ')

  // Handle keyboard interactions
  const handleKeyDown = (event: React.KeyboardEvent<HTMLButtonElement>) => {
    // Allow space key to activate button (Enter is handled by default)
    if (event.key === ' ') {
      event.preventDefault()
      if (!disabled && !loading) {
        event.currentTarget.click()
      }
    }
    
    // Call original onKeyDown if provided
    props.onKeyDown?.(event)
  }

  return (
    <button
      ref={ref}
      disabled={disabled || loading}
      className={combinedClassName}
      onKeyDown={handleKeyDown}
      aria-disabled={disabled || loading}
      {...props}
    >
      {loading ? (
        <>
          <LoadingSpinner size={size} />
          <span>{children}</span>
        </>
      ) : (
        <>
          {leftIcon && <span className="flex-shrink-0" aria-hidden="true">{leftIcon}</span>}
          <span>{children}</span>
          {rightIcon && <span className="flex-shrink-0" aria-hidden="true">{rightIcon}</span>}
        </>
      )}
    </button>
  )
})

Button.displayName = 'Button'

export default Button

// Button Group Component for related actions
interface ButtonGroupProps {
  children: ReactNode
  className?: string
  spacing?: 'none' | 'sm' | 'md'
  orientation?: 'horizontal' | 'vertical'
}

export function ButtonGroup({ 
  children, 
  className = '',
  spacing = 'sm',
  orientation = 'horizontal'
}: ButtonGroupProps) {
  const spacingStyles = {
    none: '',
    sm: orientation === 'horizontal' ? 'space-x-2' : 'space-y-2',
    md: orientation === 'horizontal' ? 'space-x-4' : 'space-y-4',
  }

  const orientationStyles = {
    horizontal: 'flex flex-row',
    vertical: 'flex flex-col',
  }

  const combinedClassName = [
    orientationStyles[orientation],
    spacingStyles[spacing],
    className,
  ].filter(Boolean).join(' ')

  return (
    <div className={combinedClassName}>
      {children}
    </div>
  )
}

// Icon Button for minimal actions
interface IconButtonProps extends Omit<ButtonProps, 'leftIcon' | 'rightIcon' | 'children'> {
  icon: ReactNode
  'aria-label': string
}

export function IconButton({ 
  icon, 
  variant = 'ghost',
  size = 'md',
  className = '',
  ...props 
}: IconButtonProps) {
  // Ensure minimum touch target size for accessibility
  const sizeStyles = {
    xs: 'w-11 h-11 p-1',    // Increased to meet 44px minimum
    sm: 'w-11 h-11 p-1.5',  // Increased to meet 44px minimum
    md: 'w-10 h-10 p-2',
    lg: 'w-12 h-12 p-2.5',
    xl: 'w-14 h-14 p-3',
  }

  return (
    <Button
      variant={variant}
      size={size}
      className={`${sizeStyles[size]} ${className}`}
      {...props}
    >
      <span aria-hidden="true">{icon}</span>
    </Button>
  )
}

// Floating Action Button for primary actions
interface FloatingActionButtonProps extends Omit<ButtonProps, 'variant' | 'size'> {
  size?: 'md' | 'lg'
  position?: 'bottom-right' | 'bottom-left' | 'top-right' | 'top-left'
}

export function FloatingActionButton({
  size = 'lg',
  position = 'bottom-right',
  className = '',
  children,
  ...props
}: FloatingActionButtonProps) {
  const positionStyles = {
    'bottom-right': 'fixed bottom-6 right-6',
    'bottom-left': 'fixed bottom-6 left-6',
    'top-right': 'fixed top-6 right-6',
    'top-left': 'fixed top-6 left-6',
  }

  const sizeStyles = {
    md: 'w-12 h-12',
    lg: 'w-14 h-14',
  }

  const combinedClassName = [
    positionStyles[position],
    sizeStyles[size],
    'rounded-full shadow-lg hover:shadow-xl',
    'transition-all duration-200',
    'z-50',
    className,
  ].filter(Boolean).join(' ')

  return (
    <Button
      variant="primary"
      className={combinedClassName}
      {...props}
    >
      {children}
    </Button>
  )
}