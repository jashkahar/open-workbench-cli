# Template System

This document provides a comprehensive guide to the Open Workbench CLI's dynamic template system, including how to create, configure, and use templates.

## 🎯 Overview

The Open Workbench CLI uses a sophisticated dynamic template system that allows for:

- **Conditional Logic**: Show/hide parameters and files based on user choices
- **Parameter Validation**: Custom validation rules with regex patterns
- **Post-Scaffolding Actions**: Automatic cleanup and setup after project creation
- **Parameter Groups**: Organized parameter collection for better UX
- **Multiple Parameter Types**: String, boolean, select, and multiselect inputs
- **Smart Command Integration**: Works seamlessly with the smart command system

## 📁 Template Structure

### Directory Layout

Each template follows this structure:

```
templates/
└── template-name/
    ├── template.json          # Template manifest (required)
    ├── package.json           # Template files
    ├── src/
    │   └── components/
    │       └── App.tsx
    ├── tests/
    │   └── App.test.tsx
    └── README.md
```

### Template Manifest (`template.json`)

The template manifest is the heart of the template system. It defines:

- Template metadata (name, description)
- Parameter definitions and validation rules
- Post-scaffolding actions
- Conditional logic

```json
{
  "name": "Template Display Name",
  "description": "Template description for users",
  "parameters": [
    {
      "name": "ProjectName",
      "prompt": "What is your project name?",
      "group": "Project Details",
      "type": "string",
      "required": true,
      "validation": {
        "regex": "^[a-z0-9-]+$",
        "errorMessage": "Project name can only contain lowercase letters, numbers, and hyphens."
      }
    }
  ],
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "file-to-delete.js",
        "condition": "IncludeFeature == false"
      }
    ],
    "commands": [
      {
        "command": "npm install",
        "description": "Installing dependencies...",
        "condition": "InstallDeps == true"
      }
    ]
  }
}
```

## 📋 Parameter System

### Parameter Types

#### 1. String Parameters

Basic text input with optional validation.

```json
{
  "name": "ProjectName",
  "prompt": "What is your project name?",
  "type": "string",
  "required": true,
  "default": "my-project",
  "helpText": "This will be used as the directory name and package name.",
  "validation": {
    "regex": "^[a-z0-9-]+$",
    "errorMessage": "Project name can only contain lowercase letters, numbers, and hyphens."
  }
}
```

#### 2. Boolean Parameters

Yes/No questions with default values.

```json
{
  "name": "IncludeTesting",
  "prompt": "Include testing framework?",
  "type": "boolean",
  "default": true,
  "helpText": "Add Jest or Vitest for testing your components."
}
```

#### 3. Select Parameters

Single-choice dropdown with predefined options.

```json
{
  "name": "TestingFramework",
  "prompt": "Which testing framework would you like?",
  "type": "select",
  "options": ["Jest", "Vitest"],
  "default": "Jest",
  "condition": "IncludeTesting == true"
}
```

#### 4. Multiselect Parameters

Multiple-choice selection for complex configurations.

```json
{
  "name": "Features",
  "prompt": "Which features would you like to include?",
  "type": "multiselect",
  "options": ["Authentication", "Database", "API", "Admin Panel"],
  "default": ["Authentication"],
  "helpText": "Select all features you want in your application."
}
```

### Parameter Groups

Organize parameters into logical groups for better UX:

```json
{
  "parameters": [
    {
      "name": "ProjectName",
      "group": "Project Details",
      "type": "string",
      "required": true
    },
    {
      "name": "Owner",
      "group": "Project Details",
      "type": "string",
      "required": true
    },
    {
      "name": "IncludeTesting",
      "group": "Development Tools",
      "type": "boolean",
      "default": true
    },
    {
      "name": "IncludeDocker",
      "group": "Deployment",
      "type": "boolean",
      "default": true
    }
  ]
}
```

### Parameter Validation

#### Regex Validation

Use regex patterns for string validation:

```json
{
  "name": "ProjectName",
  "type": "string",
  "validation": {
    "regex": "^[a-z0-9-]+$",
    "errorMessage": "Project name can only contain lowercase letters, numbers, and hyphens."
  }
}
```

#### Required Field Validation

Ensure critical parameters are provided:

```json
{
  "name": "ProjectName",
  "type": "string",
  "required": true,
  "helpText": "This is required for project setup."
}
```

#### Custom Validation

Implement custom validation logic:

```json
{
  "name": "Port",
  "type": "string",
  "validation": {
    "regex": "^[0-9]+$",
    "errorMessage": "Port must be a number.",
    "custom": {
      "min": 1024,
      "max": 65535,
      "message": "Port must be between 1024 and 65535."
    }
  }
}
```

## 🔄 Conditional Logic

### Parameter Conditions

Show/hide parameters based on other parameter values:

```json
{
  "parameters": [
    {
      "name": "IncludeTesting",
      "type": "boolean",
      "default": true
    },
    {
      "name": "TestingFramework",
      "type": "select",
      "options": ["Jest", "Vitest"],
      "condition": "IncludeTesting == true"
    },
    {
      "name": "IncludeE2E",
      "type": "boolean",
      "default": false,
      "condition": "IncludeTesting == true"
    }
  ]
}
```

### File Conditions

Generate or delete files based on parameter values:

```json
{
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "src/tests/App.test.tsx",
        "condition": "IncludeTesting == false"
      },
      {
        "path": "tailwind.config.js",
        "condition": "IncludeTailwind == false"
      }
    ]
  }
}
```

### Command Conditions

Execute commands conditionally:

```json
{
  "postScaffold": {
    "commands": [
      {
        "command": "npm install",
        "description": "Installing dependencies...",
        "condition": "InstallDeps == true"
      },
      {
        "command": "git init",
        "description": "Initializing git repository...",
        "condition": "InitGit == true"
      }
    ]
  }
}
```

## 🚀 Post-Scaffolding Actions

### File Operations

#### Delete Files

Remove files based on conditions:

```json
{
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "unused-config.js",
        "condition": "IncludeFeature == false"
      },
      {
        "path": "src/components/Example.tsx",
        "condition": "IncludeExamples == false"
      }
    ]
  }
}
```

#### Rename Files

Rename files based on parameters:

```json
{
  "postScaffold": {
    "filesToRename": [
      {
        "from": "src/App.tsx",
        "to": "src/{{ProjectName}}.tsx",
        "condition": "CustomAppName == true"
      }
    ]
  }
}
```

### Command Execution

#### Install Dependencies

```json
{
  "postScaffold": {
    "commands": [
      {
        "command": "npm install",
        "description": "Installing dependencies...",
        "condition": "InstallDeps == true",
        "cwd": "."
      }
    ]
  }
}
```

#### Initialize Git

```json
{
  "postScaffold": {
    "commands": [
      {
        "command": "git init",
        "description": "Initializing git repository...",
        "condition": "InitGit == true"
      },
      {
        "command": "git add .",
        "description": "Adding files to git...",
        "condition": "InitGit == true"
      },
      {
        "command": "git commit -m \"Initial commit\"",
        "description": "Creating initial commit...",
        "condition": "InitGit == true"
      }
    ]
  }
}
```

#### Custom Scripts

```json
{
  "postScaffold": {
    "commands": [
      {
        "command": "echo 'Project {{ProjectName}} created successfully!'",
        "description": "Displaying success message...",
        "condition": "ShowMessage == true"
      }
    ]
  }
}
```

## 🎨 Template Examples

### Next.js Full-Stack Template

```json
{
  "name": "Next.js Full-Stack",
  "description": "Production-ready Next.js application with TypeScript, testing, and Docker",
  "parameters": [
    {
      "name": "ProjectName",
      "prompt": "What is your project name?",
      "group": "Project Details",
      "type": "string",
      "required": true,
      "validation": {
        "regex": "^[a-z0-9-]+$",
        "errorMessage": "Project name can only contain lowercase letters, numbers, and hyphens."
      }
    },
    {
      "name": "Owner",
      "prompt": "Who is the project owner?",
      "group": "Project Details",
      "type": "string",
      "required": true
    },
    {
      "name": "IncludeTesting",
      "prompt": "Include testing framework?",
      "group": "Development Tools",
      "type": "boolean",
      "default": true
    },
    {
      "name": "TestingFramework",
      "prompt": "Which testing framework?",
      "group": "Development Tools",
      "type": "select",
      "options": ["Jest", "Vitest"],
      "default": "Jest",
      "condition": "IncludeTesting == true"
    },
    {
      "name": "IncludeDocker",
      "prompt": "Include Docker configuration?",
      "group": "Deployment",
      "type": "boolean",
      "default": true
    },
    {
      "name": "IncludeTailwind",
      "prompt": "Include Tailwind CSS?",
      "group": "Styling",
      "type": "boolean",
      "default": true
    }
  ],
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "src/app/globals.css",
        "condition": "IncludeTailwind == false"
      },
      {
        "path": "tailwind.config.js",
        "condition": "IncludeTailwind == false"
      },
      {
        "path": "postcss.config.js",
        "condition": "IncludeTailwind == false"
      }
    ],
    "commands": [
      {
        "command": "npm install",
        "description": "Installing dependencies...",
        "condition": "InstallDeps == true"
      },
      {
        "command": "git init",
        "description": "Initializing git repository...",
        "condition": "InitGit == true"
      }
    ]
  }
}
```

### FastAPI Basic Template

