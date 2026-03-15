# example-project

An example Go kick-starter project for AI-enhanced development.

This repository demonstrates how to start a service with:

- a traditional Go project layout (`cmd`, `internal`, `api`, `configs`, `scripts`, `test`)
- a working stdlib HTTP endpoint (`GET /hello`)
- OpenAPI-first contract placement in `api/openapi/`
- explicit agent guardrails in `AGENTS.md` (tests, formatting, and contract updates)
- reusable AI skills in `.opencode/skill/` for testing workflows

The goal is to give you a clean baseline that is immediately runnable, testable, and friendly for both human and AI contributors.

## Quick start

```bash
go run ./cmd/app
```

Then open:

- `http://localhost:8080/hello`

## Development workflow

```bash
go fmt ./...
go test ./...
```

Or with make:

```bash
make fmt
make test
make check
```

## Project structure

- `cmd/app/`: service entrypoint and startup wiring
- `internal/app/`: core application/domain logic
- `internal/httpserver/`: HTTP routing and handlers (stdlib `net/http`)
- `api/openapi/`: OpenAPI contracts for API behavior
- `.opencode/skill/`: project-specific skills for AI-assisted workflows

## API contract

An example OpenAPI contract for the hello endpoint is available at:

- `api/openapi/hello.yaml`
