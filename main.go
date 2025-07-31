// Package main provides the Open Workbench CLI, a command-line tool for scaffolding
// modern web applications with pre-configured templates and best practices.
//
// The CLI supports multiple execution modes:
//   - Interactive mode for guided project creation
//   - Non-interactive CLI mode with command-line flags
//
// Features:
//   - Dynamic template system with conditional logic
//   - Parameter validation and grouping
//   - Post-scaffolding actions
//   - Cross-platform support
//
// Usage:
//
//	open-workbench-cli        # Interactive mode
//	open-workbench-cli create # CLI mode with flags
package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/jashkahar/open-workbench-cli/internal/templating"
)

// templatesFS embeds the templates directory into the binary.
// This allows the CLI to be distributed as a single executable
// without requiring external template files.
//
//go:embed templates
var templatesFS embed.FS

// main is the entry point for the Open Workbench CLI application.
// It handles command-line argument parsing and routes to the appropriate
// execution mode based on the provided arguments.
//
// Supported commands:
//   - No arguments: Interactive mode with template selection
//   - "create": CLI mode with command-line flags
func main() {
	// Check for command-line arguments
	if len(os.Args) > 1 {
		// Handle subcommands
		switch os.Args[1] {
		case "create":
			runCLICreate()
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands:")
			fmt.Println("  open-workbench-cli          # Interactive mode")
			fmt.Println("  open-workbench-cli create   # CLI mode with flags")
			fmt.Println()
			fmt.Println("Run 'open-workbench-cli create --help' for detailed CLI usage")
			fmt.Println("Run 'open-workbench-cli' for interactive mode")
			os.Exit(1)
		}
		return
	}

	// Default behavior: run interactive mode
	runInteractiveScaffold()
}

// runInteractiveScaffold handles the interactive survey flow.
// This function provides a user-friendly way to create projects by allowing
// template selection and collecting parameters through the survey library.
//
// This mode is useful for guided project creation with template selection.
func runInteractiveScaffold() {
	// Discover available templates
	templates, err := templating.DiscoverTemplates(templatesFS)
	if err != nil {
		fmt.Printf("❌ Could not discover templates: %v\n", err)
		os.Exit(1)
	}

	if len(templates) == 0 {
		fmt.Println("❌ No templates found")
		os.Exit(1)
	}

	// Create template options for selection
	var templateOptions []string
	templateMap := make(map[string]string)
	for _, template := range templates {
		templateOptions = append(templateOptions, fmt.Sprintf("%s - %s", template.Name, template.Description))
		templateMap[fmt.Sprintf("%s - %s", template.Name, template.Description)] = template.Name
	}

	// Prompt user to select a template
	var selectedTemplateOption string
	templateQuestion := &survey.Select{
		Message: "Choose a template:",
		Options: templateOptions,
	}
	err = survey.AskOne(templateQuestion, &selectedTemplateOption)
	if err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			fmt.Println("\nOperation cancelled.")
			os.Exit(0)
		}
		fmt.Printf("❌ Could not select template: %v\n", err)
		os.Exit(1)
	}

	// Get the actual template name
	selectedTemplate := templateMap[selectedTemplateOption]
	fmt.Printf("🚀 Selected template: %s\n", selectedTemplate)

	// Load the selected template's manifest
	templateInfo, err := templating.GetTemplateInfo(templatesFS, selectedTemplate)
	if err != nil {
		fmt.Printf("❌ %s\n", templating.FormatErrorForUser(err))
		os.Exit(1)
	}

	// Collect parameters using the dynamic system
	parameterValues, err := promptForTemplateParameters(templateInfo.Manifest)
	if err != nil {
		fmt.Printf("❌ Could not collect template parameters: %v\n", err)
		fmt.Println("\nPossible solutions:")
		fmt.Println("• Check that all required parameters are provided")
		fmt.Println("• Ensure parameter values are valid")
		os.Exit(1)
	}

	// Execute scaffolding with the collected parameters
	scaffoldAndApplyDynamic(selectedTemplate, parameterValues)
}

