# REST API Basics Lab

## Purpose

This lab will be a small API-only backend for practicing the HTTP boundary of
one `books` resource: JSON, CRUD, validation, consistent errors, middleware,
and focused handler tests.

The lab is intentionally independent from `book-social`. It uses an in-memory
store, has no users or database, and does not share domain code with the
applied project.

Current status: **Steps 0–2 complete.** The API skeleton, healthcheck, and v1
read routes are implemented and verified.

## Planned v1 project structure

The lab will stay in one small `main` package until a later stage has evidence
that another package is useful.

```text
labs/rest-api/
├── README.md
├── go.mod                  # Step 1: standalone module
└── cmd/api/
    ├── main.go             # Step 1: config, dependencies, and HTTP server
    ├── routes.go           # Step 1: Chi router; later book routes
    ├── helpers.go          # Step 1: JSON response envelope helper
    ├── health.go           # Step 1: GET /health
    ├── store.go            # Step 1: Book and concurrency-safe memory store
    ├── books.go            # Steps 2–3: read and write handlers and DTOs
    ├── errors.go           # Steps 2 and 4: minimal, then complete API errors
    ├── middleware.go       # Step 4: recovery and/or request logging
    └── handlers_test.go    # Step 5: focused HTTP handler tests
```

Files marked for later steps will not be created as empty placeholders during
Step 1.

## Step 1 implementation

### Module and workspace

- Initialize module `github.com/LeeDark/go-web-labs/labs/rest-api` with the Go
  version already used by the repository workspace.
- Add only `github.com/go-chi/chi/v5` as a direct dependency.
- Add `./labs/rest-api` to the root `go.work` file.
- Do not introduce a Makefile, configuration library, environment loader, or
  another framework.

### Application wiring

- Accept an optional `-addr` flag with `:4000` as the local default.
- Use the standard `log/slog` package for startup and server errors.
- Keep shared dependencies in a small `application` value containing the
  logger and concrete `*bookStore`.
- Initialize one deterministic seed book in `newBookStore()`: ID `1`, using the
  title, author, and description from the contract example; the next ID is `2`.
- Protect the store with `sync.RWMutex`; the HTTP server may call it from
  multiple goroutines.
- Keep the concrete store local to the lab. Do not add an interface or a
  service/repository layer before Stage 4.

### HTTP skeleton

- Build the router in `routes.go` with Chi and register only `GET /health` in
  this step.
- Return the contract response `{"data":{"status":"available"}}` through a
  little JSON writer shared by later handlers.
- Configure an explicit `http.Server` with address, handler, standard-library
  error logging, `5s` read-header, `10s` read, `10s` write, and `60s` idle
  timeouts.
- Keep graceful shutdown, recovery, request logging, custom `404`/`405`
  responses, and the `/books` routes out of this step.

### Run and verification

From `labs/rest-api`:

```sh
go run ./cmd/api
go test ./...
go vet ./...
```

Manual health check while the server is running:

```sh
curl -i http://localhost:4000/health
```

Expected response:

```http
HTTP/1.1 200 OK
Content-Type: application/json

{"data":{"status":"available"}}
```

The JSON writer may append a final newline. Header ordering and additional
standard HTTP headers are not part of the contract.

### Step 1 Definition of Done

- [x] The new module is present in `go.work` and has no dependencies other than
  Chi plus the standard library.
- [x] `go run ./cmd/api` starts the server with no database or environment
  setup.
- [x] `GET /health` returns `200`, the JSON content type, and the documented data
  envelope.
- [x] The application owns an initialized, concurrency-safe in-memory book
  store.
- [x] No `/books` handler, speculative layering, middleware, or later-stage
  feature is implemented.
- [x] `gofmt`, `go test ./...`, and `go vet ./...` succeed in the lab module.
- [x] The README contains commands that have been verified against the
  implemented skeleton.

## Step 2 implementation plan

### Store reads

- Add `list()` and `get(id)` methods directly to the concrete `bookStore`; do
  not introduce a repository interface.
- Protect both operations with `RLock`/`RUnlock`.
- Return `Book` values rather than pointers into the store, so handlers cannot
  mutate stored records accidentally.
- Make `list()` return a non-nil empty slice when there are no books and sort
  the copied values by ascending `id`. Map iteration order must not leak into
  the API response.
- Make `get(id)` return `(Book, bool)`, where the boolean reports whether the
  record exists. The in-memory store has no internal failure to expose yet.

### Routes and handlers

- Add `GET /books` and `GET /books/{id}` directly to the existing Chi router.
  Do not register trailing-slash variants or add redirect middleware.
