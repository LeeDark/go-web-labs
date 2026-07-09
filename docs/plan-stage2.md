# Stage 2 Plan: Let's Go Further API Core

## Mode

Book Study Mode.

Use this plan to read and practice the first 9 chapters of *Let's Go Further* step-by-step.
The goal is to convert the source material into small practical tasks for `go-web-labs`,
without letting the book blindly dictate future project architecture.

## Scope

Source:

- Private study copy of *Let's Go Further* by Alex Edwards.
- Repository roadmap: `PLAN.md`, Stage 2 — `Let's Go Further`: API Core
- AI workflow guide: `docs/ai/ai-augmented-development-workflow.md`

Target folder:

- `books/lets-go-further`

Current stage goal:

- Practice API project structure.
- Practice JSON responses and JSON requests.
- Practice validation.
- Practice PostgreSQL setup and migrations.
- Practice CRUD.
- Practice filtering, sorting, and pagination.
- Keep reusable notes for future Go backend API work.

Out of scope for this plan:

- Chapters 10+.
- Rate limiting.
- Graceful shutdown.
- User registration, email, activation, authentication, authorization.
- CORS, metrics, deployment, and production hardening.
- Rebuilding an applied portfolio project inside this study lab.
- Introducing a different framework.
- Large rewrites of existing Stage 1 `books/lets-go` work.

## Book Study Rules

- [ ] Summarize each chapter in my own words.
- [ ] Extract key concepts instead of copying large source fragments.
- [ ] Convert each topic into one or more small practice tasks.
- [ ] Classify each learned idea as `apply now`, `postpone`, or `ignore`.
- [ ] Connect each useful pattern to future Go backend API work.
- [ ] Use AI for explanation, review, planning, and focused edits.
- [ ] Do not use AI as a substitute for actually implementing and testing the code.
- [ ] Switch to Tutor, Pair Programmer, or Manager Mode only when implementation work needs it.

## Repository Setup Checklist

- [x] Create `books/lets-go-further`.
- [x] Create a small README for the stage.
- [x] Create a Go module for the practice API if the implementation starts.
- [x] Add the module to `go.work` if needed.
- [ ] Record run commands in the README.
- [ ] Record test commands in the README.
- [ ] Keep private book files out of public documentation.
- [ ] Keep implementation close to the book while noting where future applied projects may differ.

## Per-Chapter Study Loop

For every chapter below:

- [ ] Read the chapter.
- [ ] Write a short summary in my own words.
- [ ] List the main concepts.
- [ ] Implement or sketch the smallest useful practice task.
- [ ] Run the relevant command or test when possible.
- [ ] Record what worked, what failed, and what to revisit.
- [ ] Classify patterns as `apply now`, `postpone`, or `ignore`.

## Chapter 1: Introduction

### Summary

Chapter 1 introduces the Greenlight project: a JSON API for managing movie data.
The chapter is mostly orientation. It explains the project goal, the expected
API shape, the use of PostgreSQL for persistence, and the tools/background that
will be useful while following the book.

For this repository, Greenlight is a study project. The goal is to extract
practical API patterns for future Go backend work, not to treat the book
architecture as the only correct design for every project.

### Topics

- Greenlight API purpose.
- API endpoint shape.
- PostgreSQL as persistent storage.
- Prerequisites and expected background.
- Go version and supporting tools.

### Main Concepts

- Greenlight is a REST-style JSON API.
- The core resource is `movies`.
- PostgreSQL will be used as persistent storage.
- The book assumes basic comfort with Go, HTTP, SQL, and terminal tools.
- The book gradually moves from API basics toward production concerns.
- Not every later topic belongs to the current Stage 2 scope.

### API Shape Noted

Core API endpoints for the first API Core pass:

- `GET /v1/healthcheck`
- `GET /v1/movies`
- `POST /v1/movies`
- `GET /v1/movies/:id`
- `PATCH /v1/movies/:id`
- `DELETE /v1/movies/:id`

Later endpoints for users, authentication, tokens, metrics, and debug behavior
are useful, but outside the first API Core pass.

### Checklist