// runCLICreate handles non-interactive project creation using command-line flags.
// This function parses command-line arguments and creates projects without
// any interactive prompts, making it suitable for automation and scripting.
//
// Usage: open-workbench-cli create <template> <project-name> [flags]
func runCLICreate() {
	// Check for help flag first
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--help" || os.Args[i] == "-h" {
			fmt.Println("Usage: open-workbench-cli create <template> <project-name> [flags]")
			fmt.Println()
			fmt.Println("Arguments:")
			fmt.Println("  template      Template to use (nextjs-full-stack, react-typescript, etc.)")
			fmt.Println("  project-name  Name of the project to create")
			fmt.Println()
			fmt.Println("Flags:")
			fmt.Println("  --owner string           Project owner (required)")
			fmt.Println("  --no-testing            Disable testing framework")
			fmt.Println("  --no-tailwind           Disable Tailwind CSS")
			fmt.Println("  --no-docker             Disable Docker configuration")
			fmt.Println("  --no-install-deps       Skip dependency installation")
			fmt.Println("  --no-git                Skip Git repository initialization")
			fmt.Println("  --testing-framework     Testing framework (Jest/Vitest)")
			fmt.Println("  --help                   Show this help message")
			fmt.Println()
			fmt.Println("Examples:")
			fmt.Println("  open-workbench-cli create nextjs-full-stack my-app --owner=\"John Doe\"")
			fmt.Println("  open-workbench-cli create react-typescript my-react-app --no-testing --no-tailwind")
			fmt.Println("  open-workbench-cli create express-api my-api --owner=\"Dev Team\" --docker")
			return
		}
	}

	// Check if we have enough arguments
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing required arguments")
		fmt.Println("Usage: open-workbench-cli create <template> <project-name> --owner=\"Your Name\"")
		fmt.Println()
		fmt.Println("Run 'open-workbench-cli create --help' for detailed usage and examples")
		os.Exit(1)
	}

	templateName := os.Args[2]
	projectName := os.Args[3]

	// Simple flag parsing
	owner := ""
	includeTesting := true
	includeTailwind := true
	includeDocker := true
	installDeps := true
	initGit := true
	testingFramework := "Jest"

	// Parse remaining arguments as flags
	for i := 4; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch {
		case arg == "--help" || arg == "-h":
			fmt.Println("Usage: open-workbench-cli create <template> <project-name> [flags]")
			fmt.Println()
			fmt.Println("Arguments:")
			fmt.Println("  template      Template to use (nextjs-full-stack, react-typescript, etc.)")
			fmt.Println("  project-name  Name of the project to create")
			fmt.Println()
			fmt.Println("Flags:")
			fmt.Println("  --owner string           Project owner (required)")
			fmt.Println("  --no-testing            Disable testing framework")
			fmt.Println("  --no-tailwind           Disable Tailwind CSS")
			fmt.Println("  --no-docker             Disable Docker configuration")
			fmt.Println("  --no-install-deps       Skip dependency installation")
			fmt.Println("  --no-git                Skip Git repository initialization")
			fmt.Println("  --testing-framework     Testing framework (Jest/Vitest)")
			fmt.Println("  --help                   Show this help message")
			fmt.Println()
			fmt.Println("Examples:")
			fmt.Println("  open-workbench-cli create nextjs-full-stack my-app --owner=\"John Doe\"")
			fmt.Println("  open-workbench-cli create react-typescript my-react-app --no-testing --no-tailwind")
			fmt.Println("  open-workbench-cli create express-api my-api --owner=\"Dev Team\" --docker")
			return
		case strings.HasPrefix(arg, "--owner="):
			owner = strings.TrimPrefix(arg, "--owner=")
		case arg == "--no-testing":
			includeTesting = false
		case arg == "--no-tailwind":
			includeTailwind = false
		case arg == "--no-docker":
			includeDocker = false
		case arg == "--no-install-deps":
			installDeps = false
		case arg == "--no-git":
			initGit = false
		case strings.HasPrefix(arg, "--testing-framework="):
			testingFramework = strings.TrimPrefix(arg, "--testing-framework=")
		}
	}

	// Validate required flags
	if owner == "" {
		fmt.Println("Error: --owner flag is required")
		fmt.Println("Usage: open-workbench-cli create <template> <project-name> --owner=\"Your Name\"")
		fmt.Println()
		fmt.Println("Run 'open-workbench-cli create --help' for detailed usage and examples")
		os.Exit(1)
	}

	// Validate project name format
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(projectName) {
		fmt.Println("Error: Project name can only contain lowercase letters, numbers, and hyphens")
		os.Exit(1)
	}

	// Load the template manifest to validate it exists
	_, err := templating.GetTemplateInfo(templatesFS, templateName)
	if err != nil {
		fmt.Printf("❌ %s\n", templating.FormatErrorForUser(err))
		os.Exit(1)
	}

	// Create parameter values from flags
	parameterValues := map[string]interface{}{
		"ProjectName":      projectName,
		"Owner":            owner,
		"IncludeTesting":   includeTesting,
		"TestingFramework": testingFramework,
		"IncludeTailwind":  includeTailwind,
		"IncludeDocker":    includeDocker,
		"InstallDeps":      installDeps,
		"InitGit":          initGit,
	}

	// Execute scaffolding
	scaffoldAndApplyDynamic(templateName, parameterValues)
}

