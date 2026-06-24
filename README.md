# Introduction

This **Golang** microservice is a part of my project [**sayeed1999/Ride-Sharing-Platform**](https://github.com/sayeed1999/Ride-Sharing-Platform). This microservice contains the core algorithms e.g ride matchmaking, trip transition checking of the ride sharing app.

The main app **Ride-Sharing-Platform**, which is built with **.NET**, uses this service as its microservice to do the core algorithm heavy lifting for it.

## Documentation

- [Agent guide](./docs/AGENTS.md) — single source of truth for agents and contributors
- [`AGENTS.md`](./AGENTS.md) — pointer to `docs/AGENTS.md` for coding agents
- [`CLAUDE.md`](./CLAUDE.md) — pointer to `docs/AGENTS.md` for Claude Code
- [Trip spec](./docs/specs/TRIP.md) — canonical trip behavior

## Project Deployment Guide

## Set environment variables

To set the env's properly, run from bash terminal: -

```noset
export Server__Host=0.0.0.0
export Server__Port=8080

export POSTGRES_USER=user
export POSTGRES_PASSWORD=password
export POSTGRES_HOST=0.0.0.0
export POSTGRES_PORT=5432
export POSTGRES_DB=ride_sharing_db

export PGADMIN_PASSWORD=admin
export PGADMIN_EMAIL=admin@local.com

export REQUIRE_ROLE_ON_REGISTRATION=true
```

Alternatively, you can run the following from bash:

```bash
chmod +x scripts/setenv.sh
. ./scripts/setenv.sh
```

What this does - first line changes the file’s permissions to allow it to be run directly like a program. second line runs the script in your current shell. thus keeps exported environment variables active.

[**Note:** While not running with **docker compose**, omit the first part **RideProcessingService__**.
While running docker compose, docker will omit the prefix **RideProcessingService__** for you.]

## Run the service without Docker

To run the project directly from terminal: -
Open a terminal from this directory and run `go run ./cmd/app/.`

The api will be running on `localhost:8080`.

## Build Docker Image

To manually build the Docker image, run from terminal: -

```bash
docker build -t ride-sharing-golang-api -f Dockerfile .
```

## Launch Container using Dockerfile

To manually run a container for this image, run for terminal: -

```bash
docker run --rm -it -p 8080:8080 ride-sharing-golang-api
```

The api will be running on `localhost:8080`.

## Launch Container using Docker Compose

To run through Docker Compose file, run from terminal: -

```noset
docker compose up -d
```

To stop the running containers, run: -

```noset
docker-compose down
```

Troubleshooting: see [`docs/engineering/migrations.md`](./docs/engineering/migrations.md) for dirty database migration fixes.
