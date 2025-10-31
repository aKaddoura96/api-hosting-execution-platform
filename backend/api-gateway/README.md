# API Gateway

Main HTTP router for the platform. Handles authentication, rate limiting, and routing to execution containers.

## Features

- Request routing to API endpoints
- JWT-based authentication
- API key validation
- Rate limiting per user/endpoint
- Request logging and metrics

## Running locally

```bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Run the gateway
go run main.go
```

## Endpoints

- `GET /health` - Health check
- `GET /api/v1/endpoints` - List all available API endpoints (coming soon)
- `POST /api/v1/execute/{endpoint}` - Execute an API (coming soon)

## TODO

- [ ] Implement JWT authentication middleware
- [ ] Add API key validation
- [ ] Implement rate limiting with Redis
- [ ] Add request routing to execution containers
- [ ] Implement request logging
