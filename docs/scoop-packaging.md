# Scoop packaging plan

HermesScan should support Scoop after the release asset contract is stable.

This plan is for a future Scoop manifest. It does not publish a bucket yet.

## Package shape

The Scoop package should install the existing Windows release binary directly from GitHub Releases:

- package name: `hermesscan`
- executable shim: `hermesscan.exe`
- amd64 asset: `hermesscan-windows-amd64.exe`
- arm64 asset: `hermesscan-windows-arm64.exe`
- release URL: `https://github.com/hermesscan/hermesscan/releases/download/v<version>/<asset>`
- checksum source: release `checksums.txt`

Do not introduce a Windows zip archive only for Scoop unless Scoop requires it. The current `.exe` release assets are simple and match the manual install path.

## Manifest fields

Initial manifest shape:

```json
{
  "version": "0.8.0",
  "description": "Static analyzer for build scripts, CI scripts, and pipeline definitions.",
  "homepage": "https://github.com/hermesscan/hermesscan",
  "license": "MIT",
  "architecture": {
    "64bit": {
      "url": "https://github.com/hermesscan/hermesscan/releases/download/v0.8.0/hermesscan-windows-amd64.exe#/hermesscan.exe",
      "hash": "<sha256 from checksums.txt>"
    },
    "arm64": {
      "url": "https://github.com/hermesscan/hermesscan/releases/download/v0.8.0/hermesscan-windows-arm64.exe#/hermesscan.exe",
      "hash": "<sha256 from checksums.txt>"
    }
  },
  "bin": "hermesscan.exe",
  "checkver": {
    "github": "https://github.com/hermesscan/hermesscan"
  },
  "autoupdate": {
    "architecture": {
      "64bit": {
        "url": "https://github.com/hermesscan/hermesscan/releases/download/v$version/hermesscan-windows-amd64.exe#/hermesscan.exe"
      },
      "arm64": {
        "url": "https://github.com/hermesscan/hermesscan/releases/download/v$version/hermesscan-windows-arm64.exe#/hermesscan.exe"
      }
    }
  }
}
```

Confirm the final `autoupdate` hash behavior against Scoop tooling before publishing. If Scoop cannot derive hashes from `checksums.txt`, the update script should download `checksums.txt`, extract the matching Windows hashes, and update the manifest explicitly.

## Release requirements

Every release intended for Scoop should include:

- `hermesscan-windows-amd64.exe`
- `hermesscan-windows-arm64.exe`
- `checksums.txt`

The current release workflow already publishes these assets. v0.9 and later releases also publish `hermesscan.spdx.json`; Scoop does not need to install the SBOM, but the SBOM should remain part of the release artifact set.

## Update workflow

The update workflow should:

1. Read the target version.
2. Download release `checksums.txt`.
3. Extract the hashes for `hermesscan-windows-amd64.exe` and `hermesscan-windows-arm64.exe`.
4. Update manifest `version`, URLs, and hashes.
5. Run Scoop manifest validation.
6. Smoke-test `scoop install` from the local manifest when Scoop is available.

Keep the first implementation manual or semi-automated. Do not publish to a public bucket until the manifest has been tested against at least one real release.

## Open decisions

- Use a project-owned bucket or contribute to an existing bucket.
- Decide whether to keep the manifest in this repository first or maintain it in a bucket repository only.
- Decide whether release automation should open a manifest update PR after publishing a new release.
