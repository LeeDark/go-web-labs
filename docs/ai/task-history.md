# Task History

## Stage 1 - Close / refresh Let's Go

Status: completed with TODO exercises

Goal:

Close or refresh the `books/lets-go` study project as a clean learning artifact.

Completed work:

- reviewed the current implementation;
- ran tests where possible;
- fixed a small signup handler control-flow bug;
- created `books/lets-go/README.md`;
- added guided exercises as TODO backlog;
- created `AGENTS.md`;
- created `docs/ai/repository-context.md`.

Notes:

- Do not rewrite the project.
- Do not start REST API labs in this task.
- Do not start `Let's Go Further` implementation in this task.
- Full tests need local socket permission and MySQL access.

## Stage 2 - Let's Go Further API Core

Status: in progress (Chapters 1-8 complete)

Goal:

Work through the API Core portion of *Let's Go Further* using the Greenlight
study project.

Completed work (2026-07-06 to 2026-08-08):

- created the `books/lets-go-further` Go module and Greenlight API skeleton;
- implemented routing, JSON responses, JSON errors, panic recovery, strict
  JSON request parsing, and movie validation;
- configured PostgreSQL DSN settings and the connection pool;
- added versioned migrations for the `movies` table and its constraints;
- implemented PostgreSQL-backed create, read, update, partial update, delete, optimistic locking,
  and request timeouts;
- added focused unit, HTTP-contract, and opt-in PostgreSQL integration tests during Stage 5;
- added and updated Stage 2 learning notes and the implementation README.

Remaining work:

- Chapter 9: movie list endpoint, filters, sorting, pagination, and metadata;
- README cleanup and Stage 2 review.

## Stage 3 - REST API Basics Lab

Status: completed (Steps 0-5)

Completed work (2026-07-21 to 2026-07-26):

- created the isolated `labs/rest-api` Go module with a Chi-based JSON API and
  a concurrency-safe in-memory `books` store;
- implemented documented CRUD routes, strict JSON input, validation details,
  stable error envelopes, router errors, and recovery middleware;
- added seven focused `httptest` handler tests and complete run/test/curl
  documentation, including in-memory storage limitations;
- verified the lab with `gofmt`, `go test ./...`, `go vet ./...`, and
  `git diff --check` using a temporary Go build cache in the sandbox.

Notes:

- The sandbox does not permit opening a local listening socket, so live curl
  checks were represented by `httptest` coverage here.
- The repository root has no `Makefile`; the unrelated Snippetbox Makefile was
  not run for this Stage 3 checkpoint.

## Stage 5 - Testing lab

Status: completed (Steps 0-6 complete)

Current checkpoint (2026-08-07):

- audited existing test boundaries and selected focused gaps;
- documented the testing strategy and reproducible commands;
- added service, HTTP-contract, and opt-in PostgreSQL integration tests;
- documented the disposable API Core test database and its `GREENLIGHT_TEST_DB_DSN` guard;
- added a concrete Stage 5 test matrix and reusable review checklist;
- strengthened no-call HTTP assertions, validation boundaries, cleanup, constraints, and the
  optimistic-locking scenario after Review Mode;
- ordinary module test suites pass without PostgreSQL; the integration slice skips clearly when
  the DSN is absent;
- verified the test-only migration ledger and database contract with two successful targeted runs
  against a local disposable database; its identifiers and credentials remain private.
- evaluated `testify/assert` and `testify/require` in one focused table-driven HTTP test; standard
  library assertions and manual fakes remain the default, with no `testify/mock` or `testify/suite`.

Remaining work:

- Stage 5 handoff completed after final module verification, DB evidence, documentation, and
  focused Testify evaluation.
