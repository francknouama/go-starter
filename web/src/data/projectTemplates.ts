import type { ProjectConfig, ProjectType, Architecture, Framework, LoggerType, DatabaseDriver, DatabaseORM, AuthType } from '../types'

export interface ProjectTemplate {
  id: string
  name: string
  description: string
  category: 'api' | 'cli' | 'microservice' | 'fullstack' | 'library' | 'serverless'
  complexity: 'beginner' | 'intermediate' | 'advanced' | 'expert'
  
  // Visual representation
  icon: string // emoji or icon identifier
  previewImage?: string
  color: string // hex color for the card
  
  // Template metadata
  tags: string[]
  useCase: string
  techStack: string[]
  
  // Auto-configuration
  config: Partial<ProjectConfig>
  
  // Architecture info for preview
  architecture: {
    diagram: string // ASCII art or description
    components: string[]
    patterns: string[]
  }
  
  // Getting started guide
  quickStart: {
    commands: string[]
    nextSteps: string[]
    learnMore: string
  }
  
  // Statistics
  popularity: number // 1-10 scale
  estimatedSetupTime: string
  recommendedFor: string[]
}

// 🎉 HISTORIC ACHIEVEMENT: 12 Production-Ready Blueprints - 100% Coverage! 🎉
export const PROJECT_TEMPLATES: ProjectTemplate[] = [
  // CLI Templates - Production Ready ✅
  {
    id: 'cli-simple',
    name: 'Simple CLI',
    description: 'Lightweight CLI tool with basic commands, perfect for quick utilities and learning Go CLI development',
    category: 'cli',
    complexity: 'beginner',
    icon: '⚡',
    color: '#10B981',
    tags: ['CLI', 'Cobra', 'Simple', 'Learning', 'Utilities'],
    useCase: 'Quick command-line utilities, prototypes, and learning Go CLI development',
    techStack: ['Cobra', 'Slog', 'Basic Testing', 'Makefile'],
    config: {
      projectType: 'cli',
      framework: 'cobra',
      architecture: 'standard',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Main        │────│    Commands     │────│    Output       │
│   Entry Point   │    │   (cmd/)        │    │   (stdout)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['Root Command', 'Subcommands', 'Flag Handling', 'Basic Logging'],
      patterns: ['Command Pattern', 'Simple Architecture']
    },
    quickStart: {
      commands: [
        'go build -o mytool cmd/main.go',
        './mytool --help',
        './mytool version'
      ],
      nextSteps: [
        'Add your custom commands',
        'Configure command flags and options',
        'Add input validation',
        'Extend with additional subcommands'
      ],
      learnMore: 'Simple CLI with 8 files - perfect for learning'
    },
    popularity: 7,
    estimatedSetupTime: '2-5 minutes',
    recommendedFor: ['Go beginners', 'Quick utilities', 'Learning projects']
  },
  
  {
    id: 'cli-standard',
    name: 'Standard CLI',
    description: 'Production-ready CLI application with comprehensive testing, configuration, and deployment features',
    category: 'cli',
    complexity: 'intermediate',
    icon: '💻',
    color: '#3B82F6',
    tags: ['CLI', 'Production', 'Testing', 'CI/CD', 'Configuration'],
    useCase: 'Professional command-line tools and production CLI applications',
    techStack: ['Cobra', 'Viper Config', 'Comprehensive Tests', 'GitHub Actions'],
    config: {
      projectType: 'cli',
      framework: 'cobra',
      architecture: 'standard',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Config      │────│    Commands     │────│     Logic       │
│   Management    │    │   (Cobra)       │    │   (Business)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│     Tests       │    │     CI/CD       │    │    Logging      │
│  (Unit & Int)   │    │   (Actions)     │    │    (Slog)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['Command Structure', 'Configuration', 'Business Logic', 'Testing Suite', 'CI/CD Pipeline'],
      patterns: ['Command Pattern', 'Configuration Pattern', 'Layered Architecture']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'go build -o mytool cmd/main.go',
        'go test ./...',
        './mytool --help'
      ],
      nextSteps: [
        'Configure application settings',
        'Add comprehensive tests',
        'Set up CI/CD pipeline',
        'Deploy to production'
      ],
      learnMore: 'Full-featured CLI with 29 files - production ready'
    },
    popularity: 9,
    estimatedSetupTime: '10-15 minutes',
    recommendedFor: ['Production CLIs', 'DevOps tools', 'Professional development']
  },

  // Web API Templates - Production Ready ✅
  {
    id: 'web-api-standard',
    name: 'Standard Web API',
    description: 'Production-ready REST API with middleware, database integration, and comprehensive testing',
    category: 'api',
    complexity: 'intermediate',
    icon: '🌐',
    color: '#8B5CF6',
    tags: ['REST', 'API', 'Database', 'Middleware', 'Testing'],
    useCase: 'Building standard REST APIs and web services',
    techStack: ['Gin', 'PostgreSQL', 'GORM', 'Middleware'],
    config: {
      projectType: 'web-api',
      framework: 'gin',
      architecture: 'standard',
      logger: 'slog',
      features: {
        database: {
          driver: 'postgres',
          orm: 'gorm'
        }
      }
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Router   │────│   Controllers   │────│   Services      │
│    (Gin)        │    │                 │    │   (Business)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Middleware    │    │   Database      │    │   Models        │
│  (Auth, CORS)   │    │  (PostgreSQL)   │    │   (GORM)        │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['HTTP Router', 'Controllers', 'Services', 'Database', 'Middleware'],
      patterns: ['MVC Pattern', 'Repository Pattern', 'Middleware Pattern']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'docker-compose up -d postgres',
        'go run cmd/migrate/main.go',
        'go run cmd/server/main.go'
      ],
      nextSteps: [
        'Configure database connection',
        'Add your API endpoints',
        'Set up authentication',
        'Deploy to production'
      ],
      learnMore: 'Standard REST API with ~25 files'
    },
    popularity: 10,
    estimatedSetupTime: '15-20 minutes',
    recommendedFor: ['REST APIs', 'Web services', 'CRUD applications']
  },

  {
    id: 'web-api-clean',
    name: 'Clean Architecture API',
    description: 'Enterprise-grade API following clean architecture principles with dependency injection and domain modeling',
    category: 'api',
    complexity: 'advanced',
    icon: '🏗️',
    color: '#EF4444',
    tags: ['Clean Architecture', 'DDD', 'Enterprise', 'Scalable', 'Maintainable'],
    useCase: 'Enterprise applications requiring clean architecture and domain-driven design',
    techStack: ['Gin', 'Clean Architecture', 'Dependency Injection', 'Domain Models'],
    config: {
      projectType: 'web-api',
      framework: 'gin',
      architecture: 'clean',
      logger: 'slog',
      features: {
        database: {
          driver: 'postgres',
          orm: 'gorm'
        }
      }
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frameworks    │────│   Controllers   │────│   Use Cases     │
│  (HTTP/DB)      │    │   (Handlers)    │    │  (Business)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐              │              ┌─────────────────┐
│   External      │              │              │    Domain       │
│   Interfaces    │              │              │    Entities     │
└─────────────────┘              │              └─────────────────┘
                                  │                       │
                         ┌─────────────────┐    ┌─────────────────┐
                         │  Repositories   │────│   Data Models   │
                         │  (Interfaces)   │    │   (Database)    │
                         └─────────────────┘    └─────────────────┘`,
      components: ['Domain Layer', 'Use Cases', 'Controllers', 'Repositories', 'External Interfaces'],
      patterns: ['Clean Architecture', 'Dependency Inversion', 'Domain-Driven Design', 'SOLID Principles']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'docker-compose up -d postgres',
        'go run cmd/migrate/main.go',
        'go run cmd/server/main.go'
      ],
      nextSteps: [
        'Study the clean architecture layers',
        'Implement domain entities and use cases',
        'Add repository implementations',
        'Configure dependency injection'
      ],
      learnMore: 'Clean Architecture API with ~40 files'
    },
    popularity: 8,
    estimatedSetupTime: '25-30 minutes',
    recommendedFor: ['Enterprise applications', 'Complex domains', 'Team projects']
  },

  {
    id: 'web-api-echo',
    name: 'Echo Web API',
    description: 'High-performance REST API built with Echo framework, optimized for speed and scalability',
    category: 'api',
    complexity: 'intermediate', 
    icon: '🚀',
    color: '#F59E0B',
    tags: ['Echo', 'Performance', 'REST', 'Fast', 'Scalable'],
    useCase: 'High-performance web APIs and services requiring optimal throughput',
    techStack: ['Echo', 'High Performance', 'Middleware', 'JSON Binding'],
    config: {
      projectType: 'web-api',
      framework: 'echo',
      architecture: 'standard',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Echo Router   │────│   Handlers      │────│   Services      │
│  (HTTP/HTTPS)   │    │  (Controllers)  │    │  (Business)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Middleware    │    │   Validation    │    │    Database     │
│ (Performance)   │    │   (Binding)     │    │  (Optional)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['Echo Router', 'Handlers', 'Services', 'Middleware', 'Validation'],
      patterns: ['MVC Pattern', 'Middleware Chain', 'Fast Routing']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'go run cmd/server/main.go',
        'curl http://localhost:8080/health'
      ],
      nextSteps: [
        'Add your API routes',
        'Configure middleware',
        'Implement handlers',
        'Add validation and testing'
      ],
      learnMore: 'Echo API with ~25 files - optimized for performance'
    },
    popularity: 8,
    estimatedSetupTime: '10-15 minutes',
    recommendedFor: ['Performance-critical APIs', 'High-throughput services', 'Microservices']
  },

  {
    id: 'web-api-fiber',
    name: 'Fiber Web API',
    description: 'Ultra-fast REST API using Fiber framework inspired by Express.js, perfect for rapid development',
    category: 'api',
    complexity: 'intermediate',
    icon: '⚡',
    color: '#06B6D4',
    tags: ['Fiber', 'Fast', 'Express-like', 'Modern', 'Developer-friendly'],
    useCase: 'Rapid API development with Express.js-like syntax and ultra-fast performance',
    techStack: ['Fiber', 'Ultra Fast', 'Express-like', 'Modern Go'],
    config: {
      projectType: 'web-api',
      framework: 'fiber',
      architecture: 'standard',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Fiber App     │────│   Routes        │────│   Handlers      │
│  (Express-like) │    │  (REST/JSON)    │    │  (Controllers)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Middleware    │    │   JSON/XML      │    │    Models       │
│ (Built-in)      │    │   Parsing       │    │  (Data)         │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['Fiber App', 'Routes', 'Handlers', 'Middleware', 'JSON Parsing'],
      patterns: ['Express Pattern', 'Middleware Stack', 'Fast Routing']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'go run cmd/server/main.go',
        'curl http://localhost:3000/api/health'
      ],
      nextSteps: [
        'Add your API endpoints',
        'Configure Fiber middleware',
        'Implement business logic',
        'Add database integration'
      ],
      learnMore: 'Fiber API with ~25 files - Express.js-like syntax'
    },
    popularity: 9,
    estimatedSetupTime: '10-15 minutes',
    recommendedFor: ['Rapid prototyping', 'Express.js developers', 'Ultra-fast APIs']
  },

  // Serverless Templates - Production Ready ✅
  {
    id: 'lambda-standard',
    name: 'Lambda Function',
    description: 'AWS Lambda function with event handling, logging, and comprehensive error handling',
    category: 'serverless',
    complexity: 'beginner',
    icon: '⚡',
    color: '#FF6B35',
    tags: ['AWS Lambda', 'Serverless', 'Event-driven', 'Cost-effective'],
    useCase: 'Serverless functions for event processing and lightweight APIs',
    techStack: ['AWS Lambda', 'AWS SDK', 'Event Processing', 'Cost Optimization'],
    config: {
      projectType: 'lambda-standard',
      framework: 'lambda',
      architecture: 'standard',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Event Source  │────│   Lambda        │────│   Response      │
│  (API Gateway)  │    │   Function      │    │   (JSON)        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CloudWatch    │    │   Error         │    │   AWS Services  │
│   (Logging)     │    │   Handling      │    │  (S3, DDB, etc) │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['Lambda Handler', 'Event Processing', 'Error Handling', 'AWS Integration'],
      patterns: ['Event-driven', 'Serverless', 'Function as a Service']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'go build -o main main.go',
        'sam local start-api',
        'sam deploy --guided'
      ],
      nextSteps: [
        'Configure AWS credentials',
        'Set up event sources',
        'Add business logic',
        'Deploy to AWS'
      ],
      learnMore: 'Lambda function with ~15 files - serverless ready'
    },
    popularity: 8,
    estimatedSetupTime: '10-15 minutes',
    recommendedFor: ['Event processing', 'Cost-effective APIs', 'Serverless architecture']
  },

  {
    id: 'lambda-proxy',
    name: 'Lambda API Proxy',
    description: 'API Gateway Lambda proxy integration with routing, middleware, and REST API patterns',
    category: 'serverless',
    complexity: 'intermediate',
    icon: '🌐',
    color: '#8B5CF6',
    tags: ['API Gateway', 'Lambda Proxy', 'REST', 'Serverless API'],
    useCase: 'Serverless REST APIs with API Gateway integration and Lambda proxy',
    techStack: ['API Gateway', 'Lambda Proxy', 'HTTP Routing', 'Serverless Framework'],
    config: {
      projectType: 'lambda-proxy',
      framework: 'lambda',
      architecture: 'standard',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  API Gateway    │────│  Lambda Proxy   │────│   Handlers      │
│  (HTTP Routes)  │    │  (Router)       │    │  (Business)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CORS/Auth     │    │   Request       │    │   Response      │
│  (Middleware)   │    │   Processing    │    │   Formatting    │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['API Gateway', 'Lambda Proxy', 'HTTP Routing', 'Middleware', 'Response Handling'],
      patterns: ['API Gateway Pattern', 'Lambda Proxy Integration', 'REST API']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'sam build',
        'sam local start-api',
        'sam deploy --guided'
      ],
      nextSteps: [
        'Configure API Gateway',
        'Add REST endpoints',
        'Set up CORS and authentication',
        'Deploy serverless API'
      ],
      learnMore: 'Lambda Proxy with ~20 files - serverless REST API'
    },
    popularity: 7,
    estimatedSetupTime: '15-20 minutes',
    recommendedFor: ['Serverless REST APIs', 'API Gateway integration', 'Cloud-native apps']
  },

  // Library Templates - Production Ready ✅
  {
    id: 'library-standard',
    name: 'Go Library',
    description: 'Well-structured Go library with comprehensive documentation, examples, and testing suite',
    category: 'library',
    complexity: 'beginner',
    icon: '📚',
    color: '#10B981',
    tags: ['Library', 'Package', 'Documentation', 'Testing', 'Reusable'],
    useCase: 'Creating reusable Go packages and libraries for the community',
    techStack: ['Standard Library', 'Comprehensive Tests', 'Godoc', 'Examples'],
    config: {
      projectType: 'library',
      framework: 'standard',
      architecture: 'standard',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Public API    │────│   Core Logic    │────│   Internal      │
│  (Exported)     │    │  (Business)     │    │   Helpers       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Examples      │    │   Tests         │    │   Documentation │
│  (Usage)        │    │  (Unit/Int)     │    │   (Godoc)       │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['Public API', 'Core Logic', 'Tests', 'Examples', 'Documentation'],
      patterns: ['Library Pattern', 'API Design', 'Test-driven Development']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'go test ./...',
        'go run examples/main.go',
        'godoc -http=:6060'
      ],
      nextSteps: [
        'Design your public API',
        'Implement core functionality',
        'Write comprehensive tests',
        'Create usage examples'
      ],
      learnMore: 'Go library with ~10 files - ready for publishing'
    },
    popularity: 6,
    estimatedSetupTime: '5-10 minutes',
    recommendedFor: ['Open source libraries', 'Internal packages', 'Code reuse']
  },

  // Enterprise Templates - Production Ready ✅
  {
    id: 'grpc-gateway',
    name: 'gRPC Gateway',
    description: 'Dual HTTP/gRPC API with protocol buffers, enhanced interceptors, and production-ready features',
    category: 'microservice',
    complexity: 'advanced',
    icon: '🔧',
    color: '#EF4444',
    tags: ['gRPC', 'HTTP', 'Protocol Buffers', 'Microservices', 'Enterprise'],
    useCase: 'Enterprise services requiring both gRPC and HTTP APIs with unified middleware',
    techStack: ['gRPC', 'Protocol Buffers', 'gRPC Gateway', 'Enhanced Interceptors'],
    config: {
      projectType: 'grpc-gateway',
      framework: 'grpc',
      architecture: 'clean',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │────│  gRPC Gateway   │────│   gRPC Server   │
│  (REST/JSON)    │    │  (Translation)  │    │  (Protobuf)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   OpenAPI       │    │   Interceptors  │    │   Services      │
│  (Swagger)      │    │  (Middleware)   │    │  (Business)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['gRPC Server', 'HTTP Gateway', 'Protocol Buffers', 'Interceptors', 'OpenAPI'],
      patterns: ['gRPC Pattern', 'API Gateway', 'Protocol Translation', 'Microservices']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'make proto',
        'go run cmd/server/main.go',
        'curl http://localhost:8080/v1/health'
      ],
      nextSteps: [
        'Define protocol buffers',
        'Implement gRPC services',
        'Configure HTTP gateway',
        'Add interceptors and middleware'
      ],
      learnMore: 'gRPC Gateway with 45 files - dual HTTP/gRPC APIs'
    },
    popularity: 8,
    estimatedSetupTime: '25-30 minutes',
    recommendedFor: ['Enterprise APIs', 'Microservices', 'High-performance systems']
  },

  {
    id: 'monolith',
    name: 'Monolithic Application',
    description: 'Full-stack monolithic application with background jobs, multi-layer caching, and performance monitoring',
    category: 'fullstack',
    complexity: 'intermediate',
    icon: '🏢',
    color: '#8B5CF6',
    tags: ['Monolith', 'Full-stack', 'Background Jobs', 'Caching', 'Monitoring'],
    useCase: 'Complete web applications with integrated frontend, backend, and background processing',
    techStack: ['Full-stack', 'Background Jobs', 'Multi-layer Caching', 'Performance Monitoring'],
    config: {
      projectType: 'monolith',
      framework: 'gin',
      architecture: 'layered',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Layer     │────│   Service       │────│   Data Layer    │
│  (Controllers)  │    │   Layer         │    │  (Database)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Background     │    │   Caching       │    │   Monitoring    │
│  Jobs           │    │   Layer         │    │   Layer         │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['Web Layer', 'Service Layer', 'Data Layer', 'Background Jobs', 'Caching'],
      patterns: ['Layered Architecture', 'Background Processing', 'Caching Strategy']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'docker-compose up -d',
        'go run cmd/migrate/main.go',
        'go run cmd/server/main.go'
      ],
      nextSteps: [
        'Configure database and cache',
        'Set up background job processing',
        'Add monitoring and metrics',
        'Deploy full-stack application'
      ],
      learnMore: 'Monolith with 72 files - production web apps'
    },
    popularity: 7,
    estimatedSetupTime: '30-40 minutes',
    recommendedFor: ['Full-stack apps', 'Integrated systems', 'Traditional architectures']
  },

  {
    id: 'microservice-standard',
    name: 'Enterprise Microservice',
    description: 'Enterprise-grade gRPC microservice with OpenTelemetry, rate limiting, and resilience patterns',
    category: 'microservice',
    complexity: 'advanced',
    icon: '🔧',
    color: '#06B6D4',
    tags: ['Microservice', 'gRPC', 'OpenTelemetry', 'Rate Limiting', 'Resilience'],
    useCase: 'Enterprise microservices with comprehensive observability and resilience features',
    techStack: ['gRPC', 'OpenTelemetry', 'Rate Limiting', 'Circuit Breaker', 'Service Mesh'],
    config: {
      projectType: 'microservice',
      framework: 'grpc',
      architecture: 'clean',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   gRPC Server   │────│   Business      │────│   Data Layer    │
│  (Service Mesh) │    │   Logic         │    │  (Repository)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Observability  │    │  Resilience     │    │   Security      │
│ (OpenTelemetry) │    │  Patterns       │    │  (Rate Limit)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘`,
      components: ['gRPC Service', 'Business Logic', 'Data Layer', 'Observability', 'Resilience'],
      patterns: ['Microservice', 'Circuit Breaker', 'Rate Limiting', 'Observability']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'make proto',
        'docker-compose up -d',
        'go run cmd/server/main.go'
      ],
      nextSteps: [
        'Configure OpenTelemetry',
        'Set up service discovery',
        'Add resilience patterns',
        'Deploy to Kubernetes'
      ],
      learnMore: 'Microservice with 47 files - enterprise gRPC services'
    },
    popularity: 8,
    estimatedSetupTime: '35-45 minutes',
    recommendedFor: ['Enterprise microservices', 'Distributed systems', 'Cloud-native apps']
  }
]

