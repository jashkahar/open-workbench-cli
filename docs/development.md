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
git clone https://github.com/jashkahar/open-workbench-cli.git
cd open-workbench-cli
```

#### 2. Install Dependencies

```bash
go mod tidy
```

#### 3. Verify Setup

```bash
# Build the project
go build -o open-workbench-cli main.go

# Run tests
go test ./...

# Run the CLI
./open-workbench-cli --help
```

## 🏗️ Project Structure

### Core Files

```
open-workbench-cli/
├── main.go                   # Application entry point
├── tui.go                    # Terminal User Interface
├── types.go                  # Shared type definitions
├── go.mod                    # Go module dependencies
├── go.sum                    # Dependency checksums
├── internal/
│   ├── templating/           # Dynamic templating system
│   │   ├── discovery.go      # Template discovery
│   │   ├── parameters.go     # Parameter processing
│   │   └── processor.go      # Template processing
│   └── template/             # Legacy system (deprecated)
├── templates/                # Template directory
├── docs/                     # Documentation
├── .github/                  # GitHub workflows
├── .goreleaser.yml          # Release configuration
└── CONTRIBUTING.md          # Contribution guidelines
```

### Key Components

#### Main Application (`main.go`)

- Entry point and command-line argument parsing
- Orchestrates different execution modes
- Coordinates the scaffolding process

#### Terminal UI (`tui.go`)

- Interactive template selection interface
- Uses Bubble Tea framework
- Handles user navigation and selection

#### Templating System (`internal/templating/`)

- **Discovery**: Template discovery and validation
- **Parameters**: Parameter collection and validation
- **Processor**: Template processing and file operations

## 🔧 Development Workflow

### 1. Making Changes

#### Code Style Guidelines

- Follow Go conventions and best practices
- Use `gofmt` for code formatting
- Write clear, descriptive comments
- Use meaningful variable and function names

#### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestFunctionName
```

#### Building

```bash
# Build for current platform
go build -o open-workbench-cli main.go

# Build for multiple platforms
GOOS=windows GOARCH=amd64 go build -o open-workbench-cli.exe main.go
GOOS=darwin GOARCH=amd64 go build -o open-workbench-cli-darwin main.go
GOOS=linux GOARCH=amd64 go build -o open-workbench-cli-linux main.go
```

### 2. Testing Your Changes

#### Manual Testing

```bash
# Test TUI mode
go run main.go ui

# Test interactive mode
go run main.go

# Test with specific template
go run main.go ui  # Then select a template
```

#### Template Testing

```bash
# Create a test project
go run main.go ui
# Select a template and complete the process

# Verify the generated project
cd generated-project
# Check that files are created correctly
# Verify template variables are substituted
# Test that the project builds/runs
```

### 3. Debugging

#### Debug Logging

The CLI includes debug logging for development:

```go
fmt.Printf("DEBUG: Selected template: '%s'\n", selectedTemplate)
fmt.Printf("DEBUG: Template manifest loaded: %s - %s\n", templateInfo.Name, templateInfo.Description)
fmt.Printf("DEBUG: Collected parameters: %+v\n", parameterValues)
```

#### Common Debug Scenarios

1. **Template Not Found**

   - Check template directory structure
   - Verify `template.json` exists and is valid JSON
   - Check template name in discovery logic

2. **Parameter Issues**

   - Verify parameter definitions in `template.json`
   - Check parameter validation logic
   - Test conditional parameter visibility

3. **File Processing Errors**
   - Check file paths and permissions
   - Verify template syntax in files
   - Test template variable substitution

### 4. Adding New Features

#### Feature Development Process

1. **Create a Feature Branch**

   ```bash
   git checkout -b feature/new-feature-name
   ```

2. **Implement the Feature**

   - Write code following project conventions
   - Add tests for new functionality
   - Update documentation

3. **Test Thoroughly**

   - Run existing tests
   - Add new tests for your feature
   - Test manually with different scenarios

