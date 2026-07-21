# REST API Basics Lab

## Purpose

This lab will be a small API-only backend for practicing the HTTP boundary of
one `books` resource: JSON, CRUD, validation, consistent errors, middleware,
and focused handler tests.

The lab is intentionally independent from `book-social`. It uses an in-memory
store, has no users or database, and does not share domain code with the
applied project.

Current status: **Step 0 complete — the v1 API contract is fixed; no server has
been implemented yet.** Run and test commands will be added with the skeleton
and tests in later steps.

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
