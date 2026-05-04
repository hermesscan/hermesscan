# Changelog

## 0.10.0 - Unreleased

### Added

- Started the v0.10.0 development cycle with adoption polish, distribution readiness, release-integrity rules, and configuration workflow planning.

## 0.9.0 - 2026-05-04

### Added

- Added file-level rule support with `requiredFilePattern` for conservative absence checks.
- Added `HMS0017` to flag release workflows that publish release assets without an SBOM signal.
- Added SBOM generation to the release workflow and included the SBOM in release checksum and smoke-test verification.
- Added HermesScan SARIF and Markdown release evidence to release assets, checksum coverage, and release smoke verification.
- Added workflow contract tests for SBOM/report release assets and checksum coverage.
- Added a Scoop packaging plan covering manifest shape, release assets, checksum source, and update workflow.
- Added a local Scoop manifest prototype and manifest refresh helper.
- Documented Scoop download/hash validation commands for the prototype manifest.
- Documented successful local Scoop install, shim, and runtime validation for the prototype manifest.
- Documented raw GitHub Scoop manifest installation before publishing a bucket.
- Added a complete release assurance workflow example that publishes binaries, checksums, SBOMs, and HermesScan reports together.
- Added v0.9.0 release checklist and draft release notes.

## 0.8.0 - 2026-05-03

### Added

- Added `hermesscan rules validate` to check publish-ready rule catalog metadata and regex syntax.
- Added comma-delimited multi-rule support to the GitHub Action `rule` input.
- Added tests that keep the repository rule catalog in sync with the embedded default catalog.
- Added default-rule precision tests for common false-positive cases.
- Added a rule authoring guide for rule fields, contextual precision patterns, and validation workflow.
- Added a GitHub Action workflow example that combines SARIF upload, a report artifact, and a separate summary gate.
- Expanded release smoke testing to verify every published binary asset, checksum entries, and native CLI execution.
- Expanded the baseline adoption guide with advisory, gated, and baseline reduction examples.
- Added copy-paste checksum verification examples for Windows, Linux, and macOS installs.
- Added v0.8.0 release checklist and draft release notes.

### Changed

- Refined the Docker Compose startup rule to ignore commands with explicit project-name context.
- Refined the GitHub Actions broad cache key rule to flag static `runner.os` cache keys while ignoring lockfile, hash, ref, SHA, and matrix-specific keys.
- Refined the package-install cache rule to ignore package installs with explicit runner-temp or run-scoped cache isolation.

## 0.7.0 - 2026-05-02

### Added

- Added `scan --rule RULE_ID` to scan with one or more selected rules.
- Added `.hermesscan.json` `enabledRules`, `categories`, and `tags` filters.
- Added `hermesscan rules docs` to generate Markdown rule reference documentation.
- Added `hermesscan rules categories` and `hermesscan rules tags` inventory commands.
- Added GitHub Actions rules for self-hosted runner usage and overly broad cache keys.

### Changed

- Refined the PostgreSQL fixed-port rule to focus on exposed/bound port contexts instead of every standalone `5432` token.
- Lowered the package-install cache rule to `Low` because it is advisory unless paired with shared cache or manual parallelism evidence.
- Regenerated `docs/rules.md` from the active rule catalog.

## 0.6.1 - Post-release polish

### Fixed

- Updated HermesScan SARIF upload workflow permissions by adding `actions: read`.
- Updated GitHub SARIF upload action from `github/codeql-action/upload-sarif@v3` to `github/codeql-action/upload-sarif@v4`.
- Updated GitHub Action default HermesScan CLI version to `0.6.1`.

### Added

- Added public preview and false-positive guidance.
- Added issue templates for false positives, bugs, and feature requests.
- Added baseline adoption and generated rule reference documentation.
- Added v0.7.0 milestone planning notes.


## 0.6.0 - Phase 6

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
