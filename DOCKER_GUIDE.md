# Docker Deployment Guide

## Overview

The entire API Platform runs in Docker containers for maximum simplicity and portability.

## Architecture

```
┌─────────────────────────────────────────────────┐
│              Docker Network                      │
│           (api-platform)                        │
│                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│  │PostgreSQL│  │  Redis   │  │ Frontend │     │
│  │  :5432   │  │  :6379   │  │  :3000   │     │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘     │
│       │             │               │           │
│  ┌────┴─────────────┴───────────────┘           │
│  │                                               │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │  │   API    │  │ Executor │  │Analytics │  │
│  │  │ Gateway  │  │  :8081   │  │  :8082   │  │
│  │  │  :8080   │  │          │  │          │  │
│  │  └──────────┘  └──────────┘  └──────────┘  │
│  │                                               │
└──┴───────────────────────────────────────────────┘
```

## Services

| Service | Container Name | Port | Description |
|---------|---------------|------|-------------|
| **PostgreSQL** | api-platform-postgres | 5432 | Primary database |
| **Redis** | api-platform-redis | 6379 | Cache & sessions |
| **API Gateway** | api-platform-gateway | 8080 | Main HTTP server |
| **Executor** | api-platform-executor | 8081 | Container orchestration |
| **Analytics** | api-platform-analytics | 8082 | Metrics & analytics |
| **Frontend** | api-platform-frontend | 3000 | Next.js React app |

## Quick Start

### 1. Start Everything

```powershell
# Using the helper script (recommended)
.\docker-start.ps1

# Or directly with docker-compose
docker-compose up --build -d
```

### 2. Apply Migrations

Migrations are automatically applied by the `docker-start.ps1` script.

Manual migration:
```powershell
Get-Content scripts/migrations/001_initial_schema.sql | docker exec -i api-platform-postgres psql -U apiplatform -d apiplatform
```

### 3. Access Services

- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Executor**: http://localhost:8081
- **Analytics**: http://localhost:8082

## Common Commands

### Start Services
```bash
# Start all services (build if needed)
docker-compose up --build

# Start in detached mode
docker-compose up -d

# Start specific service
docker-compose up api-gateway
```

### Stop Services
```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api-gateway

# Last 100 lines
docker-compose logs --tail=100 api-gateway
```

### Restart Services
```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart api-gateway
```

### Rebuild Services
```bash
# Rebuild all
docker-compose build

# Rebuild specific service
docker-compose build api-gateway

# Rebuild and restart
docker-compose up --build -d
```

## Environment Variables

Set in `docker-compose.yml`:

### API Gateway
- `PORT=8080`
- `DATABASE_URL` - PostgreSQL connection
- `REDIS_URL` - Redis connection
- `JWT_SECRET` - JWT signing key
- `EXECUTOR_URL` - Executor service URL
- `ANALYTICS_URL` - Analytics service URL

### Executor
- `PORT=8081`
- `DATABASE_URL` - PostgreSQL connection

### Analytics
- `PORT=8082`
- `DATABASE_URL` - PostgreSQL connection

### Frontend
- `NEXT_PUBLIC_API_URL` - API Gateway URL

## Development Workflow

### 1. Make Code Changes

Edit your code in the `backend/` or `frontend/` directories.

### 2. Rebuild Service

```bash
# Rebuild the changed service
docker-compose build api-gateway

# Restart with new image
docker-compose up -d api-gateway
```

### 3. View Logs

```bash
docker-compose logs -f api-gateway
```

## Troubleshooting

### Service won't start

```bash
# Check service status
docker-compose ps

# View logs
docker-compose logs api-gateway

# Restart service
docker-compose restart api-gateway
```

### Database connection issues

```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Test database connection
docker exec -it api-platform-postgres psql -U apiplatform -d apiplatform

# View PostgreSQL logs
docker-compose logs postgres
```

### Port already in use

```bash
# Find process using port
netstat -ano | findstr :8080

# Kill process (PowerShell)
taskkill /PID <PID> /F

# Or stop the service
docker-compose down
```

### Build failures

```bash
# Clean rebuild
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Reset everything

```bash
# WARNING: This deletes all data
docker-compose down -v
docker-compose up --build -d
```

## Production Deployment

### 1. Update Environment Variables

Edit `docker-compose.yml` or create `docker-compose.prod.yml`:

```yaml
environment:
  JWT_SECRET: ${JWT_SECRET}  # Use secrets management
  DATABASE_URL: ${DATABASE_URL}  # Use managed database
```

### 2. Use Production Images

```bash
# Build production images
docker-compose -f docker-compose.yml build

# Tag for registry
docker tag api-platform-gateway:latest your-registry/api-gateway:v1.0.0

# Push to registry
docker push your-registry/api-gateway:v1.0.0
```

### 3. Deploy

```bash
# Deploy to server
docker-compose -f docker-compose.prod.yml up -d
```

## Health Checks

All services include health checks:

```bash
# Check service health
docker ps

# API Gateway
curl http://localhost:8080/health

# Executor
curl http://localhost:8081/health

# Analytics
curl http://localhost:8082/health
```

## Volumes

### PostgreSQL Data
- Volume: `postgres_data`
- Location: Managed by Docker
- Backup: `docker exec api-platform-postgres pg_dump -U apiplatform apiplatform > backup.sql`

### Redis Data
- Volume: `redis_data`
- Location: Managed by Docker

## Networking

All services communicate via the `api-platform` Docker network:

- Services use service names as hostnames (e.g., `postgres`, `redis`, `api-gateway`)
- Internal communication uses Docker DNS
- External access through exposed ports

## Performance Tips

1. **Use BuildKit**: Set `DOCKER_BUILDKIT=1` for faster builds
2. **Layer Caching**: Order Dockerfile commands from least to most frequently changing
3. **Multi-stage Builds**: Already implemented for smaller images
4. **Resource Limits**: Add to docker-compose.yml if needed:

```yaml
services:
  api-gateway:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
```

## Monitoring

### Container Stats
```bash
docker stats
```

### Disk Usage
```bash
docker system df
```

### Clean Up
```bash
# Remove unused images
docker image prune

# Remove all unused resources
docker system prune -a
```

## Security

1. **Change default passwords** in production
2. **Use secrets management** for sensitive data
3. **Run containers as non-root** (already implemented in Dockerfiles)
4. **Keep images updated**: Regularly rebuild with latest base images
5. **Use specific image tags**: Avoid `latest` in production

## Backup & Restore

### Backup Database
```bash
docker exec api-platform-postgres pg_dump -U apiplatform apiplatform > backup.sql
```

### Restore Database
```bash
cat backup.sql | docker exec -i api-platform-postgres psql -U apiplatform -d apiplatform
```

## Support

For issues, check:
1. Service logs: `docker-compose logs`
2. Container status: `docker-compose ps`
3. Network connectivity: `docker network inspect api-platform`
