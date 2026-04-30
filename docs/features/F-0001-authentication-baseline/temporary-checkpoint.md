# Temporary Checkpoint: F-0001 Authentication Baseline

This is a temporary handoff note for resuming guided implementation. It is not a formal feature artifact.

## Resume Prompt

```text
$instructor resume F-0001 authentication baseline at Step 1I: add valid JWT secret fixture
```

## Project Context

Repository: `online-shop`

Feature package:

- `docs/features/F-0001-authentication-baseline/README.md`
- `docs/features/F-0001-authentication-baseline/specs/requirements.md`
- `docs/features/F-0001-authentication-baseline/specs/technical-spec.md`

Relevant ADRs:

- `docs/adr/0001-authentication-session-architecture.md`
- `docs/adr/0002-database-migrations-and-staff-bootstrap.md`

Implementation is currently in the config groundwork portion of F-0001. The broader feature will eventually add auth migrations, `internal/auth`, auth routes, JWT access tokens, refresh sessions, staff bootstrap CLI, and integration tests. Do not jump there yet; continue one small config step at a time.

## Session Style

The user wants `$instructor` pacing:

- Review the user's changes before giving the next step.
- Give one small portion at a time.
- Do not silently implement code unless explicitly asked.
- Pause at good atomic commit boundaries.

## Current Code State

We are in Step 1: config groundwork.

Completed:

- Added `Config.AppEnv`, `Config.Auth`, and `Config.Database`.
- Added `AuthConfig` and `DatabaseConfig`.
- Generalized env key binding through `configKeys()`.
- Added default auth and database pool settings.
- Added required `APP_ENV`, `AUTH_CUSTOMER_PORTAL_ORIGIN`, and `AUTH_STAFF_PORTAL_ORIGIN` validation.
- Added uppercase `.env` alias support.
- Added readable `.env` test helper in `internal/config/config_test.go`.
- Started JWT secret decoding in `internal/config/config.go`.

Current WIP files:

- `internal/config/config.go`
- `internal/config/config_test.go`
- `docs/features/F-0001-authentication-baseline/temporary-checkpoint.md`

Expected current test state:

```bash
go test ./internal/config -count=1
```

fails because test fixtures do not yet include `AUTH_JWT_SECRET`.

Observed failure before this checkpoint:

```text
config: AUTH_JWT_SECRET must decoded to at least 32 bytes
```

That failure is expected at this point. The wording has a typo and should be cleaned up after the fixture is added.

## Current Design Decisions

- Config validation should stay explicit with guard clauses, not a generic validation package.
- Request validation later can use Gin binding / `go-playground/validator`.
- `DATABASE_URL` remains top-level; pool knobs live in `DatabaseConfig`.
- Real `.env` files and `.env.example` should use uppercase env-style keys such as `AUTH_CUSTOMER_PORTAL_ORIGIN`.
- Viper aliases are used so uppercase `.env` keys map to nested app config keys.
- JWT secret is externally configured as base64url text without padding, then decoded into `cfg.Auth.JWTSecret []byte`.
- JWT secret validation must check decoded byte length, not encoded string length.

## Next Step

Step 1I: add a valid JWT secret fixture.

Use this value:

```text
AUTH_JWT_SECRET=MDEyMzQ1Njc4OWFiY2RlZjAxMjM0NTY3ODlhYmNkZWY
```

It decodes to the 32-byte string:

```text
0123456789abcdef0123456789abcdef
```

Add the value to:

- `setRequiredAuthEnv(t)` in `internal/config/config_test.go`
- `requiredDotenv()` in `internal/config/config_test.go`

Then run:

```bash
go test ./internal/config -count=1
```

## Follow-Up After Step 1I

After tests are green, clean up these JWT secret error messages in `internal/config/config.go` if they still have the typo/current wording:

```go
"config: AUTH_JWT_SECRET must be base64url without padding: %w"
"config: AUTH_JWT_SECRET must decode to at least 32 bytes"
```

Then add focused tests for:

- invalid base64url JWT secret
- JWT secret that decodes to fewer than 32 bytes

After those pass, review whether config groundwork is a good atomic commit boundary.

Suggested WIP commit message:

```text
[WIP] Add authentication config groundwork
```
