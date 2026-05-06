# gohive

A multi-service Go monorepo with shared core libraries and a unified maintenance CLI.

## Layout

| Path | Purpose |
| --- | --- |
| `core/` | Shared libraries (config, logger, errors, middleware, api helpers) |
| `pkg/` | Infrastructure adapters (mysql, postgres, redis, kafka) |
| `models/` | Domain types: `entity` (hand-written), `dto`, and generated `query/model` |
| `migrations/` | Goose-managed SQL migrations |
| `scripts/` | Maintenance CLI (`scripts/cmd`) plus per-task packages (`migrate`, `seed`, `fixorder`, `gen`) |
| `demo-api/`, `demo-grpc/`, `demo-ws/`, `demo-worker-*/` | Service entry points |

## Common tasks

All tasks are wrapped in [`justfile`](justfile); run `just` to see the list.

```sh
# Build / run
just build-all
just run demo-api
just run-cfg demo-api ./demo-api/config.development.toml

# Database migrations (goose)
just migrate-up
just migrate-down
just migrate-status
just migrate-create add_user_avatar      # sql migration
just migrate-create backfill_orders go   # go migration

# Seed / one-off fixes
just script seed
just script fix-order --dry-run

# Generate query + model code from the live DB schema
just gen
```

The maintenance CLI itself:

```sh
go run ./scripts/cmd --help
go run ./scripts/cmd migrate up -c ./demo-api/config.toml
```

## Configuration

Each service has its own `config.toml` (+ `config.<env>.toml` overrides), loaded by
[`core/config`](core/config/config.go) via Viper. Environment variables prefixed `APP_` override
file values (e.g. `APP_DATABASE_HOST=db.local`).

Maintenance scripts default to `./demo-api/config.toml` for the `[database]` and `[log]` sections;
override with `-c <path>`.

## Code generation vs hand-written entities

`./models/entity` holds hand-written GORM entities used by service code. `just gen` regenerates
[gorm.io/gen](https://gorm.io/gen) query interfaces and model structs into `./models/query` and
`./models/query/model` — a separate package, so it never collides with `entity`. Pick the layer
that fits the call site; do not mix them in the same query.
