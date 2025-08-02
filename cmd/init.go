package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/jashkahar/open-workbench-platform/internal/templating"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Open Workbench project",
	Long: `Initialize a new Open Workbench project in the current directory.

This command will:
  1. Check that the current directory is safe to initialize
  2. Prompt for project details
  3. Create the project structure
  4. Generate a workbench.yaml manifest file

Example:
  om init`,
	RunE: runInit,
}

func init() {
	// Command registration will be done in Execute function
}

// runInit executes the init command logic
func runInit(cmd *cobra.Command, args []string) error {
	// Step 1: Safety check - verify the current directory is empty or contains only hidden files
	if err := checkDirectorySafety(); err != nil {
		return err
	}

	// Step 2: Prompt for project name
	projectName, err := promptForProjectName()
	if err != nil {
		return err
	}

	// Step 3: Prompt for first service details
	serviceName, templateName, err := promptForFirstService()
	if err != nil {
		return err
	}

	// Step 4: Create directories
	if err := createProjectDirectories(projectName, serviceName); err != nil {
		return err
	}

	// Step 5: Run the scaffolder
	servicePath := filepath.Join(projectName, serviceName)
	if err := scaffoldService(templateName, servicePath); err != nil {
		return err
	}

	// Step 6: Create and write workbench.yaml
	if err := createWorkbenchManifest(projectName, serviceName, templateName); err != nil {
		return err
	}

	// Step 7: Print success message
	printSuccessMessage(projectName, serviceName)

	return nil
}

// checkDirectorySafety verifies that the current directory is safe to initialize
func checkDirectorySafety() error {
	// Validate current directory safety
	if err := ValidateDirectorySafety("."); err != nil {
		return fmt.Errorf("directory safety check failed: %w", err)
	}

	entries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read current directory: %w", err)
	}

	// Check if directory is empty or contains only hidden files
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			return fmt.Errorf("directory is not empty. Please run 'om init' in an empty directory or a directory containing only hidden files (like .git)")
		}
	}

	return nil
}

// promptForProjectName prompts the user for a project name
func promptForProjectName() (string, error) {
	var projectName string
	prompt := &survey.Input{
		Message: "What is your project name?",
		Help:    "This will be used as the directory name and in the workbench.yaml manifest",
	}
	err := survey.AskOne(prompt, &projectName, survey.WithValidator(survey.Required))
	if err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			fmt.Println("\nOperation cancelled.")
			os.Exit(0)
		}
		return "", fmt.Errorf("failed to get project name: %w", err)
	}

	// Validate and sanitize project name
	sanitizedName, err := ValidateAndSanitizeName(projectName, nil)
	if err != nil {
		return "", err
	}

	// Check for suspicious patterns
	if err := CheckForSuspiciousPatterns(sanitizedName); err != nil {
		return "", err
	}

	return sanitizedName, nil
}

// promptForFirstService prompts the user for the first service details
func promptForFirstService() (string, string, error) {
	// Discover available templates
	templates, err := templating.DiscoverTemplates(templatesFS)
	if err != nil {
		return "", "", fmt.Errorf("could not discover templates: %w", err)
	}

	if len(templates) == 0 {
		return "", "", fmt.Errorf("no templates found")
	}

	// Create template options for selection
	var templateOptions []string
	templateMap := make(map[string]string)
	for _, template := range templates {
		templateOptions = append(templateOptions, fmt.Sprintf("%s - %s", template.Name, template.Description))
		templateMap[fmt.Sprintf("%s - %s", template.Name, template.Description)] = template.Name
	}

	// Prompt for template selection
	var selectedTemplateOption string
	templateQuestion := &survey.Select{
		Message: "Choose a template for your first service:",
		Options: templateOptions,
		Help:    "This will be used to scaffold your first service",
	}
	err = survey.AskOne(templateQuestion, &selectedTemplateOption)
	if err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			fmt.Println("\nOperation cancelled.")
			os.Exit(0)
		}
		return "", "", fmt.Errorf("could not select template: %w", err)
	}

	selectedTemplate := templateMap[selectedTemplateOption]

	// Validate template name for security
	if err := ValidateTemplateName(selectedTemplate); err != nil {
		return "", "", fmt.Errorf("invalid template name: %w", err)
	}

	// Prompt for service name
	var serviceName string
	servicePrompt := &survey.Input{
		Message: "What is your service name?",
		Default: "frontend",
		Help:    "This will be used as the service directory name",
	}
	err = survey.AskOne(servicePrompt, &serviceName, survey.WithValidator(survey.Required))
	if err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			fmt.Println("\nOperation cancelled.")
			os.Exit(0)
		}
		return "", "", fmt.Errorf("could not get service name: %w", err)
	}

	// Validate and sanitize service name
	sanitizedServiceName, err := ValidateAndSanitizeName(serviceName, nil)
	if err != nil {
		return "", "", err
	}

	// Check for suspicious patterns
	if err := CheckForSuspiciousPatterns(sanitizedServiceName); err != nil {
		return "", "", err
	}

	return sanitizedServiceName, selectedTemplate, nil
}

