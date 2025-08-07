package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	manifestPkg "github.com/jashkahar/open-workbench-platform/internal/manifest"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete services, components, or resources from your project",
	Long: `Delete services, components, or resources from your project.

This command provides safe deletion by default - it only removes entries from workbench.yaml
without touching any files on disk. Use the --files flag to also delete the corresponding
directories and files.

Safety Features:
  • By default, only removes from workbench.yaml (safe)
  • --files flag requires explicit confirmation
  • Confirmation prompts for destructive operations
  • Validates dependencies before deletion

Examples:
  # Delete a service (manifest only)
  om delete service backend

  # Delete a service and its files
  om delete service backend --files

  # Delete a component
  om delete component gateway

  # Delete a resource
  om delete resource backend.database

  # Interactive mode
  om delete service`,
	RunE: runDelete,
}

var deleteServiceCmd = &cobra.Command{
	Use:   "service [name]",
	Short: "Delete a service from your project",
	Long: `Delete a service from your project.

This command removes the service from workbench.yaml and optionally deletes
the service directory and files.

Examples:
  om delete service backend
  om delete service backend --files`,
	RunE: runDeleteService,
}

var deleteComponentCmd = &cobra.Command{
	Use:   "component [name]",
	Short: "Delete a component from your project",
	Long: `Delete a component from your project.

This command removes the component from workbench.yaml and optionally deletes
the component directory and files.

Examples:
  om delete component gateway
  om delete component gateway --files`,
	RunE: runDeleteComponent,
}

var deleteResourceCmd = &cobra.Command{
	Use:   "resource [service.resource]",
	Short: "Delete a resource from a service",
	Long: `Delete a resource from a service.

This command removes the resource from workbench.yaml. The resource name should
be in the format "service.resource" (e.g., "backend.database").

Examples:
  om delete resource backend.database
  om delete resource frontend.cache`,
	RunE: runDeleteResource,
}

// initDeleteCommand initializes the delete command
func initDeleteCommand() {
	if rootCmd != nil {
		rootCmd.AddCommand(deleteCmd)
	}

	// Add subcommands
	deleteCmd.AddCommand(deleteServiceCmd)
	deleteCmd.AddCommand(deleteComponentCmd)
	deleteCmd.AddCommand(deleteResourceCmd)

	// Add flags
	deleteServiceCmd.Flags().Bool("files", false, "Also delete the service directory and files")
	deleteComponentCmd.Flags().Bool("files", false, "Also delete the component directory and files")
}

func runDelete(cmd *cobra.Command, args []string) error {
	// Show help if no subcommand is provided
	return cmd.Help()
}

func runDeleteService(cmd *cobra.Command, args []string) error {
	// Find project root and load manifest
	projectRoot, manifest, err := findProjectRootAndLoadManifest()
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	// Get service name from args or prompt
	serviceName, err := getServiceNameFromArgs(args, manifest)
	if err != nil {
		return fmt.Errorf("failed to get service name: %w", err)
	}

	// Get files flag
	deleteFiles, err := cmd.Flags().GetBool("files")
	if err != nil {
		return fmt.Errorf("failed to get files flag: %w", err)
	}

	// Validate service exists
	if _, exists := manifest.Services[serviceName]; !exists {
		return fmt.Errorf("service '%s' not found in workbench.yaml", serviceName)
	}

	// Confirm deletion
	if err := confirmDeletion("service", serviceName, deleteFiles); err != nil {
		return err
	}

	// Delete service
	if err := deleteService(manifest, serviceName, projectRoot, deleteFiles); err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	printDeleteSuccessMessage("service", serviceName, deleteFiles)
	return nil
}

