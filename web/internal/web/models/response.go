package models

import "time"

// APIResponse is a generic API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents an API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta contains metadata about the response
type Meta struct {
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"requestId,omitempty"`
	Version   string    `json:"version,omitempty"`
}

// GenerationResponse represents the response from project generation
type GenerationResponse struct {
	ProjectName   string           `json:"projectName"`
	ModulePath    string           `json:"modulePath"`
	Type          string           `json:"type"`
	Architecture  string           `json:"architecture"`
	Framework     string           `json:"framework"`
	Files         []GeneratedFile  `json:"files"`
	DownloadURL   string           `json:"downloadUrl,omitempty"`   // For ZIP downloads
	DownloadToken string           `json:"downloadToken,omitempty"` // Token for secure downloads
	ExpiresAt     *time.Time       `json:"expiresAt,omitempty"`     // Download expiration
	Stats         *GenerationStats `json:"stats,omitempty"`
}

// GeneratedFile represents a generated file in the project
type GeneratedFile struct {
	Path     string `json:"path"`
	Content  string `json:"content,omitempty"` // Only included in preview mode
	Size     int64  `json:"size"`
	IsDir    bool   `json:"isDir"`
	Mode     string `json:"mode,omitempty"`    // File permissions
	ModTime  string `json:"modTime,omitempty"` // Modification time
}

// GenerationStats contains statistics about the generation process
type GenerationStats struct {
	TotalFiles     int           `json:"totalFiles"`
	TotalSize      int64         `json:"totalSize"`
	GenerationTime time.Duration `json:"generationTime"`
	BlueprintUsed  string        `json:"blueprintUsed"`
	TemplateCount  int           `json:"templateCount"`
}

// PreviewResponse represents the response from project preview
type PreviewResponse struct {
	ProjectName  string          `json:"projectName"`
	ModulePath   string          `json:"modulePath"`
	Type         string          `json:"type"`
	Architecture string          `json:"architecture"`
	Framework    string          `json:"framework"`
	FileTree     *FileTreeNode   `json:"fileTree"`
	Files        []GeneratedFile `json:"files"`
	Stats        *PreviewStats   `json:"stats"`
}

// FileTreeNode represents a node in the file tree structure
type FileTreeNode struct {
	Name     string          `json:"name"`
	Path     string          `json:"path"`
	IsDir    bool            `json:"isDir"`
	Size     int64           `json:"size,omitempty"`
	Children []*FileTreeNode `json:"children,omitempty"`
}

// PreviewStats contains statistics about the preview
type PreviewStats struct {
	TotalFiles     int    `json:"totalFiles"`
	TotalDirs      int    `json:"totalDirs"`
	EstimatedSize  int64  `json:"estimatedSize"`
	BlueprintUsed  string `json:"blueprintUsed"`
	TemplateCount  int    `json:"templateCount"`
	Complexity     string `json:"complexity"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	System    *SystemInfo       `json:"system,omitempty"`
}

// SystemInfo contains system information
type SystemInfo struct {
	GoVersion    string  `json:"goVersion"`
	OS           string  `json:"os"`
	Arch         string  `json:"arch"`
	NumCPU       int     `json:"numCPU"`
	MemoryUsage  int64   `json:"memoryUsage"`  // in bytes
	GoroutineCount int   `json:"goroutineCount"`
	Uptime       string  `json:"uptime"`
}

// ConfigResponse represents the configuration response
type ConfigResponse struct {
	ProjectTypes   []ProjectTypeInfo   `json:"projectTypes"`
	Architectures  []ArchitectureInfo  `json:"architectures"`
	Frameworks     []FrameworkInfo     `json:"frameworks"`
	Loggers        []LoggerInfo        `json:"loggers"`
	Databases      []DatabaseInfo      `json:"databases"`
	Complexities   []ComplexityInfo    `json:"complexities"`
	DefaultValues  *DefaultConfig      `json:"defaultValues"`
}

// ProjectTypeInfo contains information about a project type
type ProjectTypeInfo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Icon         string   `json:"icon,omitempty"`
	Category     string   `json:"category,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Complexity   string   `json:"complexity,omitempty"`
	EstimatedTime string  `json:"estimatedTime,omitempty"`
}

// ArchitectureInfo contains information about an architecture pattern
type ArchitectureInfo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Complexity   string   `json:"complexity"`
	UseCase      string   `json:"useCase,omitempty"`
	Benefits     []string `json:"benefits,omitempty"`
	Drawbacks    []string `json:"drawbacks,omitempty"`
}

// FrameworkInfo contains information about a framework
type FrameworkInfo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ProjectTypes []string `json:"projectTypes"`
	Performance  string   `json:"performance,omitempty"`
	Popularity   string   `json:"popularity,omitempty"`
}

// LoggerInfo contains information about a logger
type LoggerInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Performance  string `json:"performance,omitempty"`
	Features     string `json:"features,omitempty"`
	Package      string `json:"package,omitempty"`
}

// DatabaseInfo contains information about a database driver
type DatabaseInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"` // sql, nosql, cache
	ORMs        []string `json:"orms,omitempty"`
}

// ComplexityInfo contains information about complexity levels
type ComplexityInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	FileCount    string `json:"fileCount,omitempty"`
	Features     string `json:"features,omitempty"`
	Recommended  string `json:"recommended,omitempty"`
}

// DefaultConfig contains default configuration values
type DefaultConfig struct {
	GoVersion    string `json:"goVersion"`
	Logger       string `json:"logger"`
	Architecture string `json:"architecture"`
	Complexity   string `json:"complexity"`
	Framework    string `json:"framework,omitempty"`
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"requestId,omitempty"`
}

// WebSocket message types
const (
	WSMessageTypeProgress    = "progress"
	WSMessageTypeComplete    = "complete"
	WSMessageTypeError       = "error"
	WSMessageTypePreview     = "preview"
	WSMessageTypeStatus      = "status"
)

// ProgressData represents progress information sent via WebSocket
type ProgressData struct {
	Stage       string  `json:"stage"`
	Progress    float64 `json:"progress"`    // 0.0 to 1.0
	Message     string  `json:"message"`
	CurrentFile string  `json:"currentFile,omitempty"`
	TotalFiles  int     `json:"totalFiles,omitempty"`
	ProcessedFiles int  `json:"processedFiles,omitempty"`
}

// NewSuccessResponse creates a successful API response
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Timestamp: time.Now(),
		},
	}
}

// NewErrorResponse creates an error API response
func NewErrorResponse(code, message, details string) *APIResponse {
	return &APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: &Meta{
			Timestamp: time.Now(),
		},
	}
}

// NewValidationErrorResponse creates a validation error response
func NewValidationErrorResponse(message string) *APIResponse {
	return NewErrorResponse("VALIDATION_ERROR", "Validation failed", message)
}

// NewInternalErrorResponse creates an internal server error response
func NewInternalErrorResponse(message string) *APIResponse {
	return NewErrorResponse("INTERNAL_ERROR", "Internal server error", message)
}

// NewNotFoundErrorResponse creates a not found error response
func NewNotFoundErrorResponse(resource string) *APIResponse {
	return NewErrorResponse("NOT_FOUND", "Resource not found", resource)
}