// Template filtering and search utilities
export const getTemplatesByCategory = (category: ProjectTemplate['category']) => {
  return PROJECT_TEMPLATES.filter(template => template.category === category)
}

export const getTemplatesByComplexity = (complexity: ProjectTemplate['complexity']) => {
  return PROJECT_TEMPLATES.filter(template => template.complexity === complexity)
}

export const searchTemplates = (query: string) => {
  const searchTerm = query.toLowerCase()
  return PROJECT_TEMPLATES.filter(template => 
    template.name.toLowerCase().includes(searchTerm) ||
    template.description.toLowerCase().includes(searchTerm) ||
    template.tags.some(tag => tag.toLowerCase().includes(searchTerm)) ||
    template.useCase.toLowerCase().includes(searchTerm)
  )
}

export const getPopularTemplates = (limit = 6) => {
  return PROJECT_TEMPLATES
    .sort((a, b) => b.popularity - a.popularity)
    .slice(0, limit)
}

export const getRecommendedTemplates = (userLevel: 'beginner' | 'intermediate' | 'advanced' | 'expert') => {
  const complexityMap = {
    beginner: ['beginner', 'intermediate'],
    intermediate: ['beginner', 'intermediate', 'advanced'],
    advanced: ['intermediate', 'advanced', 'expert'],
    expert: ['advanced', 'expert']
  }
  
  return PROJECT_TEMPLATES.filter(template => 
    complexityMap[userLevel].includes(template.complexity)
  ).sort((a, b) => b.popularity - a.popularity)
}