```json
{
  "name": "FastAPI Basic",
  "description": "FastAPI backend with automatic API documentation",
  "parameters": [
    {
      "name": "ProjectName",
      "prompt": "What is your project name?",
      "type": "string",
      "required": true
    },
    {
      "name": "Owner",
      "prompt": "Who is the project owner?",
      "type": "string",
      "required": true
    },
    {
      "name": "Database",
      "prompt": "Which database would you like to use?",
      "type": "select",
      "options": ["SQLite", "PostgreSQL", "MongoDB"],
      "default": "SQLite"
    },
    {
      "name": "IncludeAuth",
      "prompt": "Include authentication?",
      "type": "boolean",
      "default": true
    }
  ],
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "src/auth/",
        "condition": "IncludeAuth == false"
      }
    ],
    "commands": [
      {
        "command": "python -m venv venv",
        "description": "Creating virtual environment...",
        "condition": "CreateVenv == true"
      },
      {
        "command": "pip install -r requirements.txt",
        "description": "Installing Python dependencies...",
        "condition": "InstallDeps == true"
      }
    ]
  }
}
```

## 🔧 Advanced Features

### Template Inheritance

Create base templates that can be extended:

```json
{
  "name": "Base Template",
  "description": "Base template with common configurations",
  "parameters": [
    {
      "name": "ProjectName",
      "type": "string",
      "required": true
    }
  ],
  "extends": "base-template"
}
```

### Dynamic File Generation

Generate files based on parameter values:

```json
{
  "parameters": [
    {
      "name": "Database",
      "type": "select",
      "options": ["SQLite", "PostgreSQL", "MongoDB"]
    }
  ],
  "files": [
    {
      "path": "src/database/{{Database}}.py",
      "condition": "Database != null"
    }
  ]
}
```

### Environment Configuration

Generate environment files:

```json
{
  "postScaffold": {
    "envFiles": [
      {
        "path": ".env.example",
        "template": "env.example.tmpl",
        "condition": "IncludeEnv == true"
      }
    ]
  }
}
```

## 🧪 Testing Templates

### Template Validation

Test your templates before using them:

```bash
# List all templates
om list-templates

# Test template processing
go test ./internal/templating/ -v

# Test specific template
go test ./internal/templating/ -run TestProcessTemplate
```

### Template Testing Structure

```go
func TestTemplateProcessing(t *testing.T) {
    tests := []struct {
        name     string
        template string
        params   map[string]interface{}
        expect   string
    }{
        {
            name:     "basic template",
            template: "Hello {{ProjectName}}!",
            params:   map[string]interface{}{"ProjectName": "World"},
            expect:   "Hello World!",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := processTemplate(tt.template, tt.params)
            if result != tt.expect {
                t.Errorf("expected %s, got %s", tt.expect, result)
            }
        })
    }
}
```

## 🔒 Security Considerations

### Template Security

- **Validate all inputs**: Never trust user input
- **Sanitize file paths**: Prevent path traversal attacks
- **Validate template names**: Ensure safe template names
- **Check file operations**: Validate all file operations

### Security Best Practices

```json
{
  "parameters": [
    {
      "name": "ProjectName",
      "type": "string",
      "validation": {
        "regex": "^[a-z0-9-]+$",
        "errorMessage": "Invalid project name"
      }
    }
  ],
  "security": {
    "allowedPaths": ["./src", "./public"],
    "forbiddenPatterns": ["../", "..\\", "javascript:"]
  }
}
```

## 🚀 Performance Optimization

### Template Caching

Templates are cached for better performance:

```go
// Template caching
var templateCache = make(map[string]*Template)

func getTemplate(name string) (*Template, error) {
    if cached, exists := templateCache[name]; exists {
        return cached, nil
    }
    // Load and cache template
}
```

### Lazy Loading

Load templates only when needed:

```go
func loadTemplate(name string) (*Template, error) {
    // Load template from embedded filesystem
    // Parse template manifest
    // Validate template structure
    return template, nil
}
```

## 📚 Best Practices

### Template Design

1. **Clear Naming**: Use descriptive template names
2. **Comprehensive Documentation**: Document all parameters
3. **Sensible Defaults**: Provide good default values
4. **Progressive Enhancement**: Start simple, add complexity
5. **Error Handling**: Provide clear error messages

### Parameter Design

1. **Logical Grouping**: Group related parameters
2. **Clear Prompts**: Use descriptive prompts
3. **Helpful Validation**: Provide useful error messages
4. **Conditional Logic**: Use conditions sparingly
5. **Default Values**: Provide sensible defaults

### Post-Scaffolding Actions

1. **Minimal Commands**: Keep commands simple
2. **Clear Descriptions**: Describe what each command does
3. **Error Handling**: Handle command failures gracefully
4. **Conditional Execution**: Use conditions appropriately
5. **User Feedback**: Provide clear progress messages

## 🔮 Future Enhancements

### Planned Features

1. **Template Marketplace**: Community template sharing
2. **Template Versioning**: Version control for templates
3. **Advanced Conditions**: Complex conditional logic
4. **Template Composition**: Combine multiple templates
5. **Plugin System**: Extensible template system

### Template Ecosystem

1. **Community Templates**: User-contributed templates
2. **Template Validation**: Automated template testing
3. **Template Documentation**: Comprehensive documentation
4. **Template Examples**: Real-world examples
5. **Template Support**: Community support for templates

---

**Maintainer**: Jash Kahar  
**Last Updated**: August 3, 2025
