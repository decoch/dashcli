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

Config file candidate:

- `$(os.UserConfigDir())/dashcli/config.json`

### Profile behavior (`--profile`)

Profiles allow switching Redash environments (for example: `prod`, `stg`) without rewriting flags.

Proposed config shape:

```json
{
  "default_profile": "prod",
  "profiles": {
    "prod": {
      "base_url": "https://redash.example.com",
      "api_key_env": "REDASH_API_KEY_PROD"
    },
    "stg": {
      "base_url": "https://redash-stg.example.com",
      "api_key_env": "REDASH_API_KEY_STG"
    }
  }
}
```

Profile resolution rules:

1. Selected profile name: `--profile` > `default_profile` > `default`
2. Base URL: `--base-url` > `REDASH_BASE_URL` > selected profile `base_url`
3. API key: `--api-key` > env var referenced by selected profile `api_key_env` > `REDASH_API_KEY`

This makes profile-specific credentials win by default, while `REDASH_API_KEY` remains a global fallback.

If the selected profile does not exist, exit with usage error (`2`).

## Command behavior notes

- `dash version`: prints CLI binary version only and exits; it does not call Redash APIs.
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
- `internal/config`: verify precedence and profile resolution
- `internal/output`: verify text/JSON rendering behavior and stable machine output

## Deliverables for architecture phase

- `docs/architecture.md`: package layout and execution flow
- `docs/libraries.md`: library selection (adopt / optional / reject)
