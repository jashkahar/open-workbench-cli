# Open Workbench Platform

A simple command-line tool that helps you create modern web applications quickly. It provides ready-to-use templates with best practices so you can focus on building your app instead of setting up the basics.

## What is it?

Open Workbench Platform takes the hassle out of starting new projects. Instead of spending hours setting up configurations, installing dependencies, and creating boilerplate code, you can get a fully working application in minutes.

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
# Start a new project
om init

# Add a service to your project
om add service

# See what templates are available
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

## What You Get

- **Ready-to-run applications** with all dependencies and configurations
- **Production-ready templates** with testing, Docker, and deployment setup
- **Smart command system** that adapts to how you want to work
- **Multiple frameworks** including React, Next.js, FastAPI, Express, and Vue
- **Docker support** for easy deployment and development

## Available Templates

- **nextjs-full-stack**: Complete Next.js app with TypeScript and testing
- **fastapi-basic**: FastAPI backend with automatic API documentation
- **react-typescript**: Modern React app with Vite and TypeScript
- **express-api**: Node.js Express API with TypeScript
- **vue-nuxt**: Vue.js Nuxt application
- **nginx-gateway**: Nginx reverse proxy for microservices
- **redis-cache**: Redis cache service

## Documentation

For detailed guides and technical information, see the [docs/](docs/) directory:

- [User Guide](docs/user-guide.md) - Complete usage guide
- [Architecture](docs/architecture.md) - How the system works
- [Development](docs/development.md) - Contributing and development setup
- [Template System](docs/template-system.md) - Creating custom templates

## License

MIT License - see [LICENSE](LICENSE) for details.
