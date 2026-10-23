# Fork features

Features that exist only in `tkilb/slk` and not in `gammons/slk`. Maintained by
hand and by the `rebase-upstream` skill. The `upstream-watch` workflow reads the
`Upstream-watch:` lines below.

## Entry format

Each feature is a level-2 heading `## <id>: <title>` followed by these bullets.
The id is a short kebab-case slug, also used in the `Fork-Feature:` commit
trailer.

```
## <id>: <title>
- Status: active | superseded-by-upstream | dropped
- Added: YYYY-MM-DD
- Summary: one or two sentences on what it does.
- Files: main paths touched (helps resolve rebase conflicts)
- Upstream-watch: comma-separated keywords upstream would plausibly use
```

When upstream ships an equivalent, set `Status: superseded-by-upstream`, add an
`Upstream:` bullet with the PR or commit, remove the fork code, and move the
entry to "Retired".

## Active

None yet. The fork currently matches upstream.

## Retired

None.
