// Package templating provides the core templating system for the Open Workbench CLI.
// This package implements dynamic template discovery, parameter processing, and
// file generation capabilities with support for conditional logic and validation.
package templating

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateProcessor handles dynamic template processing with conditional logic.
// This struct manages the complete template processing workflow, including
// file generation, variable substitution, and post-scaffolding actions.
type TemplateProcessor struct {
	manifest *TemplateManifest      // The template manifest containing configuration
	values   map[string]interface{} // Collected parameter values for substitution
	progress *ProgressReporter      // Progress reporter for user feedback
}

// NewTemplateProcessor creates a new template processor.
// This function initializes a template processor with the given manifest and
// parameter values, preparing it for template processing operations.
//
// Parameters:
//   - manifest: The template manifest containing template configuration
//   - values: A map of parameter names to their collected values
//   - verbose: Whether to show detailed progress information
//
// Returns:
//   - A pointer to the initialized TemplateProcessor
func NewTemplateProcessor(manifest *TemplateManifest, values map[string]interface{}, verbose bool) *TemplateProcessor {
	return &TemplateProcessor{
		manifest: manifest,
		values:   values,
		progress: NewProgressReporter(0, verbose), // Will be updated with actual steps
	}
}

// ProcessTemplate processes a template string with the provided values.
// This function applies Go template processing to a string, substituting
// variables and executing conditional logic based on the collected parameters.
//
// Parameters:
//   - content: The template string to process
//
// Returns:
//   - The processed string with variables substituted
//   - An error if template processing fails
func (tp *TemplateProcessor) ProcessTemplate(content string) (string, error) {
	// Create a template with custom functions for enhanced processing
	tmpl, err := template.New("content").Funcs(tp.getTemplateFunctions()).Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute the template with the collected parameter values
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tp.values); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ProcessFileName processes a filename template and returns the processed filename.
// This function applies template processing to file names, allowing for dynamic
// file naming based on parameter values.
//
// Parameters:
//   - fileName: The template string for the filename
//
// Returns:
//   - The processed filename
//   - An error if template processing fails
func (tp *TemplateProcessor) ProcessFileName(fileName string) (string, error) {
	processed, err := tp.ProcessTemplate(fileName)
	if err != nil {
		return "", err
	}

	// If the processed filename is empty, it means the file should be skipped
	if strings.TrimSpace(processed) == "" {
		return "", nil
	}

	return processed, nil
}

// getTemplateFunctions returns custom template functions.
// This function provides additional template functions beyond the standard
// Go template functions, enabling more sophisticated template processing.
//
// Returns:
//   - A template.FuncMap containing custom template functions
func (tp *TemplateProcessor) getTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		// Equality comparison function
		"eq": func(a, b interface{}) bool {
			return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
		},
		// Inequality comparison function
		"ne": func(a, b interface{}) bool {
			return fmt.Sprintf("%v", a) != fmt.Sprintf("%v", b)
		},
		// Array contains check function
		"contains": func(slice []string, item string) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
		// String manipulation functions
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
		"trim":  strings.TrimSpace,
	}
}

// ScaffoldProject scaffolds a complete project from a template.
// This function performs the main scaffolding operation, copying template files
// and processing them with parameter substitution to create a new project.
//
// Parameters:
//   - templateFS: The embedded filesystem containing template files
//   - templateName: The name of the template to use
//   - destDir: The destination directory for the scaffolded project
//
// Returns:
//   - An error if scaffolding fails
func (tp *TemplateProcessor) ScaffoldProject(templateFS fs.FS, templateName, destDir string) error {
	sourceDir := fmt.Sprintf("templates/%s", templateName)

	// Start progress reporting
	tp.progress.StartOperation("Scaffolding project")

	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return NewFileSystemError("create destination directory", destDir, err)
	}

	// Walk through the template directory and process each file
	return fs.WalkDir(templateFS, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the template.json file itself
		if d.Name() == "template.json" {
			return nil
		}

		// Calculate the relative path from the template root
		relPath := path[len(sourceDir):]
		if relPath == "" {
			return nil // Skip the root directory
		}

		// Process the filename template
		processedFileName, err := tp.ProcessFileName(d.Name())
		if err != nil {
			return NewTemplateProcessingError(templateName, fmt.Sprintf("Failed to process filename '%s'", d.Name()), err)
		}

		// If filename is empty after processing, skip this file/directory
		if processedFileName == "" {
			return nil
		}

		// Calculate the destination path
		destPath := filepath.Join(destDir, relPath)
		destPath = filepath.Join(filepath.Dir(destPath), processedFileName)

		if d.IsDir() {
			// Create directory with appropriate permissions
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return NewFileSystemError("create directory", destPath, err)
			}
		} else {
			// Process and write the file
			if err := tp.processAndWriteFile(templateFS, path, destPath); err != nil {
				return err
			}
		}
		return nil
	})
}

