# Auth Service

A secure, production ready(trying to make) authentication service built with Go. Provides JWT-based authentication with access/refresh token rotation, rate limiting, and comprehensive security middleware.

## Features

### Current (v1.0)

- User registration with validation  
- Login with username or email  
- JWT access tokens (short-lived) + refresh tokens (long-lived)   
- Refresh token rotation with reuse detection   
- Token revocation on logout   
- Password change with session revocation   
- Role-based access control (user/admin)   
- Rate limiting per IP on auth endpoints   
- Database connection pooling   
- JSON-structured logging   
- Request ID tracing   
- Panic recovery middleware   
- CORS support   
- Secure headers (HSTS, CSP, X-Frame-Options, etc.)   

### Roadmap

See the [Roadmap](#roadmap) section for planned features.

---

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL 15+
- `goose` — database migration tool
  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

### 1. Clone and configure

```bash
git clone https://github.com/Dharshan2208/auth.git
cd auth
cp .env.example .env
```

Edit `.env` with your settings (see [Configuration](#configuration)).

### 2. Set up the database

```bash
createdb authdb
make migrate-up
```

### 3. Run the service

```bash
go run ./cmd/server
```

The server starts on `http://localhost:8080`.  
Swagger UI is available at `http://localhost:8080/swagger/`.

---

## Configuration

All configuration is via environment variables. Copy `.env.example` to `.env` and customize.

| Variable | Default | Required | Description |
|----------|---------|----------|-------------|
| `PORT` | `8080` | No | HTTP listen port |
| `DATABASE_URL` | — | **Yes** | PostgreSQL connection string |
| `JWT_SECRET` | — | **Yes** | HMAC signing key for JWT tokens (min 32 bytes recommended) |
| `ACCESS_TOKEN_TTL` | `15m` | No | Access token lifetime (Go duration format) |
| `REFRESH_TOKEN_TTL` | `168h` | No | Refresh token lifetime (Go duration format) |

### Example `.env`

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/authdb?sslmode=disable
JWT_SECRET=your-strong-256-bit-secret-key-here
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h
```

---

## Project Structure

```
cmd/server/main.go          — Entry point, server bootstrap
internal/
├── auth/
│   ├── jwt.go              — JWT access + refresh token generation
│   ├── token.go            — Token hashing (SHA-256)
│   ├── password.go         — bcrypt hashing and verification
│   └── policy.go           — Username and password validation rules
├── config/
│   └── config.go           — Environment variable loading
├── handlers/
│   ├── handler.go          — Handler struct and constructor
│   ├── auth.go             — Signup, Login, Logout, Refresh, ChangePassword
│   ├── pages.go            — Profile, Admin
│   ├── health.go           — Health check
│   └── docs_types.go       — Swagger documentation types
├── httpx/
│   └── json.go             — JSON encode/decode helpers
├── middleware/
│   ├── auth.go             — JWT authentication middleware
│   ├── ratelimit.go        — Per-IP sliding window rate limiter
│   ├── logging.go          — Structured request logging with request ID
│   ├── cors.go             — CORS headers
│   ├── secureheaders.go    — Security headers (HSTS, CSP, etc.)
│   ├── recovery.go         — Panic recovery
│   └── response.go         — ResponseWriter wrapper for status capture
├── models/
│   └── user.go             — User model
├── router/
│   └── routes.go           — Route registration
└── storage/
    └── postgres.go         — PostgreSQL queries (users, sessions)
migrations/                  — Goose SQL migrations
docs/                        — Generated Swagger/OpenAPI docs
```

---

## Security Architecture

### Token Flow

```
┌──────────┐      ┌──────────────┐       ┌──────────┐
│  Client  │ ──▶ │  Auth Service│ ───▶ │ Database │
└──────────┘      └──────────────┘       └──────────┘
     │                   │                   │
     │  1. POST /login   │                   │
     │ ◀──── tokens ────│                   │
     │                   │                   │
     │  2. GET /profile  │                   │
     │    (Bearer JWT)   │                   │
     │ ◀───── data ──── │                   │
     │                   │                   │
     │  3. POST /refresh │                   │
     │ ◀── new tokens ──│ ── rotate ──────▶│
     │                   │                   │
     │  4. POST /logout  │                   │
     │                   │ ── revoke ──────▶│
```

### Key Security Measures

- **Password storage**: bcrypt with default cost factor
- **Access tokens**: Short-lived (default 15 min), signed with HMAC-SHA256
- **Refresh tokens**: Stored as SHA-256 hashes in the database (not plaintext)
- **Token rotation**: Each refresh invalidates the old token and issues a new pair
- **Reuse detection**: If a rotated token is reused, all sessions for that device are revoked
- **Rate limiting**: Per-IP on all auth-modifying endpoints
- **Input validation**: Email RFC validation, username regex, password complexity
- **HTTP hardening**: HSTS, CSP, X-Frame-Options, X-Content-Type-Options, Referrer-Policy, Permissions-Policy
- **Panic recovery**: All panics are caught, logged, and return 500
- **Request timeouts**: Read (10s), ReadHeader (5s), Write (10s), Idle (60s)
- **Body size limit**: 1MB max via `MaxBytesReader`

---

## Development

### Database Migrations

```bash
make migrate-up          # Apply all pending migrations
make migrate-down        # Roll back the last migration
make migrate-status      # Show migration status
make migrate-reset       # Roll back all migrations
make migration name=<name>  # Create a new migration file
```

### Useful Commands

```bash
go run ./cmd/server                        # Run the server
go build -o bin/auth ./cmd/server          # Build binary
go vet ./...                               # Static analysis
go test ./...                              # Run tests (none yet)
```

---

## Roadmap

The following features are planned for future releases, organized by priority.

### Phase 2 — Core Features

- [ ] **Email verification** — Send verification email on signup, prevent login until email is confirmed
- [ ] **Password reset flow** — `POST /forgot-password` and `POST /reset-password` endpoints
- [ ] **Persistent Rate Limiting Storage** - Right now in memory is being used and have to move to redis

### Phase 3 — Security Hardening

- [ ] **Refresh token expiry cleanup** — Periodic job to delete expired session records from the database
- [ ] **CSRF protection** — Token-based CSRF protection for cookie-based auth flows

### Phase 4 — Developer & User Experience

- [ ] **Session management** — `GET /api/v1/sessions` to list active sessions, `DELETE /api/v1/sessions/:id` to revoke a specific session
- [ ] **Profile update** — `PUT /api/v1/profile` to update email, username, etc.
- [ ] **Prometheus metrics** — `/metrics` endpoint with request count, latency, error rate

### Phase 5 — Advanced Features

- [ ] **OAuth2 / OIDC providers** — Login with Google, GitHub, etc.
- [ ] **Account locking/unlocking** — Admin endpoints to lock/unlock user accounts
- [ ] **Dockerisation** — Dockerfile + docker-compose.yml with PostgreSQL
- [ ] **Testing** — Unit tests, integration tests, attack/stress tests

---

## License

MIT
