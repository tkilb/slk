// internal/ui/mode_sidebar_filter.go
//
// Fork feature sidebar-search: `/` while the sidebar has focus filters the
// sidebar's channels by name as the user types. Enter keeps the filter and
// returns to Normal mode so j/k/Enter work on the matches; Esc clears it. The
// filter also clears when focus leaves the sidebar (sidebar.SetFocused).
//
// The query lives in the sidebar (sidebar.Model.Filter), not on App.
package ui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// forkSlashSearch routes `/` by the focused panel. It reports false when the
// panel's `/` is not a fork behavior and the stock handler should run.
func (a *App) forkSlashSearch() bool {
	if a.focusedPanel != PanelSidebar {
		return false
	}
	a.statusbar.SetSearch("/" + a.sidebar.Filter())
	a.SetMode(ModeSidebarFilter)
	return true
}

func handleSidebarFilterMode(a *App, msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, a.keys.Escape):
		a.sidebar.ClearFilter()
		a.restoreSearchSegment()
		a.SetMode(ModeNormal)
		return nil
	case key.Matches(msg, a.keys.Enter):
		a.restoreSearchSegment()
		a.SetMode(ModeNormal)
		return nil
	}
	query := a.sidebar.Filter()
	switch s := normalizeFinderKey(msg); s {
	case "backspace":
		if r := []rune(query); len(r) > 0 {
			query = string(r[:len(r)-1])
		}
	case "space":
		query += " "
	default:
		// Single printable rune only, as in handleSearchMode.
		if len([]rune(s)) != 1 {
			return nil
		}
		query += s
	}
	a.sidebar.SetFilter(query)
	a.sidebar.SelectFirstChannel()
	a.statusbar.SetSearch("/" + query)
	return nil
}

// restoreSearchSegment puts back the status-line segment of an active
// channel search (if any) after the filter prompt used it.
func (a *App) restoreSearchSegment() {
	if s := a.search; s != nil {
		a.statusbar.SetSearch(fmt.Sprintf("/%s  %d/%d", s.query, s.idx+1, len(s.matches)))
		return
	}
	a.statusbar.SetSearch("")
}

// clearSidebarFilterOnEsc is Normal mode's Esc hook: it clears a committed
// sidebar filter and reports whether it did, so the stock Esc steps skip.
func (a *App) clearSidebarFilterOnEsc() bool {
	if a.focusedPanel != PanelSidebar || a.sidebar.Filter() == "" {
		return false
	}
	a.sidebar.ClearFilter()
	return true
}
