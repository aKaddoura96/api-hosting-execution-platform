# Deployment Guide - MVP Phase 1 Complete! 🎉

## What's Working Now

✅ **Complete MVP with Container Execution**

### Backend Services
1. **API Gateway** (Port 8080) - Main HTTP server
   - Authentication (JWT)
   - API CRUD operations
   - Code uploads
   - Deploy/Stop/Status endpoints

2. **Executor Service** (Port 8081) - Container orchestration
   - Docker container management
   - Python/Node.js/Go runtime support
   - Resource limits (512MB, 1 CPU)
   - Auto-restart containers

3. **Database** - PostgreSQL with full schema
4. **Cache** - Redis (ready for rate limiting)

### Frontend
- Next.js dashboard
- Authentication UI
- API management
- Marketplace

## Quick Start (3 Services)

### Terminal 1: Database
```powershell
docker-compose up
```

### Terminal 2: API Gateway
```powershell
cd backend/api-gateway
cp .env.example .env
# Edit .env - set JWT_SECRET!
go run main.go
# Runs on http://localhost:8080
```

### Terminal 3: Executor Service
```powershell
cd backend/executor
cp .env.example .env  
go run main.go
# Runs on http://localhost:8081
```

### Terminal 4: Frontend (optional)
```powershell
cd frontend
npm install
npm run dev
# Runs on http://localhost:3000
```

## Testing the Complete Flow

### 1. Sign Up
```bash
POST http://localhost:8080/api/v1/auth/signup
{
    "email": "dev@example.com",
    "password": "Test123!",
    "name": "Developer",
    "role": "developer"
}
# Response includes token
```

### 2. Create API
```bash
POST http://localhost:8080/api/v1/apis
Authorization: Bearer {token}
{
    "name": "hello-api",
    "runtime": "python",
    "visibility": "public"
}
# Response includes API ID
```

### 3. Upload Code
```bash
POST http://localhost:8080/api/v1/apis/{api_id}/upload
Authorization: Bearer {token}
Content-Type: multipart/form-data
# Upload a Python file
```

**Create test.py:**
```python
def handler(request):
    return {"message": "Hello from API Platform!"}
```

### 4. Deploy API  🚀
```bash
POST http://localhost:8080/api/v1/apis/{api_id}/deploy
Authorization: Bearer {token}
# Pulls Docker image, creates container, starts it
```

### 5. Check Status
```bash
GET http://localhost:8080/api/v1/apis/{api_id}/status
# Response: {"status": "running", "container_id": "..."}
```

### 6. Stop API
```bash
POST http://localhost:8080/api/v1/apis/{api_id}/stop
Authorization: Bearer {token}
```

## Using Postman

Import the collection and add new requests:

**Deploy API:**
- Method: POST
- URL: `{{base_url}}/api/v1/apis/{{api_id}}/deploy`
- Headers: `Authorization: Bearer {{auth_token}}`

**Get Status:**
- Method: GET
- URL: `{{base_url}}/api/v1/apis/{{api_id}}/status`

**Stop API:**
- Method: POST
- URL: `{{base_url}}/api/v1/apis/{{api_id}}/stop`
- Headers: `Authorization: Bearer {{auth_token}}`

## Docker Commands

### List Running API Containers
```powershell
docker ps --filter "name=api-"
```

### View Container Logs
```powershell
docker logs api-{api_id}
```

### Stop All API Containers
```powershell
docker ps --filter "name=api-" -q | ForEach-Object { docker stop $_ }
```

### Check Images
```powershell
docker images | Select-String "python|node|golang"
```

## Environment Variables

### API Gateway (.env)
```env
PORT=8080
DATABASE_URL=postgres://apiplatform:dev_password@localhost:5432/apiplatform?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=super-secret-change-me-in-production
EXECUTOR_URL=http://localhost:8081
```

### Executor (.env)
```env
PORT=8081
DATABASE_URL=postgres://apiplatform:dev_password@localhost:5432/apiplatform?sslmode=disable
```

## Troubleshooting

### "Cannot connect to Docker daemon"
- Ensure Docker Desktop is running
- Check: `docker ps`

### "Failed to pull image"
- Check internet connection
- Images: python:3.11-slim, node:18-alpine, golang:1.21-alpine

### "Port already in use"
- API Gateway: Change PORT in .env
- Executor: Change PORT in executor/.env

### "Database connection failed"
- Check: `docker-compose ps`
- Restart: `docker-compose restart postgres`

### "Code path not found"
- Upload code before deploying
- Check uploads/ directory exists

## What's Next (Phase 2)

- [ ] **Request Routing** - Proxy API requests to containers
- [ ] **Analytics** - Track API usage and metrics
- [ ] **Billing** - Usage metering and transactions
- [ ] **Auto-scale** - Multiple container instances
- [ ] **Port Mapping** - Dynamic port allocation
- [ ] **Health Checks** - Container liveness probes

## Production Considerations

Before deploying to production:

1. **Security**
   - Generate strong JWT_SECRET
   - Enable HTTPS
   - Review Docker security
   - Implement rate limiting

2. **Infrastructure**
   - Kubernetes for orchestration
   - Load balancer for API Gateway
   - Managed PostgreSQL
   - Redis cluster

3. **Monitoring**
   - Container metrics
   - API logs aggregation
   - Error tracking
   - Performance monitoring

4. **Scaling**
   - Horizontal API Gateway scaling
   - Container resource optimization
   - Database read replicas
   - CDN for frontend

## Current Limitations

- Containers don't expose external ports yet (Phase 2)
- No request proxying to containers (Phase 2)
- No auto-scaling (Phase 2)
- No usage analytics (Phase 2)
- No billing integration (Phase 3)

## Success Criteria ✅

- [x] User signup/login
- [x] Create and manage APIs
- [x] Upload code
- [x] Deploy to Docker containers
- [x] Stop/status containers
- [x] Database persistence
- [x] Frontend dashboard
- [x] API marketplace

---

**MVP Phase 1: COMPLETE** ✅

The platform can now:
- Register users
- Create APIs
- Upload code
- Deploy to isolated containers
- Manage container lifecycle

**Ready for Phase 2: Request routing and analytics!**
