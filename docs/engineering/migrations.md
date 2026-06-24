# Migrations

SQL migrations live in `internal/database/migrations/`. Apply with `scripts/migrate.sh`.

Uses **golang-migrate** — tracks state in `public.schema_migrations`.

## Dirty database

If a migration fails midway, you may see a "dirty database" error.

1. Open pgAdmin (or any PostgreSQL client) and connect to the app database
2. Open `public.schema_migrations`
3. Find the failed row — `dirty = true`
4. Set `version` to the previous migration (e.g. `6` → `5`) and `dirty` to `false`
5. Run the down migration, then run up again
