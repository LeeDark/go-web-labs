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

## Follow-up: v2 as a review-driven increment

v2 is deferred until practical API work identifies a small, proven improvement. It may include no
more than two or three focused refinements, such as clearer DTO/response separation, more precise
error codes, additional PATCH cases, or negative tests. It must stay an in-memory `books` lab and
avoid a database, new domain area, or premature layering refactor.

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
