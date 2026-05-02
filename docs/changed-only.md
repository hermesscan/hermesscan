# Changed-file scans

HermesScan supports scanning only files changed according to Git. This is useful for pull-request workflows where you want fast feedback and fewer findings.

## Local examples

Scan files changed from `HEAD`:

```powershell
.\hermesscan.exe scan . --changed-only --summary --no-fail
```

Scan files changed from a base branch:

```powershell
.\hermesscan.exe scan . --changed-only --changed-base origin/main --summary --no-fail
```

## GitHub Actions example

Use `fetch-depth: 0` so Git has enough history to compare against the base ref.

```yaml
name: HermesScan changed files

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

      - name: Run HermesScan on changed files
        uses: hermesscan/hermesscan@v0.7.0
        with:
          path: .
          changed-only: 'true'
          changed-base: origin/main
          format: summary
          fail-on: high
```

## Notes

- `--changed-only` requires Git to be available.
- A shallow checkout may not contain the base ref. Use `fetch-depth: 0` in GitHub Actions.
- Changed-file scans are best for pull-request gates. Periodic full scans are still useful for drift detection.
