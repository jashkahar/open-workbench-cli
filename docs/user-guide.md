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
# Start the interactive TUI
om ui

# Or use simple interactive mode
om
```

## 🎯 Using the CLI

### Available Commands

| Command     | Description                    |
| ----------- | ------------------------------ |
| `om`        | Interactive mode (recommended) |
| `om create` | CLI mode with flags            |

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

**Features:**

- Non-interactive project creation
- Command-line flags for all options
- Suitable for CI/CD and automation
- Comprehensive help system

**Available Flags:**

| Flag                  | Description                        |
| --------------------- | ---------------------------------- |
| `--owner`             | Project owner (required)           |
| `--no-testing`        | Disable testing framework          |
| `--no-tailwind`       | Disable Tailwind CSS               |
| `--no-docker`         | Disable Docker configuration       |
| `--no-install-deps`   | Skip dependency installation       |
| `--no-git`            | Skip Git repository initialization |
| `--testing-framework` | Testing framework (Jest/Vitest)    |
| `--help`              | Show help message                  |

**Examples:**

```bash
# Create a Next.js project with all features
om create nextjs-full-stack my-app --owner="John Doe"

# Create a React project without testing
om create react-typescript my-react-app --owner="Dev Team" --no-testing

# Create a FastAPI project without git initialization
om create fastapi-basic my-api --owner="Backend Team" --no-git

# Get help for CLI mode
om create --help
```

### Parameter Collection

After selecting a template, the CLI will collect parameters:

#### Parameter Types

1. **String Parameters**

   - Text input with optional validation
   - Example: Project name, owner name

2. **Boolean Parameters**

   - Yes/No questions with defaults
   - Example: "Include testing framework?"

3. **Select Parameters**

   - Single-choice dropdown
   - Example: "Which testing framework?" → [Jest, Vitest]

4. **Multiselect Parameters**
   - Multiple-choice selection
   - Example: "Which features?" → [Auth, Database, API]

#### Parameter Groups

Parameters are organized into logical groups:

- **Project Details**: Basic project information
- **Testing & Quality**: Testing framework and tools
- **Styling**: CSS frameworks and styling options
- **Deployment**: Docker, CI/CD configuration
- **Final Steps**: Post-setup actions

#### Validation

The CLI validates your input:

- **Required fields**: Must be provided
- **Format validation**: Regex patterns for specific formats
- **Custom error messages**: Clear guidance on what's wrong

### Example Workflow

```bash
$ om ui

🚀 Starting Open Workbench UI...

Please choose a project template:
┌─────────────────────────────────────────────────────────┐
│  🎨 Next.js Production-Grade                          │
│  A fully-featured Next.js application with testing,   │
│  linting, and optional CI/CD.                         │
│                                                       │
│  ⚡ FastAPI Basic                                     │
│  A FastAPI backend template with automatic API        │
│  documentation.                                        │
│                                                       │
│  🎯 React TypeScript                                  │
│  A modern React application with Vite and TypeScript. │
└─────────────────────────────────────────────────────────┘

[Select template with arrow keys, press Enter to confirm]

📋 Project Details
──────────────────
? Project Name: my-awesome-app
? Project Owner: jashkahar

📋 Testing & Quality
────────────────────
? Include a testing framework? Yes
? Which testing framework? Jest

📋 Deployment
─────────────
? Include Docker configuration? Yes

📋 Styling
──────────
? Include Tailwind CSS? Yes

📋 Final Steps
──────────────
? Install dependencies after setup? Yes

📂 Scaffolding project in './my-awesome-app'...
✏️  Applying templates...
🔧 Executing post-scaffolding actions...
------------------------------------
✅ Success! Your new project 'my-awesome-app' is ready.

Next steps:
1. cd my-awesome-app
2. npm install (if not already done)
3. npm run dev
```

## 📋 Available Templates

### 🎨 Next.js Production-Grade

A comprehensive Next.js template with modern tooling:

**Features:**

- TypeScript with strict configuration
- Testing framework (Jest or Vitest)
- ESLint and Prettier for code quality
- Tailwind CSS for styling
- Docker configuration
- GitHub Actions for CI/CD
- Husky for git hooks

**Parameters:**

- Project name and owner
- Testing framework selection
- Docker inclusion
- Tailwind CSS inclusion
- Automatic dependency installation

### ⚡ FastAPI Basic

A FastAPI backend template with automatic documentation:

**Features:**

- FastAPI with automatic OpenAPI docs
- Uvicorn ASGI server
- Python virtual environment setup
- Requirements.txt management
- Hot reload development server

**Parameters:**

- Project name
- Database support (optional)
- Database type selection

### 🎯 React TypeScript

A modern React application with Vite:

**Features:**

- React 18 with TypeScript
- Vite for fast development
- ESLint and Prettier
- Component library structure
- Modern tooling setup

**Parameters:**

- Project name and owner
- Testing framework selection
- Styling framework choice

### 🚀 Express API

A Node.js Express API template:

**Features:**

- Express.js with TypeScript
- Jest testing setup
- Swagger/OpenAPI documentation
- Middleware configuration
- Error handling

**Parameters:**

- Project name and owner
- Testing framework
- Documentation inclusion

### 🟢 Vue Nuxt

A Vue.js Nuxt application:

**Features:**

- Nuxt 3 with TypeScript
- Auto-imports for components
- Server-side rendering ready
- Modern Vue 3 composition API
- Built-in routing

**Parameters:**

- Project name and owner
- Styling framework
- Testing setup

## 🔧 Template Customization

### Understanding Template Variables

Templates use Go template syntax for variable substitution:

```go
// Basic variables
{{.ProjectName}}
{{.Owner}}

