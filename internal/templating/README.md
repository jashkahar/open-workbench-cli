# Internal Templating Package

This package contains the core templating system for the Open Workbench CLI. It provides dynamic template discovery, parameter processing, and file generation capabilities.

## 📁 Package Structure

```
internal/templating/
├── discovery.go      # Template discovery and validation
├── parameters.go     # Parameter collection and validation
├── processor.go      # Template processing and file operations
└── README.md        # This file
```

## 🎯 Overview

The templating package implements a sophisticated dynamic template system that supports:

- **Template Discovery**: Automatic discovery of available templates
- **Parameter Processing**: Dynamic parameter collection with validation
- **Conditional Logic**: Show/hide parameters and files based on conditions
- **Template Processing**: File generation with variable substitution
- **Post-Scaffolding Actions**: Automatic cleanup and setup

## 🔧 Core Components

### Discovery (`discovery.go`)

Handles template discovery and validation.

#### Key Functions

- `DiscoverTemplates(fs.FS) ([]TemplateInfo, error)`: Find all available templates
- `LoadTemplateManifest(fs.FS, string) (*TemplateManifest, error)`: Load template manifest
- `ValidateTemplate(fs.FS, string) error`: Validate template structure
- `GetTemplateInfo(fs.FS, string) (*TemplateInfo, error)`: Get specific template info

#### Usage Example

```go
// Discover all templates
templates, err := templating.DiscoverTemplates(templatesFS)
if err != nil {
    log.Fatalf("Failed to discover templates: %v", err)
}

// Load specific template
manifest, err := templating.LoadTemplateManifest(templatesFS, "nextjs-full-stack")
if err != nil {
    log.Fatalf("Failed to load template: %v", err)
}
```

### Parameters (`parameters.go`)

Handles parameter collection, validation, and processing.

#### Key Components

- `ParameterProcessor`: Main parameter processing logic
- `Parameter`: Parameter definition structure
- `Validation`: Validation rules structure

#### Key Functions

- `NewParameterProcessor(*TemplateManifest) *ParameterProcessor`: Create processor
- `GetVisibleParameters() []Parameter`: Get parameters based on conditions
- `ValidateParameter(Parameter, interface{}) error`: Validate parameter values
- `GetParameterGroups() map[string][]Parameter`: Organize parameters by groups
- `evaluateCondition(string) (bool, error)`: Evaluate conditional logic

#### Usage Example

```go
// Create parameter processor
processor := templating.NewParameterProcessor(manifest)

// Get visible parameters
visibleParams := processor.GetVisibleParameters()

// Validate parameter
err := processor.ValidateParameter(param, value)
if err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

// Get parameter groups
groups := processor.GetParameterGroups()
for groupName, params := range groups {
    fmt.Printf("Group: %s\n", groupName)
    for _, param := range params {
        fmt.Printf("  - %s\n", param.Name)
    }
}
```

### Processor (`processor.go`)

Handles template processing and file operations.

#### Key Components

- `TemplateProcessor`: Main template processing logic
- `PostScaffold`: Post-scaffolding actions definition

#### Key Functions

- `NewTemplateProcessor(*TemplateManifest, map[string]interface{}) *TemplateProcessor`: Create processor
- `ScaffoldProject(fs.FS, string, string) error`: Main scaffolding function
- `ProcessTemplate(string) (string, error)`: Process template content
- `ProcessFileName(string) (string, error)`: Process filename template
- `ExecutePostScaffoldActions(string) error`: Execute post-scaffolding actions

#### Usage Example

```go
// Create template processor
processor := templating.NewTemplateProcessor(manifest, parameterValues)

// Scaffold project
err := processor.ScaffoldProject(templatesFS, templateName, destDir)
if err != nil {
    log.Fatalf("Failed to scaffold project: %v", err)
}

// Execute post-scaffolding actions
err = processor.ExecutePostScaffoldActions(destDir)
if err != nil {
    log.Fatalf("Failed to execute post-scaffolding actions: %v", err)
}
```

## 📋 Data Structures

### TemplateManifest

Represents a template's configuration:

```go
type TemplateManifest struct {
    Name         string        `json:"name"`
    Description  string        `json:"description"`
    Parameters   []Parameter   `json:"parameters"`
    PostScaffold *PostScaffold `json:"postScaffold,omitempty"`
}
```

### Parameter

Represents a single parameter:

```go
type Parameter struct {
    Name       string      `json:"name"`
    Prompt     string      `json:"prompt"`
    HelpText   string      `json:"helpText,omitempty"`
    Group      string      `json:"group,omitempty"`
    Type       string      `json:"type"`
    Required   bool        `json:"required,omitempty"`
    Default    any         `json:"default,omitempty"`
    Options    []string    `json:"options,omitempty"`
    Condition  string      `json:"condition,omitempty"`
    Validation *Validation `json:"validation,omitempty"`
}
```

### Validation

Represents validation rules:

```go
type Validation struct {
    Regex        string `json:"regex"`
    ErrorMessage string `json:"errorMessage"`
}
```

### PostScaffold

Represents post-scaffolding actions:

```go
type PostScaffold struct {
    FilesToDelete []FileAction    `json:"filesToDelete,omitempty"`
    Commands      []CommandAction `json:"commands,omitempty"`
}
```

## 🔄 Workflow

### 1. Template Discovery

```go
// Discover available templates
templates, err := templating.DiscoverTemplates(templatesFS)
if err != nil {
    return err
}

// Validate each template
for _, template := range templates {
    err := templating.ValidateTemplate(templatesFS, template.Name)
    if err != nil {
        log.Printf("Warning: Template %s is invalid: %v", template.Name, err)
    }
}
```

### 2. Parameter Collection

```go
// Load template manifest
manifest, err := templating.LoadTemplateManifest(templatesFS, templateName)
if err != nil {
    return err
}

// Create parameter processor
processor := templating.NewParameterProcessor(manifest)

// Collect parameters
values := make(map[string]interface{})
for _, param := range processor.GetVisibleParameters() {
    value, err := promptForParameter(param)
    if err != nil {
        return err
    }

    // Validate parameter
    err = processor.ValidateParameter(param, value)
    if err != nil {
        return err
    }

    values[param.Name] = value
    processor.SetValue(param.Name, value)
}
```

### 3. Template Processing

```go
// Create template processor
processor := templating.NewTemplateProcessor(manifest, values)

// Scaffold project
err := processor.ScaffoldProject(templatesFS, templateName, destDir)
if err != nil {
    return err
}

// Execute post-scaffolding actions
err = processor.ExecutePostScaffoldActions(destDir)
if err != nil {
    return err
}
```

## 🛠️ Template Functions

The template processor provides several built-in functions:

### Comparison Functions

- `eq(a, b interface{}) bool`: Equality comparison
- `ne(a, b interface{}) bool`: Inequality comparison

### String Functions

- `lower(s string) string`: Convert to lowercase
- `upper(s string) string`: Convert to uppercase
- `title(s string) string`: Title case
- `trim(s string) string`: Trim whitespace

### Array Functions

- `contains(slice []string, item string) bool`: Check if array contains item

### Usage in Templates

```go
{{if eq .Framework "React"}}
// React-specific code
{{end}}

{{if contains .Features "API"}}
// API-related code
{{end}}

{{lower .ProjectName}}-config.js
```

## 🔍 Conditional Logic

### Condition Syntax

The system supports simple conditional logic:

#### Equality Conditions

```
"IncludeTesting == true"
"Framework == 'React'"
```

#### Inequality Conditions

```
"TestingFramework != 'Jest'"
"Framework != 'Vue'"
```

### Condition Evaluation

```go
// Evaluate condition
result, err := processor.evaluateCondition("IncludeTesting == true")
if err != nil {
    return err
}

if result {
    // Show parameter or include file
}
```

## ✅ Validation System

### String Validation

```go
// Validate string parameter
param := Parameter{
    Name: "ProjectName",
    Type: "string",
    Validation: &Validation{
        Regex: "^[a-z0-9-]+$",
        ErrorMessage: "Project name can only contain lowercase letters, numbers, and hyphens.",
    },
}

err := processor.ValidateParameter(param, "my-project")
if err != nil {
    return err
}
```

