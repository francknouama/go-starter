package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/internal/web/models"
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