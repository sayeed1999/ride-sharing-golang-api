# ADR-001: Pluggable auth module

**Status:** accepted

## Context

Auth must be reusable across projects without ride-sharing coupling.

## Decision

- Self-contained module: domain → repository → service → handler, wired via `module.go`
- Password hashing: salt + bcrypt
- `RequireRoleOnRegistration` config flag — when true, role required at signup; when false, optional or assigned later

## Consequences

- Host projects toggle role behavior via env without forking
- Trip module handles customer/driver signup; auth stays unaware of those types
