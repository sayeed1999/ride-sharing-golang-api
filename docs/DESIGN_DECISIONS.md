# Design Decisions

This document captures concise architecture decisions for important modules.

## Auth module

A compact, pluggable auth module (domain → repository → service → handler) provides register/login, password hashing (salt + bcrypt) and optional role assignment.

- Pluggable: module is self-contained and has no ride-sharing-specific dependencies; wire it via `module.go`.
- Role feature-flag: `RequireRoleOnRegistration` (config) enforces picking/assigning a role at registration when enabled; otherwise role is optional.

### Usage (when to set the flag)

- true = enforce role at signup (student/teacher, talent/client, multi-tenant).
- false = make role optional or assign later (consumer apps, invite/admin flows).

Keep the flag in `config.Config` (env toggle) so host projects can change behavior without forking.

## E2E Tests

### Facade Design Pattern

In E2E testing, the setup process contains a series of tasks:

1. Spin up a real docker container
2. Seed necessary data
3. Setup routes for test env
4. Expose routes via a http server

We used facade pattern via `test_app.go` to wrap up these tasks into a single point and used `NewTestApp()` method from tests to simplify the process.

## Trip and Trip Request Separation

We have decided to use separate tables for `trip_requests` and `trips` for the following reasons:

- **Data Integrity**: A `trip` must have a `driver_id`, but a `trip_request` does not. Having a single table would require the `driver_id` to be nullable, which could lead to data inconsistency.
- **Clearer State Management**: The status of a `trip_request` (e.g., `pending`, `searching`) is distinct from the status of a `trip` (e.g., `in_progress`, `completed`). Separating the tables allows for a clearer and more robust state machine for each entity.
- **Scalability**: This separation allows for independent scaling of the `trip_requests` and `trips` tables. For example, we might want to archive old trip requests more aggressively than completed trips.
