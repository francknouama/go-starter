package prompts

import (
	"fmt"
	"slices"
)

// Supported Go versions (latest 3 major versions)
var supportedGoVersions = []string{"auto", "1.23", "1.22", "1.21"}

// GetSupportedGoVersions returns the list of supported Go versions
func GetSupportedGoVersions() []string {
	return supportedGoVersions
}

// mapSelectionToVersion converts user selection to version string
func mapSelectionToVersion(selection string) string {
	switch selection {
	case "Auto-detect (recommended)":
		return "auto"
	case "Go 1.23 (latest)":
		return "1.23"
	case "Go 1.22":
		return "1.22"
	case "Go 1.21":
		return "1.21"
	default:
		return "auto"
	}
}


// ValidateGoVersion validates if the provided Go version is supported
func ValidateGoVersion(version string) error {
	if slices.Contains(supportedGoVersions, version) {
		return nil
	}
	return fmt.Errorf("unsupported Go version: %s. Supported versions: %v", version, supportedGoVersions)
}
