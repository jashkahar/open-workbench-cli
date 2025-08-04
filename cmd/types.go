package cmd

// WorkbenchManifest represents the evolved structure of workbench.yaml
type WorkbenchManifest struct {
	APIVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   ProjectMetadata        `yaml:"metadata"`
	Components map[string]Component   `yaml:"components,omitempty"`
	Services   map[string]Service     `yaml:"services"`
}

// ProjectMetadata contains project-level information
type ProjectMetadata struct {
	Name string `yaml:"name"`
}

// Component represents a shared project component (like a gateway)
type Component struct {
	Template string   `yaml:"template"`
	Path     string   `yaml:"path"`
	Ports    []string `yaml:"ports,omitempty"`
}

// Service represents a service in the project with its configuration
type Service struct {
	Template    string                 `yaml:"template"`
	Path        string                 `yaml:"path"`
	Port        int                    `yaml:"port,omitempty"`
	Resources   map[string]Resource    `yaml:"resources,omitempty"`
	Environment map[string]string      `yaml:"environment,omitempty"`
}

// Resource represents a service-owned resource (like a database)
type Resource struct {
	Type    string            `yaml:"type"`
	Version string            `yaml:"version,omitempty"`
	Config  map[string]string `yaml:"config,omitempty"`
}
