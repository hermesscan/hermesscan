# Baseline adoption guide

Baselines let an existing repository adopt HermesScan without fixing every current finding before enabling a CI gate.

## When to use a baseline

Use a baseline when a repository already has known findings and you want to block only new findings.

Do not use a baseline as a permanent ignore list. Treat it as technical debt that should shrink over time.

## Create a baseline

```powershell
.\hermesscan.exe scan . --create-baseline .\.hermesscan-baseline.json --no-fail
```

Linux/macOS:

```bash
./hermesscan scan . --create-baseline ./.hermesscan-baseline.json --no-fail
```

Commit the baseline if the team agrees that the current findings are accepted temporarily.

## Use a baseline in CI

```powershell
.\hermesscan.exe scan . --baseline .\.hermesscan-baseline.json --fail-on high
```

With the GitHub Action:

```yaml
- name: Run HermesScan
  uses: hermesscan/hermesscan@v0.7.0
  with:
    path: .
    baseline: .hermesscan-baseline.json
    fail-on: high
```

## Refresh a baseline

Only refresh the baseline intentionally after reviewing findings.

```powershell
.\hermesscan.exe scan . --create-baseline .\.hermesscan-baseline.json --no-fail
```

Review the diff before committing the updated baseline.

## Recommended adoption flow

1. Run HermesScan in advisory mode.
2. Review high-severity findings first.
3. Suppress intentional findings inline with comments.
4. Create a baseline for remaining accepted findings.
5. Enable `--fail-on high` in CI.
6. Reduce the baseline over time.
