package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gammons/slk/internal/ui/sidebar"
)

func sidebarFilterApp(t *testing.T) *App {
	t.Helper()
	a := newTestApp(t, withSize(120, 40), withChannels(
		sidebar.ChannelItem{ID: "C1", Name: "general", Type: "channel"},
		sidebar.ChannelItem{ID: "C2", Name: "random", Type: "channel"},
		sidebar.ChannelItem{ID: "D1", Name: "alice", Type: "dm"},
	))
	a.SetSidebarCollapseWhenUnfocused(true)
	a.focusedPanel = PanelSidebar
	a.View()
	return a
}

func typeSidebarFilter(a *App, s string) {
	for _, r := range s {
		a.Update(keyPress(r))
	}
}

func TestSlashOnSidebarOpensSidebarFilter(t *testing.T) {
	a := sidebarFilterApp(t)
	a.Update(keyPress('/'))
	if a.mode != ModeSidebarFilter {
		t.Fatalf("mode = %v, want ModeSidebarFilter", a.mode)
	}
}

func TestSlashOnMessagesStillStartsChannelSearch(t *testing.T) {
	a := sidebarFilterApp(t)
	a.focusedPanel = PanelMessages
	a.Update(keyPress('/'))
	if a.mode != ModeSearch {
		t.Fatalf("mode = %v, want ModeSearch", a.mode)
	}
}

func TestSidebarFilterNarrowsAsYouType(t *testing.T) {
	a := sidebarFilterApp(t)
	a.Update(keyPress('/'))
	typeSidebarFilter(a, "ran")
	if got := a.sidebar.VisibleItems(); len(got) != 1 || got[0].ID != "C2" {
		t.Fatalf("visible = %+v, want only C2", got)
	}
	if a.sidebar.SelectedID() != "C2" {
		t.Errorf("cursor should sit on the first match, got %q", a.sidebar.SelectedID())
	}
	a.Update(keyCode(tea.KeyBackspace))
	if got := a.sidebar.Filter(); got != "ra" {
		t.Errorf("filter after backspace = %q, want ra", got)
	}
}

func TestSidebarFilterEnterKeepsFilterInNormalMode(t *testing.T) {
	a := sidebarFilterApp(t)
	a.Update(keyPress('/'))
	typeSidebarFilter(a, "ran")
	a.Update(keyCode(tea.KeyEnter))
	if a.mode != ModeNormal {
		t.Fatalf("mode = %v, want ModeNormal", a.mode)
	}
	if a.sidebar.Filter() != "ran" || a.sidebar.SelectedID() != "C2" {
		t.Errorf("filter=%q selected=%q, want ran/C2", a.sidebar.Filter(), a.sidebar.SelectedID())
	}
}

func TestSidebarFilterEscInPromptClears(t *testing.T) {
	a := sidebarFilterApp(t)
	a.Update(keyPress('/'))
	typeSidebarFilter(a, "ran")
	a.Update(keyCode(tea.KeyEscape))
	if a.mode != ModeNormal || a.sidebar.Filter() != "" {
		t.Errorf("mode=%v filter=%q, want normal and empty", a.mode, a.sidebar.Filter())
	}
	if a.sidebar.SelectedID() != "C2" {
		t.Errorf("cursor should stay on the match, got %q", a.sidebar.SelectedID())
	}
}

func TestSidebarFilterEscInNormalModeClears(t *testing.T) {
	a := sidebarFilterApp(t)
	a.Update(keyPress('/'))
	typeSidebarFilter(a, "ran")
	a.Update(keyCode(tea.KeyEnter))
	a.Update(keyCode(tea.KeyEscape))
	if a.sidebar.Filter() != "" {
		t.Errorf("Esc in normal mode should clear the sidebar filter, got %q", a.sidebar.Filter())
	}
}

func TestSidebarFilterClearsWhenFocusLeaves(t *testing.T) {
	a := sidebarFilterApp(t)
	a.Update(keyPress('/'))
	typeSidebarFilter(a, "ran")
	a.Update(keyCode(tea.KeyEnter))
	a.focusedPanel = PanelMessages
	a.View()
	if a.sidebar.Filter() != "" {
		t.Errorf("filter should clear on blur, got %q", a.sidebar.Filter())
	}
}

func TestSidebarFilterPromptIsStaleAfterFocusMoves(t *testing.T) {
	a := sidebarFilterApp(t)
	a.Update(keyPress('/'))
	a.focusedPanel = PanelMessages // e.g. a mouse click
	a.View()
	a.Update(keyPress('x'))
	if a.mode != ModeNormal || a.sidebar.Filter() != "" {
		t.Errorf("mode=%v filter=%q, want normal and no filter", a.mode, a.sidebar.Filter())
	}
}
