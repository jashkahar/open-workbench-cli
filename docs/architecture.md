# System Architecture

This document provides an in-depth look at the Open Workbench Platform's architecture, core logic, and design patterns.

## High-Level System Overview

Open Workbench Platform is a CLI tool that automates the creation of multi-service applications. The system consists of several key components:

- **Command Layer** (`cmd/`): CLI commands and user interaction
- **Templating Engine** (`internal/templating/`): Dynamic template processing
- **Manifest System** (`internal/manifest/`): Project configuration management
- **Generator System** (`internal/generator/`): Deployment configuration generation
- **Security Layer** (`cmd/security.go`): Input validation and sanitization

## Core Components

### Command Layer (`cmd/`)

The command layer is built using the Cobra framework and provides the following commands:

#### `om init`
- **Purpose**: Initialize a new Open Workbench project
- **Process**: 
  1. Validates directory safety
  2. Prompts for project details
  3. Creates project structure
  4. Generates `workbench.yaml` manifest
- **Key Files**: `cmd/init.go`

#### `om add service`
- **Purpose**: Add a new service to an existing project
- **Modes**: Interactive and direct (with flags)
- **Process**:
  1. Loads existing `workbench.yaml`
  2. Validates service name uniqueness
  3. Scaffolds service using template
  4. Updates manifest file
- **Key Files**: `cmd/add_service.go`

#### `om add component`
- **Purpose**: Add shared infrastructure components
- **Process**: Similar to add service but for components
- **Key Files**: `cmd/add_service.go` (shared logic)

#### `om add resource`
- **Purpose**: Add an infrastructure resource (database, cache, storage, MQ) to a service
- **Process**:
  1. Loads existing `workbench.yaml`
  2. Collects resource type and parameters (interactive/direct)
  3. Updates manifest file under the selected service
- **Key Files**: `cmd/add_resource.go`, `internal/resources` (blueprints)

#### `om compose`
- **Purpose**: Generate deployment configurations
- **Targets**: Docker Compose (Terraform prototype is currently disabled)
- **Process**:
  1. Loads `workbench.yaml`
  2. Selects target (docker)
  3. Generates configuration files
- **Key Files**: `cmd/compose.go`

#### `om ls`
- **Purpose**: List project services and components
- **Process**: Reads and displays `workbench.yaml` contents
- **Key Files**: `cmd/ls.go`

#### `om delete`
- **Purpose**: Remove services or components
- **Process**: Updates manifest and removes files
- **Key Files**: `cmd/delete.go`

### Templating Engine (`internal/templating/`)

The templating engine is the core of the system, providing dynamic template processing with conditional logic.

#### Key Components

**TemplateProcessor** (`processor.go`):
- Processes template files with Go template syntax
- Handles conditional logic (`{{ if .Condition }}`)
- Manages post-scaffolding actions
- Provides progress reporting

**Parameter System** (`parameters.go`):
- Defines parameter types (string, boolean, select, multiselect)
- Handles parameter validation and grouping
- Manages interactive prompting

**Discovery System** (`discovery.go`):
- Discovers available templates
- Loads template manifests (`template.json`)
- Validates template structure

#### Template Processing Flow

1. **Discovery**: Find available templates in embedded filesystem
2. **Parameter Collection**: Interactive or direct parameter gathering
3. **Validation**: Validate parameters against template requirements
4. **Processing**: Apply template with collected parameters
5. **Post-Actions**: Execute conditional file deletions and commands

### Manifest System (`internal/manifest/`)

The manifest system manages project configuration through `workbench.yaml` files.

#### WorkbenchManifest Structure

```yaml
apiVersion: openworkbench.io/v1alpha1
kind: Project
metadata:
  name: project-name
services:
  frontend:
    template: nextjs-full-stack
    path: ./frontend
    port: 3000
  backend:
    template: fastapi-basic
    path: ./backend
    port: 8000
components:
  gateway:
    template: nginx-gateway
    path: ./gateway
    ports: ["8080:80"]
```

#### Key Features

- **Service Management**: Track all services in the project
- **Component Support**: Shared infrastructure components
- **Environment Configuration**: Multi-environment deployment support
- **Resource Tracking**: Service-specific resources (databases, etc.)

### Generator System (`internal/generator/`)

The generator system creates deployment configurations from the manifest.

#### Docker Generator (`generator/docker/`)
- Generates `docker-compose.yml` files
- Handles service networking
- Manages environment variables
- Supports volume mounts

#### Terraform Generator (`generator/terraform/`)
- (Temporarily disabled) Future support for generating Terraform configurations

### Security Layer (`cmd/security.go`)

The security layer provides input validation and sanitization.

#### Key Features