// scaffoldAndApplyDynamic orchestrates the complete project scaffolding process
// using the new dynamic template system. This function handles:
//   - Template manifest loading and validation
//   - Project directory creation
//   - Template file processing with parameter substitution
//   - Post-scaffolding actions execution
//
// Parameters:
//   - templateName: The name of the template to use for scaffolding
//   - parameterValues: A map of parameter names to their collected values
//
// The function provides detailed progress feedback and comprehensive error handling.
func scaffoldAndApplyDynamic(templateName string, parameterValues map[string]interface{}) {
	// Load and validate the template manifest
	templateInfo, err := templating.GetTemplateInfo(templatesFS, templateName)
	if err != nil {
		// Use the new error formatting system
		fmt.Printf("❌ %s\n", templating.FormatErrorForUser(err))
		os.Exit(1)
	}

	// Extract the project name from collected parameters
	// This is a required parameter for all templates
	projectName, ok := parameterValues["ProjectName"].(string)
	if !ok {
		fmt.Println("❌ ProjectName parameter is required")
		fmt.Println("\nPossible solutions:")
		fmt.Println("• Ensure you provided a project name during setup")
		fmt.Println("• Check that the template configuration is correct")
		os.Exit(1)
	}

	destDir := projectName

	fmt.Printf("📂 Scaffolding project in './%s'...\n", destDir)

	// Create a template processor with the manifest and parameter values
	// Use verbose mode for detailed progress reporting
	processor := templating.NewTemplateProcessor(templateInfo.Manifest, parameterValues, true)

	// Execute the main scaffolding process
	// This handles file copying, template processing, and variable substitution
	err = processor.ScaffoldProject(templatesFS, templateName, destDir)
	if err != nil {
		fmt.Printf("❌ %s\n", templating.FormatErrorForUser(err))
		os.Exit(1)
	}

	// Execute post-scaffolding actions such as file deletion and command execution
	err = processor.ExecutePostScaffoldActions(destDir)
	if err != nil {
		fmt.Printf("❌ %s\n", templating.FormatErrorForUser(err))
		os.Exit(1)
	}

	// Provide success feedback to the user
	fmt.Println("------------------------------------")
	fmt.Printf("✅ Success! Your new project '%s' is ready.\n", projectName)
}

// promptForTemplateParameters collects all required parameters for a template
// using the dynamic parameter system. This function:
//   - Groups parameters by their defined groups for better UX
//   - Prompts for each parameter with appropriate validation
//   - Handles conditional parameter visibility
//   - Returns a map of parameter names to their collected values
//
// Parameters:
//   - manifest: The template manifest containing parameter definitions
//
// Returns:
//   - A map of parameter names to their collected values
//   - An error if parameter collection fails
func promptForTemplateParameters(manifest *templating.TemplateManifest) (map[string]interface{}, error) {
	// Create a parameter processor to handle parameter logic
	parameterProcessor := templating.NewParameterProcessor(manifest)
	values := make(map[string]interface{})

	// Get parameters organized by their groups for better user experience
	groups := parameterProcessor.GetParameterGroups()

	// Define the order in which groups should be processed
	groupOrder := []string{"Project Details", "Testing & Quality", "Styling", "Deployment", "Final Steps"}

	// Process each parameter group in the defined order
	for _, groupName := range groupOrder {
		params, exists := groups[groupName]
		if !exists || len(params) == 0 {
			continue
		}

		// Display group header with visual separator
		fmt.Printf("\n📋 %s\n", groupName)
		fmt.Println(strings.Repeat("-", len(groupName)+4))

		// Prompt for each parameter in the group
		for _, param := range params {
			value, err := promptForParameter(param)
			if err != nil {
				return nil, err
			}

			// Store the collected value and update the processor state
			values[param.Name] = value
			parameterProcessor.SetValue(param.Name, value)
		}
	}

	// Process any remaining groups that weren't in the predefined order
	for groupName, params := range groups {
		// Skip groups we've already processed
		alreadyProcessed := false
		for _, processedGroup := range groupOrder {
			if groupName == processedGroup {
				alreadyProcessed = true
				break
			}
		}
		if alreadyProcessed {
			continue
		}

		// Skip empty groups
		if len(params) == 0 {
			continue
		}

		// Display group header with visual separator
		fmt.Printf("\n📋 %s\n", groupName)
		fmt.Println(strings.Repeat("-", len(groupName)+4))

		// Prompt for each parameter in the group
		for _, param := range params {
			value, err := promptForParameter(param)
			if err != nil {
				return nil, err
			}

			// Store the collected value and update the processor state
			values[param.Name] = value
			parameterProcessor.SetValue(param.Name, value)
		}
	}

	return values, nil
}

