// Package main provides the Open Workbench CLI, a command-line tool for scaffolding
// modern web applications with pre-configured templates and best practices.
package main

import (
	"embed"

	"github.com/jashkahar/open-workbench-platform/cmd"
)

// templatesFS embeds the templates directory into the binary.
// This allows the CLI to be distributed as a single executable
// without requiring external template files.
//
//go:embed templates
var templatesFS embed.FS

func main() {
	cmd.Execute(templatesFS)
}
