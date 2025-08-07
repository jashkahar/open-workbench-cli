# Redis Cache Component

A production-ready Redis cache component for session storage and application caching.

## Features

- **Session Storage**: Store user sessions and authentication data
- **Application Caching**: Cache frequently accessed data
- **Memory Management**: Configurable memory limits and eviction policies
- **Persistence**: Automatic data persistence with configurable intervals
- **Health Checks**: Built-in health check via Redis PING
- **Docker Ready**: Containerized with Alpine Linux
- **Security**: Optional password protection

## Configuration

### Parameters

- `ProjectName`: Project name for the cache (default: "cache")
- `Owner`: Project owner (default: "Open Workbench")
- `Port`: Port for the Redis server (default: 6379)
- `Password`: Redis password (optional, default: "")

### Usage

This component is designed to be used as a shared infrastructure component in your Open Workbench project.

## Docker

Build and run:

```bash
docker build -t redis-cache .
docker run -p 6379:6379 redis-cache
```

## Health Check

The Redis server responds to PING commands for health checks.

## Connection

Connect to Redis using:

```bash
redis-cli -h localhost -p 6379
```

## Next Steps

1. Configure memory limits based on your needs
2. Set up password authentication if required
3. Configure persistence settings
4. Set up monitoring and alerting
