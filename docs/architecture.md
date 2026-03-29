# dashcli architecture

## 1. Package layout

```text
cmd/dash/
  main.go                 # thin entrypoint

internal/app/
  run.go                  # Run(ctx, args, stdout, stderr) int
  root.go                 # root command + global flags
  cmd_version.go          # version subcommand
  cmd_auth.go             # keyring auth commands
  cmd_me.go               # current-user subcommand
  cmd_query.go            # query subcommands
  cmd_dashboard.go        # dashboard subcommands
  cmd_datasource.go       # datasource subcommands
  cmd_job.go              # job subcommands

internal/redash/
  client.go               # HTTP client + auth transport
  me.go                   # /api/me
  queries.go              # /api/queries
  dashboards.go           # /api/dashboards
  jobs.go                 # /api/jobs
  datasources.go          # /api/data_sources
  errors.go               # API error mapping

internal/config/
  config.go               # file/env/flag merge
  profile.go              # profile resolution

internal/secrets/
  keyring.go              # OS keyring wrapper

internal/output/
  output.go               # text/json renderer

internal/exitcode/
  exitcode.go             # code mapping
```

## 2. Execution flow

1. `cmd/dash/main.go` calls `app.Run(...)`
2. `internal/app` parses flags and resolves profile
3. `internal/config` resolves profile and credentials (`flag > keyring > env`, with config-backed profile metadata)
4. Build `internal/redash/client` (timeout, auth header)
5. Execute the subcommand (`auth` commands bypass API/client requirements)
6. `internal/output` renders text/json
7. `internal/exitcode` maps and returns the exit code

## 3. Design decisions from references

### From `gogcli`

- Keep global flag design explicit (output mode, verbose, auth-related)
- Treat machine-readable output as first-class
- Keep entrypoint minimal and route behavior through a testable internal runner
- Centralize error formatting and map domain errors to stable exit codes

### From `sonoscli`

- Split command tree by subcommand responsibility
- Keep dependency creation at function boundaries for test-time replacement

### From `blucli`

- Use `Run(ctx, args, stdout, stderr)` for high testability
- Separate `internal/output` to keep output behavior consistent

## 4. Error model

- Usage error (missing required args, invalid flags)
- API error (4xx/5xx + Redash error payload)
- Network error (timeout, DNS, TLS)

Normalize these into internal error types and map them to exit codes.

## 5. Extensibility policy

- Add new APIs via `internal/redash/<resource>.go`
- Add new command handlers via `internal/app/cmd_<resource>.go`
- Limit output format changes to `internal/output`