// processAndWriteFile processes a single file and writes it to the destination.
// This function reads a template file, processes it with parameter substitution,
// and writes the result to the destination location.
//
// Parameters:
//   - templateFS: The embedded filesystem containing the source file
//   - sourcePath: The path to the source file in the template
//   - destPath: The destination path for the processed file
//
// Returns:
//   - An error if file processing fails
func (tp *TemplateProcessor) processAndWriteFile(templateFS fs.FS, sourcePath, destPath string) error {
	// Read the source file content
	content, err := fs.ReadFile(templateFS, sourcePath)
	if err != nil {
		return NewFileSystemError("read source file", sourcePath, err)
	}

	// Process the file content with template substitution
	processedContent, err := tp.ProcessTemplate(string(content))
	if err != nil {
		return NewTemplateProcessingError("", fmt.Sprintf("Failed to process file content: %s", sourcePath), err)
	}

	// Ensure the destination directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return NewFileSystemError("create destination directory", destDir, err)
	}

	// Write the processed content to the destination file
	if err := os.WriteFile(destPath, []byte(processedContent), 0644); err != nil {
		return NewFileSystemError("write destination file", destPath, err)
	}

	return nil
}

// ExecutePostScaffoldActions executes post-scaffolding actions.
// This function performs cleanup and setup actions after the main scaffolding
// is complete, such as file deletion and command execution.
//
// Parameters:
//   - projectDir: The directory containing the scaffolded project
//
// Returns:
//   - An error if post-scaffolding actions fail
func (tp *TemplateProcessor) ExecutePostScaffoldActions(projectDir string) error {
	// Skip if no post-scaffolding actions are defined
	if tp.manifest.PostScaffold == nil {
		return nil
	}

	// Start progress reporting for post-scaffolding actions
	tp.progress.StartOperation("Executing post-scaffolding actions")

	// Execute file deletions based on conditions
	if err := tp.executeFileDeletions(projectDir); err != nil {
		return NewTemplateProcessingError("", "Failed to execute file deletions", err)
	}

	// Execute commands based on conditions
	if err := tp.executeCommands(projectDir); err != nil {
		return NewTemplateProcessingError("", "Failed to execute commands", err)
	}

	return nil
}

// executeFileDeletions executes file deletion actions based on conditions.
// This function removes files and directories based on conditional logic
// defined in the template manifest.
//
// Parameters:
//   - projectDir: The directory containing the scaffolded project
//
// Returns:
//   - An error if file deletion fails
func (tp *TemplateProcessor) executeFileDeletions(projectDir string) error {
	// Skip if no file deletion actions are defined
	if tp.manifest.PostScaffold.FilesToDelete == nil {
		return nil
	}

	// Process each file deletion action
	for i, fileAction := range tp.manifest.PostScaffold.FilesToDelete {
		// Report progress
		tp.progress.ReportProgress(fmt.Sprintf("Checking file deletion: %s", fileAction.Path), i+1, len(tp.manifest.PostScaffold.FilesToDelete))

		// Evaluate the condition for this file deletion
		shouldDelete, err := tp.evaluateCondition(fileAction.Condition)
		if err != nil {
			return NewTemplateProcessingError("", fmt.Sprintf("Failed to evaluate condition for file deletion '%s'", fileAction.Path), err)
		}

		// Delete the file if the condition is met
		if shouldDelete {
			filePath := filepath.Join(projectDir, fileAction.Path)
			if err := os.RemoveAll(filePath); err != nil {
				return NewFileSystemError("delete file", filePath, err)
			}
			tp.progress.CompleteStep(fmt.Sprintf("Deleted file: %s", fileAction.Path), true, "")
		}
	}

	return nil
}

