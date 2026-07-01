# Go Web Topics Index

This file is a topic inventory for Go web learning. It is not the execution roadmap.

Use `PLAN.md` for stage order and current priorities. Use this file to collect topics, vocabulary, and reference areas that may become future notes, labs, or checklist items.

Priority labels:

- Current: active or near-active repo focus.
- Next: planned soon after the current focus.
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
- HTTP methods and status codes — Current
- Headers, content types, content negotiation — Next
- URL structure, query strings, path params — Current
- DNS basics — Reference only
- Web servers (nginx, apache, caddy)
- Reverse proxy basics — Later

## Application types
- MPA — Current
- Server-rendered Go app — Current
- API-only Backend — Next
- REST API backend — Next
- SPA — Reference only
- Hybrid MPA + JSON endpoints — Later
- Rendering layer: SSR/CSR, Hybrid — Reference only

## API Design
- REST — Next
- Resource naming — Next
- Request/response shape conventions — Next
- Error response format — Next
- OpenAPI, API docs — Later
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
- Request parsing: path, query, body — Current
- Response writing: status, headers, body — Current
- Content type handling — Next
- Stateful HTTP, sessions, cookies
- HTML templates and dynamic templates, MPA
- Template layouts and partials — Current
- HTML forms
- Form validation and redisplay — Current
- Flash messages — Current
- HTML/static embedding
- Static files — Current
- JSON requests/responses, SPA
- JSON decoding limits and unknown fields — Next
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
  - Secure cookie flags — Current / Later
  - Session fixation basics — Later
  - Input validation vs output escaping — Current
  - Secrets/config handling — Later
  - Request body size limits — Later
  - Rate limiting — Later
  - Basic abuse protection — Later
  - Safe logging of sensitive data — Later

## Web frameworks
- net/http — Current base stack
- chi — Current primary router/framework
- echo — Later comparison
- gin — Later comparison
- fiber — Reference only
- gorilla/mux — Reference only

## Supporting Topics

## Application structure
- Skeleton
  - `cmd/` and `internal/` package layout — Next
  - Handler dependencies and app struct — Current / Next
  - Layering: Handlers, Services, Repositories — Next
  - Dependency Injection — Next
  - DTOs vs domain models — Next
  - Repository interfaces — Next
  - Configuration — Later
- Engineering hygiene
  - Logging — Later
  - Errors, Error handling — Next
  - Error boundaries between layers — Next
  - Quality control, linters — Later
  - Building, versioning — Later

## Quality & Operations
- Local verification — Current
  - Golden path local run command
  - README run/test commands
  - `go test ./...`
- Testing — Stage 5
  - Table-driven tests
  - Handler tests
  - Service tests
  - Integration test basics
- Production operations — Stage 8
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
- Small, reviewable changes — Current
- Debugging with logs and tests — Current
- Reading existing code before refactoring — Current
- Documentation discipline — Current
- SOLID — Reference only
- YAGNI — Current
- KISS — Current

## Adjacent topics
- Frontend: HTML, CSS, JavaScript — Reference only
- Git — Current
- HTMX — Reference only / later comparison
- Templ — Reference only / later comparison
- Docker basics — Later
- Makefile/task runner basics — Later
- NoSQL — Reference only
- Security and hashing algorithms — Later

## AI
- AI-assisted code review — Current
- AI-assisted documentation drafts — Current
- AI-assisted review checklist — Current
- AI-assisted README cleanup — Current
- AI-assisted test idea generation — Later
- AI-assisted test case brainstorming — Later
- Prompting workflow notes — Later
- Human verification of AI output — Current
