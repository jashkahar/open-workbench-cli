package compose

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	// Create a test workbench.yaml content
	workbenchContent := `apiVersion: openworkbench.io/v1alpha1
kind: Project
metadata:
  name: "test-project"

components:
  gateway:
    template: nginx-gateway
    path: ./components/gateway
    ports: ["8080:80"]

services:
  frontend:
    template: nextjs-golden-path
    path: ./frontend
    port: 3000
    environment:
      NEXT_PUBLIC_API_URL: "http://${components.gateway.name}:${components.gateway.port}/api"

  backend:
    template: fastapi-basic
    path: ./backend
    port: 8000
    resources:
      database:
        type: postgres
        version: "15"
    environment:
      DATABASE_URL: "postgres://${services.backend.resources.database.user}:${services.backend.resources.database.password}@${services.backend.resources.database.name}:5432/${services.backend.resources.database.dbname}"
`

	// Write test workbench.yaml
	err := os.WriteFile("test-workbench.yaml", []byte(workbenchContent), 0644)
	require.NoError(t, err)
	defer os.Remove("test-workbench.yaml")

	// Load the project
	project, err := LoadWorkbenchProject("test-workbench.yaml")
	require.NoError(t, err)
	assert.Equal(t, "test-project", project.Metadata.Name)

	// Create generator
	generator := NewGenerator(project)

	// Generate docker-compose config
	config, err := generator.Generate()
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Verify basic structure
	assert.Equal(t, "3.8", config.Version)
	assert.NotEmpty(t, config.Services)
	assert.NotEmpty(t, config.Networks)

	// Verify services exist
	assert.Contains(t, config.Services, "gateway")
	assert.Contains(t, config.Services, "frontend")
	assert.Contains(t, config.Services, "backend")
	assert.Contains(t, config.Services, "backend-database")

	// Verify network configuration
	assert.Contains(t, config.Services["frontend"].Networks, "workbench_net")
	assert.Contains(t, config.Services["backend"].Networks, "workbench_net")

	// Verify basic service configuration
	frontendService := config.Services["frontend"]
	assert.NotNil(t, frontendService.Build)
	assert.Equal(t, "./frontend", frontendService.Build.Context)

	backendService := config.Services["backend"]
	assert.NotNil(t, backendService.Build)
	assert.Equal(t, "./backend", backendService.Build.Context)
}

func TestGenerator_GenerateEnvFile(t *testing.T) {
	// Create a test project with resources
	project := &WorkbenchProject{
		APIVersion: "openworkbench.io/v1alpha1",
		Kind:       "Project",
		Metadata: ProjectMetadata{
			Name: "test-project",
		},
		Services: map[string]Service{
			"backend": {
				Template: "fastapi-basic",
				Path:     "./backend",
				Resources: map[string]Resource{
					"database": {
						Type:    "postgres",
						Version: "15",
					},
				},
			},
			"frontend": {
				Template: "nextjs-golden-path",
				Path:     "./frontend",
				Resources: map[string]Resource{
					"cache": {
						Type:    "redis",
						Version: "7",
					},
				},
			},
		},
	}

	generator := NewGenerator(project)

	// Generate environment variables
	envVars, err := generator.GenerateEnvFile()
	require.NoError(t, err)
	assert.NotEmpty(t, envVars)

	// Verify backend database credentials
	assert.Equal(t, "backend_user", envVars["backend_database_user"])
	assert.Equal(t, "password123", envVars["backend_database_password"])
	assert.Equal(t, "backend_database", envVars["backend_database_name"])
	assert.Equal(t, "backend_database_db", envVars["backend_database_dbname"])

	// Verify frontend cache credentials
	assert.Equal(t, "password123", envVars["frontend_cache_password"])
}

func TestPrerequisiteChecker_CheckDockerCompose(t *testing.T) {
	checker := NewPrerequisiteChecker()

	// This test will pass if docker is available, fail gracefully if not
	err := checker.CheckDockerCompose()
	if err != nil {
		t.Logf("Docker Compose not available (expected in CI): %v", err)
	}

	// Verify command detection works
	cmd := checker.GetDockerComposeCommand()
	assert.NotEmpty(t, cmd)
	assert.Contains(t, cmd, "docker")
}

func TestSaveAndLoadFiles(t *testing.T) {
	// Test saving docker-compose.yml
	config := &DockerComposeConfig{
		Version: "3.8",
		Services: map[string]DockerComposeService{
			"test": {
				Build: &BuildConfig{
					Context: "./test",
				},
				Networks: []string{"workbench_net"},
			},
		},
		Networks: map[string]interface{}{
			"workbench_net": map[string]string{
				"driver": "bridge",
			},
		},
	}

	err := SaveDockerCompose(config, "test-docker-compose.yml")
	require.NoError(t, err)
	defer os.Remove("test-docker-compose.yml")

	// Test saving .env file
	envVars := map[string]string{
		"TEST_VAR":    "test_value",
		"ANOTHER_VAR": "another_value",
	}

	err = SaveEnvFile(envVars, "test.env")
	require.NoError(t, err)
	defer os.Remove("test.env")

	// Test saving .env.example file
	err = SaveEnvExampleFile(envVars, "test.env.example")
	require.NoError(t, err)
	defer os.Remove("test.env.example")

	// Verify files were created
	_, err = os.Stat("test-docker-compose.yml")
	assert.NoError(t, err)

	_, err = os.Stat("test.env")
	assert.NoError(t, err)

	_, err = os.Stat("test.env.example")
	assert.NoError(t, err)
}
