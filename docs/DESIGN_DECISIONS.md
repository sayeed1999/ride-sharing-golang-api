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
