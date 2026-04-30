# ADR-0002: Database Migrations and Staff Bootstrap

**Status:** Accepted
**Date:** 2026-04-30
**Last Updated:** 2026-04-30
**Author:** adr-author
**Reviewers:** Not required for this project workflow

---

## Table of Contents

1. Context
2. Scope
3. Constraints
4. Options Considered
5. Decision
6. Migration Strategy
7. Staff Bootstrap Strategy
8. Operational Safeguards
9. Consequences
10. Open Questions
11. Decision Log

---

## 1. Context

The authentication baseline introduces persistent identity and session tables. These tables must be created through a deliberate migration process that works locally, in CI, and in production without being hidden inside API startup behavior.

The MVP also needs a way to create at least one staff user before staff self-management exists. Public staff registration is not allowed, staff setup links are deferred, and backdoor HTTP endpoints are explicitly out of scope. A bootstrap path is still required so staff credentials login can be tested and used.

This ADR documents the migration tool and the temporary staff bootstrap approach so future developers understand why the project uses explicit migrations and a short-lived operational command instead of automatic startup migrations, raw SQL scripts, or a maintenance endpoint.

---

## 2. Scope

**In scope:**

- Migration tooling for the Go/PostgreSQL API.
- Whether migrations run automatically on API startup.
- Staff bootstrap mechanism for MVP.
- Production safeguards for staff bootstrap.
- Lifecycle expectation for the temporary bootstrap command.

**Out of scope (addressed elsewhere or deferred):**

- Full staff account management.
- Staff setup links.
- Email delivery.
- Password reset.
- Fine-grained permissions.
- CI/CD pipeline implementation details.
- Infrastructure-as-code.

---

## 3. Constraints

- Database migrations must be explicit and intentional.
- API startup must not automatically mutate the production schema.
- Staff accounts must not be publicly registerable.
- Staff bootstrap must not use a network-facing backdoor or maintenance endpoint.
- Staff bootstrap must not commit a real password or precomputed password hash in a migration.
- Staff bootstrap must be removable after real staff creation/setup exists.
- The application uses Go and PostgreSQL.

---

## 4. Options Considered

### Option A — Goose SQL Migrations ✅ Selected

Use Goose with SQL-first migrations stored in `migrations/`.

| Aspect             | Assessment                                                 |
| ------------------ | ---------------------------------------------------------- |
| Go fit             | Strong, with CLI and library support                       |
| Reviewability      | SQL migrations are explicit and easy to review             |
| Local workflow     | Simple `goose -dir migrations postgres "$DATABASE_URL" up` |
| Future flexibility | Can embed migrations or use Go migrations later if needed  |
| Familiarity        | Less familiar to the user than `migrate`, but acceptable   |

Selected because it fits a Go/PostgreSQL project and supports intentional SQL-first migrations with minimal ceremony.

### Option B — golang-migrate/migrate ❌ Rejected

Use the widely known `migrate` tool with paired up/down SQL files.

| Aspect       | Assessment                                                          |
| ------------ | ------------------------------------------------------------------- |
| Familiarity  | High across many teams                                              |
| Simplicity   | Strong                                                              |
| Go fit       | Strong                                                              |
| Review style | Paired files can be slightly noisier than one-file Goose migrations |

Rejected because Goose was selected for this project after comparing both tools. `migrate` remains a valid alternative if the team later standardizes on it.

### Option C — Automatic Startup Migrations ❌ Rejected

Run migrations from the API process during startup.

| Aspect              | Assessment                                                    |
| ------------------- | ------------------------------------------------------------- |
| Convenience         | High for simple local development                             |
| Production safety   | Risky when multiple instances start or deploy concurrently    |
| Operational clarity | Weak because schema mutation is hidden inside server boot     |
| Failure mode        | Startup health becomes coupled to database migration behavior |

Rejected because production schema changes should be deliberate operational actions.

### Option D — Versioned Temporary Staff Bootstrap CLI ✅ Selected

Create a temporary command under `cmd/create-staff-user` to create a staff user with shared application logic.

| Aspect        | Assessment                                                                |
| ------------- | ------------------------------------------------------------------------- |
| Reviewability | Versioned code can be reviewed                                            |
| Security      | No network-facing bootstrap surface                                       |
| Consistency   | Can reuse password hashing, validation, repository, and transaction logic |
| Lifecycle     | Must be removed after staff self-management ships                         |

Selected because it creates a safe MVP path without adding a backdoor endpoint or ad hoc local script.

### Option E — Raw Local Script Outside Version Control ❌ Rejected

Use an untracked local script or manual SQL to create staff users.

| Aspect        | Assessment                                                       |
| ------------- | ---------------------------------------------------------------- |
| Initial speed | Fast                                                             |
| Repeatability | Poor                                                             |
| Reviewability | None                                                             |
| Security      | Risk of inconsistent hashing, role assignment, or leaked secrets |

