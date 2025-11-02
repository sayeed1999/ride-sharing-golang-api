# Introduction

This **Golang** microservice is a part of my project [**sayeed1999/Ride-Sharing-Platform**](https://github.com/sayeed1999/Ride-Sharing-Platform). This microservice contains the core algorithms e.g ride matchmaking, trip transition checking of the ride sharing app.

The main app **Ride-Sharing-Platform**, which is built with **.NET**, uses this service as its microservice to do the core algorithm heavy lifting for it.

## Find Project Documentations here

- [Contribution Guidelines](./CONTRIBUTING.md)
- [BRD - Business Requirement Document](./docs/BRD.md)
- [PRD - Product Requirement Document](./docs/PRD.md)
- [Project Structure](./docs/PROJECT_STRUCTURE.md)
- [Design Decisions](./docs/DESIGN_DECISIONS.md)
- [Testing](./docs/TESTING.md)

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
docker build -t ride-sharing-golang-api -f deployments/docker/Dockerfile .
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
docker compose -f deployments/docker/docker-compose.yml up -d
```

To stop the running containers, run: -

```noset
docker-compose -f deployments/docker/docker-compose.yml down
```

## Troubleshooting

### Dirty Database Migration

<i>We are using **golang-migrate** for managing database migrations.</i>

In the event of a migration failure that results in a "dirty database" error, you may need to manually intervene to resolve the issue. This typically happens when a migration script fails midway, leaving the database in an inconsistent state. To fix this, you can follow these steps:

1.  Open a PostgreSQL administration tool (e.g., pgAdmin).
2.  Connect to the application database.
3.  Navigate to the `public.schema_migrations` table.
4.  You will see a row for each migration that has been applied. The `version` column corresponds to the migration file number, and the `dirty` column indicates whether the migration failed.
5.  To fix a dirty migration, you will need to manually edit the corresponding row in the `schema_migrations` table. For example, if migration version 6 is dirty, you would change the `version` from 6 to 5 and the `dirty` column from `true` to `false`.
6.  After manually fixing the `schema_migrations` table, you can run the down migration to revert the changes and then run the up migration again to apply the changes correctly.
