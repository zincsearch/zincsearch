# AGENTS.md

ZincSearch — lightweight Elasticsearch alternative. Go 1.25 server (`github.com/zincsearch/zincsearch`) with an embedded React/TypeScript/Vite console.

## Layout

- `cmd/zincsearch/main.go` — entrypoint (telemetry, gin server, shutdown); swag annotation root.
- `embed.go` — `//go:embed frontend/dist`; the UI is compiled into the binary.
- `pkg/config` — `config.Global`: embedded `default.toml` → optional `./conf/zinc.toml` → env vars. Field tag `toml:"zinc_server_port"` is both the TOML key and, uppercased, the env var. Durations in ns, sizes in bytes. No `.env`, no `ZINC_CONFIG_FILE`.
- `pkg/routes` — `/api/...` (native) and `/es/...` (Elasticsearch-compatible), auth middleware, pprof (debug only), swagger.
- `pkg/handlers/{auth,index,document,search}` — gin handlers.
- `pkg/core` — indexes, shards, WAL-backed writes, search.
- `pkg/bluge`, `pkg/uquery` — Bluge wrappers, query DSL / aggregations.
- `pkg/metadata` — persistence: bolt (default), badger, etcd.
- `pkg/meta`, `pkg/errors`, `pkg/zutils`, `pkg/wal`, `pkg/ider`, `pkg/upgrade` — types, helpers, WAL, ids, on-disk upgrades.
- `docker/` — Dockerfiles (`.ci` nightly, `.release` goreleaser); `sh/` — scripts, runnable from anywhere.
- `docs/` — generated swagger; `test/api`, `test/benchmark` — HTTP and benchmark tests; `web/` — legacy Vue UI, not built.

## Build / run

Node.js 24 + npm; keep `frontend/package-lock.json` in sync.

```shell
cd frontend && npm ci && npm run build && cd ..   # REQUIRED before any go build/test
./sh/build.sh                                     # frontend + versioned static binary
go build -o zincsearch cmd/zincsearch/main.go
ZINC_FIRST_ADMIN_USER=admin ZINC_FIRST_ADMIN_PASSWORD=... go run cmd/zincsearch/main.go  # :4080
cd frontend && npm run dev                        # UI dev server
./sh/develop.sh                                   # reflex live-reload
```

Version via `-ldflags -X .../pkg/meta.{Version,CommitHash,BuildDate}`.

## Test / lint

```shell
./sh/test.sh [TestFoo|bench]  # go test -v ./... (sets env, cleans data dirs); arg → -test.run, or benchmarks
./sh/coverage.sh              # -race -coverpkg=./...; fails under 70%
./sh/golangci-lint.sh         # golangci-lint v2, .golangci.yaml
./sh/swagger.sh               # regenerate docs/ after swag annotation changes
cd frontend && npm run build  # TS check + assets
```

CI (GitHub Actions): `ci.yml` — lint, vet, build, test on ubuntu/macos/windows + coverage to Codecov; `nightly.yml` — `ghcr.io/zincsearch/zincsearch-dev:nightly`; `release.yml` — goreleaser on `v*` tags (`.goreleaser.yml`).

## Conventions

- Apache-2.0 "Copyright 2022 Zinc Labs Inc. and Contributors" header on every Go file.
- Imports: stdlib, third-party, `github.com/zincsearch/zincsearch/...`.
- `pkg/zutils/json` instead of `encoding/json`; logging via `github.com/rs/zerolog/log`.
- New settings: field in `pkg/config` struct with `toml:"zinc_..."` tag + default in `default.toml`.
- Tests: `*_test.go` beside code, `testify/assert`, table-driven with `t.Run`. Name index fixtures after the test (`TestIndex_Index.index_1`) — indexes are global and collide.
- `test/api` tests share the gin engine from `init.go` (`server()`, `request(method, path, body)`, basic auth `admin`).

## Pitfalls

- `go build`/`go test`/gopls fail with `pattern frontend/dist: no matching files found` until the UI is built; rebuild after frontend changes.
- Tests create `data/` under `pkg/` and `test/`; `sh/test.sh` cleans them, plain `go test` does not.
- First run needs `ZINC_FIRST_ADMIN_USER`/`ZINC_FIRST_ADMIN_PASSWORD`; never commit a `./conf/zinc.toml` with real credentials.
- `ZINC_METADATA_STORAGE=etcd` tests need a running etcd with auth (`test/README.md`).
- Behavior changes usually need both `/api` and `/es` routes updated.
- Coverage gate is 70% — add tests with behavior changes.
