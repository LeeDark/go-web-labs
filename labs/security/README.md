# Security Lab

## Purpose

`labs/security` contains concise source-study notes, checklists, and small examples for Stage 7. It
is not a second application or a shared authentication package for `book-social`.

The current wave is **Stage 7A**. The `book-social` v0.2.5 foundation is applied at commit `41a8ddb`;
v0.2.6 completes the user-facing authentication flow. Stage 7B (CORS, API rate limiting, and
API-specific controls) does not begin until a stable `/api/*` contract exists.

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

## Applied foundation and remaining flow contract

The accepted v0.2.5–v0.2.6 scope limits the first user-facing slice to registration, login, logout, and
`GET /me`; use DB-backed opaque sessions, renew the session on authentication, invalidate it on
logout, keep `HttpOnly`/`SameSite=Lax` cookie policy, and exclude profile, activation, recovery, API
auth, roles/RBAC routes, CORS, rate limiting, JWT, and OpenAPI.

The v0.2.5 foundation at `41a8ddb` implements the migrations, password and session services,
repositories, cookie manager, current-user context, authentication guard, and global
`http.CrossOriginProtection`. It deliberately has no production auth routes or auth-aware UI; those
belong to v0.2.6.

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
when it appears, anonymous, owner, and non-owner cases must be tested separately. Navigation, hidden
fields, database roles, and `is_admin` are never authorization controls by themselves.

Session lifecycle checkpoint: `book-social` uses a DB-backed opaque session with cookie
`book_social_session`; the database stores only the token hash, `user_id`, and lifecycle timestamps.
Successful registration/login creates a new session, missing/invalid/expired tokens mean anonymous,
logout invalidates the row and clears the cookie, the lifetime is seven days without sliding renewal,
and expired-row cleanup is lazy. Session middleware excludes `/static/*` and `/healthz`.

Password policy checkpoint: use bcrypt from one maintained auth package, with at least 12 Unicode
characters, at most 72 UTF-8 bytes, and matching confirmation, and persist only the adaptive hash.
Plaintext passwords and hashes stay out of DTOs, page models, responses, errors, logs, metrics, and
fixtures. Unknown login and wrong password share the neutral `Invalid login or password` outcome;
password reset/recovery and activation are deferred.

Error/logging checkpoint: client responses use stable safe outcomes (`422` for validation/neutral
credentials, redirect for anonymous session, generic `500` for internal failures); logs may contain
request ID, route, operation, safe actor state, and typed error class only. Passwords, hashes, raw
session/CSRF tokens, cookie or authorization headers, submitted credentials, and private content are
denylisted. API rate limiting, CORS, login throttling, and account lockout require separate scope.

Contract review checkpoint: the accepted contract explicitly separates the applied v0.2.5
foundation from the planned v0.2.6 forms/UI. Commit `41a8ddb` is evidence for the foundation only,
not for a working registration/login/logout/`GET /me` flow.

## *Let's Go* Chapter 10: User authentication

Chapter 10 builds a complete authentication flow on top of the stateful sessions from Chapter 8. In
the book's example, anonymous users can view snippets and sign up, while only authenticated users can
create snippets.

### What the chapter covers

- **Routes and forms.** Signup, login, and logout have separate routes; every state-changing action,
  including logout, uses `POST`. Forms validate input and render field-level errors, while credential
  failures use a neutral non-field error.
- **Users model and password encryption.** The users table stores an ID, name, email, and bcrypt hash,
  never a plaintext password. Signup validates required fields, email format, password length, and
  duplicate email. Bcrypt cost is selected against workload and user-facing latency rather than copied
  as an unmeasured constant.
- **Login and session fixation.** `Authenticate` loads the hash by email and compares it with bcrypt;
  an unknown email and a wrong password produce the same invalid-credentials outcome. On success the
  session ID is renewed, the authenticated user ID is stored in the session, and the handler redirects.
- **Logout.** A protected `POST` renews the session ID, removes the authentication marker, adds a
  flash message, and redirects home. Logout must end server-side access, not merely change navigation.
- **Authorization.** Authentication status is exposed to template data, but hiding a link is not an
  access control. Separate middleware requires authentication for protected routes and adds
  `Cache-Control: no-store`; anonymous requests are redirected to login.
