# F-0001: Authentication Baseline

**Status:** In Progress
**Created:** 2026-04-30
**Last Updated:** 2026-04-30

---

## Summary

This feature establishes the MVP authentication and role foundation for the online shop API. It supports first-party credentials registration and login for customers, credentials login for manually bootstrapped staff users, separate customer and staff sessions, JWT access tokens, rotating refresh sessions, seeded `customer` and `staff` roles, and a non-HTTP staff bootstrap path.

---

## Phase Tracker

| Phase     | Status      | Owner            | Notes                                                                                                |
| --------- | ----------- | ---------------- | ---------------------------------------------------------------------------------------------------- |
| PO        | ✅ Accepted | product-owner    |                                                                                                      |
| Tech Lead | ✅ Accepted | tech-lead        | Technical spec written.                                                                              |
| DevOps    | ⚠️ Likely   | devops-engineer  | Complete before Dev; auth secrets, cookie/origin config, and CI test behavior affect implementation. |
| ADR       | ✅ Accepted | adr-author       | Auth/session architecture and migration/bootstrap strategy documented.                               |
| Dev       | ⏳ Pending  | senior-developer |                                                                                                      |
| QA        | ⏳ Pending  | qa-engineer      |                                                                                                      |
| Review    | ⏳ Pending  | code-reviewer    |                                                                                                      |

**Status legend:**

- ✅ Accepted — complete and signed off
- ⏳ Pending — not yet started
- 🔄 In Progress — agent currently working this phase
- ⚠️ Blocked — waiting on dependency or open question
- ⚪ Optional — may be skipped; owner decides at handover
- ❌ Skipped — explicitly bypassed with reason

---

## Artifact Index

| Artifact                                            | Status      | Path                                                      |
| --------------------------------------------------- | ----------- | --------------------------------------------------------- |
| requirements.md                                     | ✅ Accepted | specs/requirements.md                                     |
| technical-spec.md                                   | ✅ Accepted | specs/technical-spec.md                                   |
| infrastructure-design.md                            | ⚠️ Likely   | specs/infrastructure-design.md                            |
| test-plan.md                                        | ⏳ Pending  | specs/test-plan.md                                        |
| ADR-0001-authentication-session-architecture.md     | ✅ Accepted | ../../adr/0001-authentication-session-architecture.md     |
| ADR-0002-database-migrations-and-staff-bootstrap.md | ✅ Accepted | ../../adr/0002-database-migrations-and-staff-bootstrap.md |
| ADR(s)                                              | ✅ Accepted | ../../adr/                                                |
| review-{date}.md                                    | ⏳ Pending  | reviews/                                                  |

---

## Dependencies

No cross-feature dependencies.

---

## Open Blockers

| #   | Blocker                                                                                                                              | Raised By           | Priority |
| --- | ------------------------------------------------------------------------------------------------------------------------------------ | ------------------- | -------- |
| 1   | Complete DevOps phase before Dev so cookie, CORS, secret, CI, rate-limit, and database pool configuration expectations are explicit. | PO/Tech Lead review | High     |

---

## Decision Log

Chronological record of phase completions, handovers, and significant decisions.
Never delete entries.

| Date       | Event                      | Notes                                                                                                  |
| ---------- | -------------------------- | ------------------------------------------------------------------------------------------------------ |
| 2026-04-30 | PO phase accepted          | Feature bootstrapped                                                                                   |
| 2026-04-30 | ADR-0001 accepted          | Authentication session architecture documented                                                         |
| 2026-04-30 | ADR-0002 accepted          | Database migration and staff bootstrap strategy documented                                             |
| 2026-04-30 | Tech Lead phase accepted   | technical-spec.md written                                                                              |
| 2026-04-30 | Tech Lead refinement       | Use timestamped Goose migration names with descriptive titles                                          |
| 2026-04-30 | Tech Lead refinement       | Standardize successful JSON responses with a data envelope                                             |
| 2026-04-30 | PO/Tech Lead review update | Address security, cookie, API contract, migration, and implementation feedback before DevOps/Dev phase |
