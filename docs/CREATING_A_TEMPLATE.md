# Creating Templates

This guide explains how to create and contribute templates for the Open Workbench Platform.

## Template Types

Open Workbench supports three types of templates:

### 1. Service Templates

- **Purpose**: Individual application services (frontend, backend, API)
- **Examples**: Next.js app, FastAPI backend, React frontend
- **Location**: `templates/[template-name]/`

### 2. Component Templates

- **Purpose**: Shared infrastructure components (gateways, load balancers)
- **Examples**: Nginx gateway, Redis cache, Load balancer
- **Location**: `templates/[template-name]/`

### 3. Resource Templates

- **Purpose**: Service-specific resources (databases, storage)
- **Examples**: PostgreSQL database, Redis cache, S3 bucket
- **Location**: `templates/[template-name]/`

## Template Structure

Each template follows this directory structure:

```
template-name/
├── template.json         # Template manifest
├── README.md             # Template documentation
├── src/                  # Source files (optional)
├── tests/                # Test files (optional)
├── Dockerfile            # Container configuration (optional)
├── package.json          # Node.js dependencies (optional)
├── requirements.txt      # Python dependencies (optional)
└── [other files]         # Template-specific files
```

## Template Manifest (`template.json`)

The `template.json` file defines the template's behavior and parameters.

### Basic Structure

```json
{
  "name": "Template Display Name",
  "description": "Template description",
  "parameters": [
    {
      "name": "ParameterName",
      "prompt": "User prompt text",
      "group": "Parameter Group",
      "type": "string|boolean|select|multiselect",
      "required": true,
      "default": "default_value",
      "validation": {
        "regex": "^[a-z0-9-]+$",
        "errorMessage": "Validation error message"
      },
      "condition": "ConditionExpression",
      "options": ["option1", "option2"],
      "helpText": "Help text for the parameter"
    }
  ],
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "file-to-delete.txt",
        "condition": "DeleteCondition"
      }
    ],
    "commands": [
      {
        "command": "npm install",
        "description": "Installing dependencies...",
        "condition": "InstallCondition"
      }
    ]
  }
}
```

### Parameter Types

#### String Parameters

```json
{
  "name": "ProjectName",
  "prompt": "Project Name:",
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
  "prompt": "Include testing framework?",
  "type": "boolean",
  "default": true
}
```

#### Select Parameters

```json
{
  "name": "Framework",
  "prompt": "Select framework:",
  "type": "select",
  "default": "React",
  "options": ["React", "Vue", "Angular"],
  "condition": "IncludeFrontend == true"
}
```

#### Multiselect Parameters

```json
{
  "name": "Features",
  "prompt": "Select features:",
  "type": "multiselect",
  "options": ["Authentication", "Database", "API", "Testing"]
}
```

### Conditional Logic

Parameters can be conditionally shown based on other parameters:

```json
{
  "name": "TestingFramework",
  "prompt": "Which testing framework?",
  "type": "select",
  "options": ["Jest", "Vitest"],
  "condition": "IncludeTesting == true"
}
```

### Post-Scaffold Actions

#### File Deletions

```json
{
  "filesToDelete": [
    {
      "path": "jest.config.js",
      "condition": "TestingFramework != 'Jest'"
    },
    {
      "path": "tests/",
      "condition": "IncludeTesting == false"
    }
  ]
}
```

#### Commands

```json
{
  "commands": [
    {
      "command": "npm install",
      "description": "Installing dependencies...",
      "condition": "InstallDeps == true"
    },
    {
      "command": "git init",
      "description": "Initializing Git repository...",
      "condition": "InitGit == true"
    }
  ]
}
```

## Template Files

### Go Template Syntax

Template files use Go template syntax for dynamic content:

#### Basic Variable Substitution

```go
{{ .ProjectName }}
{{ .Owner }}
```

#### Conditional Logic

```go
{{ if .IncludeTesting }}
// Testing configuration
{{ end }}

{{ if eq .Framework "React" }}
// React-specific code
{{ else if eq .Framework "Vue" }}
// Vue-specific code
{{ end }}
```

#### Loops

```go
{{ range .Features }}
// Feature: {{ . }}
{{ end }}
```

#### Custom Functions

```go
{{ .ProjectName | ToLower }}
{{ .Description | Truncate 100 }}
```

### File Naming

Files can have dynamic names using template syntax:

```
{{ .ProjectName }}/src/{{ .Framework | ToLower }}/index.{{ if eq .Framework "React" }}tsx{{ else }}vue{{ end }}
```

### Conditional Files

Files can be conditionally included by using empty names:

```
{{ if .IncludeDocker }}Dockerfile{{ end }}
{{ if .IncludeTesting }}tests/{{ end }}
```

## Workbench.yaml Schema

The `workbench.yaml` file is automatically generated and updated by the system:

```yaml
apiVersion: openworkbench.io/v1alpha1
kind: Project
metadata:
  name: project-name
environments:
  dev:
    provider: aws
    region: us-west-2
services:
  frontend:
    template: nextjs-full-stack
    path: ./frontend
    port: 3000
    environment:
      NODE_ENV: development
  backend:
    template: fastapi-basic
    path: ./backend
    port: 8000
    resources:
      database:
        type: postgresql
        version: "14"
components:
  gateway:
    template: nginx-gateway
    path: ./gateway
    ports: ["80", "443"]
```

### Schema Fields

#### Metadata

- `name`: Project name (required)

#### Services

- `template`: Template name used
- `path`: Service directory path
- `port`: Service port (optional)
- `environment`: Environment variables (optional)
- `resources`: Service-specific resources (optional)

#### Components

- `template`: Template name used
- `path`: Component directory path
- `ports`: List of ports (optional)

#### Environments

- `provider`: Cloud provider (aws, gcp, azure)
- `region`: Cloud region
- `config`: Environment-specific configuration

## Testing Your Template

### Local Testing

1. **Create your template** in the `templates/` directory
2. **Build the CLI**:
   ```bash
   go build -o bin/om main.go
   ```
3. **Test initialization**:
   ```bash
   ./bin/om init
   ```
4. **Test service addition**:
   ```bash
   ./bin/om add service --template your-template-name
   ```

### Template Validation

The system validates templates automatically:

- **Parameter validation**: Ensures all required parameters are defined
- **File existence**: Checks that referenced files exist
- **Syntax validation**: Validates Go template syntax
- **Condition validation**: Ensures conditional logic is valid

### Common Issues

#### Template Not Found

- Ensure template directory exists in `templates/`
- Check template name spelling
- Verify `template.json` exists

#### Parameter Errors

- Ensure all required parameters are provided
- Check parameter validation rules
- Verify parameter types match expectations

#### File Generation Errors

- Check Go template syntax
- Ensure file paths are valid
- Verify conditional logic

## Best Practices

### Template Design

1. **Clear Naming**: Use descriptive template and parameter names
2. **Logical Grouping**: Group related parameters together
3. **Sensible Defaults**: Provide useful default values
4. **Comprehensive Documentation**: Include clear README files

### Parameter Design

1. **Minimal Required Parameters**: Only require essential parameters
2. **Clear Prompts**: Use descriptive prompt text
3. **Helpful Validation**: Provide clear validation error messages
4. **Logical Conditions**: Use conditional logic sparingly

### File Organization

1. **Standard Structure**: Follow the standard template structure
2. **Clear Documentation**: Include README files
3. **Proper Dependencies**: Include all necessary dependency files
4. **Testing Support**: Include test files when appropriate

### Security Considerations

1. **Input Validation**: Validate all user inputs
2. **Path Safety**: Ensure file paths are safe
3. **Command Safety**: Be careful with post-scaffold commands
4. **Template Security**: Validate template content

## Contributing Templates

### Submission Process

1. **Create your template** following the guidelines above
2. **Test thoroughly** using local builds
3. **Document clearly** with README files
4. **Submit pull request** with clear description

### Template Review

Templates are reviewed for:

- **Functionality**: Does the template work correctly?
- **Documentation**: Is the template well-documented?
- **Best Practices**: Does it follow established patterns?
- **Security**: Are there any security concerns?

### Template Maintenance

- **Keep templates updated** with framework versions
- **Test regularly** with new CLI versions
- **Respond to issues** and user feedback
- **Update documentation** as needed

## Example Templates

### Next.js Full Stack Template

```json
{
  "name": "Next.js Full Stack",
  "description": "A fully-featured Next.js application with testing and deployment.",
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
      "name": "IncludeTesting",
      "prompt": "Include testing framework?",
      "type": "boolean",
      "default": true
    },
    {
      "name": "IncludeDocker",
      "prompt": "Include Docker configuration?",
      "type": "boolean",
      "default": true
    }
  ],
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "tests/",
        "condition": "IncludeTesting == false"
      },
      {
        "path": "Dockerfile",
        "condition": "IncludeDocker == false"
      }
    ],
    "commands": [
      {
        "command": "npm install",
        "description": "Installing dependencies...",
        "condition": "true"
      }
    ]
  }
}
```

### FastAPI Basic Template

```json
{
  "name": "FastAPI Basic",
  "description": "A FastAPI backend with automatic API documentation.",
  "parameters": [
    {
      "name": "ProjectName",
      "prompt": "Project Name:",
      "type": "string",
      "required": true
    },
    {
      "name": "IncludeDatabase",
      "prompt": "Include database configuration?",
      "type": "boolean",
      "default": true
    }
  ],
  "postScaffold": {
    "commands": [
      {
        "command": "pip install -r requirements.txt",
        "description": "Installing Python dependencies...",
        "condition": "true"
      }
    ]
  }
}
```

## Resources

- **[System Architecture](ARCHITECTURE.md)**: Understanding the template system
- **[Contributing Guide](../CONTRIBUTING.md)**: General contribution guidelines
- **[Available Templates](../templates/)**: Existing template examples
