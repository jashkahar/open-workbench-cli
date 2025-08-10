package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	manifestPkg "github.com/jashkahar/open-workbench-platform/internal/manifest"
	"github.com/jashkahar/open-workbench-platform/internal/resources"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var addResourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Add a new resource to a service",
	Long: `Add a new resource (database, cache, etc.) to a service in your project.

This command allows you to declaratively add infrastructure dependencies to your services.
Resources are automatically configured for both Docker Compose and Terraform targets.

Interactive Mode (no parameters):
  om add resource

Direct Mode (with parameters):
  om add resource --service backend --type postgres-db --name user_database

Examples:
  # Interactive mode - prompts for all details
  om add resource

  # Direct mode with all parameters
  om add resource --service backend --type postgres-db --name user_database

  # Direct mode with minimal parameters (others will be prompted)
  om add resource --service frontend --type redis-cache

Available resource types:
  • postgres-db - PostgreSQL Database
  • mysql-db - MySQL Database
  • mongodb - MongoDB Database
  • redis-cache - Redis Cache
  • memcached - Memcached Cache
  • rabbitmq - RabbitMQ Message Queue`,
	RunE: runAddResource,
}

// initAddResourceCommand initializes the add resource command
func initAddResourceCommand() {
	addCmd.AddCommand(addResourceCmd)

	// Add flags for the resource command (optional for interactive mode)
	addResourceCmd.Flags().String("service", "", "Service name (optional - will prompt if not provided)")
	addResourceCmd.Flags().String("type", "", "Resource type (optional - will prompt if not provided)")
	addResourceCmd.Flags().String("name", "", "Resource name (optional - will prompt if not provided)")
}

func runAddResource(cmd *cobra.Command, args []string) error {
	// Find project root and load manifest
	projectRoot, manifest, err := findProjectRootAndLoadManifest()
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	// Get parameters from flags or prompt user
	serviceName, resourceType, resourceName, err := getResourceParameters(cmd, manifest)
	if err != nil {
		return fmt.Errorf("failed to get resource parameters: %w", err)
	}

	// Validate the resource configuration
	if err := validateResourceConfiguration(manifest, serviceName, resourceType, resourceName); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get resource blueprint
	resourceRegistry := resources.NewRegistry()
	blueprint, err := resourceRegistry.Get(resourceType)
	if err != nil {
		return fmt.Errorf("failed to get resource blueprint: %w", err)
	}

	// Collect resource parameters
	resourceConfig, err := collectResourceParameters(blueprint)
	if err != nil {
		return fmt.Errorf("failed to collect resource parameters: %w", err)
	}

	// Add resource to manifest
	if err := addResourceToManifest(manifest, serviceName, resourceName, resourceType, resourceConfig); err != nil {
		return fmt.Errorf("failed to add resource to manifest: %w", err)
	}

	// Save updated manifest
	if err := saveWorkbenchManifest(manifest, projectRoot); err != nil {
		return fmt.Errorf("failed to save workbench.yaml: %w", err)
	}

	// Print success message with next steps and actionable guidance
	printAddResourceSuccessMessage(serviceName, resourceName, resourceType, blueprint, resourceConfig)

	return nil
}

