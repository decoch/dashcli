# Contributing to dashcli

## Development setup

Requirements: Go 1.25+

```bash
git clone https://github.com/decoch/dashcli.git
cd dashcli
go build -o dashcli ./cmd/dashcli
```

Available make targets:

| Command | Description |
|---|---|
| `make test` | Run all tests |
| `make vet` | Run `go vet` |
| `make lint` | Run golangci-lint |
| `make fmt` | Format all Go files |
| `make fmt-check` | Check formatting without modifying files |
| `make ci` | Run fmt-check, vet, lint, and test |

Run `make ci` before submitting a PR.

## Code style

- Lint rules are defined in `.golangci.yml` (govet, staticcheck, errcheck, gosec, and others)
- Format with `gofmt` (enforced by CI)
- Keep changes focused — fix one thing per PR

## Testing

- Tests live alongside their package (`package_test.go`)
- Use the standard library `testing` package only — no external test libraries
- Add tests for any behavior change
- For API-layer tests, use `httptest.Server` with a round-trip stub (see `internal/redash/client_test.go`)
- For command-layer tests, stub package-level function variables (see `internal/app/cmd_auth_test.go`)

## Submitting a PR

1. Fork and create a branch
2. Make your changes
3. Run `make ci`
4. Open a PR against `main`

See the [Contributing section in the README](README.md#contributing) for the short version.