- Create `books.go` with `listBooksHandler` and `showBookHandler`.
- Add a reusable positive `int64` route-ID parser to `helpers.go`, using
  `chi.URLParam(r, "id")` and base-10 parsing.
- Return list and single-book responses through the existing `data` envelope
  and `writeJSON()` helper.
- Keep the handlers thin: parse HTTP input, call the concrete store, map the
  result to HTTP, and write JSON.

### Minimal read errors

- Create `errors.go` now with the common error shape and one small JSON error
  writer. Step 4 will extend the same file instead of replacing the contract.
- Return `400 invalid_id` when `{id}` is non-numeric, zero, negative, or outside
  the `int64` range.
- Return `404 book_not_found` when a valid positive ID is absent.
- Use the fixed messages `Book ID must be a positive integer` and
  `Book not found`. Error codes are the client-facing stable part.
- Log response encoding or write failures with request method and URI. Do not
  attempt to write a second response after the output may already have started.
- Leave custom router-level `404`/`405`, `internal_error`, recovery, and request
  logging for Step 4.

### Verification in Pair Programmer Mode

Run from `labs/rest-api`:

```sh
gofmt -w ./cmd/api/*.go
go test ./...
go vet ./...
go run ./cmd/api
```

While the server is running, verify:

```sh
curl -i http://localhost:4000/books
curl -i http://localhost:4000/books/1
curl -i http://localhost:4000/books/999
curl -i http://localhost:4000/books/0
curl -i http://localhost:4000/books/not-a-number
```

Expected outcomes:

| Request                   | Status | Response requirement                          |
|---------------------------|-------:|-----------------------------------------------|
| `GET /books`              |  `200` | `data` is an array ordered by ascending `id`. |
| `GET /books/1`            |  `200` | `data` is the deterministic seed book.        |
| `GET /books/999`          |  `404` | `error.code` is `book_not_found`.             |
| `GET /books/0`            |  `400` | `error.code` is `invalid_id`.                 |
| `GET /books/not-a-number` |  `400` | `error.code` is `invalid_id`.                 |

Every response in this table must use `Content-Type: application/json`.

### Step 2 Definition of Done

- [x] Both read routes are registered with Chi and no write route is added.
- [x] Store reads are concurrency-safe and do not expose the internal map.
- [x] List output is a non-nil array in deterministic ascending-ID order.
- [x] Existing, missing, malformed, and non-positive IDs follow the documented
  response envelopes and statuses.
- [x] No request-body decoder, validation DTO, middleware, SQL, or speculative
  layer is introduced.
- [x] `gofmt`, `go test ./...`, and `go vet ./...` succeed.
- [x] All five documented `curl` checks match the contract.

## v1 API contract

### Resource

A book response has this shape:

```json
{
    "id": 1,
    "title": "The Left Hand of Darkness",
    "author": "Ursula K. Le Guin",
    "description": "A science fiction novel."
}
```

| Field         | JSON type | Rules                                                                            |
|---------------|-----------|----------------------------------------------------------------------------------|
| `id`          | number    | Server-generated positive integer; immutable.                                    |
| `title`       | string    | Required; trimmed; 1–200 Unicode characters after trimming.                      |
| `author`      | string    | Required; trimmed; 1–120 Unicode characters after trimming.                      |
| `description` | string    | Optional; trimmed; at most 1000 Unicode characters; defaults to an empty string. |

The lab identifies books by numeric `id`. A human-readable `slug` belongs to
the future applied API in `book-social` and is not part of this contract.

### Input DTOs

Create request:

```json
{
    "title": "The Left Hand of Darkness",
    "author": "Ursula K. Le Guin",
    "description": "A science fiction novel."
}
```

- `title` and `author` are required.
- `description` may be omitted and then defaults to `""`.
- `id` is not accepted from the client.

Partial update request:

```json
{
    "title": "A new title"
}
```

- Every field is optional, but at least one of `title`, `author`, or
  `description` must be present.
- Omitted fields stay unchanged.
- A supplied `title` or `author` must not be blank.
- `description` may be `""` to clear it.
- `null` is not a valid value for any field.
- Sending an unchanged value is valid and returns the current resource.

For both DTOs, the API accepts one JSON object, rejects unknown fields and
trailing JSON values, and limits the body to 1 MiB. `POST` and `PATCH` require
`Content-Type: application/json`; a `charset` parameter is allowed.

### Response envelopes

A successful response with a representation uses a top-level `data` member:

```json
{
    "data": {
        "id": 1,
        "title": "The Left Hand of Darkness",
        "author": "Ursula K. Le Guin",
        "description": "A science fiction novel."
    }
}
```

