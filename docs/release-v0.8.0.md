# HermesScan v0.8.0 release notes

HermesScan v0.8.0 focuses on rule precision, rule-catalog validation, GitHub Action ergonomics, and release confidence.

## Highlights

- Added `hermesscan rules validate` for publish-ready rule catalog checks.
- Added tests that keep embedded default rules synchronized with `rules/hermes.rules.json`.
- Added precision coverage for common default-rule false-positive cases.
- Refined Docker Compose project-name detection to respect explicit project-name context.
- Refined broad GitHub Actions cache-key detection to ignore lockfile, hash, ref, SHA, and matrix-specific keys.
- Refined package-install cache findings to ignore runner-temp or run-scoped cache isolation.
- Added comma-delimited multi-rule support to the GitHub Action `rule` input.
- Added a complete GitHub Action reporting workflow for SARIF upload, retained report artifacts, and a separate summary gate.
- Added a rule authoring guide for rule fields, contextual precision patterns, and validation workflow.
- Expanded release smoke testing across all published binary assets.

## Upgrade notes

After the `v0.8.0` release is published, pin the GitHub Action to:

```yaml
uses: hermesscan/hermesscan@v0.8.0
```

The action downloads the matching `0.8.0` CLI by default unless the `version` input is overridden.

Run `rules validate` before publishing custom rule catalogs:

```powershell
.\hermesscan.exe rules validate
```

## Validation

Validation used before tagging:

```powershell
go test .\...
go vet .\...
go build -ldflags "-X main.version=0.8.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
.\hermesscan.exe rules validate
.\hermesscan.exe rules docs --output .\docs\rules.md
.\hermesscan.exe scan . --summary --exclude "examples/**" --no-fail
.\hermesscan.exe scan .\examples --summary --no-fail
```

Expected self-scan result, excluding intentionally risky examples:

```text
HermesScan: 0 findings
```

## GitHub Release body

Use this text for the GitHub Release body if generated notes are not sufficient:

```markdown
HermesScan v0.8.0 improves rule precision, rule-catalog validation, GitHub Action workflows, and release confidence.

### Added

- `hermesscan rules validate` for rule catalog validation.
- Embedded rule catalog sync tests.
- Default-rule precision tests for common false-positive cases.
- Comma-delimited multi-rule support for the GitHub Action `rule` input.
- A rule authoring guide.
- A complete GitHub Action reporting workflow for SARIF upload, report artifacts, and a separate summary gate.
- Expanded release smoke testing for all published binaries and checksums.

### Changed

- Refined Docker Compose project-name detection.
- Refined GitHub Actions broad cache-key detection.
- Refined package-install cache isolation handling.
```
