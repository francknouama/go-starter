/**
 * Mock API Service - For development without backend
 */

import type { 
  HealthResponse, 
  Blueprint, 
  ProjectConfig, 
  PreviewResponse, 
  GenerationRequest, 
  GenerationResponse,
  FileNode 
} from './api'

// Mock data
const mockHealth: HealthResponse = {
  status: 'healthy',
  timestamp: new Date().toISOString(),
  version: '1.0.0-mock',
  uptime: '1h 30m',
  checks: {
    database: 'healthy',
    memory: 'healthy',
    disk: 'healthy'
  }
}

const mockBlueprints: Blueprint[] = [
  {
    id: 'web-api-standard',
    name: 'Standard Web API',
    description: 'Basic REST API with standard architecture',
    category: 'web-api',
    difficulty: 'beginner',
    features: ['REST', 'JSON', 'Middleware'],
    architectures: ['standard'],
    frameworks: ['gin', 'echo', 'fiber'],
    estimatedFiles: 15,
    tags: ['api', 'rest', 'web']
  },
  {
    id: 'cli-simple',
    name: 'Simple CLI',
    description: 'Basic command-line application',
    category: 'cli',
    difficulty: 'beginner',
    features: ['Commands', 'Flags'],
    architectures: ['standard'],
    frameworks: ['cobra'],
    estimatedFiles: 8,
    tags: ['cli', 'command']
  }
]

const mockDefaultConfig: ProjectConfig = {
  projectName: 'my-go-project',
  moduleUrl: 'github.com/user/my-go-project',
  projectType: 'web-api',
  architecture: 'standard',
  framework: 'gin',
  logger: 'slog',
  goVersion: '1.21'
}

