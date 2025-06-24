package templates

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

// templatesFS embeds all template files from the templates directory
// Note: Due to Go embed limitations, we need to use relative paths from the package location
// For now, we'll use the filesystem directly for development
var templatesFS embed.FS

// GetTemplatesFS returns the filesystem for templates
func GetTemplatesFS() fs.FS {
	// For development, use the actual filesystem
	// In production, this would use the embedded filesystem
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	templatesDir := filepath.Join(projectRoot, "templates")
	
	// Check if templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		panic("templates directory not found at: " + templatesDir)
	}
	
	return os.DirFS(templatesDir)
}