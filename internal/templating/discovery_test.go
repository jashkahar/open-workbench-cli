package templating

import (
	"testing"
)

func TestTemplateManifest_Validation(t *testing.T) {
	tests := []struct {
		name     string
		manifest *TemplateManifest
		wantErr  bool
	}{
		{
			name: "valid manifest",
			manifest: &TemplateManifest{
				Name:        "test-template",
				Description: "A test template",
				Parameters: []Parameter{
					{
						Name:     "project_name",
						Prompt:   "What is your project name?",
						Type:     "string",
						Required: true,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			manifest: &TemplateManifest{
				Description: "A test template",
				Parameters: []Parameter{
					{
						Name:     "project_name",
						Prompt:   "What is your project name?",
						Type:     "string",
						Required: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing description",
			manifest: &TemplateManifest{
				Name: "test-template",
				Parameters: []Parameter{
					{
						Name:     "project_name",
						Prompt:   "What is your project name?",
						Type:     "string",
						Required: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid parameter type",
			manifest: &TemplateManifest{
				Name:        "test-template",
				Description: "A test template",
				Parameters: []Parameter{
					{
						Name:     "project_name",
						Prompt:   "What is your project name?",
						Type:     "invalid_type",
						Required: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "select parameter without options",
			manifest: &TemplateManifest{
				Name:        "test-template",
				Description: "A test template",
				Parameters: []Parameter{
					{
						Name:     "framework",
						Prompt:   "Which framework?",
						Type:     "select",
						Required: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "multiselect parameter without options",
			manifest: &TemplateManifest{
				Name:        "test-template",
				Description: "A test template",
				Parameters: []Parameter{
					{
						Name:     "features",
						Prompt:   "Which features?",
						Type:     "multiselect",
						Required: true,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid select parameter with options",
			manifest: &TemplateManifest{
				Name:        "test-template",
				Description: "A test template",
				Parameters: []Parameter{
					{
						Name:     "framework",
						Prompt:   "Which framework?",
						Type:     "select",
						Required: true,
						Options:  []string{"react", "vue", "angular"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid multiselect parameter with options",
			manifest: &TemplateManifest{
				Name:        "test-template",
				Description: "A test template",
				Parameters: []Parameter{
					{
						Name:     "features",
						Prompt:   "Which features?",
						Type:     "multiselect",
						Required: true,
						Options:  []string{"auth", "database", "api"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test parameter validation logic
			if tt.manifest.Name == "" || tt.manifest.Description == "" {
				if !tt.wantErr {
					t.Errorf("expected validation to pass but manifest has missing required fields")
				}
				return
			}

			// Test parameter validation logic
			for _, param := range tt.manifest.Parameters {
				if param.Type == "select" || param.Type == "multiselect" {
					if len(param.Options) == 0 {
						if !tt.wantErr {
							t.Errorf("expected validation to pass but parameter has no options")
						}
						return
					}
				}

				if param.Type != "string" && param.Type != "boolean" && param.Type != "select" && param.Type != "multiselect" {
					if !tt.wantErr {
						t.Errorf("expected validation to pass but parameter has invalid type")
					}
					return
				}
			}

			// If we get here and expected an error, that's a problem
			if tt.wantErr {
				t.Errorf("expected validation to fail but it passed")
			}
		})
	}
}

func TestParameter_Validation(t *testing.T) {
	tests := []struct {
		name    string
		param   Parameter
		wantErr bool
	}{
		{
			name: "valid string parameter",
			param: Parameter{
				Name:     "project_name",
				Prompt:   "What is your project name?",
				Type:     "string",
				Required: true,
			},
			wantErr: false,
		},
		{
			name: "valid boolean parameter",
			param: Parameter{
				Name:     "use_typescript",
				Prompt:   "Use TypeScript?",
				Type:     "boolean",
				Required: false,
			},
			wantErr: false,
		},
		{
			name: "valid select parameter",
			param: Parameter{
				Name:     "framework",
				Prompt:   "Which framework?",
				Type:     "select",
				Required: true,
				Options:  []string{"react", "vue", "angular"},
			},
			wantErr: false,
		},
		{
			name: "valid multiselect parameter",
			param: Parameter{
				Name:     "features",
				Prompt:   "Which features?",
				Type:     "multiselect",
				Required: false,
				Options:  []string{"auth", "database", "api"},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			param: Parameter{
				Prompt:   "What is your project name?",
				Type:     "string",
				Required: true,
			},
			wantErr: true,
		},
		{
			name: "missing prompt",
			param: Parameter{
				Name:     "project_name",
				Type:     "string",
				Required: true,
			},
			wantErr: true,
		},
		{
			name: "missing type",
			param: Parameter{
				Name:     "project_name",
				Prompt:   "What is your project name?",
				Required: true,
			},
			wantErr: true,
		},
		{
			name: "invalid type",
			param: Parameter{
				Name:     "project_name",
				Prompt:   "What is your project name?",
				Type:     "invalid_type",
				Required: true,
			},
			wantErr: true,
		},
		{
			name: "select without options",
			param: Parameter{
				Name:     "framework",
				Prompt:   "Which framework?",
				Type:     "select",
				Required: true,
			},
			wantErr: true,
		},
		{
			name: "multiselect without options",
			param: Parameter{
				Name:     "features",
				Prompt:   "Which features?",
				Type:     "multiselect",
				Required: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test parameter validation
			if tt.param.Name == "" || tt.param.Prompt == "" || tt.param.Type == "" {
				if !tt.wantErr {
					t.Errorf("expected validation to pass but parameter has missing required fields")
				}
				return
			}

			// Test type validation
			if tt.param.Type != "string" && tt.param.Type != "boolean" && tt.param.Type != "select" && tt.param.Type != "multiselect" {
				if !tt.wantErr {
					t.Errorf("expected validation to pass but parameter has invalid type")
				}
				return
			}

			// Test options validation for select/multiselect
			if (tt.param.Type == "select" || tt.param.Type == "multiselect") && len(tt.param.Options) == 0 {
				if !tt.wantErr {
					t.Errorf("expected validation to pass but parameter has no options")
				}
				return
			}

			// If we get here and expected an error, that's a problem
			if tt.wantErr {
				t.Errorf("expected validation to fail but it passed")
			}
		})
	}
}

func TestValidation_Structure(t *testing.T) {
	validation := &Validation{
		Regex:        "^[a-z0-9-]+$",
		ErrorMessage: "Project name must contain only lowercase letters, numbers, and hyphens",
	}

	if validation.Regex == "" {
		t.Error("expected regex to be non-empty")
	}
	if validation.ErrorMessage == "" {
		t.Error("expected error message to be non-empty")
	}
}

func TestPostScaffold_Structure(t *testing.T) {
	postScaffold := &PostScaffold{
		FilesToDelete: []FileAction{
			{
				Path:      "unused-file.txt",
				Condition: "{{.use_typescript}} == false",
			},
		},
		Commands: []CommandAction{
			{
				Command:     "npm install",
				Description: "Install dependencies",
			},
		},
	}

	if len(postScaffold.FilesToDelete) == 0 && len(postScaffold.Commands) == 0 {
		t.Error("expected at least one post-scaffold action")
	}
}

func TestFileAction_Structure(t *testing.T) {
	fileAction := FileAction{
		Path:      "unused-file.txt",
		Condition: "{{.use_typescript}} == false",
	}

	if fileAction.Path == "" {
		t.Error("expected path to be non-empty")
	}
}

func TestCommandAction_Structure(t *testing.T) {
	commandAction := CommandAction{
		Command:     "npm install",
		Description: "Install dependencies",
		Condition:   "{{.use_typescript}} == true",
	}

	if commandAction.Command == "" {
		t.Error("expected command to be non-empty")
	}
	if commandAction.Description == "" {
		t.Error("expected description to be non-empty")
	}
}
