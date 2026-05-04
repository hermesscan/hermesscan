# HermesScan v0.9.0 release notes

HermesScan v0.9.0 focuses on SBOM-aware release assurance and the first Windows package-manager installation path.

## Highlights

- Added file-level rule support with `requiredFilePattern` for conservative absence checks.
- Added `HMS0017` to flag release workflows that publish release assets without an SBOM signal.
- Added SBOM generation to the release workflow and included the SBOM in release checksum and smoke-test verification.
- Added HermesScan SARIF and Markdown release evidence to release assets, checksum coverage, and release smoke verification.
- Added workflow contract tests for SBOM and report release assets.
- Added a complete SBOM and release assurance guide.
- Added a local Scoop manifest prototype and manifest refresh helper.
- Documented raw GitHub Scoop manifest installation before publishing a Scoop bucket.

## Upgrade notes

After the `v0.9.0` release is published, pin the GitHub Action to:

```yaml
uses: hermesscan/hermesscan@v0.9.0
```

The action downloads the matching `0.9.0` CLI by default unless the `version` input is overridden.

Release consumers can verify binaries and release evidence against `checksums.txt`. Starting with v0.9.0, the release workflow is expected to publish:

- six native binaries
- `hermesscan.spdx.json`
- `hermesscan.sarif`
- `hermesscan.md`
- `checksums.txt`

## Scoop

Until HermesScan has a dedicated Scoop bucket, install from the raw manifest:

```powershell
scoop install https://raw.githubusercontent.com/hermesscan/hermesscan/main/packaging/scoop/hermesscan.json
hermesscan version
```

The manifest should be refreshed after the `v0.9.0` release assets exist.

## Validation

Validation used before tagging:

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

Expected self-scan result, excluding intentionally risky examples:

```text
HermesScan: 0 findings
```

## GitHub Release body

Use this text for the GitHub Release body if generated notes are not sufficient:

```markdown
HermesScan v0.9.0 adds SBOM-aware release assurance and a prototype Scoop installation path.

### Added

- `requiredFilePattern` rule support for conservative file-level absence checks.
- `HMS0017` to flag release workflows that publish assets without an SBOM signal.
- Release workflow SBOM generation with checksum and smoke-test coverage.
- HermesScan SARIF and Markdown release evidence assets.
- Workflow contract tests for SBOM/report release assets and checksum coverage.
- SBOM and release assurance guide with a complete GitHub Actions example.
- Local Scoop manifest prototype and refresh helper.
- Raw GitHub Scoop manifest install documentation.

### Release assets

- Six native binaries for Windows, Linux, and macOS on amd64/arm64.
- `hermesscan.spdx.json`
- `hermesscan.sarif`
- `hermesscan.md`
- `checksums.txt`
```
