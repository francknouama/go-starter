import { CheckCircleIcon, ExclamationCircleIcon } from '@heroicons/react/20/solid'

interface FormValidationSummaryProps {
  isValid: boolean
  projectName: string
  moduleUrl: string
}

export default function FormValidationSummary({ 
  isValid, 
  projectName, 
  moduleUrl 
}: FormValidationSummaryProps) {
  if (!projectName && !moduleUrl) {
    return null // Don't show summary when form is empty
  }

  return (
    <div className={`
      rounded-lg p-3 mb-4 border transition-all duration-200
      ${isValid 
        ? 'bg-green-50 border-green-200' 
        : 'bg-yellow-50 border-yellow-200'
      }
    `}>
      <div className="flex items-start space-x-2">
        {isValid ? (
          <CheckCircleIcon className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
        ) : (
          <ExclamationCircleIcon className="w-5 h-5 text-yellow-600 flex-shrink-0 mt-0.5" />
        )}
        <div className="flex-1">
          <p className={`text-sm font-medium ${isValid ? 'text-green-800' : 'text-yellow-800'}`}>
            {isValid ? 'Ready to generate!' : 'Almost ready...'}
          </p>
          <p className={`text-xs mt-1 ${isValid ? 'text-green-700' : 'text-yellow-700'}`}>
            {isValid 
              ? `Your project "${projectName}" will be created with module "${moduleUrl}"`
              : 'Please complete all required fields with valid values'
            }
          </p>
        </div>
      </div>
    </div>
  )
}