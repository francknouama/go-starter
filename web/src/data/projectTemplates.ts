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

export const PROJECT_TEMPLATES: ProjectTemplate[] = [
  // API Templates
  {
    id: 'blog-api',
    name: 'Blog API',
    description: 'Complete REST API for a blog platform with authentication, posts, comments, and media management',
    category: 'api',
    complexity: 'intermediate',
    icon: '📝',
    color: '#3B82F6',
    tags: ['REST', 'CRUD', 'Authentication', 'JWT', 'PostgreSQL'],
    useCase: 'Building a content management system or blog platform',
    techStack: ['Gin', 'PostgreSQL', 'GORM', 'JWT', 'Zap Logger'],
    config: {
      projectType: 'web-api',
      framework: 'gin',
      architecture: 'clean',
      logger: 'zap',
      features: {
        database: {
          driver: 'postgres',
          orm: 'gorm'
        },
        authentication: {
          type: 'jwt',
          providers: []
        }
      }
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   HTTP Router   │────│   Controllers    │────│   Use Cases     │
│  (Gin/Fiber)    │    │                  │    │  (Blog Logic)   │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                         │
                                                ┌─────────────────┐
                                                │  Repositories   │
                                                │ (Data Access)   │
                                                └─────────────────┘
                                                         │
                                                ┌─────────────────┐
                                                │   PostgreSQL    │
                                                │   Database      │
                                                └─────────────────┘`,
      components: ['User Management', 'Post CRUD', 'Comment System', 'Media Upload', 'Auth Middleware'],
      patterns: ['Clean Architecture', 'Repository Pattern', 'JWT Authentication', 'Input Validation']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'docker-compose up -d postgres',
        'go run cmd/migrate/main.go',
        'go run cmd/server/main.go'
      ],
      nextSteps: [
        'Configure your database connection in .env',
        'Run migrations to create tables',
        'Test the API endpoints with the included Postman collection',
        'Customize the blog post schema for your needs'
      ],
      learnMore: 'Check out docs/API_GUIDE.md for endpoint documentation'
    },
    popularity: 8,
    estimatedSetupTime: '15-20 minutes',
    recommendedFor: ['Content platforms', 'Personal blogs', 'CMS backends']
  },
  
  {
    id: 'ecommerce-api',
    name: 'E-commerce API',
    description: 'Production-ready e-commerce backend with products, orders, payments, inventory, and customer management',
    category: 'api',
    complexity: 'advanced',
    icon: '🛒',
    color: '#10B981',
    tags: ['E-commerce', 'Payments', 'Inventory', 'Orders', 'Microservice-ready'],
    useCase: 'Building online stores, marketplaces, or retail platforms',
    techStack: ['Echo', 'PostgreSQL', 'Redis', 'Stripe API', 'Logrus'],
    config: {
      projectType: 'web-api',
      framework: 'echo',
      architecture: 'ddd',
      logger: 'logrus',
      features: {
        database: {
          driver: 'postgres',
          orm: 'gorm'
        },
        authentication: {
          type: 'jwt',
          providers: []
        }
      }
    },
    architecture: {
      diagram: `
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  Products   │  │   Orders    │  │  Payments   │  │  Customers  │
│   Domain    │  │   Domain    │  │   Domain    │  │   Domain    │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
       │                 │                 │                 │
┌─────────────────────────────────────────────────────────────────┐
│                    Shared Infrastructure                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │ PostgreSQL  │  │    Redis    │  │  Event Bus  │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘`,
      components: ['Product Catalog', 'Order Management', 'Payment Processing', 'Inventory Tracking', 'Customer Accounts'],
      patterns: ['Domain-Driven Design', 'Event Sourcing', 'CQRS', 'Saga Pattern', 'API Gateway']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'docker-compose up -d postgres redis',
        'go run cmd/migrate/main.go',
        'go run cmd/server/main.go'
      ],
      nextSteps: [
        'Configure payment provider (Stripe) credentials',
        'Set up Redis for session management and caching',
        'Configure email notifications for orders',
        'Review the domain models and customize for your business'
      ],
      learnMore: 'See docs/ECOMMERCE_GUIDE.md for business logic documentation'
    },
    popularity: 9,
    estimatedSetupTime: '25-30 minutes',
    recommendedFor: ['Online stores', 'B2B marketplaces', 'Retail platforms']
  },

  {
    id: 'saas-backend',
    name: 'SaaS Backend',
    description: 'Multi-tenant SaaS backend with subscription management, billing, user workspaces, and admin dashboard',
    category: 'api',
    complexity: 'expert',
    icon: '🏢',
    color: '#8B5CF6',
    tags: ['Multi-tenant', 'Subscriptions', 'Billing', 'Admin Panel', 'Scalable'],
    useCase: 'Building software-as-a-service applications with subscription models',
    techStack: ['Fiber', 'PostgreSQL', 'Redis', 'Stripe', 'Slog'],
    config: {
      projectType: 'web-api',
      framework: 'fiber',
      architecture: 'hexagonal',
      logger: 'slog',
      features: {
        database: {
          driver: 'postgres',
          orm: 'ent'
        },
        authentication: {
          type: 'oauth2',
          providers: ['google', 'github']
        }
      }
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Tenant A      │    │   Tenant B      │    │   Tenant C      │
│   Workspace     │    │   Workspace     │    │   Workspace     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
┌─────────────────────────────────────────────────────────────────┐
│                      Multi-Tenant Core                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │    Auth     │  │   Billing   │  │    Admin    │            │
│  │   Service   │  │   Service   │  │    Panel    │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
└─────────────────────────────────────────────────────────────────┘
                               │
                     ┌─────────────────┐
                     │   Shared DB     │
                     │  (Row Level     │
                     │   Security)     │
                     └─────────────────┘`,
      components: ['Tenant Management', 'Subscription Billing', 'User Workspaces', 'Admin Dashboard', 'Multi-tenant Auth'],
      patterns: ['Hexagonal Architecture', 'Multi-tenancy', 'Event-driven', 'Domain Events', 'Ports & Adapters']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'docker-compose up -d postgres redis',
        'go run cmd/migrate/main.go --tenant=system',
        'go run cmd/server/main.go'
      ],
      nextSteps: [
        'Configure OAuth providers (Google, GitHub)',
        'Set up Stripe webhook endpoints for billing',
        'Configure row-level security policies for multi-tenancy',
        'Customize tenant onboarding flow'
      ],
      learnMore: 'Read docs/MULTI_TENANT_GUIDE.md for architecture details'
    },
    popularity: 7,
    estimatedSetupTime: '45-60 minutes',
    recommendedFor: ['SaaS platforms', 'B2B applications', 'Enterprise software']
  },

  // CLI Templates
  {
    id: 'devops-cli',
    name: 'DevOps CLI Tool',
    description: 'Command-line tool for managing AWS resources, Kubernetes clusters, and deployment pipelines',
    category: 'cli',
    complexity: 'intermediate',
    icon: '⚙️',
    color: '#F59E0B',
    tags: ['DevOps', 'AWS', 'Kubernetes', 'CI/CD', 'Infrastructure'],
    useCase: 'Automating infrastructure management and deployment workflows',
    techStack: ['Cobra', 'AWS SDK', 'K8s Client', 'Viper Config', 'Zerolog'],
    config: {
      projectType: 'cli',
      framework: 'cobra',
      architecture: 'standard',
      logger: 'zerolog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Commands      │    │   Providers     │    │   Resources     │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │ aws       │  │────│  │    AWS    │  │────│  │    EC2    │  │
│  │ k8s       │  │    │  │    SDK    │  │    │  │   EKS     │  │
│  │ deploy    │  │    │  │           │  │    │  │   S3      │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
└─────────────────┘    │  ┌───────────┐  │    │  ┌───────────┐  │
                       │  │    K8s    │  │────│  │   Pods    │  │
                       │  │  Client   │  │    │  │ Services  │  │
                       │  └───────────┘  │    │  └───────────┘  │
                       └─────────────────┘    └─────────────────┘`,
      components: ['AWS Resource Management', 'Kubernetes Operations', 'Deployment Automation', 'Config Management', 'Progress Tracking'],
      patterns: ['Command Pattern', 'Provider Pattern', 'Configuration Management', 'Progress Indicators']
    },
    quickStart: {
      commands: [
        'go build -o devops-tool cmd/main.go',
        './devops-tool --help',
        './devops-tool config init',
        './devops-tool aws list-instances'
      ],
      nextSteps: [
        'Configure AWS credentials (aws configure)',
        'Set up kubeconfig for Kubernetes access',
        'Create configuration profiles for different environments',
        'Add custom commands for your specific workflow'
      ],
      learnMore: 'Check docs/DEVOPS_COMMANDS.md for all available operations'
    },
    popularity: 8,
    estimatedSetupTime: '10-15 minutes',
    recommendedFor: ['DevOps engineers', 'SRE teams', 'Infrastructure automation']
  },

  {
    id: 'file-processor',
    name: 'File Processor',
    description: 'Batch file processing tool with progress tracking, parallel processing, and format conversion',
    category: 'cli',
    complexity: 'beginner',
    icon: '📁',
    color: '#6366F1',
    tags: ['File Processing', 'Batch Operations', 'Parallel', 'Progress'],
    useCase: 'Processing large batches of files with transformations and validations',
    techStack: ['Cobra', 'Concurrent Processing', 'Progress Bars', 'Slog'],
    config: {
      projectType: 'cli',
      framework: 'cobra',
      architecture: 'standard',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   File Input    │────│   Processor     │────│   File Output   │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │  Scanner  │  │    │  │   Worker  │  │    │  │  Writer   │  │
│  │  Filter   │  │    │  │   Pool    │  │    │  │ Validator │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                               │
                    ┌─────────────────┐
                    │   Progress      │
                    │   Tracking      │
                    └─────────────────┘`,
      components: ['File Scanner', 'Worker Pool', 'Progress Tracking', 'Format Converters', 'Validation Engine'],
      patterns: ['Worker Pool Pattern', 'Pipeline Pattern', 'Observer Pattern', 'Strategy Pattern']
    },
    quickStart: {
      commands: [
        'go build -o file-processor cmd/main.go',
        './file-processor process --input ./data --output ./processed',
        './file-processor convert --format json --input ./csv-files',
        './file-processor validate --schema ./schema.json ./data'
      ],
      nextSteps: [
        'Configure processing rules in config.yaml',
        'Add custom file processors for your formats',
        'Set up scheduling for automated batch jobs',
        'Add monitoring and alerting for long-running jobs'
      ],
      learnMore: 'See docs/PROCESSING_GUIDE.md for supported formats and operations'
    },
    popularity: 6,
    estimatedSetupTime: '5-10 minutes',
    recommendedFor: ['Data processing', 'ETL workflows', 'Batch operations']
  },

  {
    id: 'code-generator',
    name: 'Code Generator',
    description: 'Template-based code generator with custom templates, variable substitution, and multi-language support',
    category: 'cli',
    complexity: 'advanced',
    icon: '🚀',
    color: '#EF4444',
    tags: ['Code Generation', 'Templates', 'Developer Tools', 'Multi-language'],
    useCase: 'Generating boilerplate code, scaffolding projects, or creating consistent code patterns',
    techStack: ['Cobra', 'Go Templates', 'YAML Config', 'Git Integration', 'Zap'],
    config: {
      projectType: 'cli',
      framework: 'cobra',
      architecture: 'clean',
      logger: 'zap'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Templates     │────│   Generator     │────│   Generated     │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │    Go     │  │    │  │ Template  │  │    │  │  Source   │  │
│  │   React   │  │    │  │  Engine   │  │    │  │   Code    │  │
│  │  Python   │  │    │  └───────────┘  │    │  └───────────┘  │
│  └───────────┘  │    │  ┌───────────┐  │    │  ┌───────────┐  │
└─────────────────┘    │  │ Variable  │  │    │  │   Docs    │  │
                       │  │ Resolver  │  │    │  │   Tests   │  │
                       │  └───────────┘  │    │  └───────────┘  │
                       └─────────────────┘    └─────────────────┘`,
      components: ['Template Engine', 'Variable Resolution', 'Multi-language Support', 'Git Integration', 'Validation'],
      patterns: ['Template Method', 'Strategy Pattern', 'Chain of Responsibility', 'Factory Pattern']
    },
    quickStart: {
      commands: [
        'go build -o codegen cmd/main.go',
        './codegen init --template-dir ./templates',
        './codegen generate api --name UserService --output ./generated',
        './codegen list-templates'
      ],
      nextSteps: [
        'Create custom templates in the templates/ directory',
        'Configure variable definitions in template.yaml files',
        'Set up Git repositories for template sharing',
        'Add validation rules for generated code'
      ],
      learnMore: 'Read docs/TEMPLATE_GUIDE.md for creating custom templates'
    },
    popularity: 7,
    estimatedSetupTime: '15-20 minutes',
    recommendedFor: ['Developer tools', 'Code scaffolding', 'Project generators']
  },

  // Microservice Templates
  {
    id: 'user-service',
    name: 'User Service',
    description: 'Microservice for user authentication, profile management, and access control with gRPC APIs',
    category: 'microservice',
    complexity: 'intermediate',
    icon: '👤',
    color: '#06B6D4',
    tags: ['Microservice', 'gRPC', 'Authentication', 'User Management', 'Distributed'],
    useCase: 'User management microservice in a distributed system architecture',
    techStack: ['gRPC', 'PostgreSQL', 'JWT', 'Protocol Buffers', 'Logrus'],
    config: {
      projectType: 'microservice',
      framework: 'grpc',
      architecture: 'clean',
      logger: 'logrus',
      features: {
        database: {
          driver: 'postgres',
          orm: 'gorm'
        },
        authentication: {
          type: 'jwt',
          providers: []
        }
      }
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   gRPC Server   │────│   User Service  │────│   Database      │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │    Auth   │  │    │  │   Auth    │  │    │  │   Users   │  │
│  │ Profile   │  │    │  │  Domain   │  │    │  │ Profiles  │  │
│  │  Admin    │  │    │  │  Service  │  │    │  │   Roles   │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │
┌─────────────────┐    ┌─────────────────┐
│  Service Mesh   │    │     Events      │
│  (Istio/Envoy)  │    │   (Message      │
│                 │    │     Bus)        │
└─────────────────┘    └─────────────────┘`,
      components: ['gRPC Server', 'User Domain', 'Auth Service', 'Profile Management', 'Role-based Access'],
      patterns: ['Microservice', 'Domain-Driven Design', 'Event Sourcing', 'Circuit Breaker', 'Service Discovery']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'make proto-gen',
        'docker-compose up -d postgres',
        'go run cmd/server/main.go'
      ],
      nextSteps: [
        'Configure service discovery (Consul/etcd)',
        'Set up monitoring and tracing (Jaeger)',
        'Configure message bus for events (NATS/Kafka)',
        'Add integration tests with test containers'
      ],
      learnMore: 'Check docs/MICROSERVICE_GUIDE.md for deployment patterns'
    },
    popularity: 8,
    estimatedSetupTime: '20-25 minutes',
    recommendedFor: ['Distributed systems', 'Microservice architecture', 'User management']
  },

  {
    id: 'notification-service',
    name: 'Notification Service',
    description: 'Multi-channel notification service supporting email, SMS, push notifications, and webhooks',
    category: 'microservice',
    complexity: 'advanced',
    icon: '📧',
    color: '#84CC16',
    tags: ['Notifications', 'Multi-channel', 'Queue', 'Templates', 'Retry Logic'],
    useCase: 'Centralized notification handling for distributed applications',
    techStack: ['gRPC', 'Redis', 'Message Queue', 'Template Engine', 'Zap'],
    config: {
      projectType: 'microservice',
      framework: 'grpc',
      architecture: 'event-driven',
      logger: 'zap',
      features: {
        database: {
          driver: 'redis',
          orm: 'redis-client'
        }
      }
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Event Bus     │────│  Notification   │────│   Channels      │
│  ┌───────────┐  │    │    Service      │    │  ┌───────────┐  │
│  │   Events  │  │    │  ┌───────────┐  │    │  │   Email   │  │
│  │   Queue   │  │    │  │ Template  │  │    │  │    SMS    │  │
│  └───────────┘  │    │  │  Engine   │  │    │  │   Push    │  │
└─────────────────┘    │  └───────────┘  │    │  │ Webhook   │  │
                       │  ┌───────────┐  │    │  └───────────┘  │
                       │  │   Retry   │  │    └─────────────────┘
                       │  │  Handler  │  │
                       │  └───────────┘  │
                       └─────────────────┘`,
      components: ['Event Processing', 'Template Engine', 'Channel Handlers', 'Retry Logic', 'Delivery Tracking'],
      patterns: ['Event-driven Architecture', 'Template Method', 'Strategy Pattern', 'Circuit Breaker', 'Dead Letter Queue']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'docker-compose up -d redis nats',
        'make proto-gen',
        'go run cmd/server/main.go'
      ],
      nextSteps: [
        'Configure email provider (SendGrid/AWS SES)',
        'Set up SMS provider (Twilio/AWS SNS)',
        'Configure push notification services (FCM/APNS)',
        'Create notification templates for your use cases'
      ],
      learnMore: 'See docs/NOTIFICATION_PATTERNS.md for delivery patterns'
    },
    popularity: 7,
    estimatedSetupTime: '25-30 minutes',
    recommendedFor: ['Notification systems', 'Event-driven architecture', 'Communication platforms']
  },

  // Library Templates
  {
    id: 'sdk-library',
    name: 'SDK Library',
    description: 'Well-structured SDK library with comprehensive documentation, examples, and testing',
    category: 'library',
    complexity: 'intermediate',
    icon: '📚',
    color: '#8B5CF6',
    tags: ['SDK', 'Library', 'Documentation', 'Examples', 'Testing'],
    useCase: 'Creating client libraries and SDKs for APIs or services',
    techStack: ['Standard Library', 'Comprehensive Tests', 'Examples', 'Godoc'],
    config: {
      projectType: 'library',
      framework: 'standard',
      architecture: 'standard',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Public API    │────│   Core Logic    │────│   HTTP Client   │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │  Client   │  │    │  │ Business  │  │    │  │  Request  │  │
│  │  Config   │  │    │  │   Logic   │  │    │  │  Handler  │  │
│  │ Methods   │  │    │  └───────────┘  │    │  │   Auth    │  │
│  └───────────┘  │    │  ┌───────────┐  │    │  └───────────┘  │
└─────────────────┘    │  │   Models  │  │    └─────────────────┘
                       │  │    DTOs   │  │
                       │  └───────────┘  │
                       └─────────────────┘`,
      components: ['Public API', 'HTTP Client', 'Authentication', 'Error Handling', 'Response Models'],
      patterns: ['Builder Pattern', 'Factory Pattern', 'Strategy Pattern', 'Decorator Pattern']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'go test ./...',
        'go run examples/basic/main.go',
        'godoc -http=:6060'
      ],
      nextSteps: [
        'Customize the client configuration options',
        'Add your API endpoints and models',
        'Write comprehensive examples for common use cases',
        'Set up CI/CD for automated testing and releases'
      ],
      learnMore: 'Read docs/SDK_DEVELOPMENT.md for API design best practices'
    },
    popularity: 6,
    estimatedSetupTime: '10-15 minutes',
    recommendedFor: ['API clients', 'SDK development', 'Third-party integrations']
  },

  // Serverless Templates
  {
    id: 'serverless-api',
    name: 'Serverless API',
    description: 'AWS Lambda-based REST API with API Gateway, DynamoDB, and event-driven architecture',
    category: 'serverless',
    complexity: 'intermediate',
    icon: '⚡',
    color: '#F97316',
    tags: ['Serverless', 'AWS Lambda', 'API Gateway', 'DynamoDB', 'Event-driven'],
    useCase: 'Building scalable serverless APIs with pay-per-use pricing',
    techStack: ['AWS Lambda', 'API Gateway', 'DynamoDB', 'CloudFormation', 'Slog'],
    config: {
      projectType: 'lambda-proxy',
      framework: 'lambda',
      architecture: 'event-driven',
      logger: 'slog'
    },
    architecture: {
      diagram: `
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  API Gateway    │────│   Lambda Proxy  │────│   DynamoDB      │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │   Routes  │  │    │  │ Handlers  │  │    │  │   Tables  │  │
│  │    Auth   │  │    │  │ Business  │  │    │  │   Indexes │  │
│  │Validation │  │    │  │   Logic   │  │    │  └───────────┘  │
│  └───────────┘  │    │  └───────────┘  │    └─────────────────┘
└─────────────────┘    └─────────────────┘
         │                       │
┌─────────────────┐    ┌─────────────────┐
│  CloudWatch     │    │   EventBridge   │
│  (Monitoring)   │    │    (Events)     │
└─────────────────┘    └─────────────────┘`,
      components: ['API Gateway Routes', 'Lambda Handlers', 'DynamoDB Operations', 'Event Processing', 'CloudWatch Monitoring'],
      patterns: ['Serverless', 'Event-driven', 'Single Responsibility', 'Circuit Breaker', 'Bulkhead Pattern']
    },
    quickStart: {
      commands: [
        'go mod tidy',
        'sam build',
        'sam local start-api',
        'sam deploy --guided'
      ],
      nextSteps: [
        'Configure AWS credentials and region',
        'Customize the DynamoDB table schema',
        'Set up CloudWatch alarms for monitoring',
        'Add authentication with AWS Cognito or custom authorizer'
      ],
      learnMore: 'Check docs/SERVERLESS_DEPLOYMENT.md for AWS setup'
    },
    popularity: 8,
    estimatedSetupTime: '15-20 minutes',
    recommendedFor: ['Serverless APIs', 'Event-driven systems', 'Cost-optimized backends']
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
  { id: 'api', name: 'Web APIs', description: 'REST APIs and web services', icon: '🌐' },
  { id: 'cli', name: 'CLI Tools', description: 'Command-line applications', icon: '💻' },
  { id: 'microservice', name: 'Microservices', description: 'Distributed services', icon: '🔧' },
  { id: 'library', name: 'Libraries', description: 'Reusable packages', icon: '📚' },
  { id: 'serverless', name: 'Serverless', description: 'Cloud functions and APIs', icon: '⚡' },
  { id: 'fullstack', name: 'Full-Stack', description: 'Complete applications', icon: '🎯' }
] as const

export const COMPLEXITY_LEVELS = [
  { id: 'beginner', name: 'Beginner', description: 'Simple projects to get started', color: '#10B981' },
  { id: 'intermediate', name: 'Intermediate', description: 'Moderate complexity with common patterns', color: '#3B82F6' },
  { id: 'advanced', name: 'Advanced', description: 'Complex architecture patterns', color: '#F59E0B' },
  { id: 'expert', name: 'Expert', description: 'Enterprise-grade solutions', color: '#EF4444' }
] as const