# Architecture Overview

This document provides a comprehensive overview of the Open Workbench CLI architecture, including system design, components, data flow, and technical decisions.

## 🏗️ System Architecture

### High-Level Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Input    │    │   Command       │    │   Template      │    │   Output        │
│                 │    │   System        │    │   Processing    │    │                 │
│ • CLI Args      │───▶│ • Cobra        │───▶│ • Discovery     │───▶│ • Project       │
│ • Interactive   │    │ • Security      │    │ • Parameters    │    │ • Files         │
│ • Smart Mode    │    │ • Validation    │    │ • Processing    │    │ • Docker        │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Core Components

#### 1. Command System (`cmd/`)

**Purpose**: Modern CLI framework with security, testing, and smart mode detection

**Responsibilities**:

- Provide structured command hierarchy using Cobra
- Implement comprehensive security validation for all inputs
- Handle smart mode detection (interactive/direct/partial)
- Manage project initialization with manifests
- Generate Docker Compose configurations

**Key Components**:

- **Root Command** (`cmd/root.go`): Main CLI setup and command registration
- **Init Command** (`cmd/init.go`): Project initialization with `workbench.yaml` manifests
- **Add Service** (`cmd/add_service.go`): Smart service addition with mode detection
- **Compose Command** (`cmd/compose.go`): Docker Compose generation
- **Security** (`cmd/security.go`): Comprehensive security utilities and validation
- **Types** (`cmd/types.go`): YAML manifest type definitions

**Key Functions**:

- `Execute()`: Main CLI execution with embedded filesystem
- `runInit()`: Project initialization workflow
- `runAddService()`: Smart service addition with mode detection
- `runCompose()`: Docker Compose generation
- `ValidateAndSanitizeName()`: Input validation and sanitization
- `ValidateAndSanitizePath()`: Path security validation
- `CheckForSuspiciousPatterns()`: Malicious pattern detection

**Smart Command System**:

- **Mode Detection**: Automatically switches between interactive and direct modes
- **Interactive Mode**: No parameters → prompts for all details
- **Direct Mode**: All parameters provided → uses provided parameters
- **Partial Mode**: Some parameters → uses provided, prompts for missing

**Security Features**:

