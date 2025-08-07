package generator

import (
	"github.com/jashkahar/open-workbench-platform/internal/manifest"
)

// Generator represents a deployment target generator (Docker, Terraform, etc.)
type Generator interface {
	// Name returns the unique identifier for this generator
	Name() string

	// Description returns a human-readable description of this generator
	Description() string

	// Generate creates the deployment configuration for the given manifest
	Generate(manifest *manifest.WorkbenchManifest) error

	// Validate checks if the manifest is compatible with this generator
	Validate(manifest *manifest.WorkbenchManifest) error
}

// GeneratorResult represents the output of a generator
type GeneratorResult struct {
	Files    map[string][]byte // filename -> content
	Messages []string          // informational messages
	Warnings []string          // warning messages
	Errors   []string          // error messages
}

// Registry manages all available generators
type Registry interface {
	// Register adds a new generator to the registry
	Register(generator Generator) error

	// Get retrieves a generator by name
	Get(name string) (Generator, error)

	// List returns all available generators
	List() []Generator

	// Names returns all available generator names
	Names() []string
}