- [x] Capture the purpose of the book project in my own words.
- [x] List the core API endpoints that will matter for chapters 1-9.
- [x] Identify which endpoints belong to the API Core scope.
- [x] Check the local Go version.
- [x] Check local availability of `curl`.
- [x] Decide whether to install optional tools only when needed.
- [x] Note the difference between studying Greenlight and building portfolio projects.

### Classification

- [x] Apply now: API endpoint map and study scope.
- [x] Postpone: production deployment endpoint concerns.
- [x] Ignore for now: anything from later user/auth/deployment chapters.

### Personal Takeaway

This chapter defines the boundary for the work. I am studying a
production-style API project, but for now I only need the API core: structure,
JSON handling, validation, database access, CRUD, filtering, sorting, and
pagination.

## Chapter 2: Getting Started

### Summary

Chapter 2 sets up the first working version of the Greenlight API. The project
gets a small Go module, an API entry point under `cmd/api`, a basic HTTP server,
configuration from command-line flags, structured logging with `slog`, and an
`application` struct for sharing dependencies with handlers.

The chapter also introduces the first API routes: a healthcheck endpoint and
placeholder movie endpoints. Routing is moved into a dedicated `routes.go` file
and handled with `httprouter`, which gives clean path parameters, method-aware
routing, automatic `405 Method Not Allowed` responses, and useful `OPTIONS`
behavior.

### Topics

- Project setup and skeleton structure.
- Basic HTTP server.
- Configuration with command-line flags.
- Structured logging.
- Dependency injection through an application struct.
- API endpoints and RESTful routing.
- Route parameter parsing.

### Main Concepts

- Keep executable application code under `cmd/api`.
- Add `internal`, `migrations`, `bin`, and deployment-related folders only when
  they have a real purpose in this repository.
- Use command-line flags for simple environment-specific configuration.
- Store configuration and shared dependencies on an `application` struct.
- Define handlers and helpers as methods on `application` for consistency.
- Use server timeouts from the beginning.
- Keep route declarations in one place with `routes.go`.
- Use `/v1` path prefixing for visible API versioning.
- Prefer REST-style resource URLs like `/v1/movies/:id`.
- Use HTTP methods semantically: `GET` reads, `POST` creates, `PUT` replaces,
  `PATCH` partially updates, and `DELETE` deletes.
- Extract repeated route parameter parsing into a small helper.

### API Shape Noted

Implemented in this chapter:

- `GET /v1/healthcheck`
- `POST /v1/movies`
- `GET /v1/movies/:id`

Planned for later API Core chapters:

- `GET /v1/movies`
- `PATCH /v1/movies/:id`
- `DELETE /v1/movies/:id`

### Checklist

- [x] Create the initial project skeleton under `books/lets-go-further`.
- [x] Add `cmd/api`.
- [x] Add `internal` only when there is code that belongs there.
- [x] Add `migrations` only when database work begins.
- [x] Add a basic `main.go`.
- [x] Add configuration fields for port and environment.
- [x] Add a structured logger.
- [x] Add an application struct for shared dependencies.
- [x] Start a basic HTTP server.
- [x] Add `GET /v1/healthcheck`.
- [x] Add route definitions in a dedicated routes file.
- [x] Add placeholder movie routes.
- [x] Add a helper for reading positive integer `id` route parameters.
- [x] Verify the healthcheck endpoint with `curl`.
- [x] Record the run command in the README.

### Classification

- [x] Apply now: project skeleton, config, logger, application struct, healthcheck.
- [x] Apply now: RESTful route shape for movie resources.
- [x] Postpone: production server scripts and deployment folders.
- [x] Ignore for now: framework comparison.

### Personal Takeaway

This chapter establishes the first useful API skeleton. The most reusable
patterns are the `config` struct, the `application` dependency container,
centralized routes, explicit server timeouts, REST-style endpoint naming, and
the small helper for parsing positive integer URL IDs.

The current responses are intentionally plain text. JSON responses and JSON
error behavior belong to Chapter 3.

## Chapter 3: Sending JSON Responses

### Summary

Chapter 3 changes the API from plain-text responses to structured JSON
responses. The chapter starts with a simple fixed-format JSON string, then moves
to Go's `encoding/json` package and a reusable `writeJSON()` helper.

