# Go Web Labs Plan

## Purpose

`go-web-labs` is a learning and reference repository for Go web development.

It supports **Offer 1: Go Backend APIs & Integrations**, especially the Basic level:

- small REST/API endpoint tasks;
- handler/service/repository structure;
- validation and error responses;
- tests;
- OpenAPI/API documentation;
- basic API security;
- simple external API integrations;
- local run and handoff notes.

This repository is **not** a second `book-social` implementation.

```text
go-web-labs  = study, notes, small isolated labs, reusable patterns
book-social  = applied project and portfolio proof
```

The goal is not to learn everything forever. That way lies browser tabs, abandoned folders, and the quiet rustle of unpaid enlightenment. The goal is to make learning produce reusable patterns and portfolio proof.

---

## Guiding principles

1. Learn patterns in `go-web-labs`, apply selected patterns in `book-social`.
2. Prefer small, focused labs over large duplicate projects.
3. Keep Chi as the primary router/framework for now.
4. Treat Echo and Gin as later comparison topics, not current priorities.
5. Use `Let's Go Further` as a practical reference while building API features in `book-social`.
6. Each meaningful folder should have a small README with:
   - purpose;
   - run commands;
   - test commands;
   - learned patterns;
   - TODOs.
7. Every larger stage should answer:
   - what to do;
   - example;
   - output;
   - offer mapping;
   - relation to `book-social`.

---

## Current roadmap

## Stage 0 — Repository consolidation

**Folder:** repository root

**Goal:** make the repository understandable as a learning path and portfolio-support repo.

### What to do

- Update the root `README.md`.
- Clarify the relation with `book-social`.
- Maintain this `PLAN.md`.
- Add folder-level READMEs where useful.
- Add TODO sections instead of trying to complete every exercise immediately.
- Add or maintain:
  - `AGENTS.md`;
  - `repository-context.md`;
  - optionally `task-history.md`.

### Example

```text
README.md
PLAN.md
AGENTS.md
repository-context.md
books/lets-go/README.md
```

### Output

- Clean root README.
- Updated PLAN.
- Clear repo map.
- AI-agent context files for future Codex work.

### Offer mapping

- Better handoff quality.
- Better client-facing documentation habits.
- Less chaos when returning to older work, a shocking innovation.

### Estimate

2–4 hours.

### Priority

High.

---

## Stage 1 — Close / refresh `Let's Go`

**Folder:** `books/lets-go`

**Goal:** close the existing `Let's Go` work as a stable study artifact.

This stage is not about rewriting Snippetbox until it becomes a monument to perfection. It is about finishing, reviewing, documenting, and moving on.

### What to do

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
- Do not start `Let's Go Further` in this task.
- Do not rewrite the app.

### Example

README structure:

```md
# Let's Go / Snippetbox

## Purpose

Study project based on "Let's Go" by Alex Edwards.

## What this project demonstrates

- net/http basics
- routing
- middleware
- HTML templates
- forms
- sessions/cookies
- database-backed web application
- authentication basics
- testing basics

## How to run

...

## How to test

...

## Useful patterns learned

...

## Remaining TODO / Guided exercises

...
```

### Output

- Reviewed `Let's Go` project.
- README with run/test commands and learned patterns.
- TODO list for optional exercises.

### Offer mapping

- HTTP fundamentals.
- MPA basics.
- Middleware basics.
- Basic testing.
- Better explanation of old code during client handoff.

### Applied in `book-social`

- Handler organization.
- Template rendering.
- Middleware chain.
- Sessions/cookies foundation.
- Integration/acceptance test thinking.

### Estimate

3–6 hours for review + README.

Optional exercises: separate backlog, not a priority.

### Priority

Current focus.

---

## Stage 2 — REST API basics

**Folder:** `labs/rest-api`

**Goal:** capture basic REST/API patterns that directly support Offer 1 Basic.

This can be a tiny isolated lab, but it should not grow into a second `book-social`. If a feature is becoming large, apply it in `book-social` instead.

