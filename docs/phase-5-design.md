# HermesScan Phase 5 Design

Phase 5 focuses on real repository usability instead of new scanner theory.

## Features

- `rules list` now prints an aligned table.
- `rules show RULE_ID` explains one rule in detail.
- `scan --category NAME` filters the loaded rule set by category.
- `scan --tag NAME` filters the loaded rule set by tag.
- `scan --changed-only` uses Git to restrict candidate files to changed files.
- `scan --changed-base REF` selects the Git base used by changed-only scans.
- `scan --github-annotations` emits GitHub Actions workflow command annotations.
- `scripts/Install-HermesScan.ps1` installs a local Windows executable and can add the destination directory to the user PATH.

## Changed-file behavior

`--changed-only` currently shells out to Git and runs:

```text
git -C <root> diff --name-only HEAD
```

When `--changed-base` is supplied, the base replaces `HEAD`.

## Annotation behavior

GitHub annotation levels map as follows:

| HermesScan severity | GitHub level |
|---|---|
| Critical | error |
| High | error |
| Medium | warning |
| Low | notice |
| Info | notice |

## Remaining improvement areas

- Add config-level category/tag filters.
- Add changed-file support for staged-only and untracked files.
- Add better context-aware rules to reduce false positives.
- Add a first-party GitHub Action wrapper.
