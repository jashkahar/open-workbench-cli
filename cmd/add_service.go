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

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new component to your project.",
}

var addServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Add a new service to your project.",
	Long: `Add a new service to your project. This command is smart and will automatically 
switch between interactive and direct modes based on whether you provide parameters.

Interactive Mode (no parameters):
  om add service

Direct Mode (with parameters):
  om add service --name frontend --template react-typescript --params ProjectName=my-app,Owner=John

Examples:
  # Interactive mode - prompts for all details
  om add service

  # Direct mode with all parameters
  om add service --name frontend --template react-typescript --params ProjectName=my-app,Owner=John,IncludeTesting=true,IncludeTailwind=true

  # Direct mode with minimal parameters (others will be prompted)
  om add service --name backend --template fastapi-basic

Available templates: react-typescript, nextjs-full-stack, fastapi-basic, express-api, vue-nuxt`,
	RunE: runAddService,
}

// Top-level command to list available templates
var listTemplatesCmd = &cobra.Command{
	Use:   "list-templates",
	Short: "List available templates and their parameters.",
	Long:  `Lists all available templates and their parameters to help you understand what options are available for the add service command.`,
	RunE:  runListTemplates,
}

// initAddCommands initializes all add commands
func initAddCommands() {
	addCmd.AddCommand(addServiceCmd)

	// Add flags for the service command (optional for interactive mode)
	addServiceCmd.Flags().String("name", "", "Service name (optional - will prompt if not provided)")
	addServiceCmd.Flags().String("template", "", "Template name (optional - will prompt if not provided)")
	addServiceCmd.Flags().StringToString("params", nil, "Template parameters as key=value pairs (e.g., --params IncludeTesting=true,Framework=React)")
}

func init() {
	// Commands are now initialized in initAddCommands() function
}

// runAddService executes the add service command logic - smart mode detection
func runAddService(cmd *cobra.Command, args []string) error {
	// Check if we're in direct mode (parameters provided)
	nameFlag, _ := cmd.Flags().GetString("name")
	templateFlag, _ := cmd.Flags().GetString("template")
	paramsFlag, _ := cmd.Flags().GetStringToString("params")

	isDirectMode := nameFlag != "" || templateFlag != "" || len(paramsFlag) > 0

	if isDirectMode {
		// Direct mode - use provided parameters
		return runAddServiceDirect(cmd, args)
	} else {
		// Interactive mode - prompt for all details
		return runAddServiceInteractive(cmd, args)
	}
}

