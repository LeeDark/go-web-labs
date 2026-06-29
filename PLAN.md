# Go Web Labs Plan

## Purpose

`go-web-labs` is a learning and reference repository for Go web development.

It should support **Offer 1: Go Backend APIs & Integrations**, especially the Basic level:

- small REST/API endpoint tasks;
- handler/service/repository structure;
- validation and error responses;
- tests;
- OpenAPI/API documentation;
- basic API security;
- simple external API integrations.

This repository is **not** a second `book-social` implementation. The split is:

```text
go-web-labs  = study, notes, small isolated labs, reusable patterns
book-social  = applied project and portfolio proof
```

## Guiding principles

1. Learn patterns in `go-web-labs`, apply them in `book-social`.
2. Prefer small, focused labs over large duplicate projects.
3. Keep Chi as the primary router/framework for now.
4. Treat Echo and Gin as later comparison topics, not current priorities.
5. Use `Let's Go Further` as a practical reference while building API features in `book-social`.
6. Each meaningful folder should have a small README with purpose, run commands, learned patterns, and TODOs.

---

# Current roadmap

## Stage 0 — Repository consolidation

**Goal:** make the repository understandable as a learning path and portfolio-support repo.

### Tasks

- Update the root `README.md`.
- Clarify the relation with `book-social`.
- Add/maintain this `PLAN.md`.
- Add folder-level READMEs where useful.
- Add TODO sections instead of trying to complete every exercise immediately.

### Output

- Clean root README.
- Updated PLAN.
- Clear repo map.

### Estimate

2–4 hours.

---

## Stage 1 — Close / refresh `Let's Go`

**Folder:** `books/lets-go`

**Goal:** close the existing `Let's Go` work as a stable study artifact.

This stage is not about rewriting Snippetbox until it becomes a monument to perfection. It is about finishing, reviewing, documenting, and moving on like a responsible adult, which is apparently still legal.

### Tasks

- Ask Codex to review the current implementation.
- Fix only important issues.
- Add or improve `books/lets-go/README.md`.
- Document learned patterns:
  - `net/http` basics;
  - routing;
  - middleware;
  - templates;
  - forms;
  - sessions/cookies;
  - database access;
  - basic testing.
- Add guided exercises as TODO, not as a blocker.

### Output

- Reviewed `Let's Go` project.
- README with run/test commands and learned patterns.
- TODO list for optional exercises.

### Estimate

3–6 hours for review + README.

Optional exercises: separate backlog, not a priority.

---

## Stage 2 — REST API basics

**Folder:** `labs/rest-api`

**Goal:** capture basic REST/API patterns that directly support Offer 1 Basic.

This can be a tiny isolated lab, but it should not grow into a second `book-social`. If a feature is becoming large, apply it in `book-social` instead.

### Topics

- HTTP routes.
- JSON requests/responses.
- Status codes.
- Error response format.
- Validation basics.
- Request context.
- Basic middleware.

### Minimal lab idea

```text
GET    /health
GET    /books
GET    /books/{id}
POST   /books
PATCH  /books/{id}
DELETE /books/{id}
```

Use in-memory storage first. SQLite can be added only if it stays small.

### Output

- Small REST API lab.
- README with curl examples.
- Notes about patterns later applied in `book-social`.

### Estimate

8–14 hours.

### Priority

Medium.

This is useful, but the stronger proof should eventually be the read-only JSON API slice in `book-social` v0.2.

---

## Stage 3 — Handler → Service → Repository

**Folder:** `labs/layered-api`

**Goal:** practice backend layering without turning the lab into architecture cosplay.

### Topics

- Handler layer.
- Service/application layer.
- Repository interface.
- In-memory repository for tests.
- Request/response DTOs vs domain models.
- Dependency injection without framework magic.

### Output

- Small layered API example.
- README explaining where each responsibility belongs.
- Notes to compare with `book-social` structure.

### Estimate

8–16 hours.

### Priority

Medium.

If `book-social` v0.2 already covers this well, this lab can be shortened to notes + tiny example.

---

## Stage 4 — Testing lab

**Folder:** `labs/testing`

**Goal:** support the client-facing promise: “I can add or fix backend behavior and cover it with tests.”

### Topics

- Table-driven tests.
- Service unit tests.
- Handler tests with `httptest`.
- Fake/in-memory repositories.
- Basic integration tests.
- Test DB bootstrap notes.

### Output

- 5–10 focused tests.
- README with test commands.
- Notes about testing patterns applied in `book-social`.

### Estimate

8–16 hours.

### Priority

High after `Let's Go` closure.

---

## Stage 5 — `Let's Go Further` applied reading

**Folder:** `books/lets-go-further`

**Goal:** read `Let's Go Further` in a way that directly helps `book-social` and Offer 1 Basic.

Do not read it as an isolated pilgrimage. Use it as a practical reference while implementing real tasks. Books are useful. Worship is extra and usually billable to nobody.

## Pass 1 — API Core

### Topics

