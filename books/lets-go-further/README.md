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

## How to Run

TODO: add after Chapter 2 server setup.

## How to Test

TODO: add when tests are introduced.

## TODO

- Add basic API server.
- Add healthcheck endpoint.
- Record run command.
