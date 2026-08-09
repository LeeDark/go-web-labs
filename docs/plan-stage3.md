# Stage 3 Plan: REST API Basics

## Status

**v1 complete (Steps 0–5).**

Stage 3 produced a small, standalone REST API lab in `labs/rest-api`. It is a practice and reference
project for common Go backend work: adding or fixing an endpoint, validating JSON input, and
returning predictable HTTP responses.

The lab deliberately remains independent of applied projects. It uses in-memory storage only and
does not contain users, authentication, a database, or a second application domain.

## Goal

Build one small, understandable, and testable API-only backend around a single
`books` resource. The completed v1 demonstrates JSON CRUD, validation, consistent errors,
lightweight middleware, documented manual checks, and focused handler tests.

## Delivered in v1

- [x] Defined the API contract: resource fields, input DTOs, success and error envelopes, status
  codes, and verification scenarios.
- [x] Created `labs/rest-api` as a separate Go module in the repository workspace.
- [x] Added a Chi router, explicit HTTP server timeouts, `GET /health`, and a concurrency-safe
  in-memory book store with deterministic seed data.
- [x] Implemented `GET /books` and `GET /books/{id}` with deterministic list ordering, route-ID
  validation, and not-found responses.
- [x] Implemented `POST /books`, `PATCH /books/{id}`, and `DELETE /books/{id}`. IDs are
  server-generated; PATCH is partial and DELETE returns `204 No Content`.
- [x] Hardened the HTTP boundary: JSON media-type checks, a 1 MiB request limit, one strict JSON
  object, unknown-field and trailing-value rejection, Unicode field limits, and validation details.
- [x] Added consistent JSON errors for invalid input, missing routes/resources, unsupported methods,
  oversized requests, unsupported media types, and safe internal failures.
- [x] Added recovery middleware that logs panics without exposing internal details to clients.
- [x] Added seven focused `httptest` handler tests covering CRUD behavior and HTTP-boundary
  failures.
- [x] Documented the runnable API contract, commands, and `curl` scenarios in
  `labs/rest-api/README.md`.

## API surface

```text
GET    /health
GET    /books
GET    /books/{id}
POST   /books
PATCH  /books/{id}
DELETE /books/{id}
```

Successful representation responses use a `data` envelope. Errors use an
`error` envelope with a stable machine-readable `code`, a human-readable
`message`, and validation field details when applicable.

```json
{
    "data": {
        "id": 1,
        "title": "Example Book"
    }
}
```

```json
{
    "error": {
        "code": "validation_failed",
        "message": "Request validation failed"
    }
}
```

## Scope boundary

Stage 3 intentionally does not add:

- SQL, migrations, Docker, or OpenAPI;
- authentication, authorization, CORS, or rate limiting;
- extra resources or user flows;
- service/repository interfaces or package splitting before there is a concrete need.

These topics belong to later, focused stages. The lab is a reference, not a shared package: reuse
its ideas and test cases only after reviewing them for the receiving application.

## Follow-up: v2 technical-debt review

**Status: completed on 2026-08-09 for the selected contract-regression scope.**

The response-mapping and atomic-PATCH candidates remain deferred until a
proven need arises.

v2 is a small, review-driven maintenance increment for this same lab. It is not a second Stage 3,
and it is not an automatic follow-on from a version number in `book-social`. Start it only when a
current API slice or a focused review demonstrates a concrete gap worth carrying into a real
application. Keep the lab independent: it must neither import nor mirror `book-social` code or
models.

The completed review selected contract-regression coverage only: handler tests now assert the full
public book representation for list, read, create, and PATCH, and prove that an invalid PATCH or
oversized request does not mutate the store. The HTTP API behavior is unchanged.

### Why a v2 review is still useful

The original private plan listed DTO separation, PATCH semantics, error codes, deterministic data,
and negative tests as possible v2 work. The v1 implementation has already closed much of that list:

- create and PATCH input DTOs are separate from request handling;
- PATCH distinguishes omitted fields from supplied values, rejects `null`, and validates field
  limits;
- JSON errors, validation errors, routing errors, and safe internal errors have stable codes;
- the store has a deterministic seed and ordered list results;
- focused handler tests cover CRUD and representative HTTP-boundary failures.

