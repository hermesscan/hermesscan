# HermesScan Phase 2 Design

Phase 2 turns the Phase 1G.1 prototype into a more practical repository scanner.

## Goals

- Keep the scanner small and understandable.
- Add repository-level configuration.
- Add controlled suppression for intentional findings.
- Make output paths work naturally on Windows, Linux, and macOS.
- Prepare SARIF output for GitHub Code Scanning.
- Improve CI usability without turning HermesScan into a CI platform.

## Configuration

HermesScan automatically looks for `.hermesscan.json` under the scan root.

Create a starter config with one of the adoption profiles:

```bash
hermesscan init --profile minimal
hermesscan init --profile ci
hermesscan init --profile supply-chain
```

The default `ci` profile keeps the normal high-severity gate. The `minimal` profile is advisory by default, and the `supply-chain` profile filters to the `supply-chain` rule category.

The configuration schema is published at `schemas/hermesscan-config.schema.json`.

Supported settings:

| Setting | Purpose |
|---|---|
| `rules` | Relative or absolute path to the JSON rule catalog |
| `include` | Glob-like file patterns to include in scanning |
| `exclude` | Glob-like file patterns to exclude from scanning |
| `enabledRules` | Rule IDs to include |
| `disabledRules` | Rule IDs to skip entirely |
| `categories` | Rule categories to include |
| `tags` | Rule tags to include |
| `severityOverrides` | Per-rule severity remapping |
| `failOn` | Default fail threshold |
| `minSeverity` | Default report filtering threshold |
| `suppressionsEnabled` | Enable or disable inline suppressions |

## Suppressions

Suppressions are intended for deliberate exceptions, fixtures, or transitional work.

Supported markers:

```text
hermesscan:disable-next-line HMS0001
hermesscan:disable-line HMS0001
hermesscan:disable-file HMS0001
```

If no rule ID is specified, the suppression applies to all HermesScan rules for that scope.

## Output behavior

When `--output` points to a nested path, HermesScan creates the parent directory automatically.

Example:

```bash
hermesscan scan . --format sarif --output reports/code-scanning/hermes.sarif
```

## Phase 2 boundaries

Phase 2 intentionally does not add AST parsing, .gitignore interpretation, or custom rule packs. Those are better candidates for later phases.
