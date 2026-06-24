# Conventions

## Layering

- HTTP: Handler → Service → Repository
- Background jobs: Scheduler → Service → Repository
- DI through each module's `di.go`
- Early-return error handling
- Module-sliced under `internal/modules/*` — not a monolithic `internal/handlers` layout

## Cross-module signup

- Auth module is isolated (separate DB); trip module orchestrates customer/driver signup
- Compensating actions on failure (e.g. delete auth user if customer creation fails) — no distributed transaction

## Middleware (trip module)

See `internal/modules/trip/routes.go` for chains:

- `customer_middleware` — validates customer, injects into context
- `driver_middleware` — validates driver, injects into context
- `trip_request_middleware` — loads trip request, verifies ownership (**after** `customer_middleware`)
- `trip_middleware` — loads trip, optional ownership check on `/trips/:trip_id`

## Trip status transitions

- Spec: [`docs/specs/TRIP.md`](../specs/TRIP.md) §4
- Code: `internal/modules/trip/domain/status_transitions.go`
- Use `from.CanTransitionTo(to)` before updates — extend the map when spec changes, no ad-hoc status checks

## Unit tests

- `testify` for assertions and mocks; `x.go` → `x_test.go`
- `t.Run()` subtests; happy path + critical failures
- Mocks in `<module>/repository/mocks/`
- Shared fixtures in `test_helpers_test.go`
- Strict mock matching only when field values are core behavior; else `mock.Anything`
- Assert calls, non-calls, and error/status outcomes

## Do NOT

- Implement trip behavior outside [`docs/specs/TRIP.md`](../specs/TRIP.md)
- Assume monolith handler/service layout