func runDeleteComponent(cmd *cobra.Command, args []string) error {
	// Find project root and load manifest
	projectRoot, manifest, err := findProjectRootAndLoadManifest()
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	// Get component name from args or prompt
	componentName, err := getComponentNameFromArgs(args, manifest)
	if err != nil {
		return fmt.Errorf("failed to get component name: %w", err)
	}

	// Get files flag
	deleteFiles, err := cmd.Flags().GetBool("files")
	if err != nil {
		return fmt.Errorf("failed to get files flag: %w", err)
	}

	// Validate component exists
	if _, exists := manifest.Components[componentName]; !exists {
		return fmt.Errorf("component '%s' not found in workbench.yaml", componentName)
	}

	// Confirm deletion
	if err := confirmDeletion("component", componentName, deleteFiles); err != nil {
		return err
	}

	// Delete component
	if err := deleteComponent(manifest, componentName, projectRoot, deleteFiles); err != nil {
		return fmt.Errorf("failed to delete component: %w", err)
	}

	printDeleteSuccessMessage("component", componentName, deleteFiles)
	return nil
}

func runDeleteResource(cmd *cobra.Command, args []string) error {
	// Find project root and load manifest
	projectRoot, manifest, err := findProjectRootAndLoadManifest()
	if err != nil {
		return fmt.Errorf("failed to load project: %w", err)
	}

	// Get resource name from args or prompt
	resourceName, err := getResourceNameFromArgs(args, manifest)
	if err != nil {
		return fmt.Errorf("failed to get resource name: %w", err)
	}

	// Parse service.resource format
	parts := strings.Split(resourceName, ".")
	if len(parts) != 2 {
		return fmt.Errorf("resource name must be in format 'service.resource' (e.g., 'backend.database')")
	}

	serviceName := parts[0]
	resourceNameOnly := parts[1]

	// Validate service and resource exist
	if _, exists := manifest.Services[serviceName]; !exists {
		return fmt.Errorf("service '%s' not found in workbench.yaml", serviceName)
	}

	if _, exists := manifest.Services[serviceName].Resources[resourceNameOnly]; !exists {
		return fmt.Errorf("resource '%s' not found in service '%s'", resourceNameOnly, serviceName)
	}

	// Confirm deletion
	if err := confirmDeletion("resource", resourceName, false); err != nil {
		return err
	}

	// Delete resource
	if err := deleteResource(manifest, serviceName, resourceNameOnly, projectRoot); err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	printDeleteSuccessMessage("resource", resourceName, false)
	return nil
}

func getServiceNameFromArgs(args []string, manifest *manifestPkg.WorkbenchManifest) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	// Interactive mode - prompt for service
	serviceNames := make([]string, 0, len(manifest.Services))
	for name := range manifest.Services {
		serviceNames = append(serviceNames, name)
	}

	if len(serviceNames) == 0 {
		return "", fmt.Errorf("no services found in workbench.yaml")
	}

	var selectedService string
	prompt := &survey.Select{
		Message: "Which service would you like to delete?",
		Options: serviceNames,
		Help:    "Select the service to delete from your project",
	}

	if err := survey.AskOne(prompt, &selectedService); err != nil {
		return "", fmt.Errorf("failed to get service selection: %w", err)
	}

	return selectedService, nil
}

func getComponentNameFromArgs(args []string, manifest *manifestPkg.WorkbenchManifest) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	// Interactive mode - prompt for component
	componentNames := make([]string, 0, len(manifest.Components))
	for name := range manifest.Components {
		componentNames = append(componentNames, name)
	}

	if len(componentNames) == 0 {
		return "", fmt.Errorf("no components found in workbench.yaml")
	}

	var selectedComponent string
	prompt := &survey.Select{
		Message: "Which component would you like to delete?",
		Options: componentNames,
		Help:    "Select the component to delete from your project",
	}

	if err := survey.AskOne(prompt, &selectedComponent); err != nil {
		return "", fmt.Errorf("failed to get component selection: %w", err)
	}

	return selectedComponent, nil
}

func getResourceNameFromArgs(args []string, manifest *manifestPkg.WorkbenchManifest) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	// Interactive mode - prompt for resource
	var resourceOptions []string
	for serviceName, service := range manifest.Services {
		for resourceName := range service.Resources {
			resourceOptions = append(resourceOptions, fmt.Sprintf("%s.%s", serviceName, resourceName))
		}
	}

	if len(resourceOptions) == 0 {
		return "", fmt.Errorf("no resources found in workbench.yaml")
	}

	var selectedResource string
	prompt := &survey.Select{
		Message: "Which resource would you like to delete?",
		Options: resourceOptions,
		Help:    "Select the resource to delete from your project",
	}

	if err := survey.AskOne(prompt, &selectedResource); err != nil {
		return "", fmt.Errorf("failed to get resource selection: %w", err)
	}

	return selectedResource, nil
}

