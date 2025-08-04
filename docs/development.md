# Development Guide

This guide provides comprehensive information for developers who want to contribute to the Open Workbench CLI project, including setup, development workflow, and best practices.

## 🚀 Getting Started

### Prerequisites

- **Go 1.21 or later**: The project uses Go modules and modern Go features
- **Git**: For version control and collaboration
- **Docker**: For testing Docker Compose functionality
- **Your preferred IDE**: VS Code, GoLand, Vim, etc.
- **Terminal**: For running commands and testing

### Development Environment Setup

#### 1. Clone the Repository

```bash
git clone https://github.com/jashkahar/open-workbench-platform.git
cd open-workbench-platform
```

#### 2. Install Dependencies

```bash
go mod tidy
```

#### 3. Verify Setup

```bash
# Build the project
go build -o om main.go

# Run tests
go test ./...

# Run the CLI
./om --help

# Test the new init command
./om init --help

# Test smart command system
./om add service --help
```

## 🏗️ Project Structure

### Core Files

```
om/
├── main.go                   # Application entry point with embedded FS
├── cmd/                      # Command implementations
│   ├── root.go              # Root command setup with Cobra
│   ├── init.go              # om init command implementation
│   ├── add_service.go       # Smart service addition with mode detection
│   ├── compose.go           # Docker Compose generation
│   ├── types.go             # YAML manifest type definitions
│   ├── security.go          # Security utilities and validation
│   ├── security_test.go     # Security tests (100% coverage)
│   └── init_test.go         # Init command tests
├── internal/
│   ├── templating/           # Dynamic templating system
│   │   ├── discovery.go      # Template discovery
│   │   ├── parameters.go     # Parameter processing
│   │   ├── processor.go      # Template processing
│   │   └── progress.go       # Progress tracking
│   └── compose/              # Docker Compose system
│       ├── generator.go      # Compose generation
│       ├── prerequisites.go  # Docker validation
│       └── types.go          # Compose types
├── templates/                # Template directory (embedded)
├── docs/                     # Documentation
├── go.mod                    # Go module dependencies
├── go.sum                    # Dependency checksums
├── .github/                  # GitHub workflows
├── .goreleaser.yml          # Release configuration
└── CONTRIBUTING.md          # Contribution guidelines
```

### Key Components

#### Command System (`cmd/`)

- **Root Command** (`cmd/root.go`): Main CLI setup with Cobra framework
- **Init Command** (`cmd/init.go`): Project initialization with `workbench.yaml` manifests
- **Add Service** (`cmd/add_service.go`): Smart service addition with mode detection
- **Compose Command** (`cmd/compose.go`): Docker Compose generation
- **Security** (`cmd/security.go`): Comprehensive security utilities and validation
- **Types** (`cmd/types.go`): YAML manifest type definitions

#### Main Application (`main.go`)

- Entry point with embedded filesystem
- Routes to command system
- Handles application lifecycle

#### Templating System (`internal/templating/`)

- **Discovery**: Template discovery and validation
- **Parameters**: Parameter collection and validation
- **Processor**: Template processing and file operations
- **Progress**: Progress tracking and user feedback

#### Compose System (`internal/compose/`)

- **Generator**: Docker Compose configuration generation
- **Prerequisites**: Docker environment validation
- **Types**: Compose-specific type definitions

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific test files
go test ./cmd/security_test.go
go test ./cmd/init_test.go

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./cmd/
```

### Test Coverage

**Current Coverage**: 100% for security components, comprehensive for core functionality

**Test Categories**:

- **Unit Tests**: Individual function testing
- **Integration Tests**: Command system testing
- **Security Tests**: Security validation testing
- **Template Tests**: Template processing testing
- **Compose Tests**: Docker Compose generation testing

### Security Testing

**Comprehensive Security Test Suite**:

```bash
# Run security tests
go test ./cmd/security_test.go -v