// Dynamic file structure generation based on configuration
const generateFileStructure = (config: ProjectConfig): FileNode[] => {
  const projectName = config.projectName || 'my-go-project'
  
  const baseFiles: FileNode[] = [
    { name: 'main.go', type: 'file', path: '/main.go', size: 1024 },
    { name: 'go.mod', type: 'file', path: '/go.mod', size: 256 },
    { name: 'go.sum', type: 'file', path: '/go.sum', size: 512 },
    { name: 'README.md', type: 'file', path: '/README.md', size: 2048 },
  ]

  // Add Dockerfile for web-api projects
  if (config.projectType === 'web-api') {
    baseFiles.push({ name: 'Dockerfile', type: 'file', path: '/Dockerfile', size: 512 })
  }

  // Command directory structure based on project type
  const cmdChildren: FileNode[] = []
  if (config.projectType === 'web-api') {
    cmdChildren.push({ name: 'server.go', type: 'file', path: '/cmd/server.go', size: 2048 })
  } else if (config.projectType === 'cli') {
    cmdChildren.push(
      { name: 'root.go', type: 'file', path: '/cmd/root.go', size: 1536 },
      { name: 'version.go', type: 'file', path: '/cmd/version.go', size: 512 }
    )
  }

  // Internal directory structure
  const internalChildren: FileNode[] = [
    { name: 'config.go', type: 'file', path: '/internal/config.go', size: 1024 },
  ]

  // Add project-type specific files
  if (config.projectType === 'web-api') {
    // Framework-specific files
    if (config.framework === 'gin') {
      internalChildren.push(
        { name: 'handler.go', type: 'file', path: '/internal/handler.go', size: 3072 },
        { name: 'middleware.go', type: 'file', path: '/internal/middleware.go', size: 1536 },
        { name: 'router.go', type: 'file', path: '/internal/router.go', size: 2048 }
      )
    } else if (config.framework === 'echo') {
      internalChildren.push(
        { name: 'handlers.go', type: 'file', path: '/internal/handlers.go', size: 3072 },
        { name: 'middleware.go', type: 'file', path: '/internal/middleware.go', size: 1536 },
        { name: 'server.go', type: 'file', path: '/internal/server.go', size: 2048 }
      )
    } else if (config.framework === 'fiber') {
      internalChildren.push(
        { name: 'routes.go', type: 'file', path: '/internal/routes.go', size: 3072 },
        { name: 'middleware.go', type: 'file', path: '/internal/middleware.go', size: 1536 },
        { name: 'app.go', type: 'file', path: '/internal/app.go', size: 2048 }
      )
    }
  }

  // Logger-specific files
  if (config.logger && config.logger !== 'slog') {
    internalChildren.push({
      name: 'logger.go',
      type: 'file',
      path: '/internal/logger.go',
      size: 1024
    })
  }

  // Add database files if database is enabled
  const hasDatabase = config.features?.database?.driver
  if (hasDatabase) {
    const dbChildren: FileNode[] = [
      { name: 'connection.go', type: 'file', path: '/internal/database/connection.go', size: 1536 },
      { name: 'migrations.go', type: 'file', path: '/internal/database/migrations.go', size: 2048 },
      { name: 'models.go', type: 'file', path: '/internal/database/models.go', size: 2560 }
    ]

    // Add driver-specific file
    const driver = config.features.database.driver
    dbChildren.push({
      name: `${driver}.go`,
      type: 'file',
      path: `/internal/database/${driver}.go`,
      size: 1024
    })

    // Add ORM-specific file
    const orm = config.features.database.orm
    if (orm !== 'sqlx') { // sqlx doesn't need a separate file
      dbChildren.push({
        name: `${orm}_models.go`,
        type: 'file',
        path: `/internal/database/${orm}_models.go`,
        size: 1536
      })
    }

    internalChildren.push({
      name: 'database',
      type: 'directory',
      path: '/internal/database',
      children: dbChildren
    })
  }

  // Add authentication files if authentication is enabled
  const hasAuth = config.features?.authentication?.type
  if (hasAuth) {
    const authChildren: FileNode[] = [
      { name: 'middleware.go', type: 'file', path: '/internal/auth/middleware.go', size: 2048 },
      { name: 'service.go', type: 'file', path: '/internal/auth/service.go', size: 3072 }
    ]

    // Add auth type-specific file
    const authType = config.features.authentication.type
    authChildren.push({
      name: `${authType}.go`,
      type: 'file',
      path: `/internal/auth/${authType}.go`,
      size: 1536
    })

    // Add config file for OAuth2
    if (authType === 'oauth2') {
      authChildren.push({
        name: 'oauth_config.go',
        type: 'file',
        path: '/internal/auth/oauth_config.go',
        size: 1024
      })
    }

    // Add token validation for JWT
    if (authType === 'jwt') {
      authChildren.push({
        name: 'token_validator.go',
        type: 'file',
        path: '/internal/auth/token_validator.go',
        size: 1280
      })
    }

    internalChildren.push({
      name: 'auth',
      type: 'directory',
      path: '/internal/auth',
      children: authChildren
    })
  }

  // Architecture-specific structure adjustments
  if (config.architecture === 'clean') {
    // Clean Architecture: Add domain and usecase layers
    internalChildren.push(
      {
        name: 'domain',
        type: 'directory',
        path: '/internal/domain',
        children: [
          { name: 'entities.go', type: 'file', path: '/internal/domain/entities.go', size: 2048 },
          { name: 'repositories.go', type: 'file', path: '/internal/domain/repositories.go', size: 1536 }
        ]
      },
      {
        name: 'usecase',
        type: 'directory',
        path: '/internal/usecase',
        children: [
          { name: 'interfaces.go', type: 'file', path: '/internal/usecase/interfaces.go', size: 1024 },
          { name: 'user_usecase.go', type: 'file', path: '/internal/usecase/user_usecase.go', size: 2560 }
        ]
      }
    )
  } else if (config.architecture === 'hexagonal') {
    // Hexagonal Architecture: Add ports and adapters
    internalChildren.push(
      {
        name: 'ports',
        type: 'directory',
        path: '/internal/ports',
        children: [
          { name: 'primary.go', type: 'file', path: '/internal/ports/primary.go', size: 1536 },
          { name: 'secondary.go', type: 'file', path: '/internal/ports/secondary.go', size: 1536 }
        ]
      },
      {
        name: 'adapters',
        type: 'directory',
        path: '/internal/adapters',
        children: [
          { name: 'http.go', type: 'file', path: '/internal/adapters/http.go', size: 2048 },
          { name: 'persistence.go', type: 'file', path: '/internal/adapters/persistence.go', size: 2048 }
        ]
      }
    )
  } else if (config.architecture === 'ddd') {
    // Domain-Driven Design: Add domain services and value objects
    internalChildren.push(
      {
        name: 'domain',
        type: 'directory',
        path: '/internal/domain',
        children: [
          { name: 'aggregates.go', type: 'file', path: '/internal/domain/aggregates.go', size: 3072 },
          { name: 'value_objects.go', type: 'file', path: '/internal/domain/value_objects.go', size: 2048 },
          { name: 'domain_services.go', type: 'file', path: '/internal/domain/domain_services.go', size: 2560 }
        ]
      },
      {
        name: 'application',
        type: 'directory',
        path: '/internal/application',
        children: [
          { name: 'commands.go', type: 'file', path: '/internal/application/commands.go', size: 2048 },
          { name: 'queries.go', type: 'file', path: '/internal/application/queries.go', size: 2048 }
        ]
      }
    )
  }

  // Test files
  const testsChildren: FileNode[] = [
    { name: 'main_test.go', type: 'file', path: '/tests/main_test.go', size: 1024 },
    { name: 'integration_test.go', type: 'file', path: '/tests/integration_test.go', size: 2048 }
  ]

  // Add feature-specific test files
  if (hasDatabase) {
    testsChildren.push({
      name: 'database_test.go',
      type: 'file',
      path: '/tests/database_test.go',
      size: 1536
    })
  }

  if (hasAuth) {
    testsChildren.push({
      name: 'auth_test.go',
      type: 'file',
      path: '/tests/auth_test.go',
      size: 2048
    })
  }

  // Build final structure
  const structure: FileNode[] = [
    {
      name: projectName,
      type: 'directory',
      path: '/',
      children: [
        ...baseFiles,
        {
          name: 'cmd',
          type: 'directory',
          path: '/cmd',
          children: cmdChildren
        },
        {
          name: 'internal',
          type: 'directory',
          path: '/internal',
          children: internalChildren
        },
        {
          name: 'tests',
          type: 'directory',
          path: '/tests',
          children: testsChildren
        }
      ]
    }
  ]

  return structure
}

