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
go build -ldflags "-X main.version=0.7.0" -o hermesscan ./cmd/hermesscan
./hermesscan scan ./examples --summary --no-fail
```

Windows PowerShell:

```powershell
go test .\...
go vet .\...
go build -ldflags "-X main.version=0.7.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe scan .\examples --summary --no-fail
```

## Rule changes

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

Add or update tests when changing rules or scanner behavior.

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
