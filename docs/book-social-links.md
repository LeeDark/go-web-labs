# Links between `go-web-labs` and `book-social`

This document maps reusable study patterns to concrete applied work in `book-social`. It records
current evidence and deferred links; it is not a promise that every lab feature belongs in the
application.

## Applied baseline

Latest accepted applied ref: `book-social` tag `v0.2.4` at commit `cba82f0`.

The release contains the normalized catalog read side and the v0.2.4 HTTP Foundation. Its current
routes are MPA/catalog routes such as `/`, `/books`, `/books/{slug}`, `/authors/{slug}`, and
`/healthz`. It does not yet contain the planned read-only `/api/*` JSON slice.

Verification:

```bash
GOCACHE=/tmp/book-social-go-cache go test ./...
```

Result: PASS. PostgreSQL parity is opt-in through `BOOK_SOCIAL_POSTGRES_TEST_DSN` and must be
reported separately from the default test run.

## Pattern map

| `go-web-labs` source                                                                                  | Applied `book-social` evidence                                                              | Status and boundary                                                                                                    |
|-------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------|
| `labs/rest-api` — JSON envelopes, validation, strict HTTP boundary, predictable errors                | `internal/http/response`, catalog handlers, and MPA error behavior at `v0.2.4`              | Reusable principles are applied; the application remains MPA and does not yet expose the lab's JSON `/api/*` contract. |
| `labs/layered-api` — handler → service → repository responsibilities                                  | `internal/modules/books/{handler,service,repository}.go` at `v0.2.4`                        | Applied and tested in the catalog module; keep interfaces/use cases domain-specific.                                   |
| `labs/testing` and Stage 5 matrix — unit, HTTP-contract, DB integration boundaries                    | `internal/*_test.go`, `internal/testutil`, SQLite integration and opt-in PostgreSQL helpers | Applied with SQLite as the normal test path; PostgreSQL evidence requires the opt-in DSN.                              |
| `books/lets-go-further` — migrations, PostgreSQL, optimistic concurrency, API notes                   | `db/sqlite`, `db/postgresql`, normalized catalog migrations, repository parity              | Selected database patterns applied; Greenlight's movie API contract is a study reference, not the application domain.  |
| Stage 9A — middleware order, request IDs, logging, recovery, headers, cache policy, timeout, shutdown | `internal/app`, `internal/http/middleware`, and `internal/app/server.go` at `v0.2.4`        | HTTP foundation is applied and tested; later auth-specific security remains in v0.2.5–v0.2.6.                          |
| Stage 6 OpenAPI                                                                                       | No current `/api/*` receiver slice                                                          | Deferred until a real JSON API exists; do not list OpenAPI as current applied evidence.                                |
| Stage 7A auth/MPA security                                                                            | v0.2.5 Auth Foundation and v0.2.6 are the current planned releases                          | Current work; sessions, CSRF, password handling, and protected routes are not v0.2.4 evidence yet.                     |

## What belongs in `book-social`

Apply a pattern when it solves an existing product or maintenance need:

- catalog read models and repository behavior;
- migration and seed discipline;
- clear handler/service/repository boundaries;
- HTTP lifecycle and middleware behavior;
- focused tests around catalog and HTTP contracts;
- later authentication, private library, and API work as their releases define them.

## What remains in `go-web-labs`

Keep isolated examples and source-study details in the lab repository when they would duplicate or
distort the application:

- in-memory REST examples and deliberately strict negative cases;
- Greenlight's movie domain and book-specific source exercises;
- framework comparisons and testing experiments;
- disposable-database teaching notes and matrix/checklist templates.

## Active and deferred applied links

| `book-social` work                               | Related lab/stage                             | Status or gate                                      |
|--------------------------------------------------|-----------------------------------------------|-----------------------------------------------------|
| v0.2.5–v0.2.6 authentication                     | Stage 7A, Learn Go with Tests TDD foundations | Current work; auth/session scope and tests accepted |
| Real read-only `/api/*` slice                    | Stage 6 OpenAPI, Stage 7B API security        | At least one stable JSON API contract               |
| Catalog filtering, sorting, pagination           | Stage 2 Chapter 9                             | `book-social` v0.4 discovery scope                  |
| Library transactions, privacy, PostgreSQL parity | Stages 7–8 and Stage 11 as justified          | Concrete library use cases and disposable DB path   |

## Evidence rule

When a pattern is claimed as applied, record the `book-social` tag or commit, source path, receiving
path, verification command, and any adaptation. A planned roadmap item or a study-project endpoint
is not an applied `book-social` capability.