### What to do

- Build a minimal API-only backend.
- Use `net/http` or Chi.
- Return JSON responses.
- Accept JSON requests.
- Use consistent error responses.
- Add basic validation.
- Add simple middleware.
- Use in-memory storage first.
- Add README with curl examples.

### Example

```text
GET    /health
GET    /books
GET    /books/{id}
POST   /books
PATCH  /books/{id}
DELETE /books/{id}
```

Possible response shape:

```json
{
  "data": {
    "id": 1,
    "title": "Example Book"
  }
}
```

Possible error shape:

```json
{
  "error": {
    "code": "validation_failed",
    "message": "Request validation failed"
  }
}
```

### Output

- Small REST API lab.
- README with curl examples.
- Notes about patterns later applied in `book-social`.

### Offer mapping

- Add endpoint.
- Fix endpoint.
- Improve handler behavior.
- Add validation.
- Improve error response format.

### Applied in `book-social`

The stronger proof should eventually be a read-only JSON API slice in `book-social` v0.2:

```text
GET /api/books
GET /api/books/{slug}
GET /api/authors/{slug}
GET /api/genres/{slug}
```

### Estimate

8–14 hours.

### Priority

Medium.

---

## Stage 3 — Handler → Service → Repository

**Folder:** `labs/layered-api`

**Goal:** practice backend layering without turning the lab into architecture cosplay.

### What to do

- Split code into handler, service, and repository layers.
- Keep handlers focused on HTTP.
- Keep services focused on use cases and business rules.
- Keep repositories focused on data access.
- Use request/response DTOs separately from domain models.
- Add in-memory repository for tests.
- Keep dependency injection boring and explicit.

### Example

```text
cmd/api/
internal/http/handlers/
internal/http/middleware/
internal/books/
  domain.go
  service.go
  repository.go
  memory_repository.go
```

Possible interface:

```go
type BookRepository interface {
    ListBooks(ctx context.Context) ([]Book, error)
    GetBook(ctx context.Context, id int64) (Book, error)
    CreateBook(ctx context.Context, input CreateBookInput) (Book, error)
}
```

### Output

- Small layered API example.
- README explaining where each responsibility belongs.
- Notes to compare with `book-social` structure.

### Offer mapping

- Fix service logic.
- Add business rule.
- Add DB-backed endpoint.
- Improve maintainability around a small feature.

### Applied in `book-social`

- Catalog service.
- Catalog repository.
- Book details read model.
- Author details read model.
- Future user/auth services.

### Estimate

8–16 hours.

### Priority

Medium.

If `book-social` v0.2 already covers this well, this lab can be shortened to notes + tiny example.

---

## Stage 4 — Testing lab

**Folder:** `labs/testing`

**Goal:** support the client-facing promise: “I can add or fix backend behavior and cover it with tests.”

### What to do

- Add table-driven tests.
- Add service unit tests.
- Add handler tests with `httptest`.
- Add fake/in-memory repositories.
- Add basic integration tests.
- Add test DB bootstrap notes.
- Document test commands.

### Example

Minimal tests:

```text
GET /books -> 200
GET /books/{existing-id} -> 200
GET /books/{missing-id} -> 404
POST /books valid -> 201
POST /books invalid -> 400
PATCH /books/{id} -> 200
DELETE /books/{id} -> 204
```

### Output

- 5–10 focused tests.
- README with test commands.
- Notes about testing patterns applied in `book-social`.

### Offer mapping

- Add tests for an existing route.
- Fix bug and cover it with tests.
- Improve QA around a small backend change.
- Make small client work safer.

### Applied in `book-social`

- v0.1 integration/acceptance tests.
- v0.2 repository tests.
- v0.2 handler tests.
- Auth flow tests later.

### Estimate

8–16 hours.

### Priority

High after `Let's Go` closure.

---

## Stage 5 — `Let's Go Further`: API Core

**Folder:** `books/lets-go-further`