func getResourceParameters(cmd *cobra.Command, manifest *manifestPkg.WorkbenchManifest) (string, string, string, error) {
	// Get service name
	serviceName, err := cmd.Flags().GetString("service")
	if err != nil {
		return "", "", "", err
	}

	if serviceName == "" {
		// Interactive mode - prompt for service
		serviceNames := make([]string, 0, len(manifest.Services))
		for name := range manifest.Services {
			serviceNames = append(serviceNames, name)
		}

		if len(serviceNames) == 0 {
			return "", "", "", fmt.Errorf("no services found in workbench.yaml")
		}

		var selectedService string
		prompt := &survey.Select{
			Message: "Which service should this resource belong to?",
			Options: serviceNames,
			Help:    "Select the service that will use this resource",
		}

		if err := survey.AskOne(prompt, &selectedService); err != nil {
			return "", "", "", fmt.Errorf("failed to get service selection: %w", err)
		}

		serviceName = selectedService
	}

	// Get resource type
	resourceType, err := cmd.Flags().GetString("type")
	if err != nil {
		return "", "", "", err
	}

	if resourceType == "" {
		// Interactive mode - prompt for resource type
		resourceRegistry := resources.NewRegistry()

		// Group by category
		categories := resourceRegistry.Categories()
		var options []string
		for _, category := range categories {
			categoryBlueprints := resourceRegistry.ListByCategory(category)
			options = append(options, fmt.Sprintf("--- %s ---", strings.Title(category)))
			for _, blueprint := range categoryBlueprints {
				options = append(options, fmt.Sprintf("%s - %s", blueprint.Name, blueprint.Description))
			}
		}

		var selectedOption string
		prompt := &survey.Select{
			Message: "Which type of resource would you like to add?",
			Options: options,
			Help:    "Select the type of resource to add to your service",
		}

		if err := survey.AskOne(prompt, &selectedOption); err != nil {
			return "", "", "", fmt.Errorf("failed to get resource type selection: %w", err)
		}

		// Extract resource type from selection
		parts := strings.Split(selectedOption, " - ")
		if len(parts) >= 1 {
			resourceType = parts[0]
		}
	}

	// Get resource name
	resourceName, err := cmd.Flags().GetString("name")
	if err != nil {
		return "", "", "", err
	}

	if resourceName == "" {
		// Interactive mode - prompt for resource name
		var name string
		prompt := &survey.Input{
			Message: "What should this resource be named?",
			Help:    "Enter a descriptive name for this resource (e.g., user_database, cache_store)",
		}

		if err := survey.AskOne(prompt, &name); err != nil {
			return "", "", "", fmt.Errorf("failed to get resource name: %w", err)
		}

		resourceName = name
	}

	return serviceName, resourceType, resourceName, nil
}

func validateResourceConfiguration(manifest *manifestPkg.WorkbenchManifest, serviceName, resourceType, resourceName string) error {
	// Validate service exists
	if _, exists := manifest.Services[serviceName]; !exists {
		return fmt.Errorf("service '%s' not found in workbench.yaml", serviceName)
	}

	// Validate resource type
	resourceRegistry := resources.NewRegistry()
	if _, err := resourceRegistry.Get(resourceType); err != nil {
		return fmt.Errorf("invalid resource type '%s': %w", resourceType, err)
	}

	// Validate resource name format
	if resourceName == "" {
		return fmt.Errorf("resource name cannot be empty")
	}

	// Check for duplicate resource name in the service
	service := manifest.Services[serviceName]
	if _, exists := service.Resources[resourceName]; exists {
		return fmt.Errorf("resource '%s' already exists in service '%s'", resourceName, serviceName)
	}

	return nil
}

func collectResourceParameters(blueprint resources.ResourceBlueprint) (map[string]string, error) {
	config := make(map[string]string)

	// Set default values
	for _, param := range blueprint.Parameters {
		if param.Default != nil {
			config[param.Name] = fmt.Sprintf("%v", param.Default)
		}
	}

	// Collect required parameters
	for _, param := range blueprint.Parameters {
		if param.Required {
			var value string

			switch param.Type {
			case "select":
				var selected string
				prompt := &survey.Select{
					Message: fmt.Sprintf("%s:", param.Description),
					Options: param.Options,
					Help:    fmt.Sprintf("Select %s for %s", param.Description, blueprint.Name),
				}

				if err := survey.AskOne(prompt, &selected); err != nil {
					return nil, fmt.Errorf("failed to get %s: %w", param.Name, err)
				}
				value = selected

			case "string":
				var input string
				prompt := &survey.Input{
					Message: fmt.Sprintf("%s:", param.Description),
					Help:    fmt.Sprintf("Enter %s for %s", param.Description, blueprint.Name),
				}

				if err := survey.AskOne(prompt, &input); err != nil {
					return nil, fmt.Errorf("failed to get %s: %w", param.Name, err)
				}
				value = input

			case "number":
				var input int
				prompt := &survey.Input{
					Message: fmt.Sprintf("%s:", param.Description),
					Help:    fmt.Sprintf("Enter %s for %s", param.Description, blueprint.Name),
				}

				if err := survey.AskOne(prompt, &input); err != nil {
					return nil, fmt.Errorf("failed to get %s: %w", param.Name, err)
				}
				value = fmt.Sprintf("%d", input)

			default:
				return nil, fmt.Errorf("unsupported parameter type: %s", param.Type)
			}

			config[param.Name] = value
		}
	}

	return config, nil
}

