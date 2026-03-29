# dashcli

`dashcli` is a Go-based CLI for calling the Redash API.

This repository is currently in the architecture/spec phase.

## Scope (v0)

- Stable Redash API calls from CLI
- Script-friendly JSON output
- Testable internal structure for long-term maintenance

Planned initial commands include:

- `dash version`
- `dash me`
- `dash query list|get|run`
- `dash job get|wait`
- `dash dashboard list|get`
- `dash datasource list`

## Documentation

- `docs/spec.md` — product and CLI specification
- `docs/architecture.md` — package layout and execution flow
- `docs/libraries.md` — dependency selection (adopt / optional / reject)

## Reference projects

Architecture decisions are based on patterns from:

- `steipete/gogcli`
- `steipete/sonoscli`
- `steipete/blucli`
