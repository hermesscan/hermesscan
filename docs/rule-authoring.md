# Rule authoring guide

HermesScan rules are JSON records loaded from `rules/hermes.rules.json` and embedded from `internal/rules/defaults/hermes.rules.json`.
Rules are line-oriented regular expressions with optional context filters for reducing known false positives.

## Rule fields

Required fields:

| Field | Purpose |
|---|---|
| `id` | Stable rule identifier, for example `HMS0016`. Do not reuse retired IDs. |
| `name` | Short display name used in reports. |
| `severity` | One of `Info`, `Low`, `Medium`, `High`, or `Critical`. |
| `category` | Broad group such as `reliability`, `isolation`, `cache`, or `supply-chain`. |
| `tags` | Searchable rule labels. Include domain tags such as `docker`, `github-actions`, or `package-manager`. |
| `fileTypes` | Candidate file types where the rule applies. |
| `pattern` | Primary Go regular expression matched against one line at a time. |
| `description` | What risk the rule identifies. |
| `recommendation` | What the user should do instead. |

Optional precision fields:

| Field | Purpose |
|---|---|
| `excludePattern` | Go regular expression matched against the same line. If it matches, the finding is skipped. |
| `contextBeforePattern` | Go regular expression matched against previous lines. If it matches within the configured window, the finding is skipped. |
| `contextBeforeLines` | Number of previous lines to inspect when `contextBeforePattern` is set. Must be at least `1`. |
| `triggerFilePattern` | Go regular expression that must appear somewhere in the same file before `pattern` matches are reported. |
| `requiredFilePattern` | Go regular expression that must appear somewhere in the same file. If `pattern` matches but this file-level pattern is absent, HermesScan reports the finding once for that file. |

## Pattern choice

Use only `pattern` when the risk signal is specific enough by itself.

Use `excludePattern` when a safe or reviewed form appears on the same line as the risky command. Examples:

- `HMS0007` ignores `docker compose up` when `-p`, `--project-name`, or `COMPOSE_PROJECT_NAME=...` appears on the same line.
- `HMS0010` ignores package installs with same-line cache isolation such as `npm ci --cache "$RUNNER_TEMP/npm-cache"`.
- `HMS0016` ignores cache keys that contain specificity signals such as `hashFiles(...)`, `matrix.`, `github.ref`, or lockfile names.

Use `contextBeforePattern` when setup on a nearby previous line makes the matched command safer. Examples:

- `HMS0007` ignores `docker compose up` when `COMPOSE_PROJECT_NAME` is set shortly before the command.
- `HMS0010` ignores `pip install` when `PIP_CACHE_DIR` is set shortly before the install and points at a run-scoped temporary location.

Avoid using context windows as a substitute for parsing a full file. If a rule needs deep YAML, shell, or PowerShell semantics, keep the regex rule conservative and add parser-backed scanner behavior later.

Use `triggerFilePattern` when a line-level risk is meaningful only if another signal appears anywhere in the same file. For example, `HMS0019` flags `permissions: write-all` only when the workflow also appears to publish release assets.

Use `requiredFilePattern` for conservative absence checks where a local trigger is meaningful only if another file-level signal is missing. For example, `HMS0017` flags release workflows that publish binaries, checksums, or release assets when the same workflow does not mention an SBOM, SPDX, CycloneDX, or Syft output. `HMS0018` uses the same pattern to flag release workflows that publish assets without checksum generation. Keep these rules advisory unless the absence signal is very specific.

## Severity

Prefer conservative severity while a rule is broad:

| Severity | Guidance |
|---|---|
| `Critical` | Pattern is almost always dangerous and should normally block CI. |
| `High` | Strong reliability or supply-chain risk with low expected noise. |
| `Medium` | Risky pattern that frequently needs review or local policy. |
| `Low` | Advisory signal, hygiene issue, or pattern that needs surrounding evidence. |
| `Info` | Educational note or optional improvement. |

Lower severity is appropriate when a finding is useful but incomplete without stronger context. `HMS0010` is `Low` because package installation alone does not prove a shared-cache collision.

## Change checklist

When adding or changing a rule:

1. Update `rules/hermes.rules.json`.
2. Update `internal/rules/defaults/hermes.rules.json` with the same rule catalog change.
3. Add or update precision tests in `internal/scanner/rule_precision_test.go`.
4. Run `go test ./...`.
5. Build a local CLI and run `hermesscan rules validate`.
6. Regenerate `docs/rules.md` with `hermesscan rules docs --output docs/rules.md`.
7. Run `.\scripts\Test-HermesScanQuality.ps1` on Windows.

The test suite includes a catalog sync check, so rule changes fail if the repository catalog and embedded default catalog drift.

## Example rule

```json
{
  "id": "HMS0016",
  "name": "GitHub Actions broad cache key",
  "severity": "Medium",
  "category": "cache",
  "tags": ["github-actions", "cache"],
  "fileTypes": ["yaml"],
  "pattern": "(?i)^\\s*key:\\s*[^\\r\\n]*\\$\\{\\{\\s*runner\\.os\\s*\\}\\}[^\\r\\n]*$",
  "excludePattern": "(?i)hashFiles\\s*\\(|package-lock\\.json|github\\.(?:ref|sha)|matrix\\.",
  "description": "Broad cache keys based mostly on runner OS can cause unrelated branches or dependency states to reuse the same writable cache namespace.",
  "recommendation": "Include dependency lockfile hashes, matrix dimensions, or source revision inputs in cache keys."
}
```

This rule flags short cache keys based on `runner.os`, but it skips keys with lockfile hashes, matrix dimensions, or source revision context.
