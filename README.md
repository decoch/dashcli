# dashcli

`dashcli` is a Go-based CLI for Redash.

It is designed for script-friendly automation, clear command behavior, and secure API key handling via OS keyring.

## Features

- Redash API operations from terminal commands
- JSON output mode for automation (`--json`)
- Single-instance default workflow (no profile management)
- Base URL and API key storage in OS keyring (`dashcli auth`)
- Clear precedence rules for credentials and configuration
- Small, testable Go codebase

## Command Surface

- `dashcli version`
- `dashcli auth set|delete|status`
- `dashcli query list|get|run`
- `dashcli query create`
- `dashcli query update <id>`
- `dashcli query archive <id>`
- `dashcli query results <id>`
- `dashcli job get|wait`
- `dashcli dashboard list|get`
- `dashcli datasource list`
- `dashcli datasource schema <id>`
- `dashcli sql run`

Run `dashcli --help` for full command and flag documentation.

`dashcli auth ...` commands work without `--base-url` because they manage local keyring state.

## Installation

### Build from source

Requirements:

- Go 1.25+

```bash
git clone https://github.com/decoch/dashcli.git
cd dashcli
go build -o dashcli ./cmd/dashcli
./dashcli version
```

### Install via `go install`

```bash
go install github.com/decoch/dashcli/cmd/dashcli@latest
dashcli version
```

### Download binary (recommended)

Download the latest release from GitHub Releases:
https://github.com/decoch/dashcli/releases

macOS (Apple Silicon):

```bash
curl -L https://github.com/decoch/dashcli/releases/latest/download/dashcli_Darwin_arm64.tar.gz | tar xz
sudo mv dashcli /usr/local/bin/
```

macOS (Intel):

```bash
curl -L https://github.com/decoch/dashcli/releases/latest/download/dashcli_Darwin_amd64.tar.gz | tar xz
sudo mv dashcli /usr/local/bin/
```

Linux (amd64):

```bash
curl -L https://github.com/decoch/dashcli/releases/latest/download/dashcli_Linux_amd64.tar.gz | tar xz
sudo mv dashcli /usr/local/bin/
```

## Quick Start

```bash
# 1. Store base URL and API key in keyring
dashcli auth set
# You will be prompted for:
# - Base URL (e.g. https://your-redash.example.com)
# - API key

# 2. Use
dashcli query list
dashcli --json datasource list
```

## Authentication and Configuration

### API key precedence

`dashcli` resolves API key in this order:

1. `--api-key`
2. keyring `api_key`
3. `REDASH_API_KEY`

### Base URL precedence

1. `--base-url`
2. keyring `base_url`
3. `REDASH_BASE_URL`

### Keyring entries

`dashcli auth set` stores two entries in keyring service `dashcli`:

- `base_url`
- `api_key`

`dashcli auth status` checks whether both values are present.

`dashcli auth delete` removes both values.

## Security Notes

- Prefer `dashcli auth set` (keyring) over `--api-key` for day-to-day use.
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
