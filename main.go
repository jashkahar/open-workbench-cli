package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
)

// TemplateData holds the values that will be injected into the templates.
type TemplateData struct {
	ProjectName string
	Owner       string
}

// main is the entry point of our CLI application.
func main() {
	fmt.Println("🚀 Welcome to the Open Workbench CLI!")
	fmt.Println("------------------------------------")

	// 1. Get project details from the user.
	data, err := promptForProjectDetails()
	if err != nil {
		log.Fatalf("❌ Could not get project details: %v", err)
	}

	// Define source (template) and destination (output) directories.
	// This is hardcoded for our Phase 1 goal.
	sourceDir := "templates/nextjs-golden-path"
	destDir := data.ProjectName

	// 2. Scaffold the project structure.
	fmt.Printf("📂 Scaffolding project in './%s'...\n", destDir)
	err = scaffoldProject(sourceDir, destDir)
	if err != nil {
		log.Fatalf("❌ Failed to scaffold project: %v", err)
	}

	// 3. Process the templates in the new project directory.
	fmt.Println("✏️  Applying templates...")
	err = applyTemplates(destDir, data)
	if err != nil {
		log.Fatalf("❌ Failed to apply templates: %v", err)
	}

	fmt.Println("------------------------------------")
	fmt.Printf("✅ Success! Your new project '%s' is ready.\n", data.ProjectName)
	fmt.Println("\nNext steps:")
	fmt.Printf("1. cd %s\n", data.ProjectName)
	fmt.Println("2. npm install")
	fmt.Println("3. npm run dev")
}

// promptForProjectDetails uses the survey library to create an interactive prompt
// for collecting the necessary project information from the user.
func promptForProjectDetails() (*TemplateData, error) {
	data := &TemplateData{}

	// Define the set of questions to ask the user.
	questions := []*survey.Question{
		{
			Name:     "ProjectName",
			Prompt:   &survey.Input{Message: "What is your project name?"},
			Validate: survey.Required,
		},
		{
			Name:   "Owner",
			Prompt: &survey.Input{Message: "Who is the owner of this project?"},
			Validate: survey.Required,
		},
	}

	// Ask the questions and store the answers in our struct.
	err := survey.Ask(questions, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// scaffoldProject copies the entire directory structure from a source to a destination.
// It preserves the file hierarchy.
func scaffoldProject(sourceDir, destDir string) error {
	// Check if the destination directory already exists to prevent overwriting.
	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' already exists", destDir)
	}

	// Walk through the source directory.
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create the corresponding path in the destination directory.
		outPath := filepath.Join(destDir, path[len(sourceDir):])

		if info.IsDir() {
			// If it's a directory, create it.
			return os.MkdirAll(outPath, info.Mode())
		}
		
		// If it's a file, copy its contents.
		// We copy first, then template. This is simpler than templating in memory.
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(outPath, fileBytes, info.Mode())
	})
}

// applyTemplates walks through the newly created project directory and processes
// any file with Go template syntax.
func applyTemplates(destDir string, data *TemplateData) error {
	return filepath.Walk(destDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// We only want to process regular files, not directories.
		if info.IsDir() {
			return nil
		}

		// Read the content of the file.
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Create a new template and parse the file content.
		// Using the file's name for the template is a good practice.
		tmpl, err := template.New(info.Name()).Parse(string(content))
		if err != nil {
			// This error might occur if the file has malformed template syntax.
			// For now, we'll log it and continue, as some files might not be templates.
			// A more robust solution might check for `{{` before trying to parse.
			fmt.Printf("⚠️  Could not parse template for %s: %v. Skipping.\n", path, err)
			return nil
		}

		// Create a new file to write the processed template to.
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Execute the template with the user's data and write the output to the file.
		return tmpl.Execute(file, data)
	})
}