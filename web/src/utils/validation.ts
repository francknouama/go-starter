// Common validation patterns
export const ValidationPatterns = {
  // Project name: letters, numbers, hyphens, underscores
  projectName: /^[a-zA-Z0-9-_]+$/,
  
  // Go module path: domain/user/project format
  goModule: /^[a-zA-Z0-9.-]+\/[a-zA-Z0-9._\/-]+$/,
  
  // Semantic version
  version: /^\d+\.\d+(\.\d+)?$/,
  
  // Email
  email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  
  // URL
  url: /^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$/,
}

// Validation error messages
export const ValidationMessages = {
  projectName: {
    required: 'Project name is required',
    pattern: 'Use only letters, numbers, hyphens, and underscores',
    minLength: 'Project name must be at least 2 characters',
    maxLength: 'Project name must be less than 50 characters',
  },
  goModule: {
    required: 'Module path is required',
    pattern: 'Must be a valid Go module path (e.g., github.com/user/project)',
    example: 'Example: github.com/username/project-name',
  },
  version: {
    pattern: 'Must be a valid version (e.g., 1.21)',
  },
}

// Custom validation functions
export const CustomValidators = {
  projectName: (value: string): string | null => {
    if (!value) return ValidationMessages.projectName.required
    if (!ValidationPatterns.projectName.test(value)) return ValidationMessages.projectName.pattern
    if (value.length < 2) return ValidationMessages.projectName.minLength
    if (value.length > 50) return ValidationMessages.projectName.maxLength
    return null
  },
  
  goModule: (value: string): string | null => {
    if (!value) return ValidationMessages.goModule.required
    if (!ValidationPatterns.goModule.test(value)) {
      return `${ValidationMessages.goModule.pattern}. ${ValidationMessages.goModule.example}`
    }
    return null
  },
  
  uniqueName: (existingNames: string[]) => (value: string): string | null => {
    if (existingNames.includes(value)) {
      return 'This name is already taken'
    }
    return null
  },
}