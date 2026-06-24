# Project structure

Go service for auth and ride-flow logic.

## Tech stack

- Go 1.24, Gin, PostgreSQL + GORM, `golang-migrate`
- JWT (`github.com/golang-jwt/jwt/v5`), config via `godotenv` + `config.LoadConfig()`
- Tests: `testify`, `testcontainers-go`

## Layout

```text
cmd/app/main.go                       # entrypoint
config/config.go                      # env config
internal/
  database/                           # DB init, migrations/*.sql
  modules/<module>/                   # auth, trip, etc.
    domain/ repository/ service/ handler/ di/ module.go routes.go
  pkg/                                # app-specific (middleware, password)
pkg/                                  # shared (jwt, middleware, test_helper)
tests/e2e/                            # end-to-end API tests
```

## Commands

```bash
go run ./cmd/app/.
go test ./...
go build -o bin/app ./cmd/app/.
docker compose up -d
docker compose down
scripts/migrate.sh                    # apply migrations
```
