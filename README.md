# Grove Base

> Official starter template for [Grove](https://github.com/caiolandgraf/grove) — a production-ready Go foundation wiring together GORM, fuego, Atlas, gest and OpenTelemetry.

## Stack

| Tool | Role |
|------|------|
| [fuego](https://github.com/go-fuego/fuego) | HTTP router + automatic OpenAPI 3.1 |
| [GORM](https://gorm.io) | ORM & generic repository layer |
| [Atlas](https://atlasgo.io) | Schema migration engine |
| [gest](https://github.com/caiolandgraf/gest) | Jest-inspired testing framework for Go |
| [SCS](https://github.com/alexedwards/scs) | Session management (Redis-backed) |
| [OpenTelemetry](https://opentelemetry.io) | Distributed tracing + Prometheus metrics |

---

## Quick Start

This project is scaffolded automatically by the Grove CLI:

```sh
grove setup my-api
cd my-api
cp .env.example .env
grove dev
```

Your API will be running at `http://localhost:8080`.
The OpenAPI / Swagger UI (Scalar) is available at `http://localhost:8080/swagger`.

---

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── app/                     # Shared singletons (DB, Redis, Session, Metrics)
│   ├── config/                  # Infrastructure initializers (DB, Redis, OTel, etc.)
│   ├── controllers/             # fuego route handlers
│   ├── database/                # Generic GORM repository
│   ├── dto/                     # Request and response types
│   ├── middleware/              # HTTP middlewares (CORS, session, observability)
│   ├── models/                  # GORM models
│   ├── routes/                  # Route registration
│   └── tests/                   # gest spec files
│       ├── main.go              # gest entrypoint (auto-created by grove make:test)
│       └── user_spec.go         # Example spec
├── migrations/                  # Atlas SQL migrations
├── infra/                       # Observability stack config (Prometheus, Grafana, Loki, Jaeger)
├── .env.example                 # Committed environment template
├── atlas.hcl                    # Atlas configuration
├── docker-compose.yml           # Full observability stack
└── grove.toml                   # Grove dev server configuration
```

---

## Development

### Prerequisites

- Go 1.22+
- [Atlas CLI](https://atlasgo.io/docs/getting-started/installation)
- Docker (for the database and observability stack)

### Running locally

```sh
# 1. Start infrastructure (Postgres, Redis, Jaeger, Grafana, Prometheus, Loki)
docker compose up -d db redis

# 2. Copy and fill in environment variables
cp .env.example .env

# 3. Run migrations
grove migrate

# 4. Start the dev server with built-in hot reload
grove dev
```

---

## Commands

### Generators

```sh
grove make:model <Name>          # Scaffold a GORM model
grove make:model <Name> -r       # Full resource (model + migration + controller + DTO)
grove make:controller <Name>     # Scaffold a fuego controller
grove make:dto <Name>            # Scaffold request/response DTOs
grove make:middleware <Name>     # Scaffold an HTTP middleware
grove make:migration <name>      # Generate a SQL migration via Atlas diff
grove make:resource <Name>       # Alias for make:model -r
```

### Testing

```sh
grove make:test <Name>   # Scaffold a new gest spec in internal/tests/
grove test               # Run all specs
grove test -c            # Run specs + per-suite coverage report
grove test -w            # Watch mode — re-run on every save
grove test -wc           # Watch mode + coverage
```

### Server & Build

```sh
grove dev          # Hot reload (built-in, no external tools required)
grove dev:air      # Hot reload via Air
grove build        # Compile binary to ./bin/app
```

### Database

```sh
grove migrate              # Apply pending migrations
grove migrate:rollback     # Rollback last migration
grove migrate:status       # Show migration status
grove migrate:fresh        # Drop all tables and re-apply ⚠️
grove migrate:hash         # Rehash atlas.sum
```

---

## Typical Workflow

```sh
# 1. Scaffold a full resource
grove make:resource Post

# 2. Edit the model and DTO, then apply the migration
grove migrate

# 3. Register the routes in internal/routes/routes.go

# 4. Write tests
grove make:test Post

# 5. Run the suite
grove test -c
```

---

## Testing with gest

Spec files live in `internal/tests/` and self-register via `init()`:

```go
// internal/tests/post_spec.go
package main

import "github.com/caiolandgraf/gest/gest"

func init() {
    s := gest.Describe("Post")

    s.It("should have a valid title", func(t *gest.T) {
        t.Expect("Hello").Not().ToBe("")
    })

    gest.Register(s)
}
```

> **Note:** gest uses `_spec.go` instead of `_test.go` because the Go toolchain reserves `_test.go` for `go test`. gest runs via `go run`, so any other suffix works.

---

## Observability

Bring up the full stack:

```sh
docker compose up -d
```

| Service | URL |
|---------|-----|
| API | http://localhost:8080 |
| OpenAPI (Scalar) | http://localhost:8080/swagger |
| Prometheus metrics | http://localhost:8080/metrics |
| Prometheus UI | http://localhost:9090 |
| Grafana | http://localhost:3000 (admin/admin) |
| Jaeger UI | http://localhost:16686 |
| Loki | http://localhost:3100 |

---

## Environment Variables

Copy `.env.example` to `.env` and adjust:

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `grove_user` | Database user |
| `DB_PASSWORD` | `grove_password` | Database password |
| `DB_NAME` | `grove_db` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode |
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_PASSWORD` | _(empty)_ | Redis password |
| `OTEL_SERVICE_NAME` | `grove-app` | Service name in traces |
| `OTEL_EXPLOERER_OTLP_ENDPOINT` | `localhost:4318` | OTLP HTTP endpoint |
| `LOG_LEVEL` | `info` | Log level (`debug`, `info`, `warn`, `error`) |

---

## License

MIT © Caio Landgraf