- Project structure.
- JSON responses.
- JSON requests.
- Validation.
- SQL migrations.
- CRUD.
- Filtering, sorting, pagination.

### Applied in `book-social`

- Catalog read models.
- `/api/books`.
- `/api/books/{slug}`.
- `/api/authors/{slug}`.
- `/api/genres/{slug}`.
- Error responses.
- API tests.

### Estimate

25–45 hours.

## Pass 2 — Security & Auth

### Topics

- Rate limiting.
- Authentication.
- Permission-based authorization.
- CORS.
- Secure error handling.

### Applied in `book-social`

- v0.2 auth foundation.
- Sessions/cookies.
- Login/register/logout.
- Current user middleware.
- Basic protected route.

### Estimate

20–35 hours.

## Pass 3 — Operations

### Topics

- Background tasks.
- Graceful shutdown.
- Metrics.
- Build/audit/versioning.
- Deployment basics.

### Applied in `book-social`

- HTTP middleware foundation.
- Graceful shutdown.
- Metrics later.
- Make/CI quality flow.

### Estimate

15–30 hours.

---

## Stage 6 — OpenAPI lab

**Folder:** `labs/openapi`

**Goal:** document actual API contracts, not imaginary HTML endpoints wearing an OpenAPI costume.

### Rule

OpenAPI should be added after there is a real `/api/*` slice in `book-social`.

### Topics

- `openapi.yaml`.
- Paths.
- Request schemas.
- Response schemas.
- Error response schema.
- Examples.
- Tags.

### Applied in `book-social`

- Document only JSON API endpoints:
  - `GET /api/books`;
  - `GET /api/books/{slug}`;
  - `GET /api/authors/{slug}`;
  - `GET /api/genres/{slug}`.

### Output

- Minimal OpenAPI spec.
- README with usage notes.

### Estimate

6–12 hours.

---

## Stage 7 — API security basic

**Folder:** `labs/security`

**Goal:** cover REST API + OpenAPI + API Security basics for interviews and small freelance tasks.

### Topics

- Request size limits.
- Secure headers.
- CORS.
- CSRF for MPA forms.
- Authentication basics.
- Authorization checks.
- Rate limiting.
- Safe error messages.
- Avoiding secret leaks in logs.

### Applied in `book-social`

- v0.2 HTTP middleware foundation.
- Sessions/cookies.
- Login/register/logout forms.
- CSRF TODO or minimal implementation.
- Basic protected route.

### Output

- Security checklist.
- Small examples.
- Links to `book-social` implementation notes.

### Estimate

10–20 hours.

---

## Stage 8 — Integration / External API lab

**Folder:** `labs/integrations`

**Goal:** prepare for common Upwork tasks: connect an external API, fix a webhook, add timeout/retry behavior.

### Topics

- HTTP client wrapper.
- Context timeouts.
- Safe retry basics.
- Fake external API.
- Webhook receiver.
- Basic signature verification mock.
- Tests with fake server.

### Output

- Small integration demo.
- Tests for success, timeout, and external error.
- README with notes.

### Estimate

10–18 hours.

---

# Relation to `book-social`

`book-social` is the main applied proof for Go Web skills.

## `book-social` v0.1

Finish as baseline:

- Integration / acceptance tests introduction.
- Templ spike.
- HTMX spike.
- Decision: keep Go templates or move to Templ.

## `book-social` v0.2

Use `go-web-labs` topics while implementing:

1. Quality + DB foundation.
2. SQLite/PostgreSQL + migrations.
3. Schema v0.2.
4. Catalog read models.
5. Catalog MPA update.
6. HTTP middleware + graceful shutdown.
7. Sessions/cookies.
8. User module.
9. Login/register/logout.
10. Optional read-only JSON API slice:
    - `GET /api/books`;
    - `GET /api/books/{slug}`;
    - `GET /api/authors/{slug}`;
    - `GET /api/genres/{slug}`.
11. OpenAPI only for `/api/*`.

---

# Practical priority

## Current focus

1. Close `Let's Go` with review + README.
2. Finish `book-social` v0.1.
3. Start `Let's Go Further` notes.
4. Build `book-social` v0.2 foundation.
5. Add read-only JSON API slice in `book-social` after catalog read models.
6. Add OpenAPI/security labs after real API endpoints exist.

## Not now

- Full Echo/Gin comparison.
- Large standalone REST project duplicating `book-social`.
- Complex auth in labs.
- Full API write operations before `book-social` auth and validation are ready.
- Microservices topics in this repo.

---

# Suggested repository map

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
```

---

# Offer 1 Basic proof checklist

This repository should help demonstrate:

- [ ] clean HTTP handlers;
- [ ] REST/JSON API basics;
- [ ] handler/service/repository layering;
- [ ] validation and error responses;
- [ ] unit tests;
- [ ] handler tests;
- [ ] basic integration tests;
- [ ] OpenAPI contract for real API endpoints;
- [ ] basic API security checklist;
- [ ] small external API integration patterns;
- [ ] links to applied implementation in `book-social`.
