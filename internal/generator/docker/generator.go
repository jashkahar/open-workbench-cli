package docker

import (
	"fmt"
	"os"
	"strings"

	"github.com/jashkahar/open-workbench-platform/internal/compose"
	"github.com/jashkahar/open-workbench-platform/internal/manifest"
)

// Generator implements the Generator interface for Docker Compose
type Generator struct{}

// NewGenerator creates a new Docker generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Name returns the unique identifier for this generator
func (g *Generator) Name() string {
	return "docker"
}

// Description returns a human-readable description of this generator
func (g *Generator) Description() string {
	return "Generate Docker Compose configuration for local development"
}

// Validate checks if the manifest is compatible with this generator
func (g *Generator) Validate(manifest *manifest.WorkbenchManifest) error {
	if manifest == nil {
		return fmt.Errorf("manifest cannot be nil")
	}

	if manifest.Metadata.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if len(manifest.Services) == 0 {
		return fmt.Errorf("at least one service is required")
	}

	return nil
}

// Generate creates the Docker Compose configuration for the given manifest
func (g *Generator) Generate(manifest *manifest.WorkbenchManifest) error {
	// Validate the manifest
	if err := g.Validate(manifest); err != nil {
		return fmt.Errorf("manifest validation failed: %w", err)
	}

	// Check prerequisites
	fmt.Println("🔍 Checking prerequisites...")
	checker := compose.NewPrerequisiteChecker()
	if err := checker.CheckAllPrerequisites(); err != nil {
		return fmt.Errorf("prerequisite check failed: %w", err)
	}
	fmt.Println("✅ Prerequisites satisfied")

	// Convert manifest to compose.WorkbenchProject
	project := convertManifestToProject(manifest)

	// Create generator
	generator := compose.NewGenerator(project)

	fmt.Println("🔧 Generating Docker Compose configuration...")

	// Generate docker-compose.yml
	config, err := generator.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate docker-compose configuration: %w", err)
	}

	// Save docker-compose.yml
	if err := compose.SaveDockerCompose(config, "docker-compose.yml"); err != nil {
		return fmt.Errorf("failed to save docker-compose.yml: %w", err)
	}

	fmt.Println("✅ Generated docker-compose.yml")

	// Generate environment files
	fmt.Println("🔐 Generating environment files...")

	envVars, err := generator.GenerateEnvFile()
	if err != nil {
		return fmt.Errorf("failed to generate environment variables: %w", err)
	}

	// Save .env file
	if err := compose.SaveEnvFile(envVars, ".env"); err != nil {
		return fmt.Errorf("failed to save .env file: %w", err)
	}

	// Save .env.example file
	if err := compose.SaveEnvExampleFile(envVars, ".env.example"); err != nil {
		return fmt.Errorf("failed to save .env.example file: %w", err)
	}

	fmt.Println("✅ Generated .env and .env.example files")

	// Update .gitignore to include .env
	if err := updateGitignore(); err != nil {
		fmt.Printf("⚠️  Warning: Could not update .gitignore: %s\n", err)
	}

	// Print success message with instructions
	printComposeSuccessMessage(checker.GetDockerComposeCommand())

	return nil
}

// convertManifestToProject converts manifest.WorkbenchManifest to compose.WorkbenchProject
func convertManifestToProject(manifest *manifest.WorkbenchManifest) *compose.WorkbenchProject {
	project := &compose.WorkbenchProject{
		APIVersion: manifest.APIVersion,
		Kind:       manifest.Kind,
		Metadata: compose.ProjectMetadata{
			Name: manifest.Metadata.Name,
		},
		Components: make(map[string]compose.Component),
		Services:   make(map[string]compose.Service),
	}

	// Convert components
	for name, component := range manifest.Components {
		project.Components[name] = compose.Component{
			Template: component.Template,
			Path:     component.Path,
			Ports:    component.Ports,
		}
	}

	// Convert services
	for name, service := range manifest.Services {
		project.Services[name] = compose.Service{
			Template:    service.Template,
			Path:        service.Path,
			Port:        service.Port,
			Environment: service.Environment,
			Resources:   make(map[string]compose.Resource),
		}

		// Convert resources
		for resourceName, resource := range service.Resources {
			project.Services[name].Resources[resourceName] = compose.Resource{
				Type:    resource.Type,
				Version: resource.Version,
				Config:  resource.Config,
			}
		}
	}

	return project
}

func updateGitignore() error {
	gitignorePath := ".gitignore"

	// Read existing .gitignore
	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Check if .env is already in .gitignore
	contentStr := string(content)
	if contains(contentStr, ".env") {
		return nil // Already exists
	}

	// Append .env to .gitignore
	newContent := contentStr
	if len(newContent) > 0 && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += "\n# Environment variables\n.env\n"

	return os.WriteFile(gitignorePath, []byte(newContent), 0644)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func printComposeSuccessMessage(dockerComposeCmd string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ Successfully built your local environment configuration!")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n📁 Generated files:")
	fmt.Println("  • docker-compose.yml - Main configuration file")
	fmt.Println("  • .env - Environment variables with default credentials")
	fmt.Println("  • .env.example - Template for environment variables")

	fmt.Println("\n🔑 Security notes:")
	fmt.Println("  • Default credentials are in .env file - review and change them")
	fmt.Println("  • .env file is automatically added to .gitignore")
	fmt.Println("  • Use .env.example as a template for production deployments")

	fmt.Println("\n🚀 To start your application, run:")
	fmt.Printf("  %s -f docker-compose.yml up --build\n", dockerComposeCmd)

	fmt.Println("\n📋 Additional commands:")
	fmt.Printf("  %s -f docker-compose.yml down     # Stop all services\n", dockerComposeCmd)
	fmt.Printf("  %s -f docker-compose.yml logs     # View service logs\n", dockerComposeCmd)
	fmt.Printf("  %s -f docker-compose.yml ps       # List running services\n", dockerComposeCmd)

	fmt.Println("\n💡 Tips:")
	fmt.Println("  • The generated docker-compose.yml is human-readable and editable")
	fmt.Println("  • Make changes to workbench.yaml and re-run 'om compose' to regenerate")
	fmt.Println("  • Each service owns its own resources (databases, etc.)")
	fmt.Println("  • Services are automatically networked together")

	fmt.Println("\n🎉 Your local development environment is ready!")
}
