# Baseline adoption guide

Baselines let an existing repository adopt HermesScan without fixing every current finding before enabling a CI gate.

## When to use a baseline

Use a baseline when a repository already has known findings and you want to block only new findings.

Do not use a baseline as a permanent ignore list. Treat it as technical debt that should shrink over time.

## Create a baseline

Start with an advisory scan so the team can review the current findings without failing CI:

```powershell
.\hermesscan.exe scan . --summary --no-fail
```

For a reviewable report:

```powershell
.\hermesscan.exe scan . --format markdown --output .\reports\hermes-scan.md --no-fail
```

After the team reviews current findings, create the baseline:

```powershell
.\hermesscan.exe scan . --create-baseline .\.hermesscan-baseline.json --no-fail
```

Linux/macOS:

```bash
./hermesscan scan . --create-baseline ./.hermesscan-baseline.json --no-fail
```

Commit the baseline if the team agrees that the current findings are accepted temporarily.

Review the generated file before committing it. The baseline should contain findings the team has explicitly accepted for now, not findings that should be fixed before the gate is enabled.

## Use a baseline in CI

```powershell
.\hermesscan.exe scan . --baseline .\.hermesscan-baseline.json --fail-on high
```

With the GitHub Action:

```yaml
- name: Run HermesScan
  uses: hermesscan/hermesscan@v0.8.0
  with:
    path: .
    baseline: .hermesscan-baseline.json
    fail-on: high
```

For early rollout, run the Action in advisory mode first:

```yaml
- name: Run HermesScan advisory scan
  uses: hermesscan/hermesscan@v0.8.0
  with:
    path: .
    format: summary
    no-fail: 'true'
```

Then switch to the baseline gate:

```yaml
- name: Run HermesScan baseline gate
  uses: hermesscan/hermesscan@v0.8.0
  with:
    path: .
    baseline: .hermesscan-baseline.json
    format: summary
    fail-on: high
```

Keep reporting and gating separate when you also upload SARIF or artifacts. Generate reports with `no-fail: 'true'`, then run a final summary gate with `baseline` and `fail-on`.

## Refresh a baseline

Only refresh the baseline intentionally after reviewing findings.

```powershell
.\hermesscan.exe scan . --create-baseline .\.hermesscan-baseline.json --no-fail
```

Review the diff before committing the updated baseline.

To shrink a baseline safely, create the next version in a temporary file and compare it with the committed baseline:

```powershell
.\hermesscan.exe scan . --create-baseline .\.hermesscan-baseline.next.json --no-fail
git diff --no-index .\.hermesscan-baseline.json .\.hermesscan-baseline.next.json
```

If the next baseline removes entries because findings were fixed, replace the committed baseline:

```powershell
Move-Item -Path .\.hermesscan-baseline.next.json -Destination .\.hermesscan-baseline.json -Force
```

Do not refresh the baseline just to make a new finding disappear. Review the finding, fix it or suppress it intentionally, then update the baseline only if the remaining finding is accepted technical debt.

## Recommended adoption flow

1. Run HermesScan in advisory mode.
2. Review high-severity findings first.
3. Suppress intentional findings inline with comments.
4. Create a baseline for remaining accepted findings.
5. Enable `--fail-on high` in CI.
6. Lower the gate to `--fail-on medium` only after high findings are handled.
7. Reduce the baseline over time.
