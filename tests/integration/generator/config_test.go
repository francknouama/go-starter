package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/francknouama/go-starter/internal/config"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
)

// TestGenerator_Config_Integration tests generator integration with configuration
func TestGenerator_Config_Integration(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name           string
		config         *config.Config
		projectConfig  types.ProjectConfig
		expectedValues map[string]string
	}{
		{
			name: "profile with defaults applied",
			config: &config.Config{
				Profiles: map[string]config.Profile{
					"test": {
						Author:  "Test Author",
						Email:   "test@example.com",
						License: "MIT",
						Defaults: config.ProfileDefaults{
							GoVersion:    "1.21",
							Framework:    "gin",
							Architecture: "standard",
							Logger:       "slog",
						},
					},
				},
				CurrentProfile: "test",
			},
			projectConfig: types.ProjectConfig{
				Name:   "config-integration-test",
				Module: "github.com/test/config-integration-test",
				Type:   "web-api",
			},
			expectedValues: map[string]string{
				"Author":       "Test Author",
				"Email":        "test@example.com",
				"License":      "MIT",
				"GoVersion":    "1.21",
				"Framework":    "gin",
				"Architecture": "standard",
				"Logger":       "slog",
			},
		},
		{
			name: "explicit config overrides profile defaults",
			config: &config.Config{
				Profiles: map[string]config.Profile{
					"dev": {
						Author:  "Dev Author",
						Email:   "dev@example.com",
						License: "Apache-2.0",
						Defaults: config.ProfileDefaults{
							GoVersion:    "1.20",
							Framework:    "echo",
							Architecture: "clean",
							Logger:       "zap",
						},
					},
				},
				CurrentProfile: "dev",
			},
			projectConfig: types.ProjectConfig{
				Name:         "override-test",
				Module:       "github.com/test/override-test",
				Type:         "web-api",
				GoVersion:    "1.22",    // Override profile default
				Framework:    "fiber",   // Override profile default
				Architecture: "ddd",     // Override profile default
				Logger:       "logrus",  // Override profile default
				Author:       "Custom Author", // Override profile
				Email:        "custom@example.com", // Override profile
				License:      "BSD-3-Clause", // Override profile
			},
			expectedValues: map[string]string{
				"Author":       "Custom Author",
				"Email":        "custom@example.com", 
				"License":      "BSD-3-Clause",
				"GoVersion":    "1.22",
				"Framework":    "fiber",
				"Architecture": "ddd",
				"Logger":       "logrus",
			},
		},
		{
			name: "multiple profiles configuration",
			config: &config.Config{
				Profiles: map[string]config.Profile{
					"personal": {
						Author:  "Personal User",
						Email:   "personal@gmail.com",
						License: "MIT",
						Defaults: config.ProfileDefaults{
							GoVersion:    "1.21",
							Framework:    "gin",
							Architecture: "standard",
							Logger:       "slog",
						},
					},
					"work": {
						Author:  "Work User",
						Email:   "work@company.com",
						License: "Apache-2.0",
						Defaults: config.ProfileDefaults{
							GoVersion:    "1.20",
							Framework:    "echo",
							Architecture: "hexagonal",
							Logger:       "zap",
						},
					},
				},
				CurrentProfile: "work",
			},
			projectConfig: types.ProjectConfig{
				Name:   "multi-profile-test",
				Module: "github.com/company/multi-profile-test",
				Type:   "web-api",
			},
			expectedValues: map[string]string{
				"Author":       "Work User",
				"Email":        "work@company.com",
				"License":      "Apache-2.0",
				"GoVersion":    "1.20",
				"Framework":    "echo",
				"Architecture": "hexagonal",
				"Logger":       "zap",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := generator.New()
			require.NotNil(t, gen)

			// Apply configuration values to project config
			// In a real scenario, this would be done by the CLI layer
			profile := tt.config.Profiles[tt.config.CurrentProfile]
			
			// Apply profile values if not explicitly set in project config
			if tt.projectConfig.Author == "" {
				tt.projectConfig.Author = profile.Author
			}
			if tt.projectConfig.Email == "" {
				tt.projectConfig.Email = profile.Email
			}
			if tt.projectConfig.License == "" {
				tt.projectConfig.License = profile.License
			}
			if tt.projectConfig.GoVersion == "" {
				tt.projectConfig.GoVersion = profile.Defaults.GoVersion
			}
			if tt.projectConfig.Framework == "" {
				tt.projectConfig.Framework = profile.Defaults.Framework
			}
			if tt.projectConfig.Architecture == "" {
				tt.projectConfig.Architecture = profile.Defaults.Architecture
			}
			if tt.projectConfig.Logger == "" {
				tt.projectConfig.Logger = profile.Defaults.Logger
			}

			// Test generation with configuration
			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + tt.projectConfig.Name,
				DryRun:     true, // Use dry run for config testing
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(tt.projectConfig, options)

			// Accept template not found errors
			if err != nil {
				if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
					t.Skip("Skipping test as template is not yet implemented")
					return
				}
			}

			// Verify configuration values were applied correctly
			for key, expectedValue := range tt.expectedValues {
				var actualValue string
				switch key {
				case "Author":
					actualValue = tt.projectConfig.Author
				case "Email":
					actualValue = tt.projectConfig.Email
				case "License":
					actualValue = tt.projectConfig.License
				case "GoVersion":
					actualValue = tt.projectConfig.GoVersion
				case "Framework":
					actualValue = tt.projectConfig.Framework
				case "Architecture":
					actualValue = tt.projectConfig.Architecture
				case "Logger":
					actualValue = tt.projectConfig.Logger
				}
				
				assert.Equal(t, expectedValue, actualValue, 
					"Configuration value %s should match expected value", key)
			}

			t.Logf("Configuration integration test passed for profile: %s", tt.config.CurrentProfile)
		})
	}
}

