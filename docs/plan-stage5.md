# Stage 5 Plan: Testing Lab

## Status

**Planned.**

Stage 5 will make the repository's testing practice explicit: add or fix a small backend behavior
and cover the relevant risk with a focused test.

`labs/testing` will hold the testing strategy, test matrix, and test-database notes. Tests
themselves will live beside the code they exercise, such as
`labs/rest-api`, `labs/layered-api`, or the API Core study project. This avoids building a third
copy of the same API solely for testing.

## Starting point

Stages 1, 3, and 4 already include focused HTTP handler, service, repository, router, middleware,
and in-memory-adapter tests. Stage 5 will audit that coverage first and close only meaningful gaps
rather than duplicate happy paths to inflate a test count.

## Scope

- [ ] Audit existing tests in the REST API, layered API, Snippetbox, and API Core, then choose 5–10
  new focused checks for uncovered risks.
- [ ] Use table-driven service tests for business rules and edge cases.
- [ ] Use small fakes or in-memory adapters for observable service behavior.
- [ ] Use `httptest` for handler and router HTTP contracts.
- [x] Add one opt-in integration slice for a real database/migration risk.
- [ ] Document commands, fixtures, test boundaries, and troubleshooting in
  `labs/testing/README.md`.
- [ ] Add safe test-database bootstrap notes for the API Core study project.

## Test-level selection

| Question                                          | Preferred test                                          |
|---------------------------------------------------|---------------------------------------------------------|
| Does a business rule or error mapping work?       | Service unit test with a small fake repository.         |
| Are HTTP status, headers, and JSON correct?       | Handler or router test using `httptest`.                |
| Do migrations, SQL, and the database model agree? | Opt-in integration test with a dedicated test database. |

The chosen test should be the closest level that can expose the risk. Avoid repeating the same happy
path at every layer.

## Work plan

### 0. Establish the baseline

- [x] Inventory current tests in the REST API, layered API, Snippetbox, and API Core.
- [x] Create a behavior → test-level → existing-test matrix.
- [x] Select 5–10 additional checks based on concrete gaps, not coverage alone.
- [x] Confirm that the planned integration test has real SQL or migration value.

#### Baseline audit — 2026-08-07

| Area                          | Behavior / risk                                                                                      | Existing test level and coverage                                     | Decision                                                                                                                                  |
|-------------------------------|------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| REST API                      | Book CRUD, validation, malformed JSON, route/method errors, and `204` responses                      | Route-level `httptest` coverage in `cmd/api/handlers_test.go`        | No happy-path duplication; add only a safe-panic-body assertion if this API is changed.                                                   |
| Layered API                   | Create, update, delete business rules and in-memory repository concurrency                           | Service unit tests with a manual fake; repository unit tests         | Add uncovered service error propagation checks.                                                                                           |
| Layered API                   | HTTP request parsing, status/error mapping, safe internal errors, and stateful router flow           | Handler tests with a fake service; small route and middleware tests  | Add invalid-ID and generic-error boundary checks where handlers currently lack them.                                                      |
| Snippetbox (`Let's Go`)       | MPA handlers, CSRF/form validation, security headers, template formatting, and user-model existence  | `httptest.NewTLSServer` handler tests plus a MySQL-backed model test | Existing coverage is sufficient for this stage; document its current integration-test boundary and avoid adding another database pattern. |
| API Core (`Let's Go Further`) | Movie validation, HTTP helpers/handlers, PostgreSQL CRUD, migrations, arrays, and optimistic locking | No tests yet                                                         | Primary source of Stage 5 gaps: add focused unit/HTTP tests and one opt-in PostgreSQL integration slice.                                  |

Selected additional checks (nine total):

1. Layered API service: `List` propagates a repository error.
2. Layered API service: `Get` propagates an unexpected repository error.
3. Layered API handler: `Get` converts an unexpected service error to a safe `500` JSON response.
4. Layered API handler: `Update` and `Delete` reject an invalid route ID without calling the
   service.
5. API Core validator: table-driven boundaries for a valid movie and invalid
   required/year/runtime/genre values.
6. API Core HTTP helper: malformed, unknown-field, empty, and multi-value JSON input produces the
   documented bad-request path.
7. API Core handler: validation failure returns `422` without calling the database model.
8. API Core integration: migrations plus `MovieModel.Insert` and `Get` preserve generated fields and
   PostgreSQL genre arrays.
9. API Core integration: a stale-version `MovieModel.Update` maps to `ErrEditConflict`.

The API Core integration candidates have real SQL and migration value: they exercise PostgreSQL
arrays, database-generated fields, check constraints, and optimistic-locking semantics that a fake
cannot reproduce. The initial implementation will use a dedicated `GREENLIGHT_TEST_DB_DSN`, skip
when it is absent, and operate only on an explicitly disposable database.

#### External reference audit — 2026-08-07

The following projects are read-only references for testing patterns. Their code, test suites, and
test gaps are outside this stage's implementation scope and do not add to the selected nine checks.

