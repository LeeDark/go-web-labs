# Let's Go Further

Study project based on *Let's Go Further* by Alex Edwards.

## Purpose

This folder is for practicing production-style Go API patterns through the
Greenlight project.

Current focus: API Core.

## Current Scope

- Project structure
- HTTP routing
- JSON responses
- JSON request parsing
- Validation
- PostgreSQL
- Migrations
- CRUD
- Filtering, sorting, and pagination

Out of current scope:

- Authentication
- Authorization
- Rate limiting
- Metrics
- Deployment
- Production hardening

## Chapter Notes

### Chapter 1: Introduction

Greenlight is introduced as a JSON API for movie data. The chapter defines the
project direction, expected tools, and the high-level API shape.

The main value for this repository is orientation: use Greenlight to learn
reusable API patterns, while keeping the study project separate from future
portfolio applications.

Expected early endpoints:

- `GET /v1/healthcheck`
- Movie CRUD endpoints under `/v1/movies`

### Chapter 2: Getting Started

Chapter 2 sets up the first runnable API skeleton.

Implemented pieces:

- `cmd/api` application entry point
- `config` struct with `port` and `env` command-line flags
- `application` struct for handler dependencies
- structured logger using `log/slog`
- HTTP server with basic timeout settings
- centralized routes in `routes.go`
- `GET /v1/healthcheck`
- placeholder movie routes:
  - `POST /v1/movies`
  - `GET /v1/movies/:id`
- `readIDParam()` helper for positive integer URL IDs

Useful pattern:

Keep route definitions, handlers, and small helpers separated early. This keeps
`main.go` focused on wiring configuration, dependencies, logging, and the HTTP
server.

## How to Run

From this folder:

```sh
go run ./greenlight/cmd/api
```

With custom flags:

```sh
go run ./greenlight/cmd/api -port=3030 -env=production
```

Example requests:

```sh
curl -i localhost:4000/v1/healthcheck
curl -i -X POST localhost:4000/v1/movies
curl -i localhost:4000/v1/movies/123
```

## How to Test

TODO: add when tests are introduced.

## TODO

- Verify the healthcheck endpoint with `curl`.
- Convert plain-text responses to JSON in Chapter 3.
- Add JSON error responses in Chapter 3.
