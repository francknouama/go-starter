//go:build !noembedtemplates
// +build !noembedtemplates

package main

import (
	"embed"
)

// TemplatesFS embeds all template files from the templates directory
// The all: prefix ensures dot files and other special files are included
//
//go:embed all:blueprints
var TemplatesFS embed.FS
