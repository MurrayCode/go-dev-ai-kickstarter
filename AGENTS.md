# AGENTS

This repository follows AI-driven development with strict verification on every change.

## Mandatory rules for every agent change

1. For all application code changes, add or update relevant unit tests.

2. For application code changes that touch external dependencies, persistence, networking, or cross-component behavior, add or update integration tests where necessary.

3. For every API contract change, add or update the relevant OpenAPI contract in `api/openapi/`.

4. Run formatting after making edits:

```bash
go fmt ./...
```

5. Run all tests before finishing:

```bash
go test ./...
```

6. If formatting or tests fail, fix issues and rerun both commands until they pass.

7. Do not finalize work when the repository is in a failing state.

## Recommended command sequence

```bash
go fmt ./...
go test ./...
```

Or use:

```bash
make check
```

## Standard Go project structure for this repository

- `cmd/`: Application entrypoints. Each subdirectory builds one binary.
  - `cmd/app/`: Main executable for this service (`main.go`).
- `internal/`: Private application code that should not be imported by external modules.
  - `internal/app/`: Core domain/application logic.
  - `internal/httpserver/`: HTTP transport layer, routing, and handlers.
- `pkg/`: Public reusable libraries intended for external consumption (use only when APIs are intentionally public).
- `api/`: API contracts and interface definitions (for example OpenAPI specs, protobuf files, JSON schemas), not runtime server handler implementation.
- `configs/`: Versioned configuration templates and examples.
- `scripts/`: Developer and CI helper scripts.
- `test/`: Cross-package test assets, end-to-end/integration harnesses, and shared test data.
- `.opencode/skill/`: Project-specific AI skills and workflows used by agents.

## Placement guidance

1. Put executable startup code in `cmd/<binary-name>/main.go`.
2. Put runtime application code in `internal/...` by default.
3. Put handler/server code in `internal/httpserver` unless there is a strong reason to split by feature.
4. Put API specifications in `api/`.
5. Move code to `pkg/` only when you intentionally want a stable, public import path.
