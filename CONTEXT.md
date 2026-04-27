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

### Middleware Usage Importance

- Middlewares are heavily used to centralize MUST-HAVE checks from similar domain endpoints e.g.,
  - `customer_middleware` checks valid customer or not and injects `customer` in context. Used in customer endpoints
  - `driver_middleware` checks valid driver or not and injects `driver` in context. Used in driver endpoints (TODO)
  - `trip_request_middleware` checks valid trip_request or not, then (MOST IMPORTANT)
    - fetches the authenticated customer from context, & matches the trip_request belongs to this customer or not to stop fradulent attacks
    *Note: here it makes `trip_requst_middleware` completely dependent on `customer_middleware`, so must be added in call chain after it*

## Do NOT

- Do not assume a generic `internal/handlers` + `internal/services` monolith structure; this codebase is module-sliced under `internal/modules/*`.
- 