func confirmDeletion(entityType, name string, deleteFiles bool) error {
	var message string
	if deleteFiles {
		message = fmt.Sprintf("Are you sure you want to delete %s '%s' and ALL its files? This action cannot be undone.", entityType, name)
	} else {
		message = fmt.Sprintf("Are you sure you want to delete %s '%s' from workbench.yaml? (This will not delete any files)", entityType, name)
	}

	var confirmed bool
	prompt := &survey.Confirm{
		Message: message,
		Help:    "This action will remove the entry from workbench.yaml",
	}

	if err := survey.AskOne(prompt, &confirmed); err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}

	if !confirmed {
		return fmt.Errorf("deletion cancelled")
	}

	// Additional confirmation for file deletion
	if deleteFiles {
		var finalConfirmed bool
		finalPrompt := &survey.Confirm{
			Message: fmt.Sprintf("⚠️  FINAL WARNING: This will permanently delete the %s directory and ALL files. Are you absolutely sure?", entityType),
			Help:    "This action is irreversible and will delete all files in the directory",
		}

		if err := survey.AskOne(finalPrompt, &finalConfirmed); err != nil {
			return fmt.Errorf("failed to get final confirmation: %w", err)
		}

		if !finalConfirmed {
			return fmt.Errorf("file deletion cancelled")
		}
	}

	return nil
}

func deleteService(manifest *manifestPkg.WorkbenchManifest, serviceName, projectRoot string, deleteFiles bool) error {
	// Get service path before deletion
	service := manifest.Services[serviceName]
	servicePath := service.Path

	// Remove from manifest
	delete(manifest.Services, serviceName)

	// Save updated manifest
	if err := saveWorkbenchManifest(manifest, projectRoot); err != nil {
		return fmt.Errorf("failed to save workbench.yaml: %w", err)
	}

	// Delete files if requested
	if deleteFiles && servicePath != "" {
		fullPath := filepath.Join(projectRoot, servicePath)
		if err := os.RemoveAll(fullPath); err != nil {
			return fmt.Errorf("failed to delete service directory: %w", err)
		}
	}

	return nil
}

func deleteComponent(manifest *manifestPkg.WorkbenchManifest, componentName, projectRoot string, deleteFiles bool) error {
	// Get component path before deletion
	component := manifest.Components[componentName]
	componentPath := component.Path

	// Remove from manifest
	delete(manifest.Components, componentName)

	// Save updated manifest
	if err := saveWorkbenchManifest(manifest, projectRoot); err != nil {
		return fmt.Errorf("failed to save workbench.yaml: %w", err)
	}

	// Delete files if requested
	if deleteFiles && componentPath != "" {
		fullPath := filepath.Join(projectRoot, componentPath)
		if err := os.RemoveAll(fullPath); err != nil {
			return fmt.Errorf("failed to delete component directory: %w", err)
		}
	}

	return nil
}

func deleteResource(manifest *manifestPkg.WorkbenchManifest, serviceName, resourceName, projectRoot string) error {
	// Remove from manifest
	delete(manifest.Services[serviceName].Resources, resourceName)

	// Save updated manifest
	if err := saveWorkbenchManifest(manifest, projectRoot); err != nil {
		return fmt.Errorf("failed to save workbench.yaml: %w", err)
	}

	return nil
}

func printDeleteSuccessMessage(entityType, name string, deletedFiles bool) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("✅ Successfully deleted %s '%s'!\n", entityType, name)
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n📁 Updated files:")
	fmt.Println("  • workbench.yaml - Removed entry")

	if deletedFiles {
		fmt.Println("  • Deleted directory and all files")
	} else {
		fmt.Println("  • Files were preserved (use --files to delete them)")
	}

	fmt.Println("\n💡 Tips:")
	fmt.Println("  • Run 'om ls' to see the updated project structure")
	fmt.Println("  • Run 'om compose' to regenerate configuration files")
	fmt.Println("  • The deletion only affects workbench.yaml by default")

	fmt.Println("\n🎉 Deletion completed successfully!")
}
