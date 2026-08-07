# Testing Lab

## Purpose

This folder documents the testing method used across the Go web labs. It does
not contain a second API or a shared test package: tests live beside the code
whose behavior they protect.

Stage 5 focuses on choosing the closest useful test boundary. A passing test
count is not a goal by itself.

## Test boundaries

| Boundary | Question it answers | Preferred tools | Current examples |
|---|---|---|---|
| Unit / service | Does a business rule or domain error mapping work? | Standard `testing`, table-driven cases, small manual fakes | `labs/layered-api/internal/books/service_test.go` |
| Adapter / repository | Does a concrete in-memory adapter preserve its local behavior? | Standard `testing`; concurrency checks when the adapter owns shared state | `labs/layered-api/internal/books/memory_repository_test.go` |
| Handler / router | Does an HTTP request produce the promised status, headers, and JSON? | `httptest.NewRequest`, `httptest.NewRecorder`, fake services | `labs/rest-api/cmd/api/handlers_test.go`, `labs/layered-api/internal/http/handlers/books_test.go` |
| Route integration | Are routing, middleware, handler, and real in-memory dependencies wired together correctly? | In-process `httptest` against the router | `labs/layered-api/internal/http/router/router_test.go` |
| Database integration | Do migrations, SQL, driver behavior, and database constraints agree? | Dedicated disposable database; explicit setup and cleanup | Planned for API Core only; Snippetbox has an existing MySQL example |

Choose the lowest level that can expose the risk. Do not repeat a happy path at
every layer merely to increase coverage.

## Test commands

Run each command from the named module directory. In restricted environments,
set the writable cache shown below.

| Module | Working directory | Command | Notes |
|---|---|---|---|
| REST API | `labs/rest-api` | `GOCACHE=/tmp/go-web-labs-go-cache go test ./...` | In-memory route-level tests; no external services. |
| Layered API | `labs/layered-api` | `GOCACHE=/tmp/go-web-labs-go-cache go test ./...` | Unit, handler, router, middleware, and composition checks; no external services. |
| Snippetbox | `books/lets-go` | `GOCACHE=/tmp/go-web-labs-go-cache go test ./...` | Includes `httptest.NewTLSServer` and a MySQL-backed model test. It needs a reachable local test MySQL instance. |
| API Core | `books/lets-go-further` | `GOCACHE=/tmp/go-web-labs-go-cache go test ./...` | Current baseline has no tests. Future ordinary tests must not require PostgreSQL. |

Run `GOCACHE=/tmp/go-web-labs-go-cache go vet ./...` from an affected module
after a code change. For documentation-only changes, run `git diff --check`.

## Deterministic fixtures and cleanup

- Prefer fixed input values and explicit expected output; do not depend on test
  execution order, current time, random data, or a shared process state.
- Construct a fresh in-memory store, fake, handler, or router for each test
  unless the test intentionally verifies a single stateful flow.
- A fake records only the call inputs that the assertion needs. Configure its
  returned values and errors directly in the test.
- Keep database fixtures isolated. Setup must name the exact disposable target,
  and cleanup must remove only data created by that test.
- Never connect to, create, drop, migrate, or clean a production database.

The existing Snippetbox model test uses a fixed local MySQL test database and
SQL setup/teardown files. It is a historical example, not the template for the
API Core slice. The API Core integration test will instead require an explicit
`GREENLIGHT_TEST_DB_DSN`; when it is unset, the test must skip with a clear
message.

## Practical choices

### Table-driven tests

Use a table when one behavior has several meaningful inputs with the same
assertion shape: validation boundaries, duplicate cases, invalid IDs, or error
mapping. Give each case a descriptive name and create independent state inside
each subtest.

### Small fakes

Use a manual fake to test a service or handler without a database, router, or
unrelated package. This keeps the test focused on observable behavior. Avoid a
mocking framework unless interaction verification is genuinely necessary.

### `httptest.NewRecorder`

Use a recorder with `httptest.NewRequest` for a handler or in-process router.
Assert client-facing contract details: status, JSON envelope and error code,
`Content-Type`, `Location`, `Allow`, or an empty `204` body where applicable.
Unexpected errors must return safe client messages and never leak internal
details.

### `httptest.NewServer`

Use `httptest.NewServer` or `httptest.NewTLSServer` only when a real HTTP client
or server behavior is part of the risk. Most handler and router contracts are
clearer and faster with a recorder. Such tests may require a runtime that
permits local listening sockets.

## Stage 5 decisions

- Standard-library assertions and manual fakes are the default.
- `testify/assert` and `testify/require` may be evaluated in one focused test
  later in this stage; do not rewrite the suite or introduce `testify/suite`
  for consistency.
- The planned PostgreSQL integration test belongs next to API Core database
  code, not in this folder, and must remain opt-in.