Rejected because staff bootstrap touches real credentials and authorization state.

### Option F — HTTP Maintenance or Backdoor Endpoint ❌ Rejected

Expose a temporary internal endpoint to create a staff user.

| Aspect                 | Assessment                                                   |
| ---------------------- | ------------------------------------------------------------ |
| Convenience            | Easy to call remotely                                        |
| Attack surface         | High-value network-facing staff creation path                |
| Lifecycle risk         | Temporary endpoints are easy to forget                       |
| Protection requirement | Requires another authentication mechanism before auth exists |

Rejected because the CLI bootstrap covers the need without exposing a network path.

---

## 5. Decision

Use Goose SQL migrations and run migrations intentionally outside API startup.

Use a temporary versioned `cmd/create-staff-user` command for MVP staff bootstrap, with production safeguards, no HTTP maintenance endpoint, shared service/repository logic, and a requirement to remove the command after real staff self-management exists.

---

## 6. Migration Strategy

Migrations live in:

```text
migrations/
```

The intended command shape is:

```bash
goose -dir migrations postgres "$DATABASE_URL" up
goose -dir migrations postgres "$DATABASE_URL" down
goose -dir migrations postgres "$DATABASE_URL" status
```

Local development should expose Make targets for common migration actions. API startup must not invoke Goose. Testcontainers-based integration tests may run migrations inside test setup against disposable PostgreSQL containers.

Auth migrations may enable shared PostgreSQL extensions such as `citext`, but the auth down migration must not drop those extensions. Later migrations may add unrelated columns that depend on the same extension, so the auth down path should only remove auth-owned tables, indexes, constraints, and seed data in reverse dependency order.

---

## 7. Staff Bootstrap Strategy

The MVP staff bootstrap command lives in:

```text
cmd/create-staff-user
```

The command handles only operator-facing orchestration:

- parse flags
- read `DATABASE_URL`
- read `APP_ENV`
- display target database host and database name
- prompt for password
- perform production confirmation
- call shared application service logic

Shared service/repository logic handles:

- email normalization
- password policy validation
- Argon2id password hashing
- user creation
- credentials account creation
- `staff` role lookup and assignment
- transaction boundaries
- duplicate email handling

The command must refuse to overwrite an existing email. It must never log or echo the password.

---

## 8. Operational Safeguards

When `APP_ENV=production`, the command must require:

- `--confirm-production`
- printed target database host and database name
- interactive confirmation by typing the target email
- interactive password entry rather than a command-line password flag

The command must not expose a corresponding HTTP endpoint. The following endpoint classes are explicitly disallowed:

```text
/maintenance/create-staff
/internal/bootstrap
/debug/create-staff
hidden auth bypass routes
```

The command is temporary. The feature that introduces real staff creation and setup must remove this command or explicitly replace it with audited break-glass tooling.

---

## 9. Consequences

### Positive

- Migrations are explicit, reviewable, and not hidden in API startup.
- Goose gives the Go project a straightforward SQL migration workflow.
- Staff bootstrap is repeatable and versioned without exposing an HTTP attack surface.
- Password hashing and role assignment can reuse application logic.
- Production bootstrap is possible without encouraging manual SQL.

### Negative / Tradeoffs Accepted

- Operators must remember to run migrations intentionally.
- Goose may be less familiar to contributors who have primarily used `migrate`.
- The temporary command must be maintained until staff self-management exists.
- Production bootstrap still requires direct database connectivity and careful operator handling.
- Removing the command later must be tracked; otherwise temporary tooling may linger.

---

## 10. Open Questions

No open questions at time of writing.

---

## 11. Decision Log

| Date       | Decision                                         | Rationale                                                                    |
| ---------- | ------------------------------------------------ | ---------------------------------------------------------------------------- |
| 2026-04-30 | Use Goose SQL migrations                         | Fits Go/PostgreSQL and supports explicit SQL-first migration review.         |
| 2026-04-30 | Do not run migrations on API startup             | Avoids production startup races and hidden schema mutation.                  |
| 2026-04-30 | Use `cmd/create-staff-user` for bootstrap        | Provides a reviewable non-HTTP path for first staff user creation.           |
| 2026-04-30 | Reject HTTP maintenance/backdoor endpoints       | Avoids exposing staff creation as a network attack surface.                  |
| 2026-04-30 | Make bootstrap production-capable but guarded    | Avoids ad hoc production SQL while requiring explicit operator confirmation. |
| 2026-04-30 | Share service/repository logic from the command  | Keeps hashing, validation, transactions, and role assignment consistent.     |
| 2026-04-30 | Remove command after staff self-management ships | Prevents temporary bootstrap tooling from becoming a second admin system.    |