The list endpoint returns an array in the same envelope. Books are ordered by
ascending `id`, so the in-memory result is deterministic.

```json
{
    "data": [
        {
            "id": 1,
            "title": "The Left Hand of Darkness",
            "author": "Ursula K. Le Guin",
            "description": "A science fiction novel."
        }
    ]
}
```

All error responses use an `error` member with a stable machine-readable
`code` and a human-readable `message`. Validation errors may also include a
`fields` object.

```json
{
    "error": {
        "code": "validation_failed",
        "message": "Request validation failed",
        "fields": {
            "title": "must be provided"
        }
    }
}
```

Clients may depend on `code`, but should not parse `message` or field messages.
Internal errors are logged by the server and are never exposed in the response.
Every response with a JSON body uses `Content-Type: application/json`.

`DELETE` is the one successful operation without an envelope: it returns
`204 No Content` with an empty body.

### Endpoints and statuses

| Method and path      | Success | Expected errors                          | Notes                                                 |
|----------------------|--------:|------------------------------------------|-------------------------------------------------------|
| `GET /health`        |   `200` | `500`                                    | Returns `{"data":{"status":"available"}}`.            |
| `GET /books`         |   `200` | `500`                                    | Returns `data` as an array, including when empty.     |
| `GET /books/{id}`    |   `200` | `400`, `404`, `500`                      | `id` must be a positive base-10 integer.              |
| `POST /books`        |   `201` | `400`, `413`, `415`, `422`, `500`        | Returns the created book and `Location: /books/{id}`. |
| `PATCH /books/{id}`  |   `200` | `400`, `404`, `413`, `415`, `422`, `500` | Applies only supplied fields.                         |
| `DELETE /books/{id}` |   `204` | `400`, `404`, `500`                      | Successful response has no body.                      |

Unsupported methods return `405 Method Not Allowed` in the error envelope and
include the router-generated `Allow` header where applicable.

Error codes are fixed as follows:

| Status | Code                     | Used when                                                                                                                  |
|-------:|--------------------------|----------------------------------------------------------------------------------------------------------------------------|
|  `400` | `invalid_id`             | A route `id` is not a positive base-10 integer.                                                                            |
|  `400` | `invalid_json`           | The body is empty, malformed, not an object, has unknown fields or trailing JSON, or contains a wrong JSON type or `null`. |
|  `404` | `book_not_found`         | No book exists for the supplied valid `id`.                                                                                |
|  `405` | `method_not_allowed`     | The route exists but does not accept the HTTP method.                                                                      |
|  `413` | `request_too_large`      | A JSON request body exceeds 1 MiB.                                                                                         |
|  `415` | `unsupported_media_type` | A JSON request has an unsupported or missing `Content-Type`.                                                               |
|  `422` | `validation_failed`      | The decoded DTO violates field rules or a PATCH object has no supported fields.                                            |
|  `500` | `internal_error`         | An unexpected server failure occurs.                                                                                       |

## Step 0 verification scenarios

These scenarios define the behavior to verify manually with `curl` and, where
noted later, with focused `httptest` cases:

1. `GET /health` returns `200`, JSON content type, and `status: available`.
2. `GET /books` returns `200` and a deterministically ordered JSON array.
3. `GET /books/{existing-id}` returns `200`; a missing positive ID returns the
   `404 book_not_found` envelope; a non-positive or non-numeric ID returns
   `400 invalid_id`.
4. A valid `POST /books` returns `201`, a `Location` header, and the normalized
   created book with a generated ID.
5. A create request with a missing or blank required field returns
   `422 validation_failed` and field details without changing the store.
6. A `PATCH` containing only `title` returns `200` and preserves `author` and
   `description`.
7. An empty PATCH object returns `422 validation_failed`; a malformed body,
   unknown field, wrong JSON type, or `null` returns `400 invalid_json`.
8. `DELETE /books/{existing-id}` returns an empty `204`; a later read or second
   delete returns `404 book_not_found`.
9. An oversized body returns `413`, and a missing or unsupported request media
   type returns `415`, both in the common error envelope.
10. An unsupported method returns `405` in the common error envelope; an
    unexpected server failure returns only the safe `500 internal_error`
    response.

The v1 automated suite will select 5–7 focused cases from this list while the
full list remains the manual contract checklist.

## Scope boundary

The v1 lab will use Go, `net/http`, Chi, and a local in-memory store. It will
not add SQL, migrations, Docker, OpenAPI, authentication, authorization, CORS,
rate limiting, extra resources, or speculative service/repository packages.

A later review may improve this same example, but it must remain a small lab.
Applied catalog endpoints in `book-social` will use their own models and
`/api/books/{slug}` contract rather than importing this code.