// executeCommands executes command actions based on conditions.
// This function runs shell commands based on conditional logic defined
// in the template manifest, such as dependency installation.
//
// Parameters:
//   - projectDir: The directory containing the scaffolded project
//
// Returns:
//   - An error if command execution fails
func (tp *TemplateProcessor) executeCommands(projectDir string) error {
	// Skip if no command actions are defined
	if tp.manifest.PostScaffold.Commands == nil {
		return nil
	}

	// Process each command action
	for i, commandAction := range tp.manifest.PostScaffold.Commands {
		shouldExecute := true

		// Report progress
		tp.progress.ReportProgress(fmt.Sprintf("Checking command: %s", commandAction.Description), i+1, len(tp.manifest.PostScaffold.Commands))

		// Evaluate the condition for this command if specified
		if commandAction.Condition != "" {
			var err error
			shouldExecute, err = tp.evaluateCondition(commandAction.Condition)
			if err != nil {
				return NewTemplateProcessingError("", fmt.Sprintf("Failed to evaluate condition for command '%s'", commandAction.Command), err)
			}
		}

		// Execute the command if the condition is met
		if shouldExecute {
			if err := tp.executeCommand(commandAction, projectDir); err != nil {
				return NewCommandExecutionError(commandAction.Command, commandAction.Description, err)
			}
		}
	}

	return nil
}

// evaluateCondition evaluates a condition string against current values.
// This function implements a simple condition evaluator that supports equality
// and inequality comparisons for boolean and string values.
//
// Parameters:
//   - condition: The condition string to evaluate (e.g., "IncludeTesting == true")
//
// Returns:
//   - true if the condition is met, false otherwise
//   - An error if the condition cannot be evaluated
func (tp *TemplateProcessor) evaluateCondition(condition string) (bool, error) {
	// This is a simplified condition evaluator
	// In a production system, you might want to use a proper expression parser

	condition = strings.TrimSpace(condition)

	// Handle simple equality conditions like "IncludeTesting == true"
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) != 2 {
			return false, fmt.Errorf("invalid condition format: %s", condition)
		}

		paramName := strings.TrimSpace(parts[0])
		expectedValue := strings.TrimSpace(parts[1])

		// Get the actual value for the parameter
		actualValue, exists := tp.values[paramName]
		if !exists {
			return false, nil
		}

		// Convert expected value to appropriate type and compare
		switch expectedValue {
		case "true":
			return actualValue == true, nil
		case "false":
			return actualValue == false, nil
		default:
			// String comparison
			return fmt.Sprintf("%v", actualValue) == expectedValue, nil
		}
	}

	// Handle inequality conditions like "TestingFramework != 'Jest'"
	if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		if len(parts) != 2 {
			return false, fmt.Errorf("invalid condition format: %s", condition)
		}

		paramName := strings.TrimSpace(parts[0])
		expectedValue := strings.TrimSpace(parts[1])

		// Get the actual value for the parameter
		actualValue, exists := tp.values[paramName]
		if !exists {
			return true, nil // If parameter doesn't exist, condition is true
		}

		// Remove quotes from expected value
		expectedValue = strings.Trim(expectedValue, "'\"")

		return fmt.Sprintf("%v", actualValue) != expectedValue, nil
	}

	return false, fmt.Errorf("unsupported condition format: %s", condition)
}

// executeCommand executes a single command.
// This function runs a shell command in the project directory.
// Currently, this is a placeholder implementation that logs the command.
// In a real implementation, you would use os/exec to run the command.
//
// Parameters:
//   - commandAction: The command action to execute
//   - projectDir: The directory to execute the command in
//
// Returns:
//   - An error if command execution fails
func (tp *TemplateProcessor) executeCommand(commandAction CommandAction, projectDir string) error {
	// Report command execution
	tp.progress.ReportCommandExecution(commandAction.Command, commandAction.Description)

	// This is a placeholder for command execution
	// In a real implementation, you would use os/exec to run the command
	// For now, we'll just print what would be executed

	fmt.Printf("Would execute: %s\n", commandAction.Command)
	fmt.Printf("Description: %s\n", commandAction.Description)

	// TODO: Implement actual command execution using os/exec
	// cmd := exec.Command("sh", "-c", commandAction.Command)
	// cmd.Dir = projectDir
	// return cmd.Run()

	// Report successful completion
	tp.progress.ReportCommandResult(commandAction.Command, true, "")
	return nil
}
