# v0.8.0 release checklist

Use this checklist before creating the `v0.8.0` tag.

## Local validation

```powershell
go test .\...
go vet .\...
go build -ldflags "-X main.version=0.8.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
.\hermesscan.exe rules validate
.\hermesscan.exe rules docs --output .\docs\rules.md
.\hermesscan.exe scan . --summary --exclude "examples/**" --no-fail
.\hermesscan.exe scan .\examples --summary --no-fail
```

Expected version output:

```text
HermesScan 0.8.0
```

Expected self-scan result, excluding intentionally risky examples:

```text
HermesScan: 0 findings
```

The examples scan should continue to report the intentional fixture findings.

## Release artifact validation

Build every target locally before tagging:

```powershell
.\scripts\Build-HermesScan.ps1 -AllTargets -Version 0.8.0
```

Confirm `dist/` contains:

```text
hermesscan-windows-amd64.exe
hermesscan-windows-arm64.exe
hermesscan-linux-amd64
hermesscan-linux-arm64
hermesscan-darwin-amd64
hermesscan-darwin-arm64
```

The local build script also writes per-file `.sha256` files. The GitHub release workflow publishes `checksums.txt`.

## Documentation checks

- `README.md` describes the `0.8.0` development focus.
- `CHANGELOG.md` has the finalized `0.8.0` entry and date.
- `docs/rules.md` has been regenerated from the active rule catalog.
- `docs/github-action.md` includes multi-rule, SARIF, artifact, and gate examples.
- `docs/baseline-adoption.md` includes advisory, gated, and baseline reduction examples.
- `docs/install.md` includes copy-paste checksum verification examples.
- `docs/release-v0.8.0.md` matches the release contents.

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
git commit -m "Prepare HermesScan v0.8.0"
git push
```

## Tag

Create the tag only after CI passes on `main`.

```powershell
git tag v0.8.0
git push origin v0.8.0
```

## Release verification

After the release workflow completes:

- Confirm all six release binaries exist.
- Confirm `checksums.txt` exists.
- Confirm the release workflow validation step passed.
- Run the `Release smoke test` workflow for `0.8.0`.
- Verify the GitHub Action example works against `hermesscan/hermesscan@v0.8.0`.
- Verify one manual install path using the checksum instructions in `docs/install.md`.
