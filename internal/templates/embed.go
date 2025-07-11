package templates

import (
	"io/fs"
)

// templatesFS holds the embedded filesystem set by the main package
var templatesFS fs.FS

// SetTemplatesFS sets the embedded filesystem (called from main package or tests)
func SetTemplatesFS(fs fs.FS) {
	templatesFS = fs
}

// GetTemplatesFS returns the filesystem for templates
func GetTemplatesFS() fs.FS {
	if templatesFS == nil {
		panic("templates filesystem not initialized - ensure SetTemplatesFS is called from main")
	}

	// Check if we need to strip the "blueprints" prefix
	// For embedded FS from root, we need to strip it
	// For test DirFS pointing directly to blueprints, we don't
	if _, err := fs.Stat(templatesFS, "blueprints"); err == nil {
		// This is likely the embedded FS with "blueprints" directory
		subFS, err := fs.Sub(templatesFS, "blueprints")
		if err != nil {
			panic("failed to create sub-filesystem for blueprints: " + err.Error())
		}
		return subFS
	}

	// This is likely a DirFS pointing directly to templates directory
	return templatesFS
}