- **CSRF.** SameSite cookies are only a defensive layer. The book adds token middleware (`nosurf`),
  places the token in hidden form fields, and verifies it server-side for state-changing requests.
  CSRF middleware covers dynamic routes, including logout; static files are excluded.

### Applying it to `book-social`

The chapter is a source model for bcrypt, session renewal when privilege state changes,
authorization middleware, and token-based CSRF. `book-social` adapts the browser-request protection
to the standard-library `http.CrossOriginProtection` boundary instead of copying `nosurf`; the
Snippetbox routes, MySQL schema, flash messages, and redirect to `/snippet/create` are not our contract.

For the current slice, keep only registration/login/logout and protected `GET /me`: successful
registration and login create a new DB-backed session and redirect to `/me`; anonymous `/me` redirects
to `/login`; logout must pass the cross-origin protection boundary and invalidate the session. The
Stage 7A contract already fixes the password policy, neutral login error,
`HttpOnly`/`SameSite=Lax` cookie, no `next` parameter, and the exclusion of `/static/*` and
`/healthz`. Profile, recovery, roles/RBAC, and private-resource ownership remain out of scope.

### Testable invariants

1. The database contains only an adaptive password hash; plaintext never reaches responses, logs, or
   fixtures.
2. Login renews session identity before recording the current user; logout invalidates the old identity.
3. Invalid credentials do not reveal whether an email exists.
4. A protected route checks server-side identity, not navigation, hidden fields, or a client-supplied ID.
5. Unsafe cross-origin browser requests are refused, same-origin requests succeed, and GET/HEAD do
   not change authentication state.

## *Let's Go* Chapter 11: Using request context

Chapter 11 replaces repeated session checks with a request-scoped authentication decision. The goal is
to validate the session identity once in middleware and make the result available to later middleware,
handlers, and templates during the same request.

### What the chapter teaches

- **Context is attached to the request.** A handler starts from `r.Context()`, derives a context with
  `context.WithValue`, and passes a copied request created by `r.WithContext(ctx)` to the next handler.
  The original request is not mutated in place.
- **Values require type-safe keys and assertions.** String keys risk collisions with other packages;
  the book defines a private `contextKey` type and a package-owned key. Reads use a checked type
  assertion and fail closed when the value is absent or has the wrong type.
- **Authentication is validated once.** Middleware reads the authenticated user ID from the session,
  verifies that the user still exists in the database, and places the authenticated state in context.
  An absent or deleted user is treated as anonymous; a database failure is an internal error.
- **Authorization remains a separate decision.** The context-backed authentication helper feeds the
  protected-route middleware, which redirects anonymous users and marks protected responses
  `Cache-Control: no-store`. A context value does not itself grant access.
- **Context scope is deliberately narrow.** Request context is for request-scoped data that travels
  through the handler chain. It is not a general dependency container for databases, loggers,
  template caches, or other application-lifetime services.

### Applying it to `book-social`

The book's boolean `isAuthenticated` value is a useful minimum example, but the applied contract
needs a validated current-user boundary for `GET /me`. Middleware may load the user associated with
the DB-backed session once and pass only request-scoped, non-secret identity data downstream. It must
not expose raw session tokens, cookies, passwords, or client-supplied user IDs through context.

Use a package-owned typed key, keep the lookup and failure behavior explicit, and make the safe fallback
anonymous. `requireAuthentication` should consume this validated request state; templates may use it
for navigation, but navigation remains presentation rather than authorization. The middleware chain
must stay limited to dynamic routes that need session state, with `/static/*` and `/healthz` excluded.

### Testable invariants

1. A valid session whose user was deleted is treated as anonymous on the next request.
2. The database user lookup runs once per request boundary, not once per helper/template call.
3. Missing, malformed, or wrong-type context values fail closed and never authorize a request.
4. Context carries only request-scoped identity/authentication state; secrets and long-lived dependencies
   remain outside it.
5. The protected `/me` route and navigation agree with the same validated current-user state.

## Next study step

The contract and *Let's Go* Chapters 8, 10, and 11 are studied and aligned. The v0.2.5 Auth
Foundation and its HTTP boundaries are applied at `book-social` commit `41a8ddb`. The next applied
step is v0.2.6: wire registration/login/logout/`GET /me`, session lifecycle, forms, errors, flashes,
and auth-aware navigation without claiming that this flow already exists.
