package compose

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Generator handles the translation of workbench.yaml to docker-compose.yml
type Generator struct {
	project *WorkbenchProject
}

// NewGenerator creates a new generator instance
func NewGenerator(project *WorkbenchProject) *Generator {
	return &Generator{
		project: project,
	}
}

// Generate creates the docker-compose.yml configuration
func (g *Generator) Generate() (*DockerComposeConfig, error) {
	config := &DockerComposeConfig{
		Version:  "3.8",
		Services: make(map[string]DockerComposeService),
		Volumes:  make(map[string]interface{}),
		Networks: map[string]interface{}{
			"workbench_net": map[string]string{
				"driver": "bridge",
			},
		},
	}

	// Process components first
	for name, component := range g.project.Components {
		service := g.createComponentService(name, component)
		config.Services[name] = service
	}

	// Process services
	for name, service := range g.project.Services {
		dockerService := g.createService(name, service)
		config.Services[name] = dockerService

		// Add service-owned resources
		for resourceName, resource := range service.Resources {
			resourceService := g.createResourceService(name, resourceName, resource)
			config.Services[resourceServiceName(name, resourceName)] = resourceService

			// Add volume for the resource
			volumeName := fmt.Sprintf("%s_%s_data", name, resourceName)
			config.Volumes[volumeName] = nil
		}
	}

	// Resolve dependencies and environment variables
	g.resolveDependencies(config)
	g.resolveEnvironmentVariables(config)

	return config, nil
}

// createComponentService creates a Docker Compose service for a component
func (g *Generator) createComponentService(name string, component Component) DockerComposeService {
	service := DockerComposeService{
		Build: &BuildConfig{
			Context: component.Path,
		},
		Networks: []string{"workbench_net"},
	}

	if len(component.Ports) > 0 {
		service.Ports = component.Ports
	}

	return service
}

