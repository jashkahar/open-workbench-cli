package compose

// WorkbenchProject represents the evolved workbench.yaml structure
// that supports components, services with resources, and environment variables.
type WorkbenchProject struct {
	APIVersion string               `yaml:"apiVersion"`
	Kind       string               `yaml:"kind"`
	Metadata   ProjectMetadata      `yaml:"metadata"`
	Components map[string]Component `yaml:"components,omitempty"`
	Services   map[string]Service   `yaml:"services"`
}

// ProjectMetadata contains project-level metadata
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

// DockerComposeService represents a service in the generated docker-compose.yml
type DockerComposeService struct {
	Build       *BuildConfig `yaml:"build,omitempty"`
	Image       string       `yaml:"image,omitempty"`
	Ports       []string     `yaml:"ports,omitempty"`
	Environment []string     `yaml:"environment,omitempty"`
	EnvFile     []string     `yaml:"env_file,omitempty"`
	Networks    []string     `yaml:"networks,omitempty"`
	DependsOn   []string     `yaml:"depends_on,omitempty"`
	Volumes     []string     `yaml:"volumes,omitempty"`
}

// BuildConfig represents the build configuration for a service
type BuildConfig struct {
	Context string `yaml:"context"`
}

// DockerComposeConfig represents the complete docker-compose.yml structure
type DockerComposeConfig struct {
	Version  string                          `yaml:"version"`
	Services map[string]DockerComposeService `yaml:"services"`
	Volumes  map[string]interface{}          `yaml:"volumes,omitempty"`
	Networks map[string]interface{}          `yaml:"networks,omitempty"`
}
