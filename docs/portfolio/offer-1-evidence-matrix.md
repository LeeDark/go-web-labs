# Offer 1 evidence matrix

This document is the working evidence register for **Offer 1 — Go Backend APIs & Integrations**. It
records what is implemented, where it lives, how it was verified, and whether a pattern can move
between repositories.

The matrix distinguishes a reusable pattern from an implemented capability in the receiving project.
A source example, a planned contract, or a branch without an accepted milestone is not evidence of a
finished receiving capability.

## Evidence levels

Offer levels are cumulative client-facing scopes:

| Level        | Current state                       | Claim rule                                                                                                          |
|--------------|-------------------------------------|---------------------------------------------------------------------------------------------------------------------|
| Basic        | Partially evidenced                 | Claim only the focused endpoint/API, validation, errors, tests, and handoff behavior that has a concrete reference. |
| Advanced     | In progress in `go-service-starter` | Claim after the accepted `proof/advanced` milestone and target tag `v0.2-advanced`.                                 |
| Professional | Planned                             | Claim after the accepted `proof/professional` milestone and target tag `v0.3-professional`.                         |
| Ultimate     | Planned                             | Claim after the accepted `proof/ultimate` milestone and target tag `v1.0-ultimate`.                                 |

## Status vocabulary

| Status                 | Meaning                                                                                        |
|------------------------|------------------------------------------------------------------------------------------------|
| `candidate`            | A source pattern was identified, but receiving implementation or evidence is absent.           |
| `transferable-now`     | The pattern can be reused without material adaptation and has a concrete source reference.     |
| `adaptation-required`  | Reuse is plausible, but the receiving architecture, contract, or infrastructure needs changes. |
| `verified-in-receiver` | The receiving implementation exists at a specific ref and the listed verification passes.      |
| `project-specific`     | Keep it in the source project; transferring it would add the wrong domain or architecture.     |
| `blocked`              | A required milestone, test, environment, or decision is missing.                               |

## Required evidence fields

Every non-empty matrix row should identify:

- Offer level;
- capability or pattern;
- source repository, branch/tag, commit, and path;
- receiving repository, branch/tag, commit, and path;
- exact verification command and result;
- adaptation required;
- current status;
- known limitations and the next evidence step.

Use immutable commit or milestone-tag references for portfolio claims. A moving branch may be useful
for a candidate, but it is not sufficient for `verified-in-receiver`.

## Initial matrix

The refs below are the review baseline as of 2026-08-09:

- `go-web-labs@main` at `231be88`;
- `book-social@main` and tag `v0.2.4` at `cba82f0`;
- `go-service-starter@proof/advanced` at `af9a1ee`;
- `go-service-starter@v0.1-basic` at `b3849d3`.

## Basic inventory result

The Basic inventory was checked against the immutable refs above. The verification commands passed:

| Repository/ref                              | Command                                                     | Result | Evidence boundary                                                                                                                                                                                                  |
|---------------------------------------------|-------------------------------------------------------------|--------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `book-social@v0.2.4` (`cba82f0`)            | `GOCACHE=/tmp/book-social-go-cache go test ./...`           | PASS   | Catalog, HTTP foundation, unit, HTTP, and SQLite-backed tests are confirmed. PostgreSQL parity remains opt-in and requires `BOOK_SOCIAL_POSTGRES_TEST_DSN`; this check alone does not claim a live PostgreSQL run. |
| `go-service-starter@v0.1-basic` (`b3849d3`) | `GOCACHE=/tmp/go-service-starter-basic-cache go test ./...` | PASS   | Basic HTTP, response, handler, routing, config, and in-memory store behavior are confirmed. The Basic tag has no CRUD persistence, service/repository boundary, Docker, or OpenAPI proof.                          |

This establishes two different kinds of Basic evidence:

- `book-social` is the applied MPA/catalog proof at `v0.2.4`;
- `go-service-starter` is the focused REST Basic Proof at `v0.1-basic`;
- `go-web-labs` supplies the reusable study and lab patterns at `231be88`.

The inventory is complete for these refs, but it does not upgrade any row to
`verified-in-receiver` unless the receiving capability exists at the cited ref. Most transfers into
`go-service-starter` therefore remain `adaptation-required`.

