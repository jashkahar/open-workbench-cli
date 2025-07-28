package template

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplateProcessor handles dynamic template processing with conditional logic
type TemplateProcessor struct {
	manifest *TemplateManifest
	values   map[string]interface{}
}

// NewTemplateProcessor creates a new template processor
func NewTemplateProcessor(manifest *TemplateManifest, values map[string]interface{}) *TemplateProcessor {
	return &TemplateProcessor{
		manifest: manifest,
		values:   values,
	}
}

// ProcessTemplate processes a template string with the provided values
func (tp *TemplateProcessor) ProcessTemplate(content string) (string, error) {
	// Create a template with custom functions
	tmpl, err := template.New("content").Funcs(tp.getTemplateFunctions()).Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tp.values); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ProcessFileName processes a filename template and returns the processed filename
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

// getTemplateFunctions returns custom template functions
func (tp *TemplateProcessor) getTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"eq": func(a, b interface{}) bool {
			return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
		},
		"ne": func(a, b interface{}) bool {
			return fmt.Sprintf("%v", a) != fmt.Sprintf("%v", b)
		},
		"contains": func(slice []string, item string) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
		"trim":  strings.TrimSpace,
	}
}

// ScaffoldProject scaffolds a complete project from a template
func (tp *TemplateProcessor) ScaffoldProject(templateFS fs.FS, templateName, destDir string) error {
	sourceDir := filepath.Join("templates", templateName)

	// Create destination directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Walk through the template directory
	return fs.WalkDir(templateFS, sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the template.json file itself
		if d.Name() == "template.json" {
			return nil
		}

		// Calculate relative path from template root
		relPath := path[len(sourceDir):]
		if relPath == "" {
			return nil // Skip the root directory
		}

		// Process the filename
		processedFileName, err := tp.ProcessFileName(d.Name())
		if err != nil {
			return fmt.Errorf("failed to process filename %s: %w", d.Name(), err)
		}

		// If filename is empty after processing, skip this file/directory
		if processedFileName == "" {
			return nil
		}

		// Calculate destination path
		destPath := filepath.Join(destDir, relPath)
		destPath = filepath.Join(filepath.Dir(destPath), processedFileName)

		if d.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, 0755)
		} else {
			// Process and write file
			return tp.processAndWriteFile(templateFS, path, destPath)
		}
	})
}

// processAndWriteFile processes a single file and writes it to the destination
func (tp *TemplateProcessor) processAndWriteFile(templateFS fs.FS, sourcePath, destPath string) error {
	// Read source file
	content, err := fs.ReadFile(templateFS, sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", sourcePath, err)
	}

	// Process the content
	processedContent, err := tp.ProcessTemplate(string(content))
	if err != nil {
		return fmt.Errorf("failed to process file content %s: %w", sourcePath, err)
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Write processed content to destination
	if err := os.WriteFile(destPath, []byte(processedContent), 0644); err != nil {
		return fmt.Errorf("failed to write destination file %s: %w", destPath, err)
	}

	return nil
}

// ExecutePostScaffoldActions executes post-scaffolding actions
func (tp *TemplateProcessor) ExecutePostScaffoldActions(projectDir string) error {
	if tp.manifest.PostScaffold == nil {
		return nil
	}

	// Execute file deletions
	if err := tp.executeFileDeletions(projectDir); err != nil {
		return fmt.Errorf("failed to execute file deletions: %w", err)
	}

	// Execute commands
	if err := tp.executeCommands(projectDir); err != nil {
		return fmt.Errorf("failed to execute commands: %w", err)
	}

	return nil
}

// executeFileDeletions executes file deletion actions based on conditions
func (tp *TemplateProcessor) executeFileDeletions(projectDir string) error {
	if tp.manifest.PostScaffold.FilesToDelete == nil {
		return nil
	}

	for _, fileAction := range tp.manifest.PostScaffold.FilesToDelete {
		shouldDelete, err := tp.evaluateCondition(fileAction.Condition)
		if err != nil {
			return fmt.Errorf("failed to evaluate condition for file deletion %s: %w", fileAction.Path, err)
		}

		if shouldDelete {
			filePath := filepath.Join(projectDir, fileAction.Path)
			if err := os.RemoveAll(filePath); err != nil {
				return fmt.Errorf("failed to delete file %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// executeCommands executes command actions based on conditions
func (tp *TemplateProcessor) executeCommands(projectDir string) error {
	if tp.manifest.PostScaffold.Commands == nil {
		return nil
	}

	for _, commandAction := range tp.manifest.PostScaffold.Commands {
		shouldExecute := true

		if commandAction.Condition != "" {
			var err error
			shouldExecute, err = tp.evaluateCondition(commandAction.Condition)
			if err != nil {
				return fmt.Errorf("failed to evaluate condition for command %s: %w", commandAction.Command, err)
			}
		}

		if shouldExecute {
			if err := tp.executeCommand(commandAction, projectDir); err != nil {
				return fmt.Errorf("failed to execute command %s: %w", commandAction.Command, err)
			}
		}
	}

	return nil
}

// evaluateCondition evaluates a condition string against current values
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

		actualValue, exists := tp.values[paramName]
		if !exists {
			return false, nil
		}

		// Convert expected value to appropriate type
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

// executeCommand executes a single command
func (tp *TemplateProcessor) executeCommand(commandAction CommandAction, projectDir string) error {
	// This is a placeholder for command execution
	// In a real implementation, you would use os/exec to run the command
	// For now, we'll just print what would be executed

	fmt.Printf("Would execute: %s\n", commandAction.Command)
	fmt.Printf("Description: %s\n", commandAction.Description)

	// TODO: Implement actual command execution using os/exec
	// cmd := exec.Command("sh", "-c", commandAction.Command)
	// cmd.Dir = projectDir
	// return cmd.Run()

	return nil
}
