package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jashkahar/open-workbench-platform/internal/compose"
	"github.com/spf13/cobra"
)

var composeCmd = &cobra.Command{
	Use:   "compose",
	Short: "Generate Docker Compose configuration from workbench.yaml",
	Long: `Generate a complete Docker Compose configuration from your workbench.yaml file.

This command will:
1. Check for required prerequisites (Docker, Docker Compose)
2. Parse your workbench.yaml file
3. Generate docker-compose.yml with proper service networking
4. Create .env and .env.example files with secure defaults
5. Provide clear instructions for starting your application

The generated docker-compose.yml file is clean, human-readable, and follows
Docker Compose best practices. You can modify it directly or regenerate it
by running this command again.`,
	RunE: runCompose,
}

// initComposeCommand registers the compose command with the root command
func initComposeCommand() {
	if rootCmd != nil {
		rootCmd.AddCommand(composeCmd)
	}
}

func runCompose(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Checking prerequisites...")

	// Check prerequisites
	checker := compose.NewPrerequisiteChecker()
	if err := checker.CheckAllPrerequisites(); err != nil {
		fmt.Printf("❌ Prerequisite check failed: %s\n\n", err)
		fmt.Println("📋 Installation instructions:")
		fmt.Println(checker.GetPlatformSpecificInstructions())
		return err
	}

	fmt.Println("✅ Prerequisites satisfied")

	// Find workbench.yaml
	workbenchPath := "workbench.yaml"
	if _, err := os.Stat(workbenchPath); os.IsNotExist(err) {
		return fmt.Errorf("workbench.yaml not found in current directory. Please run this command from your project root")
	}

	fmt.Println("📖 Loading workbench.yaml...")

	// Load and parse workbench.yaml
	project, err := compose.LoadWorkbenchProject(workbenchPath)
	if err != nil {
		return fmt.Errorf("failed to load workbench.yaml: %w", err)
	}

	fmt.Printf("✅ Loaded project: %s\n", project.Metadata.Name)

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
