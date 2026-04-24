# Ride Sharing Golang API

## Overview

This Go service powers auth and ride-flow logic for a larger ride-sharing platform.

## Tech Stack

- **Language:** Go 1.24
- **HTTP:** Gin
- **Database:** PostgreSQL via GORM, with `golang-migrate` SQL migrations
- **Auth:** JWT (`github.com/golang-jwt/jwt/v5`) with HMAC signing
- **Config:** environment variables via `godotenv` + `config.LoadConfig()`
- **Testing:** `testify` + `testcontainers-go` for e2e/integration-style tests

## Project Structure

```text
cmd/
  app/main.go                         # app entrypoint, wires config/db/modules
config/
  config.go                           # env-based config loading
internal/
  database/
    database.go                       # DB init/close + migration execution helpers
    migrations/*.sql                  # schema changes (auth + trip schemas)
  modules/
    auth/                             # auth domain: users, roles, signup/login
      domain/
      repository/
      service/
      handler/http/
      module.go
    trip/                             # ride domain: customer/driver/trip request APIs
      domain/
      dto/
      repository/
      service/
      handler/
      di/
      module.go
      routes.go
  pkg/
    middleware/                       # domain-aware middlewares (customer/trip request loading)
    password/                         # bcrypt + salt helpers
pkg/
  jwt/                                # shared JWT service used by auth module
  middleware/                         # shared auth middleware (bearer parsing + claim extraction)
  test_helper/                        # HTTP test utilities/assertion helpers
tests/
  e2e/                                # end-to-end API workflow tests
```

## Commands

- `go run ./cmd/app/.` # run locally
- `go test ./...` # run all tests
- `go build -o bin/app ./cmd/app/.` # build binary
- `docker compose -f deployments/docker/docker-compose.yml up -d` # local stack
- `docker compose -f deployments/docker/docker-compose.yml down` # stop local stack

## Conventions Observed in Code

- Handler -> Service -> Repository layering is consistently used in `auth` and `trip` modules.
- Constructors follow `NewX(...)`; route/module wiring happens in `module.go`, `routes.go`, and `di` containers.
- HTTP responses commonly use stable JSON keys (`error`, `message`, `token`, `customer`, `driver`, `trip_request`).
- Error handling is mostly early-return style; handlers frequently return raw `err.Error()` for server-side failures, while some auth paths normalize messages (e.g., invalid credentials).
- `pkg/middleware.AuthMiddleware` sets `x-user-email` in Gin context from JWT `sub`, and downstream trip handlers/middleware rely on that contract.
- Cross-module consistency during signup uses compensating actions (delete auth user when customer/driver creation fails), instead of a single distributed DB transaction.

## Consistent Domain Terms

- **Actors:** `user`, `role`, `customer`, `driver`
- **Trip request lifecycle:** `trip_request`, `origin`, `destination`, `status`
- **Vehicle:** `vehicle_type`, `vehicle_registration` with values like `bike`, `cng`, `car`
- **Cross-module identity link:** `auth_user_id`
- **Transition checker vocabulary:** `CanTransition` and status-transition validation in `trip-processor`

## Gaps vs CLAUDE.md

- No `CLAUDE.md` file exists in this repository right now, so there is no written CLAUDE guidance to compare against implementation.
- Because of that, the practical source of truth is the module wiring and tests (`cmd/app/main.go`, `internal/modules/*`, `tests/e2e/*`).

## Current Focus

The active implementation emphasis is auth + trip-request workflows, including endpoint-level transition checks and full e2e ride flow validation.

## Do NOT

- Do not assume a generic `internal/handlers` + `internal/services` monolith structure; this codebase is module-sliced under `internal/modules/*`.
- Do not introduce abstractions that break existing JWT claim/context conventions (`sub` -> `x-user-email`) without updating middleware and handlers together.
- Do not treat `trip-processor` as DB-backed business logic; it is currently an isolated transition-validation module.