| Level        | Capability / pattern                                                         | Source ref and path                                                                                             | Receiving ref and path                                                                                        | Verification                                                                                                               | Adaptation / limitation                                                                                                                | Status                 |
|--------------|------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|------------------------|
| Basic        | REST JSON CRUD boundary, strict request shape, validation, and stable errors | `go-web-labs@231be88:labs/rest-api`                                                                             | `go-service-starter@v0.1-basic`: `internal/http/response`, `internal/http/handlers`                           | `GOCACHE=/tmp/go-service-starter-basic-cache go test ./...` — PASS                                                         | Receiver currently exposes only read-only in-memory `GET`; create/update/delete and strict body parsing are not implemented.           | `adaptation-required`  |
| Basic        | Applied MPA catalog and HTTP foundation                                      | `book-social@cba82f0:docs/roadmap.md`, `internal/app`, `internal/http/middleware`                               | `book-social@v0.2.4`: same paths                                                                              | `GOCACHE=/tmp/book-social-go-cache go test ./...` — PASS                                                                   | Confirmed applied evidence, but not a JSON `/api/*` proof; PostgreSQL parity is opt-in.                                                | `verified-in-receiver` |
| Basic        | Handler → service/repository separation                                      | `go-web-labs@231be88:labs/layered-api`                                                                          | `book-social@cba82f0`: `internal/modules/books`; future `go-service-starter@proof/advanced`: `internal/items` | `GOCACHE=/tmp/book-social-go-cache go test ./...` — PASS; no service boundary in Basic tag                                 | Adapt interfaces from the source-consuming package and preserve the `items` domain.                                                    | `adaptation-required`  |
| Basic        | Unit, HTTP-contract, and opt-in PostgreSQL test matrix                       | `go-web-labs@231be88:labs/testing`, `docs/plan-stage5.md`                                                       | `book-social@cba82f0`: `internal/*_test.go`, `internal/testutil`                                              | `GOCACHE=/tmp/book-social-go-cache go test ./...` — PASS; PostgreSQL additionally requires `BOOK_SOCIAL_POSTGRES_TEST_DSN` | `go-service-starter` Basic has unit/`httptest` tests but no real PostgreSQL integration slice.                                         | `adaptation-required`  |
| Basic        | README, test matrix, review checklist, and safe handoff guidance             | `go-web-labs@231be88:labs/testing/README.md`, `labs/testing/test-matrix.md`, `labs/testing/review-checklist.md` | `go-service-starter@v0.1-basic`: `README.md`; `proof/advanced`: `docs/migrations.md`                          | `git show --check`; documented `go test ./...` commands                                                                    | Receiving docs need an Advanced status section, architecture notes, and handoff checklist.                                             | `adaptation-required`  |
| Advanced     | PostgreSQL migration workflow                                                | `go-service-starter@af9a1ee:docs/migrations.md`, `migrations/000001_create_items_table.*`                       | Same `proof/advanced` ref and paths                                                                           | `go test ./...`; disposable PostgreSQL migration smoke is documented but not represented by a real integration test        | Workflow exists; CRUD schema/application integration and accepted `v0.2-advanced` milestone are still missing.                         | `adaptation-required`  |
| Advanced     | PATCH and optimistic version conflict                                        | `go-web-labs@231be88:books/lets-go-further/greenlight/internal/data`, Stage 2/5 notes                           | `go-service-starter@af9a1ee:docs/advanced-contract.md`; no receiving implementation yet                       | Source Greenlight tests pass; receiver has contract only                                                                   | Adapt SQL, error codes, handler mapping, and tests to `items`; do not claim receiver support from the contract alone.                  | `candidate`            |
| Advanced     | PostgreSQL pool lifecycle and bounded startup ping                           | `go-service-starter@af9a1ee:internal/storage/postgres/postgres.go`                                              | `go-service-starter@af9a1ee:internal/storage/postgres`                                                        | `GOCACHE=/tmp/go-service-starter-go-cache go test ./...`                                                                   | This is already implemented in the receiver, but only pool startup/error behavior is covered; it does not prove the full Advanced API. | `verified-in-receiver` |
| Professional | Middleware, request ID/logging, graceful shutdown, and operational handoff   | `book-social@cba82f0:internal/app`, `internal/http/middleware`, `internal/app/server.go`                        | `go-service-starter@proof/advanced`: no receiving middleware/lifecycle implementation yet                     | `GOCACHE=/tmp/book-social-go-cache go test ./...` verifies the source project                                              | Standard-library middleware and lifecycle design must be adapted to the starter; target is `proof/professional`.                       | `candidate`            |
| Professional | OpenAPI/API examples and Compose delivery                                    | `go-web-labs@231be88:PLAN.md` Stage 6/8/9; `book-social@cba82f0` has no real `/api/*` slice yet                 | `go-service-starter@proof/advanced`: planned in `PLAN.md`, not implemented                                    | No receiving verification yet                                                                                              | Requires an actual API contract, Docker Compose, and accepted Professional milestone.                                                  | `blocked`              |
| All levels   | MPA templates, sessions/CSRF, and book-social domain fixtures                | `book-social@cba82f0:internal/web`, auth roadmap, catalog fixtures                                              | None                                                                                                          | Not applicable                                                                                                             | These are project-specific and should not be copied into the REST starter.                                                             | `project-specific`     |

## Reading the matrix

The matrix supports three separate questions:

1. **What was learned?** Use the source ref/path.
2. **What was implemented elsewhere?** Require a receiving ref/path and a passing verification
   command.
3. **What can be reused?** Read the adaptation and limitation fields; `candidate` is not a portfolio
   claim.

When a receiving milestone is accepted, replace the provisional branch ref with its immutable tag or
commit, update the verification result, and add a short note to the corresponding level proof
document. Keep source-only patterns and project-specific decisions in the matrix even when they are
not transferable; they explain the boundary of the portfolio claim.

## Next review actions

1. Keep the Basic inventory anchored to the accepted `book-social` v0.2.4 and `go-service-starter`
   `v0.1-basic` refs; update it only when those evidence refs change.
2. Re-review the Advanced rows after `proof/advanced` implements CRUD, service/repository
   boundaries, and focused PostgreSQL/HTTP tests.
3. Create a level-specific proof document only after its target milestone/tag is accepted.
