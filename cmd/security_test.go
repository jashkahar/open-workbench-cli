package cmd

import (
	"os"
	"testing"
)

func TestValidateAndSanitizePath(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    string
	}{
		{"valid path", "my-project", false, "my-project"},
		{"valid path with hyphens", "my-awesome-project", false, "my-awesome-project"},
		{"empty path", "", true, ""},
		{"whitespace only", "   ", true, ""},
		{"path traversal attempt", "../etc/passwd", true, ""},
		{"path traversal attempt 2", "..\\windows\\system32", true, ""},
		{"absolute path", "/home/user", false, "\\home\\user"},
		{"absolute path windows", "C:\\Users\\user", true, ""},
		{"too long path", string(make([]byte, 300)), true, ""},
		{"windows reserved name", "con", true, ""},
		{"windows reserved name 2", "prn", true, ""},
		{"windows reserved name 3", "aux", true, ""},
		{"windows reserved name 4", "nul", true, ""},
		{"windows reserved name 5", "com1", true, ""},
		{"windows reserved name 6", "lpt1", true, ""},
		{"suspicious characters", "project<script>", false, "project<script>"},
		{"suspicious characters 2", "project<script>alert(1)</script>", false, "project<script>alert(1)<\\script>"},
		{"normal path with dots", "my.project", false, "my.project"},
		{"path with numbers", "project123", false, "project123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateAndSanitizePath(tt.input, nil)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none for input: %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %s: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("expected %s, got %s for input: %s", tt.expected, result, tt.input)
				}
			}
		})
	}
}

func TestValidateAndSanitizeName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    string
	}{
		{"valid name", "myproject", false, "myproject"},
		{"valid name with hyphens", "my-awesome-project", false, "my-awesome-project"},
		{"valid name with numbers", "project123", false, "project123"},
		{"empty name", "", true, ""},
		{"whitespace only", "   ", true, ""},
		{"too long name", string(make([]byte, 60)), true, ""},
		{"uppercase letters", "MyProject", true, ""},
		{"mixed case", "My-Awesome-Project", true, ""},
		{"starts with number", "123project", true, ""},
		{"ends with hyphen", "project-", true, ""},
		{"contains uppercase", "Project", true, ""},
		{"contains special chars", "project@#$", true, ""},
		{"contains spaces", "my project", true, ""},
		{"contains underscores", "my_project", true, ""},
		{"windows reserved name", "con", true, ""},
		{"windows reserved name 2", "prn", true, ""},
		{"windows reserved name 3", "aux", true, ""},
		{"windows reserved name 4", "nul", true, ""},
		{"windows reserved name 5", "com1", true, ""},
		{"windows reserved name 6", "lpt1", true, ""},
		{"suspicious patterns", "javascript:alert(1)", true, ""},
		{"suspicious patterns 2", "data:text/html,<script>alert(1)</script>", true, ""},
		{"suspicious patterns 3", "eval(alert(1))", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateAndSanitizeName(tt.input, nil)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none for input: %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %s: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("expected %s, got %s for input: %s", tt.expected, result, tt.input)
				}
			}
		})
	}
}

func TestValidateDirectorySafety(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-dir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test valid directory
	err = ValidateDirectorySafety(tempDir)
	if err != nil {
		t.Errorf("expected no error for valid directory, got: %v", err)
	}

	// Test non-existent directory
	err = ValidateDirectorySafety("/non/existent/directory")
	if err == nil {
		t.Error("expected error for non-existent directory")
	}

	// Test file instead of directory
	tempFile, err := os.CreateTemp("", "test-file")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	err = ValidateDirectorySafety(tempFile.Name())
	if err == nil {
		t.Error("expected error for file instead of directory")
	}
}

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"valid template name", "nextjs-full-stack", false},
		{"valid template name 2", "react-typescript", false},
		{"valid template name 3", "express-api", false},
		{"empty name", "", true},
		{"whitespace only", "   ", true},
		{"too long name", string(make([]byte, 110)), true},
		{"contains path traversal", "../templates/malicious", true},
		{"contains forward slash", "template/name", true},
		{"contains backslash", "template\\name", true},
		{"contains special chars", "template@#$", true},
		{"contains spaces", "template name", true},
		{"suspicious patterns", "javascript:alert(1)", true},
		{"suspicious patterns 2", "data:text/html,<script>alert(1)</script>", true},
		{"suspicious patterns 3", "eval(alert(1))", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplateName(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none for input: %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %s: %v", tt.input, err)
				}
			}
		})
	}
}

func TestCheckForSuspiciousPatterns(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{"normal input", "my-project", false},
		{"javascript pattern", "javascript:alert(1)", true},
		{"data pattern", "data:text/html,<script>alert(1)</script>", true},
		{"vbscript pattern", "vbscript:msgbox('hello')", true},
		{"onload pattern", "onload=alert(1)", true},
		{"onerror pattern", "onerror=alert(1)", true},
		{"script tag", "<script>alert(1)</script>", true},
		{"eval pattern", "eval(alert(1))", true},
		{"exec pattern", "exec(system('rm -rf /'))", true},
		{"system pattern", "system('rm -rf /')", true},
		{"mixed case", "JavaScript:alert(1)", true},
		{"mixed case 2", "OnLoad=alert(1)", true},
		{"normal text with script", "This is a normal text with <script>alert(1)</script>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckForSuspiciousPatterns(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none for input: %s", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %s: %v", tt.input, err)
				}
			}
		})
	}
}

func TestSanitizeForFileSystem(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal input", "my-project", "my-project"},
		{"with angle brackets", "project<script>", "project"},
		{"with colons", "project:name", "projectname"},
		{"with quotes", "project\"name", "projectname"},
		{"with pipes", "project|name", "projectname"},
		{"with question marks", "project?name", "projectname"},
		{"with asterisks", "project*name", "projectname"},
		{"with backslashes", "project\\name", "projectname"},
		{"with forward slashes", "project/name", "projectname"},
		{"mixed special chars", "project<script>alert(1)</script>", "project(1)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeForFileSystem(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s for input: %s", tt.expected, result, tt.input)
			}
		})
	}
}

func TestDefaultSecurityConfig(t *testing.T) {
	config := DefaultSecurityConfig()

	if config.MaxPathLength != 255 {
		t.Errorf("expected MaxPathLength to be 255, got %d", config.MaxPathLength)
	}

	if config.MaxNameLength != 50 {
		t.Errorf("expected MaxNameLength to be 50, got %d", config.MaxNameLength)
	}

	if config.AllowedCharacters == nil {
		t.Error("expected AllowedCharacters to be set")
	}

	if len(config.ForbiddenPatterns) == 0 {
		t.Error("expected ForbiddenPatterns to be set")
	}
}

func TestSecurityConfigValidation(t *testing.T) {
	config := DefaultSecurityConfig()

	// Test path validation with custom config
	validPath := "my-project"
	result, err := ValidateAndSanitizePath(validPath, config)
	if err != nil {
		t.Errorf("unexpected error for valid path: %v", err)
	}
	if result != validPath {
		t.Errorf("expected %s, got %s", validPath, result)
	}

	// Test name validation with custom config
	validName := "myproject"
	result, err = ValidateAndSanitizeName(validName, config)
	if err != nil {
		t.Errorf("unexpected error for valid name: %v", err)
	}
	if result != validName {
		t.Errorf("expected %s, got %s", validName, result)
	}
}
