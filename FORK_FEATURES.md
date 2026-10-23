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

## sidebar-focus-collapse: Sidebar sections collapse unless focused
- Status: active
- Added: 2026-10-23
- Summary: Every sidebar section is collapsed while the sidebar is unfocused
  (the active channel row stays visible, headers keep their unread badges) and
  expanded while the user cursors through it. Manual toggles are disabled in
  this mode. The cursor returns to the active channel after a blur.
- Files: internal/ui/sidebar/model.go, internal/ui/sidebar/collapse_unfocused_test.go,
  internal/ui/fork_sidebar.go, cmd/slk/main.go
- Upstream-watch: collapse sidebar, auto-collapse, sidebar focus, collapse sections

## Retired

None.
