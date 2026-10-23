# Agent Guide

Go 1.24 + Gin API with PostgreSQL (GORM, golang-migrate) and Redis. Code is organized by module under `internal/modules/<module>`, not by global handler/service layers.

Docs map: [`docs/AGENTS.md`](docs/AGENTS.md).

## Source of truth

- Trip behavior: [`docs/specs/TRIP.md`](docs/specs/TRIP.md). Driver location: [`docs/specs/DRIVER_LOCATION.md`](docs/specs/DRIVER_LOCATION.md). Update the spec first when behavior changes; do not invent trip rules from scratch.
- Architecture, checklists, migrations: [`docs/engineering/`](docs/engineering/).
- Decisions and rationale: [`docs/adr/`](docs/adr/).

## Commands

```bash
. ./scripts/setenv.sh                 # exports local Postgres/Redis/JWT settings
docker compose up -d db redis         # dependencies only; add api for full stack
go run ./cmd/app/.                    # API on localhost:8080; runs pending migrations at boot
go build ./cmd/app/.

go test ./internal/...                # unit tests only, no Docker
go test ./tests/e2e -v                # e2e; Docker must be running
go test ./tests/e2e -run TestName -v  # one e2e test
go test ./internal/modules/trip/service -run TestName -v
go test ./...                         # everything; e2e starts containers
```

There is no Makefile, CI workflow, or lint config. Use `gofmt`, `go vet ./...`, and the tests above.

## Configuration and setup

- Env names use double underscores: `Server__Host`, `Server__Port`, `Redis__URL`. Database vars are `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`; auth signing uses `JWT_SECRET`.
- `config.LoadConfig()` loads `.env` from the current working directory and then reads OS env; defaults are in [`config/config.go`](config/config.go). `.env` is gitignored.
- `docker compose` reads variable values from the shell, so export them before `docker compose up`; there is no `env_file` entry.
- `scripts/migrate.sh` and `scripts/migrate-down.sh` source `scripts/setenv.sh`, which overwrites existing `POSTGRES_*` values with local defaults.
- Migrations run automatically on app start. Files live in `internal/database/migrations/` as `000NNN_name.up.sql` + `.down.sql`; both are expected. Tables use separate `auth` and `trip` Postgres schemas, and GORM models map to them via `TableName()` (for example `trip.trip_requests`).

## Architecture

- Entrypoint [`cmd/app/main.go`](cmd/app/main.go): load config → connect/migrate Postgres → connect Redis → register module routes on Gin.
- Auth module owns `/register`, `/login`, and auth users. Trip module owns `/customers/signup`, `/drivers/signup`, trip-request, trip, and location routes; it calls `auth.NewUserService` for cross-module signup and uses compensating deletes on failure.
- Routes are registered through each module's `ExposeRoutes`. The app and e2e tests mount everything at `/`; the `/api/trip` base path in specs is a deployment concern, not a prefix added in code.
- Trip wiring is in [`internal/modules/trip/routes.go`](internal/modules/trip/routes.go) and [`di/di_container.go`](internal/modules/trip/di/di_container.go). Trip handlers live in flat `trip/handler/`; auth uses `auth/handler/http/`. Follow the existing module's layout.
- Auth middleware sets `x-user-email` from the JWT `sub` claim. Customer/driver middleware sets `customer`/`driver`; trip middleware sets `trip`. Handlers read these keys with `c.MustGet(...)`.
- Trip status changes must go through `CanTransitionTo` in [`domain/status_transitions.go`](internal/modules/trip/domain/status_transitions.go). On driver accept, freeze the trip request and create the trip; all later lifecycle updates touch the trip only. There is one trip per trip request.

## Testing

- Unit tests live beside services (`internal/modules/trip/service/*_test.go`) and use hand-written testify mocks in `repository/mocks/`.
- Service constructors accept a `*gorm.DB`, but unit tests use unexported `new...Service` constructors with the `noopTxRunner` in [`test_helpers_test.go`](internal/modules/trip/service/test_helpers_test.go). Reuse its fixtures and mock setup instead of standing up a database.
- E2E tests use `setup.NewTestApp(ctx, t, requireRole)`; it starts `postgres:18-alpine` and `redis:7-alpine` containers per test, runs migrations, and must be paired with `defer testApp.CleanUp(ctx, t)`. Docker is required; e2e is the slow part of `go test ./...`.
