package ui

// SetSidebarCollapseWhenUnfocused turns on the fork's focus-driven sidebar
// collapse (fork feature sidebar-focus-collapse): every section collapses
// while the sidebar is unfocused and expands while the user is cursoring
// through it. The sidebar learns about focus on each render (renderSidebar).
func (a *App) SetSidebarCollapseWhenUnfocused(on bool) {
	a.sidebar.SetCollapseWhenUnfocused(on)
}