func addResourceToManifest(manifest *manifestPkg.WorkbenchManifest, serviceName, resourceName, resourceType string, config map[string]string) error {
	// Get the service and create a copy
	service := manifest.Services[serviceName]

	// Initialize resources map if it doesn't exist
	if service.Resources == nil {
		service.Resources = make(map[string]manifestPkg.Resource)
	}

	// Add the resource
	service.Resources[resourceName] = manifestPkg.Resource{
		Type:   resourceType,
		Config: config,
	}

	// Update the service in the manifest
	manifest.Services[serviceName] = service

	return nil
}

func saveWorkbenchManifest(manifest *manifestPkg.WorkbenchManifest, projectRoot string) error {
	// Marshal manifest to YAML
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write to workbench.yaml
	workbenchPath := projectRoot + "/workbench.yaml"
	if err := os.WriteFile(workbenchPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write workbench.yaml: %w", err)
	}

	return nil
}

func printAddResourceSuccessMessage(serviceName, resourceName, resourceType string, blueprint resources.ResourceBlueprint, cfg map[string]string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ Successfully added resource to your project!")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("\n📦 Resource Details:")
	fmt.Printf("\n  • Service: %s", serviceName)
	fmt.Printf("\n  • Resource: %s (%s)", resourceName, resourceType)
	fmt.Printf("\n  • Type: %s", blueprint.Description)
	fmt.Printf("\n  • Category: %s", strings.Title(blueprint.Category))

	fmt.Println("\n\n📁 Updated files:")
	fmt.Println("  • workbench.yaml - Added resource configuration")

	fmt.Println("\n🚀 Next steps:")
	fmt.Println("  1) Run: om compose --target docker")
	fmt.Println("  2) Start: docker compose up --build")
	fmt.Println("  3) Add optional init/config files if needed (see below)")

	// Actionable guidance (concise)
	fmt.Println("\n🧭 Guidance:")
	if v := cfg["port"]; v != "" {
		fmt.Printf("  - Port: %s (host port maps to container)\n", v)
	}
	if v := cfg["databaseName"]; v != "" {
		fmt.Printf("  - Database: %s\n", v)
	}
	if u := cfg["username"]; u != "" {
		fmt.Printf("  - Username: %s\n", u)
	}
	switch strings.ToLower(resourceType) {
	case "postgres", "postgres-db":
		fmt.Println("  - Example connection: postgres://<user>:<password>@localhost:<port>/<db>?sslmode=disable")
		fmt.Println("  - Mount init SQL: ./<service>/init:/docker-entrypoint-initdb.d (optional)")
	case "mysql", "mysql-db":
		fmt.Println("  - Example connection: mysql://<user>:<password>@localhost:<port>/<db>")
		fmt.Println("  - Mount init SQL: ./<service>/init:/docker-entrypoint-initdb.d (optional)")
	case "mongodb", "mongo":
		fmt.Println("  - Example connection: mongodb://<user>:<password>@localhost:<port>/<db>")
	case "redis", "redis-cache":
		fmt.Println("  - Example connection: redis://:<password>@localhost:<port>")
	}
	fmt.Println("  - Volume name: <service>_<resource>_data (defined in docker-compose.yml)")

	fmt.Println("\n💡 Tips:")
	fmt.Println("  • Resources are automatically networked with their services")
	fmt.Println("  • Environment variables will be generated for service-to-resource communication")
	fmt.Println("  • You can add multiple resources to the same service")
	fmt.Println("  • Use 'om ls' to view your project structure")

	fmt.Println("\n🎉 Your resource has been added successfully!")
}
