# Requirements: Authentication Baseline

**Feature ID:** F-0001
**Status:** Accepted
**Author:** product-owner
**Date:** 2026-04-30

---

## 1. Context

The online shop API currently exposes public health endpoints but has no authentication or authorization foundation. The system needs a baseline identity layer before customer portal, staff portal, order, product, and account-specific workflows can safely exist.

This feature establishes first-party credentials authentication in the Go API for two portal contexts: a Next.js customer portal and a React staff portal. It intentionally focuses on the smallest useful foundation: customer registration and login, staff login for manually seeded staff accounts, separate portal sessions, role-based access checks, and refresh-token-backed sessions.

The result is not a complete identity platform. It is the baseline that future features can build on for Google login, email verification, password reset, staff account setup, profile data, fine-grained permissions, and MFA.

---

## 2. Scope

**In scope:**

- Customer credentials registration.
- Customer credentials login, refresh, logout, and current-user lookup.
- Staff credentials login, refresh, logout, and current-user lookup.
- Separate customer and staff portal sessions.
- JWT access tokens with stored refresh sessions.
- Refresh token rotation and simple previous-token reuse detection.
- Seeded system roles for `customer` and `staff`.
- Role assignment storage through `roles` and `user_roles`.
- Manual or command-based staff bootstrap for MVP.
- Basic per-instance rate limiting for customer registration, customer login, and staff login.
- Integration testing against real PostgreSQL using Testcontainers.

**Out of scope:**

- Google OAuth login.
- Account linking across providers.
- Email verification.
- Password reset.
- Staff account creation from the staff portal.
- Staff password setup links.
- Email delivery.
- User profile tables or profile management.
- Fine-grained permissions and `role_permissions`.
- MFA.
- Logout-all-devices.
- Redis-backed JWT blacklist or immediate access-token revocation.
- Backdoor or maintenance HTTP endpoints for staff bootstrap.
- Distributed, cross-instance, or risk-based rate limiting.

---

## 3. Constraints

- Authentication must be owned by the existing Go API.
- The API must support a Next.js customer portal and a React staff portal on separate subdomains that call the API on another subdomain.
- User email addresses must be globally unique.
- Customer and staff sessions must be separate.
- Staff authentication must be credentials-only in this MVP.
- Staff accounts must not be publicly registerable.
- Staff bootstrap must not use a network-facing backdoor or maintenance endpoint.
- Database migrations must be explicit and intentional, not automatically run on API startup.
- Access tokens must use HS256 for MVP.
- Passwords must be hashed with Argon2id.
- Auth integration tests must run under normal `go test ./...` using Testcontainers.
- Credentials passwords must be 8 to 128 characters and must not contain leading or trailing whitespace.
- Customer registration, customer login, and staff login endpoints must be rate limited.

---

## 4. Dependencies

No cross-feature dependencies.

---

## 5. User Stories

### US-001: Customer Credentials Registration

**As a** customer
**I want** to register with my email and password
**So that** I can create an account for the customer portal

#### Acceptance Criteria

- [ ] Given an unused email and a valid password, when the customer registers, then the system creates a user with the `customer` role and returns a successful registration response.
- [ ] Given an email that already belongs to any user, when the customer attempts to register, then the system rejects the request with a safe conflict response.
- [ ] Given a password shorter than 8 characters, longer than 128 characters, or containing leading or trailing whitespace, when the customer attempts to register, then the system rejects the request without creating a user.

---

### US-002: Customer Credentials Login

**As a** customer
**I want** to log in with email and password
**So that** I can access customer portal features as myself

#### Acceptance Criteria

- [ ] Given a customer account with correct credentials, when the customer logs in through the customer auth endpoint, then the system returns a customer access token and sets a customer refresh session cookie.
- [ ] Given incorrect credentials, when the customer attempts to log in, then the system returns a generic login failure response.
- [ ] Given a user without the `customer` role, when that user attempts customer login, then the system does not grant customer portal access.

---

### US-003: Staff Credentials Login

**As a** staff user
**I want** to log in with email and password
**So that** I can access the staff portal

#### Acceptance Criteria

- [ ] Given a manually bootstrapped staff account with correct credentials, when the staff user logs in through the staff auth endpoint, then the system returns a staff access token and sets a staff refresh session cookie.
- [ ] Given a customer-only account, when the user attempts staff login, then the system does not grant staff portal access.
- [ ] Given incorrect credentials, an inactive user, a missing credentials account, or a missing staff role, when a staff login is attempted, then the system returns a generic login failure response.

---

### US-004: Separate Portal Sessions

**As a** user who may interact with different portals
**I want** customer and staff sessions to remain separate
**So that** one portal session does not overwrite or leak into another portal

