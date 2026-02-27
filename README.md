# Go Project Base

A production-ready Go REST API boilerplate built with modern tooling and best practices.

## Tech Stack

| Category             | Technology                                                                                                                                               |
| -------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Language**         | [Go 1.25+](https://go.dev/)                                                                                                                              |
| **HTTP Framework**   | [Fuego](https://github.com/go-fuego/fuego)                                                                                                               |
| **Database**         | [PostgreSQL](https://www.postgresql.org/) via [GORM](https://gorm.io/)                                                                                   |
| **Cache / Sessions** | [Redis](https://redis.io/) via [Redigo](https://github.com/gomodule/redigo) + [SCS](https://github.com/alexedwards/scs)                                  |
| **Migrations**       | [Atlas](https://atlasgo.io/) (schema-as-code from GORM models)                                                                                           |
| **Tracing**          | [OpenTelemetry](https://opentelemetry.io/) + [Jaeger](https://www.jaegertracing.io/)                                                                     |
| **Metrics**          | [Prometheus](https://prometheus.io/) via OTel Prometheus Exporter                                                                                        |
| **Logs**             | [`log/slog`](https://pkg.go.dev/log/slog) → [Loki](https://grafana.com/oss/loki/) via [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/) |
| **Dashboards**       | [Grafana](https://grafana.com/) (pre-configured with all data sources)                                                                                   |
| **Logging**          | [`log/slog`](https://pkg.go.dev/log/slog) (structured JSON logs)                                                                                         |
| **API Docs**         | OpenAPI 3 auto-generated + [Scalar UI](https://github.com/scalar/scalar)                                                                                 |
| **Hot Reload**       | [Air](https://github.com/air-verse/air)                                                                                                                  |
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
                   │   GET /metrics        logs/app.log    OTLP HTTP    │
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
| **Logs**    | App → `slog` JSON → `logs/app.log` → Promtail → Loki → Grafana |

---

## Project Structure

```
go-project-base/
├── cmd/
│   ├── api/
│   │   └── main.go                # Application entrypoint
│   └── scalar/
│       └── scalar.go              # Scalar API docs UI handler
├── internal/
│   ├── config/
│   │   ├── database.go            # PostgreSQL + GORM setup
│   │   ├── logger.go              # slog JSON logger (stdout + file)
│   │   ├── metrics.go             # Prometheus metrics via OTel exporter
│   │   ├── otel.go                # OpenTelemetry tracing setup
│   │   ├── redis.go               # Redis connection pool
│   │   └── session.go             # SCS session manager
│   ├── container/
│   │   └── container.go           # Dependency injection container
│   ├── controllers/
│   │   ├── auth-controller.go     # Auth endpoints (login, register, logout)
│   │   └── users-controller.go    # User CRUD endpoints
│   ├── dto/
│   │   ├── auth-dto.go            # Auth request/response types
│   │   ├── common-dto.go          # Shared types
│   │   ├── health-dto.go          # Health check types
│   │   └── user-dto.go            # User request/response types
│   ├── helpers/
│   │   ├── jsonutils/             # JSON utilities
│   │   └── validator/             # Request validation
│   ├── middlewares/
│   │   ├── cors.go                # CORS middleware
│   │   └── session.go             # Session middleware
│   ├── models/
│   │   └── user.go                # GORM User model
│   ├── repositories/
│   │   └── user-repository.go     # User data access layer
│   ├── routes/
│   │   ├── health.go              # Health check routes
│   │   └── routes.go              # Route registration
│   └── services/
│       ├── auth-service.go        # Auth business logic
│       └── user-service.go        # User business logic
├── infra/
│   ├── grafana/
│   │   ├── dashboards/
│   │   │   └── go-project-base.json  # Pre-built dashboard
│   │   └── provisioning/
│   │       ├── dashboards/
│   │       │   └── dashboards.yml    # Dashboard provisioning
│   │       └── datasources/
│   │           └── datasources.yml   # Prometheus, Loki, Jaeger
│   ├── loki-config.yml            # Loki storage config
│   ├── prometheus.yml             # Prometheus scrape config
│   └── promtail-config.yml        # Promtail log collection
├── logs/                          # App log files (tailed by Promtail)
├── migrations/                    # Atlas SQL migrations
├── doc/
│   └── openapi.json               # Generated OpenAPI spec
├── .air.toml                      # Air hot reload config
├── atlas.hcl                      # Atlas migration config
├── docker-compose.yml             # Full infrastructure stack
├── jaeger-ui.json                 # Jaeger UI config (dark mode)
├── Makefile                       # Dev commands
└── go.mod                         # Go module definition
```

### Architecture

The project follows a **layered architecture** with clear separation of concerns:

```
Request → Routes → Middlewares → Controllers → Services → Repositories → Database
```

| Layer            | Responsibility                                              |
| ---------------- | ----------------------------------------------------------- |
| **Routes**       | Maps HTTP endpoints to controller handlers                  |
| **Middlewares**  | Cross-cutting concerns (CORS, sessions, auth)               |
| **Controllers**  | Parses requests, validates input, returns responses         |
| **Services**     | Business logic and orchestration                            |
| **Repositories** | Data access abstraction over GORM                           |
| **Models**       | GORM entity definitions (also used by Atlas for migrations) |
| **Container**    | Wires all dependencies together (poor man's DI)             |

---

## Prerequisites

- **Go** 1.25+
- **Docker** & **Docker Compose** (for PostgreSQL, Redis, Jaeger, Grafana, Prometheus, Loki)
- **Atlas CLI** — [install guide](https://atlasgo.io/getting-started#installation)
- **Air** _(optional, for hot reload)_ — `go install github.com/air-verse/air@latest`

---

## Getting Started

### 1. Clone the repository

```
git clone https://github.com/caiolandgraf/go-project-base.git
cd go-project-base
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
| `DB_USER`                      | `postgres`              | PostgreSQL user                   |
| `DB_PASSWORD`                  | `postgres`              | PostgreSQL password               |
| `DB_NAME`                      | `mcs_dctfweb_sender`    | PostgreSQL database name          |
| `DB_SSLMODE`                   | `disable`               | PostgreSQL SSL mode               |
| `REDIS_HOST`                   | `localhost`             | Redis host                        |
| `REDIS_PORT`                   | `6379`                  | Redis port                        |
| `REDIS_PASSWORD`               |                         | Redis password (optional)         |
| `OTEL_SERVICE_NAME`            | `go-project-base`       | OpenTelemetry service name        |
| `OTEL_EXPLOERER_OTLP_ENDPOINT` | `localhost:4318`        | OTLP HTTP collector endpoint      |
| `BASE_URL`                     | `http://localhost:8080` | Base URL for Scalar API docs      |
| `APP_NAME`                     | `Go Project Base`       | Application name (Scalar UI)      |
| `LOG_LEVEL`                    | `info`                  | Log level (debug/info/warn/error) |

### 3. Start infrastructure

```
docker compose up -d
```

This starts **PostgreSQL**, **Redis**, **Jaeger**, **Prometheus**, **Loki**, **Promtail**, and **Grafana**.

### 4. Install dependencies

```
make install
```

### 5. Run database migrations

```
make migrate-up
```

### 6. Start the server

```
# Standard
make run

# With hot reload (recommended for development)
make dev
```

The server starts at **http://localhost:8080**.

---

## API Endpoints

### Documentation

Once running, visit the interactive API docs powered by Scalar:

- **Scalar UI**: [http://localhost:8080/swagger](http://localhost:8080/swagger)
- **OpenAPI JSON**: [http://localhost:8080/swagger/openapi.json](http://localhost:8080/swagger/openapi.json)

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

## Makefile Commands

```
make help             # Show all available commands
make install          # Download and tidy Go modules
make run              # Run the application
make dev              # Run with Air hot reload
make test             # Run all tests
make migrate-create   # Create a new migration (usage: make migrate-create name=add_books_table)
make migrate-up       # Apply pending migrations
make migrate-status   # Check migration status
make migrate-hash     # Rehash migration directory
make db-reset         # Drop and recreate database, then migrate
```

---

## Migrations with Atlas

This project uses [Atlas](https://atlasgo.io/) with a **GORM provider** — your GORM models in `internal/models/` are the single source of truth for the database schema.

### Create a new migration

1. Add or modify a model in `internal/models/`
2. Generate the migration:

```
make migrate-create name=describe_your_change
```

3. Review the generated SQL in `migrations/`
4. Apply it:

```
make migrate-up
```

### Check migration status

```
make migrate-status
```

---

## Observability

### Dashboards & UIs

| Tool           | URL                                              | Credentials       |
| -------------- | ------------------------------------------------ | ----------------- |
| **Grafana**    | [http://localhost:3000](http://localhost:3000)   | `admin` / `admin` |
| **Jaeger**     | [http://localhost:16686](http://localhost:16686) | —                 |
| **Prometheus** | [http://localhost:9090](http://localhost:9090)   | —                 |

### Logging (`slog` → Loki)

The project uses Go's standard library **`log/slog`** with a JSON handler for structured logging. Logs are written to both **stdout** and **`logs/app.log`** (collected by Promtail and shipped to Loki).

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
{"time":"2025-07-17T10:00:00Z","level":"INFO","msg":"OpenTelemetry initialized","service":"go-project-base","endpoint":"localhost:4318"}
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
{job="go-project-base"} | json
```

Filter by level:

```
{job="go-project-base"} | json | level = "ERROR"
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

Grafana comes with a **pre-provisioned dashboard** (`Go Project Base`) that includes:

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

- **File naming**: kebab-case (`user-controller.go`, `auth-dto.go`)
- **Package naming**: lowercase, single word (`controllers`, `services`)
- **Error handling**: Errors are wrapped with `fmt.Errorf("context: %w", err)` and propagated up
- **Configuration**: Environment variables with sensible defaults via `getEnv(key, default)`
- **Models**: GORM models are the single source of truth; Atlas generates migrations from them
- **Logging**: Always use `slog` with structured key-value pairs — never `fmt.Println` or `log.Println`
- **Observability**: All infrastructure configs live in `infra/`; Grafana is pre-provisioned on `docker compose up`

---

## License

This project is provided as a boilerplate/template. Use it freely for your own projects.
