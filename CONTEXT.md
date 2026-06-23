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
    <module_name>>/                   # auth/trip/payment/etc.
      domain/
      dto/
      repository/
      service/
      handler/
      di/
      module.go
      routes.go
  pkg/                # packages that are specific to application business
    middleware/                       # domain-aware middlewares (customer/trip request loading)
    password/                         # bcrypt + salt helpers
pkg/.                 # packages that are not business specific and can be shared across projects
  jwt/                                # shared JWT service used by auth module
  middleware/                         # shared auth middleware (bearer parsing + claim extraction)
  test_helper/                        # HTTP test utilities/assertion helpers
tests/                # tests outside unit tests, e.g integraton or e2e or others
  e2e/                                # end-to-end API workflow tests
```

## Commands

- `go run ./cmd/app/.` # run locally
- `go test ./...` # run all tests
- `go build -o bin/app ./cmd/app/.` # build binary
- `docker compose -f deployments/docker/docker-compose.yml up -d` # local stack
- `docker compose -f deployments/docker/docker-compose.yml down` # stop local stack

## Conventions Followed in Code

- Auth module is completely isolated with separate db, because we may want to change whole auth impl later without breaking ride sharing business.
- Customer signup and driver signup is separated because in real ride sharing they might have separate app. And it happens like trip module reaches to auth modules for signup. But for auth, it doesn't know about customer and driver, so it depends on trip module to send necessary details and do post signup activities.
- For endpoints => Handler -> Service -> Repository layering is consistently used in every modules.
- For schedulers/background jobs => Scheduler -> Service -> Repository layering will be followed.
- Proper DI is mandatory through `di.go` in each module to make the whole business unit testable.
- Error handling is mostly early-return style to reduce nested blocks.
- Cross-module consistency during signup uses compensating actions (delete auth user when customer/driver creation fails), instead of a single distributed DB transaction.

### Unit Test + Mocking Conventions

- Use `testify` for assertions, mocks, and test suites.
- Keep test files aligned with source files (`x.go` -> `x_test.go`).
- Prefer `t.Run()` subtests for scenarios; include happy path plus critical failure/edge path(s).
- Prioritize critical business behavior over 100% coverage.
- Keep reusable mocks under `<module>/repository/mocks/` (or closest dependency-level `mocks/` package).
- Keep common setup, fixtures, and shared test constants in `test_helpers_test.go` (or equivalent helper files).
- Prefer readable, low-noise tests: reuse fixtures/constants instead of repeated hardcoded values.
- Mock matching rule:
  - use strict argument matching only when field-level values are core behavior
  - otherwise prefer `mock.Anything` to reduce brittleness
- Always assert key side effects in important paths:
  - expected calls
  - expected non-calls
  - error/status outcomes
- Unit tests should validate business logic through mocked dependencies, not DB/network integration.

### Middleware Usage Importance

- Middlewares centralize MUST-HAVE checks for similar domain endpoints (see `internal/modules/trip/routes.go` for chains):
  - `customer_middleware` — validates customer, injects `customer` into context. Used on customer trip-request routes.
  - `driver_middleware` — validates driver, injects `driver` into context. Used on driver trip-request and trip start/complete routes.
  - `trip_request_middleware` — loads trip request by ID and verifies it belongs to the authenticated customer. **Must run after `customer_middleware`.**
  - `trip_middleware` — loads trip by ID and sets `trip` in context. Optionally verifies driver/customer ownership when those actors are already in context. Used on all `/trips/:trip_id` routes.

## Trip module

- **Implementation spec (canonical):** [`docs/specs/TRIP.md`](./docs/specs/TRIP.md) — endpoints, statuses, transitions, invariants.
- Do not implement trip behavior from other docs; follow the spec.

### Status transition enforcement

Allowed transitions are defined in TRIP.md §4. Code enforces them centrally:

- **File:** `internal/modules/trip/domain/status_transitions.go`
- **Shape:** `map[Status][]Status` for both `TripRequestStatus` and `TripStatus`, mirroring §4.
- **Check:** `from.CanTransitionTo(to)` before any status update in services.
- **Rule:** Do not add ad-hoc `if status == X` checks in handlers or services; extend the map when §4 changes.

## Do NOT

- Do not assume a generic `internal/handlers` + `internal/services` monolith structure; this codebase is module-sliced under `internal/modules/*`.
