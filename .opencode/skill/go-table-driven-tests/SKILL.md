---
name: go-table-driven-tests
description: |-
  Write and refactor Go tests into table-driven form, including subtests, clear case names,
  deterministic fixtures, edge-case coverage, and parallel-safe patterns. Use for adding new
  unit tests, converting repetitive tests, or improving test readability and coverage in Go packages.
  Use proactively when multiple similar assertions appear or when new behavior needs many input/output cases.

  Examples:
  - user: "Add tests for this parser" -> create table-driven tests with named cases and edge inputs
  - user: "These tests are repetitive" -> refactor to a case table with t.Run subtests
  - user: "Cover error paths too" -> expand table with failure cases and explicit expected errors
  - user: "Make tests easier to read" -> normalize setup, expected values, and case naming
---
# Go Table-Driven Tests

Use this workflow for Go test authoring in `*_test.go` files.

## Goals

- Prefer one table for one behavior.
- Keep test case names explicit and behavior-focused.
- Cover happy path, edge cases, and invalid input paths.
- Keep assertions deterministic and stable.

## Workflow

1. Identify the unit under test and expected behavior.
2. Define a local `testCase` struct with only needed fields.
3. Build a `[]testCase` literal with readable `name` values.
4. Loop over cases and run each with `t.Run(tc.name, func(t *testing.T) { ... })`.
5. Compare outputs and errors explicitly, without hidden control flow.
6. Add regression case names when fixing bugs.

## Recommended Pattern

```go
func TestThing(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "empty input", input: "", want: 0, wantErr: true},
		{name: "single token", input: "a", want: 1, wantErr: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Thing(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Thing() error = %v, wantErr %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Fatalf("Thing() = %d, want %d", got, tc.want)
			}
		})
	}
}
```

## Conventions

- Use `testCases` or `cases` consistently in a file.
- Keep expected values inline unless fixtures are large.
- For error checks, prefer boolean expectation or `errors.Is`.
- Use `t.Helper()` in shared assertion helpers.
- Use `t.Parallel()` only when test data and globals are isolated.

## Quality Checks

- Ensure each case name describes intent, not just input.
- Ensure at least one boundary or invalid case exists.
- Ensure failure messages include function name and got/want.
- Run:

```bash
go fmt ./...
go test ./...
```
