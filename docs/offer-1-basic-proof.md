# Offer 1 Basic proof

## Scope

This document records the currently supported proof for **Offer 1 Basic — a small Go backend fix,
REST endpoint, or integration task**.

It is evidence for focused work, not a claim that the repositories are production-ready or that the
Advanced, Professional, or Ultimate levels are complete.

## Evidence baseline

The baseline uses immutable refs reviewed on 2026-08-09:

| Repository           | Ref                       | Role in the proof                                                      |
|----------------------|---------------------------|------------------------------------------------------------------------|
| `go-web-labs`        | `main` at `231be88`       | Study projects, isolated labs, testing guidance, and reusable patterns |
| `book-social`        | `v0.2.4` at `cba82f0`     | Applied MPA/catalog and HTTP-foundation proof                          |
| `go-service-starter` | `v0.1-basic` at `b3849d3` | Focused REST Basic Proof                                               |

## What this demonstrates

- REST resource and route design;
- JSON response and error envelopes;
- strict request-boundary and validation patterns in the study lab;
- handler/service/repository separation as a reusable design pattern;
- focused unit, HTTP-contract, and integration-test thinking;
- migration and disposable-database safety notes;
- HTTP lifecycle, middleware, recovery, headers, request IDs, and timeout decisions in the applied
  MPA project;
- clear README, verification, and handoff documentation.

The strongest current applied evidence is `book-social` v0.2.4. The strongest focused REST evidence
is `go-service-starter` v0.1-basic. They demonstrate different parts of the Basic offer and should
not be described as one identical implementation.

## Verification

The reviewed commands passed against the cited refs:

```bash
# book-social v0.2.4
GOCACHE=/tmp/book-social-go-cache go test ./...

# go-service-starter v0.1-basic
GOCACHE=/tmp/go-service-starter-basic-cache go test ./...
```

The `book-social` PostgreSQL parity tests are opt-in and require
`BOOK_SOCIAL_POSTGRES_TEST_DSN`. A normal `go test ./...` run does not claim that a live PostgreSQL
database was exercised.

## Example Basic tasks supported

- add or correct a small JSON/HTTP endpoint;
- add request validation and stable error responses;
- fix a route, status code, or response envelope;
- add focused handler or service tests;
- diagnose a small repository or migration issue;
- improve local run, test, or handoff instructions;
- review a small Go backend change for boundary and test gaps.

## Evidence links

- [Offer 1 evidence matrix](portfolio/offer-1-evidence-matrix.md)
- [Testing lab README](../labs/testing/README.md)
- [Testing matrix](../labs/testing/test-matrix.md)
- [Testing review checklist](../labs/testing/review-checklist.md)
- [Backend endpoint task checklist](checklists/backend-endpoint-task.md)
- [API review checklist](checklists/api-review.md)
- [Backend handoff notes](checklists/handoff-notes.md)
- [REST API lab](../labs/rest-api/README.md)
- [Layered API lab](../labs/layered-api/README.md)
- [Let's Go Further study README](../books/lets-go-further/README.md)

Applied-project references:

- `book-social` tag `v0.2.4`: catalog read models, migration-first database workflow, HTTP
  foundation, middleware, lifecycle, and tests;
- `go-service-starter` tag `v0.1-basic`: health endpoint, read-only item endpoint, JSON helpers,
  validation, deterministic in-memory storage, and `httptest` coverage.

## Explicit limitations

This Basic proof does not claim:

- a completed database-backed CRUD service in `go-service-starter`;
- an accepted `go-service-starter` `v0.2-advanced` milestone;
- a real `/api/*` JSON slice in `book-social`;
- OpenAPI for an applied API;
- authentication, sessions, CSRF, or authorization;
- Docker/Compose delivery as a production platform;
- Professional or Ultimate portfolio readiness.

These claims require later accepted refs and are tracked in the evidence matrix.

## Handoff checklist

Before presenting a Basic example, provide:

- the exact repository ref;
- the changed file or package path;
- the run and test commands;
- expected HTTP behavior and error cases;
- known limitations and deferred work;
- whether database or external-service access is optional or required.