// promptForParameter prompts for a single parameter using the survey library.
// This function handles different parameter types (string, boolean, select, multiselect)
// and provides appropriate validation and user interface elements.
//
// Parameters:
//   - param: The parameter definition containing type, validation rules, and UI options
//
// Returns:
//   - The collected parameter value as an interface{}
//   - An error if parameter collection fails
func promptForParameter(param templating.Parameter) (interface{}, error) {
	var questions []*survey.Question

	// Create appropriate survey question based on parameter type
	switch param.Type {
	case "string":
		questions = []*survey.Question{
			{
				Name:     param.Name,
				Prompt:   &survey.Input{Message: param.Prompt},
				Validate: createStringValidator(param),
			},
		}
	case "boolean":
		// Extract default value for boolean parameters
		defaultValue := false
		if param.Default != nil {
			if boolVal, ok := param.Default.(bool); ok {
				defaultValue = boolVal
			}
		}
		questions = []*survey.Question{
			{
				Name: param.Name,
				Prompt: &survey.Confirm{
					Message: param.Prompt,
					Default: defaultValue,
				},
			},
		}
	case "select":
		questions = []*survey.Question{
			{
				Name: param.Name,
				Prompt: &survey.Select{
					Message: param.Prompt,
					Options: param.Options,
					Default: param.Default,
				},
			},
		}
	case "multiselect":
		questions = []*survey.Question{
			{
				Name: param.Name,
				Prompt: &survey.MultiSelect{
					Message: param.Prompt,
					Options: param.Options,
				},
			},
		}
	default:
		return nil, fmt.Errorf("unsupported parameter type: %s", param.Type)
	}

	// Display help text if provided for the parameter
	if param.HelpText != "" {
		fmt.Printf("💡 %s\n", param.HelpText)
	}

	// Execute the survey and collect the result
	result := make(map[string]interface{})
	err := survey.Ask(questions, &result)
	if err != nil {
		// Handle user interruption (Ctrl+C) gracefully
		if errors.Is(err, terminal.InterruptErr) {
			fmt.Println("\nOperation cancelled.")
			os.Exit(0)
		}
		return nil, err
	}

	return result[param.Name], nil
}

// createStringValidator creates a survey validator function for string parameters.
// This function handles both required field validation and custom regex validation
// as defined in the parameter's validation rules.
//
// Parameters:
//   - param: The parameter definition containing validation rules
//
// Returns:
//   - A survey.Validator function that validates string input
func createStringValidator(param templating.Parameter) survey.Validator {
	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			// Check if the field is required and empty
			if param.Required && str == "" {
				return fmt.Errorf("this field is required")
			}

			// Apply custom regex validation if specified
			if param.Validation != nil && param.Validation.Regex != "" {
				matched, err := regexp.MatchString(param.Validation.Regex, str)
				if err != nil {
					return fmt.Errorf("invalid validation pattern")
				}
				if !matched {
					// Use custom error message if provided, otherwise use generic message
					if param.Validation.ErrorMessage != "" {
						return fmt.Errorf(param.Validation.ErrorMessage)
					}
					return fmt.Errorf("value does not match required pattern")
				}
			}
		}
		return nil
	}
}

