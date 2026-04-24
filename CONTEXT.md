# Ride Sharing Golang API

## Overview

[One paragraph — what it does, who uses it, what it's NOT]

## Tech Stack

- **Language:** Go 1.22
- **Framework:** [Gin / Echo / Fiber / net/http]
- **Database:** [Postgres / MySQL] via [sqlx / GORM / pgx]
- **Auth:** [JWT / sessions]
- **Config:** [env vars / viper / godotenv]
- **Hosting:** [Railway / Fly.io / GCP]

## Project Structure

```
cmd/
  main.go           # entrypoint
internal/
  handlers/         # HTTP handlers
  services/         # business logic
  repository/       # DB layer
  models/           # structs
pkg/                # shared utilities
config/             # config loading
```

## Commands

- go run ./cmd/main.go     # run locally
- go test ./...            # run all tests
- go build -o bin/app ./cmd/main.go
- make migrate-up          # run migrations
- make migrate-down

## Conventions

- Errors: always wrap with fmt.Errorf("context: %w", err)
- No logic in handlers — delegate to service layer
- All DB access through repository interfaces only

## Current Focus

[What you're building right now]

## Do NOT

[Your guardrails]
