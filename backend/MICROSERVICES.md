# Backend Microservices Architecture

## Overview

The backend is organized as a Go workspace with one module per service. Each service exposes an HTTP API using Gin, loads configuration from `.env`, and keeps service-owned code under `internal/`.

The **API Gateway** (port 8080) is the single external entrypoint and uses reverse proxying to route requests to microservices.

## Quick Links

- 📖 [Service Setup Guide](SERVICE_SETUP.md) — How to create and configure a service
- ⚡ [Quick Reference](QUICK_REFERENCE.md) — Common commands and endpoints
- 📋 [Service Details](SERVICES.md) — API endpoints for each service

## Directory Structure

```text
backend/
├── .env.example
├── .env
├── Makefile
├── go.work
├── app/
│   ├── api-gateway/
│   │   ├── cmd/server/main.go
│   │   └── internal/
│   │       ├── config/
│   │       └── http/
│   ├── auth-service/
│   │   ├── cmd/server/main.go
│   │   └── internal/
│   │       ├── auth/
│   │       ├── config/
│   │       └── http/
│   ├── employee-service/
│   ├── branch-service/
│   ├── member-service/
│   ├── membership-service/
│   ├── attendance-service/
│   ├── notification-service/
│   └── payment-service/
├── packages/
│   ├── config/       # shared dotenv/env helpers
│   └── httpserver/   # shared Gin router setup
└── protos/
    └── auth/v1/auth.proto
```

## Services

| Service | Default Port | Resource Route |
| --- | ---: | --- |
| API Gateway | `8080` | `/api/v1/*` |
| Auth | `5001` | `/auth` |
| Employee | `5002` | `/employees` |
| Branch | `5003` | `/branches` |
| Member | `5004` | `/members` |
| Membership | `5005` | `/memberships` |
| Attendance | `5006` | `/attendance` |
| Notification | `5007` | `/notifications` |
| Payment | `5008` | `/payments` |

Every service also exposes:

```text
GET /health
GET /
```

Auth currently exposes:

```text
POST /auth/login
GET  /auth/validate
POST /auth/refresh
```

The gateway also provides convenience routes:

```text
POST /login
GET  /validate
POST /refresh
```

## Configuration

Copy or edit `backend/.env`. Each service loads `.env` from its local directory or the backend root.

Important variables:

```env
GIN_MODE=debug

API_GATEWAY_PORT=8080
AUTH_SERVICE_PORT=5001
EMPLOYEE_SERVICE_PORT=5002
BRANCH_SERVICE_PORT=5003
MEMBER_SERVICE_PORT=5004
MEMBERSHIP_SERVICE_PORT=5005
ATTENDANCE_SERVICE_PORT=5006
NOTIFICATION_SERVICE_PORT=5007
PAYMENT_SERVICE_PORT=5008

AUTH_SERVICE_URL=http://localhost:5001
EMPLOYEE_SERVICE_URL=http://localhost:5002
BRANCH_SERVICE_URL=http://localhost:5003
MEMBER_SERVICE_URL=http://localhost:5004
MEMBERSHIP_SERVICE_URL=http://localhost:5005
ATTENDANCE_SERVICE_URL=http://localhost:5006
NOTIFICATION_SERVICE_URL=http://localhost:5007
PAYMENT_SERVICE_URL=http://localhost:5008
```

## Development

Install or refresh dependencies:

```bash
cd backend
make deps
```

Build all services:

```bash
cd backend
make build-all
```

Run one service:

```bash
cd backend
make run-auth-service
make run-api-gateway
```

Example auth request through the gateway:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass"}'
```

Example proxied service request:

```bash
curl http://localhost:8080/api/v1/employees/employees
```

## Adding a Service

1. Create `backend/app/<name>-service`.
2. Add `cmd/server/main.go`.
3. Put service-owned code in `internal/config`, `internal/http`, and domain-specific packages.
4. Add the module to `go.work`.
5. Add env vars to `.env.example`.
6. Add the service to `SERVICES` in `Makefile`.
7. Add the service URL to the gateway config if it should be externally proxied.

📖 **See [SERVICE_SETUP.md](SERVICE_SETUP.md) for detailed step-by-step instructions.**

## Documentation

| File | Purpose |
|------|---------|
| [MICROSERVICES.md](MICROSERVICES.md) | Architecture & overview (this file) |
| [SERVICE_SETUP.md](SERVICE_SETUP.md) | How to set up a new microservice |
| [QUICK_REFERENCE.md](QUICK_REFERENCE.md) | Commands, endpoints, common issues |
| [SERVICES.md](SERVICES.md) | API documentation per service |
