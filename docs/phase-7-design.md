# Phase 7 Design Notes - HermesScan 0.7.0

HermesScan 0.7.0 focuses on rule-quality and adoption polish after the public `v0.6.1` baseline.

## Goals

- Reduce obvious false positives in high-noise rules.
- Make rule discovery and documentation easier.
- Allow targeted scans by rule ID, category, and tag.
- Support configuration-based rule subsets for repositories that want a narrower policy.

## Changes

- Added `scan --rule RULE_ID` for one or more selected rules.
- Added `.hermesscan.json` `enabledRules`, `categories`, and `tags` filters.
- Added `hermesscan rules docs` to generate Markdown rule references from the active rule catalog.
- Added `hermesscan rules categories` and `hermesscan rules tags` inventory commands.
- Refined `HMS0002` to focus on exposed PostgreSQL port contexts.
- Lowered `HMS0010` to `Low` because package installation alone is advisory.
- Added `HMS0015` for self-hosted GitHub Actions runners.
- Added `HMS0016` for overly broad GitHub Actions cache keys.

## Non-goals

- Full YAML parsing.
- Full shell parsing.
- Dynamic execution or runtime tracing.
- Replacement for GitHub CodeQL, ShellCheck, or actionlint.

## Future work

- Add structured parser-backed rules where regex rules become too noisy.
- Add `rules validate` for external rule packs.
- Add explicit rule confidence metadata.
- Add rule documentation generated during release.