The important design idea is consistency. Successful responses should use the
same JSON-writing path, response bodies should have a predictable envelope, and
client-facing errors should also be JSON. Internal error details belong in logs;
clients should receive clear and stable messages.

### Topics

- Fixed-format JSON responses.
- JSON encoding.
- Encoding structs.
- Response envelopes.
- Custom JSON output.
- JSON error messages.

### Main Concepts

- Set `Content-Type: application/json` for JSON responses.
- JSON can be written as plain text, but `json.Marshal()` is safer and more
  reusable for dynamic responses.
- Encode the response before writing headers, so encoding errors can still be
  handled cleanly.
- Centralize JSON response writing in a helper that accepts status, data, and
  optional headers.
- Add a trailing newline to JSON responses for easier terminal reading.
- Exported struct fields are required for JSON encoding.
- Use struct tags to define the public API shape instead of exposing Go field
  names directly.
- Use `json:"-"` for fields that should never be exposed.
- Use `omitzero` when zero values should be left out of the response.
- Use `omitempty` mainly when empty slices or maps should be omitted.
- Wrap response data in envelopes such as `movie`, `error`, or `system_info`.
- Envelopes make responses consistent and leave room for metadata later.
- Implement `MarshalJSON()` on small domain types when the default JSON shape is
  not the API shape that clients should see.
- Keep detailed server errors in logs and return generic JSON messages for
  unexpected failures.
- Configure router-level `404` and `405` handlers so routing errors also use
  the API's JSON error format.
- Use panic recovery middleware to return a JSON `500` response for handler
  panics.
- Panic recovery middleware does not catch panics from background goroutines.

### API Shape Noted

Chapter 3 keeps the same early endpoints, but changes their response format:

- `GET /v1/healthcheck` returns JSON with status and system information.
- `GET /v1/movies/:id` returns a dummy movie inside a `movie` envelope.
- Not-found and method-not-allowed cases return JSON error envelopes.

### Checklist

- [x] Replace plain-text healthcheck output with JSON.
- [x] Add a reusable JSON response helper.
- [x] Use response envelopes where they improve consistency.
- [x] Create a movie response shape.
- [x] Encode structs rather than manually building JSON strings.
- [x] Add custom JSON formatting only where it has a concrete purpose.
- [x] Add reusable JSON error response helpers.
- [x] Return JSON for common API errors.
- [x] Add custom router error handlers for JSON `404` and `405` responses.
- [x] Add panic recovery middleware.
- [x] Verify compilation with `go test ./...`.
- [x] Record response behavior in the README or notes.

### Classification

- [x] Apply now: JSON helper, response envelopes, structured errors.
- [x] Apply now: consistent `Content-Type` behavior.
- [x] Apply now: custom JSON for the `Runtime` domain type.
- [x] Apply now: JSON `404`, `405`, and `500` behavior.
- [x] Postpone: broader JSON customization until a domain type needs it.
- [x] Ignore for now: cosmetic response formatting not useful for APIs.

### Personal Takeaway

This chapter introduces the first API response contract. The key pattern is not
just "return JSON", but "return JSON through one consistent path". For future Go
backend work, the reusable pieces are the `writeJSON()` helper, response
envelopes, small domain-specific JSON types, centralized error helpers, and
router-level JSON error handling.

One code-style note from the chapter: the basic `range` loop for copying headers
is clearer than using `maps.Insert()` here, even though both can work.

## Chapter 4: Parsing JSON Requests

### Summary

Chapter 4 adds the request side of JSON APIs. The main flow is: decode a JSON
body into a small input struct, convert it into a domain value, validate it, and
return clear JSON errors when the request is bad.

### Topics

- JSON decoding.
- Bad request handling.
- Input restrictions.
- Custom JSON decoding.
- Input validation.

### Main Concepts

- Use `json.Decoder` for request bodies.
- Decode into a pointer destination.
- Decode into an input struct, not directly into the domain model.
- Use a reusable `readJSON()` helper for consistent request parsing.
- Return clear `400 Bad Request` errors for malformed, empty, oversized,
  unknown-field, or multi-value JSON bodies.
