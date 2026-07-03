# Backend Quick Reference

## Common Commands

### Setup

```bash
cd backend
make deps           # Install all dependencies
make build-all      # Build all services
```

### Running Services

**Terminal 1 — Gateway (required for proxying):**
```bash
make run-api-gateway
```

**Terminal 2 — Auth Service:**
```bash
make run-auth-service
```

**Terminal 3 — Employee Service:**
```bash
make run-employee-service
```

**Any other service:**
```bash
make run-{service-name}
```

### Building

```bash
make build-all              # Build all services
make build-auth-service     # Build one service
```

## API Endpoints

### Via Gateway (Recommended)

**Auth (Port 8080):**
```bash
POST   http://localhost:8080/login
GET    http://localhost:8080/validate?token=xxx
POST   http://localhost:8080/refresh
```

**Employee (Port 8080):**
```bash
GET    http://localhost:8080/api/v1/employee/employees
POST   http://localhost:8080/api/v1/employee/employees
GET    http://localhost:8080/api/v1/employee/employees/:id
```

**Payment (Port 8080):**
```bash
GET    http://localhost:8080/api/v1/payment/payments
POST   http://localhost:8080/api/v1/payment/payments
```

**Other Services:**
```bash
http://localhost:8080/api/v1/{service}/*
```

### Direct (Debug Only)

**Auth Service (Port 5001):**
```bash
POST   http://localhost:5001/auth/login
GET    http://localhost:5001/auth/validate?token=xxx
POST   http://localhost:5001/auth/refresh
```

**Employee Service (Port 5002):**
```bash
GET    http://localhost:5002/employees
POST   http://localhost:5002/employees
```

## Service Ports

| Service | Port |
|---------|------|
| API Gateway | 8080 |
| Auth | 5001 |
| Employee | 5002 |
| Branch | 5003 |
| Member | 5004 |
| Membership | 5005 |
| Attendance | 5006 |
| Notification | 5007 |
| Payment | 5008 |

## Configuration

**Environment File:** `backend/.env`

**Key Variables:**
```env
GIN_MODE=debug                                          # or "release"
API_GATEWAY_PORT=8080
AUTH_SERVICE_URL=http://localhost:5001
EMPLOYEE_SERVICE_URL=http://localhost:5002
AUTH_DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth_db
```

**Per-Service Database:**
```env
{SERVICE_UPPER}_DATABASE_URL=postgres://user:pass@host/db_name
```

## Testing Example

### 1. Start Services

```bash
# Terminal 1
make run-api-gateway

# Terminal 2
make run-auth-service

# Terminal 3
make run-employee-service
```

### 2. Test Auth

```bash
# Via gateway
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass"}'

# Direct
curl -X POST http://localhost:5001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass"}'
```

### 3. Test Employee

```bash
# Via gateway
curl http://localhost:8080/api/v1/employee/employees

# Direct
curl http://localhost:5002/employees
```

### 4. Create Employee

```bash
curl -X POST http://localhost:8080/api/v1/employee/employees \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

## File Structure per Service

```
app/{service}-service/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── config/config.go         # Load env config
│   ├── http/routes.go           # HTTP handlers
│   └── {service}/service.go     # Business logic
├── go.mod
└── go.sum
```

## Common Issues

| Issue | Solution |
|-------|----------|
| `connection refused` | Service not running; check `lsof -i :PORT` |
| `502 Bad Gateway` | Service URL wrong in `.env` or service down |
| `database connection failed` | Check `DATABASE_URL` env var or PostgreSQL |
| Port already in use | Kill process: `lsof -i :PORT \| grep LISTEN \| awk '{print $2}' \| xargs kill -9` |

## Debug Tips

**Check if port is in use:**
```bash
lsof -i :5001
```

**Kill process on port:**
```bash
lsof -i :5001 | grep LISTEN | awk '{print $2}' | xargs kill -9
```

**View all service ports:**
```bash
lsof -i | grep LISTEN | grep -E ":(5|8)"
```

**Health check:**
```bash
curl http://localhost:8080/health
curl http://localhost:5001/health
```

## Adding a New Service

1. Create directory: `mkdir -p app/myservice-service/cmd/server`
2. Copy structure from `auth-service` or `employee-service`
3. Implement config, routes, and business logic
4. Add to `Makefile` `SERVICES` variable
5. Add env vars to `.env` and `.env.example`
6. Run: `make run-myservice-service`
7. Test via gateway: `curl http://localhost:8080/api/v1/myservice/`

## Gateway Flow

```
Client Request
    ↓
Gateway (8080)
    ↓
Route Matching (/api/v1/{service}/*)
    ↓
Lookup Service URL from Config
    ↓
Reverse Proxy to Service
    ↓
Service Response
    ↓
Return to Client
```

## See Also

- [SERVICE_SETUP.md](SERVICE_SETUP.md) — Detailed setup guide
- [MICROSERVICES.md](MICROSERVICES.md) — Architecture overview