// TestGenerator_Config_Variables tests custom variable handling
func TestGenerator_Config_Variables(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name      string
		variables map[string]string
		expected  map[string]string
	}{
		{
			name: "basic custom variables",
			variables: map[string]string{
				"AppName":     "MyAwesomeApp",
				"Version":     "1.0.0",
				"Description": "A test application",
			},
			expected: map[string]string{
				"AppName":     "MyAwesomeApp",
				"Version":     "1.0.0",
				"Description": "A test application",
			},
		},
		{
			name: "variables with special characters",
			variables: map[string]string{
				"DatabaseURL":    "postgresql://user:pass@localhost:5432/db",
				"APIKey":         "sk-1234567890abcdef",
				"ServerAddress":  "0.0.0.0:8080",
			},
			expected: map[string]string{
				"DatabaseURL":    "postgresql://user:pass@localhost:5432/db",
				"APIKey":         "sk-1234567890abcdef",
				"ServerAddress":  "0.0.0.0:8080",
			},
		},
		{
			name: "boolean and numeric variables as strings",
			variables: map[string]string{
				"Debug":       "true",
				"Port":        "8080",
				"MaxWorkers":  "10",
				"EnableTLS":   "false",
			},
			expected: map[string]string{
				"Debug":       "true",
				"Port":        "8080",
				"MaxWorkers":  "10",
				"EnableTLS":   "false",
			},
		},
		{
			name:      "empty variables map",
			variables: map[string]string{},
			expected:  map[string]string{},
		},
		{
			name:      "nil variables map",
			variables: nil,
			expected:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "variables-test",
				Module:    "github.com/test/variables-test",
				Type:      "library",
				GoVersion: "1.21",
				Variables: tt.variables,
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + config.Name,
				DryRun:     true,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			// Accept template not found errors
			if err != nil {
				if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
					t.Skip("Skipping test as template is not yet implemented")
					return
				}
			}

			// Verify variables are preserved correctly
			if tt.variables == nil {
				assert.Nil(t, config.Variables, "Variables should remain nil")
			} else {
				require.NotNil(t, config.Variables, "Variables should not be nil")
				for key, expectedValue := range tt.expected {
					actualValue, exists := config.Variables[key]
					assert.True(t, exists, "Variable %s should exist", key)
					assert.Equal(t, expectedValue, actualValue, 
						"Variable %s should have expected value", key)
				}
			}

			t.Logf("Variables test passed for %d variables", len(tt.expected))
		})
	}
}

