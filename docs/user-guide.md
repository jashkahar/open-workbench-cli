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

# Add services with smart mode detection
om add service                    # Interactive mode
om add service --name frontend --template react-typescript --params ProjectName=my-app,Owner=John  # Direct mode

# List available templates
om list-templates

# Start the interactive TUI
om ui

# Or use simple interactive mode
om
```

## 🎯 Using the CLI

### Available Commands

| Command             | Description                   | Mode               |
| ------------------- | ----------------------------- | ------------------ |
| `om init`           | Initialize new project        | Interactive        |
| `om add service`    | Add service (smart detection) | Interactive/Direct |
| `om list-templates` | List available templates      | Direct             |
| `om`                | Interactive mode              | Interactive        |

### Smart Command System

The CLI features intelligent mode detection that automatically adapts to your needs:

#### Interactive Mode

```bash
om add service
```

- Prompts for all details interactively
- Perfect for exploration and learning
- Guided parameter collection with validation

#### Direct Mode

```bash
om add service --name frontend --template react-typescript --params ProjectName=my-app,Owner=John,IncludeTesting=true
```

- Uses all provided parameters
- Perfect for automation and scripting
- No prompts for provided parameters

#### Partial Direct Mode

```bash
om add service --name backend --template fastapi-basic
```

- Uses provided parameters
- Prompts only for missing parameters
- Best of both worlds

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

### Smart Service Addition

The `om add service` command intelligently adapts to your needs:

#### Interactive Mode

```bash
om add service
```

**Features:**

- Template selection from all available templates
- Organized parameter collection with grouping
- Comprehensive validation and error handling
- Optional git initialization and dependency installation

#### Direct Mode

```bash
om add service --name frontend --template react-typescript --params ProjectName=my-app,Owner=John,IncludeTesting=true,IncludeTailwind=true
```

**Features:**

- All parameters specified via command line
- Perfect for automation and CI/CD pipelines
- No interactive prompts
- Fast execution

#### Partial Direct Mode

```bash
om add service --name backend --template fastapi-basic
```

**Features:**

- Uses provided parameters
- Prompts only for missing parameters
- Best balance of speed and flexibility

### Template Discovery

Use `om list-templates` to explore available templates:

```bash
om list-templates
```

**Output includes:**

- Template names and descriptions
- Available parameters for each template
- Parameter types and default values
- Usage examples

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

The CLI performs comprehensive security checks:

- **Path Traversal Protection**: Blocks `../` and `..\` attacks
- **Malicious Pattern Detection**: Prevents JavaScript injection, command injection
- **Cross-Platform Security**: Windows reserved names, absolute path prevention
- **Directory Safety Checks**: Validates permissions, accessibility, symbolic links
- **Template Security**: Secure template name validation and content verification

## 📋 Template System

### Available Templates

#### 🎨 nextjs-full-stack

A production-ready Next.js application with TypeScript, testing, Tailwind CSS, Docker, and CI/CD setup.

#### ⚡ fastapi-basic

A FastAPI backend template with Python best practices, virtual environments, and API documentation.

#### 🎯 react-typescript

A modern React application with Vite, TypeScript, and modern tooling.

#### 🚀 express-api

A Node.js Express API template with TypeScript, testing, and documentation.

#### 🟢 vue-nuxt

A Vue.js Nuxt application with TypeScript, auto-imports, and SSR configuration.

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

## 🛠️ Advanced Usage

### Working with Projects

#### Adding Multiple Services

```bash
# Add a frontend service
om add service --name frontend --template react-typescript

# Add a backend service
om add service --name backend --template fastapi-basic

# Add an API service
om add service --name api --template express-api
```

#### Project Structure Example

```
my-project/
├── workbench.yaml
├── frontend/          # React TypeScript frontend
├── backend/           # FastAPI backend
└── api/              # Express API
```

### Automation Examples

#### CI/CD Pipeline

```bash
# Automated project creation
om init --name my-app --template nextjs-full-stack --service frontend

# Automated service addition
om add service --name backend --template fastapi-basic --params ProjectName=my-app,Owner=CI,IncludeTesting=true,IncludeDocker=true
```

#### Scripting

```bash
#!/bin/bash
# Create a full-stack project
om init
om add service --name frontend --template react-typescript --params ProjectName=my-app,Owner=Dev,IncludeTesting=true
om add service --name backend --template fastapi-basic --params ProjectName=my-app,Owner=Dev,IncludeTesting=true
```

## 🐛 Troubleshooting

### Common Issues

#### Permission Errors

```bash
# Ensure the directory is writable
chmod 755 my-project
```

#### Template Not Found

```bash
# List available templates
om list-templates

# Check template name spelling
om add service --template nextjs-full-stack
```

#### Parameter Validation Errors

```bash
# Use valid project names (lowercase, hyphens only)
om add service --name my-project --template react-typescript

# Avoid special characters and reserved names
# ❌ Don't use: my_project, MyProject, con, aux
# ✅ Use: my-project, project123, frontend
```

### Getting Help

```bash
# General help
om --help

# Command-specific help
om init --help
om add service --help
om list-templates --help
```

## 📚 Examples

### Quick Start Examples

#### Create a Full-Stack Project

```bash
# Initialize project
om init

# Add frontend
om add service --name frontend --template react-typescript

# Add backend
om add service --name backend --template fastapi-basic
```

#### Create a Monorepo

```bash
# Initialize project
om init

# Add multiple services
om add service --name web --template nextjs-full-stack
om add service --name api --template express-api
om add service --name admin --template react-typescript
```

#### Automation Example

```bash
# Create project with all parameters
om init
om add service --name frontend --template react-typescript --params ProjectName=my-app,Owner=Team,IncludeTesting=true,IncludeTailwind=true,IncludeDocker=true,InstallDeps=true,InitGit=true
om add service --name backend --template fastapi-basic --params ProjectName=my-app,Owner=Team,IncludeTesting=true,IncludeDocker=true,InstallDeps=true,InitGit=true
```

## 🎯 Best Practices

### Project Organization

1. **Use descriptive service names**: `frontend`, `backend`, `api`, `admin`
2. **Follow naming conventions**: lowercase with hyphens
3. **Group related services**: Keep frontend and backend in same project
4. **Use workbench.yaml**: Let the manifest manage your project structure

### Security

1. **Validate inputs**: Always use the CLI's built-in validation
2. **Check permissions**: Ensure directories are writable
3. **Review generated code**: Always review scaffolded code before deployment
4. **Use secure names**: Avoid special characters and reserved names

### Development Workflow

1. **Start with `om init`**: Create a managed project
2. **Add services incrementally**: Use `om add service` for each component
3. **Use interactive mode for exploration**: Learn templates with `om add service`
4. **Use direct mode for automation**: Script with parameters
5. **List templates regularly**: Use `om list-templates` to discover options

---

**Next Steps**: Check out the [Architecture Guide](architecture.md) for technical details or the [Development Guide](development.md) for contributing to the project.
