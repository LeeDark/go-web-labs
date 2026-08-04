# Stage 4 Plan: Handler → Service → Repository

## Status

**v1 complete.** Step 0 confirmed that the full lab remains useful because `book-social`
demonstrates layering for MPA catalog reads but not this JSON create flow, in-memory adapter, or
duplicate rule.

Stage 4 will create a small standalone lab in `labs/layered-api`. Its purpose is to practice clear
backend boundaries without turning a learning project into an over-engineered architecture exercise.

It is not a rewrite of `labs/rest-api` and will not import its code. Stage 3 is the reference for
the HTTP boundary; this lab focuses on where HTTP, application logic, and persistence
responsibilities belong.

## Goal and scope

The lab will use one `books` resource and three use cases: list, get by ID, and create. It will
demonstrate:

- [x] Thin HTTP handlers that translate requests and responses.
- [x] An application service that owns use cases and business rules.
- [x] A small repository interface and a concurrency-safe in-memory implementation.
- [x] Separate HTTP DTOs and domain types.
- [x] Explicit, ordinary dependency injection in `main`.
- [x] Focused service unit tests and HTTP handler tests.
- [x] A README that explains responsibility boundaries and runnable checks.

The create use case will reject a duplicate normalized `title` + `author`
combination. This gives the service layer one real business rule that does not belong in a handler
or an in-memory map.

## Architecture boundary

| Layer      | Responsibility                                                           | Must not do                                             |
|------------|--------------------------------------------------------------------------|---------------------------------------------------------|
| Handler    | Parse HTTP input, call a service, map results to JSON and status codes.  | Contain business rules or access a repository directly. |
| Service    | Execute use cases, normalize business input, enforce the duplicate rule. | Depend on HTTP or JSON details.                         |
| Repository | Read and store domain books.                                             | Decide HTTP behavior or validate request DTOs.          |
| `main`     | Construct the repository, service, handlers, and router.                 | Hide dependencies in globals or a DI framework.         |

The planned structure stays deliberately small:

```text
labs/layered-api/
├── cmd/api/
│   ├── main.go
│   └── routes.go
└── internal/
    ├── books/
    │   ├── domain.go
    │   ├── service.go
    │   ├── repository.go
    │   └── memory_repository.go
    └── http/
        ├── handlers/
        └── middleware/
```

Files will be created only when they have a concrete purpose.

## API contract

```text
GET   /health
GET   /books
GET   /books/{id}
POST  /books
```

- [x] Define the exact resource fields, request/response DTOs, JSON envelopes, statuses, and error
  codes in the lab README before implementation.
- [x] Return `201 Created` and `Location: /books/{id}` for a valid create.
- [x] Return `415 unsupported_media_type` when create does not use `application/json`.
- [x] Return `422 validation_failed` for blank required fields after trimming.
- [x] Return `409 duplicate_book` for a duplicate normalized title/author pair.
- [x] Return `404 book_not_found` for an unknown valid ID.
- [x] Keep malformed JSON and invalid IDs at the HTTP boundary as `400` errors.
- [x] Return only a safe `500 internal_error` for unexpected repository errors.

## Work plan

### 0. Confirm the scope

- [x] Compare the current applied-project structure with this learning goal.
- [x] Decide whether a full lab is still needed, or whether notes plus a tiny example would provide
  enough evidence.
- [x] Limit v1 to one resource, three use cases, and one business rule.
- [x] Record the definition of done and verification scenarios.

### 1. Create the minimal application skeleton

- [x] Create the standalone module and add it to the workspace.
- [x] Use Chi and an explicit HTTP server setup with timeouts.
- [x] Wire the memory repository, service, and handlers explicitly in `main`.
- [x] Add `/health` and documented run/test commands.

### 2. Implement the domain and repository

- [x] Define the small `Book` domain type and a create input without HTTP tags.
- [x] Define only repository operations the service needs: list, get, and create.
- [x] Implement a concurrency-safe in-memory repository with deterministic seed data and
  server-owned IDs.
- [x] Keep persistence behavior separate from HTTP responses.

### 3. Implement service use cases

- [x] Implement list and get through the service.
- [x] Implement create with trimming, required-field validation, and duplicate detection.
- [x] Keep error mapping on a domain/application boundary, not in HTTP code.
- [x] Add table-driven service tests for valid creation, validation failure, duplicates, missing
  records, and repository failures.

### 4. Implement the HTTP adapter

- [x] Add separate request and response DTOs.
- [x] Implement list, get, and create handlers that depend on the service only.
- [x] Map domain outcomes to the documented JSON status/error contract.
- [x] Apply only the needed HTTP-boundary protections: strict JSON input, positive-ID parsing, safe
  errors, and recovery.

### 5. Verify and document

- [x] Add `httptest` coverage for list, existing/missing get, valid/invalid create, and duplicate
  `409` behavior.
- [x] Confirm response envelopes, `Content-Type`, `Location`, and safe `500`
  behavior.
- [x] Document the dependency direction, responsibility boundaries, commands,
  `curl` examples, limitations, and learned patterns.
- [x] Run `gofmt`, `go test ./...`, `go vet ./...`, and `git diff --check`.

## Out of scope

- [ ] Do not add SQL, migrations, Docker, caching, queues, or external services.
- [ ] Do not add users, authentication, roles, or extra resources.
- [ ] Do not add PATCH, DELETE, pagination, filtering, or OpenAPI.
- [ ] Do not introduce a generic repository, DI container, service locator, event bus, CQRS, or
  interfaces without a direct consumer.
- [ ] Do not copy this structure into an applied project without reviewing its actual use cases.

## Definition of done

- [x] The lab runs from its README without external infrastructure.
- [x] The focused test suite passes.
- [x] One create flow can be traced cleanly from HTTP handler to service to repository and back.
- [x] The duplicate rule is covered by a service test and is absent from the handler and repository.
- [x] The README makes the responsibility boundary understandable to a reviewer.

## Handoff to later work

Stage 4 will provide clear seams for the testing lab: service unit tests, handler tests, and an
in-memory adapter. It also provides criteria for deciding when a real application needs a catalog
service, repository, or read model; it does not prescribe that its package layout must be copied
wholesale.
