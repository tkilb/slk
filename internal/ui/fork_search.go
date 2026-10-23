// internal/ui/fork_search.go
//
// Fork features sidebar-search and thread-search: `/`, Esc, n and N depend on
// which panel has focus. Normal mode calls these hooks first; each reports
// whether it handled the key, and false falls through to the stock behavior
// (channel search on the messages pane).
package ui

// forkSlashSearch routes `/` by the focused panel.
func (a *App) forkSlashSearch() bool {
	switch a.focusedPanel {
	case PanelSidebar:
		a.statusbar.SetSearch("/" + a.sidebar.Filter())
		a.SetMode(ModeSidebarFilter)
		return true
	case PanelThread:
		if a.threadPanel.IsEmpty() {
			return true // nothing to search; swallow the key as before
		}
		a.statusbar.SetSearch("/" + a.threadPanel.SearchQuery())
		a.SetMode(ModeThreadSearch)
		return true
	}
	return false
}

// forkClearOnEsc clears a committed sidebar filter or thread search, so the
// first Esc undoes the search and only a second Esc does the stock thing
// (such as closing the thread).
func (a *App) forkClearOnEsc() bool {
	switch {
	case a.focusedPanel == PanelSidebar && a.sidebar.Filter() != "":
		a.sidebar.ClearFilter()
		return true
	case a.focusedPanel == PanelThread && a.threadPanel.SearchActive():
		a.threadPanel.ClearSearch()
		a.restoreSearchSegment()
		return true
	}
	return false
}

// forkSearchStep handles n (dir 1) and N (dir -1) on the thread pane and
// reports whether a thread search took the key.
func (a *App) forkSearchStep(dir int) bool {
	if a.focusedPanel != PanelThread || !a.threadPanel.SearchActive() {
		return false
	}
	a.threadPanel.StepSearch(dir)
	a.statusbar.SetSearch(a.threadSearchSegment())
	return true
}