// =============================================================================
// LEGACY FUNCTIONS FOR BACKWARD COMPATIBILITY
// These functions will be removed once the new dynamic system is fully implemented
// =============================================================================

// TemplateData represents the legacy template data structure used by the old
// templating system. This struct is maintained for backward compatibility
// and will be removed in future versions.
type TemplateData struct {
	ProjectName string // The name of the project to be created
	Owner       string // The owner/maintainer of the project
}

// scaffoldAndApply contains the legacy scaffolding logic that uses the old
// template system. This function is kept for backward compatibility but
// should not be used in new code. Use scaffoldAndApplyDynamic instead.
//
// Parameters:
//   - templateName: The name of the template to use
//   - data: The legacy template data structure
func scaffoldAndApply(templateName string, data *TemplateData) {
	sourceDir := "templates/" + templateName
	destDir := data.ProjectName

	fmt.Printf("📂 Scaffolding project in './%s'...\n", destDir)
	err := scaffoldProject(templatesFS, sourceDir, destDir)
	if err != nil {
		log.Fatalf("❌ Failed to scaffold project: %v", err)
	}

	fmt.Println("✏️  Applying templates...")
	err = applyTemplates(destDir, data)
	if err != nil {
		log.Fatalf("❌ Failed to apply templates: %v", err)
	}

	fmt.Println("------------------------------------")
	fmt.Printf("✅ Success! Your new project '%s' is ready.\n", data.ProjectName)
}

// promptForProjectDetails uses the survey library to collect basic project details
// for the legacy template system. This function is maintained for backward
// compatibility and will be removed in future versions.
//
// Returns:
//   - A TemplateData struct with collected project information
//   - An error if parameter collection fails
func promptForProjectDetails() (*TemplateData, error) {
	data := &TemplateData{}
	questions := []*survey.Question{
		{
			Name:     "ProjectName",
			Prompt:   &survey.Input{Message: "What is your project name?"},
			Validate: survey.Required,
		},
		{
			Name:     "Owner",
			Prompt:   &survey.Input{Message: "Who is the owner of this project?"},
			Validate: survey.Required,
		},
	}
	err := survey.Ask(questions, data)
	if err != nil {
		// Handle user interruption (Ctrl+C) gracefully
		if errors.Is(err, terminal.InterruptErr) {
			fmt.Println("\nOperation cancelled.")
			os.Exit(0)
		}
		return nil, err
	}
	return data, nil
}

// scaffoldProject copies the entire directory structure from a source to a destination
// using the embedded filesystem. This function is part of the legacy template system
// and will be removed in future versions.
//
// Parameters:
//   - sourceFS: The embedded filesystem containing templates
//   - sourceDir: The source directory path within the filesystem
//   - destDir: The destination directory on the local filesystem
//
// Returns:
//   - An error if the scaffolding operation fails
func scaffoldProject(sourceFS fs.FS, sourceDir, destDir string) error {
	// Check if destination directory already exists
	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", destDir)
	}

	// Walk through the source directory and copy files
	return fs.WalkDir(sourceFS, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate the relative path from source directory
		outPath := filepath.Join(destDir, path[len(sourceDir):])

		if d.IsDir() {
			// Create directory with appropriate permissions
			return os.MkdirAll(outPath, 0755)
		}

		// Read and write file with appropriate permissions
		fileBytes, err := fs.ReadFile(sourceFS, path)
		if err != nil {
			return err
		}
		return os.WriteFile(outPath, fileBytes, 0644)
	})
}

// applyTemplates walks through the newly created project directory and processes
// Go templates in all files. This function is part of the legacy template system
// and will be removed in future versions.
//
// Parameters:
//   - destDir: The destination directory containing the scaffolded project
//   - data: The template data to use for variable substitution
//
// Returns:
//   - An error if template processing fails
func applyTemplates(destDir string, data *TemplateData) error {
	return filepath.Walk(destDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Parse and execute Go template
		tmpl, err := template.New(info.Name()).Parse(string(content))
		if err != nil {
			// Log warning but continue processing other files
			fmt.Printf("⚠️  Could not parse template for %s: %v. Skipping.\n", path, err)
			return nil
		}

		// Write processed content back to file
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		return tmpl.Execute(file, data)
	})
}
