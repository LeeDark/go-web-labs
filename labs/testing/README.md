# Testing Lab

## Purpose

This folder documents the testing method used across the Go web labs. It does not contain a second
API or a shared test package: tests live beside the code whose behavior they protect.

Stage 5 focuses on choosing the closest useful test boundary. A passing test count is not a goal by
itself.

Related review artifacts:

- [`test-matrix.md`](test-matrix.md) maps the selected risks to concrete tests and evidence status.
- [`review-checklist.md`](review-checklist.md) provides a reusable boundary and safety checklist.

## Test boundaries

| Boundary             | Question it answers                                                                         | Preferred tools                                                           | Current examples                                                                                  |
|----------------------|---------------------------------------------------------------------------------------------|---------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------|
| Unit / service       | Does a business rule or domain error mapping work?                                          | Standard `testing`, table-driven cases, small manual fakes                | `labs/layered-api/internal/books/service_test.go`                                                 |
| Adapter / repository | Does a concrete in-memory adapter preserve its local behavior?                              | Standard `testing`; concurrency checks when the adapter owns shared state | `labs/layered-api/internal/books/memory_repository_test.go`                                       |
| Handler / router     | Does an HTTP request produce the promised status, headers, and JSON?                        | `httptest.NewRequest`, `httptest.NewRecorder`, fake services              | `labs/rest-api/cmd/api/handlers_test.go`, `labs/layered-api/internal/http/handlers/books_test.go` |
| Route integration    | Are routing, middleware, handler, and real in-memory dependencies wired together correctly? | In-process `httptest` against the router                                  | `labs/layered-api/internal/http/router/router_test.go`                                            |
| Database integration | Do migrations, SQL, driver behavior, and database constraints agree?                        | Dedicated disposable database; explicit setup and cleanup                 | Opt-in API Core PostgreSQL test; Snippetbox has an existing MySQL example                         |

Choose the lowest level that can expose the risk. Do not repeat a happy path at every layer merely
to increase coverage.

## Test commands

Run each command from the named module directory. In restricted environments, set the writable cache
shown below.

| Module      | Working directory       | Command                                           | Notes                                                                                                           |
|-------------|-------------------------|---------------------------------------------------|-----------------------------------------------------------------------------------------------------------------|
| REST API    | `labs/rest-api`         | `GOCACHE=/tmp/go-web-labs-go-cache go test ./...` | In-memory route-level tests; no external services.                                                              |
| Layered API | `labs/layered-api`      | `GOCACHE=/tmp/go-web-labs-go-cache go test ./...` | Unit, handler, router, middleware, and composition checks; no external services.                                |
| Snippetbox  | `books/lets-go`         | `GOCACHE=/tmp/go-web-labs-go-cache go test ./...` | Includes `httptest.NewTLSServer` and a MySQL-backed model test. It needs a reachable local test MySQL instance. |
| API Core    | `books/lets-go-further` | `GOCACHE=/tmp/go-web-labs-go-cache go test ./...` | Ordinary tests need no PostgreSQL; the opt-in integration test skips unless `GREENLIGHT_TEST_DB_DSN` is set.    |

Run `GOCACHE=/tmp/go-web-labs-go-cache go vet ./...` from an affected module after a code change.
For documentation-only changes, run `git diff --check`.

## Deterministic fixtures and cleanup

- Prefer fixed input values and explicit expected output; do not depend on test execution order,
  current time, random data, or a shared process state.
- Construct a fresh in-memory store, fake, handler, or router for each test unless the test
  intentionally verifies a single stateful flow.
- A fake records only the call inputs that the assertion needs. Configure its returned values and
  errors directly in the test.
- Keep database fixtures isolated. Setup must name the exact disposable target, and cleanup must
  remove only data created by that test.
- Never connect to, create, drop, migrate, or clean a production database.

The existing Snippetbox model test uses a fixed local MySQL test database and SQL setup/teardown
files. It is a historical example, not the template for the API Core slice. The API Core integration
test will instead require an explicit
`GREENLIGHT_TEST_DB_DSN`; when it is unset, the test must skip with a clear message.

### API Core PostgreSQL integration

The API Core integration test is opt-in. Provision a dedicated disposable database through the
developer's normal local PostgreSQL setup, then explicitly provide its name and DSN. The repository
does not create roles or databases during ordinary test runs, and it never drops a database.

Prerequisites:

- PostgreSQL and its client utilities (`createdb`) are installed locally.
- You have a local PostgreSQL role that can create a database, or a local administrator can do that
  one-time setup for you.
- The database is non-production, disposable, and its name ends in `_test`.

For a typical local Linux installation using peer authentication, an administrator can provision a
role matching the current OS user once:

```bash
sudo -u postgres createuser --createdb "$USER"
```

Other PostgreSQL installations may use a different administrator account, host, port, or credential
method. That local setup is outside the repository; do not add credentials to the repository.

Create the disposable database and run the test from the API Core module:

```bash
export GREENLIGHT_TEST_DB_NAME='greenlight_test'
export GREENLIGHT_TEST_DB_ROLE="$USER"
createdb --owner="$GREENLIGHT_TEST_DB_ROLE" "$GREENLIGHT_TEST_DB_NAME"
export GREENLIGHT_TEST_DB_DSN="host=/var/run/postgresql dbname=$GREENLIGHT_TEST_DB_NAME user=$GREENLIGHT_TEST_DB_ROLE sslmode=disable"
cd books/lets-go-further
go test ./greenlight/internal/data \
  -run '^TestMovieModelInsertAndGetWithPostgres$' -count=1 -v
# Run the same command a second time to verify repeatable migrations.
go test ./greenlight/internal/data \
  -run '^TestMovieModelInsertAndGetWithPostgres$' -count=1 -v
```

The example uses a Linux local PostgreSQL socket and peer authentication. Set
`GREENLIGHT_TEST_DB_ROLE` to the PostgreSQL role your environment provides, and adapt `host`,
`port`, and SSL settings in the DSN when needed. If the database already exists, reuse it only after
confirming that it is the same disposable test database; do not point the variables at a production
or shared database.

Before the test starts, it verifies that the DSN points to exactly
`GREENLIGHT_TEST_DB_NAME`, requires the `_test` suffix, and logs database owner and connected user.
It applies each checked-in migration once, recording the filename and SHA-256 checksum in the
test-only `stage5_test_migrations` ledger. Later runs verify the checksum and skip completed
migrations; a changed migration requires a fresh disposable database.

The test verifies a database constraint, generated fields, the PostgreSQL `text[]` round trip, and
optimistic locking, then deletes its exact fixture by generated ID. Never point either variable at a
production or shared database. Without `GREENLIGHT_TEST_DB_DSN`, ordinary `go test ./...` remains
database-free and the integration test is skipped. Use standard PostgreSQL credential handling
(`.pgpass`, `PGPASSFILE`, local socket, or an environment variable); do not commit credentials.

## Practical choices

### Table-driven tests

Use a table when one behavior has several meaningful inputs with the same assertion shape:
validation boundaries, duplicate cases, invalid IDs, or error mapping. Give each case a descriptive
name and create independent state inside each subtest.

### Small fakes

Use a manual fake to test a service or handler without a database, router, or unrelated package.
This keeps the test focused on observable behavior. Avoid a mocking framework unless interaction
verification is genuinely necessary.

### `httptest.NewRecorder`

Use a recorder with `httptest.NewRequest` for a handler or in-process router. Assert client-facing
contract details: status, JSON envelope and error code,
`Content-Type`, `Location`, `Allow`, or an empty `204` body where applicable. Unexpected errors must
return safe client messages and never leak internal details.

### `httptest.NewServer`

Use `httptest.NewServer` or `httptest.NewTLSServer` only when a real HTTP client or server behavior
is part of the risk. Most handler and router contracts are clearer and faster with a recorder. Such
tests may require a runtime that permits local listening sockets.

## Stage 5 decisions

- Standard-library assertions and manual fakes are the default.
- `testify/assert` and `testify/require` may be evaluated in one focused test later in this stage;
  do not rewrite the suite or introduce `testify/suite`
  for consistency.
- The PostgreSQL integration test belongs next to API Core database code, not in this folder, and
  remains opt-in.

### Reference-derived decisions

- A test must verify the exact disposable database identity before any destructive setup or cleanup;
  a DSN alone is not sufficient authorization to reset a schema.
- Use database integration tests for SQL, driver, and migration risks that fakes cannot reproduce;
  keep business-rule and HTTP-contract checks at their narrower boundaries.
- Docker Compose, Testcontainers, and CI infrastructure are appropriate when a real multi-service or
  CI reproducibility need exists. They are not a required baseline for focused local tests.
