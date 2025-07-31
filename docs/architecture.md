# Architecture Overview

This document provides a comprehensive overview of the Open Workbench CLI architecture, including system design, components, data flow, and technical decisions.

## 🏗️ System Architecture

### High-Level Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   User Input    │    │   Template      │    │   Output        │
│                 │    │   Processing    │    │                 │
│ • CLI Args      │───▶│ • Discovery     │───▶│ • Project       │
│ • TUI           │    │ • Parameters    │    │ • Files         │
│ • Interactive   │    │ • Processing    │    │ • Commands      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Core Components

#### 1. Main Application (`main.go`)

**Purpose**: Entry point and orchestration layer

**Responsibilities**:

- Parse command-line arguments
- Route to appropriate execution mode (TUI, interactive, non-interactive)
- Coordinate the scaffolding process
- Handle user interaction and error reporting

**Key Functions**:

- `main()`: Application entry point

- `runInteractiveScaffold()`: Interactive mode execution
- `runCLICreate()`: CLI mode execution
- `scaffoldAndApplyDynamic()`: Core scaffolding logic

#### 2. Terminal User Interface (`tui.go`)

**Purpose**: Interactive template selection interface

**Responsibilities**:

- Present available templates in a user-friendly interface
- Handle user navigation and selection
- Integrate with template discovery system
- Provide visual feedback during selection

**Key Components**:

- `model`: TUI state management
- `item`: Template representation in the list
- `runTUI()`: Main TUI execution function

**Dependencies**:

#### 3. CLI Mode (`main.go`)

**Purpose**: Non-interactive command-line interface

**Responsibilities**:

- Parse command-line arguments and flags
- Validate required parameters
- Execute scaffolding without user interaction
- Provide comprehensive help and error messages

**Key Functions**:

- `runCLICreate()`: Main CLI mode execution
- Flag parsing and validation
- Parameter value mapping from flags
- Error handling with help guidance

**Features**:

- Support for all template options via flags
- Optional git initialization (`--no-git`)
- Optional dependency installation (`--no-install-deps`)
- Conditional feature flags (`--no-testing`, `--no-tailwind`, etc.)
- Comprehensive help system with examples

#### 3. Dynamic Templating System (`internal/templating/`)

##### Discovery (`discovery.go`)

**Purpose**: Template discovery and validation

**Responsibilities**:

- Scan embedded filesystem for available templates
- Load and parse template manifests
- Validate template structure and parameters
- Provide template metadata

**Key Functions**:

- `DiscoverTemplates()`: Find all available templates
- `LoadTemplateManifest()`: Parse template.json files
- `ValidateTemplate()`: Validate template structure
- `GetTemplateInfo()`: Get specific template information

##### Parameters (`parameters.go`)

**Purpose**: Parameter collection and validation

**Responsibilities**:

- Collect user input for template parameters
- Validate parameter values against rules
- Handle conditional parameter visibility
- Group parameters for better UX

**Key Components**:

- `ParameterProcessor`: Main parameter processing logic
- `Parameter`: Parameter definition structure
- `Validation`: Validation rules structure

**Key Functions**:

- `GetVisibleParameters()`: Return parameters based on conditions
- `ValidateParameter()`: Validate parameter values
- `GetParameterGroups()`: Organize parameters by groups
- `evaluateCondition()`: Evaluate conditional logic

##### Processor (`processor.go`)

**Purpose**: Template processing and file operations

**Responsibilities**:

- Process template files with parameter substitution
- Handle conditional file inclusion/exclusion
- Execute post-scaffolding actions
- Manage file and directory operations

**Key Components**:

- `TemplateProcessor`: Main template processing logic
- `PostScaffold`: Post-scaffolding actions definition

**Key Functions**:

- `ScaffoldProject()`: Main scaffolding function
- `ProcessTemplate()`: Template content processing
- `ProcessFileName()`: Filename template processing
- `ExecutePostScaffoldActions()`: Execute post-scaffolding actions

## 🔄 Data Flow

### 1. Template Discovery Flow

```
1. Application Start
   ↓
2. DiscoverTemplates() scans embedded filesystem
   ↓
3. LoadTemplateManifest() for each template directory
   ↓
4. ValidateTemplate() checks structure and parameters
   ↓
5. Return TemplateInfo array
```

### 2. User Interaction Flow

```
1. User selects execution mode (TUI/Interactive)
   ↓
2. Template selection (TUI: visual, Interactive: default)
   ↓
3. Load template manifest
   ↓
4. Collect parameters with validation
   ↓
5. Process template with parameters
   ↓
6. Execute post-scaffolding actions
   ↓
7. Report success/failure
```

### 3. Template Processing Flow

```
1. TemplateProcessor created with manifest and values
   ↓
2. Walk template directory structure
   ↓
3. Process each file:
   - Process filename template
   - Process file content template
   - Write to destination
   ↓
4. Execute post-scaffolding actions:
   - Conditional file deletion
   - Command execution
```

