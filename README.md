# Go + Web Laboratory

Learning and reference repository for Go web development.

This repository supports the learning path for **Go Backend APIs & Integrations**. It is focused on
practical Go web skills that can be reused in pet projects, portfolio work, interviews, and small
freelance tasks.

The main applied project is `book-social`. This repository is the lab; `book-social` is the place
where the patterns become a real application.

```text
go-web-labs  = study, notes, small isolated labs, reusable patterns
book-social  = applied project and portfolio proof
```

## Current priority

Finished:

1. Repository consolidation.
2. Close / refresh `Let's Go`.
3. REST API basics.
4. Handler → Service → Repository.
5. Testing lab.

Current paired work:

1. `book-social` v0.2.4 HTTP Foundation with Stage 9A follow-up from `Let's Go Further`.

Next paired work:

1. `book-social` v0.2.5–v0.2.6 authentication with MPA security basics and Learn Go with Tests TDD
   foundations.

Deferred by product trigger:

1. `Let's Go Further` API Core Chapter 9 and review resume with `book-social` v0.4 catalog
   filtering, sorting, and pagination.

Later:

1. PostgreSQL and Docker foundations.
2. Production API topics, external integrations, test automation, bridge notes, and portfolio proof.

### Parallel learning tracks

- **Learn Go with Tests:** start TDD foundations with `book-social` v0.2.5–v0.2.6, then take later
  units only when their task needs them; it is not a course-completion gate.
- **PostgreSQL:** SQL/ACID/MVCC, schema, migrations, transactions, pool lifecycle, and optimistic
  concurrency, tied to API Core and the PostgreSQL/Docker foundation stage.
- **Docker:** images, multi-stage Go builds, build context, networking, volumes, diagnostics, and
  safe cleanup, tied to containerized API work and later test automation.

## Books

### Let's Go

[books/lets-go](books/lets-go/README.md)

Study project based on *Let's Go* by Alex Edwards.

Focus:

- `net/http` fundamentals;
- routing;
- middleware;
- HTML templates;
- forms;
- sessions and cookies;
- database-backed web applications;
- basic testing.

Status: completed as a stable study artifact; optional exercises remain TODO backlog.

### Let's Go Further

[books/lets-go-further](books/lets-go-further)

Study project based on *Let's Go Further* by Alex Edwards.

Focus:

- RESTful JSON APIs;
- project structure;
- JSON requests and responses;
- validation;
- SQL migrations;
- CRUD;
- filtering, sorting, and pagination;
- authentication and authorization;
- CORS;
- rate limiting;
- graceful shutdown;
- metrics;
- build, audit, versioning, and deployment basics.

Status: API Core Chapters 1–8 are complete. Chapter 9 and the API Core review are deferred until
`book-social` v0.4 needs catalog filtering, sorting, and pagination. Chapter 11 graceful shutdown
is complete as the first Stage 9A study unit alongside v0.2.4; remaining production topics stay
product-triggered.

## Frameworks and routers

### Chi

Primary router for now.

Use Chi for practical routing, middleware, handler tests, and project structure experiments.

### Echo

Later comparison topic. Not a current priority.

### Gin

Later comparison topic. Not a current priority.

## Go Backend Roadmap

| Repo folder                                        | Skill                                | Purpose                                                                                             |
|----------------------------------------------------|--------------------------------------|-----------------------------------------------------------------------------------------------------|
| [books/lets-go](books/lets-go)                     | Go web fundamentals                  | Completed study artifact with reusable MPA and testing patterns.                                    |
| [books/lets-go-further](books/lets-go-further)     | REST API and production API patterns | API Core active; Chapters 1–8 complete, with production topics later.                               |
| [labs/rest-api](labs/rest-api)                     | REST API basics                      | Completed focused lab for routes, JSON, validation, and error responses.                            |
| [labs/layered-api](labs/layered-api)               | Handler → Service → Repository       | Completed layering example to compare with `book-social`.                                           |
| [labs/testing](labs/testing)                       | Testing basics                       | Completed strategy, unit/HTTP patterns, and opt-in PostgreSQL integration guidance.                 |
| `books/lets-go-further`, `docs/`, focused API code | PostgreSQL and Docker foundations    | Stage 8: migrations, `pgx`/pool behavior, a Go image, local networking, volumes, and safe cleanup.  |
| `labs/testing`, `docs/checklists/`                 | Test automation and delivery         | Stage 11: justified CI, deeper tests, and Compose/Testcontainers decisions.                         |
| [labs/openapi](labs/openapi)                       | OpenAPI                              | Document real JSON API endpoints after `/api/*` exists in `book-social`.                            |
| [labs/security](labs/security)                     | API security basic                   | Security checklist and small examples: CORS, CSRF, auth, authorization, rate limiting, safe errors. |
| [labs/integrations](labs/integrations)             | External API integrations            | HTTP clients, timeouts, fake external APIs, webhook basics, tests.                                  |

## Relation to `book-social`

`book-social` is the main applied proof for this learning path.

Important patterns from this repository should eventually be applied there:

- handler/service/repository structure;
- database migrations;
- catalog read models;
- tests;
- middleware;
- sessions/cookies;
- authentication;
- read-only JSON API endpoints;
- OpenAPI for `/api/*`;
- API security notes.

## Offer 1 Basic mapping

This repository supports the Basic freelance offer:

> Small Go backend fix, REST API endpoint, or integration task.

The learning path should demonstrate:

- adding or fixing REST endpoints;
- improving handler/service logic;
- adding validation;
- improving error responses;
- adding focused tests;
- documenting API contracts;
- handling simple integrations;
- understanding basic API security risks.

## Portfolio evidence and handoff

- [Offer 1 Basic proof](docs/offer-1-basic-proof.md)
- [Offer 1 evidence matrix](docs/portfolio/offer-1-evidence-matrix.md)
- [Links to applied `book-social` work](docs/book-social-links.md)
- [Backend endpoint task checklist](docs/checklists/backend-endpoint-task.md)
- [API review checklist](docs/checklists/api-review.md)
- [Backend handoff notes checklist](docs/checklists/handoff-notes.md)

## Suggested repository structure

```text
books/
  lets-go/
  lets-go-further/
labs/
  rest-api/
  layered-api/
  testing/
  openapi/
  security/
  integrations/
docs/
  book-social-links.md
  offer-1-basic-proof.md
  portfolio/
    offer-1-evidence-matrix.md
  checklists/
    backend-endpoint-task.md
    api-review.md
    handoff-notes.md
AGENTS.md
PLAN.md
README.md
docs/ai/repository-context.md
docs/ai/task-history.md
```

## Notes

This repository should stay small and useful. If a lab starts turning into a full product, move that
work to `book-social` or create a dedicated demo project.
