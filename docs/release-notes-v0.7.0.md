# HermesScan v0.7.0 release notes

HermesScan `v0.7.0` focuses on rule quality and adoption polish.

## Highlights

- Added single-rule scanning with `--rule`.
- Added generated Markdown rule documentation with `hermesscan rules docs`.
- Added rule inventory commands: `rules categories` and `rules tags`.
- Added config-level filters for enabled rules, categories, and tags.
- Refined the PostgreSQL fixed-port rule to focus on exposed/bound port contexts.
- Lowered the package-install cache rule to `Low` to reduce noise.
- Added GitHub Actions checks for self-hosted runners and broad cache keys.

## Upgrade notes

The GitHub Action examples now use:

```yaml
uses: hermesscan/hermesscan@v0.7.0
```

The action downloads the `0.7.0` CLI by default unless the `version` input is overridden.

## Validation

Validation used before tagging:

```powershell
go test .\...
go vet .\...
go build -ldflags "-X main.version=0.7.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe scan . --summary --exclude "examples/**" --no-fail
```
