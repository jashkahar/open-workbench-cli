package terraform

import (
	"os"
	"path/filepath"
	"testing"

	manifestPkg "github.com/jashkahar/open-workbench-platform/internal/manifest"
)

func TestNewGenerator(t *testing.T) {
	generator := NewGenerator()
	if generator == nil {
		t.Fatal("NewGenerator() returned nil")
	}
}

func TestGenerator_Name(t *testing.T) {
	generator := NewGenerator()
	name := generator.Name()
	if name != "terraform" {
		t.Errorf("expected name 'terraform', got '%s'", name)
	}
}

func TestGenerator_Description(t *testing.T) {
	generator := NewGenerator()
	description := generator.Description()
	if description == "" {
		t.Error("expected non-empty description")
	}
}

func TestGenerator_Validate(t *testing.T) {
	generator := NewGenerator()

	tests := []struct {
		name     string
		manifest *manifestPkg.WorkbenchManifest
		wantErr  bool
	}{
		{
			name: "valid manifest",
			manifest: &manifestPkg.WorkbenchManifest{
				Metadata: manifestPkg.ProjectMetadata{
					Name: "test-project",
				},
				Services: map[string]manifestPkg.Service{
					"frontend": {
						Template: "react-typescript",
						Path:     "frontend",
						Port:     3000,
					},
				},
				Environments: map[string]manifestPkg.Environment{
					"production": {
						Provider: "aws",
						Region:   "us-east-1",
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "nil manifest",
			manifest: nil,
			wantErr:  true,
		},
		{
			name: "missing project name",
			manifest: &manifestPkg.WorkbenchManifest{
				Metadata: manifestPkg.ProjectMetadata{},
				Services: map[string]manifestPkg.Service{
					"frontend": {
						Template: "react-typescript",
						Path:     "frontend",
					},
				},
				Environments: map[string]manifestPkg.Environment{
					"production": {
						Provider: "aws",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no services",
			manifest: &manifestPkg.WorkbenchManifest{
				Metadata: manifestPkg.ProjectMetadata{
					Name: "test-project",
				},
				Services: map[string]manifestPkg.Service{},
				Environments: map[string]manifestPkg.Environment{
					"production": {
						Provider: "aws",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no environments",
			manifest: &manifestPkg.WorkbenchManifest{
				Metadata: manifestPkg.ProjectMetadata{
					Name: "test-project",
				},
				Services: map[string]manifestPkg.Service{
					"frontend": {
						Template: "react-typescript",
						Path:     "frontend",
					},
				},
				Environments: map[string]manifestPkg.Environment{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generator.Validate(tt.manifest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_Generate(t *testing.T) {
	generator := NewGenerator()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "terraform-test")
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

	manifest := &manifestPkg.WorkbenchManifest{
		Metadata: manifestPkg.ProjectMetadata{
			Name: "test-project",
		},
		Services: map[string]manifestPkg.Service{
			"frontend": {
				Template: "react-typescript",
				Path:     "frontend",
				Port:     3000,
			},
			"backend": {
				Template: "express-api",
				Path:     "backend",
				Port:     8080,
			},
		},
		Components: map[string]manifestPkg.Component{
			"gateway": {
				Template: "nginx-gateway",
				Path:     "gateway",
				Ports:    []string{"80", "443"},
			},
		},
		Environments: map[string]manifestPkg.Environment{
			"production": {
				Provider: "aws",
				Region:   "us-east-1",
				Config: map[string]string{
					"instance_type": "t3.micro",
				},
			},
		},
	}

	err = generator.Generate(manifest)
	if err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}

	// Verify that terraform directory was created
	terraformDir := "terraform"
	if _, err := os.Stat(terraformDir); os.IsNotExist(err) {
		t.Fatal("terraform directory was not created")
	}

	// Verify that required files were generated
	requiredFiles := []string{"main.tf", "variables.tf", "outputs.tf", "terraform.tfvars.example"}
	for _, file := range requiredFiles {
		filePath := filepath.Join(terraformDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("required file %s was not generated", file)
		}
	}
}

func TestGenerator_generateServiceResources(t *testing.T) {
	generator := NewGenerator()

	service := manifestPkg.Service{
		Template: "react-typescript",
		Path:     "frontend",
		Port:     3000,
		Environment: map[string]string{
			"NODE_ENV": "production",
		},
	}

	content := generator.generateServiceResources("frontend", service)

	// Verify that the generated content contains expected elements
	expectedElements := []string{
		"# Service: frontend",
		"resource \"aws_ecs_service\" \"frontend\"",
		"resource \"aws_ecs_task_definition\" \"frontend\"",
		"resource \"aws_lb_target_group\" \"frontend\"",
		"container_port   = 3000",
		"containerPort = 3000",
		"port     = 3000",
	}

	for _, element := range expectedElements {
		if !contains(content, element) {
			t.Errorf("generated content missing expected element: %s", element)
		}
	}
}

func TestGenerator_generateComponentResources(t *testing.T) {
	generator := NewGenerator()

	component := manifestPkg.Component{
		Template: "nginx-gateway",
		Path:     "gateway",
		Ports:    []string{"80", "443"},
	}

	content := generator.generateComponentResources("gateway", component)

	// Verify that the generated content contains expected elements
	expectedElements := []string{
		"# Component: gateway",
		"resource \"aws_ecs_service\" \"gateway\"",
		"resource \"aws_ecs_task_definition\" \"gateway\"",
	}

	for _, element := range expectedElements {
		if !contains(content, element) {
			t.Errorf("generated content missing expected element: %s", element)
		}
	}
}

func TestGenerator_generateVariablesTf(t *testing.T) {
	generator := NewGenerator()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "terraform-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manifest := &manifestPkg.WorkbenchManifest{
		Metadata: manifestPkg.ProjectMetadata{
			Name: "test-project",
		},
		Services: map[string]manifestPkg.Service{
			"frontend": {
				Template: "react-typescript",
				Path:     "frontend",
				Port:     3000,
			},
		},
		Components: map[string]manifestPkg.Component{
			"gateway": {
				Template: "nginx-gateway",
				Path:     "gateway",
			},
		},
	}

	// Create servicesForEnv map for testing
	servicesForEnv := map[string]manifestPkg.Service{
		"frontend": manifest.Services["frontend"],
	}

	err = generator.generateVariablesTf(manifest, tempDir, servicesForEnv)
	if err != nil {
		t.Fatalf("generateVariablesTf() failed: %v", err)
	}

	// Verify that variables.tf was created
	variablesFile := filepath.Join(tempDir, "variables.tf")
	if _, err := os.Stat(variablesFile); os.IsNotExist(err) {
		t.Fatal("variables.tf was not generated")
	}

	// Read and verify content
	content, err := os.ReadFile(variablesFile)
	if err != nil {
		t.Fatalf("failed to read variables.tf: %v", err)
	}

	contentStr := string(content)
	expectedElements := []string{
		"variable \"aws_region\"",
		"variable \"project_name\"",
		"variable \"frontend_desired_count\"",
		"variable \"frontend_cpu\"",
		"variable \"frontend_memory\"",
		"variable \"frontend_image\"",
		"variable \"gateway_desired_count\"",
		"variable \"gateway_cpu\"",
		"variable \"gateway_memory\"",
		"variable \"gateway_image\"",
	}

	for _, element := range expectedElements {
		if !contains(contentStr, element) {
			t.Errorf("variables.tf missing expected element: %s", element)
		}
	}
}

func TestGenerator_generateOutputsTf(t *testing.T) {
	generator := NewGenerator()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "terraform-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manifest := &manifestPkg.WorkbenchManifest{
		Metadata: manifestPkg.ProjectMetadata{
			Name: "test-project",
		},
		Services: map[string]manifestPkg.Service{
			"frontend": {
				Template: "react-typescript",
				Path:     "frontend",
				Port:     3000,
			},
		},
	}

	// Create servicesForEnv map for testing
	servicesForEnv := map[string]manifestPkg.Service{
		"frontend": manifest.Services["frontend"],
	}

	err = generator.generateOutputsTf(manifest, tempDir, servicesForEnv)
	if err != nil {
		t.Fatalf("generateOutputsTf() failed: %v", err)
	}

	// Verify that outputs.tf was created
	outputsFile := filepath.Join(tempDir, "outputs.tf")
	if _, err := os.Stat(outputsFile); os.IsNotExist(err) {
		t.Fatal("outputs.tf was not generated")
	}

	// Read and verify content
	content, err := os.ReadFile(outputsFile)
	if err != nil {
		t.Fatalf("failed to read outputs.tf: %v", err)
	}

	contentStr := string(content)
	expectedElements := []string{
		"output \"vpc_id\"",
		"output \"alb_dns_name\"",
		"output \"ecs_cluster_name\"",
		"output \"frontend_service_name\"",
		"output \"frontend_task_definition_arn\"",
	}

	for _, element := range expectedElements {
		if !contains(contentStr, element) {
			t.Errorf("outputs.tf missing expected element: %s", element)
		}
	}
}

func TestGenerator_generateTfvarsExample(t *testing.T) {
	generator := NewGenerator()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "terraform-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manifest := &manifestPkg.WorkbenchManifest{
		Metadata: manifestPkg.ProjectMetadata{
			Name: "test-project",
		},
		Services: map[string]manifestPkg.Service{
			"frontend": {
				Template: "react-typescript",
				Path:     "frontend",
				Port:     3000,
			},
		},
		Components: map[string]manifestPkg.Component{
			"gateway": {
				Template: "nginx-gateway",
				Path:     "gateway",
			},
		},
	}

	// Create servicesForEnv map for testing
	servicesForEnv := map[string]manifestPkg.Service{
		"frontend": manifest.Services["frontend"],
	}

	err = generator.generateTfvarsExample(manifest, tempDir, servicesForEnv)
	if err != nil {
		t.Fatalf("generateTfvarsExample() failed: %v", err)
	}

	// Verify that terraform.tfvars.example was created
	tfvarsFile := filepath.Join(tempDir, "terraform.tfvars.example")
	if _, err := os.Stat(tfvarsFile); os.IsNotExist(err) {
		t.Fatal("terraform.tfvars.example was not generated")
	}

	// Read and verify content
	content, err := os.ReadFile(tfvarsFile)
	if err != nil {
		t.Fatalf("failed to read terraform.tfvars.example: %v", err)
	}

	contentStr := string(content)
	expectedElements := []string{
		"aws_region = \"us-east-1\"",
		"project_name = \"test-project\"",
		"# frontend service configuration",
		"frontend_desired_count = 1",
		"frontend_cpu = 256",
		"frontend_memory = 512",
		"frontend_image = \"nginx:alpine\"",
		"# gateway component configuration",
		"gateway_desired_count = 1",
		"gateway_cpu = 256",
		"gateway_memory = 512",
		"gateway_image = \"nginx:alpine\"",
	}

	for _, element := range expectedElements {
		if !contains(contentStr, element) {
			t.Errorf("terraform.tfvars.example missing expected element: %s", element)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
