# Open Workbench Platform

A command-line tool for quickly scaffolding modern web applications with pre-configured templates and best practices.

## What is it?

Open Workbench Platform helps developers bootstrap new projects using carefully crafted templates. It eliminates repetitive setup and ensures you start with a solid foundation following industry best practices.

## Quick Start

### Installation

```bash
# Homebrew (macOS)
brew install jashkahar/tap/om

# Scoop (Windows)
scoop bucket add jashkahar https://github.com/jashkahar/scoop-bucket
scoop install om

# Or download from GitHub Releases
# https://github.com/jashkahar/open-workbench-platform/releases
```

### Basic Usage

```bash
# Initialize a new project
om init

# Add a service to your project
om add service

# List available templates
om list-templates
```

### Examples

```bash
# Create a Next.js frontend
om add service --name frontend --template nextjs-full-stack

# Create a FastAPI backend
om add service --name backend --template fastapi-basic

# Create a React app
om add service --name app --template react-typescript
```

## Available Templates

- **nextjs-full-stack**: Production-ready Next.js with TypeScript, testing, and Docker
- **fastapi-basic**: FastAPI backend with automatic API documentation
- **react-typescript**: Modern React app with Vite and TypeScript
- **express-api**: Node.js Express API with TypeScript
- **vue-nuxt**: Vue.js Nuxt application with SSR

## Documentation

For detailed information, see the [docs/](docs/) directory:

- [User Guide](docs/user-guide.md) - Complete usage guide
- [Architecture](docs/architecture.md) - System design and components
- [Development](docs/development.md) - Contributing and development setup
- [Template System](docs/template-system.md) - Creating custom templates

## License

MIT License - see [LICENSE](LICENSE) for details.
