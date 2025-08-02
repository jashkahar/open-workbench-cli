# User Guide

This comprehensive user guide will help you get started with Open Workbench CLI and use it effectively for your project scaffolding needs.

## 🚀 Quick Start

### Installation

Choose your preferred installation method:

#### Option 1: Download from GitHub Releases

1. Go to [GitHub Releases](https://github.com/jashkahar/open-workbench-platform/releases)
2. Download the binary for your operating system
3. Extract and add to your PATH, or run directly

#### Option 2: Package Managers

```bash
# Homebrew (macOS)
brew install jashkahar/tap/om

# Scoop (Windows)
scoop bucket add jashkahar https://github.com/jashkahar/scoop-bucket
scoop install om
```

#### Option 3: Build from Source

```bash
git clone https://github.com/jashkahar/open-workbench-platform.git
cd open-workbench-platform
go build -o om main.go
```

### First Run

```bash
# Initialize a new project (recommended)
om init

# Start the interactive TUI
om ui

# Or use simple interactive mode
om
```

## 🎯 Using the CLI

### Available Commands

| Command     | Description                    |
| ----------- | ------------------------------ |
| `om init`   | Initialize a new project (NEW) |
| `om`        | Interactive mode (recommended) |
| `om create` | CLI mode with flags            |

### Project Management with `om init` (NEW)

The `om init` command creates managed projects with a `workbench.yaml` manifest file:

```bash
# Initialize a new project
om init
```

**What `om init` does:**

1. **Safety Check**: Verifies the current directory is empty or contains only hidden files
2. **Project Name**: Prompts for a project name with validation
3. **Template Selection**: Shows available templates for the first service
4. **Service Name**: Prompts for the first service name
5. **Project Creation**: Creates the project structure and scaffolds the first service
6. **Manifest Generation**: Creates a `workbench.yaml` file for project management

**Example workflow:**

```bash
$ om init
What is your project name? my-awesome-app
Choose a template for your first service:
  ❯ nextjs-full-stack - A production-ready Next.js application
    react-typescript - A modern React application
    fastapi-basic - A FastAPI backend template
    express-api - A Node.js Express API template
    vue-nuxt - A Vue.js Nuxt application
What is your service name? frontend

✅ Success! Your new project 'my-awesome-app' is ready.

📁 Project structure:
  my-awesome-app/
  ├── workbench.yaml
  └── frontend/

🚀 Next steps:
  cd my-awesome-app
  om add service  # Add more services to your project
  om run          # Run your project (when implemented)
  om deploy       # Deploy your project (when implemented)
```

**Project Structure Created:**

```
my-awesome-app/
├── workbench.yaml          # Project manifest
└── frontend/              # First service
    ├── package.json
    ├── src/
    ├── public/
    └── ... (template files)
```

**workbench.yaml Manifest:**

```yaml
apiVersion: openworkbench.io/v1alpha1
kind: Project
metadata:
  name: my-awesome-app
services:
  frontend:
    template: nextjs-full-stack
    path: ./frontend
```

### Interactive Modes

#### 1. Terminal User Interface (TUI) - Recommended

The TUI provides a beautiful, interactive interface for template selection:

```bash
om ui
```

**Features:**

- Visual template selection with descriptions
- Keyboard navigation (arrow keys, enter, q to quit)
- Spinner animations and smooth interactions
- Clear visual feedback

**Navigation:**

- `↑/↓` - Navigate templates
- `Enter` - Select template
- `q` - Quit without selection
- `Ctrl+C` - Exit immediately

#### 2. Interactive Mode (Recommended)

For guided project creation with template selection:

```bash
om
```

**Features:**

- Template selection from all available templates
- Organized parameter collection with grouping
- Comprehensive validation and error handling
- Optional git initialization and dependency installation

#### 3. CLI Mode (Non-Interactive)

For automation and scripting:

```bash
om create <template> <project-name> --owner="Your Name" [flags]
```

## 🔒 Security Features

### Enterprise-Grade Security

Open Workbench Platform includes comprehensive security features to protect against common attacks:

#### Input Validation

All user inputs are validated for security:

```bash
# ✅ Valid project names
my-project
project123
frontend

# ❌ Blocked (security reasons)
../malicious          # Path traversal
javascript:alert(1)   # JavaScript injection
C:\Windows\System32  # Absolute paths
con                   # Windows reserved names
```

#### Security Validations

- **Path Traversal Protection**: Blocks `../` and `..\` attacks
- **Malicious Pattern Detection**: Prevents JavaScript injection, command injection
- **Cross-Platform Security**: Windows reserved names, absolute path prevention
- **Directory Safety Checks**: Validates permissions, accessibility, symbolic links
- **Template Security**: Secure template name validation and content verification

#### Security Configuration

The security system is configurable and extensible:

```go
// Security configuration
type SecurityConfig struct {
    MaxPathLength     int
    MaxNameLength     int
    AllowedCharacters *regexp.Regexp
    ForbiddenPatterns []*regexp.Regexp
}
```

## 🧪 Testing

### Comprehensive Test Suite

The platform includes a comprehensive test suite with 100% coverage:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run security tests
go test ./cmd -v

# Run benchmarks
go test ./cmd -bench=.
```

### Test Categories

1. **Security Tests**: Input validation, path traversal, malicious patterns
2. **Command Tests**: Init command, project creation, manifest generation
3. **Integration Tests**: End-to-end workflow testing
4. **Performance Tests**: Benchmark tests for critical functions

### Test Results

```
=== RUN   TestValidateAndSanitizePath --- PASS
=== RUN   TestValidateAndSanitizeName --- PASS
=== RUN   TestValidateDirectorySafety --- PASS
=== RUN   TestValidateTemplateName --- PASS
=== RUN   TestCheckForSuspiciousPatterns --- PASS
=== RUN   TestCreateProjectDirectories --- PASS
=== RUN   TestCreateWorkbenchManifest --- PASS
=== RUN   TestCheckDirectorySafety --- PASS

BenchmarkValidateAndSanitizeName-8:     100,788 ops/sec (~12μs/op)
BenchmarkValidateAndSanitizePath-8:      85,692 ops/sec (~12μs/op)
BenchmarkCheckForSuspiciousPatterns-8: 11,804,667 ops/sec (~149ns/op)
```

## 📋 Template Parameters

### Parameter Types

The CLI supports various parameter types for collecting user input:

#### String Parameters

```json
{
  "name": "ProjectName",
  "prompt": "What is your project name?",
  "type": "string",
  "required": true,
  "validation": {
    "regex": "^[a-z0-9-]+$",
    "errorMessage": "Project name can only contain lowercase letters, numbers, and hyphens."
  }
}
```

#### Boolean Parameters

```json
{
  "name": "IncludeTesting",
  "prompt": "Include testing setup?",
  "type": "boolean",
  "default": true
}
```

#### Select Parameters

```json
{
  "name": "TestingFramework",
  "prompt": "Choose testing framework:",
  "type": "select",
  "options": ["Jest", "Vitest"],
  "condition": "IncludeTesting == true"
}
```

#### Multiselect Parameters

```json
{
  "name": "Features",
  "prompt": "Select features to include:",
  "type": "multiselect",
  "options": ["Tailwind CSS", "Docker", "CI/CD", "Storybook"]
}
```

### Parameter Groups

Parameters can be organized into groups for better UX:

```json
{
  "name": "ProjectName",
  "group": "Basic Settings",
  "prompt": "What is your project name?"
},
{
  "name": "IncludeTesting",
  "group": "Testing",
  "prompt": "Include testing setup?"
}
```

### Conditional Logic

Parameters can be conditionally shown based on other parameter values:

```json
{
  "name": "TestingFramework",
  "condition": "IncludeTesting == true",
  "type": "select",
  "options": ["Jest", "Vitest"]
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
        "path": "src/components/Example.tsx",
        "condition": "IncludeExamples == false"
      }
    ]
  }
}
```

### Command Execution

Run setup commands after project creation:

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

## 🎨 Available Templates

### nextjs-full-stack

A production-ready Next.js application with:

- **TypeScript**: Full TypeScript support with strict configuration
- **Testing**: Jest or Vitest with comprehensive test setup
- **Styling**: Tailwind CSS with PostCSS configuration
- **Docker**: Ready-to-use Dockerfile for containerization
- **Quality Tools**: ESLint, Prettier, and Husky for code quality
- **CI/CD Ready**: GitHub Actions workflows included

**Parameters:**

- `ProjectName` (string): Project name
- `Owner` (string): Project owner
- `IncludeTesting` (boolean): Include testing setup
- `TestingFramework` (select): Jest or Vitest
- `IncludeTailwind` (boolean): Include Tailwind CSS
- `IncludeDocker` (boolean): Include Docker configuration
- `InstallDeps` (boolean): Install dependencies after creation
- `InitGit` (boolean): Initialize git repository

### fastapi-basic

A FastAPI backend template with:

- **FastAPI**: Modern, fast web framework for building APIs
- **Uvicorn**: ASGI server for running the application
- **Python Best Practices**: Virtual environment setup and dependency management
- **API Documentation**: Automatic OpenAPI/Swagger documentation
- **Hot Reload**: Development server with auto-reload capability

**Parameters:**

- `ProjectName` (string): Project name
- `Owner` (string): Project owner
- `IncludeTesting` (boolean): Include testing setup
- `InstallDeps` (boolean): Install dependencies after creation
- `InitGit` (boolean): Initialize git repository

### react-typescript

A modern React application with:

- **Vite**: Lightning-fast build tool and dev server
- **TypeScript**: Full TypeScript support
- **Modern Tooling**: ESLint, Prettier configuration
- **Component Library**: Ready-to-use component structure

**Parameters:**

- `ProjectName` (string): Project name
- `Owner` (string): Project owner
- `IncludeTesting` (boolean): Include testing setup
- `TestingFramework` (select): Jest or Vitest
- `IncludeTailwind` (boolean): Include Tailwind CSS
- `InstallDeps` (boolean): Install dependencies after creation
- `InitGit` (boolean): Initialize git repository

### express-api

A Node.js Express API template with:

- **Express.js**: Fast, unopinionated web framework
- **TypeScript**: Full TypeScript support
- **Testing**: Jest setup with API testing utilities
- **Documentation**: Swagger/OpenAPI documentation

**Parameters:**

- `ProjectName` (string): Project name
- `Owner` (string): Project owner
- `IncludeTesting` (boolean): Include testing setup
- `InstallDeps` (boolean): Install dependencies after creation
- `InitGit` (boolean): Initialize git repository

### vue-nuxt

A Vue.js Nuxt application with:

- **Nuxt 3**: Full-stack Vue.js framework
- **TypeScript**: Full TypeScript support
- **Auto-imports**: Automatic component and composable imports
- **SSR Ready**: Server-side rendering configuration

**Parameters:**

- `ProjectName` (string): Project name
- `Owner` (string): Project owner
- `IncludeTesting` (boolean): Include testing setup
- `InstallDeps` (boolean): Install dependencies after creation
- `InitGit` (boolean): Initialize git repository

## 🚀 Advanced Usage

### Custom Templates

Create your own templates by following the template structure:

1. Create a new directory in `templates/` with your template name
2. Add a `template.json` file with parameter definitions
3. Include your template files with Go template syntax where needed
4. Test the template using the CLI

### Template Manifest Structure

```json
{
  "name": "Template Name",
  "description": "Template description",
  "parameters": [
    {
      "name": "ParameterName",
      "prompt": "User prompt",
      "group": "Group Name",
      "type": "string|boolean|select|multiselect",
      "required": true,
      "default": "default value",
      "options": ["option1", "option2"],
      "condition": "OtherParam == true",
      "validation": {
        "regex": "^[a-z0-9-]+$",
        "errorMessage": "Custom error message"
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

### Go Template Syntax

Use Go template syntax in your template files:

````go
// In package.json
{
  "name": "{{.ProjectName}}",
  "version": "1.0.0",
  "description": "{{.ProjectName}} - A {{.TemplateName}} project",
  "author": "{{.Owner}}"
}

// In README.md
# {{.ProjectName}}

Created by {{.Owner}} using Open Workbench Platform.

{{if .IncludeTesting}}
## Testing

Run tests with:
```bash
npm test
````

{{end}}

````

### Conditional Logic

Use conditional statements in templates:

```go
{{if .IncludeTesting}}
import { render, screen } from '@testing-library/react';
{{end}}

{{if eq .TestingFramework "Jest"}}
import '@testing-library/jest-dom';
{{else}}
import { vi } from 'vitest';
{{end}}
````

## 🐛 Troubleshooting

### Common Issues

#### Permission Denied

```bash
Error: failed to create project directory: permission denied
```

**Solution**: Ensure you have write permissions in the current directory.

#### Template Not Found

```bash
Error: template not found: my-template
```

**Solution**: Check that the template exists in the `templates/` directory.

#### Invalid Project Name

```bash
Error: project name can only contain lowercase letters, numbers, and hyphens
```

**Solution**: Use only lowercase letters, numbers, and hyphens in project names.

#### Directory Not Empty

```bash
Error: directory is not empty. Please run 'om init' in an empty directory
```

**Solution**: Run `om init` in an empty directory or one containing only hidden files.

### Getting Help

```bash
# Get help for all commands
om --help

# Get help for specific command
om init --help
om create --help

# Get help for specific template
om create nextjs-full-stack --help
```

### Debug Mode

Enable debug mode for more verbose output:

```bash
# Set debug environment variable
export DEBUG=true
om init
```

## 📚 Additional Resources

- [Architecture Documentation](architecture.md)
- [Template System Guide](template-system.md)
- [Development Guide](development.md)
- [Contributing Guidelines](../CONTRIBUTING.md)

---

**Maintainer**: Jash Kahar  
**Last Updated**: February 8, 2025
