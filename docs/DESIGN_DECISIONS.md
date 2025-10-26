# Design Decisions

This document captures concise architecture decisions for important modules.

## Auth module
A compact, pluggable auth module (domain → repository → usecase → handler) provides register/login, password hashing (salt + bcrypt) and optional role assignment.

- Pluggable: module is self-contained and has no ride-sharing-specific dependencies; wire it via `module.go`.
- Role feature-flag: `RequireRoleOnRegistration` (config) enforces picking/assigning a role at registration when enabled; otherwise role is optional.

### Usage (when to set the flag)

- true = enforce role at signup (student/teacher, talent/client, multi-tenant).
- false = make role optional or assign later (consumer apps, invite/admin flows).

Keep the flag in `config.Config` (env toggle) so host projects can change behavior without forking.


## Facade Design Pattern in E2E Tests

In E2E testing, the setup process contains a series of tasks:
1. Spin up a real docker container
2. Seed necessary data
3. Setup routes for test env
4. Expose routes via a http server

We used facade pattern via `test_app.go` to wrap up these tasks into a single point and used `NewTestApp()` method from tests to simplify the process.
