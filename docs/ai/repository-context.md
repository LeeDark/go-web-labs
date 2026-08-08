# Repository Context

## Name

`go-web-labs`

## Purpose

A Go Web learning and practice repository.

The repository is used to study Go web development from simple MPA applications to REST/API backend patterns.

It supports future portfolio and freelance positioning around small Go backend tasks:

- adding or fixing endpoints;
- improving handlers and middleware;
- adding tests;
- documenting APIs;
- improving local run instructions;
- preparing small integration tasks.

## Main sections

### `books/lets-go`

Study project based on *Let's Go*.

Purpose:

- Go web fundamentals;
- MPA;
- templates;
- forms;
- sessions/cookies;
- database-backed web application;
- testing basics.

Current application:

- `books/lets-go/snippetbox`

### `books/lets-go-further`

Active study project based on *Let's Go Further*.

Purpose:

- RESTful JSON APIs;
- API structure;
- validation;
- migrations;
- authentication;
- permissions;
- rate limiting;
- graceful shutdown;
- metrics;
- deployment and production topics.

### `labs/*`

Focused labs for reusable skills:

- REST API Basic;
- layered API;
- testing;
- OpenAPI;
- API security;
- integrations.

### `frameworks/*`

Small framework and tool experiments.

Existing examples include:

- Echo quick-start notes;
- Templ hello-world spike.

### Parallel learning tracks

- Learn Go with Tests supplies just-in-time TDD, dependency-design, concurrency/context, and
  refactoring practice for later stages; it is not a gate before implementation.
- PostgreSQL supplies SQL/ACID/MVCC, schema, migrations, `pgx`, transactions, pool lifecycle, and
  optimistic-concurrency practice for Stage 2 and Stage 8.
- Docker supplies image/build, networking, volume, diagnostics, and cleanup practice for Stage 8 and
  later delivery automation.

## Relation to book-social

`go-web-labs` is the learning laboratory.

`book-social` is the applied project where selected patterns are used in a real product-like codebase.

Do not duplicate the full `book-social` domain here.

Use `go-web-labs` for small exercises, notes, and reusable patterns.

## Current roadmap status

Completed:

- Stage 0: repository consolidation;
- Stage 1: `books/lets-go` closure;
- Stage 3: REST API basics;
- Stage 4: Handler → Service → Repository;
- Stage 5: testing lab.

Active:

- Stage 2: `books/lets-go-further` API Core. Chapters 1–8 are complete; Chapter 9 and the stage
  review remain.

Queued:

- Stage 6: OpenAPI, after a real `/api/*` slice exists in `book-social`;
- Stage 7: API/MPA security, after a stable contract and authentication use case exist.

Later stages include PostgreSQL and Docker foundations, production API topics, external integrations,
test automation and delivery, a bridge to `go-microservices-starter`, and an Offer 1 portfolio
package. Learn Go with Tests is a parallel just-in-time learning track, not a stage gate.