# Run security benchmarks
go test -bench=Benchmark ./cmd/security_test.go
```

**Security Test Coverage**:

- Path traversal attack prevention
- Malicious pattern detection
- Cross-platform security validation
- Directory safety testing
- Template security validation

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

## 🔒 Security Development

### Security Guidelines

When contributing to the project, follow these security guidelines:

#### 1. Input Validation

Always validate user inputs:

```go
// ✅ Good: Validate all inputs
func ValidateProjectName(name string) error {
    if strings.Contains(name, "../") {
        return errors.New("path traversal not allowed")
    }
    return nil
}

// ❌ Bad: No validation
func CreateProject(name string) error {
    // Direct use without validation
    return os.Mkdir(name, 0755)
}
```

#### 2. Path Safety

Ensure safe path operations:

```go
// ✅ Good: Safe path handling
func SafePathJoin(base, name string) (string, error) {
    if strings.Contains(name, "..") {
        return "", errors.New("unsafe path")
    }
    return filepath.Join(base, name), nil
}
```

#### 3. Cross-Platform Security

Consider Windows and Unix security:

```go
// ✅ Good: Cross-platform security
var windowsReserved = []string{"con", "prn", "aux", "nul"}

func ValidateName(name string) error {
    for _, reserved := range windowsReserved {
        if strings.EqualFold(name, reserved) {
            return errors.New("reserved name")
        }
    }
    return nil
}
```

### Security Testing Requirements

All security-related code must have:

1. **Unit Tests**: Test all validation functions
2. **Edge Cases**: Test boundary conditions
3. **Malicious Inputs**: Test attack vectors
4. **Cross-Platform**: Test on Windows and Unix

## 🚀 Development Workflow

### 1. Feature Development

#### Adding New Commands

1. **Create Command File**:

```go
// cmd/new_command.go
package cmd

import (
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "new",
    Short: "New command description",
    RunE:  runNewCommand,
}

func runNewCommand(cmd *cobra.Command, args []string) error {
    // Implementation
    return nil
}
```

2. **Register Command**:

```go
// In cmd/root.go
func init() {
    rootCmd.AddCommand(newCmd)
}
```

3. **Add Tests**:

```go
// cmd/new_command_test.go
func TestRunNewCommand(t *testing.T) {
    // Test implementation
}
```

#### Adding New Templates

1. **Create Template Directory**:

```
templates/
└── new-template/
    ├── template.json
    ├── package.json
    └── src/
```

2. **Define Template Manifest**:

```json
{
  "name": "New Template",
  "description": "Template description",
  "parameters": [
    {
      "name": "ProjectName",
      "type": "string",
      "required": true
    }
  ]
}
```

3. **Test Template**:

```bash
# Test template processing
go test ./internal/templating/ -run TestProcessTemplate
```

### 2. Smart Command System

#### Mode Detection Logic

The smart command system automatically detects the mode based on provided parameters:

```go
func runAddService(cmd *cobra.Command, args []string) error {
    // Check if any parameters are provided
    name, _ := cmd.Flags().GetString("name")
    template, _ := cmd.Flags().GetString("template")
    params, _ := cmd.Flags().GetStringToString("params")

    if name != "" || template != "" || len(params) > 0 {
        // Direct mode - use provided parameters
        return runAddServiceDirect(cmd, args)
    } else {
        // Interactive mode - prompt for all parameters
        return runAddServiceInteractive(cmd, args)
    }
}
```

#### Adding New Modes

To add new modes to the smart command system:

1. **Implement Mode Detection**:

```go
func detectMode(cmd *cobra.Command) CommandMode {
    // Add your mode detection logic
    return InteractiveMode
}
```

2. **Add Mode Handler**:

```go
func runAddServiceWithMode(cmd *cobra.Command, args []string) error {
    mode := detectMode(cmd)

    switch mode {
    case InteractiveMode:
        return runAddServiceInteractive(cmd, args)
    case DirectMode:
        return runAddServiceDirect(cmd, args)
    case PartialMode:
        return runAddServicePartial(cmd, args)
    default:
        return runAddServiceInteractive(cmd, args)
    }
}
```

### 3. Docker Compose Integration

#### Adding Compose Features

1. **Extend Generator**:

```go
// internal/compose/generator.go
func (g *Generator) GenerateAdvanced() (*ComposeConfig, error) {
    // Add advanced compose features
    return &ComposeConfig{}, nil
}
```

2. **Add Prerequisites**:

```go
// internal/compose/prerequisites.go
func (c *Checker) CheckDockerVersion() error {
    // Check Docker version requirements
    return nil
}
```

3. **Test Compose Features**:

```go
// internal/compose/generator_test.go
func TestGenerateAdvanced(t *testing.T) {
    // Test advanced compose generation
}
```

## 📋 Code Quality

### Code Style

Follow Go conventions and project standards:

1. **Formatting**: Use `gofmt` or `goimports`
2. **Linting**: Use `golangci-lint`
3. **Documentation**: Document all exported functions
4. **Error Handling**: Use proper error wrapping

### Code Review Checklist

Before submitting a pull request, ensure:

- [ ] All tests pass
- [ ] Security tests pass
- [ ] Code is properly formatted
- [ ] Documentation is updated
- [ ] Error handling is comprehensive
- [ ] Cross-platform compatibility
- [ ] Performance considerations

### Performance Guidelines

1. **Efficient Validation**: Use regex for pattern matching
2. **Lazy Loading**: Load templates only when needed
3. **Caching**: Cache validation results when appropriate
4. **Memory Management**: Avoid unnecessary allocations

## 🔧 Debugging

### Debug Mode

Enable debug logging:

```bash
# Set debug environment variable
export OM_DEBUG=1

