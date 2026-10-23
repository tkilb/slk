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

func handleSidebarFilterMode(a *App, msg tea.KeyMsg) tea.Cmd {
	if a.focusedPanel != PanelSidebar {
		// Focus moved by a non-key path (a mouse click): the prompt is stale.
		a.restoreSearchSegment()
		a.SetMode(ModeNormal)
		return handleNormalMode(a, msg)
	}
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
	query, ok := editQuery(a.sidebar.Filter(), msg)
	if !ok {
		return nil
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
