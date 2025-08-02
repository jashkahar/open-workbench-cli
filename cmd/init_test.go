package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidProjectName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid name", "myproject", true},
		{"valid name with hyphens", "my-awesome-project", true},
		{"valid name with numbers", "project123", true},
		{"empty name", "", false},
		{"starts with number", "123project", false},
		{"ends with hyphen", "project-", false},
		{"contains uppercase", "Project", false},
		{"contains special chars", "project@#$", false},
		{"contains spaces", "my project", false},
		{"contains underscores", "my_project", false},
		{"single letter", "a", true},
		{"single number", "1", false},
		{"single hyphen", "-", false},
		{"only hyphens", "---", false},
		{"mixed valid", "a1-b2-c3", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidProjectName(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for input: %s", tt.expected, result, tt.input)
			}
		})
	}
}

func TestCreateProjectDirectories(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-init")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	tests := []struct {
		name        string
		projectName string
		serviceName string
		expectError bool
	}{
		{"valid names", "my-project", "frontend", false},
		{"valid names with numbers", "project123", "service456", false},
		{"empty project name", "", "frontend", true},
		{"empty service name", "my-project", "", true},
		{"invalid project name", "../malicious", "frontend", true},
		{"invalid service name", "my-project", "../malicious", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createProjectDirectories(tt.projectName, tt.serviceName)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none for project: %s, service: %s", tt.projectName, tt.serviceName)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for project: %s, service: %s: %v", tt.projectName, tt.serviceName, err)
				}

				// Verify directories were created
				projectPath := filepath.Join(tempDir, tt.projectName)
				if _, err := os.Stat(projectPath); os.IsNotExist(err) {
					t.Errorf("project directory was not created: %s", projectPath)
				}

				servicePath := filepath.Join(projectPath, tt.serviceName)
				if _, err := os.Stat(servicePath); os.IsNotExist(err) {
					t.Errorf("service directory was not created: %s", servicePath)
				}
			}
		})
	}
}

func TestCreateWorkbenchManifest(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-manifest")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	tests := []struct {
		name         string
		projectName  string
		serviceName  string
		templateName string
		expectError  bool
	}{
		{"valid manifest", "my-project", "frontend", "react-typescript", true},
		{"valid manifest with numbers", "project123", "service456", "nextjs-full-stack", true},
		{"empty project name", "", "frontend", "react-typescript", false},
		{"empty service name", "my-project", "", "react-typescript", true},
		{"empty template name", "my-project", "frontend", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createWorkbenchManifest(tt.projectName, tt.serviceName, tt.templateName)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none for project: %s, service: %s, template: %s",
						tt.projectName, tt.serviceName, tt.templateName)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for project: %s, service: %s, template: %s: %v",
						tt.projectName, tt.serviceName, tt.templateName, err)
				}

				// Verify manifest file was created
				manifestPath := filepath.Join(tt.projectName, "workbench.yaml")
				if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
					t.Errorf("manifest file was not created: %s", manifestPath)
				}
			}
		})
	}
}

func TestCheckDirectorySafety(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-safety")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Test empty directory (should pass)
	err = checkDirectorySafety()
	if err != nil {
		t.Errorf("expected no error for empty directory, got: %v", err)
	}

	// Test directory with hidden file (should pass)
	hiddenFile := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(hiddenFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create hidden file: %v", err)
	}

	err = checkDirectorySafety()
	if err != nil {
		t.Errorf("expected no error for directory with hidden file, got: %v", err)
	}

	// Test directory with visible file (should fail)
	visibleFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(visibleFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create visible file: %v", err)
	}

	err = checkDirectorySafety()
	if err == nil {
		t.Error("expected error for directory with visible file")
	}

	// Clean up
	os.Remove(visibleFile)
}

func TestPromptForProjectName(t *testing.T) {
	// This test would require mocking the survey library
	// For now, we'll test the validation logic separately
	// The actual prompt testing would require integration tests
	t.Skip("Skipping prompt test - requires survey library mocking")
}

func TestPromptForFirstService(t *testing.T) {
	// This test would require mocking the survey library and templatesFS
	// For now, we'll test the validation logic separately
	// The actual prompt testing would require integration tests
	t.Skip("Skipping prompt test - requires survey library and templatesFS mocking")
}

func TestScaffoldService(t *testing.T) {
	// This test would require mocking the templating system
	// For now, we'll test the validation logic separately
	// The actual scaffolding testing would require integration tests
	t.Skip("Skipping scaffold test - requires templating system mocking")
}

func TestPrintSuccessMessage(t *testing.T) {
	// Test that the function doesn't panic
	// This is a simple smoke test
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printSuccessMessage panicked: %v", r)
		}
	}()

	printSuccessMessage("test-project", "frontend")
}

// Integration test helper functions
func TestIntegrationInitCommand(t *testing.T) {
	// This would be a full integration test
	// It would require setting up a mock environment
	// and testing the entire init command flow
	t.Skip("Skipping integration test - requires full environment setup")
}

// Benchmark tests
func BenchmarkValidateAndSanitizeName(b *testing.B) {
	validName := "my-awesome-project-123"
	for i := 0; i < b.N; i++ {
		_, err := ValidateAndSanitizeName(validName, nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkValidateAndSanitizePath(b *testing.B) {
	validPath := "my-awesome-project/subdirectory"
	for i := 0; i < b.N; i++ {
		_, err := ValidateAndSanitizePath(validPath, nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

func BenchmarkCheckForSuspiciousPatterns(b *testing.B) {
	normalInput := "my-normal-project-name"
	for i := 0; i < b.N; i++ {
		err := CheckForSuspiciousPatterns(normalInput)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