# Run with debug output
./om init
```

### Common Issues

#### Template Processing Issues

```bash
# Check template discovery
./om list-templates

# Test template processing
go test ./internal/templating/ -v
```

#### Security Validation Issues

```bash
# Test security functions
go test ./cmd/security_test.go -v

# Run security benchmarks
go test -bench=Benchmark ./cmd/security_test.go
```

#### Docker Compose Issues

```bash
# Test compose generation
go test ./internal/compose/ -v

# Check Docker prerequisites
docker --version
docker-compose --version
```

## 🚀 Release Process

### Pre-Release Checklist

1. **Test Coverage**: Ensure 100% coverage for security components
2. **Cross-Platform Testing**: Test on Windows, macOS, and Linux
3. **Security Audit**: Review all security-related changes
4. **Documentation**: Update all documentation
5. **Performance**: Run benchmarks to ensure no regressions

### Release Steps

1. **Update Version**: Update version in relevant files
2. **Run Tests**: Execute full test suite
3. **Build Binaries**: Build for all target platforms
4. **Create Release**: Create GitHub release
5. **Update Documentation**: Update version-specific docs

### Version Management

Follow semantic versioning:

- **Major**: Breaking changes
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes and security updates

## 🤝 Contributing

### Contribution Guidelines

1. **Fork the Repository**: Create your own fork
2. **Create Feature Branch**: Use descriptive branch names
3. **Write Tests**: Include tests for new features
4. **Update Documentation**: Update relevant documentation
5. **Submit Pull Request**: Include detailed description

### Development Setup

```bash
# Fork and clone
git clone https://github.com/your-username/open-workbench-platform.git
cd open-workbench-platform

# Add upstream remote
git remote add upstream https://github.com/jashkahar/open-workbench-platform.git

# Create feature branch
git checkout -b feature/new-command

# Make changes and test
go test ./...

# Commit and push
git commit -m "feat: add new command"
git push origin feature/new-command
```

### Code of Conduct

- Be respectful and inclusive
- Focus on technical discussions
- Provide constructive feedback
- Follow project conventions

## 📚 Additional Resources

### Documentation

- [User Guide](user-guide.md): Complete usage guide
- [Architecture](architecture.md): System design and components
- [Template System](template-system.md): Creating custom templates

### External Resources

- [Go Documentation](https://golang.org/doc/)
- [Cobra Framework](https://github.com/spf13/cobra)
- [Docker Documentation](https://docs.docker.com/)
- [Security Best Practices](https://owasp.org/)

### Community

- **Issues**: Report bugs and request features
- **Discussions**: Ask questions and share ideas
- **Pull Requests**: Contribute code and improvements

---

**Maintainer**: Jash Kahar  
**Last Updated**: August 3, 2025
