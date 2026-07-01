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

# Web

## Foundations
Reference only / cheat sheet.
- Networking
- Internet, TCP/IP
- HTTP protocol
- Web servers (nginx, apache, caddy)

## Application types
- MPA — Current
- API-only Backend — Next
- SPA — Reference only
- Rendering layer: SSR/CSR, Hybrid — Reference only

## API Design
- REST — Next
- OpenAPI, API docs — Later
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
- Stateful HTTP, sessions, cookies
- HTML templates and dynamic templates, MPA
- HTML Forms
- HTML/static embedding
- JSON requests/responses, SPA
- CRUD operations
- Validation
- Filtering, sorting, and pagination
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
  - Rate limiting — Later
  - Basic abuse protection — Later

## Web frameworks
- chi — Current primary router/framework
- echo — Later comparison
- gin — Later comparison
- fiber — Reference only
- gorilla/mux — Reference only

# Supporting Topics

## Application structure
- Skeleton
  - Layering: Handlers, Services, Repositories — Next
  - Dependency Injection — Next
  - Configuration — Later
- Engineering hygiene
  - Logging — Later
  - Errors, Error handling — Next
  - Quality control, linters — Later
  - Building, versioning — Later

## Quality & Operations
- Testing — Stage 5
  - Table-driven tests
  - Handler tests
  - Service tests
  - Integration test basics
- Production operations — Stage 8
  - Graceful shutdown
  - Health checks
  - Metrics basics
  - Build, audit, versioning
  - Deployment and hosting basics
- Later / reference
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

## Skills
- SOLID — Reference only
- YAGNI — Current
- KISS — Current

## Think
- Frontend: HTML, CSS, JavaScript — Reference only
- Git — Current
- NoSQL — Reference only
- Security and hashing algorithms — Later

## AI
- AI-assisted code review — Current
- AI-assisted documentation drafts — Current
- AI-assisted test idea generation — Later
- Prompting workflow notes — Later
