# Go Stock Engine

## Purpose

This is a Go API for managing product stock and ranking replenishment
priorities. Keep the priority calculation explainable and deterministic.
Do not change its formula or the meaning of its policy factors without an
explicit business requirement.

## Architecture

- `cmd/`: application composition, environment configuration, and startup.
- `internal/domain/`: business entities, validation, errors, pagination, and
  replenishment-priority rules. It must not depend on HTTP, GORM, or PostgreSQL.
- `internal/application/`: use cases that coordinate domain objects and the
  repository interface.
- `internal/domain/repository/`: persistence contracts used by use cases.
- `internal/infrastructure/`: GORM/PostgreSQL implementation and seed data.
- `internal/presentation/http/`: Gin routes, request parsing, and HTTP error mapping.
- `docs/`: generated Swagger artifacts; regenerate them when the HTTP contract changes.

Keep dependencies flowing inward: presentation and infrastructure depend on
application/domain, never the reverse.

## Domain and API rules

- Create and update product stock through the domain validation paths; preserve
  `domain.Error` and its error codes through the application layer.
- Keep HTTP status mapping centralized in `internal/presentation/http/handlers.go`.
- `PriorityPolicy` is configured through environment variables in `cmd`; do
  not hard-code business factors in handlers or repository code.
- Preserve pagination validation and the existing `/restock/priorities` API
  contract unless a requested change explicitly includes it.

## Development conventions

- Use idiomatic Go, `gofmt`, and focused changes. Avoid new dependencies unless
  they are necessary for the requested behavior.
- Add domain tests next to domain code; add use-case tests under
  `internal/application/tests` using the existing repository mock pattern.
- Add or update E2E tests in `tests/e2e` only for HTTP-level behavior.
- Update Swagger comments and regenerate `docs/` for route, request, response,
  or parameter changes:

  ```bash
  swag init -g cmd/app/main.go -o docs --parseInternal
  ```

## Verification

Run the narrowest relevant check first:

```bash
go test ./internal/domain
go test ./internal/application/tests
```

E2E tests require the API running with PostgreSQL:

```bash
docker compose up -d
go test ./tests/e2e/ -v -count=1
```

Run the complete suite before handing off a cross-layer change:

```bash
go test ./...
```

For local execution, configure the environment described in `README.md`, then
run `go run ./cmd/app`; Docker uses `docker compose up --build`.
