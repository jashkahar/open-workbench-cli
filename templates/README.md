# Templates Directory

This directory contains all the available project templates for the Open Workbench CLI. Each template is a complete, production-ready project structure that can be customized through parameters.

## 📁 Directory Structure

```
templates/
├── README.md                    # This file
├── nextjs-golden-path/         # Next.js production template
│   ├── template.json           # Template manifest
│   ├── package.json            # Package configuration
│   ├── Dockerfile              # Docker configuration
│   ├── src/
│   │   └── app/
│   │       └── page.tsx        # Main page component
│   ├── tests/                  # Test files
│   └── README.md               # Template-specific README
├── fastapi-basic/              # FastAPI backend template
│   ├── template.json           # Template manifest
│   ├── main.py                 # FastAPI application
│   ├── requirements.txt        # Python dependencies
│   └── README.md               # Template-specific README
├── react-typescript/           # React TypeScript template
│   ├── template.json           # Template manifest
│   ├── package.json            # Package configuration
│   ├── src/
│   │   ├── App.tsx            # Main App component
│   │   └── main.tsx           # Entry point
│   └── README.md               # Template-specific README
├── express-api/                # Express.js API template
│   ├── template.json           # Template manifest
│   ├── package.json            # Package configuration
│   ├── src/
│   │   └── index.ts            # Express application
│   └── README.md               # Template-specific README
└── vue-nuxt/                   # Vue.js Nuxt template
    ├── template.json           # Template manifest
    ├── package.json            # Package configuration
    ├── pages/
    │   └── index.vue           # Main page
    └── README.md               # Template-specific README
```

## 🎯 Available Templates

### 🎨 nextjs-golden-path

A production-ready Next.js application with comprehensive tooling and best practices.

**Features:**

- Next.js 14 with App Router
- TypeScript with strict configuration
- Testing framework (Jest or Vitest)
- ESLint and Prettier for code quality
- Tailwind CSS for styling
- Docker configuration for containerization
- GitHub Actions for CI/CD
- Husky for git hooks
- Comprehensive error handling
- SEO optimization

**Parameters:**

- `ProjectName`: Project name (required)
- `Owner`: Project owner (required)
- `IncludeTesting`: Include testing framework (boolean, default: true)
- `TestingFramework`: Testing framework choice (select: Jest/Vitest)
- `IncludeDocker`: Include Docker configuration (boolean, default: true)
- `IncludeTailwind`: Include Tailwind CSS (boolean, default: true)
- `InstallDeps`: Install dependencies after setup (boolean, default: true)

**Post-Scaffolding Actions:**

- Initialize git repository
- Install dependencies (if enabled)
- Remove unused configuration files based on choices

**Use Case:** Full-stack React applications, production web apps, SEO-focused websites

### ⚡ fastapi-basic

A FastAPI backend template with automatic API documentation and modern Python practices.

**Features:**

- FastAPI with automatic OpenAPI/Swagger documentation
- Uvicorn ASGI server for production
- Python virtual environment setup
- Requirements.txt dependency management
- Hot reload development server
- CORS configuration
- Error handling middleware
- Health check endpoint
- Environment variable configuration

**Parameters:**

- `ProjectName`: Project name (required)
- `IncludeDatabase`: Include database support (boolean, default: false)
- `DatabaseType`: Database type (select: PostgreSQL/SQLite/MySQL)

**Post-Scaffolding Actions:**

- Create virtual environment
- Install Python dependencies

**Use Case:** Backend APIs, microservices, RESTful services

### 🎯 react-typescript

A modern React application with Vite and TypeScript for fast development.

**Features:**

- React 18 with TypeScript
- Vite for lightning-fast development
- ESLint and Prettier configuration
- Component library structure
- Modern tooling setup
- Hot module replacement
- Optimized build configuration
- TypeScript strict mode

**Parameters:**

