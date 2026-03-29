# dashcli spec

## Goal

Design a Go-based Redash CLI (`dash`) with the following primary goals.

- Provide a stable CLI for calling the Redash API
- Standardize script-friendly JSON output
- Keep the structure highly testable for long-term maintainability

## Non-goals (v0)

- Do not cover every Redash API endpoint from day one
- Do not build an interactive TUI
- Do not over-optimize in ways that change behavior incompatibly

## Reference repositories and extracted patterns

This spec is based on the following repositories.

- `steipete/gogcli`
- `steipete/sonoscli`
- `steipete/blucli`

Extracted architecture patterns:

- Keep `cmd/<bin>/main.go` thin (entrypoint only)
- Separate responsibilities under `internal/*`
- Expose a testable entrypoint such as `Run(ctx, args, stdout, stderr)`
- Isolate output concerns in `internal/output`
- Manage common global flags like `--json` / `--debug` / `--timeout` consistently

## v0 Command surface

### Global flags

- `--base-url` (for example: `https://redash.example.com`)
- `--api-key`
- `--json`
- `--timeout`
- `--debug`
- `--profile` (switch between multiple environments)

### Commands (first slice)

- `dash version`
- `dash me`
- `dash query list`
- `dash query get <id>`
- `dash query run <id>`
- `dash job get <job-id>`
- `dash job wait <job-id>`
- `dash dashboard list`
- `dash dashboard get <slug-or-id>`
- `dash datasource list`

## API/auth and config policy

- Use Redash API keys for authentication
- API key precedence:
  1. `--api-key`
  2. `REDASH_API_KEY`
  3. local config file
- Base URL precedence:
  1. `--base-url`
  2. `REDASH_BASE_URL`
  3. local config file

Config file candidate:

- `$(os.UserConfigDir())/dashcli/config.json`

## Output policy

- Default: human-readable text
- `--json`: machine-readable JSON to stdout
- Emit errors to stderr consistently
- Define exit codes clearly
  - `0`: success
  - `1`: runtime/API error
  - `2`: usage error

## Testing strategy

- `internal/redash`: API contract tests with `httptest.Server`
- `internal/app`: verify `Run(...)` I/O and exit codes
- `internal/config`: verify precedence and profile resolution

## Deliverables for architecture phase

- `docs/architecture.md`: package layout and execution flow
- `docs/libraries.md`: library selection (adopt / optional / reject)
