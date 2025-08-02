# Open Workbench Platform

![Status: Production Ready](https://img.shields.io/badge/status-production%20ready-green)
![Go](https://img.shields.io/badge/Go-1.24%2B-blue)
![License: MIT](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)
![Security](https://img.shields.io/badge/security-enterprise%20grade-brightgreen)
![Tests](https://img.shields.io/badge/tests-100%25%20coverage-brightgreen)

🚀 A powerful, secure command-line tool for scaffolding modern web applications with pre-configured templates and best practices.

<!-- TODO: Add demo GIF or asciinema recording here -->

## What is Open Workbench Platform?

Open Workbench Platform is a Go-based command-line interface that helps developers quickly bootstrap new projects using carefully crafted templates. It eliminates the repetitive setup process and ensures you start with a solid foundation following industry best practices.

### Key Features

- **🎯 Dynamic Template System**: Advanced templating with conditional logic, parameter validation, and post-scaffolding actions
- **🖥️ Terminal User Interface (TUI)**: Beautiful interactive interface for template selection and configuration
- **📋 Parameter Groups**: Organized parameter collection with grouping and conditional visibility
- **✅ Validation & Error Handling**: Comprehensive input validation with custom regex patterns and error messages
- **🔧 Post-Scaffolding Actions**: Automatic file cleanup, dependency installation, and git initialization
- **🌐 Cross-Platform**: Works on Windows, macOS, and Linux
- **📦 Multiple Installation Methods**: Homebrew, Scoop, GitHub Releases, and source builds
- **🔒 Enterprise-Grade Security**: Comprehensive input validation, path traversal protection, and malicious pattern detection
- **🧪 Comprehensive Testing**: 100% test coverage with security-focused test suites
- **📁 Project Management**: New `om init` command for creating managed projects with `workbench.yaml` manifests

## Quick Start

### Installation

#### From GitHub Releases

```bash
# Download from https://github.com/jashkahar/open-workbench-platform/releases
# Extract and add to PATH
```

#### Package Managers (Recommended)

```bash
# Homebrew (macOS)
brew install jashkahar/tap/om

# Scoop (Windows)
scoop bucket add jashkahar https://github.com/jashkahar/scoop-bucket
scoop install om
```

### Usage

#### Interactive Mode (Recommended)

```bash
# Initialize a new project with interactive prompts
om init

# Interactive template creation
om
```

#### CLI Mode with Flags

```bash
# Create a project with specific template
om create <template> <project-name> --owner="Your Name" [flags]

# Get help
om --help
om init --help
```

### CLI Mode Examples

```bash
# Create a Next.js project with all features
om create nextjs-full-stack my-app --owner="John Doe"

# Create a React project without testing
om create react-typescript my-react-app --owner="Dev Team" --no-testing

# Create a FastAPI project without git initialization
om create fastapi-basic my-api --owner="Backend Team" --no-git

# Get help for CLI mode
om create --help
```

## New: Project Management with `om init`

The new `om init` command creates managed projects with a `workbench.yaml` manifest file:

```bash
# Initialize a new project
om init

# This will:
# 1. Check directory safety
# 2. Prompt for project name
# 3. Select first service template
# 4. Create project structure
# 5. Generate workbench.yaml manifest
```

### Project Structure Created by `om init`

```
my-project/
├── workbench.yaml          # Project manifest
└── frontend/              # First service
    ├── package.json
    ├── src/
    └── ... (template files)
```

### workbench.yaml Manifest

```yaml
apiVersion: openworkbench.io/v1alpha1
kind: Project
metadata:
  name: my-project
services:
  frontend:
    template: nextjs-full-stack
    path: ./frontend
```

## Security Features

### 🔒 Enterprise-Grade Security

- **Input Validation**: Comprehensive validation for project names, paths, and template names
- **Path Traversal Protection**: Blocks `../` and `..\` attacks
- **Malicious Pattern Detection**: Prevents JavaScript injection, command injection, and other attacks
- **Cross-Platform Security**: Windows reserved names, absolute path prevention
- **Directory Safety Checks**: Validates permissions, accessibility, and symbolic links
- **Template Security**: Secure template name validation and content verification

### Security Validations

```bash
# ✅ Valid project names
my-project
project123
frontend

# ❌ Blocked (security reasons)
../malicious
javascript:alert(1)
C:\Windows\System32
con (Windows reserved)
```

## Available Templates

### 🎨 nextjs-full-stack

A production-ready Next.js application with:

- **TypeScript**: Full TypeScript support with strict configuration
- **Testing**: Jest or Vitest with comprehensive test setup
- **Styling**: Tailwind CSS with PostCSS configuration
- **Docker**: Ready-to-use Dockerfile for containerization
- **Quality Tools**: ESLint, Prettier, and Husky for code quality
- **CI/CD Ready**: GitHub Actions workflows included

### ⚡ fastapi-basic

A FastAPI backend template with:

- **FastAPI**: Modern, fast web framework for building APIs
- **Uvicorn**: ASGI server for running the application
- **Python Best Practices**: Virtual environment setup and dependency management
- **API Documentation**: Automatic OpenAPI/Swagger documentation
- **Hot Reload**: Development server with auto-reload capability

### 🎯 react-typescript

A modern React application with:

- **Vite**: Lightning-fast build tool and dev server
- **TypeScript**: Full TypeScript support
- **Modern Tooling**: ESLint, Prettier configuration
- **Component Library**: Ready-to-use component structure

### 🚀 express-api

A Node.js Express API template with:

- **Express.js**: Fast, unopinionated web framework
- **TypeScript**: Full TypeScript support
- **Testing**: Jest setup with API testing utilities
- **Documentation**: Swagger/OpenAPI documentation

### 🟢 vue-nuxt

A Vue.js Nuxt application with:

- **Nuxt 3**: Full-stack Vue.js framework
- **TypeScript**: Full TypeScript support
- **Auto-imports**: Automatic component and composable imports
- **SSR Ready**: Server-side rendering configuration

## Dynamic Template System

The CLI uses an advanced dynamic templating system that supports:

### Parameter Types

- **String**: Text input with validation
- **Boolean**: Yes/No questions with defaults
- **Select**: Single-choice dropdown
- **Multiselect**: Multiple-choice selection

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

### Validation

Custom validation with regex patterns and error messages:

```json
{
  "name": "ProjectName",
  "validation": {
    "regex": "^[a-z0-9-]+$",
    "errorMessage": "Project name can only contain lowercase letters, numbers, and hyphens."
  }
}
```

### Post-Scaffolding Actions

Automatic cleanup and setup after project creation:

- **File Deletion**: Remove files based on conditions
- **Command Execution**: Run setup commands like `npm install` or `git init`

## Project Structure

```
om/
├── main.go                   # Main CLI application entry point
├── cmd/                      # Command implementations
│   ├── root.go              # Root command setup
│   ├── init.go              # om init command
│   ├── types.go             # YAML manifest types
│   ├── security.go          # Security utilities
│   ├── security_test.go     # Security tests
│   └── init_test.go         # Init command tests
├── internal/
│   ├── templating/           # Dynamic templating system
│   │   ├── discovery.go      # Template discovery and validation
│   │   ├── parameters.go     # Parameter processing and validation
│   │   └── processor.go      # Template processing and file operations
│   └── template/             # Legacy template system (deprecated)
├── templates/                # Template directory
│   ├── nextjs-full-stack/    # Next.js full-stack template
│   ├── fastapi-basic/        # FastAPI backend template
│   ├── react-typescript/     # React + TypeScript template
│   ├── express-api/          # Express.js API template
│   └── vue-nuxt/            # Vue.js Nuxt template
├── docs/                     # Documentation
│   ├── README.md            # Documentation overview and index
│   ├── user-guide.md        # Complete user guide
│   ├── architecture.md      # System architecture
│   ├── template-system.md   # Template system guide
│   ├── development.md       # Development guide
│   ├── DOCUMENTATION_UPDATES.md # Documentation updates summary
│   └── VERSION_UPDATE_SUMMARY.md # Version update summary
├── CONTRIBUTING.md           # Contribution guidelines
├── LICENSE                   # MIT License
└── .goreleaser.yml          # Release automation
```

## Development

### Prerequisites

- Go 1.21 or later
- Git

### Building from Source

```bash
git clone https://github.com/jashkahar/open-workbench-platform.git
cd open-workbench-platform
go mod tidy
go build -o om main.go
```

### Running in Development

```bash
# Run with TUI
go run main.go ui

# Run simple interactive mode
go run main.go

# Test the new init command
go run main.go init
```

### Testing

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

### Adding New Templates

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

## Architecture

### Core Components

1. **Main Application (`main.go`)**

   - Command-line argument parsing
   - Entry point for different modes (TUI, interactive, non-interactive)
   - Orchestrates the scaffolding process

2. **Command System (`cmd/`)**

   - **Root Command**: Main CLI setup with Cobra framework
   - **Init Command**: Project initialization with `workbench.yaml` manifests
   - **Security**: Comprehensive security utilities and validation
   - **Types**: YAML manifest type definitions

3. **Terminal User Interface (`tui.go`)**

   - Beautiful interactive template selection
   - Uses Bubble Tea for TUI framework
   - Integrates with template discovery system

4. **Dynamic Templating System (`internal/templating/`)**

   - **Discovery**: Template discovery and validation
   - **Parameters**: Parameter collection, validation, and processing
   - **Processor**: Template processing and file operations

5. **Security System (`cmd/security.go`)**
   - Input validation and sanitization
   - Path traversal protection
   - Malicious pattern detection
   - Cross-platform security checks

### Data Flow

1. **Template Discovery**: CLI discovers available templates from embedded filesystem
2. **Template Selection**: User selects template via TUI or command line
3. **Parameter Collection**: Dynamic parameter collection with validation
4. **Template Processing**: Files are processed with parameter substitution
5. **Post-Scaffolding**: Conditional file deletion and command execution
6. **Security Validation**: All inputs validated for security threats

## Release Automation

This project uses [GoReleaser](https://goreleaser.com/) for automated releases:

- **Multi-platform builds**: Windows, macOS, and Linux (AMD64 and ARM64)
- **Package managers**: Homebrew and Scoop support
- **GitHub Releases**: Automatic release creation with changelog
- **Checksums**: SHA256 checksums for all binaries

### Release Process

1. Create and push a new tag: `git tag v0.5.0 && git push origin v0.5.0`
2. GitHub Actions automatically builds and releases
3. Binaries are available on GitHub Releases
4. Homebrew and Scoop packages are updated automatically

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new features (aim for 100% coverage)
- Update documentation for any changes
- Use conventional commit messages
- Ensure security validation for all user inputs

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/jashkahar/open-workbench-platform/issues) page
2. Create a new issue with detailed information
3. Join our community discussions

## Roadmap

### Completed ✅

- [x] Dynamic template system with conditional logic
- [x] Terminal User Interface (TUI)
- [x] Parameter validation and grouping
- [x] Post-scaffolding actions
- [x] Multiple template types (Next.js, FastAPI, React, Express, Vue)
- [x] Homebrew and Scoop installation support
- [x] Release automation with GoReleaser
- [x] Non-interactive CLI mode with command-line flags
- [x] Optional git initialization
- [x] Comprehensive error handling with help guidance
- [x] Template selection in interactive mode
- [x] **NEW**: `om init` command for project management
- [x] **NEW**: `workbench.yaml` manifest system
- [x] **NEW**: Enterprise-grade security features
- [x] **NEW**: Comprehensive testing suite (100% coverage)
- [x] **NEW**: Cross-platform security validation

### In Progress 🚧

- [ ] Template preview functionality
- [ ] Plugin system for custom templates

### Planned 📋

- [ ] CI/CD integration templates
- [ ] Cloud deployment templates
- [ ] Template validation and testing framework
- [ ] Template marketplace
- [ ] IDE integration plugins
- [ ] Template versioning and updates
- [ ] Advanced security audit features
- [ ] Security compliance reporting

---

**Maintainer**: Jash Kahar  
**Last Updated**: February 8, 2025  
**Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.
