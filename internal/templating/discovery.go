// Package templating provides the core templating system for the Open Workbench CLI.
// This package implements dynamic template discovery, parameter processing, and
// file generation capabilities with support for conditional logic and validation.
package templating

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
)

// TemplateManifest represents the structure of a template.json file.
// This struct defines the complete configuration for a template, including
// metadata, parameters, validation rules, and post-scaffolding actions.
//
// The manifest is loaded from JSON files in each template directory and
// provides the blueprint for template processing and parameter collection.
type TemplateManifest struct {
	Name         string        `json:"name"`                   // Display name for the template
	Description  string        `json:"description"`            // Human-readable description
	Parameters   []Parameter   `json:"parameters"`             // List of parameters to collect
	PostScaffold *PostScaffold `json:"postScaffold,omitempty"` // Post-processing actions
}

// Parameter represents a single parameter that the user needs to provide.
// This struct defines the configuration for collecting user input during
// the template scaffolding process, including validation rules and UI options.
type Parameter struct {
	Name       string      `json:"name"`                 // Unique parameter identifier
	Prompt     string      `json:"prompt"`               // User-facing question text
	HelpText   string      `json:"helpText,omitempty"`   // Additional help information
	Group      string      `json:"group,omitempty"`      // Group for organizing parameters
	Type       string      `json:"type"`                 // Parameter type (string, boolean, select, multiselect)
	Required   bool        `json:"required,omitempty"`   // Whether parameter is mandatory
	Default    any         `json:"default,omitempty"`    // Default value for the parameter
	Options    []string    `json:"options,omitempty"`    // Available options for select/multiselect
	Condition  string      `json:"condition,omitempty"`  // Conditional visibility rule
	Validation *Validation `json:"validation,omitempty"` // Validation rules for the parameter
}

// Validation represents validation rules for string parameters.
// This struct defines regex patterns and error messages for validating
// user input during parameter collection.
type Validation struct {
	Regex        string `json:"regex"`        // Regular expression pattern for validation
	ErrorMessage string `json:"errorMessage"` // Custom error message for validation failures
}

// PostScaffold represents actions to perform after scaffolding.
// This struct defines cleanup and setup actions that should be executed
// after the main template processing is complete.
type PostScaffold struct {
	FilesToDelete []FileAction    `json:"filesToDelete,omitempty"` // Files to delete based on conditions
	Commands      []CommandAction `json:"commands,omitempty"`      // Commands to execute after scaffolding
}

// FileAction represents a file or directory to delete based on a condition.
// This struct defines conditional file deletion rules that can be applied
// after template processing based on user parameter choices.
type FileAction struct {
	Path      string `json:"path"`      // Path to the file or directory to delete
	Condition string `json:"condition"` // Condition that must be true for deletion
}

// CommandAction represents a command to execute after scaffolding.
// This struct defines shell commands that should be run after template
// processing, such as dependency installation or git initialization.
type CommandAction struct {
	Command     string `json:"command"`             // Shell command to execute
	Description string `json:"description"`         // Human-readable description of the command
	Condition   string `json:"condition,omitempty"` // Optional condition for execution
}

// TemplateInfo represents metadata about a discovered template.
// This struct contains information about a template discovered during
// the template discovery process, including its manifest and metadata.
type TemplateInfo struct {
	Name        string            // The template name (directory name)
	Description string            // Human-readable description from manifest
	Path        string            // Path to the template directory
	Manifest    *TemplateManifest // Parsed template manifest
}

// LoadTemplateManifest loads and parses a template.json file from the embedded filesystem.
// This function reads the JSON manifest file for a specific template and validates
// its structure to ensure it contains all required fields.
//
// Parameters:
//   - templateFS: The embedded filesystem containing template files
//   - templateName: The name of the template directory to load
//
// Returns:
//   - A pointer to the parsed TemplateManifest
//   - An error if the manifest cannot be loaded or parsed
func LoadTemplateManifest(templateFS fs.FS, templateName string) (*TemplateManifest, error) {
	// Construct the path to the template manifest file
	// Use forward slashes for embedded filesystem paths
	manifestPath := fmt.Sprintf("templates/%s/template.json", templateName)

	fmt.Printf("DEBUG: Trying to load manifest from: %s\n", manifestPath)

	// Read the manifest file from the embedded filesystem
	manifestBytes, err := fs.ReadFile(templateFS, manifestPath)
	if err != nil {
		fmt.Printf("DEBUG: Failed to read manifest: %v\n", err)
		return nil, fmt.Errorf("failed to read template manifest: %w", err)
	}

	fmt.Printf("DEBUG: Successfully read manifest bytes: %d bytes\n", len(manifestBytes))

	// Parse the JSON manifest into the TemplateManifest struct
	var manifest TemplateManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		fmt.Printf("DEBUG: Failed to parse manifest JSON: %v\n", err)
		return nil, fmt.Errorf("failed to parse template manifest: %w", err)
	}

	fmt.Printf("DEBUG: Successfully parsed manifest: %s - %s\n", manifest.Name, manifest.Description)

	// Validate that all required fields are present
	if manifest.Name == "" {
		return nil, fmt.Errorf("template manifest missing required field: name")
	}
	if manifest.Description == "" {
		return nil, fmt.Errorf("template manifest missing required field: description")
	}
	if len(manifest.Parameters) == 0 {
		return nil, fmt.Errorf("template manifest missing required field: parameters")
	}

	return &manifest, nil
}

