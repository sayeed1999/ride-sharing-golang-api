# Testing

This project contains two primary test tiers used during development and CI:

- Unit tests: fast, isolated tests that exercise individual usecases and helpers. These use in-memory mock data (the `repository/mocks` package) so tests stay fast and deterministic.
- End-to-end (E2E) tests: spin up a real PostgreSQL container and exercise real HTTP endpoints to simulate real behavior and provide confidence that the API works.

Why this structure

- Unit tests give quick feedback and validate business logic (e.g., password hashing, role rules).
- E2E tests provide confidence that migrations, GORM models, and password hashing + verification work against a real database; these run inside Docker via testcontainers so CI can run them reproducibly.

Key patterns used

- Real docker container is spin up for E2E tests on runtime and destroyed once the tests are completed.
- In-memory mock data is used in the unit tests.

Test tooling

- testcontainers-go — starts/stops real database containers from tests and returns connection details.
- testify — assertions and helpers (require/assert).

How tests are organised

- `internal/modules/.../usecase/*_test.go` — unit tests using `repository/mocks`.
- `tests/e2e/*` — full E2E using testcontainers and the real database.

Run tests locally

Ensure Docker is running for E2E tests.

Run all tests (unit + E2E):

```bash
go test ./...
```

Run the E2E suite (requires Docker):

```bash
# from repo root
go test ./tests/e2e -v
```
