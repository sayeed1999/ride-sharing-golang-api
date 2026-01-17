#!/bin/bash

# Install golang-migrate if not installed
if ! [ -x "$(command -v migrate)" ]; then
  echo 'Error: migrate is not installed.'
  echo 'Installing...'
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# set env before running migrations
. ./scripts/setenv.sh

# Run down migration for last one migration
migrate -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" -path internal/database/migrations down 1
