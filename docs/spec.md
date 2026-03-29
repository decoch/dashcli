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

## Language/runtime

- Go `1.25+` (see `go.mod`)

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

### Commands (first slice)

- `dash version`
- `dash auth set`
- `dash auth delete`
- `dash auth status`
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
- Prefer keyring-backed secrets over command-line flags for local development
- Reject `http://` base URLs to prevent plaintext API key transport
- Default product mode is single-instance (no profile switching)

Credential/base URL resolution rules:

1. Base URL: `--base-url` > keyring `base_url` > `REDASH_BASE_URL`
2. API key: `--api-key` > keyring `api_key` > `REDASH_API_KEY`

If `--api-key` is used and wins resolution, print a warning on stderr about insecure CLI argument usage.

`dash auth set` stores both values in keyring service `dashcli`.

## Command behavior notes

- `dash version`: prints CLI binary version only and exits; it does not call Redash APIs.
- `dash auth set`: prompts for base URL and API key, then stores both in OS keyring.
- `dash auth delete`: removes stored base URL and API key from OS keyring.
- `dash auth status`: prints whether base URL and API key are stored; key-not-found is a normal success state.
- `dash me`: calls the current-user endpoint and prints a compact user summary in text mode (`id`, `name`, `email`, `is_admin`); `--json` prints the API response payload.

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
