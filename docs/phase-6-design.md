# Phase 6 design notes

Phase 6 moves HermesScan from local CLI utility toward easier GitHub adoption and release packaging.

## Goals

- Normalize GitHub annotation paths to forward slashes for portable workflow output.
- Provide a first-party composite GitHub Action wrapper.
- Document Windows, Linux, and macOS installation.
- Document SARIF and pull-request annotation workflows.
- Replace the placeholder Go module path with the project module path.
- Add contribution guidance for rules, tests, and release checks.

## Non-goals

- No marketplace publishing automation yet.
- No Homebrew tap yet.
- No package manager installers yet.

## Release shape

The release workflow builds these assets:

```text
hermesscan-windows-amd64.exe
hermesscan-windows-arm64.exe
hermesscan-linux-amd64
hermesscan-linux-arm64
hermesscan-darwin-amd64
hermesscan-darwin-arm64
checksums.txt
```

The composite action downloads those release assets and invokes the binary with the requested scan options.
