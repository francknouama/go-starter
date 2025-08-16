package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/francknouama/go-starter/web/internal/web/models"
)

// GetConfig returns the default configuration options for the project generator
func GetConfig(c *gin.Context) {
	config := &models.ConfigResponse{
		ProjectTypes: []models.ProjectTypeInfo{
			{
				ID:          "web-api",
				Name:        "Web API",
				Description: "REST API with web framework",
				Icon:        "globe",
				Category:    "web-api",
				Tags:        []string{"rest", "api", "web"},
				Complexity:  "intermediate",
				EstimatedTime: "5-10 minutes",
			},
			{
				ID:          "cli",
				Name:        "CLI Application",
				Description: "Command-line interface application",
				Icon:        "terminal",
				Category:    "cli",
				Tags:        []string{"cli", "command-line", "tool"},
				Complexity:  "beginner",
				EstimatedTime: "2-5 minutes",
			},
			{
				ID:          "library",
				Name:        "Go Library",
				Description: "Reusable Go package or library",
				Icon:        "package",
				Category:    "library",
				Tags:        []string{"library", "package", "reusable"},
				Complexity:  "beginner",
				EstimatedTime: "3-7 minutes",
			},
			{
				ID:          "lambda",
				Name:        "AWS Lambda",
				Description: "Serverless function for AWS Lambda",
				Icon:        "cloud",
				Category:    "serverless",
				Tags:        []string{"serverless", "lambda", "aws"},
				Complexity:  "intermediate",
				EstimatedTime: "5-8 minutes",
			},
			{
				ID:          "lambda-proxy",
				Name:        "Lambda API Proxy",
				Description: "API Gateway proxy integration with Lambda",
				Icon:        "cloud-lightning",
				Category:    "serverless",
				Tags:        []string{"serverless", "api-gateway", "proxy"},
				Complexity:  "advanced",
				EstimatedTime: "8-12 minutes",
			},
			{
				ID:          "microservice",
				Name:        "Microservice",
				Description: "Distributed microservice with gRPC",
				Icon:        "grid",
				Category:    "microservices",
				Tags:        []string{"microservice", "grpc", "distributed"},
				Complexity:  "advanced",
				EstimatedTime: "10-15 minutes",
			},
			{
				ID:          "monolith",
				Name:        "Monolithic App",
				Description: "Full-featured monolithic application",
				Icon:        "layers",
				Category:    "monolith",
				Tags:        []string{"monolith", "full-stack", "web-app"},
				Complexity:  "expert",
				EstimatedTime: "15-25 minutes",
			},
			{
				ID:          "workspace",
				Name:        "Go Workspace",
				Description: "Multi-module Go workspace project",
				Icon:        "folder",
				Category:    "workspace",
				Tags:        []string{"workspace", "multi-module", "monorepo"},
				Complexity:  "expert",
				EstimatedTime: "10-20 minutes",
			},
		},
		Architectures: []models.ArchitectureInfo{
			{
				ID:          "standard",
				Name:        "Standard",
				Description: "Simple layered architecture",
				Complexity:  "beginner",
				UseCase:     "Simple applications and prototypes",
				Benefits:    []string{"Easy to understand", "Quick to implement", "Low complexity"},
				Drawbacks:   []string{"Limited scalability", "Tight coupling possible"},
			},
			{
				ID:          "clean",
				Name:        "Clean Architecture",
				Description: "Uncle Bob's Clean Architecture pattern",
				Complexity:  "intermediate",
				UseCase:     "Business-critical applications",
				Benefits:    []string{"Testable", "Framework independent", "Clear separation"},
				Drawbacks:   []string{"More complex", "Higher learning curve"},
			},
			{
				ID:          "ddd",
				Name:        "Domain-Driven Design",
				Description: "Domain-focused architecture with bounded contexts",
				Complexity:  "advanced",
				UseCase:     "Complex business domains",
				Benefits:    []string{"Domain focus", "Business alignment", "Bounded contexts"},
				Drawbacks:   []string{"High complexity", "Requires domain expertise"},
			},
			{
				ID:          "hexagonal",
				Name:        "Hexagonal Architecture",
				Description: "Ports and adapters pattern",
				Complexity:  "advanced",
				UseCase:     "Highly testable applications",
				Benefits:    []string{"Highly testable", "Technology agnostic", "Clear boundaries"},
				Drawbacks:   []string{"Complex setup", "Over-engineering risk"},
			},
		},
		Frameworks: []models.FrameworkInfo{
			{
				ID:          "gin",
				Name:        "Gin",
				Description: "Fast HTTP web framework",
				ProjectTypes: []string{"web-api", "microservice", "monolith"},
				Performance: "excellent",
				Popularity:  "very-high",
			},
			{
				ID:          "echo",
				Name:        "Echo",
				Description: "High performance, minimalist web framework",
				ProjectTypes: []string{"web-api", "microservice", "monolith"},
				Performance: "excellent",
				Popularity:  "high",
			},
			{
				ID:          "fiber",
				Name:        "Fiber",
				Description: "Express.js inspired web framework",
				ProjectTypes: []string{"web-api", "microservice", "monolith"},
				Performance: "excellent",
				Popularity:  "high",
			},
			{
				ID:          "chi",
				Name:        "Chi",
				Description: "Lightweight, composable HTTP router",
				ProjectTypes: []string{"web-api", "microservice"},
				Performance: "good",
				Popularity:  "medium",
			},
			{
				ID:          "cobra",
				Name:        "Cobra",
				Description: "CLI application framework",
				ProjectTypes: []string{"cli"},
				Performance: "good",
				Popularity:  "very-high",
			},
		},
		Loggers: []models.LoggerInfo{
			{
				ID:          "slog",
				Name:        "slog",
				Description: "Go standard library structured logger",
				Performance: "good",
				Features:    "structured logging, levels, context",
				Package:     "log/slog",
			},
			{
				ID:          "zap",
				Name:        "Zap",
				Description: "Uber's high-performance logger",
				Performance: "excellent",
				Features:    "zero allocation, high performance, structured",
				Package:     "go.uber.org/zap",
			},
			{
				ID:          "logrus",
				Name:        "Logrus",
				Description: "Structured logger with hooks",
				Performance: "good",
				Features:    "hooks, formatters, structured logging",
				Package:     "github.com/sirupsen/logrus",
			},
			{
				ID:          "zerolog",
				Name:        "Zerolog",
				Description: "Zero allocation JSON logger",
				Performance: "excellent",
				Features:    "zero allocation, chainable API, fast",
				Package:     "github.com/rs/zerolog",
			},
		},
		Databases: []models.DatabaseInfo{
			{
				ID:          "postgres",
				Name:        "PostgreSQL",
				Description: "Advanced open-source relational database",
				Type:        "sql",
				ORMs:        []string{"gorm", "sqlx", "sqlc", "ent"},
			},
			{
				ID:          "mysql",
				Name:        "MySQL",
				Description: "Popular open-source relational database",
				Type:        "sql",
				ORMs:        []string{"gorm", "sqlx", "sqlc"},
			},
			{
				ID:          "sqlite",
				Name:        "SQLite",
				Description: "Lightweight embedded SQL database",
				Type:        "sql",
				ORMs:        []string{"gorm", "sqlx", "sqlc"},
			},
			{
				ID:          "mongodb",
				Name:        "MongoDB",
				Description: "Document-oriented NoSQL database",
				Type:        "nosql",
				ORMs:        []string{"mongo-driver"},
			},
			{
				ID:          "redis",
				Name:        "Redis",
				Description: "In-memory data structure store",
				Type:        "cache",
				ORMs:        []string{"go-redis", "redigo"},
			},
		},
		Complexities: models.GetComplexityLevels(),
		DefaultValues: &models.DefaultConfig{
			GoVersion:    "1.21",
			Logger:       "slog",
			Architecture: "standard",
			Complexity:   "standard",
			Framework:    "", // Will be set based on project type
		},
	}
	
	response := models.NewSuccessResponse(config)
	c.JSON(http.StatusOK, response)
}

