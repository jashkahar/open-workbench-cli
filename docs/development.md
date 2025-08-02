# Development Guide

This guide provides comprehensive information for developers who want to contribute to the Open Workbench CLI project, including setup, development workflow, and best practices.

## 🚀 Getting Started

### Prerequisites

- **Go 1.21 or later**: The project uses Go modules and modern Go features
- **Git**: For version control and collaboration
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
```

## 🏗️ Project Structure

### Core Files

```
om/
├── main.go                   # Application entry point with embedded FS
├── tui.go                    # Terminal User Interface
├── cmd/                      # Command implementations
│   ├── root.go              # Root command setup with Cobra
│   ├── init.go              # om init command implementation
│   ├── types.go             # YAML manifest type definitions
│   ├── security.go          # Security utilities and validation
│   ├── security_test.go     # Security tests (100% coverage)
│   └── init_test.go         # Init command tests
├── internal/
│   ├── templating/           # Dynamic templating system
│   │   ├── discovery.go      # Template discovery
│   │   ├── parameters.go     # Parameter processing
│   │   └── processor.go      # Template processing
│   └── template/             # Legacy system (deprecated)
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
- **Security** (`cmd/security.go`): Comprehensive security utilities and validation
- **Types** (`cmd/types.go`): YAML manifest type definitions

#### Main Application (`main.go`)

- Entry point with embedded filesystem
- Routes to command system
- Handles application lifecycle

#### Terminal UI (`tui.go`)

- Interactive template selection interface
- Uses Bubble Tea framework
- Handles user navigation and selection

#### Templating System (`internal/templating/`)

- **Discovery**: Template discovery and validation
- **Parameters**: Parameter collection and validation
- **Processor**: Template processing and file operations

#### Security System (`cmd/security.go`)

- **Input Validation**: Comprehensive validation for all inputs
- **Path Security**: Path traversal protection and sanitization
- **Malicious Pattern Detection**: JavaScript injection, command injection prevention
- **Cross-Platform Security**: Windows reserved names, absolute path prevention

## 🔧 Development Workflow

### 1. Making Changes

#### Code Style Guidelines

- Follow Go conventions and best practices
- Use `gofmt` for code formatting
- Write clear, descriptive comments
- Use meaningful variable and function names
- **Security First**: Always validate user inputs
- **Test Coverage**: Aim for 100% test coverage

#### Security Guidelines

When making changes, always consider security:

```go
// ✅ Good: Validate all user inputs
func processUserInput(input string) error {
    sanitized, err := ValidateAndSanitizeName(input, nil)
    if err != nil {
        return fmt.Errorf("invalid input: %w", err)
    }
    // Process sanitized input...
}

// ❌ Bad: Direct use of user input
func processUserInput(input string) error {
    // Direct use without validation
    return os.WriteFile(input, data, 0644)
}
```

#### Testing Requirements

**100% Test Coverage Required** for:

- Security functions
- Command functions
- Core logic functions

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

### 2. Adding New Commands

#### Command Structure

```go
// cmd/newcommand.go
package cmd

import (
    "github.com/spf13/cobra"
)

var newCommand = &cobra.Command{
    Use:   "newcommand",
    Short: "Description of the new command",
    Long: `Detailed description of the new command.

    This command does something useful.

    Example:
      om newcommand --flag value`,
    RunE: runNewCommand,
}

func runNewCommand(cmd *cobra.Command, args []string) error {
    // Validate inputs first
    if err := validateInputs(args); err != nil {
        return err
    }

    // Implement command logic
    return nil
}

func validateInputs(args []string) error {
    // Always validate user inputs
    for _, arg := range args {
        if _, err := ValidateAndSanitizePath(arg, nil); err != nil {
            return err
        }
    }
    return nil
}
```

#### Registering Commands

```go
// In cmd/root.go
func Execute(fs embed.FS) {
    // ... existing setup ...

    // Add new command
    rootCmd.AddCommand(newCommand)

    // ... rest of setup ...
}
```