**Goal:** read and implement the API core parts of `Let's Go Further` in a way that directly helps `book-social` and Offer 1 Basic.

Do not read it as an isolated pilgrimage. Use it as a practical reference while implementing real tasks. Books are useful. Worship is extra and usually billable to nobody.

### What to do

- Study API project structure.
- Implement JSON responses.
- Implement JSON requests.
- Add validation.
- Study SQL migrations.
- Practice CRUD.
- Practice filtering, sorting, and pagination.
- Keep notes about patterns worth applying in `book-social`.

### Example

Possible API endpoints from the book-style lab:

```text
GET    /v1/healthcheck
GET    /v1/movies
GET    /v1/movies/{id}
POST   /v1/movies
PATCH  /v1/movies/{id}
DELETE /v1/movies/{id}
```

Possible `book-social` equivalent:

```text
GET /api/books
GET /api/books/{slug}
GET /api/authors/{slug}
GET /api/genres/{slug}
```

### Output

- `books/lets-go-further` project or notes.
- README with learned API patterns.
- Notes linking selected patterns to `book-social`.

### Offer mapping

- Add JSON endpoint.
- Add pagination/filtering.
- Add DB query.
- Add migration.
- Fix validation.
- Improve API behavior.

### Applied in `book-social`

- Catalog read models.
- Read-only JSON API slice.
- Error responses.
- API tests.
- Migrations strategy.
- Filtering/sorting/pagination for catalog.

### Estimate

25–45 hours.

### Priority

High after `Let's Go` closure and `book-social` v0.1 baseline.

---

## Stage 6 — OpenAPI lab

**Folder:** `labs/openapi`

**Goal:** document actual API contracts, not imaginary HTML endpoints wearing an OpenAPI costume.

### Rule

OpenAPI should be added after there is a real `/api/*` slice in `book-social`.

### What to do

- Add minimal `openapi.yaml`.
- Describe paths.
- Describe request schemas.
- Describe response schemas.
- Describe common error response schema.
- Add examples.
- Add tags.
- Add README with usage notes.

### Example

Only document JSON API endpoints:

```text
GET /api/books
GET /api/books/{slug}
GET /api/authors/{slug}
GET /api/genres/{slug}
```

Possible files:

```text
labs/openapi/openapi.yaml
labs/openapi/README.md
```

### Output

- Minimal OpenAPI spec.
- README with usage notes.
- Optional notes on how this maps to `book-social`.

### Offer mapping

- Document existing API.
- Add endpoint with contract.
- Improve frontend/client handoff.
- Improve API examples for Upwork/client tasks.

### Applied in `book-social`

- OpenAPI for `/api/*` only.
- MPA endpoint docs stay in `docs/endpoints.md`.

### Estimate

6–12 hours.

### Priority

Medium, after real API endpoints exist.

---

## Stage 7 — API security basic

**Folder:** `labs/security`

**Goal:** cover REST API + OpenAPI + API Security basics for interviews and small freelance tasks.

### What to do

- Create an API security checklist.
- Add small examples where useful.
- Cover request size limits.
- Cover secure headers.
- Cover CORS basics.
- Cover CSRF for MPA forms.
- Cover authentication basics.
- Cover authorization checks.
- Cover rate limiting.
- Cover safe error messages.
- Cover logging without leaking secrets.

### Example

Minimal protected resource rule:

```text
Anonymous users can read public resources.
Authenticated users can read their own private resources.
Users cannot update resources owned by another user.
Admin users can manage all resources.
```

Possible endpoints:

```text
POST /tokens/authentication
GET  /users/me
GET  /books/{id}
PATCH /books/{id}
```

### Output

- Security checklist.
- Small examples.
- Links to `book-social` implementation notes.

### Offer mapping

- Add basic auth check.
- Fix insecure endpoint.
- Add rate limiting.
- Add CORS config.
- Review API security risks.

### Applied in `book-social`