| Reference project    | Reusable observation                                                                                                                                       | Stage 5 use                                                                                                                                       |
|----------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| `book-social`        | Manual repository fakes support service tests; `httptest` covers handler rendering and application routes; SQLite integration tests use dedicated helpers. | Keep fakes small and use a real database only for persistence behavior that cannot be reproduced otherwise.                                       |
| `go-service-starter` | Table-driven handler tests assert status, JSON, and `Content-Type`; route registration and in-memory store behavior have distinct tests.                   | Keep HTTP-contract, router, and in-memory behavior tests at separate, narrow levels.                                                              |
| `mslab-wire-nats`    | Event-envelope JSON round trips are tested as contracts; health/readiness handlers are tested without starting services or infrastructure.                 | Treat message contracts as their own boundary in future event-driven work; this is not needed for the current API-focused Stage 5 implementation. |

Do not modify these projects, run their test suites, or include their uncovered behavior in this
repository's test-count target without a separate task.

### 1. Document the testing strategy

- [x] Create `labs/testing/README.md`.
- [x] Describe unit/service, handler/router, and integration/database boundaries.
- [x] Document module-specific test commands and their working directories.
- [x] Define deterministic fixture and cleanup rules.
- [x] Explain when to use table-driven tests, fakes, `httptest.NewRecorder`, and
  `httptest.NewServer`.
- [x] Add a concrete test matrix and reusable testing review checklist.

### 2. Close unit-test gaps

- [x] Add table rows only for uncovered validation, duplicate, missing-record, state, or
  repository-error behavior.
- [x] Keep fakes small, configurable, and readable; capture only inputs that a test must observe.
- [x] Keep service tests independent of JSON, Chi, a live server, and external state.

### 3. Close HTTP-contract gaps

- [x] Add `httptest` cases for route/method behavior, status, JSON envelope,
  `Content-Type`, `Location`, error codes, and empty `204` responses where needed.
- [x] Use a fake service when testing HTTP-to-application mapping.
- [x] Add one in-process route-level test only when it protects an untested
  router/middleware/handler integration point.
- [x] Verify that safe `500` responses never expose internal details.

### 4. Add an opt-in database integration slice

- [x] Choose one database-backed API Core operation that cannot be reliably tested with a fake.
- [x] Put the integration test next to that database code.
- [x] Require a dedicated `GREENLIGHT_TEST_DB_DSN`; without it, the test must skip with a clear
  message.
- [x] Document manual creation and migration of a disposable, non-production database.
- [x] Isolate or clean test data predictably, with the exact cleanup target known before running
  tests.
- [x] Ensure ordinary `go test ./...` does not require a database.

#### Deferred database follow-up

These High findings may be completed after Step 5, but must remain visible before the final Stage 5
handoff:

- [ ] Make the migration setup repeatable for reruns against the same disposable database, or
  document and verify an exact fresh-database lifecycle.
- [ ] Confirm the concrete test database name, owner, purpose, and disposable non-production status;
  then run the targeted PostgreSQL integration test and record runtime evidence separately from a
  skipped result.

### 5. Explore Testify with a real test case

- [ ] Compare standard-library assertions (`if`/`t.Fatalf`) with
  `testify/assert` and `testify/require` in one small existing or new table-driven test.
- [ ] Evaluate readability, failure output, `require`'s early exit, and whether the additional
  dependency helps that specific case.
- [ ] Read about `testify/mock`, but retain the readable manual fakes from Stages 3–4 unless
  interaction verification is genuinely required.
- [ ] Do not rewrite the existing suite or introduce `testify/suite` merely for consistency.
- [ ] Record the decision in `labs/testing/README.md`: stdlib and fakes remain the default;
  `assert`/`require` are optional readability tools.

### 6. Verify and hand off

- [ ] Update only README files whose actual commands or behavior changed.
- [ ] Run the narrowest affected tests first, then `go test ./...` and
  `go vet ./...` for each affected module.
- [ ] Run the integration command only against an explicitly configured disposable database; report
  a skipped integration test accurately.
- [ ] Run `git diff --check` and record pre-existing failures separately.

## Out of scope

- [ ] Do not rewrite existing tests solely for stylistic consistency.
- [ ] Do not create a duplicate API or second domain in `labs/testing`.
- [ ] Do not add a broad assertion or mocking framework without a concrete need.
- [ ] Do not add browser E2E, load, fuzz, mutation, CI, Docker, or Testcontainers work in this
  stage. Revisit Testcontainers only as backlog no earlier than Stages 9–10, when a real integration
  or local-infrastructure need exists.
- [ ] Do not connect to, create, drop, or clean a production database.

## Definition of done

- [ ] The selected 5–10 new tests close documented behavior gaps; PostgreSQL checks remain
  runtime-unverified until they run against a confirmed disposable database.
- [x] Unit, HTTP-contract, and opt-in integration boundaries are clear in the testing README.
- [ ] Commands are reproducible without hidden infrastructure; the opt-in migration setup still
  needs a repeatable rerun strategy.
- [x] The integration test is safe to skip without a dedicated test database.
- [ ] No duplicate application, unsafe database cleanup, or unnecessary testing framework was
  introduced.

## Handoff

The reusable outcome is a method, not a shared test package: choose the test boundary before
implementing a feature, use deterministic fixtures, and assert safe client-facing errors. Future
authentication tests wait for a real authentication use case; Stage 6 can reuse stable test examples
when documenting a real API contract.