// runAddServiceInteractive executes the add service command in interactive mode
func runAddServiceInteractive(cmd *cobra.Command, args []string) error {
	// Step 1: Find project root and load manifest
	projectRoot, manifest, err := findProjectRootAndLoadManifest()
	if err != nil {
		return err
	}

	// Step 2: Prompt for new service details
	serviceName, templateName, err := promptForNewService()
	if err != nil {
		return err
	}

	// Step 3: Perform critical safety checks
	if err := performSafetyChecks(manifest, projectRoot, serviceName); err != nil {
		return err
	}

	// Step 4: Create service directory
	servicePath := filepath.Join(projectRoot, serviceName)
	if err := os.MkdirAll(servicePath, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Step 5: Run the scaffolder
	if err := scaffoldService(templateName, servicePath, true, "", ""); err != nil {
		// Clean up the created directory if scaffolding fails
		os.RemoveAll(servicePath)
		return fmt.Errorf("failed to scaffold service: %w", err)
	}

	// Step 6: Update workbench.yaml (atomic update)
	if err := updateWorkbenchManifest(manifest, serviceName, templateName, projectRoot); err != nil {
		// Clean up the created directory if manifest update fails
		os.RemoveAll(servicePath)
		return fmt.Errorf("failed to update workbench.yaml: %w", err)
	}

	// Step 7: Print success message
	printAddServiceSuccessMessage(serviceName, templateName)

	return nil
}

// runAddServiceDirect executes the add service command with direct parameter specification
func runAddServiceDirect(cmd *cobra.Command, args []string) error {
	// Step 1: Find project root and load manifest
	projectRoot, manifest, err := findProjectRootAndLoadManifest()
	if err != nil {
		return err
	}

	// Step 2: Get parameters from command line flags
	serviceName, templateName, params, err := getDirectServiceParameters(cmd)
	if err != nil {
		return err
	}

	// Step 3: Perform critical safety checks
	if err := performSafetyChecks(manifest, projectRoot, serviceName); err != nil {
		return err
	}

	// Step 4: Validate template and parameters
	if err := validateTemplateAndParameters(templateName, params); err != nil {
		return err
	}

	// Step 5: Create service directory
	servicePath := filepath.Join(projectRoot, serviceName)
	if err := os.MkdirAll(servicePath, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Step 6: Run the scaffolder with direct parameters
	if err := scaffoldServiceDirect(templateName, servicePath, params); err != nil {
		// Clean up the created directory if scaffolding fails
		os.RemoveAll(servicePath)
		return fmt.Errorf("failed to scaffold service: %w", err)
	}

	// Step 7: Update workbench.yaml (atomic update)
	if err := updateWorkbenchManifest(manifest, serviceName, templateName, projectRoot); err != nil {
		// Clean up the created directory if manifest update fails
		os.RemoveAll(servicePath)
		return fmt.Errorf("failed to update workbench.yaml: %w", err)
	}

	// Step 8: Print success message
	printAddServiceSuccessMessage(serviceName, templateName)

	return nil
}

// runListTemplates lists all available templates and their parameters
func runListTemplates(cmd *cobra.Command, args []string) error {
	// Discover available templates
	templates, err := templating.DiscoverTemplates(templatesFS)
	if err != nil {
		return fmt.Errorf("could not discover templates: %w", err)
	}

	if len(templates) == 0 {
		fmt.Println("No templates found.")
		return nil
	}

	fmt.Println("Available Templates:")
	fmt.Println("===================")
	fmt.Println()

	for i, template := range templates {
		fmt.Printf("%d. %s\n", i+1, template.Name)
		fmt.Printf("   Description: %s\n", template.Description)
		fmt.Printf("   Template ID: %s\n", template.Name)

		if template.Manifest != nil && len(template.Manifest.Parameters) > 0 {
			fmt.Printf("   Parameters:\n")
			for _, param := range template.Manifest.Parameters {
				required := ""
				if param.Required {
					required = " (required)"
				}

				fmt.Printf("     - %s%s (%s)\n", param.Name, required, param.Type)
				if param.HelpText != "" {
					fmt.Printf("       %s\n", param.HelpText)
				}

				if param.Type == "select" && len(param.Options) > 0 {
					fmt.Printf("       Options: %s\n", strings.Join(param.Options, ", "))
				}

				if param.Default != nil {
					fmt.Printf("       Default: %v\n", param.Default)
				}

				if param.Condition != "" {
					fmt.Printf("       Condition: %s\n", param.Condition)
				}
			}
		}

		fmt.Println()
	}

	fmt.Println("Usage Examples:")
	fmt.Println("===============")
	fmt.Println()
	fmt.Println("# Interactive mode:")
	fmt.Println("om add service")
	fmt.Println()
	fmt.Println("# Direct mode with all parameters:")
	fmt.Println("om add service --name frontend --template react-typescript \\")
	fmt.Println("  --params ProjectName=my-app,Owner=John,IncludeTesting=true,IncludeTailwind=true")
	fmt.Println()
	fmt.Println("# Direct mode with minimal parameters (others will be prompted):")
	fmt.Println("om add service --name backend --template fastapi-basic")

	return nil
}

// findProjectRootAndLoadManifest finds the project root by searching for workbench.yaml
// and loads the manifest file
func findProjectRootAndLoadManifest() (string, *WorkbenchManifest, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	// Search for workbench.yaml in current directory and parent directories
	projectRoot, err := findWorkbenchYaml(currentDir)
	if err != nil {
		return "", nil, fmt.Errorf("could not find workbench.yaml in current or parent directories: %w", err)
	}

	// Load and parse the workbench.yaml file
	manifestPath := filepath.Join(projectRoot, "workbench.yaml")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read workbench.yaml: %w", err)
	}

	var manifest WorkbenchManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return "", nil, fmt.Errorf("failed to parse workbench.yaml: %w", err)
	}

	return projectRoot, &manifest, nil
}

