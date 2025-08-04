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

# Generate Docker Compose configuration
om compose
```

## 🎯 Using the CLI

### Available Commands

| Command             | Description                   | Mode               |
| ------------------- | ----------------------------- | ------------------ |
| `om init`           | Initialize new project        | Interactive        |
| `om add service`    | Add service (smart detection) | Interactive/Direct |
| `om add component`  | Add infrastructure component  | Interactive/Direct |
| `om list-templates` | List available templates      | Direct             |
| `om compose`        | Generate Docker Compose       | Direct             |

### Smart Command System

The CLI features intelligent mode detection that automatically adapts to your needs:

#### Interactive Mode

```bash
om add service
```

- Prompts for all details interactively
- Perfect for exploration and learning
- Guided parameter collection with validation
- Clear help text and examples

#### Direct Mode

```bash
om add service --name frontend --template react-typescript --params ProjectName=my-app,Owner=John,IncludeTesting=true
```

- Uses all provided parameters
- Perfect for automation and scripting
- No prompts for provided parameters
- Fast execution for experienced users

#### Partial Direct Mode

```bash
om add service --name backend --template fastapi-basic
```

- Uses provided parameters
- Prompts only for missing parameters
- Best of both worlds
- Flexible for different use cases

## 🏗️ Project Management

### Project Initialization

The `om init` command creates a new project with a structured approach:

```bash
om init
```

This will:

1. Check directory safety (empty or hidden files only)
2. Prompt for project name and first service
3. Create project structure with `workbench.yaml` manifest
4. Scaffold the first service with all configurations

### Adding Services

Use the smart `om add service` command to add new services:

```bash
# Interactive mode
om add service

# Direct mode with all parameters
om add service --name frontend --template nextjs-full-stack --params ProjectName=my-app,Owner=John,IncludeTesting=true,IncludeDocker=true

# Partial mode - prompts for missing parameters
om add service --name backend --template fastapi-basic
```

### Adding Components

Add infrastructure components like gateways, caches, or load balancers:

```bash
# Interactive mode
om add component

# Direct mode
om add component --name gateway --template nginx-gateway
```

### Docker Compose Generation

Generate complete Docker Compose configurations from your project:

```bash
om compose
```

This will:

1. Check for Docker and Docker Compose prerequisites
2. Parse your `workbench.yaml` file
3. Generate `docker-compose.yml` with proper networking
4. Create `.env` and `.env.example` files
5. Update `.gitignore` for Docker files

## 📋 Template System

### Available Templates

#### Frontend Templates

- **nextjs-full-stack**: Complete Next.js application with TypeScript, testing, and Docker
- **react-typescript**: Modern React app with Vite, TypeScript, and Tailwind CSS
- **vue-nuxt**: Vue.js Nuxt application with SSR and PWA support

#### Backend Templates

- **fastapi-basic**: FastAPI backend with automatic API documentation
- **express-api**: Node.js Express API with TypeScript and authentication

#### Infrastructure Templates

- **nginx-gateway**: Nginx reverse proxy for microservices
- **redis-cache**: Redis cache service with configuration

### Template Parameters

Each template supports customizable parameters:

#### Common Parameters

- `ProjectName`: Project name (required)
- `Owner`: Project owner (required)
- `IncludeTesting`: Include testing framework (boolean)
- `IncludeDocker`: Include Docker configuration (boolean)
- `IncludeTailwind`: Include Tailwind CSS (boolean)

#### Template-Specific Parameters

- **nextjs-full-stack**: Testing framework choice (Jest/Vitest)
- **fastapi-basic**: Database choice (SQLite/PostgreSQL/MongoDB)
- **express-api**: Authentication type (JWT/Basic)

### Parameter Validation

The system includes comprehensive validation:

- **Regex patterns** for string validation
- **Required field checking** with clear error messages
- **Type validation** for booleans, selects, and multiselects
- **Cross-field validation** for dependent parameters

## 🔧 Advanced Features

### Post-Scaffolding Actions

Templates can perform automatic actions after creation:

```json
{
  "postScaffold": {
    "commands": [
      {
        "command": "npm install",
        "description": "Installing dependencies...",
        "condition": "InstallDeps == true"
      }
    ],
    "filesToDelete": [
      {
        "path": "unused-config.js",
        "condition": "IncludeFeature == false"
      }
    ]
  }
}
```

### Conditional Logic

Templates support conditional file generation and parameter display:

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
    }
  ]
}
```

### Security Features

The CLI includes comprehensive security measures:

- **Path traversal protection** against `../` attacks
- **Malicious pattern detection** for injection attacks
- **Cross-platform security** for Windows and Unix systems
- **Directory safety checks** for permissions and accessibility
- **Template validation** to prevent malicious templates

## 🐳 Docker Integration

### Docker Compose Generation

The `om compose` command generates production-ready Docker configurations:

```bash
om compose
```

Generated files:

- `docker-compose.yml`: Complete service orchestration
- `.env`: Environment variables with secure defaults
- `.env.example`: Template for environment configuration
- Updated `.gitignore`: Docker-specific exclusions

### Docker Features

- **Service networking** with proper isolation
- **Environment variable management** with secure defaults
- **Volume mounting** for development and production
- **Health checks** for service monitoring
- **Multi-stage builds** for optimized images

## 🔍 Troubleshooting

### Common Issues

#### Template Not Found

```bash
# List available templates
om list-templates

# Check template name spelling
om add service --name myapp --template nextjs-fullstack  # Wrong
om add service --name myapp --template nextjs-full-stack # Correct
```

#### Parameter Validation Errors

```bash
# Check parameter format
om add service --params ProjectName=my-app,Owner=John  # Correct
om add service --params "ProjectName=my app,Owner=John" # Wrong (spaces)
```

#### Docker Compose Issues

```bash
# Check prerequisites
om compose

# Verify Docker installation
docker --version
docker-compose --version
```

### Getting Help

```bash
# Command help
om --help
om add service --help
om compose --help

# Template information
om list-templates
```

## 📚 Examples

### Complete Project Setup

```bash
# 1. Initialize project
om init

# 2. Add frontend
om add service --name frontend --template nextjs-full-stack

# 3. Add backend
om add service --name backend --template fastapi-basic

# 4. Add cache
om add component --name cache --template redis-cache

# 5. Generate Docker Compose
om compose

# 6. Start the application
docker-compose up
```

### Microservices Setup

```bash
# Create API gateway
om add component --name gateway --template nginx-gateway

# Create multiple services
om add service --name user-service --template fastapi-basic
om add service --name product-service --template express-api
om add service --name frontend --template react-typescript

# Generate and run
om compose
docker-compose up
```

### Quick Development Setup

```bash
# Fast setup for development
om init
om add service --name app --template react-typescript --params ProjectName=dev-app,Owner=Developer,IncludeTesting=false,IncludeDocker=false
```

## 🚀 Best Practices

### Project Organization

1. **Use descriptive names** for services and components
2. **Group related services** in the same project
3. **Use the manifest file** (`workbench.yaml`) for project documentation
4. **Version control** your project structure

### Template Selection

1. **Start simple** with basic templates for learning
2. **Use production templates** for real projects
3. **Consider infrastructure** needs early in the project
4. **Plan for scaling** with microservices architecture

### Development Workflow

1. **Initialize project** with `om init`
2. **Add core services** first
3. **Add infrastructure** components as needed
4. **Generate Docker Compose** for deployment
5. **Customize configurations** for your specific needs
