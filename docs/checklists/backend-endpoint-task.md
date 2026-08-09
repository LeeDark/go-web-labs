# Backend endpoint task checklist

Use this checklist for a small Go backend fix or endpoint addition. It is a review and handoff aid,
not a replacement for the receiving project's plan or acceptance criteria.

## Scope and contract

- [ ] State the resource, user/operator outcome, and explicit out-of-scope work.
- [ ] Identify the method, path, path/query/body inputs, and response shape.
- [ ] Define success status, headers, and error statuses before implementation.
- [ ] Define not-found, validation, conflict, unsupported-method, and internal-error behavior where
      applicable.
- [ ] Record ownership, authentication, authorization, and privacy rules if the resource is private.

## Implementation boundaries

- [ ] Keep request parsing and HTTP status mapping in the handler layer.
- [ ] Keep domain/use-case rules in the service layer when a service boundary exists.
- [ ] Keep SQL/driver-specific behavior in the repository or storage layer.
- [ ] Add interfaces from the consuming package only when a real seam is required.
- [ ] Keep errors safe for clients and detailed enough for server logs without leaking secrets.

## Validation and data

- [ ] Reject malformed, oversized, unknown, or trailing JSON values as required by the contract.
- [ ] Validate path and query parameters explicitly.
- [ ] Validate domain fields separately from JSON decoding.
- [ ] Check migration, constraint, transaction, timeout, and cancellation effects when persistence
      is involved.
- [ ] Use a disposable database for migration or integration checks; never run destructive cleanup
      against a development or production database.

## Tests

- [ ] Add or update focused unit tests for domain and service rules.
- [ ] Add `httptest` coverage for success, validation, not-found, conflict, and internal-error paths
      relevant to the change.
- [ ] Add repository/integration coverage only where database behavior is part of the risk.
- [ ] Verify repeated execution, stale-version behavior, and cleanup for stateful tests.
- [ ] Run the narrowest relevant check first, then the module's full test suite when justified.

## Verification and handoff

- [ ] Run format, tests, and vet (or the receiving project's equivalent commands).
- [ ] Record exact commands, repository ref, and result in the change notes.
- [ ] Update API examples or contract documentation.
- [ ] Record migrations, configuration, environment variables, and external dependencies.
- [ ] List known limitations, deferred improvements, and rollback/cleanup notes.
- [ ] Confirm the working tree and generated artifacts are in the expected state.
