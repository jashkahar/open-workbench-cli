package resources

// ResourceBlueprint defines the configuration for a resource type
type ResourceBlueprint struct {
	// Basic information
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"` // database, cache, storage, etc.

	// Docker Compose configuration
	DockerComposeSnippet string `json:"dockerComposeSnippet"`

	// Terraform configuration
	TerraformModule string `json:"terraformModule"`

	// Resource-specific parameters
	Parameters []ResourceParameter `json:"parameters,omitempty"`

	// Dependencies
	DependsOn []string `json:"dependsOn,omitempty"`
}

// ResourceParameter defines a parameter for a resource
type ResourceParameter struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"` // string, number, boolean, select
	Required    bool     `json:"required"`
	Default     any      `json:"default,omitempty"`
	Options     []string `json:"options,omitempty"` // for select type
}

// ResourceInstance represents an actual resource instance in a service
type ResourceInstance struct {
	Type    string            `yaml:"type"`
	Version string            `yaml:"version,omitempty"`
	Config  map[string]string `yaml:"config,omitempty"`
}

// ResourceValidation represents validation rules for a resource
type ResourceValidation struct {
	// Validation rules can be added here as needed
	MinVersion string `json:"minVersion,omitempty"`
	MaxVersion string `json:"maxVersion,omitempty"`
}
