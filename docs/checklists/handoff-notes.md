# Backend handoff notes checklist

Use this checklist to make a small Go backend change runnable and reviewable by another developer
or client.

## What changed

- [ ] State the user/operator outcome in one or two sentences.
- [ ] List changed endpoints, packages, migrations, and configuration keys.
- [ ] Link the exact repository branch, tag, or commit.
- [ ] Separate implemented behavior from planned or deferred behavior.

## Run and verify

- [ ] Document prerequisites and the working directory for each command.
- [ ] Document local run, test, format, vet, and migration commands.
- [ ] Include one successful request and the most important failure examples.
- [ ] State whether PostgreSQL, Docker, external services, or special environment variables are
      required.
- [ ] State whether integration tests are opt-in and how the disposable test data is owned and
      cleaned up.

## Risks and operations

- [ ] List schema, transaction, timeout, concurrency, security, and compatibility risks relevant to
      the change.
- [ ] Document safe cleanup and rollback steps; do not provide destructive commands without a
      validated target.
- [ ] Record known limitations and the next recommended step.
- [ ] Confirm that logs and examples contain no credentials, tokens, or private data.

## Evidence

- [ ] Add the source/receiving paths and immutable refs to the evidence matrix.
- [ ] Record exact verification results and any skipped environment-dependent checks.
- [ ] Confirm the README and API contract describe the current behavior.
