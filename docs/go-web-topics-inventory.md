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
- HTTP Protocol
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
- Routing, API Endpoints
- Middleware
- HTTP Request Context
- Stateful HTTP, Sessions, Cookies
- HTML Templates & Dynamic Templates, MPA
- HTML Forms
- HTML/Static Embedding
- JSON Requests/Responses, SPA
- CRUD, Operations
- Validation
- Filtering, Sorting, Pagination, etc
- Think about
  - Redirects
  - File uploads, multipart
  - Error responses
  - Idempotency basics
  - API versioning basics
  - HTML escaping, Template safety

## Security
- Security core
  - HTTPS/TLS — Later
  - CORS — Later
  - CSRF — Later
- Authentication — Later
- Authorization — Later
- Security, Protection, Resilience
  - Rate limiting — Later
  - Basic abuse protection — Later

## Web frameworks
- chi
- echo
- gin
- fiber
- gorilla/mux

## Advanced Topics

# Supporting Topics

## Application structure
- Skeleton
  - Layering: Handlers, Services, Repositories — Next
  - Dependency Injection — Next
  - Configuration — Later
- Engineering hygiene
  - Logging — Later
  - Errors, Error handling — Next
  - Quality Control, Linters — Later
  - Building, Versioning — Later

## Quality & Operations
- Current / Next
  - Testing
- Later
  - Deployment, Hosting
  - Caching basics
  - Graceful Shutdown
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
  - Databases packages, sqlc, Query builders
  - ORM, Gorm

## Skills
- SOLID — Reference only
- YAGNI — Current
- KISS — Current

## Think
- Frontend: HTML, CSS, JavaScript — Reference only
- Git — Current
- NoSQL — Reference only
- Security & Hashing Algorithms — Later

## AI
- AI Assisted Coding & Prompt Engineering Roadmap
- 