### 3. Adding Security Features

#### Security Validation

```go
// Example: Adding new security validation
func ValidateNewInput(input string) error {
    // Check for empty input
    if strings.TrimSpace(input) == "" {
        return fmt.Errorf("input cannot be empty")
    }

    // Check length
    if len(input) > 100 {
        return fmt.Errorf("input too long (max 100 characters)")
    }

    // Check for malicious patterns
    if err := CheckForSuspiciousPatterns(input); err != nil {
        return err
    }

    return nil
}
```

#### Security Tests

```go
// Example: Adding security tests
func TestValidateNewInput(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectError bool
    }{
        {"valid input", "normal-input", false},
        {"empty input", "", true},
        {"malicious input", "javascript:alert(1)", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateNewInput(tt.input)
            if tt.expectError && err == nil {
                t.Errorf("expected error but got none")
            }
            if !tt.expectError && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

### 4. Testing Guidelines

#### Test Categories

1. **Unit Tests**: Individual function testing
2. **Integration Tests**: End-to-end workflow testing
3. **Security Tests**: Security validation testing
4. **Performance Tests**: Benchmark testing

#### Test Structure

```go
// Example test structure
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectError bool
        expected    string
    }{
        {"valid case", "input", false, "expected"},
        {"error case", "bad-input", true, ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)

            if tt.expectError {
                if err == nil {
                    t.Errorf("expected error but got none")
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if result != tt.expected {
                    t.Errorf("expected %s, got %s", tt.expected, result)
                }
            }
        })
    }
}
```

#### Benchmark Tests

```go
// Example benchmark test
func BenchmarkFunctionName(b *testing.B) {
    input := "test-input"
    for i := 0; i < b.N; i++ {
        _, err := FunctionName(input)
        if err != nil {
            b.Fatalf("unexpected error: %v", err)
        }
    }
}
```

### 5. Adding New Templates

#### Template Structure

```
templates/
├── my-new-template/
│   ├── template.json          # Template manifest
│   ├── package.json           # Template file
│   ├── src/
│   │   └── main.ts           # Template file
│   └── README.md             # Template file
```

#### Template Manifest

```json
{
  "name": "My New Template",
  "description": "Description of the template",
  "parameters": [
    {
      "name": "ProjectName",
      "prompt": "What is your project name?",
      "type": "string",
      "required": true,
      "validation": {
        "regex": "^[a-z0-9-]+$",
        "errorMessage": "Project name can only contain lowercase letters, numbers, and hyphens."
      }
    },
    {
      "name": "IncludeTesting",
      "prompt": "Include testing setup?",
      "type": "boolean",
      "default": true
    }
  ],
  "postScaffold": {
    "filesToDelete": [
      {
        "path": "example.js",
        "condition": "IncludeTesting == false"
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

#### Template Testing

```bash
# Test the new template
go run main.go create my-new-template test-project --owner="Test User"

# Verify the output
ls test-project/
cat test-project/package.json
```

## 🔒 Security Development

### Security Checklist

When contributing code, ensure:

- [ ] All user inputs are validated
- [ ] Path traversal attacks are prevented
- [ ] Malicious patterns are detected
- [ ] Cross-platform security is considered
- [ ] Security tests are included
- [ ] Error messages don't leak sensitive information

### Security Testing

```bash
# Run security tests
go test ./cmd -run TestSecurity

# Run with security focus
go test ./cmd -v -run "TestValidate|TestSecurity"
```

### Common Security Patterns

#### Input Validation

```go
// Always validate user inputs
func processInput(input string) error {
    sanitized, err := ValidateAndSanitizeName(input, nil)
    if err != nil {
        return err
    }
    // Process sanitized input...
}
```

#### Path Security

```go
// Always validate paths
func processPath(path string) error {
    cleanPath, err := ValidateAndSanitizePath(path, nil)
    if err != nil {
        return err
    }
    // Process clean path...
}
```

#### Template Security

```go
// Always validate template names
func processTemplate(templateName string) error {
    if err := ValidateTemplateName(templateName); err != nil {
        return err
    }
    // Process template...
}
```

## 🧪 Testing Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific test
go test ./cmd -run TestValidateAndSanitizeName

# Run with coverage
go test ./cmd -cover

# Run benchmarks
go test ./cmd -bench=.

# Run tests with verbose output
go test ./cmd -v
```

### Test Coverage Requirements

- **Security Functions**: 100% coverage required
- **Command Functions**: 100% coverage required
- **Core Logic**: 100% coverage required
- **Error Paths**: 100% coverage required

### Adding Tests

#### Unit Tests

```go
func TestNewFunction(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectError bool
        expected    string
    }{
        {"valid input", "test", false, "test"},
        {"empty input", "", true, ""},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := NewFunction(tt.input)
            // Test logic...
        })
    }
}
```

#### Integration Tests

```go
func TestEndToEndWorkflow(t *testing.T) {
    // Setup test environment
    tempDir, err := os.MkdirTemp("", "test")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir)

    // Change to temp directory
    originalDir, err := os.Getwd()
    if err != nil {
        t.Fatalf("failed to get current directory: %v", err)
    }
    defer os.Chdir(originalDir)

    if err := os.Chdir(tempDir); err != nil {
        t.Fatalf("failed to change directory: %v", err)
    }

    // Test the workflow
    // ... test implementation
}
```

## 🚀 Performance Development

### Benchmarking

```bash
# Run benchmarks
go test ./cmd -bench=.