#### Acceptance Criteria

- [ ] Given a customer login, when the session is created, then it is stored as a customer portal session and uses the customer refresh cookie.
- [ ] Given a staff login, when the session is created, then it is stored as a staff portal session and uses the staff refresh cookie.
- [ ] Given a customer access token, when it is used against staff-only authentication context, then the system rejects it because the portal context does not match.
- [ ] Given a staff access token, when it is used against customer-only authentication context, then the system rejects it because the portal context does not match.

---

### US-005: Refresh Token Rotation

**As a** signed-in user
**I want** my session to refresh securely
**So that** I can stay signed in without storing long-lived access tokens in the portal app

#### Acceptance Criteria

- [ ] Given a valid refresh token cookie, when the user refreshes, then the system issues a new access token, rotates the refresh token, and stores only the new refresh token hash.
- [ ] Given a previously rotated refresh token, when it is reused, then the system revokes only that refresh session for MVP and requires login again.
- [ ] Given an expired or revoked refresh session, when refresh is attempted, then the system rejects the request and does not issue a new access token.

---

### US-006: Logout Current Session

**As a** signed-in user
**I want** to log out of my current portal session
**So that** the refresh token for this browser can no longer be used

#### Acceptance Criteria

- [ ] Given a valid customer refresh session, when customer logout is requested, then the system revokes that customer session and clears the customer refresh cookie.
- [ ] Given a valid staff refresh session, when staff logout is requested, then the system revokes that staff session and clears the staff refresh cookie.
- [ ] Given logout succeeds, when the old refresh token is used again, then the system rejects the refresh attempt.

---

### US-007: Current User Lookup

**As a** signed-in user
**I want** the portal to fetch my current authenticated identity
**So that** the UI can render account and role-aware state

#### Acceptance Criteria

- [ ] Given a valid customer access token, when `/me` is requested, then the response includes the current user identity, customer portal context, and assigned roles.
- [ ] Given a valid staff access token, when `/me` is requested, then the response includes the current user identity, staff portal context, and assigned roles.
- [ ] Given a missing, expired, or invalid access token, when `/me` is requested, then the system returns an unauthorized response.

---

### US-008: Staff Bootstrap

**As an** operator
**I want** a non-HTTP way to create the first staff user
**So that** the MVP can support staff login without exposing a temporary backdoor endpoint

#### Acceptance Criteria

- [ ] Given an operator intentionally runs the staff bootstrap path with a unique email and valid password, when the command completes, then the system creates a staff user, credentials account, and `staff` role assignment.
- [ ] Given the target email already exists, when staff bootstrap is attempted, then the system refuses to overwrite the existing user.
- [ ] Given production environment is targeted, when staff bootstrap is attempted, then the system requires a `--confirm-production` flag and typed target-email confirmation before creating the staff user.

---

### US-009: Auth Schema and Role Foundation

**As a** developer
**I want** a durable auth schema with seeded roles
**So that** future authorization features can build on the MVP without replacing the foundation

#### Acceptance Criteria

- [ ] Given migrations are run intentionally, when the auth schema is applied, then the database contains `users`, `accounts`, `sessions`, `roles`, and `user_roles` tables.
- [ ] Given role seed data is applied, when the system starts using authorization checks, then `customer` and `staff` roles are available as seeded system records.
- [ ] Given future permissions are out of scope, when the schema is created, then it does not include permission tables in this feature.

---

### US-010: Rate-Limited Authentication Attempts

**As a** system operator
**I want** registration and login attempts to be rate limited
**So that** the authentication baseline has basic protection against brute force and credential stuffing

#### Acceptance Criteria

- [ ] Given repeated customer login attempts from the same client and email within the configured window, when the limit is exceeded, then the system rejects additional attempts with a rate limit response.
- [ ] Given repeated staff login attempts from the same client and email within the configured window, when the limit is exceeded, then the system rejects additional attempts with a rate limit response.
- [ ] Given repeated customer registration attempts from the same client within the configured window, when the limit is exceeded, then the system rejects additional attempts with a rate limit response.

---

## 6. Non-Goals

- This feature will not implement Google OAuth login.
- This feature will not implement customer email verification or password reset.
- This feature will not implement staff account creation, staff setup links, or email delivery.
- This feature will not implement profiles, fine-grained permissions, MFA, logout-all, or Redis token blacklisting.
- This feature will not expose a backdoor or maintenance HTTP endpoint for staff bootstrap.
- This feature will not automatically run migrations during API startup.
- This feature will not implement distributed, cross-instance, or risk-based rate limiting.
- This feature will not implement role revocation workflows or role revocation audit history.

---

## 7. Open Questions

No open questions at time of writing.
