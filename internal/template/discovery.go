package template

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
)

// TemplateManifest represents the structure of a template.json file
type TemplateManifest struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Parameters   []Parameter   `json:"parameters"`
	PostScaffold *PostScaffold `json:"postScaffold,omitempty"`
}

// Parameter represents a single parameter that the user needs to provide
type Parameter struct {
	Name       string      `json:"name"`
	Prompt     string      `json:"prompt"`
	HelpText   string      `json:"helpText,omitempty"`
	Group      string      `json:"group,omitempty"`
	Type       string      `json:"type"`
	Required   bool        `json:"required,omitempty"`
	Default    any         `json:"default,omitempty"`
	Options    []string    `json:"options,omitempty"`
	Condition  string      `json:"condition,omitempty"`
	Validation *Validation `json:"validation,omitempty"`
}

// Validation represents validation rules for string parameters
type Validation struct {
	Regex        string `json:"regex"`
	ErrorMessage string `json:"errorMessage"`
}

// PostScaffold represents actions to perform after scaffolding
type PostScaffold struct {
	FilesToDelete []FileAction    `json:"filesToDelete,omitempty"`
	Commands      []CommandAction `json:"commands,omitempty"`
}

// FileAction represents a file or directory to delete based on a condition
type FileAction struct {
	Path      string `json:"path"`
	Condition string `json:"condition"`
}

// CommandAction represents a command to execute after scaffolding
type CommandAction struct {
	Command     string `json:"command"`
	Description string `json:"description"`
	Condition   string `json:"condition,omitempty"`
}

// TemplateInfo represents metadata about a discovered template
type TemplateInfo struct {
	Name        string
	Description string
	Path        string
	Manifest    *TemplateManifest
}

// LoadTemplateManifest loads and parses a template.json file from the embedded filesystem
func LoadTemplateManifest(templateFS fs.FS, templateName string) (*TemplateManifest, error) {
	manifestPath := filepath.Join("templates", templateName, "template.json")

	manifestBytes, err := fs.ReadFile(templateFS, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template manifest: %w", err)
	}

	var manifest TemplateManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse template manifest: %w", err)
	}

	// Validate required fields
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

// DiscoverTemplates finds all available templates in the embedded filesystem
func DiscoverTemplates(templateFS fs.FS) ([]TemplateInfo, error) {
	var templates []TemplateInfo

	entries, err := fs.ReadDir(templateFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			templateName := entry.Name()

			// Try to load the manifest for this template
			manifest, err := LoadTemplateManifest(templateFS, templateName)
			if err != nil {
				// Skip templates with invalid manifests for now
				// In the future, we might want to log this or handle it differently
				continue
			}

			templateInfo := TemplateInfo{
				Name:        templateName,
				Description: manifest.Description,
				Path:        filepath.Join("templates", templateName),
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

// GetTemplateInfo returns information about a specific template
func GetTemplateInfo(templateFS fs.FS, templateName string) (*TemplateInfo, error) {
	manifest, err := LoadTemplateManifest(templateFS, templateName)
	if err != nil {
		return nil, err
	}

	return &TemplateInfo{
		Name:        templateName,
		Description: manifest.Description,
		Path:        filepath.Join("templates", templateName),
		Manifest:    manifest,
	}, nil
}

// ValidateTemplate checks if a template is valid and complete
func ValidateTemplate(templateFS fs.FS, templateName string) error {
	manifest, err := LoadTemplateManifest(templateFS, templateName)
	if err != nil {
		return err
	}

	// Validate parameters
	for i, param := range manifest.Parameters {
		if param.Name == "" {
			return fmt.Errorf("parameter %d missing required field: name", i)
		}
		if param.Prompt == "" {
			return fmt.Errorf("parameter %d missing required field: prompt", i)
		}
		if param.Type == "" {
			return fmt.Errorf("parameter %d missing required field: type", i)
		}

		// Validate parameter type
		switch param.Type {
		case "string", "boolean", "select", "multiselect":
			// Valid types
		default:
			return fmt.Errorf("parameter %s has invalid type: %s", param.Name, param.Type)
		}

		// Validate select/multiselect parameters have options
		if (param.Type == "select" || param.Type == "multiselect") && len(param.Options) == 0 {
			return fmt.Errorf("parameter %s of type %s must have options", param.Name, param.Type)
		}
	}

	return nil
}
