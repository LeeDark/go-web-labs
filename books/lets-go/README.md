# Let's Go / Snippetbox

## Purpose

Study project based on *Let's Go* by Alex Edwards.

This project is kept as a learning artifact for Go web fundamentals and as a reference for future small Go backend tasks.

## What this project demonstrates

- `net/http` routing and handlers
- Middleware chains
- HTML templates
- Static assets and embedded files
- Forms and validation
- Sessions, cookies, and CSRF protection
- MySQL-backed models
- Authentication basics
- Handler tests and model integration tests
- Secure server configuration basics

## Repository location

- Module: `books/lets-go`
- Application: `books/lets-go/snippetbox`
- Entrypoint: `books/lets-go/snippetbox/cmd/web`
- Templates and static assets: `books/lets-go/snippetbox/ui`
- Models and validators: `books/lets-go/snippetbox/internal`

## How to run

Run from the application directory:

```bash
cd books/lets-go/snippetbox
go run ./cmd/web
```

The application starts an HTTPS server on `:4000` by default.

Useful flags:

```bash
go run ./cmd/web -addr=":4000" -dsn="web:pass@/snippetbox?parseTime=true"
```

The server expects local TLS files at:

```text
./tls/cert.pem
./tls/key.pem
```

This repo currently includes local cert files. If they need to be regenerated, run this from `books/lets-go/snippetbox`:

```bash
go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

## How to test

Run all tests from the module root:

```bash
cd books/lets-go
go test ./...
```

In restricted environments, use a temporary Go cache:

```bash
GOCACHE=/tmp/go-web-labs-go-cache go test ./...
```

Some tests open local sockets with `httptest.NewTLSServer`; model integration tests connect to MySQL unless skipped with `-short`.

Useful focused checks:

```bash
GOCACHE=/tmp/go-web-labs-go-cache go test ./snippetbox/cmd/web -run 'TestPing$|TestCommonHeaders|TestHumanDate'
GOCACHE=/tmp/go-web-labs-go-cache go test -short ./snippetbox/internal/models
```

Coverage helper:

```bash
cd books/lets-go/snippetbox
make cover
```

## Configuration

Runtime flags:

- `-addr`: HTTP network address, default `:4000`
- `-static-dir`: path to static assets, default `./ui/static`
- `-dsn`: MySQL data source name, default `web:pass@/snippetbox?parseTime=true`

The current routes use embedded UI files from `snippetbox/ui`, so templates and static assets are compiled into the binary.

## Database setup

The application expects a MySQL database named `snippetbox` and a user matching the runtime DSN.

Default application DSN:

```text
web:pass@/snippetbox?parseTime=true
```

Model integration tests expect a separate MySQL database and user:

```text
test_web:Pass_Test1@/test_snippetbox?parseTime=true&multiStatements=true
```

The test schema is managed by:

- `snippetbox/internal/models/testdata/setup.sql`
- `snippetbox/internal/models/testdata/teardown.sql`

No production migration tool is configured for this study project.

## Useful patterns learned

- Handlers are methods on an `application` struct with shared dependencies.
- Route setup is centralized in `cmd/web/routes.go`.
- Middleware is composed with `alice`.
- Templates are parsed once into a cache and rendered through helper methods.
- Error responses are centralized in helper methods.
- Models keep SQL access behind small interfaces for handler tests.
- Test doubles live under `internal/models/mocks`.
- Sessions are stored in MySQL through `scs/mysqlstore`.
- CSRF protection is applied to dynamic routes with `nosurf`.

## Review notes

- The project is an MPA Snippetbox application, not a REST API lab.
- Static files and templates are embedded through `snippetbox/ui/efs.go`.
- Local TLS certs are required because the server uses `ListenAndServeTLS`.
- Full tests need permission to open local sockets; model integration tests need local MySQL.
- `go test ./...` was attempted during refresh and failed in the sandbox because local socket creation and MySQL access are not permitted.

## Remaining TODO / Guided exercises

Chapter 16 guided exercises are intentionally kept as backlog, not implemented in this refresh.

- TODO: Copy the exact Chapter 16 guided exercise titles from the local book notes into this section.
- TODO: Pick one guided exercise per follow-up task.
- TODO: Keep each exercise implementation small and covered by focused tests where practical.
- TODO: Do not mix Chapter 16 exercises with REST API labs or `Let's Go Further` work.

## Status

Completed with TODO exercises.

The existing Snippetbox implementation has been reviewed, a small signup handler flow bug was fixed, and the project now has run, test, configuration, database, and follow-up notes.
