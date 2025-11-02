# API Hosting & Execution Platform 🚀

A MENA-first platform where developers can host REST APIs or scripts serverless-style, manage versions, and optionally sell access — with built-in billing, authentication, and analytics.

## 🎯 Vision

Enable developers to:
- **Host APIs**: Deploy Python, Node.js, or Go code as containerized APIs
- **Monetize**: Sell API access with built-in billing and subscription management
- **Track Usage**: Real-time analytics and execution metrics
- **Scale**: Serverless-style automatic scaling and resource management

## ✨ Features

### Currently Available ✅
- ✅ User authentication & authorization (JWT)
- ✅ API creation & management
- ✅ Code upload & storage
- ✅ Container deployment (Python/Node/Go runtimes)
- ✅ Deployment status tracking
- ✅ Public marketplace
- ✅ Analytics infrastructure
- ✅ RESTful API Gateway

### Coming Soon 🚧
- Request proxying to deployed containers
- Usage-based billing & payments
- Consumer dashboard
- Rate limiting & quotas

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                    Frontend                          │
│              (Next.js 14 + React)                   │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│               API Gateway (Go)                       │
│          Port 8080 - HTTP/REST                      │
│    Auth • API Management • Deployment               │
└──┬──────────┬──────────┬─────────────┬─────────────┘
   │          │          │             │
┌──▼────┐ ┌──▼─────┐ ┌─▼────────┐ ┌──▼──────────┐
│Executor│ │Analytics│ │PostgreSQL│ │   Redis     │
│(Go)    │ │  (Go)   │ │          │ │             │
│8081    │ │  8082   │ │   5432   │ │    6379     │
└────────┘ └─────────┘ └──────────┘ └─────────────┘
```

## 🚀 Quick Start

### Prerequisites
- **Docker Desktop** (that's it!)

### Start Everything with Docker

```powershell
# One command to start everything!
.\docker-start.ps1

# Or use docker-compose directly
docker-compose up --build

# Stop all services
docker-compose down
```

### Development Mode (Without Docker)

```powershell
# If you want to run services locally for development
.\scripts\start.ps1

# Stop all services
.\scripts\stop.ps1
```

### Access the Platform
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **API Docs**: See POSTMAN_TESTING_GUIDE.md

## 📚 Documentation

- **[Getting Started](./GETTING_STARTED.md)** - Setup & first API
- **[Architecture](./ARCHITECTURE.md)** - System design
- **[Development](./DEVELOPMENT.md)** - Contributing guide
- **[Deployment](./DEPLOYMENT_GUIDE.md)** - Production deployment
- **[API Testing](./POSTMAN_TESTING_GUIDE.md)** - Postman guide

## 🧪 Quick Test

```bash
# 1. Signup
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@example.com","password":"pass123","name":"Dev","role":"developer"}'

# 2. Login (save token)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"dev@example.com","password":"pass123"}'

# 3. Create API
curl -X POST http://localhost:8080/api/v1/apis \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"my-api","description":"Test","runtime":"python","visibility":"public"}'
```

## 📊 Tech Stack

**Backend**: Go, PostgreSQL, Redis, Docker  
**Frontend**: Next.js 14, React, Tailwind CSS  
**DevOps**: Docker Compose

## 🤝 Contributing

See [DEVELOPMENT.md](./DEVELOPMENT.md) for setup and contribution guidelines.

## 📝 License

MIT License

---

**Built with ❤️ for the MENA developer community**
