// Core types for the Go Starter Web UI

export type DisclosureMode = 'basic' | 'advanced'

export type ProjectType = 
  | 'cli' 
  | 'web-api' 
  | 'library' 
  | 'lambda' 
  | 'lambda-proxy'
  | 'event-driven'
  | 'microservice'
  | 'monolith'
  | 'workspace'

export type Architecture = 
  | 'standard'
  | 'clean'
  | 'ddd'
  | 'hexagonal'
  | 'event-driven'

export type Framework = 
  | 'gin'
  | 'echo'
  | 'fiber'
  | 'chi'
  | 'cobra'

export type LoggerType = 
  | 'slog'
  | 'zap'
  | 'logrus'
  | 'zerolog'

export type DatabaseDriver = 
  | 'postgres'
  | 'mysql'
  | 'mongodb'
  | 'sqlite'
  | 'redis'

export type DatabaseORM = 
  | 'gorm'
  | 'sqlx'
  | 'sqlc'
  | 'ent'

export type AuthType = 
  | 'jwt'
  | 'oauth2'
  | 'session'
  | 'api-key'

export type CloudProvider = 
  | 'aws'
  | 'gcp'
  | 'azure'

export interface ProjectConfig {
  // Basic configuration
  projectName: string
  moduleUrl: string
  goVersion: string
  projectType: ProjectType
  
  // Framework and architecture
  framework: Framework
  architecture: Architecture
  logger: LoggerType
  
  // Advanced configuration
  features?: {
    database?: {
      driver: DatabaseDriver
      orm: DatabaseORM
    }
    authentication?: {
      type: AuthType
      providers: string[]
    }
    testing?: {
      framework: string
      coverage: boolean
    }
    deployment?: {
      targets: string[]
    }
    advanced?: {
      [key: string]: {
        enabled: boolean
        options?: { [key: string]: any }
      }
    }
    blueprintOptions?: {
      [key: string]: any
    }
  }
}

export interface Blueprint {
  id: string
  name: string
  description: string
  type: ProjectType
  architecture: Architecture
  complexity: 'simple' | 'standard' | 'advanced' | 'expert'
  fileCount: number
  dependencies: string[]
  features: string[]
}

export interface GenerationStatus {
  id: string
  status: 'pending' | 'generating' | 'completed' | 'error'
  progress: number
  filesGenerated: number
  totalFiles: number
  currentFile?: string
  error?: string
  generationTime?: string
}

export interface GeneratedFile {
  path: string
  content: string
  size: number
  type: string
  language?: string
}

export interface ProjectGeneration {
  id: string
  config: ProjectConfig
  status: GenerationStatus
  files: GeneratedFile[]
  downloadUrl?: string
  expiresAt?: string
}

export interface PreviewUpdate {
  type: 'file_added' | 'file_updated' | 'error' | 'complete'
  path?: string
  content?: string
  error?: string
  progress?: number
}

// WebSocket Message Types
export interface WSMessage {
  type: WSMessageType
  data: any
  timestamp: string
  requestId?: string
}

export type WSMessageType = 
  | 'progress'
  | 'complete'
  | 'error'
  | 'preview'
  | 'status'
  | 'file_tree'
  | 'file_content'
  | 'preview_start'
  | 'preview_complete'

export interface WSProgressData {
  stage: string
  progress: number // 0.0 to 1.0
  message: string
  currentFile?: string
  totalFiles?: number
  processedFiles?: number
}

export interface WSStatusData {
  message: string
  timestamp: string
}

export interface WSFileTreeNode {
  name: string
  path: string
  isDir: boolean
  size?: number
  children?: WSFileTreeNode[]
}

export interface WSFileContent {
  path: string
  content: string
  size: number
  isDir: boolean
  mode?: string
  modTime?: string
}

export interface WSPreviewCompleteData {
  projectName: string
  modulePath: string
  type: string
  architecture: string
  framework: string
  totalFiles: number
  blueprintUsed: string
}

export interface ValidationError {
  field: string
  message: string
  severity: 'error' | 'warning'
}

export interface BlueprintListResponse {
  blueprints: Blueprint[]
}

export interface GenerateProjectRequest {
  blueprint: string
  config: ProjectConfig
  options: {
    memoryMode: boolean
    includeExamples: boolean
  }
}

export interface GenerateProjectResponse {
  id: string
  status: 'completed' | 'error'
  filesGenerated: number
  generationTime: string
  downloadUrl: string
  expiresAt: string
  files: Array<{
    path: string
    size: number
    type: string
  }>
}

// WebSocket Connection State
export interface WSConnectionState {
  connected: boolean
  connecting: boolean
  error?: string
  lastReconnectAttempt?: number
  reconnectAttempts: number
}

// Real-time Preview State
export interface RealtimePreviewState {
  isGenerating: boolean
  progress: WSProgressData | null
  fileTree: WSFileTreeNode | null
  selectedFile: WSFileContent | null
  files: Map<string, WSFileContent>
  status: string
  error?: string
}