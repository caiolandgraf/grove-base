# Grove Base

> Official starter template for [Grove](https://github.com/caiolandgraf/grove) — a production-ready Go REST API foundation with modular architecture, full observability, Atlas migrations, and [gest](https://github.com/caiolandgraf/gest) testing.

## Tech Stack

| Category             | Technology                                                                                                                                               |
| -------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **CLI / Dev**        | [Grove](https://github.com/caiolandgraf/grove) (hot reload, generators, migrations, tests)                                                               |
| **Language**         | [Go 1.25+](https://go.dev/)                                                                                                                              |
| **HTTP Framework**   | [Fuego](https://github.com/go-fuego/fuego)                                                                                                               |
| **Database**         | [PostgreSQL](https://www.postgresql.org/) via [GORM](https://gorm.io/)                                                                                   |
| **Cache / Sessions** | [Redis](https://redis.io/) via [Redigo](https://github.com/gomodule/redigo) + [SCS](https://github.com/alexedwards/scs)                                  |
| **Migrations**       | [Atlas](https://atlasgo.io/) (schema-as-code from GORM models)                                                                                           |
| **Testing**          | [gest](https://github.com/caiolandgraf/gest) (Jest-style) + integration tests with PostgreSQL                                                            |
| **Tracing**          | [OpenTelemetry](https://opentelemetry.io/) + [Jaeger](https://www.jaegertracing.io/)                                                                     |
| **Metrics**          | [Prometheus](https://prometheus.io/) via OTel Prometheus Exporter                                                                                        |
| **Logs**             | [`log/slog`](https://pkg.go.dev/log/slog) → [Loki](https://grafana.com/oss/loki/) via [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/) |
| **Dashboards**       | [Grafana](https://grafana.com/) (pre-configured with all data sources)                                                                                   |
| **API Docs**         | OpenAPI 3 auto-generated + [Scalar UI](https://github.com/scalar/scalar)                                                                                 |
| **Containers**       | [Docker Compose](https://docs.docker.com/compose/)                                                                                                       |

---

## Observability Architecture

```
                          ┌──────────────────────────────────────────────────┐
                          │                   Grafana :3000                  │
                          │  ┌──────────┐  ┌──────────┐  ┌──────────────┐    │
                          │  │Dashboards│  │  Explore  │  │    Alerts    │   │
                          │  └────┬─────┘  └────┬─────┘  └──────┬───────┘    │
                          └───────┼─────────────┼───────────────┼────────────┘
                                  │             │               │
                 ┌────────────────┼─────────────┼───────────────┼────────────┐
                 │                ▼             ▼               ▼            │
                 │  ┌──────────────────┐ ┌───────────┐ ┌──────────────┐      │
                 │  │ Prometheus :9090 │ │ Loki :3100│ │ Jaeger :16686│      │
                 │  └───────┬──────────┘ └─────┬─────┘ └──────┬───────┘      │
                 │          │ scrape            │ push          │ OTLP       │
                 │          │                   │               │            │
                 │          │            ┌──────┴──────┐        │            │
                 │          │            │  Promtail   │        │            │
                 │          │            └──────┬──────┘        │            │
                 │          │                   │ tail          │            │
                 └──────────┼───────────────────┼──────────────┼─────────────┘
                            │                   │              │
                   ┌────────┼───────────────────┼──────────────┼────────┐
                   │        ▼                   ▼              ▲        │
                   │   GET /metrics     .grove/logs/app.log   OTLP HTTP    │
                   │        │                   │              │        │
                   │        └───────────────────┴──────────────┘        │
                   │                   Go App :8080                     │
                   │           slog (JSON) + OTel + Prometheus          │
                   └────────────────────────────────────────────────────┘
```

| Signal      | Flow                                                           |
| ----------- | -------------------------------------------------------------- |
| **Traces**  | App → OTLP HTTP → Jaeger → Grafana                             |
| **Metrics** | Prometheus scrapes `GET /metrics` → Grafana                    |
| **Logs**    | App → `slog` JSON → `.grove/logs/app.log` → Promtail → Loki → Grafana |

---

## Project Structure

```
grove-base/
├── cmd/
│   ├── api/
│   │   └── main.go                # Application entrypoint
│   ├── atlas/
│   │   └── main.go                # Atlas GORM schema loader
│   ├── seed/
│   │   └── main.go                # Database seeder CLI
│   └── scalar/
│       └── scalar.go              # Scalar API docs UI handler
├── internal/
│   ├── app/                       # Shared infrastructure (config, DB, helpers, middleware, types)
│   │   ├── config/
│   │   │   ├── env.go             # Typed config (caarlos0/env) + .env discovery
│   │   │   ├── database.go        # PostgreSQL + GORM setup
│   │   │   ├── logger.go          # slog JSON logger (stdout + file)
│   │   │   ├── metrics.go         # Prometheus metrics via OTel exporter
│   │   │   ├── otel.go            # OpenTelemetry tracing setup
│   │   │   ├── redis.go           # Redis connection pool
│   │   │   └── session.go         # SCS session manager
│   │   ├── database/
│   │   │   ├── repository.go      # Generic Repository[T] (Eloquent-like base)
│   │   │   ├── registry.go        # Atlas migration registry (auto via model init())
│   │   │   ├── migrations/        # Atlas SQL migrations
│   │   │   └── seeders/           # Database seeders
│   │   ├── helpers/
│   │   │   ├── jsonutils/         # JSON utilities
│   │   │   └── validator/         # Request validation
│   │   ├── middleware/
│   │   │   ├── cors.go            # CORS middleware
│   │   │   └── session.go         # Session middleware
│   │   ├── types/
│   │   │   ├── common-dto.go      # Shared HTTP types (errors, messages)
│   │   │   └── health-dto.go      # Health check types
│   │   ├── router/
│   │   │   ├── doc.go             # Declarative OpenAPI route documentation
│   │   │   ├── options.go         # Doc → Fuego options
│   │   │   └── register.go        # router.Get/Post/Put/Delete wrappers
│   │   └── app.go                 # Package doc
│   ├── modules/
│   │   ├── auth/                  # Auth domain (dto, service, controller, docs)
│   │   ├── users/                 # Users domain (model+repo, dto, service, controller, docs)
│   │   ├── module.go              # Module interface + Boot
│   │   └── register.go            # Module registry
│   └── routes/
│       ├── health.go              # Health check routes
│       └── routes.go              # Global routes + module mounting
│   └── tests/                     # gest test suites (*_test.go)
├── infra/
│   ├── grafana/
│   │   ├── dashboards/
│   │   │   └── grove-base.json  # Pre-built dashboard
│   │   └── provisioning/
│   │       ├── dashboards/
│   │       │   └── dashboards.yml    # Dashboard provisioning
│   │       └── datasources/
│   │           └── datasources.yml   # Prometheus, Loki, Jaeger
│   ├── loki-config.yml            # Loki storage config
│   ├── prometheus.yml             # Prometheus scrape config
│   └── promtail-config.yml        # Promtail log collection
├── .github/
│   └── workflows/
│       ├── ci.yml                 # Lint, test, build, migration checks
│       ├── cd.yml                 # Docker image build & push to GHCR
│       └── release.yml            # Tag-based binaries + Docker release
├── atlas.hcl                      # Atlas migration config
├── docker-compose.yml             # Full infrastructure stack
├── Dockerfile                     # Multi-stage production image
├── grove.toml                     # Grove dev server configuration
└── go.mod                         # Go module definition
```

### Architecture

The project follows a **modular MSC architecture** — each domain is self-contained:

```
Request → Routes → Middlewares → Module (Controller → Service → Repository) → Database
```

| Layer            | Responsibility                                              |
| ---------------- | ----------------------------------------------------------- |
| **Routes**       | Global middleware, health, mounts all modules               |
| **Modules**      | Per-domain package: model+repo, dto, service, controller, docs |
| **Router**       | Declarative OpenAPI documentation per endpoint              |
| **App/database** | Generic `Repository[T]` base (Eloquent-like) embedded per module |
| **App/types**    | Shared HTTP response types (errors, health, messages)       |

Each module wires itself (`New` + `Wire`) and registers routes via `Mount`. Add new domains in `internal/modules/register.go`. Models auto-register for Atlas via `init()` in each module's `model.go`.

---

## Prerequisites

- **Go** 1.25+
- **[Grove CLI](https://github.com/caiolandgraf/grove)** — `go install github.com/caiolandgraf/grove@latest`
- **[gest CLI](https://github.com/caiolandgraf/gest)** _(optional, prettier test output)_ — `go install github.com/caiolandgraf/gest/v2/cmd/gest@latest`
- **Docker** & **Docker Compose** (PostgreSQL, Redis, observability stack)
- **Atlas CLI** — [install guide](https://atlasgo.io/getting-started#installation)

---

## Getting Started

Scaffold a new project with Grove:

```
grove setup my-api
cd my-api
cp .env.example .env
docker compose up -d
grove migrate
grove dev
```

Or clone this repository directly:

```
git clone https://github.com/caiolandgraf/grove-base.git
cd grove-base
```

### 2. Configure environment variables

Copy the example file and adjust values as needed:

```
cp .env.example .env
```

Required environment variables:

| Variable                       | Default                 | Description                       |
| ------------------------------ | ----------------------- | --------------------------------- |
| `DB_HOST`                      | `localhost`             | PostgreSQL host                   |
| `DB_PORT`                      | `5432`                  | PostgreSQL port                   |
| `DB_USER`                      | `grove_user`            | PostgreSQL user                   |
| `DB_PASSWORD`                  | `grove_password`        | PostgreSQL password               |
| `DB_NAME`                      | `grove_db`       | PostgreSQL database name          |
| `DB_SSLMODE`                   | `disable`               | PostgreSQL SSL mode               |
| `REDIS_HOST`                   | `localhost`             | Redis host                        |
| `REDIS_PORT`                   | `6379`                  | Redis port                        |
| `REDIS_PASSWORD`               |                         | Redis password (optional)         |
| `OTEL_SERVICE_NAME`            | `grove-base`       | OpenTelemetry service name        |
| `OTEL_ENABLED`                 | `true`                  | Enable OpenTelemetry tracing      |
| `OTEL_EXPLOERER_OTLP_ENDPOINT` | `localhost:4318`        | OTLP HTTP collector endpoint      |
| `OTEL_TRACE_SAMPLE_RATIO`      | `1.0`                   | Trace sampling ratio (0.0–1.0)    |
| `METRICS_ENABLED`              | `true`                  | Enable Prometheus metrics         |
| `CORS_ALLOWED_ORIGINS`         | `http://localhost`      | Comma-separated allowed origins   |
| `BASE_URL`                     | `http://localhost:8080` | Base URL for Scalar API docs      |
| `SERVER_ADDR`                  | `:8080`                 | HTTP listen address               |
| `ENVIRONMENT`                  | `development`           | App environment (affects cookies) |
| `APP_NAME`                     | `Grove Base`            | Application name (Scalar UI)      |
| `LOG_LEVEL`                    | `info`                  | Log level (debug/info/warn/error) |
| `LOG_FILE`                     | `.grove/logs/app.log`   | Log file for Promtail (empty = stdout only) |

### 3. Start infrastructure

```
docker compose up -d
```

This starts **PostgreSQL**, **Redis**, **Jaeger**, **Prometheus**, **Loki**, **Promtail**, and **Grafana**.

### 4. Run migrations

```
grove migrate
```

### 5. Seed the database (optional)

```
grove db:seed
```

### 6. Start the server

```
grove dev
```

The server starts at **http://localhost:8080** with built-in hot reload (configured in `grove.toml`).

---

## API Endpoints

### Documentation

Once running, visit the interactive API docs powered by Scalar:

- **Scalar UI**: [http://localhost:8080/swagger](http://localhost:8080/swagger)
- **OpenAPI JSON**: [http://localhost:8080/swagger/openapi.json](http://localhost:8080/swagger/openapi.json) (generated at runtime by Fuego)

### Available Routes

| Method   | Path                      | Description             | Auth   |
| -------- | ------------------------- | ----------------------- | ------ |
| `GET`    | `/`                       | Health check (simple)   | Public |
| `GET`    | `/health`                 | Health check (detailed) | Public |
| `GET`    | `/metrics`                | Prometheus metrics      | Public |
| `POST`   | `/api/v1/auth/register`   | Register a new user     | Public |
| `POST`   | `/api/v1/auth/login`      | Login                   | Public |
| `POST`   | `/api/v1/auth/logout`     | Logout                  | Public |
| `GET`    | `/api/v1/auth/me`         | Current user info       | Public |
| `GET`    | `/api/v1/users`           | List all users          | Public |
| `POST`   | `/api/v1/users`           | Create user             | Public |
| `GET`    | `/api/v1/users/{user_id}` | Get user by ID          | Public |
| `PUT`    | `/api/v1/users/{user_id}` | Update user             | Public |
| `DELETE` | `/api/v1/users/{user_id}` | Delete user             | Public |

---

## Grove Commands

### Generators

```
grove make:resource <Name>       # Full module scaffold (model + migration + controller + DTO)
grove make:model <Name>          # Scaffold a GORM model in internal/modules/
grove make:controller <Name>     # Scaffold a fuego controller
grove make:dto <Name>            # Scaffold request/response DTOs
grove make:middleware <Name>     # Scaffold HTTP middleware
grove make:migration <name>      # Generate SQL migration via Atlas diff
grove make:test <Name>           # Scaffold gest test in internal/tests/
```

### Testing

```
grove test              # Run all tests (gest CLI if installed, else go test)
grove test -c           # Tests + per-suite coverage
grove test -w           # Watch mode
grove test -wc          # Watch + coverage
```

### Server & Build

```
grove dev               # Hot reload (built-in, no Air required)
grove dev:air           # Hot reload via Air
grove build                        # Compile binary (default: ./bin/app)
grove build -o .grove/bin/app      # Recommended: keep builds under .grove/
```

### Database

```
grove migrate              # Apply pending migrations
grove migrate:rollback     # Rollback last migration
grove migrate:status       # Show migration status
grove migrate:fresh        # Drop all tables and re-apply (dev only)
grove migrate:hash         # Rehash atlas.sum
grove db:seed              # Run database seeders
```

---

## Testing with gest

Test files live in `internal/tests/` as standard `_test.go` files. Each suite uses gest's `Describe` / `It` / `Expect` API:

```go
func TestUser(t *testing.T) {
    s := gest.Describe("User")
    s.It("should have valid fields", func(t *gest.T) {
        t.Expect(user.Name).ToBe("John Doe")
    })
    s.Run(t)
}
```

**Unit tests** use in-memory mocks (`users_service_test.go`).

**Integration tests** (`integration_test.go`) hit a real PostgreSQL database — run `grove migrate` first, or use `go test -short ./...` to skip them.

```
grove test -c           # recommended
go test ./...           # plain output
go test -short ./...    # skip integration tests
```

---

## Migrations with Atlas

This project uses [Atlas](https://atlasgo.io/) in **Program Mode** — models self-register via `init()` and `cmd/atlas` loads all modules automatically.

### Create a new migration

1. Add or modify a model in `internal/modules/<domain>/model.go` (with `database.Register(&YourModel{})` in `init()`)
2. Register the module in `internal/modules/register.go` (HTTP mount — model registration happens via import)
3. Generate the migration:

```
grove make:migration describe_your_change
```

4. Review the generated SQL in `internal/app/database/migrations/`
5. Apply it:

```
grove migrate
```

### Check migration status

```
grove migrate:status
```

---

## Database Seeders

Seeders live in `internal/app/database/seeders/`. Run with:

```
grove db:seed
```

Default seeded users:

| Email                    | Password   |
| ------------------------ | ---------- |
| `admin@grove.local` | `admin123` |
| `user@grove.local`  | `user1234` |

---

## Observability

### Dashboards & UIs

| Tool           | URL                                              | Credentials       |
| -------------- | ------------------------------------------------ | ----------------- |
| **Grafana**    | [http://localhost:3000](http://localhost:3000)   | `admin` / `admin` |
| **Jaeger**     | [http://localhost:16686](http://localhost:16686) | —                 |
| **Prometheus** | [http://localhost:9090](http://localhost:9090)   | —                 |

### Logging (`slog` → Loki)

Logs are written to **stdout** and, by default, **`.grove/logs/app.log`** (created at runtime, gitignored). Promtail tails that file and ships to Loki.

#### Configuration

Set the `LOG_LEVEL` environment variable:

| Level   | Description                                    |
| ------- | ---------------------------------------------- |
| `debug` | Verbose output, includes source file locations |
| `info`  | General operational messages (default)         |
| `warn`  | Warning conditions                             |
| `error` | Error conditions only                          |

#### Example output

```
{"time":"2025-07-17T10:00:00Z","level":"INFO","msg":"Server starting","addr":":8080"}
{"time":"2025-07-17T10:00:00Z","level":"INFO","msg":"Database connected successfully","host":"localhost","port":"5432","database":"mcs_dctfweb_sender"}
{"time":"2025-07-17T10:00:00Z","level":"INFO","msg":"Redis connected successfully","host":"localhost","port":"6379"}
{"time":"2025-07-17T10:00:00Z","level":"INFO","msg":"OpenTelemetry initialized","service":"grove-base","endpoint":"localhost:4318"}
```

#### GORM integration

GORM queries are routed through a custom `slog` adapter:

- Normal queries → `DEBUG` level
- Slow queries (>200ms) → `WARN` level
- Query errors → `ERROR` level

All query logs include `component=gorm`, elapsed time, affected rows, and the SQL statement.

#### Querying logs in Grafana

Open Grafana → Explore → select **Loki** data source:

```
{job="grove-base"} | json
```

Filter by level:

```
{job="grove-base"} | json | level = "ERROR"
```

### Metrics (Prometheus)

The app exposes a `GET /metrics` endpoint powered by the **OTel Prometheus exporter**. This automatically collects metrics from the `otelhttp` middleware and Go runtime.

Available metrics include:

| Metric                                 | Type      | Description                   |
| -------------------------------------- | --------- | ----------------------------- |
| `http_server_request_duration_seconds` | Histogram | HTTP request latency by route |
| `http_server_active_requests`          | Gauge     | Currently in-flight requests  |
| `go_goroutines`                        | Gauge     | Number of active goroutines   |
| `go_memstats_alloc_bytes`              | Gauge     | Allocated heap memory         |
| `go_gc_duration_seconds`               | Summary   | GC pause durations            |

### Tracing (Jaeger)

Full **OpenTelemetry** distributed tracing:

- **HTTP requests** are traced via `otelhttp` middleware
- **GORM queries** are automatically traced via `otelgorm`
- **Trace context propagation** via W3C TraceContext + Baggage headers

All traces are exported via OTLP HTTP to Jaeger. The Jaeger UI supports **dark mode** (toggle in the top-right corner).

#### Correlation: Logs ↔ Traces

In Grafana, the Loki data source is configured with **derived fields** that extract `traceID` from JSON logs and link them directly to Jaeger. Click on a trace ID in any log line to jump to the full trace view.

### Pre-built Dashboard

Grafana comes with a **pre-provisioned dashboard** (`Grove Base`) that includes:

| Panel               | Data Source | Description                               |
| ------------------- | ----------- | ----------------------------------------- |
| HTTP Request Rate   | Prometheus  | Requests per second by route and method   |
| HTTP Latency (p95)  | Prometheus  | 95th percentile latency by route          |
| HTTP Error Rate     | Prometheus  | 4xx/5xx errors per second                 |
| Active Requests     | Prometheus  | Currently in-flight requests              |
| Total Requests (5m) | Prometheus  | Total requests in the last 5 minutes      |
| Error Rate %        | Prometheus  | 5xx errors as a percentage of all traffic |
| Avg Latency         | Prometheus  | Average response time                     |
| Memory Usage        | Prometheus  | Alloc, Sys, Heap In-Use                   |
| Goroutines          | Prometheus  | Active goroutine count over time          |
| GC Cycles           | Prometheus  | Garbage collection frequency              |
| Error Logs          | Loki        | Live stream of ERROR-level log entries    |
| Log Volume by Level | Loki        | Stacked chart of logs by level over time  |
| All Logs            | Loki        | Full log stream with JSON parsing         |

---

## Docker Compose Services

| Service      | Image                             | Ports                              | Purpose                    |
| ------------ | --------------------------------- | ---------------------------------- | -------------------------- |
| `db`         | `postgres:latest`                 | `5432`                             | Primary database           |
| `redis`      | `redis:7-alpine`                  | `6379`                             | Session store & cache      |
| `jaeger`     | `jaegertracing/all-in-one:latest` | `16686` (UI), `4317`/`4318` (OTLP) | Distributed tracing        |
| `loki`       | `grafana/loki:3.5.0`              | `3100`                             | Log aggregation            |
| `promtail`   | `grafana/promtail:3.5.0`          | —                                  | Log collection agent       |
| `prometheus` | `prom/prometheus:v3.4.1`          | `9090`                             | Metrics storage & querying |
| `grafana`    | `grafana/grafana:11.6.0`          | `3000`                             | Dashboards & visualization |

---

## Project Conventions

- **Modules**: One package per domain under `internal/modules/` (model+repo, dto, service, controller, docs)
- **Shared infra**: Config, database, helpers, middleware, router, and types live under `internal/app/`
- **Wiring**: `New` (testable, accepts interfaces) + `Wire` (production, accepts `*gorm.DB`) per layer
- **Routes**: Document each endpoint in `docs.go` using `router.Doc`; register module in `modules/register.go`
- **Migrations**: Model self-registers in `init()`; SQL files live in `internal/app/database/migrations/`; Atlas loads all modules via `cmd/atlas`
- **Error handling**: Errors are wrapped with `fmt.Errorf("context: %w", err)` and propagated up
- **Configuration**: Typed `config.Env` struct via `caarlos0/env`; call `config.Load()` at startup
- **Logging**: Always use `slog` with structured key-value pairs — never `fmt.Println` or `log.Println`
- **Observability**: All infrastructure configs live in `infra/`; Grafana is pre-provisioned on `docker compose up`

---

## CI/CD

The project ships with **GitHub Actions** pipelines and a production **Dockerfile** out of the box.

### Continuous Integration (`ci.yml`)

Runs on every push and pull request to `main`:

| Job          | What it does                                              |
| ------------ | --------------------------------------------------------- |
| **Lint**     | `golangci-lint` with project config (`.golangci.yml`)     |
| **Test**     | `go test -race` with PostgreSQL and Redis service containers |
| **Build**    | Compiles `cmd/api` and `cmd/atlas`, uploads API artifact  |
| **Migrations** | `atlas migrate validate` + `atlas migrate lint`         |

Run the same checks locally:

```
golangci-lint run ./...
grove test
grove build
atlas migrate validate --env local
```

Install golangci-lint: [golangci-lint install](https://golangci-lint.run/welcome/install/)

### Continuous Delivery (`cd.yml`)

On push to `main` or version tags (`v*`), builds a multi-stage Docker image and publishes it to **GitHub Container Registry**:

```
ghcr.io/<owner>/grove-base:latest
ghcr.io/<owner>/grove-base:<sha>
ghcr.io/<owner>/grove-base:v1.0.0   # on tag
```

Build and run locally:

```
docker build -t grove-base .
docker run --rm -p 8080:8080 --env-file .env grove-base
```

> The app reads configuration from environment variables. `.env` is optional (used for local development only).

### Release (`release.yml`)

On version tags (`v*`), builds cross-platform binaries (linux/darwin, amd64/arm64) and publishes a GitHub Release, plus a semver-tagged Docker image to GHCR.

### Dependabot

`.github/dependabot.yml` keeps Go modules and GitHub Actions up to date with weekly PRs.

---

## License

This project is provided as a boilerplate/template. Use it freely for your own projects.
