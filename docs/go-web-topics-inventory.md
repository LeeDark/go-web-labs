# Go Web Topics Index

This file is a topic inventory for Go web learning. It is not the execution roadmap.

Use `PLAN.md` for stage order and current priorities. Use this file to collect topics, vocabulary, and reference areas that may become future notes, labs, or checklist items.

Priority labels:

- Finished: covered by Stage 0 or Stage 1.
- Active: part of the current parallel work.
- Stage 2: book learning track, `Let's Go Further` API Core.
- Offer 1 foundation: Stages 3-7, freelance path labs and notes.
- Later: useful, but not immediate.
- Reference only: awareness topic, not planned implementation in this repository.

## Sources
- https://roadmap.sh/golang
- https://roadmap.sh/backend
- https://softaims.com/roadmap/go
- https://softaims.com/roadmap/backend

## Web

## Foundations
Reference only / cheat sheet.
- Networking
- Internet, TCP/IP
- HTTP protocol
- HTTP methods and status codes — Finished / Offer 1 foundation
- Headers, content types, content negotiation — Stage 2 / Offer 1 foundation
- URL structure, query strings, path params — Finished / Offer 1 foundation
- DNS basics — Reference only
- Web servers (nginx, apache, caddy)
- Reverse proxy basics — Later

## Application types
- MPA — Finished
- Server-rendered Go app — Finished
- API-only Backend — Stage 2 / Offer 1 foundation
- REST API backend — Stage 2 / Offer 1 foundation
- SPA — Reference only
- Hybrid MPA + JSON endpoints — Later
- Rendering layer: SSR/CSR, Hybrid — Reference only

## API Design
- REST — Stage 2 / Offer 1 foundation
- Resource naming — Stage 2 / Offer 1 foundation
- Request/response shape conventions — Stage 2 / Offer 1 foundation
- Error response format — Stage 2 / Offer 1 foundation
- OpenAPI, API docs — Offer 1 foundation
- Pagination styles: limit/offset, cursor — Later
- Filtering and sorting conventions — Later
- Backward-compatible API changes — Later
- GraphQL — Reference only
- SOAP — Reference only
- JSON-API — Reference only
- JSON-RPC — Reference only
- gRPC — Later, in `go-microservices-starter`

## Go HTTP Stack
- net/http
- Routing and API endpoints
- Middleware
- HTTP request context
- Request parsing: path, query, body — Finished / Offer 1 foundation
- Response writing: status, headers, body — Finished / Offer 1 foundation
- Content type handling — Stage 2 / Offer 1 foundation
- Stateful HTTP, sessions, cookies
- HTML templates and dynamic templates, MPA
- Template layouts and partials — Finished
- HTML forms
- Form validation and redisplay — Finished
- Flash messages — Finished
- HTML/static embedding
- Static files — Finished
- JSON requests/responses, SPA
- JSON decoding limits and unknown fields — Stage 2 / Offer 1 foundation
- CRUD operations
- Validation
- Filtering, sorting, and pagination
- Timeouts: server, read/write, request context — Later
- Problem details / structured errors — Later
- Think about
  - Redirects
  - File uploads, multipart
  - Error responses
  - Idempotency basics
  - API versioning basics
  - HTML escaping, template safety

## Security
- Security core
  - HTTPS/TLS — Later
  - CORS — Later
  - CSRF — Later
- Authentication — Later
- Authorization — Later
- Security, protection, resilience
  - Password hashing — Later
  - Secure cookie flags — Finished / Later
  - Session fixation basics — Later
  - Input validation vs output escaping — Finished / Offer 1 foundation
  - Secrets/config handling — Later
  - Request body size limits — Later
  - Rate limiting — Later
  - Basic abuse protection — Later
  - Safe logging of sensitive data — Later

## Web frameworks
- net/http — Finished base stack / Offer 1 foundation
- chi — Active primary router/framework
- echo — Later comparison
- gin — Later comparison
- fiber — Reference only
- gorilla/mux — Reference only

## Supporting Topics

## Application structure
- Skeleton
  - `cmd/` and `internal/` package layout — Offer 1 foundation
  - Handler dependencies and app struct — Finished / Offer 1 foundation
  - Layering: Handlers, Services, Repositories — Offer 1 foundation
  - Dependency Injection — Offer 1 foundation
  - DTOs vs domain models — Offer 1 foundation
  - Repository interfaces — Offer 1 foundation
  - Configuration — Later
- Engineering hygiene
  - Logging — Later
  - Errors, Error handling — Stage 2 / Offer 1 foundation
  - Error boundaries between layers — Offer 1 foundation
  - Quality control, linters — Later
  - Building, versioning — Later

## Quality & Operations
- Local verification — Finished / Active
  - Golden path local run command
  - README run/test commands
  - `go test ./...`
- Testing — Offer 1 foundation, Stage 5
  - Table-driven tests
  - Handler tests
  - Service tests
  - Integration test basics
- Production operations — Later, Stage 8
  - Graceful shutdown
  - Graceful shutdown signal handling
  - Health checks
  - Request logging and request ID
  - Metrics basics
  - Build, audit, versioning
  - Deployment and hosting basics
- Later / reference
  - Race detector
  - `go vet`
  - Caching basics
  - Observability

## Relational Databases
- SQL fundamentals — Later
  - ACID, Normalization
  - Schema, Migrations
  - MySQL, PostgreSQL, SQLite
- Performance basics — Later
  - Transactions
  - Indexes
  - Profiling
- Data access — Later
  - database/sql, sqlc, query builders
  - ORM, GORM
  - Connection pool basics
  - Context-aware queries
  - NULL handling
  - SQL injection prevention
  - Test database strategy
  - Migrations rollback policy

## Skills
- Small, reviewable changes — Active
- Debugging with logs and tests — Active
- Reading existing code before refactoring — Active
- Documentation discipline — Active
- SOLID — Reference only
- YAGNI — Active
- KISS — Active

## Adjacent topics
- Frontend: HTML, CSS, JavaScript — Reference only
- Git — Active
- HTMX — Reference only / later comparison
- Templ — Reference only / later comparison
- Docker basics — Later
- Makefile/task runner basics — Later
- NoSQL — Reference only
- Security and hashing algorithms — Later

## AI
- AI-assisted code review — Active
- AI-assisted documentation drafts — Active
- AI-assisted review checklist — Active
- AI-assisted README cleanup — Active
- AI-assisted test idea generation — Later
- AI-assisted test case brainstorming — Later
- Prompting workflow notes — Later
- Human verification of AI output — Active