- `ProjectName`: Project name (required)
- `Owner`: Project owner (required)
- `IncludeTesting`: Include testing framework (boolean, default: true)
- `TestingFramework`: Testing framework choice (select: Jest/Vitest)
- `IncludeStyling`: Include styling framework (boolean, default: true)
- `StylingFramework`: Styling framework choice (select: Tailwind/CSS Modules/Styled Components)

**Post-Scaffolding Actions:**

- Initialize git repository
- Install dependencies

**Use Case:** Frontend applications, component libraries, single-page applications

### 🚀 express-api

A Node.js Express API template with TypeScript and comprehensive tooling.

**Features:**

- Express.js with TypeScript
- Jest testing setup with API testing utilities
- Swagger/OpenAPI documentation
- Middleware configuration
- Error handling and logging
- CORS configuration
- Environment variable management
- Request validation
- Rate limiting

**Parameters:**

- `ProjectName`: Project name (required)
- `Owner`: Project owner (required)
- `IncludeTesting`: Include testing framework (boolean, default: true)
- `IncludeDocumentation`: Include API documentation (boolean, default: true)
- `IncludeValidation`: Include request validation (boolean, default: true)

**Post-Scaffolding Actions:**

- Initialize git repository
- Install dependencies

**Use Case:** Node.js APIs, backend services, RESTful APIs

### 🟢 vue-nuxt

A Vue.js Nuxt application with modern Vue 3 features and SSR capabilities.

**Features:**

- Nuxt 3 with Vue 3 Composition API
- TypeScript support
- Auto-imports for components and composables
- Server-side rendering (SSR) ready
- Built-in routing and state management
- Modern CSS with PostCSS
- ESLint and Prettier configuration
- SEO optimization

**Parameters:**

- `ProjectName`: Project name (required)
- `Owner`: Project owner (required)
- `IncludeStyling`: Include styling framework (boolean, default: true)
- `StylingFramework`: Styling framework choice (select: Tailwind/CSS Modules/Styled Components)
- `IncludeTesting`: Include testing framework (boolean, default: true)

**Post-Scaffolding Actions:**

- Initialize git repository
- Install dependencies

**Use Case:** Vue.js applications, SSR applications, content-heavy websites

## 📋 Template Manifest Structure

Each template includes a `template.json` file that defines:

### Basic Information

```json
{
  "name": "Template Display Name",
  "description": "Template description for users"
}
```

### Parameters

```json
{
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
  ]
}
```

### Post-Scaffolding Actions

```json
{
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

## 🔧 Template File Processing

### Template Variables

Template files use Go template syntax for variable substitution:

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

### File Name Templates

File names can also be templated:

```
{{if .IncludeTesting}}tests/{{end}}App.test.tsx
{{if .IncludeDocker}}Dockerfile{{end}}
{{lower .ProjectName}}-config.js
```

### Template Functions

Available template functions:

| Function   | Description           | Example                        |
| ---------- | --------------------- | ------------------------------ |
| `eq`       | Equality comparison   | `{{eq .Framework "React"}}`    |
| `ne`       | Inequality comparison | `{{ne .Framework "Vue"}}`      |
| `contains` | Array contains check  | `{{contains .Features "API"}}` |
| `lower`    | Convert to lowercase  | `{{lower .ProjectName}}`       |
| `upper`    | Convert to uppercase  | `{{upper .ProjectName}}`       |
| `title`    | Title case            | `{{title .ProjectName}}`       |
| `trim`     | Trim whitespace       | `{{trim .ProjectName}}`        |

## 🛠️ Creating Custom Templates

### 1. Create Template Directory

```bash
mkdir templates/my-custom-template
cd templates/my-custom-template
```

### 2. Create Template Manifest

Create `template.json`:

```json
{
  "name": "My Custom Template",
  "description": "Description of your template",
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
      "name": "IncludeFeature",
      "prompt": "Include advanced features?",
      "type": "boolean",
      "default": false
    }
  ],
  "postScaffold": {
    "commands": [
      {
        "command": "npm install",
        "description": "Installing dependencies..."
      }
    ]
  }
}
```

### 3. Add Template Files

Create your template files with Go template syntax:

```typescript
// src/App.tsx
import React from 'react';

