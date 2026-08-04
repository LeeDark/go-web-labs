# Layered API Lab

Small Go API lab for practicing the handler → service → repository boundary. It is intentionally
separate from `labs/rest-api` and is not a template to copy unchanged into `book-social`.

## Scope and definition of done

The contract has one `books` resource and five use cases: list, get by ID, create, partial update,
and delete. It is complete when the API runs without external infrastructure, its focused checks
pass, and a write flow can be traced from HTTP to service to repository and back.

`Book` has server-owned `id`, `title`, and `author`. HTTP DTOs are private to the handler package;
the domain type has no JSON tags. The service trims title and author, requires both values, and
rejects a case-insensitive duplicate title/author pair.

## Dependency direction

```text
HTTP handler -> books.Service -> books.BookRepository -> MemoryRepository
                    ^
                    |
                 main wires dependencies
```

- The handler parses HTTP input and maps application errors to JSON and status codes.
- The service owns validation and the duplicate rule; it knows neither JSON nor `http.Request`.
- The repository stores and loads domain values; it has no HTTP behavior.
- `main` is the explicit composition root.

## Routes and responses

| Route                | Success            | Important failures                                                                                                                    |
|----------------------|--------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| `GET /health`        | `200`              | —                                                                                                                                     |
| `GET /books`         | `200`              | `500 internal_error`                                                                                                                  |
| `GET /books/{id}`    | `200`              | `400 invalid_id`, `404 book_not_found`                                                                                                |
| `POST /books`        | `201` + `Location` | `400 invalid_request`, `415 unsupported_media_type`, `422 validation_failed`, `409 duplicate_book`                                    |
| `PATCH /books/{id}`  | `200`              | `400 invalid_id`/`invalid_request`, `415 unsupported_media_type`, `422 validation_failed`, `409 duplicate_book`, `404 book_not_found` |
| `DELETE /books/{id}` | `204` with no body | `400 invalid_id`, `404 book_not_found`                                                                                                |

Successful book responses use `{"data": ...}`. Failures use
`{"error":{"code":"...","message":"..."}}`. Details of unexpected repository errors are logged and
never returned to the client.

`POST` and `PATCH` require `Content-Type: application/json` (parameters such as `charset=utf-8` are
accepted). Missing, malformed, or other media types return `415 unsupported_media_type`.

`PATCH` accepts an object with at least one of `title` or `author`. An omitted field is preserved;
an explicitly empty string is merged and then fails validation with `422`. An empty object, a
top-level `null`, a field with `null`, unknown fields, malformed JSON, or a trailing JSON value
returns `400 invalid_request`. The resulting normalized title/author pair is still checked
case-insensitively for duplicates.

## Run and verify

Run these commands from this folder:

```sh
go run ./cmd/api
go test ./...
go vet ./...
```

Examples, with the server running on its default `:4001`:

```sh
curl http://localhost:4001/books
curl http://localhost:4001/books/1
curl -i -X POST http://localhost:4001/books \
  -H 'Content-Type: application/json' \
  -d '{"title":"Go in Action","author":"William Kennedy"}'
curl -i -X PATCH http://localhost:4001/books/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"The Go Programming Language, Second Edition"}'
curl -i -X DELETE http://localhost:4001/books/1
```

## Limitations and learned patterns

This lab deliberately excludes SQL, migrations, auth, pagination, filtering, OpenAPI, generic
repositories, and DI containers. `MemoryRepository` is safe for concurrent access but is only a
learning adapter. The reusable lesson is to introduce a service when a concrete use case needs
business behavior, not merely to add another package.