## 📁 File Structure

### Embedded Filesystem

```
templates/
├── nextjs-full-stack/
│   ├── template.json          # Template manifest
│   ├── package.json           # Template file
│   ├── src/
│   │   └── app/
│   │       └── page.tsx       # Template file
│   └── Dockerfile             # Template file
├── fastapi-basic/
│   ├── template.json          # Template manifest
│   ├── main.py               # Template file
│   └── requirements.txt       # Template file
└── [other templates...]
```

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

## 🔧 Technical Decisions

### 1. Embedded Filesystem

**Decision**: Use Go's `embed` directive for template storage

**Rationale**:

- Single binary distribution
- No external file dependencies
- Version control for templates
- Cross-platform compatibility

**Implementation**:

```go
//go:embed templates
var templatesFS embed.FS
```

### 2. JSON-Based Template Manifests

**Decision**: Use JSON for template configuration

**Rationale**:

- Human-readable format
- Easy to edit and version control
- No compilation required
- Flexible structure for future extensions

### 3. Conditional Logic System

**Decision**: Simple string-based condition evaluation

**Rationale**:

- Lightweight implementation
- Easy to understand and debug
- Sufficient for current use cases
- Extensible for future enhancements

**Implementation**:

```go
// Simple equality conditions
"IncludeTesting == true"
"TestingFramework != 'Jest'"
```

### 4. Template Processing Strategy

**Decision**: Process both filenames and content

**Rationale**:

- Dynamic file structure based on parameters
- Flexible template system
- Support for conditional file inclusion
- Maintains template simplicity

### 5. Post-Scaffolding Actions

**Decision**: JSON-defined post-processing actions

**Rationale**:

- Declarative configuration
- Template-specific customization
- Conditional execution
- Extensible action types

## 🛡️ Error Handling

### Error Categories

1. **User Input Errors**

   - Invalid parameter values
   - Missing required parameters
   - Validation failures

2. **Template Errors**

   - Invalid template manifests
   - Missing template files
   - Template processing errors

3. **System Errors**
   - File system errors
   - Permission issues
   - Network errors (for future features)

### Error Handling Strategy

1. **Graceful Degradation**

   - Continue processing when possible
   - Provide clear error messages
   - Allow user recovery

2. **Validation First**

   - Validate templates before processing
   - Validate parameters before use
   - Early error detection

3. **User-Friendly Messages**
   - Clear, actionable error messages
   - Context-specific guidance
   - Debug information when appropriate

## 🔮 Future Architecture Considerations

### 1. Plugin System

**Planned**: Extensible plugin architecture

**Design**:

- Plugin interface for custom actions
- Plugin discovery and loading
- Plugin configuration management
- Plugin versioning and updates

### 2. Template Marketplace

**Planned**: External template distribution

**Design**:

- Template registry system
- Template versioning
- Template validation and testing
- Community template submission

### 3. Advanced Condition Engine

**Planned**: More sophisticated conditional logic

**Design**:

- Expression parser for complex conditions
- Support for mathematical operations
- String manipulation functions
- Boolean logic operators

### 4. Template Caching

**Planned**: Performance optimization

**Design**:

- Template manifest caching
- Compiled template caching
- Incremental template updates
- Cache invalidation strategies

## 📊 Performance Considerations

### Current Optimizations

1. **Embedded Filesystem**

   - No disk I/O for template loading
   - Fast template discovery
   - Memory-efficient storage

2. **Lazy Loading**

   - Load templates only when needed
   - Minimal memory footprint
   - Fast startup time

3. **Efficient Processing**
   - Stream-based file processing
   - Minimal memory allocations
   - Optimized template parsing

### Future Optimizations

1. **Parallel Processing**

   - Concurrent template processing
   - Parallel file operations
   - Background validation

2. **Caching Strategies**

   - Template manifest caching
   - Compiled template caching
   - Parameter validation caching

3. **Memory Management**
   - Streaming for large templates
   - Memory pool for allocations
   - Garbage collection optimization

## 🔍 Monitoring and Observability

### Current Observability

1. **Debug Logging**

   - Template discovery logs
   - Parameter processing logs
   - File operation logs

2. **Error Reporting**
   - Detailed error messages
   - Stack traces for debugging
   - Context information

### Planned Observability

1. **Metrics Collection**

   - Template usage statistics
   - Performance metrics
   - Error rate tracking

2. **Telemetry**

   - Usage analytics (opt-in)
   - Performance monitoring
   - Error reporting

3. **Health Checks**
   - Template validation status
   - System resource usage
   - Dependency health

---

**Last Updated**: 07/29/2025  
**Version**: v0.5.0  
**Maintainers**: [Project Maintainers]
