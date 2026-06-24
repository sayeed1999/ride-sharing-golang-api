# Testing

## Tiers

| Tier | Where | What |
|------|-------|------|
| Unit | `internal/modules/.../service/*_test.go` | Business logic with `repository/mocks` |
| E2E | `tests/e2e/*` | Real PostgreSQL (testcontainers) + HTTP |

E2E setup uses facade pattern — see [ADR-002](../adr/002-e2e-test-facade.md) (`NewTestApp()` in `test_app.go`).

## Tooling

- `testcontainers-go` — database containers for e2e
- `testify` — assertions (`require` / `assert`)

## Run

Docker must be running for e2e.

```bash
go test ./...              # all tests
go test ./tests/e2e -v     # e2e only
```

## When adding code

- Service changes → unit tests with mocked repos
- User-facing flows → e2e in `tests/e2e/<feature>_test.go`
- Run `go test ./...` before finishing
