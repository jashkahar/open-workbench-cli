package cmd

// WorkbenchManifest represents the structure of workbench.yaml
type WorkbenchManifest struct {
	APIVersion string             `yaml:"apiVersion"`
	Kind       string             `yaml:"kind"`
	Metadata   ProjectMetadata    `yaml:"metadata"`
	Services   map[string]Service `yaml:"services"`
}

// ProjectMetadata contains project-level information
type ProjectMetadata struct {
	Name string `yaml:"name"`
}

// Service represents a service in the project
type Service struct {
	Template string `yaml:"template"`
	Path     string `yaml:"path"`
}
