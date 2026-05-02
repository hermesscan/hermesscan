# v0.7.0 release checklist

Use this checklist before creating the `v0.7.0` tag.

## Local validation

```powershell
go test .\...
go vet .\...
go build -ldflags "-X main.version=0.7.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
.\hermesscan.exe rules list
.\hermesscan.exe rules docs --output .\docs\rules.md
.\hermesscan.exe scan . --summary --exclude "examples/**" --no-fail
```

Expected self-scan result:

```text
HermesScan: 0 findings
```

## Documentation checks

- `README.md` describes the new `0.7.0` capabilities.
- `docs/rules.md` has been regenerated from the active rule catalog.
- `docs/github-action.md` references `v0.7.0` examples.
- `docs/install.md` references `0.7.0` release binaries.
- `CHANGELOG.md` has a `0.7.0` entry.

## Git checks

```powershell
git status
git diff --stat
```

Confirm no generated binaries are staged:

```powershell
git status --short
```

Do not commit:

```text
hermesscan.exe
hermesscan
dist/
reports/
*.sarif
```

## Commit

```powershell
git add .
git commit -m "Prepare HermesScan v0.7.0"
git push
```

## Tag

Create the tag only after CI passes on `main`.

```powershell
git tag v0.7.0
git push origin v0.7.0
```

## Release verification

After the release workflow completes:

- Confirm all release binaries exist.
- Confirm `checksums.txt` exists.
- Run the release smoke workflow for `0.7.0`.
- Verify the GitHub Action example works against `hermesscan/hermesscan@v0.7.0`.