// Conditional content
{{if .IncludeTesting}}
import { render, screen } from '@testing-library/react';
{{end}}

// Template functions
{{lower .ProjectName}}
{{upper .ProjectName}}
{{title .ProjectName}}
```

### Template Functions

The CLI provides several template functions:

| Function   | Description           | Example                        |
| ---------- | --------------------- | ------------------------------ |
| `eq`       | Equality comparison   | `{{eq .Framework "React"}}`    |
| `ne`       | Inequality comparison | `{{ne .Framework "Vue"}}`      |
| `contains` | Array contains check  | `{{contains .Features "API"}}` |
| `lower`    | Convert to lowercase  | `{{lower .ProjectName}}`       |
| `upper`    | Convert to uppercase  | `{{upper .ProjectName}}`       |
| `title`    | Title case            | `{{title .ProjectName}}`       |
| `trim`     | Trim whitespace       | `{{trim .ProjectName}}`        |

### Conditional Logic

Templates support conditional logic:

```go
{{if .IncludeTesting}}
// Testing configuration
{{end}}

{{if .IncludeDocker}}
# Dockerfile content
FROM node:18-alpine
{{end}}
```

## 🛠️ Advanced Usage

### Non-Interactive Mode (Coming Soon)

For automation and scripting:

```bash
# Basic usage
om create --name my-project --template nextjs-full-stack

# With parameters
om create \
  --name my-project \
  --template nextjs-full-stack \
  --param "IncludeTesting=true" \
  --param "TestingFramework=Jest" \
  --param "IncludeDocker=true"
```

### Template Preview (Coming Soon)

Preview template output without creating files:

```bash
om preview nextjs-full-stack --params '{"ProjectName":"test"}'
```

### Custom Templates

Create your own templates:

1. **Create template directory**

   ```bash
   mkdir -p templates/my-custom-template
   ```

2. **Create template manifest**

   ```json
   {
     "name": "My Custom Template",
     "description": "Description of your template",
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

3. **Add template files**
   - Create your template files
   - Use Go template syntax for variables
   - Test the template

## 🔍 Troubleshooting

### Common Issues

#### 1. Template Not Found

**Problem**: CLI can't find available templates

**Solutions**:

- Check if templates are properly embedded
- Verify template directory structure
- Ensure `template.json` exists and is valid JSON

#### 2. Parameter Validation Errors

**Problem**: Input validation fails

**Solutions**:

- Check the validation rules in the template
- Follow the format requirements (e.g., lowercase for project names)
- Read the error message for specific guidance

#### 3. File Permission Errors

**Problem**: Can't create project directory or files

**Solutions**:

- Check directory permissions
- Ensure you have write access to the current directory
- Try running in a different directory

#### 4. Template Processing Errors

**Problem**: Template variables not substituted correctly

**Solutions**:

- Check template syntax in files
- Verify parameter names match template variables
- Look for syntax errors in Go template syntax

### Debug Mode

Enable debug output for troubleshooting:

```bash
# Run with debug logging
go run main.go ui

# Check debug output for:
# - Template discovery
# - Parameter collection
# - File processing
# - Post-scaffolding actions
```

### Getting Help

1. **Check Documentation**: Review this user guide and other docs
2. **Search Issues**: Look for similar issues on GitHub
3. **Create Issue**: Open a new issue with detailed information
4. **Community**: Ask questions in GitHub Discussions

## 📚 Best Practices

### Project Naming

- Use lowercase letters, numbers, and hyphens only
- Avoid spaces and special characters
- Keep names descriptive but concise
- Follow your organization's naming conventions

### Template Selection

- **Next.js**: For full-stack React applications
- **FastAPI**: For Python backend APIs
- **React TypeScript**: For frontend-only React apps
- **Express API**: For Node.js backend APIs
- **Vue Nuxt**: For Vue.js applications

### Parameter Configuration

- **Testing**: Include testing for production projects
- **Docker**: Include for containerized deployments
- **Styling**: Choose based on your team's preferences
- **Dependencies**: Let the CLI install dependencies automatically

### Post-Setup Steps

After project creation:

1. **Review the generated code**

   - Check file structure
   - Verify configuration files
   - Review package.json or requirements.txt

2. **Install dependencies**

   ```bash
   cd your-project
   npm install  # or pip install -r requirements.txt
   ```

3. **Start development**

   ```bash
   npm run dev  # or python main.py
   ```

4. **Run tests**

   ```bash
   npm test  # or pytest
   ```

5. **Initialize git repository**
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   ```

## 🔮 Future Features

### Planned Enhancements

- **Template Marketplace**: Browse and install community templates
- **Template Versioning**: Manage template updates
- **Plugin System**: Extend CLI functionality
- **IDE Integration**: VS Code and other IDE plugins
- **Cloud Deployment**: Direct deployment to cloud platforms

### Feature Requests

Have an idea for a new feature?

1. **Check existing issues** to avoid duplicates
2. **Create a feature request** with detailed description
3. **Provide use cases** and examples
4. **Consider contributing** the feature yourself

## 📞 Support

### Getting Help

- **Documentation**: Start with this user guide
- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and discussions
- **Community**: Join our community channels

### Contributing

Want to help improve the CLI?

- **Report bugs** with detailed information
- **Suggest features** with use cases
- **Contribute code** following our guidelines
- **Improve documentation** for better user experience

### Feedback

We value your feedback:

- **User experience**: What works well, what doesn't
- **Template suggestions**: New templates you'd like to see
- **Feature requests**: Functionality you need
- **Documentation**: What's unclear or missing

---

**Last Updated**: 07/29/2025  
**Version**: v0.5.0  
**Maintainers**: [Project Maintainers]
