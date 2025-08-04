# Nginx Gateway Component

A production-ready Nginx reverse proxy gateway for load balancing and routing.

## Features

- **Reverse Proxy**: Route requests to backend services
- **Load Balancing**: Distribute traffic across multiple instances
- **SSL/TLS Support**: Optional SSL termination
- **Health Checks**: Built-in health check endpoint
- **Security Headers**: Comprehensive security headers
- **Gzip Compression**: Optimized for performance
- **Docker Ready**: Containerized with Alpine Linux

## Configuration

### Parameters

- `ProjectName`: Project name for the gateway (default: "gateway")
- `Owner`: Project owner (default: "Open Workbench")
- `Port`: Port for the nginx gateway (default: 80)
- `IncludeSSL`: Include SSL configuration (default: false)

### Usage

This component is designed to be used as a shared infrastructure component in your Open Workbench project.

## Docker

Build and run:

```bash
docker build -t nginx-gateway .
docker run -p 80:80 nginx-gateway
```

## Health Check

The gateway exposes a health check endpoint at `/health` that returns a 200 status when healthy.

## Next Steps

1. Configure routing rules in `default.conf`
2. Add SSL certificates if needed
3. Set up load balancing to backend services
4. Configure environment-specific settings
