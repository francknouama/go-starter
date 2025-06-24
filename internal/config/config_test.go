package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig

	// Check that default profile exists
	if _, exists := config.Profiles["default"]; !exists {
		t.Error("Default profile should exist")
	}

	// Check current profile is set
	if config.CurrentProfile != "default" {
		t.Errorf("Expected current profile to be 'default', got '%s'", config.CurrentProfile)
	}

	// Check default values
	defaultProfile := config.Profiles["default"]
	if defaultProfile.License != "MIT" {
		t.Errorf("Expected default license to be 'MIT', got '%s'", defaultProfile.License)
	}

	if defaultProfile.Defaults.GoVersion != "1.21" {
		t.Errorf("Expected default Go version to be '1.21', got '%s'", defaultProfile.Defaults.GoVersion)
	}

	if defaultProfile.Defaults.Framework != "gin" {
		t.Errorf("Expected default framework to be 'gin', got '%s'", defaultProfile.Defaults.Framework)
	}
}

func TestConfigLoad_NoFile(t *testing.T) {
	// Test loading config when no file exists
	config, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("Expected no error when config file doesn't exist, got: %v", err)
	}

	// Should return default config
	if config.CurrentProfile != "default" {
		t.Errorf("Expected current profile to be 'default', got '%s'", config.CurrentProfile)
	}

	if _, exists := config.Profiles["default"]; !exists {
		t.Error("Default profile should exist")
	}
}

