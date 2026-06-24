# ADR-003: Trip request / trip separation

**Status:** accepted

## Context

A trip_request requires only `customer_id`, but a trip requires both `customer_id` and `driver_id`. Their lifecycles differ (matching vs ride execution).

## Decision

Use separate tables: `trip_requests` and `trips`.

- Data integrity — no nullable `driver_id` on a combined table
- Clearer state machines per entity
- Independent scaling and archival

Behavior (statuses, endpoints, transitions): [`docs/specs/TRIP.md`](../specs/TRIP.md)

## Consequences

- Handoff on driver accept: freeze `trip_request`, create `trip`
- Post-accept lifecycle updates `trip` only
