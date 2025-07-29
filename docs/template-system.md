# Template System

This document provides a comprehensive guide to the Open Workbench CLI's dynamic template system, including how to create, configure, and use templates.

## 🎯 Overview

The Open Workbench CLI uses a sophisticated dynamic template system that allows for:

- **Conditional Logic**: Show/hide parameters and files based on user choices
- **Parameter Validation**: Custom validation rules with regex patterns
- **Post-Scaffolding Actions**: Automatic cleanup and setup after project creation
- **Parameter Groups**: Organized parameter collection for better UX
- **Multiple Parameter Types**: String, boolean, select, and multiselect inputs

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
  "helpText": "This will add Jest or Vitest configuration to your project."
}
```

#### 3. Select Parameters

Single-choice dropdown with predefined options.

```json
{
  "name": "TestingFramework",
  "prompt": "Which testing framework?",
  "type": "select",
  "default": "Jest",
  "options": ["Jest", "Vitest"],
  "condition": "IncludeTesting == true"
}
```

#### 4. Multiselect Parameters

Multiple-choice selection.

```json
{
  "name": "Features",
  "prompt": "Which features would you like to include?",
  "type": "multiselect",
  "options": ["Authentication", "Database", "API", "UI Components"],
  "default": ["API"]
}
```

### Parameter Properties

| Property     | Type    | Required | Description                                           |
| ------------ | ------- | -------- | ----------------------------------------------------- |
| `name`       | string  | Yes      | Unique parameter identifier                           |
| `prompt`     | string  | Yes      | User-facing question                                  |
| `type`       | string  | Yes      | Parameter type (string, boolean, select, multiselect) |
| `group`      | string  | No       | Group for organizing parameters                       |
| `required`   | boolean | No       | Whether parameter is required (default: false)        |
| `default`    | any     | No       | Default value                                         |
| `options`    | array   | No       | Available options for select/multiselect              |
| `condition`  | string  | No       | Conditional visibility rule                           |
| `helpText`   | string  | No       | Additional help text                                  |
| `validation` | object  | No       | Validation rules                                      |

### Parameter Groups

Parameters can be organized into groups for better UX:

```json
{
  "parameters": [
    {
      "name": "ProjectName",
      "group": "Project Details",
      "type": "string"
    },
    {
      "name": "Owner",
      "group": "Project Details",
      "type": "string"
    },
    {
      "name": "IncludeTesting",
      "group": "Testing & Quality",
      "type": "boolean"
    },
    {
      "name": "IncludeDocker",
      "group": "Deployment",
      "type": "boolean"
    }
  ]
}
```

## 🔄 Conditional Logic

### Condition Syntax

The template system supports simple conditional logic:

#### Equality Conditions

```json
{
  "condition": "IncludeTesting == true"
}
```

#### Inequality Conditions

```json
{
  "condition": "TestingFramework != 'Jest'"
}
```

#### String Comparisons

```json
{
  "condition": "Framework == 'React'"
}
```

### Conditional Parameters

Parameters can be conditionally shown based on other parameter values:

```json
{
  "parameters": [
    {
      "name": "IncludeTesting",
      "prompt": "Include testing framework?",
      "type": "boolean",
      "default": true
    },
    {
      "name": "TestingFramework",
      "prompt": "Which testing framework?",
      "type": "select",
      "options": ["Jest", "Vitest"],
      "condition": "IncludeTesting == true"
    }
  ]
}
```

### Conditional Files

Files can be conditionally included or excluded:

```json
{
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "jest.config.js",
        "condition": "TestingFramework != 'Jest'"
      },
      {
        "path": "vitest.config.js",
        "condition": "TestingFramework != 'Vitest'"
      },
      {
        "path": "tests/",
        "condition": "IncludeTesting == false"
      }
    ]
  }
}
```

## ✅ Validation System

### String Validation

String parameters can have custom validation rules:

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

### Built-in Validations

The system provides several built-in validations:

#### Project Name Validation

```json
{
  "validation": {
    "regex": "^[a-z0-9-]+$",
    "errorMessage": "Project name can only contain lowercase letters, numbers, and hyphens."
  }
}
```

#### Email Validation

```json
{
  "validation": {
    "regex": "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
    "errorMessage": "Please enter a valid email address."
  }
}
```

#### URL Validation

```json
{
  "validation": {
    "regex": "^https?://.+",
    "errorMessage": "Please enter a valid URL starting with http:// or https://"
  }
}
```

## 🔧 Post-Scaffolding Actions

### File Deletion

Remove files based on conditions:

```json
{
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "jest.config.js",
        "condition": "TestingFramework != 'Jest'"
      },
      {
        "path": "Dockerfile",
        "condition": "IncludeDocker == false"
      },
      {
        "path": "tailwind.config.js",
        "condition": "IncludeTailwind == false"
      }
    ]
  }
}
```

### Command Execution

Run commands after scaffolding:

```json
{
  "postScaffold": {
    "commands": [
      {
        "command": "git init",
        "description": "Initializing Git repository...",
        "condition": "InitGit == true"
      },
      {
        "command": "npm install",
        "description": "Installing dependencies...",
        "condition": "InstallDeps == true"
      },
      {
        "command": "npm run test",
        "description": "Running tests...",
        "condition": "IncludeTesting == true"
      }
    ]
  }
}
```

### Optional Git Initialization

The `InitGit` parameter allows users to control whether git initialization occurs:

```json
{
  "name": "InitGit",
  "prompt": "Initialize Git repository?",
  "group": "Final Steps",
  "type": "boolean",
  "default": true,
  "helpText": "This will run 'git init' to initialize a new Git repository."
}
```

This parameter is available in CLI mode via the `--no-git` flag.

## 📝 Template File Processing

### Template Variables

Template files can use Go template syntax for variable substitution:

#### Basic Variables

```go
{{.ProjectName}}
{{.Owner}}
{{.IncludeTesting}}
```

#### Conditional Content

```go
{{if .IncludeTesting}}
import { render, screen } from '@testing-library/react';
{{end}}
```

#### Conditional Files

```go
{{if .IncludeDocker}}
# Dockerfile content
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
{{end}}
```

#### Template Functions

The system provides several template functions:

```go
{{eq .Framework "React"}}     // Equality comparison
{{ne .Framework "Vue"}}       // Inequality comparison
{{contains .Features "API"}}  // Array contains check
{{lower .ProjectName}}        // Convert to lowercase
{{upper .ProjectName}}        // Convert to uppercase
{{title .ProjectName}}        // Title case
{{trim .ProjectName}}         // Trim whitespace
```

### File Name Templates

File names can also be templated:

```
{{if .IncludeTesting}}tests/{{end}}App.test.tsx
{{if .IncludeDocker}}Dockerfile{{end}}
{{lower .ProjectName}}-config.js
```

## 🎨 Template Examples

### Next.js Template

```json
{
  "name": "Next.js Production-Grade",
  "description": "A fully-featured Next.js application with testing, linting, and optional CI/CD.",
  "parameters": [
    {
      "name": "ProjectName",
      "prompt": "Project Name:",
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
      "prompt": "Project Owner:",
      "group": "Project Details",
      "type": "string",
      "required": true
    },
    {
      "name": "IncludeTesting",
      "prompt": "Include a testing framework?",
      "group": "Testing & Quality",
      "type": "boolean",
      "default": true
    },
    {
      "name": "TestingFramework",
      "prompt": "Which testing framework?",
      "group": "Testing & Quality",
      "type": "select",
      "default": "Jest",
      "options": ["Jest", "Vitest"],
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
    },
    {
      "name": "InstallDeps",
      "prompt": "Install dependencies after setup?",
      "group": "Final Steps",
      "type": "boolean",
      "default": true,
      "helpText": "This will run 'npm install' automatically for you."
    }
  ],
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "jest.config.js",
        "condition": "TestingFramework != 'Jest'"
      },
      {
        "path": "vitest.config.js",
        "condition": "TestingFramework != 'Vitest'"
      },
      {
        "path": "tests/",
        "condition": "IncludeTesting == false"
      },
      {
        "path": "Dockerfile",
        "condition": "IncludeDocker == false"
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
        "command": "git init",
        "description": "Initializing Git repository..."
      },
      {
        "command": "npm install",
        "description": "Installing project dependencies...",
        "condition": "InstallDeps == true"
      }
    ]
  }
}
```

### FastAPI Template

```json
{
  "name": "FastAPI Basic",
  "description": "A FastAPI backend template with automatic API documentation.",
  "parameters": [
    {
      "name": "ProjectName",
      "prompt": "Project Name:",
      "type": "string",
      "required": true,
      "validation": {
        "regex": "^[a-z0-9-]+$",
        "errorMessage": "Project name can only contain lowercase letters, numbers, and hyphens."
      }
    },
    {
      "name": "IncludeDatabase",
      "prompt": "Include database support?",
      "type": "boolean",
      "default": false
    },
    {
      "name": "DatabaseType",
      "prompt": "Which database?",
      "type": "select",
      "options": ["PostgreSQL", "SQLite", "MySQL"],
      "condition": "IncludeDatabase == true"
    }
  ],
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "database/",
        "condition": "IncludeDatabase == false"
      }
    ],
    "commands": [
      {
        "command": "python -m venv venv",
        "description": "Creating virtual environment..."
      },
      {
        "command": "pip install -r requirements.txt",
        "description": "Installing dependencies..."
      }
    ]
  }
}
```

## 🛠️ Best Practices

### Template Design

1. **Clear Parameter Names**: Use descriptive, consistent parameter names
2. **Logical Grouping**: Group related parameters together
3. **Sensible Defaults**: Provide helpful default values
4. **Comprehensive Validation**: Validate user input appropriately
5. **Clear Help Text**: Provide helpful guidance for complex parameters

### Conditional Logic

1. **Simple Conditions**: Keep conditions simple and readable
2. **Consistent Naming**: Use consistent parameter names in conditions
3. **Test Conditions**: Thoroughly test conditional logic
4. **Document Dependencies**: Document parameter dependencies

### File Organization

1. **Logical Structure**: Organize files logically
2. **Conditional Files**: Use conditional file inclusion sparingly
3. **Clear Naming**: Use clear, descriptive file names
4. **Minimal Dependencies**: Minimize template dependencies

### Validation Rules

1. **Appropriate Validation**: Validate only what's necessary
2. **Clear Error Messages**: Provide helpful error messages
3. **Test Validation**: Test validation rules thoroughly
4. **Common Patterns**: Use common validation patterns

## 🔍 Debugging Templates

### Common Issues

1. **Template Not Found**: Check template directory structure
2. **Invalid JSON**: Validate template.json syntax
3. **Missing Parameters**: Ensure all required parameters are defined
4. **Condition Errors**: Check condition syntax and parameter names
5. **File Not Found**: Verify file paths in template

### Debug Tools

1. **Template Validation**: Use the CLI's built-in validation
2. **Parameter Testing**: Test parameter collection separately
3. **File Processing**: Check file processing step by step
4. **Condition Testing**: Test conditional logic independently

### Debug Commands

```bash
# Validate template structure
open-workbench-cli validate template-name

# Test parameter collection
open-workbench-cli test-params template-name

# Preview template output
open-workbench-cli preview template-name --params '{"ProjectName":"test"}'
```

## 🔮 Future Enhancements

### Planned Features

1. **Advanced Conditions**: More sophisticated conditional logic
2. **Template Inheritance**: Template composition and inheritance
3. **Custom Validators**: User-defined validation functions
4. **Template Versioning**: Template version management
5. **Template Marketplace**: External template distribution

### Extension Points

1. **Plugin System**: Custom template processors
2. **Custom Functions**: User-defined template functions
3. **External Data**: Fetch data from external sources
4. **Template Composition**: Combine multiple templates

---

**Last Updated**: 07/29/2025  
**Version**: v0.5.0  
**Maintainers**: [Project Maintainers]
