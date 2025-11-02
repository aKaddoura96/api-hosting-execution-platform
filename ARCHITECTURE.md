# Architecture Overview

## System Design

### High-Level Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                         Frontend Layer                        │
│                    Next.js 14 + React + Tailwind             │
│           Auth Context • API Client • UI Components          │
└─────────────────────────┬────────────────────────────────────┘
                          │ HTTP/REST
                          │
┌─────────────────────────▼────────────────────────────────────┐
│                      API Gateway (Go)                         │
│                        Port 8080                              │
├───────────────────────────────────────────────────────────────┤
│  • Authentication (JWT)                                       │
│  • API Management (CRUD)                                      │
│  • Code Upload                                                │
│  • Deployment Orchestration                                   │
│  • Request Routing (Future)                                   │
└──┬────────┬────────┬────────────────┬────────────────────────┘
   │        │        │                │
   │        │        │                │
┌──▼────┐ ┌▼─────┐ ┌▼──────────┐  ┌──▼──────────┐
│Executor│ │Analyt│ │PostgreSQL │  │   Redis     │
│Service │ │ics   │ │           │  │             │
│        │ │      │ │  Primary  │  │  Cache &    │
│Port    │ │Port  │ │  Database │  │  Sessions   │
│8081    │ │8082  │ │           │  │             │
└────────┘ └──────┘ └───────────┘  └─────────────┘
    │
    │ Docker API
    │
┌───▼─────────────────────────────────────────────────┐
│              Docker Engine                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │Python API│  │Node.js   │  │  Go API  │        │
│  │Container │  │Container │  │Container │        │
│  └──────────┘  └──────────┘  └──────────┘        │
└─────────────────────────────────────────────────────┘
```

## Components

### 1. Frontend (Next.js 14)

**Purpose**: User interface for developers and consumers

**Technology**:
- Next.js 14 (App Router)
- React 18
- Tailwind CSS
- TypeScript

**Key Features**:
- Authentication UI (login/signup)
- Developer dashboard
- API creation & management
- Code upload interface
- Marketplace browsing
- Analytics visualization (future)

**File Structure**:
```
frontend/
├── app/
│   ├── login/page.tsx
│   ├── signup/page.tsx
│   ├── dashboard/page.tsx
│   └── marketplace/page.tsx
└── lib/
    ├── api.ts              # API client
    └── auth-context.tsx    # Auth state management
