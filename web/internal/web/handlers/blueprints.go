package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/web/internal/web/models"
)

// ListBlueprints returns all available project blueprints
func ListBlueprints(c *gin.Context) {
	// Initialize template registry
	registry := templates.NewRegistry()
	
	// Get query parameters for filtering
	category := c.Query("category")
	projectType := c.Query("type")
	complexity := c.Query("complexity")
	architecture := c.Query("architecture")
	framework := c.Query("framework")
	search := c.Query("search")
	
	// Get all templates from registry
	allTemplates := registry.List()
	
	// Convert to blueprint info
	var blueprints []models.BlueprintInfo
	for _, tmpl := range allTemplates {
		blueprint := convertTemplateToBlueprint(tmpl)
		
		// Apply filters
		if category != "" && blueprint.Category != category {
			continue
		}
		if projectType != "" && blueprint.Type != projectType {
			continue
		}
		if complexity != "" && blueprint.Complexity != complexity {
			continue
		}
		if architecture != "" && blueprint.Architecture != architecture {
			continue
		}
		if framework != "" && blueprint.Framework != framework {
			continue
		}
		if search != "" && !matchesSearch(blueprint, search) {
			continue
		}
		
		blueprints = append(blueprints, blueprint)
	}
	
	// Generate statistics
	stats := generateBlueprintStats(blueprints)
	
	// Generate categories
	categories := models.GetBlueprintCategories()
	updateCategoryCounts(categories, blueprints)
	
	response := &models.BlueprintListResponse{
		Blueprints: blueprints,
		Categories: categories,
		Stats:      stats,
	}
	
	apiResponse := models.NewSuccessResponse(response)
	c.JSON(http.StatusOK, apiResponse)
}

