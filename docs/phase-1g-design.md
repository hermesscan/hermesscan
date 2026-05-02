# HermesScan Phase 1G.1 Design

Phase 1G.1 keeps HermesScan small while improving its CI usability.

## Goals

- Keep the native Go CLI easy to read and modify.
- Add report modes useful in local development and CI.
- Add SARIF output as the foundation for GitHub code scanning integration.
- Preserve simple JSON rule loading.
- Avoid external dependencies for now.

## Commands

```bash
hermesscan scan .
hermesscan scan . --summary
hermesscan scan . --quiet --fail-on high
hermesscan scan . --min-severity medium
hermesscan scan . --format markdown --output hermes-scan-report.md
hermesscan scan . --format json --output hermes-scan-report.json
hermesscan scan . --format sarif --output hermes-scan.sarif
```

## Output modes

| Mode | Purpose |
|---|---|
| `console` | Full human-readable local output |
| `summary` | Compact counts for CI logs |
| `markdown` | Artifact report |
| `json` | Machine-readable result |
| `sarif` | GitHub code scanning integration foundation |

## Severity behavior

`--min-severity` filters the report view only. `--fail-on` evaluates the full scan result so a quiet or filtered report does not accidentally hide blocking findings.

## Current limitations

- Rules are regex-based.
- SARIF output is minimal but valid JSON and follows the SARIF 2.1.0 shape.
- `--no-color` is accepted for CLI compatibility, but console color is not currently emitted.
- File discovery does not yet honor `.gitignore`.
