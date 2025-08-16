package handlers

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/web/internal/web/models"
)

// Cache for temporary ZIP files
var zipCache = make(map[string]*CachedZip)

// CachedZip represents a cached ZIP file
type CachedZip struct {
	Data      []byte
	ExpiresAt time.Time
	FileName  string
}

// GenerateProject generates a project and returns it as a ZIP file or file structure
func GenerateProject(c *gin.Context) {
	var request models.GenerationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response := models.NewValidationErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate request
	if err := request.Validate(); err != nil {
		response := models.NewValidationErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Set defaults
	request.SetDefaults()

	// Convert to internal config
	config := convertToProjectConfig(request)

	// Initialize generator
	gen := generator.New()

	// Initialize template registry
	registry := templates.NewRegistry()

	// Determine blueprint ID
	blueprintID := determineBlueprintID(config)

	// Get template
	tmpl, err := registry.Get(blueprintID)
	if err != nil {
		// Try fallback: look for templates by type
		templatesByType := registry.GetByType(config.Type)
		if len(templatesByType) > 0 {
			tmpl = templatesByType[0]
		} else {
			response := models.NewErrorResponse(
				"BLUEPRINT_NOT_FOUND",
				"Blueprint not found",
				fmt.Sprintf("Blueprint '%s' not found", blueprintID),
			)
			c.JSON(http.StatusNotFound, response)
			return
		}
	}

	startTime := time.Now()

	// Handle dry run mode
	if request.DryRun {
		handleDryRun(c, gen, config, tmpl, startTime)
		return
	}

	// Generate project files in memory
	filesMap, err := gen.GenerateInMemory(&config, tmpl.ID)
	if err != nil {
		response := models.NewInternalErrorResponse(fmt.Sprintf("Failed to generate project: %v", err))
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Handle different output formats
	switch request.OutputFormat {
	case "zip":
		handleZipDownload(c, filesMap, config, tmpl, startTime)
	case "files":
		handleFileStructure(c, filesMap, config, tmpl, startTime)
	default:
		// Default to ZIP
		handleZipDownload(c, filesMap, config, tmpl, startTime)
	}
}

// handleDryRun handles dry run mode
func handleDryRun(c *gin.Context, gen *generator.Generator, config types.ProjectConfig, tmpl types.Template, startTime time.Time) {
	// Generate preview using existing preview functionality
	previewData, err := generatePreview(gen, config, tmpl)
	if err != nil {
		response := models.NewInternalErrorResponse(err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Create dry run response
	stats := &models.GenerationStats{
		TotalFiles:     previewData.Stats.TotalFiles,
		TotalSize:      previewData.Stats.EstimatedSize,
		GenerationTime: time.Since(startTime),
		BlueprintUsed:  tmpl.ID,
		TemplateCount:  len(tmpl.Files),
	}

	dryRunResponse := &models.GenerationResponse{
		ProjectName:  config.Name,
		ModulePath:   config.Module,
		Type:         config.Type,
		Architecture: config.Architecture,
		Framework:    config.Framework,
		Files:        previewData.Files,
		Stats:        stats,
	}

	response := models.NewSuccessResponse(dryRunResponse)
	c.JSON(http.StatusOK, response)
}

// handleZipDownload creates a ZIP file and returns download URL
func handleZipDownload(c *gin.Context, filesMap map[string][]byte, config types.ProjectConfig, tmpl types.Template, startTime time.Time) {
	// Create ZIP file in memory
	zipData, err := createZipFromFiles(filesMap, config.Name)
	if err != nil {
		response := models.NewInternalErrorResponse(fmt.Sprintf("Failed to create ZIP: %v", err))
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Generate unique token for download
	token := generateDownloadToken()
	expiresAt := time.Now().Add(30 * time.Minute) // 30 minutes expiry

	// Cache the ZIP file
	fileName := fmt.Sprintf("%s.zip", config.Name)
	zipCache[token] = &CachedZip{
		Data:      zipData,
		ExpiresAt: expiresAt,
		FileName:  fileName,
	}

	// Clean up expired cache entries
	go cleanupExpiredCache()

	// Create response
	generationTime := time.Since(startTime)
	stats := &models.GenerationStats{
		TotalFiles:     len(filesMap),
		TotalSize:      int64(len(zipData)),
		GenerationTime: generationTime,
		BlueprintUsed:  tmpl.ID,
		TemplateCount:  len(tmpl.Files),
	}

	response := &models.GenerationResponse{
		ProjectName:   config.Name,
		ModulePath:    config.Module,
		Type:          config.Type,
		Architecture:  config.Architecture,
		Framework:     config.Framework,
		Files:         convertFilesMapToGenerated(filesMap),
		DownloadURL:   fmt.Sprintf("/api/v1/download/%s", token),
		DownloadToken: token,
		ExpiresAt:     &expiresAt,
		Stats:         stats,
	}

	apiResponse := models.NewSuccessResponse(response)
	c.JSON(http.StatusOK, apiResponse)
}

// handleFileStructure returns the file structure without ZIP
func handleFileStructure(c *gin.Context, filesMap map[string][]byte, config types.ProjectConfig, tmpl types.Template, startTime time.Time) {
	generationTime := time.Since(startTime)
	
	// Convert files to response format
	files := convertFilesMapToGenerated(filesMap)
	
	// Calculate total size
	totalSize := int64(0)
	for _, file := range files {
		totalSize += file.Size
	}

	stats := &models.GenerationStats{
		TotalFiles:     len(files),
		TotalSize:      totalSize,
		GenerationTime: generationTime,
		BlueprintUsed:  tmpl.ID,
		TemplateCount:  len(tmpl.Files),
	}

	response := &models.GenerationResponse{
		ProjectName:  config.Name,
		ModulePath:   config.Module,
		Type:         config.Type,
		Architecture: config.Architecture,
		Framework:    config.Framework,
		Files:        files,
		Stats:        stats,
	}

	apiResponse := models.NewSuccessResponse(response)
	c.JSON(http.StatusOK, apiResponse)
}

// DownloadZip handles ZIP file downloads
func DownloadZip(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Download token is required"})
		return
	}

	// Get cached ZIP
	cachedZip, exists := zipCache[token]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Download not found or expired"})
		return
	}

	// Check expiry
	if time.Now().After(cachedZip.ExpiresAt) {
		delete(zipCache, token)
		c.JSON(http.StatusGone, gin.H{"error": "Download has expired"})
		return
	}

	// Set headers for file download
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", cachedZip.FileName))
	c.Header("Content-Length", fmt.Sprintf("%d", len(cachedZip.Data)))
	c.Header("Cache-Control", "no-cache")

	// Send file
	c.Data(http.StatusOK, "application/zip", cachedZip.Data)

	// Clean up token after successful download
	delete(zipCache, token)
}

// GenerateAndDownloadProject generates and immediately downloads a project as ZIP
func GenerateAndDownloadProject(c *gin.Context) {
	var request models.GenerationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response := models.NewValidationErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate request
	if err := request.Validate(); err != nil {
		response := models.NewValidationErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Set defaults and force ZIP output
	request.SetDefaults()
	request.OutputFormat = "zip"

	// Convert to internal config
	config := convertToProjectConfig(request)

	// Initialize generator and registry
	gen := generator.New()
	registry := templates.NewRegistry()

	// Determine blueprint ID and get template
	blueprintID := determineBlueprintID(config)
	tmpl, err := registry.Get(blueprintID)
	if err != nil {
		templatesByType := registry.GetByType(config.Type)
		if len(templatesByType) > 0 {
			tmpl = templatesByType[0]
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Blueprint not found"})
			return
		}
	}

	// Generate project files in memory
	filesMap, err := gen.GenerateInMemory(&config, tmpl.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate project"})
		return
	}

	// Create ZIP file
	zipData, err := createZipFromFiles(filesMap, config.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ZIP"})
		return
	}

	// Set headers and send ZIP directly
	fileName := fmt.Sprintf("%s.zip", config.Name)
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Header("Content-Length", fmt.Sprintf("%d", len(zipData)))
	c.Header("Cache-Control", "no-cache")

	c.Data(http.StatusOK, "application/zip", zipData)
}

// createZipFromFiles creates a ZIP archive from a file map
func createZipFromFiles(filesMap map[string][]byte, projectName string) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// Sort files for consistent ZIP structure
	var filePaths []string
	for path := range filesMap {
		filePaths = append(filePaths, path)
	}
	sort.Strings(filePaths)

	// Add files to ZIP
	for _, filePath := range filePaths {
		content := filesMap[filePath]
		
		// Create ZIP entry with project name as root directory
		zipPath := filepath.Join(projectName, filePath)
		
		// Ensure we use forward slashes in ZIP paths
		zipPath = strings.ReplaceAll(zipPath, string(os.PathSeparator), "/")
		
		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("failed to create ZIP entry for %s: %w", zipPath, err)
		}

		if _, err := writer.Write(content); err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("failed to write content for %s: %w", zipPath, err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close ZIP writer: %w", err)
	}

	return buf.Bytes(), nil
}

// convertFilesMapToGenerated converts a file map to GeneratedFile slice
func convertFilesMapToGenerated(filesMap map[string][]byte) []models.GeneratedFile {
	var files []models.GeneratedFile
	
	for path, content := range filesMap {
		file := models.GeneratedFile{
			Path:    path,
			Content: string(content),
			Size:    int64(len(content)),
			IsDir:   false,
			Mode:    "0644",
			ModTime: time.Now().Format(time.RFC3339),
		}
		files = append(files, file)
	}

	// Sort files by path
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	return files
}

// generateDownloadToken generates a secure download token
func generateDownloadToken() string {
	// Generate UUID
	id := uuid.New()
	
	// Add timestamp for uniqueness
	timestamp := time.Now().Unix()
	
	// Create hash
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%d", id.String(), timestamp)))
	
	// Return first 32 characters of hex
	return fmt.Sprintf("%x", hash)[:32]
}

// cleanupExpiredCache removes expired cache entries
func cleanupExpiredCache() {
	now := time.Now()
	for token, cached := range zipCache {
		if now.After(cached.ExpiresAt) {
			delete(zipCache, token)
		}
	}
}

// GetDownloadStatus returns the status of a download token
func GetDownloadStatus(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		response := models.NewValidationErrorResponse("Download token is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	cachedZip, exists := zipCache[token]
	if !exists {
		response := models.NewNotFoundErrorResponse("Download not found")
		c.JSON(http.StatusNotFound, response)
		return
	}

	// Check if expired
	if time.Now().After(cachedZip.ExpiresAt) {
		delete(zipCache, token)
		response := models.NewErrorResponse("DOWNLOAD_EXPIRED", "Download has expired", "")
		c.JSON(http.StatusGone, response)
		return
	}

	// Return status
	status := map[string]interface{}{
		"token":     token,
		"fileName":  cachedZip.FileName,
		"size":      len(cachedZip.Data),
		"expiresAt": cachedZip.ExpiresAt,
		"status":    "ready",
	}

	response := models.NewSuccessResponse(status)
	c.JSON(http.StatusOK, response)
}