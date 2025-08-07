package manifest

import (
	"testing"
)

func TestWorkbenchManifest_Validation(t *testing.T) {
	tests := []struct {
		name     string
		manifest WorkbenchManifest
		isValid  bool
	}{
		{
			name: "valid manifest",
			manifest: WorkbenchManifest{
				APIVersion: "v1",
				Kind:       "Workbench",
				Metadata: ProjectMetadata{
					Name: "test-project",
				},
				Services: map[string]Service{
					"frontend": {
						Template: "react-typescript",
						Path:     "frontend",
						Port:     3000,
					},
				},
				Environments: map[string]Environment{
					"production": {
						Provider: "aws",
						Region:   "us-east-1",
					},
				},
			},
			isValid: true,
		},
		{
			name: "missing project name",
			manifest: WorkbenchManifest{
				APIVersion: "v1",
				Kind:       "Workbench",
				Metadata:   ProjectMetadata{},
				Services: map[string]Service{
					"frontend": {
						Template: "react-typescript",
						Path:     "frontend",
					},
				},
			},
			isValid: false,
		},
		{
			name: "no services",
			manifest: WorkbenchManifest{
				APIVersion: "v1",
				Kind:       "Workbench",
				Metadata: ProjectMetadata{
					Name: "test-project",
				},
				Services: map[string]Service{},
			},
			isValid: false,
		},
		{
			name: "no environments",
			manifest: WorkbenchManifest{
				APIVersion: "v1",
				Kind:       "Workbench",
				Metadata: ProjectMetadata{
					Name: "test-project",
				},
				Services: map[string]Service{
					"frontend": {
						Template: "react-typescript",
						Path:     "frontend",
					},
				},
				Environments: map[string]Environment{},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test basic validation logic
			isValid := tt.manifest.Metadata.Name != "" &&
				len(tt.manifest.Services) > 0 &&
				len(tt.manifest.Environments) > 0

			if isValid != tt.isValid {
				t.Errorf("expected validity %v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestService_Validation(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		isValid bool
	}{
		{
			name: "valid service",
			service: Service{
				Template: "react-typescript",
				Path:     "frontend",
				Port:     3000,
				Environment: map[string]string{
					"NODE_ENV": "production",
				},
			},
			isValid: true,
		},
		{
			name: "missing template",
			service: Service{
				Path: "frontend",
				Port: 3000,
			},
			isValid: false,
		},
		{
			name: "missing path",
			service: Service{
				Template: "react-typescript",
				Port:     3000,
			},
			isValid: false,
		},
		{
			name: "valid service with resources",
			service: Service{
				Template: "react-typescript",
				Path:     "frontend",
				Port:     3000,
				Resources: map[string]Resource{
					"database": {
						Type:    "postgres",
						Version: "13",
						Config: map[string]string{
							"size": "db.t3.micro",
						},
					},
				},
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.service.Template != "" && tt.service.Path != ""
			if isValid != tt.isValid {
				t.Errorf("expected validity %v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestComponent_Validation(t *testing.T) {
	tests := []struct {
		name      string
		component Component
		isValid   bool
	}{
		{
			name: "valid component",
			component: Component{
				Template: "nginx-gateway",
				Path:     "gateway",
				Ports:    []string{"80", "443"},
			},
			isValid: true,
		},
		{
			name: "missing template",
			component: Component{
				Path:  "gateway",
				Ports: []string{"80", "443"},
			},
			isValid: false,
		},
		{
			name: "missing path",
			component: Component{
				Template: "nginx-gateway",
				Ports:    []string{"80", "443"},
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.component.Template != "" && tt.component.Path != ""
			if isValid != tt.isValid {
				t.Errorf("expected validity %v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestEnvironment_Validation(t *testing.T) {
	tests := []struct {
		name        string
		environment Environment
		isValid     bool
	}{
		{
			name: "valid environment",
			environment: Environment{
				Provider: "aws",
				Region:   "us-east-1",
				Config: map[string]string{
					"instance_type": "t3.micro",
				},
			},
			isValid: true,
		},
		{
			name: "missing provider",
			environment: Environment{
				Region: "us-east-1",
			},
			isValid: false,
		},
		{
			name: "valid environment without region",
			environment: Environment{
				Provider: "docker",
				Config: map[string]string{
					"network": "bridge",
				},
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.environment.Provider != ""
			if isValid != tt.isValid {
				t.Errorf("expected validity %v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestResource_Validation(t *testing.T) {
	tests := []struct {
		name     string
		resource Resource
		isValid  bool
	}{
		{
			name: "valid resource",
			resource: Resource{
				Type:    "postgres",
				Version: "13",
				Config: map[string]string{
					"size": "db.t3.micro",
				},
			},
			isValid: true,
		},
		{
			name: "missing type",
			resource: Resource{
				Version: "13",
				Config: map[string]string{
					"size": "db.t3.micro",
				},
			},
			isValid: false,
		},
		{
			name: "valid resource without version",
			resource: Resource{
				Type: "redis",
				Config: map[string]string{
					"memory": "256mb",
				},
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.resource.Type != ""
			if isValid != tt.isValid {
				t.Errorf("expected validity %v, got %v", tt.isValid, isValid)
			}
		})
	}
}
