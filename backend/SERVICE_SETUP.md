# Service Setup & Integration Guide

This guide explains how to set up each microservice and integrate it with the API Gateway.

## Table of Contents

1. [Service Architecture](#service-architecture)
2. [Service Structure](#service-structure)
3. [Setting Up a New Service](#setting-up-a-new-service)
4. [Service Details](#service-details)
5. [API Gateway Integration](#api-gateway-integration)
6. [Environment Configuration](#environment-configuration)
7. [Testing Services](#testing-services)

---

## Service Architecture

```
┌─────────────────────────────────────────────┐
│         External Clients (Port 8080)        │
└────────────────┬────────────────────────────┘
                 │ HTTP/REST
                 ▼
    ┌────────────────────────────────┐
    │      API Gateway               │
    │   (Reverse Proxy & Router)     │
    │   Port: 8080                   │
    └─┬──────────┬──────────┬────────┘
      │          │          │
    ◄─┴──────────┴──────────┴───────► HTTP calls to services
      │          │          │
      ▼          ▼          ▼
  ┌────────┐ ┌────────┐ ┌────────┐
  │ Auth   │ │Employee│ │Payment │
  │:5001   │ │:5002   │ │:5008   │
  └────────┘ └────────┘ └────────┘
      │          │          │
      └──────────┴──────────┘
           DB (Optional)
```

**Key Points:**
- **API Gateway** (`localhost:8080`) is the single entry point
- Each service runs on its own port
- Services are called via HTTP reverse proxy
- Services are independent — each has its own DB

---

## Service Structure

Every service follows this standardized structure:

```
app/{service}-service/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Load environment config
│   ├── http/
│   │   └── routes.go            # HTTP handlers & routes
│   ├── {domain}/
│   │   └── service.go           # Business logic
│   └── database/                # (optional) DB models
├── go.mod
├── go.sum
└── README.md
```

### File Responsibilities

**`cmd/server/main.go`** — Bootstrap
- Load config from environment
- Connect to database (if needed)
- Create router & register routes
- Start HTTP server

**`internal/config/config.go`** — Configuration
- Parse environment variables
- Provide defaults
- Validate required config

**`internal/http/routes.go`** — HTTP API
- Define request/response types
- Register route handlers with Gin
- Handle serialization/validation

**`internal/{domain}/service.go`** — Business Logic
- Implement core domain logic
- Use database/external calls
- Return errors, not HTTP status codes

---

## Setting Up a New Service

### 1. Create Service Directory

```bash
cd backend
mkdir -p app/myservice-service/cmd/server
mkdir -p app/myservice-service/internal/config
mkdir -p app/myservice-service/internal/http
mkdir -p app/myservice-service/internal/myservice
```

### 2. Create `go.mod`

```bash
cd app/myservice-service
go mod init github.com/panduputragit/gym/backend/app/myservice-service
go get google.golang.org/grpc
go get github.com/gin-gonic/gin
```

### 3. Implement `internal/config/config.go`

```go
package config

import (
	"os"
	"strconv"
)

type Config struct {
	Name        string
	Port        string
	GinMode     string
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		Name:        "myservice",
		Port:        os.Getenv("MYSERVICE_SERVICE_PORT"),
		GinMode:     os.Getenv("GIN_MODE"),
		DatabaseURL: os.Getenv("MYSERVICE_DATABASE_URL"),
	}
}
```

### 4. Implement `internal/http/routes.go`

```go
package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/panduputragit/gym/backend/app/myservice-service/internal/myservice"
)

type Handler struct {
	service *myservice.Service
}

func NewHandler(svc *myservice.Service) *Handler {
	return &Handler{service: svc}
}

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	group := router.Group("/myresource")
	group.GET("", handler.ListAll)
	group.POST("", handler.Create)
	group.GET("/:id", handler.Get)
}

func (h *Handler) ListAll(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) Create(c *gin.Context) {
	// Parse request, call service, return response
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	item, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}
```

### 5. Implement `internal/myservice/service.go`

```go
package myservice

import "context"

type Service struct {
	// Inject dependencies: DB, cache, other services
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) List(ctx context.Context) ([]interface{}, error) {
	// TODO: Query DB or call other services
	return []interface{}{}, nil
}

func (s *Service) Get(ctx context.Context, id string) (interface{}, error) {
	// TODO: Fetch one item
	return nil, nil
}
```

### 6. Implement `cmd/server/main.go`

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/panduputragit/gym/backend/app/myservice-service/internal/config"
	myhttp "github.com/panduputragit/gym/backend/app/myservice-service/internal/http"
	"github.com/panduputragit/gym/backend/app/myservice-service/internal/myservice"
	"github.com/panduputragit/gym/backend/packages/database"
	"github.com/panduputragit/gym/backend/packages/httpserver"
)

func main() {
	cfg := config.Load()

	// Optional: Connect to DB
	db, err := database.ConnectOptional(context.Background(), database.Config{URL: cfg.DatabaseURL})
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	if db != nil {
		defer db.Close()
		fmt.Printf("%s connected to database\n", cfg.Name)
	}

	router := httpserver.NewRouter(cfg.Name, cfg.GinMode)
	myhttp.RegisterRoutes(router, myhttp.NewHandler(myservice.NewService()))

	addr := ":" + cfg.Port
	fmt.Printf("%s listening on %s\n", cfg.Name, addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
```

### 7. Add to `Makefile`

Edit `backend/Makefile` and add your service to `SERVICES`:

```makefile
SERVICES := api-gateway auth-service employee-service myservice-service ...
```

### 8. Add to `go.work`

Add your module to `backend/go.work`:

```
go 1.26.4

use (
	./packages
	./app/api-gateway
	./app/auth-service
	./app/myservice-service
	...
)
```

### 9. Update `.env` and `.env.example`

Add configuration:

```env
MYSERVICE_SERVICE_PORT=5009
MYSERVICE_SERVICE_URL=http://localhost:5009
MYSERVICE_DATABASE_URL=postgres://postgres:postgres@localhost:5432/myservice_db?sslmode=disable
```

### 10. Test

```bash
cd backend
make deps           # Install dependencies
make run-myservice-service  # Start the service
curl http://localhost:5009/myresource  # Test directly
```

---

## Service Details

### Auth Service (Port 5001)

**Endpoints:**
- `POST /auth/login` — User login
- `GET /auth/validate` — Validate token
- `POST /auth/refresh` — Refresh access token

**Gateway Aliases:**
- `POST /login` → `POST /auth/login`
- `GET /validate` → `GET /auth/validate`
- `POST /refresh` → `POST /auth/refresh`

**Example:**

```bash
# Direct call to service
curl -X POST http://localhost:5001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass"}'

# Via gateway (preferred)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass"}'
```

---

### Employee Service (Port 5002)

**Endpoints:**
- `GET /employees` — List all employees
- `POST /employees` — Create employee
- `GET /employees/:id` — Get employee by ID
- `PUT /employees/:id` — Update employee
- `DELETE /employees/:id` — Delete employee

**Gateway Route:**
- `/api/v1/employee/*` → `http://localhost:5002/*`

**Example:**

```bash
# Via gateway
curl http://localhost:8080/api/v1/employee/employees
curl -X POST http://localhost:8080/api/v1/employee/employees \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Direct
curl http://localhost:5002/employees
```

---

### Payment Service (Port 5008)

**Endpoints:**
- `GET /payments` — List payments
- `POST /payments` — Create payment
- `GET /payments/:id` — Get payment status

**Gateway Route:**
- `/api/v1/payment/*` → `http://localhost:5008/*`

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/payment/payments \
  -H "Content-Type: application/json" \
  -d '{"amount":100,"currency":"USD"}'
```

---

### Other Services

| Service | Port | Route | Key Endpoints |
|---------|------|-------|---------------|
| Branch | 5003 | `/api/v1/branch/` | `/branches`, `/branches/:id` |
| Member | 5004 | `/api/v1/member/` | `/members`, `/members/:id` |
| Membership | 5005 | `/api/v1/membership/` | `/memberships`, `/memberships/:id` |
| Attendance | 5006 | `/api/v1/attendance/` | `/attendance`, `/attendance/:id` |
| Notification | 5007 | `/api/v1/notification/` | `/notifications`, `/notifications/:id` |

---

## API Gateway Integration

### How It Works

The API Gateway (`api-gateway/internal/http/proxy.go`) uses **reverse proxy** to forward requests:

1. Client sends: `GET http://localhost:8080/api/v1/employee/employees`
2. Gateway receives request
3. Gateway strips `/api/v1/employee` prefix → path becomes `/employees`
4. Gateway forwards to `http://localhost:5002/employees`
5. Service responds
6. Gateway returns response to client

### Proxy Rules

#### 1. Service Proxying

Service URLs from config are proxied under `/api/v1/{service}/`:

```
/api/v1/employee/*  → http://localhost:5002/*
/api/v1/payment/*   → http://localhost:5008/*
/api/v1/branch/*    → http://localhost:5003/*
...
```

#### 2. Auth Shortcuts

Special routes for auth service:

```
POST /login    → POST http://localhost:5001/auth/login
GET /validate  → GET http://localhost:5001/auth/validate
POST /refresh  → POST http://localhost:5001/auth/refresh
```

### Adding a Service to Gateway

1. **Register in `go.mod` and `go.work`** — Service must be running
2. **Set env variable** in `.env`:
   ```env
   MYSERVICE_SERVICE_URL=http://localhost:5009
   ```
3. **Gateway loads config** in `internal/config/config.go`:
   ```go
   func (c *Config) ServiceURLs() map[string]string {
       return map[string]string{
           "auth":       c.AuthServiceURL,
           "employee":   c.EmployeeServiceURL,
           "myservice":  c.MyserviceServiceURL,  // Add here
       }
   }
   ```

### Example Flow

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/employee/employees \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice"}'
```

**Gateway Processing:**
```
1. Receive: POST /api/v1/employee/employees
2. Match route: /api/v1/{service}/*
3. Extract service: "employee"
4. Look up URL: serviceURLs["employee"] = "http://localhost:5002"
5. Strip prefix: /api/v1/employee/employees → /employees
6. Forward: POST http://localhost:5002/employees
7. Receive response from service
8. Return to client
```

---

## Environment Configuration

### Structure

All services use a shared env file: `backend/.env`

Each service loads this file and extracts its config using `Load()` in `internal/config/config.go`.

### Port Convention

```
API Gateway:     8080
Auth Service:    5001
Employee:        5002
Branch:          5003
Member:          5004
Membership:      5005
Attendance:      5006
Notification:    5007
Payment:         5008
```

### Database Convention

Each service gets its own database:

```env
AUTH_DATABASE_URL=postgres://user:pass@localhost/auth_db
EMPLOYEE_DATABASE_URL=postgres://user:pass@localhost/employee_db
PAYMENT_DATABASE_URL=postgres://user:pass@localhost/payment_db
```

### Example `.env`

```env
GIN_MODE=debug

# Gateway
API_GATEWAY_PORT=8080

# Service Ports
AUTH_SERVICE_PORT=5001
EMPLOYEE_SERVICE_PORT=5002
BRANCH_SERVICE_PORT=5003
MEMBER_SERVICE_PORT=5004
MEMBERSHIP_SERVICE_PORT=5005
ATTENDANCE_SERVICE_PORT=5006
NOTIFICATION_SERVICE_PORT=5007
PAYMENT_SERVICE_PORT=5008

# Service URLs (for Gateway to call)
AUTH_SERVICE_URL=http://localhost:5001
EMPLOYEE_SERVICE_URL=http://localhost:5002
BRANCH_SERVICE_URL=http://localhost:5003
MEMBER_SERVICE_URL=http://localhost:5004
MEMBERSHIP_SERVICE_URL=http://localhost:5005
ATTENDANCE_SERVICE_URL=http://localhost:5006
NOTIFICATION_SERVICE_URL=http://localhost:5007
PAYMENT_SERVICE_URL=http://localhost:5008

# Databases
AUTH_DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth_db
EMPLOYEE_DATABASE_URL=postgres://postgres:postgres@localhost:5432/employee_db
BRANCH_DATABASE_URL=postgres://postgres:postgres@localhost:5432/branch_db
MEMBER_DATABASE_URL=postgres://postgres:postgres@localhost:5432/member_db
MEMBERSHIP_DATABASE_URL=postgres://postgres:postgres@localhost:5432/membership_db
ATTENDANCE_DATABASE_URL=postgres://postgres:postgres@localhost:5432/attendance_db
NOTIFICATION_DATABASE_URL=postgres://postgres:postgres@localhost:5432/notification_db
PAYMENT_DATABASE_URL=postgres://postgres:postgres@localhost:5432/payment_db
```

---

## Testing Services

### 1. Start All Services

**Terminal 1 — Gateway:**
```bash
cd backend && make run-api-gateway
```

**Terminal 2 — Auth:**
```bash
cd backend && make run-auth-service
```

**Terminal 3 — Employee:**
```bash
cd backend && make run-employee-service
```

### 2. Test Auth Service

**Direct:**
```bash
curl -X POST http://localhost:5001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

**Via Gateway:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password"}'
```

### 3. Test Employee Service

**Direct:**
```bash
curl http://localhost:5002/employees
```

**Via Gateway:**
```bash
curl http://localhost:8080/api/v1/employee/employees
```

### 4. Test All Service Health

Each service exposes:
- `GET /` — Service info
- `GET /health` — Health check

```bash
curl http://localhost:8080/health        # Gateway health
curl http://localhost:5001/health        # Auth health
curl http://localhost:5002/health        # Employee health
```

### 5. Debugging

**Check if service is running:**
```bash
lsof -i :5001  # Check if auth service port is in use
lsof -i :8080  # Check if gateway port is in use
```

**View logs:**
Services print to stdout. Redirect to file for analysis:
```bash
make run-auth-service > logs/auth.log 2>&1 &
```

**Verify gateway config:**
```bash
curl http://localhost:8080/  # Should respond
```

---

## Common Issues

### Service won't connect

**Error:** `connection refused`

**Solution:**
- Ensure service is running: `lsof -i :PORT`
- Check `.env` has correct service URLs
- Verify port is not blocked by firewall

### Gateway returns 502 Bad Gateway

**Cause:** Service is down or URL is wrong

**Solution:**
- Check `.env` SERVICE_URL matches running service
- Restart service: `make run-{service}`
- Check service logs for errors

### Database connection fails

**Error:** `connection refused` (DB)

**Solution:**
- Ensure PostgreSQL is running
- Check DATABASE_URL in `.env`
- Verify database exists

---

## Checklist for New Service

- [ ] Create service directory structure
- [ ] Implement `config.go`
- [ ] Implement `routes.go` with HTTP handlers
- [ ] Implement `service.go` with business logic
- [ ] Implement `main.go` (cmd/server)
- [ ] Add to `go.work`
- [ ] Add to Makefile `SERVICES`
- [ ] Add to `.env` and `.env.example`
- [ ] Add gateway config in `api-gateway/internal/config/config.go`
- [ ] Test direct: `curl http://localhost:{PORT}/`
- [ ] Test via gateway: `curl http://localhost:8080/api/v1/{service}/`
- [ ] Add to documentation with port & endpoints
