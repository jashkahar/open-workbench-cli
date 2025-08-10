package cmd

import (
	"fmt"
	"strings"

	manifestPkg "github.com/jashkahar/open-workbench-platform/internal/manifest"
	"github.com/jashkahar/open-workbench-platform/internal/resources"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List project structure and components",
	Long: `Display a high-level, human-readable overview of your project's architecture.

This command reads your workbench.yaml file and displays a formatted tree structure
showing all services, components, and resources in your project.

Examples:
  # List project structure
  om ls

  # List with detailed information
  om ls --detailed

The output includes:
  • Project name and metadata
  • Services with their templates and resources
  • Components with their templates
  • Resource types and configurations
  • Environment configurations (if any)`,
	RunE: runLs,
}

// initLsCommand initializes the ls command
func initLsCommand() {
	if rootCmd != nil {
		rootCmd.AddCommand(lsCmd)
	}

	// Add detailed flag
	lsCmd.Flags().Bool("detailed", false, "Show detailed information including resource configurations")
}

func runLs(cmd *cobra.Command, args []string) error {
	// Find project root and load manifest
	_, manifest, err := findProjectRootAndLoadManifest()
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	// Get detailed flag
	detailed, err := cmd.Flags().GetBool("detailed")
	if err != nil {
		return fmt.Errorf("failed to get detailed flag: %w", err)
	}

	// Print project header
	printProjectHeader(manifest)

	// Print environments if any
	if len(manifest.Environments) > 0 {
		printEnvironments(manifest.Environments)
	}

	// Print components
	if len(manifest.Components) > 0 {
		printComponents(manifest.Components, detailed)
	}

	// Print services
	if len(manifest.Services) > 0 {
		printServices(manifest.Services, detailed)
	}

	// Print summary
	printSummary(manifest)

	return nil
}

func printProjectHeader(manifest *manifestPkg.WorkbenchManifest) {
	fmt.Println("📁 Project Structure")
	fmt.Println("===================")
	fmt.Printf("Project: %s\n", manifest.Metadata.Name)
	fmt.Printf("API Version: %s\n", manifest.APIVersion)
	fmt.Printf("Kind: %s\n", manifest.Kind)
	fmt.Println()
}

func printEnvironments(environments map[string]manifestPkg.Environment) {
	fmt.Println("🌍 Environments")
	fmt.Println("---------------")
	for name, env := range environments {
		fmt.Printf("  %s:\n", name)
		fmt.Printf("    Provider: %s\n", env.Provider)
		if env.Region != "" {
			fmt.Printf("    Region: %s\n", env.Region)
		}
		if len(env.Config) > 0 {
			fmt.Printf("    Config: %v\n", env.Config)
		}
	}
	fmt.Println()
}

func printComponents(components map[string]manifestPkg.Component, detailed bool) {
	fmt.Println("📦 Components")
	fmt.Println("--------------")
	for name, component := range components {
		fmt.Printf("  📦 %s (%s)\n", name, component.Template)
		if detailed {
			fmt.Printf("    Path: %s\n", component.Path)
			if len(component.Ports) > 0 {
				fmt.Printf("    Ports: %s\n", strings.Join(component.Ports, ", "))
			}
		}
	}
	fmt.Println()
}

func printServices(services map[string]manifestPkg.Service, detailed bool) {
	fmt.Println("🚀 Services")
	fmt.Println("------------")
	for name, service := range services {
		fmt.Printf("  💻 %s (%s)\n", name, service.Template)
		if detailed {
			fmt.Printf("    Path: %s\n", service.Path)
			if service.Port != 0 {
				fmt.Printf("    Port: %d\n", service.Port)
			}
			if len(service.Environment) > 0 {
				fmt.Printf("    Environment Variables: %d\n", len(service.Environment))
			}
		}

		// Print resources for this service
		if len(service.Resources) > 0 {
			printServiceResources(name, service.Resources, detailed)
		}
	}
	fmt.Println()
}

func printServiceResources(serviceName string, serviceResources map[string]manifestPkg.Resource, detailed bool) {
	resourceRegistry := resources.NewRegistry()

	for resourceName, resource := range serviceResources {
		// Get resource blueprint for description
		blueprint, err := resourceRegistry.Get(resource.Type)
		description := resource.Type
		if err == nil {
			description = blueprint.Description
		}

		// Choose appropriate emoji based on resource category
		emoji := "🔧"
		if err == nil {
			switch blueprint.Category {
			case "database":
				emoji = "🐘"
			case "cache":
				emoji = "⚡"
			case "storage":
				emoji = "📦"
			case "message-queue":
				emoji = "📨"
			}
		}

		fmt.Printf("    %s %s (%s)\n", emoji, resourceName, description)

		if detailed {
			if resource.Version != "" {
				fmt.Printf("      Version: %s\n", resource.Version)
			}
			if len(resource.Config) > 0 {
				fmt.Printf("      Config: %v\n", resource.Config)
			}
		}
	}
}

func printSummary(manifest *manifestPkg.WorkbenchManifest) {
	fmt.Println("📊 Summary")
	fmt.Println("----------")

	// Count components
	componentCount := len(manifest.Components)
	if componentCount > 0 {
		fmt.Printf("Components: %d\n", componentCount)
	}

	// Count services
	serviceCount := len(manifest.Services)
	fmt.Printf("Services: %d\n", serviceCount)

	// Count total resources
	totalResources := 0
	for _, service := range manifest.Services {
		totalResources += len(service.Resources)
	}
	if totalResources > 0 {
		fmt.Printf("Resources: %d\n", totalResources)
	}

	// Count environments
	environmentCount := len(manifest.Environments)
	if environmentCount > 0 {
		fmt.Printf("Environments: %d\n", environmentCount)
	}

	fmt.Println()
	fmt.Println("💡 Use 'om ls --detailed' for more information")
	fmt.Println("💡 Use 'om compose --target docker' to generate local configuration")
}