- **Directory Safety**: Prevents overwriting existing projects
- **Parameter Validation**: Validates user inputs against schemas
- **Path Sanitization**: Prevents path traversal attacks
- **Template Validation**: Ensures template integrity

## Command Reference

### `om init`

Initialize a new Open Workbench project.

**Flags:**
- None (interactive mode only)

**Process:**
1. Validates current directory is empty or contains only hidden files
2. Prompts for project name and first service
3. Creates project structure
4. Generates initial `workbench.yaml`

### `om add service`

Add a new service to the project.

**Flags:**
- `--name`: Service name (optional)
- `--template`: Template name (optional)
- `--params`: Key-value parameters (optional)

**Modes:**
- **Interactive**: No flags provided, prompts for all details
- **Direct**: Flags provided, minimal prompting

### `om add component`

Add a shared component to the project.

**Flags:**
- `--name`: Component name (optional)
- `--template`: Template name (optional)
- `--params`: Key-value parameters (optional)

### `om compose`

Generate deployment configuration.

**Flags:**
- `--target`: Deployment target (docker)
- `--env`: Environment name (reserved for Terraform)

### `om ls`

List project services and components.

**Flags:**
- `--detailed`: Show detailed information including paths, ports, env vars, and resource configs

### `om delete`

Remove services, components, or resources.

Subcommands and flags:
- `om delete service [name]` — remove a service from `workbench.yaml`
  - Flags: `--files` (also delete the service directory and files)
- `om delete component [name]` — remove a component from `workbench.yaml`
  - Flags: `--files` (also delete the component directory and files)
- `om delete resource service.resource` — remove a resource from a service
  - Example: `om delete resource backend.database`

## Security Architecture

### Input Validation

The system implements multiple layers of input validation:

1. **Parameter Validation**: Each parameter has validation rules (regex, required, etc.)
2. **Directory Safety**: Prevents overwriting existing projects
3. **Path Sanitization**: Prevents path traversal attacks
4. **Template Validation**: Ensures template integrity

### Safety Checks

- **Directory Safety**: Validates target directories are safe for initialization
- **Service Uniqueness**: Ensures service names are unique within projects
- **Template Existence**: Validates templates exist before processing
- **Parameter Completeness**: Ensures all required parameters are provided

## Data Flow

### Project Initialization Flow

1. **User runs `om init`**
2. **Directory Safety Check**: Validates current directory
3. **Project Name Collection**: Prompts for project name
4. **First Service Setup**: Prompts for service details
5. **Directory Creation**: Creates project structure
6. **Service Scaffolding**: Processes template with parameters
7. **Manifest Creation**: Generates `workbench.yaml`
8. **Success Feedback**: Reports completion

### Service Addition Flow

1. **User runs `om add service`**
2. **Manifest Loading**: Loads existing `workbench.yaml`
3. **Service Details Collection**: Interactive or direct parameter collection
4. **Safety Validation**: Checks service name uniqueness
5. **Template Processing**: Scaffolds service using template
6. **Manifest Update**: Updates `workbench.yaml`
7. **Success Feedback**: Reports completion

### Template Processing Flow

1. **Template Discovery**: Loads template from embedded filesystem
2. **Parameter Collection**: Gathers parameters interactively or from flags
3. **Parameter Validation**: Validates against template requirements
4. **File Processing**: Processes each file in template
5. **Conditional Logic**: Applies conditional file deletions
6. **Post-Actions**: Executes post-scaffolding commands
7. **Progress Reporting**: Provides user feedback throughout

## Design Patterns

### Command Pattern
Each CLI command is implemented as a separate Cobra command with clear separation of concerns.

### Template Method Pattern
The templating engine uses a template method pattern for processing templates with customizable steps.

### Strategy Pattern
The generator system uses strategy pattern for different deployment targets (Docker, Terraform).

### Factory Pattern
The template processor uses factory pattern for creating different parameter types.

## Error Handling

The system implements comprehensive error handling:

1. **Graceful Degradation**: Continues processing when possible
2. **User-Friendly Messages**: Clear error messages for users
3. **Validation Errors**: Specific validation error messages
4. **Recovery Mechanisms**: Automatic cleanup on failures

## Performance Considerations

1. **Embedded Templates**: Templates are embedded in binary for fast access
2. **Lazy Loading**: Templates are loaded only when needed
3. **Minimal I/O**: Efficient file operations
4. **Memory Management**: Proper cleanup of resources

## Extension Points

The system is designed for extensibility:

1. **Template System**: Easy to add new templates
2. **Generator System**: Easy to add new deployment targets
3. **Parameter Types**: Easy to add new parameter types
4. **Command System**: Easy to add new commands 