export const TEMPLATE_CATEGORIES = [
  { id: 'api', name: 'Web APIs', description: 'REST APIs and web services', icon: '🌐', count: 4 },
  { id: 'cli', name: 'CLI Tools', description: 'Command-line applications', icon: '💻', count: 2 },
  { id: 'microservice', name: 'Microservices', description: 'Distributed services', icon: '🔧', count: 2 },
  { id: 'library', name: 'Libraries', description: 'Reusable packages', icon: '📚', count: 1 },
  { id: 'serverless', name: 'Serverless', description: 'Cloud functions and APIs', icon: '⚡', count: 2 },
  { id: 'fullstack', name: 'Full-Stack', description: 'Complete applications', icon: '🎯', count: 1 }
] as const

export const COMPLEXITY_LEVELS = [
  { id: 'beginner', name: 'Beginner', description: 'Simple projects to get started', color: '#10B981', count: 3 },
  { id: 'intermediate', name: 'Intermediate', description: 'Moderate complexity with common patterns', color: '#3B82F6', count: 6 },
  { id: 'advanced', name: 'Advanced', description: 'Complex architecture patterns', color: '#F59E0B', count: 3 },
  { id: 'expert', name: 'Expert', description: 'Enterprise-grade solutions', color: '#EF4444', count: 0 }
] as const