- v0.2 HTTP middleware foundation.
- Sessions/cookies.
- Login/register/logout forms.
- CSRF TODO or minimal implementation.
- Basic protected route.
- Safe error messages.

### Estimate

10–20 hours.

### Priority

Medium.

---

## Stage 8 — `Let's Go Further`: Production API Topics

**Folder:** `books/lets-go-further`

**Goal:** use the production-oriented parts of `Let's Go Further` to support reliable API behavior and better handoff quality.

This is the missing “make it less toy-like” layer. Not enterprise. Not cloud-native ceremonial robes. Just enough production discipline to avoid embarrassing yourself in front of a server.

### What to do

- Study and document background tasks.
- Study and document graceful shutdown.
- Study and document metrics basics.
- Study and document build/audit/versioning.
- Study and document deployment basics.
- Decide what belongs in `book-social` v0.2 and what should wait.
- Add notes, not a huge framework.

### Example

Possible topics/files:

```text
books/lets-go-further/README.md
docs/checklists/graceful-shutdown.md
docs/checklists/build-audit-versioning.md
docs/checklists/metrics-basics.md
```

Possible `book-social` examples:

```text
GET /health
GET /metrics
graceful shutdown on SIGTERM
request ID middleware
request logging middleware
```

### Output

- Production API notes.
- Small checklist files.
- Clear decision about what is applied now vs later.
- Links to `book-social` tasks where relevant.

### Offer mapping

- Improve reliability of a Go backend.
- Add graceful shutdown.
- Add health endpoint.
- Add basic metrics.
- Add build/test/audit commands.
- Improve local verification and handoff notes.

### Applied in `book-social`

- v0.2 HTTP middleware foundation.
- Graceful shutdown.
- Health endpoint.
- Optional metrics later.
- Make/CI quality flow.
- Release notes and tag discipline.

### Estimate

15–30 hours.

### Priority

Medium after API Core and v0.2 foundation.

---

## Stage 9 — Integration / External API lab

**Folder:** `labs/integrations`

**Goal:** prepare for common Upwork tasks: connect an external API, fix a webhook, add timeout/retry behavior.

### What to do

- Build a small HTTP client wrapper.
- Use context timeouts.
- Add safe retry basics.
- Create a fake external API.
- Create a webhook receiver.
- Add basic signature verification mock.
- Add tests with fake server.

### Example

Possible endpoints:

```text
POST /webhooks/payment
GET  /external/books/{isbn}
POST /notifications/email
```

Possible test scenarios:

```text
external API success -> 200
external API timeout -> controlled error
external API 500 -> controlled error
invalid webhook signature -> 401
valid webhook signature -> 204
```

### Output

- Small integration demo.
- Tests for success, timeout, and external error.
- README with notes.

### Offer mapping

- Add simple external API integration.
- Fix webhook handler.
- Add timeout to external dependency call.
- Improve error handling around integration.

### Applied in `book-social`

- Future book metadata import.
- Future email/notification integration.
- Future webhook-style exercises.
- Better external dependency handling.

### Estimate

10–18 hours.

### Priority

Medium.

---

## Stage 10 — Bridge to `go-microservices-starter`

**Folder:** `docs/bridge-to-go-microservices-starter.md`

**Goal:** connect Go Web skills to the microservices/gRPC learning path without dragging microservices into this repository too early.

`go-microservices-starter` is the second layer, not the first. For Basic-level paid work, the useful parts are service contracts, local infra, health checks, config, logging, simple workers, and simple service-to-service calls.

Implementation belongs in `go-microservices-starter`; this repository should only contain notes and links.

### What to do

- Add a small note file linking Go Web patterns to `go-microservices-starter`.
- Identify what moves from REST/API practice to service-oriented practice.
- Do not implement gRPC here.
- Do not add NATS/Kafka/Temporal here.
- Keep this repo focused on Go Web.

### Example

Possible note file:

```text
docs/bridge-to-go-microservices-starter.md
```

Suggested content:

