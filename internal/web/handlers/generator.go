package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/web/models"
	"github.com/francknouama/go-starter/pkg/types"
)

type GeneratorHandler struct {
	// In-memory storage for generated projects (in production, use Redis or similar)
	projects map[string]*models.GeneratedProject
	mutex    sync.RWMutex
}

func NewGeneratorHandler() *GeneratorHandler {
	handler := &GeneratorHandler{
		projects: make(map[string]*models.GeneratedProject),
	}

	// Start cleanup goroutine
	go handler.cleanupExpiredProjects()

	return handler
}

// ValidateConfig validates project configuration
func (h *GeneratorHandler) ValidateConfig(c *gin.Context) {
	var req models.ValidateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "INVALID_REQUEST",
		})
		return
	}

	errors := validateProjectConfig(req.Config)
	
	c.JSON(http.StatusOK, models.ValidateConfigResponse{
		Valid:  len(errors) == 0,
		Errors: errors,
	})
}

// GenerateProject generates a new project
func (h *GeneratorHandler) GenerateProject(c *gin.Context) {
	var req models.GenerateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "INVALID_REQUEST",
		})
		return
	}

	// Validate configuration
	if errors := validateProjectConfig(req.Config); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Configuration validation failed",
			"code":   "VALIDATION_FAILED",
			"errors": errors,
		})
		return
	}

	// Generate unique ID for this generation
	generationID := uuid.New().String()

	// Convert web config to internal config
	config := &types.ProjectConfig{
		Name:         req.Config.ProjectName,
		Module:       req.Config.ModuleURL,
		Type:         req.Config.ProjectType,
		Framework:    req.Config.Framework,
		Architecture: req.Config.Architecture,
		Logger:       req.Config.Logger,
		GoVersion:    req.Config.GoVersion,
	}

	// Generate project in memory
	startTime := time.Now()
	gen := generator.New()
	
	// For web mode, we generate to a temporary in-memory buffer
	files, err := gen.GenerateInMemory(config, req.Blueprint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate project",
			"code":  "GENERATION_FAILED",
		})
		return
	}

	generationTime := time.Since(startTime)

	// Create ZIP archive
	zipBuffer, err := createZipArchive(files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create archive",
			"code":  "ARCHIVE_FAILED",
		})
		return
	}

	// Store generated project
	project := &models.GeneratedProject{
		ID:             generationID,
		Config:         req.Config,
		Files:          convertToWebFiles(files),
		ZipData:        zipBuffer.Bytes(),
		GenerationTime: generationTime,
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(24 * time.Hour), // Expire after 24 hours
	}

	h.mutex.Lock()
	h.projects[generationID] = project
	h.mutex.Unlock()

	// Return response
	fileList := make([]models.GeneratedFileInfo, 0, len(files))
	for path, content := range files {
		fileList = append(fileList, models.GeneratedFileInfo{
			Path: path,
			Size: len(content),
			Type: getFileType(path),
		})
	}

	c.JSON(http.StatusOK, models.GenerateProjectResponse{
		ID:             generationID,
		Status:         "completed",
		FilesGenerated: len(files),
		GenerationTime: generationTime.String(),
		DownloadURL:    fmt.Sprintf("/api/v1/download/%s", generationID),
		ExpiresAt:      project.ExpiresAt.Format(time.RFC3339),
		Files:          fileList,
	})
}

// DownloadProject handles project download
func (h *GeneratorHandler) DownloadProject(c *gin.Context) {
	projectID := c.Param("id")

	h.mutex.RLock()
	project, exists := h.projects[projectID]
	h.mutex.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Project not found or expired",
			"code":  "PROJECT_NOT_FOUND",
		})
		return
	}

	// Check if project has expired
	if time.Now().After(project.ExpiresAt) {
		h.mutex.Lock()
		delete(h.projects, projectID)
		h.mutex.Unlock()

		c.JSON(http.StatusNotFound, gin.H{
			"error": "Project has expired",
			"code":  "PROJECT_EXPIRED",
		})
		return
	}

	// Set headers for file download
	filename := fmt.Sprintf("%s.zip", project.Config.ProjectName)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Length", fmt.Sprintf("%d", len(project.ZipData)))

	// Send the ZIP file
	c.Data(http.StatusOK, "application/zip", project.ZipData)
}

// CleanupProject manually cleans up a project
func (h *GeneratorHandler) CleanupProject(c *gin.Context) {
	projectID := c.Param("id")

	h.mutex.Lock()
	delete(h.projects, projectID)
	h.mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Project cleaned up successfully",
	})
}

// Helper functions

func validateProjectConfig(config models.ProjectConfig) []models.ValidationError {
	var errors []models.ValidationError

	if config.ProjectName == "" {
		errors = append(errors, models.ValidationError{
			Field:    "project_name",
			Message:  "Project name is required",
			Severity: "error",
		})
	}

	if config.ModuleURL == "" {
		errors = append(errors, models.ValidationError{
			Field:    "module_url",
			Message:  "Module URL is required",
			Severity: "error",
		})
	}

	if config.ProjectType == "" {
		errors = append(errors, models.ValidationError{
			Field:    "project_type",
			Message:  "Project type is required",
			Severity: "error",
		})
	}

	return errors
}

func createZipArchive(files map[string][]byte) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	for path, content := range files {
		f, err := writer.Create(path)
		if err != nil {
			return nil, err
		}
		
		_, err = f.Write(content)
		if err != nil {
			return nil, err
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func convertToWebFiles(files map[string][]byte) []models.GeneratedFile {
	result := make([]models.GeneratedFile, 0, len(files))
	
	for path, content := range files {
		result = append(result, models.GeneratedFile{
			Path:     path,
			Content:  string(content),
			Size:     len(content),
			Type:     getFileType(path),
			Language: getLanguage(path),
		})
	}

	return result
}

func getFileType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".go":
		return "go"
	case ".yaml", ".yml":
		return "yaml"
	case ".json":
		return "json"
	case ".md":
		return "markdown"
	case ".toml":
		return "toml"
	case ".sh":
		return "shell"
	case ".sql":
		return "sql"
	default:
		return "text"
	}
}

func getLanguage(path string) string {
	return getFileType(path)
}

// cleanupExpiredProjects removes expired projects from memory
func (h *GeneratorHandler) cleanupExpiredProjects() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		
		h.mutex.Lock()
		for id, project := range h.projects {
			if now.After(project.ExpiresAt) {
				delete(h.projects, id)
			}
		}
		h.mutex.Unlock()
	}
}