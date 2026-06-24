# Add an HTTP endpoint

Checklist for a new endpoint in `internal/modules/<module>/`:

1. **Migration** — `internal/database/migrations/000NNN_name.up.sql` + `.down.sql`; run `scripts/migrate.sh`
2. **Domain** — structs in `domain/` (no HTTP/DB coupling)
3. **Repository** — interface + PostgreSQL implementation in `repository/`
4. **Service** — business logic in `service/`; unit test with mocked repo
5. **HTTP** — DTOs, handler, route in `handler/http/`
6. **Wire up** — DI in `di.go`; register routes from `module.go` via handler package
7. **Test** — unit tests for service; e2e if user-facing flow

If the feature has defined behavior, write or update a spec in `docs/specs/` first.