```

### 2. API Gateway (Go - Port 8080)

**Purpose**: Main HTTP server, request routing, authentication

**Responsibilities**:
- User authentication & authorization
- API CRUD operations
- Code file handling
- Deployment coordination
- Request proxying (future)
- Rate limiting (future)

**Handlers**:
- `/api/v1/auth/*` - Authentication
- `/api/v1/apis/*` - API management
- `/api/v1/marketplace/*` - Public APIs
- `/api/v1/apis/:id/upload` - Code upload
- `/api/v1/apis/:id/deploy` - Trigger deployment
- `/api/v1/apis/:id/status` - Check status

**Key Files**:
```
backend/api-gateway/
├── main.go
├── handlers/
│   ├── auth.go       # Signup, login, JWT
│   ├── api.go        # CRUD, upload
│   └── deploy.go     # Deployment coordination
└── middleware/
    └── auth.go       # JWT validation
```

### 3. Executor Service (Go - Port 8081)

**Purpose**: Container orchestration and lifecycle management

**Responsibilities**:
- Docker container deployment
- Container lifecycle (start/stop/restart)
- Resource allocation
- Runtime selection (Python/Node/Go)
- Health monitoring

**Container Configuration**:
- **Python**: `python:3.11-slim`
- **Node.js**: `node:18-alpine`
- **Go**: `golang:1.21-alpine`
- **Resources**: 512MB RAM, 1 CPU per container
- **Restart Policy**: unless-stopped

**Endpoints**:
- `POST /deploy` - Deploy API to container
- `POST /stop/:id` - Stop container
- `GET /status/:id` - Get container status

**Key Files**:
```
backend/executor/
├── main.go
└── container/
    └── manager.go    # Docker SDK integration
```

### 4. Analytics Service (Go - Port 8082)

**Purpose**: Usage tracking and metrics collection

**Responsibilities**:
- Execution logging
- Performance metrics
- Usage statistics
- Billing data aggregation
- Error tracking

**Metrics Tracked**:
- Request count
- Response time
- Success/error rate
- Resource usage
- API popularity

**Endpoints**:
- `POST /log` - Log execution
- `GET /stats/:id` - Get statistics
- `GET /history/:id` - Execution history

### 5. Database Layer

#### PostgreSQL (Port 5432)

**Purpose**: Primary data store

**Tables**:
```sql
users           - User accounts
apis            - API definitions
api_keys        - Authentication keys
executions      - Execution logs
usage           - Daily usage aggregation
transactions    - Billing records
subscriptions   - User subscriptions
```

**Indexes**:
- `idx_users_email` - Fast user lookup
- `idx_apis_user_id` - User's APIs
- `idx_apis_visibility` - Public APIs
- `idx_executions_api_id` - API metrics

#### Redis (Port 6379)

**Purpose**: Caching and sessions (future)

**Use Cases**:
- Session storage
- API response caching
- Rate limiting counters
- Real-time metrics

## Shared Libraries

### backend/shared/

**Purpose**: Common code across services

**Structure**:
```
shared/
├── auth/
│   └── jwt.go          # Token generation/validation
├── database/
│   └── db.go           # Database connection
├── models/
│   ├── user.go         # Data models
│   ├── api.go
│   ├── execution.go
│   └── billing.go
└── repository/
    ├── user_repo.go    # Database queries
    ├── api_repo.go
    └── execution_repo.go
```

## Data Flow

### 1. User Signup/Login
```
Browser → API Gateway → JWT Generation → PostgreSQL
        ← JWT Token ←
```

### 2. API Creation
```
Browser → API Gateway → Validate → PostgreSQL
        ← API Object ←
```

### 3. Code Upload
```
Browser → API Gateway → Save File → Update DB
        ← Success ←
```

### 4. Deployment
```
Browser → API Gateway → Executor → Docker Engine
        ← Status ←      ← Create Container ←
                        → Update DB
```

### 5. API Execution (Future)
```
Consumer → API Gateway → Proxy → Container
         ← Response ←    ← Execute ←
         → Analytics → Log Execution
```

## Security

### Authentication
- JWT-based authentication
- Password hashing with bcrypt
- Token expiration (24 hours)
- Role-based access (developer/consumer)

### Authorization
- User-specific API access
- Ownership verification
- API key validation (future)

### Container Security
- Read-only code mounts
- Resource limits
- Network isolation (future)
- Secret management (future)

## Scalability

### Current
- Microservices architecture
- Container-based deployment
- Horizontal scaling ready

### Future
- Load balancing
- Multiple executor nodes
- Database read replicas
- CDN for static assets
- Message queue for async tasks

## Monitoring & Observability

### Implemented
- Health check endpoints
- Basic logging
- Execution tracking

### Planned
- Structured logging
- Metrics collection (Prometheus)
- Distributed tracing
- Error tracking (Sentry)
- Performance monitoring

## Deployment

### Development
- Docker Compose for local services
- Hot reload for frontend
- Direct Go binary execution

### Production (Planned)
- Kubernetes orchestration
- Managed PostgreSQL (RDS/Cloud SQL)
- Managed Redis (ElastiCache/MemoryStore)
- Container registry (ECR/GCR)
- Load balancer
- Auto-scaling groups

## Technology Decisions

### Why Go?
- High performance
- Excellent concurrency
- Small binary size
- Great standard library
- Docker SDK support

### Why Next.js?
- Server-side rendering
- File-based routing
- Built-in API routes
- Great developer experience
- Production optimizations

### Why PostgreSQL?
- Robust and reliable
- JSONB support
- Full-text search
- Strong consistency
- Rich ecosystem

### Why Docker?
- Language-agnostic
- Isolation
- Resource control
- Easy deployment
- Industry standard

## Future Enhancements

1. **Request Proxying**: Route consumer requests to deployed containers
2. **Auto-scaling**: Scale containers based on load
3. **Multi-region**: Deploy closer to users
4. **CI/CD Integration**: GitHub Actions for API deployment
5. **API Versioning**: Support multiple versions
6. **WebSocket Support**: Real-time APIs
7. **GraphQL Support**: GraphQL gateway
8. **Serverless Functions**: FaaS-style execution
