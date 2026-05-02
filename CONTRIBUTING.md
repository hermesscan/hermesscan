# Contributing to HermesScan

Thanks for helping improve HermesScan.

## Development prerequisites

- Go 1.22 or later
- Git
- Windows PowerShell 5.1 only if working on the helper scripts

## Local validation

Run these before opening a pull request:

```bash
go test ./...
go vet ./...
go build -ldflags "-X main.version=0.8.0" -o hermesscan ./cmd/hermesscan
./hermesscan rules validate
./hermesscan scan ./examples --summary --no-fail
```

Windows PowerShell:

```powershell
go test .\...
go vet .\...
go build -ldflags "-X main.version=0.8.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe rules validate
.\hermesscan.exe scan .\examples --summary --no-fail
```

## Rule changes

See [docs/rule-authoring.md](docs/rule-authoring.md) for the full rule model, contextual precision fields, and validation workflow.

Rules live in:

```text
rules/hermes.rules.json
```

Each rule should include:

- stable `id`
- clear `name`
- `severity`
- `category`
- useful `tags`
- file types
- regex pattern
- description
- recommendation

Rule-change checklist:

- Update both `rules/hermes.rules.json` and `internal/rules/defaults/hermes.rules.json`.
- Add or update precision tests in `internal/scanner/rule_precision_test.go`.
- Run `hermesscan rules validate`.
- Regenerate `docs/rules.md`.
- Run the local validation commands above.

## Rule severity guidance

| Severity | Use for |
|---|---|
| Critical | Extremely dangerous patterns that should almost always block CI |
| High | Strong CI reliability or supply-chain risks |
| Medium | Risky patterns that frequently need review |
| Low | Reproducibility, hygiene, or advisory findings |
| Info | Educational notes or optional improvements |

## False positives

Prefer rule precision over broad noisy matching. If a rule is intentionally broad, document why and keep its severity conservative.

## PowerShell scripts

PowerShell helper scripts must remain compatible with Windows PowerShell 5.1.

Use:

- `[CmdletBinding(SupportsShouldProcess = $true)]`
- comment-based help
- `-WhatIf` examples for state-changing operations
- no PowerShell 7-only syntax
- no `Invoke-Expression`
