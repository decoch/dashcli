# dashcli architecture

## 1. Package layout

```text
cmd/dashcli/
  main.go                 # thin entrypoint

internal/app/
  run.go                  # Run(ctx, args, stdout, stderr) int
  root.go                 # root command + global flags
  cmd_version.go          # version subcommand
  cmd_auth.go             # keyring auth commands
  cmd_query.go            # query subcommands
  cmd_dashboard.go        # dashboard subcommands
  cmd_datasource.go       # datasource subcommands
  cmd_job.go              # job subcommands
  cmd_sql.go              # ad-hoc SQL execution
  helpers.go              # shared output helpers

internal/redash/
  client.go               # HTTP client + auth transport
  queries.go              # /api/queries
  dashboards.go           # /api/dashboards
  jobs.go                 # /api/jobs
  datasources.go          # /api/data_sources
  errors.go               # API error mapping

internal/secrets/
  keyring.go              # OS keyring wrapper

internal/output/
  output.go               # text/json renderer

internal/exitcode/
  exitcode.go             # code mapping
```

## 2. Execution flow

1. `cmd/dashcli/main.go` calls `app.Run(...)`
2. `internal/app` parses flags and resolves credentials (`flag > keyring > env`)
3. `internal/secrets` provides keyring-backed `base_url` and `api_key`
4. Build `internal/redash/client` (timeout, auth header)
5. Execute the subcommand (`auth` commands bypass API/client requirements)
6. `internal/output` renders text/json
7. `internal/exitcode` maps and returns the exit code

## 3. Design principles

- **Thin entrypoint**: `cmd/dashcli/main.go` only calls `app.Run(...)`. All logic lives under `internal/`.
- **Testable runner**: `Run(ctx, args, stdout, stderr) int` makes the full CLI testable without subprocess invocation.
- **Layered responsibility**: `internal/app` owns CLI concerns; `internal/redash` owns API concerns. Neither bleeds into the other.
- **Explicit credential resolution**: flags → keyring → env, evaluated in order, with a warning when the insecure flag path wins.
- **Machine-readable output as first-class**: every command supports `--json`; the shape is stable and script-friendly.
- **Stable exit codes**: errors are classified at the boundary and mapped to fixed codes (0/1/2) before returning to the OS.
- **Output isolation**: all rendering goes through `internal/output` so format behavior is consistent and easy to test.

## 4. Error model

- Usage error (missing required args, invalid flags)
- API error (4xx/5xx + Redash error payload)
- Network error (timeout, DNS, TLS)

Normalize these into internal error types and map them to exit codes.

## 5. Extensibility policy

- Add new APIs via `internal/redash/<resource>.go`
- Add new command handlers via `internal/app/cmd_<resource>.go`
- Limit output format changes to `internal/output`
