# Development Guide

## Getting Started

### Prerequisites
- Go 1.21+ installed
- Node.js 18+ installed
- Docker Desktop running
- PostgreSQL client (psql) for migrations

### Initial Setup

```powershell
# 1. Clone repository
git clone https://github.com/aKaddoura96/api-hosting-execution-platform.git
cd api-hosting-execution-platform

# 2. Start Docker services
docker-compose up -d

# 3. Run migrations
.\scripts\run-migrations.ps1

# 4. Start backend services (3 separate terminals)
cd backend/api-gateway && go run main.go
cd backend/executor && go run main.go
cd backend/analytics && go run main.go

# 5. Start frontend
cd frontend && npm install && npm run dev
```

## Project Structure

```
├── backend/
│   ├── api-gateway/         # Main HTTP server (port 8080)
│   ├── executor/            # Container management (port 8081)
│   ├── analytics/           # Metrics (port 8082)
│   └── shared/              # Common code
├── frontend/                # Next.js app
├── scripts/                 # Utilities
│   ├── migrations/          # SQL migrations
│   ├── start.ps1            # Start all services
│   └── stop.ps1             # Stop all services
└── docker-compose.yml       # Local dev services
```

## Development Workflow

### 1. Backend Development

**Hot Reload**: Use `go run` for automatic recompilation

```powershell
cd backend/api-gateway
go run main.go
```

**Building**:
```powershell
go build -o api-gateway.exe main.go
```

**Testing**:
```powershell
go test ./...
```

### 2. Frontend Development

```powershell
cd frontend
npm run dev        # Development server
npm run build      # Production build
npm run lint       # Linting
```

### 3. Database Migrations

**Creating a new migration**:
```sql
-- scripts/migrations/00X_description.sql
ALTER TABLE apis ADD COLUMN new_field VARCHAR(255);
```

**Running migrations**:
```powershell
.\scripts\run-migrations.ps1
```

## Code Style

### Go
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Add comments for exported functions
- Use structured logging

### TypeScript/React
- Use TypeScript for type safety
- Follow React hooks best practices
- Use Tailwind for styling
- Componentize reusable UI

## Adding Features

### Adding a New API Endpoint

**1. Define handler** (`backend/api-gateway/handlers/`):
```go
func (h *Handler) NewEndpoint(w http.ResponseWriter, r *http.Request) {
    log.Printf("INFO: NewEndpoint called")
    // Implementation
}
```

**2. Add route** (`backend/api-gateway/main.go`):
```go
protected.HandleFunc("/new-endpoint", handler.NewEndpoint).Methods("GET")
```

**3. Update Postman collection** if needed

### Adding a New Database Table

**1. Create migration** (`scripts/migrations/00X_new_table.sql`):
```sql
CREATE TABLE new_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**2. Add model** (`backend/shared/models/`):
```go
type NewModel struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
}
```

**3. Add repository** (`backend/shared/repository/`):
```go
func (r *Repo) Create(model *NewModel) error {
    // Implementation
}
```

## Testing

### Manual Testing with Postman
1. Import `API-Platform.postman_collection.json`
2. Select "API Platform" environment
3. Run requests in order

### Automated Testing (Future)
```powershell
# Backend
go test ./...

# Frontend
npm test
```

## Debugging

### Backend Logs
All services log to stdout. Check logs for errors:
```powershell
# If running with start.ps1, logs are in terminal windows
# Or use Docker logs
docker-compose logs postgres
```

### Frontend Debugging
- Use React DevTools
- Check browser console
- Use Next.js debug mode: `npm run dev`

## Common Issues

### Port Already in Use
```powershell
# Find process using port
netstat -ano | findstr :8080
# Kill process
taskkill /PID <PID> /F
```

### Database Connection Failed
- Ensure Docker is running: `docker ps`
- Check connection string in `.env`
- Verify PostgreSQL is accessible: `docker-compose logs postgres`

### Frontend Won't Start
```powershell
# Clear node_modules and reinstall
rm -r node_modules
npm install
```

## Environment Variables

### Backend (.env files in each service)
```env
PORT=8080
DATABASE_URL=postgres://apiplatform:dev_password@localhost:5432/apiplatform?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key-change-in-production
EXECUTOR_URL=http://localhost:8081
ANALYTICS_URL=http://localhost:8082
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Git Workflow

1. Create feature branch: `git checkout -b feature/name`
2. Make changes
3. Commit: `git commit -m "feat: description"`
4. Push: `git push origin feature/name`
5. Create Pull Request

### Commit Message Format
```
feat: Add new feature
fix: Bug fix
docs: Documentation changes
style: Code style changes
refactor: Code refactoring
test: Adding tests
chore: Maintenance tasks
```

## Deployment

See [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md) for production deployment instructions.

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Next.js Documentation](https://nextjs.org/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Docker Documentation](https://docs.docker.com/)

## Getting Help

- Check existing issues on GitHub
- Review documentation files
- Ask in discussions
