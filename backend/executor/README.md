# Executor Service

Container orchestration service for deploying and managing API containers.

## Features

- **Container Deployment** - Deploy user code in isolated Docker containers
- **Runtime Support** - Python, Node.js, and Go runtimes
- **Resource Limits** - 512MB memory, 1 CPU per container
- **Auto Restart** - Containers restart automatically unless stopped
- **Status Monitoring** - Check container health and logs
- **Lifecycle Management** - Start, stop, and remove containers

## How It Works

1. **Deploy Request** - API Gateway sends deploy request with API ID
2. **Code Loading** - Executor mounts uploaded code into container
3. **Image Pull** - Downloads appropriate runtime image if needed
4. **Container Start** - Starts container with resource limits
5. **Status Update** - Updates database with container ID and status

## API Endpoints

### Deploy API
```bash
POST /deploy
Body: {"api_id": "uuid"}
Response: {"status": "success", "container_id": "...", "message": "API deployed successfully"}
```

### Stop API
```bash
POST /stop/{api_id}
Response: {"status": "success", "message": "API stopped successfully"}
```

### Get Status
```bash
GET /status/{api_id}
Response: {"api_id": "...", "status": "running", "container_id": "..."}
```

### Health Check
```bash
GET /health
Response: {"status": "healthy", "service": "executor"}
```

## Runtime Images

- **Python**: `python:3.11-slim`
- **Node.js**: `node:18-alpine`
- **Go**: `golang:1.21-alpine`

## Resource Limits

- Memory: 512MB per container
- CPU: 1 core per container
- Restart Policy: unless-stopped

## Running Locally

```bash
# Ensure Docker is running
docker ps

# Set environment
cp .env.example .env

# Download dependencies
go mod download

# Run executor
go run main.go
```

## Environment Variables

```env
PORT=8081
DATABASE_URL=postgres://apiplatform:dev_password@localhost:5432/apiplatform?sslmode=disable
```

## Container Naming

Containers are named: `api-{api_id}`

## Volume Mounts

Code directory is mounted read-only at `/app/code` in container.

## Docker Requirements

- Docker must be installed and running
- Docker socket accessible at `/var/run/docker.sock`
- Sufficient resources for container limits

## Security

- Containers run in isolated environments
- Code mounted read-only
- Resource limits prevent abuse
- No network access between containers (by default)

## TODO

- [ ] Add Python Flask/FastAPI runtime wrapper
- [ ] Add Node.js Express runtime wrapper  
- [ ] Add Go HTTP server runtime wrapper
- [ ] Implement port mapping and routing
- [ ] Add container health checks
- [ ] Implement log streaming
- [ ] Add metrics collection
- [ ] Support custom Dockerfiles
- [ ] Implement auto-scaling
