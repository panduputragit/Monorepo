# Services Documentation

Detailed documentation for each microservice.

## Table of Contents

1. [API Gateway](#api-gateway)
2. [Auth Service](#auth-service)
3. [Employee Service](#employee-service)
4. [Branch Service](#branch-service)
5. [Member Service](#member-service)
6. [Membership Service](#membership-service)
7. [Attendance Service](#attendance-service)
8. [Notification Service](#notification-service)
9. [Payment Service](#payment-service)

---

## API Gateway

**Port:** 8080  
**Module:** `app/api-gateway`  
**Role:** HTTP reverse proxy and request router

### Overview

The API Gateway is the single entry point for all external HTTP requests. It:
- Routes requests to appropriate microservices
- Performs reverse proxying via HTTP
- Provides shortcut routes for frequently used auth endpoints
- Loads service URLs from environment configuration

### Architecture

```
Client → :8080 → Gateway → Service URL Lookup → Reverse Proxy → Service
```

### Endpoints

All routes are prefixed with `/api/v1/` except for auth shortcuts:

**Auth Shortcuts (no prefix):**
```
POST   /login       → POST   http://auth:5001/auth/login
GET    /validate    → GET    http://auth:5001/auth/validate
POST   /refresh     → POST   http://auth:5001/auth/refresh
```

**Service Proxying:**
```
/api/v1/{service}/*  → http://{service-host}:{service-port}/*
```

### Examples

```bash
# Auth via shortcut
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"pass"}'

# Employee via proxy
curl http://localhost:8080/api/v1/employee/employees

# Payment via proxy
curl -X POST http://localhost:8080/api/v1/payment/payments \
  -H "Content-Type: application/json" \
  -d '{"amount":100}'
```

### Configuration

**Environment Variables:**
```env
API_GATEWAY_PORT=8080
GIN_MODE=debug
AUTH_SERVICE_URL=http://localhost:5001
EMPLOYEE_SERVICE_URL=http://localhost:5002
BRANCH_SERVICE_URL=http://localhost:5003
MEMBER_SERVICE_URL=http://localhost:5004
MEMBERSHIP_SERVICE_URL=http://localhost:5005
ATTENDANCE_SERVICE_URL=http://localhost:5006
NOTIFICATION_SERVICE_URL=http://localhost:5007
PAYMENT_SERVICE_URL=http://localhost:5008
```

### Implementation Details

**Key Files:**
- `cmd/server/main.go` — Loads config and starts server
- `internal/config/config.go` — Parses service URLs from env
- `internal/http/proxy.go` — Reverse proxy logic

**How It Routes:**
1. Receive request: `POST /api/v1/employee/employees`
2. Parse path: service = `employee`
3. Look up URL: `serviceURLs["employee"]` = `http://localhost:5002`
4. Strip prefix: path becomes `/employees`
5. Forward to `http://localhost:5002/employees`
6. Return service response to client

---

## Auth Service

**Port:** 5001  
**Module:** `app/auth-service`  
**Database:** `AUTH_DATABASE_URL` (optional)

### Overview

Handles user authentication, token management, and validation.

### Endpoints

All routes are under `/auth`:

```
POST   /auth/login              # User login
GET    /auth/validate           # Validate token
POST   /auth/refresh            # Refresh access token
GET    /                        # Service info
GET    /health                  # Health check
```

### API Reference

#### Login

**Request:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600
}
```

**Error (401):**
```json
{"error": "invalid credentials"}
```

#### Validate Token

**Request:**
```bash
# Via header
curl http://localhost:8080/validate \
  -H "Authorization: Bearer {token}"

# Via query param
curl http://localhost:8080/validate?token={token}
```

**Response (200):**
```json
{
  "valid": true,
  "user_id": "123",
  "email": "user@example.com"
}
```

#### Refresh Token

**Request:**
```bash
curl -X POST http://localhost:8080/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }'
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600
}
```

### Configuration

**Environment Variables:**
```env
AUTH_SERVICE_PORT=5001
AUTH_SERVICE_URL=http://localhost:5001
AUTH_DATABASE_URL=postgres://user:pass@localhost:5432/auth_db
GIN_MODE=debug
```

### Implementation Details

**Key Files:**
- `internal/auth/service.go` — Authentication business logic
- `internal/http/routes.go` — HTTP handlers
- `internal/config/config.go` — Configuration

**Service Methods:**
- `Login(ctx, email, password)` — Authenticate user
- `ValidateToken(ctx, token)` — Verify token validity
- `Refresh(ctx, refreshToken)` — Issue new access token

---

## Employee Service

**Port:** 5002  
**Module:** `app/employee-service`  
**Database:** `EMPLOYEE_DATABASE_URL` (optional)

### Overview

Manages employee records and employee-related operations.

### Endpoints

All routes are under `/employees` (via gateway: `/api/v1/employee/`):

```
GET    /employees               # List all employees
POST   /employees               # Create employee
GET    /employees/:id           # Get employee by ID
PUT    /employees/:id           # Update employee
DELETE /employees/:id           # Delete employee
GET    /                        # Service info
GET    /health                  # Health check
```

### API Reference

#### List Employees

**Request:**
```bash
curl http://localhost:8080/api/v1/employee/employees
```

**Response (200):**
```json
{
  "data": [
    {
      "id": "1",
      "name": "John Doe",
      "email": "john@example.com",
      "department": "Engineering"
    }
  ]
}
```

#### Create Employee

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/employee/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith",
    "email": "alice@example.com",
    "department": "Sales"
  }'
```

**Response (201):**
```json
{
  "id": "2",
  "name": "Alice Smith",
  "email": "alice@example.com",
  "department": "Sales"
}
```

#### Get Employee

**Request:**
```bash
curl http://localhost:8080/api/v1/employee/employees/1
```

**Response (200):**
```json
{
  "id": "1",
  "name": "John Doe",
  "email": "john@example.com",
  "department": "Engineering"
}
```

**Error (404):**
```json
{"error": "employee not found"}
```

#### Update Employee

**Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/employee/employees/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "department": "Management"
  }'
```

**Response (200):**
```json
{
  "id": "1",
  "name": "John Smith",
  "email": "john@example.com",
  "department": "Management"
}
```

#### Delete Employee

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/employee/employees/1
```

**Response (204):** No content

### Configuration

**Environment Variables:**
```env
EMPLOYEE_SERVICE_PORT=5002
EMPLOYEE_SERVICE_URL=http://localhost:5002
EMPLOYEE_DATABASE_URL=postgres://user:pass@localhost:5432/employee_db
GIN_MODE=debug
```

---

## Branch Service

**Port:** 5003  
**Module:** `app/branch-service`  
**Database:** `BRANCH_DATABASE_URL` (optional)

### Overview

Manages gym branch/location information.

### Endpoints

```
GET    /branches               # List all branches
POST   /branches               # Create branch
GET    /branches/:id           # Get branch
PUT    /branches/:id           # Update branch
DELETE /branches/:id           # Delete branch
```

### Gateway Route

```
/api/v1/branch/*  → http://localhost:5003/*
```

### Example

```bash
curl http://localhost:8080/api/v1/branch/branches
curl -X POST http://localhost:8080/api/v1/branch/branches \
  -H "Content-Type: application/json" \
  -d '{"name":"Central Branch","city":"Jakarta"}'
```

---

## Member Service

**Port:** 5004  
**Module:** `app/member-service`  
**Database:** `MEMBER_DATABASE_URL` (optional)

### Overview

Manages gym member information and member profiles.

### Endpoints

```
GET    /members               # List all members
POST   /members               # Register member
GET    /members/:id           # Get member
PUT    /members/:id           # Update member
DELETE /members/:id           # Delete member
```

### Gateway Route

```
/api/v1/member/*  → http://localhost:5004/*
```

---

## Membership Service

**Port:** 5005  
**Module:** `app/membership-service`  
**Database:** `MEMBERSHIP_DATABASE_URL` (optional)

### Overview

Manages membership plans, pricing, and subscriptions.

### Endpoints

```
GET    /memberships           # List plans
POST   /memberships           # Create plan
GET    /memberships/:id       # Get plan
PUT    /memberships/:id       # Update plan
DELETE /memberships/:id       # Delete plan
```

### Gateway Route

```
/api/v1/membership/*  → http://localhost:5005/*
```

---

## Attendance Service

**Port:** 5006  
**Module:** `app/attendance-service`  
**Database:** `ATTENDANCE_DATABASE_URL` (optional)

### Overview

Tracks member gym check-ins and attendance.

### Endpoints

```
GET    /attendance            # List attendance records
POST   /attendance            # Record check-in
GET    /attendance/:id        # Get attendance record
```

### Gateway Route

```
/api/v1/attendance/*  → http://localhost:5006/*
```

---

## Notification Service

**Port:** 5007  
**Module:** `app/notification-service`  
**Database:** `NOTIFICATION_DATABASE_URL` (optional)

### Overview

Handles notifications (email, SMS, push).

### Endpoints

```
GET    /notifications         # List notifications
POST   /notifications         # Send notification
GET    /notifications/:id     # Get notification status
```

### Gateway Route

```
/api/v1/notification/*  → http://localhost:5007/*
```

---

## Payment Service

**Port:** 5008  
**Module:** `app/payment-service`  
**Database:** `PAYMENT_DATABASE_URL` (optional)

### Overview

Manages payments, invoicing, and billing.

### Endpoints

```
GET    /payments              # List payments
POST   /payments              # Create payment
GET    /payments/:id          # Get payment status
PUT    /payments/:id          # Update payment
```

### Gateway Route

```
/api/v1/payment/*  → http://localhost:5008/*
```

### API Reference

#### Create Payment

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/payment/payments \
  -H "Content-Type: application/json" \
  -d '{
    "member_id": "123",
    "amount": 150000,
    "currency": "IDR",
    "type": "membership"
  }'
```

**Response (201):**
```json
{
  "id": "p_456",
  "member_id": "123",
  "amount": 150000,
  "currency": "IDR",
  "status": "pending",
  "created_at": "2026-07-03T10:30:00Z"
}
```

#### Get Payment Status

**Request:**
```bash
curl http://localhost:8080/api/v1/payment/payments/p_456
```

**Response (200):**
```json
{
  "id": "p_456",
  "member_id": "123",
  "amount": 150000,
  "currency": "IDR",
  "status": "completed",
  "paid_at": "2026-07-03T10:32:00Z"
}
```

---

## Service Interaction Example

### Scenario: Member Signs Up → Gets Charged → Receives Notification

```
1. POST /login (Auth Service)
   → Generate token

2. POST /members (Member Service, with token from auth)
   → Create member record

3. POST /payments (Payment Service)
   → Charge membership fee

4. POST /notifications (Notification Service)
   → Send welcome email

5. POST /attendance (Attendance Service)
   → Record first check-in
```

**Implementation Notes:**
- Each service can call other services via their HTTP URLs
- Pass auth token in Authorization header
- Handle errors gracefully (service timeout, 500, etc.)
- Log inter-service calls for debugging

---

## See Also

- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) — Quick commands
- [SERVICE_SETUP.md](SERVICE_SETUP.md) — Detailed setup guide
- [MICROSERVICES.md](MICROSERVICES.md) — Architecture overview