// createService creates a Docker Compose service for a regular service
func (g *Generator) createService(name string, service Service) DockerComposeService {
	dockerService := DockerComposeService{
		Build: &BuildConfig{
			Context: service.Path,
		},
		EnvFile:  []string{"./.env"},
		Networks: []string{"workbench_net"},
	}

	// Add port mapping if specified
	if service.Port > 0 {
		dockerService.Ports = []string{fmt.Sprintf("%d:%d", service.Port, service.Port)}
	}

	// Add environment variables
	if len(service.Environment) > 0 {
		for key, value := range service.Environment {
			dockerService.Environment = append(dockerService.Environment, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return dockerService
}

// createResourceService creates a Docker Compose service for a resource (like a database)
func (g *Generator) createResourceService(serviceName, resourceName string, resource Resource) DockerComposeService {
	dockerService := DockerComposeService{
		Image:    fmt.Sprintf("%s:%s", resource.Type, resource.Version),
		EnvFile:  []string{"./.env"},
		Networks: []string{"workbench_net"},
	}

	// Add volume for data persistence
	volumeName := fmt.Sprintf("%s_%s_data", serviceName, resourceName)
	dockerService.Volumes = []string{fmt.Sprintf("%s:/var/lib/%s/data", volumeName, resource.Type)}

	return dockerService
}

// resolveDependencies analyzes environment variables to determine service dependencies
func (g *Generator) resolveDependencies(config *DockerComposeConfig) {
	for serviceName, service := range config.Services {
		dependencies := g.extractDependencies(serviceName, service)
		if len(dependencies) > 0 {
			service.DependsOn = dependencies
			config.Services[serviceName] = service
		}
	}
}

// extractDependencies extracts service dependencies from environment variables
func (g *Generator) extractDependencies(serviceName string, service DockerComposeService) []string {
	var dependencies []string

	for _, envVar := range service.Environment {
		// Look for patterns like ${services.backend.name} or ${components.gateway.name}
		matches := regexp.MustCompile(`\$\{([^}]+)\}`).FindAllString(envVar, -1)

		for _, match := range matches {
			// Extract the variable name without ${}
			varName := strings.Trim(match, "${}")
			parts := strings.Split(varName, ".")

			if len(parts) >= 2 {
				// Check if it's a service reference
				if parts[0] == "services" && len(parts) >= 3 {
					referencedService := parts[1]
					if referencedService != serviceName {
						dependencies = append(dependencies, referencedService)
					}
				}
				// Check if it's a component reference
				if parts[0] == "components" && len(parts) >= 3 {
					referencedComponent := parts[1]
					dependencies = append(dependencies, referencedComponent)
				}
			}
		}
	}

	return dependencies
}

// resolveEnvironmentVariables resolves environment variable references
func (g *Generator) resolveEnvironmentVariables(config *DockerComposeConfig) {
	for serviceName, service := range config.Services {
		resolvedEnv := make([]string, 0, len(service.Environment))

		for _, envVar := range service.Environment {
			resolvedVar := g.resolveEnvironmentVariable(envVar, serviceName)
			resolvedEnv = append(resolvedEnv, resolvedVar)
		}

		service.Environment = resolvedEnv
		config.Services[serviceName] = service
	}
}

// resolveEnvironmentVariable resolves a single environment variable
func (g *Generator) resolveEnvironmentVariable(envVar, serviceName string) string {
	// Replace ${services.service.resources.resource.property} patterns
	re := regexp.MustCompile(`\$\{services\.([^.]+)\.resources\.([^.]+)\.([^}]+)\}`)
	envVar = re.ReplaceAllStringFunc(envVar, func(match string) string {
		parts := strings.Split(strings.Trim(match, "${}"), ".")
		if len(parts) >= 4 {
			service := parts[1]
			resource := parts[2]
			property := parts[3]

			// Generate default values based on the property
			switch property {
			case "user":
				return fmt.Sprintf("%s_%s_user", service, resource)
			case "password":
				return fmt.Sprintf("%s_%s_password", service, resource)
			case "name":
				return fmt.Sprintf("%s_%s", service, resource)
			case "dbname":
				return fmt.Sprintf("%s_%s_db", service, resource)
			}
		}
		return match
	})

	// Replace ${components.component.property} patterns
	re = regexp.MustCompile(`\$\{components\.([^.]+)\.([^}]+)\}`)
	envVar = re.ReplaceAllStringFunc(envVar, func(match string) string {
		parts := strings.Split(strings.Trim(match, "${}"), ".")
		if len(parts) >= 2 {
			component := parts[1]
			property := parts[2]

			switch property {
			case "name":
				return component
			case "port":
				// Extract port from component configuration
				if comp, exists := g.project.Components[component]; exists && len(comp.Ports) > 0 {
					portMapping := comp.Ports[0]
					portParts := strings.Split(portMapping, ":")
					if len(portParts) >= 2 {
						return portParts[0]
					}
				}
			}
		}
		return match
	})

	return envVar
}

// resourceServiceName generates the name for a resource service
func resourceServiceName(serviceName, resourceName string) string {
	return fmt.Sprintf("%s-%s", serviceName, resourceName)
}

// GenerateEnvFile generates the .env file with default credentials
func (g *Generator) GenerateEnvFile() (map[string]string, error) {
	envVars := make(map[string]string)

	// Generate default credentials for each service's resources
	for serviceName, service := range g.project.Services {
		for resourceName, resource := range service.Resources {
			prefix := fmt.Sprintf("%s_%s", serviceName, resourceName)

			// Generate default credentials based on resource type
			switch resource.Type {
			case "postgres":
				envVars[fmt.Sprintf("%s_user", prefix)] = fmt.Sprintf("%s_user", serviceName)
				envVars[fmt.Sprintf("%s_password", prefix)] = "password123"
				envVars[fmt.Sprintf("%s_name", prefix)] = fmt.Sprintf("%s_%s", serviceName, resourceName)
				envVars[fmt.Sprintf("%s_dbname", prefix)] = fmt.Sprintf("%s_%s_db", serviceName, resourceName)
			case "mysql":
				envVars[fmt.Sprintf("%s_user", prefix)] = fmt.Sprintf("%s_user", serviceName)
				envVars[fmt.Sprintf("%s_password", prefix)] = "password123"
				envVars[fmt.Sprintf("%s_name", prefix)] = fmt.Sprintf("%s_%s", serviceName, resourceName)
				envVars[fmt.Sprintf("%s_dbname", prefix)] = fmt.Sprintf("%s_%s_db", serviceName, resourceName)
			case "redis":
				envVars[fmt.Sprintf("%s_password", prefix)] = "password123"
			}
		}
	}

	return envVars, nil
}

// LoadWorkbenchProject loads a workbench.yaml file
func LoadWorkbenchProject(filePath string) (*WorkbenchProject, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workbench.yaml: %w", err)
	}

	var project WorkbenchProject
	if err := yaml.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to parse workbench.yaml: %w", err)
	}

	return &project, nil
}

// SaveDockerCompose saves the docker-compose.yml file
func SaveDockerCompose(config *DockerComposeConfig, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal docker-compose config: %w", err)
	}

	// Add header comment
	header := "# THIS FILE IS AUTO-GENERATED BY 'om compose'.\n# For permanent changes, modify your workbench.yaml and re-run the command.\n\n"
	data = append([]byte(header), data...)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write docker-compose.yml: %w", err)
	}

	return nil
}

// SaveEnvFile saves the .env file
func SaveEnvFile(envVars map[string]string, filePath string) error {
	var lines []string
	for key, value := range envVars {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	data := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	return nil
}

// SaveEnvExampleFile saves the .env.example file
func SaveEnvExampleFile(envVars map[string]string, filePath string) error {
	var lines []string
	for key := range envVars {
		lines = append(lines, fmt.Sprintf("%s=", key))
	}

	data := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
		return fmt.Errorf("failed to write .env.example file: %w", err)
	}

	return nil
}
