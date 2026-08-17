# Security Lab

## Purpose

`labs/security` contains concise source-study notes, checklists, and small examples for Stage 7. It
is not a second application or a shared authentication package for `book-social`.

The current wave is **Stage 7A**: `book-social` v0.2.5–v0.2.6 authentication alongside the first TDD
foundations unit. Stage 7B (CORS, API rate limiting, and API-specific controls) does not begin until
a stable `/api/*` contract exists.

The applied v0.2.5–v0.2.6 scope currently limits the first slice to registration, login, logout,
and `GET /me`; use DB-backed opaque sessions, renew the session on authentication, invalidate it on
logout, keep `HttpOnly`/`SameSite=Lax` cookie policy, and exclude profile, activation, recovery,
API auth, roles/RBAC routes, CORS, rate limiting, JWT, and OpenAPI.

Actor model checkpoint: the first release has `anonymous` and `authenticated user` actors only.
`GET /me` is the first protected route and is available only to the current user established by
server-side session identity. `admin`/`is_admin` does not create a bypass or route; ownership rules
for the future private-library resource remain deferred until that resource has an explicit contract.

Route/form contract checkpoint: the first slice is `GET/POST /register`, `GET/POST /login`,
CSRF-protected `POST /logout`, and protected `GET /me`. Successful registration/login create a new
session and redirect to `/me`; invalid credentials return a neutral `422`; anonymous `/me` requests
redirect to `/login`; `next` is intentionally absent and GET/HEAD never change state.

Ownership checkpoint: `GET /me` authorizes only the current authenticated identity. The first release
has no private library resource, so owner/non-owner behavior is deferred to the resource contract;
when it appears, anonymous, owner, and non-owner cases must be tested separately. Navigation,
hidden fields, database roles, and `is_admin` are never authorization controls by themselves.

Session lifecycle checkpoint: `book-social` uses a DB-backed opaque session with cookie
`book_social_session`; the database stores only the token hash, `user_id`, and lifecycle timestamps.
Successful registration/login creates a new session, missing/invalid/expired tokens mean anonymous,
logout invalidates the row and clears the cookie, the lifetime is seven days without sliding renewal,
and expired-row cleanup is lazy. Session middleware excludes `/static/*` and `/healthz`.

Password policy checkpoint: use bcrypt from one maintained auth package, with passwords 12–128
characters and matching confirmation, and persist only the adaptive hash. Plaintext passwords and
hashes stay out of DTOs, page models, responses, errors, logs, metrics, and fixtures. Unknown login
and wrong password share the neutral `Invalid login or password` outcome; password reset/recovery and
activation are deferred.

Error/logging checkpoint: client responses use stable safe outcomes (`422` for validation/neutral
credentials, redirect for anonymous session, generic `500` for internal failures); logs may contain
request ID, route, operation, safe actor state, and typed error class only. Passwords, hashes, raw
session/CSRF tokens, cookie or authorization headers, submitted credentials, and private content are
denylisted. API rate limiting, CORS, login throttling, and account lockout require separate scope.

Contract review checkpoint: steps 0–6 are consolidated in the applied `book-social` plan. The accepted
contract is limited to the MPA auth foundation and explicitly separates v0.2.5 persistence/service
work from v0.2.6 forms/UI. It is a pre-implementation contract, not evidence that the flow already
exists.

## *Let's Go* Chapter 8: Stateful HTTP

Source: Alex Edwards, *Let's Go*, 2nd edition (2025), Chapter 8: “Choosing a session manager”,
“Setting up the session manager”, and “Working with session data”.

### Summary

HTTP does not retain the state between requests. A session connects multiple requests from one
browser user: a cookie carries a random session token/ID, and a server-side store uses it to find
session data and its expiry. The token itself should not contain a flash message, user identity, or
other session payload.

| Approach                                              | When it fits                                                             | Stage 7A constraint                                                                                           |
|-------------------------------------------------------|--------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------|
| Client-side session data in a signed/encrypted cookie | You specifically need client-side data and accept its constraints.       | Do not use it for private session data without an explicit decision on size, exposure, and lifecycle.         |
| Server-side session store                             | You need controlled invalidation, expiry, and private server-side state. | For an auth flow, prefer a manager that can renew the session ID after login to reduce session-fixation risk. |

For Snippetbox, the book chooses `alexedwards/scs`: the manager receives a store, lifetime, and the
`LoadAndSave` middleware. The middleware loads a session by cookie before a handler and saves
changes afterward. This is a useful reference, not a ready-made `book-social` decision: the book's
MySQL table and 12-hour lifetime do not transfer automatically.

### What applies to `book-social`

- Choose a server-side store and session manager after the auth contract is fixed. With a
  server-side model, the manager must be able to renew or rotate the session ID after successful
  authentication.
- Keep only an opaque, high-entropy identifier in the cookie. User IDs and other data remain in the
  server-side session store; passwords, CSRF values, and private content never enter logs or
  templates.
- Define the lifecycle explicitly: creation after successful login, renewal at authentication,
  expiry/idle policy, deletion or invalidation at logout, and the outcome for an expired or missing
  session.
- Apply load/save middleware only to dynamic MPA routes that read or change session state. Static
  files and the health endpoint should not create session-store work without a reason.
- Pass session state to handlers through the request context supplied by middleware, not through a
  global variable, query parameter, or trusted client header.
- Use one-time `Pop` semantics for flash messages; use ordinary `Get` only when a value must survive
  the next request.

### Decisions that need separate treatment

The chapter provides a model, but it does not replace the Stage 7A security decisions:

- the concrete `book-social` store and its migration/cleanup policy;
- cookie `HttpOnly`, `Secure`, `SameSite`, name, path, and development-versus-HTTPS behavior;
- absolute and idle expiration, renewal semantics, and concurrent-session policy;
- login/register/logout routes, redirect and error contract, and CSRF integration;
- the ownership/authorization boundary for a private resource;
- retention and cleanup of expired server-side sessions, plus operational observability without
  recording session IDs.

Record these decisions in the auth contract before implementation; do not infer them solely from a
chosen library.

## Minimal flow and testable invariants

```text
successful login
  → create/renew server-side session
  → response sets opaque cookie
  → protected route loads current user from session
  → logout invalidates session and clears its browser state
```

| Risk                                         | Minimum verification before handoff                                |
|----------------------------------------------|--------------------------------------------------------------------|
| An anonymous user obtains a private resource | HTTP test: redirect or refusal without a session.                  |
| An ID fixed before login remains valid       | Test: successful login renews the session identifier.              |
| Logout does not end access                   | HTTP flow: login → protected action → logout → refusal.            |
| A flash message appears repeatedly           | Test: set flash → first render shows it → refresh does not.        |
| A secret reaches the client or logs          | Handler/log assertion: no password, raw session ID, or CSRF value. |
| A static or health request creates a session | Router/middleware test: no session cookie or store write.          |

Use `httptest`/`ServeHTTP` and a deterministic test store when testing the HTTP/session boundary. A
real database integration test is needed only after a concrete persistence risk exists, and it must
use an explicitly disposable database.

## Next study step

Before v0.2.5 code, describe the session contract: actors, routes, session contents, store, cookie
policy, lifetime, renewal, logout invalidation, and exact test cases. Then study *Let's Go* Chapter
10 for password handling, authorization, and CSRF. Do not add a login form or user model based on
Chapter 8 alone.
