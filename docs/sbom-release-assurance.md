# SBOM and release assurance

HermesScan is not an SBOM generator. It complements SBOM tooling by checking the CI and release workflows that generate, retain, and publish SBOM artifacts.

SBOM tools describe what is in released software. HermesScan checks whether the release workflow has reliability and supply-chain patterns that affect whether those artifacts can be trusted and found later.

## Recommended release outputs

A release workflow should publish these artifacts together:

- release binaries or packages
- `checksums.txt`
- an SPDX or CycloneDX SBOM
- scan reports such as SARIF or Markdown when applicable

The checksum manifest should cover the SBOM as well as the binaries so consumers can verify every released artifact.

## Complete GitHub Actions example

Use one release job to produce the binaries, HermesScan reports, SBOM, and checksum manifest before publishing the release. The report steps use `--no-fail` so SARIF and Markdown evidence is still created; the separate summary gate decides whether the release may continue.

```yaml
name: Release

on:
  push:
    tags:
      - 'v*.*.*'
  workflow_dispatch:
    inputs:
      version:
        description: 'Version to build, for example 0.9.0'
        required: true
        default: '0.9.0'

permissions:
  contents: write
  actions: read
  security-events: write

jobs:
  release:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22.x'

      - name: Determine version
        id: version
        shell: bash
        run: |
          if [ "${{ github.event_name }}" = "workflow_dispatch" ]; then
            VERSION="${{ github.event.inputs.version }}"
          else
            VERSION="${GITHUB_REF_NAME#v}"
          fi
          echo "version=$VERSION" >> "$GITHUB_OUTPUT"

      - name: Test
        run: go test ./...

      - name: Vet
        run: go vet ./...

      - name: Build release binaries
        shell: bash
        run: |
          set -euo pipefail
          mkdir -p dist
          build() {
            local goos="$1"
            local goarch="$2"
            local output="$3"
            GOOS="$goos" GOARCH="$goarch" go build -ldflags "-X main.version=${{ steps.version.outputs.version }}" -o "dist/$output" ./cmd/hermesscan
          }
          build windows amd64 hermesscan-windows-amd64.exe
          build windows arm64 hermesscan-windows-arm64.exe
          build linux amd64 hermesscan-linux-amd64
          build linux arm64 hermesscan-linux-arm64
          build darwin amd64 hermesscan-darwin-amd64
          build darwin arm64 hermesscan-darwin-arm64

      - name: Validate release CLI
        run: |
          ./dist/hermesscan-linux-amd64 version
          ./dist/hermesscan-linux-amd64 rules validate
          ./dist/hermesscan-linux-amd64 scan ./examples --summary --no-fail

      - name: Generate HermesScan SARIF report
        run: |
          ./dist/hermesscan-linux-amd64 scan . \
            --format sarif \
            --output dist/hermesscan.sarif \
            --exclude 'dist/**' \
            --no-fail

      - name: Upload HermesScan SARIF
        uses: github/codeql-action/upload-sarif@v4
        with:
          sarif_file: dist/hermesscan.sarif

      - name: Generate HermesScan Markdown report
        run: |
          ./dist/hermesscan-linux-amd64 scan . \
            --format markdown \
            --output dist/hermesscan.md \
            --exclude 'dist/**' \
            --no-fail

      - name: Generate source SBOM
        uses: anchore/sbom-action@v0
        with:
          path: .
          format: spdx-json
          output-file: dist/hermesscan.spdx.json
          upload-artifact: false
          upload-release-assets: false

      - name: Generate checksums
        shell: bash
        run: |
          set -euo pipefail
          cd dist
          sha256sum hermesscan-* hermesscan.spdx.json hermesscan.sarif hermesscan.md > checksums.txt

      - name: Store release evidence
        uses: actions/upload-artifact@v4
        with:
          name: hermesscan-release-evidence
          path: dist/*

      - name: Gate release with HermesScan
        run: ./dist/hermesscan-linux-amd64 scan . --summary --fail-on high --exclude 'dist/**'

      - name: Publish release
        if: startsWith(github.ref, 'refs/tags/')
        uses: softprops/action-gh-release@v2
        with:
          files: dist/*
          generate_release_notes: true
```

Adapt the build commands and asset names for non-Go projects, but keep the order: build, generate reports, generate SBOM, create checksums that cover every release artifact, retain the evidence, run the blocking gate, and then publish.

## HermesScan rule coverage

`HMS0017` flags GitHub release workflows that appear to publish release assets or checksum manifests without any SBOM, SPDX, CycloneDX, or Syft signal in the same workflow file.

This is an advisory rule. It does not prove a project lacks an SBOM across all automation. It identifies workflows that are likely responsible for release assets and should usually publish an SBOM with those assets.

Existing supply-chain rules also support SBOM adoption:

- `HMS0009` flags mutable GitHub Action references such as `@main`.
- `HMS0013` flags risky `pull_request_target` usage.
- `HMS0014` flags broad `permissions: write-all`.
- `HMS0016` flags overly broad cache keys that can affect reproducible release jobs.

## Adoption guidance

Start with `HMS0017` as advisory:

```powershell
.\hermesscan.exe scan . --rule HMS0017 --summary --no-fail
```

After the release workflow publishes an SBOM and checksum manifest together, include `HMS0017` in regular supply-chain scans:

```powershell
.\hermesscan.exe scan . --category supply-chain --summary --fail-on high
```

Keep SBOM generation close to the release job. If SBOMs are generated in a separate workflow, document how they are attached to the same release and retained with the release assets.
