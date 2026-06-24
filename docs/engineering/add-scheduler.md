# Add a scheduler job

Checklist for a background job in `internal/modules/<module>/`:

1. **Spec** — if behavior is defined (e.g. trip request expiry in [`TRIP.md`](../specs/TRIP.md)), follow it first
2. **Service** — business logic in `service/` (status transitions, batch queries); unit test with mocked repo
3. **Repository** — query/update methods the job needs in `repository/`
4. **Scheduler** — tick loop or cron in `scheduler/`; calls service only — no HTTP, no direct DB access
5. **Wire up** — construct scheduler in `di.go`; start from `module.go` (or `cmd/app/main.go` if app-wide)
6. **Test** — unit tests on service; optional integration test for scheduler tick with mocked service

Layering: **Scheduler → Service → Repository**
