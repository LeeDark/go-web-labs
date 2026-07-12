# Let's Go Further

Study project based on *Let's Go Further* by Alex Edwards.

## Purpose

This folder is for practicing production-style Go API patterns through the
Greenlight project.

Current focus: API Core.

## Current Scope

- Project structure
- HTTP routing
- JSON responses
- JSON request parsing
- Validation
- PostgreSQL
- Migrations
- CRUD
- Filtering, sorting, and pagination

Out of current scope:

- Authentication
- Authorization
- Rate limiting
- Metrics
- Deployment
- Production hardening

## Chapter Notes

### Chapter 1: Introduction

Greenlight is introduced as a JSON API for movie data. The chapter defines the
project direction, expected tools, and the high-level API shape.

The main value for this repository is orientation: use Greenlight to learn
reusable API patterns, while keeping the study project separate from future
portfolio applications.

Expected early endpoints:

- `GET /v1/healthcheck`
- Movie CRUD endpoints under `/v1/movies`

### Chapter 2: Getting Started

Chapter 2 sets up the first runnable API skeleton.

Implemented pieces:

- `cmd/api` application entry point
- `config` struct with `port` and `env` command-line flags
- `application` struct for handler dependencies
- structured logger using `log/slog`
- HTTP server with basic timeout settings
- centralized routes in `routes.go`
- `GET /v1/healthcheck`
- placeholder movie routes:
  - `POST /v1/movies`
  - `GET /v1/movies/:id`
- `readIDParam()` helper for positive integer URL IDs

Useful pattern:

Keep route definitions, handlers, and small helpers separated early. This keeps
`main.go` focused on wiring configuration, dependencies, logging, and the HTTP
server.

### Chapter 3: Sending JSON Responses

Chapter 3 replaces plain-text API responses with JSON responses.

Implemented pieces:

- reusable `writeJSON()` helper
- `envelope` type for consistent response shapes
- JSON healthcheck response
- dummy movie JSON response under a `movie` envelope
- `Movie` response struct with JSON tags
- hidden internal fields with `json:"-"`
- omitted empty or zero-value fields where appropriate
- custom `Runtime` JSON formatting
- reusable JSON error response helpers
- custom JSON `404 Not Found` and `405 Method Not Allowed` responses
- panic recovery middleware for JSON `500 Internal Server Error` responses

Useful pattern:

Encode the response before writing headers, then send every successful and error
response through one small JSON helper. This keeps the API response contract
consistent and makes later metadata or headers easier to add.

### Chapter 4: Parsing JSON Requests

Chapter 4 adds JSON request parsing and validation for `POST /v1/movies`.

Implemented pieces:

- reusable `readJSON()` helper
- strict request body handling
- `400 Bad Request` JSON errors for invalid JSON input
- custom `Runtime` JSON decoding
- reusable `internal/validator` package
- movie validation with `422 Unprocessable Entity` responses

Useful pattern:

Decode into a small input struct, copy into the domain model, then validate the
domain value before continuing.

### Chapter 5: Database Setup and Configuration

Chapter 5 introduces PostgreSQL and configures the API database connection pool.

Implemented pieces:

- `GREENLIGHT_DB_DSN` / `-db-dsn` configuration
- `lib/pq` PostgreSQL driver
- `openDB()` helper with `PingContext()`
- connection pool settings for open, idle, and idle-timeout connections

Useful pattern:

Keep database credentials outside code, verify the connection at startup, and
configure pool limits explicitly.

### Chapter 6: SQL Migrations

Chapter 6 adds versioned SQL migrations for the Greenlight database schema.

Implemented pieces:

- `movies` table migration
- rollback migration for the `movies` table
- movie check-constraint migration
- rollback migration for the constraints

Useful pattern:

Keep schema changes as ordered `up` and `down` files so the database can be
recreated or rolled back predictably.

## How to Run

From this folder:

```sh
go run ./greenlight/cmd/api
```

With custom flags:

```sh
go run ./greenlight/cmd/api -port=3030 -env=production
```

Example requests:

```sh
curl -i localhost:4000/v1/healthcheck
curl -i -X POST localhost:4000/v1/movies
curl -i localhost:4000/v1/movies/123
```

## How to Test

TODO: add when tests are introduced.

## TODO

- Verify Chapter 6 migrations locally.
- Add CRUD notes in Chapter 7.
- Connect movie handlers to database persistence in Chapter 7.