### Built-in Validations

#### Project Name Validation

```go
Validation{
    Regex: "^[a-z0-9-]+$",
    ErrorMessage: "Project name can only contain lowercase letters, numbers, and hyphens.",
}
```

#### Email Validation

```go
Validation{
    Regex: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
    ErrorMessage: "Please enter a valid email address.",
}
```

## 🔧 Post-Scaffolding Actions

### File Deletion

```go
PostScaffold{
    FilesToDelete: []FileAction{
        {
            Path:      "jest.config.js",
            Condition: "TestingFramework != 'Jest'",
        },
        {
            Path:      "Dockerfile",
            Condition: "IncludeDocker == false",
        },
    },
}
```

### Command Execution

```go
PostScaffold{
    Commands: []CommandAction{
        {
            Command:     "git init",
            Description: "Initializing Git repository...",
        },
        {
            Command:     "npm install",
            Description: "Installing dependencies...",
            Condition:   "InstallDeps == true",
        },
    },
}
```

## 🧪 Testing

### Unit Tests

```go
func TestParameterValidation(t *testing.T) {
    param := Parameter{
        Name: "ProjectName",
        Type: "string",
        Validation: &Validation{
            Regex: "^[a-z0-9-]+$",
        },
    }

    processor := NewParameterProcessor(&TemplateManifest{})

    // Test valid input
    err := processor.ValidateParameter(param, "valid-project")
    if err != nil {
        t.Errorf("Valid input should not error: %v", err)
    }

    // Test invalid input
    err = processor.ValidateParameter(param, "Invalid Project")
    if err == nil {
        t.Error("Invalid input should error")
    }
}
```

### Integration Tests

```go
func TestTemplateProcessing(t *testing.T) {
    manifest := &TemplateManifest{
        Name: "Test Template",
        Description: "Test description",
    }

    values := map[string]interface{}{
        "ProjectName": "test-project",
        "Owner": "test-owner",
    }

    processor := NewTemplateProcessor(manifest, values)

    content := "Project: {{.ProjectName}}, Owner: {{.Owner}}"
    expected := "Project: test-project, Owner: test-owner"

    result, err := processor.ProcessTemplate(content)
    if err != nil {
        t.Fatalf("Failed to process template: %v", err)
    }

    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

## 🔍 Debugging

### Debug Logging

The package includes debug logging for development:

```go
fmt.Printf("DEBUG: Loading manifest from: %s\n", manifestPath)
fmt.Printf("DEBUG: Successfully read manifest bytes: %d bytes\n", len(manifestBytes))
fmt.Printf("DEBUG: Successfully parsed manifest: %s - %s\n", manifest.Name, manifest.Description)
```

### Common Issues

1. **Template Not Found**

   - Check template directory structure
   - Verify `template.json` exists and is valid JSON
   - Check template name in discovery logic

2. **Parameter Issues**

   - Verify parameter definitions in `template.json`
   - Check parameter validation logic
   - Test conditional parameter visibility

3. **File Processing Errors**
   - Check file paths and permissions
   - Verify template syntax in files
   - Test template variable substitution

## 🔮 Future Enhancements

### Planned Features

1. **Advanced Conditions**: More sophisticated conditional logic
2. **Custom Validators**: User-defined validation functions
3. **Template Inheritance**: Template composition and inheritance
4. **External Data**: Fetch data from external sources
5. **Template Versioning**: Template version management

### Extension Points

1. **Plugin System**: Custom template processors
2. **Custom Functions**: User-defined template functions
3. **Template Composition**: Combine multiple templates
4. **Condition Engine**: Advanced expression parser

## 📚 Related Documentation

- [Template System Documentation](../../docs/template-system.md)
- [Architecture Overview](../../docs/architecture.md)
- [Development Guide](../../docs/development.md)
- [User Guide](../../docs/user-guide.md)

---

**Last Updated**: August 3, 2025  
**Version**: v0.5.0  
**Maintainers**: [Project Maintainers]
