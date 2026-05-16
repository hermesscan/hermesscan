# Baseline adoption guide

Baselines let an existing repository adopt HermesScan without fixing every current finding before enabling a CI gate.

## Start with a profile

Use a profile that matches the first adoption step:

| Profile | Use it when |
|---|---|
| `minimal` | You want advisory local scans before enabling a gate. |
| `ci` | You are ready to gate on high findings across the full catalog. |
| `supply-chain` | You want to start with supply-chain rules only. |

```powershell
.\hermesscan.exe init --profile minimal
```

Move to `ci` or `supply-chain` when the team is ready to enforce a gate.

## When to use a baseline

Use a baseline when a repository already has known findings and you want to block only new findings.

Do not use a baseline as a permanent ignore list. Treat it as technical debt that should shrink over time.

## 1. Run an advisory scan

Start with an advisory scan so the team can review the current findings without failing CI:

```powershell
.\hermesscan.exe scan . --summary --no-fail
```

For a reviewable report:

```powershell
.\hermesscan.exe scan . --format markdown --output .\reports\hermes-scan.md --no-fail
```

For an advisory-only config, use the `minimal` profile and keep CI non-blocking while findings are reviewed.

## 2. Review findings

Triage findings before creating a baseline:

- Fix findings that are real and practical to address now.
- Use inline suppressions for intentional, local exceptions that should stay next to the code.
- Use a baseline only for reviewed findings that are accepted temporarily across the repository.

## 3. Prefer inline suppressions for local exceptions

Use inline suppressions when a finding is intentionally safe at a specific line or in a specific file and the reason belongs with the code:

```bash
# hermesscan:disable-next-line HMS0001 -- fixture delay is intentional
sleep 30
```

```bash
sleep 30 # hermesscan:disable-line HMS0001 -- fixture delay is intentional
```

Use file-level suppression only when the whole file is intentionally outside the rule's normal expectation:

```bash
# hermesscan:disable-file HMS0001 -- dedicated fixture file
```

Prefer a baseline instead when the repository has many reviewed existing findings that should be paid down over time rather than annotated one by one.

## 4. Create a reviewed baseline

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

## 5. Use a baseline in CI

```powershell
.\hermesscan.exe scan . --baseline .\.hermesscan-baseline.json --fail-on high
```

With the GitHub Action:

```yaml
- name: Run HermesScan
  uses: hermesscan/hermesscan@v0.9.0
  with:
    path: .
    baseline: .hermesscan-baseline.json
    fail-on: high
```

For early rollout, run the Action in advisory mode first:

```yaml
- name: Run HermesScan advisory scan
  uses: hermesscan/hermesscan@v0.9.0
  with:
    path: .
    format: summary
    no-fail: 'true'
```

Then switch to the baseline gate:

```yaml
- name: Run HermesScan baseline gate
  uses: hermesscan/hermesscan@v0.9.0
  with:
    path: .
    baseline: .hermesscan-baseline.json
    format: summary
    fail-on: high
```

Keep reporting and gating separate when you also upload SARIF or artifacts. Generate reports with `no-fail: 'true'`, then run a final summary gate with `baseline` and `fail-on`.

## 6. Reduce the baseline

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

1. Initialize `minimal` for advisory scans, or `supply-chain` if the team wants a narrower first scope.
2. Review high-severity findings first.
3. Fix real issues and suppress intentional local exceptions inline with reasons.
4. Create a baseline only for remaining reviewed findings accepted temporarily.
5. Move to `ci` or `supply-chain` gating with `--fail-on high`.
6. Lower the gate to `--fail-on medium` only after high findings are handled.
7. Reduce the baseline over time instead of refreshing it to hide new findings.
