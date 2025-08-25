package templates

import (
	"io/fs"

	"github.com/francknouama/go-starter/pkg/embedfs"
)

// templatesFS holds a custom filesystem set by tests
var templatesFS fs.FS

// SetTemplatesFS sets a custom filesystem (primarily used by tests)
func SetTemplatesFS(fs fs.FS) {
	templatesFS = fs
}

// GetTemplatesFS returns the filesystem for templates
func GetTemplatesFS() fs.FS {
	// If a custom filesystem is set (for tests), use it
	if templatesFS != nil {
		// Check if we need to strip the "blueprints" prefix
		if _, err := fs.Stat(templatesFS, "blueprints"); err == nil {
			// This is likely an embedded FS with "blueprints" directory
			subFS, err := fs.Sub(templatesFS, "blueprints")
			if err != nil {
				panic("failed to create sub-filesystem for blueprints: " + err.Error())
			}
			return subFS
		}
		// This is likely a DirFS pointing directly to templates directory
		return templatesFS
	}

	// Use the shared embedded/filesystem source
	return embedfs.GetBlueprintsFS()
}
