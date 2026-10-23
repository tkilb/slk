// internal/ui/mode_thread_search.go
//
// Fork feature thread-search: `/` with the thread pane focused searches the
// parent and loaded replies as the user types, with no server call. Enter
// keeps the search so n/N step through matches; Esc clears it. The query and
// matches live in the thread model (thread.Model.SetSearch).
package ui

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func handleThreadSearchMode(a *App, msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, a.keys.Escape):
		a.threadPanel.ClearSearch()
		a.restoreSearchSegment()
		a.SetMode(ModeNormal)
		return nil
	case key.Matches(msg, a.keys.Enter):
		if a.threadPanel.SearchActive() {
			a.statusbar.SetSearch(a.threadSearchSegment())
		} else {
			a.restoreSearchSegment()
		}
		a.SetMode(ModeNormal)
		return nil
	}
	query := a.threadPanel.SearchQuery()
	switch s := normalizeFinderKey(msg); s {
	case "backspace":
		if r := []rune(query); len(r) > 0 {
			query = string(r[:len(r)-1])
		}
	case "space":
		query += " "
	default:
		if len([]rune(s)) != 1 {
			return nil
		}
		query += s
	}
	a.threadPanel.SetSearch(query)
	if a.threadPanel.SearchActive() {
		a.statusbar.SetSearch(a.threadSearchSegment())
	} else {
		a.statusbar.SetSearch("/")
	}
	return nil
}

// threadSearchSegment is the status-line text for the active thread search:
// "/query  i/N", or "/query  no matches".
func (a *App) threadSearchSegment() string {
	idx, n := a.threadPanel.SearchStatus()
	if n == 0 {
		return "/" + a.threadPanel.SearchQuery() + "  no matches"
	}
	return fmt.Sprintf("/%s  %d/%d", a.threadPanel.SearchQuery(), idx, n)
}