// findWorkbenchYaml searches for workbench.yaml in the given directory and its parents
func findWorkbenchYaml(startDir string) (string, error) {
	current := startDir
	for {
		// Check if workbench.yaml exists in current directory
		manifestPath := filepath.Join(current, "workbench.yaml")
		if _, err := os.Stat(manifestPath); err == nil {
			return current, nil
		}

		// Move to parent directory
		parent := filepath.Dir(current)
		if parent == current {
			// We've reached the root directory
			return "", fmt.Errorf("workbench.yaml not found")
		}
		current = parent
	}
}

// promptForNewService prompts the user for the new service details
func promptForNewService() (string, string, error) {
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
		Message: "Choose a template for your new service:",
		Options: templateOptions,
		Help:    "This will be used to scaffold your new service",
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
		Default: "backend",
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

// getDirectServiceParameters extracts parameters from command line flags
func getDirectServiceParameters(cmd *cobra.Command) (string, string, map[string]interface{}, error) {
	serviceName, err := cmd.Flags().GetString("name")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to get service name: %w", err)
	}

	templateName, err := cmd.Flags().GetString("template")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to get template name: %w", err)
	}

	// Parse parameters from string map
	paramStrings, err := cmd.Flags().GetStringToString("params")
	if err != nil {
		return "", "", nil, fmt.Errorf("failed to get parameters: %w", err)
	}

	// Convert string parameters to interface{} map
	params := make(map[string]interface{})
	for key, value := range paramStrings {
		// Try to convert to appropriate type
		if value == "true" || value == "false" {
			params[key] = value == "true"
		} else if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			// Handle array values (e.g., "[item1,item2]")
			items := strings.Trim(value, "[]")
			if items == "" {
				params[key] = []string{}
			} else {
				params[key] = strings.Split(items, ",")
			}
		} else {
			params[key] = value
		}
	}

	// If service name is not provided, prompt for it
	if serviceName == "" {
		servicePrompt := &survey.Input{
			Message: "What is your service name?",
			Default: "backend",
			Help:    "This will be used as the service directory name",
		}
		err = survey.AskOne(servicePrompt, &serviceName, survey.WithValidator(survey.Required))
		if err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				fmt.Println("\nOperation cancelled.")
				os.Exit(0)
			}
			return "", "", nil, fmt.Errorf("could not get service name: %w", err)
		}
	}

	// Validate and sanitize service name
	sanitizedServiceName, err := ValidateAndSanitizeName(serviceName, nil)
	if err != nil {
		return "", "", nil, err
	}

	// Check for suspicious patterns
	if err := CheckForSuspiciousPatterns(sanitizedServiceName); err != nil {
		return "", "", nil, err
	}

	// If template name is not provided, prompt for it
	if templateName == "" {
		// Discover available templates
		templates, err := templating.DiscoverTemplates(templatesFS)
		if err != nil {
			return "", "", nil, fmt.Errorf("could not discover templates: %w", err)
		}

		if len(templates) == 0 {
			return "", "", nil, fmt.Errorf("no templates found")
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
			Message: "Choose a template for your new service:",
			Options: templateOptions,
			Help:    "This will be used to scaffold your new service",
		}
		err = survey.AskOne(templateQuestion, &selectedTemplateOption)
		if err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				fmt.Println("\nOperation cancelled.")
				os.Exit(0)
			}
			return "", "", nil, fmt.Errorf("could not select template: %w", err)
		}

		templateName = templateMap[selectedTemplateOption]
	}

	// Validate template name for security
	if err := ValidateTemplateName(templateName); err != nil {
		return "", "", nil, fmt.Errorf("invalid template name: %w", err)
	}

	return sanitizedServiceName, templateName, params, nil
}

