# ADR-002: E2E test facade

**Status:** accepted

## Context

E2E setup requires multiple steps: Docker container, seed data, test routes, HTTP server.

## Decision

Wrap setup in a facade via `test_app.go` — tests call `NewTestApp()` for a single entry point.

## Consequences

- Tests stay simple; setup changes live in one place
- See `docs/engineering/testing.md` for how to run e2e tests