```md
# Bridge to go-microservices-starter

## From REST service to gRPC service

1. REST service first.
2. Add gRPC equivalent in `go-microservices-starter`.
3. Keep service/application/repository boundaries.
4. Add local infrastructure.
5. Add README and handoff notes.
```

Possible target work in `go-microservices-starter`:

```text
define .proto
generate Go code
implement one gRPC service
add basic client
add tests
add health/readiness endpoint
document local run command
```

### Output

- Bridge note in this repository.
- Clear separation of responsibilities between repos.
- Future task list for `go-microservices-starter`.

### Offer mapping

- REST endpoint fix.
- gRPC method fix.
- Small integration.
- Service starter cleanup.
- Local development setup.
- Tests around existing API behavior.

### Applied in `book-social`

Not directly. `book-social` remains the applied Go Web/MPA/API project.

### Applied in `go-microservices-starter`

- gRPC basics.
- Service contracts.
- Local Docker Compose.
- Health checks.
- Config/logging.
- Simple worker/background task.
- Simple service-to-service call.

### Estimate

4–8 hours for notes and planning.

Future implementation in `go-microservices-starter`: 12–24 hours for a small focused service-boundary task.

### Priority

Low/medium.

Do this after the Go Web foundation is useful. Do not run toward microservices just because the word sounds employable.

---

## Stage 11 — Offer 1 Basic Portfolio Package

**Folder:** `docs/`

**Goal:** package the learning and applied work into a clear portfolio proof for **Offer 1 Basic — Small Go Backend Fix / API Endpoint**.

This is the stage where learning stops hiding in folders and starts explaining itself. Revolutionary stuff: documentation that says why the code exists.

### What to do

- Create or update `docs/offer-1-basic-proof.md`.
- Create or update `docs/book-social-links.md`.
- Create or update `docs/checklists/backend-endpoint-task.md`.
- Link relevant `go-web-labs` folders.
- Link applied `book-social` tasks.
- Add short case-study style notes.
- Keep it concise.

### Example

Possible files:

```text
docs/offer-1-basic-proof.md
docs/book-social-links.md
docs/checklists/backend-endpoint-task.md
docs/checklists/api-review.md
docs/checklists/handoff-notes.md
```

Possible case study structure:

```md
# Small Go Backend API Proof

## What this demonstrates

- REST endpoint design
- JSON request/response handling
- validation
- error responses
- handler/service/repository split
- DB access
- tests
- OpenAPI
- basic API security

## Learning sources

- `books/lets-go`
- `books/lets-go-further`
- `labs/rest-api`
- `labs/testing`
- `labs/openapi`
- `labs/security`

## Applied proof

- `book-social` v0.1 baseline
- `book-social` v0.2 catalog read models
- `book-social` read-only JSON API slice
- `book-social` tests and docs

## Freelance tasks this supports

- add endpoint
- fix endpoint
- add validation
- add tests
- document API
- add small integration
```

### Output

- Portfolio proof document.
- Links from root README.
- Clear “what this repo proves” explanation.
- Clear mapping from study to applied project.

### Offer mapping

- Small Go backend fix.
- REST API endpoint task.
- API integration task.
- Basic tests and docs.
- Local verification and handoff notes.

### Applied in `book-social`

- Link to concrete features, not vague claims.
- Use `book-social` as the main implementation proof.
- Use `go-web-labs` as learning and pattern proof.

### Estimate

6–12 hours.

### Priority

High after:

1. `Let's Go` closure;
2. basic testing lab or notes;
3. first meaningful `book-social` v0.2 proof.

---

## Relation to `book-social`

`book-social` is the main applied proof for Go Web skills.

### `book-social` v0.1

Finish as baseline:

- Integration / acceptance tests introduction.
- Templ spike.
- HTMX spike.
- Decision: keep Go templates or move to Templ.

### `book-social` v0.2

Use `go-web-labs` topics while implementing:

