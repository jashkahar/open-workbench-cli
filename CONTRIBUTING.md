# Contributing to Open Workbench Platform

Thank you for your interest in contributing to Open Workbench Platform! This guide will help you get started.

## Development Setup

### Prerequisites

- **Go 1.24 or later**
- **Git**
- **Docker** (for testing deployment generation)

### Local Development Environment

1. **Clone the repository:**

   ```bash
   git clone https://github.com/jashkahar/open-workbench-platform.git
   cd open-workbench-platform
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Build the project:**

   ```bash
   go build -o bin/om main.go
   ```

4. **Test your build:**
   ```bash
   ./bin/om --help
   ```

## Running Tests

### All Tests

```bash
go test ./...
```

### Tests with Coverage

```bash
go test -cover ./...
```

### Specific Package Tests

```bash
go test ./cmd/...
go test ./internal/templating/...
go test ./internal/manifest/...
```

### Integration Tests

```bash
go test -v ./cmd/init_test.go
go test -v ./cmd/security_test.go
```

## Code Style Guidelines

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for code formatting
- Keep functions focused and concise
- Add comments for exported functions and types

### File Organization

- **Commands**: Place in `cmd/` directory
- **Internal Logic**: Place in `internal/` directory
- **Templates**: Place in `templates/` directory
- **Documentation**: Place in `docs/` directory

### Naming Conventions

- **Functions**: Use camelCase for function names
- **Variables**: Use camelCase for variable names
- **Constants**: Use UPPER_SNAKE_CASE for constants
- **Types**: Use PascalCase for type names

## Pull Request Process

### Before Submitting

1. **Ensure tests pass:**

   ```bash
   go test ./...
   ```

2. **Format your code:**

   ```bash
   go fmt ./...
   ```

3. **Check for linting issues:**

   ```bash
   golangci-lint run
   ```

4. **Update documentation** if your changes affect user-facing functionality

### Creating a Pull Request

1. **Fork the repository** on GitHub
2. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes** following the code style guidelines
4. **Commit your changes:**
   ```bash
   git commit -m "Add feature: brief description"
   ```
5. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```
6. **Create a pull request** with a clear description

### Pull Request Guidelines

- **Title**: Clear, descriptive title
- **Description**: Explain what the PR does and why
- **Testing**: Describe how you tested the changes
- **Breaking Changes**: Note any breaking changes
- **Documentation**: Update relevant documentation

## Areas for Contribution

### High Priority

- **Bug fixes** in existing functionality
- **Documentation improvements** for unclear sections
- **Test coverage** for untested code
- **Performance improvements** for slow operations

### Medium Priority

- **New templates** for popular frameworks
- **Additional deployment targets** (Kubernetes, etc.)
- **Enhanced parameter types** (file upload, etc.)
- **Better error messages** and user feedback

### Low Priority

- **UI improvements** for interactive prompts
- **Additional validation rules** for parameters
- **More template examples** and documentation
- **CI/CD improvements** and automation

## Template Contributions

### Creating New Templates

1. **Study existing templates** in the `templates/` directory
2. **Follow the template structure** guidelines in [CREATING_A_TEMPLATE.md](docs/CREATING_A_TEMPLATE.md)
3. **Test your template** thoroughly before submitting
4. **Include comprehensive documentation** with your template

### Template Review Criteria

- **Functionality**: Does the template work correctly?
- **Documentation**: Is the template well-documented?
- **Best Practices**: Does it follow established patterns?
- **Security**: Are there any security concerns?

## Reporting Issues

### Bug Reports

When reporting bugs, please include:

- **Clear description** of the problem
- **Steps to reproduce** the issue
- **Expected behavior** vs actual behavior
- **Environment details** (OS, Go version, etc.)
- **Error messages** and stack traces

### Feature Requests

When requesting features, please include:

- **Clear description** of the feature
- **Use case** and motivation
- **Proposed implementation** if you have ideas
- **Impact** on existing functionality

## Getting Help

### Before Asking

1. **Check existing issues** for similar problems
2. **Read the documentation** in the `docs/` directory
3. **Try the latest version** from the main branch
4. **Search the codebase** for relevant code

### Asking Questions

- **Be specific** about your problem
- **Include relevant code** and error messages
- **Describe what you've tried** already
- **Use clear, concise language**

## Code of Conduct

### Our Standards

- **Be respectful** of others
- **Be collaborative** and helpful
- **Be constructive** in feedback
- **Be inclusive** of diverse perspectives

### Enforcement

- **Report violations** to project maintainers
- **Investigation** of reported violations
- **Appropriate action** for confirmed violations

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **Major**: Breaking changes
- **Minor**: New features (backward compatible)
- **Patch**: Bug fixes (backward compatible)

### Release Steps

1. **Update version** in relevant files
2. **Update changelog** with new features/fixes
3. **Run full test suite** to ensure stability
4. **Create release** on GitHub with release notes
5. **Update documentation** if needed

## Acknowledgments

Thank you to all contributors who have helped make Open Workbench Platform better!

---

**Questions?** Open an issue or reach out to the maintainers.