// DiscoverTemplates finds all available templates in the embedded filesystem.
// This function scans the templates directory and loads the manifest for each
// template directory, providing a complete list of available templates with
// their metadata and configuration.
//
// Parameters:
//   - templateFS: The embedded filesystem containing template directories
//
// Returns:
//   - A slice of TemplateInfo structs containing template metadata
//   - An error if template discovery fails
func DiscoverTemplates(templateFS fs.FS) ([]TemplateInfo, error) {
	var templates []TemplateInfo

	// Read the templates directory to find all template subdirectories
	entries, err := fs.ReadDir(templateFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	// Process each entry in the templates directory
	for _, entry := range entries {
		// Only process directories (skip files)
		if entry.IsDir() {
			templateName := entry.Name()

			// Try to load the manifest for this template
			manifest, err := LoadTemplateManifest(templateFS, templateName)
			if err != nil {
				// Skip templates with invalid manifests for now
				// In the future, we might want to log this or handle it differently
				continue
			}

			// Create TemplateInfo struct with template metadata
			templateInfo := TemplateInfo{
				Name:        templateName,
				Description: manifest.Description,
				Path:        fmt.Sprintf("templates/%s", templateName),
				Manifest:    manifest,
			}

			templates = append(templates, templateInfo)
		}
	}

	// Sort templates by name for consistent ordering
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

// GetTemplateInfo returns information about a specific template.
// This function loads and validates a template's manifest, providing
// complete information about the template including its configuration.
//
// Parameters:
//   - templateFS: The embedded filesystem containing template files
//   - templateName: The name of the template to load
//
// Returns:
//   - A pointer to TemplateInfo containing template metadata
//   - An error if the template cannot be loaded
func GetTemplateInfo(templateFS fs.FS, templateName string) (*TemplateInfo, error) {
	// Load the template manifest
	manifest, err := LoadTemplateManifest(templateFS, templateName)
	if err != nil {
		return nil, err
	}

	// Create and return TemplateInfo struct
	return &TemplateInfo{
		Name:        templateName,
		Description: manifest.Description,
		Path:        fmt.Sprintf("templates/%s", templateName),
		Manifest:    manifest,
	}, nil
}

// ValidateTemplate checks if a template is valid and complete.
// This function performs comprehensive validation of a template's manifest,
// ensuring that all required fields are present and parameter definitions
// are valid according to the template system's rules.
//
// Parameters:
//   - templateFS: The embedded filesystem containing template files
//   - templateName: The name of the template to validate
//
// Returns:
//   - An error if the template is invalid, nil if valid
func ValidateTemplate(templateFS fs.FS, templateName string) error {
	// Load the template manifest
	manifest, err := LoadTemplateManifest(templateFS, templateName)
	if err != nil {
		return err
	}

	// Validate each parameter in the manifest
	for i, param := range manifest.Parameters {
		// Check for required parameter fields
		if param.Name == "" {
			return fmt.Errorf("parameter %d missing required field: name", i)
		}
		if param.Prompt == "" {
			return fmt.Errorf("parameter %d missing required field: prompt", i)
		}
		if param.Type == "" {
			return fmt.Errorf("parameter %d missing required field: type", i)
		}

		// Validate parameter type against supported types
		switch param.Type {
		case "string", "boolean", "select", "multiselect":
			// Valid types - no action needed
		default:
			return fmt.Errorf("parameter %s has invalid type: %s", param.Name, param.Type)
		}

		// Validate that select/multiselect parameters have options defined
		if (param.Type == "select" || param.Type == "multiselect") && len(param.Options) == 0 {
			return fmt.Errorf("parameter %s of type %s must have options", param.Name, param.Type)
		}
	}

	return nil
}