- Use `UnmarshalJSON()` for domain-specific JSON input such as `"107 mins"`.
- Treat validation separately from decoding.
- Return `422 Unprocessable Entity` for business-rule validation failures.
- Keep reusable validation helpers in `internal/validator`.
- Keep movie-specific validation close to the movie domain type.

### Checklist

- [x] Add request body decoding for create movie.
- [x] Add a reusable JSON request decoding helper.
- [x] Return clear JSON errors for malformed requests.
- [x] Restrict unexpected input fields.
- [x] Handle empty or invalid request bodies.
- [x] Add custom decoding where the domain type needs it.
- [x] Add validation for title, year, runtime, and genres.
- [x] Keep validation separate enough to reuse.
- [x] Record request behavior in notes.

### Classification

- [x] Apply now: JSON decoding helper, bad request handling, validation.
- [x] Apply now: strict request input behavior.
- [x] Apply now: custom runtime decoding.
- [x] Postpone: persistence and real create responses until database chapters.
- [x] Ignore for now: accepting loose client input for convenience.

### Personal Takeaway

Chapter 4 defines the input contract. Decoding protects the API from invalid
JSON shape; validation protects the domain rules. Keeping those steps separate
makes handlers easier to reason about.

## Chapter 5: Database Setup and Configuration

### Topics

- PostgreSQL setup.
- Database connection.
- Database connection pool configuration.

### Checklist

- [ ] Decide the local PostgreSQL setup for this stage.
- [ ] Create a local database for the practice API.
- [ ] Add database DSN configuration.
- [ ] Connect to PostgreSQL from the API.
- [ ] Ping the database during startup.
- [ ] Add connection pool settings.
- [ ] Pass database dependencies through the application structure.
- [ ] Document required environment variables or flags.
- [ ] Document local database setup commands.
- [ ] Verify the API starts with a working database connection.

### Classification

- [ ] Apply now: PostgreSQL connection and pool configuration.
- [ ] Apply now: explicit local database setup notes.
- [ ] Postpone: production database setup.
- [ ] Ignore for now: non-PostgreSQL database alternatives.

## Chapter 6: SQL Migrations

### Topics

- Purpose of SQL migrations.
- Creating and applying migrations.
- Tracking schema changes.

### Checklist

- [ ] Choose the migration tool used by the book for this stage.
- [ ] Add migration files under `migrations`.
- [ ] Create the movies table migration.
- [ ] Add a rollback migration.
- [ ] Run migrations locally.
- [ ] Verify the resulting schema.
- [ ] Document migration commands.
- [ ] Add migration commands to the README or Makefile only when useful.
- [ ] Record common migration failure cases.

### Classification

- [ ] Apply now: versioned SQL migrations.
- [ ] Apply now: local migration commands.
- [ ] Postpone: deployment migration process.
- [ ] Ignore for now: ORM-based schema management.

## Chapter 7: CRUD Operations

### Topics

- Movie model setup.
- Creating a movie.
- Fetching a movie.
- Updating a movie.
- Deleting a movie.

### Checklist

- [ ] Add the movie model.
- [ ] Add a model layer for database access.
- [ ] Add create movie database logic.
- [ ] Connect `POST /v1/movies` to database persistence.
- [ ] Add fetch movie by ID logic.
- [ ] Connect `GET /v1/movies/{id}` to database reads.
- [ ] Add update movie logic.
- [ ] Connect movie update route to the model layer.
- [ ] Add delete movie logic.
- [ ] Connect movie delete route to the model layer.
- [ ] Map not-found database cases to API responses.
- [ ] Map validation/database errors to clear API responses.
- [ ] Verify CRUD behavior with `curl`.
- [ ] Add focused tests if the implementation shape makes them practical.

### Classification

- [ ] Apply now: model layer, CRUD handlers, error mapping.
- [ ] Apply now: separation between handlers and database logic.
- [ ] Postpone: generic repository abstractions unless duplication proves it is needed.
- [ ] Ignore for now: implementing unrelated resources.

## Chapter 8: Advanced CRUD Operations

### Topics

- Partial updates.
- Optimistic concurrency control.
- SQL query timeouts.