// Calculate file count and size from structure
const calculateStats = (structure: FileNode[]): { fileCount: number; totalSize: number } => {
  let fileCount = 0
  let totalSize = 0
  
  const traverse = (nodes: FileNode[]) => {
    nodes.forEach(node => {
      if (node.type === 'file') {
        fileCount++
        totalSize += node.size || 0
      } else if (node.children) {
        traverse(node.children)
      }
    })
  }
  
  traverse(structure)
  return { fileCount, totalSize }
}

// Generate mock preview with static fallback
const generateMockPreview = (config?: ProjectConfig): PreviewResponse => {
  const fileStructure: FileNode[] = config ? generateFileStructure(config) : [
    {
      name: 'my-go-project',
      type: 'directory' as const,
      path: '/',
      children: [
        { name: 'main.go', type: 'file' as const, path: '/main.go', size: 1024 },
        { name: 'go.mod', type: 'file' as const, path: '/go.mod', size: 256 },
        { name: 'README.md', type: 'file' as const, path: '/README.md', size: 2048 },
      ]
    }
  ]
  
  const { fileCount, totalSize } = calculateStats(fileStructure)
  
  return {
    fileStructure,
    estimatedSize: `${(totalSize / 1024).toFixed(1)} KB`,
    fileCount
  }
}

// Mock API functions with delays to simulate network
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

export const mockApi = {
  health: {
    async getHealth(): Promise<HealthResponse> {
      await delay(200)
      return mockHealth
    },

    async getSimpleHealth(): Promise<{ status: string }> {
      await delay(100)
      return { status: 'ok' }
    },

    async getMetrics(): Promise<Record<string, unknown>> {
      await delay(150)
      return { requests: 100, uptime: 3600 }
    },

    async getStatus(): Promise<{ status: string }> {
      await delay(100)
      return { status: 'operational' }
    }
  },

  config: {
    async getDefaultConfig(): Promise<ProjectConfig> {
      await delay(300)
      return { ...mockDefaultConfig }
    },

    async getProjectTypeDetails(type: string): Promise<{ type: string; description: string; features?: string[] }> {
      await delay(200)
      return { type, description: `Details for ${type}` }
    },

    async getFrameworks(): Promise<string[]> {
      await delay(150)
      return ['gin', 'echo', 'fiber', 'chi']
    },

    async getArchitectures(): Promise<string[]> {
      await delay(150)
      return ['standard', 'clean', 'hexagonal', 'ddd']
    }
  },

  blueprints: {
    async getBlueprints(filters?: Record<string, string>): Promise<Blueprint[]> {
      await delay(400)
      return [...mockBlueprints]
    },

    async getBlueprintById(id: string): Promise<Blueprint> {
      await delay(200)
      const blueprint = mockBlueprints.find(b => b.id === id)
      if (!blueprint) throw new Error('Blueprint not found')
      return { ...blueprint }
    },

    async getBlueprintsByCategory(category: string): Promise<Blueprint[]> {
      await delay(300)
      return mockBlueprints.filter(b => b.category === category)
    },

    async validateBlueprintConfig(id: string, config: ProjectConfig): Promise<{ valid: boolean; errors?: string[] }> {
      await delay(250)
      return { valid: true }
    }
  },

  projects: {
    async generatePreview(config: ProjectConfig): Promise<PreviewResponse> {
      await delay(800)
      return generateMockPreview(config)
    },

    async generateProject(request: GenerationRequest): Promise<GenerationResponse> {
      await delay(2000)
      return {
        projectId: 'mock-project-123',
        downloadUrl: '/api/download/mock-project-123',
        fileCount: 15,
        estimatedSize: '25.4 KB'
      }
    },

    async generateAndDownload(request: GenerationRequest): Promise<Blob> {
      await delay(3000)
      // Create a fake ZIP blob
      const content = 'Mock ZIP file content'
      return new Blob([content], { type: 'application/zip' })
    },

    async downloadProject(token: string): Promise<Blob> {
      await delay(1000)
      const content = 'Mock ZIP file content'
      return new Blob([content], { type: 'application/zip' })
    },

    async getDownloadStatus(token: string): Promise<{ status: string; progress?: number }> {
      await delay(100)
      return { status: 'ready', progress: 100 }
    }
  },

  ws: {
    connect(): Promise<void> {
      console.log('Mock WebSocket: Connected')
      return Promise.resolve()
    },

    disconnect(): void {
      console.log('Mock WebSocket: Disconnected')
    },

    send(message: unknown): void {
      console.log('Mock WebSocket: Sent', message)
    },

    subscribe(event: string, callback: (data: unknown) => void): () => void {
      console.log('Mock WebSocket: Subscribed to', event)
      return () => console.log('Mock WebSocket: Unsubscribed from', event)
    }
  }
}

export default mockApi