// GetBlueprint returns details for a specific blueprint
func GetBlueprint(c *gin.Context) {
	blueprintID := c.Param("id")
	if blueprintID == "" {
		response := models.NewValidationErrorResponse("Blueprint ID is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	// Initialize template registry
	registry := templates.NewRegistry()
	
	// Get template
	tmpl, err := registry.Get(blueprintID)
	if err != nil {
		response := models.NewNotFoundErrorResponse("Blueprint not found")
		c.JSON(http.StatusNotFound, response)
		return
	}
	
	// Convert to blueprint info
	blueprint := convertTemplateToBlueprint(tmpl)
	
	apiResponse := models.NewSuccessResponse(blueprint)
	c.JSON(http.StatusOK, apiResponse)
}

// GetBlueprintsByCategory returns blueprints filtered by category
func GetBlueprintsByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		response := models.NewValidationErrorResponse("Category is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	// Initialize template registry
	registry := templates.NewRegistry()
	allTemplates := registry.List()
	
	// Filter by category and convert to blueprint info
	var blueprints []models.BlueprintInfo
	for _, tmpl := range allTemplates {
		blueprint := convertTemplateToBlueprint(tmpl)
		if blueprint.Category == category {
			blueprints = append(blueprints, blueprint)
		}
	}
	
	response := map[string]interface{}{
		"category":   category,
		"blueprints": blueprints,
		"count":      len(blueprints),
	}
	
	apiResponse := models.NewSuccessResponse(response)
	c.JSON(http.StatusOK, apiResponse)
}

// GetBlueprintValidation validates a blueprint configuration
func GetBlueprintValidation(c *gin.Context) {
	blueprintID := c.Param("id")
	if blueprintID == "" {
		response := models.NewValidationErrorResponse("Blueprint ID is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	var request models.GenerationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response := models.NewValidationErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	// Validate the request
	if err := request.Validate(); err != nil {
		validationResponse := map[string]interface{}{
			"valid":  false,
			"errors": []string{err.Error()},
		}
		
		apiResponse := models.NewSuccessResponse(validationResponse)
		c.JSON(http.StatusOK, apiResponse)
		return
	}
	
	// Additional blueprint-specific validation could go here
	
	validationResponse := map[string]interface{}{
		"valid":   true,
		"errors":  []string{},
		"message": "Configuration is valid",
	}
	
	apiResponse := models.NewSuccessResponse(validationResponse)
	c.JSON(http.StatusOK, apiResponse)
}

// convertTemplateToBlueprint converts a template to blueprint info
func convertTemplateToBlueprint(tmpl types.Template) models.BlueprintInfo {
	// Determine complexity based on estimated file count or explicit metadata
	complexity := determineComplexity(tmpl)
	
	// Extract architecture from template ID
	architecture := extractArchitecture(tmpl.ID, tmpl.Type)
	
	// Extract framework from metadata or template
	framework := extractFramework(tmpl)
	
	// Generate tags based on template properties
	tags := generateTags(tmpl)
	
	// Estimate file count and size
	fileCount := len(tmpl.Files)
	estimatedSize := estimateProjectSize(tmpl)
	
	// Extract features
	features := extractFeatures(tmpl)
	
	// Convert dependencies
	dependencies := convertDependencies(tmpl.Dependencies)
	
	// Extract variables
	variables := convertVariables(tmpl.Variables)
	
	// Generate examples
	examples := generateExamples(tmpl)
	
	return models.BlueprintInfo{
		ID:           tmpl.ID,
		Name:         tmpl.Name,
		Description:  tmpl.Description,
		Category:     mapTypeToCategory(tmpl.Type),
		Type:         tmpl.Type,
		Architecture: architecture,
		Framework:    framework,
		Version:      tmpl.Version,
		Tags:         tags,
		Icon:         getIconForType(tmpl.Type),
		
		Complexity:   complexity,
		MinGoVersion: "1.21", // Default minimum Go version
		
		FileCount:     fileCount,
		EstimatedSize: estimatedSize,
		Features:      features,
		Dependencies:  dependencies,
		
		Author:        tmpl.Author,
		License:       tmpl.License,
		Documentation: "", // Could be extracted from metadata
		Repository:    "", // Could be extracted from metadata
		
		Variables: variables,
		Examples:  examples,
		
		OS:        []string{"linux", "darwin", "windows"}, // Default to all platforms
		Platforms: []string{"amd64", "arm64"},             // Default to common architectures
	}
}

// determineComplexity determines the complexity level of a template
func determineComplexity(tmpl types.Template) string {
	fileCount := len(tmpl.Files)
	
	// Check if complexity is explicitly defined in metadata
	if complexity, ok := tmpl.Metadata["complexity"].(string); ok {
		return complexity
	}
	
	// Determine based on file count and template ID
	switch {
	case strings.Contains(tmpl.ID, "simple"):
		return "simple"
	case fileCount <= 15:
		return "simple"
	case fileCount <= 40:
		return "standard"
	case fileCount <= 80:
		return "advanced"
	default:
		return "expert"
	}
}

// extractArchitecture extracts architecture pattern from template ID
func extractArchitecture(templateID, templateType string) string {
	// Remove type prefix to get architecture
	if strings.HasPrefix(templateID, templateType+"-") {
		arch := strings.TrimPrefix(templateID, templateType+"-")
		// Validate it's a known architecture
		knownArchs := []string{"clean", "ddd", "hexagonal", "event-driven"}
		for _, known := range knownArchs {
			if arch == known {
				return arch
			}
		}
	}
	return "standard"
}

// extractFramework extracts framework from template metadata
func extractFramework(tmpl types.Template) string {
	if framework, ok := tmpl.Metadata["framework"].(string); ok {
		return framework
	}
	
	// Infer from template ID or dependencies
	for _, dep := range tmpl.Dependencies {
		switch dep.Module {
		case "github.com/gin-gonic/gin":
			return "gin"
		case "github.com/labstack/echo/v4":
			return "echo"
		case "github.com/gofiber/fiber/v2":
			return "fiber"
		case "github.com/go-chi/chi/v5":
			return "chi"
		case "github.com/spf13/cobra":
			return "cobra"
		}
	}
	
	return ""
}

// generateTags generates tags for a template
func generateTags(tmpl types.Template) []string {
	tags := []string{}
	
	// Add type-based tags
	switch tmpl.Type {
	case "web-api":
		tags = append(tags, "api", "web", "rest", "http")
	case "cli":
		tags = append(tags, "cli", "command-line", "tool")
	case "library":
		tags = append(tags, "library", "package", "reusable")
	case "lambda":
		tags = append(tags, "serverless", "lambda", "aws")
	case "microservice":
		tags = append(tags, "microservice", "grpc", "distributed")
	}
	
	// Add architecture-based tags
	architecture := extractArchitecture(tmpl.ID, tmpl.Type)
	if architecture != "standard" {
		tags = append(tags, architecture)
	}
	
	// Add framework-based tags
	framework := extractFramework(tmpl)
	if framework != "" {
		tags = append(tags, framework)
	}
	
	return tags
}

// estimateProjectSize estimates the total size of generated project
func estimateProjectSize(tmpl types.Template) int64 {
	// Rough estimation based on file count and type
	baseSize := int64(len(tmpl.Files) * 1024) // 1KB per file on average
	
	// Adjust based on project type
	switch tmpl.Type {
	case "web-api":
		baseSize *= 2 // Web APIs tend to have larger files
	case "microservice":
		baseSize *= 3 // Microservices have more boilerplate
	case "monolith":
		baseSize *= 5 // Monoliths are typically larger
	}
	
	return baseSize
}

// extractFeatures extracts feature list from template
func extractFeatures(tmpl types.Template) []string {
	features := []string{}
	
	// Extract from dependencies
	for _, dep := range tmpl.Dependencies {
		switch {
		case strings.Contains(dep.Module, "postgres"):
			features = append(features, "PostgreSQL")
		case strings.Contains(dep.Module, "mysql"):
			features = append(features, "MySQL")
		case strings.Contains(dep.Module, "mongo"):
			features = append(features, "MongoDB")
		case strings.Contains(dep.Module, "redis"):
			features = append(features, "Redis")
		case strings.Contains(dep.Module, "gorm"):
			features = append(features, "GORM ORM")
		case strings.Contains(dep.Module, "jwt"):
			features = append(features, "JWT Authentication")
		case strings.Contains(dep.Module, "grpc"):
			features = append(features, "gRPC")
		case strings.Contains(dep.Module, "prometheus"):
			features = append(features, "Metrics")
		case strings.Contains(dep.Module, "jaeger") || strings.Contains(dep.Module, "opentelemetry"):
			features = append(features, "Tracing")
		}
	}
	
	// Extract from file names
	for _, file := range tmpl.Files {
		if strings.Contains(file.Destination, "docker") {
			features = append(features, "Docker")
		}
		if strings.Contains(file.Destination, "kubernetes") {
			features = append(features, "Kubernetes")
		}
		if strings.Contains(file.Destination, "middleware") {
			features = append(features, "Middleware")
		}
	}
	
	return uniqueStrings(features)
}

// convertDependencies converts template dependencies to dependency info
func convertDependencies(deps []types.Dependency) []models.DependencyInfo {
	dependencies := make([]models.DependencyInfo, len(deps))
	for i, dep := range deps {
		dependencies[i] = models.DependencyInfo{
			Name:        extractPackageName(dep.Module),
			Version:     dep.Version,
			Type:        "go-module",
			Optional:    false, // Field doesn't exist in types.Dependency
			Description: "",    // Field doesn't exist in types.Dependency
		}
	}
	return dependencies
}

// convertVariables converts template variables to variable info
func convertVariables(vars []types.TemplateVariable) map[string]models.VariableInfo {
	variables := make(map[string]models.VariableInfo)
	for _, v := range vars {
		variables[v.Name] = models.VariableInfo{
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
			Required:    v.Required,
			Default:     fmt.Sprintf("%v", v.Default),
		}
	}
	return variables
}

// generateExamples generates example configurations for a template
func generateExamples(tmpl types.Template) []models.ExampleInfo {
	examples := []models.ExampleInfo{}
	
	// Generate basic example
	basicConfig := map[string]interface{}{
		"name":         "my-" + tmpl.Type,
		"modulePath":   "github.com/user/my-" + tmpl.Type,
		"type":         tmpl.Type,
		"architecture": "standard",
		"logger":       "slog",
	}
	
	examples = append(examples, models.ExampleInfo{
		Name:        "Basic " + tmpl.Name,
		Description: "Basic configuration for " + tmpl.Type,
		Config:      basicConfig,
		UseCase:     "Getting started with " + tmpl.Type,
	})
	
	// Generate advanced example if applicable
	if !strings.Contains(tmpl.ID, "simple") {
		advancedConfig := map[string]interface{}{
			"name":         "advanced-" + tmpl.Type,
			"modulePath":   "github.com/company/advanced-" + tmpl.Type,
			"type":         tmpl.Type,
			"architecture": extractArchitecture(tmpl.ID, tmpl.Type),
			"logger":       "zap",
			"features": map[string]interface{}{
				"database": map[string]interface{}{
					"driver": "postgres",
					"orm":    "gorm",
				},
			},
		}
		
		examples = append(examples, models.ExampleInfo{
			Name:        "Advanced " + tmpl.Name,
			Description: "Production-ready configuration with database",
			Config:      advancedConfig,
			UseCase:     "Production deployment",
		})
	}
	
	return examples
}

// Helper functions

func mapTypeToCategory(projectType string) string {
	switch projectType {
	case "web-api":
		return "web-api"
	case "cli":
		return "cli"
	case "library":
		return "library"
	case "lambda", "lambda-proxy":
		return "serverless"
	case "microservice":
		return "microservices"
	case "monolith":
		return "monolith"
	case "workspace":
		return "workspace"
	default:
		return "other"
	}
}

func getIconForType(projectType string) string {
	switch projectType {
	case "web-api":
		return "globe"
	case "cli":
		return "terminal"
	case "library":
		return "package"
	case "lambda", "lambda-proxy":
		return "cloud"
	case "microservice":
		return "grid"
	case "monolith":
		return "layers"
	case "workspace":
		return "folder"
	default:
		return "file"
	}
}

func matchesSearch(blueprint models.BlueprintInfo, search string) bool {
	search = strings.ToLower(search)
	return strings.Contains(strings.ToLower(blueprint.Name), search) ||
		strings.Contains(strings.ToLower(blueprint.Description), search) ||
		containsTag(blueprint.Tags, search)
}

func containsTag(tags []string, search string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), search) {
			return true
		}
	}
	return false
}

func generateBlueprintStats(blueprints []models.BlueprintInfo) *models.BlueprintListStats {
	stats := &models.BlueprintListStats{
		TotalBlueprints: len(blueprints),
		Complexity:      make(map[string]int),
		Types:           make(map[string]int),
	}
	
	categories := make(map[string]bool)
	
	for _, blueprint := range blueprints {
		stats.Complexity[blueprint.Complexity]++
		stats.Types[blueprint.Type]++
		categories[blueprint.Category] = true
	}
	
	stats.Categories = len(categories)
	return stats
}

func updateCategoryCounts(categories []models.CategoryInfo, blueprints []models.BlueprintInfo) {
	categoryCounts := make(map[string]int)
	for _, blueprint := range blueprints {
		categoryCounts[blueprint.Category]++
	}
	
	for i := range categories {
		categories[i].Count = categoryCounts[categories[i].ID]
	}
}

func extractPackageName(module string) string {
	parts := strings.Split(module, "/")
	return parts[len(parts)-1]
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}