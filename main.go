package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// TemplateData struct remains the same
type TemplateData struct {
	ProjectName string
	Owner       string
}

func main() {
	if len(os.Args) > 1 {
		// This is a simple way to handle subcommands without a complex library.
		switch os.Args[1] {
		case "ui":
			runUIAndScaffold()
		case "create":
			// For non-interactive mode with flags, e.g., `open-workbench-cli create --name ...`
			// (Implementation would require a flag parsing library like `flag` or `cobra`)
			fmt.Println("Non-interactive 'create' command not yet fully implemented.")
			fmt.Println("For now, use the 'ui' command or run without arguments.")
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			fmt.Println("Available commands: 'ui'")
		}
		return
	}

	// If no command is provided, run the simple interactive survey (our original logic)
	runSimpleInteractiveScaffold()
}

// runUIAndScaffold handles the TUI flow
func runUIAndScaffold() {
	fmt.Println("🚀 Starting Open Workbench UI...")
	selectedTemplate, err := runTUI()
	if err != nil {
		log.Fatalf("❌ Could not start TUI: %v", err)
	}

	fmt.Printf("DEBUG: Selected template: '%s'\n", selectedTemplate)

	if selectedTemplate == "" {
		// User quit the TUI
		fmt.Println("No template selected, exiting...")
		return
	}

	// After selecting a template, we get the rest of the details
	projectData, err := promptForProjectDetails()
	if err != nil {
		log.Fatalf("❌ Could not get project details: %v", err)
	}

	// Perform scaffolding
	scaffoldAndApply(selectedTemplate, projectData)
}

// runSimpleInteractiveScaffold handles the original survey flow
func runSimpleInteractiveScaffold() {
	// For this simple version, we hardcode the first template.
	// A better version would ask the user to choose here as well.
	const defaultTemplate = "nextjs-golden-path"
	fmt.Printf("🚀 Welcome! Using default template: %s\n", defaultTemplate)

	projectData, err := promptForProjectDetails()
	if err != nil {
		log.Fatalf("❌ Could not get project details: %v", err)
	}

	scaffoldAndApply(defaultTemplate, projectData)
}

// scaffoldAndApply contains the shared logic for scaffolding and templating
func scaffoldAndApply(templateName string, data *TemplateData) {
	sourceDir := filepath.Join("templates", templateName)
	destDir := data.ProjectName

	fmt.Printf("📂 Scaffolding project in './%s'...\n", destDir)
	err := scaffoldProject(sourceDir, destDir)
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

// promptForProjectDetails uses the survey library for interactive mode.
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
		if errors.Is(err, terminal.InterruptErr) { // Handle Ctrl+C
			fmt.Println("\nOperation cancelled.")
			os.Exit(0)
		}
		return nil, err
	}
	return data, nil
}

// scaffoldProject copies the entire directory structure from a source to a destination.
func scaffoldProject(sourceDir, destDir string) error {
	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", destDir)
	}
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		outPath := filepath.Join(destDir, path[len(sourceDir):])
		if info.IsDir() {
			return os.MkdirAll(outPath, info.Mode())
		}
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(outPath, fileBytes, info.Mode())
	})
}

// applyTemplates walks through the newly created project directory and processes templates.
func applyTemplates(destDir string, data *TemplateData) error {
	return filepath.Walk(destDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		tmpl, err := template.New(info.Name()).Parse(string(content))
		if err != nil {
			fmt.Printf("⚠️  Could not parse template for %s: %v. Skipping.\n", path, err)
			return nil
		}
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		return tmpl.Execute(file, data)
	})
}
