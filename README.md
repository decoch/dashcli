# dashcli

`dashcli` is a Go-based CLI for Redash.

It is designed for script-friendly automation, clear command behavior, and secure API key handling via OS keyring.

## Features

- Redash API operations from terminal commands
- JSON output mode for automation (`--json`)
- Profile-based environment switching (`--profile`)
- API key storage in OS keyring (`dash auth`)
- Clear precedence rules for credentials and configuration
- Small, testable Go codebase

## Command Surface

- `dash version`
- `dash auth set|delete|status`
- `dash me`
- `dash query list|get|run`
- `dash job get|wait`
- `dash dashboard list|get`
- `dash datasource list`

Run `dash --help` for full command and flag documentation.

## Installation

### Build from source

Requirements:

- Go 1.22+

```bash
git clone https://github.com/decoch/dashcli.git
cd dashcli
go build -o dash ./cmd/dash
./dash version
```

### Install via `go install`

```bash
go install github.com/decoch/dashcli/cmd/dash@latest
dash version
```

## Quick Start

1. Set base URL and API key (recommended via keyring):

```bash
echo "<YOUR_REDASH_API_KEY>" | dash auth set --profile prod
```

2. Call Redash:

```bash
dash --profile prod --base-url https://redash.example.com me
dash --profile prod --base-url https://redash.example.com query list
```

3. Use JSON output for scripts:

```bash
dash --json --profile prod --base-url https://redash.example.com datasource list
```

## Authentication and Configuration

### API key precedence

`dashcli` resolves API key in this order:

1. `--api-key`
2. keyring secret for selected profile
3. profile `api_key_env`
4. `REDASH_API_KEY`

### Base URL precedence

1. `--base-url`
2. `REDASH_BASE_URL`
3. selected profile `base_url`

### Profile selection

1. `--profile`
2. `default_profile` in config
3. `default` (when profiles exist)

### Config file

- Path: `$(os.UserConfigDir())/dashcli/config.json`

Example:

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

## Security Notes

- Prefer `dash auth set` (keyring) over `--api-key` for day-to-day use.
- `--api-key` prints a warning because command-line args can leak via shell history/process list.
- `http://` base URLs are rejected to prevent unencrypted API key transmission.

## Development

```bash
make fmt
make lint
make test
make ci
```

CI runs on pushes and pull requests to `main`.

## Contributing

Issues and pull requests are welcome.

- Keep changes focused and minimal.
- Add/update tests for behavior changes.
- Run `make ci` before submitting.

## License

MIT — see `LICENSE`.

## Docs

- `docs/spec.md`
- `docs/architecture.md`
- `docs/libraries.md`