func TestConfigSaveLoad(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "go-starter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "config.yaml")

	// Create test config
	config := &Config{
		Profiles: map[string]Profile{
			"test": {
				Author:  "Test Author",
				Email:   "test@example.com",
				License: "Apache-2.0",
				Defaults: ProfileDefaults{
					GoVersion:    "1.20",
					Framework:    "echo",
					Architecture: "clean",
				},
			},
		},
		CurrentProfile: "test",
	}

	// Save config
	err = config.Save(configFile)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedConfig, err := Load(configFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config
	if loadedConfig.CurrentProfile != "test" {
		t.Errorf("Expected current profile to be 'test', got '%s'", loadedConfig.CurrentProfile)
	}

	testProfile, exists := loadedConfig.Profiles["test"]
	if !exists {
		t.Fatal("Test profile should exist")
	}

	if testProfile.Author != "Test Author" {
		t.Errorf("Expected author to be 'Test Author', got '%s'", testProfile.Author)
	}

	if testProfile.Email != "test@example.com" {
		t.Errorf("Expected email to be 'test@example.com', got '%s'", testProfile.Email)
	}

	if testProfile.License != "Apache-2.0" {
		t.Errorf("Expected license to be 'Apache-2.0', got '%s'", testProfile.License)
	}
}

func TestGetCurrentProfile(t *testing.T) {
	config := &Config{
		Profiles: map[string]Profile{
			"test": {
				Author: "Test Author",
			},
		},
		CurrentProfile: "test",
	}

	profile, err := config.GetCurrentProfile()
	if err != nil {
		t.Fatalf("Failed to get current profile: %v", err)
	}

	if profile.Author != "Test Author" {
		t.Errorf("Expected author to be 'Test Author', got '%s'", profile.Author)
	}
}

func TestGetCurrentProfile_NotFound(t *testing.T) {
	config := &Config{
		Profiles:       map[string]Profile{},
		CurrentProfile: "nonexistent",
	}

	_, err := config.GetCurrentProfile()
	if err == nil {
		t.Error("Expected error when current profile doesn't exist")
	}
}

func TestSetCurrentProfile(t *testing.T) {
	config := &Config{
		Profiles: map[string]Profile{
			"test1": {},
			"test2": {},
		},
		CurrentProfile: "test1",
	}

	err := config.SetCurrentProfile("test2")
	if err != nil {
		t.Fatalf("Failed to set current profile: %v", err)
	}

	if config.CurrentProfile != "test2" {
		t.Errorf("Expected current profile to be 'test2', got '%s'", config.CurrentProfile)
	}
}

func TestSetCurrentProfile_NotFound(t *testing.T) {
	config := &Config{
		Profiles: map[string]Profile{
			"test1": {},
		},
		CurrentProfile: "test1",
	}

	err := config.SetCurrentProfile("nonexistent")
	if err == nil {
		t.Error("Expected error when setting nonexistent profile")
	}

	// Current profile should remain unchanged
	if config.CurrentProfile != "test1" {
		t.Errorf("Current profile should remain 'test1', got '%s'", config.CurrentProfile)
	}
}

func TestAddProfile(t *testing.T) {
	config := &Config{
		Profiles: map[string]Profile{},
	}

	profile := Profile{
		Author: "New Author",
	}

	config.AddProfile("new", profile)

	addedProfile, exists := config.Profiles["new"]
	if !exists {
		t.Error("Profile should be added")
		return
	}

	if addedProfile.Author != "New Author" {
		t.Errorf("Expected author to be 'New Author', got '%s'", addedProfile.Author)
	}
}

func TestRemoveProfile(t *testing.T) {
	config := &Config{
		Profiles: map[string]Profile{
			"default": {},
			"test":    {},
		},
		CurrentProfile: "test",
	}

	err := config.RemoveProfile("test")
	if err != nil {
		t.Fatalf("Failed to remove profile: %v", err)
	}

	if _, exists := config.Profiles["test"]; exists {
		t.Error("Profile should be removed")
	}

	// Current profile should switch to default
	if config.CurrentProfile != "default" {
		t.Errorf("Current profile should switch to 'default', got '%s'", config.CurrentProfile)
	}
}

func TestRemoveProfile_Default(t *testing.T) {
	config := &Config{
		Profiles: map[string]Profile{
			"default": {},
		},
		CurrentProfile: "default",
	}

	err := config.RemoveProfile("default")
	if err == nil {
		t.Error("Should not be able to remove default profile")
	}
}

func TestRemoveProfile_NotFound(t *testing.T) {
	config := &Config{
		Profiles: map[string]Profile{
			"default": {},
		},
		CurrentProfile: "default",
	}

	err := config.RemoveProfile("nonexistent")
	if err == nil {
		t.Error("Expected error when removing nonexistent profile")
	}
}

func TestApplyDefaults(t *testing.T) {
	config := &Config{
		Profiles: map[string]Profile{
			"test": {
				Author: "Test Author",
				// License and Defaults will be empty, should get defaults
			},
		},
		CurrentProfile: "test",
	}

	err := config.applyDefaults()
	if err != nil {
		t.Fatalf("Failed to apply defaults: %v", err)
	}

	testProfile := config.Profiles["test"]

	// Should have default license
	if testProfile.License != "MIT" {
		t.Errorf("Expected default license 'MIT', got '%s'", testProfile.License)
	}

	// Should have default Go version
	if testProfile.Defaults.GoVersion != "1.21" {
		t.Errorf("Expected default Go version '1.21', got '%s'", testProfile.Defaults.GoVersion)
	}

	// Should preserve existing author
	if testProfile.Author != "Test Author" {
		t.Errorf("Expected author to be preserved as 'Test Author', got '%s'", testProfile.Author)
	}
}

func TestValidateProfile(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
		wantErr bool
	}{
		{
			name: "valid profile",
			profile: Profile{
				Author:  "Test Author",
				Email:   "test@example.com",
				License: "MIT",
				Defaults: ProfileDefaults{
					GoVersion:    "1.21",
					Framework:    "gin",
					Architecture: "standard",
					Database: DatabaseDefaults{
						Driver: "postgres",
						ORM:    "gorm",
					},
					Auth: AuthDefaults{
						Type: "jwt",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid go version",
			profile: Profile{
				Defaults: ProfileDefaults{
					GoVersion: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid framework",
			profile: Profile{
				Defaults: ProfileDefaults{
					GoVersion: "1.21",
					Framework: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid architecture",
			profile: Profile{
				Defaults: ProfileDefaults{
					GoVersion:    "1.21",
					Framework:    "gin",
					Architecture: "invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProfile("test", tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
