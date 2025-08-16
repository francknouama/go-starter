// Help content for various concepts in the go-starter UI

export const HelpContent = {
  // Project Types
  projectTypes: {
    'web-api': 'A RESTful API server for building backend services. Includes HTTP routing, middleware, and JSON handling.',
    'cli': 'Command-line application with argument parsing, subcommands, and terminal output formatting.',
    'library': 'Reusable Go package that can be imported by other projects. Focuses on public API design.',
    'lambda': 'AWS Lambda serverless function optimized for event processing and cloud deployment.',
    'microservice': 'Small, independently deployable service designed for distributed systems and container orchestration.',
  },

  // Architectures
  architectures: {
    'standard': 'Simple, straightforward structure perfect for most projects. Easy to understand and maintain.',
    'clean': 'Clean Architecture separates business logic from external dependencies using layers and interfaces.',
    'ddd': 'Domain-Driven Design focuses on modeling complex business domains with rich domain models.',
    'hexagonal': 'Hexagonal Architecture (Ports & Adapters) isolates core logic from external systems for maximum flexibility.',
  },

  // Frameworks
  frameworks: {
    'gin': 'Fast and lightweight HTTP web framework with excellent performance and minimal memory footprint.',
    'echo': 'High performance, extensible web framework with robust middleware support.',
    'fiber': 'Express-inspired web framework built on Fasthttp for extreme performance.',
    'chi': 'Lightweight, idiomatic HTTP router that stays close to standard library patterns.',
    'cobra': 'Industry-standard CLI framework used by Kubernetes, Docker, and other major projects.',
  },

  // Loggers
  loggers: {
    'slog': 'Go\'s standard structured logging library. Simple, performant, and officially supported.',
    'zap': 'Uber\'s blazing fast, structured logger with zero allocations in hot paths.',
    'logrus': 'Popular structured logger with extensive features and community support.',
    'zerolog': 'Zero-allocation JSON logger focused on performance and efficiency.',
  },

  // Database Drivers
  databases: {
    'postgres': 'PostgreSQL - Advanced open-source relational database with excellent Go support.',
    'mysql': 'MySQL/MariaDB - Popular relational database for web applications.',
    'mongodb': 'MongoDB - Document-oriented NoSQL database for flexible schemas.',
    'sqlite': 'SQLite - Embedded database perfect for development and small applications.',
    'redis': 'Redis - In-memory data structure store used for caching and real-time features.',
  },

  // Authentication Types
  authentication: {
    'jwt': 'JSON Web Tokens - Stateless authentication tokens ideal for APIs and microservices.',
    'oauth2': 'OAuth 2.0 - Industry standard for authorization, enables "Login with Google/GitHub".',
    'session': 'Session-based auth using cookies, traditional approach for web applications.',
    'api-key': 'API Key authentication for machine-to-machine communication.',
  },

  // General Concepts
  concepts: {
    'module-path': 'Go module path uniquely identifies your project. Format: domain.com/user/project',
    'go-version': 'The Go language version your project will use. Higher versions have more features.',
    'deployment': 'How and where your application will run in production.',
    'orm': 'Object-Relational Mapping - Libraries that simplify database interactions.',
    'middleware': 'Code that runs before/after HTTP handlers for cross-cutting concerns.',
    'progressive-disclosure': 'UI pattern that reveals complexity gradually as users need it.',
  },

  // Quick Tips
  tips: {
    'project-name': 'Use lowercase letters, numbers, and hyphens. This becomes your directory name.',
    'module-naming': 'Follow Go conventions: github.com/username/projectname',
    'version-selection': 'Use the latest stable Go version unless you have specific requirements.',
    'architecture-choice': 'Start with Standard architecture and evolve as your project grows.',
    'framework-selection': 'Gin is a safe default. Consider Chi for standard library compatibility.',
  },

  // Progressive Disclosure Help
  disclosure: {
    'basic-mode': 'Basic mode shows only essential options, perfect for beginners and quick projects. You can always switch to Advanced mode later.',
    'advanced-mode': 'Advanced mode provides full customization including database, authentication, and deployment options. Recommended for complex projects.',
    'mode-switching': 'You can switch between Basic and Advanced modes at any time. Your current configuration will be preserved.',
  },

  // Keyboard Shortcuts
  shortcuts: {
    'question-mark': 'Press ? to show keyboard shortcuts overlay',
    'escape': 'Press Escape to close any open modals or overlays',
    'alt-h': 'Press Alt+H to open the quick start guide',
    'alt-t': 'Press Alt+T to toggle help tooltips on all form fields',
    'numbers': 'Press 1-5 to quickly select project types (CLI, Web API, Library, Lambda, Microservice)',
  },

  // Accessibility
  accessibility: {
    'keyboard-navigation': 'All interactive elements can be accessed using keyboard navigation with Tab and Enter keys.',
    'screen-readers': 'Screen reader support with proper ARIA labels and announcements for state changes.',
    'high-contrast': 'Interface adapts to high contrast mode preferences for better visibility.',
    'focus-indicators': 'Clear focus indicators on all interactive elements for keyboard navigation.',
  },

  // Error Messages
  errors: {
    'invalid-project-name': 'Project name must contain only lowercase letters, numbers, and hyphens.',
    'invalid-module-path': 'Module path must be a valid Go module identifier (e.g., github.com/user/project).',
    'required-field': 'This field is required to generate your project.',
    'connection-error': 'Unable to connect to the generation service. Please check your internet connection.',
  }
}

// Helper function to get help content
export function getHelpContent(category: keyof typeof HelpContent, key: string): string {
  const categoryContent = HelpContent[category]
  if (categoryContent && key in categoryContent) {
    return categoryContent[key as keyof typeof categoryContent]
  }
  return 'No help available for this item.'
}