function App() {
  return (
    <div className="App">
      <h1>Welcome to {{.ProjectName}}</h1>
      {{if .IncludeFeature}}
      <p>Advanced features are enabled!</p>
      {{end}}
    </div>
  );
}

export default App;
```

### 4. Test Your Template

```bash
# Build the CLI
go build -o open-workbench-cli main.go

# Test your template
./open-workbench-cli ui
# Select your template and test the process
```

## 📚 Template Best Practices

### Template Design

1. **Clear Structure**: Organize files logically
2. **Comprehensive Documentation**: Include detailed README files
3. **Modern Tooling**: Use current best practices and tools
4. **Flexible Configuration**: Allow customization through parameters
5. **Production Ready**: Include testing, linting, and deployment configs

### Parameter Design

1. **Descriptive Names**: Use clear, consistent parameter names
2. **Logical Grouping**: Group related parameters together
3. **Sensible Defaults**: Provide helpful default values
4. **Comprehensive Validation**: Validate user input appropriately
5. **Clear Help Text**: Provide helpful guidance for complex parameters

### File Organization

1. **Standard Structure**: Follow framework conventions
2. **Conditional Files**: Use conditional inclusion sparingly
3. **Clear Naming**: Use descriptive file and directory names
4. **Minimal Dependencies**: Minimize template dependencies

### Post-Scaffolding Actions

1. **Essential Setup**: Include basic setup commands
2. **Conditional Actions**: Use conditions for optional actions
3. **Clear Descriptions**: Provide clear action descriptions
4. **Error Handling**: Handle potential command failures

## 🔍 Template Testing

### Manual Testing

1. **Template Discovery**: Verify template appears in CLI
2. **Parameter Collection**: Test all parameter types and validation
3. **File Generation**: Check that files are created correctly
4. **Template Processing**: Verify variable substitution works
5. **Post-Scaffolding**: Test conditional actions and commands
6. **Project Functionality**: Ensure generated project works

### Automated Testing

```bash
# Test template discovery
go test -run TestDiscoverTemplates

# Test parameter validation
go test -run TestParameterValidation

# Test template processing
go test -run TestTemplateProcessing
```

## 🔮 Future Templates

### Planned Templates

1. **Angular**: Full-featured Angular application
2. **Svelte**: Modern Svelte application
3. **Go API**: Go backend API template
4. **Rust API**: Rust backend API template
5. **Python Web**: Django or Flask web application
6. **Mobile**: React Native or Flutter templates
7. **Desktop**: Electron or Tauri application
8. **Microservices**: Multi-service architecture
9. **Full-Stack**: Complete full-stack applications
10. **CMS**: Content management system templates

### Template Categories

- **Frontend**: React, Vue, Angular, Svelte
- **Backend**: Node.js, Python, Go, Rust, Java
- **Full-Stack**: Next.js, Nuxt, SvelteKit
- **Mobile**: React Native, Flutter
- **Desktop**: Electron, Tauri
- **Microservices**: Multi-service architectures
- **Specialized**: CMS, E-commerce, Admin panels

## 📞 Contributing Templates

### Template Guidelines

1. **Follow Standards**: Use current best practices
2. **Include Documentation**: Comprehensive README and comments
3. **Test Thoroughly**: Ensure template works correctly
4. **Keep Updated**: Maintain with latest framework versions
5. **Consider Users**: Design for different skill levels

### Submission Process

1. **Create Template**: Follow the template creation guide
2. **Test Locally**: Ensure template works correctly
3. **Document**: Include comprehensive documentation
4. **Submit PR**: Create pull request with template
5. **Review**: Address feedback and make improvements

### Template Review Criteria

- **Functionality**: Template works correctly
- **Documentation**: Clear and comprehensive docs
- **Best Practices**: Follows current standards
- **Flexibility**: Appropriate parameter options
- **Maintainability**: Clean, well-organized code

---

**Last Updated**: 07/29/2025  
**Version**: v0.5.0  
**Maintainers**: [Project Maintainers]