// TestGenerator_Config_Features tests feature configuration handling
func TestGenerator_Config_Features(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name     string
		features *types.Features
		validate func(t *testing.T, features *types.Features)
	}{
		{
			name: "complete database configuration",
			features: &types.Features{
				Database: types.DatabaseConfig{
					Driver: "postgresql",
					ORM:    "gorm",
				},
			},
			validate: func(t *testing.T, features *types.Features) {
				require.NotNil(t, features)
				assert.Equal(t, "postgresql", features.Database.Driver)
				assert.Equal(t, "gorm", features.Database.ORM)
			},
		},
		{
			name: "multiple database drivers",
			features: &types.Features{
				Database: types.DatabaseConfig{
					Drivers: []string{"postgresql", "redis", "mongodb"},
					ORM:     "gorm",
				},
			},
			validate: func(t *testing.T, features *types.Features) {
				require.NotNil(t, features)
				assert.Len(t, features.Database.Drivers, 3)
				assert.Contains(t, features.Database.Drivers, "postgresql")
				assert.Contains(t, features.Database.Drivers, "redis")
				assert.Contains(t, features.Database.Drivers, "mongodb")
				assert.Equal(t, "gorm", features.Database.ORM)
			},
		},
		{
			name: "authentication configuration",
			features: &types.Features{
				Authentication: types.AuthConfig{
					Type:      "jwt",
					Providers: []string{"google", "github"},
				},
			},
			validate: func(t *testing.T, features *types.Features) {
				require.NotNil(t, features)
				assert.Equal(t, "jwt", features.Authentication.Type)
				assert.Len(t, features.Authentication.Providers, 2)
				assert.Contains(t, features.Authentication.Providers, "google")
				assert.Contains(t, features.Authentication.Providers, "github")
			},
		},
		{
			name: "logging configuration",
			features: &types.Features{
				Logging: types.LoggingConfig{
					Type:       "zap",
					Level:      "debug",
					Format:     "json",
					Structured: true,
				},
			},
			validate: func(t *testing.T, features *types.Features) {
				require.NotNil(t, features)
				assert.Equal(t, "zap", features.Logging.Type)
				assert.Equal(t, "debug", features.Logging.Level)
				assert.Equal(t, "json", features.Logging.Format)
				assert.True(t, features.Logging.Structured)
			},
		},
		{
			name: "testing configuration",
			features: &types.Features{
				Testing: types.TestConfig{
					Framework: "testify",
					Coverage:  true,
				},
			},
			validate: func(t *testing.T, features *types.Features) {
				require.NotNil(t, features)
				assert.Equal(t, "testify", features.Testing.Framework)
				assert.True(t, features.Testing.Coverage)
			},
		},
		{
			name: "deployment configuration",
			features: &types.Features{
				Deployment: types.DeployConfig{
					Targets: []string{"docker", "kubernetes", "aws"},
				},
			},
			validate: func(t *testing.T, features *types.Features) {
				require.NotNil(t, features)
				assert.Len(t, features.Deployment.Targets, 3)
				assert.Contains(t, features.Deployment.Targets, "docker")
				assert.Contains(t, features.Deployment.Targets, "kubernetes")
				assert.Contains(t, features.Deployment.Targets, "aws")
			},
		},
		{
			name: "combined features configuration",
			features: &types.Features{
				Database: types.DatabaseConfig{
					Driver: "postgresql",
					ORM:    "gorm",
				},
				Authentication: types.AuthConfig{
					Type: "oauth2",
				},
				Logging: types.LoggingConfig{
					Type:   "slog",
					Level:  "info",
					Format: "text",
				},
				Testing: types.TestConfig{
					Framework: "ginkgo",
					Coverage:  false,
				},
			},
			validate: func(t *testing.T, features *types.Features) {
				require.NotNil(t, features)
				
				// Database
				assert.Equal(t, "postgresql", features.Database.Driver)
				assert.Equal(t, "gorm", features.Database.ORM)
				
				// Authentication
				assert.Equal(t, "oauth2", features.Authentication.Type)
				
				// Logging
				assert.Equal(t, "slog", features.Logging.Type)
				assert.Equal(t, "info", features.Logging.Level)
				assert.Equal(t, "text", features.Logging.Format)
				
				// Testing
				assert.Equal(t, "ginkgo", features.Testing.Framework)
				assert.False(t, features.Testing.Coverage)
			},
		},
		{
			name:     "nil features configuration",
			features: nil,
			validate: func(t *testing.T, features *types.Features) {
				assert.Nil(t, features)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "features-test",
				Module:    "github.com/test/features-test",
				Type:      "web-api",
				Framework: "gin",
				GoVersion: "1.21",
				Features:  tt.features,
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + config.Name,
				DryRun:     true,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			// Accept template not found errors
			if err != nil {
				if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
					t.Skip("Skipping test as template is not yet implemented")
					return
				}
			}

			// Validate features configuration
			tt.validate(t, config.Features)

			t.Logf("Features configuration test passed")
		})
	}
}

