package embedfs

import (
	"io/fs"
	"os"
)

// Global variable to store embedded blueprints (set by root binary)
var embeddedBlueprintsFS fs.FS

// SetBlueprintsFS sets the embedded blueprints filesystem (called by root binary)
func SetBlueprintsFS(fs fs.FS) {
	embeddedBlueprintsFS = fs
}

// GetBlueprintsFS returns the blueprints filesystem
// For cmd/go-starter binary, it reads from filesystem
// For root binary, it uses embedded blueprints
func GetBlueprintsFS() fs.FS {
	// If embedded blueprints are available (from root binary), use them
	if embeddedBlueprintsFS != nil {
		return embeddedBlueprintsFS
	}
	
	// First try to read from filesystem (development/cmd binary)
	if _, err := os.Stat("blueprints"); err == nil {
		return os.DirFS("blueprints")
	}
	
	// Try from parent directory (when running from subdirectory)
	if _, err := os.Stat("../blueprints"); err == nil {
		return os.DirFS("../blueprints")
	}
	
	// Try from project root (when running from deeper subdirectories)
	if _, err := os.Stat("../../blueprints"); err == nil {
		return os.DirFS("../../blueprints")
	}
	
	// If no filesystem blueprints found, return nil
	// The caller should handle this case
	return nil
}