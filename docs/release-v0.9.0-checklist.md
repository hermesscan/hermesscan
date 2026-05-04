# v0.9.0 release checklist

Use this checklist before creating the `v0.9.0` tag.

## Local validation

```powershell
go test .\...
go vet .\...
go build -ldflags "-X main.version=0.9.0" -o .\hermesscan.exe .\cmd\hermesscan
.\hermesscan.exe version
.\hermesscan.exe rules validate
.\hermesscan.exe rules docs --output .\docs\rules.md
.\hermesscan.exe scan . --summary --fail-on high --exclude "examples/**" --exclude "dist/**"
.\hermesscan.exe scan .\examples --summary --no-fail
```

Expected version output:

```text
HermesScan 0.9.0
```

Expected self-scan result, excluding intentionally risky examples:

```text
HermesScan: 0 findings
```

The examples scan should continue to report the intentional fixture findings.

## Release artifact validation

Build every target locally before tagging:

```powershell
.\scripts\Build-HermesScan.ps1 -AllTargets -Version 0.9.0
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

The local build script writes per-file `.sha256` files. The GitHub release workflow publishes the final `checksums.txt`, source SBOM, and HermesScan scan reports.

## Documentation checks

- `README.md` describes the `0.9.0` development focus.
- `CHANGELOG.md` has the finalized `0.9.0` entry and date.
- `action.yml` defaults the `version` input to `0.9.0` before tagging.
- README GitHub Action examples point at `hermesscan/hermesscan@v0.9.0` before tagging.
- `docs/rules.md` has been regenerated from the active rule catalog.
- `docs/sbom-release-assurance.md` includes the complete release-assurance workflow example.
- `docs/scoop-packaging.md` describes the manifest shape, checksum source, and update workflow.
- `packaging/scoop/hermesscan.json` still points at the latest published release until `v0.9.0` assets exist.
- `docs/release-v0.9.0.md` matches the release contents.

## Git checks

```powershell
git status
git diff --stat
```

Confirm no generated binaries or local reports are staged:

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
git commit -m "Prepare HermesScan v0.9.0"
git push
```

## Tag

Create the tag only after CI passes on `main`.

```powershell
git tag v0.9.0
git push origin v0.9.0
```

## Release verification

After the release workflow completes:

- Confirm all six release binaries exist.
- Confirm `hermesscan.spdx.json` exists.
- Confirm `hermesscan.sarif` exists.
- Confirm `hermesscan.md` exists.
- Confirm `checksums.txt` includes every binary, SBOM, SARIF, and Markdown report asset.
- Confirm the release workflow validation and HermesScan gate steps passed.
- Run the `Release smoke test` workflow for `0.9.0`.
- Verify the GitHub Action example works against `hermesscan/hermesscan@v0.9.0`.
- Verify one manual install path using the checksum instructions in `docs/install.md`.

## Scoop verification

After `v0.9.0` release assets are available, refresh and validate the local Scoop manifest:

```powershell
.\scripts\Update-ScoopManifest.ps1 -Version 0.9.0
$null = Get-Content .\packaging\scoop\hermesscan.json -Raw | ConvertFrom-Json
scoop download .\packaging\scoop\hermesscan.json --force --no-update-scoop --arch 64bit
scoop download .\packaging\scoop\hermesscan.json --force --no-update-scoop --arch arm64
if ((scoop list 2>$null) -match '^hermesscan\s') { scoop uninstall hermesscan }
scoop install .\packaging\scoop\hermesscan.json
hermesscan version
```

Expected version output:

```text
HermesScan 0.9.0
```

Commit the manifest refresh after validation.
