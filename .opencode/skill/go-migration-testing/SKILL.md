---
name: go-migration-testing
description: |-
  Design and implement Go migration tests that verify schema creation, data backfills, constraints,
  and rollback safety using real databases (commonly with testcontainers-go). Use for validating SQL
  migration chains, preventing drift between app expectations and schema, and catching destructive
  migration behavior before release. Use proactively when adding/changing migrations or when production
  incidents were caused by schema incompatibility.

  Examples:
  - user: "Test these new SQL migrations" -> run up/down in a disposable DB and assert schema + data outcomes
  - user: "Validate rollback behavior" -> apply migration, insert representative data, run down, verify integrity
  - user: "Catch migration drift in CI" -> add integration tests that assert expected tables, indexes, and constraints
  - user: "Backfill should be safe" -> test idempotency and correctness of data transformation steps
---
# Go Migration Testing

Use this workflow for migration verification in integration tests.

## Goals

- Prove migrations are correct, repeatable, and safe.
- Verify both structural and data-level effects.
- Detect non-idempotent or destructive behavior early.
- Keep migration tests deterministic for CI.

## What to Validate

- Up migration creates/changes expected schema objects.
- Constraints and indexes exist and behave as intended.
- Data backfills transform rows correctly.
- Down migration works when supported by team policy.
- Re-running migrations is safe or fails with expected errors.

## Recommended Test Structure

1. Start disposable database fixture (often testcontainers-go).
2. Apply baseline migrations to known version.
3. Seed representative pre-migration data.
4. Apply target migration(s).
5. Assert schema and data expectations.
6. Optionally run down migration and assert rollback state.
7. Cleanup via `t.Cleanup`.

## Example Pattern

```go
func TestMigrations_UserEmailIndex(t *testing.T) {
	ctx := context.Background()
	fx := testutil.StartPostgres(t)

	if err := migrateTo(ctx, fx.DSN, "202603150900_base"); err != nil {
		t.Fatalf("migrate baseline: %v", err)
	}

	if _, err := fx.DB.ExecContext(ctx, `
		INSERT INTO users (id, email) VALUES
		(1, 'a@example.com'),
		(2, 'b@example.com')
	`); err != nil {
		t.Fatalf("seed users: %v", err)
	}

	if err := migrateTo(ctx, fx.DSN, "202603151000_add_users_email_idx"); err != nil {
		t.Fatalf("migrate target: %v", err)
	}

	if !indexExists(ctx, t, fx.DB, "users", "users_email_idx") {
		t.Fatalf("expected index users_email_idx to exist")
	}
}
```

## Conventions

- Use timestamped migration IDs in assertions to avoid ambiguity.
- Keep migration tests focused: one behavior per test when possible.
- Use helper functions for schema checks (`tableExists`, `indexExists`, `constraintExists`).
- Validate backfills with realistic sample rows and edge values.
- Make rollback tests explicit only if the project supports down migrations.

## Safety Guardrails

- Pin database image versions in fixtures.
- Avoid global shared DB state across tests unless reset is guaranteed.
- Avoid sleep-based waits; use readiness checks from fixture helpers.
- Fail with actionable messages including migration ID and object name.

## CI Strategy

- Run migration tests in the default integration suite or `make check` path.
- Keep long-running migration tests in a separate target when needed.
- Surface SQL and migration-tool errors directly in test failure output.

## Quality Checks

- Ensure each new migration has a corresponding migration test.
- Ensure assertions cover both schema and data when applicable.
- Ensure tests pass from a clean database state.
- Run:

```bash
go fmt ./...
go test ./...
```
