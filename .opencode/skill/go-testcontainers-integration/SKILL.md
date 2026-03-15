---
name: go-testcontainers-integration
description: |-
  Build and maintain Go integration tests using testcontainers-go with deterministic setup,
  readiness checks, isolated test data, and reliable cleanup. Use for new integration test suites,
  migrating flaky external-service tests, or adding database/queue/cache coverage in Go services.
  Use proactively when tests depend on Postgres, Redis, Kafka, or any service that should run in containers.

  Examples:
  - user: "Add integration tests for our repository layer" -> start a Postgres container, run migrations, verify CRUD paths
  - user: "Our Redis tests are flaky" -> add readiness strategy, per-test isolation, and cleanup hooks
  - user: "Test against real dependencies locally and in CI" -> wire container lifecycle to go test and CI-friendly defaults
  - user: "Move docker-compose tests to Go" -> replace ad-hoc scripts with testcontainers-go fixtures
---
# Go Integration Testing with testcontainers-go

Use this workflow when authoring `*_integration_test.go` or package-local integration tests that require real dependencies.

## Goals

- Keep tests deterministic and reproducible on local machines and CI.
- Minimize flakiness via explicit wait strategies and timeouts.
- Isolate data between tests.
- Guarantee cleanup even when assertions fail.

## Workflow

1. Define integration test scope (single dependency vs multi-service).
2. Start containers with explicit image tags and wait conditions.
3. Build runtime config from mapped host/port values.
4. Run schema setup or seed data in test setup.
5. Execute behavior-focused test cases.
6. Terminate containers and remove temporary resources.

## Recommended Pattern

```go
package repository_test

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

func TestUserRepository_Integration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "app",
			"POSTGRES_USER":     "app",
			"POSTGRES_PASSWORD": "app",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * time.Second),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() {
		_ = postgres.Terminate(context.Background())
	})

	host, err := postgres.Host(ctx)
	if err != nil {
		t.Fatalf("postgres host: %v", err)
	}

	port, err := postgres.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("postgres mapped port: %v", err)
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

	// Run migrations/fixtures, then execute integration assertions.
}
```

## Conventions

- Prefer pinned image tags (`postgres:16-alpine`), not `latest`.
- Use `context.WithTimeout` to bound startup and test duration.
- Use `t.Cleanup` for container termination and resource cleanup.
- Prefer package-level helper fixtures for repeated container setup.
- Keep integration tests behind explicit naming or build tags when suites are heavy.

## Flakiness Guardrails

- Use readiness checks (`wait.ForLog`, `wait.ForListeningPort`, health checks).
- Avoid arbitrary sleeps.
- Ensure each test owns its data (unique schema/table prefixes or full reset).
- Keep assertions tolerant only where eventual consistency is expected.

## CI Notes

- Ensure Docker is available in CI runner.
- Keep startup timeouts realistic for shared runners.
- Keep test logs useful (`t.Logf`) for container startup/debug failures.
- Split long-running suites from fast unit tests when needed.

## Quality Checks

- Verify container startup failure messages are actionable.
- Verify cleanup always runs.
- Verify test data isolation between cases.
- Run:

```bash
go fmt ./...
go test ./...
```
