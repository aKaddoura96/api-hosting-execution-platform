# Getting Started

This guide will help you set up and run the API Platform locally.

## Prerequisites

- **Go 1.21+** ([Download](https://go.dev/dl/))
- **Node.js 18+** ([Download](https://nodejs.org/))
- **Docker & Docker Compose** ([Download](https://www.docker.com/get-started))
- **PostgreSQL client** (psql) for running migrations

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/aKaddoura96/api-hosting-execution-platform.git
cd api-hosting-execution-platform
```

### 2. Start Database Services

```bash
# Start PostgreSQL and Redis with Docker Compose
docker-compose up -d

# Verify services are running
docker-compose ps
```

### 3. Run Database Migrations

**On Windows (PowerShell):**
```powershell
.\scripts\run-migrations.ps1
```

**On Linux/Mac:**
```bash
chmod +x scripts/run-migrations.sh
./scripts/run-migrations.sh
```

### 4. Set Up Backend

```bash
cd backend/api-gateway

# Copy environment file
cp .env.example .env

# Edit .env and set JWT_SECRET to a secure random string
# JWT_SECRET=your-secure-random-string-here

# Download dependencies
go mod download

# Download shared module dependencies
cd ../shared
go mod download
cd ../api-gateway

# Run the API Gateway
go run main.go
```

The API Gateway will start on `http://localhost:8080`

### 5. Set Up Frontend

Open a new terminal:

```bash
cd frontend

# Install dependencies
npm install

# Copy environment file
cp .env.example .env.local

# Run development server
npm run dev
```

The frontend will start on `http://localhost:3000`

## Usage

### 1. Create an Account

1. Navigate to `http://localhost:3000`
2. Click "Get Started" or "Sign up"
3. Choose your role:
   - **Developer**: Host and monetize APIs
   - **Consumer**: Use APIs from the marketplace
4. Complete registration

### 2. Create Your First API (Developers)

1. Sign in and go to Dashboard
2. Click "+ Create API"
3. Fill in:
   - **Name**: e.g., "weather-api"
   - **Description**: What your API does
   - **Runtime**: python, nodejs, or go
   - **Visibility**: private, public, or paid
4. Click "Create"

### 3. Browse the Marketplace (Consumers)

1. Navigate to "Marketplace"
2. Browse public APIs
3. View API details, endpoints, and documentation
4. Use the provided endpoint URLs to integrate

## Project Structure

```
├── backend/
│   ├── api-gateway/       # Main HTTP server
│   ├── executor/          # Container execution (TODO)
│   ├── billing/           # Usage tracking (TODO)
│   ├── analytics/         # Metrics collection (TODO)
│   └── shared/            # Common libraries
├── frontend/              # Next.js React app
│   ├── app/               # Pages (Next.js 14 app router)
│   └── lib/               # API client, auth context
├── scripts/               # Database migrations
└── docker-compose.yml     # Local dev services
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/signup` - Create account
- `POST /api/v1/auth/login` - Sign in
- `GET /api/v1/auth/me` - Get current user (requires auth)

### API Management (Requires Auth)
- `GET /api/v1/apis` - List my APIs
- `POST /api/v1/apis` - Create new API
- `GET /api/v1/apis/:id` - Get API details
- `DELETE /api/v1/apis/:id` - Delete API
- `POST /api/v1/apis/:id/upload` - Upload code

### Marketplace (Public)
- `GET /api/v1/marketplace/apis` - List public APIs
- `GET /api/v1/marketplace/apis/:id` - Get public API details

## Environment Variables

### Backend (.env)
```env
PORT=8080
DATABASE_URL=postgres://apiplatform:dev_password@localhost:5432/apiplatform?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=change-this-to-a-secure-random-string
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Testing the Platform

### 1. Test Authentication
```bash
# Sign up
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User","role":"developer"}'

# Sign in
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 2. Test API Creation
```bash
# Create an API (use token from login)
curl -X POST http://localhost:8080/api/v1/apis \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"name":"test-api","description":"My test API","runtime":"python","visibility":"public"}'
```

## Troubleshooting

### Database Connection Issues
- Ensure Docker containers are running: `docker-compose ps`
- Check database logs: `docker-compose logs postgres`
- Verify connection string in `.env`

### Frontend Not Connecting to Backend
- Verify API Gateway is running on port 8080
- Check CORS settings in `backend/api-gateway/main.go`
- Ensure `NEXT_PUBLIC_API_URL` is set correctly

### Migration Failures
- Ensure PostgreSQL is running
- Check database credentials match docker-compose.yml
- Manually connect: `psql postgres://apiplatform:dev_password@localhost:5432/apiplatform`

## Next Steps

- [x] ✅ Authentication & API management
- [ ] Implement code execution engine
- [ ] Add analytics and usage tracking
- [ ] Integrate billing system
- [ ] Deploy to production

## Need Help?

- Check the [README.md](README.md) for architecture details
- Review [WARP.md](WARP.md) for development commands
- Open an issue on GitHub

---

**Built with ❤️ for the MENA developer community**
