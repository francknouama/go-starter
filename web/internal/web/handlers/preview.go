package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/web/internal/web/models"
)

// PreviewProject generates a preview of the project structure without creating files
func PreviewProject(c *gin.Context) {
	var request models.PreviewRequest
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
	config := convertToProjectConfig(request.GenerationRequest)

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

	// Generate preview
	previewData, err := generatePreview(gen, config, tmpl)
	if err != nil {
		response := models.NewInternalErrorResponse(err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Build response
	previewResponse := &models.PreviewResponse{
		ProjectName:  config.Name,
		ModulePath:   config.Module,
		Type:         config.Type,
		Architecture: config.Architecture,
		Framework:    config.Framework,
		FileTree:     buildFileTree(previewData.Files),
		Files:        previewData.Files,
		Stats:        previewData.Stats,
	}

	response := models.NewSuccessResponse(previewResponse)
	c.JSON(http.StatusOK, response)
}

// generatePreview generates preview data for a project
func generatePreview(gen *generator.Generator, config types.ProjectConfig, tmpl types.Template) (*PreviewData, error) {
	// Generate files in memory to get the structure
	filesMap, err := gen.GenerateInMemory(&config, tmpl.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate preview: %w", err)
	}

	// Convert to preview format
	var files []models.GeneratedFile
	totalSize := int64(0)
	fileCount := 0
	dirCount := 0

	// Track directories
	dirs := make(map[string]bool)

	for filePath, content := range filesMap {
		// Add parent directories
		dir := filepath.Dir(filePath)
		for dir != "." && dir != "/" {
			dirs[dir] = true
			dir = filepath.Dir(dir)
		}

		// Create file info
		file := models.GeneratedFile{
			Path:    filePath,
			Content: string(content), // Include content for preview
			Size:    int64(len(content)),
			IsDir:   false,
			Mode:    "0644",
			ModTime: time.Now().Format(time.RFC3339),
		}

		files = append(files, file)
		totalSize += file.Size
		fileCount++
	}

	// Add directory entries
	for dir := range dirs {
		files = append(files, models.GeneratedFile{
			Path:    dir,
			Content: "",
			Size:    0,
			IsDir:   true,
			Mode:    "0755",
			ModTime: time.Now().Format(time.RFC3339),
		})
		dirCount++
	}

	// Sort files by path
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	// Create stats
	stats := &models.PreviewStats{
		TotalFiles:     fileCount,
		TotalDirs:      dirCount,
		EstimatedSize:  totalSize,
		BlueprintUsed:  tmpl.ID,
		TemplateCount:  len(tmpl.Files),
		Complexity:     config.Variables["complexity"],
	}

	if stats.Complexity == "" {
		stats.Complexity = "standard"
	}

	return &PreviewData{
		Files: files,
		Stats: stats,
	}, nil
}

// buildFileTree builds a hierarchical file tree structure
func buildFileTree(files []models.GeneratedFile) *models.FileTreeNode {
	root := &models.FileTreeNode{
		Name:     "root",
		Path:     "",
		IsDir:    true,
		Children: []*models.FileTreeNode{},
	}

	// Build tree structure
	for _, file := range files {
		if file.Path == "" {
			continue
		}

		parts := strings.Split(filepath.Clean(file.Path), string(filepath.Separator))
		current := root

		// Navigate/create path to file
		currentPath := ""
		for i, part := range parts {
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = filepath.Join(currentPath, part)
			}

			// Look for existing child
			var found *models.FileTreeNode
			for _, child := range current.Children {
				if child.Name == part {
					found = child
					break
				}
			}

			if found == nil {
				// Create new node
				isDir := i < len(parts)-1 || file.IsDir
				newNode := &models.FileTreeNode{
					Name:     part,
					Path:     currentPath,
					IsDir:    isDir,
					Children: []*models.FileTreeNode{},
				}

				if !isDir {
					newNode.Size = file.Size
				}

				current.Children = append(current.Children, newNode)
				current = newNode
			} else {
				current = found
			}
		}
	}

	// Sort children recursively
	sortFileTreeNode(root)

	return root
}

// sortFileTreeNode sorts the children of a file tree node recursively
func sortFileTreeNode(node *models.FileTreeNode) {
	if node.Children == nil {
		return
	}

	// Sort children: directories first, then files, both alphabetically
	sort.Slice(node.Children, func(i, j int) bool {
		a, b := node.Children[i], node.Children[j]
		
		// Directories come before files
		if a.IsDir && !b.IsDir {
			return true
		}
		if !a.IsDir && b.IsDir {
			return false
		}
		
		// Within same type, sort alphabetically
		return a.Name < b.Name
	})

	// Recursively sort children
	for _, child := range node.Children {
		sortFileTreeNode(child)
	}
}

// PreviewData holds the preview generation data
type PreviewData struct {
	Files []models.GeneratedFile
	Stats *models.PreviewStats
}

// convertToProjectConfig converts a GenerationRequest to internal ProjectConfig
func convertToProjectConfig(req models.GenerationRequest) types.ProjectConfig {
	config := types.ProjectConfig{
		Name:         req.Name,
		Module:       req.ModulePath,
		Type:         req.Type,
		Architecture: req.Architecture,
		Framework:    req.Framework,
		GoVersion:    req.GoVersion,
		Logger:       req.Logger,
		Variables:    make(map[string]string),
	}

	// Add complexity to variables for blueprint selection
	if req.Complexity != "" {
		config.Variables["complexity"] = req.Complexity
	}

	// Add advanced flag
	if req.Advanced {
		config.Variables["advanced"] = "true"
	}

	// Convert features if present
	if req.Features != nil {
		config.Features = &types.Features{}

		if req.Features.Database != nil {
			config.Features.Database = types.DatabaseConfig{
				Driver: req.Features.Database.Driver,
				ORM:    req.Features.Database.ORM,
			}
		}

		if req.Features.Authentication != nil {
			config.Features.Authentication = types.AuthConfig{
				Type:      req.Features.Authentication.Type,
				Providers: req.Features.Authentication.Providers,
			}
		}

		if req.Features.Deployment != nil {
			config.Features.Deployment = types.DeployConfig{
				Targets: req.Features.Deployment.Targets,
			}
		}

		if req.Features.Testing != nil {
			config.Features.Testing = types.TestConfig{
				Framework: req.Features.Testing.Framework,
				Coverage:  req.Features.Testing.Coverage,
			}
		}

		if req.Features.Logging != nil {
			config.Features.Logging = types.LoggingConfig{
				Type:       req.Logger, // Use logger from request
				Level:      req.Features.Logging.Level,
				Format:     req.Features.Logging.Format,
				Structured: true, // Default to structured logging
			}
		}

		if req.Features.Monitoring != nil {
			config.Features.Monitoring = types.MonitorConfig{
				Metrics: req.Features.Monitoring.Metrics,
				Tracing: req.Features.Monitoring.Tracing,
			}
		}
	}

	return config
}

// determineBlueprintID determines the blueprint ID based on the configuration
func determineBlueprintID(config types.ProjectConfig) string {
	// Check for complexity-based blueprint selection (e.g., cli-simple)
	if complexity, exists := config.Variables["complexity"]; exists && complexity == "simple" && config.Type == "cli" {
		return "cli-simple"
	}

	// Use architecture-based selection for other types
	if config.Architecture != "" && config.Architecture != "standard" {
		return fmt.Sprintf("%s-%s", config.Type, config.Architecture)
	}

	return config.Type
}