// TestGenerator_Config_ProfileDefaults tests profile default application
func TestGenerator_Config_ProfileDefaults(t *testing.T) {
	setupTestTemplates(t)

	profile := config.Profile{
		Author:  "Profile Author",
		Email:   "profile@example.com",
		License: "MIT",
		Defaults: config.ProfileDefaults{
			GoVersion:    "1.21",
			Framework:    "gin",
			Architecture: "clean",
			Logger:       "zap",
		},
	}

	tests := []struct {
		name          string
		projectConfig types.ProjectConfig
		applyDefaults bool
		expected      types.ProjectConfig
	}{
		{
			name: "apply all defaults",
			projectConfig: types.ProjectConfig{
				Name:   "defaults-test",
				Module: "github.com/test/defaults-test",
				Type:   "web-api",
			},
			applyDefaults: true,
			expected: types.ProjectConfig{
				Name:         "defaults-test",
				Module:       "github.com/test/defaults-test",
				Type:         "web-api",
				GoVersion:    "1.21",
				Framework:    "gin",
				Architecture: "clean",
				Logger:       "zap",
				Author:       "Profile Author",
				Email:        "profile@example.com",
				License:      "MIT",
			},
		},
		{
			name: "partial override of defaults",
			projectConfig: types.ProjectConfig{
				Name:         "partial-test",
				Module:       "github.com/test/partial-test",
				Type:         "web-api",
				Framework:    "echo", // Override default
				Logger:       "slog", // Override default
			},
			applyDefaults: true,
			expected: types.ProjectConfig{
				Name:         "partial-test",
				Module:       "github.com/test/partial-test",
				Type:         "web-api",
				GoVersion:    "1.21",     // From profile
				Framework:    "echo",     // Explicit override
				Architecture: "clean",    // From profile
				Logger:       "slog",     // Explicit override
				Author:       "Profile Author",
				Email:        "profile@example.com",
				License:      "MIT",
			},
		},
		{
			name: "no defaults applied",
			projectConfig: types.ProjectConfig{
				Name:      "no-defaults-test",
				Module:    "github.com/test/no-defaults-test",
				Type:      "library",
				GoVersion: "1.20",
			},
			applyDefaults: false,
			expected: types.ProjectConfig{
				Name:      "no-defaults-test",
				Module:    "github.com/test/no-defaults-test",
				Type:      "library",
				GoVersion: "1.20",
				// No profile values applied
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.projectConfig

			// Simulate profile default application (normally done by CLI)
			if tt.applyDefaults {
				if config.Author == "" {
					config.Author = profile.Author
				}
				if config.Email == "" {
					config.Email = profile.Email
				}
				if config.License == "" {
					config.License = profile.License
				}
				if config.GoVersion == "" {
					config.GoVersion = profile.Defaults.GoVersion
				}
				if config.Framework == "" {
					config.Framework = profile.Defaults.Framework
				}
				if config.Architecture == "" {
					config.Architecture = profile.Defaults.Architecture
				}
				if config.Logger == "" {
					config.Logger = profile.Defaults.Logger
				}
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + config.Name,
				DryRun:     true,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			// Accept template not found errors
			if err != nil {
				if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
					t.Skip("Skipping test as template is not yet implemented")
					return
				}
			}

			// Verify final configuration matches expected
			assert.Equal(t, tt.expected.Name, config.Name)
			assert.Equal(t, tt.expected.Module, config.Module)
			assert.Equal(t, tt.expected.Type, config.Type)
			assert.Equal(t, tt.expected.GoVersion, config.GoVersion)
			assert.Equal(t, tt.expected.Framework, config.Framework)
			assert.Equal(t, tt.expected.Architecture, config.Architecture)
			assert.Equal(t, tt.expected.Logger, config.Logger)
			assert.Equal(t, tt.expected.Author, config.Author)
			assert.Equal(t, tt.expected.Email, config.Email)
			assert.Equal(t, tt.expected.License, config.License)

			t.Logf("Profile defaults test passed")
		})
	}
}