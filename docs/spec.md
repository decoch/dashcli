# dashcli spec

## Goal

Design a Go-based Redash CLI (`dashcli`) with the following primary goals.

- Provide a stable CLI for calling the Redash API
- Standardize script-friendly JSON output
- Keep the structure highly testable for long-term maintainability

## Non-goals (v0)

- Do not cover every Redash API endpoint from day one
- Do not build an interactive TUI
- Do not over-optimize in ways that change behavior incompatibly

## Language/runtime

- Go `1.25+` (see `go.mod`)

## Architecture patterns

- Keep `cmd/<bin>/main.go` thin (entrypoint only)
- Separate responsibilities under `internal/*`
- Expose a testable entrypoint such as `Run(ctx, args, stdout, stderr)`
- Isolate output concerns in `internal/output`
- Manage common global flags like `--json` / `--timeout` consistently

## v0 Command surface

### Global flags

- `--base-url` (for example: `https://redash.example.com`)
- `--api-key`
- `--json`
- `--timeout`
- `--user-agent` (default: `dashcli`)

### Commands (first slice)

- `dashcli version`
- `dashcli auth set`
- `dashcli auth delete`
- `dashcli auth status`
- `dashcli query list`
- `dashcli query get <id>`
- `dashcli query run <id>`
- `dashcli query create`
- `dashcli query update <id>`
- `dashcli query archive <id>`
- `dashcli query result <id>`
- `dashcli query-result get <id>`
- `dashcli query-result create`
- `dashcli job get <job-id>`
- `dashcli job wait <job-id>`
- `dashcli dashboard list`
- `dashcli dashboard get <slug-or-id>`
- `dashcli datasource list`
- `dashcli datasource schema <id>`

## API/auth and config policy

- Use Redash API keys for authentication
- Prefer keyring-backed secrets over command-line flags for local development
- Reject `http://` base URLs to prevent plaintext API key transport
- Default product mode is single-instance (no profile switching)

Credential/base URL resolution rules:

1. Base URL: `--base-url` > keyring `base_url` > `REDASH_BASE_URL`
2. API key: `--api-key` > keyring `api_key` > `REDASH_API_KEY`

If `--api-key` is used and wins resolution, print a warning on stderr about insecure CLI argument usage.

`dashcli auth set` stores both values in keyring service `dashcli`.

## Command behavior notes

- `dashcli version`: prints CLI binary version only and exits; it does not call Redash APIs.
- `dashcli auth set`: prompts for base URL and API key, then stores both in OS keyring.
- `dashcli auth delete`: removes stored base URL and API key from OS keyring.
- `dashcli auth status`: prints whether base URL and API key are stored; key-not-found is a normal success state.
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
- `internal/secrets`: verify keyring read/write/remove and auth command behavior
- `internal/output`: verify text/JSON rendering behavior and stable machine output

## Project docs

- `docs/architecture.md`: package layout and execution flow
- `docs/libraries.md`: library selection (adopt / optional / reject)
