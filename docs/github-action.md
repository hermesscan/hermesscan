# GitHub Action usage

HermesScan includes a composite GitHub Action wrapper in `action.yml`.

## Action inputs

| Input | Default | Description |
|---|---:|---|
| `path` | `.` | Path to scan. |
| `version` | `0.8.0` | HermesScan release binary version to download. Use `latest` only when you intentionally want the newest release binary. |
| `repository` | `hermesscan/hermesscan` | Repository that hosts HermesScan release binaries. |
| `format` | `summary` | Report format: `console`, `summary`, `markdown`, `json`, `sarif`, or `github`. |
| `output` | empty | Optional report output file. |
| `config` | empty | Optional `.hermesscan.json` path. |
| `baseline` | empty | Optional baseline file. |
| `fail-on` | `high` | Fail when findings meet this severity: `info`, `low`, `medium`, `high`, `critical`, or `none`. |
| `min-severity` | empty | Only report findings at or above this severity. |
| `category` | empty | Restrict scanning to one rule category. |
| `tag` | empty | Restrict scanning to rules with one tag. |
| `rule` | empty | Restrict scanning to one rule ID or a comma-delimited list such as `HMS0001,HMS0010`. |
| `changed-only` | `false` | Scan only files changed according to Git. |
| `changed-base` | empty | Base ref or commit for changed-file scans. |
| `github-annotations` | `false` | Emit GitHub Actions annotations. |
| `no-fail` | `false` | Always return success even when findings are detected. |

After `v0.8.0` is published, the action wrapper is pinned by `uses: hermesscan/hermesscan@v0.8.0`, and the downloaded CLI defaults to version `0.8.0`. Override `version` only when you intentionally want a different release binary.

## Basic pull-request gate

```yaml
name: HermesScan

on:
  pull_request:
    branches:
      - main

permissions:
  contents: read

jobs:
  hermesscan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Run HermesScan
        uses: hermesscan/hermesscan@v0.8.0
        with:
          path: .
          format: summary
          fail-on: high
```

## PR annotations

```yaml
- name: Run HermesScan annotations
  uses: hermesscan/hermesscan@v0.8.0
  with:
    path: .
    github-annotations: 'true'
    changed-only: 'true'
    changed-base: origin/main
    fail-on: high
```

You can also request annotation output by setting `format: github`; `github-annotations: 'true'` is the clearer form and is preferred in examples.

## Changed-file pull-request scan

```yaml
- name: Run HermesScan on changed files
  uses: hermesscan/hermesscan@v0.8.0
  with:
    path: .
    changed-only: 'true'
    changed-base: origin/main
    format: summary
    fail-on: high
```

Use `actions/checkout` with `fetch-depth: 0` when using changed-file scans so the base ref is available locally.

## Rule-focused scans

Scan with one rule:

```yaml
- name: Run one HermesScan rule
  uses: hermesscan/hermesscan@v0.8.0
  with:
    path: .
    rule: HMS0002
    format: summary
    no-fail: 'true'
```

Scan with multiple selected rules:

```yaml
- name: Run selected HermesScan rules
  uses: hermesscan/hermesscan@v0.8.0
  with:
    path: .
    rule: HMS0002,HMS0013
    format: summary
    no-fail: 'true'
```

## SARIF for GitHub Code Scanning

```yaml
name: HermesScan SARIF

on:
  pull_request:
    branches:
      - main
  push:
    branches:
      - main

permissions:
  contents: read
  security-events: write
  actions: read

jobs:
  hermesscan-sarif:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Generate HermesScan SARIF
        uses: hermesscan/hermesscan@v0.8.0
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

## SARIF plus report artifact

Use a non-failing reporting pass for SARIF and retained reports, then run a separate summary gate that controls the job result.

```yaml
name: HermesScan reporting

on:
  pull_request:
    branches:
      - main
  push:
    branches:
      - main

permissions:
  contents: read
  security-events: write
  actions: read

jobs:
  hermesscan:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Generate HermesScan SARIF
        uses: hermesscan/hermesscan@v0.8.0
        with:
          path: .
          format: sarif
          output: reports/hermes-scan.sarif
          no-fail: 'true'

      - name: Upload SARIF
        uses: github/codeql-action/upload-sarif@v4
        with:
          sarif_file: reports/hermes-scan.sarif

      - name: Generate HermesScan Markdown report
        uses: hermesscan/hermesscan@v0.8.0
        with:
          path: .
          format: markdown
          output: reports/hermes-scan.md
          no-fail: 'true'

      - name: Upload HermesScan report artifact
        uses: actions/upload-artifact@v4
        with:
          name: hermesscan-report
          path: reports/

      - name: Run HermesScan gate
        uses: hermesscan/hermesscan@v0.8.0
        with:
          path: .
          format: summary
          fail-on: high
```

## Use a baseline

```yaml
- name: Run HermesScan with baseline
  uses: hermesscan/hermesscan@v0.8.0
  with:
    path: .
    baseline: .hermesscan-baseline.json
    fail-on: high
```

## Pinning

For production use, pin the action to a release tag or commit SHA.

```yaml
uses: hermesscan/hermesscan@v0.8.0
```

Avoid floating references such as `@main` for required gates.
