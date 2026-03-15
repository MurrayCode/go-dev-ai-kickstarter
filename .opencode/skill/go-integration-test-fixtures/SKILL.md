---
name: go-integration-test-fixtures
description: |-
  Create reusable Go integration-test fixture packages (for example `internal/testutil`) that centralize
  testcontainers-go setup, runtime config wiring, schema/bootstrap steps, and teardown behavior.
  Use for reducing duplicated container boot code, standardizing test environments, and making integration
  tests easier to read and maintain across packages. Use proactively when multiple integration tests repeat
  similar setup logic for databases, caches, queues, or external services.

  Examples:
  - user: "Our integration tests all start Postgres differently" -> extract a shared `internal/testutil` fixture with one setup path
  - user: "Reduce duplicate Redis setup in tests" -> build reusable fixture constructor and cleanup hook
  - user: "Make integration tests cleaner" -> move container + migration setup into helper APIs used by tests
  - user: "Standardize test env config" -> return typed fixture config (DSN, host, ports, clients) from one helper
---
# Go Integration Test Fixtures

Use this workflow to create maintainable fixture helpers for integration tests.

## Goals

- Keep test files focused on behavior, not infrastructure bootstrapping.
- Reuse a consistent setup/teardown pattern across packages.
- Make fixture APIs typed, small, and explicit.
- Preserve isolation and determinism.

## Recommended Layout

- `internal/testutil/` for shared fixture helpers.
- `internal/testutil/postgres.go`, `redis.go`, etc. per dependency.
- `internal/testutil/migrations.go` for schema/bootstrap helpers.
- Keep production code free from test-only dependencies.

## Fixture API Pattern

```go
package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresFixture struct {
	DB  *sql.DB
	DSN string
}

func StartPostgres(t *testing.T) *PostgresFixture {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "app",
				"POSTGRES_USER":     "app",
				"POSTGRES_PASSWORD": "app",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(context.Background()) })

	host, err := ctr.Host(ctx)
	if err != nil {
		t.Fatalf("postgres host: %v", err)
	}
	port, err := ctr.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("postgres port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://app:app@%s:%s/app?sslmode=disable", host, port.Port())
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("ping db: %v", err)
	}

	return &PostgresFixture{DB: db, DSN: dsn}
}
```

## Usage Pattern in Tests

1. Start fixture at test start (`fx := testutil.StartPostgres(t)`).
2. Run migrations/bootstrap using `fx.DB`.
3. Build repository/service under test with fixture config.
4. Execute assertions.
5. Rely on `t.Cleanup` for teardown.

## Conventions

- Keep fixture constructors side-effect limited to setup only.
- Return typed fixture structs, not loose maps.
- Use explicit timeouts and wait strategies.
- Prefer per-test fixtures unless suite-level reuse is justified.
- If using suite-level reuse, reset state between tests.

## Anti-Patterns

- Starting containers directly in every test function.
- Hidden global mutable state in fixtures.
- Sleep-based readiness checks.
- Unpinned container image tags.

## Quality Checks

- Ensure fixture functions call `t.Helper()`.
- Ensure cleanup is registered for every external resource.
- Ensure test data does not leak between cases.
- Run:

```bash
go fmt ./...
go test ./...
```
