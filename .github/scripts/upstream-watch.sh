#!/usr/bin/env bash
# Report upstream commits that may implement a feature listed in FORK_FEATURES.md.
#
# For each active feature, search upstream commit messages between the merge base
# with HEAD and upstream/main for any of its `Upstream-watch:` keywords. A match
# is a prompt for a human or the rebase-upstream skill to look, not a verdict.
#
# Env:
#   UPSTREAM_URL  upstream clone URL (default https://github.com/gammons/slk.git)
#   DRY_RUN=1     print findings instead of opening issues (no gh needed)
set -euo pipefail

UPSTREAM_URL="${UPSTREAM_URL:-https://github.com/gammons/slk.git}"
FEATURES_FILE="${FEATURES_FILE:-FORK_FEATURES.md}"
LABEL="upstream-overlap"

if git remote get-url upstream >/dev/null 2>&1; then
  git remote set-url upstream "$UPSTREAM_URL"
else
  git remote add upstream "$UPSTREAM_URL"
fi
git fetch --quiet upstream main

base="$(git merge-base HEAD upstream/main)"

# Emit "<id>\t<keyword>" for each keyword of each active feature.
features() {
  awk '
    /^## Active/ { on = 1; next }
    /^## /       && !/^## [a-z0-9-]+: / { on = 0 }
    on && /^## [a-z0-9-]+: / { id = $2; sub(/:$/, "", id); active = 0 }
    on && /^- Status: *active/ { active = 1 }
    on && active && /^- Upstream-watch:/ {
      line = $0; sub(/^- Upstream-watch: */, "", line)
      n = split(line, kw, /, */)
      for (i = 1; i <= n; i++) if (kw[i] != "") print id "\t" kw[i]
    }
  ' "$FEATURES_FILE"
}

# Print the matching upstream commits for one feature id, grouped by keyword.
matches_for() {
  local id="$1" kw hits
  while IFS=$'\t' read -r fid kw; do
    [ "$fid" = "$id" ] || continue
    hits="$(git log "$base..upstream/main" --regexp-ignore-case --fixed-strings \
      --grep="$kw" --format='- %h %s' || true)"
    if [ -n "$hits" ]; then
      printf 'Keyword `%s`:\n%s\n' "$kw" "$hits"
    fi
  done < <(features)
}

ids="$(features | cut -f1 | sort -u)"
opened=0

for id in $ids; do
  found="$(matches_for "$id")"
  [ -n "$found" ] || continue
  opened=1
  title="Upstream may now support fork feature: $id"
  body="Upstream commits since the merge base match this feature's \`Upstream-watch:\` keywords. Review them. If equivalent, run the \`rebase-upstream\` skill and retire the feature in \`FORK_FEATURES.md\`. Keyword matches are heuristic.

$found"
  if [ "${DRY_RUN:-0}" = "1" ]; then
    printf '%s\n%s\n\n' "$title" "$body"
    continue
  fi
  gh label create "$LABEL" --color FBCA04 \
    --description "Upstream may now support a fork feature" 2>/dev/null || true
  existing="$(gh issue list --label "$LABEL" --state open --search "\"$title\" in:title" --json number --jq '.[0].number // empty')"
  if [ -n "$existing" ]; then
    echo "Issue #$existing already open for $id; skipping."
  else
    gh issue create --title "$title" --label "$LABEL" --body "$body"
  fi
done

[ "$opened" = "1" ] || echo "No upstream overlap found."
