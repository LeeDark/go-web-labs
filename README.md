# Go + Web Laboratory

Learning and reference repository for Go web development.

This repository supports the learning path for **Go Backend APIs & Integrations**. It is focused on practical Go web skills that can be reused in pet projects, portfolio work, interviews, and small freelance tasks.

The main applied project is `book-social`. This repository is the lab; `book-social` is the place where the patterns become a real application.

```text
go-web-labs  = study, notes, small isolated labs, reusable patterns
book-social  = applied project and portfolio proof
```

## Current priority

1. Close / refresh `Let's Go`.
2. Finish `book-social` v0.1 as an external checkpoint.
3. Work through `Let's Go Further` API Core.
4. Build REST API basics.
5. Practice Handler → Service → Repository structure.
6. Add testing lab or testing notes.
7. Add OpenAPI and API security notes after real API endpoints exist.
8. Add production API notes, integrations, bridge notes, and portfolio proof after the core stages.

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

Current goal: close this work as a stable study artifact with review notes, README, run/test commands, and optional exercises as TODO.

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

Use the API core parts after closing `Let's Go`. Use the production-oriented parts later for production notes and small checklists.

## Frameworks and routers

### Chi

Primary router for now.

Use Chi for practical routing, middleware, handler tests, and project structure experiments.

### Echo

Later comparison topic. Not a current priority.

### Gin

Later comparison topic. Not a current priority.

## Go Backend Roadmap

| Repo folder | Skill | Purpose |
|---|---|---|
| [books/lets-go](books/lets-go) | Go web fundamentals | Close the existing `Let's Go` work and document reusable patterns. |
| [books/lets-go-further](books/lets-go-further) | REST API and production API patterns | Use API core topics after `Let's Go`, then add production notes and checklists later. |
| [labs/rest-api](labs/rest-api) | REST API basics | Small optional API lab for routes, JSON, validation, and error responses. Do not duplicate `book-social`. |
| [labs/layered-api](labs/layered-api) | Handler → Service → Repository | Small layering example or notes that can be compared with `book-social`. |
| [labs/testing](labs/testing) | Testing basics | Unit tests, handler tests, fake repositories, and basic integration testing patterns. |
| [labs/openapi](labs/openapi) | OpenAPI | Document real JSON API endpoints after `/api/*` exists in `book-social`. |
| [labs/security](labs/security) | API security basic | Security checklist and small examples: CORS, CSRF, auth, authorization, rate limiting, safe errors. |
| [labs/integrations](labs/integrations) | External API integrations | HTTP clients, timeouts, fake external APIs, webhook basics, tests. |

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
  checklists/
AGENTS.md
repository-context.md
PLAN.md
README.md
```

## Notes

This repository should stay small and useful. If a lab starts turning into a full product, move that work to `book-social` or create a dedicated demo project.