1. Quality + DB foundation.
2. SQLite/PostgreSQL + migrations.
3. Schema v0.2.
4. Catalog read models.
5. Catalog MPA update.
6. MPA endpoint docs.
7. HTTP middleware + graceful shutdown.
8. Sessions/cookies.
9. User module.
10. Login/register/logout.
11. Optional read-only JSON API slice:
    - `GET /api/books`;
    - `GET /api/books/{slug}`;
    - `GET /api/authors/{slug}`;
    - `GET /api/genres/{slug}`.
12. OpenAPI only for `/api/*`.

### Recommended order

Do not add the JSON API slice before catalog read models are stable.

```text
v0.1 baseline
→ v0.2 quality + DB foundation
→ v0.2 schema
→ v0.2 catalog read models
→ v0.2 MPA catalog update
→ read-only JSON API slice
→ OpenAPI for /api/*
→ security/auth improvements
```

---

## Practical priority

### Current focus

1. Close `Let's Go` with review + README.
2. Finish `book-social` v0.1.
3. Start `Let's Go Further` notes.
4. Build `book-social` v0.2 foundation.
5. Add read-only JSON API slice in `book-social` after catalog read models.
6. Add OpenAPI/security labs after real API endpoints exist.
7. Add portfolio proof package after there is enough applied proof.

### Not now

- Full Echo/Gin comparison.
- Large standalone REST project duplicating `book-social`.
- Complex auth in labs.
- Full API write operations before `book-social` auth and validation are ready.
- Microservices implementation in this repo.
- NATS/Kafka/Temporal in `go-web-labs`.
- Perfect framework abstraction, because that is how small repos become sad enterprises.

---

## Suggested repository map

```text
books/
  lets-go/
    README.md
  lets-go-further/
    README.md
labs/
  rest-api/
    README.md
  layered-api/
    README.md
  testing/
    README.md
  openapi/
    README.md
    openapi.yaml
  security/
    README.md
  integrations/
    README.md
frameworks/
  echo-lab/
    README.md
  templ/
    README.md
docs/
  book-social-links.md
  bridge-to-go-microservices-starter.md
  offer-1-basic-proof.md
  checklists/
    backend-endpoint-task.md
    api-review.md
    graceful-shutdown.md
    build-audit-versioning.md
    handoff-notes.md
```

---

## Offer 1 Basic proof checklist

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
- [ ] production API notes: graceful shutdown, metrics, build/audit/versioning;
- [ ] bridge notes to `go-microservices-starter`;
- [ ] links to applied implementation in `book-social`;
- [ ] concise portfolio proof document for Offer 1 Basic.

---

## Stage estimates summary

| Stage | Name                                      |                                 Estimate | Priority                |
|------:|-------------------------------------------|-----------------------------------------:|-------------------------|
|     0 | Repository consolidation                  |                                     2–4h | High                    |
|     1 | Close / refresh `Let's Go`                |                                     3–6h | Current                 |
|     2 | REST API basics                           |                                    8–14h | Medium                  |
|     3 | Handler → Service → Repository            |                                    8–16h | Medium                  |
|     4 | Testing lab                               |                                    8–16h | High                    |
|     5 | `Let's Go Further`: API Core              |                                   25–45h | High                    |
|     6 | OpenAPI lab                               |                                    6–12h | Medium                  |
|     7 | API security basic                        |                                   10–20h | Medium                  |
|     8 | `Let's Go Further`: Production API Topics |                                   15–30h | Medium                  |
|     9 | Integration / External API lab            |                                   10–18h | Medium                  |
|    10 | Bridge to `go-microservices-starter`      | 4–8h notes, 12–24h future implementation | Low/Medium              |
|    11 | Offer 1 Basic Portfolio Package           |                                    6–12h | High after proof exists |

---

## Short positioning summary

`go-web-labs` should prove that the fundamentals were studied intentionally.

`book-social` should prove that the fundamentals were applied in a real project.

`go-microservices-starter` should become the next layer after the Go Web foundation is useful, not the place where every unresolved thought goes to start a distributed systems support group.
