# Open Workbench Platform

![Go](https://img.shields.io/badge/Go-1.24%2B-blue)
![License: MIT](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)

**Open Workbench Platform** is your all-in-one CLI for rapidly bootstrapping, developing, and managing modern, production-ready applications—locally and in the cloud.

### Why Open Workbench?

- **Tired of spending days wiring up boilerplate, Docker, and cloud configs for every new project?**
- **Frustrated by the complexity of multi-service (monorepo, microservices, full-stack) setups?**
- **Want to go from idea to running code in minutes, not weeks?**

Open Workbench solves the "blank canvas" problem by automating the tedious setup for multi-service applications. It lets you focus on writing code, not wiring up infrastructure.

### Who is it for?

- **Developers & teams** who want to quickly scaffold, run, and iterate on real-world apps
- **Startups & hackathons** needing to move fast with best practices baked in
- **Anyone** who wants a smooth path from local dev to cloud deployment

### What makes it unique?

- **One command to bootstrap**: Instantly create a new project with batteries-included templates (Node, Python, React, Vue, etc.)
- **Multi-service made easy**: Add APIs, frontends, databases, gateways, and more—no manual Docker or YAML
- **Local & cloud ready**: Generate Docker Compose for local dev, or Terraform for cloud infra, from the same config
- **Consistent, repeatable environments**: Share, version, and reproduce your entire stack with a single YAML file

**In short:** Open Workbench is the fastest way to go from zero to a running, production-grade app—locally or in the cloud.

## Try it in 60 seconds
```
# macOS / Linux
brew install jashkahar/tap/open-workbench-platform

# Windows (PowerShell)
scoop bucket add jashkahar https://github.com/jashkahar/scoop-bucket.git
scoop install open-workbench-platform

# If that does not work
go install github.com/jashkahar/open-workbench-platform@latest

# Quickstart (interactive)
om init
om add service --template fastapi-basic
om compose --target docker

# requires docker desktop installed
docker compose up --build

# Then open your browser (example)
# http://localhost:8000
```

## 🚀 Quick Start

### Installation

**macOS / Linux (Homebrew):**

```bash
brew install jashkahar/tap/open-workbench-platform
```

**Windows (Scoop):**

```bash
scoop bucket add jashkahar https://github.com/jashkahar/scoop-bucket.git
scoop install open-workbench-platform
```

### Usage

1. **Initialize a new project:**

   ```bash
   om init
   ```

   This creates a `workbench.yaml` file to define your project structure.

2. **Add a backend service:**

   ```bash
   om add service
   ```

   This adds services to your project. Available templates include:

   - `express-api`: Node.js Express API
   - `fastapi-basic`: Python FastAPI
   - `nextjs-full-stack`: Next.js full-stack app
   - `react-typescript`: React with TypeScript
   - `vue-nuxt`: Vue.js with Nuxt
   - `nginx-gateway`: Nginx reverse proxy
   - `redis-cache`: Redis cache service

3. **Add an infrastructure resource to a service (e.g., Postgres):**

   ```bash
   # Interactive
   om add resource

   # Direct
   om add resource --service backend --type postgres-db --name database
   ```

4. **Generate your local environment:**

   ```bash
   om compose
   ```

   This generates the `docker-compose.yml` file needed to run your application.

   **Available flags:**

   - `--target`: Specify deployment target (`docker`)
   - `--env`: Environment name (`dev`, `staging`, `prod`)

   **Examples:**

   ```bash
   # Interactive (prompts for target)
   om compose

   # Direct: generate Docker Compose for local development
   om compose --target docker
   ```

   **Deployment targets:**

   - Docker: Generates `docker-compose.yml` for local development.

5. **List your services:**

   ```bash
   om ls
   ```

   Shows all services in your project and their current status.

6. **Run your application:**
   ```bash
   docker compose -f docker-compose.yml up --build
   ```

### Additional commands

- `om list-templates`: List available templates and their parameters.

## 📚 Learn More

For a full command reference and details on the architecture, please see our **[Full Documentation](docs/README.md)**.

## 🤝 Contributing

We welcome contributions! Please check out our **[Contributing Guide](CONTRIBUTING.md)**.
