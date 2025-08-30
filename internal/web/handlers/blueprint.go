package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/internal/web/models"
	"github.com/francknouama/go-starter/pkg/types"
)

type BlueprintHandler struct {
	registry *templates.Registry
}

func NewBlueprintHandler() *BlueprintHandler {
	registry := templates.NewRegistry()
	return &BlueprintHandler{
		registry: registry,
	}
}

// ListBlueprints returns all available blueprints
func (h *BlueprintHandler) ListBlueprints(c *gin.Context) {
	templates := h.registry.List()
	blueprints := make([]models.Blueprint, 0, len(templates))

	for _, template := range templates {
		blueprint := models.Blueprint{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			Type:        template.Type,
			Complexity:  getComplexityLevel(template.ID),
			FileCount:   len(template.Files),
		}

		// Add features from template features
		features := make([]string, len(template.Features))
		for i, feature := range template.Features {
			features[i] = feature.Name
		}
		blueprint.Features = features

		// Add dependencies from template dependencies
		dependencies := make([]string, len(template.Dependencies))
		for i, dep := range template.Dependencies {
			dependencies[i] = dep.Module
		}
		blueprint.Dependencies = dependencies

		blueprints = append(blueprints, blueprint)
	}

	c.JSON(http.StatusOK, models.BlueprintListResponse{
		Blueprints: blueprints,
	})
}

// GetBlueprint returns details for a specific blueprint
func (h *BlueprintHandler) GetBlueprint(c *gin.Context) {
	blueprintID := c.Param("id")
	
	template, err := h.registry.Get(blueprintID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Blueprint not found",
			"code":  "BLUEPRINT_NOT_FOUND",
		})
		return
	}

	blueprint := models.Blueprint{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		Type:        template.Type,
		Complexity:  getComplexityLevel(template.ID),
		FileCount:   len(template.Files),
	}

	// Add features from template features
	features := make([]string, len(template.Features))
	for i, feature := range template.Features {
		features[i] = feature.Name
	}
	blueprint.Features = features

	// Add dependencies from template dependencies
	dependencies := make([]string, len(template.Dependencies))
	for i, dep := range template.Dependencies {
		dependencies[i] = dep.Module
	}
	blueprint.Dependencies = dependencies

	// Add file list
	files := make([]models.BlueprintFile, 0, len(template.Files))
	for _, file := range template.Files {
		files = append(files, models.BlueprintFile{
			Source:      file.Source,
			Destination: file.Destination,
			Condition:   file.Condition,
		})
	}

	// Convert template variables to map
	variables := make(map[string]interface{})
	for _, v := range template.Variables {
		variables[v.Name] = map[string]interface{}{
			"type":        v.Type,
			"description": v.Description,
			"default":     v.Default,
			"required":    v.Required,
			"choices":     v.Choices,
			"validation":  v.Validation,
		}
	}

	response := models.BlueprintDetailResponse{
		Blueprint: blueprint,
		Files:     files,
		Variables: variables,
	}

	c.JSON(http.StatusOK, response)
}

// getComplexityLevel determines the complexity level based on blueprint name
func getComplexityLevel(name string) string {
	switch name {
	case "cli-simple":
		return "simple"
	case "cli", "library-standard", "lambda-standard":
		return "standard"
	case "web-api-clean", "web-api-ddd", "microservice-standard":
		return "advanced"
	case "web-api-hexagonal", "grpc-gateway":
		return "expert"
	default:
		return "standard"
	}
}

// ValidateBlueprintConfig validates a configuration for a specific blueprint
func (h *BlueprintHandler) ValidateBlueprintConfig(c *gin.Context) {
	blueprintID := c.Param("type")
	
	var req models.ValidateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"code":  "INVALID_REQUEST",
		})
		return
	}

	// Get blueprint to validate against
	template, err := h.registry.Get(blueprintID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Blueprint not found",
			"code":  "BLUEPRINT_NOT_FOUND",
		})
		return
	}

	// Validate configuration against blueprint requirements
	errors := h.validateAgainstBlueprint(req.Config, template)
	
	c.JSON(http.StatusOK, models.ValidateConfigResponse{
		Valid:  len(errors) == 0,
		Errors: errors,
	})
}

