# Changelog

## 0.6.1 - Phase 6

- Normalized GitHub Actions annotation paths to forward slashes.
- Added first-party composite GitHub Action wrapper (`action.yml`).
- Added installation and GitHub Action documentation.
- Added `CONTRIBUTING.md`.
- Updated module path to `github.com/hermesscan/hermesscan`.
- Updated release, CI, and README guidance for release adoption.

## 0.5.0 - Phase 5

- Added `rules show RULE_ID`.
- Improved `rules list` table formatting.
- Added `scan --category` and `scan --tag` filters.
- Added `scan --changed-only` and `scan --changed-base`.
- Added `scan --github-annotations`.
- Added Windows install helper script.
- Updated docs for Phase 5 workflows.


## 0.4.0 - Phase 4

### Added

- Added `--create-baseline` to write current findings to a baseline file.
- Added `--baseline` to suppress existing findings from future scans.
- Added stable finding fingerprints to JSON output.
- Added rule `category` and `tags` metadata.
- Added `HMS0013` for `pull_request_target` workflow risks.
- Added `HMS0014` for `permissions: write-all` workflow risks.
- Added baseline suppressed counts to console, summary, Markdown, and JSON reports.

### Changed

- `rules list` now includes rule category.
- Updated documentation with baseline adoption workflow.


## 0.3.0 - Phase 3

### Added

- Added `--no-fail` to suppress failing exit codes even when `failOn` or `--fail-on` is configured.
- Added repeatable `--include` and `--exclude` scan filters.
- Added release workflow for Windows, Linux, and macOS amd64/arm64 binaries.
- Added build-time version injection using `go build -ldflags "-X main.version=<version>"`.
- Added SHA256 checksum generation to `scripts/Build-HermesScan.ps1`.
- Added default config excludes for `reports/**`, `coverage/**`, `tmp/**`, and `.git/**`.
- Added GitHub Code Scanning SARIF workflow guidance.
- Added tests for CLI include/exclude parsing and include filtering.

### Changed

- Updated README with Windows PowerShell examples and PATH guidance.
- Updated `.hermesscan.example.json` to match Phase 3 defaults.
- Updated GitHub Actions workflow to produce SARIF and run a separate gate.

## 0.2.0 - Phase 2

### Added

- Added `.hermesscan.json` configuration support.
- Added `hermesscan init`.
- Added `hermesscan rules list`.
- Added inline suppressions: `disable-next-line`, `disable-line`, and `disable-file`.
- Added automatic parent directory creation for `--output`.
- Added disabled rules and severity overrides.

## 0.1.1 - Phase 1G.1

### Added

- Added `--summary`, `--quiet`, `--min-severity`, `--no-color`, and SARIF output.

## 0.1.0 - Phase 1G

### Added

- Initial Go-native CLI.
- JSON rule loading.
- File discovery.
- Regex scanning.
- Console, Markdown, and JSON reports.
- Severity gate via `--fail-on`.
