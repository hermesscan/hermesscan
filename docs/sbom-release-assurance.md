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

## GitHub Actions example

```yaml
permissions:
  contents: write
  actions: read

steps:
  - name: Checkout repository
    uses: actions/checkout@v4

  - name: Build release binaries
    run: ./scripts/build-release.sh

  - name: Generate SBOM
    uses: anchore/sbom-action@v0
    with:
      path: .
      format: spdx-json
      output-file: dist/project.spdx.json
      upload-artifact: false
      upload-release-assets: false

  - name: Generate checksums
    run: |
      cd dist
      sha256sum project-* project.spdx.json > checksums.txt

  - name: Publish release
    uses: softprops/action-gh-release@v2
    with:
      files: dist/*
```

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
