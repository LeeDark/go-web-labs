# Stage 5 Test Matrix

This matrix records why each focused Stage 5 check exists, the nearest useful test boundary, and the
evidence currently available. It complements `README.md`; it is not a coverage target or a shared
test suite.

Status meanings:

- **Verified** — the relevant ordinary Go test suite has run successfully.
- **Opt-in pending** — the test exists and skips safely, but has not run against PostgreSQL because
  `GREENLIGHT_TEST_DB_DSN` was not configured.

| ID | Area        | Risk / behavior                                                                                                       | Boundary                  | Evidence                                                                                                                     | Status         |
|---:|-------------|-----------------------------------------------------------------------------------------------------------------------|---------------------------|------------------------------------------------------------------------------------------------------------------------------|----------------|
|  1 | Layered API | `List` must preserve an unexpected repository error.                                                                  | Service unit              | `TestServiceListPropagatesRepositoryError`                                                                                   | Verified       |
|  2 | Layered API | `Get` must preserve an unexpected repository error.                                                                   | Service unit              | `TestServiceGet/propagates_repository_error`                                                                                 | Verified       |
|  3 | Layered API | `Get` must map an unexpected service error to a safe `500` JSON response.                                             | Handler                   | `TestBooksHandlerGet/internal_error`                                                                                         | Verified       |
|  4 | Layered API | Invalid IDs for `Update` and `Delete` must return `400` without calling the service.                                  | Handler with fake service | `TestBooksHandlerUpdateRejectsInvalidIDWithoutCallingService`, `TestBooksHandlerDeleteRejectsInvalidIDWithoutCallingService` | Verified       |
|  5 | API Core    | Movie validation must cover valid input and required, length, year, runtime, and genre boundaries.                    | Unit                      | `TestValidateMovie`                                                                                                          | Verified       |
|  6 | API Core    | Malformed, unknown-field, empty, and multi-value JSON must follow the documented bad-request path.                    | Handler/helper            | `TestCreateMovieHandlerRejectsInvalidJSON`                                                                                   | Verified       |
|  7 | API Core    | Validation failure must return a `422` JSON error before database access.                                             | Handler                   | `TestCreateMovieHandlerReturnsValidationErrorBeforeDatabaseAccess`                                                           | Verified       |
|  8 | API Core    | Migrations, generated fields, constraints, and PostgreSQL genre arrays must agree with `MovieModel.Insert` and `Get`. | Database integration      | `TestMovieModelInsertAndGetWithPostgres`                                                                                     | Opt-in pending |
|  9 | API Core    | A genuinely stale movie version must map to `ErrEditConflict`.                                                        | Database integration      | `TestMovieModelInsertAndGetWithPostgres`                                                                                     | Opt-in pending |

## Audited areas without new Stage 5 checks

- REST API already covers CRUD, strict JSON, validation, method/path errors, safe panic handling,
  `Location`, and empty `204` responses at the route boundary.
- Snippetbox already demonstrates MPA handler tests and a historical MySQL-backed model-test
  pattern. Stage 5 documents that boundary rather than duplicating it.
- Layered API already has in-memory repository, middleware, router, and composition-root tests.

## Residual evidence gap

Checks 8 and 9 are implemented but are not runtime-verified until the targeted command is run
against an explicitly confirmed disposable PostgreSQL database. A skipped integration test is not
recorded as passed database evidence.