// createProjectDirectories creates the project and service directories
func createProjectDirectories(projectName, serviceName string) error {
	// Validate and sanitize paths
	sanitizedProjectName, err := ValidateAndSanitizePath(projectName, nil)
	if err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}

	sanitizedServiceName, err := ValidateAndSanitizePath(serviceName, nil)
	if err != nil {
		return fmt.Errorf("invalid service name: %w", err)
	}

	// Create project directory
	if err := os.MkdirAll(sanitizedProjectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create service directory
	servicePath := filepath.Join(sanitizedProjectName, sanitizedServiceName)
	if err := os.MkdirAll(servicePath, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	return nil
}

// scaffoldService runs the scaffolding process for the service
func scaffoldService(templateName, servicePath string) error {
	// Load the template manifest
	templateInfo, err := templating.GetTemplateInfo(templatesFS, templateName)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	// For now, we'll use default parameters
	// In the future, this could prompt for template-specific parameters
	parameterValues := map[string]interface{}{
		"ProjectName":      filepath.Base(servicePath),
		"Owner":            "Open Workbench",
		"IncludeTesting":   true,
		"TestingFramework": "Jest",
		"IncludeTailwind":  true,
		"IncludeDocker":    true,
		"InstallDeps":      false, // Don't install deps during init
		"InitGit":          false, // Don't init git during init
	}

	// Create a template processor
	processor := templating.NewTemplateProcessor(templateInfo.Manifest, parameterValues, false)

	// Execute the scaffolding process
	err = processor.ScaffoldProject(templatesFS, templateName, servicePath)
	if err != nil {
		return fmt.Errorf("failed to scaffold service: %w", err)
	}

	// Execute post-scaffolding actions
	err = processor.ExecutePostScaffoldActions(servicePath)
	if err != nil {
		return fmt.Errorf("failed to execute post-scaffold actions: %w", err)
	}

	return nil
}

// createWorkbenchManifest creates and writes the workbench.yaml file
func createWorkbenchManifest(projectName, serviceName, templateName string) error {
	manifest := WorkbenchManifest{
		APIVersion: "openworkbench.io/v1alpha1",
		Kind:       "Project",
		Metadata: ProjectMetadata{
			Name: projectName,
		},
		Services: map[string]Service{
			serviceName: {
				Template: templateName,
				Path:     filepath.Join(".", serviceName),
			},
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(&manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write to file
	manifestPath := filepath.Join(projectName, "workbench.yaml")
	err = os.WriteFile(manifestPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write workbench.yaml: %w", err)
	}

	return nil
}

// printSuccessMessage prints a success message with next steps
func printSuccessMessage(projectName, serviceName string) {
	fmt.Println("------------------------------------")
	fmt.Printf("✅ Success! Your new project '%s' is ready.\n", projectName)
	fmt.Println()
	fmt.Printf("📁 Project structure:\n")
	fmt.Printf("  %s/\n", projectName)
	fmt.Printf("  ├── workbench.yaml\n")
	fmt.Printf("  └── %s/\n", serviceName)
	fmt.Println()
	fmt.Println("🚀 Next steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  om add service  # Add more services to your project")
	fmt.Println("  om run          # Run your project (when implemented)")
	fmt.Println("  om deploy       # Deploy your project (when implemented)")
}

// isValidProjectName validates that a project/service name follows the required format
func isValidProjectName(name string) bool {
	if name == "" {
		return false
	}

	// Check for valid characters: lowercase letters, numbers, and hyphens only
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
			return false
		}
	}

	// Must start with a letter
	if len(name) > 0 && (name[0] < 'a' || name[0] > 'z') {
		return false
	}

	// Must end with a letter or number
	if len(name) > 0 {
		lastChar := name[len(name)-1]
		if !((lastChar >= 'a' && lastChar <= 'z') || (lastChar >= '0' && lastChar <= '9')) {
			return false
		}
	}

	return true
}
