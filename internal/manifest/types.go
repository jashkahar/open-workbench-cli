package manifest

// WorkbenchManifest represents the evolved structure of workbench.yaml
type WorkbenchManifest struct {
	APIVersion   string                 `yaml:"apiVersion"`
	Kind         string                 `yaml:"kind"`
	Metadata     ProjectMetadata        `yaml:"metadata"`
	Environments map[string]Environment `yaml:"environments,omitempty"`
	Components   map[string]Component   `yaml:"components,omitempty"`
	Services     map[string]Service     `yaml:"services"`
}

// ProjectMetadata contains project-level information
type ProjectMetadata struct {
	Name string `yaml:"name"`
}

// Environment represents a deployment environment configuration
type Environment struct {
	Provider string            `yaml:"provider"` // aws, gcp, azure, etc.
	Region   string            `yaml:"region,omitempty"`
	Config   map[string]string `yaml:"config,omitempty"`
}

// Component represents a shared project component (like a gateway)
type Component struct {
	Template string   `yaml:"template"`
	Path     string   `yaml:"path"`
	Ports    []string `yaml:"ports,omitempty"`
}

// Service represents a service in the project with its configuration
type Service struct {
	Template    string              `yaml:"template"`
	Path        string              `yaml:"path"`
	Port        int                 `yaml:"port,omitempty"`
	Resources   map[string]Resource `yaml:"resources,omitempty"`
	Environment map[string]string   `yaml:"environment,omitempty"`
}

// Resource represents a service-owned resource (like a database)
type Resource struct {
	Type    string            `yaml:"type"`
	Version string            `yaml:"version,omitempty"`
	Config  map[string]string `yaml:"config,omitempty"`
}