4. **Update Documentation**

   - Update relevant documentation files
   - Add examples if applicable
   - Update README if needed

5. **Submit a Pull Request**
   - Create a detailed PR description
   - Include testing instructions
   - Reference any related issues

#### Adding New Templates

1. **Create Template Directory**

   ```bash
   mkdir templates/new-template-name
   ```

2. **Create Template Manifest**

   ```json
   {
     "name": "New Template",
     "description": "Description of the template",
     "parameters": [
       {
         "name": "ProjectName",
         "prompt": "Project Name:",
         "type": "string",
         "required": true
       }
     ]
   }
   ```

3. **Add Template Files**

   - Create the actual template files
   - Use Go template syntax for variables
   - Test template variable substitution

4. **Test the Template**
   ```bash
   go run main.go ui
   # Select your new template and test the process
   ```

#### Adding New Parameter Types

1. **Update Parameter Types**

   - Add new type to `Parameter` struct
   - Update validation logic in `parameters.go`
   - Add UI handling in `main.go`

2. **Update Template Processing**

   - Add processing logic in `processor.go`
   - Update template functions if needed

3. **Add Tests**
   - Test parameter validation
   - Test template processing
   - Test UI interaction

## 🧪 Testing Strategy

### Unit Tests

#### Testing Template Discovery

```go
func TestDiscoverTemplates(t *testing.T) {
    templates, err := templating.DiscoverTemplates(templatesFS)
    if err != nil {
        t.Fatalf("Failed to discover templates: %v", err)
    }

    if len(templates) == 0 {
        t.Error("No templates found")
    }
}
```

#### Testing Parameter Processing

```go
func TestParameterValidation(t *testing.T) {
    param := templating.Parameter{
        Name: "ProjectName",
        Type: "string",
        Validation: &templating.Validation{
            Regex: "^[a-z0-9-]+$",
        },
    }

    processor := templating.NewParameterProcessor(&templating.TemplateManifest{})

    // Test valid input
    err := processor.ValidateParameter(param, "valid-project")
    if err != nil {
        t.Errorf("Valid input should not error: %v", err)
    }

    // Test invalid input
    err = processor.ValidateParameter(param, "Invalid Project")
    if err == nil {
        t.Error("Invalid input should error")
    }
}
```

#### Testing Template Processing

```go
func TestTemplateProcessing(t *testing.T) {
    manifest := &templating.TemplateManifest{
        Name: "Test Template",
        Description: "Test description",
    }

    values := map[string]interface{}{
        "ProjectName": "test-project",
        "Owner": "test-owner",
    }

    processor := templating.NewTemplateProcessor(manifest, values)

    content := "Project: {{.ProjectName}}, Owner: {{.Owner}}"
    expected := "Project: test-project, Owner: test-owner"

    result, err := processor.ProcessTemplate(content)
    if err != nil {
        t.Fatalf("Failed to process template: %v", err)
    }

    if result != expected {
        t.Errorf("Expected %s, got %s", expected, result)
    }
}
```

### Integration Tests

#### Testing Full Workflow

```go
func TestFullScaffoldingWorkflow(t *testing.T) {
    // Test the complete workflow from template selection to project creation
    // This would test the integration between all components
}
```

### Manual Testing Checklist

- [ ] TUI mode works correctly
- [ ] Interactive mode works correctly
- [ ] Template discovery finds all templates
- [ ] Parameter collection works for all types
- [ ] Template processing substitutes variables correctly
- [ ] Post-scaffolding actions execute properly
- [ ] Error handling works for invalid inputs
- [ ] Help text and validation messages are clear

## 🔍 Code Quality

### Linting and Formatting

```bash
# Format code
go fmt ./...

# Run linter (if you have golangci-lint installed)
golangci-lint run

# Check for common issues
go vet ./...
```

### Code Review Checklist

