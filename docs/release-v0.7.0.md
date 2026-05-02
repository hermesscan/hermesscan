# HermesScan v0.7.0 release notes

HermesScan v0.7.0 is a rule-quality and adoption-polish release built on the public `v0.6.1` baseline.

## Highlights

- Added `scan --rule RULE_ID` for focused rule validation.
- Added `rules docs`, `rules categories`, and `rules tags` commands.
- Added generated rule reference documentation in `docs/rules.md`.
- Added config-level `enabledRules`, `categories`, and `tags` filters.
- Refined `HMS0002` so PostgreSQL port findings focus on exposed or bound port contexts.
- Lowered `HMS0010` package-install cache findings to `Low` because the rule is advisory without stronger concurrency evidence.
- Added GitHub Actions rules for self-hosted runners and broad cache keys.

## Validation checklist

Before tagging `v0.7.0`, run:

```powershell
go test .\...
go vet .\...
go build -ldflags "-X main.version=0.7.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
.\hermesscan.exe rules docs --output .\docsules.md
.\hermesscan.exe scan . --summary --exclude "examples/**" --no-fail
```

Expected version output:

```text
HermesScan 0.7.0
```

Expected self-scan result, excluding intentionally risky examples:

```text
HermesScan: 0 findings across ... files
```

## Release notes

Use this text for the GitHub Release body if generated notes are not sufficient:

```markdown
HermesScan v0.7.0 improves rule precision and adoption workflows.

### Added

- `scan --rule RULE_ID` for focused scans.
- `rules docs` to generate Markdown rule documentation.
- `rules categories` and `rules tags` inventory commands.
- Config filters for `enabledRules`, `categories`, and `tags`.
- New GitHub Actions rules for self-hosted runners and broad cache keys.

### Changed

- Refined fixed PostgreSQL port detection to focus on exposed/bound port contexts.
- Lowered package-install cache risk from Medium to Low.
- Updated docs and examples for v0.7.0.
```
