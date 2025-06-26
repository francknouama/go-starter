package templates

import (
	"io/fs"
)

// templatesFS holds the embedded filesystem set by the main package
var templatesFS fs.FS

// SetTemplatesFS sets the embedded filesystem (called from main package)
func SetTemplatesFS(fs fs.FS) {
	templatesFS = fs
}

// GetTemplatesFS returns the filesystem for templates
func GetTemplatesFS() fs.FS {
	if templatesFS == nil {
		panic("templates filesystem not initialized - ensure SetTemplatesFS is called from main")
	}
	
	// When embedded from module root, we need to strip the "templates" prefix
	subFS, err := fs.Sub(templatesFS, "templates")
	if err != nil {
		panic("failed to create sub-filesystem for templates: " + err.Error())
	}
	return subFS
}