# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Commands

**Backend (Go):**
```bash
cd backend/api-gateway
go run main.go              # Run API gateway
go test ./...              # Run tests
go build -o bin/gateway    # Build binary
```

**Frontend (React/Next.js):**
```bash
cd frontend
npm install                # Install dependencies
npm run dev               # Development server
npm run build             # Production build
npm run lint              # Lint code
```

## High-level architecture

A developer marketplace and API hosting platform with MENA-first positioning.

**Tech Stack:**
- **Backend:** Go (API gateway, billing, execution orchestration)
- **Frontend:** React/Next.js (marketplace, dashboards)
- **Data:** PostgreSQL (metadata, billing), Redis (caching, rate limiting)
- **Execution:** Docker/Firecracker (sandboxed API containers)
- **Queue:** NATS/RabbitMQ (async jobs, logs)
- **Payments:** Stripe/PayTabs/Moyasar

**Core Components:**

1. **API Gateway** - Request routing, auth, rate limiting
2. **Executor** - Container lifecycle, sandboxed execution
3. **Billing Service** - Usage tracking, payments, developer payouts
4. **Analytics Service** - Logs, metrics, dashboards
5. **Marketplace Frontend** - Public API catalog, developer/consumer dashboards

**MVP Focus (Phase 1):**
- Code upload → containerized endpoint
- Basic auth (API keys)
- Request logging
- Simple analytics dashboard

## Repository structure

```
├── backend/
│   ├── api-gateway/      # Main HTTP router, auth, rate limiting
│   ├── executor/         # Container execution engine
│   ├── billing/          # Usage tracking, payments, payouts
│   ├── analytics/        # Metrics collection, logging
│   └── shared/           # Common libs (DB, auth, config)
├── frontend/             # Next.js marketplace + dashboards
├── docs/                 # Architecture docs, API specs
├── infrastructure/       # Docker, K8s, Terraform configs
└── scripts/              # Dev setup, migrations
```

## Notes for future updates

- When the implementation and tooling are added, document:
  - How to run the service locally (entry command), build, lint/typecheck, test (all and single-test invocation), and any seed/migration steps.
  - Any local dependencies (e.g., queue, DB, tracing backend) and how to start them (e.g., via Docker Compose) if applicable.
- If CLAUDE/Cursor/Copilot rules are introduced, summarize key constraints here.
