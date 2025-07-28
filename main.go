// Package main provides the Open Workbench CLI, a command-line tool for scaffolding
// modern web applications with pre-configured templates and best practices.
//
// The CLI supports multiple execution modes:
//   - Interactive TUI mode for template selection
//   - Simple interactive mode with default template
//   - Non-interactive mode (planned)
//
// Features:
//   - Dynamic template system with conditional logic
//   - Parameter validation and grouping
//   - Post-scaffolding actions
//   - Cross-platform support
//
// Usage:
//
//	open-workbench-cli        # Simple interactive mode
//	open-workbench-cli ui     # Enhanced TUI mode
//	open-workbench-cli create # Non-interactive mode (coming soon)
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
//   - No arguments: Simple interactive mode with default template
//   - "ui": Enhanced Terminal User Interface mode
//   - "create": Non-interactive mode (planned)
func main() {
	if len(os.Args) > 1 {
		// Handle subcommands without requiring a complex flag parsing library
		switch os.Args[1] {
		case "ui":
			runUIAndScaffold()
		case "create":
			// TODO: Implement non-interactive mode with flag parsing
			// This would allow usage like: open-workbench-cli create --name my-project --template nextjs
			fmt.Println("Non-interactive 'create' command not yet fully implemented.")
			fmt.Println("For now, use the 'ui' command or run without arguments.")
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands: 'ui'")
		}
		return
	}

	// Default behavior: run simple interactive mode with default template
	runSimpleInteractiveScaffold()
}

// runUIAndScaffold orchestrates the complete TUI-based scaffolding workflow.
// This function handles the enhanced Terminal User Interface flow, including:
//   - Template selection via TUI
//   - Parameter collection with validation
//   - Project scaffolding with dynamic template processing
//   - Post-scaffolding actions execution
//
// The function provides detailed debug output during development and
// graceful error handling with user-friendly messages.
func runUIAndScaffold() {
	fmt.Println("🚀 Starting Open Workbench UI...")

	// Launch the Terminal User Interface for template selection
	selectedTemplate, err := runTUI()
	if err != nil {
		log.Fatalf("❌ Could not start TUI: %v", err)
	}

	fmt.Printf("DEBUG: Selected template: '%s'\n", selectedTemplate)

	// Handle case where user quits the TUI without selecting a template
	if selectedTemplate == "" {
		fmt.Println("No template selected, exiting...")
		return
	}

	// Load and validate the selected template's manifest
	fmt.Println("DEBUG: Loading template manifest...")
	templateInfo, err := templating.GetTemplateInfo(templatesFS, selectedTemplate)
	if err != nil {
		log.Fatalf("❌ Could not load template manifest: %v", err)
	}
	fmt.Printf("DEBUG: Template manifest loaded: %s - %s\n", templateInfo.Name, templateInfo.Description)

	// Collect user parameters using the dynamic parameter system
	fmt.Println("DEBUG: Collecting template parameters...")
	parameterValues, err := promptForTemplateParameters(templateInfo.Manifest)
	if err != nil {
		log.Fatalf("❌ Could not collect template parameters: %v", err)
	}
	fmt.Printf("DEBUG: Collected parameters: %+v\n", parameterValues)

	// Execute the scaffolding process with collected parameters
	scaffoldAndApplyDynamic(selectedTemplate, parameterValues)
}

// runSimpleInteractiveScaffold handles the legacy interactive survey flow.
// This function provides a simpler alternative to the TUI mode by using
// a hardcoded default template (nextjs-golden-path) and collecting
// parameters through the survey library.
//
// This mode is useful for quick project creation without template selection.
func runSimpleInteractiveScaffold() {
	// Use a hardcoded default template for simplicity
	// TODO: In the future, this could present a simple template selection
	const defaultTemplate = "nextjs-golden-path"
	fmt.Printf("🚀 Welcome! Using default template: %s\n", defaultTemplate)

	// Load the default template's manifest
	templateInfo, err := templating.GetTemplateInfo(templatesFS, defaultTemplate)
	if err != nil {
		log.Fatalf("❌ Could not load template manifest: %v", err)
	}

	// Collect parameters using the dynamic system
	parameterValues, err := promptForTemplateParameters(templateInfo.Manifest)
	if err != nil {
		log.Fatalf("❌ Could not collect template parameters: %v", err)
	}

	// Execute scaffolding with the collected parameters
	scaffoldAndApplyDynamic(defaultTemplate, parameterValues)
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
		log.Fatalf("❌ Could not load template manifest: %v", err)
	}

	// Extract the project name from collected parameters
	// This is a required parameter for all templates
	projectName, ok := parameterValues["ProjectName"].(string)
	if !ok {
		log.Fatalf("❌ ProjectName parameter is required")
	}

	destDir := projectName

	fmt.Printf("📂 Scaffolding project in './%s'...\n", destDir)

	// Create a template processor with the manifest and parameter values
	processor := templating.NewTemplateProcessor(templateInfo.Manifest, parameterValues)

	// Execute the main scaffolding process
	// This handles file copying, template processing, and variable substitution
	err = processor.ScaffoldProject(templatesFS, templateName, destDir)
	if err != nil {
		log.Fatalf("❌ Failed to scaffold project: %v", err)
	}

	fmt.Println("✏️  Applying templates...")
	// Note: Template processing is now handled by the ScaffoldProject method

	// Execute post-scaffolding actions such as file deletion and command execution
	fmt.Println("🔧 Executing post-scaffolding actions...")
	err = processor.ExecutePostScaffoldActions(destDir)
	if err != nil {
		log.Fatalf("❌ Failed to execute post-scaffolding actions: %v", err)
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

	// Process each parameter group
	for groupName, params := range groups {
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