// GetBlueprintDefaults returns default configuration for a blueprint
func (h *BlueprintHandler) GetBlueprintDefaults(c *gin.Context) {
	blueprintType := c.Param("type")
	
	template, err := h.registry.Get(blueprintType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Blueprint not found",
			"code":  "BLUEPRINT_NOT_FOUND",
		})
		return
	}

	// Build default configuration based on blueprint type
	defaultConfig := h.buildDefaultConfig(blueprintType, template)
	
	c.JSON(http.StatusOK, gin.H{
		"config": defaultConfig,
	})
}

// GetBlueprintOptions returns available options for a blueprint
func (h *BlueprintHandler) GetBlueprintOptions(c *gin.Context) {
	blueprintType := c.Param("type")
	
	template, err := h.registry.Get(blueprintType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Blueprint not found",
			"code":  "BLUEPRINT_NOT_FOUND",
		})
		return
	}

	// Build available options for this blueprint
	options := h.buildBlueprintOptions(blueprintType, template)
	
	c.JSON(http.StatusOK, gin.H{
		"options": options,
	})
}

// Helper methods

func (h *BlueprintHandler) validateAgainstBlueprint(config models.ProjectConfig, template types.Template) []models.ValidationError {
	var errors []models.ValidationError

	// Basic validation
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

	if config.GoVersion == "" {
		errors = append(errors, models.ValidationError{
			Field:    "go_version",
			Message:  "Go version is required",
			Severity: "error",
		})
	}

	// Blueprint-specific validation
	switch template.Type {
	case "web-api":
		if config.Framework == "" {
			errors = append(errors, models.ValidationError{
				Field:    "framework",
				Message:  "Framework is required for web API projects",
				Severity: "error",
			})
		}
	case "grpc":
		if config.Architecture == "" {
			errors = append(errors, models.ValidationError{
				Field:    "architecture",
				Message:  "Architecture pattern is recommended for gRPC services",
				Severity: "warning",
			})
		}
	}

	return errors
}

func (h *BlueprintHandler) buildDefaultConfig(blueprintType string, template types.Template) models.ProjectConfig {
	config := models.ProjectConfig{
		ProjectName: "my-project",
		ModuleURL:   "github.com/user/my-project",
		GoVersion:   "1.21",
		ProjectType: blueprintType,
		Logger:      "slog",
	}

	// Set defaults based on blueprint type
	switch template.Type {
	case "web-api":
		config.Framework = "gin"
		config.Architecture = "standard"
	case "cli":
		config.Framework = "cobra"
	case "grpc":
		config.Architecture = "clean"
	case "lambda":
		config.Framework = "aws-lambda-go"
	case "microservice":
		config.Framework = "grpc"
		config.Architecture = "hexagonal"
	case "monolith":
		config.Framework = "gin"
		config.Architecture = "mvc"
	}

	return config
}

func (h *BlueprintHandler) buildBlueprintOptions(blueprintType string, template types.Template) map[string]interface{} {
	options := make(map[string]interface{})

	// Common options
	options["go_versions"] = []string{"1.21", "1.22", "1.23"}
	options["loggers"] = []string{"slog", "zap", "logrus", "zerolog"}

	// Blueprint-specific options
	switch template.Type {
	case "web-api":
		options["frameworks"] = []string{"gin", "echo", "fiber", "chi"}
		options["architectures"] = []string{"standard", "clean", "hexagonal", "ddd"}
		options["databases"] = map[string]interface{}{
			"drivers": []string{"postgres", "mysql", "sqlite", "mongodb"},
			"orms":    []string{"gorm", "ent", "sqlx"},
		}
	case "cli":
		options["frameworks"] = []string{"cobra", "urfave/cli"}
		options["complexities"] = []string{"simple", "standard", "advanced"}
	case "grpc":
		options["architectures"] = []string{"clean", "hexagonal", "ddd"}
		options["features"] = []string{"gateway", "reflection", "health-check"}
	case "lambda":
		options["runtimes"] = []string{"provided.al2", "go1.x"}
		options["triggers"] = []string{"api-gateway", "sqs", "s3", "dynamodb"}
	case "microservice":
		options["frameworks"] = []string{"grpc", "http"}
		options["architectures"] = []string{"hexagonal", "clean", "ddd"}
		options["patterns"] = []string{"cqrs", "event-sourcing", "saga"}
	case "monolith":
		options["frameworks"] = []string{"gin", "echo", "fiber"}
		options["architectures"] = []string{"mvc", "layered", "modular"}
		options["frontend"] = []string{"htmx", "vue", "react", "none"}
	}

	return options
}