### Checklist

- [ ] Implement partial update behavior with `PATCH`.
- [ ] Preserve existing values when request fields are omitted.
- [ ] Validate partial update inputs.
- [ ] Add version-based optimistic concurrency control.
- [ ] Return a clear conflict response when concurrent updates collide.
- [ ] Add query context timeouts for database operations.
- [ ] Make timeout behavior visible in code without overcomplicating handlers.
- [ ] Verify partial updates with `curl`.
- [ ] Verify stale update behavior if practical.
- [ ] Record where this pattern may help future client API work.

### Classification

- [ ] Apply now: PATCH semantics and partial updates.
- [ ] Apply now: optimistic concurrency control.
- [ ] Apply now: SQL query timeouts.
- [ ] Postpone: broad transaction patterns until needed.
- [ ] Ignore for now: over-engineered generic update systems.

## Chapter 9: Filtering, Sorting, and Pagination

### Topics

- Query string parameter parsing.
- Query parameter validation.
- Listing data.
- Filtering lists.
- Full-text search.
- Sorting lists.
- Paginating lists.
- Pagination metadata.

### Checklist

- [ ] Add a query string parsing helper.
- [ ] Add validation for filters.
- [ ] Add a movie list endpoint.
- [ ] Implement title/genre filtering.
- [ ] Implement full-text search if it fits the book implementation.
- [ ] Add explicit allowed sort fields.
- [ ] Reject unsupported sort values.
- [ ] Add page and page size parameters.
- [ ] Add limit and offset behavior in SQL.
- [ ] Return pagination metadata.
- [ ] Verify list behavior with `curl`.
- [ ] Record examples for filtering, sorting, and pagination.
- [ ] Note reusable patterns for future catalog/read API work.

### Classification

- [ ] Apply now: validated filtering, sorting, pagination.
- [ ] Apply now: response metadata for list endpoints.
- [ ] Postpone: complex search ranking beyond the book scope.
- [ ] Ignore for now: unbounded list endpoints.

## Cross-Chapter Deliverables

- [ ] `books/lets-go-further/README.md` explains purpose, run commands, test commands, and learned patterns.
- [ ] Notes capture chapter summaries in my own words.
- [ ] API can be run locally.
- [ ] Healthcheck endpoint works.
- [ ] Movie CRUD endpoints work.
- [ ] Filtering, sorting, and pagination endpoint works.
- [ ] Database migrations are documented and repeatable.
- [ ] JSON request and response behavior is consistent.
- [ ] Validation failures return useful JSON responses.
- [ ] Any skipped topic is explicitly marked as postponed or ignored.

## Suggested Study Sessions

- [x] Session 1: Chapter 1 and local prerequisites.
- [x] Session 2: Chapter 2 project skeleton and healthcheck.
- [x] Session 3: Chapter 2 routing and route parameter helper.
- [x] Session 4: Chapter 3 JSON response helpers.
- [x] Session 5: Chapter 3 JSON error responses.
- [x] Session 6: Chapter 4 JSON request decoding.
- [x] Session 7: Chapter 4 validation.
- [ ] Session 8: Chapter 5 PostgreSQL connection.
- [ ] Session 9: Chapter 6 migrations.
- [ ] Session 10: Chapter 7 create and fetch movie.
- [ ] Session 11: Chapter 7 update and delete movie.
- [ ] Session 12: Chapter 8 partial updates.
- [ ] Session 13: Chapter 8 concurrency control and timeouts.
- [ ] Session 14: Chapter 9 list endpoint and filters.
- [ ] Session 15: Chapter 9 sorting, pagination, and metadata.
- [ ] Session 16: README cleanup and stage review.

## Definition of Done

- [ ] First 9 chapters have been read.
- [ ] Each chapter has a short summary.
- [ ] Each chapter has at least one practical task or explicit reason to skip.
- [ ] The practice API implements the API Core behavior from chapters 1-9.
- [ ] Run and test commands are documented.
- [ ] Useful patterns are mapped to future Go backend API work.
- [ ] Postponed topics are clearly separated from current work.
- [ ] No chapters beyond 9 are started as part of this plan.
