package models

import (
	"time"
)

// ProjectConfig represents the web UI project configuration
type ProjectConfig struct {
	ProjectName  string            `json:"project_name" binding:"required"`
	ModuleURL    string            `json:"module_url" binding:"required"`
	GoVersion    string            `json:"go_version" binding:"required"`
	ProjectType  string            `json:"project_type" binding:"required"`
	Framework    string            `json:"framework"`
	Architecture string            `json:"architecture"`
	Logger       string            `json:"logger"`
	Database     *DatabaseConfig   `json:"database,omitempty"`
	Auth         *AuthConfig       `json:"authentication,omitempty"`
	Deployment   *DeploymentConfig `json:"deployment,omitempty"`
	Features     *FeaturesConfig   `json:"features,omitempty"`
}

type DatabaseConfig struct {
	Driver string `json:"driver"`
	ORM    string `json:"orm"`
}

type AuthConfig struct {
	Type      string   `json:"type"`
	Providers []string `json:"providers,omitempty"`
}

type DeploymentConfig struct {
	Targets       []string `json:"targets"`
	CloudProvider string   `json:"cloud_provider,omitempty"`
}

type FeaturesConfig struct {
	Testing    bool `json:"testing"`
	Monitoring bool `json:"monitoring"`
	Logging    bool `json:"logging"`
	Caching    bool `json:"caching"`
}

// Generation request/response models

type ValidateConfigRequest struct {
	Config ProjectConfig `json:"config" binding:"required"`
}

type ValidateConfigResponse struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Severity string `json:"severity"` // "error" or "warning"
}

type GenerateProjectRequest struct {
	Blueprint string        `json:"blueprint" binding:"required"`
	Config    ProjectConfig `json:"config" binding:"required"`
	Options   GenerationOptions `json:"options"`
}

type GenerationOptions struct {
	MemoryMode      bool `json:"memory_mode"`
	IncludeExamples bool `json:"include_examples"`
}

type GenerateProjectResponse struct {
	ID             string                `json:"id"`
	Status         string                `json:"status"`
	FilesGenerated int                   `json:"files_generated"`
	GenerationTime string                `json:"generation_time"`
	DownloadURL    string                `json:"download_url"`
	ExpiresAt      string                `json:"expires_at"`
	Files          []GeneratedFileInfo   `json:"files"`
}

type GeneratedFileInfo struct {
	Path string `json:"path"`
	Size int    `json:"size"`
	Type string `json:"type"`
}

// Internal models for storing generated projects

type GeneratedFile struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Size     int    `json:"size"`
	Type     string `json:"type"`
	Language string `json:"language,omitempty"`
}

type GeneratedProject struct {
	ID             string          `json:"id"`
	Config         ProjectConfig   `json:"config"`
	Files          []GeneratedFile `json:"files"`
	ZipData        []byte          `json:"-"` // Don't serialize in JSON
	GenerationTime time.Duration   `json:"generation_time"`
	CreatedAt      time.Time       `json:"created_at"`
	ExpiresAt      time.Time       `json:"expires_at"`
}

// WebSocket message types

type PreviewUpdate struct {
	Type     string `json:"type"` // "file_added", "file_updated", "error", "complete"
	Path     string `json:"path,omitempty"`
	Content  string `json:"content,omitempty"`
	Error    string `json:"error,omitempty"`
	Progress int    `json:"progress,omitempty"`
}

type GenerationStatus struct {
	ID             string `json:"id"`
	Status         string `json:"status"` // "pending", "generating", "completed", "error"
	Progress       int    `json:"progress"`
	FilesGenerated int    `json:"files_generated"`
	TotalFiles     int    `json:"total_files"`
	CurrentFile    string `json:"current_file,omitempty"`
	Error          string `json:"error,omitempty"`
}