- [ ] Code follows Go conventions
- [ ] Functions are well-documented
- [ ] Error handling is appropriate
- [ ] Tests are included for new functionality
- [ ] No obvious performance issues
- [ ] Security considerations addressed
- [ ] Documentation is updated

### Performance Considerations

- **Template Processing**: Use efficient string operations
- **File Operations**: Minimize disk I/O
- **Memory Usage**: Avoid unnecessary allocations
- **Concurrency**: Consider parallel processing for large templates

## 🚀 Release Process

### Pre-Release Checklist

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Version is updated in relevant files
- [ ] Changelog is updated
- [ ] Release notes are prepared

### Creating a Release

1. **Update Version**

   ```bash
   # Update version in main.go or version.go
   ```

2. **Create Release Tag**

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **Monitor Release**
   - Check GitHub Actions for build status
   - Verify binaries are uploaded
   - Test installation from package managers

### Release Automation

The project uses GoReleaser for automated releases:

- **Multi-platform builds**: Windows, macOS, Linux
- **Package managers**: Homebrew, Scoop
- **GitHub Releases**: Automatic release creation
- **Checksums**: SHA256 checksums for all binaries

## 🤝 Contributing Guidelines

### Before Contributing

1. **Check Existing Issues**: Look for existing issues or discussions
2. **Discuss Changes**: Open an issue for significant changes
3. **Follow the Style Guide**: Use consistent code style
4. **Write Tests**: Include tests for new functionality

### Pull Request Process

1. **Create Feature Branch**: Use descriptive branch names
2. **Make Changes**: Follow the development workflow
3. **Test Thoroughly**: Ensure all tests pass
4. **Update Documentation**: Update relevant docs
5. **Submit PR**: Include clear description and testing instructions

### Code Review Process

1. **Automated Checks**: CI/CD runs tests and linting
2. **Manual Review**: Maintainers review the code
3. **Address Feedback**: Make requested changes
4. **Merge**: Once approved, changes are merged

## 🔧 Development Tools

### Recommended Tools

- **VS Code**: With Go extension
- **GoLand**: Full-featured Go IDE
- **Delve**: Go debugger
- **golangci-lint**: Linting tool
- **gofmt**: Code formatter

### Useful Commands

```bash
# Run with debug output
go run main.go ui

# Build with debug symbols
go build -gcflags=all="-N -l" -o open-workbench-cli main.go

# Profile the application
go test -cpuprofile cpu.prof -bench .
go tool pprof cpu.prof

# Memory profiling
go test -memprofile mem.prof -bench .
go tool pprof mem.prof
```

### Debugging Tips

1. **Use Debug Logging**: Add debug prints for development
2. **Test Incrementally**: Test small changes frequently
3. **Use Delve**: For complex debugging scenarios
4. **Check Templates**: Verify template syntax and structure
5. **Monitor File Operations**: Check file creation and permissions

## 📚 Learning Resources

### Go Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go by Example](https://gobyexample.com/)

### Project-Specific Resources

- [Template System Documentation](./template-system.md)
- [Architecture Overview](./architecture.md)
- [API Reference](./api-reference.md)

### Related Technologies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): TUI framework
- [Survey](https://github.com/AlecAivazis/survey): Interactive prompts
- [GoReleaser](https://goreleaser.com/): Release automation

## 🆘 Getting Help

### Internal Resources

- **Documentation**: Check the docs directory
- **Issues**: Search existing GitHub issues
- **Discussions**: Use GitHub Discussions for questions

### External Resources

- **Go Community**: [Go Forum](https://forum.golangbridge.org/)
- **Stack Overflow**: Tag questions with `go`
- **Reddit**: r/golang community

### Contact Maintainers

- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and discussions
- **Email**: For private or sensitive matters

---

**Last Updated**: [Current Date]  
**Version**: [Current Version]  
**Maintainers**: [Project Maintainers]
