import { useState, useEffect } from 'react'
import { ExclamationCircleIcon, CheckCircleIcon } from '@heroicons/react/20/solid'

interface ValidationRule {
  required?: boolean
  pattern?: RegExp
  minLength?: number
  maxLength?: number
  custom?: (value: string) => string | null // Returns error message or null if valid
}

interface ValidatedInputProps {
  label: string
  value: string
  onChange: (value: string) => void
  placeholder?: string
  helpText?: string
  validation?: ValidationRule
  validateOnBlur?: boolean
  showSuccessState?: boolean
  autoComplete?: string
  type?: 'text' | 'email' | 'url'
}

export default function ValidatedInput({
  label,
  value,
  onChange,
  placeholder,
  helpText,
  validation,
  validateOnBlur = true,
  showSuccessState = true,
  autoComplete,
  type = 'text'
}: ValidatedInputProps) {
  const [error, setError] = useState<string | null>(null)
  const [touched, setTouched] = useState(false)
  const [isValid, setIsValid] = useState(false)

  // Validate the input
  const validate = (inputValue: string): string | null => {
    if (!validation) return null

    // Required validation
    if (validation.required && !inputValue.trim()) {
      return 'This field is required'
    }

    // Pattern validation
    if (validation.pattern && inputValue && !validation.pattern.test(inputValue)) {
      return 'Invalid format'
    }

    // Min length validation
    if (validation.minLength && inputValue.length < validation.minLength) {
      return `Must be at least ${validation.minLength} characters`
    }

    // Max length validation
    if (validation.maxLength && inputValue.length > validation.maxLength) {
      return `Must be no more than ${validation.maxLength} characters`
    }

    // Custom validation
    if (validation.custom) {
      return validation.custom(inputValue)
    }

    return null
  }

  // Validate on value change
  useEffect(() => {
    if (touched || value) {
      const validationError = validate(value)
      setError(validationError)
      setIsValid(!validationError && value.length > 0)
    }
  }, [value, touched])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e.target.value)
  }

  const handleBlur = () => {
    setTouched(true)
  }

  const getInputClasses = () => {
    const baseClasses = 'w-full px-3 py-2 border rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2'
    
    if (error && touched) {
      return `${baseClasses} border-red-300 text-red-900 placeholder-red-300 focus:ring-red-500 focus:border-red-500`
    }
    
    if (isValid && showSuccessState && touched) {
      return `${baseClasses} border-green-300 text-green-900 placeholder-green-300 focus:ring-green-500 focus:border-green-500`
    }
    
    return `${baseClasses} border-gray-300 focus:ring-primary-500 focus:border-primary-500`
  }

  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-2">
        {label}
        {validation?.required && (
          <span className="text-red-500 ml-1" aria-label="required">*</span>
        )}
      </label>
      
      <div className="relative">
        <input
          type={type}
          value={value}
          onChange={handleChange}
          onBlur={validateOnBlur ? handleBlur : undefined}
          placeholder={placeholder}
          autoComplete={autoComplete}
          className={getInputClasses()}
          aria-invalid={error && touched ? 'true' : 'false'}
          aria-describedby={error && touched ? `${label}-error` : helpText ? `${label}-help` : undefined}
        />
        
        {/* Validation icons */}
        {touched && (
          <>
            {error && (
              <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                <ExclamationCircleIcon className="h-5 w-5 text-red-500" />
              </div>
            )}
            {isValid && showSuccessState && (
              <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                <CheckCircleIcon className="h-5 w-5 text-green-500" />
              </div>
            )}
          </>
        )}
      </div>
      
      {/* Error message with live region */}
      {error && touched && (
        <div 
          className="mt-1 text-sm text-red-600 flex items-center" 
          id={`${label}-error`}
          role="alert"
          aria-live="polite"
        >
          <ExclamationCircleIcon className="h-4 w-4 mr-1 flex-shrink-0" aria-hidden="true" />
          {error}
        </div>
      )}
      
      {/* Success message with live region */}
      {isValid && showSuccessState && touched && !error && (
        <div 
          className="mt-1 text-sm text-green-600 flex items-center"
          role="status"
          aria-live="polite"
        >
          <CheckCircleIcon className="h-4 w-4 mr-1 flex-shrink-0" aria-hidden="true" />
          Valid {label.toLowerCase()}
        </div>
      )}
      
      {/* Help text */}
      {helpText && !error && !isValid && (
        <p className="mt-1 text-sm text-gray-500" id={`${label}-help`}>
          {helpText}
        </p>
      )}
    </div>
  )
}