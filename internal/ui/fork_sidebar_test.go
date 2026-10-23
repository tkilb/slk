package ui

import "testing"

func TestSidebarCollapsesUnlessFocused(t *testing.T) {
	a := newTestApp(t, withSize(120, 40))
	a.SetSidebarCollapseWhenUnfocused(true)

	a.focusedPanel = PanelMessages
	a.View()
	if !a.sidebar.IsCollapsed("Direct Messages") {
		t.Fatalf("sidebar should be collapsed while the messages pane has focus")
	}

	a.focusedPanel = PanelSidebar
	a.View()
	if a.sidebar.IsCollapsed("Direct Messages") {
		t.Errorf("sidebar should expand while it has focus")
	}

	a.focusedPanel = PanelThread
	a.View()
	if !a.sidebar.IsCollapsed("Direct Messages") {
		t.Errorf("sidebar should collapse again when focus leaves it")
	}
}
