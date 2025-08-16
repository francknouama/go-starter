// Type definitions for go-starter Web UI
// These types need to be consistent with the API service types

export type DisclosureMode = 'basic' | 'advanced'

// Additional types for the UI
export type ProjectType = 'web-api' | 'cli' | 'library' | 'lambda' | 'lambda-proxy' | 'microservice' | 'monolith' | 'workspace' | 'event-driven'
export type Architecture = 'standard' | 'clean' | 'ddd' | 'hexagonal' | 'event-driven'
export type Framework = 'gin' | 'echo' | 'fiber' | 'chi' | 'cobra' | 'gorm' | 'sqlx' | 'ent'
export type LoggerType = 'slog' | 'zap' | 'logrus' | 'zerolog'
export type DatabaseDriver = 'postgres' | 'mysql' | 'sqlite' | 'mongodb' | 'redis'
export type DatabaseORM = 'gorm' | 'sqlx' | 'sqlc' | 'ent'
export type AuthType = 'jwt' | 'oauth2' | 'session' | 'api-key'

// Re-export types from API service for convenience
export type {
  ProjectConfig,
  Blueprint,
  PreviewResponse,
  GenerationRequest,
  GenerationResponse,
  FileNode
} from './services/api'