# Run specific benchmark
go test ./cmd -bench=BenchmarkValidateAndSanitizeName

# Run with memory profiling
go test ./cmd -bench=. -benchmem
```

### Performance Guidelines

- **Security Functions**: Should complete in < 100μs
- **Command Functions**: Should complete in < 1s
- **Template Processing**: Should complete in < 10s for typical templates

### Performance Testing

```go
func BenchmarkFunctionName(b *testing.B) {
    input := "test-input"
    for i := 0; i < b.N; i++ {
        _, err := FunctionName(input)
        if err != nil {
            b.Fatalf("unexpected error: %v", err)
        }
    }
}
```

## 📝 Documentation Development

### Documentation Requirements

When adding new features:

- [ ] Update README.md with new features
- [ ] Update user guide with usage examples
- [ ] Update architecture documentation
- [ ] Add inline code comments
- [ ] Update command help text

### Documentation Standards

- Use clear, concise language
- Include code examples
- Provide security considerations
- Include testing instructions
- Update maintainer information

## 🔄 Release Process

### Pre-Release Checklist

- [ ] All tests pass
- [ ] Security tests pass
- [ ] Performance benchmarks are acceptable
- [ ] Documentation is updated
- [ ] Version is updated
- [ ] Changelog is updated

### Release Steps

```bash
# 1. Update version
# Edit main.go or version file

# 2. Run all tests
go test ./...

# 3. Build for all platforms
go build -o om main.go

# 4. Create release tag
git tag v0.6.0
git push origin v0.6.0

# 5. Verify release
# Check GitHub releases
# Test installation methods
```

## 🤝 Contributing Guidelines

### Pull Request Process

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/new-feature`
3. **Make your changes** following the guidelines above
4. **Add tests** for new functionality
5. **Update documentation** as needed
6. **Run all tests**: `go test ./...`
7. **Commit your changes**: Use conventional commit messages
8. **Push to your branch**: `git push origin feature/new-feature`
9. **Create a pull request**

### Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

Examples:

- `feat(cmd): add new init command for project management`
- `fix(security): prevent path traversal in template names`
- `test(cmd): add comprehensive tests for security functions`
- `docs(readme): update with new security features`

### Code Review Process

- All changes require review
- Security changes require security review
- Tests must pass before merge
- Documentation must be updated
- Performance impact must be considered

---

**Maintainer**: Jash Kahar  
**Last Updated**: February 8, 2025
