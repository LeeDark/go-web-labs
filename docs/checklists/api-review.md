# API review checklist

Use this checklist when reviewing a small REST or JSON API change.

## Contract

- [ ] Method and URL follow the resource model and existing versioning policy.
- [ ] Request and response JSON shapes are documented with examples.
- [ ] Success status, `Content-Type`, `Location`, caching, and other relevant headers are explicit.
- [ ] Error envelope, stable error code, message, and field details are consistent with the API.
- [ ] Unknown routes and unsupported methods have predictable responses.

## Boundary behavior

- [ ] Path, query, and body inputs have explicit parsing and size limits.
- [ ] Unknown fields, malformed JSON, empty bodies, and trailing values are handled intentionally.
- [ ] Validation is distinct from decoding and maps to the documented status.
- [ ] Internal errors do not expose SQL, DSNs, stack traces, or secrets.
- [ ] Context cancellation and request timeouts reach database or external calls.

## Data and security

- [ ] Authorization is enforced in the use-case/service boundary, not only in rendering.
- [ ] Ownership, privacy, and conflict behavior are tested where relevant.
- [ ] Transactions and optimistic concurrency are explicit for multi-step writes.
- [ ] Migrations are ordered, reversible where supported, and tested on a disposable database.
- [ ] Logs contain useful operation context without credentials, tokens, or private content.

## Evidence

- [ ] Unit and HTTP-contract tests cover the changed behavior and important failures.
- [ ] Database integration is included only when it verifies a real persistence risk.
- [ ] The exact source and receiving refs are recorded in the evidence matrix.
- [ ] Verification commands pass and known limitations are written down.