// validateTemplateAndParameters validates the template and its parameters
func validateTemplateAndParameters(templateName string, params map[string]interface{}) error {
	// Load template manifest to validate parameters
	manifest, err := templating.LoadTemplateManifest(templatesFS, templateName)
	if err != nil {
		return fmt.Errorf("failed to load template manifest: %w", err)
	}

	// Create parameter processor to validate parameters
	processor := templating.NewParameterProcessor(manifest)

	// Validate each provided parameter
	for paramName, paramValue := range params {
		// Find the parameter definition
		var paramDef templating.Parameter
		found := false
		for _, param := range manifest.Parameters {
			if param.Name == paramName {
				paramDef = param
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("unknown parameter: %s", paramName)
		}

		// Validate the parameter value
		if err := processor.ValidateParameter(paramDef, paramValue); err != nil {
			return fmt.Errorf("invalid value for parameter %s: %w", paramName, err)
		}

		// Set the value in the processor for conditional logic
		processor.SetValue(paramName, paramValue)
	}

	// Check if all required parameters are provided
	for _, param := range manifest.Parameters {
		if param.Required {
			if _, exists := params[param.Name]; !exists {
				return fmt.Errorf("required parameter missing: %s", param.Name)
			}
		}
	}

	return nil
}

// scaffoldServiceDirect scaffolds a service with direct parameter specification
func scaffoldServiceDirect(templateName, servicePath string, params map[string]interface{}) error {
	// Load template manifest
	manifest, err := templating.LoadTemplateManifest(templatesFS, templateName)
	if err != nil {
		return fmt.Errorf("failed to load template manifest: %w", err)
	}

	// Create template processor with the provided parameters
	processor := templating.NewTemplateProcessor(manifest, params, false)

	// Scaffold the project
	if err := processor.ScaffoldProject(templatesFS, templateName, servicePath); err != nil {
		return fmt.Errorf("failed to scaffold project: %w", err)
	}

	// Execute post-scaffold actions
	if err := processor.ExecutePostScaffoldActions(servicePath); err != nil {
		return fmt.Errorf("failed to execute post-scaffold actions: %w", err)
	}

	return nil
}

// performSafetyChecks performs critical safety checks before adding the service
func performSafetyChecks(manifest *WorkbenchManifest, projectRoot, serviceName string) error {
	// Check if service already exists in manifest
	if _, exists := manifest.Services[serviceName]; exists {
		return fmt.Errorf("error: a service named '%s' already exists in your project", serviceName)
	}

	// Check if directory already exists on filesystem
	servicePath := filepath.Join(projectRoot, serviceName)
	if _, err := os.Stat(servicePath); err == nil {
		return fmt.Errorf("error: a directory named '%s' already exists", serviceName)
	}

	return nil
}

// updateWorkbenchManifest updates the workbench.yaml file with the new service
func updateWorkbenchManifest(manifest *WorkbenchManifest, serviceName, templateName, projectRoot string) error {
	// Add the new service to the manifest
	manifest.Services[serviceName] = Service{
		Template: templateName,
		Path:     filepath.Join(".", serviceName),
	}

	// Marshal to YAML
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write to file
	manifestPath := filepath.Join(projectRoot, "workbench.yaml")
	err = os.WriteFile(manifestPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write workbench.yaml: %w", err)
	}

	return nil
}

// printAddServiceSuccessMessage prints a success message for adding a service
func printAddServiceSuccessMessage(serviceName, templateName string) {
	fmt.Println("------------------------------------")
	fmt.Printf("✅ Success! Service '%s' has been added to your project.\n", serviceName)
	fmt.Println()
	fmt.Printf("📁 Service details:\n")
	fmt.Printf("  Template: %s\n", templateName)
	fmt.Printf("  Path: ./%s\n", serviceName)
	fmt.Println()
	fmt.Println("🚀 Next steps:")
	fmt.Printf("  cd %s\n", serviceName)
	fmt.Println("  om add service  # Add more services to your project")
	fmt.Println("  om run          # Run your project (when implemented)")
	fmt.Println("  om deploy       # Deploy your project (when implemented)")
}
