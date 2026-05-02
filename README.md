# HermesScan

[![CI](https://github.com/hermesscan/hermesscan/actions/workflows/ci.yml/badge.svg)](https://github.com/hermesscan/hermesscan/actions/workflows/ci.yml)
[![HermesScan](https://github.com/hermesscan/hermesscan/actions/workflows/hermes-scan.yml/badge.svg)](https://github.com/hermesscan/hermesscan/actions/workflows/hermes-scan.yml)

HermesScan is a native Go static analyzer for build scripts, CI scripts, and pipeline definitions.
It detects accidental orchestration, shared-runner hazards, fixed-port risks, package-cache collision risks,
weak process lifecycle patterns, and other CI reliability smells.

HermesScan is **not** a CI platform. It is a scanner and quality gate for the scripts and workflow files that feed CI.

## Status

Current version: `0.6.1`

Phase 6 focuses on GitHub adoption and packaging polish:

- normalized GitHub Actions annotation paths
- first-party composite GitHub Action wrapper
- installation docs for Windows, Linux, and macOS
- GitHub Action and SARIF usage docs
- module path updated to `github.com/hermesscan/hermesscan`
- contribution guidance

## Quick start from source

### Windows PowerShell

```powershell
go test .\...
go build -ldflags "-X main.version=0.6.1" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
.\hermesscan.exe scan .\examples --summary --no-fail
```

### Linux/macOS

```bash
go test ./...
go build -ldflags "-X main.version=0.6.1" -o ./hermesscan ./cmd/hermesscan
./hermesscan version
./hermesscan scan ./examples --summary --no-fail
```

## Install

See [docs/install.md](docs/install.md).

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
hermesscan init
hermesscan rules list
hermesscan rules show RULE_ID
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

Basic usage. The action tag and downloaded CLI version are both `0.6.1` by default:

```yaml
- name: Run HermesScan
  uses: hermesscan/hermesscan@v0.6.1
  with:
    path: .
    format: summary
    fail-on: high
```

SARIF upload:

```yaml
- name: Generate HermesScan SARIF
  uses: hermesscan/hermesscan@v0.6.1
  with:
    path: .
    format: sarif
    output: reports/hermes-scan.sarif
    no-fail: 'true'

- name: Upload SARIF
  uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: reports/hermes-scan.sarif
```

## Configuration

Create a starter config:

```powershell
.\hermesscan.exe init
```

Example `.hermesscan.json`:

```json
{
  "rules": "",
  "include": [],
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
  "disabledRules": [],
  "severityOverrides": {
    "HMS0010": "Low"
  },
  "failOn": "high",
  "minSeverity": "",
  "suppressionsEnabled": true
}
```

Run with config. An empty `rules` value means HermesScan uses the local `rules/hermes.rules.json` file when present, otherwise it uses embedded default rules:

```powershell
.\hermesscan.exe scan . --config .hermesscan.json
```

Override config failure behavior:

```powershell
.\hermesscan.exe scan . --config .hermesscan.json --no-fail
```

## Baseline adoption workflow

Create a baseline from current findings:

```powershell
.\hermesscan.exe scan . --create-baseline .\.hermesscan-baseline.json --no-fail
```

Use the baseline to fail only on new findings:

```powershell
.\hermesscan.exe scan . --baseline .\.hermesscan-baseline.json --fail-on high
```

The baseline uses finding fingerprints based on rule id, normalized file path, line number, and matched text. If the risky line moves or changes, HermesScan treats it as a new finding that should be reviewed again.

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

## Rule inventory

```powershell
.\hermesscan.exe rules list
.\hermesscan.exe rules show HMS0001
```

## Build release binaries

PowerShell 5.1-compatible build script:

```powershell
.\scripts\Build-HermesScan.ps1 -AllTargets -Version 0.6.1
```

Outputs are written to `dist/` with `.sha256` checksum files.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Design principle

HermesScan should be advisory by default during adoption and strict only when a team intentionally turns it into a quality gate.