- Path traversal protection (`../` and `..\` attacks)
- Malicious pattern detection (JavaScript injection, command injection)
- Cross-platform security (Windows reserved names, absolute paths)
- Directory safety checks (permissions, accessibility, symbolic links)
- Template security validation

#### 2. Main Application (`main.go`)

**Purpose**: Entry point and embedded filesystem management

**Responsibilities**:

- Initialize embedded filesystem for templates
- Route to command system
- Handle application lifecycle

**Key Functions**:

- `main()`: Application entry point with embedded filesystem
- `embed` directive: Embed templates into binary

#### 3. Compose System (`internal/compose/`)

**Purpose**: Docker Compose generation and orchestration

**Responsibilities**:

- Parse `workbench.yaml` project manifests
- Generate production-ready Docker Compose configurations
- Create environment files with secure defaults
- Validate Docker prerequisites

**Key Components**:

- **Generator** (`generator.go`): Docker Compose configuration generation
- **Prerequisites** (`prerequisites.go`): Docker environment validation
- **Types** (`types.go`): Compose-specific type definitions

**Key Functions**:

- `Generate()`: Create complete docker-compose.yml
- `GenerateEnvFile()`: Generate environment variables
- `CheckAllPrerequisites()`: Validate Docker installation
- `SaveDockerCompose()`: Write configuration files

**Features**:

- Service networking with proper isolation
- Environment variable management with secure defaults
- Volume mounting for development and production
- Health checks for service monitoring
- Multi-stage builds for optimized images

#### 4. Templating System (`internal/templating/`)

**Purpose**: Dynamic template processing with conditional logic

**Responsibilities**:

- Template discovery and validation
- Parameter collection and validation
- Conditional file generation
- Post-scaffolding actions

**Key Components**:

- **Discovery** (`discovery.go`): Template discovery and validation
- **Parameters** (`parameters.go`): Parameter processing and validation
- **Processor** (`processor.go`): Template processing and file operations
- **Progress** (`progress.go`): Progress tracking and user feedback

**Key Functions**:

- `DiscoverTemplates()`: Find and validate available templates
- `CollectParameters()`: Interactive parameter collection
- `ProcessTemplate()`: Template processing with conditional logic
- `ExecutePostScaffold()`: Post-scaffolding actions

**Advanced Features**:

- Conditional parameter display based on other parameters
- Conditional file generation and deletion
- Post-scaffolding commands and actions
- Parameter validation with regex patterns
- Parameter grouping for better UX

## 🔒 Security Architecture

### Input Validation System

**Purpose**: Comprehensive security validation for all user inputs

**Components**:

- **Path Validation**: Prevents path traversal attacks
- **Name Validation**: Validates project and service names
- **Pattern Detection**: Identifies malicious patterns
- **Cross-Platform Security**: Handles Windows and Unix security

**Validation Rules**:

```go
// Path traversal protection
if strings.Contains(name, "../") || strings.Contains(name, "..\\") {
    return errors.New("path traversal not allowed")
}

// Malicious pattern detection
suspiciousPatterns := []string{
    "javascript:", "data:", "vbscript:", "onload=", "onerror=",
    "eval(", "setTimeout(", "setInterval(", "document.cookie",
}

// Windows reserved names
windowsReserved := []string{"con", "prn", "aux", "nul", "com1", "com2"}
```

### Directory Safety System

**Purpose**: Ensure safe directory operations

**Features**:

- Directory permission validation
- Symbolic link detection
- Accessibility checks
- Empty directory validation

**Safety Checks**:

```go
func ValidateDirectorySafety(path string) error {
    // Check if directory exists and is accessible
    // Validate permissions
    // Check for symbolic links
    // Ensure directory is writable
}
```

### Template Security

**Purpose**: Validate template integrity and prevent malicious templates

**Features**:

- Template name validation
- Template content verification
- Parameter validation
- File operation safety

## 🧪 Testing Architecture

### Test Coverage

**Current Coverage**: 100% for security components, comprehensive for core functionality

**Test Categories**:

- **Unit Tests**: Individual function testing
- **Integration Tests**: Command system testing
- **Security Tests**: Security validation testing
- **Template Tests**: Template processing testing

### Test Structure

```
cmd/
├── security_test.go     # Security tests (100% coverage)
├── init_test.go         # Init command tests
└── compose_test.go      # Compose command tests

internal/
├── templating/
│   └── processor_test.go # Template processing tests
└── compose/
    └── generator_test.go # Compose generation tests
```

### Security Testing

**Comprehensive Security Test Suite**:

- Path traversal attack prevention
- Malicious pattern detection
- Cross-platform security validation
- Directory safety testing
- Template security validation

## 📊 Data Flow

### Project Initialization Flow

```
1. User runs 'om init'
   ↓
2. Directory safety check
   ↓
3. Project name validation
   ↓
4. Template selection
   ↓
5. Service name validation
   ↓
6. Parameter collection
   ↓
7. Template processing
   ↓
8. Manifest creation
   ↓
9. Success feedback
```

### Smart Service Addition Flow

```
1. User runs 'om add service'
   ↓
2. Mode detection (interactive/direct/partial)
   ↓
3. Parameter collection/validation
   ↓
4. Template validation
   ↓
5. Safety checks
   ↓
6. Template processing
   ↓
7. Manifest update
   ↓
8. Success feedback
```

### Docker Compose Generation Flow

```
1. User runs 'om compose'
   ↓
2. Prerequisite checking
   ↓
3. Manifest parsing
   ↓
4. Service analysis
   ↓
5. Docker Compose generation
   ↓
6. Environment file creation
   ↓
7. Gitignore update
   ↓
8. Success feedback
```

## 🔧 Configuration Management

### Project Manifest (`workbench.yaml`)

**Purpose**: Central project configuration and service management

**Structure**:

```yaml
apiVersion: openworkbench.io/v1alpha1
kind: Project
metadata:
  name: my-project
services:
  frontend:
    template: nextjs-full-stack
    path: ./frontend
  backend:
    template: fastapi-basic
    path: ./backend
components:
  gateway:
    template: nginx-gateway
    path: ./gateway
```

### Template Manifest (`template.json`)

**Purpose**: Template configuration and parameter definitions

**Structure**:

```json
{
  "name": "Template Display Name",
  "description": "Template description",
  "parameters": [
    {
      "name": "ProjectName",
      "type": "string",
      "required": true,
      "validation": {
        "regex": "^[a-z0-9-]+$"
      }
    }
  ],
  "postScaffold": {
    "commands": [
      {
        "command": "npm install",
        "condition": "InstallDeps == true"
      }
    ]
  }
}
```

## 🚀 Performance Considerations

### Embedded Filesystem

**Benefits**:

- Single binary distribution
- No external template dependencies
- Faster template loading
- Reduced deployment complexity

**Implementation**:

```go
//go:embed templates
var templatesFS embed.FS
```

### Smart Mode Detection

**Benefits**:

- Reduced user interaction for experienced users
- Maintained simplicity for new users
- Flexible automation support
- Improved user experience

### Template Processing

**Optimizations**:

- Lazy template discovery
- Cached parameter validation
- Efficient file operations
- Progress tracking for large templates

## 🔄 Error Handling

### Comprehensive Error System

**Error Categories**:

- **Validation Errors**: Input validation failures
- **Security Errors**: Security check failures
- **Template Errors**: Template processing failures
- **System Errors**: File system and permission errors

**Error Handling Strategy**:

- Clear, actionable error messages
- Contextual help and suggestions
- Graceful degradation
- Comprehensive logging

### User-Friendly Error Messages

**Examples**:

```bash
# Clear validation error
❌ Invalid project name: "my project"
   Project names can only contain lowercase letters, numbers, and hyphens.
   Try: "my-project"

# Security error
❌ Security check failed: path traversal detected
   Project names cannot contain "../" or "..\"

# Template error
❌ Template "invalid-template" not found
   Available templates: nextjs-full-stack, react-typescript, fastapi-basic
```

## 🔮 Future Architecture

### Planned Enhancements

1. **Plugin System**: Extensible template and command system
2. **Cloud Integration**: Direct deployment to cloud platforms
3. **Advanced Orchestration**: Kubernetes and Docker Swarm support
4. **Template Marketplace**: Community template sharing
5. **Advanced Security**: Additional security layers and compliance

### Scalability Considerations

1. **Modular Design**: Easy to add new commands and features
2. **Template Ecosystem**: Extensible template system
3. **Cloud-Native**: Ready for cloud deployment
4. **Enterprise Features**: Security and compliance ready
