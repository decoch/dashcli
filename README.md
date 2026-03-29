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

`dash auth ...` commands work without `--base-url` because they operate on local keyring state.

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

### Download binary (recommended)

Download the latest release from GitHub Releases:
https://github.com/decoch/dashcli/releases

macOS (Apple Silicon):

```bash
curl -L https://github.com/decoch/dashcli/releases/latest/download/dashcli_Darwin_arm64.tar.gz | tar xz
sudo mv dash /usr/local/bin/
```

macOS (Intel):

```bash
curl -L https://github.com/decoch/dashcli/releases/latest/download/dashcli_Darwin_amd64.tar.gz | tar xz
sudo mv dash /usr/local/bin/
```

Linux (amd64):

```bash
curl -L https://github.com/decoch/dashcli/releases/latest/download/dashcli_Linux_amd64.tar.gz | tar xz
sudo mv dash /usr/local/bin/
```

## Quick Start

### Quick Start (single Redash instance)

```bash
# 1. Create config file
#    macOS:  ~/Library/Application Support/dashcli/config.json
#    Linux:  ~/.config/dashcli/config.json
{
  "base_url": "https://your-redash.example.com"
}

# 2. Store API key in keyring
dash auth set
# You will be prompted to enter your API key securely.

# 3. Use
dash me
dash query list
dash --json datasource list
```

### Advanced: multiple Redash instances

Use `--profile` to switch between instances (e.g. different companies or teams).

Config:

```json
{
  "default_profile": "projectA",
  "profiles": {
    "projectA": { "base_url": "https://redash.projecta.com" },
    "projectB": { "base_url": "https://redash.projectb.com" }
  }
}
```

Store keys per profile:

```bash
dash auth set --profile projectA
dash auth set --profile projectB
```

Switch:

```bash
dash --profile projectB query list
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

- macOS: `~/Library/Application Support/dashcli/config.json`
- Linux: `~/.config/dashcli/config.json`

Simple (single instance):

```json
{
  "base_url": "https://your-redash.example.com"
}
```

With profiles (multiple instances):

```json
{
  "default_profile": "projectA",
  "profiles": {
    "projectA": { "base_url": "https://redash.projecta.com" },
    "projectB": { "base_url": "https://redash.projectb.com" }
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
