# Grove Base — Release Notes

## v1.3.0 — gest v2, standard _test.go files

**Released:** 2026-03-04

Migrates the testing layer from gest v1's custom runner to gest v2, which is powered by the native `go test` engine. Tests are now standard `_test.go` files — full IDE support, caching, `-race` detector and real coverage out of the box.

---

### What's new

#### gest upgraded to v2.0.3

The library import path changes from `github.com/caiolandgraf/gest` to `github.com/caiolandgraf/gest/v2`. v2 drops the custom runner entirely in favour of standard Go test infrastructure.

#### `internal/tests/` now uses `_test.go` files

`user_spec.go` (`package main`) has been replaced with `user_test.go` (`package tests`). Each suite is a standard `TestXxx(t *testing.T)` function that calls `s.Run(t)` at the end — no `init()`, no `Register()`, no `RunRegistered()`.

```go
// internal/tests/user_test.go
package tests

import (
    "testing"
    "github.com/caiolandgraf/gest/v2/gest"
    "github.com/caiolandgraf/grove-base/internal/models"
)

func TestUser(t *testing.T) {
    s := gest.Describe("User")
    s.It("should have valid fields", func(t *gest.T) { ... })
    s.Run(t)
}
```

#### `internal/tests/main.go` deleted

The gest v1 entrypoint (`main.go` calling `gest.RunRegistered()`) is no longer needed and has been removed.

#### Tests run with `gest ./...` or `go test ./...`

Install the gest CLI once and run from anywhere:

```sh
go install github.com/caiolandgraf/gest/v2/cmd/gest@latest

gest ./...           # beautiful Jest-style output
gest -w ./...        # watch mode
gest -c ./...        # coverage table
go test ./...        # plain go test also works
```

#### `grove.toml` watcher exclusion removed

`internal/tests` is no longer excluded from the `grove dev` watcher — `_test.go` files don't affect the application binary and the exclusion is unnecessary.

---

### Dependencies

- `github.com/caiolandgraf/gest` v1 removed
- `github.com/caiolandgraf/gest/v2` v2.0.3 added

---

## v1.2.0 — Grove-aligned structure, test scaffold & observability rebrand

**Released:** 2026-03-03

This release fully aligns `grove-base` with the conventions and expectations of the Grove CLI. The template is now ready to be used as-is after `grove setup`, with working test examples, correct directory naming, and a consistent identity across all tooling.

---

### What's new

#### Testing scaffold included out of the box

`internal/tests/` is now part of the template with two files pre-created:

- `main.go` — the gest entrypoint, identical to what `grove make:test` would generate automatically
- `user_spec.go` — a working example spec for the `User` model covering field validation and table naming

Running `grove test` works immediately after `grove setup`, with no extra steps required.

#### `internal/middlewares/` renamed to `internal/middleware/`

The directory and package name now match the Grove CLI docs and generator output exactly. All internal imports were updated accordingly.

#### `grove.toml` added

A `grove.toml` with a fully configured `[dev]` section is now committed to the template. The `internal/tests` directory is excluded from the hot-reload watcher so a spec save never triggers an application rebuild.

#### `.grove/` added to `.gitignore`

The temporary binary directory created by `grove dev` is now properly ignored.

---

### Fixes & cleanup

- **`DB_NAME` default** changed from a leftover internal project name (`mcs_dctfweb_sender`) to `grove_db`, consistent with `.env.example`
- **Atlas `local` env URL** updated to `grove_user:grove_password@localhost:5432/grove_db`
- **OTel service name fallback** changed from `go-project-base` to `grove-app`
- **otelhttp span name** in `routes.go` updated to `grove-app`
- **Grafana dashboard** renamed from `go-project-base.json` to `grove-app.json`; all Prometheus/Loki query expressions updated to use the `grove-app` service label
- **Prometheus & Promtail** job names updated from `go-project-base` to `grove-app`
- **`README.md`** fully rewritten with Grove-oriented quick start, project structure, command reference, workflow guide, and observability table

---

### Dependencies

- `github.com/caiolandgraf/gest` bumped from `v0.1.0` → `v1.1.0` (latest) and promoted from indirect to a direct dependency