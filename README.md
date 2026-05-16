# HermesScan

[![CI](https://github.com/hermesscan/hermesscan/actions/workflows/ci.yml/badge.svg)](https://github.com/hermesscan/hermesscan/actions/workflows/ci.yml)
[![HermesScan](https://github.com/hermesscan/hermesscan/actions/workflows/hermes-scan.yml/badge.svg)](https://github.com/hermesscan/hermesscan/actions/workflows/hermes-scan.yml)

HermesScan is a native Go static analyzer for build scripts, CI scripts, and pipeline definitions.
It detects accidental orchestration, shared-runner hazards, fixed-port risks, package-cache collision risks,
weak process lifecycle patterns, and other CI reliability smells.

HermesScan is **not** a CI platform. It is a scanner and quality gate for the scripts and workflow files that feed CI.

## Status

Current development version: `0.10.0`

Latest release: `0.9.0`

> HermesScan is currently in public preview. Rules are intentionally conservative and may evolve as the scanner matures.

Version 0.10.0 focuses on adoption polish, distribution readiness, release-integrity rules, and better configuration workflows.

Version 0.9.0 focused on SBOM-aware release workflow assurance and the first Windows package-manager installation path.

Version 0.8.0 focused on rule precision and packaging polish:

- rule catalog validation with `rules validate`
- sync checks between embedded and repository rule catalogs
- default-rule precision tests for common false-positive cases
- Docker Compose project-name context handling
- GitHub Actions cache-key specificity handling
- package-install cache isolation context handling
- safer rule authoring and release validation workflows

## Quick start from source

### Windows PowerShell

```powershell
go test .\...
go build -ldflags "-X main.version=0.10.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
.\hermesscan.exe scan .\examples --summary --no-fail
```

### Linux/macOS

```bash
go test ./...
go build -ldflags "-X main.version=0.10.0" -o ./hermesscan ./cmd/hermesscan
./hermesscan version
./hermesscan scan ./examples --summary --no-fail
```

## Install

See [docs/install.md](docs/install.md).

Additional guides:

- [GitHub Action usage](docs/github-action.md)
- [Baseline adoption guide](docs/baseline-adoption.md)
- [Rule reference](docs/rules.md)
- [Rule authoring guide](docs/rule-authoring.md)
- [Changed-file scans](docs/changed-only.md)
- [SBOM and release assurance](docs/sbom-release-assurance.md)
- [Scoop packaging plan](docs/scoop-packaging.md)
- [v0.10.0 milestone](docs/milestones/v0.10.0.md)
- [v0.9.0 milestone](docs/milestones/v0.9.0.md)
- [v0.9.0 release checklist](docs/release-v0.9.0-checklist.md)
- [v0.9.0 release notes](docs/release-v0.9.0.md)
- [v0.8.0 release checklist](docs/release-v0.8.0-checklist.md)
- [v0.8.0 release notes](docs/release-v0.8.0.md)
- [v0.7.0 release notes](docs/release-v0.7.0.md)

Short Windows development note: PowerShell does not execute programs from the current directory by bare name. Use:

```powershell
.\hermesscan.exe scan .
```

After installing to a directory on PATH, this works:

```powershell
hermesscan scan .
```

## Commands

```text
hermesscan scan [path]
hermesscan init [--profile minimal|ci|supply-chain]
hermesscan rules list
hermesscan rules show RULE_ID
hermesscan rules docs [--output docs/rules.md]
hermesscan rules validate
hermesscan rules categories
hermesscan rules tags
hermesscan version
```

## Scan examples

Console report:

```powershell
.\hermesscan.exe scan .\examples
```

Summary report:

```powershell
.\hermesscan.exe scan .\examples --summary --no-fail
```

JSON report:

```powershell
.\hermesscan.exe scan .\examples --format json --output .\reports\hermes-scan.json --no-fail
```

SARIF report:

```powershell
.\hermesscan.exe scan .\examples --format sarif --output .\reports\hermes-scan.sarif --no-fail
```

GitHub Actions annotations:

```powershell
.\hermesscan.exe scan .\examples --github-annotations --no-fail
```

Equivalent explicit format:

```powershell
.\hermesscan.exe scan .\examples --format github --no-fail
```

Fail when high findings exist:

```powershell
.\hermesscan.exe scan . --fail-on high
```

Generate a report but never fail the process:

```powershell
.\hermesscan.exe scan . --config .hermesscan.json --no-fail
```

Only report high and critical findings:

```powershell
.\hermesscan.exe scan . --min-severity high
```

Include or exclude paths from the command line:

```powershell
.\hermesscan.exe scan . --include "scripts/**" --exclude "dist/**" --exclude "reports/**"
```

Scan only supply-chain rules:

```powershell
.\hermesscan.exe scan . --category supply-chain --summary --no-fail
```

Scan only rules tagged with `docker`:

```powershell
.\hermesscan.exe scan . --tag docker --summary --no-fail
```

Scan only selected rule IDs:

```powershell
.\hermesscan.exe scan . --rule HMS0002 --rule HMS0013 --summary --no-fail
```

Scan only files changed from `HEAD`:

```powershell
.\hermesscan.exe scan . --changed-only --summary --no-fail
```

Scan only files changed from a base ref:

```powershell
.\hermesscan.exe scan . --changed-only --changed-base origin/main --summary --no-fail
```

## GitHub Action

See [docs/github-action.md](docs/github-action.md).

Basic usage after the `v0.9.0` tag is published. The action tag and downloaded CLI version are both `0.9.0` by default:

```yaml
- name: Run HermesScan
  uses: hermesscan/hermesscan@v0.9.0
  with:
    path: .
    format: summary
    fail-on: high
```

Run selected rules through the action:

```yaml
- name: Run selected HermesScan rules
  uses: hermesscan/hermesscan@v0.9.0
  with:
    path: .
    rule: HMS0002,HMS0013
    format: summary
    no-fail: 'true'
```

SARIF upload:

```yaml
- name: Generate HermesScan SARIF
  uses: hermesscan/hermesscan@v0.9.0
  with:
    path: .
    format: sarif
    output: reports/hermes-scan.sarif
    no-fail: 'true'

- name: Upload SARIF
  uses: github/codeql-action/upload-sarif@v4
  with:
    sarif_file: reports/hermes-scan.sarif
```

For a complete workflow that uploads SARIF, stores a Markdown report artifact, and runs a separate summary gate, see [GitHub Action usage](docs/github-action.md#sarif-plus-report-artifact).

## Reporting false positives

HermesScan is rule-based and may flag legitimate patterns. If you find a false positive, open an issue with:

- the rule ID,
- the matched file type,
- a minimal code example,
- why the pattern is safe in your case.

Use inline suppressions or a baseline only when the finding has been reviewed and accepted. Prefer inline suppressions for intentional local exceptions; use a baseline for reviewed repository-wide debt that should shrink over time.

## Configuration

Create a starter config:

```powershell
.\hermesscan.exe init
```

Choose an adoption profile when you want a narrower starting point:

```powershell
.\hermesscan.exe init --profile minimal
.\hermesscan.exe init --profile ci
.\hermesscan.exe init --profile supply-chain
```

Profiles:

| Profile | Use case |
|---|---|
| `minimal` | Advisory local scans. It enables common excludes and suppressions, but does not configure a default fail threshold. |
| `ci` | Default CI adoption. It enables common excludes, suppressions, and a `high` fail threshold. |
| `supply-chain` | CI adoption focused on supply-chain rules by setting the `supply-chain` category filter and a `high` fail threshold. |

Example `.hermesscan.json`:

```json
{
  "rules": "",
  "exclude": [
    "dist/**",
    "build/**",
    "node_modules/**",
    "vendor/**",
    "reports/**",
    "coverage/**",
    "tmp/**",
    ".git/**"
  ],
  "failOn": "high",
  "suppressionsEnabled": true
}
```

Run with config. An empty `rules` value means HermesScan uses the local `rules/hermes.rules.json` file when present, otherwise it uses embedded default rules:

```powershell
.\hermesscan.exe scan . --config .hermesscan.json
```

JSON schemas for editor integration and external validation are available at:

- `schemas/hermesscan-config.schema.json`
- `schemas/hermesscan-rule-catalog.schema.json`

Override config failure behavior:

```powershell
.\hermesscan.exe scan . --config .hermesscan.json --no-fail
```

## Baseline adoption workflow

Start advisory:

```powershell
.\hermesscan.exe init --profile minimal
.\hermesscan.exe scan . --summary --no-fail
```

Use inline suppressions for reviewed local exceptions:

```bash
# hermesscan:disable-next-line HMS0001 -- fixture delay is intentional
sleep 30
```

Create a baseline from current findings:

```powershell
.\hermesscan.exe scan . --create-baseline .\.hermesscan-baseline.json --no-fail
```

Use the baseline to fail only on new findings:

```powershell
.\hermesscan.exe scan . --baseline .\.hermesscan-baseline.json --fail-on high
```

The baseline uses finding fingerprints based on rule id, normalized file path, line number, and matched text. If the risky line moves or changes, HermesScan treats it as a new finding that should be reviewed again.

For the full adoption sequence, including reviewed baseline creation, CI gating, and intentional baseline reduction, see [Baseline adoption guide](docs/baseline-adoption.md).

## Suppressions

Suppress the next line:

```bash
# hermesscan:disable-next-line HMS0001 -- fixture delay
sleep 30
```

Suppress the current line:

```bash
sleep 30 # hermesscan:disable-line HMS0001
```

Suppress a rule for the whole file:

```bash
# hermesscan:disable-file HMS0001
```

## Rule inventory and documentation

```powershell
.\hermesscan.exe rules list
.\hermesscan.exe rules show HMS0001
.\hermesscan.exe rules categories
.\hermesscan.exe rules tags
.\hermesscan.exe rules validate
.\hermesscan.exe rules docs --output .\docs\rules.md
```

## Generate rule reference

Regenerate the committed rule reference after rule changes:

```powershell
.\hermesscan.exe rules docs --output .\docs\rules.md
```


Commit `docs/rules.md` with releases so users can review the active rule catalog without running the CLI first.

## Release checklist

Before tagging a release:

```powershell
go test .\...
go vet .\...
go build -ldflags "-X main.version=0.10.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
.\hermesscan.exe rules validate
.\hermesscan.exe rules docs --output .\docs\rules.md
.\hermesscan.exe scan . --summary --exclude "examples/**" --no-fail
```

After publishing a release, run the `Release smoke test` workflow for the new version. It verifies all published binary, SBOM, and scan-report assets, validates `checksums.txt`, and runs native CLI smoke checks on each target.

## Build release binaries

PowerShell 5.1-compatible build script:

```powershell
.\scripts\Build-HermesScan.ps1 -AllTargets -Version 0.10.0
```

Outputs are written to `dist/` with `.sha256` checksum files.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Design principle

HermesScan should be advisory by default during adoption and strict only when a team intentionally turns it into a quality gate.