// GetProjectTypeDetails returns detailed information about a specific project type
func GetProjectTypeDetails(c *gin.Context) {
	projectType := c.Param("type")
	if projectType == "" {
		response := models.NewValidationErrorResponse("Project type is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	// This would typically fetch from a configuration service or database
	var details *models.ProjectTypeInfo
	
	switch projectType {
	case "web-api":
		details = &models.ProjectTypeInfo{
			ID:          "web-api",
			Name:        "Web API",
			Description: "REST API with web framework supporting multiple architecture patterns",
			Icon:        "globe",
			Category:    "web-api",
			Tags:        []string{"rest", "api", "web", "http", "json"},
			Complexity:  "intermediate",
			EstimatedTime: "5-10 minutes",
		}
	case "cli":
		details = &models.ProjectTypeInfo{
			ID:          "cli",
			Name:        "CLI Application",
			Description: "Command-line interface application with subcommands and configuration support",
			Icon:        "terminal",
			Category:    "cli",
			Tags:        []string{"cli", "command-line", "tool", "cobra"},
			Complexity:  "beginner",
			EstimatedTime: "2-5 minutes",
		}
	default:
		response := models.NewNotFoundErrorResponse("Project type not found")
		c.JSON(http.StatusNotFound, response)
		return
	}
	
	response := models.NewSuccessResponse(details)
	c.JSON(http.StatusOK, response)
}

// GetFrameworksForType returns frameworks compatible with a project type
func GetFrameworksForType(c *gin.Context) {
	projectType := c.Query("type")
	if projectType == "" {
		response := models.NewValidationErrorResponse("Project type query parameter is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	var frameworks []models.FrameworkInfo
	
	switch projectType {
	case "web-api":
		frameworks = []models.FrameworkInfo{
			{
				ID:          "gin",
				Name:        "Gin",
				Description: "Fast HTTP web framework with middleware support",
				ProjectTypes: []string{"web-api"},
				Performance: "excellent",
				Popularity:  "very-high",
			},
			{
				ID:          "echo",
				Name:        "Echo",
				Description: "High performance, minimalist web framework",
				ProjectTypes: []string{"web-api"},
				Performance: "excellent",
				Popularity:  "high",
			},
			{
				ID:          "fiber",
				Name:        "Fiber",
				Description: "Express.js inspired web framework built on Fasthttp",
				ProjectTypes: []string{"web-api"},
				Performance: "excellent",
				Popularity:  "high",
			},
			{
				ID:          "chi",
				Name:        "Chi",
				Description: "Lightweight, composable HTTP router",
				ProjectTypes: []string{"web-api"},
				Performance: "good",
				Popularity:  "medium",
			},
		}
	case "cli":
		frameworks = []models.FrameworkInfo{
			{
				ID:          "cobra",
				Name:        "Cobra",
				Description: "CLI application framework used by Docker, Kubernetes, and GitHub CLI",
				ProjectTypes: []string{"cli"},
				Performance: "good",
				Popularity:  "very-high",
			},
		}
	default:
		frameworks = []models.FrameworkInfo{}
	}
	
	response := models.NewSuccessResponse(map[string]interface{}{
		"projectType": projectType,
		"frameworks":  frameworks,
		"count":       len(frameworks),
	})
	
	c.JSON(http.StatusOK, response)
}

// GetArchitecturesForType returns architectures compatible with a project type
func GetArchitecturesForType(c *gin.Context) {
	projectType := c.Query("type")
	if projectType == "" {
		response := models.NewValidationErrorResponse("Project type query parameter is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	var architectures []models.ArchitectureInfo
	
	switch projectType {
	case "web-api":
		architectures = []models.ArchitectureInfo{
			{
				ID:          "standard",
				Name:        "Standard",
				Description: "Simple layered architecture with controllers, services, and repositories",
				Complexity:  "beginner",
				UseCase:     "Simple APIs and prototypes",
				Benefits:    []string{"Easy to understand", "Quick to implement", "Low complexity"},
				Drawbacks:   []string{"Limited scalability", "Tight coupling possible"},
			},
			{
				ID:          "clean",
				Name:        "Clean Architecture",
				Description: "Uncle Bob's Clean Architecture with use cases and entities",
				Complexity:  "intermediate",
				UseCase:     "Business-critical APIs with complex logic",
				Benefits:    []string{"Testable", "Framework independent", "Clear separation"},
				Drawbacks:   []string{"More complex", "Higher learning curve"},
			},
			{
				ID:          "ddd",
				Name:        "Domain-Driven Design",
				Description: "Domain-focused architecture with aggregates and value objects",
				Complexity:  "advanced",
				UseCase:     "Complex business domains with rich models",
				Benefits:    []string{"Domain focus", "Business alignment", "Bounded contexts"},
				Drawbacks:   []string{"High complexity", "Requires domain expertise"},
			},
			{
				ID:          "hexagonal",
				Name:        "Hexagonal Architecture",
				Description: "Ports and adapters pattern for maximum testability",
				Complexity:  "advanced",
				UseCase:     "APIs requiring high testability and multiple adapters",
				Benefits:    []string{"Highly testable", "Technology agnostic", "Clear boundaries"},
				Drawbacks:   []string{"Complex setup", "Over-engineering risk"},
			},
		}
	case "cli":
		architectures = []models.ArchitectureInfo{
			{
				ID:          "standard",
				Name:        "Standard",
				Description: "Simple command structure with cobra framework",
				Complexity:  "beginner",
				UseCase:     "CLI tools and utilities",
				Benefits:    []string{"Simple", "Easy to understand", "Quick to implement"},
				Drawbacks:   []string{"Limited for complex CLIs"},
			},
		}
	default:
		architectures = []models.ArchitectureInfo{}
	}
	
	response := models.NewSuccessResponse(map[string]interface{}{
		"projectType":   projectType,
		"architectures": architectures,
		"count":         len(architectures),
	})
	
	c.JSON(http.StatusOK, response)
}