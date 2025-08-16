package models

import (
	"time"
)

// BlueprintInfo represents metadata about a blueprint
type BlueprintInfo struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Category     string                 `json:"category"`
	Type         string                 `json:"type"`
	Architecture string                 `json:"architecture,omitempty"`
	Framework    string                 `json:"framework,omitempty"`
	Version      string                 `json:"version"`
	Tags         []string               `json:"tags"`
	Icon         string                 `json:"icon,omitempty"`
	
	// Complexity and requirements
	Complexity    string   `json:"complexity"`
	MinGoVersion  string   `json:"minGoVersion,omitempty"`
	Requirements  []string `json:"requirements,omitempty"`
	
	// Statistics
	FileCount     int                    `json:"fileCount"`
	EstimatedSize int64                  `json:"estimatedSize"`
	Features      []string               `json:"features"`
	Dependencies  []DependencyInfo       `json:"dependencies,omitempty"`
	
	// Metadata
	Author        string                 `json:"author,omitempty"`
	License       string                 `json:"license,omitempty"`
	Documentation string                 `json:"documentation,omitempty"`
	Repository    string                 `json:"repository,omitempty"`
	
	// Usage statistics
	Downloads     int64                  `json:"downloads,omitempty"`
	Stars         int                    `json:"stars,omitempty"`
	LastUpdated   time.Time              `json:"lastUpdated,omitempty"`
	
	// Template information
	Variables     map[string]VariableInfo `json:"variables,omitempty"`
	Examples      []ExampleInfo          `json:"examples,omitempty"`
	
	// Compatibility
	OS            []string               `json:"os,omitempty"`           // Supported operating systems
	Platforms     []string               `json:"platforms,omitempty"`    // Supported platforms/architectures
}

// DependencyInfo represents information about a blueprint dependency
type DependencyInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Type        string `json:"type"`        // "go-module", "system", "tool"
	Optional    bool   `json:"optional"`
	Description string `json:"description,omitempty"`
}

// VariableInfo represents information about a template variable
type VariableInfo struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`         // "string", "bool", "int", "select"
	Description  string   `json:"description"`
	Required     bool     `json:"required"`
	Default      string   `json:"default,omitempty"`
	Options      []string `json:"options,omitempty"`      // For select type
	Validation   string   `json:"validation,omitempty"`   // Validation pattern/rule
	Example      string   `json:"example,omitempty"`
}

// ExampleInfo represents an example project configuration
type ExampleInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	UseCase     string                 `json:"useCase,omitempty"`
}

// BlueprintListResponse represents the response for listing blueprints
type BlueprintListResponse struct {
	Blueprints []BlueprintInfo     `json:"blueprints"`
	Categories []CategoryInfo      `json:"categories"`
	Stats      *BlueprintListStats `json:"stats,omitempty"`
}

// CategoryInfo represents information about blueprint categories
type CategoryInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Count       int    `json:"count"`
	Icon        string `json:"icon,omitempty"`
}

// BlueprintListStats contains statistics about the blueprint list
type BlueprintListStats struct {
	TotalBlueprints int            `json:"totalBlueprints"`
	Categories      int            `json:"categories"`
	Complexity      map[string]int `json:"complexity"`
	Types           map[string]int `json:"types"`
}

// BlueprintFilter represents filters for blueprint listing
type BlueprintFilter struct {
	Category     string   `json:"category,omitempty"`
	Type         string   `json:"type,omitempty"`
	Complexity   string   `json:"complexity,omitempty"`
	Architecture string   `json:"architecture,omitempty"`
	Framework    string   `json:"framework,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Search       string   `json:"search,omitempty"`
	MinFileCount int      `json:"minFileCount,omitempty"`
	MaxFileCount int      `json:"maxFileCount,omitempty"`
}

// BlueprintSort represents sorting options for blueprint listing
type BlueprintSort struct {
	Field string `json:"field"`     // "name", "complexity", "fileCount", "lastUpdated", "downloads"
	Order string `json:"order"`     // "asc", "desc"
}

// GetBlueprintCategories returns the available blueprint categories
func GetBlueprintCategories() []CategoryInfo {
	return []CategoryInfo{
		{
			ID:          "web-api",
			Name:        "Web APIs",
			Description: "REST APIs and web services",
			Icon:        "globe",
		},
		{
			ID:          "cli",
			Name:        "Command Line Tools",
			Description: "CLI applications and utilities",
			Icon:        "terminal",
		},
		{
			ID:          "library",
			Name:        "Libraries",
			Description: "Reusable Go packages and libraries",
			Icon:        "package",
		},
		{
			ID:          "serverless",
			Name:        "Serverless",
			Description: "AWS Lambda and serverless functions",
			Icon:        "cloud",
		},
		{
			ID:          "microservices",
			Name:        "Microservices",
			Description: "Distributed systems and microservices",
			Icon:        "grid",
		},
		{
			ID:          "monolith",
			Name:        "Monolithic Apps",
			Description: "Full-featured monolithic applications",
			Icon:        "layers",
		},
		{
			ID:          "workspace",
			Name:        "Workspaces",
			Description: "Multi-module Go workspaces",
			Icon:        "folder",
		},
	}
}

// GetComplexityLevels returns the available complexity levels
func GetComplexityLevels() []ComplexityInfo {
	return []ComplexityInfo{
		{
			ID:          "simple",
			Name:        "Simple",
			Description: "Minimal structure for quick prototypes",
			FileCount:   "5-15 files",
			Features:    "Basic functionality, minimal dependencies",
			Recommended: "Learning, quick scripts, simple utilities",
		},
		{
			ID:          "standard",
			Name:        "Standard",
			Description: "Balanced structure for production use",
			FileCount:   "15-40 files",
			Features:    "Common patterns, standard dependencies",
			Recommended: "Most production applications",
		},
		{
			ID:          "advanced",
			Name:        "Advanced",
			Description: "Complex architecture with advanced patterns",
			FileCount:   "40-80 files",
			Features:    "Advanced patterns, extensive dependencies",
			Recommended: "Enterprise applications, complex domains",
		},
		{
			ID:          "expert",
			Name:        "Expert",
			Description: "Comprehensive structure with all features",
			FileCount:   "80+ files",
			Features:    "All patterns, complete feature set",
			Recommended: "Large-scale systems, platform development",
		},
	}
}