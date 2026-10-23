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
  (the active channel row stays visible, headers keep their unread badges).
  While it is focused, only the section holding the cursor is expanded, so
  cursoring across sections opens one and closes the other; a filter shows
  matches in every section. Manual toggles are disabled in this mode. The
  cursor returns to the active channel after a blur.
- Files: internal/ui/sidebar/model.go, internal/ui/sidebar/collapse_unfocused_test.go,
  internal/ui/fork_sidebar.go, cmd/slk/main.go
- Upstream-watch: collapse sidebar, auto-collapse, sidebar focus, collapse sections

## sidebar-search: `/` filters the sidebar when it has focus
- Status: active
- Added: 2026-10-23
- Summary: `/` is context sensitive. With the sidebar focused it opens a live
  name filter (Enter keeps it for j/k/Enter on the matches, Esc clears it, and
  it clears when focus leaves). With the messages pane focused it is the stock
  channel search. The thread pane is covered by `thread-search`.
- Files: internal/ui/mode_sidebar_filter.go, internal/ui/mode_sidebar_filter_test.go,
  internal/ui/fork_search.go (shared `/` routing), internal/ui/mode_normal.go (hook calls), internal/ui/mode.go,
  internal/ui/mode_handlers.go, internal/ui/sidebar/model.go,
  internal/ui/reducer_search_test.go and internal/ui/modekeys_test.go (fixture updates)
- Upstream-watch: sidebar search, sidebar filter, filter channels, search sidebar

## thread-search: `/` searches the thread when the thread pane has focus
- Status: active
- Added: 2026-10-23
- Summary: With the thread pane focused, `/` searches the parent and loaded
  replies as you type (client-side, word-prefix, case and accent insensitive,
  all terms must match), highlights hits, and jumps the cursor. Enter keeps the
  search so `n`/`N` step through matches with wraparound; the first Esc clears
  it and a second Esc closes the thread as before.
- Files: internal/ui/thread/search_fork.go, internal/ui/thread/search_fork_test.go,
  internal/ui/thread/model.go (state field, reset in SetThread/Clear, render terms),
  internal/ui/fork_search.go, internal/ui/mode_thread_search.go,
  internal/ui/mode_thread_search_test.go, internal/ui/mode_normal.go (hooks),
  internal/ui/mode_normal_keys_test.go (pinned `/` case updated)
- Upstream-watch: thread search, search thread, search in thread, find in thread

## sidebar-section-jump: `{` and `}` jump between sidebar categories
- Status: active
- Added: 2026-10-23
- Summary: With the sidebar focused, `}` moves the cursor to the next section
  header and `{` to the previous section's header (from a channel row, the
  section before the one it is in). Jumps stop at either end. Other panels
  ignore the keys.
- Files: internal/ui/sidebar/model.go (`JumpSection`), internal/ui/sidebar/jump_section_fork_test.go,
  internal/ui/fork_search.go, internal/ui/mode_section_jump_test.go,
  internal/ui/mode_normal.go (hook)
- Upstream-watch: next section, previous section, jump section, jump category

## nav-history-keys: `ctrl+o` goes back and `ctrl+i` goes forward
- Status: active
- Added: 2026-10-23
- Summary: Vim-style aliases for navigation history. `ctrl+o` is added to
  `ctrl+h` (back) and `ctrl+i` to `ctrl+l` / `ctrl+k` (forward); the originals
  still work and the help overlay lists the new keys. `ctrl+i` is a distinct
  key only in terminals with key disambiguation (kitty protocol: kitty,
  ghostty, wezterm, recent iTerm2 and others); elsewhere it arrives as Tab,
  which keeps its "next panel" meaning.
- Files: internal/ui/keys.go (two bindings), internal/ui/nav_keys_fork_test.go,
  internal/ui/mode_help_test.go (pinned help label updated)
- Upstream-watch: ctrl+o, ctrl+i, jump list, navigate back, nav history key

## close-dm: `x` in the sidebar and `/close` close a DM
- Status: active
- Added: 2026-10-23
- Summary: With the sidebar focused and the cursor on a DM or group DM, `x`
  closes it (Slack's conversations.close) with no confirmation; on any other
  channel it explains that only DMs can be closed. `/close` does the same for
  the conversation it is typed in. The row leaves the sidebar when Slack
  confirms, a closed DM stays in the finder, and a new message reopens it.
  Closes made from other Slack clients (im_close, mpim_close) also remove the
  row. If the active DM is closed its row goes but the pane stays until you
  navigate away.
- Files: internal/slack/close_fork.go (+ test), internal/slack/client.go (interface
  method), internal/slack/events.go (one case), internal/core/ports.go and adapters.go,
  internal/ui/close_conversation.go (+ test), internal/ui/sidebar/remove_item_fork.go (+ test),
  internal/ui/app.go (reducer chain), internal/ui/mode_normal.go (key case),
  internal/ui/slash_command.go (`/close`), cmd/slk/close_conversation.go (+ test), cmd/slk/main.go (one line),
  pinned tests updated in internal/slack/client_test.go and internal/ui/slash_builtins_test.go
- Upstream-watch: close dm, close conversation, conversations.close, /close, hide dm

## Retired

None.
