# HermesScan Phase 4 Design

Phase 4 focuses on adoption quality and rule precision.

## New capabilities

- Baseline creation with `--create-baseline`.
- Baseline filtering with `--baseline`.
- Finding fingerprints stored in JSON and SARIF/JSON output.
- Rule metadata: `category` and `tags`.
- Rule list output includes category.
- Additional GitHub Actions rules for `pull_request_target` and `permissions: write-all`.

## Baseline model

A baseline captures currently accepted findings. Future scans can suppress those exact findings so teams may adopt HermesScan without fixing every existing issue at once.

Typical flow:

```bash
hermesscan scan . --create-baseline .hermesscan-baseline.json --no-fail
hermesscan scan . --baseline .hermesscan-baseline.json --fail-on high
```

The baseline fingerprint uses rule id, normalized file path, line number, and matched text. This is intentionally conservative: if code moves or changes, the finding is treated as new and should be reviewed again.

## Rule metadata

Rules now include:

- `category`, such as `orchestration`, `isolation`, `cache`, `supply-chain`, or `reliability`.
- `tags`, such as `docker`, `github-actions`, `powershell`, or `shared-runner`.

These fields support better reporting, future filtering, and documentation generation.
