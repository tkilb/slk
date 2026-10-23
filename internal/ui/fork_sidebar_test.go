package ui

import (
	"testing"

	"github.com/gammons/slk/internal/ui/sidebar"
)

func TestSidebarCollapsesUnlessFocused(t *testing.T) {
	a := newTestApp(t, withSize(120, 40), withChannels(
		sidebar.ChannelItem{ID: "C1", Name: "general", Type: "channel"},
		sidebar.ChannelItem{ID: "D1", Name: "alice", Type: "dm"},
	))
	a.SetSidebarCollapseWhenUnfocused(true)

	a.focusedPanel = PanelMessages
	a.View()
	if !a.sidebar.IsCollapsed("Direct Messages") {
		t.Fatalf("sidebar should be collapsed while the messages pane has focus")
	}

	a.focusedPanel = PanelSidebar
	a.View() // the sidebar learns about focus on render
	a.sidebar.SelectByID("D1")
	a.View()
	if a.sidebar.IsCollapsed("Direct Messages") {
		t.Errorf("the section under the cursor should expand while the sidebar has focus")
	}
	if !a.sidebar.IsCollapsed("Channels") {
		t.Errorf("sections without the cursor should stay collapsed")
	}

	a.focusedPanel = PanelThread
	a.View()
	if !a.sidebar.IsCollapsed("Direct Messages") {
		t.Errorf("sidebar should collapse again when focus leaves it")
	}
}
