# dashcli library selection

## Selection policy

- Keep dependencies minimal and centered on the standard library
- Adopt only what is needed for API calls and CLI UX
- Prioritize patterns proven in `gogcli`, `sonoscli`, and `blucli`

## Adopt (v0)

### `github.com/spf13/cobra`

- Reason: naturally models resource-oriented commands (query/dashboard/job)
- Reason: strong built-in help and completion support
- Reference: proven command-tree usage in `sonoscli`

### Standard library (`net/http`, `encoding/json`, `context`, `log/slog`)

- Reason: the Redash API is REST + JSON, which standard packages handle well
- Reason: fewer dependencies reduce maintenance cost
- Reference: `blucli` runs stably with a lean dependency set
- Note: CLI flag parsing is handled by `cobra`/`pflag`, not Go's `flag` package directly

## Optional (need-based)

### `github.com/99designs/keyring`

- Use case: store API keys in OS keychain when required
- Reference: mature secret handling in `gogcli`
- Decision: not required for v0; add when plaintext config becomes unacceptable

### `github.com/stretchr/testify`

- Use case: improve test readability
- Decision: introduce only if standard `testing` becomes too verbose

## Reject for now

### `viper`

- Reason to reject: heavy config layer with more implicit behavior
- Alternative: implement explicit `flag > env > file` in `internal/config`

### `resty`

- Reason to reject: `net/http` is sufficient for v0 API requirements
- Alternative: keep thin shared helpers in `internal/redash/client.go`

### Colored UI libraries (for example: `termenv`)

- Reason to reject: v0 is script-first and prioritizes stable JSON/text output
- Alternative: add later behind `internal/output` if needed

## Initial `go.mod` target

```text
require (
  github.com/spf13/cobra v1.8.1
)
```

Pin an exact stable version in `go.mod` and upgrade intentionally (do not use floating placeholders).