Those items are therefore **not** a backlog to reimplement. The remaining debt is chiefly about
making the established contract easier to preserve as the example is reviewed or adapted:

1. The public response is serialized directly from the storage `Book` type. If the storage shape
   changes, an accidental API-contract change is possible.
2. `PATCH` performs a read followed by a separate update. v1 deliberately uses last-write-wins
   behavior, but a review should decide whether the lab needs an explicit atomic mutation boundary
   or a documented decision to retain that trade-off.
3. The tests assert representative failures, but do not systematically protect every documented
   contract edge or each refinement chosen for v2.

### Entry criteria

Before editing code, record a short review note that names the observed API use case or concrete
maintenance risk, identifies the affected v1 behavior, and selects **at most three** refinements
from the debt list below. If no concrete gap is found, close the review without changing the lab.

Do not use an old `book-social` release milestone as the trigger. A real, current API requirement or
an evidence-backed review finding is the trigger; any applied API work remains a separate task.

### Candidate refinements

Choose only the items justified by the entry review:

| Debt                         | Smallest acceptable refinement                                                                                                                       | Required proof                                                                                                                          |
|------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| Response/storage coupling    | Add a private response DTO or mapping function and keep the existing JSON representation unchanged.                                                  | Handler tests assert the complete response shape for list, read, create, and PATCH.                                                     |
| PATCH read–update gap        | Either move the read/apply/update sequence into one lock-protected store operation, or document why last-write-wins remains sufficient for this lab. | A focused store or handler test demonstrates the chosen behavior; no optimistic concurrency, version fields, or database is introduced. |
| Contract-regression coverage | Add table-driven negative and boundary cases that correspond to the selected refinement or an uncovered documented rule.                             | Tests assert status, JSON error code, relevant headers, and lack of unintended store mutation where applicable.                         |

More precise error codes are not a default v2 task: change them only if a real client ambiguity is
identified, preserve the common error envelope, and document any contract change explicitly.

### Implementation sequence

1. Review the current lab contract, handlers, store, and handler tests against the triggering use
   case. Write down the selected one to three refinements and the invariants that must not change.
2. Make the smallest local change in `labs/rest-api`; retain one `main` package and the concrete
   in-memory store.
3. Add or tighten only the focused tests needed to prove the selected behavior. Do not increase test
   count merely to satisfy a quota.
4. Update `labs/rest-api/README.md` only if behavior, a documented trade-off, or a manual scenario
   changes. Add a concise “what changed and why” note for v2.
5. Review the diff as a standalone increment. Do not combine it with `book-social` work or a
   later-stage refactor.

### Boundaries

v2 must remain an in-memory `books` lab. It must not add a database, migrations, Docker, OpenAPI,
authentication, authorization, CORS, rate limiting, another resource, a new user flow, package
splitting, or service/repository interfaces. Do not turn the PATCH refinement into versioning or
optimistic concurrency; those are separate production/API concerns.

Any applied increment in `book-social` is separately scoped. A future read API uses its own catalog
model and `/api/books/{slug}` contract. Write operations require a defined use case plus validation,
authentication, authorization, ownership rules, pagination where relevant, and an error contract;
they must not be copied from this CRUD lab for symmetry.

### Definition of done and verification

v2 is done when the review note shows a proven reason for every selected refinement, the diff is
small and independently reviewable, the public contract remains unchanged unless explicitly
documented, and the README/tests reflect the final behavior.

Run from `labs/rest-api`:

```sh
gofmt -w ./cmd/api/*.go
go test ./...
go vet ./...
git diff --check
```

For any changed endpoint, rerun the corresponding documented `curl` scenario from
[`labs/rest-api/README.md`](../labs/rest-api/README.md). In restricted environments, use the focused
`httptest` suite and report the skipped listener-based check.

## Definition of done and verification

v1 is complete when the lab can be run and checked from its README without an external database, its
contract is documented, and its focused tests pass.

Run from `labs/rest-api`:

```sh
go test ./...
go vet ./...
go run ./cmd/api
```

The complete endpoint contract and manual `curl` examples are maintained in
[`labs/rest-api/README.md`](../labs/rest-api/README.md).
