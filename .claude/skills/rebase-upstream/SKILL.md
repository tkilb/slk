---
name: rebase-upstream
description: Sync this fork with gammons/slk by rebasing main onto upstream/main, resolving conflicts, checking whether upstream now supports any fork feature, and updating FORK_FEATURES.md. Use when the user asks to sync, update, rebase, or pull in upstream changes.
disable-model-invocation: true
---

# Rebase the fork onto upstream

Fork `main` carries custom commits on top of `upstream/main` (`gammons/slk`).
History stays linear: rebase, never merge. Read `FORK_FEATURES.md` first.

## Confirmation gates

Ask the user before each of these, once per run: `git fetch upstream`,
the rebase itself, and the force push. Never push to `upstream`.

## Steps

1. **Preflight.** Working tree must be clean (`git status --short`). Work on a
   branch, not `main`: `git switch -c sync/upstream-<YYYYMMDD> main`. If the
   `upstream` remote is missing, add it:
   `git remote add upstream https://github.com/gammons/slk.git`.
2. **Fetch.** `git fetch upstream main`.
3. **Preview.** Report `git log --oneline main..upstream/main | wc -l` new
   upstream commits and the fork-only commits (`git log --oneline
   upstream/main..main`). Stop if upstream has nothing new.
4. **Overlap check before rebasing.** For each `Status: active` feature in
   `FORK_FEATURES.md`, search the new upstream commits (`git log
   main..upstream/main`, plus `git diff --stat` on the feature's `Files:`) for
   the same capability. Judge by reading the diff, not only keywords. Report
   each feature as: still unique, partially covered, or now supported upstream.
5. **Rebase.** `git rebase upstream/main`. On conflict:
   - Identify the fork commit (`Fork-Feature:` trailer) and its entry in
     `FORK_FEATURES.md`.
   - Keep upstream's structure and reapply the fork behavior on top of it.
     Honor the upstream invariants in `AGENTS.md` (reducers, mode handler
     table, no I/O in `internal/ui`, shared helpers).
   - Do not resolve a conflict by dropping upstream changes.
   - If a commit is wholly superseded, use `git rebase --skip` and retire the
     feature (step 7).
6. **Verify.** Run `go build ./...`, `go vet ./...`, `gofmt -l .` (must print
   nothing), then `go test ./... -race`. Fix fork-caused failures. Report
   failures that also occur on `upstream/main` instead of fixing them.
7. **Update the list.** In `FORK_FEATURES.md`: refresh `Files:` if code moved;
   for a feature now supported upstream, set `Status: superseded-by-upstream`,
   add `Upstream: <PR or commit URL>`, remove the fork code in its own commit
   (`Fork-Feature:` trailer kept), and move the entry under "Retired". Commit as
   `docs: update fork features after upstream sync`.
8. **Publish.** Show the user `git range-diff upstream/main main <branch>`
   summary. After confirmation: fast-forward `main` to the branch, then
   `git push --force-with-lease origin main`. Never plain `--force`.
9. **Report.** Lead with the result: upstream commits absorbed, conflicts
   resolved, features retired or at risk, test status. Offer to close any open
   `upstream-overlap` issue for retired features.

## Recovery

If the rebase goes wrong: `git rebase --abort`, or return to the old tip with
`git reset --hard ORIG_HEAD` on the sync branch. `main` is untouched until
step 8.
