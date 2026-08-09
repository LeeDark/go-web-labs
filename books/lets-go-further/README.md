# Let's Go Further

Study project based on *Let's Go Further* by Alex Edwards.

## Purpose

This folder is for practicing production-style Go API patterns through the
Greenlight project.

Current focus: production API topics. API Core Chapters 1–8 are complete, and
Chapter 11 is complete as the first Stage 9A unit. Chapter 9 remains deferred.

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
- Graceful shutdown

Out of current scope:

- Authentication
- Authorization
- Rate limiting
- Metrics
- Deployment
- Further production hardening

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

### Chapter 7: CRUD Operations

Chapter 7 connects the movie endpoints to PostgreSQL and completes the first
database-backed CRUD pass.

Implemented pieces:

- `data.Models` as the application-level collection of concrete data models
- `MovieModel.Insert()`, `Get()`, `Update()`, and `Delete()` SQL methods
- PostgreSQL array conversion for movie genres
- `sql.ErrNoRows` mapping to `data.ErrRecordNotFound`
- `POST /v1/movies` with validation, `201 Created`, and a `Location` header
- `GET /v1/movies/:id` backed by the database
- initial `PUT /v1/movies/:id` full replacement, reworked as `PATCH` in
  Chapter 8
- version increment on update
- `DELETE /v1/movies/:id` with affected-row not-found detection
- safe `404` and `500` handler mappings for model errors

Useful pattern:

Keep request decoding and HTTP status handling in handlers, and keep SQL,
driver-specific values, and database error translation in the model layer. A
small concrete model is enough here; a generic repository abstraction would
not improve this study step.

Review follow-ups:

- Align validation messages with their rules: the year 1888 is accepted even
  though the message says “greater than 1888”, and the genre rule allows up to
  5 even though its message says “more than 1 genres”.
- Return the error from `ResponseWriter.Write()` in `writeJSON()` so handler
  logging can observe failed response writes.

### Chapter 8: Advanced CRUD Operations

Chapter 8 replaces full movie updates with partial updates and protects records
from stale concurrent writes.

Implemented pieces:

- `PATCH /v1/movies/:id` for partial updates
- pointer-based input fields, preserving stored values when fields are omitted
- validation after merging the partial input into the stored movie
- three-second context timeouts for every database operation
- atomic version increment during updates
- SQL optimistic locking using both movie ID and version
- `409 Conflict` response for a stale or conflicting update
- optional `X-Expected-Version` round-trip version check before the update

Useful pattern:

Use a pointer-based PATCH input type to distinguish omitted scalar fields from
zero values. For data integrity, perform the final version comparison in the
same SQL `UPDATE`; a preliminary request-header comparison is helpful to the
client but cannot replace the atomic database check.

Review note:

If a movie is deleted after the handler reads it but before the update query,
the update affects no row and returns `409 Conflict`. This is a safe and
intentional simplification for the study project; distinguishing it from a
stale version would need another database read.

### Chapter 11: Graceful Shutdown

Chapter 11 is complete. The API now moves HTTP server setup and lifecycle into
`cmd/api/server.go` and handles `SIGINT` and `SIGTERM` with a graceful
shutdown sequence.

Implemented pieces:

- `app.serve()` owns the `http.Server` lifecycle;
- a buffered signal channel receives `SIGINT` and `SIGTERM`;
- `srv.Shutdown()` stops accepting new requests and allows in-flight requests
  up to 30 seconds to finish;
- `http.ErrServerClosed` is treated as the expected result of shutdown;
- shutdown errors are returned to `main()` and logged consistently;
- startup and completed-shutdown events are logged.

Short review:

The implementation follows the chapter's intended signal → shutdown → error
channel flow and keeps `main()` focused on application wiring. The 30-second
shutdown context protects the process from waiting forever. The next production
follow-up is to coordinate background goroutines explicitly; `http.Server`
shutdown does not wait for them. Chapter 13.5 covers that extension. A focused
shutdown test is also a useful future exercise.

## Local PostgreSQL and Migrations

The API opens a PostgreSQL connection during startup, so create the local role
and database before running it. In a `psql` session with sufficient privileges:

```sql
CREATE ROLE greenlight WITH LOGIN PASSWORD 'greenlight';
CREATE DATABASE greenlight OWNER greenlight;
\c greenlight
CREATE EXTENSION IF NOT EXISTS citext;
```

Set the DSN for both the API and the migration tool:

```sh
export GREENLIGHT_DB_DSN='postgres://greenlight:greenlight@localhost/greenlight?sslmode=disable'
```

This project uses the `migrate` CLI from `golang-migrate`. Install the pinned
version once (ensure your Go bin directory is on `PATH`):

```sh
go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1
```

From this folder, apply the migrations or roll back the latest migration:

```sh
migrate -path=./greenlight/migrations -database="$GREENLIGHT_DB_DSN" up
migrate -path=./greenlight/migrations -database="$GREENLIGHT_DB_DSN" down 1
```

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

curl -i -X POST localhost:4000/v1/movies \
  -H 'Content-Type: application/json' \
  -d '{"title":"Casablanca","year":1942,"runtime":"102 mins","genres":["drama"]}'

curl -i localhost:4000/v1/movies/1

curl -i -X PATCH localhost:4000/v1/movies/1 \
  -H 'Content-Type: application/json' \
  -d '{"genres":["drama","romance"]}'

curl -i -X PATCH localhost:4000/v1/movies/1 \
  -H 'Content-Type: application/json' \
  -H 'X-Expected-Version: 1' \
  -d '{"title":"Casablanca (restored)"}'

curl -i -X DELETE localhost:4000/v1/movies/1
```

## How to Test

From this folder:

```sh
go test ./...
go vet ./...
```

These commands currently provide compilation and static checks. Focused model
or handler tests have not been added yet.

## TODO

- Align the Chapter 7 year and genre validation messages with their rules.
- Propagate response write errors from `writeJSON()`.
- Add focused model or handler tests when an appropriate test boundary is in
  place.
