# Task History

## Stage 1 - Close / refresh Let's Go

Status: completed with TODO exercises

Goal:

Close or refresh the `books/lets-go` study project as a clean learning artifact.

Completed work:

- reviewed the current implementation;
- ran tests where possible;
- fixed a small signup handler control-flow bug;
- created `books/lets-go/README.md`;
- added guided exercises as TODO backlog;
- created `AGENTS.md`;
- created `docs/ai/repository-context.md`.

Notes:

- Do not rewrite the project.
- Do not start REST API labs in this task.
- Do not start `Let's Go Further` implementation in this task.
- Full tests need local socket permission and MySQL access.

## Stage 2 - Let's Go Further API Core

Status: in progress (Chapters 1-6 complete)

Goal:

Work through the API Core portion of *Let's Go Further* using the Greenlight
study project.

Completed work (2026-07-06 to 2026-07-15):

- created the `books/lets-go-further` Go module and Greenlight API skeleton;
- implemented routing, JSON responses, JSON errors, panic recovery, strict
  JSON request parsing, and movie validation;
- configured PostgreSQL DSN settings and the connection pool;
- added versioned migrations for the `movies` table and its constraints;
- added and updated Stage 2 learning notes and the implementation README.

Remaining work:

- Chapters 7-9: movie persistence, CRUD, filtering, sorting, and pagination;
- verify the Chapter 6 migrations against a local PostgreSQL instance;
- add tests when the relevant chapters introduce them.
