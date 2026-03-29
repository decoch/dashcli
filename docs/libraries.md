# dashcli library selection

## Selection policy

- Keep dependencies minimal and centered on the standard library
- Adopt only what is needed for API calls and CLI UX
- Prioritize stability and explicitness over convenience abstractions

## Adopt (v0)

### `github.com/spf13/cobra`

- Reason: naturally models resource-oriented commands (query/dashboard/job)
- Reason: strong built-in help and completion support
- Reason: subcommand structure (resource verb pattern) maps naturally to Redash's REST API

### `github.com/99designs/keyring`

- Reason: secure API key storage in OS keychain/keyring
- Reason: avoids leaking long-lived keys via shell history and process args
- Reason: cross-platform (macOS Keychain, Linux Secret Service, Windows Credential Store)

### Standard library (`net/http`, `encoding/json`, `context`, `log/slog`)

- Reason: the Redash API is REST + JSON, which standard packages handle well
- Reason: fewer dependencies reduce maintenance cost
- Reason: minimizes transitive dependencies and keeps the build reproducible
- Note: CLI flag parsing is handled by `cobra`/`pflag`, not Go's `flag` package directly

## Optional (need-based)

### `github.com/stretchr/testify`

- Use case: improve test readability
- Decision: introduce only if standard `testing` becomes too verbose

## Reject for now

### `viper`

- Reason to reject: heavy config layer with more implicit behavior
- Alternative: keep explicit `flag > keyring > env` resolution in `internal/app`

### `resty`

- Reason to reject: `net/http` is sufficient for v0 API requirements
- Alternative: keep thin shared helpers in `internal/redash/client.go`

### Colored UI libraries (for example: `termenv`)

- Reason to reject: v0 is script-first and prioritizes stable JSON/text output
- Alternative: add later behind `internal/output` if needed

## Initial `go.mod` target

```text
require (
  github.com/spf13/cobra v1.10.2
  github.com/99designs/keyring v1.2.2
)
```

Pin an exact stable version in `go.mod` and upgrade intentionally (do not use floating placeholders).
