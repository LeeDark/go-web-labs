# Testing Review Checklist

Use this checklist when adding or reviewing a focused Go web test. Apply only the sections relevant
to the selected boundary.

## Select the boundary

- [ ] State the behavior or regression risk in one sentence.
- [ ] Search for an existing test that already protects the behavior.
- [ ] Choose the lowest boundary that can expose the risk: unit/service, handler/router, or database.
- [ ] Avoid duplicating the same happy path at multiple layers.
- [ ] Keep the test beside the code it exercises.

## Unit and service tests

- [ ] Use table rows only when cases share the same setup and assertion shape.
- [ ] Keep fakes small and configure returned values or errors directly.
- [ ] Capture only inputs and call counts needed for observable assertions.
- [ ] Use `errors.Is` when an error may be wrapped.
- [ ] Verify rejected commands do not call a dependency or mutate state when that is part of the
  contract.
- [ ] Keep JSON, routers, listeners, databases, time, and network out of a unit test unless they are
  the behavior under test.

## Handler and router tests

- [ ] Use `httptest.NewRequest` and `httptest.NewRecorder` for in-process HTTP contracts.
- [ ] Assert the exact status and relevant headers: `Content-Type`, `Location`, `Allow`, or an empty
  `204` body.
- [ ] Decode JSON and assert its envelope and fields instead of relying only on substrings.
- [ ] Use a fake service for HTTP-to-application mapping and record call counts when asserting that
  no call occurred.
- [ ] Verify unexpected errors return a safe `500` without internal, database, or credential detail.
- [ ] Add a route-level test only for a real router/middleware/handler integration risk.

## Database integration tests

- [ ] Record the exact database name, owner, purpose, and disposable status before running.
- [ ] Require a dedicated test DSN and skip clearly when it is absent.
- [ ] Require the configured database name explicitly, verify it after connecting, and enforce a
  project test-name convention such as `_test`.
- [ ] Keep any database-creation helper separate from test execution: require an admin DSN, prompt
  only for non-secret metadata, refuse overwrite, and do not provide automatic deletion by default.
- [ ] Never use production credentials or a shared database.
- [ ] Make migration setup repeatable or document the required fresh-database lifecycle exactly.
- [ ] Record applied migration names and checksums when the test owns migration setup.
- [ ] Register cleanup immediately after creating data and target only known test-owned rows or
  objects.
- [ ] Exercise behavior a fake cannot reproduce: SQL mapping, arrays, constraints, generated fields,
  or optimistic locking.
- [ ] Model stale writes with two copies of the same persisted version and a successful intervening
  update.
- [ ] Report a skipped test as skipped, not as passed integration evidence.

## Verification and handoff

- [ ] Run the narrowest affected package first.
- [ ] Run `go test ./...` and `go vet ./...` from each affected module.
- [ ] Run `git diff --check`.
- [ ] Separate pre-existing failures, environment limitations, and skipped integration checks.
- [ ] Update the public and private Stage plan evidence without overstating completion.
- [ ] Keep implementation, documentation, and commits focused on one test boundary.
