# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Commands

- No build/lint/test tooling is configured yet (only `README.md` is present). Update this section once a stack is chosen and scripts are added (e.g., in `Makefile`, `package.json`, `pyproject.toml`, or CI).

## High-level architecture (from README intent)

This repo is a platform to host, route, and execute APIs/functions with isolation, scalability, and observability. The big-picture components to expect as the code is implemented:

- Routing and Versioning
  - HTTP/gRPC entrypoints, endpoint registry, versioned routes, health endpoint.
- Execution Engine (pluggable runtimes)
  - Abstraction to run user code in isolated contexts (process/container/VM/runtime); resource limits; timeouts; cold/warm start policy.
- AuthN/Z, Rate Limiting, Quotas
  - Project- and environment-scoped policies; tenant-aware tokens/keys; per-endpoint limits.
- Async Jobs & Scheduling
  - Queue integration for background/long-running executions; scheduled triggers; retries and DLQ.
- Observability
  - Structured logs, metrics, traces; request/execution correlation IDs; per-tenant dashboards.
- Multi-tenancy & Secrets
  - Projects, environments, and scoped secret management; audit trail.
- Configuration & Persistence
  - Backing stores for routing tables, executions, and policy state; config via files/env.

## Repository structure (planned)

- `src/` — service code and modules (not created yet)
- `docs/` — design docs and runbooks

## Notes for future updates

- When the implementation and tooling are added, document:
  - How to run the service locally (entry command), build, lint/typecheck, test (all and single-test invocation), and any seed/migration steps.
  - Any local dependencies (e.g., queue, DB, tracing backend) and how to start them (e.g., via Docker Compose) if applicable.
- If CLAUDE/Cursor/Copilot rules are introduced, summarize key constraints here.
