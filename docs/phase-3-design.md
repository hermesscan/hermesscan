# HermesScan Phase 3 Design Notes

Phase 3 turns the prototype into a release-ready developer tool.

## Goals

- Improve local developer experience on Windows, Linux, and macOS.
- Make GitHub Code Scanning integration practical through SARIF.
- Make release artifacts easy to produce and verify.
- Preserve advisory adoption mode through `--no-fail`.
- Add include/exclude command-line filters for quick experimentation.

## New CLI behavior

```text
hermesscan scan . --no-fail
hermesscan scan . --include "scripts/**" --exclude "dist/**"
hermesscan scan . --format sarif --output reports/hermes-scan.sarif
```

`--no-fail` overrides both `.hermesscan.json` `failOn` and the command-line `--fail-on` setting.
This is useful for generating reports in CI without blocking a workflow.

## Release readiness

The release workflow builds these artifacts:

```text
hermesscan-windows-amd64.exe
hermesscan-windows-arm64.exe
hermesscan-linux-amd64
hermesscan-linux-arm64
hermesscan-darwin-amd64
hermesscan-darwin-arm64
checksums.txt
```

The PowerShell build script also writes per-file `.sha256` checksum files.

## Windows behavior

Local Windows runs should use explicit relative execution:

```powershell
.\hermesscan.exe scan .
```

Bare `hermesscan` works only after the executable is installed in a directory listed in `$env:Path`.
